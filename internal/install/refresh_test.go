package install

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

// TestParseSectionWhitelist covers the helper backing the live refresh's
// empty-whitelist safety check: only the binary directives of the requested
// [watch] section are returned, never directives of a following section.
func TestParseSectionWhitelist(t *testing.T) {
	conf := `[watch /a]
/usr/bin/one
need_encryption: true
/usr/bin/two READ,WRITE

[watch /b]
/usr/bin/three
`
	got := ParseSectionWhitelist(conf, "/a")
	want := []string{"/usr/bin/one", "/usr/bin/two"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseSectionWhitelist(/a) = %v, want %v", got, want)
	}

	if got := ParseSectionWhitelist(conf, "/b"); !reflect.DeepEqual(got, []string{"/usr/bin/three"}) {
		t.Errorf("ParseSectionWhitelist(/b) = %v, want [/usr/bin/three]", got)
	}

	if got := ParseSectionWhitelist(conf, "/missing"); got != nil {
		t.Errorf("ParseSectionWhitelist(/missing) = %v, want nil", got)
	}
}

// TestLiveEmptyWhitelistRejected is the regression test for the live
// refresh's fail-closed contract: a re-scan that comes back empty for a
// previously-populated encrypted resource must be refused (the vault is
// locked, or this installer binary is not the running daemon's) instead of
// persisting a silently shrunk whitelist.
func TestLiveEmptyWhitelistRejected(t *testing.T) {
	if !LiveEmptyWhitelistRejected(true, true, 0, 3) {
		t.Error("live empty re-scan of a previously-populated encrypted resource must be rejected")
	}
	if LiveEmptyWhitelistRejected(true, false, 0, 3) {
		t.Error("the stopped flow has its own unlock/re-lock contract: not a live rejection")
	}
	if LiveEmptyWhitelistRejected(false, true, 0, 3) {
		t.Error("non-encrypted resources may legitimately end up empty (uninstalled binary)")
	}
	if LiveEmptyWhitelistRejected(true, true, 2, 3) {
		t.Error("a non-empty re-scan must never be rejected")
	}
	if LiveEmptyWhitelistRejected(true, true, 0, 0) {
		t.Error("a section that was already empty stays in deny-everything mode: nothing to protect")
	}
}

func TestRefreshCatalogReportsVersionBump(t *testing.T) {
	home := t.TempDir()
	root := filepath.Join(home, ".config", "discord")
	for _, p := range []string{filepath.Join(root, "sentry"), filepath.Join(root, "0.0.2")} {
		if err := os.MkdirAll(p, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	newBin := filepath.Join(root, "0.0.2", "Discord")
	if err := os.WriteFile(newBin, []byte("x"), 0o755); err != nil {
		t.Fatal(err)
	}
	oldBin := filepath.Join(root, "0.0.1", "Discord")
	conf := "[watch " + root + "]\nneed_encryption: false\n" + oldBin + "\n"
	confPath := filepath.Join(t.TempDir(), "daemon.conf")
	if err := os.WriteFile(confPath, []byte(conf), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := daemonconfig.Load(confPath)
	if err != nil {
		t.Fatal(err)
	}
	plain := func(_ *daemonconfig.Resource, expand func() []BinaryRule) ([]BinaryRule, bool, error) {
		return expand(), false, nil
	}
	text, changes, err := RefreshCatalog(conf, cfg, []User{{Name: "u", Home: home}},
		RefreshOptions{Live: true, Scan: plain})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(text, newBin) || strings.Contains(text, oldBin) {
		t.Errorf("the refreshed section must name the new version only:\n%s", text)
	}
	want := []SectionChange{{Section: root, Admitted: []string{newBin}, Dropped: []string{oldBin}}}
	if !reflect.DeepEqual(changes, want) {
		t.Errorf("changes = %+v, want %+v", changes, want)
	}
}

func TestRefreshCatalogUserSectionsOnly(t *testing.T) {
	dir := t.TempDir()
	conf := "[watch " + dir + "]\nneed_encryption: false\n/usr/bin/true\n"
	confPath := filepath.Join(t.TempDir(), "daemon.conf")
	if err := os.WriteFile(confPath, []byte(conf), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := daemonconfig.Load(confPath)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = RefreshCatalog(conf, cfg, nil, RefreshOptions{Live: true})
	if !errors.Is(err, ErrNoCatalogMatch) {
		t.Fatalf("a config without catalog sections: err = %v, want ErrNoCatalogMatch", err)
	}
}
