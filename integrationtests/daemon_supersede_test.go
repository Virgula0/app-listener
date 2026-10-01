package integrationtests

import (
	"regexp"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

const (
	supersedeProbe  = "/exploits/supersede_probe"
	supersedeCtl    = "/tmp/ctl"
	supersedeMarker = "TOP-SECRET-SUPERSEDE-71C4"
)

// probeAsk hands the running supersede_probe one command and returns its answer.
func (s *IntegrationSuite) probeAsk(c testcontainers.Container, cmd string) string {
	out := supersedeCtl + "/" + cmd + ".out"
	s.exec(c, []string{"sh", "-c", "rm -f " + out + " && touch " + supersedeCtl + "/" + cmd})
	got, ok := s.awaitFile(c, out, 20*time.Second)
	s.Require().Truef(ok, "the probe did not answer %q; daemon log:\n%s", cmd, s.readDaemonLog(c))
	return got
}

// startSupersedeProbe guards /protected for a root-owned /opt/app/bin/app (supersede_probe) and
// starts it holding the secret, as an app running before its package is upgraded.
func (s *IntegrationSuite) startSupersedeProbe(c testcontainers.Container) {
	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	for _, f := range []string{"supersede_probe", "swap_reader", "swap_benign"} {
		s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/"+f), "/exploits/"+f, 0o755), "copy "+f)
	}
	s.exec(c, []string{"sh", "-c", "mkdir -p /protected /etc/app-listener /opt/app/bin " + supersedeCtl +
		" && printf '" + supersedeMarker + "' > /protected/secret && cp " + supersedeProbe + " /opt/app/bin/app"})
	s.startDaemon(c, `[watch /protected]
need_encryption: false
/opt/app/bin/app`)
	s.exec(c, []string{"sh", "-c", "nohup /opt/app/bin/app hold /protected/secret " + supersedeCtl +
		" >/dev/null 2>&1 &"})
	s.Require().Contains(s.probeAsk(c, "read"), "STOLEN|"+supersedeMarker, "baseline: the whitelisted app reads its resource")
}

// The upgrade's rename over the binary supersedes the old inode the probe still runs.
func (s *IntegrationSuite) upgradeSupersedeProbe(c testcontainers.Container) {
	code, out := s.exec(c, []string{"sh", "-c", "cp " + swapBenignPath +
		" /opt/app/bin/.app.new && mv -f /opt/app/bin/.app.new /opt/app/bin/app 2>&1"})
	s.Require().Equalf(0, code, "root's package-style replacement must go through: %s", out)
}

// assertSupersedeSplit: the app started before the upgrade, and its fork, keep reading; the old
// image re-exec'd after it (held fd, /proc/self/exe) and a non-whitelisted exec do not.
func (s *IntegrationSuite) assertSupersedeSplit(c testcontainers.Container, when string) {
	s.Require().Containsf(s.probeAsk(c, "read"), "STOLEN|"+supersedeMarker,
		"%s: the app started before the upgrade lost its resource", when)
	s.Require().Containsf(s.probeAsk(c, "fork"), "STOLEN|"+supersedeMarker,
		"%s: a fork (no exec) of the running app lost its resource", when)
	for _, cmd := range []string{"reexec", "procexe", "other"} {
		got := s.probeAsk(c, cmd)
		s.Require().NotContainsf(got, supersedeMarker, "%s: %s after the upgrade read the resource: %s\ndaemon log:\n%s",
			when, cmd, got, s.readDaemonLog(c))
	}
}

func (s *IntegrationSuite) TestDaemon_Superseded_OldProcessKeepsAccessNewExecRefused() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.startSupersedeProbe(c)

	// Control: before the upgrade the re-exec is the whitelisted image like any other launch.
	s.Require().Contains(s.probeAsk(c, "reexec"), "STOLEN|"+supersedeMarker,
		"control: the re-exec reads while the binary is current")

	s.upgradeSupersedeProbe(c)
	s.assertSupersedeSplit(c, "after the upgrade")

	s.reloadAndAwait(c)
	s.Require().Contains(s.readDaemonLog(c), reloadEnd, "the reload must commit")
	s.assertSupersedeSplit(c, "after a reload")

	s.exec(c, []string{"sh", "-c", "pkill -f 'app/bin/app hold'; pkill -f 'app-listener daemon' || true"})
}

// A freed whitelisted inode's number goes to the next file created, any user's: the reused number
// must not inherit the whitelist, neither before the daemon prunes the old key nor after.
func (s *IntegrationSuite) TestDaemon_Bypass_ReusedInodeNumberInheritsWhitelist() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	s.copySwapFixtures(c)

	const marker = "TOP-SECRET-REUSED-INO-2D9A"
	build := `
set -e
mkdir -p /protected /etc/app-listener /mnt/img && printf '` + marker + `' > /protected/secret
dd if=/dev/zero of=/tmp/img bs=1M count=16 status=none
mkfs.ext4 -q -F /tmp/img
for i in $(seq 0 15); do [ -e /dev/loop$i ] || mknod /dev/loop$i b 7 "$i"; done
mount -o loop /tmp/img /mnt/img
mkdir -p /mnt/img/bin /mnt/img/drop && chmod 1777 /mnt/img/drop
cp ` + swapReaderPath + ` /mnt/img/bin/app
echo MOUNTED`
	code, out := s.exec(c, []string{"sh", "-c", build})
	if code != 0 || !strings.Contains(out, "MOUNTED") {
		s.T().Skipf("ext4 loop mount unavailable here (exit %d): %s", code, out)
	}
	defer s.exec(c, []string{"sh", "-c", "umount /mnt/img 2>/dev/null; losetup -D 2>/dev/null; true"})

	s.startDaemon(c, `[watch /protected]
need_encryption: false
/mnt/img/bin/app`)
	_, out = s.exec(c, []string{"/mnt/img/bin/app", "/protected/secret"})
	s.Require().Containsf(out, "STOLEN|"+marker, "baseline: the whitelisted binary reads its resource: %s", out)

	old := s.inodeOf(c, "/mnt/img/bin/app")
	code, out = s.exec(c, []string{"sh", "-c", "rm -f /mnt/img/bin/app && " + nobodyRun + "cp " + swapReaderPath +
		" /mnt/img/drop/evil && stat -c 'ino=%i' /mnt/img/drop/evil"})
	s.Require().Equalf(0, code, "uninstall + an unprivileged copy: %s", out)
	m := regexp.MustCompile(`ino=(\d+)`).FindStringSubmatch(out)
	s.Require().NotNil(m, out)
	s.Require().Equalf(old, m[1], "fixture: ext4 must hand the freed number to the next file (old %s)", old)

	for _, when := range []string{"before the prune", "after the prune"} {
		_, out = s.exec(c, []string{"sh", "-c", nobodyRun + "/mnt/img/drop/evil /protected/secret 2>&1"})
		s.Require().NotContainsf(out, marker, "%s: a file reusing a freed whitelisted inode number read the "+
			"resource: %s\ndaemon log:\n%s", when, out, s.readDaemonLog(c))
		if when == "before the prune" {
			s.Require().Truef(s.awaitLog(c, "forgot freed binary inode "+old, 45*time.Second),
				"the freed key was not pruned; daemon log:\n%s", s.readDaemonLog(c))
		}
	}

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}
