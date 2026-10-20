// The `app-listener update` command self-updates from the latest GitHub release of --channel
// (default stable): skip if not newer; download binary + sha256 + Ed25519 signature and verify
// (embedded public key, checksum, GitHub asset digest); show the changelog and confirm (--yes
// skips, non-terminal stdin aborts); atomically replace the binary and restart the daemon per the
// installer's stop/start contract. The signing key lives only in a GitHub Actions secret.
package update

import (
	"cmp"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
	"debug/elf"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/certificates"
	"github.com/Virgula0/app-listener/internal/headless"
	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/systemd"
	"github.com/Virgula0/app-listener/internal/wizard"
)

const (
	// githubAPIBase is the GitHub REST API endpoint used to list releases.
	githubAPIBase = "https://api.github.com"
	// Only repo the updater trusts: its workflow signs assets with the key
	// embedded in this binary, so others cannot be trusted.
	defaultRepo = "Virgula0/app-listener"
	// Matches the release workflow's pre-YYYYMMDD-HHMMSS-<sha7> tags; HHMMSS is optional, since
	// installs from releases cut before it was added run tags without it.
	preVersionPattern = `^pre-(\d{8})(?:-\d{6})?-([0-9a-f]{7})$`
	// Matches vX.Y.Z semver stable tags.
	stableVersionPattern = `^v?(\d+)\.(\d+)\.(\d+)$`
	// updateHTTPTimeout bounds a single GitHub API/download request.
	updateHTTPTimeout = 60 * time.Second
)

const (
	// The accepted --channel values (stable tracks non-prereleases).
	channelStable     = "stable"
	channelPreRelease = "pre-release"
)

var preVersionRe = regexp.MustCompile(preVersionPattern)

var stableVersionRe = regexp.MustCompile(stableVersionPattern)

// yesFlag skips the changelog viewer and the confirmation prompts (a changed unit is accepted).
var yesFlag bool

// channelFlag selects the tracked release channel (default channelStable).
var channelFlag string

func init() {
	UpdateCmd.Flags().BoolVarP(&yesFlag, "yes", "y", false,
		"Don't show the changelog or ask (scripts; a changed daemon unit is accepted, its diff logged)")
	UpdateCmd.Flags().StringVar(&channelFlag, "channel", channelStable, "Release channel: stable or pre-release")
}

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the installed binaries to the latest signed release",
	Long: `Update /usr/local/sbin/app-listener (and app-listener-gui, when installed) to the
latest release of --channel. The download is checked against the release's
Ed25519 signature first; the changelog and a changed daemon unit are shown and
asked before anything changes, then the daemon is restarted. Without a
terminal it refuses unless --yes.`,
	Example: `  sudo app-listener update
  sudo app-listener update --channel pre-release`,
	Args: cobra.NoArgs,
	RunE: runUpdate,
}

// releaseAssets names one architecture's release assets (.github/workflows/_build-release.yml);
// check tells the build apart (headless.Check, or headless.CheckGUI for the desktop asset).
type releaseAssets struct {
	binary, checksum, signature string
	machine                     elf.Machine
	check                       func(path string) error
}

// assetsFor maps a Go architecture to its headless release assets. amd64 keeps the unsuffixed
// names: older binaries' `update` looks them up.
func assetsFor(goarch string) (releaseAssets, error) {
	return archAssets(goarch, "app-listener", headless.Check)
}

// guiAssetsFor is assetsFor for the desktop build (app-listener-gui[-arm64]).
func guiAssetsFor(goarch string) (releaseAssets, error) {
	return archAssets(goarch, "app-listener-gui", headless.CheckGUI)
}

func archAssets(goarch, base string, check func(string) error) (releaseAssets, error) {
	var bin string
	var machine elf.Machine
	switch goarch {
	case "amd64":
		bin, machine = base, elf.EM_X86_64
	case "arm64":
		bin, machine = base+"-arm64", elf.EM_AARCH64
	default:
		return releaseAssets{}, fmt.Errorf("no release is published for linux/%s: build from source (make build)", goarch)
	}
	return releaseAssets{binary: bin, checksum: bin + ".sha256", signature: bin + ".sha256.sig", machine: machine,
		check: check}, nil
}

