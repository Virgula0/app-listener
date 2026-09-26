package guard

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/suite"
	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
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

// TestComputeBinaryEntryLargeFile makes sure the streamed hash (pooled
// 128 KiB buffer, multiple reads) matches a one-shot hash for a file well
// past the buffer size.
func (s *guardUnitTest) TestComputeBinaryEntryLargeFile() {
	dir := s.T().TempDir()
	binaryPath := filepath.Join(dir, "big")
	content := make([]byte, 700*1024+123)
	for i := range content {
		content[i] = byte(i*7 + 3)
	}
	s.Require().NoError(os.WriteFile(binaryPath, content, 0644))

	entry, err := ComputeBinaryEntry(binaryPath)
	s.Require().NoError(err)
	s.Require().Equal(sha256.Sum256(content), entry.Hash)
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
	s.Require().Equal("readonly", modeString(ModeReadOnly))
}

func (s *guardUnitTest) TestGuardModeKey() {
	for _, m := range []Mode{ModeBlacklist, ModeWhitelist, ModeReadOnly} {
		k, err := guardModeKey(m)
		s.Require().NoError(err)
		s.Require().Equal(uint64(m), k)
	}
	_, err := guardModeKey(Mode(99))
	s.Require().Error(err)
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
	// READ bit set, and OPEN + STAT implicitly allowed (a reader stat()s and
	// open()s the file before reading it).
	s.Require().NotEqual(uint32(0), mask&(1<<uint(ebpf.EventRead)))
	s.Require().NotEqual(uint32(0), mask&(1<<uint(ebpf.EventOpen)))
	s.Require().NotEqual(uint32(0), mask&(1<<uint(ebpf.EventStat)))
}

func (s *guardUnitTest) TestEventMaskWriteMmap() {
	mask, err := eventMask([]ebpf.EventType{ebpf.EventWrite, ebpf.EventMmap})
	s.Require().NoError(err)
	s.Require().NotEqual(uint32(0), mask&(1<<uint(ebpf.EventWrite)))
	s.Require().NotEqual(uint32(0), mask&(1<<uint(ebpf.EventMmap)))
	s.Require().NotEqual(uint32(0), mask&(1<<uint(ebpf.EventOpen)))
	s.Require().NotEqual(uint32(0), mask&(1<<uint(ebpf.EventStat)))
	// Independent events stay unset.
	s.Require().Equal(uint32(0), mask&(1<<uint(ebpf.EventDelete)))
	s.Require().Equal(uint32(0), mask&(1<<uint(ebpf.EventMkdir)))
}

