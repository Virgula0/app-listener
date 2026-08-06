package daemon

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/usecase"
)

func deniedEvent() usecase.DaemonEvent {
	return usecase.DaemonEvent{
		Resource: "/home/alice/.ssh",
		Event: guard.GuardEvent{FileEvent: ebpf.FileEvent{
			PID: 1234, UID: 1000, Comm: "ssh",
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
	for _, only := range []bool{false, true} {
		var buf bytes.Buffer
		ev := deniedEvent()
		if !writeEvent(&buf, only, &ev) {
			t.Fatalf("denied event must always be written (blockedOnly=%v)", only)
		}
		line := buf.String()
		if !strings.HasPrefix(line, "<4>DAEMON DENIED|") {
			t.Errorf("line %q: missing syslog warning marker", line)
		}
		for _, want := range []string{
			"/home/alice/.ssh",
			"OPEN",
			"ssh",
			"/home/alice/.ssh/authorized_keys",
			"1234",
			"1000",
		} {
			if !strings.Contains(line, want) {
				t.Errorf("line %q: missing field %q", line, want)
			}
		}
	}
}

func TestWriteEventAllowed(t *testing.T) {
	var buf bytes.Buffer
	ev := allowedEvent()
	if !writeEvent(&buf, false, &ev) {
		t.Fatal("allowed event must be written when blockedOnly is false")
	}
	if line := buf.String(); !strings.HasPrefix(line, "<6>DAEMON|") {
		t.Errorf("line %q: missing syslog info prefix", line)
	}
}

func TestWriteEventAllowedSuppressed(t *testing.T) {
	var buf bytes.Buffer
	ev := allowedEvent()
	if writeEvent(&buf, true, &ev) {
		t.Fatal("allowed event must be dropped under --blocked-only")
	}
	if buf.Len() != 0 {
		t.Errorf("nothing should be written, got %q", buf.String())
	}
}
