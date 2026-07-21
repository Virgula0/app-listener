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

func (s *monitorUnitTest) TestInWatchDir() {
	m := &Monitor{watchDir: "/tmp/watch"}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{name: "exact match", path: "/tmp/watch", expected: true},
		{name: "child path", path: "/tmp/watch/file.txt", expected: true},
		{name: "nested child", path: "/tmp/watch/sub/dir/file.txt", expected: true},
		{name: "non matching", path: "/var/log/syslog", expected: false},
		{name: "similar prefix match", path: "/tmp/watch_extra", expected: true},
		{name: "empty path", path: "", expected: false},
		{name: "root path", path: "/", expected: false},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result := m.inWatchDir(tt.path)
			s.Require().Equal(tt.expected, result)
		})
	}
}

func (s *monitorUnitTest) TestMatchesFilterNonRecursive() {
	m := &Monitor{watchDir: "/watch", recursive: false}

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
	m := &Monitor{watchDir: "/watch", recursive: true, depth: 0}

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{name: "watch dir itself", path: "/watch", expected: true},
		{name: "direct child", path: "/watch/file.txt", expected: true},
		{name: "nested child", path: "/watch/sub/file.txt", expected: true},
		{name: "deeply nested", path: "/watch/a/b/c/d/e.txt", expected: true},
		{name: "outside watch (not rejected by matchesFilter alone)", path: "/other/file.txt", expected: true},
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
	m := &Monitor{watchDir: "/watch", recursive: true, depth: 2}

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
		{name: "outside watch (exceeds depth due to rel path)", path: "/other/file.txt", expected: false},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			ev := &FileEvent{Path: tt.path}
			result := m.matchesFilter(ev)
			s.Require().Equal(tt.expected, result)
		})
	}
}
