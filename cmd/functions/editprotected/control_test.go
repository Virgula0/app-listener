package editprotected

import (
	"bufio"
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