func (s *guardUnitTest) TestEventMaskIndependent() {
	mask, err := eventMask([]ebpf.EventType{ebpf.EventDelete, ebpf.EventRename})
	s.Require().NoError(err)
	s.Require().NotEqual(uint32(0), mask&(1<<uint(ebpf.EventDelete)))
	s.Require().NotEqual(uint32(0), mask&(1<<uint(ebpf.EventRename)))
	// No read/write/mmap → OPEN and STAT must NOT be implied.
	s.Require().Equal(uint32(0), mask&(1<<uint(ebpf.EventOpen)))
	s.Require().Equal(uint32(0), mask&(1<<uint(ebpf.EventStat)))
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

// TestIsGuardedBinaryConfiguredSymlink is the gh/gcloud/goland regression:
// the config whitelists a symlink (e.g. /home/angelo/.local/bin/gh) while
// /proc/PID/exe reports the resolved real path. Guard creation records the
// canonical target in canonicalPaths, so the real path must match.
func (s *guardUnitTest) TestIsGuardedBinaryConfiguredSymlink() {
	tmpDir := s.T().TempDir()
	targetPath := filepath.Join(tmpDir, "opt", "gh")
	s.Require().NoError(os.MkdirAll(filepath.Dir(targetPath), 0o755))
	s.Require().NoError(os.WriteFile(targetPath, []byte("data"), 0755))

	linkPath := filepath.Join(tmpDir, "bin", "gh")
	s.Require().NoError(os.MkdirAll(filepath.Dir(linkPath), 0o755))
	s.Require().NoError(os.Symlink(targetPath, linkPath))

	g := &Guard{
		binaries: []BinaryEntry{
			{Path: linkPath, Comm: "gh"},
		},
		canonicalPaths: map[string]string{linkPath: targetPath},
	}

	s.Require().True(g.isGuardedBinary(targetPath),
		"resolved exe path must match a symlinked whitelist entry")
	s.Require().False(g.isGuardedBinary(filepath.Join(tmpDir, "opt", "other")),
		"unrelated resolved path must not match")
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

func (s *guardUnitTest) TestCanonicalBinaryPath() {
	dir := s.T().TempDir()
	real := filepath.Join(dir, "real")
	s.Require().NoError(os.WriteFile(real, []byte("x"), 0o644))
	link := filepath.Join(dir, "link")
	s.Require().NoError(os.Symlink(real, link))

	// a symlink resolves to its target for whitelist bookkeeping
	s.Require().Equal(real, canonicalBinaryPath(link))
	// a broken symlink falls back to the literal path (fail-closed)
	broken := filepath.Join(dir, "missing")
	s.Require().Equal(broken, canonicalBinaryPath(broken))
	// a plain path is returned unchanged
	s.Require().Equal(real, canonicalBinaryPath(real))
}

func (s *guardUnitTest) TestResolveDeferred() {
	dir := s.T().TempDir()
	binaryPath := filepath.Join(dir, "tool")
	s.Require().NoError(os.WriteFile(binaryPath, []byte("#!/bin/sh\nexit 0"), 0o755))

	g := &Guard{
		path: dir,
		deferred: []deferredBinary{
			{rule: daemonconfig.BinaryRule{Path: binaryPath, Events: []ebpf.EventType{ebpf.EventRead}}},
			{rule: daemonconfig.BinaryRule{Path: filepath.Join(dir, "still-locked"), Events: []ebpf.EventType{ebpf.EventRead}}},
		},
	}

	resolved, events, stillDeferred := g.resolveDeferred()
	s.Require().Len(resolved, 1, "only the readable rule resolves")
	s.Require().Equal(binaryPath, resolved[0].Path)
	s.Require().NotEqual([32]byte{}, resolved[0].Hash, "resolved entry carries the binary hash")
	s.Require().Equal(filepath.Base(binaryPath), resolved[0].Comm, "comm derives from the binary basename")
	s.Require().Equal([]ebpf.EventType{ebpf.EventRead}, events[resolved[0].Path], "events are keyed by the canonical path")
	s.Require().Len(events, 1, "only resolved rules contribute events")
	s.Require().Len(stillDeferred, 1, "the unreadable rule stays deferred")
	s.Require().Equal(filepath.Join(dir, "still-locked"), stillDeferred[0].rule.Path)
	s.Require().Equal(1, stillDeferred[0].attempts, "attempt count is incremented")
}

// ----------------------------------------------------------------------
// BPF-dependent Integration Tests
// These tests execute real BPF logic and require root privileges.
// They are skipped if not running as root.
// ----------------------------------------------------------------------

// helperChild is the entrypoint taken when the copied test binary is
// exec'd by a guard test: it reads its seed file, and the guard's BPF
// whitelist decides whether the open succeeds. Exit 0 means the access
// was allowed, 1 that it was denied (the binaries are copies of the test
// binary itself, which is static, so they run inside the container).
func helperChild(seedPath string) int {
	data, err := os.ReadFile(seedPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "helper: reading %s: %v\n", seedPath, err)
		return 1
	}
	if len(data) == 0 {
		fmt.Fprintln(os.Stderr, "helper: empty seed")
		return 1
	}
	return 0
}

// copySelf copies the running test binary (a static ELF) to dst.
func copySelf(dst string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(exe)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o755)
}

// runTool execs the copied binary against seed; an error means the
// guard denied its open of the seed file.
func runTool(tool, seed string) error {
	return exec.Command(tool, "-helper-child", seed).Run()
}

// guardedTree creates a whitelist-mode guard over root with the given
// whitelisted binaries, eagerly populated, and started.
func (s *guardUnitTest) newGuardedTree(root string, binaries []BinaryEntry, pending []daemonconfig.BinaryRule) *Guard {
	opts := []GuardOption{WithEagerPopulate()}
	if len(pending) > 0 {
		opts = append(opts, WithPendingBinaries(pending))
	}
	g, err := NewGuard(root, ModeWhitelist, binaries, true, 0, opts...)
	s.Require().NoError(err, "building whitelist guard")
	s.Require().NoError(g.Start(), "starting whitelist guard")
	return g
}

// TestResolvePendingBinariesWhitelist exercises the deferred-whitelist
// flow that fixed the Discord startup bug: a binary unreadable at
// guard-build time stays out of the whitelist (denied in whitelist mode)
// until ResolvePendingBinaries runs after the "unlock", and is then
// allowed. Fail-closed before, allowed after.
func (s *guardUnitTest) TestResolvePendingBinariesWhitelist() {
	if os.Getuid() != 0 {
		s.T().Skip("Skipping BPF test: requires root")
	}

	root := s.T().TempDir()
	seed := filepath.Join(root, "seed.txt")
	s.Require().NoError(os.WriteFile(seed, []byte("data"), 0o644))

	tool := filepath.Join("/tmp", fmt.Sprintf("guard-tool-%d", os.Getpid()))
	s.Require().NoError(copySelf(tool))
	defer os.Remove(tool)

	g := s.newGuardedTree(root, nil, []daemonconfig.BinaryRule{{Path: tool}})
	defer g.Stop()

	// Before resolve: not whitelisted → the open of seed is denied.
	s.Require().Error(runTool(tool, seed), "unresolved pending binary must be denied")

	s.Require().NoError(g.ResolvePendingBinaries())

	s.Require().NoError(runTool(tool, seed), "resolved binary must be allowed")
}

