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
