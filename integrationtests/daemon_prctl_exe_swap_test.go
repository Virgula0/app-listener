package integrationtests

import "strings"

const prctlExeSwap = "/exploits/prctl_exe_swap"

// TestDaemon_Bypass_PrctlExeFileSwap: the guard keys a caller on task->mm->exe_file. A non
// whitelisted process repoints its own exe_file at a whitelisted binary's inode via
// prctl(PR_SET_MM, ...) (CRIU's restore primitive, reachable unprivileged via an unshared user
// namespace) and then reads the guarded secret. Without a gate the guard reads the whitelisted
// inode and ALLOWS it — a full identity spoof (confirmed on 7.0.11: the swap lands, the secret
// leaks). The guard_task_prctl LSM hook denies PR_SET_MM_MAP / PR_SET_MM_EXE_FILE, so the swap
// fails and the exploit is judged on its real, non-whitelisted exe inode. The secret must never
// come back.
func (s *IntegrationSuite) TestDaemon_Bypass_PrctlExeFileSwap() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/prctl_exe_swap"), prctlExeSwap, 0o755),
		"copy prctl_exe_swap")
	s.copySwapFixtures(c)

	const marker = "TOP-SECRET-PRCTL-EXE-9F3C"
	s.exec(c, []string{"sh", "-c", "mkdir -p /protected /etc/app-listener /opt/app/bin" +
		" && printf '" + marker + "' > /protected/secret && cp " + swapReaderPath + " /opt/app/bin/reader"})

	// The exploit and the whitelisted binary must be distinct inodes or the swap is a no-op.
	s.requireDistinctExes(c, prctlExeSwap, "/opt/app/bin/reader")

	s.startDaemon(c, `[watch /protected]
need_encryption: false
/opt/app/bin/reader`)

	// Baseline: the whitelisted binary reads; the non-whitelisted exploit, not swapping, is denied.
	_, out := s.exec(c, []string{"/opt/app/bin/reader", "/protected/secret"})
	s.Require().Containsf(out, "STOLEN|"+marker, "baseline: the whitelisted binary reads its resource: %s", out)
	_, out = s.exec(c, []string{prctlExeSwap, "--noswap", "/opt/app/bin/reader", "/protected/secret"})
	s.Require().Containsf(out, "DENIED|", "baseline: the non-whitelisted exploit must be denied: %s", out)
	s.Require().NotContains(out, marker, "baseline: the non-whitelisted exploit read the secret")

	// Attack: swap exe_file to the whitelisted inode, then read. The secret must not leak.
	_, out = s.exec(c, []string{prctlExeSwap, "/opt/app/bin/reader", "/protected/secret"})
	s.T().Logf("prctl exe-swap bypass attempt:\n%s", out)
	s.Require().NotContainsf(out, marker,
		"an exe_file swap via prctl(PR_SET_MM) spoofed a whitelisted identity and read the secret:\n%s\ndaemon log:\n%s",
		out, s.readDaemonLog(c))
	// The swap itself may still land (the kernel allows it where exe_file and the mapped text differ,
	// e.g. overlayfs); what guard_task_prctl guarantees is that the swap cannot happen while guarded,
	// so the secret never leaks. If the SWAP line reports success the hook did not deny the prctl.
	if strings.Contains(out, "SWAP|PR_SET_MM") {
		s.T().Logf("NOTE: the exe_file swap was NOT denied by guard_task_prctl (SWAP reported ok); "+
			"it only did not leak because the guard judged the real exe inode. Investigate the hook.\n%s", out)
	}

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}
