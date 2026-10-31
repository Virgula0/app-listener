package integrationtests

import (
	"fmt"
	"strings"
	"testing"
)

// Bypass: AF_UNIX sockets bound inside a guarded tree (agent / RPC sockets). connect(2) finds the
// listener by inode and opens nothing, so only a connect/send/recv gate denies a non-whitelisted
// caller. The servers bind before the guard starts, as a long-running agent does.
func (s *IntegrationSuite) TestGuard_Bypass_UnixConnect() {
	c := s.guardContainer()
	// pooled: terminated at suite end

	s.exec(c, []string{"mkdir", "-p", "/watch", "/exploits"})
	host := absPath("./exploits/unix_connect")
	for _, name := range []string{"unix_connect", "unix_connect_ok"} {
		s.Require().NoError(c.CopyFileToContainer(s.ctx, host, "/exploits/"+name, 0755), "copy %s", name)
	}
	for _, t := range []string{"stream", "seqpacket", "dgram"} {
		s.capturePID(c, fmt.Sprintf("nohup /exploits/unix_connect serve /watch/%s.sock %s >/tmp/serve-%s.log 2>&1 &", t, t, t))
	}
	defer s.exec(c, []string{"sh", "-c", "pkill -f 'unix_connect serve' || true"})
	code, out := s.exec(c, []string{"sh", "-c",
		"for i in $(seq 50); do [ -S /watch/stream.sock ] && [ -S /watch/seqpacket.sock ] && [ -S /watch/dgram.sock ] && exit 0; sleep 0.1; done; exit 1"})
	s.Require().Equalf(0, code, "servers never bound: %s", out)

	s.startGuardStd(c, "/watch", "-w", "/exploits/unix_connect_ok")
	defer s.stopGuard(c)

	code, out = s.exec(c, []string{"/exploits/unix_connect_ok", "stream", "/watch/stream.sock"})
	s.Require().Equalf(0, code, "whitelisted connect must reach the server: %s", out)

	cases := []struct {
		name string
		args []string
	}{
		{"stream", []string{"/exploits/unix_connect", "stream", "/watch/stream.sock"}},
		{"seqpacket", []string{"/exploits/unix_connect", "seqpacket", "/watch/seqpacket.sock"}},
		{"dgram", []string{"/exploits/unix_connect", "dgram", "/watch/dgram.sock"}},
		{"uring", []string{"/exploits/unix_connect", "uring", "/watch/stream.sock"}},
		{"inherit", []string{"/exploits/unix_connect_ok", "inherit", "/watch/stream.sock", "/exploits/unix_connect"}},
		{"inheritsplice", []string{"/exploits/unix_connect_ok", "inheritsplice", "/watch/stream.sock", "/exploits/unix_connect"}},
		{"passfd", []string{"/exploits/unix_connect_ok", "passfd", "/watch/stream.sock", "/exploits/unix_connect"}},
	}
	cursor := s.readGuardLog(c)
	for _, tc := range cases {
		s.T().Run(tc.name, func(t *testing.T) {
			code, out := s.exec(c, tc.args)
			if code == 2 {
				t.Skipf("unsupported here: %s", out)
			}
			s.Require().Falsef(strings.Contains(out, "UNIX-SOCKET-SECRET"),
				"%s reached the guarded socket: %s", tc.name, out)
			s.Require().Equalf(1, code, "%s: want a kernel denial (exit 1), got %d: %s", tc.name, code, out)

			after := s.readGuardLog(c)
			events := guardDeltaEvents(cursor, after)
			cursor = after
			s.requireBlockedEvent(events, "CONNECT", "unix_connect")
		})
	}
}

