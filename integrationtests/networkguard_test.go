package integrationtests

import (
	"fmt"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

// Network-guard state for the pooled container (see the guard helpers for
// the PID-scoped lifecycle rationale).
var (
	netGuardPID     int
	netGuardLogPath string
)

func (s *IntegrationSuite) startNetworkGuardStd(c testcontainers.Container, extraFlags ...string) {
	flags := strings.Join(extraFlags, " ")
	netGuardLogPath = fmt.Sprintf("/tmp/netguard-%d.log", time.Now().UnixNano())
	cmd := fmt.Sprintf("nohup /app-listener network-guard %s --headless > %s 2>&1 &", flags, netGuardLogPath)
	netGuardPID = s.capturePID(c, cmd)

	// Wait until the guard finished attaching its LSM hooks (not just a fixed
	// sleep: under host load attach can take a few seconds, and tests that
	// rely on immediate blocking would race with it).
	deadline := time.Now().Add(20 * time.Second)
	for {
		if strings.Contains(s.readNetGuardLog(c), "network guard started") {
			break
		}
		if time.Now().After(deadline) {
			s.T().Fatalf("network-guard did not finish attaching in time; log tail:\n%s",
				netGuardTail(s.readNetGuardLog(c)))
			break
		}
		time.Sleep(250 * time.Millisecond)
	}

	codeCheck, outCheck := s.exec(c, []string{"sh", "-c",
		fmt.Sprintf("kill -0 %d 2>/dev/null && echo alive || echo dead", netGuardPID)})
	s.Require().Equalf(0, codeCheck, "network-guard process not running after start: %s", outCheck)
}

func (s *IntegrationSuite) readNetGuardLog(c testcontainers.Container) string {
	// Only the tail is read: in whitelist mode blocked host daemons can grow
	// the guard log very large, and streaming the whole file through a docker
	// exec deadlocks the exec (docker's relay buffers while the exec command
	// is still running).
	_, out := s.exec(c, []string{"sh", "-c", fmt.Sprintf("tail -c 262144 %s", netGuardLogPath)})
	return out
}

func (s *IntegrationSuite) stopNetGuard(c testcontainers.Container) {
	if netGuardPID != 0 {
		s.exec(c, []string{"sh", "-c", fmt.Sprintf("kill %d 2>/dev/null || true", netGuardPID)})
		if !s.awaitGone(c, netGuardPID, 3*time.Second) {
			s.exec(c, []string{"sh", "-c", "pkill -f 'app-[l]istener network-guard' || true"})
			s.awaitGone(c, netGuardPID, 3*time.Second)
		}
		netGuardPID = 0
	} else {
		s.exec(c, []string{"sh", "-c", "pkill -f 'app-[l]istener' || true"})
	}
	// Pooled-container invariant: no LIVE stray app-listener may survive.
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		code, _ := s.exec(c, []string{"sh", "-c", noLiveAppListenerProcs})
		if code == 0 {
			return
		}
		s.exec(c, []string{"sh", "-c", "pkill -f 'app-[l]istener' || true"})
		time.Sleep(200 * time.Millisecond)
	}
	s.Require().Fail("netguard cleanup", "stray app-listener processes survived stopNetGuard")
}

// netGuardBlockedEvent extracts a blocked (Blocked=true) event of the given
// type/comm from the guard log.
func netGuardHasBlockedEvent(logContent, expectedComm, expectedType string) bool {
	for _, line := range strings.Split(logContent, "\n") {
		idx := strings.Index(line, "NETGUARD|")
		if idx < 0 {
			continue
		}
		rest := line[idx+len("NETGUARD|"):]
		parts := strings.SplitN(rest, "|", 10)
		if len(parts) < 10 {
			continue
		}
		if strings.TrimSpace(parts[0]) == expectedType && parts[1] == expectedComm && strings.TrimSpace(parts[9]) == "true" {
			return true
		}
	}
	return false
}

