package integrationtests

import (
	"regexp"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
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

// narrowConfig: perl may do anything in /outer, but only grep (READ) and tee/mkdir (writers) inside
// its [narrow] subtrees. perl stands for gpg: whitelisted for the tree, never for the secret keys.
const narrowConfig = `[watch /outer]
need_encryption: false
/usr/bin/perl
/usr/bin/grep
/usr/bin/tee
/usr/bin/mkdir

[narrow /outer/keys]
/usr/bin/grep READ
/usr/bin/tee
/usr/bin/mkdir

[narrow /outer/empty]
/usr/bin/grep READ`

const narrowMarker = "NARROW-SECRET-MARKER-0123456789"

// perlOp runs one perl statement as the tree-only binary and returns its output; `$!` names the
// errno, so a denial reads "Operation not permitted".
func (s *IntegrationSuite) perlOp(c testcontainers.Container, stmt string) string {
	_, out := s.exec(c, []string{"perl", "-e", stmt})
	return out
}

func (s *IntegrationSuite) assertNarrowHolds(c testcontainers.Container, when string) {
	out := s.perlOp(c, `open(F,"<","/outer/pub") or die "open: $!\n"; print <F>`)
	s.Require().Containsf(out, "pub", "%s: perl must still read its tree: %s", when, out)
	_, out = s.exec(c, []string{"grep", "-r", "", "/outer/keys"})
	s.Require().Containsf(out, narrowMarker, "%s: the narrow's own binary must read it: %s", when, out)

	for name, stmt := range map[string]string{
		"read a key":         `open(F,"<","/outer/keys/secret") or die "open: $!\n"; print <F>`,
		"read a later file":  `open(F,"<","/outer/keys/d1/d2/d3/late") or die "open: $!\n"; print <F>`,
		"list the keys":      `opendir(D,"/outer/keys") or die "open: $!\n"; print join(",",readdir D)`,
		"stat the narrow":    `my @s = stat("/outer/keys") or die "stat: $!\n"; print "mtime=$s[9]\n"`,
		"create a key":       `open(F,">","/outer/keys/planted") or die "open: $!\n"; print "made\n"`,
		"move the narrow":    `rename("/outer/keys","/outer/k2") or die "rename: $!\n"; print "moved\n"`,
		"remove an empty":    `rmdir("/outer/empty") or die "rmdir: $!\n"; print "removed\n"`,
		"move a key out":     `rename("/outer/keys/secret","/outer/stolen") or die "rename: $!\n"; print "moved\n"`,
		"move a file in":     `rename("/outer/pub","/outer/keys/pub") or die "rename: $!\n"; print "moved\n"`,
		"unlink a key":       `unlink("/outer/keys/secret") or die "unlink: $!\n"; print "unlinked\n"`,
		"hardlink a key out": `link("/outer/keys/secret","/outer/linked") or die "link: $!\n"; print "linked\n"`,
	} {
		out := s.perlOp(c, stmt)
		s.Require().Containsf(out, "Operation not permitted", "%s: perl must not %s: %s", when, name, out)
		s.Require().NotContainsf(out, narrowMarker, "%s: %s leaked the key: %s", when, name, out)
	}
	_, out = s.exec(c, []string{"sh", "-c", "ls -d /outer/keys /outer/empty /outer/pub; cat /outer/keys/secret 2>&1"})
	s.Require().NotContainsf(out, "No such file", "%s: a denied operation went through: %s", when, out)
}

// A [narrow] subtree follows only its own whitelist: a binary whitelisted for the tree around it can't
// read, list, stat, create, move, remove or hard-link anything inside, files made after the start and
// deeper than the scan included, and a reload keeps the subtree's owner.
func (s *IntegrationSuite) TestDaemon_Narrow_SubtreeFollowsItsOwnWhitelist() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"sh", "-c", "mkdir -p /outer/keys /outer/empty /etc/app-listener && echo pub > /outer/pub && " +
		"echo " + narrowMarker + " > /outer/keys/secret"})
	s.startDaemon(c, narrowConfig)
	code, out := s.exec(c, []string{"sh", "-c", "mkdir -p /outer/keys/d1/d2/d3 && echo " + narrowMarker +
		" | tee /outer/keys/d1/d2/d3/late"})
	s.Require().Equalf(0, code, "the narrow's writers must write it: %s", out)
	s.assertNarrowHolds(c, "after start")

	s.sigDaemon(c, "HUP")
	const reloadDone = "configuration reloaded without dropping protection"
	done := false
	for deadline := time.Now().Add(daemonShutdownTimeout); time.Now().Before(deadline); {
		if strings.Contains(s.readDaemonLog(c), reloadDone) {
			done = true
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	s.Require().Truef(done, "reload did not complete, log:\n%s", s.readDaemonLog(c))
	s.assertNarrowHolds(c, "after a reload")
	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// A [narrow] granting what its parent doesn't, or naming a missing path, refuses to start: never a
// silent fallback to the parent's wider whitelist.
func (s *IntegrationSuite) TestDaemon_Narrow_ConfigRefused() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.exec(c, []string{"sh", "-c", "mkdir -p /outer/keys /etc/app-listener"})
	for name, narrow := range map[string]string{
		"wider":   "[narrow /outer/keys]\n/usr/bin/cat\n",
		"missing": "[narrow /outer/gone]\n/usr/bin/grep\n",
	} {
		s.exec(c, []string{"sh", "-c", "cat > /etc/app-listener/daemon.conf <<'EOF'\n[watch /outer]\n" +
			"need_encryption: false\n/usr/bin/grep READ\n\n" + narrow + "EOF"})
		_, out := s.exec(c, []string{"sh", "-c",
			"timeout 60 /app-listener daemon --config /etc/app-listener/daemon.conf --headless 2>&1; echo rc=$?"})
		s.Require().Containsf(out, "[narrow", "%s: the daemon must name the section: %s", name, out)
		s.Require().NotContainsf(out, "rc=0", "%s: must refuse: %s", name, out)
		s.Require().NotContainsf(out, "rc=124", "%s: must refuse at load, not run: %s", name, out)
	}
}
