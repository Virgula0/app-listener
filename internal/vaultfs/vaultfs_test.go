package vaultfs

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func inodeOf(t *testing.T, path string) uint64 {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	return info.Sys().(*syscall.Stat_t).Ino
}

func TestIsBinaryContent(t *testing.T) {
	if IsBinaryContent([]byte("plain text\n")) {
		t.Error("plain text classified as binary")
	}
	if !IsBinaryContent([]byte{'a', 0, 'b'}) {
		t.Error("NUL byte not classified as binary")
	}
	if !IsBinaryContent([]byte{0}) {
		t.Error("leading NUL not classified as binary")
	}
	// Only the first BinaryScanBytes window is scanned: a NUL beyond it
	// is not seen, so the file is treated as text.
	late := make([]byte, BinaryScanBytes+16) // all-NUL beyond the window
	copy(late, strings.Repeat("x", BinaryScanBytes))
	late[BinaryScanBytes] = 0
	if IsBinaryContent(late) {
		t.Error("NUL beyond the scan window classified the file as binary")
	}
}

func openRootT(t *testing.T, dir string) *os.Root {
	t.Helper()
	r, err := os.OpenRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = r.Close() })
	return r
}

func TestWriteFileKeepMeta(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "secret.cfg")
	if err := os.WriteFile(path, []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := WriteFileKeepMeta(openRootT(t, dir), "secret.cfg", []byte("new-content")); err != nil {
		t.Fatalf("WriteFileKeepMeta: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "new-content" {
		t.Errorf("content = %q, want new-content", data)
	}
	info, err := os.Lstat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("mode = %v, want 0600 (the original, not the temp default)", info.Mode().Perm())
	}
}

func TestWriteFileKeepMetaRefusesSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	if err := os.WriteFile(target, []byte("t"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "link")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	if err := WriteFileKeepMeta(openRootT(t, dir), "link", []byte("hijack")); err == nil {
		t.Fatal("writing through a symlink was not refused")
	}
}

func TestWriteFileInPlace(t *testing.T) {
	path := filepath.Join(t.TempDir(), "registry.vdf")
	if err := os.WriteFile(path, []byte("old content"), 0o600); err != nil {
		t.Fatal(err)
	}
	before := inodeOf(t, path)

	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := WriteFileInPlace(f, []byte("new")); err != nil {
		t.Fatalf("WriteFileInPlace: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "new" {
		t.Errorf("content = %q, want %q", data, "new")
	}
	if got := inodeOf(t, path); got != before {
		t.Errorf("WriteFileInPlace changed the inode: %d -> %d", before, got)
	}
}

func TestOpenVaultRefusesSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	if err := os.WriteFile(target, []byte("t"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "link")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	if v, err := Open(link); err == nil {
		v.Close()
		t.Fatal("a symlinked vault root was not refused")
	}
}

// TestIsOctalMode covers the chmod validator.
func TestIsOctalMode(t *testing.T) {
	valid := []string{"0", "7", "644", "0644", "0777", "0000", "7000"}
	for _, v := range valid {
		if !IsOctalMode(v) {
			t.Errorf("IsOctalMode(%q) = false, want true", v)
		}
	}
	invalid := []string{"", "8", "9", "06444", "0x1ff", "644 ", " -1", "abcd", "064A"}
	for _, v := range invalid {
		if IsOctalMode(v) {
			t.Errorf("IsOctalMode(%q) = true, want false", v)
		}
	}
}

// TestDefaultNewFileMode models the conventional 0644/0755 defaults.
func TestDefaultNewFileMode(t *testing.T) {
	if got := DefaultNewFileMode(false); got != 0o644 {
		t.Errorf("file default = %o, want 0644", got)
	}
	if got := DefaultNewFileMode(true); got != 0o755 {
		t.Errorf("dir default = %o, want 0755", got)
	}
}

// TestUnixModeToFileMode covers the conversion from conventional unix mode
// numbers (with special bits) to Go's os.FileMode flags. A raw cast is not
// equivalent: it maps setuid/setgid/sticky onto undefined FileMode bits
// that chmod(2) ignores, silently reducing 4755 to 0755.
func TestUnixModeToFileMode(t *testing.T) {
	cases := []struct {
		in   uint32
		want os.FileMode
	}{
		{0o644, 0o644},
		{0o755, 0o755},
		{0o0000, 0},
		{0o4000, os.ModeSetuid},
		{0o2000, os.ModeSetgid},
		{0o1000, os.ModeSticky},
		{0o4755, os.ModeSetuid | 0o755},
		{0o2755, os.ModeSetgid | 0o755},
		{0o1777, os.ModeSticky | 0o777},
		{0o7777, os.ModeSetuid | os.ModeSetgid | os.ModeSticky | 0o777},
	}
	for _, tc := range cases {
		if got := UnixModeToFileMode(tc.in); got != tc.want {
			t.Errorf("UnixModeToFileMode(%04o) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

// TestFileModeToUnixModeRoundTrip verifies both converters invert each
// other for every combination that fits in four octal digits.
func TestFileModeToUnixModeRoundTrip(t *testing.T) {
	for _, mode := range []uint32{0o0000, 0o0644, 0o755, 0o1777, 0o2755, 0o4755, 0o6755, 0o7755, 0o7777} {
		if got := FileModeToUnixMode(UnixModeToFileMode(mode)); got != mode {
			t.Errorf("round trip of %04o = %04o", mode, got)
		}
		if back := UnixModeToFileMode(FileModeToUnixMode(UnixModeToFileMode(mode))); back != UnixModeToFileMode(mode) {
			t.Errorf("inverse round trip diverged for %04o", mode)
		}
	}
}

// WriteFileKeepMeta Lstat-checks `path` but creates its temp file at the predictable sibling
// path+".app_listener.edit" with O_CREATE|O_TRUNC and no O_EXCL/O_NOFOLLOW. A symlink planted at
// that name is followed, so root truncates and rewrites the link target — and then chmods/chowns
// it to the edited file's owner.
func TestWriteFileKeepMetaRefusesPlantedTempSibling(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "id_rsa")
	victim := filepath.Join(root, "victim")
	mustWrite(t, target, "original key")
	mustWrite(t, victim, "VICTIM")

	if err := os.Symlink(victim, target+".app_listener.edit"); err != nil {
		t.Fatal(err)
	}

	err := WriteFileKeepMeta(openRootT(t, root), "id_rsa", []byte("edited by root"))

	// Security property: the symlink target must be untouched. The write may still succeed by
	// replacing the real target atomically (the temp name is our namespace, cleared with O_EXCL),
	// but it must never follow the planted link.
	if got := mustRead(t, victim); got != "VICTIM" {
		t.Fatalf("root wrote through the planted temp-sibling symlink: victim = %q (err %v)", got, err)
	}
	if err == nil {
		if got := mustRead(t, target); got != "edited by root" {
			t.Fatalf("write reported success but the real target was not updated: %q", got)
		}
	}
}
