// The `app-listener update` command keeps an installed daemon binary up to
// date with the latest GitHub pre-release. It is fully non-interactive:
//
//  1. reads the installed binary's embedded version (--version)
//  2. lists the repository pre-releases from the GitHub API and picks the
//     newest one (by published_at)
//  3. skips when the installed tag is not older
//  4. downloads the release binary, its sha256 checksum and the Ed25519
//     signature of that checksum
//  5. verifies the signature against the public key embedded in the binary
//     (certificates/app-listener-release.pub), the checksum against the
//     downloaded binary, and the GitHub-provided asset digest when present
//  6. replaces /usr/local/sbin/app-listener atomically and restarts the
//     systemd daemon, mirroring the installer's stop/start contract
//
// The signing private key never lives in this repository: it is stored as
// the RELEASE_SIGNING_KEY GitHub Actions secret only.
package update

import (
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
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/certificates"
	"github.com/Virgula0/app-listener/internal/systemd"
)

const (
	// githubAPIBase is the GitHub REST API endpoint used to list releases.
	githubAPIBase = "https://api.github.com"
	// releaseAssetBinary / Checksum / Signature are the release asset
	// names produced by .github/workflows/release.yml.
	releaseAssetBinary    = "app-listener"
	releaseAssetChecksum  = "app-listener.sha256"
	releaseAssetSignature = "app-listener.sha256.sig"
	// preVersionPattern matches the release tag format
	// pre-YYYYMMDD-<sha7> produced by the release workflow.
	preVersionPattern = `^pre-(\d{8})-([0-9a-f]{7})$`
	// updateHTTPTimeout bounds a single GitHub API/download request.
	updateHTTPTimeout = 60 * time.Second
)

var preVersionRe = regexp.MustCompile(preVersionPattern)

// repoFlag overrides the GitHub repository used by the updater.
var repoFlag string

func init() {
	UpdateCmd.Flags().StringVar(&repoFlag, "repo", "Virgula0/app-listener",
		"GitHub repository (owner/name) to fetch pre-releases from")
}

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Non-interactive self-update from the latest GitHub pre-release",
	Long: `Non-interactive, root-only self-updater for the daemon mode.

The command checks the latest pre-release of the repository (the
pre-YYYYMMDD-<sha> builds produced by the release workflow), compares it
with the version embedded in the installed binary at /usr/local/sbin/
app-listener, and when a newer one exists:

  1. downloads the release binary and its sha256 checksum
  2. verifies the Ed25519 signature of the checksum against the public key
     compiled into this binary (certificates/app-listener-release.pub),
     then the checksum against the downloaded binary, and the
     GitHub-provided asset digest when present — a failed verification
     aborts the update
  3. replaces the installed binary atomically
  4. restarts the systemd daemon (it is stopped first so the watch
     directories are verified locked again before the new binary starts)

The command is fully non-interactive: it logs its decisions and exits.
Use --repo to point at a different GitHub repository.`,
	Args: cobra.NoArgs,
	RunE: runUpdate,
}

// githubRelease mirrors the fields of the GitHub releases API relevant to
// the updater.
type githubRelease struct {
	TagName     string `json:"tag_name"`
	Prerelease  bool   `json:"prerelease"`
	PublishedAt string `json:"published_at"`
	Assets      []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
		Digest             string `json:"digest"`
	} `json:"assets"`
}

