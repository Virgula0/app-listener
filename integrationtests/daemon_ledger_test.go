package integrationtests

import (
	"regexp"
	"strconv"
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

var readyHashRe = regexp.MustCompile(`hashed (\d+) binaries \(\d+ MiB\), reused (\d+) hashes`)

// One binary whitelisted by two resources is read once per run; the hash both guards share must
// not blind either verifier: an in-place rewrite still demotes it in both.
func (s *IntegrationSuite) TestDaemon_SharedHash_InPlaceRewriteDemotesEveryResource() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	s.copySwapFixtures(c)
	s.installFakeSystemctl(c)

	const bin = "/opt/shared/reader"
	secrets := []string{"/p1/secret", "/p2/secret"}
	// The sleep lets the copy's ctime settle (guard settleAge): a fresher file is never shared.
	code, out := s.exec(c, []string{"sh", "-c", "mkdir -p /p1 /p2 /opt/shared /etc/app-listener && for f in " +
		strings.Join(secrets, " ") + "; do printf '" + ledgerMarker + "' > $f; done && cp " + swapReaderPath +
		" " + bin + " && sleep 2 2>&1"})
	s.Require().Equalf(0, code, "setup: %s", out)
	config := "[watch /p1]\nneed_encryption: false\n" + bin + "\n\n[watch /p2]\nneed_encryption: false\n" + bin
	s.startDaemon(c, config)
	s.Require().Truef(s.awaitLog(c, "daemon: ready", 30*time.Second), "no ready line; daemon log:\n%s", s.readDaemonLog(c))

	m := readyHashRe.FindStringSubmatch(s.readDaemonLog(c))
	s.Require().NotNilf(m, "no hash counters on the ready line; daemon log:\n%s", s.readDaemonLog(c))
	reused, _ := strconv.Atoi(m[2])
	s.Require().GreaterOrEqualf(reused, 1, "the second resource re-read %s (hashed %s, reused %s)", bin, m[1], m[2])
	for _, sec := range secrets {
		_, out := s.exec(c, []string{bin, sec})
		s.Require().Containsf(out, "STOLEN|"+ledgerMarker, "baseline: %s reads %s: %s", bin, sec, out)
	}

	// Same inode, new content: the verifier (10s tick) must re-read it in both guards.
	code, out = s.exec(c, []string{"sh", "-c", "i1=$(stat -c %i " + bin + "); printf x >> " + bin +
		"; [ \"$i1\" = \"$(stat -c %i " + bin + ")\" ] && echo SAME-INODE"})
	s.Require().Containsf(out, "SAME-INODE", "fixture: rewrite must keep the inode (%d): %s", code, out)
	s.Require().Truef(s.awaitLogCount(c, "was modified in place", 2, 40*time.Second),
		"the rewrite was not demoted in both resources; daemon log:\n%s", s.readDaemonLog(c))
	for _, sec := range secrets {
		_, out := s.exec(c, []string{"sh", "-c", "timeout 10 " + bin + " " + sec + " 2>&1"})
		s.Require().NotContainsf(out, "STOLEN|"+ledgerMarker, "the rewritten binary still reads %s", sec)
	}

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

var (
	cachedRe = regexp.MustCompile(`reused \d+ hashes, (\d+) from the hash cache`)
	sha256Re = regexp.MustCompile(`[0-9a-f]{64}`)
)

// readyCached waits for this start's ready line and returns its hash-cache hit count.
func (s *IntegrationSuite) readyCached(c testcontainers.Container) int {
	s.Require().Truef(s.awaitLog(c, "daemon: ready", 60*time.Second), "no ready line; daemon log:\n%s", s.readDaemonLog(c))
	m := cachedRe.FindStringSubmatch(s.readDaemonLog(c))
	s.Require().NotNilf(m, "no hash-cache counter on the ready line; daemon log:\n%s", s.readDaemonLog(c))
	n, _ := strconv.Atoi(m[1])
	return n
}

// A binary on a local filesystem is read once across restarts (the hash cache in binaries.db),
// but a rewrite while the daemon is stopped, size and mtime restored, is hashed again: the ledger
// sees the new content and refuses it.
func (s *IntegrationSuite) TestDaemon_HashCache_RestartReusesAndStoppedRewriteMisses() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	s.copySwapFixtures(c)
	s.installFakeSystemctl(c)

	const bin, overlayBin = "/mnt/img/bin/app", "/opt/overlay/app"
	// -I 256: ext4 keeps a birth time only in large inodes. The sleep lets the copy's ctime settle.
	build := `
set -e
mkdir -p /protected /etc/app-listener /mnt/img && printf '` + ledgerMarker + `' > /protected/secret
dd if=/dev/zero of=/tmp/img bs=1M count=64 status=none
mkfs.ext4 -q -F -I 256 /tmp/img
for i in $(seq 0 15); do [ -e /dev/loop$i ] || mknod /dev/loop$i b 7 "$i"; done
mount -o loop /tmp/img /mnt/img
mkdir -p /mnt/img/bin /opt/overlay && cp ` + swapReaderPath + ` ` + bin + ` && cp ` + swapReaderPath + ` ` + overlayBin + `
sleep 2
echo MOUNTED`
	code, out := s.exec(c, []string{"sh", "-c", build})
	if code != 0 || !strings.Contains(out, "MOUNTED") {
		s.T().Skipf("ext4 loop mount unavailable here (exit %d): %s", code, out)
	}
	defer s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon'; sleep 1; umount /mnt/img 2>/dev/null; losetup -D 2>/dev/null; true"})

	config := "[watch /protected]\nneed_encryption: false\n" + bin + "\n" + overlayBin
	s.startDaemon(c, config)
	s.readyCached(c)
	_, out = s.exec(c, []string{bin, "/protected/secret"})
	s.Require().Containsf(out, "STOLEN|"+ledgerMarker, "baseline: %s reads: %s", bin, out)

	s.restartUnconfirmed(c, config, nil)
	// The container root is overlayfs: its binary is never vouched for across restarts.
	s.Require().Equalf(1, s.readyCached(c), "want %s from the hash cache and %s hashed again", bin, overlayBin)
	_, out = s.exec(c, []string{bin, "/protected/secret"})
	s.Require().Containsf(out, "STOLEN|"+ledgerMarker, "after a cached restart %s must still read: %s", bin, out)

	// Same inode, size and mtime; only ctime moves.
	var sum string
	s.restartUnconfirmed(c, config, func() {
		code, out := s.exec(c, []string{"sh", "-c", "set -e; B=" + bin + "; i1=$(stat -c %i $B); m=$(stat -c %y $B); " +
			"printf X | dd of=$B bs=1 seek=$(($(stat -c %s $B)-1)) conv=notrunc status=none; touch -d \"$m\" $B; " +
			"[ \"$i1\" = \"$(stat -c %i $B)\" ] && sha256sum $B"})
		s.Require().Equalf(0, code, "fixture: in-place rewrite while stopped: %s", out)
		sum = sha256Re.FindString(out) // docker exec frames can prefix the output
		s.Require().NotEmptyf(sum, "fixture: no sha256 in %q", out)
	})
	line := s.unconfirmedLine(c, bin)
	s.Require().NotEmptyf(line, "a rewrite while stopped was admitted from the hash cache; daemon log:\n%s", s.readDaemonLog(c))
	s.Require().Containsf(line, "sha256="+sum, "the refusal must carry the rewritten content's hash: %s", line)
	s.assertNeverStolen(c, bin, "/protected/secret", ledgerMarker)
}
