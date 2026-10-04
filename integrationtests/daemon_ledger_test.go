package integrationtests

import (
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

const (
	ledgerMarker = "TOP-SECRET-LEDGER-6B3D"
	unconfirmed  = "binary-unconfirmed"
)

// restartUnconfirmed stops the daemon, runs whileStopped, and starts it without the harness
// confirming daemon.conf's binaries: what changed meanwhile is judged by the binary ledger alone.
func (s *IntegrationSuite) restartUnconfirmed(c testcontainers.Container, config string, whileStopped func()) {
	s.sigDaemon(c, "TERM")
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not exit after SIGTERM")
	if whileStopped != nil {
		whileStopped()
	}
	s.launchDaemonUnconfirmed(c)
	s.awaitDaemonUp(c, config)
}

// unconfirmedLine returns the daemon's binary-unconfirmed line for path, "" if none.
func (s *IntegrationSuite) unconfirmedLine(c testcontainers.Container, path string) string {
	for _, l := range strings.Split(s.readDaemonLog(c), "\n") {
		if strings.Contains(l, unconfirmed) && strings.Contains(l, "path="+path+" ") {
			return l
		}
	}
	return ""
}

// awaitReadAs polls readAs(bin) for the refresh secret: admission follows a reload or re-sync.
func (s *IntegrationSuite) awaitReadAs(c testcontainers.Container, bin string) bool {
	for dl := time.Now().Add(60 * time.Second); time.Now().Before(dl); time.Sleep(time.Second) {
		if strings.Contains(s.readAs(c, bin), "STOLEN|"+refreshMarker) {
			return true
		}
	}
	return false
}

func (s *IntegrationSuite) trustBinaries(c testcontainers.Container, args ...string) string {
	code, out := s.exec(c, append([]string{"/app-listener", "trust-binaries"}, args...))
	s.Require().Equalf(0, code, "trust-binaries %v: %s", args, out)
	return out
}

// Nothing guards a whitelisted path while the daemon is stopped: a binary replaced meanwhile, in a
// home (root's) or a user-owned dir outside /home, is refused at the next start until
// `trust-binaries` confirms it, and the confirmation persists. A root-placed one needs none.
func (s *IntegrationSuite) TestDaemon_Ledger_BinaryReplacedWhileStoppedRefusedUntilConfirmed() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	s.copySwapFixtures(c)
	s.installFakeSystemctl(c)

	const (
		secret  = "/protected/secret"
		inHome  = "/root/app/reader"
		userDir = "/opt/userapp/reader"
		sysBin  = "/usr/local/bin/sysreader"
	)
	code, out := s.exec(c, []string{"sh", "-c", "mkdir -p /protected /etc/app-listener /root/app /opt/userapp && printf '" +
		ledgerMarker + "' > " + secret + " && for b in " + inHome + " " + userDir + " " + sysBin + "; do cp " +
		swapReaderPath + " $b; done && chown -R 65534 /opt/userapp 2>&1"})
	s.Require().Equalf(0, code, "setup: %s", out)
	config := `[watch /protected]
need_encryption: false
` + inHome + `
` + userDir + `
` + sysBin
	s.startDaemon(c, config)
	for _, b := range []string{inHome, userDir, sysBin} {
		_, out := s.exec(c, []string{b, secret})
		s.Require().Containsf(out, "STOLEN|"+ledgerMarker, "baseline: %s reads its resource: %s", b, out)
	}

	// A new inode with new content at each line: root's in its home and at a root-placed name, the
	// user's in the dir it owns.
	replace := "cp " + swapReaderPath + " $(dirname $B)/.new && printf x >> $(dirname $B)/.new && mv -f $(dirname $B)/.new $B"
	s.restartUnconfirmed(c, config, func() {
		for _, cmd := range []string{
			"B=" + inHome + "; " + replace,
			"B=" + sysBin + "; " + replace,
			nobodyRun + "sh -c " + shQuote("B="+userDir+"; "+replace),
		} {
			code, out := s.exec(c, []string{"sh", "-c", cmd + " 2>&1"})
			s.Require().Equalf(0, code, "fixture: replacing while stopped (%s): %s", cmd, out)
		}
	})

	for _, b := range []string{inHome, userDir} {
		s.Require().NotEmptyf(s.unconfirmedLine(c, b), "%s changed while stopped was not reported; daemon log:\n%s", b,
			s.readDaemonLog(c))
	}
	s.Require().Empty(s.unconfirmedLine(c, sysBin), "a root-placed replacement needs no confirmation")
	_, out = s.exec(c, []string{sysBin, secret})
	s.Require().Containsf(out, "STOLEN|"+ledgerMarker, "the root-placed replacement must be admitted: %s", out)
	_, out = s.exec(c, []string{userDir, secret})
	s.Require().NotContainsf(out, ledgerMarker, "the user's replacement was admitted unconfirmed: %s", out)
	s.assertNeverStolen(c, inHome, secret, ledgerMarker)

	list := s.trustBinaries(c, "--list")
	s.Require().Containsf(list, inHome, "trust-binaries --list must show the refused binaries: %s", list)
	s.Require().Contains(list, userDir)
	s.Require().NotContains(list, sysBin)
	out = s.trustBinaries(c, "--yes")
	s.Require().Containsf(out, "confirmed", "trust-binaries --yes: %s", out)
	for _, b := range []string{inHome, userDir} {
		ok, out := s.awaitStolen(c, b, secret, ledgerMarker, 60*time.Second)
		s.Require().Truef(ok, "%s not admitted after trust-binaries: %s\ndaemon log:\n%s", b, out, s.readDaemonLog(c))
	}
	s.Require().Contains(s.trustBinaries(c, "--list"), "No refused binaries", "confirmed rows must leave the list")

	s.restartUnconfirmed(c, config, nil)
	s.Require().NotContainsf(s.readDaemonLog(c), unconfirmed, "the confirmation did not persist across a restart")
	for _, b := range []string{inHome, userDir, sysBin} {
		_, out := s.exec(c, []string{b, secret})
		s.Require().Containsf(out, "STOLEN|"+ledgerMarker, "after a restart %s must read: %s", b, out)
	}

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// A version the updater installed while the daemon ran is recorded and admitted again after a
// restart; a catalog glob match created while it was stopped is refused at the startup refresh,
// though its bytes are the client's own, until confirmed.
func (s *IntegrationSuite) TestDaemon_Ledger_GlobMatchPlantedWhileStoppedRefused() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.installFakeSystemctl(c)
	s.startRefreshDaemon(c)

	updated := s.installDiscordVersion(c, "0.0.2")
	s.Require().Truef(s.awaitReadAs(c, updated), "the updater's version was not admitted live; daemon log:\n%s",
		s.readDaemonLog(c))

	plant := discordDir + "/0.0.7/Discord"
	s.restartUnconfirmed(c, "[watch "+discordDir+"/sentry]", func() {
		code, out := s.exec(c, []string{"sh", "-c", "mkdir -p " + discordDir + "/0.0.7 && cp /usr/bin/dash " + plant + " 2>&1"})
		s.Require().Equalf(0, code, "fixture: planting while stopped: %s", out)
	})
	for dl := time.Now().Add(60 * time.Second); time.Now().Before(dl) && s.unconfirmedLine(c, plant) == ""; {
		time.Sleep(500 * time.Millisecond)
	}
	s.Require().NotEmptyf(s.unconfirmedLine(c, plant), "the startup refresh did not refuse the plant; daemon log:\n%s",
		s.readDaemonLog(c))
	s.Require().NotContainsf(s.readAs(c, plant), refreshMarker, "a match planted while stopped read the resource")
	s.Require().Empty(s.unconfirmedLine(c, updated), "the updater's recorded version needs no confirmation")
	s.Require().Containsf(s.readAs(c, updated), "STOLEN|"+refreshMarker, "the updater's version must survive a restart")
	s.Require().Contains(s.readAs(c, discordClient), "STOLEN|"+refreshMarker, "control: the configured client reads")

	s.trustBinaries(c, "--yes")
	s.Require().Truef(s.awaitReadAs(c, plant), "confirmed, the match must be admitted; daemon log:\n%s",
		s.readDaemonLog(c))

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}
