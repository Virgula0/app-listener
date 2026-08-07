package protected

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Virgula0/app-listener/internal/fscrypt"
)

// TestExistingEncryptedDirsSkipsPlain verifies that the catalog re-scan
// only returns directories that currently carry an fscrypt policy: plain
// directories and regular files are skipped. The paths are plain here, so
// the result must be empty — and no path may error.
func TestExistingEncryptedDirsSkipsPlain(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "d")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(base, "f")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := ExistingEncryptedDirs(fscrypt.New(), []string{dir, file})
	if err != nil {
		t.Fatalf("ExistingEncryptedDirs: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("encrypted = %v, want none for plain paths", got)
	}
}

// TestExistingEncryptedDirsNilPaths is a defensive check that the scan
// tolerates an empty path set.
func TestExistingEncryptedDirsNilPaths(t *testing.T) {
	got, err := ExistingEncryptedDirs(fscrypt.New(), nil)
	if err != nil {
		t.Fatalf("ExistingEncryptedDirs(nil): %v", err)
	}
	if len(got) != 0 {
		t.Errorf("got = %v, want empty", got)
	}
}