// waitForNetGuardBlockedEvent polls the guard log until the exact blocked
// event (type + comm + Blocked=true) appears.
func (s *IntegrationSuite) waitForNetGuardBlockedEvent(c testcontainers.Container, expectedComm, expectedType string, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if netGuardHasBlockedEvent(s.readNetGuardLog(c), expectedComm, expectedType) {
			return
		}
		time.Sleep(250 * time.Millisecond)
	}
	full := s.readNetGuardLog(c)
	s.T().Fatalf("timed out waiting for blocked event (%s/%s); log tail:\n%s\ncontext files:\n%s",
		expectedType, expectedComm, netGuardTail(full), s.guardContextFiles(c))
}

// netGuardBlockedEventCount returns how many blocked (Blocked=true) events of
// the given type/comm appear in the guard log.
func netGuardBlockedEventCount(logContent, comm, typ string) int {
	count := 0
	for _, line := range strings.Split(logContent, "\n") {
		idx := strings.Index(line, "NETGUARD|")
		if idx < 0 {
			continue
		}
		rest := line[idx+len("NETGUARD|"):]
		parts := strings.SplitN(rest, "|", 10)
		if len(parts) < 10 {
			continue
		}
		if strings.TrimSpace(parts[0]) == typ && parts[1] == comm && strings.TrimSpace(parts[9]) == "true" {
			count++
		}
	}
	return count
}

// waitForNetGuardEventCount polls until at least min blocked events of the
// given type/comm are logged.
func (s *IntegrationSuite) waitForNetGuardEventCount(c testcontainers.Container, comm, typ string, min int, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if netGuardBlockedEventCount(s.readNetGuardLog(c), comm, typ) >= min {
			return
		}
		time.Sleep(250 * time.Millisecond)
	}
	full := s.readNetGuardLog(c)
	found := netGuardBlockedEventCount(full, comm, typ)
	s.T().Fatalf("timed out waiting for >=%d blocked %s events (comm=%s), found=%d, logBytes=%d; "+
		"log head/raw:\n%q\ncontext files:\n%s",
		min, typ, comm, found, len(full), netGuardTail(full), s.guardContextFiles(c))
}

// guardContextFiles dumps the marker-related net_tester logs so a
// timing/collision failure shows which phase the test binary reached.
func (s *IntegrationSuite) guardContextFiles(c testcontainers.Container) string {
	var b strings.Builder
	for _, f := range []string{"/tmp/delayed.log", "/tmp/send.log", "/tmp/recv.log", "/tmp/resolved.log"} {
		_, out := s.exec(c, []string{"sh", "-c", "cat " + f + " 2>/dev/null || true"})
		if strings.TrimSpace(out) != "" {
			fmt.Fprintf(&b, "[%s]\n%s\n", f, out)
		}
	}
	return b.String()
}

func netGuardTail(log string) string {
	lines := strings.Split(log, "\n")
	if len(lines) > 30 {
		lines = lines[len(lines)-30:]
	}
	return strings.Join(lines, "\n")
}

// guardNetTypesForComm returns the distinct (blocked) event types logged for a
// given comm.
func guardNetTypesForComm(logContent, comm string) []string {
	seen := make(map[string]bool)
	var types []string
	for _, line := range strings.Split(logContent, "\n") {
		idx := strings.Index(line, "NETGUARD|")
		if idx < 0 {
			continue
		}
		rest := line[idx+len("NETGUARD|"):]
		parts := strings.SplitN(rest, "|", 10)
		if len(parts) < 10 || parts[1] != comm || strings.TrimSpace(parts[9]) != "true" {
			continue
		}
		if !seen[strings.TrimSpace(parts[0])] {
			seen[strings.TrimSpace(parts[0])] = true
			types = append(types, strings.TrimSpace(parts[0]))
		}
	}
	return types
}

// newNetGuardContainer starts a privileged container with the guard binary
// and copies net_tester (+ optional extra binaries) into it.
func (s *IntegrationSuite) newNetGuardContainer(extraBinaries []string) testcontainers.Container {
	c := s.netguardContainer()
	s.Require().NoError(c.CopyFileToContainer(s.ctx, netTesterAmd64Bin, "/net_tester", 0755))
	for _, b := range extraBinaries {
		s.Require().NoError(c.CopyFileToContainer(s.ctx, netTesterAmd64Bin, b, 0755))
	}
	return c
}

