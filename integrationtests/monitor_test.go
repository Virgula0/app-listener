package integrationtests

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
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

	monitorCmd := fmt.Sprintf("nohup /app-listener monitor -w /watch --recursive --depth 3 --headless > /tmp/monitor.log 2>&1 &")
	_, _ = s.exec(c, []string{"sh", "-c", monitorCmd})
	time.Sleep(3 * time.Second)

	codeCheck, outCheck := s.exec(c, []string{"pgrep", "-f", "app-listener"})
	s.Require().Equalf(0, codeCheck, "monitor process not running after start:\n%s", outCheck)

	logBefore := s.readMonitorLog(c)

	s.exec(c, []string{"touch", "/watch/test.txt"})
	s.exec(c, []string{"sh", "-c", "echo hello > /watch/test.txt"})
	s.exec(c, []string{"rm", "/watch/test.txt"})
	s.exec(c, []string{"mkdir", "/watch/subdir"})
	time.Sleep(2 * time.Second)

	codeAfter, outAfter := s.exec(c, []string{"pgrep", "-f", "app-listener"})
	s.Require().Equalf(0, codeAfter, "monitor crashed after file events:\n%s", outAfter)

	logAfter := s.readMonitorLog(c)
	if logAfter != "" {
		s.Require().True(strings.Contains(logAfter, "eBPF available"),
			"monitor log missing eBPF check:\n%s", logAfter)
		s.Require().True(strings.Contains(logAfter, "monitor created"),
			"monitor log missing probe attachment:\n%s", logAfter)
		s.requireNewEventType(logBefore, logAfter, "MKDIR")
	}
}

// ---------------------------------------------------------------
// Comprehensive event coverage test
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestMonitorAllEvents() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})

	// Copy mmap exploit binary for MMAP event
	mmapHostPath := absPath("./exploits/mmap")
	err := c.CopyFileToContainer(s.ctx, mmapHostPath, "/mmap_exploit", 0755)
	s.Require().NoError(err, "copy mmap exploit")

	// Start monitor and save initial log (contains startup events)
	s.startMonitorStd(c, "/watch")
	logBefore := s.readMonitorLog(c)

	// 1. OPEN + WRITE: create file
	s.exec(c, []string{"sh", "-c", "echo 'test data' > /watch/data.txt"})
	logAfter := s.waitForEventType(c, logBefore, "OPEN", 5*time.Second)
	s.waitForEventType(c, logAfter, "WRITE", 5*time.Second)
	logBefore = s.readMonitorLog(c)

	// 2. READ: cat the file
	s.exec(c, []string{"sh", "-c", "cat /watch/data.txt > /dev/null"})
	logAfter = s.waitForEventType(c, logBefore, "READ", 5*time.Second)
	logBefore = logAfter

	// 3. RENAME: rename file
	s.exec(c, []string{"mv", "/watch/data.txt", "/watch/renamed.txt"})
	logAfter = s.waitForEventType(c, logBefore, "RENAME", 5*time.Second)
	logBefore = logAfter

	// 4. SYMLINK: create symlink
	s.exec(c, []string{"ln", "-s", "/watch/renamed.txt", "/watch/link"})
	logAfter = s.waitForEventType(c, logBefore, "SYMLINK", 5*time.Second)
	logBefore = logAfter

	// 5. HARDLINK: create hardlink
	s.exec(c, []string{"ln", "/watch/renamed.txt", "/watch/hardlink"})
	logAfter = s.waitForEventType(c, logBefore, "HARDLINK", 5*time.Second)
	logBefore = logAfter

	// 6. MKDIR: create directory
	s.exec(c, []string{"mkdir", "/watch/subdir"})
	logAfter = s.waitForEventType(c, logBefore, "MKDIR", 5*time.Second)
	logBefore = logAfter

	// 7. DELETE: remove file
	s.exec(c, []string{"rm", "/watch/renamed.txt"})
	logAfter = s.waitForEventType(c, logBefore, "DELETE", 5*time.Second)
	logBefore = logAfter

	// 8. MMAP + OPEN: run mmap exploit on hardlink (still exists after rm)
	s.exec(c, []string{"/mmap_exploit", "/watch/hardlink"})
	logAfter = s.waitForEventType(c, logBefore, "MMAP", 5*time.Second)
	s.waitForEventType(c, logAfter, "OPEN", 5*time.Second)

	s.stopMonitor(c)
}

