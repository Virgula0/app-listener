package install

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	inst "github.com/Virgula0/app-listener/internal/install"
)

// TestGroupCandidates verifies the installer collapses the per-path
// candidates of one resource (issue #51) into a single TUI group, ordered
// by first appearance, and that selecting a group yields every path.
func TestGroupCandidates(t *testing.T) {
	u := inst.User{Name: "alice", Home: "/home/alice"}
	cands := []inst.Candidate{
		{User: u, Entry: inst.CandidateDir{Name: "Claude Code"}, Path: "/home/alice/.claude"},
		{User: u, Entry: inst.CandidateDir{Name: "SSH"}, Path: "/home/alice/.ssh"},
		{User: u, Entry: inst.CandidateDir{Name: "Claude Code"}, Path: "/home/alice/.config/claude"},
	}
	groups := groupCandidates(cands)
	if len(groups) != 2 {
		t.Fatalf("got %d groups, want 2: %+v", len(groups), groups)
	}
	if len(groups[0].candidates) != 2 || groups[0].candidates[1].Path != "/home/alice/.config/claude" {
		t.Errorf("Claude Code group not merged in order: %+v", groups[0].candidates)
	}
	if !strings.Contains(groups[0].label, "[2 locations]") || !strings.Contains(groups[0].label, "(user alice)") {
		t.Errorf("group label = %q", groups[0].label)
	}
	if len(groups[1].candidates) != 1 || strings.Contains(groups[1].label, "locations]") {
		t.Errorf("single-path group should not show a location count: %q", groups[1].label)
	}
}

