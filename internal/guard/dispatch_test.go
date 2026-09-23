package guard

import "testing"

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