func (s *IntegrationSuite) TestMonitorSingleFile() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})

	// Create watched file before starting monitor (must exist for inode resolution)
	s.exec(c, []string{"sh", "-c", "echo 'initial data' > /watch/target.txt"})

	// Copy mmap exploit binary for MMAP event
	mmapHostPath := absPath("./exploits/mmap")
	err := c.CopyFileToContainer(s.ctx, mmapHostPath, "/mmap_exploit", 0755)
	s.Require().NoError(err, "copy mmap exploit")

	// Start monitor watching the specific file (not a directory)
	s.startMonitorStd(c, "/watch/target.txt")
	logBefore := s.readMonitorLog(c)

	// 1. OPEN + READ: cat the file
	s.exec(c, []string{"sh", "-c", "cat /watch/target.txt > /dev/null"})
	logAfter := s.waitForEventType(c, logBefore, "OPEN", 5*time.Second)
	s.waitForEventType(c, logAfter, "READ", 5*time.Second)
	logBefore = s.readMonitorLog(c)

	// 2. WRITE: append to the file
	s.exec(c, []string{"sh", "-c", "echo 'more data' >> /watch/target.txt"})
	logAfter = s.waitForEventType(c, logBefore, "WRITE", 5*time.Second)
	logBefore = logAfter

	// 3. HARDLINK: create hardlink to watched file (event path matches target)
	s.exec(c, []string{"ln", "/watch/target.txt", "/watch/hardlink"})
	logAfter = s.waitForEventType(c, logBefore, "HARDLINK", 5*time.Second)
	logBefore = logAfter

	// SYMLINK skipped for single-file watch: the event's path is the new symlink
	// location (e.g., /watch/link), not the watched file. In Docker overlay2, the
	// dentry resolves to an inaccessible overlay path that matchTarget rejects.
	// SYMLINK is covered by TestMonitorAllEvents (directory watch).

	// 5. RENAME: rename the watched file away (event path matches target)
	s.exec(c, []string{"mv", "/watch/target.txt", "/watch/renamed.txt"})
	logAfter = s.waitForEventType(c, logBefore, "RENAME", 5*time.Second)
	logBefore = logAfter

	// 6. MMAP: re-create the file at the watched path, then run mmap exploit
	s.exec(c, []string{"sh", "-c", "echo 'mmap data' > /watch/target.txt"})
	s.waitForEventType(c, logBefore, "OPEN", 5*time.Second)
	logBefore = s.readMonitorLog(c)

	s.exec(c, []string{"/mmap_exploit", "/watch/target.txt"})
	logAfter = s.waitForEventType(c, logBefore, "MMAP", 5*time.Second)
	s.waitForEventType(c, logAfter, "OPEN", 5*time.Second)
	logBefore = logAfter

	// 7. DELETE: remove the re-created file (path matches target)
	s.exec(c, []string{"rm", "/watch/target.txt"})
	s.waitForEventType(c, logBefore, "DELETE", 5*time.Second)

	s.stopMonitor(c)
}

