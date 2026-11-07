package integrationtests

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

	// FS_IOC_SETFLAGS (_IOW('f', 2, long)) with FS_IMMUTABLE_FL; tmpfs takes it since 6.0.
	code, out := s.exec(c, []string{"sh", "-c", "mkdir -p /fw /etc/app-listener && mount -t tmpfs tmpfs /fw && " +
		"echo SECRET-file-root > /fw/key && perl -e 'open F, \"<\", \"/fw/key\" or die $!; " +
		"ioctl(F, 0x40086602, pack(\"l!\", 0x10)) or die \"setflags: $!\\n\"' 2>&1"})
	if code != 0 {
		s.T().Skipf("immutable flag unavailable: %s", out)
	}
	s.startDaemon(c, "[watch /fw/key]\nneed_encryption: false\n/usr/bin/perl\n")

	_, out = s.exec(c, []string{"perl", "-e", `unlink "/fw/key" and exit 0; print "unlink: $!\n"; exit 1`})
	s.Require().Containsf(out, "not permitted", "the whitelisted unlink must reach the kernel and fail: %s", out)

	_, out = s.exec(c, []string{"sh", "-c", "cat /fw/key 2>&1"})
	s.Require().NotContainsf(out, "SECRET-file-root", "the file root became readable after a failed unlink: %s\n%s",
		out, s.readDaemonLog(c))
}
