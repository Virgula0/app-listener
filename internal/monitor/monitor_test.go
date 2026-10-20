package monitor

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/Virgula0/app-listener/internal/infrastructure"
)

type monitorUnitTest struct {
	suite.Suite
}

func TestMonitorUnitTest(t *testing.T) {
	suite.Run(t, new(monitorUnitTest))
}

func dirTarget(dir string) []ebpf.Target {
	return []ebpf.Target{{Dir: dir, IsDir: true}}
}

func (s *monitorUnitTest) TestInWatchPath() {
	m := &Monitor{targets: dirTarget("/tmp/watch")}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{name: "exact match", path: "/tmp/watch", expected: true},
		{name: "child path", path: "/tmp/watch/file.txt", expected: true},
		{name: "nested child", path: "/tmp/watch/sub/dir/file.txt", expected: true},
		{name: "non matching", path: "/var/log/syslog", expected: false},
		{name: "similar prefix no longer matches", path: "/tmp/watch_extra", expected: false},
		{name: "empty path", path: "", expected: false},
		{name: "root path", path: "/", expected: false},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result := m.inWatchPath(tt.path, "")
			s.Require().Equal(tt.expected, result)
		})
	}
}

func (s *monitorUnitTest) TestInWatchPathFileTarget() {
	m := &Monitor{targets: []ebpf.Target{{Dir: "/parent", File: "target.txt", IsDir: false}}}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{name: "exact basename", path: "target.txt", expected: true},
		{name: "full path", path: "/parent/target.txt", expected: true},
		{name: "different file", path: "/parent/other.txt", expected: false},
		{name: "empty", path: "", expected: false},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result := m.inWatchPath(tt.path, "")
			s.Require().Equal(tt.expected, result)
		})
	}
}

func (s *monitorUnitTest) TestInWatchPathDest() {
	m := &Monitor{targets: dirTarget("/watch")}

	tests := []struct {
		name     string
		path     string
		dest     string
		expected bool
	}{
		{name: "path matches", path: "/watch/file", dest: "", expected: true},
		{name: "dest matches", path: "", dest: "/watch/new", expected: true},
		{name: "neither matches", path: "/other", dest: "", expected: false},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result := m.inWatchPath(tt.path, tt.dest)
			s.Require().Equal(tt.expected, result)
		})
	}
}

func (s *monitorUnitTest) TestMatchesFilterNonRecursive() {
	m := &Monitor{targets: dirTarget("/watch"), recursive: false}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{name: "watch dir itself", path: "/watch", expected: true},
		{name: "direct child", path: "/watch/file.txt", expected: true},
		{name: "nested child", path: "/watch/sub/file.txt", expected: false},
		{name: "deeply nested", path: "/watch/a/b/c.txt", expected: false},
		{name: "outside watch", path: "/other/file.txt", expected: false},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ev := &ebpf.FileEvent{Path: tt.path}
			result := m.matchesFilter(ev)
			s.Require().Equal(tt.expected, result)
		})
	}
}

func (s *monitorUnitTest) TestMatchesFilterIntoTargetByDest() {
	m := &Monitor{targets: dirTarget("/watch"), recursive: true}
	for _, typ := range []ebpf.EventType{ebpf.EventRename, ebpf.EventHardlink} {
		s.Require().True(m.matchesFilter(&ebpf.FileEvent{Type: typ, Path: "/other/f", Dest: "/watch/f"}), typ.String())
		s.Require().False(m.matchesFilter(&ebpf.FileEvent{Type: typ, Path: "/other/f", Dest: "/other/g"}), typ.String())
	}
	// A symlink's Dest is its body: naming a target doesn't touch it.
	s.Require().False(m.matchesFilter(&ebpf.FileEvent{Type: ebpf.EventSymlink, Path: "/other/l", Dest: "/watch/f"}))
}

func (s *monitorUnitTest) TestMatchesFilterRecursiveNoDepth() {
	m := &Monitor{targets: dirTarget("/watch"), recursive: true, depth: 0}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{name: "watch dir itself", path: "/watch", expected: true},
		{name: "direct child", path: "/watch/file.txt", expected: true},
		{name: "nested child", path: "/watch/sub/file.txt", expected: true},
		{name: "deeply nested", path: "/watch/a/b/c/d/e.txt", expected: true},
		{name: "outside watch", path: "/other/file.txt", expected: false},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ev := &ebpf.FileEvent{Path: tt.path}
			result := m.matchesFilter(ev)
			s.Require().Equal(tt.expected, result)
		})
	}
}

