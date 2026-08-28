package integrationtests

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

// ---------------------------------------------------------------
// Container pool: one privileged container per test FILE, reused by every
// test in that file. Creating and destroying a container costs ~13-15s
// (docker create + binary copy + 10s stop grace + ryuk teardown), which
// dominated the suite runtime. Tests within a file are sequential and use
// per-test state isolation, so a shared container is safe:
//
//   - readiness is polled (log markers + PID liveness), never slept;
//   - guard/monitor processes are stopped PID-scoped between tests;
//   - before each hand-off the pool verifies the container is running,
//     reachable and has NO stray app-listener processes, recreating it
//     otherwise (the running test fails as before — no logic retries).
// ---------------------------------------------------------------

type pooledContainer struct {
	container testcontainers.Container
	uses      int
}

// maxPoolUses caps how many tests share one container: a long-lived
// container accumulates overlayfs whiteout/dentry-cache state from repeated
// `rm -rf`+recreate cycles on the same work paths, which eventually degrades
// the guard's kernel-side path reconstruction (events still enforced, but
// the reported path may drop components). Rotating the container keeps every
// test in a fresh overlay while still saving ~90% of the churn.
const maxPoolUses = 12

var pools map[string]*pooledContainer

// acquirePool returns the shared container for the given domain, starting or
// recreating it when missing, unhealthy, or past its usage budget.
func (s *IntegrationSuite) acquirePool(domain string) testcontainers.Container {
	if pools == nil {
		pools = make(map[string]*pooledContainer)
	}
	if p, ok := pools[domain]; ok && p != nil && p.uses < maxPoolUses && s.poolHealthy(p.container) {
		p.uses++
		return p.container
	}
	if p, ok := pools[domain]; ok && p != nil {
		_ = p.container.Terminate(s.ctx)
	}
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	pools[domain] = &pooledContainer{container: c, uses: 1}
	return c
}

// noLiveAppListenerProcs is the zombie-tolerant stray-process probe: killed
// children of the container init can linger as <defunct> entries which
// pgrep/kill -0 still see, but a zombie holds no resources and is not a
// stray. Exits 0 when no LIVE (non-zombie) app-listener process exists,
// 1 when one does.
const noLiveAppListenerProcs = `ps -eo stat,comm | awk 'BEGIN{clean=1} $1 !~ /^Z/ && $2 ~ /app-listener/ {clean=0} END{exit clean ? 0 : 1}'`

// poolHealthy reports whether the pooled container is running, reachable and
// free of stray live app-listener processes from a previous test.
func (s *IntegrationSuite) poolHealthy(c testcontainers.Container) bool {
	if c == nil {
		return false
	}
	// exec ping: must not use s.exec (Require would fail the running test).
	if code, _, err := c.Exec(s.ctx, []string{"sh", "-c", "echo ok"}); err != nil || code != 0 {
		return false
	}
	code, _, err := c.Exec(s.ctx, []string{"sh", "-c", noLiveAppListenerProcs})
	if err != nil || code != 0 {
		return false
	}
	return true
}

// teardownPools terminates every pooled container at suite end.
func (s *IntegrationSuite) teardownPools() {
	for name, p := range pools {
		if p != nil && p.container != nil {
			_ = p.container.Terminate(s.ctx)
		}
		delete(pools, name)
	}
}

// Domain accessors. Each resets the standard per-domain state so every test
// starts from a clean slate inside the shared container.

// guardContainer is the pool for guard_test.go: resets the standard guard
// work directories. Guard logs are self-resetting (startGuardStd truncates
// them via `>` on every start). A leftover mount on /watch (a failed
// tmpfs-variant test) is lazily unmounted first — files inside a submount
// report mount-relative paths in guard events.
func (s *IntegrationSuite) guardContainer() testcontainers.Container {
	c := s.acquirePool("guard")
	s.exec(c, []string{"sh", "-c", "mountpoint -q /watch && umount /watch 2>/dev/null; rm -rf /watch /protected /exploits /ctl; mkdir -p /watch"})
	return c
}

// monitorContainer is the pool for monitor_test.go: resets /watch; the
// monitor log is truncated by startMonitorStd on every start.
func (s *IntegrationSuite) monitorContainer() testcontainers.Container {
	c := s.acquirePool("monitor")
	s.exec(c, []string{"sh", "-c", "rm -rf /watch; mkdir -p /watch"})
	return c
}