func (s *IntegrationSuite) TestNetworkGuard_Blacklist_Allowed() {
	c := s.newNetGuardContainer([]string{"/other_tester"})
	// pooled: terminated at suite end

	s.startNetworkGuardStd(c, "-b", "/net_tester")

	s.exec(c, []string{"sh", "-c", "/other_tester tcp-server 8080 > /tmp/server.log 2>&1 &"})
	s.awaitNetTesterServer(c, 8080, "tcp")

	code, out := s.exec(c, []string{"sh", "-c", "/other_tester tcp-client 127.0.0.1 8080"})
	s.Require().Equalf(0, code, "non-blacklisted binary should be allowed, but failed: %s", out)

	s.stopNetGuard(c)
}

func (s *IntegrationSuite) TestNetworkGuard_Blacklist_BlockedBind() {
	c := s.newNetGuardContainer(nil)
	// pooled: terminated at suite end

	s.startNetworkGuardStd(c, "-b", "/net_tester")

	code, out := s.exec(c, []string{"sh", "-c", "/net_tester tcp-server 8081"})
	s.Require().NotEqualf(0, code, "blacklisted binary should be blocked from binding: %s", out)
	s.waitForNetGuardBlockedEvent(c, "net_tester", "BIND", 8*time.Second)

	s.stopNetGuard(c)
}

func (s *IntegrationSuite) TestNetworkGuard_Blacklist_BlockedConnect() {
	c := s.newNetGuardContainer(nil)
	// pooled: terminated at suite end

	s.startNetworkGuardStd(c, "-b", "/net_tester")

	code, out := s.exec(c, []string{"sh", "-c", "/net_tester tcp-client 127.0.0.1 9999"})
	s.Require().NotEqualf(0, code, "blacklisted binary should be blocked on connect: %s", out)
	s.waitForNetGuardBlockedEvent(c, "net_tester", "CONNECT", 8*time.Second)

	s.stopNetGuard(c)
}

func (s *IntegrationSuite) TestNetworkGuard_Whitelist_Allowed() {
	c := s.newNetGuardContainer([]string{"/other_tester"})
	// pooled: terminated at suite end

	s.startNetworkGuardStd(c, "-w", guardBinaryFlag("/net_tester"))

	s.exec(c, []string{"sh", "-c", "/net_tester tcp-server 8081 > /tmp/server.log 2>&1 &"})
	s.awaitNetTesterServer(c, 8081, "tcp")

	code, out := s.exec(c, []string{"sh", "-c", "/net_tester tcp-client 127.0.0.1 8081"})
	s.Require().Equalf(0, code, "whitelisted binary should be allowed, but failed: %s", out)

	// A non-whitelisted binary must not be able to talk to the server.
	code, out = s.exec(c, []string{"sh", "-c", "/other_tester tcp-client 127.0.0.1 8081"})
	s.Require().NotEqualf(0, code, "non-whitelisted binary should be blocked, but succeeded: %s", out)

	s.stopNetGuard(c)
}

func (s *IntegrationSuite) TestNetworkGuard_Whitelist_BlockedConnect() {
	c := s.newNetGuardContainer([]string{"/other_tester"})
	// pooled: terminated at suite end

	s.startNetworkGuardStd(c, "-w", guardBinaryFlag("/net_tester"))

	code, out := s.exec(c, []string{"sh", "-c", "/other_tester tcp-client 127.0.0.1 9999"})
	s.Require().NotEqualf(0, code, "non-whitelisted binary should be blocked on connect: %s", out)
	s.waitForNetGuardBlockedEvent(c, "other_tester", "CONNECT", 8*time.Second)

	s.stopNetGuard(c)
}

func (s *IntegrationSuite) TestNetworkGuard_Whitelist_BlockedDns() {
	c := s.newNetGuardContainer([]string{"/other_tester"})
	// pooled: terminated at suite end

	s.startNetworkGuardStd(c, "-w", guardBinaryFlag("/net_tester"))

	code, out := s.exec(c, []string{"sh", "-c", "/other_tester dns example.com"})
	s.Require().NotEqualf(0, code, "non-whitelisted binary should be blocked on DNS: %s", out)
	s.waitForNetGuardBlockedEvent(c, "other_tester", "CONNECT", 8*time.Second)

	s.stopNetGuard(c)
}

