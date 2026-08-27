// The `app-listener update` command self-updates the daemon from the latest
// GitHub release of --channel (default stable): read the embedded version,
// pick the newest release of the channel (skip when not older); download
// binary + sha256 + Ed25519 signature and verify them (embedded public key,
// checksum, GitHub asset digest); show the changelog and confirm (--yes
// skips, non-terminal stdin aborts); atomically replace the installed binary
// and restart the daemon mirroring the installer's stop/start contract. The
// signing private key lives only in a GitHub Actions secret.
package update

import (
	"cmp"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
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
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/certificates"
	"github.com/Virgula0/app-listener/internal/systemd"
	"github.com/Virgula0/app-listener/internal/wizard"
)

const (
	// githubAPIBase is the GitHub REST API endpoint used to list releases.
	githubAPIBase = "https://api.github.com"
	// Only repo the updater trusts: its workflow signs assets with the key
	// embedded in this binary, so others cannot be trusted.
	defaultRepo = "Virgula0/app-listener"
	// Release asset names produced by .github/workflows/release.yml.
	releaseAssetBinary    = "app-listener"
	releaseAssetChecksum  = "app-listener.sha256"
	releaseAssetSignature = "app-listener.sha256.sig"
	// Matches pre-YYYYMMDD-<sha7> tags from the release workflow.
	preVersionPattern = `^pre-(\d{8})-([0-9a-f]{7})$`
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

// yesFlag skips the changelog viewer and the confirmation prompt.
var yesFlag bool

// channelFlag selects the tracked release channel (default channelStable).
var channelFlag string

func init() {
	UpdateCmd.Flags().BoolVar(&yesFlag, "yes", false,
		"update without the changelog viewer and confirmation prompt")
	UpdateCmd.Flags().StringVar(&channelFlag, "channel", channelStable,
		"release channel to track: stable or pre-release")
}

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Self-update from the latest GitHub release of the selected channel",
	Long: `Root-only self-updater for the daemon mode.

The command checks the latest release of the selected --channel (stable
or pre-release, default stable), compares it with the version embedded
in the installed binary at /usr/local/sbin/app-listener, and when a
newer one exists:

  1. downloads the release binary and its sha256 checksum
  2. verifies the Ed25519 signature of the checksum against the public key
     compiled into this binary (certificates/app-listener-release.pub),
     then the checksum against the downloaded binary, and the
     GitHub-provided asset digest when present — a failed verification
     aborts the update
  3. shows the release changelog in a TUI viewer and asks for confirmation
     before anything is changed
  4. replaces the installed binary atomically
  5. restarts the systemd daemon (it is stopped first so the watch
     directories are verified locked again before the new binary starts)

The stable channel tracks non-pre-release releases (vX.Y.Z tags); the
pre-release channel tracks the pre-YYYYMMDD-<sha> builds produced by the
release workflow.

When stdin is not a terminal the command aborts without updating; use
--yes to update without the viewer and the prompt (for scripts and
cron jobs). The repository is fixed: only releases of the signing
repository are accepted.`,
	Args: cobra.NoArgs,
	RunE: runUpdate,
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
	if os.Geteuid() != 0 {
		return errors.New("update must be run as root: sudo app-listener update")
	}
	if err := validateChannel(channelFlag); err != nil {
		return err
	}
	if _, err := os.Lstat(systemd.InstallBinaryPath); err != nil {
		return fmt.Errorf("the daemon is not installed (no %s): run `app-listener install` first", systemd.InstallBinaryPath)
	}

	installed := readInstalledVersion(systemd.InstallBinaryPath)
	log.Infof("installed version: %s", installed)

	releases, err := fetchReleases(githubAPIBase, defaultRepo, updateHTTPClient())
	if err != nil {
		return err
	}
	channel := filterChannel(releases, channelFlag)
	latest, ok := pickLatestRelease(channel)
	if !ok {
		log.Infof("no %s releases found: nothing to update", channelFlag)
		return nil
	}
	log.Infof("latest %s release: %s (published %s)", channelFlag, latest.TagName, latest.PublishedAt)

	var newer bool
	if channelFlag == channelStable {
		newer = newerThanStable(installed, &latest)
	} else {
		newer = newerThanInstalled(installed, &latest, channel)
	}
	if !newer {
		log.Infof("installed version %s is up to date", installed)
		return nil
	}

	binURL, binDigest, checksumURL, sigURL, err := resolveAssets(&latest)
	if err != nil {
		return err
	}

	// Staging dir owned here: the verified download must survive until
	// applyUpdate swaps the binary (downloadAndVerify's return would clean
	// it up prematurely).
	tmp, err := os.MkdirTemp("", "app-listener-update-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	binPath, err := downloadAndVerify(tmp, &latest, binURL, binDigest, checksumURL, sigURL)
	if err != nil {
		return err
	}

	proceed, err := confirmUpdate(&latest)
	if err != nil {
		return err
	}
	if !proceed {
		log.Info("update canceled by the user: nothing was changed")
		return nil
	}

	return applyUpdate(binPath)
}

// downloadAndVerify downloads the three assets into the caller-owned staging
// dir and runs every verification (signature, checksum, digest, binary
// sanity) before anything reaches the user; progress shows on the TUI bar.
func downloadAndVerify(tmp string, r *githubRelease, binURL, binDigest, checksumURL, sigURL string) (string, error) {
	var binPath, checksumPath, sigPath string
	err := wizard.WithBottomBar(func(bar *wizard.BottomBar) error {
		var dlErr error
		binPath, checksumPath, sigPath, dlErr = downloadReleaseFiles(updateHTTPClient(), tmp, binURL, checksumURL, sigURL, bar)
		return dlErr
	})
	if err != nil {
		return "", err
	}

	err = verifyRelease(binPath, checksumPath, sigPath, binDigest)
	if err != nil {
		return "", fmt.Errorf("release verification failed: %w — the download is rejected; fix the release before updating", err)
	}
	log.Infof("release %s verified: Ed25519 signature, sha256 checksum and GitHub digest match", r.TagName)

	err = sanityCheckBinary(binPath, r.TagName)
	if err != nil {
		return "", err
	}
	log.Infof("release %s verified and downloaded to %s", r.TagName, binPath)
	return binPath, nil
}

// confirmUpdate decides whether to apply: --yes logs the changelog and
// proceeds; otherwise the notes open in a TUI viewer (non-tty aborts) and
// the user confirms.
func confirmUpdate(r *githubRelease) (bool, error) {
	if yesFlag {
		log.Infof("--yes given: skipping the changelog viewer and the confirmation prompt")
		log.Infof("changelog of %s:\n%s", r.TagName, changelogText(r))
		return true, nil
	}
	if !isTerminal(os.Stdin) {
		return false, errors.New("stdin is not a terminal: interactive confirmation is impossible; re-run with --yes to update without prompts")
	}
	if err := showChangelog(r.TagName, changelogText(r)); err != nil {
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

// resolveAssets returns the binary/checksum/signature URLs and the
// GitHub-provided binary digest when present.
func resolveAssets(r *githubRelease) (binURL, binDigest, checksumURL, sigURL string, err error) {
	asset := func(name string) (url, digest string, ok bool) {
		for i := range r.Assets {
			if r.Assets[i].Name == name {
				return r.Assets[i].BrowserDownloadURL, r.Assets[i].Digest, true
			}
		}
		return "", "", false
	}
	var ok bool
	if binURL, binDigest, ok = asset(releaseAssetBinary); !ok {
		return "", "", "", "", fmt.Errorf("release %s has no %s asset", r.TagName, releaseAssetBinary)
	}
	if checksumURL, _, ok = asset(releaseAssetChecksum); !ok {
		return "", "", "", "", fmt.Errorf("release %s has no %s asset", r.TagName, releaseAssetChecksum)
	}
	if sigURL, _, ok = asset(releaseAssetSignature); !ok {
		return "", "", "", "", fmt.Errorf("release %s has no %s asset", r.TagName, releaseAssetSignature)
	}
	return binURL, binDigest, checksumURL, sigURL, nil
}

// downloadReleaseFiles downloads the three assets into dir, reporting
// progress on bar when non-nil.
func downloadReleaseFiles(client *http.Client, dir, binURL, checksumURL, sigURL string, bar *wizard.BottomBar) (binPath, checksumPath, sigPath string, err error) {
	binPath = filepath.Join(dir, releaseAssetBinary)
	if err := downloadFile(client, binURL, binPath, bar, "Downloading app-listener"); err != nil {
		return "", "", "", err
	}
	checksumPath = filepath.Join(dir, releaseAssetChecksum)
	if err := downloadFile(client, checksumURL, checksumPath, bar, "Downloading checksum"); err != nil {
		return "", "", "", err
	}
	sigPath = filepath.Join(dir, releaseAssetSignature)
	if err := downloadFile(client, sigURL, sigPath, bar, "Downloading signature"); err != nil {
		return "", "", "", err
	}
	return binPath, checksumPath, sigPath, nil
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

// downloadFile downloads url into path atomically (temp file + rename),
// leaving mode 0700: hardened kernels refuse to exec without x bits and
// os.CreateTemp makes 0600; progress reports on bar when non-nil.
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

// verifyRelease checks the checksum's Ed25519 signature against the embedded
// key, then the checksum against the binary, then the GitHub asset digest.
func verifyRelease(binPath, checksumPath, sigPath, assetDigest string) error {
	pub, err := parsePublicKey(certificates.ReleasePublicKeyPEM)
	if err != nil {
		return fmt.Errorf("parsing the embedded release public key: %w", err)
	}
	return verifyReleaseWithKey(binPath, checksumPath, sigPath, assetDigest, pub)
}

// verifyReleaseWithKey is verifyRelease with an explicit public key.
func verifyReleaseWithKey(binPath, checksumPath, sigPath, assetDigest string, pub ed25519.PublicKey) error {
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

	expected, err := parseChecksum(checksum)
	if err != nil {
		return fmt.Errorf("parsing the signed checksum file: %w", err)
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

// parseChecksum extracts the validated sha256 hex token of a checksum file.
func parseChecksum(data []byte) (string, error) {
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return "", errors.New("the checksum file is empty")
	}
	if len(fields[0]) != 64 {
		return "", fmt.Errorf("expected a 64-hex-char sha256, got %q", fields[0])
	}
	if _, err := hex.DecodeString(fields[0]); err != nil {
		return "", fmt.Errorf("invalid sha256 hex %q: %w", fields[0], err)
	}
	return fields[0], nil
}

// sanityCheckBinary refuses non-ELF files whose --version misses the wanted tag.
func sanityCheckBinary(path, wantVersion string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	magic := make([]byte, 4)
	if _, readErr := io.ReadFull(f, magic); readErr != nil {
		return fmt.Errorf("reading %s: %w", path, readErr)
	}
	if string(magic) != "\x7fELF" {
		return fmt.Errorf("%s is not an ELF executable", path)
	}
	cmd := exec.CommandContext(context.Background(), path, "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("the downloaded binary does not run (--version failed: %v): %s", err, strings.TrimSpace(string(out)))
	}
	if !strings.Contains(string(out), wantVersion) {
		return fmt.Errorf("the downloaded binary reports %q, expected tag %q — refusing to deploy", strings.TrimSpace(string(out)), wantVersion)
	}
	return nil
}

// applyUpdate stops the daemon (verifying the lock state), atomically
// replaces the installed binary and brings the unit back to
// enabled-and-running.
func applyUpdate(binPath string) error {
	if err := systemd.DeployInstalledBinary(binPath); err != nil {
		return err
	}
	log.Info("update complete: the daemon is running the new binary")
	return nil
}
