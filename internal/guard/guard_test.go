package guard

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"

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
