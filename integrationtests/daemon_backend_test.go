package integrationtests

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"

	"github.com/Virgula0/app-listener/internal/editbackend"
)

// backendScript is requests as the JSON lines the GUI sends, ids 1..n.
func (s *IntegrationSuite) backendScript(reqs ...editbackend.Request) string {
	var b strings.Builder
	for i := range reqs {
		reqs[i].ID = uint64(i + 1)
		line, err := json.Marshal(&reqs[i])
		s.Require().NoError(err)
		b.Write(line)
		b.WriteByte('\n')
	}
	return b.String()
}

// runBackend pipes script into `edit-protected --backend` and returns its replies by id.
func (s *IntegrationSuite) runBackend(c testcontainers.Container, script string) map[uint64]editbackend.Message {
	s.exec(c, []string{"sh", "-c", "printf '%s' " + shQuote(script) + " > /tmp/backend.in"})
	code, out := s.exec(c, []string{"sh", "-c",
		"timeout 60 /app-listener edit-protected --backend < /tmp/backend.in 2>/tmp/backend.err; echo EXIT=$?"})
	s.Require().Equalf(0, code, "backend: %s", out)
	replies := map[uint64]editbackend.Message{}
	for l := range strings.SplitSeq(out, "\n") {
		var m editbackend.Message
		if strings.HasPrefix(l, "{") && json.Unmarshal([]byte(l), &m) == nil && m.Kind == editbackend.KindReply {
			replies[m.ID] = m
		}
	}
	return replies
}

func (s *IntegrationSuite) backendErr(c testcontainers.Container) string {
	_, out := s.exec(c, []string{"sh", "-c", "cat /tmp/backend.err 2>/dev/null"})
	return out
}