// githubRelease mirrors the updater-relevant fields of the releases API.
type githubRelease struct {
	TagName     string `json:"tag_name"`
	Body        string `json:"body"`
	Prerelease  bool   `json:"prerelease"`
	PublishedAt string `json:"published_at"`
	Assets      []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
		Digest             string `json:"digest"`
	} `json:"assets"`
}

// validateChannel rejects unsupported --channel values.
func validateChannel(ch string) error {
	if ch != channelStable && ch != channelPreRelease {
		return fmt.Errorf("invalid --channel %q: must be %q or %q", ch, channelStable, channelPreRelease)
	}
	return nil
}

// runUpdate drives the update flow: every check runs first; changelog and
// confirmation come last unless --yes skips them.
func runUpdate(cmd *cobra.Command, args []string) error {
	return run(channelFlag, yesFlag, systemd.GUIInstalled())
}

// RunWith updates to the latest release of the stable (or pre-release) channel with every prompt on
// ui: the desktop GUI's updater. Like the CLI it deploys the release's headless asset, never the
// running executable. withGUI also installs or updates the release's desktop build at
// systemd.GUIBinaryPath; without it that file is left as is. One run at a time: wizard.Use is
// process-wide.
//
// TODO(self-binary-guard): runs in the GUI process (not the daemon's exe inode), safe only because
// nothing it touches is daemon-guarded: with the daemon up it just execs and hashes
// systemd.InstallBinaryPath (so do install.RunWith and gui readInstalled); the swap follows
// systemctl stop. When the daemon guards its own binary, keep those reads open to root (ModeReadOnly)
// or move them after the stop; a guard outliving the stop leaves the daemon stopped on a failed
// swap, so then hand the swap to the installed binary (fixed argv, no shell) or to the daemon.
func RunWith(ui wizard.UI, preRelease, withGUI bool) error {
	if !runMu.TryLock() {
		return errors.New("an update is already in progress")
	}
	defer runMu.Unlock()
	defer wizard.Use(ui)()
	return run(channelOf(preRelease), false, withGUI)
}

var runMu sync.Mutex

func channelOf(preRelease bool) string {
	if preRelease {
		return channelPreRelease
	}
	return channelStable
}

// run is the update of channel; auto is --yes (no changelog viewer or prompts). withGUI puts the
// desktop build on the daemon's release (installing it when missing).
func run(channel string, auto, withGUI bool) error {
	assets, err := checkUpdatePreconditions(channel)
	if err != nil {
		return err
	}

	wizard.Stage("Checking the installed version")
	installed := readInstalledVersion(systemd.InstallBinaryPath)
	log.Infof("installed version: %s", installed)

	wizard.Stage("Looking up the latest " + channel + " release")
	c, err := findUpdate(channel, installed)
	if err != nil {
		return err
	}
	if c == nil {
		log.Infof("no %s releases found: nothing to update", channel)
		return nil
	}
	log.Infof("latest %s release: %s (published %s)", channel, c.latest.TagName, c.latest.PublishedAt)
	if !c.newer {
		log.Infof("installed version %s is up to date", installed)
		if withGUI {
			return syncGUI(installed)
		}
		return nil
	}
	return deploy(c, assets, installed, channel, auto, withGUI)
}

