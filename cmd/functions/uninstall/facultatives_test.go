package uninstall

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	inst "github.com/Virgula0/app-listener/internal/install"
)

func TestRemoveSSHAgentFacultativeKeepsGuardedSSHConfig(t *testing.T) {
	u := inst.User{Name: "tester", UID: uint32(os.Getuid()), GID: uint32(os.Getgid()), Home: t.TempDir(),
		Shell: "/bin/bash"}
	if err := os.Mkdir(filepath.Join(u.Home, ".ssh"), 0o700); err != nil {
		t.Fatal(err)
	}
	if _, _, err := inst.EnsureAddKeysToAgent(u); err != nil {
		t.Fatal(err)
	}
	if _, err := inst.EnsureSSHAgentEnv(u); err != nil {
		t.Fatal(err)
	}
	unit := inst.SSHAgentUnitPath(u.Home)
	if err := os.MkdirAll(filepath.Dir(unit), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(unit, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := removeFacultative(&installedFacultative{user: u, unit: true, sshGuarded: true}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(unit); !os.IsNotExist(err) {
		t.Errorf("unit still present: %v", err)
	}
	if inst.HasSSHAgentEnv(u) {
		t.Error("SSH_AUTH_SOCK block still present")
	}
	data, err := os.ReadFile(filepath.Join(u.Home, ".ssh", "config"))
	if err != nil || !strings.Contains(string(data), "AddKeysToAgent yes") {
		t.Errorf("guarded ~/.ssh/config must be left alone: %q, %v", data, err)
	}

	if err := removeFacultative(&installedFacultative{user: u}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(u.Home, ".ssh", "config")); !os.IsNotExist(err) {
		t.Errorf("unguarded: the dedicated ~/.ssh/config holding only our block must go: %v", err)
	}
}
