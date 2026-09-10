package daemon

import (
	"bufio"
	"strings"
	"testing"
	"time"
)

func TestReadAuthRequest(t *testing.T) {
	pw, err := readAuthRequest(bufio.NewReader(strings.NewReader("AUTH\nhunter2hunter2\n")))
	if err != nil {
		t.Fatal(err)
	}
	if pw != "hunter2hunter2" {
		t.Fatalf("parsed %q", pw)
	}
	for _, bad := range []string{"HELLO\npw\n", "AUTH\n\n", "AUTH\n", "AUTH pw\n"} {
		if _, err := readAuthRequest(bufio.NewReader(strings.NewReader(bad))); err == nil {
			t.Errorf("readAuthRequest(%q) should error", bad)
		}
	}
}

func TestReadSelectRequest(t *testing.T) {
	res, err := readSelectRequest(bufio.NewReader(strings.NewReader("SELECT /home/alice/.ssh\n")))
	if err != nil {
		t.Fatal(err)
	}
	if res != "/home/alice/.ssh" {
		t.Fatalf("parsed %q", res)
	}
	for _, bad := range []string{"PICK /x\n", "SELECT\n", "SELECT \n", "SELECT /x"} {
		if _, err := readSelectRequest(bufio.NewReader(strings.NewReader(bad))); err == nil {
			t.Errorf("readSelectRequest(%q) should error", bad)
		}
	}
}

func TestRecordFailureLockout(t *testing.T) {
	cs := &controlServer{}
	for i := 0; i < maxAuthFailures-1; i++ {
		cs.recordFailure()
	}
	if !cs.lockUntil.IsZero() {
		t.Fatal("locked too early")
	}
	cs.recordFailure() // trips the lockout
	if cs.lockUntil.IsZero() || time.Until(cs.lockUntil) <= 0 {
		t.Fatalf("expected an active lockout, got %v", cs.lockUntil)
	}
	if cs.failures != 0 {
		t.Fatalf("failure counter should reset after lockout, got %d", cs.failures)
	}
}

func TestEditControlSessionEndIdempotent(t *testing.T) {
	revokes := 0
	s := &editControlSession{
		resource: "/x",
		revoke:   func() error { revokes++; return nil },
		done:     make(chan struct{}),
	}
	s.end()
	s.end()
	if revokes != 1 {
		t.Fatalf("revoke called %d times, want 1", revokes)
	}
	select {
	case <-s.done:
	default:
		t.Fatal("done channel not closed")
	}
}
