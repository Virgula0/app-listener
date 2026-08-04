package integrationtests

import (
	"fmt"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

type netMonitorEvent struct {
	Type     string
	Comm     string
	Protocol string
	SrcAddr  string
	DstAddr  string
}

func parseNetMonitorEvents(logContent string) []netMonitorEvent {
	var events []netMonitorEvent
	for _, line := range strings.Split(logContent, "\n") {
		idx := strings.Index(line, "NETEVENT|")
		if idx < 0 {
			continue
		}
		rest := line[idx+len("NETEVENT|"):]
		parts := strings.SplitN(rest, "|", 10)
		if len(parts) < 10 {
			continue
		}
		events = append(events, netMonitorEvent{
			Type:     strings.TrimSpace(parts[0]),
			Comm:     strings.TrimSpace(parts[1]),
			Protocol: strings.TrimSpace(parts[2]),
			SrcAddr:  strings.TrimSpace(parts[3]),
			DstAddr:  strings.TrimSpace(parts[4]),
		})
	}
	return events
}

func netMonitorDeltaEvents(logBefore, logAfter string) []netMonitorEvent {
	evBefore := parseNetMonitorEvents(logBefore)
	evAfter := parseNetMonitorEvents(logAfter)
	if len(evAfter) > len(evBefore) {
		return evAfter[len(evBefore):]
	}
	return nil
}

func fmtEvents(events []netMonitorEvent) string {
	types := make([]string, 0, len(events))
	for _, e := range events {
		types = append(types, e.Type)
	}
	return strings.Join(types, ",")
}

func netMonitorHasTypes(events []netMonitorEvent, needed []string) bool {
	seen := map[string]bool{}
	for _, e := range events {
		seen[e.Type] = true
	}
	for _, t := range needed {
		if !seen[t] {
			return false
		}
	}
	return true
}

// requireNetMonitorTypes fails unless at least one event of each expected type
// was emitted.
func (s *IntegrationSuite) requireNetMonitorTypes(events []netMonitorEvent, expected ...string) {
	required := make(map[string]bool, len(expected))
	for _, t := range expected {
		required[t] = true
	}
	seen := make(map[string]bool)
	for _, e := range events {
		seen[e.Type] = true
	}
	for t := range required {
		s.Require().Truef(seen[t], "expected event %s not captured; captured: %s", t, fmtEvents(events))
	}
}

// requireOnlyNetMonitorTypes asserts that no event type outside the given set
// was emitted (catches unexpected/bypassing events).
func (s *IntegrationSuite) requireOnlyNetMonitorTypes(events []netMonitorEvent, allowed ...string) {
	allowedSet := make(map[string]bool, len(allowed))
	for _, t := range allowed {
		allowedSet[t] = true
	}
	for _, e := range events {
		s.Require().Truef(allowedSet[e.Type],
			"unexpected event type %s captured (allowed: %v): %v", e.Type, allowed, fmtEvents(events))
	}
}

// waitForNetMonitorTypes polls the monitor log until all expected types appear
// (as delta from logBefore), returning the delta events.
func (s *IntegrationSuite) waitForNetMonitorTypes(c testcontainers.Container, logBefore string, needed []string, timeout time.Duration) []netMonitorEvent {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		events := netMonitorDeltaEvents(logBefore, s.readNetMonitorLog(c))
		if netMonitorHasTypes(events, needed) {
			return events
		}
		time.Sleep(250 * time.Millisecond)
	}
	events := netMonitorDeltaEvents(logBefore, s.readNetMonitorLog(c))
	s.T().Fatalf("timed out waiting for event types %v; captured: %v", needed, fmtEvents(events))
	return nil
}

func (s *IntegrationSuite) startNetworkMonitorStd(c testcontainers.Container, binaries string, extraFlags ...string) {
	// Tracepoints require tracefs to be mounted inside the container.
	_, _ = s.exec(c, []string{"sh", "-c", "mkdir -p /sys/kernel/tracing && mount -t tracefs tracefs /sys/kernel/tracing 2>/dev/null; true"})

	flags := strings.Join(extraFlags, " ")
	cmd := fmt.Sprintf("nohup /app-listener network-monitor %s --headless %s > /tmp/monitor.log 2>&1 &", binaries, flags)
	s.exec(c, []string{"sh", "-c", cmd})

	time.Sleep(3 * time.Second)

	codeCheck, outCheck := s.exec(c, []string{"pgrep", "-f", "app-listener"})
	s.Require().Equalf(0, codeCheck, "network-monitor process not running after start: %s", outCheck)

	// All 11 probes/tracepoints must attach (the kretprobe attaches even without
	// tracefs, so a lower count means tracepoints were silently skipped).
	startup := s.readNetMonitorLog(c)
	s.Require().Containsf(startup, "11/11 probes/tracepoints attached",
		"not all network-monitor probes attached:\n%s", startup)
}