func (s *monitorUnitTest) TestMatchesFilterRecursiveWithDepth() {
	m := &Monitor{targets: dirTarget("/watch"), recursive: true, depth: 2}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{name: "watch dir itself", path: "/watch", expected: true},
		{name: "depth 1", path: "/watch/file.txt", expected: true},
		{name: "depth 2", path: "/watch/sub/file.txt", expected: true},
		{name: "depth 3 (exceeds)", path: "/watch/a/b/file.txt", expected: false},
		{name: "depth 4 (exceeds)", path: "/watch/a/b/c/d.txt", expected: false},
		{name: "outside watch", path: "/other/file.txt", expected: false},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ev := &ebpf.FileEvent{Path: tt.path}
			result := m.matchesFilter(ev)
			s.Require().Equal(tt.expected, result)
		})
	}
}

func (s *monitorUnitTest) TestResolveSymlinkToTarget() {
	dir := s.T().TempDir()

	watchedDir := filepath.Join(dir, "watched")
	s.Require().NoError(os.Mkdir(watchedDir, 0755))
	watchedFile := filepath.Join(watchedDir, "secret.txt")
	s.Require().NoError(os.WriteFile(watchedFile, []byte("data"), 0644))

	m := &Monitor{targets: dirTarget(watchedDir)}

	symOutside := filepath.Join(dir, "outside_link")
	s.Require().NoError(os.Symlink(watchedFile, symOutside))

	symInside := filepath.Join(watchedDir, "inside_link")
	s.Require().NoError(os.Symlink(watchedFile, symInside))

	symDir := filepath.Join(dir, "dir_link")
	s.Require().NoError(os.Symlink(watchedDir, symDir))

	otherFile := filepath.Join(dir, "other.txt")
	s.Require().NoError(os.WriteFile(otherFile, []byte("other"), 0644))

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "path already matches",
			path:     watchedFile,
			expected: watchedFile,
		},
		{
			name:     "symlink resolves to watched file",
			path:     symOutside,
			expected: watchedFile,
		},
		{
			name:     "symlink inside watched dir",
			path:     symInside,
			expected: symInside,
		},
		{
			name:     "symlink to watched dir with subpath",
			path:     filepath.Join(symDir, "secret.txt"),
			expected: watchedFile,
		},
		{
			name:     "non-matching regular file",
			path:     otherFile,
			expected: otherFile,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result := m.resolveSymlinkToTarget(tt.path)
			s.Require().Equal(tt.expected, result)
		})
	}
}

func (s *monitorUnitTest) TestPathCacheEvictionAndTTL() {
	c := newPathCache()
	c.store("/a", pathCacheEntry{resolved: "/watched/a", at: time.Now()})
	if e, ok := c.lookup("/a"); !ok || e.resolved != "/watched/a" {
		s.Failf("lookup", "got (%+v,%v)", e, ok)
	}
	// expired entry is a miss
	c.m["/a"] = pathCacheEntry{resolved: "/watched/a", at: time.Now().Add(-2 * pathCacheTTL)}
	if _, ok := c.lookup("/a"); ok {
		s.Fail("expired entry must be a miss")
	}
	// overflow clears the map
	for i := range pathCacheMax + 1 {
		c.store(fmt.Sprintf("/p%d", i), pathCacheEntry{at: time.Now()})
	}
	s.Require().LessOrEqual(len(c.m), pathCacheMax)
}

