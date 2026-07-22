package ebpf

import (
	"os"
	"path/filepath"
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
