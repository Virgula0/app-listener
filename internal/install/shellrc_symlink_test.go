package install

import (
	"os"
	"path/filepath"
	"testing"
)

// EnsureSSHAgentEnv runs as root and writes into the target user's $HOME. createOwned correctly
// uses O_EXCL|O_NOFOLLOW, but that only protects the FINAL component: mkdirAllOwned probes each
// parent with os.Stat, which follows symlinks, and returns success when the path resolves to a
// directory. A user who points a parent directory at a location outside their home therefore gets
// root to create a file there and chown it to them.
func TestEnsureSSHAgentEnvRefusesSymlinkedParent(t *testing.T) {
	u := rcUser(t, "/usr/bin/fish")
	outside := t.TempDir() // stands in for a root-owned directory such as /etc/systemd/system

	// The fish drop-in dir is $HOME/.config/fish/conf.d; the attacker owns $HOME and points the
	// last parent component at a directory outside it.
	if err := os.MkdirAll(filepath.Join(u.Home, ".config", "fish"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(u.Home, ".config", "fish", "conf.d")); err != nil {
		t.Fatal(err)
	}

	_, err := EnsureSSHAgentEnv(u)

	entries, rerr := os.ReadDir(outside)
	if rerr != nil {
		t.Fatal(rerr)
	}
	if len(entries) != 0 {
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Fatalf("root created %v outside the user's home through a symlinked parent", names)
	}
	if err == nil {
		t.Fatal("EnsureSSHAgentEnv traversed a symlinked parent; expected a refusal")
	}
}

// RemoveSSHAgentEnv checks the dedicated fish drop-in through a descriptor, then deletes it by
// path. The removal must act on the file that was checked, even if a parent directory is replaced
// in between; an unrelated same-named file elsewhere must survive.
func TestRemoveSSHAgentEnvRemovesOnlyTheCheckedFile(t *testing.T) {
	u := rcUser(t, "/usr/bin/fish")
	if _, err := EnsureSSHAgentEnv(u); err != nil {
		t.Fatal(err)
	}
	confd := filepath.Join(u.Home, ".config", "fish", "conf.d")
	name := "app-listener-ssh-agent.fish"

	other := t.TempDir()
	unrelated := filepath.Join(other, name)
	if err := os.WriteFile(unrelated, []byte("unrelated\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	defer func(prev func(string)) { testHookBeforeRCRemove = prev }(testHookBeforeRCRemove)
	testHookBeforeRCRemove = func(string) {
		if err := os.Rename(confd, confd+".moved"); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(other, confd); err != nil {
			t.Fatal(err)
		}
	}

	_, _ = RemoveSSHAgentEnv(u)

	if _, err := os.Stat(unrelated); err != nil {
		t.Fatalf("a file outside the checked directory was removed: %v", err)
	}
}
