package integrationtests

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

// ---------------------------------------------------------------
// Guard helpers
// ---------------------------------------------------------------

type guardEvent struct {
	Type    string
	Comm    string
	Path    string
	PID     int
	Dest    string
	Blocked bool
	UID     int
}

func parseGuardEvents(logContent string) []guardEvent {
	var events []guardEvent
	for _, line := range strings.Split(logContent, "\n") {
		idx := strings.Index(line, "GUARD|")
		if idx < 0 {
			continue
		}
		rest := line[idx+6:]
		parts := strings.SplitN(rest, "|", 7)
		if len(parts) < 7 {
			continue
		}
		blocked, _ := strconv.ParseBool(parts[5])
		pid, _ := strconv.Atoi(parts[3])
		uid, _ := strconv.Atoi(parts[6])
		events = append(events, guardEvent{
			Type:    parts[0],
			Comm:    parts[1],
			Path:    parts[2],
			PID:     pid,
			Dest:    parts[4],
			Blocked: blocked,
			UID:     uid,
		})
	}
	return events
}

// guardDeltaEvents returns guard events that appear in newLog but not in oldLog,
// using line-based diff.
func guardDeltaEvents(oldLog, newLog string) []guardEvent {
	oldEvents := parseGuardEvents(oldLog)
	newEvents := parseGuardEvents(newLog)

	if len(newEvents) <= len(oldEvents) {
		return nil
	}
	return newEvents[len(oldEvents):]
}

func guardEventTypes(events []guardEvent) []string {
	seen := make(map[string]bool)
	var types []string
	for _, e := range events {
		if !seen[e.Type] {
			seen[e.Type] = true
			types = append(types, e.Type)
		}
	}
	return types
}

func (s *IntegrationSuite) startGuardStd(c testcontainers.Container, guardPath string, extraFlags ...string) {
	flags := strings.Join(extraFlags, " ")
	cmd := fmt.Sprintf("nohup /app-listener guard %s --headless %s > /tmp/guard.log 2>&1 &", guardPath, flags)
	code, out := s.exec(c, []string{"sh", "-c", cmd})
	s.Require().Equalf(0, code, "starting guard: %s", out)

	_, ok := s.waitForGuardStd(c, "guard started", 10*time.Second)
	s.Require().True(ok, "guard should become ready")

	// Verify the process is alive
	codeCheck, outCheck := s.exec(c, []string{"pgrep", "-f", "app-listener"})
	s.Require().Equalf(0, codeCheck, "guard process not running after start: %s", outCheck)

	time.Sleep(2 * time.Second)
}

func (s *IntegrationSuite) waitForGuardStd(c testcontainers.Container, needle string, timeout time.Duration) (string, bool) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		_, out := s.exec(c, []string{"sh", "-c", "cat /tmp/guard.log 2>/dev/null || true"})
		if strings.Contains(out, needle) {
			return out, true
		}
		time.Sleep(500 * time.Millisecond)
	}
	_, out := s.exec(c, []string{"sh", "-c", "cat /tmp/guard.log 2>/dev/null || true"})
	return out, false
}

func (s *IntegrationSuite) readGuardLog(c testcontainers.Container) string {
	_, out := s.exec(c, []string{"sh", "-c", "cat /tmp/guard.log 2>/dev/null || true"})
	return out
}

func (s *IntegrationSuite) stopGuard(c testcontainers.Container) {
	s.exec(c, []string{"pkill", "-f", "app-listener"})
	time.Sleep(500 * time.Millisecond)
}

func (s *IntegrationSuite) requireBlockedEvent(events []guardEvent, expectedType string, comms ...string) {
	for _, e := range events {
		if e.Type == expectedType && e.Blocked {
			if len(comms) == 0 {
				return
			}
			if slices.Contains(comms, e.Comm) {
				return
			}
		}
	}
	commsStr := ""
	if len(comms) > 0 {
		commsStr = fmt.Sprintf(" (comm in %v)", comms)
	}
	s.Require().Failf("missing blocked event",
		"expected blocked GUARD|%s%s, got %v", expectedType, commsStr, events)
}

