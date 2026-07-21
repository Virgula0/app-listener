package ebpf

import (
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
