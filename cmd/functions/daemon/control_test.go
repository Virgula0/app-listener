package daemon

import (
	"bufio"
	"encoding/json"
	"net"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Virgula0/app-listener/cmd/functions/editprotected"
	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
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

func TestReadGrantRequestSelect(t *testing.T) {
	req, err := readGrantRequest(bufio.NewReader(strings.NewReader("SELECT /home/alice/.ssh\n")))
	if err != nil {
		t.Fatal(err)
	}
	if req.resource != "/home/alice/.ssh" || req.forward != nil {
		t.Fatalf("parsed %+v", req)
	}
	for _, bad := range []string{"PICK /x\n", "SELECT\n", "SELECT \n", "SELECT /x"} {
		if _, err := readGrantRequest(bufio.NewReader(strings.NewReader(bad))); err == nil {
			t.Errorf("readGrantRequest(%q) should error", bad)
		}
	}
}

func TestReadGrantRequestForward(t *testing.T) {
	in := "FORWARD allow OPEN,READ 2 1\n/home/a/.claude\n/home/a/.ssh\n/home/a/bin/my tool\n"
	req, err := readGrantRequest(bufio.NewReader(strings.NewReader(in)))
	if err != nil {
		t.Fatal(err)
	}
	fwd := req.forward
	if fwd == nil || !fwd.rule.Allow {
		t.Fatalf("parsed %+v", req)
	}
	if !slices.Equal(fwd.rule.Events, []ebpf.EventType{ebpf.EventOpen, ebpf.EventRead}) {
		t.Errorf("events = %v", fwd.rule.Events)
	}
	if !slices.Equal(fwd.resources, []string{"/home/a/.claude", "/home/a/.ssh"}) ||
		!slices.Equal(fwd.binaries, []string{"/home/a/bin/my tool"}) {
		t.Errorf("paths = %v / %v", fwd.resources, fwd.binaries)
	}

	req, err = readGrantRequest(bufio.NewReader(strings.NewReader("FORWARD block * 1 1\n/r\n/b\n")))
	if err != nil {
		t.Fatal(err)
	}
	if req.forward.rule.Allow || req.forward.rule.Events != nil {
		t.Errorf("block * parsed as %+v", req.forward.rule)
	}

	for _, bad := range []string{
		"FORWARD allow * 1 1\n/r\n",        // a binary line missing
		"FORWARD allow * 1 1\nr\n/b\n",     // relative path
		"FORWARD permit * 1 1\n/r\n/b\n",   // unknown action
		"FORWARD allow EXEC 1 1\n/r\n/b\n", // unknown event
		"FORWARD allow * 0 1\n/b\n",        // no resource
		"FORWARD allow * 1 0\n/r\n",        // no binary
		"FORWARD allow * 1 33\n/r\n",       // over the binary cap
		"FORWARD allow * 65 1\n",           // over the resource cap
		"FORWARD allow *\n",                // short header
	} {
		if _, err := readGrantRequest(bufio.NewReader(strings.NewReader(bad))); err == nil {
			t.Errorf("readGrantRequest(%q) should error", bad)
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

func TestReadGrantRequestTimeout(t *testing.T) {
	req, err := readGrantRequest(bufio.NewReader(strings.NewReader("SELECT /r\n")))
	if err != nil || req.timeout != editprotected.SessionTimeoutDefault {
		t.Fatalf("default timeout = %v, %v", req, err)
	}
	req, err = readGrantRequest(bufio.NewReader(strings.NewReader("TIMEOUT 45\nFORWARD allow * 1 1\n/r\n/b\n")))
	if err != nil || req.timeout != 45*time.Second || req.forward == nil {
		t.Fatalf("TIMEOUT 45 parsed as %+v, %v", req, err)
	}
	for _, bad := range []string{"TIMEOUT 0\nSELECT /r\n", "TIMEOUT -5\nSELECT /r\n", "TIMEOUT 10m\nSELECT /r\n",
		"TIMEOUT 86401\nSELECT /r\n", "TIMEOUT 5\n"} {
		if _, err := readGrantRequest(bufio.NewReader(strings.NewReader(bad))); err == nil {
			t.Errorf("readGrantRequest(%q) should error", bad)
		}
	}
}

func newTestSession(timeout time.Duration, events chan guard.GuardEvent) (*editControlSession, *atomic.Int32) {
	revokes := &atomic.Int32{}
	return &editControlSession{
		resource: "/r",
		revoke:   func() error { revokes.Add(1); return nil },
		done:     make(chan struct{}),
		timeout:  timeout,
		activity: make(chan struct{}, 1),
		events:   events,
	}, revokes
}

func TestRunSessionIdleTimeoutResetsOnPing(t *testing.T) {
	server, client := net.Pipe()
	defer client.Close()
	session, revokes := newTestSession(150*time.Millisecond, nil)
	finished := make(chan struct{})
	go func() {
		(&controlServer{}).runSession(server, bufio.NewReader(server), session)
		close(finished)
	}()

	// Pinging every 50ms keeps a 150ms idle timeout alive well past 150ms.
	for range 8 {
		if _, err := client.Write([]byte("PING\n")); err != nil {
			t.Fatal(err)
		}
		time.Sleep(50 * time.Millisecond)
	}
	if revokes.Load() != 0 {
		t.Fatal("an active session hit the idle timeout")
	}
	select {
	case <-finished:
	case <-time.After(2 * time.Second):
		t.Fatal("an idle session was never revoked")
	}
	if revokes.Load() != 1 {
		t.Fatalf("revoked %d times, want 1", revokes.Load())
	}
}

func TestRunSessionStreamsEventsAsActivity(t *testing.T) {
	server, client := net.Pipe()
	defer client.Close()
	events := make(chan guard.GuardEvent, 4)
	session, _ := newTestSession(time.Hour, events)
	go (&controlServer{}).runSession(server, bufio.NewReader(server), session)

	events <- guard.GuardEvent{FileEvent: ebpf.FileEvent{PID: 42, Path: "/r/secret", Type: ebpf.EventRead}, Blocked: true}
	line, err := bufio.NewReader(client).ReadString('\n')
	if err != nil {
		t.Fatal(err)
	}
	payload, ok := strings.CutPrefix(strings.TrimSpace(line), "EVENT ")
	var ev guard.GuardEvent
	if !ok || json.Unmarshal([]byte(payload), &ev) != nil || ev.PID != 42 || ev.Path != "/r/secret" || !ev.Blocked {
		t.Fatalf("streamed %q", line)
	}
	if _, err := client.Write([]byte("END\n")); err != nil {
		t.Fatal(err)
	}
	select {
	case <-session.done:
	case <-time.After(2 * time.Second):
		t.Fatal("END did not end the session")
	}
}