// TestPopulateInodesFillsMap verifies the eager inode scan registers
// every file and directory of the guarded tree in guard_inodes, so
// kernel-side own/parent inode checks protect deep files.
func (s *guardUnitTest) TestPopulateInodesFillsMap() {
	if os.Getuid() != 0 {
		s.T().Skip("Skipping BPF test: requires root")
	}

	root := s.T().TempDir()
	s.Require().NoError(os.MkdirAll(filepath.Join(root, "sub", "deep"), 0o755))
	dirs := []string{
		root,
		filepath.Join(root, "sub"),
		filepath.Join(root, "sub", "deep"),
	}
	files := []string{
		filepath.Join(root, "a.txt"),
		filepath.Join(root, "sub", "b.txt"),
		filepath.Join(root, "sub", "deep", "c.txt"),
	}
	for _, p := range files {
		s.Require().NoError(os.WriteFile(p, []byte("x"), 0o644))
	}
	nodes := append(append([]string{}, dirs...), files...)

	// Stat every node BEFORE the guard attaches: once it is live this
	// non-whitelisted test process is denied stat(2) on a guarded path
	// (inode_getattr enforcement).
	keys := make(map[string]GuardInodeKey, len(nodes))
	for _, p := range nodes {
		dev, ino, err := ebpf.StatInode(p)
		s.Require().NoError(err, "stating %s", p)
		keys[p] = GuardInodeKey{Dev: dev, Ino: ino}
	}

	g := s.newGuardedTree(root, nil, nil)
	defer g.Stop()

	for _, p := range nodes {
		var v uint32
		s.Require().Truef(g.objs().GuardInodes.Lookup(keys[p], &v) == nil,
			"inode of %s missing from guard_inodes", p)
	}

	count := 0
	it := g.objs().GuardInodes.Iterate()
	var k GuardInodeKey
	var v uint32
	for it.Next(&k, &v) {
		// guard_inodes is shared by every resource; count only this guard's rows.
		if v == g.resID {
			count++
		}
	}
	s.Require().Equal(len(nodes), count, "guard_inodes must contain exactly the tree nodes")
}

// TestSweepInodesRecreatedFileRoot verifies the cheap periodic sweep re-maps
// a single-file watch root that its application deleted and recreated (a new
// inode that nothing else would add to guard_inodes).
func (s *guardUnitTest) TestSweepInodesRecreatedFileRoot() {
	if os.Getuid() != 0 {
		s.T().Skip("Skipping BPF test: requires root")
	}

	dir := s.T().TempDir()
	fileRoot := filepath.Join(dir, "state.db")
	s.Require().NoError(os.WriteFile(fileRoot, []byte("v1"), 0o644))

	// The guard denies unlink of a guarded file by a non-whitelisted process,
	// so whitelist this test binary — mirroring the real case, where the app
	// that deletes and recreates its own state file is itself whitelisted.
	exe, err := os.Executable()
	s.Require().NoError(err)
	self, err := ComputeBinaryEntry(exe)
	s.Require().NoError(err)
	g := s.newGuardedTree(fileRoot, []BinaryEntry{self}, nil)
	defer g.Stop()

	inMap := func(path string) bool {
		dev, ino, err := ebpf.StatInode(path)
		s.Require().NoError(err)
		var v uint32
		return g.objs().GuardInodes.Lookup(GuardInodeKey{Dev: dev, Ino: ino}, &v) == nil
	}
	s.Require().True(inMap(fileRoot), "the file root must be mapped at build")

	_, oldIno, _ := ebpf.StatInode(fileRoot)
	s.Require().NoError(os.Remove(fileRoot))
	replacement := filepath.Join(dir, "state.db.new")
	s.Require().NoError(os.WriteFile(replacement, []byte("v2-recreated"), 0o644))
	s.Require().NoError(os.Rename(replacement, fileRoot))
	_, newIno, _ := ebpf.StatInode(fileRoot)
	if oldIno == newIno {
		// The container's fs reused the freed inode number — the guarded
		// entry still matches, so SweepInodes has nothing to fix and the
		// scenario this test targets did not occur.
		s.T().Skip("filesystem reused the inode number on recreate")
	}

	s.Require().False(inMap(fileRoot), "the recreated inode is not mapped yet")
	s.Require().NoError(g.SweepInodes())
	s.Require().True(inMap(fileRoot), "SweepInodes must map the recreated file root")

	// Regression for the false-DENY-on-unrelated-files / silently-unguarded-real-file bug:
	// SweepInodes must move BOTH the kernel root-confinement anchor (the resource slot,
	// root_in_chain) and g.rootKey to the new inode and evict the old one from guard_inodes.
	// Otherwise the old freed number stays "protected" forever (ReconcileInodes won't evict
	// g.rootKey) and, once the filesystem gives it to an unrelated file elsewhere, that file is
	// denied under this whitelist by coincidence while the real recreated file silently loses
	// protection (its chain no longer contains the stale root).
	cfg, cfgErr := g.resConfig()
	s.Require().NoError(cfgErr)
	gotDev, gotIno := cfg.RootDev, cfg.RootIno
	s.Require().Equal(newIno, gotIno, "the resource slot's root ino must follow the recreated file")

	g.mu.Lock()
	gotRootKey := g.rootKey
	g.mu.Unlock()
	s.Require().Equal(GuardInodeKey{Dev: gotDev, Ino: gotIno}, gotRootKey, "g.rootKey must match the resource slot")

	var v uint32
	oldKey := GuardInodeKey{Dev: gotDev, Ino: oldIno}
	s.Require().Error(g.objs().GuardInodes.Lookup(oldKey, &v), "the stale old-root inode must be evicted from guard_inodes")
}

