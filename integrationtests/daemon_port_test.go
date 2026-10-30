package integrationtests

import (
	"strings"
	"testing"

	"github.com/testcontainers/testcontainers-go"
)

const portSecret = "INET-SERVICE-SECRET"

// A [port tcp/N] section guards a loopback TCP service (bitcoind RPC, the adb server): only its
// MKNOD binary listens on the port, and only on loopback; only its CONNECT binaries reach it,
// through every connect path; a connection handed to another exe is judged by that exe.
func (s *IntegrationSuite) TestDaemon_PortGuard() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const port = "8332"
	s.exec(c, []string{"mkdir", "-p", "/exploits", "/etc/app-listener", "/watch"})
	host := absPath("./exploits/inet_connect")
	for _, name := range []string{"inet_connect", "inet_connect_ok", "inet_connect_srv"} {
		s.Require().NoError(c.CopyFileToContainer(s.ctx, host, "/exploits/"+name, 0o755),
			"copy %s (run make -C integrationtests/exploits)", name)
	}
	const always = "[watch /watch]\nneed_encryption: false\n/usr/bin/true\n"
	conf := always + "\n[launch_rules]\ndeny_args /exploits/inet_connect_ok dumpprivkey\n\n" +
		"[port tcp/" + port + "]\n/exploits/inet_connect_srv MKNOD\n/exploits/inet_connect_ok CONNECT\n"
	// The daemon self-guards /etc/app-listener: it runs on a copy the test can rewrite (the later
	// --config wins); startDaemon's /etc one only feeds confirmConfigBinaries.
	const liveConf = "/tmp/port-daemon.conf"
	writeLive := func(body string) {
		code, out := s.exec(c, []string{"sh", "-c", "cat > " + liveConf + " <<'EOF'\n" + body + "EOF"})
		s.Require().Equalf(0, code, "writing %s: %s", liveConf, out)
	}
	writeLive(conf)
	s.startDaemon(c, conf, "--config", liveConf)

	s.T().Run("listen", func(t *testing.T) {
		for what, args := range map[string][]string{
			"a non-whitelisted listener":    {"/exploits/inet_connect", "listen", "127.0.0.1", port},
			"a CONNECT-only binary":         {"/exploits/inet_connect_ok", "listen", "127.0.0.1", port},
			"the server on 0.0.0.0":         {"/exploits/inet_connect_srv", "listen", "0.0.0.0", port},
			"the server on ::":              {"/exploits/inet_connect_srv", "listen", "::", port},
			"the server on ::ffff:0.0.0.0":  {"/exploits/inet_connect_srv", "listen", "::ffff:0.0.0.0", port},
			"the server on the LAN address": {"sh", "-c", "/exploits/inet_connect_srv listen $(hostname -i | cut -d' ' -f1) " + port},
		} {
			code, out := s.exec(c, args)
			if code == 2 {
				t.Logf("%s: unsupported here: %s", what, out)
				continue
			}
			s.Require().Equalf(1, code, "%s must be refused the port (exit 1): %s\n%s", what, out, s.readDaemonLog(c))
		}
		s.requireDaemonDenial(c, "MKNOD", "inet_connect")
	})

	s.capturePID(c, "nohup /exploits/inet_connect_srv serve 127.0.0.1 "+port+" >/tmp/serve4.log 2>&1 &")
	s.capturePID(c, "nohup /exploits/inet_connect_srv serve ::1 "+port+" >/tmp/serve6.log 2>&1 &")
	code, out := s.exec(c, []string{"sh", "-c", "for i in $(seq 50); do /exploits/inet_connect_ok connect 127.0.0.1 " +
		port + " >/dev/null && exit 0; sleep 0.1; done; cat /tmp/serve4.log; exit 1"})
	s.Require().Equalf(0, code, "the whitelisted server and client never met: %s\n%s", out, s.readDaemonLog(c))

	dests := []string{"127.0.0.1", "0.0.0.0", "::ffff:127.0.0.1", "::1", "::"}
	s.T().Run("connect", func(t *testing.T) {
		for _, d := range dests {
			code, out := s.exec(c, []string{"/exploits/inet_connect_ok", "connect", d, port})
			if code == 2 {
				t.Logf("%s: no listener reachable here: %s", d, out)
				continue
			}
			s.Require().Equalf(0, code, "whitelisted connect to %s: %s", d, out)
			s.requirePortDenied(c, t, []string{"/exploits/inet_connect", "connect", d, port}, "inet_connect")
		}
		for _, mode := range []string{"uring", "fastopen"} {
			s.requirePortDenied(c, t, []string{"/exploits/inet_connect", mode, "127.0.0.1", port}, "inet_connect")
		}
		s.requirePortDenied(c, t, []string{"/exploits/inet_connect_ok", "connect", "127.0.0.1", port, "dumpprivkey"},
			"inet_connect_ok", "ARGV")
	})

	s.T().Run("handoff", func(t *testing.T) {
		for _, mode := range []string{"inherit", "inheritsplice", "inheritzc", "inheritmmap", "passfd"} {
			s.requirePortDenied(c, t, []string{"/exploits/inet_connect_ok", mode, "127.0.0.1", port, "/exploits/inet_connect"},
				"inet_connect")
		}
	})

	s.reloadAndAwait(c)
	s.Require().NotContainsf(s.readDaemonLog(c), reloadRefused, "reload refused:\n%s", s.readDaemonLog(c))
	code, out = s.exec(c, []string{"/exploits/inet_connect_ok", "connect", "127.0.0.1", port})
	s.Require().Equalf(0, code, "whitelisted connect after a reload: %s", out)
	s.requirePortDenied(c, s.T(), []string{"/exploits/inet_connect", "connect", "127.0.0.1", port}, "inet_connect")

	// Dropping the section stops guarding the port: there is no vault to leave unlocked.
	writeLive(always)
	s.reloadAndAwait(c)
	s.Require().NotContainsf(s.readDaemonLog(c), reloadRefused, "dropping a port section refused:\n%s", s.readDaemonLog(c))
	code, out = s.exec(c, []string{"/exploits/inet_connect", "connect", "127.0.0.1", port})
	s.Require().Equalf(0, code, "a dropped port section still guards the port: %s\n%s", out, s.readDaemonLog(c))

	s.exec(c, []string{"sh", "-c", "pkill -f inet_connect; pkill -f 'app-listener daemon' || true"})
}

