package guard

import (
	"bytes"
	"encoding/binary"
	"testing"
)

// bunSection lays out a `.bun` section as bun 1.3 writes it: u64 length, graph, Offsets, trailer.
func bunSection(graph []byte, flags uint32) []byte {
	off := make([]byte, bunOffsetsSize)
	binary.LittleEndian.PutUint64(off, uint64(len(graph)))
	binary.LittleEndian.PutUint32(off[bunOffsetsSize-4:], flags)
	body := append(append(append([]byte(nil), graph...), off...), bunTrailer...)
	head := make([]byte, 8)
	binary.LittleEndian.PutUint64(head, uint64(len(body)))
	return append(head, body...)
}

func TestBunGraphFlags(t *testing.T) {
	graph := bytes.Repeat([]byte{7}, 140)
	for name, tc := range map[string]struct {
		sec      []byte
		ok       bool
		autoload bool
	}{
		"default compile (0x0c)":     {bunSection(graph, 0x0c), true, true},
		"autoload disabled (0x3bff)": {bunSection(graph, 0x3bff), true, false},
		"bad trailer":                {append(bunSection(graph, 0x3bff)[:len(bunSection(graph, 0x3bff))-1], 'x'), false, true},
		"padded":                     {append(bunSection(graph, 0x3bff), 0), false, true},
		"short":                      {[]byte("tiny"), false, true},
	} {
		flags, ok := bunGraphFlags(bytes.NewReader(tc.sec), uint64(len(tc.sec)))
		if ok != tc.ok {
			t.Errorf("%s: ok = %v, want %v", name, ok, tc.ok)
		}
		if autoload := !ok || flags&bunDisableAutoloadBunfig == 0; autoload != tc.autoload {
			t.Errorf("%s: autoload = %v, want %v", name, autoload, tc.autoload)
		}
	}
	// A length prefix that disagrees with byte_count: another layout.
	sec := bunSection(graph, 0x3bff)
	binary.LittleEndian.PutUint64(sec[8+len(graph):], uint64(len(graph)-1))
	if _, ok := bunGraphFlags(bytes.NewReader(sec), uint64(len(sec))); ok {
		t.Error("a byte_count that disagrees with the length prefix must not parse")
	}
}
