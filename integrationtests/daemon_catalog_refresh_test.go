package integrationtests

import (
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

const (
	refreshMarker = "TOP-SECRET-LIVE-REFRESH-9E34"
	refreshSecret = discordDir + "/sentry/secret"
	refreshLine   = "DAEMON catalog-refresh"
)

// startRefreshDaemon guards Discord's sentry dir with root's dash copy standing in for the client.
func (s *IntegrationSuite) startRefreshDaemon(c testcontainers.Container) {
	s.exec(c, []string{"sh", "-c", "mkdir -p /etc/app-listener /exploits " + discordDir + "/sentry $(dirname " +
		discordClient + ") && cp /usr/bin/dash " + discordClient + " && printf '" + refreshMarker + "' > " +
		refreshSecret})
	s.startDaemon(c, `[watch `+discordDir+`/sentry]
need_encryption: false
`+discordClient)
	s.Require().Contains(s.readDaemonLog(c), "catalog watch:", "the catalog watch must start")
}

// readAs runs a dash-compatible binary that reads the secret with a builtin (no child exec).
func (s *IntegrationSuite) readAs(c testcontainers.Container, bin string) string {
	_, out := s.exec(c, []string{"sh", "-c", bin + " -c 'read l < " + refreshSecret + "; echo \"STOLEN|$l\"' 2>&1"})
	return out
}

// installDiscordVersion has the whitelisted client install a new version, as its updater does.
func (s *IntegrationSuite) installDiscordVersion(c testcontainers.Container, version string) string {
	bin := discordDir + "/" + version + "/Discord"
	code, out := s.exec(c, []string{discordClient, "-c", "mkdir -p " + discordDir + "/" + version + " && cat /usr/bin/dash > " +
		bin + " && chmod 755 " + bin})
	s.Require().Equalf(0, code, "the client must install its own update: %s", out)
	return bin
}

func (s *IntegrationSuite) TestDaemon_CatalogRefresh_UpdaterVersionAdmittedLive() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.startRefreshDaemon(c)
	s.Require().Contains(s.readAs(c, discordClient), "STOLEN|"+refreshMarker, "baseline")

	bin := s.installDiscordVersion(c, "0.0.2")
	s.Require().Truef(s.awaitLog(c, refreshLine, 30*time.Second), "no catalog refresh ran; daemon log:\n%s",
		s.readDaemonLog(c))
	s.Require().Truef(s.awaitLog(c, "configuration reloaded from", 30*time.Second), "the refresh did not reload")
	s.Require().Containsf(s.readAs(c, bin), "STOLEN|"+refreshMarker, "the new version was not admitted; daemon log:\n%s",
		s.readDaemonLog(c))
	_, conf := s.exec(c, []string{"cat", "/etc/app-listener/daemon.conf"})
	s.Require().Contains(conf, bin, "the refresh must persist the new version")

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// A non-writer's mkdir is allowed (version names are free) but planting the binary name is not,
// and nothing it does gets a refresh to admit anything.
func (s *IntegrationSuite) TestDaemon_CatalogRefresh_NonWriterPlantAdmitsNothing() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.startRefreshDaemon(c)
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/swap_reader"), swapReaderPath, 0o755))

	plant := discordDir + "/0.0.9/Discord"
	_, out := s.exec(c, []string{"sh", "-c", "mkdir -p " + discordDir + "/0.0.9 && cp " + swapReaderPath + " " + plant +
		"; cp " + swapReaderPath + " /tmp/Discord && mv /tmp/Discord " + plant + "; ls -l " + plant + " 2>&1"})
	s.T().Logf("plant attempts: %s", out)
	s.requireDenialLogged(c, "PLANT", "Discord")

	time.Sleep(15 * time.Second) // past the refresh debounce
	_, conf := s.exec(c, []string{"cat", "/etc/app-listener/daemon.conf"})
	s.Require().NotContainsf(conf, "0.0.9", "a non-writer's plant reached daemon.conf:\n%s", conf)
	_, out = s.exec(c, []string{"sh", "-c", plant + " " + refreshSecret + " 2>&1"})
	s.Require().NotContains(out, refreshMarker, "the plant must not read the resource")

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Filling the glob root with directories costs rescans, not a stalled daemon: the client's real
// update is still admitted.
func (s *IntegrationSuite) TestDaemon_CatalogRefresh_DirectoryFloodDoesNotStall() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.startRefreshDaemon(c)

	code, out := s.exec(c, []string{"sh", "-c", "cd " + discordDir + " && i=0; while [ $i -lt 5000 ]; do " +
		"mkdir junk$i/sub 2>/dev/null || mkdir -p junk$i/sub; i=$((i+1)); done; echo FLOODED"})
	s.Require().Truef(code == 0 && strings.Contains(out, "FLOODED"), "flood: %s", out)

	bin := s.installDiscordVersion(c, "0.0.3")
	deadline := time.Now().Add(60 * time.Second)
	for time.Now().Before(deadline) && !strings.Contains(s.readAs(c, bin), "STOLEN|"+refreshMarker) {
		time.Sleep(time.Second)
	}
	s.Require().Containsf(s.readAs(c, bin), "STOLEN|"+refreshMarker,
		"after a directory flood the update was not admitted within a minute; daemon log:\n%s", s.readDaemonLog(c))

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// A package-style replacement of a whitelisted system binary is admitted from the watch, well inside
// the 30s re-sync: the app launched right after reads its resource, and the version already running
// keeps reading.
func (s *IntegrationSuite) TestDaemon_CatalogWatch_SystemBinaryReplacedWithinDebounce() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	for _, f := range []string{"supersede_probe", "swap_reader"} {
		s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/"+f), "/exploits/"+f, 0o755), "copy "+f)
	}
	const app = "/usr/local/bin/app"
	s.exec(c, []string{"sh", "-c", "mkdir -p /protected /etc/app-listener " + supersedeCtl + " && printf '" +
		supersedeMarker + "' > /protected/secret && cp " + supersedeProbe + " " + app})
	s.startDaemon(c, `[watch /protected]
need_encryption: false
`+app)
	s.exec(c, []string{"sh", "-c", "nohup " + app + " hold /protected/secret " + supersedeCtl + " >/dev/null 2>&1 &"})
	s.Require().Contains(s.probeAsk(c, "read"), "STOLEN|"+supersedeMarker, "baseline")

	code, out := s.exec(c, []string{"sh", "-c", "cp " + swapReaderPath + " /usr/local/bin/.app.tmp && mv -f " +
		"/usr/local/bin/.app.tmp " + app + " 2>&1"})
	s.Require().Equalf(0, code, "the package-style upgrade: %s", out)
	s.Require().Truef(s.awaitLog(c, "re-synced", 10*time.Second),
		"the replacement was not re-synced within the watch's debounce; daemon log:\n%s", s.readDaemonLog(c))

	_, out = s.exec(c, []string{app, "/protected/secret"})
	s.Require().Containsf(out, "STOLEN|"+supersedeMarker, "the first launch after the upgrade must read: %s", out)
	s.Require().NotContains(s.readDaemonLog(c), "configuration reloaded", "admission must need no reload")
	s.Require().Contains(s.probeAsk(c, "read"), "STOLEN|"+supersedeMarker,
		"the version running before the upgrade must keep reading")

	s.exec(c, []string{"sh", "-c", "pkill -f 'bin/app hold'; pkill -f 'app-listener daemon' || true"})
}