// The GUI's edit-protected backend, driven over stdin as the window drives it: nothing before AUTH,
// a wrong password refused, writes confined to the selected resource, the grant revoked on end and
// on EOF (the window closing), no residual access afterwards.
func (s *IntegrationSuite) TestDaemon_EditProtected_Backend() {
	c := s.forwardSetup()
	defer c.Terminate(s.ctx)

	r := s.runBackend(c, s.backendScript(
		editbackend.Request{Op: editbackend.OpHello},
		editbackend.Request{Op: editbackend.OpRead, Path: "secret"},
		editbackend.Request{Op: editbackend.OpAuth, Password: "wrong-Pass-1234"},
		editbackend.Request{Op: editbackend.OpAuth, Password: forwardPassword, TimeoutSec: 30},
		editbackend.Request{Op: editbackend.OpSelect, Resource: "/etc"},
		editbackend.Request{Op: editbackend.OpSelect, Resource: "/protected"},
		editbackend.Request{Op: editbackend.OpList, Path: "."},
		editbackend.Request{Op: editbackend.OpWrite, Path: "secret", Data: []byte("EDITED\n")},
		editbackend.Request{Op: editbackend.OpRead, Path: "../etc/shadow"},
		editbackend.Request{Op: editbackend.OpWrite, Path: "/etc/passwd", Data: []byte("x")},
		editbackend.Request{Op: editbackend.OpRead, Path: "secret"},
		editbackend.Request{Op: editbackend.OpEnd},
	))
	errs := s.backendErr(c)
	s.Require().Lenf(r, 12, "every request must be answered: %v\nstderr:\n%s", r, errs)
	s.Require().Truef(r[1].Hello != nil && r[1].Hello.Live, "hello: %+v", r[1])
	s.Require().Containsf(r[2].Err, "not allowed", "read before auth: %+v", r[2])
	s.Require().Containsf(r[3].Err, "authentication failed", "wrong password: %+v", r[3])
	s.Require().Emptyf(r[4].Err, "auth: %+v\n%s", r[4], errs)
	s.Require().Equal([]string{"/protected"}, r[4].Resources)
	s.Require().Containsf(r[5].Err, "not a guarded directory", "select outside the list: %+v", r[5])
	s.Require().Emptyf(r[6].Err, "select: %+v\n%s", r[6], errs)
	s.Require().Emptyf(r[7].Err, "list: %+v", r[7])
	s.Require().Emptyf(r[8].Err, "write: %+v\n%s", r[8], s.readDaemonLog(c))
	s.Require().NotEmptyf(r[9].Err, "a read above the resource must be refused: %+v", r[9])
	s.Require().NotEmptyf(r[10].Err, "an absolute path must be refused: %+v", r[10])
	s.Require().NotNilf(r[11].File, "read back: %+v", r[11])
	s.Require().Equalf("EDITED\n", string(r[11].File.Data), "read back: %+v", r[11])
	s.Require().Emptyf(r[12].Err, "end: %+v", r[12])

	_, content := s.exec(c, []string{"sh", "-c", listedRead})
	s.Require().Equal("EDITED", strings.TrimSpace(content))
	s.Require().Truef(s.awaitLogCount(c, "write access REVOKED", 1, 15*time.Second),
		"the grant was not revoked at end:\n%s", s.readDaemonLog(c))
	code, out := s.exec(c, []string{"sh", "-c", "echo pwned > /protected/pwned 2>&1"})
	s.Require().NotEqualf(0, code, "a write after the session must be denied: %s", out)

	// EOF without end (the window died): the backend still narrows the grant.
	r = s.runBackend(c, s.backendScript(
		editbackend.Request{Op: editbackend.OpHello},
		editbackend.Request{Op: editbackend.OpAuth, Password: forwardPassword},
		editbackend.Request{Op: editbackend.OpSelect, Resource: "/protected"},
	))
	s.Require().Emptyf(r[3].Err, "select: %+v\n%s", r[3], s.backendErr(c))
	s.Require().Truef(s.awaitLogCount(c, "write access REVOKED", 2, 15*time.Second),
		"EOF did not revoke the grant:\n%s", s.readDaemonLog(c))

	code, out = s.exec(c, []string{"sh", "-c",
		"setpriv --reuid=65534 --regid=65534 --clear-groups /app-listener edit-protected --backend </dev/null 2>&1"})
	s.Require().NotEqualf(0, code, "non-root must be refused: %s", out)
	s.Require().Contains(out, "must be run as root")
	code, out = s.exec(c, []string{"sh", "-c",
		"/app-listener edit-protected --backend --resource /protected </dev/null 2>&1"})
	s.Require().NotEqualf(0, code, "--backend with other flags must be refused: %s", out)
	s.Require().Contains(out, "--backend takes no other flag")

	s.exec(c, []string{"sh", "-c", "pkill -TERM -f 'app-listener daemon' || true"})
}

// A FORWARD grant through the backend lasts exactly as long as its stdin: the window closing
// (EOF) revokes it.
func (s *IntegrationSuite) TestDaemon_EditProtected_Backend_ForwardEndsOnEOF() {
	c := s.forwardSetup()
	defer c.Terminate(s.ctx)

	s.requireRead(c, unlistedRead, false, "baseline: sed is not whitelisted")
	script := s.backendScript(
		editbackend.Request{Op: editbackend.OpHello},
		editbackend.Request{Op: editbackend.OpAuth, Password: forwardPassword, TimeoutSec: 60},
		editbackend.Request{Op: editbackend.OpForward, Resources: []string{"/protected"},
			Binaries: []string{fwdUnlisted}, Allow: true},
	)
	// A holder keeps the FIFO's write end open after the script is written: no EOF until it dies.
	code, out := s.exec(c, []string{"sh", "-c", "mkfifo /tmp/bk.in && " +
		"(nohup /app-listener edit-protected --backend < /tmp/bk.in > /tmp/bk.out 2>/tmp/bk.err &) && " +
		"(nohup sleep 300 > /tmp/bk.in & echo $! > /tmp/bk.holder) && " +
		"printf '%s' " + shQuote(script) + " > /tmp/bk.in"})
	s.Require().Equalf(0, code, "starting the backend: %s", out)
	s.Require().Truef(s.awaitLogCount(c, "temporary access GRANTED", 1, 15*time.Second),
		"grant never became active:\n%s", s.readDaemonLog(c))
	s.requireRead(c, unlistedRead, true, "sed while the backend holds the grant")

	s.exec(c, []string{"sh", "-c", "kill $(cat /tmp/bk.holder)"})
	s.Require().Truef(s.awaitLogCount(c, "temporary access REVOKED", 1, 15*time.Second),
		"EOF did not revoke the grant:\n%s", s.readDaemonLog(c))
	s.requireRead(c, unlistedRead, false, "sed after the backend's stdin closed")

	s.exec(c, []string{"sh", "-c", "pkill -TERM -f 'app-listener daemon' || true"})
}