// deploy downloads and verifies c's release (and its desktop build when withGUI), asks (changelog,
// unit diff) and swaps it in.
func deploy(c *candidate, assets releaseAssets, installed, channel string, auto, withGUI bool) error {
	latest, releases := c.latest, c.all
	binURL, binDigest, checksumURL, sigURL, err := resolveAssets(&latest, assets)
	if err != nil {
		return err
	}

	// Staging dir is owned here: the verified download must survive until applyUpdate swaps the
	// binary (downloadAndVerify's return would clean it up early).
	tmp, err := os.MkdirTemp("", "app-listener-update-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	wizard.Stage("Downloading and verifying " + latest.TagName)
	binPath, err := downloadAndVerify(tmp, &latest, assets, binURL, binDigest, checksumURL, sigURL)
	if err != nil {
		return err
	}
	guiPath := ""
	if withGUI {
		wizard.Stage("Downloading and verifying the desktop build of " + latest.TagName)
		if guiPath, err = fetchGUI(tmp, &latest); err != nil {
			return err
		}
	}

	wizard.Stage("Reviewing the changelog")
	proceed, err := confirmUpdate(&latest, cumulativeChangelog(installed, &latest, releases, channel), auto)
	if err != nil {
		return err
	}
	if !proceed {
		log.Info("update canceled by the user: nothing was changed")
		return nil
	}

	// Asked before anything changes: the new release's unit is read from the verified download.
	wizard.Stage("Checking the daemon unit")
	unit, err := systemd.PlanDaemonUnitUpdate(binPath, confirmUnit(auto))
	if err != nil {
		return err
	}

	wizard.Stage("Installing " + latest.TagName + " and restarting the daemon")
	if err := applyUpdate(binPath, unit); err != nil {
		return err
	}
	if guiPath == "" {
		return nil
	}
	wizard.Stage("Installing the desktop build of " + latest.TagName)
	if err := deployGUI(guiPath); err != nil {
		return fmt.Errorf("the daemon runs %s, but its desktop build was not installed: %w", latest.TagName, err)
	}
	return nil
}

// deployGUI installs the verified desktop build at systemd.GUIBinaryPath.
func deployGUI(path string) error {
	if err := systemd.DeployGUIBinary(path); err != nil {
		return err
	}
	log.Infof("desktop build installed at %s: an open window keeps running the old one until it is restarted",
		systemd.GUIBinaryPath)
	return nil
}

// fetchGUI downloads and verifies r's desktop build into dir. "" with a warning when r ships none
// for this architecture (the release builds it best-effort); a present asset that fails any check
// is an error.
func fetchGUI(dir string, r *githubRelease) (string, error) {
	assets, err := guiAssetsFor(runtime.GOARCH)
	if err != nil {
		return "", err
	}
	if !hasAsset(r, assets.binary) {
		log.Warnf("release %s has no %s asset: its desktop build is not installed", r.TagName, assets.binary)
		return "", nil
	}
	binURL, binDigest, checksumURL, sigURL, err := resolveAssets(r, assets)
	if err != nil {
		return "", err
	}
	return downloadAndVerify(dir, r, assets, binURL, binDigest, checksumURL, sigURL)
}

func hasAsset(r *githubRelease, name string) bool {
	for i := range r.Assets {
		if r.Assets[i].Name == name {
			return true
		}
	}
	return false
}

// syncGUI puts the desktop build on the installed daemon's release when the daemon itself is up to
// date: the window drives that daemon (edit-protected --backend), so the two ship together.
func syncGUI(installed string) error {
	if systemd.GUIInstalled() && readInstalledVersion(systemd.GUIBinaryPath) == installed {
		log.Infof("the desktop build at %s is on %s too", systemd.GUIBinaryPath, installed)
		return nil
	}
	if !preVersionRe.MatchString(installed) && !stableVersionRe.MatchString(installed) {
		log.Warnf("the installed daemon (%s) is not a release build: no desktop build to match it", installed)
		return nil
	}
	wizard.Stage("Downloading and verifying the desktop build of " + installed)
	r, err := fetchRelease(githubAPIBase, defaultRepo, installed, updateHTTPClient())
	if err != nil {
		return err
	}
	tmp, err := os.MkdirTemp("", "app-listener-update-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	path, err := fetchGUI(tmp, &r)
	if err != nil || path == "" {
		return err
	}
	wizard.Stage("Installing the desktop build of " + installed)
	return deployGUI(path)
}

// confirmUnit shows the daemon unit diff and asks; autoConfirm logs the diff and accepts it.
func confirmUnit(autoConfirm bool) systemd.ConfirmUnitFunc {
	if !autoConfirm {
		return inst.ConfirmOverwrite
	}
	return func(path string, installed, desired []byte) (bool, error) {
		log.Infof("--yes given: updating %s to the new release's unit:\n%s", path, inst.UnifiedDiff(string(installed), string(desired)))
		return true, nil
	}
}

// checkUpdatePreconditions refuses non-root, a bad --channel, an arch with no release, or a
// missing install; it returns this arch's release assets.
func checkUpdatePreconditions(channel string) (releaseAssets, error) {
	if os.Geteuid() != 0 {
		return releaseAssets{}, errors.New("update must be run as root: sudo app-listener update")
	}
	if err := validateChannel(channel); err != nil {
		return releaseAssets{}, err
	}
	assets, err := assetsFor(runtime.GOARCH)
	if err != nil {
		return releaseAssets{}, err
	}
	if _, err := os.Lstat(systemd.InstallBinaryPath); err != nil {
		return releaseAssets{}, fmt.Errorf("the daemon is not installed (no %s): run `app-listener install` first", systemd.InstallBinaryPath)
	}
	return assets, nil
}

// downloadAndVerify downloads the three assets into the caller-owned staging dir and runs every
// verification (signature, checksum, digest, binary sanity) before anything reaches the user;
// progress on the TUI bar.
func downloadAndVerify(tmp string, r *githubRelease, assets releaseAssets, binURL, binDigest, checksumURL, sigURL string) (string, error) {
	var binPath, checksumPath, sigPath string
	err := wizard.WithBottomBar(func(bar *wizard.BottomBar) error {
		var dlErr error
		binPath, checksumPath, sigPath, dlErr = downloadReleaseFiles(updateHTTPClient(), tmp, assets, binURL, checksumURL, sigURL, bar)
		return dlErr
	})
	if err != nil {
		return "", err
	}

	err = verifyRelease(binPath, checksumPath, sigPath, binDigest, assets.binary)
	if err != nil {
		return "", fmt.Errorf("release verification failed: %w — the download is rejected; fix the release before updating", err)
	}
	log.Infof("release %s verified: Ed25519 signature, sha256 checksum and GitHub digest match", r.TagName)

	err = sanityCheckBinary(binPath, r.TagName, assets)
	if err != nil {
		return "", err
	}
	log.Infof("release %s verified and downloaded to %s", r.TagName, binPath)
	return binPath, nil
}

// confirmUpdate: auto (--yes) logs the changelog and proceeds; otherwise the notes open in a viewer
// (a terminal one needs a tty) and the user confirms.
func confirmUpdate(r *githubRelease, notes string, auto bool) (bool, error) {
	if auto {
		log.Infof("--yes given: skipping the changelog viewer and the confirmation prompt")
		log.Infof("changelog up to %s:\n%s", r.TagName, notes)
		return true, nil
	}
	if wizard.Active() == nil && !isTerminal(os.Stdin) {
		return false, errors.New("stdin is not a terminal: interactive confirmation is impossible; re-run with --yes to update without prompts")
	}
	if err := showChangelog(r.TagName, notes); err != nil {
		return false, fmt.Errorf("showing the changelog: %w", err)
	}
	ok, err := wizard.ConfirmOnce(fmt.Sprintf("Update to %s now?", r.TagName), "Update")
	if err != nil {
		return false, err
	}
	return ok, nil
}

// isTerminal reports whether f is a character device (a real terminal).
func isTerminal(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

// resolveAssets returns this arch's binary/checksum/signature URLs and the
// GitHub-provided binary digest when present.
func resolveAssets(r *githubRelease, assets releaseAssets) (binURL, binDigest, checksumURL, sigURL string, err error) {
	asset := func(name string) (url, digest string, ok bool) {
		for i := range r.Assets {
			if r.Assets[i].Name == name {
				return r.Assets[i].BrowserDownloadURL, r.Assets[i].Digest, true
			}
		}
		return "", "", false
	}
	var ok bool
	if binURL, binDigest, ok = asset(assets.binary); !ok {
		return "", "", "", "", fmt.Errorf("release %s has no %s asset", r.TagName, assets.binary)
	}
	if checksumURL, _, ok = asset(assets.checksum); !ok {
		return "", "", "", "", fmt.Errorf("release %s has no %s asset", r.TagName, assets.checksum)
	}
	if sigURL, _, ok = asset(assets.signature); !ok {
		return "", "", "", "", fmt.Errorf("release %s has no %s asset", r.TagName, assets.signature)
	}
	return binURL, binDigest, checksumURL, sigURL, nil
}

// downloadReleaseFiles downloads the three assets into dir, reporting
// progress on bar when non-nil.
func downloadReleaseFiles(client *http.Client, dir string, assets releaseAssets, binURL, checksumURL, sigURL string, bar *wizard.BottomBar) (binPath, checksumPath, sigPath string, err error) {
	binPath = filepath.Join(dir, assets.binary)
	if err := downloadFile(client, binURL, binPath, bar, "Downloading "+assets.binary); err != nil {
		return "", "", "", err
	}
	checksumPath = filepath.Join(dir, assets.checksum)
	if err := downloadFile(client, checksumURL, checksumPath, bar, "Downloading checksum"); err != nil {
		return "", "", "", err
	}
	sigPath = filepath.Join(dir, assets.signature)
	if err := downloadFile(client, sigURL, sigPath, bar, "Downloading signature"); err != nil {
		return "", "", "", err
	}
	return binPath, checksumPath, sigPath, nil
}

// InstalledVersion is the version the installed binary reports ("unknown" when it can't tell).
func InstalledVersion() string { return readInstalledVersion(systemd.InstallBinaryPath) }

// InstalledGUIVersion is InstalledVersion for the desktop build; "" when none is installed.
func InstalledGUIVersion() string {
	if !systemd.GUIInstalled() {
		return ""
	}
	return readInstalledVersion(systemd.GUIBinaryPath)
}

// readInstalledVersion runs the installed binary with --version; failures or
// unknown output map to "unknown", treated as older than any release.
func readInstalledVersion(installPath string) string {
	cmd := exec.CommandContext(context.Background(), installPath, "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Warnf("installed binary does not report a version (%v): treating it as unknown/older", err)
		return "unknown"
	}
	line := strings.TrimSpace(string(out))
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return "unknown"
	}
	return fields[len(fields)-1]
}

// fetchReleases returns the repository releases, newest first, unfiltered.
func fetchReleases(apiBase, repo string, client *http.Client) ([]githubRelease, error) {
	url := fmt.Sprintf("%s/repos/%s/releases?per_page=100", apiBase, repo)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "app-listener-updater")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("listing releases from %s: %w", repo, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("listing releases from %s: HTTP %d: %s", repo, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var releases []githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("decoding releases of %s: %w", repo, err)
	}
	return releases, nil
}

// filterChannel keeps the selected channel's releases in server order:
// stable = non-prereleases, pre-release = prereleases.
func filterChannel(releases []githubRelease, channel string) []githubRelease {
	keep := func(r githubRelease) bool {
		if channel == channelPreRelease {
			return r.Prerelease
		}
		return !r.Prerelease
	}
	var out []githubRelease
	for i := range releases {
		if keep(releases[i]) {
			out = append(out, releases[i])
		}
	}
	return out
}

// pickLatestRelease returns the release with the most recent published_at
// timestamp.
func pickLatestRelease(releases []githubRelease) (githubRelease, bool) {
	var latest githubRelease
	found := false
	for i := range releases {
		ts, err := time.Parse(time.RFC3339, releases[i].PublishedAt)
		if err != nil {
			continue
		}
		if !found || ts.After(latestTime(&latest)) {
			latest = releases[i]
			found = true
		}
	}
	return latest, found
}

// latestTime parses published_at (zero on failure, so bad timestamps lose).
func latestTime(r *githubRelease) time.Time {
	ts, err := time.Parse(time.RFC3339, r.PublishedAt)
	if err != nil {
		return time.Time{}
	}
	return ts
}

// newerThanInstalled reports whether the latest pre-release beats the
// installed tag; non-pre-release installed versions always update.
func newerThanInstalled(installed string, latest *githubRelease, releases []githubRelease) bool {
	if !preVersionRe.MatchString(installed) {
		log.Infof("installed version %q is not a pre-release build: the latest pre-release is newer", installed)
		return true
	}
	if installed == latest.TagName {
		return false
	}
	for i := range releases {
		if releases[i].TagName == installed {
			return latestTime(latest).After(latestTime(&releases[i]))
		}
	}
	// Installed tag fell off the fetched page: treat as older, update.
	log.Infof("installed tag %s not found among the fetched pre-releases: updating", installed)
	return true
}

// newerThanStable reports whether the latest stable beats the installed
// version; non-semver installed versions always update, equal/older do not.
func newerThanStable(installed string, latest *githubRelease) bool {
	if installed == latest.TagName {
		return false
	}
	installMSIP, installOK := parseStableVersion(installed)
	if !installOK {
		log.Infof("installed version %q is not a stable release: the latest stable %s is newer", installed, latest.TagName)
		return true
	}
	latestVer, latestOK := parseStableVersion(latest.TagName)
	if !latestOK {
		log.Warnf("the latest stable release %q does not follow the vX.Y.Z tag convention: updating", latest.TagName)
		return true
	}
	return compareStableVersions(latestVer, installMSIP) > 0
}

// stableVersion holds the numeric components of a vX.Y.Z stable version.
type stableVersion struct {
	major, minor, patch uint64
}

// parseStableVersion parses a vX.Y.Z (or X.Y.Z) tag; ok is false otherwise.
func parseStableVersion(s string) (stableVersion, bool) {
	m := stableVersionRe.FindStringSubmatch(s)
	if m == nil {
		return stableVersion{}, false
	}
	v := stableVersion{}
	var err error
	if v.major, err = parseVerPart(m[1]); err != nil {
		return stableVersion{}, false
	}
	if v.minor, err = parseVerPart(m[2]); err != nil {
		return stableVersion{}, false
	}
	if v.patch, err = parseVerPart(m[3]); err != nil {
		return stableVersion{}, false
	}
	return v, true
}

// parseVerPart converts a numeric version component.
func parseVerPart(s string) (uint64, error) {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// compareStableVersions returns <0/0/>0 for a<b/a==b/a>b.
func compareStableVersions(a, b stableVersion) int {
	if a.major != b.major {
		return cmp.Compare(a.major, b.major)
	}
	if a.minor != b.minor {
		return cmp.Compare(a.minor, b.minor)
	}
	return cmp.Compare(a.patch, b.patch)
}

// updateHTTPClient returns a client bounded by updateHTTPTimeout.
func updateHTTPClient() *http.Client {
	return &http.Client{Timeout: updateHTTPTimeout}
}

// downloadFile downloads url into path atomically (temp + rename), mode 0700: hardened kernels
// refuse to exec without x bits and os.CreateTemp makes 0600. Progress goes to bar when non-nil.
func downloadFile(client *http.Client, url, path string, bar *wizard.BottomBar, label string) error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "app-listener-updater")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("downloading %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("downloading %s: HTTP %d", url, resp.StatusCode)
	}
	body := io.Reader(resp.Body)
	if bar != nil && resp.ContentLength > 0 {
		body = &progressReader{reader: resp.Body, total: resp.ContentLength, bar: bar, label: label}
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".download-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err := io.Copy(f, body); err != nil {
		f.Close()
		return fmt.Errorf("downloading %s: %w", url, err)
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Rename(f.Name(), path); err != nil {
		return err
	}
	// Executed during the sanity check; CreateTemp leaves 0600, and some
	// hardened kernels refuse to exec that even as root.
	return os.Chmod(path, 0o700)
}

// progressReader wraps a download body and reports the fraction received
// to the TUI bottom bar.
type progressReader struct {
	reader io.Reader
	total  int64
	read   int64
	bar    *wizard.BottomBar
	label  string
}

func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.reader.Read(b)
	p.read += int64(n)
	if p.bar != nil && p.total > 0 {
		p.bar.Set(p.label, float64(p.read)/float64(p.total))
	}
	return n, err
}

// Release is a downloaded, verified release: Binary is its headless asset, GUI its desktop build
// ("" when the release ships none for this architecture).
type Release struct {
	Tag, Binary, GUI string
}

// FetchVerified downloads this architecture's headless and desktop builds of tag into the
// caller-owned dir and runs every check update does on each. The GUI installer deploys the headless
// one as the daemon instead of itself (the daemon must not be the fyne build) and the desktop one
// at systemd.GUIBinaryPath.
func FetchVerified(tag, dir string) (Release, error) {
	r, err := fetchRelease(githubAPIBase, defaultRepo, tag, updateHTTPClient())
	if err != nil {
		return Release{}, err
	}
	return fetchPair(dir, &r)
}

func fetchPair(dir string, r *githubRelease) (Release, error) {
	assets, err := assetsFor(runtime.GOARCH)
	if err != nil {
		return Release{}, err
	}
	binURL, binDigest, checksumURL, sigURL, err := resolveAssets(r, assets)
	if err != nil {
		return Release{}, err
	}
	bin, err := downloadAndVerify(dir, r, assets, binURL, binDigest, checksumURL, sigURL)
	if err != nil {
		return Release{}, err
	}
	gui, err := fetchGUI(dir, r)
	if err != nil {
		return Release{}, err
	}
	return Release{Tag: r.TagName, Binary: bin, GUI: gui}, nil
}

// candidate is a channel's newest release, every release (the changelog spans them) and whether
// it beats the installed version.
type candidate struct {
	latest githubRelease
	all    []githubRelease
	newer  bool
}

// findUpdate picks channel's newest release for installed; nil when the channel has none.
func findUpdate(channel, installed string) (*candidate, error) {
	releases, err := fetchReleases(githubAPIBase, defaultRepo, updateHTTPClient())
	if err != nil {
		return nil, err
	}
	inChannel := filterChannel(releases, channel)
	latest, ok := pickLatestRelease(inChannel)
	if !ok {
		return nil, nil
	}
	c := &candidate{latest: latest, all: releases}
	if channel == channelStable {
		c.newer = newerThanStable(installed, &latest)
	} else {
		c.newer = newerThanInstalled(installed, &latest, inChannel)
	}
	return c, nil
}

// Latest returns the newest release of the stable (or pre-release) channel without downloading it,
// and whether update would install it over the installed version.
func Latest(preRelease bool, installed string) (tag string, newer bool, err error) {
	channel := channelOf(preRelease)
	c, err := findUpdate(channel, installed)
	if err != nil {
		return "", false, err
	}
	if c == nil {
		return "", false, fmt.Errorf("no %s release is published", channel)
	}
	return c.latest.TagName, c.newer, nil
}

// FetchLatest is FetchVerified for the newest release of a channel (preRelease: the pre-YYYYMMDD
// builds, else stable vX.Y.Z), picked as `update --channel` picks it.
func FetchLatest(preRelease bool, dir string) (Release, error) {
	channel := channelOf(preRelease)
	releases, err := fetchReleases(githubAPIBase, defaultRepo, updateHTTPClient())
	if err != nil {
		return Release{}, err
	}
	r, ok := pickLatestRelease(filterChannel(releases, channel))
	if !ok {
		return Release{}, fmt.Errorf("no %s release is published", channel)
	}
	log.Infof("latest %s release: %s (published %s)", channel, r.TagName, r.PublishedAt)
	return fetchPair(dir, &r)
}

// fetchRelease returns the release tagged tag.
func fetchRelease(apiBase, repo, tag string, client *http.Client) (githubRelease, error) {
	url := fmt.Sprintf("%s/repos/%s/releases/tags/%s", apiBase, repo, tag)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return githubRelease{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "app-listener-updater")
	resp, err := client.Do(req)
	if err != nil {
		return githubRelease{}, fmt.Errorf("fetching release %s of %s: %w", tag, repo, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return githubRelease{}, fmt.Errorf("fetching release %s of %s: HTTP %d: %s", tag, repo, resp.StatusCode,
			strings.TrimSpace(string(body)))
	}
	var r githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return githubRelease{}, fmt.Errorf("decoding release %s: %w", tag, err)
	}
	if r.TagName != tag {
		return githubRelease{}, fmt.Errorf("asked for release %s, GitHub returned %s", tag, r.TagName)
	}
	return r, nil
}

// verifyRelease checks the checksum's Ed25519 signature against the embedded key, that it was
// signed for assetName, then the checksum against the binary, then the GitHub asset digest.
func verifyRelease(binPath, checksumPath, sigPath, assetDigest, assetName string) error {
	pub, err := parsePublicKey(certificates.ReleasePublicKeyPEM)
	if err != nil {
		return fmt.Errorf("parsing the embedded release public key: %w", err)
	}
	return verifyReleaseWithKey(binPath, checksumPath, sigPath, assetDigest, assetName, pub)
}

// verifyReleaseWithKey is verifyRelease with an explicit public key.
func verifyReleaseWithKey(binPath, checksumPath, sigPath, assetDigest, assetName string, pub ed25519.PublicKey) error {
	bin, err := os.ReadFile(binPath)
	if err != nil {
		return err
	}
	checksum, err := os.ReadFile(checksumPath)
	if err != nil {
		return err
	}
	sig, err := os.ReadFile(sigPath)
	if err != nil {
		return err
	}
	if !ed25519.Verify(pub, checksum, sig) {
		return errors.New("the Ed25519 signature of the checksum does not match the release public key")
	}

	expected, signedName, err := parseChecksum(checksum)
	if err != nil {
		return fmt.Errorf("parsing the signed checksum file: %w", err)
	}
	// Every arch's checksum is signed by the same key: the signed name stops another arch's (or
	// the GUI's) validly signed pair from being served under this asset's name.
	if signedName != assetName {
		return fmt.Errorf("the signed checksum is for %q, not %q", signedName, assetName)
	}
	actual := sha256.Sum256(bin)
	if !strings.EqualFold(hex.EncodeToString(actual[:]), expected) {
		return fmt.Errorf("sha256 mismatch: the checksum file says %s, the downloaded binary is %s", expected, hex.EncodeToString(actual[:]))
	}

	if assetDigest != "" {
		want := strings.TrimPrefix(assetDigest, "sha256:")
		if !strings.EqualFold(want, expected) {
			return fmt.Errorf("the GitHub asset digest (%s) does not match the signed checksum (%s)", assetDigest, expected)
		}
	}
	return nil
}

// parsePublicKey parses the embedded PEM SPKI Ed25519 public key.
func parsePublicKey(pemBytes []byte) (ed25519.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("no PEM block found")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	ed, ok := pub.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("unexpected public key type %T", pub)
	}
	return ed, nil
}

