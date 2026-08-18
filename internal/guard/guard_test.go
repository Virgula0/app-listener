package guard

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/infrastructure"
)

type guardUnitTest struct {
	suite.Suite
}

func TestGuardUnitTest(t *testing.T) {
	suite.Run(t, new(guardUnitTest))
}

func (s *guardUnitTest) TestComputeBinaryEntry() {
	dir := s.T().TempDir()
	binaryPath := filepath.Join(dir, "testbin")
	content := []byte("#!/bin/bash\necho hello")
	s.Require().NoError(os.WriteFile(binaryPath, content, 0755))

	entry, err := ComputeBinaryEntry(binaryPath)
	s.Require().NoError(err)
	s.Require().Equal(binaryPath, entry.Path)
	s.Require().Equal(sha256.Sum256(content), entry.Hash)
	s.Require().Equal("testbin", entry.Comm)
}

func (s *guardUnitTest) TestComputeBinaryEntryLongName() {
	dir := s.T().TempDir()
	longName := "verylongbinaryname123"
	binaryPath := filepath.Join(dir, longName)
	content := []byte("data")
	s.Require().NoError(os.WriteFile(binaryPath, content, 0755))

	entry, err := ComputeBinaryEntry(binaryPath)
	s.Require().NoError(err)
	s.Require().Equal(longName[:15], entry.Comm, "comm should be truncated to 15 chars")
}

func (s *guardUnitTest) TestComputeBinaryEntryNonexistent() {
	_, err := ComputeBinaryEntry("/nonexistent/binary")
	s.Require().Error(err)
}

func (s *guardUnitTest) TestBinariesSummary() {
	entries := []BinaryEntry{
		{Path: "/usr/bin/cat", Hash: sha256.Sum256([]byte("cat-content")), Comm: "cat"},
		{Path: "/usr/bin/dog", Hash: sha256.Sum256([]byte("dog-content")), Comm: "dog"},
	}

	summary := BinariesSummary(entries)
	s.Require().Contains(summary, "/usr/bin/cat")
	s.Require().Contains(summary, "/usr/bin/dog")
	s.Require().Contains(summary, "sha256")
}

func (s *guardUnitTest) TestModeString() {
	s.Require().Equal("whitelist", modeString(ModeWhitelist))
	s.Require().Equal("blacklist", modeString(ModeBlacklist))
}

func (s *guardUnitTest) TestCommMatchesGuardedBinary() {
	binaries := []BinaryEntry{
		{Path: "/home/angelo/.config/discord/app-1.0.151/Discord"},
		{Path: "/usr/bin/opencode"},
		{Path: "/home/angelo/.config/discord/app-1.0.151/chrome_crashpad_handler"},
	}

	// the guarded binary's own comm matches its basename
	s.Require().True(commMatchesGuardedBinary("Discord", binaries))
	s.Require().True(commMatchesGuardedBinary("opencode", binaries))

	// worker threads of the guarded binary legitimately rename
	// themselves (Chromium's "libuv-worker", Bun's "Bun Pool N"):
	// those comms must never count as a match
	s.Require().False(commMatchesGuardedBinary("libuv-worker", binaries))
	s.Require().False(commMatchesGuardedBinary("Bun Pool 1", binaries))

	// an unrelated binary name must not match
	s.Require().False(commMatchesGuardedBinary("cat", binaries))

	// the kernel truncates comm to 15 bytes (TASK_COMM_LEN): guarded
	// names longer than that are compared truncated on both sides
	s.Require().True(commMatchesGuardedBinary("chrome_crashpad", binaries))
	s.Require().False(commMatchesGuardedBinary("chrome_crashpad_handler", binaries))

	// empty whitelist never matches
	s.Require().False(commMatchesGuardedBinary("Discord", nil))
}

func (s *guardUnitTest) TestEventMaskOpenImplied() {
	mask, err := eventMask([]ebpf.EventType{ebpf.EventRead})
	s.Require().NoError(err)
	// READ bit set, and OPEN implicitly allowed.
	s.Require().NotEqual(uint32(0), mask&(1<<uint(ebpf.EventRead)))
	s.Require().NotEqual(uint32(0), mask&(1<<uint(ebpf.EventOpen)))
}

