package integrationtests

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	amd64Bin          string
	arm64Bin          string
	netTesterAmd64Bin string
)

// hostInfraBindPaths are host binaries that the LSM network guard blocks in
// whitelist mode because its programs attach globally (they affect host
// processes too, not just the container). Bind-mounting them into the
// container at a private path lets the guard allowlist them by their real
// dev/ino, keeping docker and the host resolver functional while a whitelist
// guard is active.
var hostInfraBindPaths = []string{
	"/usr/bin/dockerd",
	"/usr/bin/containerd",
	"/usr/bin/containerd-shim-runc-v2",
	"/usr/bin/docker-proxy",
	"/usr/bin/docker",
	"/usr/lib/systemd/systemd-resolved",
}

// infraContainerPaths returns the container-side path for each bind-mounted
// host infra binary, used to build the network-guard allowlist/blocklist flag.
func infraContainerPaths() []string {
	out := make([]string, 0, len(hostInfraBindPaths))
	for _, p := range hostInfraBindPaths {
		if _, err := os.Stat(p); err != nil {
			continue
		}
		out = append(out, "/host-infra/"+path.Base(p))
	}
	return out
}

// guardBinaryFlag builds a comma-separated -b/-w value covering the main
// net_tester plus the mounted host infra binaries.
func guardBinaryFlag(mainPath string) string {
	parts := append([]string{mainPath}, infraContainerPaths()...)
	return strings.Join(parts, ",")
}

func TestMain(m *testing.M) {
	if err := os.MkdirAll("../build/test", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "mkdir build/test: %v\n", err)
		os.Exit(1)
	}

	amd64Bin = absPath("../build/test/app-listener-amd64")
	netTesterAmd64Bin = absPath("../build/test/net_tester-amd64")

	// The daemon mode links fscrypt, which needs cgo (mlock), so a pure-Go
	// CGO_ENABLED=0 build no longer compiles. The binary must stay static
	// anyway: it runs inside ubuntu:latest containers whose glibc (2.39)
	// is older than the host's, so dynamic linking would be rejected at
	// load time.
	cmd := exec.Command("go", "build", "-tags", "ci",
		"-ldflags", "-linkmode external -extldflags -static",
		"-o", amd64Bin, "..")
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "CGO_ENABLED=1")
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "build amd64 binary: %v\n", err)
		os.Exit(1)
	}

	cmdNet := exec.Command("go", "build", "-tags", "ci", "-o", netTesterAmd64Bin, "./net_tester")
	cmdNet.Stderr = os.Stderr
	cmdNet.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64")
	if err := cmdNet.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "build net_tester amd64 binary: %v\n", err)
		os.Exit(1)
	}

	arm64Bin = absPath("../build/test/app-listener-arm64")
	cmd2 := exec.Command("go", "build", "-tags", "ci", "-o", arm64Bin, "..")
	cmd2.Env = append(os.Environ(), "GOOS=linux", "GOARCH=arm64", "CGO_ENABLED=0")
	cmd2.Stderr = os.Stderr
	if err := cmd2.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "arm64 cross-build: %v (will skip arm64 tests)\n", err)
		arm64Bin = ""
	}

	os.Exit(m.Run())
}

type IntegrationSuite struct {
	suite.Suite
	ctx context.Context
}

func TestIntegrationSuite(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("Docker not available, skipping integration tests")
	}
	suite.Run(t, new(IntegrationSuite))
}

func (s *IntegrationSuite) SetupSuite() {
	s.ctx = context.Background()
}

func (s *IntegrationSuite) startContainer(
	image string,
	platform string,
	privileged bool,
	binaryPath string,
) testcontainers.Container {
	req := testcontainers.ContainerRequest{
		Image:         image,
		ImagePlatform: platform,
		Cmd:           []string{"sleep", "3600"},
		WaitingFor:    wait.ForLog("").WithPollInterval(100 * time.Millisecond).WithStartupTimeout(60 * time.Second),
	}

	if binaryPath != "" {
		req.Files = []testcontainers.ContainerFile{
			{
				HostFilePath:      binaryPath,
				ContainerFilePath: "/app-listener",
				FileMode:          0755,
			},
		}
	}

	if privileged {
		req.HostConfigModifier = func(hc *container.HostConfig) {
			hc.Privileged = true
			hc.CapAdd = []string{"BPF", "SYS_ADMIN"}
		}
		req.Mounts = testcontainers.ContainerMounts{
			{
				Source: testcontainers.GenericBindMountSource{HostPath: "/sys/kernel/btf"},
				Target: "/sys/kernel/btf",
			},
		}
		for _, hp := range hostInfraBindPaths {
			if _, err := os.Stat(hp); err != nil {
				continue
			}
			req.Mounts = append(req.Mounts, testcontainers.ContainerMount{
				Source: testcontainers.GenericBindMountSource{HostPath: hp},
				Target: testcontainers.ContainerMountTarget("/host-infra/" + path.Base(hp)),
			})
		}
	}

	c, err := testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	s.Require().NoError(err, "start container %s/%s", image, platform)
	return c
}