// forwardSetup seeds a pre-ledger edit-auth.hash: the daemon imports it into binaries.db, which it
// seals. A change goes through the daemon (AUTH with the current password, then PASSWORD); there is
// no removal.
func (s *IntegrationSuite) TestDaemon_EditProtected_ChangePasswordLive() {
	c := s.forwardSetup()
	defer c.Terminate(s.ctx)

	code, _ := s.exec(c, []string{"sh", "-c", "test -e /etc/app-listener/edit-auth.hash"})
	s.Require().NotEqualf(0, code, "the legacy edit-auth.hash was not migrated:\n%s", s.readDaemonLog(c))
	for _, f := range []string{"/etc/app-listener/binaries.db", "/etc/app-listener/binaries.db-journal"} {
		code, out := s.exec(c, []string{"sh", "-c", "cat " + f + " >/dev/null 2>&1"})
		s.Require().NotEqualf(0, code, "%s is readable by a non-app-listener process: %s", f, out)
		code, out = s.exec(c, []string{"sh", "-c", "echo x >> " + f + " 2>&1"})
		s.Require().NotEqualf(0, code, "%s is writable by a non-app-listener process: %s", f, out)
	}

	const newPassword = "Brand-New-Secret-77"
	r := s.runBackend(c, s.backendScript(
		editbackend.Request{Op: editbackend.OpHello},
		editbackend.Request{Op: editbackend.OpSetPassword, Password: "wrong-Pass-1234", NewPassword: newPassword},
		editbackend.Request{Op: editbackend.OpSetPassword, Password: forwardPassword, NewPassword: "weak"},
		editbackend.Request{Op: editbackend.OpSetPassword, Password: forwardPassword, NewPassword: newPassword},
		editbackend.Request{Op: editbackend.OpEnd},
	))
	errs := s.backendErr(c)
	s.Require().NotEmptyf(r[2].Err, "a change with a wrong current password was accepted: %+v", r[2])
	s.Require().NotEmptyf(r[3].Err, "a weak new password was accepted: %+v", r[3])
	s.Require().Emptyf(r[4].Err, "change: %+v\n%s\n%s", r[4], errs, s.readDaemonLog(c))
	s.Require().Truef(s.awaitLog(c, "edit-protected password changed", 10*time.Second),
		"the daemon did not store the change:\n%s", s.readDaemonLog(c))

	r = s.runBackend(c, s.backendScript(
		editbackend.Request{Op: editbackend.OpHello},
		editbackend.Request{Op: editbackend.OpAuth, Password: forwardPassword},
	))
	s.Require().Containsf(r[2].Err, "authentication failed", "the old password still works: %+v", r[2])
	r = s.runBackend(c, s.backendScript(
		editbackend.Request{Op: editbackend.OpHello},
		editbackend.Request{Op: editbackend.OpAuth, Password: newPassword},
		editbackend.Request{Op: editbackend.OpEnd},
	))
	s.Require().Emptyf(r[2].Err, "the new password is refused: %+v\n%s", r[2], s.backendErr(c))

	code, out := s.exec(c, []string{"sh", "-c", "/app-listener edit-protected --clear-password </dev/null 2>&1"})
	s.Require().NotEqualf(0, code, "--clear-password must not exist: %s", out)

	s.exec(c, []string{"sh", "-c", "pkill -TERM -f 'app-listener daemon' || true"})
}
