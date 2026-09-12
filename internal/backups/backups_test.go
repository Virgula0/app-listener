package backups

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDelete(t *testing.T) {
	root := t.TempDir()
	var entries []Backup
	for _, n := range []string{"a", "b"} {
		bp := filepath.Join(root, n+".app_listener.backup")
		if err := os.MkdirAll(filepath.Join(bp, "sub"), 0o700); err != nil {
			t.Fatal(err)
		}
		entries = append(entries, Backup{Path: filepath.Join(root, n), BackupPath: bp})
	}

	if err := Delete(entries); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	for _, e := range entries {
		if _, err := os.Lstat(e.BackupPath); !os.IsNotExist(err) {
			t.Errorf("%s still exists after Delete (err=%v)", e.BackupPath, err)
		}
	}
}

// TestFindSmoke: Find must not panic or error on a host with no installed
// config and no backups (it returns an empty slice).
func TestFindSmoke(t *testing.T) {
	got, err := Find()
	if err != nil {
		t.Fatalf("Find: %v", err)
	}
	for _, b := range got {
		if b.BackupPath != b.Path+".app_listener.backup" {
			t.Errorf("inconsistent entry: %+v", b)
		}
	}
}
