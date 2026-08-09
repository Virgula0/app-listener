package uninstall

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Virgula0/app-listener/internal/fscrypt"
	inst "github.com/Virgula0/app-listener/internal/install"
)

// TestIsInstallerSSHAgentUnit verifies that only a unit whose content
// matches the bundled sample is recognized as installer-installed: a
// missing file and a modified unit are not ours.
func TestIsInstallerSSHAgentUnit(t *testing.T) {
	sample, err := inst.SampleContent("ssh-agent.service")
	if err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()
	ours := filepath.Join(dir, "ours.service")
	if err := os.WriteFile(ours, sample, 0o600); err != nil {
		t.Fatal(err)
	}
	theirs := filepath.Join(dir, "theirs.service")
	if err := os.WriteFile(theirs, append([]byte("# custom\n"), sample...), 0o600); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name string
		path string
		want bool
	}{
		{"matches the bundled sample", ours, true},
		{"modified unit is not ours", theirs, false},
		{"missing unit is not ours", filepath.Join(dir, "nope"), false},
	}
	for _, c := range cases {
		if got := isInstallerSSHAgentUnit(sample, c.path); got != c.want {
			t.Errorf("%s: isInstallerSSHAgentUnit = %v, want %v", c.name, got, c.want)
		}
	}
}

// TestRemoveSSHAgentUnit verifies that reverting one per-user unit removes
// both the default.target.wants symlink and the unit file itself.
func TestRemoveSSHAgentUnit(t *testing.T) {
	home := t.TempDir()
	unitDir := filepath.Join(home, ".config", "systemd", "user")
	unitPath := filepath.Join(unitDir, "ssh-agent.service")
	if err := os.MkdirAll(unitDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(unitPath, []byte("[Unit]\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	wantsDir := filepath.Join(unitDir, "default.target.wants")
	if err := os.MkdirAll(wantsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(unitPath, filepath.Join(wantsDir, "ssh-agent.service")); err != nil {
		t.Fatal(err)
	}

	if err := removeSSHAgentUnit(userSSHAgentUnit{
		User: inst.User{Name: "alice"},
		Path: unitPath,
	}); err != nil {
		t.Fatalf("removeSSHAgentUnit: %v", err)
	}
	if _, err := os.Lstat(unitPath); !os.IsNotExist(err) {
		t.Errorf("unit file still present after revert: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(wantsDir, "ssh-agent.service")); !os.IsNotExist(err) {
		t.Errorf("wants symlink still present after revert: %v", err)
	}
}

// TestRemoveKeyAndEmptyDir verifies the --delete-key behavior: the key file
// is removed and the parent directory is removed too when it becomes empty,
// but kept when it still holds files. A missing key is a no-op success.
func TestRemoveKeyAndEmptyDir(t *testing.T) {
	// Parent becomes empty: both key and directory are removed.
	dir1 := t.TempDir()
	key1 := filepath.Join(dir1, "key")
	if err := os.WriteFile(key1, make([]byte, 32), 0o400); err != nil {
		t.Fatal(err)
	}
	if err := removeKeyAndEmptyDir(key1, dir1); err != nil {
		t.Fatalf("removeKeyAndEmptyDir(empty): %v", err)
	}
	if _, err := os.Lstat(key1); !os.IsNotExist(err) {
		t.Errorf("key still exists: %v", err)
	}
	if _, err := os.Lstat(dir1); !os.IsNotExist(err) {
		t.Errorf("empty parent directory was not removed: %v", err)
	}

	// Parent not empty: key removed, directory kept.
	dir2 := t.TempDir()
	key2 := filepath.Join(dir2, "key")
	if err := os.WriteFile(key2, make([]byte, 32), 0o400); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir2, "other"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := removeKeyAndEmptyDir(key2, dir2); err != nil {
		t.Fatalf("removeKeyAndEmptyDir(non-empty): %v", err)
	}
	if _, err := os.Lstat(key2); !os.IsNotExist(err) {
		t.Errorf("key still exists: %v", err)
	}
	if _, err := os.Lstat(dir2); err != nil {
		t.Errorf("non-empty parent dir must be kept: %v", err)
	}

	// Missing key is a no-op success.
	if err := removeKeyAndEmptyDir(filepath.Join(t.TempDir(), "missing"), dir1); err != nil {
		t.Errorf("missing key must be a no-op success: %v", err)
	}
}

// TestDecryptDirectoriesEmpty is a defensive check that the progress bar
// wrapper is never called with nothing to do.
func TestDecryptDirectoriesEmpty(t *testing.T) {
	if err := decryptDirectories(fscrypt.New(), nil); err != nil {
		t.Errorf("decryptDirectories(nil) = %v, want nil", err)
	}
	if err := decryptDirectories(fscrypt.New(), []string{}); err != nil {
		t.Errorf("decryptDirectories(empty) = %v, want nil", err)
	}
}
