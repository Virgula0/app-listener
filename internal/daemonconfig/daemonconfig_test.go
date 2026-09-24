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

// Regression for the installer bug: a binary path with spaces (Steam/Proton's ".../Proton -
// Experimental/files/bin/wineserver") had its event list parsed from the wrong field ("-
// Experimental/..." looked like an event type) because the rule line was split on whitespace. A
// double-quoted path must be taken whole.
func TestLoadQuotedPathWithSpaces(t *testing.T) {
	dir := t.TempDir()
	binDir := filepath.Join(t.TempDir(), "Proton - Experimental", "files", "bin")
	if err := os.MkdirAll(binDir, 0o700); err != nil {
		t.Fatal(err)
	}
	binPath := filepath.Join(binDir, "wineserver")
	if err := os.WriteFile(binPath, []byte("#!/bin/sh"), 0755); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(writeConfig(t, `[watch `+dir+`]
"`+binPath+`" READ,WRITE
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	r := cfg.Resources[0]
	if len(r.Binaries) != 1 {
		t.Fatalf("want 1 binary, got %+v", r.Binaries)
	}
	if r.Binaries[0].Path != binPath {
		t.Errorf("binary path = %q, want %q", r.Binaries[0].Path, binPath)
	}
	if len(r.Binaries[0].Events) != 2 || r.Binaries[0].Events[0] != ebpf.EventRead || r.Binaries[0].Events[1] != ebpf.EventWrite {
		t.Errorf("expected READ,WRITE restriction, got %+v", r.Binaries[0].Events)
	}
}

// TestLoadQuotedBareBinary covers a quoted path with no trailing event list.
func TestLoadQuotedBareBinary(t *testing.T) {
	dir := t.TempDir()
	binPath := filepath.Join(t.TempDir(), "with space", "bin")
	if err := os.MkdirAll(filepath.Dir(binPath), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(binPath, []byte("#!/bin/sh"), 0755); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(writeConfig(t, `[watch `+dir+`]
"`+binPath+`"
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources[0].Binaries) != 1 || cfg.Resources[0].Binaries[0].Path != binPath {
		t.Fatalf("binaries = %+v, want exactly %q", cfg.Resources[0].Binaries, binPath)
	}
}

// TestLoadQuotedSectionAndWatchPath covers quoting in the section header and
// in a `watch:` group directive.
func TestLoadQuotedSectionAndWatchPath(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "with space")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(dir, "also space")
	if err := os.MkdirAll(sub, 0o700); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(writeConfig(t, `[watch "`+dir+`"]
watch: "`+sub+`"
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 1 || cfg.Resources[0].Path != sub {
		t.Fatalf("resources = %+v, want exactly %q", cfg.Resources, sub)
	}
	if cfg.Resources[0].EncryptionRoot != dir {
		t.Errorf("EncryptionRoot = %q, want %q", cfg.Resources[0].EncryptionRoot, dir)
	}
}

// TestLoadUnterminatedQuoteFails ensures a malformed quoted path is a hard
// parse error, not silently truncated or misparsed.
func TestLoadUnterminatedQuoteFails(t *testing.T) {
	dir := t.TempDir()
	_, err := Load(writeConfig(t, `[watch `+dir+`]
"/no/closing/quote READ
`))
	if err == nil {
		t.Fatal("expected error for an unterminated quoted path")
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

// A grouped need_encryption: true section whose encryption root exists (locked fscrypt tree) but
// whose watch sub-paths don't resolve yet must NOT be dropped (names appear only after unlock):
// each is kept as a PathPending resource with the shared whitelist and root so the daemon can
// unlock and re-validate.
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

// A grouped watch sub-path with a symlinked intermediate component is dropped at parse time:
// validateWatchTarget Lstats only the leaf, so an intermediate symlink pointing outside the vault
// would redirect the inode-based guard onto an unrelated tree inheriting the group's whitelist and
// lifecycle.
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

func TestLoadAllowLib(t *testing.T) {
	dir := t.TempDir()
	libDir := t.TempDir()
	libA := filepath.Join(libDir, "liba.so")
	libB := filepath.Join(libDir, "lib b.so") // spaces → must be quoted
	for _, p := range []string{libA, libB} {
		if err := os.WriteFile(p, []byte("\x7fELF"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	cfg, err := Load(writeConfig(t, `[watch `+dir+`]
need_encryption: false
/usr/bin/example
allow_lib `+libA+`
allow_lib: "`+libB+`"
allow_lib `+filepath.Join(libDir, "missing.so")+`
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.Resources) != 1 {
		t.Fatalf("want 1 resource, got %d", len(cfg.Resources))
	}
	r := cfg.Resources[0]
	if len(r.AllowLibs) != 2 {
		t.Errorf("AllowLibs = %v, want the two readable libs", r.AllowLibs)
	}
	if len(r.PendingLibs) != 1 || r.PendingLibs[0] != filepath.Join(libDir, "missing.so") {
		t.Errorf("PendingLibs = %v, want the one unreadable lib deferred", r.PendingLibs)
	}
	// The binary line must not be swallowed by allow_lib handling.
	found := false
	for _, b := range append(r.Binaries, r.PendingBinaries...) {
		if b.Path == "/usr/bin/example" {
			found = true
		}
	}
	if !found {
		t.Errorf("binary /usr/bin/example missing; Binaries=%v Pending=%v", r.Binaries, r.PendingBinaries)
	}
}

func TestParseAllowLibNotADirective(t *testing.T) {
	// A binary path that merely starts with "allow_lib" must NOT be parsed as
	// the directive.
	if _, ok := parseAllowLib("/usr/bin/allow_libtool"); ok {
		t.Errorf("parseAllowLib matched a binary path with an allow_lib prefix")
	}
	if p, ok := parseAllowLib("allow_lib /usr/lib/x.so"); !ok || p != "/usr/lib/x.so" {
		t.Errorf("parseAllowLib(space) = %q,%v", p, ok)
	}
	if p, ok := parseAllowLib(`allow_lib: "/usr/lib/y.so"`); !ok || p != "/usr/lib/y.so" {
		t.Errorf("parseAllowLib(colon+quotes) = %q,%v", p, ok)
	}
}

func TestLoadLibDir(t *testing.T) {
	dir := t.TempDir()
	libDir := t.TempDir()
	missing := filepath.Join(t.TempDir(), "gone")

	cfg, err := Load(writeConfig(t, `[watch `+dir+`]
need_encryption: false
/usr/bin/example
lib_dir `+libDir+`
lib_dir: "`+missing+`"
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	// Two resources: the guarded section, plus the library directory as its
	// own read-only, unencrypted resource sharing the whitelist.
	if len(cfg.Resources) != 2 {
		t.Fatalf("want 2 resources (section + lib_dir), got %d: %+v", len(cfg.Resources), cfg.Resources)
	}
	var lib *Resource
	for i := range cfg.Resources {
		if cfg.Resources[i].Path == libDir {
			lib = &cfg.Resources[i]
		}
	}
	if lib == nil {
		t.Fatalf("lib_dir %s did not materialize; got %+v", libDir, cfg.Resources)
	}
	if !lib.ReadOnly {
		t.Errorf("lib_dir resource must be ReadOnly (world-readable, whitelist-gated writes)")
	}
	if lib.NeedEncryption {
		t.Errorf("lib_dir resource must never be encrypted")
	}
	if lib.EncryptionRoot != "" {
		t.Errorf("lib_dir resource must not join the section's encryption group, got %q", lib.EncryptionRoot)
	}
	if len(lib.Binaries)+len(lib.PendingBinaries) == 0 {
		t.Errorf("lib_dir resource must inherit the section whitelist")
	}
	// A missing lib_dir is dropped, never deferred into a phantom resource.
	for _, r := range cfg.Resources {
		if r.Path == missing {
			t.Errorf("missing lib_dir %s must be ignored, not materialized", missing)
		}
	}
}

func TestLoadLibDirSharedBySections(t *testing.T) {
	// One runtime tree named by two sections of the same app (Steam's two
	// config locations) must be guarded once, not rejected as a duplicate.
	dirA, dirB, libDir := t.TempDir(), t.TempDir(), t.TempDir()
	cfg, err := Load(writeConfig(t, `[watch `+dirA+`]
need_encryption: false
/usr/bin/example
lib_dir `+libDir+`

[watch `+dirB+`]
need_encryption: false
/usr/bin/example
lib_dir `+libDir+`
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	n := 0
	for _, r := range cfg.Resources {
		if r.Path == libDir {
			n++
		}
	}
	if n != 1 {
		t.Errorf("lib_dir named twice materialized %d resources, want exactly 1", n)
	}
}

func TestLibDirDoesNotWidenAGuardedResource(t *testing.T) {
	// A lib_dir pointing at a path that is already a real guarded resource
	// must be ignored, never folded in: folding would add this section's
	// binaries to a protected tree's whitelist.
	secret, other := t.TempDir(), t.TempDir()
	cfg, err := Load(writeConfig(t, `[watch `+secret+`]
need_encryption: false
/usr/bin/trusted

[watch `+other+`]
need_encryption: false
/usr/bin/attacker
lib_dir `+secret+`
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	for _, r := range cfg.Resources {
		if r.Path != secret {
			continue
		}
		if r.ReadOnly {
			t.Errorf("guarded resource %s was downgraded to read-only by a lib_dir", secret)
		}
		for _, b := range append(r.Binaries, r.PendingBinaries...) {
			if b.Path == "/usr/bin/attacker" {
				t.Errorf("lib_dir widened the whitelist of guarded resource %s", secret)
			}
		}
	}
}

func TestParseLibDirNotADirective(t *testing.T) {
	if _, ok := parseLibDir("/usr/bin/lib_dirtool"); ok {
		t.Errorf("parseLibDir matched a binary path with a lib_dir prefix")
	}
	if p, ok := parseLibDir("lib_dir /opt/app/lib"); !ok || p != "/opt/app/lib" {
		t.Errorf("parseLibDir(space) = %q,%v", p, ok)
	}
	if p, ok := parseLibDir(`lib_dir: "/opt/my app/lib"`); !ok || p != "/opt/my app/lib" {
		t.Errorf("parseLibDir(colon+quotes) = %q,%v", p, ok)
	}
}

func TestLibDirExcludesRestrictedWriters(t *testing.T) {
	// An event-restricted binary scopes access to the SECTION tree; read-only
	// mode cannot carry masks, so it must not become an unrestricted writer
	// of the library directory.
	dir, libDir := t.TempDir(), t.TempDir()
	cfg, err := Load(writeConfig(t, `[watch `+dir+`]
need_encryption: false
/usr/bin/free
/usr/bin/scoped READ
lib_dir `+libDir+`
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	for _, r := range cfg.Resources {
		if r.Path != libDir {
			continue
		}
		var got []string
		for _, b := range append(r.Binaries, r.PendingBinaries...) {
			if len(b.Events) != 0 {
				t.Errorf("lib_dir inherited restricted binary %s %v", b.Path, b.Events)
			}
			got = append(got, b.Path)
		}
		if len(got) != 1 || got[0] != "/usr/bin/free" {
			t.Errorf("lib_dir writers = %v, want only the unrestricted /usr/bin/free", got)
		}
	}
}

func TestEncryptionGroupsSkipsLibDirs(t *testing.T) {
	dir, libDir := t.TempDir(), t.TempDir()
	cfg, err := Load(writeConfig(t, `[watch `+dir+`]
need_encryption: false
/usr/bin/example
lib_dir `+libDir+`
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	for _, r := range cfg.EncryptionGroups() {
		if r.Path == libDir {
			t.Errorf("lib_dir %s must not be an encryption group (installer would ask/patch a non-section)", libDir)
		}
	}
	if n := len(cfg.EncryptionGroups()); n != 1 {
		t.Errorf("EncryptionGroups = %d, want 1 (the section only)", n)
	}
}

func TestParseLibBinaryNotADirective(t *testing.T) {
	// A binary path that merely starts with "lib_binary" must NOT be parsed
	// as the directive.
	if _, ok := parseLibBinary("/usr/bin/lib_binarytool"); ok {
		t.Errorf("parseLibBinary matched a binary path with a lib_binary prefix")
	}
	if v, ok := parseLibBinary("lib_binary /usr/bin/x"); !ok || v != "/usr/bin/x" {
		t.Errorf("parseLibBinary(space) = %q,%v", v, ok)
	}
	if v, ok := parseLibBinary(`lib_binary: "/usr/bin/y"`); !ok || v != `"/usr/bin/y"` {
		t.Errorf("parseLibBinary(colon+quotes) = %q,%v", v, ok)
	}
}

// TestLoadLibBinaryWritesLibDirOnly is the security property of the
// directive: a lib_binary is a writer of the section's library directories
// and NOTHING else. A runtime's own maintenance tools must be able to rebuild
// the tree they own without gaining access to the credential directory that
// happens to sit in the same section.
func TestLoadLibBinaryWritesLibDirOnly(t *testing.T) {
	secret, libDir := t.TempDir(), t.TempDir()
	writer := filepath.Join(t.TempDir(), "capsule-capture-libs")
	if err := os.WriteFile(writer, []byte("#!/bin/true\n"), 0o755); err != nil {
		t.Fatalf("writing the fake writer: %v", err)
	}

	cfg, err := Load(writeConfig(t, `[watch `+secret+`]
need_encryption: false
/usr/bin/example
lib_dir `+libDir+`
lib_binary `+writer+`
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	var protectedRes, lib *Resource
	for i := range cfg.Resources {
		switch cfg.Resources[i].Path {
		case secret:
			protectedRes = &cfg.Resources[i]
		case libDir:
			lib = &cfg.Resources[i]
		}
	}
	if protectedRes == nil || lib == nil {
		t.Fatalf("want both the section and the lib_dir resource, got %+v", cfg.Resources)
	}
	if !hasBinary(lib, writer) {
		t.Errorf("lib_binary %s must be a writer of the lib_dir; got %+v", writer, lib.Binaries)
	}
	if hasBinary(protectedRes, writer) {
		t.Errorf("lib_binary %s leaked into the PROTECTED resource's whitelist — "+
			"a library-tree writer must never reach the section's secrets", writer)
	}
	// The ordinary whitelist entry keeps reaching both, as before (it is
	// deferred here — /usr/bin/example need not exist — which is still a rule
	// on the protected resource, just an unresolved one).
	if !hasPendingBinary(protectedRes, "/usr/bin/example") {
		t.Errorf("the section whitelist must still apply to the protected resource")
	}
}

// TestLoadLibBinaryDeferredWhenUnreadable keeps the fail-closed handling of a
// binary that is not resolvable yet: deferred (denied until resolved), never
// dropped.
func TestLoadLibBinaryDeferredWhenUnreadable(t *testing.T) {
	secret, libDir := t.TempDir(), t.TempDir()
	missing := filepath.Join(t.TempDir(), "not-installed-yet")

	cfg, err := Load(writeConfig(t, `[watch `+secret+`]
need_encryption: false
/usr/bin/example
lib_dir `+libDir+`
lib_binary `+missing+`
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	for i := range cfg.Resources {
		if cfg.Resources[i].Path != libDir {
			continue
		}
		for _, p := range cfg.Resources[i].PendingBinaries {
			if p.Path == missing {
				return
			}
		}
		t.Fatalf("an unreadable lib_binary must be deferred, got pending %+v", cfg.Resources[i].PendingBinaries)
	}
	t.Fatalf("lib_dir resource missing from %+v", cfg.Resources)
}

// TestLoadLibBinaryRejectsEventRestrictions: a lib_dir is guarded read-only,
// which carries no per-binary event mask, so an event list must be refused
// rather than silently ignored.
func TestLoadLibBinaryRejectsEventRestrictions(t *testing.T) {
	secret, libDir := t.TempDir(), t.TempDir()
	_, err := Load(writeConfig(t, `[watch `+secret+`]
need_encryption: false
lib_dir `+libDir+`
lib_binary /usr/bin/example READ,WRITE
`))
	if err == nil {
		t.Fatalf("lib_binary with an event list must be rejected")
	}
	if !strings.Contains(err.Error(), "event restrictions") {
		t.Errorf("error should name the rejected restriction, got: %v", err)
	}
}

// hasPendingBinary reports whether the resource carries path as a deferred rule.
func hasPendingBinary(r *Resource, path string) bool {
	for _, b := range r.PendingBinaries {
		if b.Path == path {
			return true
		}
	}
	return false
}

// hasBinary reports whether the resource lists path as a resolved writer.
func hasBinary(r *Resource, path string) bool {
	for _, b := range r.Binaries {
		if b.Path == path {
			return true
		}
	}
	return false
}

func TestParseLibrariesSection(t *testing.T) {
	for line, want := range map[string]string{
		`[libraries]`:                 "",
		`[libraries Steam]`:           "Steam",
		`[libraries "Steam (alice)"]`: "Steam (alice)",
		"[libraries\t\"Proton GE\"]":  "Proton GE",
	} {
		got, ok := parseLibrariesSection(line)
		if !ok || got != want {
			t.Errorf("parseLibrariesSection(%q) = %q,%v; want %q,true", line, got, ok, want)
		}
	}
	for _, line := range []string{`[librariesfoo]`, `[watch /x]`, `libraries`} {
		if _, ok := parseLibrariesSection(line); ok {
			t.Errorf("parseLibrariesSection(%q) matched; want no match", line)
		}
	}
	// A near-miss header must still be refused as malformed, never merged
	// into the previous section.
	if _, err := Load(writeConfig(t, "[librariesfoo]\n")); err == nil {
		t.Errorf("a malformed [libraries...] header must be rejected")
	}
}

// TestLibrariesBlockScopesWritersToItself is the security property of the
// block form: a lib_binary writes the lib_dirs of its OWN block only. Neither
// a watch section's whitelist nor another application's block may write them.
func TestLibrariesBlockScopesWritersToItself(t *testing.T) {
	secret := t.TempDir()
	steamRuntime, otherRuntime := t.TempDir(), t.TempDir()
	mkExe := func(name string) string {
		p := filepath.Join(t.TempDir(), name)
		if err := os.WriteFile(p, []byte("#!/bin/true\n"), 0o755); err != nil {
			t.Fatalf("writing %s: %v", name, err)
		}
		return p
	}
	sectionApp, steamWriter, otherWriter := mkExe("app"), mkExe("pressure-vessel-wrap"), mkExe("other-updater")

	cfg, err := Load(writeConfig(t, `[watch `+secret+`]
need_encryption: false
`+sectionApp+`

[libraries "Steam (alice)"]
lib_dir `+steamRuntime+`
lib_binary `+steamWriter+`

[libraries "Other (alice)"]
lib_dir `+otherRuntime+`
lib_binary `+otherWriter+`
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	byPath := map[string]*Resource{}
	for i := range cfg.Resources {
		byPath[cfg.Resources[i].Path] = &cfg.Resources[i]
	}
	steam, other := byPath[steamRuntime], byPath[otherRuntime]
	if steam == nil || other == nil {
		t.Fatalf("both library blocks must materialize their lib_dir, got %+v", cfg.Resources)
	}
	if !steam.ReadOnly || steam.NeedEncryption {
		t.Errorf("a [libraries] lib_dir must be read-only and unencrypted: %+v", steam)
	}
	if !hasBinary(steam, steamWriter) {
		t.Errorf("the block's own lib_binary must write its lib_dir")
	}
	if hasBinary(steam, otherWriter) {
		t.Errorf("another application's lib_binary leaked into this block's lib_dir")
	}
	if hasBinary(steam, sectionApp) {
		t.Errorf("a watch section's whitelist leaked into a [libraries] lib_dir — blocks inherit no whitelist")
	}
	if hasBinary(byPath[secret], steamWriter) {
		t.Errorf("a lib_binary leaked into a protected resource's whitelist")
	}
}

func TestLibrariesBlockAllowLibIsShared(t *testing.T) {
	lib := filepath.Join(t.TempDir(), "plugin.so")
	if err := os.WriteFile(lib, []byte("ELF"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	missing := filepath.Join(t.TempDir(), "in-a-locked-vault.so")
	cfg, err := Load(writeConfig(t, `[libraries]
allow_lib `+lib+`
allow_lib `+missing+`
`))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	// Readable or not, both reach the trust guard: it runs after the vaults
	// are unlocked and re-stats every path itself.
	got := strings.Join(cfg.SharedAllowLibs, ",")
	if !strings.Contains(got, lib) || !strings.Contains(got, missing) {
		t.Errorf("SharedAllowLibs = %v, want both %s and %s", cfg.SharedAllowLibs, lib, missing)
	}
	if len(cfg.Resources) != 0 {
		t.Errorf("an allow_lib-only block must not create resources, got %+v", cfg.Resources)
	}
}

func TestLibrariesBlockRejectsBinaryLines(t *testing.T) {
	_, err := Load(writeConfig(t, `[libraries "Steam"]
/usr/bin/steam
`))
	if err == nil || !strings.Contains(err.Error(), "not allowed in a [libraries] block") {
		t.Fatalf("a bare binary line in a [libraries] block must be refused, got %v", err)
	}
}

// Nested guarded roots would give a file two owners in the one-owner-per-inode engine.
func TestLoadRejectsNestedRoots(t *testing.T) {
	outer := t.TempDir()
	inner := filepath.Join(outer, "inner")
	if err := os.Mkdir(inner, 0o755); err != nil {
		t.Fatal(err)
	}
	aliasParent := filepath.Join(t.TempDir(), "alias")
	if err := os.Symlink(outer, aliasParent); err != nil {
		t.Fatal(err)
	}

	for name, conf := range map[string]string{
		"watch inside watch":   "[watch " + outer + "]\nneed_encryption: false\n\n[watch " + inner + "]\nneed_encryption: false\n",
		"watch around watch":   "[watch " + inner + "]\nneed_encryption: false\n\n[watch " + outer + "]\nneed_encryption: false\n",
		"lib_dir inside watch": "[watch " + outer + "]\nneed_encryption: false\n\n[libraries \"x\"]\nlib_dir " + inner + "\n",
		"inside via symlinked parent": "[watch " + outer + "]\nneed_encryption: false\n\n[watch " +
			filepath.Join(aliasParent, "inner") + "]\nneed_encryption: false\n",
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := Load(writeConfig(t, conf)); err == nil {
				t.Fatal("a config with nested guarded roots must be rejected")
			}
		})
	}
}

func TestLoadAcceptsSiblingRoots(t *testing.T) {
	base := t.TempDir()
	a, b := filepath.Join(base, "a"), filepath.Join(base, "ab")
	for _, d := range []string{a, b} {
		if err := os.Mkdir(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	conf := "[watch " + a + "]\nneed_encryption: false\n\n[watch " + b + "]\nneed_encryption: false\n"
	if _, err := Load(writeConfig(t, conf)); err != nil {
		t.Fatalf("sibling roots sharing a name prefix must load: %v", err)
	}
}
