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
