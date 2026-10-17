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
