package integrationtests

import (
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

// awaitStolen runs bin against secret until it prints marker (true) or timeout. Each denial makes
// the daemon re-sync its whitelist (resyncMinInterval).
func (s *IntegrationSuite) awaitStolen(c testcontainers.Container, bin, secret, marker string,
	timeout time.Duration) (bool, string) {
	var out string
	for dl := time.Now().Add(timeout); time.Now().Before(dl); time.Sleep(time.Second) {
		_, out = s.exec(c, []string{"sh", "-c", bin + " " + secret + " 2>&1"})
		if strings.Contains(out, "STOLEN|"+marker) {
			return true, out
		}
	}
	return false, out
}

// A package upgrade replaces a root-owned whitelisted binary by temp file + rename: the new inode
// is admitted by the re-sync alone, with no hook and no reload.
func (s *IntegrationSuite) TestDaemon_LiveAdmission_SystemBinaryReplacedNoReload() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	s.copySwapFixtures(c)

	const marker = "TOP-SECRET-LIVE-SYSTEM-0B7E"
	s.exec(c, []string{"sh", "-c", "mkdir -p /protected /etc/app-listener /opt/app/bin && printf '" + marker +
		"' > /protected/secret && cp " + swapReaderPath + " /opt/app/bin/app"})
	s.startDaemon(c, `[watch /protected]
need_encryption: false
/opt/app/bin/app`)
	_, out := s.exec(c, []string{"/opt/app/bin/app", "/protected/secret"})
	s.Require().Containsf(out, "STOLEN|"+marker, "baseline: the whitelisted binary reads its resource: %s", out)

	old := s.inodeOf(c, "/opt/app/bin/app")
	code, out := s.exec(c, []string{"sh", "-c", "cp " + swapReaderPath + " /opt/app/bin/.app.new && " +
		"printf x >> /opt/app/bin/.app.new && mv -f /opt/app/bin/.app.new /opt/app/bin/app 2>&1"})
	s.Require().Equalf(0, code, "the upgrade: %s", out)
	s.Require().NotEqual(old, s.inodeOf(c, "/opt/app/bin/app"), "fixture: the upgrade must make a new inode")

	ok, out := s.awaitStolen(c, "/opt/app/bin/app", "/protected/secret", marker, 40*time.Second)
	s.Require().Truef(ok, "the upgraded system binary was not admitted: %s\ndaemon log:\n%s", out, s.readDaemonLog(c))
	s.Require().NotContains(s.readDaemonLog(c), "configuration reloaded", "admission must need no reload")

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// A new match of an absolute catalog glob (JetBrains' /opt/*/jbr/bin/java) is admitted live when
// root installed it; the same name in a directory a user owns is not.
func (s *IntegrationSuite) TestDaemon_LiveAdmission_NewCatalogGlobMatch() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	s.copySwapFixtures(c)

	const (
		marker = "TOP-SECRET-JBR-GLOB-5A61"
		secret = "/root/.config/JetBrains/secret"
	)
	s.exec(c, []string{"sh", "-c", "mkdir -p /etc/app-listener /root/.config/JetBrains /opt/idea/jbr/bin && printf '" +
		marker + "' > " + secret + " && cp " + swapReaderPath + " /opt/idea/jbr/bin/java"})
	s.startDaemon(c, `[watch /root/.config/JetBrains]
need_encryption: false
/opt/idea/jbr/bin/java`)

	s.exec(c, []string{"sh", "-c", "mkdir -p /opt/clion/jbr/bin /opt/evil/jbr/bin && chown 65534 /opt/evil/jbr/bin && " +
		"cp " + swapReaderPath + " /opt/clion/jbr/bin/java && " + nobodyRun + "cp " + swapReaderPath +
		" /opt/evil/jbr/bin/java"})

	ok, out := s.awaitStolen(c, "/opt/clion/jbr/bin/java", secret, marker, 40*time.Second)
	s.Require().Truef(ok, "a root-installed match of /opt/*/jbr/bin/java was not admitted: %s\ndaemon log:\n%s",
		out, s.readDaemonLog(c))
	s.Require().NotContains(s.readDaemonLog(c), "configuration reloaded", "admission must need no reload")

	// A reload rebuilds the trusted set from config paths, which do not name the match: it must
	// keep its library allowlist (TRUSTED_BINARY), or a user's LD_PRELOAD runs inside it.
	s.exec(c, []string{"mkdir", "-p", "/tmp/u"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/lib_probe.so"), "/tmp/u/evil.so", 0o755),
		"copy lib_probe.so")
	s.exec(c, []string{"chown", "-R", "65534:65534", "/tmp/u"})
	s.reloadAndAwait(c)
	ok, out = s.awaitStolen(c, "/opt/clion/jbr/bin/java", secret, marker, 10*time.Second)
	s.Require().Truef(ok, "after a reload the live-admitted match lost its resource: %s", out)
	_, out = s.exec(c, []string{"sh", "-c", "LD_PRELOAD=/tmp/u/evil.so /opt/clion/jbr/bin/java " + secret + " 2>&1"})
	s.Require().NotContainsf(out, libProbeMarker, "after a reload a user library preloaded into the live-admitted "+
		"match: %s", out)

	s.assertNeverStolen(c, "/opt/evil/jbr/bin/java", secret, marker)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}