func (s *IntegrationSuite) requireNoBlockedEvent(events []guardEvent, unexpectedType string) {
	for _, e := range events {
		if e.Type == unexpectedType && e.Blocked {
			s.Require().Failf("unexpected blocked event",
				"got unexpected blocked GUARD|%s (from %s)", unexpectedType, e.Comm)
		}
	}
}

// ---------------------------------------------------------------
// Test: guard blocks all operations on a directory
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestGuard_BlocksAll_Directory() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	// Pre-create files so their inodes are in the guard map
	s.exec(c, []string{"touch", "/watch/guarded.txt"})
	s.exec(c, []string{"mkdir", "-p", "/watch/subdir"})

	// Guard blocks all (whitelist mode with empty list → all blocked)
	s.startGuardStd(c, "/watch")
	logBefore := s.readGuardLog(c)

	// 1. cat → OPEN blocked
	code, out := s.exec(c, []string{"sh", "-c", "cat /watch/guarded.txt > /dev/null 2>&1"})
	s.Require().NotEqualf(0, code, "cat should be blocked: %s", out)

	// 2. append → WRITE blocked
	code, out = s.exec(c, []string{"sh", "-c", "echo 'data' >> /watch/guarded.txt 2>&1"})
	s.Require().NotEqualf(0, code, "append should be blocked: %s", out)

	// 3. delete → DELETE blocked
	code, out = s.exec(c, []string{"rm", "/watch/guarded.txt"})
	s.Require().NotEqualf(0, code, "rm should be blocked: %s", out)

	// 4. rename → RENAME blocked
	code, out = s.exec(c, []string{"mv", "/watch/guarded.txt", "/watch/renamed.txt"})
	s.Require().NotEqualf(0, code, "mv should be blocked: %s", out)

	// 5. mkdir → MKDIR blocked
	code, out = s.exec(c, []string{"mkdir", "/watch/subdir2"})
	s.Require().NotEqualf(0, code, "mkdir should be blocked: %s", out)

	// 6. hardlink → HARDLINK blocked
	code, out = s.exec(c, []string{"ln", "/watch/guarded.txt", "/watch/hardlink"})
	s.Require().NotEqualf(0, code, "ln should be blocked: %s", out)

	logAfter := s.readGuardLog(c)
	deltaEvents := guardDeltaEvents(logBefore, logAfter)

	s.Require().NotEmpty(deltaEvents, "expected blocked events in guard log")
	s.requireBlockedEvent(deltaEvents, "OPEN")
	s.requireBlockedEvent(deltaEvents, "DELETE")
	s.requireBlockedEvent(deltaEvents, "RENAME")
	s.requireBlockedEvent(deltaEvents, "MKDIR")
	s.requireBlockedEvent(deltaEvents, "HARDLINK")

	// Verify all events are blocked
	for _, e := range deltaEvents {
		s.Require().Truef(e.Blocked, "event %s|%s should be blocked", e.Type, e.Comm)
	}

	s.stopGuard(c)
}

func (s *IntegrationSuite) TestGuard_BlocksAll_File() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	// Pre-create the guarded file
	s.exec(c, []string{"sh", "-c", "echo 'guarded content' > /watch/target.txt"})

	s.startGuardStd(c, "/watch/target.txt")
	logBefore := s.readGuardLog(c)

	// 1. cat → OPEN blocked
	code, out := s.exec(c, []string{"sh", "-c", "cat /watch/target.txt > /dev/null 2>&1"})
	s.Require().NotEqualf(0, code, "cat should be blocked: %s", out)

	// 2. append → WRITE blocked
	code, out = s.exec(c, []string{"sh", "-c", "echo 'data' >> /watch/target.txt 2>&1"})
	s.Require().NotEqualf(0, code, "append should be blocked: %s", out)

	// 3. delete → DELETE blocked
	code, out = s.exec(c, []string{"rm", "/watch/target.txt"})
	s.Require().NotEqualf(0, code, "rm should be blocked: %s", out)

	// 4. rename → RENAME blocked
	code, out = s.exec(c, []string{"mv", "/watch/target.txt", "/watch/renamed.txt"})
	s.Require().NotEqualf(0, code, "mv should be blocked: %s", out)

	// 5. hardlink → HARDLINK blocked
	code, out = s.exec(c, []string{"ln", "/watch/target.txt", "/watch/hardlink"})
	s.Require().NotEqualf(0, code, "ln should be blocked: %s", out)

	logAfter := s.readGuardLog(c)
	deltaEvents := guardDeltaEvents(logBefore, logAfter)

	s.Require().NotEmpty(deltaEvents, "expected blocked events in guard log")
	s.requireBlockedEvent(deltaEvents, "OPEN")
	s.requireBlockedEvent(deltaEvents, "DELETE")
	s.requireBlockedEvent(deltaEvents, "RENAME")
	s.requireBlockedEvent(deltaEvents, "HARDLINK")

	for _, e := range deltaEvents {
		s.Require().Truef(e.Blocked, "event %s|%s should be blocked", e.Type, e.Comm)
	}

	s.stopGuard(c)
}

