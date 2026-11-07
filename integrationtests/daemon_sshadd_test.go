package integrationtests

import (
	"strings"
	"time"
)

// A whitelisted ssh client must not hand a guarded key to its caller: ssh-add sends the keys it reads
// to whatever agent SSH_AUTH_SOCK names (it may reach only an agent inside a guarded tree), and ssh
// and sshd print every line of a config they can't parse (ssh -F, sshd run by anyone but root).
func (s *IntegrationSuite) TestDaemon_SSHKeyDisclosure() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	code, out := s.exec(c, []string{"sh", "-c", "command -v sshd >/dev/null 2>&1 || (timeout 120 apt-get update -qq && " +
		"DEBIAN_FRONTEND=noninteractive timeout 300 apt-get install -y -qq --no-install-recommends openssh-client " +
		"openssh-server >/dev/null)"})
	if code != 0 {
		s.T().Skipf("openssh unavailable: %s", out)
	}
	const ssh = "/home/u1000/.ssh"
	asUser := func(cmd string) (int, string) {
		return s.exec(c, []string{"setpriv", "--reuid=1000", "--regid=1000", "--clear-groups", "sh", "-c", cmd})
	}
	code, out = s.exec(c, []string{"sh", "-c", "(id -u 1000 >/dev/null 2>&1 || useradd -u 1000 -m u1000) && " +
		"mkdir -p /etc/app-listener /exploits /home/u1000 && chown 1000:1000 /home/u1000"})
	s.Require().Equalf(0, code, "setup: %s", out)
	code, out = asUser("mkdir -m 700 -p " + ssh + " && ssh-keygen -q -t ed25519 -N '' -f " + ssh + "/id_ed25519")
	s.Require().Equalf(0, code, "keygen: %s", out)
	_, key := s.exec(c, []string{"sed", "-n", "2p", ssh + "/id_ed25519"})
	key = strings.TrimSpace(key)
	s.Require().NotEmpty(key)
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/unix_connect"), "/exploits/unix_connect", 0o755),
		"copy unix_connect (run make -C integrationtests/exploits)")
	s.capturePID(c, "setpriv --reuid=1000 --regid=1000 --clear-groups nohup /exploits/unix_connect serve /tmp/agent.sock "+
		"stream >/tmp/fake-agent.log 2>&1 &")

	sshd := "/usr/sbin/sshd"
	if code, _ := s.exec(c, []string{"test", "-x", "/usr/bin/sshd"}); code == 0 {
		sshd = "/usr/bin/sshd"
	}
	// binder stands for a whitelisted MKNOD binary that isn't an agent (systemd binding a user's socket
	// unit); relay/ssh-add, a confined name.
	const binder, fakeAdd = "/opt/relay/binder", "/opt/relay/ssh-add"
	code, out = s.exec(c, []string{"sh", "-c", "mkdir -p /opt/relay && cp /exploits/unix_connect " + binder +
		" && cp /exploits/unix_connect " + fakeAdd})
	s.Require().Equalf(0, code, "setup relay: %s", out)
	conf := "[watch " + ssh + "]\nneed_encryption: false\n/usr/bin/ssh\n/usr/bin/ssh-add\n/usr/bin/ssh-agent\n" +
		"/usr/bin/ssh-keygen\n" + sshd + " READ\n" + binder + " MKNOD\n" + fakeAdd + " CONNECT\n"
	s.startDaemon(c, conf)

	code, out = asUser("SSH_AUTH_SOCK=/tmp/agent.sock ssh-add " + ssh + "/id_ed25519")
	s.Require().NotEqualf(0, code, "ssh-add sent the key to an agent outside every guarded tree: %s", out)
	// unguarded_peer refuses with EPERM, a resource's whitelist with EACCES.
	s.Require().Truef(strings.Contains(strings.ToLower(out), "operation not permitted") ||
		strings.Contains(strings.ToLower(out), "permission denied"), "the connect must be refused: %s\n%s", out,
		s.readDaemonLog(c))
	s.requireDaemonDenial(c, "CONNECT", "ssh-add")

	// A listener inside the tree whose reader is no agent: the binder itself, the exe a socket unit's
	// Accept=yes hands the connection to (read, splice). ssh-add's request must reach none of them.
	asUser1000 := "setpriv --reuid=1000 --regid=1000 --clear-groups nohup "
	s.capturePID(c, asUser1000+binder+" relay "+ssh+"/r1.sock >/tmp/relay1.log 2>&1 &")
	s.capturePID(c, asUser1000+binder+" relay "+ssh+"/r2.sock /bin/cat >/tmp/relay2.log 2>&1 &")
	s.capturePID(c, asUser1000+binder+" relay "+ssh+"/r3.sock /exploits/unix_connect splicefd 0 "+
		">/tmp/relay3.log 2>&1 &")
	s.capturePID(c, asUser1000+binder+" serve "+ssh+"/d.sock dgram >/tmp/relayd.log 2>&1 &")
	for _, n := range []string{"r1", "r2", "r3"} {
		code, out = asUser("for i in $(seq 50); do grep -q relaying /tmp/relay" + n[1:] + ".log && exit 0; " +
			"sleep 0.1; done; cat /tmp/relay" + n[1:] + ".log; exit 1")
		s.Require().Equalf(0, code, "relay %s up: %s\n%s", n, out, s.readDaemonLog(c))
		asUser("SSH_AUTH_SOCK=" + ssh + "/" + n + ".sock timeout 5 ssh-add " + ssh + "/id_ed25519")
	}
	time.Sleep(time.Second)
	_, logs := s.exec(c, []string{"sh", "-c", "cat /tmp/relay1.log /tmp/relay2.log /tmp/relay3.log"})
	s.Require().NotContainsf(logs, "ssh-ed25519", "a non-agent read ssh-add's key:\n%s\n%s", logs, s.readDaemonLog(c))
	s.Require().Containsf(logs, "splice denied", "the spliced read must be refused: %s", logs)
	code, out = asUser(fakeAdd + " dgram " + ssh + "/d.sock")
	s.Require().Equalf(1, code, "a confined exe's datagram reached a non-agent: %s\n%s", out, s.readDaemonLog(c))

	for _, cmd := range []string{
		"ssh -F " + ssh + "/id_ed25519 -G x",
		"ssh -4F" + ssh + "/id_ed25519 -G x",
		"ssh -" + strings.Repeat("C", 40) + "F " + ssh + "/id_ed25519 -G x",
		sshd + " -t -f " + ssh + "/id_ed25519",
	} {
		_, out = asUser(cmd + " 2>&1")
		s.Require().NotContainsf(strings.ToLower(out), strings.ToLower(key[:24]), "%s printed the key: %s", cmd, out)
	}

	s.capturePID(c, "setpriv --reuid=1000 --regid=1000 --clear-groups nohup ssh-agent -D -a "+ssh+"/agent.sock "+
		">/tmp/agent.log 2>&1 &")
	// sh's `test` would be denied STAT in the guarded tree: retry the whitelisted ssh-add itself.
	code, out = asUser("for i in $(seq 50); do SSH_AUTH_SOCK=" + ssh + "/agent.sock ssh-add " + ssh +
		"/id_ed25519 2>/dev/null && exit 0; sleep 0.1; done; cat /tmp/agent.log; exit 1")
	s.Require().Equalf(0, code, "ssh-add to an agent inside the guarded tree: %s\n%s", out, s.readDaemonLog(c))

	s.exec(c, []string{"sh", "-c", "pkill -f 'unix_connect serve'; pkill -f /opt/relay/; pkill ssh-agent; pkill -f 'app-listener daemon' || true"})
}