// requirePortDenied runs args, which must be refused by the kernel (exit 1, no secret, a DAEMON
// DENIED line for comm, op CONNECT unless given).
func (s *IntegrationSuite) requirePortDenied(c testcontainers.Container, t *testing.T, args []string, comm string,
	op ...string) {
	code, out := s.exec(c, args)
	if code == 2 {
		t.Logf("%v: unsupported here: %s", args, out)
		return
	}
	s.Require().Falsef(strings.Contains(out, portSecret), "%v reached the guarded port: %s", args, out)
	s.Require().Equalf(1, code, "%v: want a kernel denial (exit 1), got %d: %s\nservers:\n%s\n%s", args, code, out,
		s.serverLogs(c), s.readDaemonLog(c))
	s.requireDaemonDenial(c, append(op, "CONNECT")[0], comm)
}

// requireDaemonDenial asserts the daemon logged a denial of op to comm.
func (s *IntegrationSuite) requireDaemonDenial(c testcontainers.Container, op, comm string) {
	log := s.readDaemonLog(c)
	for _, e := range parseDaemonEvents(log) {
		if e.Denied && e.Op == op && e.Comm == comm {
			return
		}
	}
	s.Require().Failf("no denial logged", "no DAEMON DENIED op=%s comm=%s:\n%s", op, comm, log)
}

func (s *IntegrationSuite) serverLogs(c testcontainers.Container) string {
	_, out := s.exec(c, []string{"sh", "-c", "tail -n 20 /tmp/serve4.log /tmp/serve6.log 2>&1"})
	return out
}