func (s *IntegrationSuite) exec(c testcontainers.Container, cmd []string) (int, string) {
	exitCode, reader, err := c.Exec(s.ctx, cmd)
	s.Require().NoError(err, "exec: %v", cmd)

	b, err := io.ReadAll(reader)
	s.Require().NoError(err, "read exec output: %v", cmd)

	return exitCode, strings.TrimSpace(string(b))
}

func (s *IntegrationSuite) logs(c testcontainers.Container) string {
	r, err := c.Logs(s.ctx)
	s.Require().NoError(err)
	defer r.Close()

	b, err := io.ReadAll(r)
	s.Require().NoError(err)
	return string(b)
}

func absPath(p string) string {
	abs, err := filepath.Abs(p)
	if err != nil {
		panic(err)
	}
	return abs
}

// readMonitorLog returns the raw content of the monitor log file.
func (s *IntegrationSuite) readMonitorLog(c testcontainers.Container) string {
	_, out := s.exec(c, []string{"sh", "-c", "cat /tmp/monitor.log 2>/dev/null || true"})
	return out
}

// newEventTypes extracts unique event type names (OPEN, READ, …) from EVENT|
// lines that appear in newLog but were not yet present in oldLog.
func newEventTypes(oldLog, newLog string) []string {
	if newLog == "" {
		return nil
	}

	// areLinesDelta is true when newLog and oldLog share a common prefix.
	areLinesDelta := oldLog != "" && len(newLog) > len(oldLog)/2

	var tail string
	if areLinesDelta {
		oldLines := strings.Split(oldLog, "\n")
		newLines := strings.Split(newLog, "\n")
		if len(newLines) > len(oldLines) {
			tail = strings.Join(newLines[len(oldLines):], "\n")
		}
	} else {
		tail = newLog
	}

	if tail == "" {
		return nil
	}

	seen := make(map[string]bool)
	var result []string
	for _, line := range strings.Split(tail, "\n") {
		if idx := strings.Index(line, "EVENT|"); idx >= 0 {
			rest := line[idx+6:]
			parts := strings.SplitN(rest, "|", 2)
			if len(parts) >= 1 && parts[0] != "" && !seen[parts[0]] {
				seen[parts[0]] = true
				result = append(result, parts[0])
			}
		}
	}
	return result
}

// requireNewEventType asserts that the event type appears in the log delta.
func (s *IntegrationSuite) requireNewEventType(oldLog, newLog, expected string) {
	events := newEventTypes(oldLog, newLog)
	s.Require().True(slices.Contains(events, expected),
		"expected EVENT|%s in new log lines, got %v", expected, events)
}

func (s *IntegrationSuite) requireNoNewEventType(oldLog, newLog, unexpected string) {
	events := newEventTypes(oldLog, newLog)
	s.Require().False(slices.Contains(events, unexpected),
		"unexpected EVENT|%s found in new log lines (got %v)", unexpected, events)
}

// waitForEventType polls the monitor log until expected event type appears or timeout.
// Returns the updated log content. Fails the test if timeout reached.
func (s *IntegrationSuite) waitForEventType(c testcontainers.Container, oldLog string, expected string, timeout time.Duration) string {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		newLog := s.readMonitorLog(c)
		events := newEventTypes(oldLog, newLog)
		for _, evt := range events {
			if evt == expected {
				return newLog
			}
		}
		time.Sleep(200 * time.Millisecond)
	}
	newLog := s.readMonitorLog(c)
	s.Require().Failf("timeout waiting for event",
		"expected EVENT|%s within %v, got events: %v", expected, timeout, newEventTypes(oldLog, newLog))
	return newLog
}
