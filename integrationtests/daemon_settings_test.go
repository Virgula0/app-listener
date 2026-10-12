package integrationtests

import (
	"strings"
	"time"
)

const settingsMarker = "TOP-SECRET-SETTINGS-41C7"

// max_inodes below the floor refuses to start (EX_CONFIG) before any guard loads; at the floor the
// map is sized to it, a tree filling more than 75% of it is warned, and enforcement is unchanged.
func (s *IntegrationSuite) TestDaemon_Settings_MaxInodesFloorAndOccupancyWarning() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	s.copySwapFixtures(c)
	s.installFakeSystemctl(c)

	const reader = "/opt/settings/reader"
	// 26,000 files: past 75% of the 32,768-row floor, within it.
	code, out := s.exec(c, []string{"sh", "-c", "mkdir -p /protected/many /etc/app-listener /opt/settings && printf '" +
		settingsMarker + "' > /protected/secret && cp " + swapReaderPath + " " + reader +
		" && cd /protected/many && seq 1 26000 | xargs touch 2>&1"})
	s.Require().Equalf(0, code, "setup: %s", out)
	watch := "[watch /protected]\nneed_encryption: false\n" + reader + "\n"

	s.exec(c, []string{"sh", "-c", "cat > /etc/app-listener/daemon.conf <<'EOF'\n[settings]\nmax_inodes: 1000\n\n" +
		watch + "EOF"})
	code, out = s.exec(c, []string{"sh", "-c", "timeout 120 /app-listener daemon --config /etc/app-listener/daemon.conf " +
		"--headless > /tmp/refused.log 2>&1; echo EXIT=$?; cat /tmp/refused.log"})
	s.Require().Containsf(out, "EXIT=78", "max_inodes below the floor must refuse to start (%d): %s", code, out)
	s.Require().Containsf(out, "max_inodes 1000 is outside", "%s", out)
	s.Require().NotContainsf(out, "guard started", "a guard attached before the config was refused: %s", out)

	config := "[settings]\nmax_inodes: 32768\n\n" + watch
	s.startDaemon(c, config)
	s.Require().Truef(s.awaitLog(c, "guard_inodes holds up to 32768 rows", 10*time.Second),
		"guard_inodes not sized by max_inodes; daemon log:\n%s", s.readDaemonLog(c))
	s.Require().Truef(s.awaitLog(c, "of 32768 rows (", 30*time.Second),
		"no occupancy warning past 75%%; daemon log:\n%s", s.readDaemonLog(c))

	_, out = s.exec(c, []string{reader, "/protected/secret"})
	s.Require().Containsf(out, "STOLEN|"+settingsMarker, "the whitelisted reader must read: %s", out)
	for _, f := range []string{"/protected/secret", "/protected/many/25999"} {
		code, out := s.exec(c, []string{"sh", "-c", "cat " + f + " 2>&1"})
		s.Require().NotEqualf(0, code, "cat read %s under a resized map: %s", f, out)
		s.Require().Truef(strings.Contains(out, "denied") || strings.Contains(out, "not permitted"),
			"cat %s: want a kernel denial, got %q", f, out)
	}

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}
