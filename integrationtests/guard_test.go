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
	{name: "splice", binary: "splice", events: []string{"OPEN"}},
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

// ---------------------------------------------------------------
// Test: non-recursive guard — deeply nested files are NOT protected
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestGuard_NonRecursive() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	// Set up a 3-level hierarchy:
	//   /watch/
	//   ├── top.txt               (level 1, direct child)
	//   └── subdir/               (level 1, direct child dir)
	//       └── deep.txt          (level 2)
	s.exec(c, []string{"mkdir", "-p", "/watch/subdir"})
	s.exec(c, []string{"sh", "-c", "echo top > /watch/top.txt"})
	s.exec(c, []string{"sh", "-c", "echo deep > /watch/subdir/deep.txt"})

	// Guard WITHOUT recursion — only the root and its direct children
	// are pre-populated in the inode map.
	s.startGuardStd(c, "/watch", "--recursive=false")
	logBefore := s.readGuardLog(c)

	// 1. Direct child file IS blocked (own inode in map)
	code, out := s.exec(c, []string{"sh", "-c", "cat /watch/top.txt > /dev/null 2>&1"})
	s.Require().NotEqualf(0, code, "direct child file should be blocked: %s", out)

	// 2. File in direct child subdir is NOT blocked (non-recursive does not
	//    add subdirectory inodes to the BPF map).
	code, out = s.exec(c, []string{"sh", "-c", "cat /watch/subdir/deep.txt > /dev/null 2>&1"})
	s.Require().Equalf(0, code, "file in subdir should NOT be blocked (non-recursive): %s", out)

	// 3. mkdir inside subdir is NOT blocked (same reason)
	code, out = s.exec(c, []string{"mkdir", "-p", "/watch/subdir/newchild"})
	s.Require().Equalf(0, code, "mkdir in subdir should NOT be blocked: %s", out)
	// Clean up
	s.exec(c, []string{"rmdir", "/watch/subdir/newchild"})

	logAfter := s.readGuardLog(c)
	deltaEvents := guardDeltaEvents(logBefore, logAfter)

	// Only the top.txt OPEN event should be present and blocked; no subdir events
	s.Require().NotEmpty(deltaEvents, "expected at least one blocked event")
	// All events should be blocked (only root-level events fire)
	for _, e := range deltaEvents {
		s.Require().Truef(e.Blocked, "event %s|%s should be blocked", e.Type, e.Comm)
	}

	s.stopGuard(c)
}

// ---------------------------------------------------------------
// Test: depth-limit guard — files beyond depth are NOT blocked
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestGuard_DepthLimit() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	// Set up a 4-level hierarchy:
	//   /watch/                   (level 0)
	//   ├── top.txt               (level 1)
	//   ├── subdir/               (level 1)
	//   │   ├── mid.txt           (level 2)
	//   │   └── inner/            (level 2)
	//   │       └── deep.txt      (level 3)
	s.exec(c, []string{"mkdir", "-p", "/watch/subdir/inner"})
	s.exec(c, []string{"sh", "-c", "echo top > /watch/top.txt"})
	s.exec(c, []string{"sh", "-c", "echo mid > /watch/subdir/mid.txt"})
	s.exec(c, []string{"sh", "-c", "echo deep > /watch/subdir/inner/deep.txt"})

	// Guard with depth=2: level 0,1,2 are guarded; level 3+ is not.
	// inner/ is NOT added to the BPF inode map (boundary dir skipped),
	// so deep.txt at level 3 has no parent inode in the map.
	s.startGuardStd(c, "/watch", "--recursive", "--depth", "2")
	logBefore := s.readGuardLog(c)

	// 1. Direct child file (level 1) — blocked
	code, out := s.exec(c, []string{"sh", "-c", "cat /watch/top.txt > /dev/null 2>&1"})
	s.Require().NotEqualf(0, code, "direct child should be blocked: %s", out)

	// 2. File at depth 2 inside subdir — blocked (own inode in map)
	code, out = s.exec(c, []string{"sh", "-c", "cat /watch/subdir/mid.txt > /dev/null 2>&1"})
	s.Require().NotEqualf(0, code, "depth-2 file should be blocked: %s", out)

	// 3. File at depth 3 inside inner/ — NOT blocked (inner/ at depth boundary
	//    is NOT in the BPF inode map, so the parent check misses).
	code, out = s.exec(c, []string{"sh", "-c", "cat /watch/subdir/inner/deep.txt > /dev/null 2>&1"})
	s.Require().Equalf(0, code, "depth-3 file should NOT be blocked: %s", out)

	logAfter := s.readGuardLog(c)
	deltaEvents := guardDeltaEvents(logBefore, logAfter)
	s.Require().NotEmpty(deltaEvents, "expected blocked events")

	s.stopGuard(c)
}

