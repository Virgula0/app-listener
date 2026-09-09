package daemon

import (
	"bufio"
	"strings"
	"testing"
	"time"
)

func TestReadHandshake(t *testing.T) {
	ok := "BEGIN /home/alice/.ssh\nhunter2hunter2\n"
	res, pw, err := readHandshake(bufio.NewReader(strings.NewReader(ok)))
	if err != nil {
		t.Fatal(err)
	}
	if res != "/home/alice/.ssh" || pw != "hunter2hunter2" {
		t.Fatalf("parsed %q / %q", res, pw)
	}

	for _, bad := range []string{
		"HELLO /x\npw\n",
		"BEGIN\npw\n",
		"BEGIN /x\n\n",
		"BEGIN /x\n",
	} {
		if _, _, err := readHandshake(bufio.NewReader(strings.NewReader(bad))); err == nil {
			t.Errorf("readHandshake(%q) should error", bad)
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
