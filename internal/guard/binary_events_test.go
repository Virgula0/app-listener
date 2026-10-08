package guard

import (
	"strings"
	"testing"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// The daemon records an event entry for EVERY whitelisted binary, nil when unrestricted: a
// ModeReadOnly guard with a whitelist (a `lib_dir`) must still start.
func TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(t *testing.T) {
	g := &Guard{mode: ModeReadOnly}
	bins := []BinaryEntry{{Path: "/usr/bin/steam"}, {Path: "/usr/bin/lsof"}}
	events := map[string][]ebpf.EventType{"/usr/bin/steam": nil, "/usr/bin/lsof": {}}
	if err := g.checkBinaryEvents(bins, events); err != nil {
		t.Fatalf("unrestricted binaries in read-only mode must be accepted, got: %v", err)
	}
}

// A real restriction still has no meaning in read-only mode and must be
// refused, not silently dropped (that would widen the binary).
func TestCheckBinaryEventsReadOnlyRejectsRestriction(t *testing.T) {
	g := &Guard{mode: ModeReadOnly}
	bins := []BinaryEntry{{Path: "/usr/bin/ssh"}}
	events := map[string][]ebpf.EventType{"/usr/bin/ssh": {ebpf.EventType(1)}}
	err := g.checkBinaryEvents(bins, events)
	if err == nil || !strings.Contains(err.Error(), "only supported in whitelist mode") {
		t.Fatalf("restricted binary in read-only mode must be refused, got: %v", err)
	}
}

// addBinaryActions validates every list before writing a row: no maps are loaded here, so any
// write would panic.
func TestAddBinaryActionsRefusesBeforeWriting(t *testing.T) {
	g := &Guard{mode: ModeReadOnly}
	bins := []BinaryEntry{{Path: "/usr/bin/true"}, {Path: "/usr/bin/ssh"}}
	events := map[string][]ebpf.EventType{"/usr/bin/ssh": {ebpf.EventType(1)}}
	if err := g.addBinaryActions(bins, events); err == nil {
		t.Fatal("restricted binary in read-only mode must be refused")
	}
}

// The chmod allowance must never reach a whitelist (secret) or blacklist tree; refused before any
// BPF slot is taken.
func TestWithChmodDropWriteNeedsReadOnly(t *testing.T) {
	for _, mode := range []Mode{ModeWhitelist, ModeBlacklist} {
		if _, err := NewGuard(t.TempDir(), mode, nil, true, 0, WithChmodDropWrite()); err == nil {
			t.Errorf("mode %s accepted WithChmodDropWrite", modeString(mode))
		}
	}
}
