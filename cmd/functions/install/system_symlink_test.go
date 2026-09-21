package install

import (
	"os"
	"path/filepath"
	"testing"

	inst "github.com/Virgula0/app-listener/internal/install"
)

// installSSHAgent runs as root and writes a systemd unit into the target user's $HOME, then chowns
// it to that user. It builds the path with filepath.Join and creates it with os.MkdirAll +
// os.WriteFile + os.Chown — all of which resolve symlinks, and MkdirAll reports success when the
// path already resolves to a directory. A user who points $HOME/.config/systemd/user at a
// root-owned directory therefore gets root to create a unit there and hand them ownership of it,
// which is a complete local privilege-escalation primitive.
func TestInstallSSHAgentRefusesSymlinkedUnitDir(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("needs a non-root uid to own the planted target")
	}
	home := t.TempDir()
	outside := t.TempDir() // stands in for /etc/systemd/system

	if err := os.MkdirAll(filepath.Join(home, ".config", "systemd"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(home, ".config", "systemd", "user")); err != nil {
		t.Fatal(err)
	}

	u := inst.User{
		Name:  "tester",
		UID:   uint32(os.Getuid()),
		GID:   uint32(os.Getgid()),
		Home:  home,
		Shell: "/bin/bash",
	}

	err := installSSHAgent(u)

	entries, rerr := os.ReadDir(outside)
	if rerr != nil {
		t.Fatal(rerr)
	}
	if len(entries) != 0 {
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Fatalf("root created %v outside the user's home through a symlinked unit dir", names)
	}
	if err == nil {
		t.Fatal("installSSHAgent traversed a symlinked unit dir; expected a refusal")
	}
}