// Directory-root counterpart of TestSweepInodesRecreatedFileRoot (same bug class: stale anchor ->
// false DENY + silently unguarded resource) for a DIRECTORY root deleted and recreated wholesale
// (in-place fscrypt migration, backup restore, app rebuilding its config dir). Before the fix only
// the single-file branch re-anchored; the directory branch tracked mtime and never re-checked its
// own inode.
func (s *guardUnitTest) TestSweepInodesRecreatedDirRoot() {
	if os.Getuid() != 0 {
		s.T().Skip("Skipping BPF test: requires root")
	}

	parent := s.T().TempDir()
	watchDir := filepath.Join(parent, "profile")
	s.Require().NoError(os.Mkdir(watchDir, 0o755))
	s.Require().NoError(os.WriteFile(filepath.Join(watchDir, "a.txt"), []byte("v1"), 0o644))

	// Deleting/recreating the guarded directory requires whitelisting this
	// test process, mirroring TestSweepInodesRecreatedFileRoot.
	exe, err := os.Executable()
	s.Require().NoError(err)
	self, err := ComputeBinaryEntry(exe)
	s.Require().NoError(err)
	g := s.newGuardedTree(watchDir, []BinaryEntry{self}, nil)
	defer g.Stop()

	inMap := func(path string) bool {
		dev, ino, err := ebpf.StatInode(path)
		s.Require().NoError(err)
		var v uint32
		return g.objs().GuardInodes.Lookup(GuardInodeKey{Dev: dev, Ino: ino}, &v) == nil
	}
	s.Require().True(inMap(watchDir), "the directory root must be mapped at build")

	_, oldIno, _ := ebpf.StatInode(watchDir)
	s.Require().NoError(os.RemoveAll(watchDir))
	s.Require().NoError(os.Mkdir(watchDir, 0o755))
	newFile := filepath.Join(watchDir, "b.txt")
	s.Require().NoError(os.WriteFile(newFile, []byte("v2-recreated"), 0o644))
	_, newIno, _ := ebpf.StatInode(watchDir)
	if oldIno == newIno {
		// The filesystem reused the freed inode number — the guarded entry
		// still matches, so SweepInodes has nothing to fix and the scenario
		// this test targets did not occur.
		s.T().Skip("filesystem reused the inode number on recreate")
	}

	s.Require().False(inMap(watchDir), "the recreated directory root is not mapped yet")
	s.Require().NoError(g.SweepInodes())
	s.Require().True(inMap(watchDir), "SweepInodes must map the recreated directory root")
	s.Require().True(inMap(newFile), "SweepInodes must map the recreated directory's new content")

	// Same coverage as TestSweepInodesRecreatedFileRoot: the root-confinement anchor
	// (the resource slot) and g.rootKey must move to the new inode and the old one be evicted from
	// guard_inodes, else it stays "protected" forever and, once reused by an unrelated
	// file/directory elsewhere, that gets denied by inode-number coincidence while the real
	// recreated directory silently loses protection.
	cfg, cfgErr := g.resConfig()
	s.Require().NoError(cfgErr)
	gotDev, gotIno := cfg.RootDev, cfg.RootIno
	s.Require().Equal(newIno, gotIno, "the resource slot's root ino must follow the recreated directory")

	g.mu.Lock()
	gotRootKey := g.rootKey
	g.mu.Unlock()
	s.Require().Equal(GuardInodeKey{Dev: gotDev, Ino: gotIno}, gotRootKey, "g.rootKey must match the resource slot")

	var v uint32
	oldKey := GuardInodeKey{Dev: gotDev, Ino: oldIno}
	s.Require().Error(g.objs().GuardInodes.Lookup(oldKey, &v), "the stale old-root inode must be evicted from guard_inodes")
}

// The guard_path_rmdir eviction fix is tested at guard level
// (TestGuard_PathRmdirEvictsInodeImmediately in integrationtests/guard_test.go): it queries the
// live guard_inodes map via bpftool as root in the privileged container, since a directory has no
// hard-link equivalent to prove eviction black-box (Linux refuses to hard-link a directory).