// ---------------------------------------------------------------
// Test: binary whitelist with recursive guard
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestGuard_Whitelist_Recursive() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch/subdir"})
	s.exec(c, []string{"sh", "-c", "echo 'recursive whitelist' > /watch/target.txt"})
	s.exec(c, []string{"sh", "-c", "echo 'subdir file' > /watch/subdir/target.txt"})

	// Whitelist cat — cat is allowed at any depth with default recursive=true
	s.startGuardStd(c, "/watch", "-w", "/usr/bin/cat")

	// 1. cat at root — allowed
	code, out := s.exec(c, []string{"sh", "-c", "cat /watch/target.txt > /dev/null 2>&1"})
	s.Require().Equalf(0, code, "cat at root should be allowed: %s", out)

	// 2. cat in subdir — allowed (recursive)
	code, out = s.exec(c, []string{"sh", "-c", "cat /watch/subdir/target.txt > /dev/null 2>&1"})
	s.Require().Equalf(0, code, "cat in subdir should be allowed: %s", out)

	// 3. rm at root — NOT whitelisted, blocked
	code, out = s.exec(c, []string{"rm", "/watch/target.txt"})
	s.Require().NotEqualf(0, code, "rm should be blocked: %s", out)

	// Clean up
	s.exec(c, []string{"rm", "-f", "/watch/target.txt", "/watch/subdir/target.txt"})
	s.stopGuard(c)
}

// ---------------------------------------------------------------
// Test: binary blacklist with depth-limited guard
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestGuard_Blacklist_Depth() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	// Hierarchy:
	//   /watch/
	//   ├── target.txt            (level 1)
	//   └── subdir/               (level 1)
	//       ├── nested.txt        (level 2)
	//       └── inner/            (level 2)
	//           └── deep.txt      (level 3)
	s.exec(c, []string{"mkdir", "-p", "/watch/subdir/inner"})
	s.exec(c, []string{"sh", "-c", "echo 'depth blacklist' > /watch/target.txt"})
	s.exec(c, []string{"sh", "-c", "echo nested > /watch/subdir/nested.txt"})
	s.exec(c, []string{"sh", "-c", "echo deep > /watch/subdir/inner/deep.txt"})

	// Blacklist rm with depth=1: only level 1 is guarded; files at depth 2+
	// leak (subdir/ added as dir inode at depth limit, but inner/ is not).
	// Since inner/ is NOT in the inode map, deep.txt's parent check passes
	// and rm succeeds.
	s.startGuardStd(c, "/watch", "-b", "/usr/bin/rm", "--recursive", "--depth", "1")

	// 1. rm at level 1 — blocked (own inode + blacklisted comm)
	code, out := s.exec(c, []string{"rm", "/watch/target.txt"})
	s.Require().NotEqualf(0, code, "rm at root should be blocked: %s", out)

	// 2. rm at level 2 (nested.txt) — NOT blocked (subdir/ at depth boundary
	//    is NOT in the BPF inode map, so the parent check misses).
	code, out = s.exec(c, []string{"rm", "/watch/subdir/nested.txt"})
	s.Require().Equalf(0, code, "rm at depth 2 should NOT be blocked: %s", out)

	// 3. rm at level 3 (deep.txt) — NOT blocked (inner/ also NOT in map)
	code, out = s.exec(c, []string{"rm", "/watch/subdir/inner/deep.txt"})
	s.Require().Equalf(0, code, "rm at depth 3 should NOT be blocked: %s", out)

	// 4. Create a NEW file at depth 3 — allowed (inner/ NOT in inode map)
	code, out = s.exec(c, []string{"sh", "-c", "touch /watch/subdir/inner/new_file.txt 2>&1"})
	s.Require().Equalf(0, code, "creating new file at depth 3 should NOT be blocked: %s", out)

	// 5. Read the new file via cat — allowed (cat not blacklisted, and
	//    inner/ NOT in inode map anyway).
	code, out = s.exec(c, []string{"sh", "-c", "cat /watch/subdir/inner/new_file.txt > /dev/null 2>&1"})
	s.Require().Equalf(0, code, "cat at depth 3 should NOT be blocked: %s", out)

	s.stopGuard(c)
}

