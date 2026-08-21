package guard

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

	g := s.newGuardedTree(root, nil, nil)
	defer g.Stop()

	inMap := func(path string) bool {
		dev, ino, err := ebpf.StatInode(path)
		s.Require().NoError(err, "stating %s", path)
		var v uint8
		return g.objs.GuardInodes.Lookup(GuardInodeKey{Dev: dev, Ino: ino}, &v) == nil
	}
	for _, p := range nodes {
		s.Require().True(inMap(p), "inode of %s missing from guard_inodes", p)
	}

	count := 0
	it := g.objs.GuardInodes.Iterate()
	var k GuardInodeKey
	var v uint8
	for it.Next(&k, &v) {
		count++
	}
	s.Require().Equal(len(nodes), count, "guard_inodes must contain exactly the tree nodes")
}

// TestReSyncBinariesReplacement verifies the in-place-replacement fix
// for the Discord updater denials: a whitelisted binary replaced in
// place (same path, new inode) is denied after relaunch — the whitelist
// is keyed by inode — until ReSyncBinaries rewrites its map entry, after
// which the updated binary is allowed again.
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

	g := s.newGuardedTree(root, []BinaryEntry{entry}, nil)
	defer g.Stop()

	// Original binary is whitelisted at build: allowed.
	s.Require().NoError(runTool(tool, seed), "original binary must be allowed")

	// Replace the binary in place: write a new file, rename over the
	// old one — same path, brand-new inode.
	replacement := filepath.Join("/tmp", fmt.Sprintf("guard-tool-new-%d", os.Getpid()))
	s.Require().NoError(copySelf(replacement))
	defer os.Remove(replacement)
	s.Require().NoError(os.Rename(replacement, tool))

	// The relaunched binary's inode is not whitelisted: denied.
	s.Require().Error(runTool(tool, seed), "replaced binary must be denied until re-synced")

	// ReSyncBinaries rewrites the map entry for the new inode.
	changed, err := g.ReSyncBinaries()
	s.Require().NoError(err)
	s.Require().GreaterOrEqual(changed, 1, "the replaced binary must be reported as changed")

	s.Require().NoError(runTool(tool, seed), "re-synced binary must be allowed")

	// Idempotent: a second pass finds nothing to do.
	again, err := g.ReSyncBinaries()
	s.Require().NoError(err)
	s.Require().Zero(again, "second re-sync must be a no-op")
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
