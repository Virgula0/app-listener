package integrationtests

import (
	"fmt"
	"slices"
	"strings"

	"github.com/testcontainers/testcontainers-go"
)

// parityOps are recordops' operations (exploits/recordops.c) and the event the guard judges each
// by. internal/recorder/coverage_test.go maps every guard LSM hook to the ops that exercise it.
var parityOps = []struct{ op, key string }{
	{"open", "OPEN"}, {"read", "READ"}, {"readv", "READ"}, {"pread", "READ"}, {"readdir", "READ"},
	{"mmap_read", "MMAP"}, {"mmap_write", "MMAP"},
	{"write", "WRITE"}, {"writev", "WRITE"}, {"pwritev", "WRITE"},
	{"ftruncate", "WRITE"}, {"otrunc", "WRITE"}, {"truncate", "ATTR"},
	{"create", "MKNOD"}, {"mknod_reg", "MKNOD"}, {"mkfifo", "MKNOD"}, {"mkdir", "MKDIR"},
	{"rmdir", "DELETE"}, {"unlink", "DELETE"},
	{"rename", "RENAME"}, {"rename_out", "RENAME"}, {"rename_in", "RENAME"},
	{"link", "HARDLINK"}, {"link_out", "HARDLINK"}, {"link_in", "HARDLINK"}, {"symlink", "SYMLINK"},
	{"chmod", "ATTR"}, {"chown", "ATTR"}, {"utimes", "ATTR"}, {"setxattr", "ATTR"}, {"removexattr", "ATTR"},
	{"getxattr", "STAT"}, {"listxattr", "STAT"}, {"stat", "STAT"}, {"readlink", "STAT"},
	{"splice_in", "READ"}, {"splice_out", "WRITE"}, {"sendfile_in", "READ"}, {"sendfile_out", "WRITE"},
	{"cfr_in", "READ"}, {"cfr_out", "WRITE"},
}

// parityReadOnly ops never write under the target: a WRITE in their mask is over-proposed. Earlier
// ops unlink files there, so a new file outside may reuse a watched inode number.
var parityReadOnly = []string{"open", "read", "readv", "pread", "readdir", "mmap_read", "stat",
	"getxattr", "listxattr", "readlink", "splice_in", "sendfile_in", "cfr_in"}

const (
	parityRec = "/rec"    // the recorded, then guarded, tree: one subdirectory per op
	parityOut = "/recout" // unguarded (cross-tree rename/link/transfer ends)
)

// The recorder and the guard agree on every operation: each op, run by its own copy of recordops
// (its own inode, so its mask is that op's alone), is recorded, and the daemon.conf --record would
// write (Classify → Plan → Render, masks on) must then admit the same op. A negative control runs
// another copy with the recorded mask minus the op's key event, which the guard must deny by
// exactly that event: it pins the key table to the guard's real semantics, so a guard change that
// judges an op by a new event fails the positive half until the monitor reports it too.
func (s *IntegrationSuite) TestRecord_GuardParity() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/recordops"), "/recordops", 0o755))
	s.Require().NoError(c.CopyFileToContainer(s.ctx, recorderTestAmd64Bin, "/recorder.test", 0o755))
	var copies, script []string
	for _, o := range parityOps {
		copies = append(copies, "cp /recordops /ops/"+o.op, "cp /recordops /neg/"+o.op)
		script = append(script, fmt.Sprintf("/ops/%s %s %s/%s %s/%s; echo \"RC %s $?\"", o.op, o.op, parityRec, o.op,
			parityOut, o.op, o.op))
	}
	code, out := s.exec(c, []string{"sh", "-c", "set -e; mkdir -p /ops /neg /etc/app-listener; " +
		strings.Join(copies, "; ")})
	s.Require().Equalf(0, code, "copies: %s", out)
	s.paritySetup(c)

	code, out = s.exec(c, []string{"env", "APPLISTENER_HARNESS_PATH=" + parityRec,
		"APPLISTENER_HARNESS_CMD=" + strings.Join(script, "\n"),
		"/recorder.test", "-test.run", "TestRecorderHarness/TestRecord", "-test.v"})
	s.Require().Equalf(0, code, "recording: %s", out)
	rcs, masks, config := parseParityRecording(out)

	var ran []string
	for _, o := range parityOps {
		switch rcs[o.op] {
		case "0":
		case "3":
			s.T().Logf("%s: unsupported on this filesystem, skipped", o.op)
			continue
		default:
			s.Require().Failf("op failed while recording", "%s rc=%q:\n%s", o.op, rcs[o.op], out)
		}
		mask, ok := masks["/ops/"+o.op]
		s.Require().Truef(ok, "%s was not recorded at all:\n%s", o.op, out)
		s.Require().Containsf(mask, o.key, "%s: the recording missed %s (mask %v)", o.op, o.key, mask)
		if slices.Contains(parityReadOnly, o.op) {
			s.Assert().NotContainsf(mask, "WRITE", "%s: WRITE over-proposed (mask %v)", o.op, mask)
		}
		ran = append(ran, o.op)
	}
	s.Require().NotEmpty(ran)

	s.paritySetup(c)
	s.startDaemon(c, config)
	for _, op := range ran {
		code, out = s.exec(c, []string{"/ops/" + op, op, parityRec + "/" + op, parityOut + "/" + op})
		s.Assert().Equalf(0, code, "%s under its recorded mask %v: %s", op, masks["/ops/"+op], out)
	}
	for _, ev := range parseDaemonEvents(s.readDaemonLog(c)) {
		s.Assert().Falsef(ev.Denied && slices.Contains(ran, ev.Comm),
			"recorded mask denied %s op=%s path=%s", ev.Comm, ev.Op, ev.Path)
	}
	s.stopParityDaemon(c)

	s.paritySetup(c)
	var neg strings.Builder
	fmt.Fprintf(&neg, "[watch %s]\nneed_encryption: false\n", parityRec)
	for _, op := range ran {
		fmt.Fprintf(&neg, "/neg/%s %s\n", op, strings.Join(withoutKey(masks["/ops/"+op], parityKey(op)), ","))
	}
	s.startDaemon(c, neg.String())
	for _, op := range ran {
		code, out = s.exec(c, []string{"/neg/" + op, op, parityRec + "/" + op, parityOut + "/" + op})
		s.Assert().NotEqualf(0, code, "%s without %s must be refused: %s", op, parityKey(op), out)
	}
	events := parseDaemonEvents(s.readDaemonLog(c))
	for _, op := range ran {
		s.Assert().Truef(slices.ContainsFunc(events, func(ev daemonEvent) bool {
			return ev.Denied && ev.Comm == op && ev.Op == parityKey(op)
		}), "%s: no DAEMON DENIED op=%s:\n%s", op, parityKey(op), tailLast(s.readDaemonLog(c), 60))
	}
	s.stopParityDaemon(c)
}

