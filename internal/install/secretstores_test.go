package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

func TestSecretStoreNotes(t *testing.T) {
	home := t.TempDir()
	write := func(rel, text string) {
		p := filepath.Join(home, rel)
		if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(text), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	write(".config/gh/hosts.yml", "github.com:\n    users:\n        alice:\n    user: alice\n")
	write(".docker/config.json", `{"auths":{"ghcr.io":{}},"credsStore":"secretservice"}`)
	write(".gitconfig", "[user]\n\thelper = nope\n[credential \"https://github.com\"]\n\thelper = !gh auth git-credential\n"+
		"[credential]\n\thelper = libsecret\n")
	u := User{Name: "alice", Home: home, UID: uint32(os.Getuid())}
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{
		{Path: filepath.Join(home, ".config/gh")}, {Path: filepath.Join(home, ".docker")},
	}}

	notes := strings.Join(SecretStoreNotes(cfg, u), "\n")
	for _, want := range []string{"--insecure-storage", "docker-credential-secretservice",
		"(!gh auth git-credential, libsecret)"} {
		if !strings.Contains(notes, want) {
			t.Errorf("notes miss %q:\n%s", want, notes)
		}
	}
	if strings.Contains(notes, "nope") {
		t.Errorf("a helper outside [credential] was read:\n%s", notes)
	}

	write(".config/gh/hosts.yml", "github.com:\n    oauth_token: gho_x\n    user: alice\n")
	cfg.Resources = cfg.Resources[:1]
	notes = strings.Join(SecretStoreNotes(cfg, u), "\n")
	if strings.Contains(notes, "gh keeps") || strings.Contains(notes, "docker") {
		t.Errorf("file-stored gh token or unguarded ~/.docker noted:\n%s", notes)
	}
}

func TestReadOwnedRegular_RefusesSymlinkAndOtherOwner(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "t")
	if err := os.WriteFile(target, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "l")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	uid := uint32(os.Getuid())
	if _, err := readOwnedRegular(link, uid); err == nil {
		t.Error("followed a symlink")
	}
	if _, err := readOwnedRegular(target, uid+1); err == nil {
		t.Error("read a file owned by another uid")
	}
	if got, err := readOwnedRegular(target, uid); err != nil || got != "x" {
		t.Errorf("readOwnedRegular = %q, %v", got, err)
	}
}