// TestSweepInodesDirRootGated verifies a directory root whose mtime has not
// moved is not re-walked (the expensive path the sweep avoids).
func (s *guardUnitTest) TestSweepInodesDirRootGated() {
	if os.Getuid() != 0 {
		s.T().Skip("Skipping BPF test: requires root")
	}

	root := s.T().TempDir()
	s.Require().NoError(os.WriteFile(filepath.Join(root, "a"), []byte("x"), 0o644))

	g := s.newGuardedTree(root, nil, nil)
	defer g.Stop()

	// First sweep records the fingerprint (mtime unchanged since build, so it
	// is a no-op) and every subsequent sweep with an unchanged root is a
	// no-op — a file created in a SUBdir must not be picked up here (BPF
	// discovery + the ancestor walk cover it), only a top-level change does.
	s.Require().NoError(g.SweepInodes())
	g.mu.Lock()
	seededMtime := g.sweepRootMtime
	g.mu.Unlock()

	s.Require().NoError(g.SweepInodes())
	g.mu.Lock()
	s.Require().Equal(seededMtime, g.sweepRootMtime, "an unchanged root must not be re-fingerprinted")
	g.mu.Unlock()
}

// The one behavior distinguishing walkLiveEntries from walkInodes: unlike PopulateInodes' tolerant
// root-vanish handling (safe there: the ancestor walk and fail-closed defenses cover
// under-collection), ReconcileInodes must never read "root briefly unreadable" as "tree empty",
// which would license evicting every entry, including the root's own protection.
func (s *guardUnitTest) TestWalkLiveEntriesRootFailureIsHardError() {
	err := walkLiveEntries("/some/path", false, 0, func(p string) error {
		if p == "/some/path" {
			return unix.ENOENT
		}
		return nil
	})
	s.Require().Error(err, "a root that fails to stat must abort the walk, not report an empty tree")
}

// TestWalkLiveEntriesToleratesVanishingChild verifies a child vanishing
// mid-walk is normal churn, not grounds to abort: it correctly means that
// entry is no longer live, so ReconcileInodes may evict its stale key.
func (s *guardUnitTest) TestWalkLiveEntriesToleratesVanishingChild() {
	dir := s.T().TempDir()
	s.Require().NoError(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0o644))
	s.Require().NoError(os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0o644))

	seen := 0
	err := walkLiveEntries(dir, false, 0, func(p string) error {
		seen++
		if filepath.Base(p) == "a.txt" {
			return unix.ENOENT
		}
		return nil
	})
	s.Require().NoError(err, "a child vanishing mid-walk must not abort the walk")
	s.Require().Equal(3, seen, "root + both children must be visited")
}

// The periodic inode GC removes a guard_inodes entry once its file is gone, leaving the watch
// root's entry untouched. Nothing else deletes from guard_inodes (PopulateInodes/SweepInodes/BPF
// mkdir-rename discovery only add), so a long-lived guard accumulates a stale (dev, ino) per
// deleted file, and on filesystems that reuse inode numbers one can collide with an unrelated file
// and false-DENY.
func (s *guardUnitTest) TestReconcileInodesEvictsStale() {
	if os.Getuid() != 0 {
		s.T().Skip("Skipping BPF test: requires root")
	}

	root := s.T().TempDir()
	stale := filepath.Join(root, "stale.txt")
	s.Require().NoError(os.WriteFile(stale, []byte("gone soon"), 0o644))

	// Deleting inside the guarded tree requires whitelisting this test
	// process, mirroring TestSweepInodesRecreatedFileRoot.
	exe, err := os.Executable()
	s.Require().NoError(err)
	self, err := ComputeBinaryEntry(exe)
	s.Require().NoError(err)
	g := s.newGuardedTree(root, []BinaryEntry{self}, nil)
	defer g.Stop()

	dev, ino, err := ebpf.StatInode(stale)
	s.Require().NoError(err)
	staleKey := GuardInodeKey{Dev: dev, Ino: ino}

	var v uint32
	s.Require().NoError(g.objs().GuardInodes.Lookup(staleKey, &v), "the file must be mapped at build")

	s.Require().NoError(os.Remove(stale))
	s.Require().NoError(g.ReconcileInodes())

	s.Require().Error(g.objs().GuardInodes.Lookup(staleKey, &v),
		"the deleted file's stale inode entry must be evicted")

	rootDev, rootIno, err := ebpf.StatInode(root)
	s.Require().NoError(err)
	s.Require().NoError(g.objs().GuardInodes.Lookup(GuardInodeKey{Dev: rootDev, Ino: rootIno}, &v),
		"ReconcileInodes must never evict the watch root's own entry")
}