func (s *IntegrationSuite) readNetMonitorLog(c testcontainers.Container) string {
	_, out := s.exec(c, []string{"sh", "-c", "cat /tmp/monitor.log 2>/dev/null || true"})
	return out
}

func (s *IntegrationSuite) stopNetMonitor(c testcontainers.Container) {
	s.exec(c, []string{"pkill", "-f", "app-listener"})
	time.Sleep(500 * time.Millisecond)
}

// newNetMonitorContainer starts a privileged container with the monitor
// binary, copies net_tester into it and starts the monitor watching it.
func (s *IntegrationSuite) newNetMonitorContainer(extraBinaries []string, extraFlags ...string) (testcontainers.Container, string) {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	s.Require().NoError(c.CopyFileToContainer(s.ctx, netTesterAmd64Bin, "/net_tester", 0755))
	for _, b := range extraBinaries {
		s.Require().NoError(c.CopyFileToContainer(s.ctx, netTesterAmd64Bin, b, 0755))
	}
	s.startNetworkMonitorStd(c, "/net_tester", extraFlags...)
	return c, s.readNetMonitorLog(c)
}

func (s *IntegrationSuite) TestNetworkMonitor_BindListen_TCP() {
	c, logBefore := s.newNetMonitorContainer(nil)
	defer c.Terminate(s.ctx)
	defer s.stopNetMonitor(c)

	s.exec(c, []string{"sh", "-c", "/net_tester tcp-server 8080 > /tmp/server.log 2>&1 &"})
	time.Sleep(1 * time.Second)

	events := s.waitForNetMonitorTypes(c, logBefore, []string{"BIND", "LISTEN"}, 8*time.Second)
	s.requireNetMonitorTypes(events, "BIND", "LISTEN")
}

func (s *IntegrationSuite) TestNetworkMonitor_Accept_TCP() {
	c, logBefore := s.newNetMonitorContainer(nil)
	defer c.Terminate(s.ctx)
	defer s.stopNetMonitor(c)

	s.exec(c, []string{"sh", "-c", "/net_tester tcp-server 8080 > /tmp/server.log 2>&1 &"})
	time.Sleep(1 * time.Second)
	s.exec(c, []string{"sh", "-c", "/net_tester tcp-client 127.0.0.1 8080"})

	events := s.waitForNetMonitorTypes(c, logBefore, []string{"ACCEPT"}, 8*time.Second)
	s.requireNetMonitorTypes(events, "ACCEPT")
}

func (s *IntegrationSuite) TestNetworkMonitor_Connect_TCP() {
	c, logBefore := s.newNetMonitorContainer(nil)
	defer c.Terminate(s.ctx)
	defer s.stopNetMonitor(c)

	// Connect to a closed port: connect(2) still fires the tracepoint.
	s.exec(c, []string{"sh", "-c", "/net_tester tcp-client 127.0.0.1 9999; true"})

	events := s.waitForNetMonitorTypes(c, logBefore, []string{"CONNECT"}, 8*time.Second)
	s.requireNetMonitorTypes(events, "CONNECT")
}

func (s *IntegrationSuite) TestNetworkMonitor_SendRecv_TCP() {
	c, logBefore := s.newNetMonitorContainer(nil)
	defer c.Terminate(s.ctx)
	defer s.stopNetMonitor(c)

	s.exec(c, []string{"sh", "-c", "/net_tester tcp-server 8080 > /tmp/server.log 2>&1 &"})
	time.Sleep(1 * time.Second)
	s.exec(c, []string{"sh", "-c", "/net_tester tcp-client 127.0.0.1 8080"})

	events := s.waitForNetMonitorTypes(c, logBefore, []string{"SEND", "RECV"}, 8*time.Second)
	s.requireNetMonitorTypes(events, "SEND", "RECV")
}

func (s *IntegrationSuite) TestNetworkMonitor_Close_TCP() {
	c, logBefore := s.newNetMonitorContainer(nil)
	defer c.Terminate(s.ctx)
	defer s.stopNetMonitor(c)

	s.exec(c, []string{"sh", "-c", "/net_tester tcp-server 8080 > /tmp/server.log 2>&1 &"})
	time.Sleep(1 * time.Second)
	s.exec(c, []string{"sh", "-c", "/net_tester tcp-client 127.0.0.1 8080"})

	events := s.waitForNetMonitorTypes(c, logBefore, []string{"CLOSE"}, 8*time.Second)
	s.requireNetMonitorTypes(events, "CLOSE")
}

func (s *IntegrationSuite) TestNetworkMonitor_Bind_UDP() {
	c, logBefore := s.newNetMonitorContainer(nil)
	defer c.Terminate(s.ctx)
	defer s.stopNetMonitor(c)

	s.exec(c, []string{"sh", "-c", "/net_tester udp-server 8081 > /tmp/server.log 2>&1 &"})
	time.Sleep(1 * time.Second)

	events := s.waitForNetMonitorTypes(c, logBefore, []string{"BIND"}, 8*time.Second)
	s.requireNetMonitorTypes(events, "BIND")
}

