package backups

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Virgula0/app-listener/internal/fscrypt"
	"github.com/Virgula0/app-listener/internal/install"
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

// Exercises find's dedup/matching against fake dependencies, never the real system config or
// catalog: the real Find() probes real catalog paths, some of which (e.g. Steam's registry.vdf) can
// be live guarded resources on a host with the daemon installed, and this test binary isn't
// whitelisted, so it would trip the guard with a spurious DENIED.
func TestFindSmoke(t *testing.T) {
	root := t.TempDir()
	withBackup := filepath.Join(root, "guarded")
	if err := os.MkdirAll(withBackup+fscrypt.BackupSuffix, 0o700); err != nil {
		t.Fatal(err)
	}
	withoutBackup := filepath.Join(root, "plain")

	got, err := find(
		filepath.Join(t.TempDir(), "daemon.conf"), // does not exist: no installed config
		func() ([]install.User, error) { return nil, nil },
		func([]install.User) []install.Candidate {
			return []install.Candidate{{Path: withBackup}, {Path: withoutBackup}}
		},
	)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if len(got) != 1 || got[0].Path != withBackup {
		t.Fatalf("find = %+v, want exactly one entry for %s", got, withBackup)
	}
	if got[0].BackupPath != withBackup+fscrypt.BackupSuffix {
		t.Errorf("inconsistent entry: %+v", got[0])
	}
}
