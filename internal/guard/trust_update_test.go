package guard

import (
	"os"
	"path/filepath"
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
	SetReplacementCheck(func(path string, o, n GuardInodeKey) bool {
		if path != bin || o != old {
			t.Errorf("check got path=%s old=%+v", path, o)
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