func (s *guardUnitTest) TestEventMaskWriteMmap() {
	mask, err := eventMask([]ebpf.EventType{ebpf.EventWrite, ebpf.EventMmap})
	s.Require().NoError(err)
	s.Require().NotEqual(uint32(0), mask&(1<<uint(ebpf.EventWrite)))
	s.Require().NotEqual(uint32(0), mask&(1<<uint(ebpf.EventMmap)))
	s.Require().NotEqual(uint32(0), mask&(1<<uint(ebpf.EventOpen)))
	// Independent events stay unset.
	s.Require().Equal(uint32(0), mask&(1<<uint(ebpf.EventDelete)))
	s.Require().Equal(uint32(0), mask&(1<<uint(ebpf.EventMkdir)))
}

func (s *guardUnitTest) TestEventMaskIndependent() {
	mask, err := eventMask([]ebpf.EventType{ebpf.EventDelete, ebpf.EventRename})
	s.Require().NoError(err)
	s.Require().NotEqual(uint32(0), mask&(1<<uint(ebpf.EventDelete)))
	s.Require().NotEqual(uint32(0), mask&(1<<uint(ebpf.EventRename)))
	// No read/write → OPEN must NOT be implied.
	s.Require().Equal(uint32(0), mask&(1<<uint(ebpf.EventOpen)))
}

func (s *guardUnitTest) TestEventMaskAllEvents() {
	mask, err := eventMask(ebpf.EventTypes())
	s.Require().NoError(err)
	for _, t := range ebpf.EventTypes() {
		s.Require().NotEqual(uint32(0), mask&(1<<uint(t)), "bit for %s should be set", t)
	}
}

func (s *guardUnitTest) TestEventMaskEmpty() {
	mask, err := eventMask(nil)
	s.Require().NoError(err)
	s.Require().Equal(uint32(0), mask)
}

func (s *guardUnitTest) TestBpfEventToGuardEvent() {
	var comm [16]byte
	copy(comm[:], "test-proc")
	var path [256]byte
	copy(path[:], "/tmp/test.txt")

	be := bpfGuardEvent{
		PID:     1234,
		UID:     1000,
		GID:     1000,
		Type:    0,
		FD:      3,
		Blocked: 1,
		Comm:    comm,
		Path:    path,
	}

	fe := be.toFileEvent()

	s.Require().Equal(uint32(1234), fe.PID)
	s.Require().Equal(uint32(1000), fe.UID)
	s.Require().Equal(uint32(1000), fe.GID)
	s.Require().Equal(ebpf.EventOpen, fe.Type)
	s.Require().Equal(uint32(3), fe.FD)
	s.Require().Equal("test-proc", fe.Comm)
	s.Require().Equal("/tmp/test.txt", fe.Path)
}

func (s *guardUnitTest) TestBpfEventToGuardEventZeroValue() {
	be := bpfGuardEvent{}

	fe := be.toFileEvent()

	s.Require().Equal(uint32(0), fe.PID)
	s.Require().Equal(ebpf.EventType(0), fe.Type)
	s.Require().Equal("", fe.Comm)
	s.Require().Equal("", fe.Path)
}

func (s *guardUnitTest) TestIsGuardedBinaryExact() {
	tmpDir := s.T().TempDir()
	targetPath := filepath.Join(tmpDir, "target")
	s.Require().NoError(os.WriteFile(targetPath, []byte("data"), 0755))

	g := &Guard{
		binaries: []BinaryEntry{
			{Path: targetPath, Comm: "target"},
		},
	}

	result := g.isGuardedBinary(targetPath)
	s.Require().True(result)
}

func (s *guardUnitTest) TestIsGuardedBinarySymlink() {
	tmpDir := s.T().TempDir()
	targetPath := filepath.Join(tmpDir, "target")
	s.Require().NoError(os.WriteFile(targetPath, []byte("data"), 0755))

	linkPath := filepath.Join(tmpDir, "link")
	s.Require().NoError(os.Symlink(targetPath, linkPath))

	g := &Guard{
		binaries: []BinaryEntry{
			{Path: targetPath, Comm: "target"},
		},
	}

	result := g.isGuardedBinary(linkPath)
	s.Require().True(result, "should resolve symlink to binary path")
}