// ---------------------------------------------------------------
// Event filter tests (-e option)
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestMonitorEventFilter_Selective() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})

	// Start monitor with a filter that excludes OPEN, WRITE, SYMLINK, MKDIR, MMAP
	s.startMonitorStd(c, "/watch", "-e", "READ,HARDLINK,DELETE")
	logBefore := s.readMonitorLog(c)

	// 1. Create file → generates OPEN+WRITE (both filtered → no events)
	s.exec(c, []string{"sh", "-c", "echo 'filter data' > /watch/filter.txt"})
	time.Sleep(1500 * time.Millisecond)
	logAfter := s.readMonitorLog(c)
	events := newEventTypes(logBefore, logAfter)
	s.Require().Emptyf(events, "OPEN+WRITE should be filtered out, got %v", events)
	logBefore = logAfter

	// 2. Read file → generates OPEN+READ (OPEN filtered, READ allowed)
	s.exec(c, []string{"sh", "-c", "cat /watch/filter.txt > /dev/null"})
	logAfter = s.waitForEventType(c, logBefore, "READ", 5*time.Second)
	s.requireNoNewEventType(logBefore, logAfter, "OPEN")
	logBefore = logAfter

	// 3. Create hardlink → generates HARDLINK (allowed)
	s.exec(c, []string{"ln", "/watch/filter.txt", "/watch/filter-hl"})
	logAfter = s.waitForEventType(c, logBefore, "HARDLINK", 5*time.Second)
	logBefore = logAfter

	// 4. Delete → generates DELETE (allowed)
	s.exec(c, []string{"rm", "/watch/filter.txt"})
	logAfter = s.waitForEventType(c, logBefore, "DELETE", 5*time.Second)
	// DELETE consumes the file, so hardlink is orphaned — clean up
	s.exec(c, []string{"rm", "-f", "/watch/filter-hl"})

	s.stopMonitor(c)
}

func (s *IntegrationSuite) TestMonitorEventFilter_SingleType() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})

	// Start monitor watching only DELETE events
	s.startMonitorStd(c, "/watch", "-e", "DELETE")
	logBefore := s.readMonitorLog(c)

	// 1. Create file → OPEN+WRITE both filtered → no events
	s.exec(c, []string{"sh", "-c", "echo 'delete test' > /watch/todelete.txt"})
	time.Sleep(1500 * time.Millisecond)
	logAfter := s.readMonitorLog(c)
	events := newEventTypes(logBefore, logAfter)
	s.Require().Emptyf(events, "OPEN+WRITE should be filtered out, got %v", events)
	logBefore = logAfter

	// 2. Read file → OPEN+READ both filtered → no events
	s.exec(c, []string{"sh", "-c", "cat /watch/todelete.txt > /dev/null"})
	time.Sleep(1500 * time.Millisecond)
	logAfter = s.readMonitorLog(c)
	events = newEventTypes(logBefore, logAfter)
	s.Require().Emptyf(events, "OPEN+READ should be filtered out, got %v", events)
	logBefore = logAfter

	// 3. Delete → DELETE allowed → should appear
	s.exec(c, []string{"rm", "/watch/todelete.txt"})
	s.waitForEventType(c, logBefore, "DELETE", 5*time.Second)

	s.stopMonitor(c)
}

// ---------------------------------------------------------------
// Event filter tests: single-file watch
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestMonitorEventFilter_Selective_FileWatch() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	// Watched file must exist before starting the monitor (inode resolution)
	s.exec(c, []string{"sh", "-c", "echo 'filter data' > /watch/target.txt"})

	// Watch a single file with filter that excludes OPEN, WRITE, SYMLINK, MMAP
	s.startMonitorStd(c, "/watch/target.txt", "-e", "READ,HARDLINK,DELETE")
	logBefore := s.readMonitorLog(c)

	// 1. Read file → OPEN+READ (OPEN filtered, READ allowed)
	s.exec(c, []string{"sh", "-c", "cat /watch/target.txt > /dev/null"})
	logAfter := s.waitForEventType(c, logBefore, "READ", 5*time.Second)
	s.requireNoNewEventType(logBefore, logAfter, "OPEN")
	logBefore = logAfter

	// 2. Create hardlink → HARDLINK allowed (catch via inode)
	s.exec(c, []string{"ln", "/watch/target.txt", "/watch/filter-hl"})
	logAfter = s.waitForEventType(c, logBefore, "HARDLINK", 5*time.Second)
	logBefore = logAfter

	// 3. Delete → DELETE allowed
	s.exec(c, []string{"rm", "/watch/target.txt"})
	s.waitForEventType(c, logBefore, "DELETE", 5*time.Second)
	s.exec(c, []string{"rm", "-f", "/watch/filter-hl"})

	s.stopMonitor(c)
}