// TestResolveSymlinkCachesResult: the EvalSymlinks resolution is memoized — a
// symlink resolved once keeps reporting its target within the TTL even after
// the link is removed (the documented staleness tradeoff that keeps the
// firehose cheap). The filtering decision is unchanged.
func (s *monitorUnitTest) TestResolveSymlinkCachesResult() {
	dir := s.T().TempDir()
	watchedDir := filepath.Join(dir, "watched")
	s.Require().NoError(os.Mkdir(watchedDir, 0755))
	watchedFile := filepath.Join(watchedDir, "secret.txt")
	s.Require().NoError(os.WriteFile(watchedFile, []byte("x"), 0644))
	link := filepath.Join(dir, "lnk")
	s.Require().NoError(os.Symlink(watchedFile, link))

	m := &Monitor{targets: dirTarget(watchedDir), paths: newPathCache()}

	s.Require().Equal(watchedFile, m.resolveSymlinkToTarget(link))

	s.Require().NoError(os.Remove(link))
	s.Require().Equal(watchedFile, m.resolveSymlinkToTarget(link), "resolution stays cached within TTL")
}

func (s *monitorUnitTest) TestCheckHardlinkByInode() {
	dir := s.T().TempDir()

	watchedDir := filepath.Join(dir, "watched")
	s.Require().NoError(os.Mkdir(watchedDir, 0755))
	watchedFile := filepath.Join(watchedDir, "target.txt")
	s.Require().NoError(os.WriteFile(watchedFile, []byte("data"), 0644))

	hardlinkPath := filepath.Join(dir, "mylink")
	s.Require().NoError(os.Link(watchedFile, hardlinkPath))

	m := &Monitor{
		targets:     dirTarget(watchedDir),
		watchInodes: make(map[uint64][]watchedInode),
		mounts:      newMountTable(),
	}
	m.addInode(watchedFile)
	dev, ino := s.kernelIdentity(m, watchedFile)

	tests := []struct {
		name     string
		dev      uint32
		ino      uint64
		path     string
		expected string
	}{
		{name: "original file matches", dev: dev, ino: ino, path: watchedFile, expected: watchedFile},
		{name: "hardlink resolves to watched file path", dev: dev, ino: ino, path: hardlinkPath, expected: watchedFile},
		{name: "other superblock", dev: dev + 1, ino: ino, path: hardlinkPath, expected: ""},
		{name: "inode not watched", dev: dev, ino: ino + 1, path: hardlinkPath, expected: ""},
		{name: "user-string event has no inode", dev: 0, ino: 0, path: hardlinkPath, expected: ""},
		{name: "path no longer that inode", dev: dev, ino: ino, path: filepath.Join(dir, "nonexistent"), expected: ""},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.Require().Equal(tt.expected, m.hardlinkByInode(tt.dev, tt.ino, tt.path))
		})
	}
}

// A watched name that left the target (renamed out, or unlinked and its inode number reused by a
// new file outside) shares no file with the event path any more.
func (s *monitorUnitTest) TestCheckHardlinkByInode_WatchedNameGone() {
	dir := s.T().TempDir()
	watchedDir := filepath.Join(dir, "watched")
	s.Require().NoError(os.Mkdir(watchedDir, 0755))
	moved, deleted := filepath.Join(watchedDir, "moved"), filepath.Join(watchedDir, "deleted")
	s.Require().NoError(os.WriteFile(moved, []byte("data"), 0644))
	s.Require().NoError(os.WriteFile(deleted, []byte("data"), 0644))

	m := &Monitor{
		targets:     dirTarget(watchedDir),
		watchInodes: make(map[uint64][]watchedInode),
		mounts:      newMountTable(),
	}
	m.populateWatchInodes()

	outside := filepath.Join(dir, "outside")
	movedDev, movedIno := s.kernelIdentity(m, moved)
	s.Require().NoError(os.Rename(moved, outside))
	s.Require().Equal("", m.hardlinkByInode(movedDev, movedIno, outside), "renamed out of the target")

	_, deletedIno := s.kernelIdentity(m, deleted)
	s.Require().NoError(os.Remove(deleted))
	fresh := filepath.Join(dir, "fresh")
	s.Require().NoError(os.WriteFile(fresh, []byte("x"), 0644))
	freshDev, freshIno := s.kernelIdentity(m, fresh)
	m.watchInodes[freshIno] = m.watchInodes[deletedIno]
	s.Require().Equal("", m.hardlinkByInode(freshDev, freshIno, fresh), "unlinked, inode number reused")
}

func (s *monitorUnitTest) TestCheckHardlinkByInodeNoInodes() {
	m := &Monitor{
		targets:     dirTarget("/tmp/watch"),
		watchInodes: nil,
	}
	s.Require().Equal("", m.hardlinkByInode(1, 42, "/some/path"))
}

