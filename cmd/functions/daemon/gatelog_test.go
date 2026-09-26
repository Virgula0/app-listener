package daemon

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/usecase"
)

func gateEvent(process, path string) usecase.DaemonEvent {
	return usecase.DaemonEvent{
		Resource: "/home/alice/.steam/registry.vdf",
		Event: guard.GuardEvent{
			FileEvent: ebpf.FileEvent{PID: 1078, Comm: "Hyprland", Path: path},
			Blocked:   true,
			Process:   process,
		},
	}
}

func testLimiter() (*gateLogLimiter, *time.Time) {
	now := time.Date(2026, 9, 19, 18, 0, 0, 0, time.UTC)
	l := newGateLogLimiter()
	l.now = func() time.Time { return now }
	return l, &now
}

// TestGateLogLimiterFoldsReadRepeats: the first metadata-only denial of a
// caller/target pair is logged, repeats are counted, and one summary line
// reports them once the window has elapsed.
func TestGateLogLimiterFoldsReadRepeats(t *testing.T) {
	l, now := testLimiter()
	ev := gateEvent("PTRACE", "pid=23837 comm=CrBrowserMain mode=READ")

	if !l.Admit(&ev) {
		t.Fatal("the first occurrence must be logged in full")
	}
	for i := 0; i < 5; i++ {
		if l.Admit(&ev) {
			t.Fatalf("repeat %d must be folded", i)
		}
	}

	var buf bytes.Buffer
	l.Flush(&buf, false)
	if buf.Len() != 0 {
		t.Fatalf("no summary before the window elapses, got %q", buf.String())
	}

	*now = now.Add(gateLogWindow)
	l.Flush(&buf, false)
	out := buf.String()
	if !strings.Contains(out, "DAEMON DENIED-REPEAT  op=PTRACE") || !strings.Contains(out, "repeats=5") {
		t.Fatalf("expected one summary with repeats=5, got %q", out)
	}
	if strings.Contains(out, "DAEMON DENIED  op=") {
		t.Errorf("a summary must not look like a single denial to parsers: %q", out)
	}
	if !l.Admit(&ev) {
		t.Error("after the window the next occurrence must be logged in full again")
	}
}

// TestGateLogLimiterNeverFoldsWhatMatters: memory access, /proc/<pid>/mem,
// traced exec and file denials are always logged in full, every time.
func TestGateLogLimiterNeverFoldsWhatMatters(t *testing.T) {
	l, _ := testLimiter()
	fileDenial := gateEvent("", "/home/alice/.steam/registry.vdf")
	for _, ev := range []usecase.DaemonEvent{
		gateEvent("PTRACE", "pid=1 comm=steam mode=ATTACH"),
		gateEvent("PROC_MEM", "pid=1 comm=steam mode=ATTACH"),
		gateEvent("TRACED_EXEC", "pid=1 comm=gdb"),
		fileDenial,
	} {
		for i := 0; i < 3; i++ {
			if !l.Admit(&ev) {
				t.Fatalf("%s %q must never be folded (occurrence %d)", ev.Event.Process, ev.Event.Path, i)
			}
		}
	}
}

// TestGateLogLimiterKeysOnTarget: a different target PROGRAM is its own first
// occurrence, logged in full.
func TestGateLogLimiterKeysOnTarget(t *testing.T) {
	l, _ := testLimiter()
	a := gateEvent("PTRACE", "pid=10 comm=steam mode=READ")
	b := gateEvent("PTRACE", "pid=11 comm=explorer.exe mode=READ")
	if !l.Admit(&a) || !l.Admit(&b) {
		t.Fatal("distinct target programs must each be logged the first time")
	}
}

// TestGateLogLimiterFoldsAcrossThreadsAndPids reproduces the field log: the
// portal inspects from worker threads with different names (pool-37,
// pool-38), and Steam restarts processes with new pids. Both are repeats of
// one caller looking at one program.
func TestGateLogLimiterFoldsAcrossThreadsAndPids(t *testing.T) {
	l, _ := testLimiter()
	first := gateEvent("PTRACE", "pid=46804 comm=steamwebhelper mode=READ")
	first.Event.PID, first.Event.Comm = 1465, "pool-37"
	otherThread := first
	otherThread.Event.Comm = "pool-38"
	otherPid := first
	otherPid.Event.Path = "pid=47111 comm=steamwebhelper mode=READ"

	if !l.Admit(&first) {
		t.Fatal("the first occurrence must be logged")
	}
	if l.Admit(&otherThread) {
		t.Error("another worker thread of the same caller process is a repeat")
	}
	if l.Admit(&otherPid) {
		t.Error("the same program under a new pid is a repeat")
	}
}

// TestGateLogLimiterFlushesOnShutdown: forced flush reports pending repeats.
func TestGateLogLimiterFlushesOnShutdown(t *testing.T) {
	l, _ := testLimiter()
	ev := gateEvent("PTRACE", "pid=10 comm=steam mode=READ")
	l.Admit(&ev)
	l.Admit(&ev)
	var buf bytes.Buffer
	l.Flush(&buf, true)
	if !strings.Contains(buf.String(), "repeats=1") {
		t.Fatalf("forced flush must report pending repeats, got %q", buf.String())
	}
}

// TestNoLogMetadataBlocks: with --no-log-metadata-blocks a metadata-only
// process-gate denial is never written — not even the first occurrence or a
// summary — while memory, process and file denials still always are.
func TestNoLogMetadataBlocks(t *testing.T) {
	l, now := testLimiter()
	meta := gateEvent("PTRACE", "pid=10 comm=steam mode=READ")
	for i := 0; i < 3; i++ {
		if admitEvent(l, &meta, true) {
			t.Fatalf("metadata-only denial %d must not be logged", i)
		}
	}
	*now = now.Add(gateLogWindow)
	var buf bytes.Buffer
	l.Flush(&buf, true)
	if buf.Len() != 0 {
		t.Errorf("no summary may be written for dropped metadata denials, got %q", buf.String())
	}
	for _, ev := range []usecase.DaemonEvent{
		gateEvent("PTRACE", "pid=1 comm=steam mode=ATTACH"),
		gateEvent("PROC_MEM", "pid=1 comm=steam mode=ATTACH"),
		gateEvent("TRACED_EXEC", "pid=1 comm=gdb"),
		gateEvent("", "/home/alice/.steam/registry.vdf"),
	} {
		if !admitEvent(l, &ev, true) {
			t.Errorf("%s %q must still be logged", ev.Event.Process, ev.Event.Path)
		}
	}
	// Without the flag the limiter's behaviour is unchanged.
	fresh, _ := testLimiter()
	if !admitEvent(fresh, &meta, false) {
		t.Error("without the flag the first metadata denial is logged")
	}
}
