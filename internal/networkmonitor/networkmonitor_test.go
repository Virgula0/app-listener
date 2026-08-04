package networkmonitor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Virgula0/app-listener/internal/infrastructure"
)

func TestComputeBinaryEntry(t *testing.T) {
	tmp := t.TempDir()
	binPath := filepath.Join(tmp, "testbin")
	if err := os.WriteFile(binPath, []byte("#!/bin/sh\necho hello"), 0755); err != nil {
		t.Fatal(err)
	}

	entry, err := ComputeBinaryEntry(binPath)
	if err != nil {
		t.Fatalf("ComputeBinaryEntry: %v", err)
	}
	if entry.Path != binPath {
		t.Errorf("expected path %q, got %q", binPath, entry.Path)
	}
	if entry.Comm != "testbin" {
		t.Errorf("expected comm testbin, got %q", entry.Comm)
	}
	if len(entry.Hash) != 32 {
		t.Errorf("expected 32-byte hash, got %d", len(entry.Hash))
	}
}

func TestComputeBinaryEntry_TruncatedComm(t *testing.T) {
	tmp := t.TempDir()
	longName := "very-long-binary-name-over-16-chars"
	binPath := filepath.Join(tmp, longName)
	if err := os.WriteFile(binPath, []byte("data"), 0755); err != nil {
		t.Fatal(err)
	}

	entry, err := ComputeBinaryEntry(binPath)
	if err != nil {
		t.Fatalf("ComputeBinaryEntry: %v", err)
	}
	if len(entry.Comm) > 15 {
		t.Errorf("comm should be truncated to 15 chars, got %q (len=%d)", entry.Comm, len(entry.Comm))
	}
}

func TestNetworkMonitor_NewFailure(t *testing.T) {
	_, err := NewNetworkMonitor(nil)
	if err == nil {
		t.Skip("NewNetworkMonitor with nil binaries may succeed if BPF is available")
	}
}

func TestNetEventTypes(t *testing.T) {
	types := ebpf.NetEventTypes()
	expected := []ebpf.NetEventType{
		ebpf.NetConnect, ebpf.NetAccept, ebpf.NetSend,
		ebpf.NetRecv, ebpf.NetClose, ebpf.NetDNS,
		ebpf.NetBind, ebpf.NetListen,
	}
	if len(types) != len(expected) {
		t.Errorf("expected %d types, got %d", len(expected), len(types))
	}
	for i, et := range expected {
		if types[i] != et {
			t.Errorf("types[%d] = %d, want %d", i, types[i], et)
		}
	}
}

func TestParseNetEventType(t *testing.T) {
	tests := []struct {
		input string
		want  ebpf.NetEventType
		ok    bool
	}{
		{"CONNECT", ebpf.NetConnect, true},
		{"connect", ebpf.NetConnect, true},
		{"ACCEPT", ebpf.NetAccept, true},
		{"SEND", ebpf.NetSend, true},
		{"RECV", ebpf.NetRecv, true},
		{"CLOSE", ebpf.NetClose, true},
		{"DNS", ebpf.NetDNS, true},
		{"BIND", ebpf.NetBind, true},
		{"LISTEN", ebpf.NetListen, true},
		{"INVALID", 0, false},
		{"", 0, false},
	}
	for _, tc := range tests {
		got, ok := ebpf.ParseNetEventType(tc.input)
		if ok != tc.ok {
			t.Errorf("ParseNetEventType(%q) ok=%v, want %v", tc.input, ok, tc.ok)
		}
		if ok && got != tc.want {
			t.Errorf("ParseNetEventType(%q) = %d, want %d", tc.input, got, tc.want)
		}
	}
}

func TestProtocolString(t *testing.T) {
	tests := []struct {
		proto uint32
		want  string
	}{
		{6, "TCP"},
		{17, "UDP"},
		{1, "ICMP"},
		{58, "ICMPv6"},
		{0, "UNKNOWN"},
		{99, "UNKNOWN"},
	}
	for _, tc := range tests {
		got := ebpf.ProtocolString(tc.proto)
		if got != tc.want {
			t.Errorf("ProtocolString(%d) = %q, want %q", tc.proto, got, tc.want)
		}
	}
}

func TestNetEventTypeString(t *testing.T) {
	tests := []struct {
		et   ebpf.NetEventType
		want string
	}{
		{ebpf.NetConnect, "CONNECT"},
		{ebpf.NetAccept, "ACCEPT"},
		{ebpf.NetSend, "SEND"},
		{ebpf.NetRecv, "RECV"},
		{ebpf.NetClose, "CLOSE"},
		{ebpf.NetDNS, "DNS"},
		{ebpf.NetBind, "BIND"},
		{ebpf.NetListen, "LISTEN"},
		{ebpf.NetEventType(99), "UNKNOWN"},
	}
	for _, tc := range tests {
		got := tc.et.String()
		if got != tc.want {
			t.Errorf("NetEventType(%d).String() = %q, want %q", tc.et, got, tc.want)
		}
	}
}

func TestBinariesSummary(t *testing.T) {
	entries := []BinaryEntry{
		{Path: "/usr/bin/test", Hash: [32]byte{0x01, 0x02}},
	}
	summary := binariesSummary(entries)
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestStatInodeKey(t *testing.T) {
	ik, err := statInodeKey("/proc/self/exe")
	if err != nil {
		t.Fatalf("statInodeKey: %v", err)
	}
	if ik.Ino == 0 {
		t.Error("expected non-zero inode")
	}
}

func TestStatInodeDevEncoding(t *testing.T) {
	dev, ino, err := ebpf.StatInode("/proc/self/exe")
	if err != nil {
		t.Fatalf("StatInode: %v", err)
	}
	if dev == 0 {
		t.Error("expected non-zero dev")
	}
	if ino == 0 {
		t.Error("expected non-zero inode")
	}
}
