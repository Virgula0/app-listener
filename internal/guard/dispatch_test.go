package guard

import (
	"testing"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

func TestDispatchFullChannelNeverDropsDenials(t *testing.T) {
	g := &Guard{path: "/r", events: make(chan GuardEvent), done: make(chan struct{})}

	g.dispatch(&GuardEvent{Blocked: true})
	if n := g.eventsDropped.Load(); n != 0 {
		t.Fatalf("a denial behind a full channel must be logged, not dropped: %d dropped", n)
	}
	g.dispatch(&GuardEvent{})
	if n := g.eventsDropped.Load(); n != 1 {
		t.Fatalf("allowed events may be dropped and counted: got %d", n)
	}
}

func TestDispatchWithoutAllowedEvents(t *testing.T) {
	g := &Guard{path: "/r", events: make(chan GuardEvent, 2), done: make(chan struct{}), dropAllowed: true}
	g.dispatch(&GuardEvent{})
	g.dispatch(&GuardEvent{Blocked: true})
	if len(g.events) != 1 || !(<-g.events).Blocked {
		t.Fatal("only the denial may be queued")
	}
	if n := g.eventsDropped.Load(); n != 0 {
		t.Fatalf("discarding unwanted allowed events is not a backlog: %d counted", n)
	}
}

func TestDeliverLabelsCrossResourceProcessGates(t *testing.T) {
	r1 := &Guard{path: "/r1", events: make(chan GuardEvent, 4), done: make(chan struct{})}
	r2 := &Guard{path: "/r2", events: make(chan GuardEvent, 4), done: make(chan struct{})}
	e := &engine{setAt: map[uint32]string{GuardMaxRes - 1: setKeyOf([]string{"/r1", "/r2"})}}
	e.slots[1], e.slots[2] = r1, r2

	e.deliver(resGlobal, &GuardEvent{Blocked: true, Process: "PTRACE"})
	if got := nextLabel(r1); got != MultipleResourceLabel {
		t.Fatalf("a GLOBAL-slot process gate must not be labelled with its reporter's path: %q", got)
	}
	e.deliver(GuardMaxRes-1, &GuardEvent{Blocked: true, Process: "PTRACE"})
	if got := nextLabel(r1); got != "/r1,/r2" {
		t.Fatalf("a taint set's process gate must name the set: %q", got)
	}
	e.deliver(resGlobal, &GuardEvent{Blocked: true, FsGate: RawDeviceResourceLabel})
	if got := nextLabel(r1); got != RawDeviceResourceLabel {
		t.Fatalf("a filesystem gate keeps its own label: %q", got)
	}
	e.deliver(2, &GuardEvent{Blocked: true, Process: "PTRACE"})
	if got := nextLabel(r2); got != "/r2" {
		t.Fatalf("a one-resource process gate is that resource's: %q", got)
	}
}

func TestSetInspectorsBounded(t *testing.T) {
	e := &engine{}
	var in []Inspector
	for i := 0; i <= maxInspectors; i++ {
		in = append(in, Inspector{Path: "/usr/bin/x", Key: GuardInodeKey{Dev: 1, Ino: uint64(i + 1)}})
	}
	if err := e.setInspectors(in, []string{"/usr/bin/x"}); err == nil {
		t.Fatal("more inspector inodes than the kernel map holds must be refused")
	}
	if err := e.setInspectors(in[:2], []string{"/usr/bin/x"}); err != nil || len(e.inspectors) != 2 {
		t.Fatalf("setInspectors: %v, %v", err, e.inspectors)
	}
	if !e.admitsExe(in[0].Key) {
		t.Fatal("an inspector key must count as admitted, or its superseded mark could be pruned")
	}
	if err := e.setInspectors(nil, nil); err != nil || len(e.inspectors) != 0 {
		t.Fatalf("an empty set must revoke every inspector: %v, %v", err, e.inspectors)
	}
}

func nextLabel(g *Guard) string {
	ev := <-g.events
	return ev.ResourceLabel(g.path)
}

func TestDispatchWithoutOwnAllowedEvents(t *testing.T) {
	g := &Guard{path: "/etc/app-listener", events: make(chan GuardEvent, 4), done: make(chan struct{})}
	WithoutOwnAllowedEvents()(g)
	own := g.quietPID
	if own == 0 {
		t.Fatal("the option must record this process's pid")
	}
	g.dispatch(&GuardEvent{FileEvent: ebpf.FileEvent{PID: own}})
	g.dispatch(&GuardEvent{FileEvent: ebpf.FileEvent{PID: own}, Blocked: true})
	g.dispatch(&GuardEvent{FileEvent: ebpf.FileEvent{PID: own + 1}})
	if len(g.events) != 2 {
		t.Fatalf("only the daemon's own allowed I/O may be discarded: %d queued", len(g.events))
	}
	if ev := <-g.events; !ev.Blocked || ev.PID != own {
		t.Fatalf("the daemon's own denial must be queued: %+v", ev)
	}
	if ev := <-g.events; ev.Blocked || ev.PID != own+1 {
		t.Fatalf("another process's allowed access must be queued: %+v", ev)
	}
	if n := g.eventsDropped.Load(); n != 0 {
		t.Fatalf("discarding own I/O is not a backlog: %d counted", n)
	}
}
