package integrationtests

import (
	"fmt"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

// TestDaemon_Bypass_HardlinkOutOfLockedVault: while the daemon is stopped no guard is attached, and
// fscrypt checks only an encrypted link DESTINATION, so the vault owner's processes can link a
// locked vault file to a name outside the tree. The guard confines a hit to names under the watch
// root, so once the next start provisions the key that outside name read the plaintext. The daemon
// must refuse to unlock a vault holding a file with a link outside it, and unlock again once the
// outside name is gone.
func (s *IntegrationSuite) TestDaemon_Bypass_HardlinkOutOfLockedVault() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.copyFscryptHarness(c)

	const marker = "HARDLINK-OUT-OF-VAULT-SECRET-5D2E"
	const mnt = "/mnt/fscrypt-dir"
	s.setupFscryptDirFilesystem(c, mnt) // skips the test if unsupported here
	defer s.exec(c, []string{"sh", "-c", "umount " + mnt + " 2>/dev/null; losetup -D 2>/dev/null; true"})

	vault := mnt + "/vault"
	drop := mnt + "/drop"
	stolen := drop + "/stolen"
	s.exec(c, []string{"sh", "-c", fmt.Sprintf(
		"mkdir -p %s %s /etc/app-listener && chmod 1777 %s && printf '%s' > %s/secret.txt && "+
			"chmod 755 %s && chmod 600 %s/secret.txt && "+
			"head -c 32 /dev/zero > /etc/app-listener/fscrypt.key && chmod 600 /etc/app-listener/fscrypt.key",
		vault, drop, drop, marker, vault, vault, vault)})
	s.harnessMigrate(c, vault)
	s.Require().True(s.harnessIsEncrypted(c, vault), "setup: the vault must be encrypted")
	// The vault owner (here nobody) is the attacker: protected_hardlinks lets an owner link its file.
	s.exec(c, []string{"sh", "-c", "chown 65534:65534 " + vault + "/*"})

	config := fmt.Sprintf("[watch %s]\nneed_encryption: true\n/usr/bin/grep", vault)
	s.startDaemon(c, config)
	_, out := s.exec(c, []string{"grep", "-c", marker, vault + "/secret.txt"})
	s.Require().Contains(out, "1", "baseline: the whitelisted reader sees the unlocked secret")
	s.sigDaemon(c, "TERM")
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not exit after SIGTERM")
	s.Require().True(s.harnessIsEncrypted(c, vault), "a graceful stop must re-lock the vault")

	// Attack, daemon stopped: link the locked file (its no-key name) out of the vault.
	code, out := s.exec(c, []string{"sh", "-c", nobodyRun + "sh -c 'ln \"$(ls -d " + vault + "/* | head -n1)\" " +
		stolen + "' && stat -c 'links=%h' " + stolen})
	s.Require().Equalf(0, code, "fixture: linking out of a locked vault must work with no guard attached: %s", out)
	s.Require().Contains(out, "links=2")

	s.launchDaemon(c)
	up := s.awaitDaemonUpOrDead(c, vault)
	_, out = s.exec(c, []string{"sh", "-c", nobodyRun + "cat " + stolen + " 2>&1"})
	s.T().Logf("daemon up=%v; outside-name read: %s", up, out)
	s.Require().NotContainsf(out, marker,
		"the outside hard link read the vault's plaintext after the daemon unlocked it; daemon log:\n%s",
		s.readDaemonLog(c))
	s.Require().False(up, "the daemon must refuse to unlock a vault with an outside hard link")
	s.Require().Contains(s.readDaemonLog(c), "hard link outside", "the refusal must say why")
	s.Require().True(s.harnessIsEncrypted(c, vault), "the refused vault stays locked")

	// Removing the outside name lets the next start unlock.
	s.exec(c, []string{"rm", "-f", stolen})
	s.launchDaemon(c)
	s.awaitDaemonUp(c, config)
	_, out = s.exec(c, []string{"grep", "-c", marker, vault + "/secret.txt"})
	s.Require().Contains(out, "1", "once the outside link is gone the vault unlocks")
	s.sigDaemon(c, "TERM")
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not exit after SIGTERM")
}

// awaitDaemonUpOrDead waits for the daemon either to start guarding resource (true) or to exit
// (false). Not any "guard started": the self-protection guards log theirs before a vault unlocks.
func (s *IntegrationSuite) awaitDaemonUpOrDead(c testcontainers.Container, resource string) bool {
	deadline := time.Now().Add(daemonShutdownTimeout)
	for time.Now().Before(deadline) {
		if strings.Contains(s.readDaemonLog(c), "guard started — guarding: "+resource) {
			return true
		}
		if code, _ := s.exec(c, []string{"sh", "-c", noLiveAppListenerProcs}); code == 0 {
			return false
		}
		time.Sleep(250 * time.Millisecond)
	}
	s.Require().Failf("daemon neither started nor exited", "log:\n%s", s.readDaemonLog(c))
	return false
}
