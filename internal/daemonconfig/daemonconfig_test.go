package daemonconfig

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// writeConfig writes content into a fresh temp file and returns its path.
func writeConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "daemon.conf")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func TestLoadBasic(t *testing.T) {
	dir := t.TempDir()
	sshPath := filepath.Join(t.TempDir(), "ssh")
	if err := os.WriteFile(sshPath, []byte("#!/bin/sh"), 0755); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(writeConfig(t, `# comment
[watch `+dir+`]

/usr/bin/ssh
`+sshPath+` READ
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 1 {
		t.Fatalf("want 1 resource, got %d", len(cfg.Resources))
	}
	r := cfg.Resources[0]
	if r.Path != dir {
		t.Errorf("path = %q, want %q", r.Path, dir)
	}
	if !r.NeedEncryption {
		t.Error("need_encryption should default to true")
	}
	if len(r.Binaries) != 2 {
		t.Fatalf("want 2 binaries, got %d", len(r.Binaries))
	}
	if r.Binaries[0].Path != "/usr/bin/ssh" || len(r.Binaries[0].Events) != 0 {
		t.Errorf("bare binary must have no event restriction: %+v", r.Binaries[0])
	}
	if len(r.Binaries[1].Events) != 1 || r.Binaries[1].Events[0] != ebpf.EventRead {
		t.Errorf("expected READ restriction, got %+v", r.Binaries[1].Events)
	}
}

func TestLoadNeedEncryptionFalse(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(writeConfig(t, `[watch `+dir+`]
need_encryption: false
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Resources[0].NeedEncryption {
		t.Error("need_encryption should be false")
	}
}

func TestLoadNeedEncryptionExplicitTrue(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(writeConfig(t, `[watch `+dir+`]
need_encryption: true
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !cfg.Resources[0].NeedEncryption {
		t.Error("need_encryption should be true")
	}
}

func TestLoadMultipleResources(t *testing.T) {
	d1, d2 := t.TempDir(), t.TempDir()
	cfg, err := Load(writeConfig(t, `[watch `+d1+`]
[watch `+d2+`]
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 2 {
		t.Fatalf("want 2 resources, got %d", len(cfg.Resources))
	}
}

func TestLoadMissingWatchPathSkipped(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(writeConfig(t, `[watch /nonexistent/dir]
[watch `+dir+`]
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 1 || cfg.Resources[0].Path != dir {
		t.Fatalf("missing path should be skipped: %+v", cfg.Resources)
	}
}

func TestLoadMissingBinaryDeferred(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(writeConfig(t, `[watch `+dir+`]
/nonexistent/binary READ
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	r := cfg.Resources[0]
	if len(r.Binaries) != 0 {
		t.Errorf("unreadable binary must not be whitelisted: %+v", r.Binaries)
	}
	if len(r.PendingBinaries) != 1 {
		t.Fatalf("unreadable binary should be deferred, got %+v", r.PendingBinaries)
	}
	if r.PendingBinaries[0].Path != "/nonexistent/binary" {
		t.Errorf("deferred path = %q, want /nonexistent/binary", r.PendingBinaries[0].Path)
	}
	if len(r.PendingBinaries[0].Events) != 1 || r.PendingBinaries[0].Events[0] != ebpf.EventRead {
		t.Errorf("deferred rule must keep its event list: %+v", r.PendingBinaries[0].Events)
	}
}

func TestLoadBinarySymlinkResolved(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(t.TempDir(), "real-bin")
	if err := os.WriteFile(target, []byte("#!/bin/sh"), 0755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(t.TempDir(), "link-bin")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(writeConfig(t, `[watch `+dir+`]
`+link+`
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	r := cfg.Resources[0]
	if len(r.Binaries) != 1 {
		t.Fatalf("want 1 binary, got %+v", r.Binaries)
	}
	if r.Binaries[0].Path != target {
		t.Errorf("symlink should resolve to real path %q, got %q", target, r.Binaries[0].Path)
	}
	if len(r.PendingBinaries) != 0 {
		t.Errorf("readable symlink must not be deferred: %+v", r.PendingBinaries)
	}
}

func TestLoadUnknownEventFails(t *testing.T) {
	dir := t.TempDir()
	_, err := Load(writeConfig(t, `[watch `+dir+`]
/usr/bin/ssh BOGUS
`))
	if err == nil {
		t.Fatal("expected error for unknown event type")
	}
}

func TestLoadMultipleEvents(t *testing.T) {
	dir := t.TempDir()
	sshPath := filepath.Join(t.TempDir(), "ssh")
	if err := os.WriteFile(sshPath, []byte("#!/bin/sh"), 0755); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(writeConfig(t, `[watch `+dir+`]
`+sshPath+` READ, WRITE,READ
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	events := cfg.Resources[0].Binaries[0].Events
	if len(events) != 2 || events[0] != ebpf.EventRead || events[1] != ebpf.EventWrite {
		t.Fatalf("expected deduplicated [READ WRITE], got %+v", events)
	}
}

func TestLoadInvalidNeedEncryptionFails(t *testing.T) {
	dir := t.TempDir()
	_, err := Load(writeConfig(t, `[watch `+dir+`]
need_encryption: maybe
`))
	if err == nil {
		t.Fatal("expected error for invalid need_encryption value")
	}
}

func TestLoadDirectiveOutsideSectionSkipped(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(writeConfig(t, `/usr/bin/ssh
[watch `+dir+`]
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 1 || cfg.Resources[0].Path != dir {
		t.Fatalf("directive outside a section should be skipped: %+v", cfg.Resources)
	}
	if len(cfg.Resources[0].Binaries) != 0 {
		t.Errorf("stray directive must not be attached to the first section: %+v", cfg.Resources[0].Binaries)
	}
}

// TestLoadSkippedSectionDirectivesIgnored is the regression test for the
// missing-directory bug: the directives of a section whose watch path does
// not exist must be ignored along with the section, never treated as
// directives outside any section.
func TestLoadSkippedSectionDirectivesIgnored(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(writeConfig(t, `[watch /nonexistent/dir]
need_encryption: true
/usr/bin/ssh READ,WRITE
[watch `+dir+`]
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 1 || cfg.Resources[0].Path != dir {
		t.Fatalf("want only the valid section, got %+v", cfg.Resources)
	}
	if len(cfg.Resources[0].Binaries) != 0 {
		t.Errorf("skipped section's binaries must not leak into the next section: %+v", cfg.Resources[0].Binaries)
	}
}

func TestLoadSkippedSectionThenValidSection(t *testing.T) {
	dir := t.TempDir()
	sshPath := filepath.Join(t.TempDir(), "ssh")
	if err := os.WriteFile(sshPath, []byte("#!/bin/sh"), 0755); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(writeConfig(t, `[watch /nonexistent/dir]
need_encryption: false
[watch `+dir+`]
`+sshPath+`
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	r := cfg.Resources[0]
	if r.Path != dir {
		t.Fatalf("path = %q, want %q", r.Path, dir)
	}
	if !r.NeedEncryption {
		t.Error("the skipped section's need_encryption: false must not apply to the next section (default is true)")
	}
	if len(r.Binaries) != 1 {
		t.Fatalf("want 1 binary, got %d", len(r.Binaries))
	}
}

func TestLoadMultipleSkippedSections(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(writeConfig(t, `[watch /nonexistent/a]
[watch /nonexistent/b]
[watch `+dir+`]
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 1 || cfg.Resources[0].Path != dir {
		t.Fatalf("want only the valid section, got %+v", cfg.Resources)
	}
}

func TestLoadUnknownEventInSkippedSectionIgnored(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(writeConfig(t, `[watch /nonexistent/dir]
/usr/bin/ssh BOGUS
[watch `+dir+`]
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 1 {
		t.Fatalf("want 1 resource, got %d", len(cfg.Resources))
	}
}

func TestLoadTrailingDirectivesSkipped(t *testing.T) {
	dir := t.TempDir()
	cfg, err := Load(writeConfig(t, `[watch `+dir+`]
need_encryption: false
/usr/bin/ssh
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 1 {
		t.Fatalf("want 1 resource, got %d", len(cfg.Resources))
	}
}

// TestLoadMalformedSectionHeaderFailsEarly is the fail-closed replacement
// of the old "treated as binary" behavior: `[watch /missing/bracket` (no
// closing bracket) used to be swallowed as a stray directive of the NEXT
// section's predecessor, silently merging sections. It is now a hard error.
func TestLoadMalformedSectionHeaderFailsEarly(t *testing.T) {
	dir := t.TempDir()
	_, err := Load(writeConfig(t, `[watch /missing/bracket
[watch `+dir+`]
`))
	if err == nil {
		t.Fatal("expected error for the malformed header")
	}
	if !strings.Contains(err.Error(), "malformed section header") {
		t.Errorf("error should mention the malformed header, got: %v", err)
	}
}

func TestLoadSectionOrderPreserved(t *testing.T) {
	d1, d2, d3 := t.TempDir(), t.TempDir(), t.TempDir()
	cfg, err := Load(writeConfig(t, `[watch `+d2+`]
[watch `+d1+`]
[watch `+d3+`]
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	var got []string
	for _, r := range cfg.Resources {
		got = append(got, r.Path)
	}
	want := []string{d2, d1, d3}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("section order not preserved: got %v, want %v", got, want)
		}
	}
}

func TestLoadDuplicateWatchFails(t *testing.T) {
	dir := t.TempDir()
	_, err := Load(writeConfig(t, `[watch `+dir+`]
[watch `+dir+`]
`))
	if err == nil {
		t.Fatal("expected error for duplicate watch path")
	}
}

func TestLoadNonexistentFile(t *testing.T) {
	if _, err := Load("/nonexistent/daemon.conf"); err == nil {
		t.Fatal("expected error for missing config file")
	}
}

// TestLoadAcceptsSingleFileWatch verifies the single-file watch support:
// a regular file is accepted as a [watch] target exactly like a directory.
func TestLoadAcceptsSingleFileWatch(t *testing.T) {
	file := filepath.Join(t.TempDir(), "secret.env")
	if err := os.WriteFile(file, []byte("KEY=value"), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(writeConfig(t, "[watch "+file+"]\n/usr/bin/cat\n"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 1 || cfg.Resources[0].Path != file {
		t.Fatalf("resources = %+v, want exactly %q", cfg.Resources, file)
	}
	if !cfg.Resources[0].NeedEncryption {
		t.Error("need_encryption should default to true for a file resource too")
	}
}

// TestLoadRefusesSymlinkWatch: guard identity is inode based, so a
// symbolic link as watch target must be skipped, not followed.
func TestLoadRefusesSymlinkWatch(t *testing.T) {
	dir := t.TempDir()
	link := filepath.Join(dir, "link")
	if err := os.Symlink(dir, link); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(writeConfig(t, "[watch "+link+"]\n/usr/bin/ssh\n"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 0 {
		t.Fatalf("symlink watch target was accepted: %+v", cfg.Resources)
	}
}

// TestLoadRefusesHardlinkedFileWatch: more than one hard link means another
// path aliases the same inode; guarding one would implicitly cover both.
func TestLoadRefusesHardlinkedFileWatch(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a")
	if err := os.WriteFile(a, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	b := filepath.Join(dir, "b")
	if err := os.Link(a, b); err != nil {
		t.Skipf("hardlinks unavailable here: %v", err)
	}

	cfg, err := Load(writeConfig(t, "[watch "+a+"]\n/usr/bin/cat\n"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 0 {
		t.Fatalf("hardlinked watch target was accepted: %+v", cfg.Resources)
	}
}

// TestLoadRefusesSpecialFileWatch covers FIFOs (the same refusal applies to
// sockets and devices): only directories and regular files are watchable.
func TestLoadRefusesSpecialFileWatch(t *testing.T) {
	dir := t.TempDir()
	fifo := filepath.Join(dir, "pipe")
	if err := syscall.Mkfifo(fifo, 0o600); err != nil {
		t.Skipf("mkfifo unavailable here: %v", err)
	}

	cfg, err := Load(writeConfig(t, "[watch "+fifo+"]\n/usr/bin/cat\n"))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 0 {
		t.Fatalf("FIFO watch target was accepted: %+v", cfg.Resources)
	}
}

// TestLoadMalformedSectionHeaderFails is the regression test for the
// cross-resource whitelist contamination: a bracketed line that is not a
// valid [watch <path>] header used to be parsed as a directive of the
// PREVIOUS section, silently moving every following binary into the wrong
// whitelist. Malformed headers are now hard parse errors.
func TestLoadMalformedSectionHeaderFails(t *testing.T) {
	dir := t.TempDir()
	sshPath := filepath.Join(t.TempDir(), "ssh")
	if err := os.WriteFile(sshPath, []byte("#!/bin/sh"), 0755); err != nil {
		t.Fatal(err)
	}
	_, err := Load(writeConfig(t, `[watch `+dir+`]
`+sshPath+`
[watch /tmp/other # comment breaks the header
`))
	if err == nil {
		t.Fatal("expected error for a malformed section header")
	}
	if !strings.Contains(err.Error(), "malformed section header") {
		t.Errorf("error should mention the malformed header, got: %v", err)
	}
}

// TestLoadWatchGroup covers the `watch:` group syntax: the section path is
// the encryption root only (it is NOT itself watched — the app's updater
// writes there), each `watch:` path becomes its own guarded resource sharing
// the group whitelist, and need_encryption applies to the whole group.
func TestLoadWatchGroup(t *testing.T) {
	dir := t.TempDir()
	sshPath := filepath.Join(t.TempDir(), "ssh")
	if err := os.WriteFile(sshPath, []byte("#!/bin/sh"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dir+"/Local Storage", 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dir+"/Cookies", []byte("sqlite"), 0600); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(writeConfig(t, `[watch `+dir+`]
watch: `+dir+`/Local Storage
watch: `+dir+`/Cookies
`+sshPath+`
need_encryption: true
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 2 {
		t.Fatalf("want 2 grouped resources, got %+v", cfg.Resources)
	}
	for _, r := range cfg.Resources {
		if r.EncryptionRoot != dir {
			t.Errorf("resource %s: EncryptionRoot = %q, want %q", r.Path, r.EncryptionRoot, dir)
		}
		if !r.NeedEncryption {
			t.Errorf("resource %s: NeedEncryption = false, want true", r.Path)
		}
		if len(r.Binaries) != 1 || r.Binaries[0].Path != sshPath {
			t.Errorf("resource %s: shared whitelist not applied: %+v", r.Path, r.Binaries)
		}
	}
	if cfg.Resources[0].Path != dir+"/Local Storage" || cfg.Resources[1].Path != dir+"/Cookies" {
		t.Errorf("unexpected watch paths: %q, %q", cfg.Resources[0].Path, cfg.Resources[1].Path)
	}
}

// TestLoadWatchGroupValidation covers the group parse errors: watch paths
// outside the section directory, duplicates, and the section path itself.
func TestLoadWatchGroupValidation(t *testing.T) {
	dir := t.TempDir()
	for _, tc := range []struct {
		name    string
		watchLn string
	}{
		{name: "outside root", watchLn: "watch: /elsewhere/data"},
		{name: "section path itself", watchLn: "watch: " + dir},
		{name: "traversal escape", watchLn: "watch: " + dir + "/../escape"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Load(writeConfig(t, "[watch "+dir+"]\n"+tc.watchLn+"\n"))
			if err == nil {
				t.Fatalf("expected error for %q", tc.watchLn)
			}
		})
	}
}

// TestLoadUngroupedUnchanged is the backward-compatibility anchor: a section
// without watch: directives keeps the historical single-resource behavior
// (section path watched, empty EncryptionRoot).
func TestLoadUngroupedUnchanged(t *testing.T) {
	dir := t.TempDir()
	sshPath := filepath.Join(t.TempDir(), "ssh")
	if err := os.WriteFile(sshPath, []byte("#!/bin/sh"), 0755); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(writeConfig(t, `[watch `+dir+`]
`+sshPath+`
need_encryption: false
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 1 {
		t.Fatalf("want 1 resource, got %d", len(cfg.Resources))
	}
	r := cfg.Resources[0]
	if r.Path != dir || r.EncryptionRoot != "" || r.NeedEncryption {
		t.Errorf("ungrouped resource changed: %+v", r)
	}
}

// TestLoadMissingSectionRootSkipsDirectives preserves the historical
// tolerance: a section whose root is missing is skipped wholesale — its
// directives (even malformed ones) are warned and ignored, never fatal.
func TestLoadMissingSectionRootSkipsDirectives(t *testing.T) {
	cfg, err := Load(writeConfig(t, `[watch /nonexistent/root]
/usr/bin/ssh BOGUS
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 0 {
		t.Errorf("want 0 resources, got %+v", cfg.Resources)
	}
}

// TestLoadGroupedEncryptedWatchPathsDeferred: a grouped need_encryption: true
// section whose encryption root exists (an unlocked-name locked fscrypt tree)
// but whose watch sub-paths do not resolve yet must NOT be dropped — the
// sub-path names only appear once the daemon unlocks the vault. Each is kept
// as a PathPending resource carrying the shared whitelist and encryption
// root, so the daemon can unlock the root and re-validate.
func TestLoadGroupedEncryptedWatchPathsDeferred(t *testing.T) {
	root := t.TempDir() // exists; the sub-paths deliberately do not
	sshPath := filepath.Join(t.TempDir(), "ssh")
	if err := os.WriteFile(sshPath, []byte("#!/bin/sh"), 0755); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(writeConfig(t, `[watch `+root+`]
watch: `+root+`/Local Storage
watch: `+root+`/Cookies
`+sshPath+`
need_encryption: true
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 2 {
		t.Fatalf("want 2 deferred grouped resources, got %+v", cfg.Resources)
	}
	for _, r := range cfg.Resources {
		if !r.PathPending {
			t.Errorf("resource %s: PathPending = false, want true", r.Path)
		}
		if r.EncryptionRoot != root {
			t.Errorf("resource %s: EncryptionRoot = %q, want %q", r.Path, r.EncryptionRoot, root)
		}
		if !r.NeedEncryption {
			t.Errorf("resource %s: NeedEncryption = false", r.Path)
		}
		if len(r.Binaries) != 1 || r.Binaries[0].Path != sshPath {
			t.Errorf("resource %s: shared whitelist not applied: %+v", r.Path, r.Binaries)
		}
	}

	// The daemon unlocks the vault -> the sub-paths appear -> ResolvePendingPaths
	// clears the flag.
	if err := os.MkdirAll(root+"/Local Storage", 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(root+"/Cookies", []byte("sqlite"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := ResolvePendingPaths(cfg); err != nil {
		t.Fatalf("ResolvePendingPaths after unlock: %v", err)
	}
	for _, r := range cfg.Resources {
		if r.PathPending {
			t.Errorf("resource %s still PathPending after the sub-path appeared", r.Path)
		}
	}
}

// TestResolvePendingPathsRejectsMissingAndSymlink: a pending sub-path that is
// still missing after the unlock, or that an attacker planted a symlink at
// during the unlock window, is a hard error — never silently dropped (that
// would leave a declared-protected tree unguarded, or let a symlink redirect
// the inode-based guard).
func TestResolvePendingPathsRejectsMissingAndSymlink(t *testing.T) {
	sshPath := filepath.Join(t.TempDir(), "ssh")
	if err := os.WriteFile(sshPath, []byte("#!/bin/sh"), 0755); err != nil {
		t.Fatal(err)
	}

	loadPending := func(root string) *Config {
		cfg, err := Load(writeConfig(t, `[watch `+root+`]
watch: `+root+`/data
`+sshPath+`
need_encryption: true
`))
		if err != nil {
			t.Fatalf("Load: %v", err)
		}
		if len(cfg.Resources) != 1 || !cfg.Resources[0].PathPending {
			t.Fatalf("want 1 pending resource, got %+v", cfg.Resources)
		}
		return cfg
	}

	// Still missing after the unlock.
	missingRoot := t.TempDir()
	if err := ResolvePendingPaths(loadPending(missingRoot)); err == nil {
		t.Error("a still-missing pending path must be a hard error")
	}

	// Symlink planted during the unlock window (deferred while missing, then
	// a symlink appears before ResolvePendingPaths runs).
	symRoot := t.TempDir()
	cfg := loadPending(symRoot)
	if err := os.Symlink(t.TempDir(), symRoot+"/data"); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	if err := ResolvePendingPaths(cfg); err == nil {
		t.Error("a pending path resolving to a symlink must be a hard error")
	}
}

// TestLoadGroupedUnencryptedMissingPathStillDropped: the defer only applies
// to need_encryption: true groups (the locked-vault case). A missing sub-path
// in a need_encryption: false group has no vault to unlock and stays dropped.
func TestLoadGroupedUnencryptedMissingPathStillDropped(t *testing.T) {
	root := t.TempDir()
	present := filepath.Join(root, "present")
	if err := os.MkdirAll(present, 0700); err != nil {
		t.Fatal(err)
	}
	sshPath := filepath.Join(t.TempDir(), "ssh")
	if err := os.WriteFile(sshPath, []byte("#!/bin/sh"), 0755); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(writeConfig(t, `[watch `+root+`]
watch: `+present+`
watch: `+filepath.Join(root, "gone")+`
`+sshPath+`
need_encryption: false
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 1 || cfg.Resources[0].Path != present {
		t.Fatalf("want only the present sub-path, got %+v", cfg.Resources)
	}
	if cfg.Resources[0].PathPending {
		t.Error("an unencrypted group's resource must never be PathPending")
	}
}

// TestLoadGroupedWatchPathRejectsSymlinkedIntermediate: a grouped watch
// sub-path whose intermediate component is a symlink is dropped at parse time
// — validateWatchTarget only Lstats the leaf, so without the component walk an
// intermediate symlink pointing outside the vault would redirect the
// inode-based guard onto an unrelated tree that would inherit the group's
// whitelist and fscrypt lifecycle.
func TestLoadGroupedWatchPathRejectsSymlinkedIntermediate(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.MkdirAll(filepath.Join(outside, "child"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "evil")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	sshPath := filepath.Join(t.TempDir(), "ssh")
	if err := os.WriteFile(sshPath, []byte("#!/bin/sh"), 0755); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(writeConfig(t, `[watch `+root+`]
watch: `+filepath.Join(root, "evil", "child")+`
`+sshPath+`
need_encryption: true
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 0 {
		t.Fatalf("a watch path via a symlinked intermediate must be dropped, got %+v", cfg.Resources)
	}
}

// TestResolvePendingPathsRejectsSymlinkedIntermediate: the same protection at
// unlock time — an intermediate component that appears as a symlink after the
// vault is unlocked is a hard error, not a silently redirected guard.
func TestResolvePendingPathsRejectsSymlinkedIntermediate(t *testing.T) {
	root := t.TempDir()
	sshPath := filepath.Join(t.TempDir(), "ssh")
	if err := os.WriteFile(sshPath, []byte("#!/bin/sh"), 0755); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(writeConfig(t, `[watch `+root+`]
watch: `+filepath.Join(root, "mid", "data")+`
`+sshPath+`
need_encryption: true
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 1 || !cfg.Resources[0].PathPending {
		t.Fatalf("want 1 pending resource, got %+v", cfg.Resources)
	}
	outside := t.TempDir()
	if err := os.MkdirAll(filepath.Join(outside, "data"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "mid")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	if err := ResolvePendingPaths(cfg); err == nil {
		t.Error("a pending path reached through a symlinked intermediate must be a hard error")
	}
}