// A connection admitted before a reload keeps its standing after it: the whitelisted holder still
// reads, a non-whitelisted exe it hands the fd to is still denied. A reload moves the resource to a
// new slot and retires the old one, so a mark naming the slot would deny the first and, once the
// slot is reused, judge the second by an unrelated whitelist.
func (s *IntegrationSuite) TestDaemon_UnixConnectMarkSurvivesReload() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/watch", "/exploits", "/etc/app-listener"})
	host := absPath("./exploits/unix_connect")
	for _, name := range []string{"unix_connect", "unix_connect_ok"} {
		s.Require().NoError(c.CopyFileToContainer(s.ctx, host, "/exploits/"+name, 0755), "copy %s", name)
	}
	s.capturePID(c, "nohup /exploits/unix_connect serve /watch/stream.sock stream >/tmp/serve.log 2>&1 &")
	code, out := s.exec(c, []string{"sh", "-c",
		"for i in $(seq 50); do [ -S /watch/stream.sock ] && exit 0; sleep 0.1; done; exit 1"})
	s.Require().Equalf(0, code, "server never bound: %s", out)

	s.startDaemon(c, "[watch /watch]\nneed_encryption: false\n/exploits/unix_connect_ok")

	holders := map[string]string{
		"self":    "",
		"inherit": " /exploits/unix_connect",
	}
	for name, exe := range holders {
		s.exec(c, []string{"sh", "-c", fmt.Sprintf(
			"(/exploits/unix_connect_ok hold /watch/stream.sock /tmp/go%s; echo $? > /tmp/hold-%s.rc) >/tmp/hold-%s.log 2>&1 &",
			exe, name, name)})
	}
	code, out = s.exec(c, []string{"sh", "-c",
		"for i in $(seq 100); do grep -q connected /tmp/hold-self.log && grep -q connected /tmp/hold-inherit.log && exit 0; " +
			"sleep 0.1; done; cat /tmp/hold-*.log; exit 1"})
	s.Require().Equalf(0, code, "whitelisted holders never connected: %s", out)

	s.reloadAndAwait(c)
	s.reloadAndAwait(c)
	s.exec(c, []string{"touch", "/tmp/go"})
	code, out = s.exec(c, []string{"sh", "-c",
		"for i in $(seq 100); do [ -f /tmp/hold-self.rc ] && [ -f /tmp/hold-inherit.rc ] && exit 0; sleep 0.1; done; exit 1"})
	s.Require().Equalf(0, code, "holders never finished: %s", out)

	_, self := s.exec(c, []string{"sh", "-c", "cat /tmp/hold-self.rc /tmp/hold-self.log"})
	s.Require().Truef(strings.HasPrefix(self, "0") && strings.Contains(self, "UNIX-SOCKET-SECRET"),
		"the whitelisted holder lost its connection across the reload: %s", self)
	_, inherit := s.exec(c, []string{"sh", "-c", "cat /tmp/hold-inherit.rc /tmp/hold-inherit.log"})
	s.Require().Falsef(strings.Contains(inherit, "UNIX-SOCKET-SECRET"),
		"a non-whitelisted exe read a pre-reload connection: %s", inherit)
	s.Require().Truef(strings.HasPrefix(inherit, "1"), "want a kernel denial (exit 1): %s", inherit)

	denied := false
	for _, e := range parseDaemonEvents(s.readDaemonLog(c)) {
		denied = denied || (e.Denied && e.Op == "CONNECT" && e.Comm == "unix_connect")
	}
	s.Require().Truef(denied, "no DAEMON DENIED op=CONNECT for the handed-off read:\n%s", s.readDaemonLog(c))

	s.exec(c, []string{"sh", "-c", "pkill -f 'unix_connect' 2>/dev/null; pkill -f 'app-listener daemon' || true"})
}

