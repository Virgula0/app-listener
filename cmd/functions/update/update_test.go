package update

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"debug/elf"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/certificates"
	"github.com/Virgula0/app-listener/internal/wizard"
)

// testKeyPair returns a fresh Ed25519 keypair and its PKIX PEM public key.
func testKeyPair(t *testing.T) (ed25519.PrivateKey, ed25519.PublicKey, []byte) {
	t.Helper()
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generating test key: %v", err)
	}
	der, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatalf("marshalling test public key: %v", err)
	}
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})
	return priv, pub, pemBytes
}

// signedChecksum writes a checksum file for data and its Ed25519 signature,
// returning the hex digest.
func signedChecksum(t *testing.T, priv ed25519.PrivateKey, name string, data []byte) (string, []byte, []byte) {
	t.Helper()
	sum := sha256.Sum256(data)
	digest := hex.EncodeToString(sum[:])
	checksum := []byte(digest + "  " + name + "\n")
	sig := ed25519.Sign(priv, checksum)
	return digest, checksum, sig
}

func TestParseChecksum(t *testing.T) {
	sum := strings.Repeat("ab", 32)
	cases := []struct {
		name     string
		data     string
		want     string
		wantName string
		wantErr  bool
	}{
		{"sha256sum format", sum + "  app-listener\n", sum, "app-listener", false},
		{"binary mode marker", sum + " *app-listener-arm64\n", sum, "app-listener-arm64", false},
		{"uppercase hex", strings.ToUpper(sum) + "  app-listener\n", strings.ToUpper(sum), "app-listener", false},
		{"empty", "", "", "", true},
		{"no name", sum + "\n", "", "", true},
		{"two lines", sum + "  app-listener\n" + sum + "  app-listener-arm64\n", "", "", true},
		{"too short", "deadbeef  app-listener\n", "", "", true},
		{"not hex", strings.Repeat("zz", 32) + "  app-listener\n", "", "", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, gotName, err := parseChecksum([]byte(c.data))
			if c.wantErr {
				if err == nil {
					t.Fatalf("parseChecksum(%q) = %q, want error", c.data, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseChecksum(%q): %v", c.data, err)
			}
			if got != c.want || gotName != c.wantName {
				t.Fatalf("parseChecksum(%q) = %q, %q, want %q, %q", c.data, got, gotName, c.want, c.wantName)
			}
		})
	}
}

func TestParsePublicKey(t *testing.T) {
	_, _, pemBytes := testKeyPair(t)

	if _, err := parsePublicKey(pemBytes); err != nil {
		t.Fatalf("parsing a valid PKIX PEM public key: %v", err)
	}

	// The embedded release key must parse.
	if _, err := parsePublicKey(certificates.ReleasePublicKeyPEM); err != nil {
		t.Fatalf("parsing the embedded release public key: %v", err)
	}

	if _, err := parsePublicKey([]byte("not a pem")); err == nil {
		t.Fatal("parsePublicKey accepted garbage")
	}

	// An RSA key must be rejected as the wrong type.
	if _, err := parsePublicKey(rsaPublicKeyPEM(t)); err == nil {
		t.Fatal("parsePublicKey accepted an RSA key")
	}
}

// rsaPublicKeyPEM returns a PKIX PEM-encoded RSA public key: parsePublicKey
// must reject it as the wrong key type.
func rsaPublicKeyPEM(t *testing.T) []byte {
	t.Helper()
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generating RSA key: %v", err)
	}
	der, err := x509.MarshalPKIXPublicKey(&rsaKey.PublicKey)
	if err != nil {
		t.Fatalf("marshalling RSA key: %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})
}

func TestVerifyRelease(t *testing.T) {
	priv, pub, _ := testKeyPair(t)

	bin := []byte("fake ELF payload")
	digest, checksum, sig := signedChecksum(t, priv, "app-listener", bin)

	dir := t.TempDir()
	write := func(name string, data []byte) string {
		t.Helper()
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, data, 0o600); err != nil {
			t.Fatalf("writing %s: %v", name, err)
		}
		return p
	}
	binPath := write("app-listener", bin)
	checksumPath := write("app-listener.sha256", checksum)
	sigPath := write("app-listener.sha256.sig", sig)

	// Valid release: signature, checksum and digest all match.
	if err := verifyReleaseWithKey(binPath, checksumPath, sigPath, "sha256:"+digest, "app-listener", pub); err != nil {
		t.Fatalf("verifyRelease on a valid release: %v", err)
	}
	// Digest absent (older GitHub API responses) is fine.
	if err := verifyReleaseWithKey(binPath, checksumPath, sigPath, "", "app-listener", pub); err != nil {
		t.Fatalf("verifyRelease without a digest: %v", err)
	}

	// Tampered binary.
	tampered := write("tampered", []byte("tampered payload"))
	if err := verifyReleaseWithKey(tampered, checksumPath, sigPath, "", "app-listener", pub); err == nil {
		t.Fatal("verifyRelease accepted a binary whose sha256 does not match the checksum")
	}

	// Signature over the wrong checksum file.
	_, otherChecksum, otherSig := signedChecksum(t, priv, "app-listener", []byte("other"))
	if err := verifyReleaseWithKey(binPath, write("other.sha256", otherChecksum), write("other.sha256.sig", otherSig), "", "app-listener", pub); err == nil {
		t.Fatal("verifyRelease accepted a signature over the wrong checksum file")
	}

	// Signature made with a different key.
	otherPriv, _, _ := testKeyPair(t)
	forgedSig := ed25519.Sign(otherPriv, checksum)
	if err := verifyReleaseWithKey(binPath, checksumPath, write("forged.sig", forgedSig), "", "app-listener", pub); err == nil {
		t.Fatal("verifyRelease accepted a signature from the wrong key")
	}

	// GitHub digest disagreeing with the signed checksum.
	if err := verifyReleaseWithKey(binPath, checksumPath, sigPath, "sha256:"+strings.Repeat("cd", 32), "app-listener", pub); err == nil {
		t.Fatal("verifyRelease accepted a GitHub digest that contradicts the signed checksum")
	}

	// Validly signed pair of another asset (the arm64 or GUI build) served under this name.
	_, armChecksum, armSig := signedChecksum(t, priv, "app-listener-arm64", bin)
	if err := verifyReleaseWithKey(binPath, write("arm.sha256", armChecksum), write("arm.sha256.sig", armSig), "", "app-listener", pub); err == nil {
		t.Fatal("verifyRelease accepted a checksum signed for another asset name")
	}
	if err := verifyReleaseWithKey(binPath, write("arm.sha256", armChecksum), write("arm.sha256.sig", armSig), "", "app-listener-arm64", pub); err != nil {
		t.Fatalf("verifyRelease on a valid arm64 release: %v", err)
	}

	// Garbage signature.
	if err := verifyReleaseWithKey(binPath, checksumPath, write("garbage.sig", []byte("garbage")), "", "app-listener", pub); err == nil {
		t.Fatal("verifyRelease accepted a garbage signature")
	}
}

func TestFetchReleases(t *testing.T) {
	// The API returns releases newest first; all of them must survive the
	// fetch untouched, preserving the server order.
	releases := []githubRelease{
		{TagName: "pre-20260201-abcdef0", Prerelease: true, PublishedAt: "2026-02-01T10:00:00Z"},
		{TagName: "v0.1.0", Prerelease: false, PublishedAt: "2025-12-01T10:00:00Z"},
		{TagName: "pre-20260101-abcdef0", Prerelease: true, PublishedAt: "2026-01-01T10:00:00Z"},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/repos/owner/repo/releases") {
			http.NotFound(w, r)
			return
		}
		if r.URL.Query().Get("per_page") != "100" {
			t.Errorf("expected per_page=100, got %q", r.URL.Query().Get("per_page"))
		}
		if r.Header.Get("User-Agent") != "app-listener-updater" {
			t.Errorf("unexpected User-Agent %q", r.Header.Get("User-Agent"))
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(releases)
	}))
	defer srv.Close()

	got, err := fetchReleases(srv.URL, "owner/repo", srv.Client())
	if err != nil {
		t.Fatalf("fetchReleases: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("fetchReleases returned %d releases, want 3", len(got))
	}
	if got[0].TagName != "pre-20260201-abcdef0" {
		t.Fatalf("first release is %q, want pre-20260201-abcdef0", got[0].TagName)
	}
}

func TestFetchReleasesHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()
	if _, err := fetchReleases(srv.URL, "owner/repo", srv.Client()); err == nil {
		t.Fatal("fetchReleases accepted an HTTP 500")
	}
}

func TestFilterChannel(t *testing.T) {
	releases := []githubRelease{
		{TagName: "pre-20260201-abcdef0", Prerelease: true},
		{TagName: "v0.1.0", Prerelease: false},
		{TagName: "pre-20260101-abcdef0", Prerelease: true},
		{TagName: "v1.2.3", Prerelease: false},
	}

	stable := filterChannel(releases, channelStable)
	if len(stable) != 2 || stable[0].TagName != "v0.1.0" || stable[1].TagName != "v1.2.3" {
		t.Fatalf("stable channel = %v, want only non-pre-releases", stable)
	}

	pre := filterChannel(releases, channelPreRelease)
	if len(pre) != 2 || pre[0].TagName != "pre-20260201-abcdef0" || pre[1].TagName != "pre-20260101-abcdef0" {
		t.Fatalf("pre-release channel = %v, want only pre-releases", pre)
	}
}

func TestPickLatestRelease(t *testing.T) {
	releases := []githubRelease{
		{TagName: "pre-20260101-aaaaaaa", Prerelease: true, PublishedAt: "2026-01-01T10:00:00Z"},
		{TagName: "v0.1.0", Prerelease: false, PublishedAt: "2026-01-02T10:00:00Z"},
		{TagName: "pre-20260101-ccccccc", Prerelease: true, PublishedAt: "not a timestamp"},
		{TagName: "pre-20260103-ddddddd", Prerelease: true, PublishedAt: "2026-01-03T10:00:00Z"},
	}
	latest, ok := pickLatestRelease(releases)
	if !ok {
		t.Fatal("pickLatestRelease found nothing")
	}
	if latest.TagName != "pre-20260103-ddddddd" {
		t.Fatalf("pickLatestRelease = %q, want pre-20260103-ddddddd", latest.TagName)
	}

	if _, ok := pickLatestRelease(nil); ok {
		t.Fatal("pickLatestRelease found a release in an empty list")
	}
	if _, ok := pickLatestRelease([]githubRelease{{TagName: "x", PublishedAt: "garbage"}}); ok {
		t.Fatal("pickLatestRelease accepted a release with a malformed timestamp")
	}
}

