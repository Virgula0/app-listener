package integrationtests

import "strings"

// ssh-add reads the keys of a guarded ~/.ssh and sends them to whatever agent SSH_AUTH_SOCK names: a
// whitelisted ssh-add may reach only an agent whose socket is inside a guarded tree.
func (s *IntegrationSuite) TestDaemon_SSHAddConfinedToGuardedAgents() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	code, out := s.exec(c, []string{"sh", "-c", "command -v ssh-add >/dev/null 2>&1 || (timeout 120 apt-get update -qq && " +
		"DEBIAN_FRONTEND=noninteractive timeout 300 apt-get install -y -qq --no-install-recommends openssh-client >/dev/null)"})
	if code != 0 {
		s.T().Skipf("openssh-client unavailable: %s", out)
	}
	code, out = s.exec(c, []string{"sh", "-c", "mkdir -p /etc/app-listener /exploits /root/.ssh && chmod 700 /root/.ssh && " +
		"ssh-keygen -q -t ed25519 -N '' -f /root/.ssh/id_ed25519"})
	s.Require().Equalf(0, code, "setup: %s", out)
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/unix_connect"), "/exploits/unix_connect", 0o755),
		"copy unix_connect (run make -C integrationtests/exploits)")
	s.capturePID(c, "nohup /exploits/unix_connect serve /tmp/agent.sock stream >/tmp/fake-agent.log 2>&1 &")

	conf := "[watch /root/.ssh]\nneed_encryption: false\n/usr/bin/ssh-add\n/usr/bin/ssh-agent\n/usr/bin/ssh-keygen\n"
	s.startDaemon(c, conf)

	code, out = s.exec(c, []string{"sh", "-c", "SSH_AUTH_SOCK=/tmp/agent.sock ssh-add /root/.ssh/id_ed25519"})
	s.Require().NotEqualf(0, code, "ssh-add sent the key to an agent outside every guarded tree: %s", out)
	s.Require().Containsf(strings.ToLower(out), "permission denied", "the connect must be refused: %s\n%s", out,
		s.readDaemonLog(c))
	s.requireDaemonDenial(c, "CONNECT", "ssh-add")

	s.capturePID(c, "nohup ssh-agent -D -a /root/.ssh/agent.sock >/tmp/agent.log 2>&1 &")
	code, out = s.exec(c, []string{"sh", "-c", "for i in $(seq 50); do [ -S /root/.ssh/agent.sock ] && " +
		"SSH_AUTH_SOCK=/root/.ssh/agent.sock exec ssh-add /root/.ssh/id_ed25519; sleep 0.1; done; cat /tmp/agent.log; exit 1"})
	s.Require().Equalf(0, code, "ssh-add to an agent inside the guarded tree: %s\n%s", out, s.readDaemonLog(c))

	s.exec(c, []string{"sh", "-c", "pkill -f 'unix_connect serve'; pkill ssh-agent; pkill -f 'app-listener daemon' || true"})
}