// A runtime section (/run/user/<uid>/ssh-agent) exists only while its user is logged in: `daemon
// login` makes it root-owned, guards it, then hands it to the user; `daemon logout` drops it once
// the mount is gone, and the next login guards the fresh mount's directory.
func (s *IntegrationSuite) TestDaemon_RuntimeSectionLogin() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	code, out := s.exec(c, []string{"sh", "-c", "command -v ssh-agent >/dev/null 2>&1 || " +
		"(timeout 120 apt-get update -qq && timeout 120 apt-get install -y -qq --no-install-recommends openssh-client >/dev/null)"})
	if code != 0 {
		s.T().Skipf("openssh-client unavailable: %s", out)
	}
	const dir = "/run/user/1000/ssh-agent"
	const sock = dir + "/socket"
	mount := "mkdir -p /run/user/1000 && mount -t tmpfs -o mode=0700,uid=1000,gid=1000 tmpfs /run/user/1000"
	code, out = s.exec(c, []string{"sh", "-c", "(id -u 1000 >/dev/null 2>&1 || useradd -u 1000 -m u1000) && " +
		"mkdir -p /watch /etc/app-listener && cp /usr/bin/ssh-add /tmp/ssh-add && " + mount})
	s.Require().Equalf(0, code, "setup: %s", out)
	asUser := func(cmd string) (int, string) {
		return s.exec(c, []string{"setpriv", "--reuid=1000", "--regid=1000", "--clear-groups", "sh", "-c", cmd})
	}

	const always = "[watch /watch]\nneed_encryption: false\n/usr/bin/true\n"
	s.exec(c, []string{"sh", "-c", "cat > /etc/app-listener/daemon.conf <<'EOF'\n" + always +
		"\n[watch " + dir + "]\nruntime: true\n/usr/bin/ssh-agent MKNOD,DELETE,STAT\n/usr/bin/rm DELETE,STAT\n" +
		"/usr/bin/ssh-add CONNECT,STAT\n/usr/bin/stat STAT\nneed_encryption: false\nEOF"})
	s.launchDaemon(c)
	s.awaitDaemonUp(c, always) // the runtime section guards nothing before the login
	// Probes inside the guarded root go through the whitelisted stat: sh's `test` is denied STAT.
	isSock := "[ \"$(stat -c %F " + sock + ")\" = socket ]"

	login := func() {
		code, out := s.exec(c, []string{"/app-listener", "daemon", "login", "1000"})
		s.Require().Equalf(0, code, "daemon login: %s\n%s", out, s.readDaemonLog(c))
		code, out = s.exec(c, []string{"stat", "-c", "%u %a", dir, dir + "/.app-listener/README"})
		s.Require().Equalf(0, code, "stat: %s", out)
		s.Require().Equalf("1000 700\n0 444", out, "the dir goes to the user, the sentinel stays root's")
	}
	login()

	code, out = asUser("(nohup /usr/bin/ssh-agent -D -a " + sock + " >/tmp/agent.log 2>&1 &) && " +
		"for i in $(seq 50); do " + isSock + " 2>/dev/null && exit 0; sleep 0.1; done; cat /tmp/agent.log; exit 1")
	s.Require().Equalf(0, code, "ssh-agent could not bind in the runtime dir: %s", out)
	code, out = asUser("SSH_AUTH_SOCK=" + sock + " /usr/bin/ssh-add -l 2>&1")
	s.Require().Equalf(1, code, "whitelisted ssh-add must reach the agent (exit 1 = no identities): %s", out)
	code, out = asUser("SSH_AUTH_SOCK=" + sock + " timeout 10 /tmp/ssh-add -l 2>&1")
	s.Require().Equalf(2, code, "a non-whitelisted ssh-add reached the agent (exit 2 = can't connect): %s", out)
	s.Require().Containsf(out, "not permitted", "want the guard's EPERM: %s", out)
	denied := false
	for _, e := range parseDaemonEvents(s.readDaemonLog(c)) {
		denied = denied || (e.Denied && e.Op == "CONNECT" && e.Comm == "ssh-add")
	}
	s.Require().Truef(denied, "no DAEMON DENIED op=CONNECT for the copied ssh-add:\n%s", s.readDaemonLog(c))

	// rm may clear a killed agent's socket, never the root: the sentinel keeps it non-empty.
	code, out = asUser("pkill -9 ssh-agent; sleep 0.2; " + isSock + " && rm -f " + sock +
		" && stat " + sock + " 2>&1 | grep -q 'No such file'")
	s.Require().Equalf(0, code, "rm must clear the stale socket: %s", out)
	asUser("rm -rf " + dir + "/.app-listener; rm -d " + dir + "; true")
	code, out = s.exec(c, []string{"stat", "-c", "%F", dir + "/.app-listener/README"})
	s.Require().Equalf(0, code, "a whitelisted rm emptied or removed the guarded root: %s", out)

	code, out = s.exec(c, []string{"sh", "-c", "umount /run/user/1000 && /app-listener daemon logout 1000"})
	s.Require().Equalf(0, code, "logout: %s", out)
	s.Require().NotContainsf(s.readDaemonLog(c), reloadRefused, "logout reload refused:\n%s", s.readDaemonLog(c))
	code, out = s.exec(c, []string{"sh", "-c", mount})
	s.Require().Equalf(0, code, "remount: %s", out)
	// A process of the user that beat the daemon to the fresh mount may have planted listeners.
	code, out = asUser("mkdir -m 0700 " + dir)
	s.Require().Equalf(0, code, "pre-creating: %s", out)
	code, out = s.exec(c, []string{"/app-listener", "daemon", "login", "1000"})
	s.Require().NotEqualf(0, code, "a runtime dir made before the daemon must fail the login: %s", out)
	s.exec(c, []string{"rmdir", dir})
	login()

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not exit")
	code, out = s.exec(c, []string{"/app-listener", "daemon", "login", "1000"})
	s.Require().Equalf(0, code, "with no daemon running a login goes ahead: %s", out)
}