// A whitelisted binary replaced in place (same path, new inode) is denied until ReSyncBinaries
// re-admits it, and re-admission needs the replacement check's approval: by path alone, anyone able
// to swap the path would inherit the whitelist entry.
func (s *guardUnitTest) TestReSyncBinariesReplacement() {
	if os.Getuid() != 0 {
		s.T().Skip("Skipping BPF test: requires root")
	}

	root := s.T().TempDir()
	seed := filepath.Join(root, "seed.txt")
	s.Require().NoError(os.WriteFile(seed, []byte("data"), 0o644))

	tool := filepath.Join("/tmp", fmt.Sprintf("guard-tool-%d", os.Getpid()))
	s.Require().NoError(copySelf(tool))
	defer os.Remove(tool)

	entry, err := ComputeBinaryEntry(tool)
	s.Require().NoError(err)
	oldDev, oldIno, err := ebpf.StatInode(tool)
	s.Require().NoError(err)
	oldKey := GuardInodeKey{Dev: oldDev, Ino: oldIno}

	g := s.newGuardedTree(root, []BinaryEntry{entry}, nil)
	defer g.Stop()

	s.Require().NoError(runTool(tool, seed), "original binary must be allowed")

	replacement := filepath.Join("/tmp", fmt.Sprintf("guard-tool-new-%d", os.Getpid()))
	s.Require().NoError(copySelf(replacement))
	defer os.Remove(replacement)
	s.Require().NoError(os.Rename(replacement, tool))
	newDev, newIno, err := ebpf.StatInode(tool)
	s.Require().NoError(err)
	newKey := GuardInodeKey{Dev: newDev, Ino: newIno}

	s.Require().Error(runTool(tool, seed), "replaced binary must be denied until re-synced")

	// No check installed (no trust guard, so no provenance): refused, still denied in-kernel.
	SetReplacementCheck(nil)
	changed, err := g.ReSyncBinaries()
	s.Require().NoError(err)
	s.Require().Zero(changed, "an unapproved replacement must not be re-admitted")
	s.Require().Error(runTool(tool, seed), "an unapproved replacement must stay denied")

	// Approved for exactly this old->new swap: re-admitted.
	SetReplacementCheck(func(path string, o, n GuardInodeKey) bool {
		return path == tool && o == oldKey && n == newKey
	})
	defer SetReplacementCheck(nil)
	changed, err = g.ReSyncBinaries()
	s.Require().NoError(err)
	s.Require().Equal(1, changed, "the approved replacement must be reported as changed")
	s.Require().NoError(runTool(tool, seed), "re-synced binary must be allowed")

	again, err := g.ReSyncBinaries()
	s.Require().NoError(err)
	s.Require().Zero(again, "second re-sync must be a no-op")
}

// newSelfAllowedGuard builds and starts a guard that, like the daemon's, allows the running test
// binary (root-gated): nested guards must still let their owner stat and scan every tree.
func (s *guardUnitTest) newSelfAllowedGuard(root string, mode Mode, binaries []BinaryEntry) *Guard {
	self, err := ComputeBinaryEntry("/proc/self/exe")
	s.Require().NoError(err)
	g, err := NewGuard(root, mode, binaries, true, 0, WithEagerPopulate(), WithSelfAllowBinary(self, nil))
	s.Require().NoError(err, "building guard over %s", root)
	s.Require().NoError(g.Start(), "starting guard over %s", root)
	return g
}

// A sealed file (whitelist, no binaries) inside a read-only tree is the shape of fscrypt.key inside
// /etc/app-listener. A rescan of the outer tree (what its periodic SweepInodes does) must not hand
// the file to the outer resource, whose read-only rule lets everyone read.
func (s *guardUnitTest) TestNestedResourceOuterRescanKeepsInnerSealed() {
	if os.Getuid() != 0 {
		s.T().Skip("Skipping BPF test: requires root")
	}

	outer := s.T().TempDir()
	sealed := filepath.Join(outer, "key")
	s.Require().NoError(os.WriteFile(sealed, []byte("data"), 0o600))

	tool := filepath.Join("/tmp", fmt.Sprintf("guard-nest-tool-%d", os.Getpid()))
	s.Require().NoError(copySelf(tool))
	defer os.Remove(tool)

	og := s.newSelfAllowedGuard(outer, ModeReadOnly, nil)
	defer og.Stop()
	ig := s.newSelfAllowedGuard(sealed, ModeWhitelist, nil)
	defer ig.Stop()

	s.Require().Error(runTool(tool, sealed), "the sealed inner file must be denied")

	s.Require().NoError(og.PopulateInodes())
	s.Require().Error(runTool(tool, sealed), "an outer rescan must not unseal the inner file")
}