func TestNewerThanInstalled(t *testing.T) {
	releases := []githubRelease{
		{TagName: "pre-20260101-aaaaaaa", Prerelease: true, PublishedAt: "2026-01-01T10:00:00Z"},
		{TagName: "pre-20260102-bbbbbbb", Prerelease: true, PublishedAt: "2026-01-02T10:00:00Z"},
	}
	latest := releases[1]

	cases := []struct {
		installed string
		want      bool
	}{
		// Non pre-release builds (dev builds, unknown) are always older.
		{"v0.1.0", true},
		{"unknown", true},
		{"pre-test", true},
		// Same tag is never newer.
		{"pre-20260102-bbbbbbb", false},
		// Tag present in the list: compare published_at.
		{"pre-20260101-aaaaaaa", true},
		// Tag older than the fetched page: update.
		{"pre-20251201-ccccccc", true},
	}
	for _, c := range cases {
		got := newerThanInstalled(c.installed, &latest, releases)
		if got != c.want {
			t.Fatalf("newerThanInstalled(%q) = %v, want %v", c.installed, got, c.want)
		}
	}
}

func TestNewerThanStable(t *testing.T) {
	latest := githubRelease{TagName: "v1.2.0"}

	cases := []struct {
		installed string
		want      bool
	}{
		// Same tag is never newer.
		{"v1.2.0", false},
		// Older stable: update.
		{"v1.1.0", true},
		{"v1.2.0-alpha", true},
		// Newer stable: never downgrade.
		{"v1.2.1", false},
		{"v2.0.0", false},
		// Non-stable versions (pre-release builds, dev, unknown) need the
		// latest stable — "update if tags differ".
		{"pre-20260101-aaaaaaa", true},
		{"pre-test", true},
		{"unknown", true},
		{"", true},
	}
	for _, c := range cases {
		got := newerThanStable(c.installed, &latest)
		if got != c.want {
			t.Fatalf("newerThanStable(%q) = %v, want %v", c.installed, got, c.want)
		}
	}
}