// A runtime dir the user made before the daemon started (an agent unit's `mkdir -p`) is guarded as
// found: the daemon plants the root-owned sentinel then too, so a TERMed agent's cleanup can't leave
// the root empty for an rmdir and a re-made, unguarded dir. An empty watch root's own rmdir is judged.
func (s *IntegrationSuite) TestDaemon_RuntimeDirFoundAtStart() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	code, out := s.exec(c, []string{"sh", "-c", "command -v ssh-agent >/dev/null 2>&1 || " +
		"(timeout 120 apt-get update -qq && timeout 120 apt-get install -y -qq --no-install-recommends openssh-client >/dev/null)"})
	if code != 0 {
		s.T().Skipf("openssh-client unavailable: %s", out)
	}
	const dir = "/run/user/1000/ssh-agent"
	const sock = dir + "/socket"
	const vault = "/home/u1000/vault"
	code, out = s.exec(c, []string{"sh", "-c", "(id -u 1000 >/dev/null 2>&1 || useradd -u 1000 -m u1000) && " +
		"mkdir -p /watch /etc/app-listener /run/user/1000 /home/u1000 && chown 1000:1000 /home/u1000 && " +
		"mount -t tmpfs -o mode=0700,uid=1000,gid=1000 tmpfs /run/user/1000"})
	s.Require().Equalf(0, code, "setup: %s", out)
	asUser := func(cmd string) (int, string) {
		return s.exec(c, []string{"setpriv", "--reuid=1000", "--regid=1000", "--clear-groups", "sh", "-c", cmd})
	}
	code, out = asUser("mkdir -m 0700 " + dir + " && mkdir -m 0700 " + vault)
	s.Require().Equalf(0, code, "pre-creating: %s", out)

	const always = "[watch /watch]\nneed_encryption: false\n/usr/bin/true\n"
	conf := always + "\n[watch " + vault + "]\nneed_encryption: false\n/usr/bin/true\n" +
		"\n[watch " + dir + "]\nruntime: true\n/usr/bin/ssh-agent MKNOD,DELETE,STAT\n" +
		"/usr/bin/ssh-add CONNECT,STAT\n/usr/bin/stat STAT\nneed_encryption: false\n"
	s.exec(c, []string{"sh", "-c", "cat > /etc/app-listener/daemon.conf <<'EOF'\n" + conf + "EOF"})
	s.launchDaemon(c)
	s.awaitDaemonUp(c, conf)

	code, out = s.exec(c, []string{"stat", "-c", "%u %a", dir + "/.app-listener/README"})
	s.Require().Equalf(0, code, "no sentinel in a runtime dir found at start: %s\n%s", out, s.readDaemonLog(c))
	s.Require().Equalf("0 444", out, "the sentinel stays root's")

	isGone := "stat " + sock + " 2>&1 | grep -q 'No such file'"
	code, out = asUser("(nohup /usr/bin/ssh-agent -D -a " + sock + " >/tmp/agent.log 2>&1 &) && " +
		"for i in $(seq 50); do " + isGone + " || break; sleep 0.1; done; pkill -TERM ssh-agent; " +
		"for i in $(seq 50); do " + isGone + " && break; sleep 0.1; done; rmdir " + dir + " 2>&1")
	s.Require().NotEqualf(0, code, "the user removed a runtime root after its agent cleaned up: %s", out)
	code, out = s.exec(c, []string{"stat", "-c", "%F", dir + "/.app-listener/README"})
	s.Require().Equalf(0, code, "the runtime root is gone: %s", out)

	code, out = asUser("rmdir " + vault + " 2>&1")
	s.Require().NotEqualf(0, code, "the user removed an empty watch root: %s", out)
	denied := false
	for _, e := range parseDaemonEvents(s.readDaemonLog(c)) {
		denied = denied || (e.Denied && e.Op == "DELETE" && e.Comm == "rmdir")
	}
	s.Require().Truef(denied, "no DAEMON DENIED op=DELETE for rmdir of the root:\n%s", s.readDaemonLog(c))

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not exit")
}