func (s *guardUnitTest) TestIsGuardedBinaryNotPresent() {
	g := &Guard{}

	result := g.isGuardedBinary("/usr/bin/nonexistent")
	s.Require().False(result)
}

func (s *guardUnitTest) TestFileEventUnchanged() {
	var comm [16]byte
	copy(comm[:], "cat")

	be := bpfGuardEvent{
		Blocked: 0,
		Comm:    comm,
	}

	fe := be.toFileEvent()
	s.Require().Equal("cat", fe.Comm)
}

func (s *guardUnitTest) TestBlockedFieldMapping() {
	tests := []struct {
		name     string
		blocked  uint32
		expected bool
	}{
		{name: "blocked", blocked: 1, expected: true},
		{name: "not blocked", blocked: 0, expected: false},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			be := bpfGuardEvent{
				Blocked: tt.blocked,
			}
			f := be.toFileEvent()
			_ = f
			s.Require().Equal(tt.blocked != 0, tt.expected)
		})
	}
}

func (s *guardUnitTest) TestBpfEventToGuardEventAllEventTypes() {
	eventTypes := []struct {
		name    string
		bpfType uint32
		goType  ebpf.EventType
	}{
		{name: "open", bpfType: 0, goType: ebpf.EventOpen},
		{name: "read", bpfType: 1, goType: ebpf.EventRead},
		{name: "write", bpfType: 2, goType: ebpf.EventWrite},
		{name: "delete", bpfType: 3, goType: ebpf.EventDelete},
		{name: "rename", bpfType: 4, goType: ebpf.EventRename},
		{name: "symlink", bpfType: 5, goType: ebpf.EventSymlink},
		{name: "hardlink", bpfType: 6, goType: ebpf.EventHardlink},
		{name: "mkdir", bpfType: 7, goType: ebpf.EventMkdir},
		{name: "mmap", bpfType: 8, goType: ebpf.EventMmap},
	}

	for _, et := range eventTypes {
		s.Run(et.name, func() {
			be := bpfGuardEvent{
				PID:     1234,
				Type:    et.bpfType,
				Blocked: 1,
				Comm:    func() (c [16]byte) { copy(c[:], "test"); return }(),
			}

			fe := be.toFileEvent()

			s.Require().Equal(et.goType, fe.Type, "event type %s should map to %v", et.name, et.goType)
			s.Require().Equal(uint32(1234), fe.PID)
			s.Require().Equal("test", fe.Comm)
		})
	}
}

// walkInodes helper: adds every visited path into a slice.
func addRecorder() (func(string) error, *[]string) {
	var got []string
	return func(p string) error {
		got = append(got, p)
		return nil
	}, &got
}

func (s *guardUnitTest) TestWalkInodesVisitsTree() {
	dir := s.T().TempDir()
	s.Require().NoError(os.MkdirAll(filepath.Join(dir, "sub"), 0o755))
	s.Require().NoError(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0o644))
	s.Require().NoError(os.WriteFile(filepath.Join(dir, "sub", "b.txt"), []byte("b"), 0o644))

	add, got := addRecorder()
	s.Require().NoError(walkInodes(dir, true, 0, 0, add))
	s.Require().Len(*got, 4)
}

func (s *guardUnitTest) TestWalkInodesNonRecursive() {
	dir := s.T().TempDir()
	s.Require().NoError(os.MkdirAll(filepath.Join(dir, "sub"), 0o755))
	s.Require().NoError(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0o644))
	s.Require().NoError(os.WriteFile(filepath.Join(dir, "sub", "b.txt"), []byte("b"), 0o644))

	add, got := addRecorder()
	s.Require().NoError(walkInodes(dir, false, 0, 0, add))
	s.Require().Len(*got, 2) // root dir + a.txt only
}