// kernelIdentity returns path's (sb->s_dev, ino) as the BPF event carries them.
func (s *monitorUnitTest) kernelIdentity(m *Monitor, path string) (uint32, uint64) {
	real, err := filepath.EvalSymlinks(path)
	s.Require().NoError(err)
	var st syscall.Stat_t
	s.Require().NoError(syscall.Lstat(real, &st))
	dev := m.mounts.devOf(real)
	s.Require().NotZero(dev)
	return dev, st.Ino
}

func (s *monitorUnitTest) TestPopulateWatchInodes_SingleFile() {
	dir := s.T().TempDir()
	targetFile := filepath.Join(dir, "target.txt")
	s.Require().NoError(os.WriteFile(targetFile, []byte("data"), 0644))

	m := &Monitor{
		targets: []ebpf.Target{
			{Dir: dir, File: "target.txt", IsDir: false, Path: targetFile},
		},
		watchInodes: make(map[uint64][]watchedInode),
	}
	m.populateWatchInodes()

	s.Require().Len(m.watchInodes, 1)
	for ino, ws := range m.watchInodes {
		s.NotZero(ino)
		s.Require().Len(ws, 1)
		s.Equal(targetFile, ws[0].path)
	}
}

func (s *monitorUnitTest) TestPopulateWatchInodes_Directory() {
	dir := s.T().TempDir()

	s.Require().NoError(os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a"), 0644))
	s.Require().NoError(os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0644))
	subDir := filepath.Join(dir, "sub")
	s.Require().NoError(os.Mkdir(subDir, 0755))
	s.Require().NoError(os.WriteFile(filepath.Join(subDir, "c.txt"), []byte("c"), 0644))

	m := &Monitor{
		targets:     dirTarget(dir),
		recursive:   true,
		depth:       0,
		watchInodes: make(map[uint64][]watchedInode),
	}
	m.populateWatchInodes()

	s.Require().Len(m.watchInodes, 3)
	for _, ws := range m.watchInodes {
		s.Contains(ws[0].path, dir)
	}
}

func (s *monitorUnitTest) TestResolveSymlinkToTarget_Chained() {
	dir := s.T().TempDir()

	watchedDir := filepath.Join(dir, "watched")
	s.Require().NoError(os.Mkdir(watchedDir, 0755))
	watchedFile := filepath.Join(watchedDir, "secret.txt")
	s.Require().NoError(os.WriteFile(watchedFile, []byte("data"), 0644))

	link2 := filepath.Join(dir, "link2")
	s.Require().NoError(os.Symlink(watchedFile, link2))
	link1 := filepath.Join(dir, "link1")
	s.Require().NoError(os.Symlink(link2, link1))

	m := &Monitor{targets: dirTarget(watchedDir)}

	result := m.resolveSymlinkToTarget(link1)
	s.Require().Equal(watchedFile, result)
}

func (s *monitorUnitTest) TestResolveSymlinkToTarget_Relative() {
	dir := s.T().TempDir()

	watchedDir := filepath.Join(dir, "watched")
	s.Require().NoError(os.Mkdir(watchedDir, 0755))
	watchedFile := filepath.Join(watchedDir, "secret.txt")
	s.Require().NoError(os.WriteFile(watchedFile, []byte("data"), 0644))

	subDir := filepath.Join(dir, "sub")
	s.Require().NoError(os.Mkdir(subDir, 0755))
	relSymlink := filepath.Join(subDir, "link")
	s.Require().NoError(os.Symlink("../watched/secret.txt", relSymlink))

	m := &Monitor{targets: dirTarget(watchedDir)}

	result := m.resolveSymlinkToTarget(relSymlink)
	s.Require().Equal(watchedFile, result)
}

