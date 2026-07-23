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
	time.Sleep(500 * time.Millisecond)
	logAfter := s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "OPEN")
	s.requireNewEventType(logBefore, logAfter, "WRITE")
	logBefore = logAfter

	// 2. READ: cat the file
	s.exec(c, []string{"sh", "-c", "cat /watch/data.txt > /dev/null"})
	time.Sleep(500 * time.Millisecond)
	logAfter = s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "READ")
	logBefore = logAfter

	// 3. RENAME: rename file
	s.exec(c, []string{"mv", "/watch/data.txt", "/watch/renamed.txt"})
	time.Sleep(500 * time.Millisecond)
	logAfter = s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "RENAME")
	logBefore = logAfter

	// 4. SYMLINK: create symlink
	s.exec(c, []string{"ln", "-s", "/watch/renamed.txt", "/watch/link"})
	time.Sleep(500 * time.Millisecond)
	logAfter = s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "SYMLINK")
	logBefore = logAfter

	// 5. HARDLINK: create hardlink
	s.exec(c, []string{"ln", "/watch/renamed.txt", "/watch/hardlink"})
	time.Sleep(500 * time.Millisecond)
	logAfter = s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "HARDLINK")
	logBefore = logAfter

	// 6. MKDIR: create directory
	s.exec(c, []string{"mkdir", "/watch/subdir"})
	time.Sleep(500 * time.Millisecond)
	logAfter = s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "MKDIR")
	logBefore = logAfter

	// 7. DELETE: remove file
	s.exec(c, []string{"rm", "/watch/renamed.txt"})
	time.Sleep(500 * time.Millisecond)
	logAfter = s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "DELETE")
	logBefore = logAfter

	// 8. MMAP + OPEN: run mmap exploit on hardlink (still exists after rm)
	s.exec(c, []string{"/mmap_exploit", "/watch/hardlink"})
	time.Sleep(2 * time.Second)
	logAfter = s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "MMAP")
	s.requireNewEventType(logBefore, logAfter, "OPEN")

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
	time.Sleep(500 * time.Millisecond)
	logAfter := s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "OPEN")
	s.requireNewEventType(logBefore, logAfter, "READ")
	logBefore = logAfter

	// 2. WRITE: append to the file
	s.exec(c, []string{"sh", "-c", "echo 'more data' >> /watch/target.txt"})
	time.Sleep(500 * time.Millisecond)
	logAfter = s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "WRITE")
	logBefore = logAfter

	// 3. HARDLINK: create hardlink to watched file (event path matches target)
	s.exec(c, []string{"ln", "/watch/target.txt", "/watch/hardlink"})
	time.Sleep(500 * time.Millisecond)
	logAfter = s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "HARDLINK")
	logBefore = logAfter

	// SYMLINK skipped for single-file watch: the event's path is the new symlink
	// location (e.g., /watch/link), not the watched file. In Docker overlay2, the
	// dentry resolves to an inaccessible overlay path that matchTarget rejects.
	// SYMLINK is covered by TestMonitorAllEvents (directory watch).

	// 5. RENAME: rename the watched file away (event path matches target)
	s.exec(c, []string{"mv", "/watch/target.txt", "/watch/renamed.txt"})
	time.Sleep(500 * time.Millisecond)
	logAfter = s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "RENAME")
	logBefore = logAfter

	// 6. MMAP: re-create the file at the watched path, then run mmap exploit
	s.exec(c, []string{"sh", "-c", "echo 'mmap data' > /watch/target.txt"})
	time.Sleep(500 * time.Millisecond)
	logAfter = s.readMonitorLog(c)
	// Ignore OPEN+WRITE from re-creation, just advance the log cursor
	logBefore = logAfter

	s.exec(c, []string{"/mmap_exploit", "/watch/target.txt"})
	time.Sleep(2 * time.Second)
	logAfter = s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "MMAP")
	s.requireNewEventType(logBefore, logAfter, "OPEN")
	logBefore = logAfter

	// 7. DELETE: remove the re-created file (path matches target)
	s.exec(c, []string{"rm", "/watch/target.txt"})
	time.Sleep(500 * time.Millisecond)
	logAfter = s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "DELETE")

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
	time.Sleep(500 * time.Millisecond)
	logAfter := s.readMonitorLog(c)
	events := newEventTypes(logBefore, logAfter)
	s.Require().Emptyf(events, "OPEN+WRITE should be filtered out, got %v", events)
	logBefore = logAfter

	// 2. Read file → generates OPEN+READ (OPEN filtered, READ allowed)
	s.exec(c, []string{"sh", "-c", "cat /watch/filter.txt > /dev/null"})
	time.Sleep(500 * time.Millisecond)
	logAfter = s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "READ")
	s.requireNoNewEventType(logBefore, logAfter, "OPEN")
	logBefore = logAfter

	// 3. Create hardlink → generates HARDLINK (allowed)
	s.exec(c, []string{"ln", "/watch/filter.txt", "/watch/filter-hl"})
	time.Sleep(500 * time.Millisecond)
	logAfter = s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "HARDLINK")
	logBefore = logAfter

	// 4. Delete → generates DELETE (allowed)
	s.exec(c, []string{"rm", "/watch/filter.txt"})
	time.Sleep(500 * time.Millisecond)
	logAfter = s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "DELETE")
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
	time.Sleep(500 * time.Millisecond)
	logAfter := s.readMonitorLog(c)
	events := newEventTypes(logBefore, logAfter)
	s.Require().Emptyf(events, "OPEN+WRITE should be filtered out, got %v", events)
	logBefore = logAfter

	// 2. Read file → OPEN+READ both filtered → no events
	s.exec(c, []string{"sh", "-c", "cat /watch/todelete.txt > /dev/null"})
	time.Sleep(500 * time.Millisecond)
	logAfter = s.readMonitorLog(c)
	events = newEventTypes(logBefore, logAfter)
	s.Require().Emptyf(events, "OPEN+READ should be filtered out, got %v", events)
	logBefore = logAfter

	// 3. Delete → DELETE allowed → should appear
	s.exec(c, []string{"rm", "/watch/todelete.txt"})
	time.Sleep(500 * time.Millisecond)
	logAfter = s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "DELETE")

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
	time.Sleep(500 * time.Millisecond)
	logAfter := s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "READ")
	s.requireNoNewEventType(logBefore, logAfter, "OPEN")
	logBefore = logAfter

	// 2. Create hardlink → HARDLINK allowed (catch via inode)
	s.exec(c, []string{"ln", "/watch/target.txt", "/watch/filter-hl"})
	time.Sleep(500 * time.Millisecond)
	logAfter = s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "HARDLINK")
	logBefore = logAfter

	// 3. Delete → DELETE allowed
	s.exec(c, []string{"rm", "/watch/target.txt"})
	time.Sleep(500 * time.Millisecond)
	logAfter = s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "DELETE")
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
	time.Sleep(500 * time.Millisecond)
	logAfter := s.readMonitorLog(c)
	events := newEventTypes(logBefore, logAfter)
	s.Require().Emptyf(events, "OPEN+READ should be filtered out, got %v", events)
	logBefore = logAfter

	// 2. Delete → DELETE allowed → should appear
	s.exec(c, []string{"rm", "/watch/target.txt"})
	time.Sleep(500 * time.Millisecond)
	logAfter = s.readMonitorLog(c)
	s.requireNewEventType(logBefore, logAfter, "DELETE")

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
