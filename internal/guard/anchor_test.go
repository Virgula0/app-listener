package guard

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

func keyOf(t *testing.T, path string) GuardInodeKey {
	t.Helper()
	dev, ino, err := ebpf.LstatInode(path)
	if err != nil {
		t.Fatal(err)
	}
	return GuardInodeKey{Dev: dev, Ino: ino}
}

func mkTree(t *testing.T) (base, secret string) {
	t.Helper()
	base = t.TempDir()
	if err := os.MkdirAll(filepath.Join(base, "victim", "sub"), 0o755); err != nil {
		t.Fatal(err)
	}
	secret = filepath.Join(base, "victim", "sub", "secret")
	if err := os.WriteFile(secret, []byte("s"), 0o600); err != nil {
		t.Fatal(err)
	}
	return base, secret
}

// A path reaching the inode through a directory symlink must yield the inode's physical ancestors,
// which is what root_in_chain judges, not the ones the path spells.
func TestInodeChainIsPhysicalThroughSymlink(t *testing.T) {
	base, secret := mkTree(t)
	if err := os.MkdirAll(filepath.Join(base, "app"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("../victim", filepath.Join(base, "app", "lib")); err != nil {
		t.Fatal(err)
	}
	key := keyOf(t, secret)

	chain, ok := inodeChain(filepath.Join(base, "app", "lib", "sub", "secret"), key, false)
	if !ok {
		t.Fatal("chain through a directory symlink must be knowable")
	}
	for _, want := range []string{secret, filepath.Join(base, "victim", "sub"), filepath.Join(base, "victim"), base} {
		if !slices.Contains(chain, keyOf(t, want)) {
			t.Errorf("chain lacks physical ancestor %s", want)
		}
	}
	for _, not := range []string{filepath.Join(base, "app"), filepath.Join(base, "app", "lib")} {
		if slices.Contains(chain, keyOf(t, not)) {
			t.Errorf("chain contains %s, which only the path spells", not)
		}
	}
	if chain[0] != key {
		t.Errorf("chain must start at the inode itself")
	}
}

func TestInodeChainDirIncludesItself(t *testing.T) {
	base, _ := mkTree(t)
	dir := filepath.Join(base, "victim")
	chain, ok := inodeChain(dir, keyOf(t, dir), false)
	if !ok || len(chain) < 2 || chain[0] != keyOf(t, dir) || chain[1] != keyOf(t, base) {
		t.Fatalf("directory chain = %v, ok=%v", chain, ok)
	}
}

// Unknowable chains must say so, so the current owner keeps its row.
func TestInodeChainRefusesUncertain(t *testing.T) {
	base, secret := mkTree(t)
	if _, ok := inodeChain(secret, GuardInodeKey{Dev: 1, Ino: 1}, false); ok {
		t.Error("a key that is not the inode at path (a raced swap) must not yield a chain")
	}
	if err := os.Link(secret, filepath.Join(base, "alias")); err != nil {
		t.Fatal(err)
	}
	if _, ok := inodeChain(secret, keyOf(t, secret), false); ok {
		t.Error("a hardlinked file has several parents: its chain is not knowable")
	}
}

// follow=false judges the link entry itself; follow=true its target.
func TestInodeChainFollow(t *testing.T) {
	base, secret := mkTree(t)
	link := filepath.Join(base, "ln")
	if err := os.Symlink(secret, link); err != nil {
		t.Fatal(err)
	}
	if _, ok := inodeChain(link, keyOf(t, secret), false); ok {
		t.Error("nofollow must not verify the target's key at the link")
	}
	chain, ok := inodeChain(link, keyOf(t, secret), true)
	if !ok || !slices.Contains(chain, keyOf(t, filepath.Join(base, "victim", "sub"))) {
		t.Errorf("follow must walk the target's physical chain: %v ok=%v", chain, ok)
	}
}

// A single-file root nested in a read-only tree (fscrypt.key under /etc/app-listener) claims its
// row through the pinned root's proc path; its chain must be the file's, not the magic link's.
func TestInodeChainOfPinnedFileRoot(t *testing.T) {
	base, secret := mkTree(t)
	root, err := openRoot(secret)
	if err != nil {
		t.Fatal(err)
	}
	defer root.close()
	chain, ok := inodeChain(root.proc, root.key, false)
	if !ok || chain[0] != keyOf(t, secret) || !slices.Contains(chain, keyOf(t, base)) {
		t.Fatalf("pinned file root chain = %v, ok=%v", chain, ok)
	}
}

func TestNearestRootClaims(t *testing.T) {
	k := func(i uint64) GuardInodeKey { return GuardInodeKey{Dev: 9, Ino: i} }
	// chain: file(1) -> dir(2) -> dir(3) -> dir(4)
	chain := []GuardInodeKey{k(1), k(2), k(3), k(4)}
	const g, o, h = uint32(5), uint32(6), uint32(7)
	isG := func(id uint32) bool { return id == g }
	for _, tc := range []struct {
		name     string
		roots    map[uint32]GuardInodeKey
		owner    uint32
		outranks bool // g outranks every other resource on a tie
		want     bool
	}{
		{"takeover: owner's root on chain, g's not", map[uint32]GuardInodeKey{g: k(99), o: k(3)}, o, true, false},
		{"same root, g outranked", map[uint32]GuardInodeKey{g: k(3), o: k(3)}, o, false, false},
		{"same root, g outranks", map[uint32]GuardInodeKey{g: k(3), o: k(3)}, o, true, true},
		{"g nested inside owner", map[uint32]GuardInodeKey{g: k(2), o: k(4)}, o, false, true},
		{"owner nested inside g", map[uint32]GuardInodeKey{g: k(4), o: k(2)}, o, true, false},
		{"stale row: owner root off chain", map[uint32]GuardInodeKey{g: k(3), o: k(99)}, o, false, true},
		{"third resource nearer than g", map[uint32]GuardInodeKey{g: k(4), o: k(99), h: k(2)}, o, false, true},
		{"link target in g's tree", map[uint32]GuardInodeKey{g: k(3)}, 0, false, true},
		{"link target outside every tree", map[uint32]GuardInodeKey{g: k(99)}, 0, false, false},
		{"link target in another tree", map[uint32]GuardInodeKey{g: k(99), h: k(2)}, 0, true, false},
		{"link target in a tree nested in g", map[uint32]GuardInodeKey{g: k(4), h: k(2)}, 0, true, false},
		{"link target, g nested in another", map[uint32]GuardInodeKey{g: k(2), h: k(4)}, 0, false, true},
		{"link target, tied root, g outranked", map[uint32]GuardInodeKey{g: k(3), h: k(3)}, 0, false, false},
		{"link target, tied root, g outranks", map[uint32]GuardInodeKey{g: k(3), h: k(3)}, 0, true, true},
	} {
		outranks := func(uint32) bool { return tc.outranks }
		if got := nearestRootClaims(chain, tc.roots, isG, tc.owner, outranks); got != tc.want {
			t.Errorf("%s: got %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestTieRank(t *testing.T) {
	lib, secret := &Guard{mode: ModeReadOnly}, &Guard{mode: ModeWhitelist}
	if tieRank(secret) <= tieRank(lib) {
		t.Error("with equal paths, a whitelist must outrank a read-only resource")
	}
	lib.canonicalRoot.Store(true)
	if tieRank(lib) <= tieRank(secret) {
		t.Error("a path that is its own root must outrank one reaching it through a symlink")
	}
	secret.canonicalRoot.Store(true)
	if tieRank(secret) <= tieRank(lib) {
		t.Error("both canonical: the whitelist must still outrank")
	}
}

func TestMountinfoShowsSbRoot(t *testing.T) {
	mi := []byte("22 1 8:1 / / rw,relatime shared:1 - ext4 /dev/sda1 rw\n" +
		"40 22 8:1 /home/u/src /mnt/bind rw shared:1 - ext4 /dev/sda1 rw\n")
	if whole, err := mountinfoShowsSbRoot(mi, 22); err != nil || !whole {
		t.Errorf("mount 22 shows its sb root: %v %v", whole, err)
	}
	if whole, err := mountinfoShowsSbRoot(mi, 40); err != nil || whole {
		t.Errorf("mount 40 is a bind of a subdirectory: %v %v", whole, err)
	}
	if _, err := mountinfoShowsSbRoot(mi, 7); err == nil {
		t.Error("an unknown mount must be an error")
	}
}
