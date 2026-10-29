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