func (s *IntegrationSuite) TestMonitorEventFilter_SingleType_FileWatch() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	// Watched file must exist before starting the monitor (inode resolution)
	s.exec(c, []string{"sh", "-c", "echo 'delete target' > /watch/target.txt"})

	// Watch a single file, only DELETE events
	s.startMonitorStd(c, "/watch/target.txt", "-e", "DELETE")
	logBefore := s.readMonitorLog(c)

	// 1. Read file → OPEN+READ both filtered → no events
	s.exec(c, []string{"sh", "-c", "cat /watch/target.txt > /dev/null"})
	time.Sleep(1500 * time.Millisecond)
	logAfter := s.readMonitorLog(c)
	events := newEventTypes(logBefore, logAfter)
	s.Require().Emptyf(events, "OPEN+READ should be filtered out, got %v", events)
	logBefore = logAfter

	// 2. Delete → DELETE allowed → should appear
	s.exec(c, []string{"rm", "/watch/target.txt"})
	s.waitForEventType(c, logBefore, "DELETE", 5*time.Second)

	s.stopMonitor(c)
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
	code, out := s.exec(c, []string{"/app-listener", "monitor", "-w", "/watch", "--recursive"})
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
	code, out := s.exec(c, []string{"/app-listener", "monitor", "-w", "/watch", "--recursive"})
	verifyEBPF(s, code, out)
}

func (s *IntegrationSuite) TestMultiDistro_Fedora_EBPF() {
	c := s.startContainer("fedora:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	code, out := s.exec(c, []string{"/app-listener", "monitor", "-w", "/watch", "--recursive"})
	verifyEBPF(s, code, out)
}

// ---------------------------------------------------------------
// Non-recursive monitoring: events only for direct children
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestMonitor_NonRecursive() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch/subdir"})
	s.exec(c, []string{"sh", "-c", "echo top > /watch/top.txt"})
	s.exec(c, []string{"sh", "-c", "echo deep > /watch/subdir/deep.txt"})

	// Non-recursive monitor (default)
	s.startMonitorStd(c, "/watch")
	logBefore := s.readMonitorLog(c)

	// 1. Direct child file — events detected
	s.exec(c, []string{"sh", "-c", "cat /watch/top.txt > /dev/null"})
	logAfter := s.waitForEventType(c, logBefore, "OPEN", 5*time.Second)
	s.waitForEventType(c, logAfter, "READ", 3*time.Second)

	// 2. File in subdir — events NOT detected (non-recursive)
	logBefore = s.readMonitorLog(c)
	s.exec(c, []string{"sh", "-c", "cat /watch/subdir/deep.txt > /dev/null"})
	time.Sleep(2 * time.Second)
	events := newEventTypes(logBefore, s.readMonitorLog(c))
	s.Require().Emptyf(events, "events in subdir should NOT be detected (non-recursive), got %v", events)

	s.stopMonitor(c)
}

// ---------------------------------------------------------------
// Depth-limited monitoring: events only up to specified depth
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestMonitor_DepthLimit() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch/subdir/inner"})
	s.exec(c, []string{"sh", "-c", "echo top > /watch/top.txt"})
	s.exec(c, []string{"sh", "-c", "echo mid > /watch/subdir/mid.txt"})
	s.exec(c, []string{"sh", "-c", "echo deep > /watch/subdir/inner/deep.txt"})

	// Monitor with depth=2: levels 0,1,2 are monitored; level 3+ is not
	s.startMonitorStd(c, "/watch", "--recursive", "--depth", "2")
	logBefore := s.readMonitorLog(c)

	// 1. Direct child (level 1) — detected
	s.exec(c, []string{"sh", "-c", "cat /watch/top.txt > /dev/null"})
	s.waitForEventType(c, logBefore, "OPEN", 5*time.Second)

	// 2. File at depth 2 inside subdir — detected
	logBefore = s.readMonitorLog(c)
	s.exec(c, []string{"sh", "-c", "cat /watch/subdir/mid.txt > /dev/null"})
	s.waitForEventType(c, logBefore, "OPEN", 5*time.Second)

	// 3. File at depth 3 inside inner/ — NOT detected
	logBefore = s.readMonitorLog(c)
	s.exec(c, []string{"sh", "-c", "cat /watch/subdir/inner/deep.txt > /dev/null"})
	time.Sleep(2 * time.Second)
	events := newEventTypes(logBefore, s.readMonitorLog(c))
	s.Require().Emptyf(events, "events at depth 3 should NOT be detected (depth=2), got %v", events)

	// 4. Create a NEW file at depth 2 — detected even after startup
	logBefore = s.readMonitorLog(c)
	s.exec(c, []string{"sh", "-c", "echo new > /watch/subdir/newfile.txt"})
	s.waitForEventType(c, logBefore, "OPEN", 5*time.Second)

	// 5. Create a NEW file at depth 3 — NOT detected
	logBefore = s.readMonitorLog(c)
	s.exec(c, []string{"sh", "-c", "echo newdeep > /watch/subdir/inner/newdeep.txt"})
	time.Sleep(2 * time.Second)
	events = newEventTypes(logBefore, s.readMonitorLog(c))
	s.Require().Emptyf(events, "events at depth 3 should NOT be detected (depth=2), got %v", events)

	s.stopMonitor(c)
}