func TestParseStableVersion(t *testing.T) {
	cases := []struct {
		tag   string
		want  stableVersion
		valid bool
	}{
		{"v0.1.0", stableVersion{0, 1, 0}, true},
		{"v1.2.3", stableVersion{1, 2, 3}, true},
		{"1.2.3", stableVersion{1, 2, 3}, true},
		{"v10.20.30", stableVersion{10, 20, 30}, true},
		// Non-stable and malformed tags.
		{"pre-20260101-aaaaaaa", stableVersion{}, false},
		{"v1.2", stableVersion{}, false},
		{"v1.2.x", stableVersion{}, false},
		{"v1.beta.3", stableVersion{}, false},
		{"v1.2.3.4", stableVersion{}, false},
		{"", stableVersion{}, false},
		{"v0.1", stableVersion{}, false},
		{"v0.1.0-rc1", stableVersion{0, 1, 0}, false},
	}
	for _, c := range cases {
		got, ok := parseStableVersion(c.tag)
		if ok != c.valid {
			t.Errorf("parseStableVersion(%q) ok = %v, want %v", c.tag, ok, c.valid)
			continue
		}
		if ok && got != c.want {
			t.Errorf("parseStableVersion(%q) = %+v, want %+v", c.tag, got, c.want)
		}
	}
}

func TestValidateChannel(t *testing.T) {
	if err := validateChannel("stable"); err != nil {
		t.Fatalf("validateChannel(stable) = %v", err)
	}
	if err := validateChannel("pre-release"); err != nil {
		t.Fatalf("validateChannel(pre-release) = %v", err)
	}
	for _, ch := range []string{"", "beta", "main", "Prerelease"} {
		if err := validateChannel(ch); err == nil {
			t.Errorf("validateChannel(%q) accepted an invalid channel", ch)
		}
	}
}

func TestDownloadFileMode0700(t *testing.T) {
	content := []byte("binary payload")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", fmt.Sprint(len(content)))
		w.Write(content)
	}))
	defer srv.Close()

	dst := filepath.Join(t.TempDir(), "app-listener")
	if err := downloadFile(srv.Client(), srv.URL, dst, nil, "Downloading app-listener"); err != nil {
		t.Fatalf("downloadFile: %v", err)
	}
	info, err := os.Stat(dst)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0o700 {
		t.Fatalf("downloaded file mode = %o, want 700 (hardened kernels refuse to exec 0600 even for root)", perm)
	}
}

func TestDownloadFileHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusInternalServerError)
	}))
	defer srv.Close()

	if err := downloadFile(srv.Client(), srv.URL, filepath.Join(t.TempDir(), "x"), nil, "label"); err == nil {
		t.Fatal("downloadFile accepted an HTTP 500")
	}
}