// ---------------------------------------------------------------
// Test: binary whitelist — only /usr/bin/cat is allowed
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestGuard_Whitelist_Binary() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	s.exec(c, []string{"sh", "-c", "echo 'whitelist test' > /watch/data.txt"})

	// Whitelist /usr/bin/cat: only cat (comm=cat) is allowed
	s.startGuardStd(c, "/watch", "-w", "/usr/bin/cat")

	// 1. cat is whitelisted → allowed
	code, out := s.exec(c, []string{"sh", "-c", "cat /watch/data.txt > /dev/null 2>&1"})
	s.Require().Equalf(0, code, "cat should be allowed: %s", out)

	// 2. rm is NOT whitelisted → blocked
	code, out = s.exec(c, []string{"rm", "/watch/data.txt"})
	s.Require().NotEqualf(0, code, "rm should be blocked: %s", out)

	// The file still exists (rm was blocked), clean up
	s.exec(c, []string{"rm", "-f", "/watch/data.txt"})
}

// ---------------------------------------------------------------
// Test: binary blacklist — only /usr/bin/rm is blocked
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestGuard_Blacklist_Binary() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	s.exec(c, []string{"sh", "-c", "echo 'blacklist test' > /watch/data.txt"})

	// Blacklist /usr/bin/rm: only rm (comm=rm) is blocked
	s.startGuardStd(c, "/watch", "-b", "/usr/bin/rm")

	// 1. cat is NOT blacklisted → allowed
	code, out := s.exec(c, []string{"sh", "-c", "cat /watch/data.txt > /dev/null 2>&1"})
	s.Require().Equalf(0, code, "cat should be allowed: %s", out)

	// 2. rm IS blacklisted → blocked
	code, out = s.exec(c, []string{"rm", "/watch/data.txt"})
	s.Require().NotEqualf(0, code, "rm should be blocked: %s", out)

	// Clean up with a shell built-in (not rm) since file still exists
	s.exec(c, []string{"sh", "-c", "> /watch/data.txt && rm -f /watch/data.txt"})

	s.stopGuard(c)
}

// ---------------------------------------------------------------
// Test: whitelist on a single file
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestGuard_Whitelist_Binary_File() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	s.exec(c, []string{"sh", "-c", "echo 'file whitelist' > /watch/target.txt"})

	s.startGuardStd(c, "/watch/target.txt", "-w", "/usr/bin/cat")

	// 1. cat is whitelisted → allowed
	code, out := s.exec(c, []string{"sh", "-c", "cat /watch/target.txt > /dev/null 2>&1"})
	s.Require().Equalf(0, code, "cat should be allowed: %s", out)

	// 2. rm is NOT whitelisted → blocked
	code, out = s.exec(c, []string{"rm", "/watch/target.txt"})
	s.Require().NotEqualf(0, code, "rm should be blocked: %s", out)

	s.exec(c, []string{"rm", "-f", "/watch/target.txt"})
	s.stopGuard(c)
}

// ---------------------------------------------------------------
// Test: blacklist on a single file
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestGuard_Blacklist_Binary_File() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	s.exec(c, []string{"sh", "-c", "echo 'file blacklist' > /watch/target.txt"})

	s.startGuardStd(c, "/watch/target.txt", "-b", "/usr/bin/rm")

	// 1. cat is NOT blacklisted → allowed
	code, out := s.exec(c, []string{"sh", "-c", "cat /watch/target.txt > /dev/null 2>&1"})
	s.Require().Equalf(0, code, "cat should be allowed: %s", out)

	// 2. rm IS blacklisted → blocked
	code, out = s.exec(c, []string{"rm", "/watch/target.txt"})
	s.Require().NotEqualf(0, code, "rm should be blocked: %s", out)

	s.exec(c, []string{"sh", "-c", "rm -f /watch/target.txt"})
	s.stopGuard(c)
}

