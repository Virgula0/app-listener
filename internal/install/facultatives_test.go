package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPathCovered(t *testing.T) {
	covered := []string{"/home/u/.ssh", "/home/u/.config/discord", ""}
	cases := []struct {
		path string
		want bool
	}{
		{"/home/u/.ssh", true},                          // exact
		{"/home/u/.ssh/keys", true},                     // inside a covered root
		{"/home/u", true},                               // parent of a covered root — would overlap
		{"/home/u/.config/discord/Local Storage", true}, // inside grouped root
		{"/home/u/.gnupg", false},                       // unrelated
		{"/home/u/.sshd", false},                        // prefix but not a path component
	}
	for _, c := range cases {
		if got := PathCovered(c.path, covered); got != c.want {
			t.Errorf("PathCovered(%q) = %v, want %v", c.path, got, c.want)
		}
	}
}

func TestLoadConfigIfPresentMissingIsEmpty(t *testing.T) {
	cfg, err := LoadConfigIfPresent(filepath.Join(t.TempDir(), "daemon.conf"))
	if err != nil || cfg == nil || len(cfg.Resources) != 0 {
		t.Fatalf("missing daemon.conf: cfg=%v err=%v, want an empty config", cfg, err)
	}
}

func TestLoadConfigIfPresentUnreadableFails(t *testing.T) {
	if _, err := LoadConfigIfPresent(t.TempDir()); err == nil {
		t.Fatal("an unreadable daemon.conf must fail, not read as empty")
	}
}

func TestHasRCBlocksTrackEnsureAndRemove(t *testing.T) {
	u := rcUser(t, "/bin/zsh")
	if HasSSHAgentEnv(u) || HasBunLauncherEnv(u) {
		t.Fatal("fresh home reports a block")
	}
	if _, err := EnsureSSHAgentEnv(u); err != nil {
		t.Fatal(err)
	}
	if _, err := EnsureBunLauncherEnv(u, []string{"opencode"}); err != nil {
		t.Fatal(err)
	}
	if !HasSSHAgentEnv(u) || !HasBunLauncherEnv(u) {
		t.Fatal("blocks not detected after Ensure")
	}
	if _, err := RemoveBunLauncherEnv(u); err != nil {
		t.Fatal(err)
	}
	if HasBunLauncherEnv(u) || !HasSSHAgentEnv(u) {
		t.Fatal("removing the Bun block must leave the ssh-agent one")
	}
	if _, err := RemoveSSHAgentShellEnv(u); err != nil {
		t.Fatal(err)
	}
	if HasSSHAgentEnv(u) {
		t.Fatal("SSH_AUTH_SOCK block still detected after removal")
	}
}

func TestRemoveSSHAgentShellEnvKeepsSSHConfig(t *testing.T) {
	u := rcUser(t, "/bin/bash")
	if err := os.Mkdir(filepath.Join(u.Home, ".ssh"), 0o700); err != nil {
		t.Fatal(err)
	}
	if _, _, err := EnsureAddKeysToAgent(u); err != nil {
		t.Fatal(err)
	}
	if _, err := RemoveSSHAgentShellEnv(u); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(u.Home, ".ssh", "config"))
	if err != nil || !strings.Contains(string(data), "AddKeysToAgent yes") {
		t.Fatalf("~/.ssh/config must be untouched: %q, %v", data, err)
	}
}

// An earlier release's unit stays ours (uninstall removes it, facultatives offers the upgrade); the
// bundled one binds inside ~/.ssh, never the unguarded runtime dir.
func TestIsInstallerSSHAgentUnit(t *testing.T) {
	sample, err := SampleContent("ssh-agent.service")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		name string
		data string
		want bool
	}{
		{"bundled", string(sample), true},
		{"earlier release", legacySSHAgentUnit, true},
		{"modified", "# custom\n" + string(sample), false},
		{"empty", "", false},
	} {
		if got := IsInstallerSSHAgentUnit([]byte(c.data)); got != c.want {
			t.Errorf("%s: IsInstallerSSHAgentUnit = %v, want %v", c.name, got, c.want)
		}
	}
	// The runtime catalog entry guards exactly this directory.
	if !strings.Contains(string(sample), "-a %t/ssh-agent/socket") {
		t.Fatalf("the bundled ssh-agent must bind in the guarded %%t/ssh-agent:\n%s", sample)
	}
}