// netguardContainer is the pool for networkguard_test.go, which already
// shares one container and one delta-parsed log today. The reset clears the
// inter-test coordination markers (/tmp/go gates the delayed net_tester
// servers: a stale marker would let a pre-attach server listen before the
// guard attaches, defeating the pre-attach assertions) and the per-test
// server logs.
func (s *IntegrationSuite) netguardContainer() testcontainers.Container {
	c := s.acquirePool("netguard")
	s.exec(c, []string{"sh", "-c", "rm -f /tmp/go /tmp/guard.log /tmp/server.log /tmp/delayed.log /tmp/send.log /tmp/recv.log /tmp/resolved.log"})
	return c
}

// netmonContainer is the pool for networkmonitor_test.go.
func (s *IntegrationSuite) netmonContainer() testcontainers.Container {
	c := s.acquirePool("netmon")
	s.exec(c, []string{"sh", "-c", "rm -f /tmp/monitor.log /tmp/server.log"})
	return c
}

// ---------------------------------------------------------------
// PID-scoped process lifecycle
// ---------------------------------------------------------------

// capturePID runs cmd (which must background a process itself) and returns
// the PID of the backgrounded process via `echo $!`.
//
// Docker exec output is a raw multiplexed stream whose 8-byte frame headers
// can leak into the captured text (the whole suite reads it with
// garbage-tolerant Contains for this reason), so the PID is extracted from
// the trailing digit run instead of parsing strictly.
func (s *IntegrationSuite) capturePID(c testcontainers.Container, cmd string) int {
	code, out := s.exec(c, []string{"sh", "-c", cmd + "\necho $!"})
	s.Require().Equalf(0, code, "starting background process: %s", out)
	m := trailingDigitsRe.FindStringSubmatch(strings.TrimSpace(out))
	if m == nil {
		s.Require().Failf("could not capture background PID", "cmd=%q output=%q", cmd, out)
	}
	pid := 0
	fmt.Sscanf(m[1], "%d", &pid)
	if pid <= 0 {
		s.Require().Failf("could not capture background PID", "cmd=%q output=%q", cmd, out)
	}
	return pid
}

// trailingDigitsRe matches the digit run at the end of a captured exec
// output (frame-header tolerant).
var trailingDigitsRe = regexp.MustCompile(`([0-9]+)\s*$`)

// awaitGone polls until the given PID no longer exists (bounded); returns
// false when it is still alive after the deadline.
func (s *IntegrationSuite) awaitGone(c testcontainers.Container, pid int, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		_, out := s.exec(c, []string{"sh", "-c", fmt.Sprintf("kill -0 %d 2>/dev/null && echo alive || echo gone", pid)})
		if strings.Contains(out, "gone") {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}

// waitPortListening polls /proc/net/tcp{,6} for a LISTEN socket on port
// (host byte order), a deterministic readiness condition for net_tester
// servers that replaces fixed sleeps.
func (s *IntegrationSuite) waitPortListening(c testcontainers.Container, port int, timeout time.Duration) bool {
	hex := fmt.Sprintf(":%04X", port)
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		code, _ := s.exec(c, []string{"sh", "-c", fmt.Sprintf(
			"awk 'NR>1 && $4==\"0A\" && $2 ~ /%s$/ {found=1} END{exit !found}' /proc/net/tcp /proc/net/tcp6 2>/dev/null", hex)})
		if code == 0 {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}

// awaitNetTesterServer waits until a net_tester server has actually bound
// its port — a LISTEN socket for TCP, a bound local address for UDP —
// replacing the fixed 1s startup sleeps with a real condition.
func (s *IntegrationSuite) awaitNetTesterServer(c testcontainers.Container, port int, proto string) {
	if proto == "udp" {
		hex := fmt.Sprintf(":%04X", port)
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			code, _ := s.exec(c, []string{"sh", "-c", fmt.Sprintf(
				"awk 'NR>1 && $2 ~ /%s$/ {found=1} END{exit !found}' /proc/net/udp /proc/net/udp6 2>/dev/null", hex)})
			if code == 0 {
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
		s.Require().Failf("udp server did not bind", "port %d was never bound", port)
		return
	}
	if !s.waitPortListening(c, port, 5*time.Second) {
		s.Require().Failf("tcp server did not listen", "port %d never entered LISTEN", port)
	}
}
