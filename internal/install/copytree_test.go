package install

import (
	"os"
	"path/filepath"
	"testing"
)

// TestCopyTreeRoundTrip copies a tree with nested directories, files,
// permissions and a symlink, then verifies everything was preserved.
func TestCopyTreeRoundTrip(t *testing.T) {
	src := t.TempDir()

	sub := filepath.Join(src, "sub")
	if err := os.MkdirAll(sub, 0o750); err != nil {
		t.Fatal(err)
	}
	secret := filepath.Join(sub, "secret.key")
	if err := os.WriteFile(secret, []byte("top-secret"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "readme.txt"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("sub/secret.key", filepath.Join(src, "link")); err != nil {
		t.Fatal(err)
	}

	dst := filepath.Join(t.TempDir(), "copy")
	if err := CopyTree(src, dst); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(filepath.Join(dst, "sub"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o750 {
		t.Errorf("sub dir perms = %o, want 750", info.Mode().Perm())
	}

	copied, err := os.ReadFile(filepath.Join(dst, "sub", "secret.key"))
	if err != nil {
		t.Fatal(err)
	}
	if string(copied) != "top-secret" {
		t.Errorf("content = %q", copied)
	}
	finfo, err := os.Stat(filepath.Join(dst, "sub", "secret.key"))
	if err != nil {
		t.Fatal(err)
	}
	if finfo.Mode().Perm() != 0o600 {
		t.Errorf("file perms = %o, want 600", finfo.Mode().Perm())
	}

	linkTarget, err := os.Readlink(filepath.Join(dst, "link"))
	if err != nil {
		t.Fatal(err)
	}
	if linkTarget != "sub/secret.key" {
		t.Errorf("symlink target = %q", linkTarget)
	}
}

// TestCopyTreeOverwriteExistingFile verifies CopyTree replaces an existing
// destination file (used to reinstall the binary over a previous one).
func TestCopyTreeOverwriteExistingFile(t *testing.T) {
	src := filepath.Join(t.TempDir(), "bin")
	if err := os.WriteFile(src, []byte("new-binary"), 0o755); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(t.TempDir(), "bin")
	if err := os.WriteFile(dst, []byte("old-binary"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := CopyTree(src, dst); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(dst)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "new-binary" {
		t.Errorf("dst content = %q, want new-binary", data)
	}
}

// TestCopyTreeProgress verifies the copy callback reports monotonically
// increasing byte counts ending at the total size of all regular files
// (symlinks contribute no bytes).
func TestCopyTreeProgress(t *testing.T) {
	src := t.TempDir()
	sub := filepath.Join(src, "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}
	// 3-byte and 2-byte files plus a symlink.
	if err := os.WriteFile(filepath.Join(src, "a"), []byte("abc"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "b"), []byte("de"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("a", filepath.Join(src, "link")); err != nil {
		t.Fatal(err)
	}

	var reports [][2]int64
	err := CopyTreeWithProgress(src, filepath.Join(t.TempDir(), "copy"), func(copied, total int64) {
		reports = append(reports, [2]int64{copied, total})
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(reports) != 2 {
		t.Fatalf("expected 2 reports (one per file), got %d: %v", len(reports), reports)
	}
	if reports[0][1] != 5 || reports[1][1] != 5 {
		t.Errorf("total = %v, want 5 bytes in all reports", reports)
	}
	if reports[1][0] != 5 {
		t.Errorf("final copied = %d, want 5", reports[1][0])
	}
	for i := 1; i < len(reports); i++ {
		if reports[i][0] < reports[i-1][0] {
			t.Errorf("progress went backwards: %v", reports)
		}
	}
}
