package daemon

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/internal/install"
)

func plantBun(t *testing.T, home string) (dir string, plants []string) {
	t.Helper()
	dir = filepath.Join(home, install.BunTmpRelDir)
	mkdirs(t, filepath.Join(dir, "sub"))
	plants = []string{filepath.Join(dir, ".bun-1000-evil.so"), filepath.Join(dir, "sub/.bun-1000-deep.so")}
	for _, p := range plants {
		if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir, plants
}

func TestOpenBunTmp_ReplacesUnvettedDir(t *testing.T) {
	home := t.TempDir()
	dir, plants := plantBun(t, home)
	u := install.User{Name: "u", Home: home}

	bt, err := openBunTmp(u, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer bt.f.Close()
	if !bt.renewed {
		t.Fatal("a dir the daemon never reserved must be replaced")
	}
	for _, p := range plants {
		if _, err := os.Lstat(p); !os.IsNotExist(err) {
			t.Errorf("%s survived the replacement (err %v)", p, err)
		}
	}
	fi, err := os.Stat(dir)
	if err != nil || fi.Mode().Perm() != 0o700 {
		t.Fatalf("fresh dir must be 0700: %v %v", fi, err)
	}
	if id, _ := fileID(bt.f); id != bt.id {
		t.Errorf("returned id %v does not match the open dir %v", bt.id, id)
	}
	left, _ := os.ReadDir(filepath.Dir(dir))
	for _, e := range left {
		if strings.Contains(e.Name(), ".stale.") {
			t.Errorf("stale dir %s left behind", e.Name())
		}
	}
}

func TestOpenBunTmp_KeepsVettedDir(t *testing.T) {
	home := t.TempDir()
	plantBun(t, home)
	u := install.User{Name: "u", Home: home}
	first, err := openBunTmp(u, nil)
	if err != nil {
		t.Fatal(err)
	}
	first.f.Close()
	own := filepath.Join(home, install.BunTmpRelDir, ".bun-1000-own.so")
	if err := os.WriteFile(own, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	again, err := openBunTmp(u, map[inodeID]bool{first.id: true})
	if err != nil {
		t.Fatal(err)
	}
	defer again.f.Close()
	if again.renewed || again.id != first.id {
		t.Errorf("a dir the last apply reserved must be kept: renewed=%v id %v -> %v", again.renewed, first.id, again.id)
	}
	if _, err := os.Stat(own); err != nil {
		t.Errorf("the Bun app's extraction in a vetted dir was removed: %v", err)
	}
}

func TestOpenBunTmp_SymlinkReplacedTargetUntouched(t *testing.T) {
	home, elsewhere := t.TempDir(), t.TempDir()
	target := filepath.Join(elsewhere, ".bun-1000-evil.so")
	if err := os.WriteFile(target, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	mkdirs(t, filepath.Join(home, filepath.Dir(install.BunTmpRelDir)))
	if err := os.Symlink(elsewhere, filepath.Join(home, install.BunTmpRelDir)); err != nil {
		t.Fatal(err)
	}
	bt, err := openBunTmp(install.User{Name: "u", Home: home}, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer bt.f.Close()
	if fi, err := os.Lstat(filepath.Join(home, install.BunTmpRelDir)); err != nil || !fi.IsDir() {
		t.Errorf("the symlink must be replaced by a real dir: %v %v", fi, err)
	}
	if _, err := os.Stat(target); err != nil {
		t.Errorf("the symlink's target must not be touched: %v", err)
	}
}