func TestProgressReader(t *testing.T) {
	payload := strings.Repeat("0123456789", 10)
	body := strings.NewReader(payload)

	var read int64
	err := wizard.WithBottomBar(func(bar *wizard.BottomBar) error {
		p := &progressReader{reader: body, total: int64(len(payload)), bar: bar, label: "Downloading"}
		out, err := io.ReadAll(p)
		if err != nil {
			return err
		}
		read = p.read
		if string(out) != payload {
			t.Fatalf("progressReader returned %d bytes, want the full %d-byte payload", len(out), len(payload))
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WithBottomBar: %v", err)
	}
	if read != int64(len(payload)) {
		t.Fatalf("progressReader counted %d bytes, want %d", read, len(payload))
	}
}

func TestReadInstalledVersion(t *testing.T) {
	dir := t.TempDir()
	bin := filepath.Join(dir, "app-listener")
	script := "#!/bin/sh\necho 'app-listener version pre-20260102-abcdef0'\n"
	if err := os.WriteFile(bin, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}

	if got := readInstalledVersion(bin); got != "pre-20260102-abcdef0" {
		t.Fatalf("readInstalledVersion = %q, want the embedded tag", got)
	}

	// A non-executable binary reports the unknown version.
	if err := os.WriteFile(bin, []byte("not executable"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := readInstalledVersion(bin); got != "unknown" {
		t.Fatalf("readInstalledVersion on a non-executable = %q, want unknown", got)
	}
}

func TestAssetsFor(t *testing.T) {
	amd, err := assetsFor("amd64")
	if err != nil {
		t.Fatal(err)
	}
	// Older binaries' `update` fetches these exact names: renaming them strands every amd64 install.
	if amd.binary != "app-listener" || amd.checksum != "app-listener.sha256" || amd.signature != "app-listener.sha256.sig" || amd.machine != elf.EM_X86_64 {
		t.Fatalf("amd64 assets = %+v", amd)
	}
	arm, err := assetsFor("arm64")
	if err != nil {
		t.Fatal(err)
	}
	if arm.binary != "app-listener-arm64" || arm.checksum != "app-listener-arm64.sha256" || arm.signature != "app-listener-arm64.sha256.sig" || arm.machine != elf.EM_AARCH64 {
		t.Fatalf("arm64 assets = %+v", arm)
	}
	if _, err := assetsFor("riscv64"); err == nil {
		t.Fatal("assetsFor accepted an architecture with no published release")
	}
}

func TestResolveAssetsPicksArch(t *testing.T) {
	r := &githubRelease{TagName: "v1.0.0"}
	for _, n := range []string{"app-listener", "app-listener.sha256", "app-listener.sha256.sig",
		"app-listener-arm64", "app-listener-arm64.sha256", "app-listener-arm64.sha256.sig"} {
		r.Assets = append(r.Assets, struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
			Digest             string `json:"digest"`
		}{Name: n, BrowserDownloadURL: "https://example.invalid/" + n, Digest: "sha256:" + n})
	}
	arm, _ := assetsFor("arm64")
	bin, digest, sum, sig, err := resolveAssets(r, arm)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(bin, "/app-listener-arm64") || digest != "sha256:app-listener-arm64" ||
		!strings.HasSuffix(sum, "/app-listener-arm64.sha256") || !strings.HasSuffix(sig, "/app-listener-arm64.sha256.sig") {
		t.Fatalf("resolveAssets(arm64) = %s %s %s %s", bin, digest, sum, sig)
	}

	r.Assets = r.Assets[:3] // a release published before arm64 builds existed
	if _, _, _, _, err := resolveAssets(r, arm); err == nil {
		t.Fatal("resolveAssets(arm64) succeeded on a release without arm64 assets")
	}
}

func TestSanityCheckBinaryRefusesOtherMachine(t *testing.T) {
	self, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	f, err := elf.Open(self)
	if err != nil {
		t.Fatal(err)
	}
	machine := f.Machine
	f.Close()

	other := elf.EM_AARCH64
	if machine == elf.EM_AARCH64 {
		other = elf.EM_X86_64
	}
	err = sanityCheckBinary(self, "v-never", other)
	if err == nil || !strings.Contains(err.Error(), "built for") {
		t.Fatalf("sanityCheckBinary with a foreign machine = %v, want a machine refusal", err)
	}

	notELF := filepath.Join(t.TempDir(), "app-listener")
	if err := os.WriteFile(notELF, []byte("#!/bin/sh\necho v1.0.0\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := sanityCheckBinary(notELF, "v1.0.0", machine); err == nil {
		t.Fatal("sanityCheckBinary accepted a non-ELF file")
	}
}
