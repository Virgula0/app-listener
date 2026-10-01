package daemon

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/install"
)

// discordRefreshFixture: a Discord section naming version 0.0.1 while only 0.0.2 is installed.
func discordRefreshFixture(t *testing.T) (r *catalogRefresher, confPath, oldBin, newBin string, reloads *int) {
	t.Helper()
	home := t.TempDir()
	root := filepath.Join(home, ".config", "discord")
	mkdirs(t, filepath.Join(root, "sentry"), filepath.Join(root, "0.0.2"))
	newBin = filepath.Join(root, "0.0.2", "Discord")
	if err := os.WriteFile(newBin, []byte("x"), 0o755); err != nil {
		t.Fatal(err)
	}
	oldBin = filepath.Join(root, "0.0.1", "Discord")
	confPath = filepath.Join(t.TempDir(), "daemon.conf")
	conf := "[watch " + root + "]\nneed_encryption: false\n" + oldBin + "\n"
	if err := os.WriteFile(confPath, []byte(conf), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := daemonconfig.Load(confPath)
	if err != nil {
		t.Fatal(err)
	}
	reloads = new(int)
	r = &catalogRefresher{
		configPath: confPath,
		cfg:        cfg,
		listUsers:  func() ([]install.User, error) { return []install.User{{Name: "u", Home: home}}, nil },
		reload: func() *daemonconfig.Config {
			*reloads++
			next, err := daemonconfig.Load(confPath)
			if err != nil {
				t.Fatal(err)
			}
			return next
		},
	}
	return r, confPath, oldBin, newBin, reloads
}

func TestRefreshOnceAppliesAndReloads(t *testing.T) {
	r, confPath, oldBin, newBin, reloads := discordRefreshFixture(t)
	if err := r.refreshOnce(r.cfg); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(confPath)
	if !strings.Contains(string(got), newBin) || strings.Contains(string(got), oldBin) {
		t.Fatalf("daemon.conf must name the installed version only:\n%s", got)
	}
	if *reloads != 1 {
		t.Fatalf("%d reloads, want 1", *reloads)
	}
}

func TestRefreshOnceSkipsAConfigChangedSinceLoad(t *testing.T) {
	r, confPath, _, _, reloads := discordRefreshFixture(t)
	edited := append(append([]byte(nil), r.cfg.Raw...), []byte("# operator edit\n")...)
	if err := os.WriteFile(confPath, edited, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := r.refreshOnce(r.cfg); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(confPath); string(got) != string(edited) || *reloads != 0 {
		t.Fatalf("an unreloaded operator edit must be left alone: reloads=%d file=\n%s", *reloads, got)
	}
}

func TestRefreshOnceRestoresTheFileWhenTheReloadFails(t *testing.T) {
	r, confPath, _, _, _ := discordRefreshFixture(t)
	before := append([]byte(nil), r.cfg.Raw...)
	r.reload = func() *daemonconfig.Config { return nil }
	if err := r.refreshOnce(r.cfg); err == nil {
		t.Fatal("a failed reload must be reported")
	}
	if got, _ := os.ReadFile(confPath); string(got) != string(before) {
		t.Fatalf("daemon.conf must keep describing the running config:\n%s", got)
	}
}

func TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(t *testing.T) {
	r, confPath, oldBin, newBin, reloads := discordRefreshFixture(t)
	// An operator's own line the catalog doesn't produce, and the installed version already listed:
	// nothing the refresh would change.
	conf := strings.Replace(string(r.cfg.Raw), oldBin, `"`+newBin+`"`+"\n/usr/bin/true", 1)
	if err := os.WriteFile(confPath, []byte(conf), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg, err := daemonconfig.Load(confPath)
	if err != nil {
		t.Fatal(err)
	}
	r.cfg = cfg
	if err := r.refreshOnce(cfg); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(confPath); string(got) != conf || *reloads != 0 {
		t.Fatalf("a refresh with nothing to admit or drop must not touch the file nor reload: reloads=%d\n%s",
			*reloads, got)
	}
}
