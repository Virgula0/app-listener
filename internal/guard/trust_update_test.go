package guard

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

func TestReSyncBinaries_RefusedReplacementKeepsOldKey(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "app")
	if err := os.WriteFile(bin, []byte("x"), 0o755); err != nil {
		t.Fatal(err)
	}
	dev, ino, err := ebpf.StatInode(bin)
	if err != nil {
		t.Fatal(err)
	}
	fresh := GuardInodeKey{Dev: dev, Ino: ino}
	old := GuardInodeKey{Dev: dev, Ino: ino + 1}
	g := &Guard{path: "/protected", binaries: []BinaryEntry{{Path: bin}},
		deployed: map[string]GuardInodeKey{bin: old}}

	var asked []GuardInodeKey
	SetReplacementCheck(func(path string, f *os.File, o, n GuardInodeKey) bool {
		if path != bin || o != old {
			t.Errorf("check got path=%s old=%+v", path, o)
		}
		// The check judges the very inode about to be admitted.
		if d, i, err := ebpf.StatFile(f); err != nil || (GuardInodeKey{Dev: d, Ino: i}) != n {
			t.Errorf("the check's fd holds %d:%d (%v), the key asked about is %+v", d, i, err, n)
		}
		asked = append(asked, n)
		return false
	})
	defer SetReplacementCheck(nil)

	changed, err := g.ReSyncBinaries()
	if err != nil || changed != 0 {
		t.Fatalf("ReSyncBinaries = %d, %v; want 0, nil", changed, err)
	}
	if g.deployed[bin] != old {
		t.Errorf("a refused replacement must keep the old key, got %+v", g.deployed[bin])
	}
	if len(asked) != 1 || asked[0] != fresh {
		t.Errorf("the check must be asked about the new inode once, got %+v", asked)
	}

	SetReplacementCheck(nil)
	if _, err := g.ReSyncBinaries(); err != nil || g.deployed[bin] != old {
		t.Errorf("no check installed must refuse too: %+v, %v", g.deployed[bin], err)
	}
}

// The prune forgets a freed key before every guard has re-synced the path's replacement: those
// guards must still judge it as a replacement of that key, or an updater's update is refused by
// all but the first guard (Proton's wineserver, whitelisted by several Steam resources).
func TestReSyncBinaries_ReplacementJudgedAgainstPrunedKey(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "app")
	if err := os.WriteFile(bin, []byte("x"), 0o755); err != nil {
		t.Fatal(err)
	}
	dev, ino, err := ebpf.StatInode(bin)
	if err != nil {
		t.Fatal(err)
	}
	old := GuardInodeKey{Dev: dev, Ino: ino + 1}
	g := &Guard{path: "/protected", binaries: []BinaryEntry{{Path: bin}},
		deployed: map[string]GuardInodeKey{bin: old}, keyPaths: map[GuardInodeKey]string{old: bin}}
	g.forgetExe(old)

	var got []GuardInodeKey
	SetReplacementCheck(func(_ string, _ *os.File, o, _ GuardInodeKey) bool {
		got = append(got, o)
		return false
	})
	defer SetReplacementCheck(nil)
	if _, err := g.ReSyncBinaries(); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got, []GuardInodeKey{old}) {
		t.Fatalf("the replacement must be judged against the pruned key %+v, got %+v", old, got)
	}
	if _, ok := g.deployed[bin]; ok {
		t.Fatal("a refused replacement must not be deployed")
	}
}

func TestRefusedReplacements_ReportsEachInodeOnce(t *testing.T) {
	var r refusedReplacements
	k := GuardInodeKey{Dev: 1, Ino: 2}
	if !r.firstTime("/a", k) || r.firstTime("/a", k) {
		t.Error("the same refusal must be reported exactly once")
	}
	if !r.firstTime("/a", GuardInodeKey{Dev: 1, Ino: 3}) {
		t.Error("a different inode at the path must be reported again")
	}
}

func fileKey(t *testing.T, f *os.File) GuardInodeKey {
	t.Helper()
	dev, ino, err := ebpf.StatFile(f)
	if err != nil {
		t.Fatal(err)
	}
	return GuardInodeKey{Dev: dev, Ino: ino}
}

// The re-sync's fd lands on root's /usr/bin/true through a link the test user placed: the name, not
// the inode, decides, so it is refused before any trust map is consulted (the zero TrustGuard has none).
func TestSystemFileRefusesUserPlacedName(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root owns the temp dir")
	}
	link := filepath.Join(t.TempDir(), "tool")
	if err := os.Symlink("/usr/bin/true", link); err != nil {
		t.Fatal(err)
	}
	f, err := ebpf.OpenConfined(link)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if (&TrustGuard{}).systemFile(link, f, fileKey(t, f)) {
		t.Fatal("a user-placed name reaching a root-owned binary was judged a system file")
	}
}

// A root-placed name now holding another inode than the one judged is refused.
func TestSystemFileRefusesAnotherInode(t *testing.T) {
	f, err := ebpf.OpenSystemPlaced("/usr/bin/false")
	if err != nil {
		t.Skip(err)
	}
	defer f.Close()
	if (&TrustGuard{}).systemFile("/usr/bin/true", f, fileKey(t, f)) {
		t.Fatal("the key judged must be the inode the name opens to")
	}
}

func TestPatternMatchesExpandOnlyRootPlacedDirs(t *testing.T) {
	if got := patternMatches("/opt/x/jbr/bin/java"); !slices.Equal(got, []string{"/opt/x/jbr/bin/java"}) {
		t.Errorf("a fixed path = %v", got)
	}
	if os.Geteuid() == 0 {
		return
	}
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "evil", "bin"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("/usr/bin/true", filepath.Join(dir, "evil", "bin", "java")); err != nil {
		t.Fatal(err)
	}
	if got := patternMatches(filepath.Join(dir, "*", "bin", "java")); got != nil {
		t.Errorf("a user-owned directory was expanded: %v", got)
	}
}