// TestEnsureInstalledBinary is the regression test for issue #44: the wizard
// never recompiles. A missing binary is deployed by copying the running
// executable into place; an already-installed binary is left untouched.
func TestEnsureInstalledBinary(t *testing.T) {
	orig := installedBinaryPath
	defer func() { installedBinaryPath = orig }()

	dir := t.TempDir()
	installedBinaryPath = filepath.Join(dir, "app-listener")

	// Missing: the running test binary is copied into place, no error.
	if err := ensureInstalledBinary(); err != nil {
		t.Fatalf("ensureInstalledBinary with no binary deployed: %v", err)
	}
	deployed, err := os.Stat(installedBinaryPath)
	if err != nil {
		t.Fatalf("binary was not deployed: %v", err)
	}
	if deployed.Mode().Perm() != 0o700 {
		t.Fatalf("deployed binary mode = %o, want 700", deployed.Mode().Perm())
	}

	// Already installed: left untouched (sentinel content survives).
	if err := os.WriteFile(installedBinaryPath, []byte("SENTINEL"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := ensureInstalledBinary(); err != nil {
		t.Fatalf("ensureInstalledBinary on a deployed binary: %v", err)
	}
	got, err := os.ReadFile(installedBinaryPath)
	if err != nil || string(got) != "SENTINEL" {
		t.Fatalf("existing binary was overwritten: content = %q (err %v)", got, err)
	}
}

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

// groupedDiscordConf builds a grouped [watch <root>] config for a synthetic
// Discord layout under home, returning the config text and its parsed form.
// The whitelist glob (%HOME%/.config/discord/*/Discord) matches app-1.0.0.
func groupedDiscordConf(t *testing.T, home string) (string, *daemonconfig.Config) {
	t.Helper()
	root := filepath.Join(home, ".config", "discord")
	lsDir := filepath.Join(root, "Local Storage")
	cookies := filepath.Join(root, "Cookies")
	appBin := filepath.Join(root, "app-1.0.0", "Discord")
	for _, d := range []string{lsDir, filepath.Dir(appBin)} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for _, f := range []string{cookies, appBin} {
		if err := os.WriteFile(f, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	confText := inst.GenerateConf([]inst.Section{{
		Path:            root,
		Allow:           []inst.BinaryRule{{Path: "/usr/bin/stale-entry"}},
		Encrypt:         true,
		ExtraWatchPaths: []string{lsDir, cookies},
	}})
	cfg, err := validateConfigText(confText)
	if err != nil {
		t.Fatalf("generated grouped config does not parse: %v", err)
	}
	return confText, cfg
}

// TestPatchCatalogSectionGroupedConfig is the regression test for the
// pacman-hook breaker: `install --update-catalog-only` addressed a grouped
// section's whitelist by its first watch sub-path, which has no [watch]
// header — SetSectionWhitelist failed "section not found" and, in the
// non-live flow, left the daemon stopped. The refresh must address the
// section by its encryption root, re-expand the shared whitelist, and
// preserve the `watch:` group structure.
func TestPatchCatalogSectionGroupedConfig(t *testing.T) {
	home := t.TempDir()
	confText, cfg := groupedDiscordConf(t, home)
	root := filepath.Join(home, ".config", "discord")
	appBin := filepath.Join(root, "app-1.0.0", "Discord")

	groups := cfg.EncryptionGroups()
	if len(groups) != 1 {
		t.Fatalf("want 1 encryption group, got %d", len(groups))
	}
	r := groups[0]
	if r.EncryptionRootOrPath() == r.Path {
		t.Fatalf("test needs a grouped resource: root %q == watch path %q", r.EncryptionRootOrPath(), r.Path)
	}

	users := []inst.User{{Name: "tester", Home: home}}
	updated, patched, err := patchCatalogSection(fscrypt.New(), confText, r, users, false)
	if err != nil {
		t.Fatalf("patchCatalogSection on a grouped config: %v", err)
	}
	if !patched {
		t.Fatal("the Discord catalog entry must match its own encryption root")
	}
	if strings.Contains(updated, "/usr/bin/stale-entry") {
		t.Errorf("stale whitelist entry survived the refresh:\n%s", updated)
	}
	if !strings.Contains(updated, appBin) {
		t.Errorf("re-expanded whitelist missing %s:\n%s", appBin, updated)
	}
	if !strings.Contains(updated, "watch: "+filepath.Join(root, "Local Storage")) ||
		!strings.Contains(updated, "watch: "+filepath.Join(root, "Cookies")) {
		t.Errorf("group watch directives were erased:\n%s", updated)
	}

	confPath := filepath.Join(t.TempDir(), "daemon.conf")
	if err := os.WriteFile(confPath, []byte(updated), 0o600); err != nil {
		t.Fatal(err)
	}
	reparsed, err := daemonconfig.Load(confPath)
	if err != nil {
		t.Fatalf("refreshed grouped config does not parse: %v", err)
	}
	if len(reparsed.Resources) != 2 {
		t.Fatalf("want 2 grouped resources after refresh, got %+v", reparsed.Resources)
	}
	for _, res := range reparsed.Resources {
		if res.EncryptionRoot != root {
			t.Errorf("resource %s lost its encryption root: %+v", res.Path, res)
		}
		if len(res.Binaries) != 1 || res.Binaries[0].Path != appBin {
			t.Errorf("resource %s: shared whitelist not applied: %+v", res.Path, res.Binaries)
		}
	}
}

// TestGroupedSectionAddressedByEncryptionRoot documents why askEncryption /
// patchCatalogSection must address a grouped section by its encryption root:
// the config text has ONE [watch <root>] header, so a text patch keyed on a
// watch sub-path fails "section not found" — the exact FATAL the installer
// hit on a grouped Discord config.
func TestGroupedSectionAddressedByEncryptionRoot(t *testing.T) {
	home := t.TempDir()
	confText, cfg := groupedDiscordConf(t, home)
	r := cfg.EncryptionGroups()[0]
	if r.EncryptionRootOrPath() == r.Path {
		t.Fatalf("test needs a grouped resource: root == watch path %q", r.Path)
	}

	if _, err := inst.SetNeedEncryption(confText, r.EncryptionRootOrPath(), false); err != nil {
		t.Errorf("SetNeedEncryption keyed on the encryption root must succeed: %v", err)
	}
	if _, err := inst.SetNeedEncryption(confText, r.Path, false); err == nil {
		t.Error("SetNeedEncryption keyed on a watch sub-path must fail (no section header): regression guard for the installer FATAL")
	}
}

// TestSetSectionWhitelistPreservesGroupStructure is the regression test for
// grouped-section refreshes: SetSectionWhitelist must replace ONLY the
// whitelist-entry lines, preserving the `watch:` group directives — a
// refresh that erased them would silently unwatch the grouped trees. The
// refreshed config must also parse back into the grouped resources with
// their shared whitelist and encryption root intact.
func TestSetSectionWhitelistPreservesGroupStructure(t *testing.T) {
	vault := t.TempDir()
	lsDir := filepath.Join(vault, "Local Storage")
	if err := os.MkdirAll(lsDir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(vault, "Cookies"), []byte("sqlite"), 0600); err != nil {
		t.Fatal(err)
	}
	appBin := filepath.Join(vault, "app-1.0.155", "Discord")
	if err := os.MkdirAll(filepath.Dir(appBin), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(appBin, []byte("ELF"), 0755); err != nil {
		t.Fatal(err)
	}

	conf := `[watch ` + vault + `]

watch: ` + lsDir + `
watch: ` + filepath.Join(vault, "Cookies") + `

` + appBin + `

need_encryption: true
`
	fresh := []inst.BinaryRule{{Path: appBin}}

	updated, err := inst.SetSectionWhitelist(conf, vault, fresh)
	if err != nil {
		t.Fatalf("SetSectionWhitelist: %v", err)
	}
	if !strings.Contains(updated, "watch: "+lsDir) ||
		!strings.Contains(updated, "watch: "+filepath.Join(vault, "Cookies")) {
		t.Errorf("group watch directives were erased:\n%s", updated)
	}
	if strings.Contains(updated, "/usr/bin/old-binary") {
		t.Errorf("stale whitelist entry survived:\n%s", updated)
	}

	// The refreshed config must still parse into the grouped resources with
	// their shared whitelist and encryption root intact.
	confPath := filepath.Join(t.TempDir(), "daemon.conf")
	if err := os.WriteFile(confPath, []byte(updated), 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := daemonconfig.Load(confPath)
	if err != nil {
		t.Fatalf("refreshed config does not parse: %v", err)
	}
	if len(cfg.Resources) != 2 {
		t.Errorf("want 2 grouped resources after refresh, got %+v", cfg.Resources)
	}
	for _, r := range cfg.Resources {
		if r.EncryptionRoot != vault {
			t.Errorf("resource %s lost its encryption root: %+v", r.Path, r)
		}
		if len(r.Binaries) != 1 || r.Binaries[0].Path != appBin {
			t.Errorf("resource %s: whitelist not applied: %+v", r.Path, r.Binaries)
		}
	}
}

// TestSoftenAutomatedRefreshErr: the pacman/apt hooks and the boot-time
// refresh unit pass --yes and must not report failure just because the
// installed config has no catalog-managed section. Interactive callers do.
func TestSoftenAutomatedRefreshErr(t *testing.T) {
	if err := softenAutomatedRefreshErr(errNoCatalogMatch, true); err != nil {
		t.Errorf("automated caller: errNoCatalogMatch must be softened to nil, got %v", err)
	}
	if err := softenAutomatedRefreshErr(errNoCatalogMatch, false); err == nil {
		t.Error("interactive caller: errNoCatalogMatch must surface")
	}
	other := errors.New("vault locked")
	if err := softenAutomatedRefreshErr(other, true); err != other {
		t.Errorf("unrelated errors must pass through even for automated callers, got %v", err)
	}
	if err := softenAutomatedRefreshErr(nil, true); err != nil {
		t.Errorf("nil in, nil out, got %v", err)
	}
}