// ---------------------------------------------------------------
// Bypass: fork+exec with inherited fd
//
// The POC opens a guarded file WITHOUT O_CLOEXEC, forks, and the
// child exec's a target binary (e.g. /usr/bin/cat) that reads from
// the inherited fd.  file_permission fires with the new binary's
// context and should block it if the binary is blacklisted.
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestGuard_Bypass_ForkExecFD() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch", "/exploits"})
	s.exec(c, []string{"sh", "-c", "echo 'fork+exec bypass target' > /watch/target.txt"})

	exploitHostPath := absPath("./exploits/fork_exec_fd")
	err := c.CopyFileToContainer(s.ctx, exploitHostPath, "/exploits/fork_exec_fd", 0755)
	s.Require().NoError(err, "copy fork_exec_fd binary")

	s.startGuardStd(c, "/watch", "-b", "/usr/bin/cat")

	code, out := s.exec(c, []string{"/exploits/fork_exec_fd", "/watch/target.txt", "/usr/bin/cat"})
	s.Require().NotEqualf(0, code,
		"fork_exec_fd bypass should be blocked (cat is blacklisted): %s", out)

	s.stopGuard(c)
}

// ---------------------------------------------------------------
// Bypass: symlink to a blacklisted binary
//
// An attacker may create a symlink to a blacklisted binary and run
// it under a different name to bypass guard:
//
//   ln -s /usr/bin/cat /tmp/myreader
//   /tmp/myreader /watch/secret.txt
//
// The guard uses exe-inode identity (resolves symlinks), so the
// symlinked binary is recognized as the same executable and blocked.
// ---------------------------------------------------------------
func (s *IntegrationSuite) TestGuard_Bypass_SymlinkBinary() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	s.exec(c, []string{"sh", "-c", "echo 'secret' > /watch/target.txt"})

	// Create a symlink to a blacklisted binary
	s.exec(c, []string{"ln", "-sf", "/usr/bin/cat", "/tmp/myreader"})

	s.startGuardStd(c, "/watch", "-b", "/usr/bin/cat")

	// Should be blocked even when accessed via the symlink
	code, out := s.exec(c, []string{"/tmp/myreader", "/watch/target.txt"})
	s.Require().NotEqualf(0, code,
		"symlink to blacklisted binary should be blocked: %s", out)

	s.stopGuard(c)
}

// ---------------------------------------------------------------
// Bypass: SCM_RIGHTS fd passing + exec
//
// The POC opens a guarded file, passes the fd to a child process
// via SCM_RIGHTS, then the child exec's cat which reads from the
// passed fd.  file_permission should block it.
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestGuard_Bypass_SCMRights() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch", "/exploits"})
	s.exec(c, []string{"sh", "-c", "echo 'SCM_RIGHTS bypass target' > /watch/target.txt"})

	exploitHostPath := absPath("./exploits/scm_rights_pass")
	err := c.CopyFileToContainer(s.ctx, exploitHostPath, "/exploits/scm_rights_pass", 0755)
	s.Require().NoError(err, "copy scm_rights_pass binary")

	s.startGuardStd(c, "/watch", "-b", "/usr/bin/cat")

	code, out := s.exec(c, []string{"/exploits/scm_rights_pass", "/watch/target.txt"})
	s.Require().NotEqualf(0, code,
		"SCM_RIGHTS bypass should be blocked (cat is blacklisted): %s", out)

	s.stopGuard(c)
}

// ---------------------------------------------------------------
// Bypass: open_by_handle_at — open file by inode handle
//
// Uses name_to_handle_at() + open_by_handle_at() to open a file
// without specifying a filesystem path.  The guard's file_open
// hook catches this because vfs_open resolves the dentry from
// the handle.  Requires CONFIG_FHANDLE + CAP_DAC_READ_SEARCH.
// Skipped if the kernel does not support name_to_handle_at.
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestGuard_Bypass_OpenByHandleAt() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch", "/exploits"})
	s.exec(c, []string{"sh", "-c", "echo 'handle bypass target' > /watch/target.txt"})

	exploitHostPath := absPath("./exploits/open_by_handle_at")
	err := c.CopyFileToContainer(s.ctx, exploitHostPath, "/exploits/open_by_handle_at", 0755)
	s.Require().NoError(err, "copy open_by_handle_at binary")

	s.startGuardStd(c, "/watch")
	logBefore := s.readGuardLog(c)

	code, out := s.exec(c, []string{"/exploits/open_by_handle_at", "/watch/target.txt"})

	logAfter := s.readGuardLog(c)
	deltaEvents := guardDeltaEvents(logBefore, logAfter)

	if code != 0 && len(deltaEvents) == 0 {
		s.T().Skipf("open_by_handle_at not supported on this kernel: %s", out)
	}

	s.Require().NotEqualf(0, code,
		"open_by_handle_at should be blocked by guard: %s", out)
	s.requireBlockedEvent(deltaEvents, "OPEN")

	s.stopGuard(c)
}

