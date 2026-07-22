package integrationtests

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// ---------------------------------------------------------------
// eBPF tests (privileged container)
//
// The monitor requires a TTY for Bubble Tea, so it exits with
// code 1 in non-TTY environments. We verify eBPF success by
// checking the log output before the TUI error.
// ---------------------------------------------------------------

func verifyEBPF(s *IntegrationSuite, code int, out string) {
	s.Require().True(strings.Contains(out, "eBPF available"),
		"eBPF check should pass:\n%s", out)
	s.Require().True(strings.Contains(out, "monitor started"),
		"monitor should start:\n%s", out)
	s.Require().True(strings.Contains(out, "monitor created"),
		"monitor should create probes:\n%s", out)
	s.Require().True(strings.Contains(out, "could not open a new TTY"),
		"should fail only due to missing TTY:\n%s", out)
}

func (s *IntegrationSuite) TestEBPF_Check() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	code, out := s.exec(c, []string{"/app-listener", "monitor", "-w", "/watch", "--recursive"})
	verifyEBPF(s, code, out)
}

func (s *IntegrationSuite) TestEBPF_CheckWithDepth() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	code, out := s.exec(c, []string{"/app-listener", "monitor", "-w", "/watch", "--recursive", "--depth", "2"})
	verifyEBPF(s, code, out)
}

func (s *IntegrationSuite) TestEBPF_FullStack() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})

	monitorCmd := fmt.Sprintf("nohup /app-listener monitor -w /watch --recursive --depth 3 > /tmp/monitor.log 2>&1 &")
	_, _ = s.exec(c, []string{"sh", "-c", monitorCmd})
	time.Sleep(3 * time.Second)

	codeCheck, outCheck := s.exec(c, []string{"pgrep", "-f", "app-listener"})
	s.Require().Equalf(0, codeCheck, "monitor process not running after start:\n%s", outCheck)

	s.exec(c, []string{"touch", "/watch/test.txt"})
	s.exec(c, []string{"sh", "-c", "echo hello > /watch/test.txt"})
	s.exec(c, []string{"rm", "/watch/test.txt"})
	s.exec(c, []string{"mkdir", "/watch/subdir"})
	time.Sleep(2 * time.Second)

	codeAfter, outAfter := s.exec(c, []string{"pgrep", "-f", "app-listener"})
	s.Require().Equalf(0, codeAfter, "monitor crashed after file events:\n%s", outAfter)

	codeLog, outLog := s.exec(c, []string{"cat", "/tmp/monitor.log"})
	if codeLog == 0 {
		s.Require().True(strings.Contains(outLog, "eBPF available"),
			"monitor log missing eBPF check:\n%s", outLog)
		s.Require().True(strings.Contains(outLog, "monitor created"),
			"monitor log missing probe attachment:\n%s", outLog)
	}
}

// ---------------------------------------------------------------
// Cross-architecture tests (QEMU)
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestCrossArch_EBPF_ARM64() {
	if arm64Bin == "" {
		s.T().Skip("arm64 binary not built")
	}

	ok := s.checkQemu("linux/arm64")
	if !ok {
		s.T().Skip("QEMU binfmt not available for arm64")
	}

	c := s.startContainer("ubuntu:latest", "linux/arm64", true, arm64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	code, out := s.exec(c, []string{"/app-listener", "monitor", "/watch", "--recursive"})
	verifyEBPF(s, code, out)
}

func (s *IntegrationSuite) checkQemu(platform string) bool {
	cmd := exec.Command("docker", "run", "--rm", "--platform", platform,
		"ubuntu:latest", "sh", "-c", "echo qemu-ok")
	return cmd.Run() == nil
}

// ---------------------------------------------------------------
// Multi-distro tests
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestMultiDistro_Alpine_EBPF() {
	c := s.startContainer("alpine:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	code, out := s.exec(c, []string{"/app-listener", "monitor", "/watch", "--recursive"})
	verifyEBPF(s, code, out)
}

func (s *IntegrationSuite) TestMultiDistro_Fedora_EBPF() {
	c := s.startContainer("fedora:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	code, out := s.exec(c, []string{"/app-listener", "monitor", "/watch", "--recursive"})
	verifyEBPF(s, code, out)
}
