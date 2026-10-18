package guard

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

func TestDiffBinaries(t *testing.T) {
	a := BinaryEntry{Path: "/a", Hash: [32]byte{1}}
	b := BinaryEntry{Path: "/b", Hash: [32]byte{2}}
	bNew := BinaryEntry{Path: "/b", Hash: [32]byte{3}}
	c := BinaryEntry{Path: "/c", Hash: [32]byte{4}}

	added, removed := diffBinaries([]BinaryEntry{a, b}, []BinaryEntry{bNew, c})
	if len(added) != 2 || added[0] != bNew || added[1] != c {
		t.Fatalf("added = %v, want [bNew c]", added)
	}
	if len(removed) != 2 || removed[0] != a || removed[1] != b {
		t.Fatalf("removed = %v, want [a b]: a re-hashed entry is removed and re-added", removed)
	}
	if added, removed := diffBinaries([]BinaryEntry{a}, []BinaryEntry{a}); added != nil || removed != nil {
		t.Fatalf("unchanged list must diff empty, got %v / %v", added, removed)
	}
}

// policyTools copies the test binary to two distinct inodes and returns their entries.
func (s *guardUnitTest) policyTools() (a, b BinaryEntry) {
	var out [2]BinaryEntry
	for i := range out {
		tool := filepath.Join("/tmp", fmt.Sprintf("guard-policy-tool-%d-%d", os.Getpid(), i))
		s.Require().NoError(copySelf(tool))
		s.T().Cleanup(func() { os.Remove(tool) })
		entry, err := ComputeBinaryEntry(tool)
		s.Require().NoError(err)
		out[i] = entry
	}
	return out[0], out[1]
}

func (s *guardUnitTest) seededRoot() (root, seed string) {
	root = s.T().TempDir()
	seed = filepath.Join(root, "seed.txt")
	s.Require().NoError(os.WriteFile(seed, []byte("data"), 0o644))
	return root, seed
}

func (s *guardUnitTest) requireAccess(seed string, allowed map[string]bool) {
	for tool, want := range allowed {
		err := runTool(tool, seed)
		if want {
			s.Require().NoError(err, "%s must be allowed", tool)
		} else {
			s.Require().Error(err, "%s must be denied", tool)
		}
	}
}

func (s *guardUnitTest) TestSetBinariesWhitelistLive() {
	if os.Getuid() != 0 {
		s.T().Skip("Skipping BPF test: requires root")
	}
	root, seed := s.seededRoot()
	a, b := s.policyTools()
	g := s.newGuardedTree(root, []BinaryEntry{a}, nil)
	defer g.Stop()
	s.requireAccess(seed, map[string]bool{a.Path: true, b.Path: false})

	s.Require().NoError(g.SetBinaries([]BinaryEntry{a, b}))
	s.requireAccess(seed, map[string]bool{a.Path: true, b.Path: true})

	s.Require().NoError(g.SetBinaries([]BinaryEntry{b}))
	s.requireAccess(seed, map[string]bool{a.Path: false, b.Path: true})

	s.Require().NoError(g.SetBinaries(nil))
	s.requireAccess(seed, map[string]bool{a.Path: false, b.Path: false})
}

func (s *guardUnitTest) TestSetBinariesBlacklistLive() {
	if os.Getuid() != 0 {
		s.T().Skip("Skipping BPF test: requires root")
	}
	root, seed := s.seededRoot()
	a, b := s.policyTools()
	g, err := NewGuard(root, ModeBlacklist, []BinaryEntry{a}, true, 0, WithEagerPopulate())
	s.Require().NoError(err)
	s.Require().NoError(g.Start())
	defer g.Stop()
	s.requireAccess(seed, map[string]bool{a.Path: false, b.Path: true})

	s.Require().NoError(g.SetBinaries([]BinaryEntry{b}))
	s.requireAccess(seed, map[string]bool{a.Path: true, b.Path: false})

	s.Require().NoError(g.SetBinaries(nil))
	s.requireAccess(seed, map[string]bool{a.Path: true, b.Path: true})
}