func (s *monitorUnitTest) TestResolveSymlinkToTarget_NonMatching() {
	dir := s.T().TempDir()

	watchedDir := filepath.Join(dir, "watched")
	s.Require().NoError(os.Mkdir(watchedDir, 0755))
	watchedFile := filepath.Join(watchedDir, "secret.txt")
	s.Require().NoError(os.WriteFile(watchedFile, []byte("data"), 0644))

	outsideDir := filepath.Join(dir, "outside")
	s.Require().NoError(os.Mkdir(outsideDir, 0755))
	outsideFile := filepath.Join(outsideDir, "other.txt")
	s.Require().NoError(os.WriteFile(outsideFile, []byte("other"), 0644))
	symToOutside := filepath.Join(dir, "link_to_outside")
	s.Require().NoError(os.Symlink(outsideFile, symToOutside))

	m := &Monitor{targets: dirTarget(watchedDir)}

	result := m.resolveSymlinkToTarget(symToOutside)
	s.Require().Equal(symToOutside, result, "non-matching symlink should return original path")
}

func (s *monitorUnitTest) TestResolveFDPath_SkipsReadWrite() {
	dir := s.T().TempDir()
	targetFile := filepath.Join(dir, "target.txt")
	s.Require().NoError(os.WriteFile(targetFile, []byte("data"), 0644))

	for _, typ := range []ebpf.EventType{ebpf.EventRead, ebpf.EventWrite} {
		ev := &ebpf.FileEvent{
			PID:  uint32(os.Getpid()),
			FD:   3,
			Type: typ,
		}
		m := &Monitor{}
		m.resolveFDPath(ev)

		s.Require().Empty(ev.Path, "resolveFDPath must not resolve fd for %s", typ)
	}
}

func (s *monitorUnitTest) TestResolveFDPath_Mmap() {
	dir := s.T().TempDir()
	targetFile := filepath.Join(dir, "target.txt")
	s.Require().NoError(os.WriteFile(targetFile, []byte("data"), 0644))

	f, err := os.Open(targetFile)
	s.Require().NoError(err)
	defer f.Close()

	ev := &ebpf.FileEvent{
		PID:  uint32(os.Getpid()),
		FD:   uint32(f.Fd()),
		Type: ebpf.EventMmap,
	}
	m := &Monitor{}
	m.resolveFDPath(ev)

	s.Require().Equal(targetFile, ev.Path)
}

func (s *monitorUnitTest) TestResolveFDPath_SkipNonFdType() {
	ev := &ebpf.FileEvent{
		PID:  1234,
		FD:   5,
		Type: ebpf.EventOpen,
		Path: "/some/path",
	}
	m := &Monitor{}
	m.resolveFDPath(ev)

	s.Require().Equal("/some/path", ev.Path)
}

func (s *monitorUnitTest) TestResolveRelativePath_Absolute() {
	m := &Monitor{}
	result := m.resolveRelativePath(uint32(os.Getpid()), "/absolute/path")
	s.Require().Equal("/absolute/path", result)
}

func (s *monitorUnitTest) TestResolveRelativePath_Relative() {
	dir := s.T().TempDir()
	s.Require().NoError(os.Chdir(dir))

	targetFile := filepath.Join(dir, "target.txt")
	s.Require().NoError(os.WriteFile(targetFile, []byte("data"), 0644))

	m := &Monitor{}
	result := m.resolveRelativePath(uint32(os.Getpid()), "target.txt")
	s.Require().Equal(targetFile, result)
}

func (s *monitorUnitTest) TestCheckHardlinkByInode_Preexisting() {
	dir := s.T().TempDir()

	watchedDir := filepath.Join(dir, "watched")
	s.Require().NoError(os.Mkdir(watchedDir, 0755))
	watchedFile := filepath.Join(watchedDir, "target.txt")
	s.Require().NoError(os.WriteFile(watchedFile, []byte("data"), 0644))

	hardlinkPath := filepath.Join(dir, "preexisting_hardlink")
	s.Require().NoError(os.Link(watchedFile, hardlinkPath))

	m := &Monitor{
		targets:     dirTarget(watchedDir),
		watchInodes: make(map[uint64][]watchedInode),
		mounts:      newMountTable(),
	}
	m.populateWatchInodes()

	dev, ino := s.kernelIdentity(m, hardlinkPath)
	s.Require().Equal(watchedFile, m.hardlinkByInode(dev, ino, hardlinkPath))
}