func (s *IntegrationSuite) TestNetworkMonitor_Connect_UDP() {
	c, logBefore := s.newNetMonitorContainer(nil)
	defer c.Terminate(s.ctx)
	defer s.stopNetMonitor(c)

	// UDP connect(2) always succeeds, even without a listener.
	s.exec(c, []string{"sh", "-c", "/net_tester udp-client 127.0.0.1 9999; true"})

	events := s.waitForNetMonitorTypes(c, logBefore, []string{"CONNECT"}, 8*time.Second)
	s.requireNetMonitorTypes(events, "CONNECT")
}

func (s *IntegrationSuite) TestNetworkMonitor_SendRecv_UDP() {
	c, logBefore := s.newNetMonitorContainer(nil)
	defer c.Terminate(s.ctx)
	defer s.stopNetMonitor(c)

	s.exec(c, []string{"sh", "-c", "/net_tester udp-server 8081 > /tmp/server.log 2>&1 &"})
	time.Sleep(1 * time.Second)
	s.exec(c, []string{"sh", "-c", "/net_tester udp-client 127.0.0.1 8081"})

	events := s.waitForNetMonitorTypes(c, logBefore, []string{"SEND", "RECV"}, 8*time.Second)
	s.requireNetMonitorTypes(events, "SEND", "RECV")
}

func (s *IntegrationSuite) TestNetworkMonitor_Close_UDP() {
	c, logBefore := s.newNetMonitorContainer(nil)
	defer c.Terminate(s.ctx)
	defer s.stopNetMonitor(c)

	s.exec(c, []string{"sh", "-c", "/net_tester udp-server 8081 > /tmp/server.log 2>&1 &"})
	time.Sleep(1 * time.Second)
	s.exec(c, []string{"sh", "-c", "/net_tester udp-client 127.0.0.1 8081"})

	events := s.waitForNetMonitorTypes(c, logBefore, []string{"CLOSE"}, 8*time.Second)
	s.requireNetMonitorTypes(events, "CLOSE")
}

func (s *IntegrationSuite) TestNetworkMonitor_DNS() {
	c, logBefore := s.newNetMonitorContainer(nil)
	defer c.Terminate(s.ctx)
	defer s.stopNetMonitor(c)

	s.exec(c, []string{"sh", "-c", "/net_tester dns example.com; true"})

	events := s.waitForNetMonitorTypes(c, logBefore, []string{"DNS"}, 8*time.Second)
	s.requireNetMonitorTypes(events, "DNS")
}

func (s *IntegrationSuite) TestNetworkMonitor_EventFilter() {
	c, logBefore := s.newNetMonitorContainer(nil, "-e", "CONNECT,SEND")
	defer c.Terminate(s.ctx)
	defer s.stopNetMonitor(c)

	s.exec(c, []string{"sh", "-c", "/net_tester tcp-server 8082 > /tmp/server.log 2>&1 &"})
	time.Sleep(1 * time.Second)
	s.exec(c, []string{"sh", "-c", "/net_tester tcp-client 127.0.0.1 8082"})

	events := s.waitForNetMonitorTypes(c, logBefore, []string{"CONNECT", "SEND"}, 8*time.Second)
	s.requireNetMonitorTypes(events, "CONNECT", "SEND")
	s.requireOnlyNetMonitorTypes(events, "CONNECT", "SEND", "DNS")
}

func (s *IntegrationSuite) TestNetworkMonitor_Identity() {
	c, logBefore := s.newNetMonitorContainer([]string{"/other_tester"})
	defer c.Terminate(s.ctx)
	defer s.stopNetMonitor(c)

	// Watch only /net_tester: a different binary must not produce events.
	s.exec(c, []string{"sh", "-c", "/other_tester tcp-server 8083 > /tmp/server.log 2>&1 &"})
	time.Sleep(1 * time.Second)
	s.exec(c, []string{"sh", "-c", "/other_tester tcp-client 127.0.0.1 8083"})
	time.Sleep(2 * time.Second)

	delta := netMonitorDeltaEvents(logBefore, s.readNetMonitorLog(c))
	s.Require().Empty(delta, "non-watched binary must not produce events, got: %v", fmtEvents(delta))
}

func (s *IntegrationSuite) TestNetworkMonitor_InvalidEventType() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.Require().NoError(c.CopyFileToContainer(s.ctx, netTesterAmd64Bin, "/net_tester", 0755))

	// An unknown event type must be rejected before anything attaches.
	code, out := s.exec(c, []string{"sh", "-c", "timeout 10 /app-listener network-monitor /net_tester -e BOGUS --headless 2>&1"})
	s.Require().NotEqualf(0, code, "invalid event type should make network-monitor fail, got: %s", out)
	s.Require().Contains(out, "unknown network event type")
}
