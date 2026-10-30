package integrationtests

import (
	"bytes"
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
	tcexec "github.com/testcontainers/testcontainers-go/exec"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	amd64Bin             string
	netTesterAmd64Bin    string
	guardTestAmd64Bin    string
	fscryptTestAmd64Bin  string
	tuiTestAmd64Bin      string
	netguardTestAmd64Bin string
	recorderTestAmd64Bin string
	secretProbeAmd64Bin  string
)

// hostInfraBindPaths are host binaries the LSM network guard blocks in whitelist mode because its
// programs attach globally (host processes too). Bind-mounting them into the container at a private
// path lets the guard allowlist them by real dev/ino, keeping docker and the host resolver working.
var hostInfraBindPaths = []string{
	"/usr/bin/dockerd",
	"/usr/bin/containerd",
	"/usr/bin/containerd-shim-runc-v2",
	"/usr/bin/docker-proxy",
	"/usr/bin/docker",
	"/usr/lib/systemd/systemd-resolved",
}

// infraContainerPaths returns the container-side path for each bind-mounted
// host infra binary, used to build the network guard allowlist/blocklist flag.
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
	guardTestAmd64Bin = absPath("../build/test/guard-amd64")
	fscryptTestAmd64Bin = absPath("../build/test/fscrypt-amd64")
	tuiTestAmd64Bin = absPath("../build/test/tui-amd64")
	netguardTestAmd64Bin = absPath("../build/test/networkguard-amd64")
	recorderTestAmd64Bin = absPath("../build/test/recorder-amd64")
	secretProbeAmd64Bin = absPath("../build/test/secret_probe-amd64")

	// Daemon mode links fscrypt, which needs cgo (mlock): CGO_ENABLED=0 no longer compiles. The
	// binary must still be static: it runs in ubuntu:latest containers whose glibc (2.39) is older
	// than the host's. osusergo+netgo avoid glibc NSS; output is buffered so benign static-link
	// warnings show only on failure.
	cmd := exec.Command("go", "build", "-tags", "ci,osusergo,netgo",
		"-ldflags", "-linkmode external -extldflags -static",
		"-o", amd64Bin, "..")
	var buildErr bytes.Buffer
	cmd.Stderr = &buildErr
	cmd.Env = append(os.Environ(), "CGO_ENABLED=1")
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "build amd64 binary: %v\n%s\n", err, buildErr.String())
		os.Exit(1)
	}
	fmt.Fprintln(os.Stderr, "built amd64 binary")

	cmdNet := exec.Command("go", "build", "-tags", "ci", "-o", netTesterAmd64Bin, "./net_tester")
	cmdNet.Stderr = os.Stderr
	cmdNet.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64")
	if err := cmdNet.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "build net_tester amd64 binary: %v\n", err)
		os.Exit(1)
	}

	cmdSecret := exec.Command("go", "build", "-o", secretProbeAmd64Bin, "./secret_probe")
	cmdSecret.Stderr = os.Stderr
	cmdSecret.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64")
	if err := cmdSecret.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "build secret_probe amd64 binary: %v\n", err)
		os.Exit(1)
	}

	// The guard BPF-gated tests (pending-binaries resolution, inode map
	// population) need root + a live BPF LSM, which only the privileged
	// test containers can provide. Compile the guard test binary as a
	// static executable (the guard package needs no cgo) and run it inside
	// one of those containers.
	cmdGuardTest := exec.Command("go", "test", "-tags", "ci", "-c",
		"-o", guardTestAmd64Bin, "../internal/guard/")
	cmdGuardTest.Stderr = os.Stderr
	cmdGuardTest.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64")
	if err := cmdGuardTest.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "build guard test binary: %v\n", err)
		os.Exit(1)
	}

	// The edit-protected editor's root-only checks (internal/tui/harness_test.go), built the same
	// way: the package needs no cgo.
	cmdTuiTest := exec.Command("go", "test", "-tags", "ci", "-c", "-o", tuiTestAmd64Bin, "../internal/tui/")
	cmdTuiTest.Stderr = os.Stderr
	cmdTuiTest.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64")
	if err := cmdTuiTest.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "build tui test binary: %v\n", err)
		os.Exit(1)
	}

	// The network guard's root-only checks (internal/networkguard/harness_test.go): no cgo either.
	cmdNetguardTest := exec.Command("go", "test", "-tags", "ci", "-c", "-o", netguardTestAmd64Bin,
		"../internal/networkguard/")
	cmdNetguardTest.Stderr = os.Stderr
	cmdNetguardTest.Env = append(os.Environ(), "CGO_ENABLED=0", "GOOS=linux", "GOARCH=amd64")
	if err := cmdNetguardTest.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "build network guard test binary: %v\n", err)
		os.Exit(1)
	}

	// The fscrypt TOCTOU tests need root plus real kernel fscrypt (directory) or the file-vault
	// AEAD (regular file): the production code, exercised directly like guardTestAmd64Bin. It links
	// google/fscrypt (cgo, mlock; see CLAUDE.md Build), so it's built like the main amd64 binary:
	// CGO_ENABLED=1, static for the container's older glibc.
	cmdFscryptTest := exec.Command("go", "test", "-tags", "ci", "-c",
		"-ldflags", "-linkmode external -extldflags -static",
		"-o", fscryptTestAmd64Bin, "../internal/fscrypt/")
	var fscryptBuildErr bytes.Buffer
	cmdFscryptTest.Stderr = &fscryptBuildErr
	cmdFscryptTest.Env = append(os.Environ(), "CGO_ENABLED=1", "GOOS=linux", "GOARCH=amd64")
	if err := cmdFscryptTest.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "build fscrypt test binary: %v\n%s\n", err, fscryptBuildErr.String())
		os.Exit(1)
	}

	// The record-parity harness (internal/recorder/harness_test.go): the recorder links fscrypt
	// (apply.go), so cgo + static like the fscrypt harness.
	cmdRecorderTest := exec.Command("go", "test", "-tags", "ci", "-c",
		"-ldflags", "-linkmode external -extldflags -static",
		"-o", recorderTestAmd64Bin, "../internal/recorder/")
	var recorderBuildErr bytes.Buffer
	cmdRecorderTest.Stderr = &recorderBuildErr
	cmdRecorderTest.Env = append(os.Environ(), "CGO_ENABLED=1", "GOOS=linux", "GOARCH=amd64")
	if err := cmdRecorderTest.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "build recorder test binary: %v\n%s\n", err, recorderBuildErr.String())
		os.Exit(1)
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

