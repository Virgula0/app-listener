package integrationtests

import "strings"

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
		"mkdir -p /etc/app-listener /exploits"})
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
	conf := "[watch " + ssh + "]\nneed_encryption: false\n/usr/bin/ssh\n/usr/bin/ssh-add\n/usr/bin/ssh-agent\n" +
		"/usr/bin/ssh-keygen\n" + sshd + " READ\n"
	s.startDaemon(c, conf)

	code, out = asUser("SSH_AUTH_SOCK=/tmp/agent.sock ssh-add " + ssh + "/id_ed25519")
	s.Require().NotEqualf(0, code, "ssh-add sent the key to an agent outside every guarded tree: %s", out)
	s.Require().Containsf(strings.ToLower(out), "permission denied", "the connect must be refused: %s\n%s", out,
		s.readDaemonLog(c))
	s.requireDaemonDenial(c, "CONNECT", "ssh-add")

	for _, cmd := range []string{
		"ssh -F " + ssh + "/id_ed25519 -G x",
		"ssh -4F" + ssh + "/id_ed25519 -G x",
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

	s.exec(c, []string{"sh", "-c", "pkill -f 'unix_connect serve'; pkill ssh-agent; pkill -f 'app-listener daemon' || true"})
}