func (s *IntegrationSuite) TestNetworkGuard_Whitelist_NoFlags() {
	c := s.newNetGuardContainer(nil)
	// pooled: terminated at suite end

	// No -w/-b: whitelist mode with no allowed binaries -> everything blocked.
	s.startNetworkGuardStd(c)

	code, out := s.exec(c, []string{"sh", "-c", "/net_tester tcp-server 8082"})
	s.Require().NotEqualf(0, code, "with no -w/-b everything should be blocked (bind), got: %s", out)
	s.waitForNetGuardBlockedEvent(c, "net_tester", "BIND", 8*time.Second)

	code, _ = s.exec(c, []string{"sh", "-c", "/net_tester tcp-client 127.0.0.1 8082"})
	s.Require().NotEqualf(0, code, "with no -w/-b everything should be blocked (connect)")
	s.waitForNetGuardBlockedEvent(c, "net_tester", "CONNECT", 8*time.Second)

	s.stopNetGuard(c)
}

func (s *IntegrationSuite) TestNetworkGuard_Whitelist_BlockedListen() {
	c := s.newNetGuardContainer([]string{"/other_tester"})
	// pooled: terminated at suite end

	// Bind a socket before the guard attaches; it sleeps until the marker file
	// appears and then calls listen(), which must be denied.
	s.exec(c, []string{"sh", "-c", "/other_tester tcp-server-delayed 8090 /tmp/go > /tmp/delayed.log 2>&1 &"})
	s.startNetworkGuardStd(c, "-w", guardBinaryFlag("/net_tester"))

	s.exec(c, []string{"touch", "/tmp/go"})
	s.waitForNetGuardBlockedEvent(c, "other_tester", "LISTEN", 8*time.Second)
	s.stopNetGuard(c)
}

func (s *IntegrationSuite) TestNetworkGuard_Whitelist_BlockedSend() {
	c := s.newNetGuardContainer([]string{"/other_tester"})
	// pooled: terminated at suite end

	// Start a sender that dials before the guard attaches and then waits for
	// the marker; its later sendmsg calls must be denied once the guard is up.
	s.exec(c, []string{"sh", "-c", "/other_tester udp-send-loop 127.0.0.1 9091 /tmp/go 300 > /tmp/send.log 2>&1 &"})
	s.startNetworkGuardStd(c, "-w", guardBinaryFlag("/net_tester"))

	s.exec(c, []string{"touch", "/tmp/go"})
	s.waitForNetGuardBlockedEvent(c, "other_tester", "SEND", 8*time.Second)
	s.stopNetGuard(c)
}

func (s *IntegrationSuite) TestNetworkGuard_Whitelist_BlockedRecv() {
	c := s.newNetGuardContainer([]string{"/other_tester"})
	// pooled: terminated at suite end

	// Bind + recv loop starts before the guard attaches and waits for the
	// marker; its later recvmsg calls must be denied.
	s.exec(c, []string{"sh", "-c", "/other_tester udp-recv-loop 9092 /tmp/go > /tmp/recv.log 2>&1 &"})
	s.startNetworkGuardStd(c, "-w", guardBinaryFlag("/net_tester"))

	s.exec(c, []string{"touch", "/tmp/go"})
	s.waitForNetGuardBlockedEvent(c, "other_tester", "RECV", 8*time.Second)
	s.stopNetGuard(c)
}

func (s *IntegrationSuite) TestNetworkGuard_NoThrottle() {
	c := s.newNetGuardContainer([]string{"/other_tester"})
	// pooled: terminated at suite end

	// Only SEND events enter the pipeline: host daemons (e.g. avahi) flood
	// recvmsg, and with --no-throttle that noise would overflow the ring
	// buffer / scroll the tail window before other_tester's events arrive.
	s.exec(c, []string{"sh", "-c", "/other_tester udp-send-loop 127.0.0.1 9095 /tmp/go 50 > /tmp/send.log 2>&1 &"})
	s.startNetworkGuardStd(c, "-w", guardBinaryFlag("/net_tester"), "-e", "SEND", "--no-throttle")

	// The sender dials before the guard attaches, then sends 10 datagrams at
	// 50ms once the marker appears. With --no-throttle every blocked send is
	// logged; the default 250ms per-(type, comm) throttle would log ~2.
	s.exec(c, []string{"touch", "/tmp/go"})
	s.waitForNetGuardEventCount(c, "other_tester", "SEND", 6, 8*time.Second)
	s.stopNetGuard(c)
}