// ---------------------------------------------------------------
// Test: guard blocks NEW files created inside a guarded directory
// (This was a critical bypass — file_open only checked the file's
// own inode, which doesn't exist for newly created files.)
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestGuard_BlocksNewFiles() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	// Pre-create a subdirectory so its inode is in the guard map
	s.exec(c, []string{"mkdir", "-p", "/watch/subdir"})

	// Guard blocks all (no whitelist)
	s.startGuardStd(c, "/watch")

	// Create a NEW file in the root guarded directory
	code, out := s.exec(c, []string{"sh", "-c", "echo 'new file' > /watch/new_file.txt 2>&1"})
	s.Require().NotEqualf(0, code, "write to new file should be blocked: %s", out)

	// Create a NEW file in a subdirectory
	code, out = s.exec(c, []string{"sh", "-c", "echo 'new sub file' > /watch/subdir/new_sub_file.txt 2>&1"})
	s.Require().NotEqualf(0, code, "write to new subdir file should be blocked: %s", out)

	// Create a new directory
	code, out = s.exec(c, []string{"mkdir", "/watch/newdir"})
	s.Require().NotEqualf(0, code, "new dir creation should be blocked: %s", out)

	// Cleanup stale empty files (VFS creates inode before file_open blocks)
	s.exec(c, []string{"rm", "-f", "/watch/new_file.txt", "/watch/subdir/new_sub_file.txt", "/watch/newdir"})

	s.stopGuard(c)
}

// ---------------------------------------------------------------
// Test: blacklist + event filter
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestGuard_Blacklist_Events() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	s.exec(c, []string{"sh", "-c", "echo 'events test' > /watch/data.txt"})

	// Blacklist rm, filter events to OPEN,WRITE,DELETE (exclude READ)
	// rm should be blocked on DELETE; cat opening the file should generate OPEN/READ but READ is excluded
	s.startGuardStd(c, "/watch", "-b", "/usr/bin/rm", "-e", "OPEN,WRITE,DELETE")
	logBefore := s.readGuardLog(c)

	// 1. cat opens the file (OPEN allowed, READ filtered)
	code, out := s.exec(c, []string{"sh", "-c", "cat /watch/data.txt > /dev/null 2>&1"})
	s.Require().Equalf(0, code, "cat should be allowed: %s", out)

	// 2. rm should be blocked (DELETE)
	code, out = s.exec(c, []string{"rm", "/watch/data.txt"})
	s.Require().NotEqualf(0, code, "rm should be blocked: %s", out)

	time.Sleep(1500 * time.Millisecond)
	logAfter := s.readGuardLog(c)
	deltaEvents := guardDeltaEvents(logBefore, logAfter)
	s.T().Logf("guard delta events: %+v", deltaEvents)

	// DELETE events from rm should be blocked
	s.requireBlockedEvent(deltaEvents, "DELETE", "rm")

	// But wait — cat was allowed, so OPEN events for cat should NOT be blocked
	for _, e := range deltaEvents {
		if e.Type == "OPEN" {
			s.Require().Falsef(e.Blocked, "OPEN from %s should not be blocked", e.Comm)
		}
	}

	s.exec(c, []string{"sh", "-c", "rm -f /watch/data.txt"})
	s.stopGuard(c)
}