// A late SetBinaries must not write into the released slot, which another guard may own.
func (s *guardUnitTest) TestSetBinariesRefusedAfterStop() {
	if os.Getuid() != 0 {
		s.T().Skip("Skipping BPF test: requires root")
	}
	root, _ := s.seededRoot()
	a, _ := s.policyTools()
	holder := s.newGuardedTree(s.T().TempDir(), nil, nil)
	defer holder.Stop()

	g := s.newGuardedTree(root, nil, nil)
	slot := g.resID
	g.Stop()
	s.Require().Error(g.SetBinaries([]BinaryEntry{a}))

	dev, ino, err := ebpf.StatInode(a.Path)
	s.Require().NoError(err)
	var action uint8
	s.Require().Error(holder.objs().GuardExeActions.Lookup(
		GuardResInodeKey{ResId: slot, Ino: GuardInodeKey{Dev: dev, Ino: ino}}, &action),
		"a refused SetBinaries must leave no row in the released slot")
}

// appendByte changes path's content (and so its hash) while keeping it a runnable ELF.
func appendByte(s *guardUnitTest, path string) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0)
	s.Require().NoError(err)
	_, err = f.Write([]byte{0})
	s.Require().NoError(err)
	s.Require().NoError(f.Close())
}

func (s *guardUnitTest) newBlacklist(root string, binaries []BinaryEntry) *Guard {
	g, err := NewGuard(root, ModeBlacklist, binaries, true, 0, WithEagerPopulate())
	s.Require().NoError(err)
	s.Require().NoError(g.Start())
	return g
}

// Re-adding a changed binary diffs as added + removed for one path: the remove must not delete
// the block row the add just wrote.
func (s *guardUnitTest) TestSetBinariesBlacklistReAddChanged() {
	if os.Getuid() != 0 {
		s.T().Skip("Skipping BPF test: requires root")
	}
	for _, replace := range []bool{false, true} {
		root, seed := s.seededRoot()
		a, _ := s.policyTools()
		g := s.newBlacklist(root, []BinaryEntry{a})
		s.requireAccess(seed, map[string]bool{a.Path: false})

		if replace {
			next := a.Path + ".new"
			s.Require().NoError(copySelf(next))
			appendByte(s, next)
			s.Require().NoError(os.Rename(next, a.Path))
		} else {
			appendByte(s, a.Path)
		}
		changed, err := ComputeBinaryEntry(a.Path)
		s.Require().NoError(err)
		s.Require().NotEqual(a.Hash, changed.Hash)

		s.Require().NoError(g.SetBinaries([]BinaryEntry{changed}), "replace=%v", replace)
		s.requireAccess(seed, map[string]bool{a.Path: false})
		g.Stop()
	}
}

func (s *guardUnitTest) sharedInodePaths(a BinaryEntry) map[string]BinaryEntry {
	sym, hard := a.Path+".sym", a.Path+".hard"
	s.Require().NoError(os.Symlink(a.Path, sym))
	s.Require().NoError(os.Link(a.Path, hard))
	s.T().Cleanup(func() { os.Remove(sym); os.Remove(hard) })
	out := make(map[string]BinaryEntry, 2)
	for kind, p := range map[string]string{"symlink": sym, "hardlink": hard} {
		e, err := ComputeBinaryEntry(p)
		s.Require().NoError(err)
		out[kind] = e
	}
	return out
}

// Two listed paths to one inode share its row: dropping one must leave the other enforced.
func (s *guardUnitTest) TestSetBinariesSharedInodeRemoveOne() {
	if os.Getuid() != 0 {
		s.T().Skip("Skipping BPF test: requires root")
	}
	for _, mode := range []Mode{ModeBlacklist, ModeWhitelist} {
		a, _ := s.policyTools()
		for kind, link := range s.sharedInodePaths(a) {
			for _, keepLink := range []bool{false, true} {
				root, seed := s.seededRoot()
				var g *Guard
				if mode == ModeBlacklist {
					g = s.newBlacklist(root, []BinaryEntry{a, link})
				} else {
					g = s.newGuardedTree(root, []BinaryEntry{a, link}, nil)
				}
				remain := a
				if keepLink {
					remain = link
				}
				s.Require().NoError(g.SetBinaries([]BinaryEntry{remain}), "%s %s", modeString(mode), kind)
				s.requireAccess(seed, map[string]bool{remain.Path: mode == ModeWhitelist})
				g.Stop()
			}
		}
	}
}