// ---------------------------------------------------------------
// Bypass: raw block device — read via debugfs on loop device
//
// This test requires loop device support in Docker and is skipped
// in automated CI.  Run manually on a real host to verify:
//
//   sudo ./integrationtests/exploits/raw_block_device /dev/mapper/cryptlvm /home/user/secret.txt
//
// Before the fix: debugfs opens the block device and reads the file
//   without triggering any guard event — VFS bypass.
//
// After the fix: the guard auto-detects the backing block device for
//   each watched path and blocks open() on that device via the
//   guard_fs_devices BPF map.  The above command should now fail
//   with the guard blocking the block device access.
// ---------------------------------------------------------------

func (s *IntegrationSuite) TestGuard_Bypass_RawBlockDevice() {
	s.T().Skip("requires loop device support and a real ext4 filesystem — run manually")
}

// ---------------------------------------------------------------
// Test: runtime mkdir — new directories under recursive guard are
// lazily discovered via file_open, so files inside are blocked.
// Uses blacklist mode: only /usr/bin/cat is blacklisted.  Uses
// gnumkdir (non-hardlinked) because /usr/bin/mkdir is hardlinked to
// cat on this Ubuntu image and would be blocked by exe-inode identity.
// ---------------------------------------------------------------
func (s *IntegrationSuite) TestGuard_RuntimeNewDir() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	s.exec(c, []string{"touch", "/watch/pre_existing.txt"})

	// Guard uses exe-inode identity (resolves symlinks).  On this Ubuntu image
	// /usr/bin/cat and /usr/bin/mkdir are hardlinked (same exe inode), so
	// blacklisting cat also blocks mkdir.  Use gnumkdir (a separate, non-hardlinked
	// binary) to create the directory.
	s.startGuardStd(c, "/watch", "--recursive", "-b", "/usr/bin/cat")
	logBefore := s.readGuardLog(c)

	// Pre-existing file: cat should be blocked
	code, out := s.exec(c, []string{"/usr/bin/cat", "/watch/pre_existing.txt"})
	s.Require().NotEqualf(0, code, "pre-existing file should be blocked for cat: %s", out)

	// Create a new subdirectory AFTER guard started (via gnumkdir, not blacklisted)
	code, out = s.exec(c, []string{"/usr/bin/gnumkdir", "/watch/runtime_dir"})
	s.Require().Equalf(0, code, "dir creation via gnumkdir should succeed: %s", out)

	// Create a file inside (shell not blacklisted → allowed)
	code, out = s.exec(c, []string{"sh", "-c", "echo secret > /watch/runtime_dir/data.txt"})
	s.Require().Equalf(0, code, "shell file creation should succeed: %s", out)

	// Lazy discovery fires during the above open → parent added to guard_inodes.
	// Now read with cat (blacklisted) → should be blocked.
	code, out = s.exec(c, []string{"/usr/bin/cat", "/watch/runtime_dir/data.txt"})
	s.Require().NotEqualf(0, code, "file in runtime-created dir should be blocked for cat: %s", out)

	logAfter := s.readGuardLog(c)
	deltaEvents := guardDeltaEvents(logBefore, logAfter)
	s.requireBlockedEvent(deltaEvents, "OPEN", "cat")

	s.stopGuard(c)
}