// ---------------------------------------------------------------
// Test: whitelist + event filter
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestGuard_Whitelist_Events() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	s.exec(c, []string{"sh", "-c", "echo 'events whitelist' > /watch/data.txt"})

	// Whitelist cat, filter to OPEN,WRITE,DELETE (exclude READ)
	// cat is whitelisted → allowed; rm is not → blocked
	s.startGuardStd(c, "/watch", "-w", "/usr/bin/cat", "-e", "OPEN,WRITE,DELETE")
	logBefore := s.readGuardLog(c)

	// 1. cat is whitelisted → allowed
	code, out := s.exec(c, []string{"sh", "-c", "cat /watch/data.txt > /dev/null 2>&1"})
	s.Require().Equalf(0, code, "cat should be allowed: %s", out)

	// 2. rm is NOT whitelisted → blocked
	code, out = s.exec(c, []string{"rm", "/watch/data.txt"})
	s.Require().NotEqualf(0, code, "rm should be blocked: %s", out)

	logAfter := s.readGuardLog(c)
	deltaEvents := guardDeltaEvents(logBefore, logAfter)

	// rm's DELETE should be blocked
	s.requireBlockedEvent(deltaEvents, "DELETE", "rm")

	// Verify OPEN from cat is NOT blocked
	for _, e := range deltaEvents {
		if e.Type == "OPEN" && e.Comm == "cat" {
			s.Require().Falsef(e.Blocked, "OPEN from whitelisted cat should not be blocked")
		}
	}

	// Can't rm since it's blocked, just stop
	s.stopGuard(c)
}

// ---------------------------------------------------------------
// Test: guard blocks all exploits
//
// The guard uses LSM hooks: file_open → OPEN, mmap_file → MMAP,
// path_unlink → DELETE, path_rename → RENAME, path_link → HARDLINK,
// path_mkdir → MKDIR. There is NO READ or WRITE LSM hook, so
// exploits that only open+read will only be blocked at OPEN.
// The mmap exploit is blocked at OPEN so MMAP never fires.
// ---------------------------------------------------------------

type guardExploitTest struct {
	name   string
	binary string
	events []string
}

var guardExploitTests = []guardExploitTest{
	{name: "copy_file_range", binary: "copy_file_range", events: []string{"OPEN"}},
	{name: "execve", binary: "execve", events: []string{"OPEN"}},
	{name: "io_uring", binary: "io_uring", events: []string{"OPEN"}},
	{name: "mmap", binary: "mmap", events: []string{"OPEN"}},
	{name: "pread64", binary: "pread64", events: []string{"OPEN"}},
	{name: "readv", binary: "readv", events: []string{"OPEN"}},
	{name: "sendfile", binary: "sendfile", events: []string{"OPEN"}},
}

func (s *IntegrationSuite) TestGuard_Exploits() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	s.exec(c, []string{"mkdir", "-p", "/exploits"})

	testFilePath := "/watch/exploit_target.txt"
	s.exec(c, []string{"sh", "-c", fmt.Sprintf("echo 'exploit target' > %s", testFilePath)})

	s.startGuardStd(c, "/watch")
	logCursor := s.readGuardLog(c)

	for _, et := range guardExploitTests {
		if et.name == "io_uring" && !s.checkKernelSupport(c) {
			s.T().Logf("skipping exploit %s (kernel may not support it)", et.name)
			continue
		}

		s.T().Run(et.name, func(t *testing.T) {
			t.Helper()

			exploitHostPath := absPath(fmt.Sprintf("./exploits/%s", et.binary))
			err := c.CopyFileToContainer(s.ctx, exploitHostPath, fmt.Sprintf("/exploits/%s", et.binary), 0755)
			s.Require().NoError(err, "copy exploit binary")

			targetPath := testFilePath
			if et.name == "execve" {
				s.exec(c, []string{"sh", "-c", "echo '#!/bin/sh\necho executed' > /watch/.exec_target && chmod +x /watch/.exec_target"})
				targetPath = "/watch/.exec_target"
			}

			args := []string{fmt.Sprintf("/exploits/%s", et.binary)}
			args = append(args, targetPath)

			code, out := s.exec(c, args)
			s.Require().NotEqualf(0, code,
				"exploit %s should be blocked by guard, got exit %d: %s", et.name, code, out)

			logAfter := s.readGuardLog(c)
			deltaEvents := guardDeltaEvents(logCursor, logAfter)
			logCursor = logAfter

			for _, expectedType := range et.events {
				s.requireBlockedEvent(deltaEvents, expectedType)
			}
		})
	}

	s.stopGuard(c)
}

func (s *IntegrationSuite) checkKernelSupport(c testcontainers.Container) bool {
	code, _ := s.exec(c, []string{"sh", "-c",
		"uname -r | grep -Eq '^4\\.(1[4-9]|[2-9][0-9])|^5\\.[0-9]|^6\\.[0-9]' && echo supported || echo unsupported"})
	return code == 0
}
