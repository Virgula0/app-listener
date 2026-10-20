package install

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	inst "github.com/Virgula0/app-listener/internal/install"
)

func facultativeUser(t *testing.T) inst.User {
	t.Helper()
	return inst.User{Name: "tester", UID: uint32(os.Getuid()), GID: uint32(os.Getgid()), Home: t.TempDir(),
		Shell: "/bin/bash"}
}

func TestListFacultativesIgnoresDaemonConf(t *testing.T) {
	u := facultativeUser(t)
	if u.UID == 0 {
		t.Skip("the ssh-agent facultative is never offered for root")
	}
	if err := os.MkdirAll(filepath.Join(u.Home, ".config", "opencode"), 0o700); err != nil {
		t.Fatal(err)
	}
	items, err := listFacultatives(&daemonconfig.Config{}, []inst.User{u, {Name: "root", Home: t.TempDir()}})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 {
		t.Fatalf("want the ssh-agent and the Bun facultative for tester and none for root, got %+v", items)
	}
	if items[0].kind != facultativeSSHAgent || items[0].sshGuarded || items[0].installed {
		t.Errorf("ssh-agent: %+v, want offered, unguarded, not installed", items[0])
	}
	if items[1].kind != facultativeBunTmpdir || !slices.Equal(items[1].launchers, []string{"opencode"}) {
		t.Errorf("Bun: %+v, want the opencode launcher from the app found in the home", items[1])
	}
}

func TestListFacultativesMarksGuardedSSH(t *testing.T) {
	u := facultativeUser(t)
	if u.UID == 0 {
		t.Skip("the ssh-agent facultative is never offered for root")
	}
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{{Path: filepath.Join(u.Home, ".ssh")}}}
	items, err := listFacultatives(cfg, []inst.User{u})
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || !items[0].sshGuarded {
		t.Fatalf("got %+v, want one ssh-agent facultative marked guarded", items)
	}
	if strings.Contains(items[0].label(), "AddKeysToAgent") {
		t.Errorf("label %q promises an edit of a guarded ~/.ssh/config", items[0].label())
	}
}

func TestInstallSSHAgentFacultative(t *testing.T) {
	for _, guarded := range []bool{false, true} {
		u := facultativeUser(t)
		if u.UID == 0 {
			t.Skip("installSSHAgent skips root")
		}
		if err := os.Mkdir(filepath.Join(u.Home, ".ssh"), 0o700); err != nil {
			t.Fatal(err)
		}
		f := facultative{kind: facultativeSSHAgent, user: u, sshGuarded: guarded}
		if err := installFacultatives(&daemonconfig.Config{}, []facultative{f}); err != nil {
			t.Fatalf("guarded=%v: %v", guarded, err)
		}
		if _, err := os.Stat(inst.SSHAgentUnitPath(u.Home)); err != nil {
			t.Errorf("guarded=%v: unit not installed: %v", guarded, err)
		}
		if !inst.HasSSHAgentEnv(u) {
			t.Errorf("guarded=%v: SSH_AUTH_SOCK block missing", guarded)
		}
		_, err := os.Stat(filepath.Join(u.Home, ".ssh", "config"))
		if guarded != os.IsNotExist(err) {
			t.Errorf("guarded=%v: ~/.ssh/config stat err=%v, want it written only when unguarded", guarded, err)
		}
	}
}