// ---------------------------------------------------------------
// Test: runtime mkdir with non-recursive guard — new directories
// are NOT auto-discovered (recursive flag is 0), so files inside
// are NOT blocked even for blacklisted binaries.
// ---------------------------------------------------------------
func (s *IntegrationSuite) TestGuard_RuntimeNewDir_NonRecursive() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	s.exec(c, []string{"touch", "/watch/top.txt"})

	// Guard uses exe-inode identity; cat and mkdir are hardlinked (same exe
	// inode), so blacklisting cat also blocks mkdir.  Use gnumkdir (separate,
	// non-hardlinked binary) to create the directory.
	s.startGuardStd(c, "/watch", "--recursive=false", "-b", "/usr/bin/cat")

	// Pre-existing file: cat blocked
	code, out := s.exec(c, []string{"/usr/bin/cat", "/watch/top.txt"})
	s.Require().NotEqualf(0, code, "pre-existing file should be blocked for cat: %s", out)

	// Create a new subdirectory after guard started (via gnumkdir, not blacklisted)
	code, out = s.exec(c, []string{"/usr/bin/gnumkdir", "/watch/runtime_dir"})
	s.Require().Equalf(0, code, "dir creation via gnumkdir should succeed: %s", out)

	code, out = s.exec(c, []string{"sh", "-c", "echo data > /watch/runtime_dir/data.txt"})
	s.Require().Equalf(0, code, "file creation should succeed: %s", out)

	// Should NOT be blocked (non-recursive → discover_guarded_parent checks
	// should_add_new_dir, which would return 0, so parent is NOT added)
	code, out = s.exec(c, []string{"/usr/bin/cat", "/watch/runtime_dir/data.txt"})
	s.Require().Equalf(0, code, "file in runtime dir should NOT be blocked (non-recursive): %s", out)

	s.stopGuard(c)
}

// ---------------------------------------------------------------
// Test: rename a directory from outside into the guarded area.
// The rename succeeds (moving into guarded area is allowed), but
// the moved directory's inode is added to guard_inodes so files
// inside become guarded.
// ---------------------------------------------------------------
func (s *IntegrationSuite) TestGuard_RuntimeRenameDir() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	s.exec(c, []string{"mkdir", "-p", "/outside"})
	s.exec(c, []string{"mkdir", "/outside/data"})
	s.exec(c, []string{"sh", "-c", "echo secret > /outside/data/file.txt"})

	// Guard watches /watch recursively, blacklist cat
	s.startGuardStd(c, "/watch", "--recursive", "-b", "/usr/bin/cat")
	logBefore := s.readGuardLog(c)

	// Rename the outside directory into the guarded area.
	// The rename itself is allowed (source was not guarded, mv not blacklisted).
	code, out := s.exec(c, []string{"mv", "/outside/data", "/watch/data"})
	s.Require().Equalf(0, code, "rename into guarded area should succeed: %s", out)
	time.Sleep(1 * time.Second)

	// Now that data/ is under /watch, its contents should be guarded.
	// The rename hook adds the moved-in directory's inode to guard_inodes.
	code, out = s.exec(c, []string{"/usr/bin/cat", "/watch/data/file.txt"})
	s.Require().NotEqualf(0, code, "file in moved-in dir should be blocked for cat: %s", out)

	logAfter := s.readGuardLog(c)
	deltaEvents := guardDeltaEvents(logBefore, logAfter)
	s.requireBlockedEvent(deltaEvents, "OPEN", "cat")

	s.stopGuard(c)
}

// ---------------------------------------------------------------
// Test: rename a file (not a directory) into the guarded area.
// Files do not get their own inode added on rename-in, but the
// file is now under a guarded parent so it IS blocked (by parent
// inode check).
// ---------------------------------------------------------------
func (s *IntegrationSuite) TestGuard_RuntimeRenameFile() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch"})
	s.exec(c, []string{"mkdir", "-p", "/outside"})
	s.exec(c, []string{"sh", "-c", "echo data > /outside/file.txt"})

	// Guard watches /watch recursively, blacklist cat
	s.startGuardStd(c, "/watch", "--recursive", "-b", "/usr/bin/cat")
	logBefore := s.readGuardLog(c)

	// Rename a regular file into the guarded area
	code, out := s.exec(c, []string{"mv", "/outside/file.txt", "/watch/file.txt"})
	s.Require().Equalf(0, code, "rename file into guarded area should succeed: %s", out)
	time.Sleep(1 * time.Second)

	// The file is now under /watch guarded parent — accessing should be blocked
	code, out = s.exec(c, []string{"/usr/bin/cat", "/watch/file.txt"})
	s.Require().NotEqualf(0, code, "file moved into guarded area should be blocked for cat: %s", out)

	logAfter := s.readGuardLog(c)
	deltaEvents := guardDeltaEvents(logBefore, logAfter)
	s.requireBlockedEvent(deltaEvents, "OPEN", "cat")

	s.stopGuard(c)
}