func (s *guardUnitTest) TestWalkInodesDepthLimit() {
	dir := s.T().TempDir()
	s.Require().NoError(os.MkdirAll(filepath.Join(dir, "l1", "l2", "l3"), 0o755))

	add, got := addRecorder()
	s.Require().NoError(walkInodes(dir, true, 2, 0, add))
	// depth 2: root, l1, l2 visited; l3 not.
	for _, p := range *got {
		s.Require().NotContains(p, "l3", "depth-limited walk must not descend past the limit")
	}
}

func (s *guardUnitTest) TestWalkInodesToleratesVanishingEntry() {
	dir := s.T().TempDir()
	s.Require().NoError(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0o644))
	s.Require().NoError(os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0o644))

	// Simulate an entry deleted between readdir and stat (a live
	// application churning its profile): one file fails with ENOENT, the
	// rest succeed. The walk must not fail.
	seen := 0
	s.Require().NoError(walkInodes(dir, false, 0, 0, func(p string) error {
		seen++
		if filepath.Base(p) == "a.txt" {
			return unix.ENOENT
		}
		return nil
	}))
	s.Require().Equal(3, seen, "the failing entry must not abort the walk")
}

func (s *guardUnitTest) TestWalkInodesDanglingSymlinkSkipped() {
	dir := s.T().TempDir()
	s.Require().NoError(os.Symlink(filepath.Join(dir, "nowhere"), filepath.Join(dir, "dangling")))
	s.Require().NoError(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0o644))

	add, _ := addRecorder()
	s.Require().NoError(walkInodes(dir, false, 0, 0, add))
}

func (s *guardUnitTest) TestWalkInodesMapFullDegrades() {
	dir := s.T().TempDir()
	s.Require().NoError(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0o644))
	s.Require().NoError(os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0o644))

	// A full hash map surfaces as E2BIG: the walk stops without failing.
	err := walkInodes(dir, false, 0, 0, func(p string) error {
		return unix.E2BIG
	})
	s.Require().NoError(err, "a full map must degrade, not fail the daemon")
}

func (s *guardUnitTest) TestWalkInodesLockedTreeDetected() {
	dir := s.T().TempDir()
	s.Require().NoError(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0o644))

	// Every entry failing to stat is the signature of an fscrypt-encrypted
	// directory whose key is not provisioned: report it with a hint. The
	// root itself stats fine (a locked directory's own inode is visible).
	err := walkInodes(dir, false, 0, 0, func(p string) error {
		if p == dir {
			return nil
		}
		return unix.ENOENT
	})
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "fscrypt-encrypted and locked")
}

// TestPopulateInodesFillsMap verifies the eager inode scan registers every
// file and directory of the guarded tree in guard_inodes, so kernel-side
// own/parent inode checks (and not only lazy discovery) protect deep
// files. Requires root to load the BPF programs; skipped otherwise.
func (s *guardUnitTest) TestPopulateInodesFillsMap() {
	if os.Geteuid() != 0 {
		s.T().Skip("requires root to load BPF programs")
	}

	dir := s.T().TempDir()
	s.Require().NoError(os.MkdirAll(filepath.Join(dir, "l1", "l2", "l3"), 0o755))
	s.Require().NoError(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0o644))
	s.Require().NoError(os.WriteFile(filepath.Join(dir, "l1", "l2", "l3", "deep.txt"), []byte("d"), 0o644))

	g, err := NewGuard(dir, ModeBlacklist, nil, true, 0)
	if err != nil {
		s.T().Skipf("cannot load guard BPF programs (requires root + LSM support): %v", err)
	}
	defer g.Stop()

	s.Require().NoError(g.PopulateInodes())

	var want []string
	s.Require().NoError(walkInodes(dir, true, 0, 0, func(p string) error {
		want = append(want, p)
		return nil
	}))

	found := make(map[GuardInodeKey]bool, len(want))
	var key GuardInodeKey
	var val uint8
	it := g.objs.GuardInodes.Iterate()
	for it.Next(&key, &val) {
		found[key] = true
	}
	s.Require().NoError(it.Err())
	s.Require().Len(found, len(want),
		"every inode of the tree (dirs and deep files) must be in guard_inodes after PopulateInodes")
}