// parseChecksum extracts the validated sha256 hex token and the file name of a one-line
// `sha256sum` checksum file.
func parseChecksum(data []byte) (sum, name string, err error) {
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return "", "", errors.New("the checksum file is empty")
	}
	if len(fields) != 2 {
		return "", "", fmt.Errorf("expected one \"<sha256>  <name>\" line, got %d fields", len(fields))
	}
	if len(fields[0]) != 64 {
		return "", "", fmt.Errorf("expected a 64-hex-char sha256, got %q", fields[0])
	}
	if _, err := hex.DecodeString(fields[0]); err != nil {
		return "", "", fmt.Errorf("invalid sha256 hex %q: %w", fields[0], err)
	}
	return fields[0], strings.TrimPrefix(fields[1], "*"), nil
}

// sanityCheckBinary refuses a file that is not an ELF for assets' machine, whose --version misses
// the wanted tag, or that assets.check refuses (the headless build where the GUI is expected, or
// the reverse).
func sanityCheckBinary(path, wantVersion string, assets releaseAssets) error {
	machine := assets.machine
	f, err := elf.Open(path)
	if err != nil {
		return fmt.Errorf("%s is not an ELF executable: %w", path, err)
	}
	got := f.Machine
	f.Close()
	if got != machine {
		return fmt.Errorf("%s is built for %s, this host runs %s — refusing to deploy", path, got, machine)
	}
	cmd := exec.CommandContext(context.Background(), path, "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("the downloaded binary does not run (--version failed: %v): %s", err, strings.TrimSpace(string(out)))
	}
	if !strings.Contains(string(out), wantVersion) {
		return fmt.Errorf("the downloaded binary reports %q, expected tag %q — refusing to deploy", strings.TrimSpace(string(out)), wantVersion)
	}
	return assets.check(path)
}

// applyUpdate stops the daemon (verifying lock state), atomically replaces the binary, writes unit
// when non-nil and brings the daemon back to enabled-and-running.
func applyUpdate(binPath string, unit []byte) error {
	if err := systemd.DeployInstalledBinary(binPath, unit); err != nil {
		return err
	}
	log.Info("update complete: the daemon is running the new binary")
	return nil
}
