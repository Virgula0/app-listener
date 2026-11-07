package integrationtests

import (
	"strings"
	"time"
)

// The guard's delete hooks run before the kernel's own checks, so a delete the guard allows can
// still fail and leave the inode in place: it must stay guarded. A whitelisted app's cleanup rmdir of
// its own non-empty root (gnome-keyring at exit) must not let anyone else create files in it.
func (s *IntegrationSuite) TestDaemon_FailedRmdirKeepsRootGuarded() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	code, out := s.exec(c, []string{"sh", "-c", "mkdir -p /watch /etc/app-listener && echo s > /watch/secret"})
	s.Require().Equalf(0, code, "setup: %s", out)
	s.startDaemon(c, "[watch /watch]\nneed_encryption: false\n/usr/bin/perl\n")

	plant := func(name string) string {
		_, out := s.exec(c, []string{"sh", "-c", "(echo x > /watch/" + name + ") 2>&1; echo rc=$?"})
		return out
	}
	s.Require().Containsf(plant("before"), "cannot create", "an unwhitelisted create must be denied")

	_, out = s.exec(c, []string{"perl", "-e", `rmdir "/watch" and exit 0; print "rmdir: $!\n"; exit 1`})
	s.Require().Containsf(out, "not empty", "the whitelisted rmdir must reach the kernel and fail: %s", out)

	out = plant("after")
	s.Require().Containsf(out, "cannot create", "the root must stay guarded after a failed rmdir: %s\n%s",
		out, s.readDaemonLog(c))
	code, _ = s.exec(c, []string{"cat", "/watch/secret"})
	s.Require().NotEqual(0, code, "the secret became readable")
}

// A failed unlink of a single-file root (immutable: EPERM from may_delete, after the guard's hook)
// must leave the file guarded.
func (s *IntegrationSuite) TestDaemon_FailedUnlinkKeepsFileRootGuarded() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	// FS_IOC_SETFLAGS (_IOW('f', 2, long)) with FS_IMMUTABLE_FL; tmpfs takes it since 6.0. perl's
	// ioctl needs a writable buffer: a bare pack() dies before the syscall.
	code, out := s.exec(c, []string{"sh", "-c", "mkdir -p /fw /etc/app-listener && mount -t tmpfs tmpfs /fw && " +
		"echo SECRET-file-root > /fw/key && perl -e 'open F, \"<\", \"/fw/key\" or die $!; " +
		"my $f = pack(\"l!\", 0x10); ioctl(F, 0x40086602, $f) or die \"setflags: $!\\n\"' 2>&1"})
	if code != 0 && strings.Contains(out, "Inappropriate ioctl") {
		s.T().Skipf("this kernel's tmpfs has no immutable flag: %s", out)
	}
	s.Require().Equalf(0, code, "setup: %s", out)
	s.startDaemon(c, "[watch /fw/key]\nneed_encryption: false\n/usr/bin/perl\n")

	_, out = s.exec(c, []string{"perl", "-e", `unlink "/fw/key" and exit 0; print "unlink: $!\n"; exit 1`})
	s.Require().Containsf(out, "not permitted", "the whitelisted unlink must reach the kernel and fail: %s", out)

	_, out = s.exec(c, []string{"sh", "-c", "cat /fw/key 2>&1"})
	s.Require().NotContainsf(out, "SECRET-file-root", "the file root became readable after a failed unlink: %s\n%s",
		out, s.readDaemonLog(c))
}

// guard_root_moved is written before the kernel's own checks too. After a failed delete of a root,
// moving its unguarded parent aside and making a fresh root at the path must not move the guard
// there: the real tree, now under the moved parent, would be left unguarded.
func (s *IntegrationSuite) TestDaemon_FailedDeleteThenAncestorSwapKeepsAnchor() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	code, out := s.exec(c, []string{"sh", "-c", "mkdir -p /cfg/app /fwm /etc/app-listener && " +
		"echo SECRET-dir-root > /cfg/app/secret && mount -t tmpfs tmpfs /fwm && mkdir /fwm/p && " +
		"echo SECRET-file-root > /fwm/p/key && perl -e 'open F, \"<\", \"/fwm/p/key\" or die $!; " +
		"my $f = pack(\"l!\", 0x10); ioctl(F, 0x40086602, $f) or die \"setflags: $!\\n\"' 2>&1"})
	if code != 0 && strings.Contains(out, "Inappropriate ioctl") {
		s.T().Skipf("this kernel's tmpfs has no immutable flag: %s", out)
	}
	s.Require().Equalf(0, code, "setup: %s", out)
	s.startDaemon(c, "[watch /cfg/app]\nneed_encryption: false\n/usr/bin/perl\n\n"+
		"[watch /fwm/p/key]\nneed_encryption: false\n/usr/bin/perl\n")

	// One command: the file root's 1s follow clears a stale record while the root is still in place.
	_, out = s.exec(c, []string{"sh", "-c", `perl -e 'rmdir "/cfg/app" or print "rmdir: $!\n"; ` +
		`unlink "/fwm/p/key" or print "unlink: $!\n"' && mv /cfg /cfg.x && mkdir -p /cfg/app && ` +
		`mv /fwm/p /fwm/p.x && mkdir /fwm/p && echo decoy > /fwm/p/key && echo SWAPPED 2>&1`})
	s.Require().Containsf(out, "rmdir: Directory not empty", "the whitelisted rmdir must reach the kernel and fail: %s", out)
	s.Require().Containsf(out, "unlink: Operation not permitted", "the whitelisted unlink must fail in the kernel: %s", out)
	s.Require().Containsf(out, "SWAPPED", "swapping the unguarded parents: %s", out)

	// One dir-root sweep (resyncSweepEvery, 30s) and many file-root follows (1s).
	deadline := time.Now().Add(45 * time.Second)
	for time.Now().Before(deadline) && strings.Count(s.readDaemonLog(c), "refusing to re-anchor") < 2 {
		time.Sleep(time.Second)
	}
	for _, f := range []string{"/cfg.x/app/secret", "/fwm/p.x/key"} {
		_, out = s.exec(c, []string{"sh", "-c", "cat " + f + " 2>&1"})
		s.Require().NotContainsf(out, "SECRET-", "%s became readable: the guard moved to the decoy: %s\n%s", f, out,
			s.readDaemonLog(c))
	}
}
