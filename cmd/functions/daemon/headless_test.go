package daemon

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/cmd/common"
	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/usecase"
)

// noPID is PID_MAX_LIMIT: pids are always below it, so /proc/<noPID>/exe never resolves on any host.
const noPID = 4194304

func deniedEvent() usecase.DaemonEvent {
	return usecase.DaemonEvent{
		Resource: "/home/alice/.ssh",
		Event: guard.GuardEvent{FileEvent: ebpf.FileEvent{
			PID: noPID, UID: 1000, Comm: "ssh",
			Type: ebpf.EventOpen, Path: "/home/alice/.ssh/authorized_keys",
		}, Blocked: true},
	}
}

func allowedEvent() usecase.DaemonEvent {
	return usecase.DaemonEvent{
		Resource: "/home/alice/.ssh",
		Event: guard.GuardEvent{FileEvent: ebpf.FileEvent{
			PID: 1234, UID: 1000, Comm: "ssh",
			Type: ebpf.EventOpen, Path: "/home/alice/.ssh/authorized_keys",
		}, Blocked: false},
	}
}

func TestWriteEventBlocked(t *testing.T) {
	uidr := common.NewUIDResolver()
	for _, only := range []bool{false, true} {
		var buf bytes.Buffer
		ev := deniedEvent()
		if !writeEvent(&buf, only, uidr, &ev) {
			t.Fatalf("denied event must always be written (blockedOnly=%v)", only)
		}
		line := buf.String()
		if !strings.HasPrefix(line, "<4>DAEMON DENIED") {
			t.Errorf("line %q: missing syslog warning marker", line)
		}
		for _, want := range []string{
			"resource=/home/alice/.ssh",
			"op=OPEN",
			"comm=ssh",
			"commFullPath=~",
			"path=/home/alice/.ssh/authorized_keys",
			"pid=4194304",
		} {
			if !strings.Contains(line, want) {
				t.Errorf("line %q: missing field %q", line, want)
			}
		}
		if !strings.Contains(line, "uid="+uidr.Resolve(1000)) {
			t.Errorf("line %q: missing resolved uid field", line)
		}
	}
}

func TestWriteEventAllowed(t *testing.T) {
	uidr := common.NewUIDResolver()
	var buf bytes.Buffer
	ev := allowedEvent()
	if !writeEvent(&buf, false, uidr, &ev) {
		t.Fatal("allowed event must be written when blockedOnly is false")
	}
	if line := buf.String(); !strings.HasPrefix(line, "<6>DAEMON ALLOWED") {
		t.Errorf("line %q: missing syslog info prefix", line)
	}
}

func TestWriteEventAllowedSuppressed(t *testing.T) {
	uidr := common.NewUIDResolver()
	var buf bytes.Buffer
	ev := allowedEvent()
	if writeEvent(&buf, true, uidr, &ev) {
		t.Fatal("allowed event must be dropped under --blocked-only")
	}
	if buf.Len() != 0 {
		t.Errorf("nothing should be written, got %q", buf.String())
	}
}

func TestEventFilterOptions(t *testing.T) {
	defer func(h, b bool) { headless, blockedOnly = h, b }(headless, blockedOnly)
	own := uint32(os.Getpid()) //nolint:gosec // test pid
	for _, tc := range []struct {
		headless, blockedOnly, otherAllowed bool
	}{
		{false, false, true},
		{true, false, true},
		{false, true, true}, // --blocked-only only filters the headless writer
		{true, true, false},
	} {
		headless, blockedOnly = tc.headless, tc.blockedOnly
		g := &guard.Guard{}
		for _, opt := range eventFilterOptions() {
			opt(g)
		}
		ev := func(pid uint32, blocked bool) *guard.GuardEvent {
			return &guard.GuardEvent{FileEvent: ebpf.FileEvent{PID: pid}, Blocked: blocked}
		}
		if !g.Reports(ev(own, true)) || !g.Reports(ev(own+1, true)) {
			t.Fatalf("%+v: denials must always be reported", tc)
		}
		if g.Reports(ev(own, false)) {
			t.Fatalf("%+v: the daemon's own allowed I/O must not be queued", tc)
		}
		if got := g.Reports(ev(own+1, false)); got != tc.otherAllowed {
			t.Fatalf("%+v: another process's allowed event reported=%v", tc, got)
		}
	}
}