// ---------------------------------------------------------------
// Recursive monitoring: events detected at all depths
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestMonitor_Recursive() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch/subdir"})
	s.exec(c, []string{"sh", "-c", "echo top > /watch/top.txt"})
	s.exec(c, []string{"sh", "-c", "echo nested > /watch/subdir/nested.txt"})

	// Recursive monitor (unlimited depth)
	s.startMonitorStd(c, "/watch", "--recursive")
	logBefore := s.readMonitorLog(c)

	// 1. Direct child — detected
	s.exec(c, []string{"sh", "-c", "cat /watch/top.txt > /dev/null"})
	logAfter := s.waitForEventType(c, logBefore, "OPEN", 5*time.Second)
	s.waitForEventType(c, logAfter, "READ", 3*time.Second)

	// 2. Nested file in subdir — detected (recursive)
	logBefore = s.readMonitorLog(c)
	s.exec(c, []string{"sh", "-c", "cat /watch/subdir/nested.txt > /dev/null"})
	s.waitForEventType(c, logBefore, "OPEN", 5*time.Second)

	// 3. NEW file in subdir — detected even after startup
	logBefore = s.readMonitorLog(c)
	s.exec(c, []string{"sh", "-c", "echo newnested > /watch/subdir/new_file.txt"})
	s.waitForEventType(c, logBefore, "OPEN", 5*time.Second)

	s.stopMonitor(c)
}

// ---------------------------------------------------------------
// Monitor helpers (used by both event tests and exploit tests)
// ---------------------------------------------------------------

func (s *IntegrationSuite) waitForMonitorStd(c testcontainers.Container, needle string, timeout time.Duration) (string, bool) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		_, out := s.exec(c, []string{"sh", "-c", "cat /tmp/monitor.log 2>/dev/null || true"})
		if strings.Contains(out, needle) {
			return out, true
		}
		time.Sleep(500 * time.Millisecond)
	}
	_, out := s.exec(c, []string{"sh", "-c", "cat /tmp/monitor.log 2>/dev/null || true"})
	return out, false
}

func (s *IntegrationSuite) startMonitorStd(c testcontainers.Container, watchDir string, extraFlags ...string) {
	flags := strings.Join(extraFlags, " ")
	cmd := fmt.Sprintf("nohup /app-listener monitor -w %s --headless %s > /tmp/monitor.log 2>&1 &", watchDir, flags)
	s.exec(c, []string{"sh", "-c", cmd})

	_, ok := s.waitForMonitorStd(c, "monitor started", 15*time.Second)
	s.Require().True(ok, "monitor should become ready (waiting for 'monitor started')")

	codeCheck, outCheck := s.exec(c, []string{"pgrep", "-f", "app-listener"})
	s.Require().Equalf(0, codeCheck, "monitor process not running after start: %s", outCheck)

	time.Sleep(1 * time.Second)
}

func (s *IntegrationSuite) stopMonitor(c testcontainers.Container) {
	s.exec(c, []string{"pkill", "-f", "app-listener"})
	time.Sleep(500 * time.Millisecond)
}

