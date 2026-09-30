package integrationtests

import (
	"strings"
	"time"
)

// fscrypt.key and edit-auth.hash are sealed guards nested in the read-only /etc/app-listener guard.
// That guard's first sweep (selfSweepEvery, 30s after attach) rescans the directory; the rescan must
// not re-assign the sealed files to the read-only resource, which lets every process read.
func (s *IntegrationSuite) TestDaemon_SelfProtection_SweepKeepsKeySealed() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const marker = "SELF-KEY-SWEEP-MARKER-0123456789" // 32 bytes, like a real key
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /exploits /etc/app-listener && echo s > /protected/secret && printf '" +
			marker + "' > /etc/app-listener/fscrypt.key && chmod 600 /etc/app-listener/fscrypt.key"})
	s.copySwapFixtures(c)
	s.startDaemon(c, `[watch /protected]
need_encryption: false
/usr/bin/sleep`)

	selfReady := false
	for dl := time.Now().Add(45 * time.Second); time.Now().Before(dl); {
		if strings.Contains(s.readDaemonLog(c), "self-protection: guarding /etc/app-listener/edit-auth.hash") {
			selfReady = true
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
	s.Require().Truef(selfReady, "self-protection guards did not attach, log:\n%s", s.readDaemonLog(c))

	// Two sweep periods.
	for dl := time.Now().Add(75 * time.Second); time.Now().Before(dl); {
		_, out := s.exec(c, []string{"sh", "-c", swapReaderPath + " /etc/app-listener/fscrypt.key 2>&1"})
		s.Require().NotContainsf(out, "STOLEN|", "fscrypt.key became readable by a non-daemon process")
		code, out := s.exec(c, []string{"sh", "-c", "cat /etc/app-listener/edit-auth.hash >/dev/null 2>&1"})
		s.Require().NotEqualf(0, code, "edit-auth.hash became readable by a non-daemon process: %s", out)
		time.Sleep(3 * time.Second)
	}

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Nested guarded roots are refused at load: the engine records one owning resource per inode, so a
// file under both roots would follow only one of their whitelists.
func (s *IntegrationSuite) TestDaemon_NestedResources_ConfigRefused() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"sh", "-c", "mkdir -p /outer/inner /etc/app-listener && echo s > /outer/inner/secret"})
	s.exec(c, []string{"sh", "-c", `cat > /etc/app-listener/daemon.conf <<'EOF'
[watch /outer]
need_encryption: false
/usr/bin/cat

[watch /outer/inner]
need_encryption: false
/usr/bin/grep
EOF`})

	code, out := s.exec(c, []string{"sh", "-c",
		"timeout 60 /app-listener daemon --config /etc/app-listener/daemon.conf --headless 2>&1; echo rc=$?"})
	s.Require().Containsf(out, "nested guarded trees are not supported", "the daemon must name the nesting: %s", out)
	s.Require().NotContainsf(out, "rc=0", "the daemon must refuse a nested config (exit %d): %s", code, out)
	s.Require().NotContainsf(out, "rc=124", "the daemon must refuse at load, not run: %s", out)
	s.Require().NotContainsf(out, "guard started", "no guard may attach for a refused config: %s", out)
}
