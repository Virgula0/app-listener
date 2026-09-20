package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func rcUser(t *testing.T, shell string) User {
	t.Helper()
	return User{Name: "tester", UID: uint32(os.Getuid()), GID: uint32(os.Getgid()), Home: t.TempDir(), Shell: shell}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestEnsureSSHAgentEnvAppendsOnceAndRemoveRestores(t *testing.T) {
	u := rcUser(t, "/usr/bin/zsh")
	rc := filepath.Join(u.Home, ".zshrc")
	const orig = "alias ll='ls -l'\nexport EDITOR=vim" // no trailing newline
	if err := os.WriteFile(rc, []byte(orig), 0o600); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 2; i++ {
		if _, err := EnsureSSHAgentEnv(u); err != nil {
			t.Fatal(err)
		}
	}
	got := readFile(t, rc)
	if strings.Count(got, rcBegin) != 1 || !strings.Contains(got, "ssh-agent.socket") {
		t.Fatalf("expected exactly one block, got:\n%s", got)
	}
	if !strings.HasPrefix(got, orig+"\n") {
		t.Fatalf("existing content altered:\n%s", got)
	}

	changed, err := RemoveSSHAgentEnv(u)
	if err != nil || len(changed) != 1 {
		t.Fatalf("remove: changed=%v err=%v", changed, err)
	}
	if got := readFile(t, rc); got != orig+"\n" {
		t.Fatalf("remove did not restore the file: %q", got)
	}
	if changed, _ := RemoveSSHAgentEnv(u); len(changed) != 0 {
		t.Fatalf("second remove must be a no-op, changed %v", changed)
	}
}

func TestEnsureSSHAgentEnvKeepsForeignSSHAuthSock(t *testing.T) {
	u := rcUser(t, "/bin/bash")
	rc := filepath.Join(u.Home, ".bashrc")
	const orig = "export SSH_AUTH_SOCK=$HOME/.my-agent.sock\n"
	if err := os.WriteFile(rc, []byte(orig), 0o600); err != nil {
		t.Fatal(err)
	}
	changed, err := EnsureSSHAgentEnv(u)
	if err != nil || len(changed) != 0 {
		t.Fatalf("changed=%v err=%v", changed, err)
	}
	if got := readFile(t, rc); got != orig {
		t.Fatalf("foreign file modified: %q", got)
	}
}

func TestEnsureSSHAgentEnvCreatesLoginShellRCOnly(t *testing.T) {
	u := rcUser(t, "/bin/bash")
	changed, err := EnsureSSHAgentEnv(u)
	if err != nil || len(changed) != 1 || filepath.Base(changed[0]) != ".bashrc" {
		t.Fatalf("changed=%v err=%v", changed, err)
	}
	if _, err := os.Stat(filepath.Join(u.Home, ".zshrc")); !os.IsNotExist(err) {
		t.Fatalf(".zshrc must not be created for a bash user: %v", err)
	}
}

func TestFishDedicatedFileRemovedWithItsBlock(t *testing.T) {
	u := rcUser(t, "/usr/bin/fish")
	changed, err := EnsureSSHAgentEnv(u)
	if err != nil || len(changed) != 1 {
		t.Fatalf("changed=%v err=%v", changed, err)
	}
	if _, err := RemoveSSHAgentEnv(u); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(changed[0]); !os.IsNotExist(err) {
		t.Fatalf("dedicated fish file should be deleted: %v", err)
	}
}

func TestEnsureSSHAgentEnvRefusesForeignOwnedTarget(t *testing.T) {
	u := rcUser(t, "/bin/bash")
	u.UID++ // rc is owned by the test user, not u
	rc := filepath.Join(u.Home, ".bashrc")
	if err := os.WriteFile(rc, []byte("x\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := EnsureSSHAgentEnv(u); err == nil {
		t.Fatal("expected refusal for a file not owned by the user")
	}
	if got := readFile(t, rc); got != "x\n" {
		t.Fatalf("file modified: %q", got)
	}
}

func TestRemoveSSHAgentEnvUnterminatedBlockIsLeftAlone(t *testing.T) {
	u := rcUser(t, "/bin/bash")
	rc := filepath.Join(u.Home, ".bashrc")
	orig := "a\n" + rcBegin + "\nexport X=1\n"
	if err := os.WriteFile(rc, []byte(orig), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := RemoveSSHAgentEnv(u); err == nil {
		t.Fatal("expected an error for an unterminated block")
	}
	if got := readFile(t, rc); got != orig {
		t.Fatalf("file modified: %q", got)
	}
}

func TestEnsureAddKeysToAgentPrependsAndRemoveRestores(t *testing.T) {
	u := rcUser(t, "/bin/bash")
	if err := os.Mkdir(filepath.Join(u.Home, ".ssh"), 0o700); err != nil {
		t.Fatal(err)
	}
	cfg := filepath.Join(u.Home, ".ssh", "config")
	const orig = "Host work\n  User me\n"
	if err := os.WriteFile(cfg, []byte(orig), 0o600); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if _, _, err := EnsureAddKeysToAgent(u); err != nil {
			t.Fatal(err)
		}
	}
	got := readFile(t, cfg)
	if !strings.HasPrefix(got, rcBegin+"\nAddKeysToAgent yes\n"+rcEnd+"\n") || !strings.HasSuffix(got, orig) ||
		strings.Count(got, rcBegin) != 1 {
		t.Fatalf("block must be first, once, original intact:\n%s", got)
	}
	if _, err := RemoveSSHAgentEnv(u); err != nil {
		t.Fatal(err)
	}
	if got := readFile(t, cfg); got != orig {
		t.Fatalf("remove did not restore the config: %q", got)
	}
}

func TestEnsureAddKeysToAgentKeepsUsersOwnSettingAndCreatesFile(t *testing.T) {
	u := rcUser(t, "/bin/bash")
	if err := os.Mkdir(filepath.Join(u.Home, ".ssh"), 0o700); err != nil {
		t.Fatal(err)
	}
	cfg := filepath.Join(u.Home, ".ssh", "config")
	const orig = "Host *\n  addkeystoagent confirm\n"
	if err := os.WriteFile(cfg, []byte(orig), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, changed, err := EnsureAddKeysToAgent(u); err != nil || changed {
		t.Fatalf("changed=%v err=%v", changed, err)
	}
	if got := readFile(t, cfg); got != orig {
		t.Fatalf("user's config modified: %q", got)
	}

	if err := os.Remove(cfg); err != nil {
		t.Fatal(err)
	}
	if _, changed, err := EnsureAddKeysToAgent(u); err != nil || !changed {
		t.Fatalf("changed=%v err=%v", changed, err)
	}
	if fi, err := os.Stat(cfg); err != nil || fi.Mode().Perm() != 0o600 {
		t.Fatalf("created config must be 0600: %v %v", fi, err)
	}
	if _, err := RemoveSSHAgentEnv(u); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(cfg); !os.IsNotExist(err) {
		t.Fatalf("created-by-us config should be deleted on removal: %v", err)
	}
}

func TestEnsureAddKeysToAgentNeverCreatesSSHDir(t *testing.T) {
	u := rcUser(t, "/bin/bash")
	if _, changed, err := EnsureAddKeysToAgent(u); err != nil || changed {
		t.Fatalf("changed=%v err=%v", changed, err)
	}
	if _, err := os.Stat(filepath.Join(u.Home, ".ssh")); !os.IsNotExist(err) {
		t.Fatalf("~/.ssh must not be created: %v", err)
	}
}