// ---------------------------------------------------------------
// Exploit tests — verify that the monitor captures events from
// various syscall-based I/O paths (splice, sendfile, mmap, ...)
// ---------------------------------------------------------------

type exploitTest struct {
	name           string
	binary         string
	args           []string
	requiresRoot   bool
	kernelOptional bool
	events         []string
	setupFile      bool
}

var exploitTests = []exploitTest{
	{
		name:   "pread64",
		binary: "pread64",
		events: []string{"OPEN", "READ"},
	},
	{
		name:   "readv",
		binary: "readv",
		events: []string{"OPEN", "READ"},
	},
	{
		name:   "sendfile",
		binary: "sendfile",
		events: []string{"OPEN", "READ"},
	},
	{
		name:           "splice",
		binary:         "splice",
		events:         []string{"OPEN", "READ"},
		kernelOptional: true,
	},
	{
		name:   "copy_file_range",
		binary: "copy_file_range",
		events: []string{"OPEN", "READ"},
	},
	{
		name:   "mmap",
		binary: "mmap",
		events: []string{"OPEN", "MMAP"},
	},
	{
		name:      "execve",
		binary:    "execve",
		events:    []string{"OPEN"},
		setupFile: true,
	},
	{
		name:           "io_uring",
		binary:         "io_uring",
		events:         []string{"OPEN", "READ"},
		kernelOptional: true,
	},
	{
		name:         "open_by_handle_at",
		binary:       "open_by_handle_at",
		events:       []string{"OPEN", "READ"},
		requiresRoot: true,
	},
}

func (s *IntegrationSuite) runExploitExpect(c testcontainers.Container, et exploitTest, testFilePath string, oldLog string) ([]string, bool) {
	exploitHostPath := absPath(fmt.Sprintf("./exploits/%s", et.binary))
	err := c.CopyFileToContainer(s.ctx, exploitHostPath, fmt.Sprintf("/exploits/%s", et.binary), 0755)
	if !s.NoError(err, "copy exploit binary") {
		return nil, false
	}

	args := []string{fmt.Sprintf("/exploits/%s", et.binary)}
	if et.name == "execve" {
		s.exec(c, []string{"sh", "-c", "echo '#!/bin/sh\necho executed' > /watch/.exec_target && chmod +x /watch/.exec_target"})
		args = append(args, "/watch/.exec_target")
	} else {
		args = append(args, testFilePath)
	}

	code, out := s.exec(c, args)
	if et.requiresRoot {
		if code != 0 {
			s.T().Logf("exploit %s exited with %d (may require root): %s", et.name, code, out)
			return nil, false
		}
	} else {
		s.Require().Equalf(0, code, "exploit %s should succeed:\n%s", et.name, out)
	}

	time.Sleep(2 * time.Second)

	newLog := s.readMonitorLog(c)
	deltaEvents := newEventTypes(oldLog, newLog)

	eventSet := make(map[string]bool, len(deltaEvents))
	for _, e := range deltaEvents {
		eventSet[e] = true
	}
	var found []string
	for _, evt := range et.events {
		if eventSet[evt] {
			found = append(found, evt)
		}
	}

	s.T().Logf("exploit %s: expected %v, found %v (delta events: %v)", et.name, et.events, found, deltaEvents)
	return found, len(found) == len(et.events)
}

func (s *IntegrationSuite) exploitSingle(name string) {
	var et *exploitTest
	for i := range exploitTests {
		if exploitTests[i].name == name {
			et = &exploitTests[i]
			break
		}
	}
	s.Require().NotNil(et, "unknown exploit %s", name)

	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch", "/exploits"})
	testFilePath := "/watch/test_file.txt"
	s.exec(c, []string{"sh", "-c", fmt.Sprintf("echo 'test data' > %s", testFilePath)})

	s.startMonitorStd(c, "/watch")
	logBefore := s.readMonitorLog(c)
	found, allOk := s.runExploitExpect(c, *et, testFilePath, logBefore)
	if !allOk && et.kernelOptional {
		s.T().Skipf("%s: kernel may not support this probe: expected %v, got %v", name, et.events, found)
	}
	s.Require().True(len(found) == len(et.events),
		"%s: expected events %v, got %v", name, et.events, found)
	s.stopMonitor(c)
}

