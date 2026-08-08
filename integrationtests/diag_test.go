package integrationtests

import (
	"time"
)

// TestDiag_GuardLSM reproduces the guard-mode drop-out observed on
// GitHub-hosted runners inside a single container and dumps the guard log
// and kernel state into the test output, so the root cause is visible
// without running the whole suite. Temporary diagnostic helper: remove it
// once the integration CI issue is understood and fixed.
func (s *IntegrationSuite) TestDiag_GuardLSM() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	// Kernel / LSM state as seen from inside the privileged container.
	for _, probe := range []string{
		"uname -r",
		"cat /sys/kernel/security/lsm 2>&1 || echo '<no lsm file>'",
		"cat /proc/sys/kernel/unprivileged_bpf_disabled 2>&1",
		"ls -l /sys/kernel/btf/vmlinux 2>&1",
		"grep CapEff /proc/self/status",
	} {
		_, out := s.exec(c, []string{"sh", "-c", probe})
		s.T().Logf("PROBE  %-40s => %s", probe, out)
	}

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	s.exec(c, []string{"sh", "-c", "echo 'guarded content' > /watch/target.txt"})
	_, out := s.exec(c, []string{"sh", "-c",
		"nohup /app-listener guard /watch --headless -l debug > /tmp/guard.log 2>&1 & echo started=$?"})
	s.T().Logf("guard launch: %s", out)

	time.Sleep(8 * time.Second)

	code, out := s.exec(c, []string{"sh", "-c",
		"cat /watch/target.txt > /dev/null 2>&1; echo rc=$?; pgrep -f 'app-listener guard' || echo '<guard process gone>'"})
	s.T().Logf("BLOCKED cat => %s (exec code %d)", out, code)

	_, logTail := s.exec(c, []string{"sh", "-c", "tail -c 65536 /tmp/guard.log 2>&1 || true"})
	s.T().Logf("--- /tmp/guard.log tail ---\n%s\n-----------------------------", logTail)

	_, dmesgOut := s.exec(c, []string{"sh", "-c",
		"if [ \"$(cat /proc/sys/kernel/dmesg_restrict 2>/dev/null)\" = 1 ] && ! dmesg >/dev/null 2>&1; then echo '<dmesg restricted>'; else dmesg | tail -60; fi"})
	s.T().Logf("----- dmesg tail -----\n%s\n-----------------------", dmesgOut)

	s.T().Logf("DIAG DONE")
}
