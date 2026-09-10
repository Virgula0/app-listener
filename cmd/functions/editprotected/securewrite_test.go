package editprotected

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteWithin_PlainFile(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "sub", "config")

	if err := writeWithin(root, dest, []byte("hello")); err != nil {
		t.Fatalf("writeWithin: %v", err)
	}
	got, err := os.ReadFile(dest)
	if err != nil || string(got) != "hello" {
		t.Fatalf("read back = %q, %v", got, err)
	}
	info, _ := os.Lstat(dest)
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("mode = %o, want 600", info.Mode().Perm())
	}

	// Overwrite works (atomic replace).
	if err := writeWithin(root, dest, []byte("world")); err != nil {
		t.Fatalf("overwrite: %v", err)
	}
	got, _ = os.ReadFile(dest)
	if string(got) != "world" {
		t.Fatalf("overwrite read back = %q", got)
	}
	// No temp file left behind.
	entries, _ := os.ReadDir(filepath.Dir(dest))
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") && strings.Contains(e.Name(), ".app_listener.put") {
			t.Fatalf("temp file left: %s", e.Name())
		}
	}
}

func TestWriteWithin_RefusesSymlinkTarget(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "victim")
	if err := os.WriteFile(outside, []byte("ORIGINAL"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(root, "config")
	if err := os.Symlink(outside, link); err != nil {
		t.Fatal(err)
	}

	err := writeWithin(root, link, []byte("PWNED"))
	if err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("writeWithin through a symlink = %v, want a symlink refusal", err)
	}
	got, _ := os.ReadFile(outside)
	if string(got) != "ORIGINAL" {
		t.Fatalf("the symlink target was written: %q", got)
	}
}

func TestWriteWithin_RefusesSymlinkDirComponent(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, "evil")); err != nil {
		t.Fatal(err)
	}

	err := writeWithin(root, filepath.Join(root, "evil", "file"), []byte("PWNED"))
	if err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("writeWithin through a symlink dir = %v, want a symlink refusal", err)
	}
	if _, statErr := os.Stat(filepath.Join(outside, "file")); statErr == nil {
		t.Fatal("a file was created outside the root via a symlinked dir component")
	}
}

func TestWriteWithin_RejectsEscape(t *testing.T) {
	root := t.TempDir()
	if err := writeWithin(root, filepath.Join(t.TempDir(), "x"), []byte("x")); err == nil {
		t.Fatal("writeWithin to a path outside root should fail")
	}
}