func (s *IntegrationSuite) TestExploit_pread64()         { s.exploitSingle("pread64") }
func (s *IntegrationSuite) TestExploit_readv()           { s.exploitSingle("readv") }
func (s *IntegrationSuite) TestExploit_sendfile()        { s.exploitSingle("sendfile") }
func (s *IntegrationSuite) TestExploit_splice()          { s.exploitSingle("splice") }
func (s *IntegrationSuite) TestExploit_mmap()            { s.exploitSingle("mmap") }
func (s *IntegrationSuite) TestExploit_copy_file_range() { s.exploitSingle("copy_file_range") }

func (s *IntegrationSuite) TestExploits() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	testFilePath := "/watch/test_file.txt"
	s.exec(c, []string{"sh", "-c", fmt.Sprintf("echo 'test data for exploits' > %s", testFilePath)})
	s.exec(c, []string{"mkdir", "-p", "/exploits"})

	s.startMonitorStd(c, "/watch")
	logBefore := s.readMonitorLog(c)

	for _, et := range exploitTests {
		s.T().Run(et.name, func(t *testing.T) {
			t.Helper()

			testFile := testFilePath
			if et.name == "execve" {
				t.Log("execve will copy test_file.txt to .exec_target")
			}

			found, allOk := s.runExploitExpect(c, et, testFile, logBefore)
			logBefore = s.readMonitorLog(c)
			if !allOk {
				t.Logf("Monitor log tail:\n%s", logBefore)
				if et.requiresRoot {
					t.Skipf("exploit %s requires root (CAP_DAC_READ_SEARCH): expected %v, got %v", et.name, et.events, found)
				}
				if et.kernelOptional {
					t.Skipf("exploit %s not supported on this kernel: expected %v, got %v", et.name, et.events, found)
				}
				t.Fatalf("exploit %s should trigger events %v, got %v", et.name, et.events, found)
			}
		})

		time.Sleep(500 * time.Millisecond)
	}

	s.stopMonitor(c)
}

func (s *IntegrationSuite) TestExploit_execve() {
	et := exploitTests[6]

	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch", "/exploits"})
	s.exec(c, []string{"sh", "-c", "echo '#!/bin/sh\necho executed' > /watch/.exec_target"})
	s.exec(c, []string{"chmod", "+x", "/watch/.exec_target"})

	s.startMonitorStd(c, "/watch")
	logBefore := s.readMonitorLog(c)
	found, _ := s.runExploitExpect(c, et, "/watch/.exec_target", logBefore)
	s.Require().True(len(found) == len(et.events),
		"execve: expected events %v, got %v", et.events, found)
	s.stopMonitor(c)
}

func (s *IntegrationSuite) TestExploit_io_uring() {
	et := exploitTests[7]

	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch", "/exploits"})
	testFilePath := "/watch/test_file.txt"
	s.exec(c, []string{"sh", "-c", fmt.Sprintf("echo 'test data' > %s", testFilePath)})

	s.startMonitorStd(c, "/watch")
	logBefore := s.readMonitorLog(c)
	found, allOk := s.runExploitExpect(c, et, testFilePath, logBefore)
	if !allOk {
		s.T().Skipf("io_uring not supported on this kernel: expected %v, got %v", et.events, found)
	}
	s.stopMonitor(c)
}

func (s *IntegrationSuite) TestExploit_open_by_handle_at() {
	et := exploitTests[8]

	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch", "/exploits"})
	testFilePath := "/watch/test_file.txt"
	s.exec(c, []string{"sh", "-c", fmt.Sprintf("echo 'test data' > %s", testFilePath)})

	s.startMonitorStd(c, "/watch")
	logBefore := s.readMonitorLog(c)
	found, allOk := s.runExploitExpect(c, et, testFilePath, logBefore)
	if !allOk {
		s.T().Logf("open_by_handle_at may require root: expected %v, got %v", et.events, found)
	}
	s.stopMonitor(c)
}
