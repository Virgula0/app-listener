package integrationtests

import (
	"regexp"
	"strings"
	"time"
)

// selfDirGuarded matches the read-only directory self guard, the last one attached.
var selfDirGuarded = regexp.MustCompile(`self-protection: guarding /etc/app-listener([^/]|$)`)

// fscrypt.key and binaries.db (the password hash) are sealed guards nested in the read-only
// /etc/app-listener guard.
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
		if selfDirGuarded.MatchString(s.readDaemonLog(c)) {
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
		code, out := s.exec(c, []string{"sh", "-c", "cat /etc/app-listener/binaries.db >/dev/null 2>&1"})
		s.Require().NotEqualf(0, code, "binaries.db became readable by a non-daemon process: %s", out)
		time.Sleep(3 * time.Second)
	}

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// systemd's ReadWritePaths=/etc/app-listener (units installed before it was dropped) bind-mounts the
// directory onto itself, so a sealed file's physical chain can't be verified past that mount root.
// The sealed guards must still own their files rather than leave them to the read-only resource.
func (s *IntegrationSuite) TestDaemon_SelfProtection_KeySealedUnderBindMount() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const marker = "SELF-KEY-BINDMNT-MARKER-01234567"
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /exploits /etc/app-listener && echo s > /protected/secret && printf '" +
			marker + "' > /etc/app-listener/fscrypt.key && chmod 600 /etc/app-listener/fscrypt.key"})
	code, out := s.exec(c, []string{"sh", "-c", "mount --bind /etc/app-listener /etc/app-listener 2>&1 && " +
		"grep -q ' /etc/app-listener /etc/app-listener ' /proc/self/mountinfo && echo bound"})
	s.Require().Truef(code == 0 && strings.Contains(out, "bound"), "bind mount failed: %s", out)
	s.copySwapFixtures(c)
	s.startDaemon(c, `[watch /protected]
need_encryption: false
/usr/bin/sleep`)

	selfReady := false
	for dl := time.Now().Add(45 * time.Second); time.Now().Before(dl); {
		if selfDirGuarded.MatchString(s.readDaemonLog(c)) {
			selfReady = true
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
	log := s.readDaemonLog(c)
	s.Require().Truef(selfReady, "self-protection guards did not attach, log:\n%s", log)
	s.Require().NotContainsf(log, "physical location could not be verified", "a sealed file was left to the "+
		"read-only resource, log:\n%s", log)

	_, out = s.exec(c, []string{"sh", "-c", swapReaderPath + " /etc/app-listener/fscrypt.key 2>&1"})
	s.Require().NotContainsf(out, "STOLEN|", "fscrypt.key is readable by a non-daemon process")
	code, out = s.exec(c, []string{"sh", "-c", "cat /etc/app-listener/fscrypt.key 2>&1"})
	s.Require().NotEqualf(0, code, "reading fscrypt.key must be denied: %s", out)
	s.Require().NotContainsf(out, marker, "fscrypt.key is readable by a non-daemon process")
	code, out = s.exec(c, []string{"sh", "-c", "cat /etc/app-listener/binaries.db >/dev/null 2>&1"})
	s.Require().NotEqualf(0, code, "binaries.db is readable by a non-daemon process: %s", out)

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

// fscrypt.key stays sealed through SIGHUP reloads: a non-whitelisted root reader races repeated
// reloads and must never see the key, since the self guards are never detached for a reload.
func (s *IntegrationSuite) TestDaemon_SelfProtection_KeySealedAcrossReloads() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const marker = "SELF-KEY-RELOAD-MARKER-012345678" // 32 bytes, like a real key
	const reloads = 5
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /etc/app-listener && echo s > /protected/secret && printf '" +
			marker + "' > /etc/app-listener/fscrypt.key && chmod 600 /etc/app-listener/fscrypt.key"})
	s.startDaemon(c, `[watch /protected]
need_encryption: false
/usr/bin/sleep`)

	selfReady := false
	for dl := time.Now().Add(45 * time.Second); time.Now().Before(dl); {
		if selfDirGuarded.MatchString(s.readDaemonLog(c)) {
			selfReady = true
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
	s.Require().Truef(selfReady, "self-protection guards did not attach, log:\n%s", s.readDaemonLog(c))

	s.exec(c, []string{"sh", "-c", "rm -f /tmp/racer.stop /tmp/racer.out /tmp/racer.n; nohup sh -c '" +
		"n=0; while [ ! -e /tmp/racer.stop ]; do cat /etc/app-listener/fscrypt.key 2>/dev/null; " +
		"n=$((n+1)); echo $n > /tmp/racer.n; done' > /tmp/racer.out 2>&1 &"})

	for i := 1; i <= reloads; i++ {
		s.sigDaemon(c, "HUP")
		done := false
		for dl := time.Now().Add(daemonShutdownTimeout); time.Now().Before(dl); {
			l := s.readDaemonLog(c)
			s.Require().NotContainsf(l, reloadRefused, "reload %d was refused:\n%s", i, l)
			if strings.Count(l, reloadEnd) >= i {
				done = true
				break
			}
			time.Sleep(100 * time.Millisecond)
		}
		s.Require().Truef(done, "reload %d did not finish, log:\n%s", i, s.readDaemonLog(c))
	}

	s.exec(c, []string{"sh", "-c", "touch /tmp/racer.stop; sleep 1"})
	_, n := s.exec(c, []string{"sh", "-c", "cat /tmp/racer.n"})
	s.Require().NotEqualf("", strings.TrimSpace(n), "the racer never ran")
	_, out := s.exec(c, []string{"sh", "-c", "cat /tmp/racer.out"})
	s.Require().NotContainsf(out, marker, "fscrypt.key was readable by a non-daemon process during a reload "+
		"(%s reads)", strings.TrimSpace(n))

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}