// guard_inodes holds one owner per inode: with nested resources (only the daemon's self guards;
// configs refuse nesting) the innermost must own a shared file whichever resource scanned last.
func (s *guardUnitTest) TestNestedResourcesInnermostOwns() {
	if os.Getuid() != 0 {
		s.T().Skip("Skipping BPF test: requires root")
	}

	outer := s.T().TempDir()
	inner := filepath.Join(outer, "inner")
	s.Require().NoError(os.Mkdir(inner, 0o755))
	secret := filepath.Join(inner, "secret")
	s.Require().NoError(os.WriteFile(secret, []byte("data"), 0o644))

	tools := map[string]string{}
	entries := map[string]BinaryEntry{}
	for _, name := range []string{"outer-only", "inner-only", "both"} {
		p := filepath.Join("/tmp", fmt.Sprintf("guard-nest-%s-%d", name, os.Getpid()))
		s.Require().NoError(copySelf(p))
		defer os.Remove(p)
		e, err := ComputeBinaryEntry(p)
		s.Require().NoError(err)
		tools[name], entries[name] = p, e
	}

	og := s.newSelfAllowedGuard(outer, ModeWhitelist, []BinaryEntry{entries["outer-only"], entries["both"]})
	defer og.Stop()
	ig := s.newSelfAllowedGuard(inner, ModeWhitelist, []BinaryEntry{entries["inner-only"], entries["both"]})
	defer ig.Stop()

	check := func(when string) {
		s.Require().NoErrorf(runTool(tools["both"], secret), "%s: a binary both resources allow must read", when)
		s.Require().NoErrorf(runTool(tools["inner-only"], secret), "%s: the inner resource must own the file", when)
		s.Require().Errorf(runTool(tools["outer-only"], secret), "%s: the outer whitelist must not apply inside the inner resource", when)
	}
	check("inner scanned last")
	s.Require().NoError(og.PopulateInodes())
	check("outer scanned last")
}

// Finding #4: SetTrusted must not empty guard_trusted_files before refilling it. A reload is the
// normal way the trusted set changes (SIGHUP catalog refresh), and while the old code cleared the
// map then re-populated it, a concurrent process was momentarily trusted for nothing — the library
// allowlist and binary write-protection enforced nothing in that window. syncMap (put-then-delete)
// keeps a binary that survives the change continuously present. A reader in a tight loop would see
// the pre-fix clear-then-fill window; post-fix it never does.
func (s *guardUnitTest) TestSetTrustedNeverDropsPersistentEntry() {
	if os.Getuid() != 0 {
		s.T().Skip("Skipping BPF test: requires root")
	}

	dir := s.T().TempDir()
	persistent := filepath.Join(dir, "persistent-bin")
	s.Require().NoError(os.WriteFile(persistent, []byte("x"), 0o755))
	dev, ino, err := ebpf.StatInode(persistent)
	s.Require().NoError(err)
	pkey := GuardInodeKey{Dev: dev, Ino: ino}

	// A pool of binaries that come and go across re-applications, so each SetTrusted is a real
	// change (adds one, drops the previous) — the persistent one is in every set.
	others := make([]string, 8)
	for i := range others {
		p := filepath.Join(dir, fmt.Sprintf("other-%d", i))
		s.Require().NoError(os.WriteFile(p, []byte("y"), 0o755))
		others[i] = p
	}

	tg, err := NewTrustGuard()
	s.Require().NoError(err)
	defer tg.Stop()
	s.Require().NoError(tg.SetTrusted([]string{persistent}, nil))

	// The persistent binary must be trusted at ALL times while SetTrusted re-applies concurrently.
	var dropped atomic.Bool
	stop := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		var v uint8
		for {
			select {
			case <-stop:
				return
			default:
			}
			if err := tg.objs.GuardTrustedFiles.Lookup(pkey, &v); err != nil {
				dropped.Store(true)
				return
			}
		}
	}()

	for i := 0; i < 300; i++ {
		s.Require().NoError(tg.SetTrusted([]string{persistent, others[i%len(others)]}, nil))
	}
	close(stop)
	wg.Wait()

	s.Require().Falsef(dropped.Load(),
		"the persistent binary vanished from guard_trusted_files during a re-apply — "+
			"SetTrusted must sync (put-then-delete), never clear before refilling")
}

// TestMain intercepts the -helper-child invocation: the copied test
// binary must not run the test suite when the guard execs it; every
// other invocation proceeds to the normal test suite.
func TestMain(m *testing.M) {
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "-helper-child" && i+1 < len(os.Args) {
			os.Exit(helperChild(os.Args[i+1]))
		}
	}
	os.Exit(m.Run())
}