// TearDownSuite terminates the per-file pooled containers (see pool_test.go).
func (s *IntegrationSuite) TearDownSuite() {
	s.teardownPools()
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

	// Containers have no securityfs mount, so /sys/kernel/security/lsm is
	// unreadable inside them. The binary's BPF LSM preflight would therefore
	// (correctly) hard-fail; the host kernel has already been verified by the
	// suite's own environment before containers start.
	req.Env = map[string]string{"APPLISTENER_ASSUME_BPF_LSM": "1"}

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
			// docker-init (tini) reaps adopted orphans: backgrounded
			// guards/monitors killed between pooled tests would otherwise
			// stay as <defunct> children of the container init (a plain
			// `sleep` never wait()s), breaking liveness checks.
			init := true
			hc.Init = &init
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
	// Multiplexed() demuxes docker's raw stdcopy stream into a memory
	// buffer before Exec returns: without it the reader is the raw
	// hijacked stream, whose frame headers leak into the captured text and
	// which can deliver empty reads under load (the transport races the
	// exec-inspect polling loop). The buffer makes output deterministic.
	exitCode, reader, err := c.Exec(s.ctx, cmd, tcexec.Multiplexed())
	s.Require().NoError(err, "exec: %v", cmd)

	b, err := io.ReadAll(reader)
	s.Require().NoError(err, "read exec output: %v", cmd)

	return exitCode, strings.TrimSpace(string(b))
}

// requireDistinctExes fails unless every path is its own dev:ino. Policy keys on the exe inode, so
// multicall coreutils (uutils, busybox: cat, head, mkdir... are one inode) make a "whitelisted"
// and a "not whitelisted" fixture the same key, and the test passes or fails for the wrong reason.
func (s *IntegrationSuite) requireDistinctExes(c testcontainers.Container, paths ...string) {
	code, out := s.exec(c, append([]string{"stat", "-Lc", "%d:%i %n"}, paths...))
	s.Require().Equalf(0, code, "stat %v: %s", paths, out)
	owner := make(map[string]string, len(paths))
	for line := range strings.SplitSeq(out, "\n") {
		id, name, _ := strings.Cut(strings.TrimSpace(line), " ")
		if prev, dup := owner[id]; dup {
			s.Require().Failf("fixture binaries share an inode",
				"%s and %s are both %s (a multicall binary?): pick binaries with their own inode", prev, name, id)
		}
		owner[id] = name
	}
	s.Require().Lenf(owner, len(paths), "stat output: %s", out)
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
	_, out := s.exec(c, []string{"sh", "-c", fmt.Sprintf("cat %s 2>/dev/null || true", monitorLogPath)})
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
	for line := range strings.SplitSeq(tail, "\n") {
		if _, rest, ok := strings.Cut(line, "EVENT|"); ok {
			parts := strings.SplitN(rest, "|", 2)
			if len(parts) >= 1 && parts[0] != "" && !seen[parts[0]] {
				seen[parts[0]] = true
				result = append(result, parts[0])
			}
		}
	}
	return result
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
		if slices.Contains(events, expected) {
			return newLog
		}
		time.Sleep(200 * time.Millisecond)
	}
	newLog := s.readMonitorLog(c)
	s.Require().Failf("timeout waiting for event",
		"expected EVENT|%s within %v, got events: %v\nraw log tail:\n%s", expected, timeout, newEventTypes(oldLog, newLog), tailLast(newLog, 20))
	return newLog
}

// tailLast returns the last n lines of s for failure diagnostics.
func tailLast(s string, n int) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, "\n")
}