func TestSetEventTypesRace(t *testing.T) {
	m := &Monitor{}
	m.eventTypes.Store(ebpf.EventTypes())

	var wg sync.WaitGroup
	wg.Go(func() {
		for range 100 {
			m.SetEventTypes([]ebpf.EventType{ebpf.EventOpen, ebpf.EventRead})
		}
	})

	for range 100 {
		m.eventTypeAllowed(ebpf.EventOpen)
		m.eventTypeAllowed(ebpf.EventRead)
		m.eventTypeAllowed(ebpf.EventWrite)
	}
	wg.Wait()
}

func TestEventTypeAllowedConcurrent(t *testing.T) {
	m := &Monitor{}
	m.eventTypes.Store([]ebpf.EventType{ebpf.EventOpen, ebpf.EventRead, ebpf.EventWrite, ebpf.EventDelete})

	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			for range 100 {
				_ = m.eventTypeAllowed(ebpf.EventOpen)
				_ = m.eventTypeAllowed(ebpf.EventWrite)
				_ = m.eventTypeAllowed(ebpf.EventMmap)
			}
		})
	}
	wg.Wait()
}

func TestSetEventTypesBeforeStartNoRace(t *testing.T) {
	m := &Monitor{}
	m.SetEventTypes([]ebpf.EventType{ebpf.EventOpen})

	types := m.eventTypes.Load().([]ebpf.EventType)
	if len(types) != 1 || types[0] != ebpf.EventOpen {
		t.Fatalf("unexpected event types: %v", types)
	}
}

func (s *monitorUnitTest) TestEventTypeAllowedAllTypes() {
	m := &Monitor{}
	m.eventTypes.Store(ebpf.EventTypes())

	for _, et := range ebpf.EventTypes() {
		s.Run(et.String(), func() {
			s.Require().True(m.eventTypeAllowed(et), "event type %s should be allowed by default", et)
		})
	}
}

func (s *monitorUnitTest) TestEventTypeAllowedMkdir() {
	m := &Monitor{}
	m.eventTypes.Store(ebpf.EventTypes())

	s.Require().True(m.eventTypeAllowed(ebpf.EventMkdir))
}

func (s *monitorUnitTest) TestResolveRelativePathMkdirRelative() {
	dir := s.T().TempDir()
	s.Require().NoError(os.Chdir(dir))

	m := &Monitor{}
	result := m.resolveRelativePath(uint32(os.Getpid()), "newdir")
	expected := filepath.Join(dir, "newdir")
	s.Require().Equal(expected, result)
}

func (s *monitorUnitTest) TestMatchesFilterMkdirNewDir() {
	m := &Monitor{targets: dirTarget("/watch"), recursive: false}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{name: "direct child new dir", path: "/watch/newdir", expected: true},
		{name: "new dir with depth", path: "/watch/sub/newdir", expected: false},
		{name: "outside watch", path: "/other/newdir", expected: false},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ev := &ebpf.FileEvent{Path: tt.path, Type: ebpf.EventMkdir}
			result := m.matchesFilter(ev)
			s.Require().Equal(tt.expected, result)
		})
	}
}

func (s *monitorUnitTest) TestMatchesFilterMkdirNewDirRecursive() {
	m := &Monitor{targets: dirTarget("/watch"), recursive: true, depth: 0}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{name: "direct child new dir", path: "/watch/newdir", expected: true},
		{name: "nested new dir", path: "/watch/a/b/newdir", expected: true},
		{name: "outside watch", path: "/other/newdir", expected: false},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ev := &ebpf.FileEvent{Path: tt.path, Type: ebpf.EventMkdir}
			result := m.matchesFilter(ev)
			s.Require().Equal(tt.expected, result)
		})
	}
}

// Every event must come from reserve_event: a probe reserving the ring buffer itself leaves exe_dev,
// exe_ino and prot as stale ring-buffer bytes, and the recorder keys accessors on them.
func TestEveryEventReservedWithIdentity(t *testing.T) {
	_, self, _, _ := runtime.Caller(0) // other tests chdir away
	src, err := os.ReadFile(filepath.Join(filepath.Dir(self), "bpf", "monitor.bpf.c"))
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(string(src), "bpf_ringbuf_reserve("); n != 1 {
		t.Fatalf("monitor.bpf.c reserves the ring buffer in %d places; only reserve_event may", n)
	}
}
