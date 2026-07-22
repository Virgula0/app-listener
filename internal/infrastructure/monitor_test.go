package ebpf

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

type monitorUnitTest struct {
	suite.Suite
}

func TestMonitorUnitTest(t *testing.T) {
	suite.Run(t, new(monitorUnitTest))
}

func dirTarget(dir string) []Target {
	return []Target{{Dir: dir, IsDir: true}}
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
	m := &Monitor{targets: []Target{{Dir: "/parent", File: "target.txt", IsDir: false}}}

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
			ev := &FileEvent{Path: tt.path}
			result := m.matchesFilter(ev)
			s.Require().Equal(tt.expected, result)
		})
	}
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
			ev := &FileEvent{Path: tt.path}
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
			ev := &FileEvent{Path: tt.path}
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

	// symlink outside → target inside watched dir
	symOutside := filepath.Join(dir, "outside_link")
	s.Require().NoError(os.Symlink(watchedFile, symOutside))

	// symlink inside watched dir → another target inside
	symInside := filepath.Join(watchedDir, "inside_link")
	s.Require().NoError(os.Symlink(watchedFile, symInside))

	// symlink to watched dir itself
	symDir := filepath.Join(dir, "dir_link")
	s.Require().NoError(os.Symlink(watchedDir, symDir))

	// regular file not inside watched dir
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
			expected: symInside, // already inside, no resolution needed
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

func (s *monitorUnitTest) TestCheckHardlinkByInode() {
	dir := s.T().TempDir()

	watchedDir := filepath.Join(dir, "watched")
	s.Require().NoError(os.Mkdir(watchedDir, 0755))
	watchedFile := filepath.Join(watchedDir, "target.txt")
	s.Require().NoError(os.WriteFile(watchedFile, []byte("data"), 0644))

	// hardlink to watched file
	hardlinkPath := filepath.Join(dir, "mylink")
	s.Require().NoError(os.Link(watchedFile, hardlinkPath))

	// Setup monitor with pre-scanned inodes
	m := &Monitor{
		targets:     dirTarget(watchedDir),
		watchInodes: make(map[string]string),
	}
	m.addInode(watchedFile)

	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "original file matches",
			path:     watchedFile,
			expected: watchedFile,
		},
		{
			name:     "hardlink resolves to watched file path",
			path:     hardlinkPath,
			expected: watchedFile,
		},
		{
			name:     "non-existent path",
			path:     filepath.Join(dir, "nonexistent"),
			expected: "",
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result := m.checkHardlinkByInode(tt.path)
			s.Require().Equal(tt.expected, result)
		})
	}
}

func (s *monitorUnitTest) TestCheckHardlinkByInodeNoInodes() {
	m := &Monitor{
		targets:     dirTarget("/tmp/watch"),
		watchInodes: nil,
	}
	result := m.checkHardlinkByInode("/some/path")
	s.Require().Equal("", result)
}

func (s *monitorUnitTest) TestPopulateWatchInodes_SingleFile() {
	dir := s.T().TempDir()
	targetFile := filepath.Join(dir, "target.txt")
	s.Require().NoError(os.WriteFile(targetFile, []byte("data"), 0644))

	m := &Monitor{
		targets: []Target{
			{Dir: dir, File: "target.txt", IsDir: false, Path: targetFile},
		},
		watchInodes: make(map[string]string),
	}
	m.populateWatchInodes()

	s.Require().Len(m.watchInodes, 1)
	for key, val := range m.watchInodes {
		s.True(len(key) > 0, "key should be dev:ino")
		s.Equal(targetFile, val)
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
		targets:   dirTarget(dir),
		recursive: true,
		depth:     0,
		watchInodes: make(map[string]string),
	}
	m.populateWatchInodes()

	// Should have 3 files (a.txt, b.txt, sub/c.txt)
	s.Require().Len(m.watchInodes, 3)
	for _, val := range m.watchInodes {
		s.Contains(val, dir)
	}
}

func (s *monitorUnitTest) TestResolveSymlinkToTarget_Chained() {
	dir := s.T().TempDir()

	watchedDir := filepath.Join(dir, "watched")
	s.Require().NoError(os.Mkdir(watchedDir, 0755))
	watchedFile := filepath.Join(watchedDir, "secret.txt")
	s.Require().NoError(os.WriteFile(watchedFile, []byte("data"), 0644))

	// link1 -> link2 -> watchedFile
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

	// symlink with relative target ../watched/secret.txt
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

	for _, typ := range []EventType{EventRead, EventWrite} {
		ev := &FileEvent{
			PID:  uint32(os.Getpid()),
			FD:   3,
			Type: typ,
		}
		m := &Monitor{}
		m.resolveFDPath(ev)

		// vfs_read/vfs_write set Path directly in BPF; resolveFDPath
		// only handles EventMmap. The path should remain empty.
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

	ev := &FileEvent{
		PID:  uint32(os.Getpid()),
		FD:   uint32(f.Fd()),
		Type: EventMmap,
	}
	m := &Monitor{}
	m.resolveFDPath(ev)

	s.Require().Equal(targetFile, ev.Path)
}

func (s *monitorUnitTest) TestResolveFDPath_SkipNonFdType() {
	ev := &FileEvent{
		PID:  1234,
		FD:   5,
		Type: EventOpen,
		Path: "/some/path",
	}
	m := &Monitor{}
	m.resolveFDPath(ev)

	// Path should remain unchanged (EventOpen not handled by resolveFDPath)
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

	// Hardlink created BEFORE monitor starts
	hardlinkPath := filepath.Join(dir, "preexisting_hardlink")
	s.Require().NoError(os.Link(watchedFile, hardlinkPath))

	// Monitor startup path: scan watched dir inodes
	m := &Monitor{
		targets:     dirTarget(watchedDir),
		watchInodes: make(map[string]string),
	}
	m.populateWatchInodes()

	// Access via hardlink should be detected
	result := m.checkHardlinkByInode(hardlinkPath)
	s.Require().Equal(watchedFile, result)
}

func TestSetEventTypesRace(t *testing.T) {
	m := &Monitor{}
	m.eventTypes.Store(EventTypes())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range 100 {
			m.SetEventTypes([]EventType{EventOpen, EventRead})
		}
	}()

	for range 100 {
		m.eventTypeAllowed(EventOpen)
		m.eventTypeAllowed(EventRead)
		m.eventTypeAllowed(EventWrite)
	}
	wg.Wait()
}

func TestEventTypeAllowedConcurrent(t *testing.T) {
	m := &Monitor{}
	m.eventTypes.Store([]EventType{EventOpen, EventRead, EventWrite, EventDelete})

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 100 {
				_ = m.eventTypeAllowed(EventOpen)
				_ = m.eventTypeAllowed(EventWrite)
				_ = m.eventTypeAllowed(EventMmap)
			}
		}()
	}
	wg.Wait()
}

func TestSetEventTypesBeforeStartNoRace(t *testing.T) {
	m := &Monitor{}
	m.SetEventTypes([]EventType{EventOpen})

	types := m.eventTypes.Load().([]EventType)
	if len(types) != 1 || types[0] != EventOpen {
		t.Fatalf("unexpected event types: %v", types)
	}
}
