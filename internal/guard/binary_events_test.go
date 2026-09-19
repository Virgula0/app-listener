package guard

import (
	"strings"
	"testing"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// Regression: the daemon records an event entry for EVERY whitelisted binary,
// nil when unrestricted. addBinaryEvents used to check the mode before the
// list length, so every ModeReadOnly guard with a whitelist (a `lib_dir`)
// failed to start with "only supported in whitelist mode".
func TestAddBinaryEventsReadOnlyAcceptsUnrestricted(t *testing.T) {
	g := &Guard{mode: ModeReadOnly}
	bins := []BinaryEntry{{Path: "/usr/bin/steam"}, {Path: "/usr/bin/lsof"}}
	events := map[string][]ebpf.EventType{"/usr/bin/steam": nil, "/usr/bin/lsof": {}}
	if err := g.addBinaryEvents(bins, events); err != nil {
		t.Fatalf("unrestricted binaries in read-only mode must be accepted, got: %v", err)
	}
}

// A real restriction still has no meaning in read-only mode and must be
// refused, not silently dropped (that would widen the binary).
func TestAddBinaryEventsReadOnlyRejectsRestriction(t *testing.T) {
	g := &Guard{mode: ModeReadOnly}
	bins := []BinaryEntry{{Path: "/usr/bin/ssh"}}
	events := map[string][]ebpf.EventType{"/usr/bin/ssh": {ebpf.EventType(1)}}
	err := g.addBinaryEvents(bins, events)
	if err == nil || !strings.Contains(err.Error(), "only supported in whitelist mode") {
		t.Fatalf("restricted binary in read-only mode must be refused, got: %v", err)
	}
}