// paritySetup rebuilds every op's fixture (by /recordops, an inode no config names).
func (s *IntegrationSuite) paritySetup(c testcontainers.Container) {
	cmd := []string{"set -e", "rm -rf " + parityRec + " " + parityOut, "mkdir -p " + parityRec + " " + parityOut}
	for _, o := range parityOps {
		cmd = append(cmd, fmt.Sprintf("/recordops setup %s/%s %s/%s", parityRec, o.op, parityOut, o.op))
	}
	code, out := s.exec(c, []string{"sh", "-c", strings.Join(cmd, "; ")})
	s.Require().Equalf(0, code, "fixture: %s", out)
}

func (s *IntegrationSuite) stopParityDaemon(c testcontainers.Container) {
	s.exec(c, []string{"sh", "-c", "pkill -TERM -f 'app-listener daemon' || true"})
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not stop")
}

// parseParityRecording reads the harness output: "RC <op> <code>" from the recorded script,
// "MASK <line> <events>" per admissible binary, and the rendered daemon.conf.
func parseParityRecording(out string) (rcs map[string]string, masks map[string][]string, config string) {
	rcs, masks = map[string]string{}, map[string][]string{}
	var conf []string
	inConfig := false
	for _, line := range strings.Split(out, "\n") {
		switch f := strings.Fields(line); {
		case line == "CONFIG-BEGIN":
			inConfig = true
		case line == "CONFIG-END":
			inConfig = false
		case inConfig:
			conf = append(conf, line)
		case len(f) == 3 && f[0] == "RC":
			rcs[f[1]] = f[2]
		case len(f) == 3 && f[0] == "MASK":
			masks[f[1]] = strings.Split(f[2], ",")
		}
	}
	return rcs, masks, strings.Join(conf, "\n")
}

func parityKey(op string) string {
	for _, o := range parityOps {
		if o.op == op {
			return o.key
		}
	}
	return ""
}

// impliedBy mirrors guard.eventMask: READ, WRITE and MMAP each grant OPEN and STAT too.
var impliedBy = map[string][]string{"OPEN": {"READ", "WRITE", "MMAP"}, "STAT": {"READ", "WRITE", "MMAP"}}

// withoutKey is mask minus key and every event that implies it. An empty event list would mean
// "every event", so a mask left empty keeps an unrelated event instead.
func withoutKey(mask []string, key string) []string {
	drop := append([]string{key}, impliedBy[key]...)
	out := slices.DeleteFunc(slices.Clone(mask), func(e string) bool { return slices.Contains(drop, e) })
	if len(out) == 0 {
		out = []string{"HARDLINK"}
	}
	return out
}