// A stopped resource's whitelist rows must leave the shared maps: resource ids are reused, so a
// surviving row would let the next resource in that slot inherit a binary it never allowed.
func (s *guardUnitTest) TestStoppedResourceWhitelistNotInheritedBySlotReuse() {
	if os.Getuid() != 0 {
		s.T().Skip("Skipping BPF test: requires root")
	}

	tool := filepath.Join("/tmp", fmt.Sprintf("guard-slot-tool-%d", os.Getpid()))
	s.Require().NoError(copySelf(tool))
	defer os.Remove(tool)
	entry, err := ComputeBinaryEntry(tool)
	s.Require().NoError(err)
	dev, ino, err := ebpf.StatInode(tool)
	s.Require().NoError(err)
	exe := GuardInodeKey{Dev: dev, Ino: ino}

	// Keeps the shared engine (and its maps) alive across the other guard's stop.
	holder := s.newGuardedTree(s.T().TempDir(), nil, nil)
	defer holder.Stop()

	allowing := s.newGuardedTree(s.T().TempDir(), []BinaryEntry{entry}, nil)
	slot := allowing.resID
	var action uint8
	s.Require().NoError(allowing.objs().GuardExeActions.Lookup(allowing.resKey(exe), &action),
		"the whitelisted binary must be registered for its resource")
	allowing.Stop()

	s.Require().Error(holder.objs().GuardExeActions.Lookup(GuardResInodeKey{ResId: slot, Ino: exe}, &action),
		"a stopped resource's whitelist row must be removed from the shared map")

	reuser := s.newGuardedTree(s.T().TempDir(), nil, nil)
	defer reuser.Stop()
	s.Require().Equal(slot, reuser.resID, "fixture: the freed slot id is reused")
	s.Require().Error(reuser.objs().GuardExeActions.Lookup(reuser.resKey(exe), &action),
		"the resource reusing the slot must not inherit the previous resource's whitelist")
}

// Findings #2/#3: a replacement guard built over the SAME root as a still-live guard (what a reload
// does before it decides whether to keep the new config) takes over the live guard's guard_inodes
// row, because claimInode does not treat a same-root replacement as inner. When that replacement is
// then stopped — a refused reload, a failed buildGuards, a rollback — Stop deletes every row it
// owns, including the one it took, so the kept resource loses its only inode row. SweepInodes won't
// re-add it (the root inode never changed) and ReconcileInodes only deletes, so the resource is
// left readable indefinitely though the old guard is still attached and "keeping" it.
func (s *guardUnitTest) TestReloadRollbackKeepsFileRootProtected() {
	if os.Getuid() != 0 {
		s.T().Skip("Skipping BPF test: requires root")
	}

	tool := filepath.Join("/tmp", fmt.Sprintf("guard-rollback-file-%d", os.Getpid()))
	s.Require().NoError(copySelf(tool))
	defer os.Remove(tool)

	root := filepath.Join(s.T().TempDir(), "secret")
	s.Require().NoError(os.WriteFile(root, []byte("data"), 0o644))

	// Both guards self-allow the test process, as the daemon's guards self-allow the daemon exe:
	// otherwise the live guard would deny this process's stat/scan while it builds the replacement,
	// which is not the behaviour under test. The separate `tool` binary (a different inode) stays
	// non-whitelisted and is the real probe.
	kept := s.newSelfAllowedGuard(root, ModeWhitelist, nil)
	defer kept.Stop()
	s.Require().Error(runTool(tool, root), "the live guard must deny a non-whitelisted reader")

	// The reload's replacement guard for the same resource, built while the old one still runs.
	repl := s.newSelfAllowedGuard(root, ModeWhitelist, nil)
	s.Require().Error(runTool(tool, root), "with the replacement active the root stays denied")

	// Reload refused / rolled back: only the replacement is stopped; kept stays attached.
	repl.Stop()
	s.Require().Error(runTool(tool, root),
		"stopping the rolled-back replacement left the kept single-file resource readable — its inode row was deleted")
}

// Directory counterpart of TestReloadRollbackKeepsFileRootProtected: the replacement's eager scan
// claims every row of the kept tree, so stopping it strips the whole tree, not just the root.
func (s *guardUnitTest) TestReloadRollbackKeepsDirTreeProtected() {
	if os.Getuid() != 0 {
		s.T().Skip("Skipping BPF test: requires root")
	}

	tool := filepath.Join("/tmp", fmt.Sprintf("guard-rollback-dir-%d", os.Getpid()))
	s.Require().NoError(copySelf(tool))
	defer os.Remove(tool)

	root := s.T().TempDir()
	inner := filepath.Join(root, "sub", "secret")
	s.Require().NoError(os.MkdirAll(filepath.Dir(inner), 0o755))
	s.Require().NoError(os.WriteFile(inner, []byte("data"), 0o644))

	// Self-allowed for the same reason as the file-root case: the probe is the separate `tool`.
	kept := s.newSelfAllowedGuard(root, ModeWhitelist, nil)
	defer kept.Stop()
	s.Require().Error(runTool(tool, inner), "the live guard must deny a non-whitelisted reader of a deep file")

	repl := s.newSelfAllowedGuard(root, ModeWhitelist, nil)
	s.Require().Error(runTool(tool, inner), "with the replacement active the deep file stays denied")

	repl.Stop()
	s.Require().Error(runTool(tool, inner),
		"stopping the rolled-back replacement left the kept directory tree readable — its inode rows were deleted")
}