// runUpdate drives the whole update flow. It never prompts.
func runUpdate(cmd *cobra.Command, args []string) error {
	if os.Geteuid() != 0 {
		return errors.New("update must be run as root: sudo app-listener update")
	}
	if _, err := os.Lstat(systemd.InstallBinaryPath); err != nil {
		return fmt.Errorf("the daemon is not installed (no %s): run `app-listener install` first", systemd.InstallBinaryPath)
	}

	installed := readInstalledVersion(systemd.InstallBinaryPath)
	log.Infof("installed version: %s", installed)

	releases, err := fetchPreReleases(githubAPIBase, repoFlag, updateHTTPClient())
	if err != nil {
		return err
	}
	latest, ok := pickLatestPreRelease(releases)
	if !ok {
		log.Info("no pre-releases found: nothing to update")
		return nil
	}
	log.Infof("latest pre-release: %s (published %s)", latest.TagName, latest.PublishedAt)

	if !newerThanInstalled(installed, latest, releases) {
		log.Infof("installed version %s is up to date", installed)
		return nil
	}

	binURL, binDigest, checksumURL, sigURL, err := resolveAssets(latest)
	if err != nil {
		return err
	}

	tmp, err := os.MkdirTemp("", "app-listener-update-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	client := updateHTTPClient()
	binPath, checksumPath, sigPath, err := downloadReleaseFiles(client, tmp, binURL, checksumURL, sigURL)
	if err != nil {
		return err
	}

	if err := verifyRelease(binPath, checksumPath, sigPath, binDigest); err != nil {
		return fmt.Errorf("release verification failed: %w — the download is rejected; fix the release before updating", err)
	}
	log.Infof("release %s verified: Ed25519 signature, sha256 checksum and GitHub digest match", latest.TagName)

	if err := sanityCheckBinary(binPath, latest.TagName); err != nil {
		return err
	}

	return applyUpdate(binPath)
}

// resolveAssets finds the binary, checksum and signature assets of a
// release, returning their download URLs and the GitHub-provided binary
// digest when present.
func resolveAssets(r githubRelease) (binURL, binDigest, checksumURL, sigURL string, err error) {
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

// downloadReleaseFiles downloads the three release assets into dir and
// returns their paths.
func downloadReleaseFiles(client *http.Client, dir, binURL, checksumURL, sigURL string) (binPath, checksumPath, sigPath string, err error) {
	binPath = filepath.Join(dir, releaseAssetBinary)
	if err := downloadFile(client, binURL, binPath); err != nil {
		return "", "", "", err
	}
	checksumPath = filepath.Join(dir, releaseAssetChecksum)
	if err := downloadFile(client, checksumURL, checksumPath); err != nil {
		return "", "", "", err
	}
	sigPath = filepath.Join(dir, releaseAssetSignature)
	if err := downloadFile(client, sigURL, sigPath); err != nil {
		return "", "", "", err
	}
	return binPath, checksumPath, sigPath, nil
}

// readInstalledVersion runs the installed binary with --version and returns
// the reported version. A binary too old to know --version is treated as an
// unknown version, which the comparer maps to "older than any
// pre-release".
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

// fetchPreReleases returns the pre-releases of the repository, newest
// first, from the GitHub releases API.
func fetchPreReleases(apiBase, repo string, client *http.Client) ([]githubRelease, error) {
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
	var pre []githubRelease
	for i := range releases {
		if releases[i].Prerelease {
			pre = append(pre, releases[i])
		}
	}
	return pre, nil
}

// pickLatestPreRelease returns the pre-release with the most recent
// published_at timestamp.
func pickLatestPreRelease(releases []githubRelease) (githubRelease, bool) {
	var latest githubRelease
	found := false
	for i := range releases {
		ts, err := time.Parse(time.RFC3339, releases[i].PublishedAt)
		if err != nil {
			continue
		}
		if !found || ts.After(latestTime(latest)) {
			latest = releases[i]
			found = true
		}
	}
	return latest, found
}

// latestTime parses the published_at of a release (zero time on failure,
// so malformed timestamps never win the comparison).
func latestTime(r githubRelease) time.Time {
	ts, err := time.Parse(time.RFC3339, r.PublishedAt)
	if err != nil {
		return time.Time{}
	}
	return ts
}

// newerThanInstalled reports whether the latest pre-release is newer than
// the installed version. The installed version is either a pre-YYYYMMDD-<sha>
// tag, or anything else (dev builds, unknown) which is always older.
func newerThanInstalled(installed string, latest githubRelease, releases []githubRelease) bool {
	if !preVersionRe.MatchString(installed) {
		log.Infof("installed version %q is not a pre-release build: the latest pre-release is newer", installed)
		return true
	}
	if installed == latest.TagName {
		return false
	}
	for i := range releases {
		if releases[i].TagName == installed {
			return latestTime(latest).After(latestTime(releases[i]))
		}
	}
	// The installed tag is no longer among the listed pre-releases: it is
	// older than the newest page of releases we fetched.
	log.Infof("installed tag %s not found among the fetched pre-releases: updating", installed)
	return true
}

// updateHTTPClient returns an HTTP client with a bounded per-request
// timeout, shared by API calls and asset downloads.
func updateHTTPClient() *http.Client {
	return &http.Client{Timeout: updateHTTPTimeout}
}

// downloadFile downloads url into path (atomically via a temp file).
func downloadFile(client *http.Client, url, path string) error {
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
	f, err := os.CreateTemp(filepath.Dir(path), ".download-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		return fmt.Errorf("downloading %s: %w", url, err)
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}

// verifyRelease checks, in order: the Ed25519 signature of the checksum
// file against the embedded public key, the sha256 checksum against the
// downloaded binary, and — when GitHub provides one — the asset digest.
func verifyRelease(binPath, checksumPath, sigPath, assetDigest string) error {
	pub, err := parsePublicKey(certificates.ReleasePublicKeyPEM)
	if err != nil {
		return fmt.Errorf("parsing the embedded release public key: %w", err)
	}
	return verifyReleaseWithKey(binPath, checksumPath, sigPath, assetDigest, pub)
}

// verifyReleaseWithKey is the verification core of verifyRelease with an
// explicit public key.
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

// parsePublicKey parses the PEM-encoded SPKI Ed25519 public key embedded
// in the binary.
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

// parseChecksum extracts the first hex token of the sha256sum-style
// checksum file.
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

// sanityCheckBinary refuses to deploy a file that is not an ELF executable
// and whose own --version does not report the release tag being installed.
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

// applyUpdate stops the daemon (verifying the resources lock), replaces
// the installed binary atomically, recreates the PATH symlink if needed and
// brings the daemon back to enabled-and-running.
func applyUpdate(binPath string) error {
	if err := systemd.StopDaemonIfRunning(); err != nil {
		return err
	}

	if err := replaceInstalledBinary(binPath, systemd.InstallBinaryPath); err != nil {
		return err
	}
	if err := systemd.EnsureBinSymlink(); err != nil {
		return err
	}

	if err := systemd.EnableAndVerify(false); err != nil {
		return err
	}
	log.Info("update complete: the daemon is running the new binary")
	return nil
}

// replaceInstalledBinary swaps the installed binary with the verified
// download, atomically and with the same 0700 mode the installer uses.
func replaceInstalledBinary(binPath, dstPath string) error {
	tmp, err := os.CreateTemp(filepath.Dir(dstPath), ".app-listener-new-*")
	if err != nil {
		return fmt.Errorf("staging the new binary: %w", err)
	}
	defer os.Remove(tmp.Name())
	in, err := os.Open(binPath)
	if err != nil {
		return err
	}
	if _, err := io.Copy(tmp, in); err != nil {
		in.Close()
		tmp.Close()
		return err
	}
	in.Close()
	if err := tmp.Chmod(0o700); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmp.Name(), dstPath); err != nil {
		return fmt.Errorf("replacing %s: %w", dstPath, err)
	}
	log.Infof("replaced %s with the new binary", dstPath)
	return nil
}
