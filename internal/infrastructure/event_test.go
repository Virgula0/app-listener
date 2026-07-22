package ebpf

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type eventUnitTest struct {
	suite.Suite
}

func TestEventUnitTest(t *testing.T) {
	suite.Run(t, new(eventUnitTest))
}

func (s *eventUnitTest) TestCstr() {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{name: "null terminated", input: []byte("hello\x00world"), expected: "hello"},
		{name: "full buffer without null", input: []byte{'a', 'b', 'c'}, expected: "abc"},
		{name: "empty buffer", input: []byte{}, expected: ""},
		{name: "starts with null", input: []byte{0, 'a', 'b'}, expected: ""},
		{name: "multiple nulls", input: []byte("ab\x00\x00c"), expected: "ab"},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			result := Cstr(tt.input)
			s.Require().Equal(tt.expected, result)
		})
	}
}

func (s *eventUnitTest) TestEventTypeString() {
	tests := []struct {
		name     string
		etype    EventType
		expected string
	}{
		{name: "open", etype: EventOpen, expected: "OPEN"},
		{name: "read", etype: EventRead, expected: "READ"},
		{name: "write", etype: EventWrite, expected: "WRITE"},
		{name: "delete", etype: EventDelete, expected: "DELETE"},
		{name: "rename", etype: EventRename, expected: "RENAME"},
		{name: "symlink", etype: EventSymlink, expected: "SYMLINK"},
		{name: "hardlink", etype: EventHardlink, expected: "HARDLINK"},
		{name: "mkdir", etype: EventMkdir, expected: "MKDIR"},
		{name: "mmap", etype: EventMmap, expected: "MMAP"},
		{name: "unknown", etype: EventType(99), expected: "UNKNOWN"},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.Require().Equal(tt.expected, tt.etype.String())
		})
	}
}

func (s *eventUnitTest) TestBpfEventToFileEvent() {
	var comm [16]byte
	copy(comm[:], "test-proc")
	var path [256]byte
	copy(path[:], "/tmp/test.txt")
	var dest [256]byte
	copy(dest[:], "/tmp/test2.txt")

	be := &BpfEvent{
		PID:  1234,
		UID:  1000,
		GID:  1000,
		Type: uint32(EventOpen),
		FD:   3,
		Comm: comm,
		Path: path,
		Dest: dest,
	}

	fe := be.ToFileEvent()

	s.Require().Equal(uint32(1234), fe.PID)
	s.Require().Equal(uint32(1000), fe.UID)
	s.Require().Equal(uint32(1000), fe.GID)
	s.Require().Equal(EventOpen, fe.Type)
	s.Require().Equal(uint32(3), fe.FD)
	s.Require().Equal("test-proc", fe.Comm)
	s.Require().Equal("/tmp/test.txt", fe.Path)
	s.Require().Equal("/tmp/test2.txt", fe.Dest)
}

func (s *eventUnitTest) TestBpfEventToFileEventZeroValue() {
	be := &BpfEvent{}
	fe := be.ToFileEvent()

	s.Require().Equal(uint32(0), fe.PID)
	s.Require().Equal(EventType(0), fe.Type)
	s.Require().Equal("", fe.Comm)
	s.Require().Equal("", fe.Path)
	s.Require().Equal("", fe.Dest)
}