func (s *IntegrationSuite) TestNetworkGuard_EventFilter() {
	c := s.newNetGuardContainer([]string{"/other_tester"})
	// pooled: terminated at suite end

	// Only CONNECT events may be reported/blocked.
	s.startNetworkGuardStd(c, "-w", guardBinaryFlag("/net_tester"), "-e", "CONNECT")

	code, _ := s.exec(c, []string{"sh", "-c", "/other_tester tcp-client 127.0.0.1 9999"})
	s.Require().NotEqualf(0, code, "non-whitelisted binary should be blocked on connect")
	s.waitForNetGuardBlockedEvent(c, "other_tester", "CONNECT", 8*time.Second)

	types := guardNetTypesForComm(s.readNetGuardLog(c), "other_tester")
	s.Require().Equal([]string{"CONNECT"}, types, "event filter should only report CONNECT events")

	s.stopNetGuard(c)
}

func (s *IntegrationSuite) TestNetworkGuard_AutoInfra() {
	c := s.newNetGuardContainer(nil)
	// pooled: terminated at suite end

	// Create fake systemd-resolved paths
	s.exec(c, []string{"mkdir", "-p", "/usr/lib/systemd", "/run/systemd/resolve"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, netTesterAmd64Bin, "/usr/lib/systemd/systemd-resolved", 0755))
	s.exec(c, []string{"touch", "/run/systemd/resolve/io.systemd.Resolve"})

	// Start systemd-resolved so it is detected as running
	s.exec(c, []string{"sh", "-c", "/usr/lib/systemd/systemd-resolved tcp-server 9000 > /tmp/resolved.log 2>&1 &"})
	s.awaitNetTesterServer(c, 9000, "tcp")

	s.startNetworkGuardStd(c, "-w", guardBinaryFlag("/net_tester"), "--auto-infra")

	// systemd-resolved is automatically allowlisted by --auto-infra.
	code, out := s.exec(c, []string{"sh", "-c", "/usr/lib/systemd/systemd-resolved tcp-client 127.0.0.1 9000"})
	s.Require().Equalf(0, code, "systemd-resolved should be automatically allowlisted, but failed: %s", out)

	// A normal unlisted binary should still fail.
	s.Require().NoError(c.CopyFileToContainer(s.ctx, netTesterAmd64Bin, "/other_tester", 0755))
	code, _ = s.exec(c, []string{"sh", "-c", "/other_tester tcp-client 127.0.0.1 9000"})
	s.Require().NotEqualf(0, code, "other binary should be blocked, but succeeded")
	s.waitForNetGuardBlockedEvent(c, "other_tester", "CONNECT", 8*time.Second)

	s.stopNetGuard(c)
}

func (s *IntegrationSuite) TestNetworkGuard_CLI_Errors() {
	c := s.newNetGuardContainer(nil)
	// pooled: terminated at suite end

	cases := []struct {
		name string
		cmd  string
		want string
	}{
		{"mutually-exclusive", "-b /net_tester -w /net_tester", "mutually exclusive"},
		{"auto-infra-blacklist", "-b /net_tester --auto-infra", "auto-infra is only available in whitelist mode"},
		{"unsafe-blacklist", "-b /net_tester --unsafe", "unsafe is only available in whitelist mode"},
		{"bad-event", "-w /net_tester -e BOGUS", "unknown network event type"},
	}
	for _, tc := range cases {
		code, out := s.exec(c, []string{"sh", "-c", "timeout 10 /app-listener network-guard " + tc.cmd + " --headless 2>&1"})
		s.Require().NotEqualf(0, code, "%s: should fail, got: %s", tc.name, out)
		s.Require().Containsf(out, tc.want, "%s: unexpected error message: %s", tc.name, out)
	}
}
