package install

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/internal/fscrypt"
)

// TestAskEncryptionSkipsNeedEncryptionFalse verifies that a resource
// declared need_encryption: false in the config never triggers the
// encryption question: the directory is not added to toEncrypt, the config
// text is not rewritten, and no interactive prompt is shown (a prompt would
// block on stdin and fail this test instead).
func TestAskEncryptionSkipsNeedEncryptionFalse(t *testing.T) {
	dir := t.TempDir()

	cfgText := "[watch]\npath = " + dir + "\nneed_encryption: false\n"
	cfg, err := validateConfigText(cfgText)
	if err != nil {
		t.Fatalf("parsing config: %v", err)
	}

	text, toEncrypt, err := askEncryption(fscrypt.New(), cfgText, cfg)
	if err != nil {
		t.Fatalf("askEncryption: %v", err)
	}
	if len(toEncrypt) != 0 {
		t.Errorf("toEncrypt = %v, want none for need_encryption: false", toEncrypt)
	}
	if text != cfgText {
		t.Errorf("config text must be left untouched, got:\n%s", text)
	}
	if strings.Contains(text, "need_encryption: true") {
		t.Errorf("config text must keep need_encryption: false, got:\n%s", text)
	}
}

// TestAskFilesystemsReadyNoPanic is a regression test for the preflight
// added with the fscrypt setup check: statting a real directory must not
// panic (os.Stat returns *syscall.Stat_t, and the preflight must accept
// exactly that type). The helper must return either nil (host filesystem
// is already initialized for fscrypt) or the classified setup error.
func TestAskFilesystemsReadyNoPanic(t *testing.T) {
	dir := t.TempDir()

	cfgText := "[watch]\npath = " + dir + "\nneed_encryption: true\n"
	cfg, err := validateConfigText(cfgText)
	if err != nil {
		t.Fatalf("parsing config: %v", err)
	}

	if err := askFilesystemsReady(fscrypt.New(), cfg); err != nil &&
		!strings.Contains(err.Error(), "fscrypt setup") {
		t.Errorf("unexpected error from preflight: %v", err)
	}
}

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
	got := parseSectionWhitelist(conf, "/a")
	want := []string{"/usr/bin/one", "/usr/bin/two"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseSectionWhitelist(/a) = %v, want %v", got, want)
	}

	if got := parseSectionWhitelist(conf, "/b"); !reflect.DeepEqual(got, []string{"/usr/bin/three"}) {
		t.Errorf("parseSectionWhitelist(/b) = %v, want [/usr/bin/three]", got)
	}

	if got := parseSectionWhitelist(conf, "/missing"); got != nil {
		t.Errorf("parseSectionWhitelist(/missing) = %v, want nil", got)
	}
}

// TestLiveEmptyWhitelistRejected is the regression test for the live
// refresh's fail-closed contract: a re-scan that comes back empty for a
// previously-populated encrypted resource must be refused (the vault is
// locked, or this installer binary is not the running daemon's) instead of
// persisting a silently shrunk whitelist.
func TestLiveEmptyWhitelistRejected(t *testing.T) {
	if !liveEmptyWhitelistRejected(true, true, 0, 3) {
		t.Error("live empty re-scan of a previously-populated encrypted resource must be rejected")
	}
	if liveEmptyWhitelistRejected(true, false, 0, 3) {
		t.Error("the stopped flow has its own unlock/re-lock contract: not a live rejection")
	}
	if liveEmptyWhitelistRejected(false, true, 0, 3) {
		t.Error("non-encrypted resources may legitimately end up empty (uninstalled binary)")
	}
	if liveEmptyWhitelistRejected(true, true, 2, 3) {
		t.Error("a non-empty re-scan must never be rejected")
	}
	if liveEmptyWhitelistRejected(true, true, 0, 0) {
		t.Error("a section that was already empty stays in deny-everything mode: nothing to protect")
	}
}

// TestApplyLiveRefreshDeliversReload is the regression test for the missing
// SIGHUP delivery: the first live implementation patched daemon.conf on
// disk and returned WITHOUT delivering the change to the running daemon,
// which kept enforcing the old whitelist while journalctl stayed silent.
// Live mode must deliver the reload exactly when the config was written,
// and skip it when the refreshed config matches the running one.
func TestApplyLiveRefreshDeliversReload(t *testing.T) {
	orig := deliverReload
	defer func() { deliverReload = orig }()

	var calls []bool
	deliverReload = func(configChanged bool) error {
		calls = append(calls, configChanged)
		return nil
	}

	if err := applyLiveRefresh(true); err != nil {
		t.Fatalf("applyLiveRefresh(true): %v", err)
	}
	if len(calls) != 1 || calls[0] != true {
		t.Fatalf("a written config must be delivered via reload, calls=%v", calls)
	}

	if err := applyLiveRefresh(false); err != nil {
		t.Fatalf("applyLiveRefresh(false): %v", err)
	}
	if len(calls) != 1 {
		t.Fatalf("an unchanged config must not trigger a reload, calls=%v", calls)
	}
}
