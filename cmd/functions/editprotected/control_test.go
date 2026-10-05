package editprotected

import (
	"bufio"
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

func TestLiveSessionWatchAndPinger(t *testing.T) {
	client, server := net.Pipe()
	s := &liveSession{conn: client, r: bufio.NewReader(client), timeout: 10 * time.Second}
	if got := s.timeoutLine(); got != "TIMEOUT 10\n" {
		t.Fatalf("timeoutLine = %q", got)
	}
	events, ended := s.Watch()

	go func() {
		_, _ = server.Write([]byte(`EVENT {"PID":7,"Path":"/r/f","Blocked":true}` + "\nEVENT not-json\n"))
	}()
	select {
	case ev := <-events:
		if ev.PID != 7 || ev.Path != "/r/f" || !ev.Blocked {
			t.Fatalf("decoded %+v", ev)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("no event decoded")
	}

	ping := s.Pinger(time.Hour)
	got := make(chan string, 1)
	go func() {
		line, _ := bufio.NewReader(server).ReadString('\n')
		got <- line
	}()
	ping()
	ping() // throttled: one PING per interval
	if line := <-got; strings.TrimSpace(line) != "PING" {
		t.Fatalf("pinger wrote %q", line)
	}

	server.Close()
	select {
	case <-ended:
	case <-time.After(2 * time.Second):
		t.Fatal("a closed connection did not end the watch")
	}
}

func TestLiveSessionConfigAndPut(t *testing.T) {
	client, server := net.Pipe()
	s := &liveSession{conn: client, r: bufio.NewReader(client)}
	sr := bufio.NewReader(server)
	go func() {
		line, _ := sr.ReadString('\n')
		if line != "CONFIG\n" {
			_, _ = server.Write([]byte("ERR bad request\n"))
			return
		}
		_, _ = server.Write([]byte("OK 5\nhello"))
		put, _ := sr.ReadString('\n')
		body := make([]byte, 3)
		_, _ = io.ReadFull(sr, body)
		if put == "PUT 3\n" && string(body) == "bye" {
			_, _ = server.Write([]byte("ERR the reload failed: boom — restored\n"))
		}
	}()
	got, err := s.Config()
	if err != nil || string(got) != "hello" {
		t.Fatalf("Config = %q, %v", got, err)
	}
	if err := s.PutConfig([]byte("bye")); err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("PutConfig err = %v, want the daemon's refusal", err)
	}
}

func TestValidateFlagsEditConfig(t *testing.T) {
	t.Cleanup(func() { editConfigFlag, forwardFlag, whitelistFlags, resourceFlags = false, false, nil, nil })
	timeoutSession = SessionTimeoutDefault
	editConfigFlag = true
	if err := validateFlags(); err != nil {
		t.Fatalf("--edit-config alone: %v", err)
	}
	forwardFlag, whitelistFlags = true, []string{"/bin/x"}
	if err := validateFlags(); err == nil {
		t.Fatal("--edit-config --forward accepted")
	}
	forwardFlag, whitelistFlags, resourceFlags = false, nil, []string{"/r"}
	if err := validateFlags(); err == nil {
		t.Fatal("--edit-config --resource accepted")
	}
}
