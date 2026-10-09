package integrationtests

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

const (
	swapReaderPath = "/exploits/swap_reader"
	swapBenignPath = "/exploits/swap_benign"
)

func (s *IntegrationSuite) copySwapFixtures(c testcontainers.Container) {
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/swap_benign"), swapBenignPath, 0755), "copy swap_benign")
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/swap_reader"), swapReaderPath, 0755), "copy swap_reader")
}

// assertNeverStolen runs bin against secret over the whole re-sync window (denial-driven and the
// periodic sweep) and fails if it ever prints the marker.
func (s *IntegrationSuite) assertNeverStolen(c testcontainers.Container, bin, secret, marker string) {
	var lastOut string
	deadline := time.Now().Add(50 * time.Second)
	for time.Now().Before(deadline) {
		_, out := s.exec(c, []string{"sh", "-c", "timeout 10 " + bin + " " + secret + " 2>&1"})
		lastOut = out
		s.Require().NotContainsf(out, "STOLEN|"+marker, "%s was admitted and read the guarded secret", bin)
		time.Sleep(2 * time.Second)
	}
	s.T().Logf("last run of %s: %q", bin, lastOut)
}

// inodeOf returns path's inode number, "" if it does not exist.
func (s *IntegrationSuite) inodeOf(c testcontainers.Container, path string) string {
	_, out := s.exec(c, []string{"sh", "-c", "stat -c 'ino=%i' " + shQuote(path) + " 2>/dev/null || true"})
	if m := regexp.MustCompile(`ino=(\d+)`).FindStringSubmatch(out); m != nil {
		return m[1]
	}
	return ""
}

// Vuln 1: protection #1 pins the binary's inode only, so moving its PARENT directory away and
// recreating it puts an attacker file at the whitelisted path, which ReSyncBinaries used to
// re-admit by path. A replacement not created by one of the binary's updaters must stay denied.
func (s *IntegrationSuite) TestDaemon_Bypass_ParentDirSwapReWhitelisted() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const marker = "TOP-SECRET-PARENT-SWAP-4C1E"
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /exploits /etc/app-listener /tmp/apps/bin && printf '" + marker +
			"' > /protected/secret && chmod 755 /protected && chmod 644 /protected/secret"})
	s.copySwapFixtures(c)
	s.exec(c, []string{"sh", "-c", "cp " + swapBenignPath + " /tmp/apps/bin/app && chmod 755 /tmp/apps/bin/app"})
	s.startDaemon(c, `[watch /protected]
need_encryption: false
/tmp/apps/bin/app`)

	code, out := s.exec(c, []string{"/tmp/apps/bin/app"})
	s.Require().Equalf(0, code, "baseline whitelisted binary should run: %s", out)
	s.Require().Contains(out, "BENIGN-APP-OK")

	// Every step runs as a non-whitelisted process; any of them may be refused.
	_, out = s.exec(c, []string{"sh", "-c",
		"mv /tmp/apps/bin /tmp/apps/old; mkdir -p /tmp/apps/bin && cp " + swapReaderPath +
			" /tmp/apps/bin/app && chmod 755 /tmp/apps/bin/app; ls -li /tmp/apps/bin 2>&1"})
	s.T().Logf("swap: %s", out)

	s.assertNeverStolen(c, "/tmp/apps/bin/app", "/protected/secret", marker)

	log := s.readDaemonLog(c)
	s.Require().Truef(strings.Contains(log, "TRUST DENIED") || strings.Contains(log, "not re-admitted"),
		"the swap must be refused or the replacement's re-admission reported, daemon log:\n%s", log)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Vuln 2: a binary whitelisted for one resource is not an updater of another resource's binary.
// Before per-owner updater bits, any TRUSTED_BINARY could unlink, rename over or rewrite any
// protected binary, and the replacement then followed the Vuln 1 re-admission path.
func (s *IntegrationSuite) TestDaemon_Bypass_CrossResourceUpdaterReplacesBinary() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const marker = "TOP-SECRET-CROSS-UPDATER-9A2F"
	// cp/mv copies keep their basenames: ubuntu's coreutils is one multi-call binary.
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /other /exploits /etc/app-listener /opt/other && printf '" + marker +
			"' > /protected/secret && chmod 755 /protected && chmod 644 /protected/secret && echo o > /other/f" +
			" && cp /usr/bin/cp /opt/other/cp && cp /usr/bin/mv /opt/other/mv"})
	s.copySwapFixtures(c)
	s.exec(c, []string{"sh", "-c", "cp " + swapBenignPath + " /tmp/app && chmod 755 /tmp/app"})
	s.startDaemon(c, `[watch /protected]
need_encryption: false
/tmp/app

[watch /other]
need_encryption: false
/opt/other/cp
/opt/other/mv`)

	code, out := s.exec(c, []string{"/tmp/app"})
	s.Require().Equalf(0, code, "baseline whitelisted binary should run: %s", out)
	before := s.inodeOf(c, "/tmp/app")
	s.Require().NotEmpty(before)

	for _, tc := range []struct {
		what string
		cmd  string
	}{
		{"unlink+create", "/opt/other/cp --remove-destination " + swapReaderPath + " /tmp/app"},
		{"rename over", "/opt/other/cp " + swapReaderPath + " /tmp/staged && /opt/other/mv -f /tmp/staged /tmp/app"},
		{"in-place rewrite", "/opt/other/cp " + swapReaderPath + " /tmp/app"},
	} {
		code, out = s.exec(c, []string{"sh", "-c", tc.cmd + " 2>&1"})
		s.Require().NotEqualf(0, code, "%s by another resource's binary must be refused: %s", tc.what, out)
		s.Require().Equalf(before, s.inodeOf(c, "/tmp/app"), "%s replaced the protected binary", tc.what)
		code, out = s.exec(c, []string{"cmp", "-s", swapBenignPath, "/tmp/app"})
		s.Require().Equalf(0, code, "%s altered the protected binary: %s", tc.what, out)
	}
	s.requireDenialLogged(c, "WRITE", "/tmp/app")

	s.assertNeverStolen(c, "/tmp/app", "/protected/secret", marker)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Bypass: planUpdaters excludes general tools by the whitelisted path's basename, but updater
// rights attach to the inode. Ubuntu's coreutils is one multi-call binary hard-linked under every
// applet name, so whitelisting a harmless applet (date) makes cp/dd/tee, the same inode, updaters
// of the resource's app binary: any process rewrites it through them, and the app keeps its inode.
func (s *IntegrationSuite) TestDaemon_Bypass_MultiCallAppletBecomesUpdater() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const (
		marker  = "TOP-SECRET-MULTICALL-UPDATER-2E7C"
		applets = "/usr/lib/cargo/bin/coreutils"
	)
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /exploits /etc/app-listener && printf '" + marker +
			"' > /protected/secret && chmod 755 /protected && chmod 644 /protected/secret"})
	code, out := s.exec(c, []string{"sh", "-c", "stat -Lc %i /usr/bin/date " + applets + "/cp 2>&1 | sort -u | wc -l"})
	if code != 0 || strings.TrimSpace(out) != "1" {
		s.T().Skipf("coreutils here is not one multi-call inode (date vs cp): %s", out)
	}
	s.copySwapFixtures(c)
	s.exec(c, []string{"sh", "-c", "cp " + swapBenignPath + " /tmp/app && chmod 755 /tmp/app"})
	s.startDaemon(c, `[watch /protected]
need_encryption: false
/tmp/app
/usr/bin/date`)

	code, out = s.exec(c, []string{"/tmp/app"})
	s.Require().Equalf(0, code, "baseline whitelisted binary should run: %s", out)
	before := s.inodeOf(c, "/tmp/app")

	// Each runs from a non-whitelisted shell; only the applet is the resource's.
	for _, tc := range []struct{ what, cmd string }{
		{"cp applet", applets + "/cp " + swapReaderPath + " /tmp/app"},
		{"dd applet", applets + "/dd if=" + swapReaderPath + " of=/tmp/app status=none"},
		{"tee applet", applets + "/tee /tmp/app < " + swapReaderPath + " > /dev/null"},
	} {
		code, out = s.exec(c, []string{"sh", "-c", tc.cmd + " 2>&1"})
		s.T().Logf("%s: rc=%d %s", tc.what, code, out)
		s.Require().Equalf(before, s.inodeOf(c, "/tmp/app"), "%s replaced the protected binary", tc.what)
		code, out = s.exec(c, []string{"cmp", "-s", swapBenignPath, "/tmp/app"})
		s.Require().Equalf(0, code,
			"%s rewrote the resource's app binary: whitelisting `date` made the multi-call coreutils inode "+
				"an updater: %s", tc.what, out)
	}

	s.assertNeverStolen(c, "/tmp/app", "/protected/secret", marker)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// The hard-link updater rule is waived for a user-writable lib_binary only: pressure-vessel
// hard-links its helpers into every var/tmp-XXXXXX and must delete those links. A root-owned
// hard-linked lib_binary (distro coreutils) stays no updater; the single-link control proves that
// refusal is the hard-link rule, not a missing owner bit.
func (s *IntegrationSuite) TestDaemon_LibBinary_HardLinkWaiverUserWritableOnly() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/rawfsop"), "/exploits/rawfsop", 0755))
	s.copySwapFixtures(c)
	code, out := s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /rt /tmp/rtw /usr/local/pv /etc/app-listener && echo s > /protected/secret" +
			" && cp /exploits/rawfsop /tmp/rtw/pv-tool && ln /tmp/rtw/pv-tool /tmp/rtw/pv-tool.link" +
			" && chown -R 65534 /tmp/rtw" +
			" && cp /exploits/rawfsop /usr/local/pv/pv-root && ln /usr/local/pv/pv-root /usr/local/pv/pv-root.link" +
			" && cp /exploits/rawfsop /usr/local/pv/pv-single" +
			" && cp " + swapBenignPath + " /tmp/app && chmod 755 /tmp/app"})
	s.Require().Equalf(0, code, "fixture setup: %s", out)
	s.startDaemon(c, `[watch /protected]
need_encryption: false
/tmp/app
lib_dir /rt
lib_binary /tmp/rtw/pv-tool
lib_binary /usr/local/pv/pv-root
lib_binary /usr/local/pv/pv-single`)

	before := s.inodeOf(c, "/tmp/app")
	code, out = s.exec(c, []string{"/usr/local/pv/pv-root", "unlink", "/tmp/app"})
	s.T().Logf("root-owned hard-linked lib_binary: rc=%d %s", code, out)
	s.Require().Equalf(before, s.inodeOf(c, "/tmp/app"),
		"a root-owned hard-linked lib_binary became an updater and removed the protected binary")
	s.requireDenialLogged(c, "WRITE", "/tmp/app")

	code, out = s.exec(c, []string{"/tmp/rtw/pv-tool", "unlink", "/tmp/rtw/pv-tool.link"})
	s.Require().Equalf(0, code, "user-writable hard-linked lib_binary must delete its own hard link: %s", out)
	s.Require().Empty(s.inodeOf(c, "/tmp/rtw/pv-tool.link"), "the hard link is still there")

	code, out = s.exec(c, []string{"/usr/local/pv/pv-single", "unlink", "/tmp/app"})
	s.Require().Equalf(0, code, "control: a single-link lib_binary of the same set must be an updater: %s", out)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// A lib_dir lets any process clear write bits of a non-root-owned regular file (Proton's
// steampipe_fixups restore) and nothing else: adding a bit, touching another bit, a directory, a
// root-owned file (dropping g+w would make it system-trusted), chown/utimes/xattrs, or a
// whitelist tree stay denied. Asserted on the resulting mode, not the probe's exit code.
func (s *IntegrationSuite) TestDaemon_LibDir_ChmodMayOnlyDropWriteBits() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	for _, p := range []string{"chmod", "chown", "utimes", "setxattr"} {
		s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/"+p), "/exploits/"+p, 0755), "copy "+p)
	}
	code, out := s.exec(c, []string{"sh", "-c", "mkdir -p /rt/sub /protected /etc/app-listener" +
		" && for f in drop addw addx suid dropr meta; do echo x > /rt/$f.so; done" +
		" && chmod 755 /rt/drop.so /rt/suid.so /rt/dropr.so /rt/meta.so && chmod 555 /rt/addw.so" +
		" && chmod 444 /rt/addx.so && chown -R 65534 /rt" +
		" && echo x > /rt/root.so && chown 0:1234 /rt/root.so && chmod 664 /rt/root.so" +
		" && echo s > /protected/secret && chown 65534 /protected/secret && chmod 644 /protected/secret"})
	s.Require().Equalf(0, code, "fixture setup: %s", out)
	s.startDaemon(c, `[libraries "Runtime"]
lib_dir /rt

[watch /protected]
need_encryption: false
/usr/bin/true`)

	modeOf := func(path string) string {
		_, out := s.exec(c, []string{"stat", "-c", "%a", path})
		return strings.TrimSpace(out)
	}
	const eperm = "Operation not permitted"
	for _, tc := range []struct{ what, path, mode, want string }{
		{"drop write bits", "/rt/drop.so", "555", "555"},
		{"add a write bit", "/rt/addw.so", "755", "555"},
		{"add exec bits", "/rt/addx.so", "555", "444"},
		{"drop write but add setuid", "/rt/suid.so", "4555", "755"},
		{"drop a read bit", "/rt/dropr.so", "311", "755"},
		{"drop write on a directory", "/rt/sub", "555", "755"},
		{"drop g+w on a root-owned file", "/rt/root.so", "644", "664"},
	} {
		_, out := s.exec(c, []string{"/exploits/chmod", tc.path, tc.mode})
		s.Require().Equalf(tc.want, modeOf(tc.path), "%s: chmod %s %s: %s", tc.what, tc.mode, tc.path, out)
		if tc.mode != tc.want {
			s.Require().Containsf(out, eperm, "%s: the refusal must be the guard's EPERM", tc.what)
		}
	}
	// stat itself is denied in a whitelist tree, so only the kernel's answer is observable.
	_, out = s.exec(c, []string{"/exploits/chmod", "/protected/secret", "444"})
	s.Require().Containsf(out, eperm, "drop write in a whitelist tree must stay denied: %s", out)

	_, before := s.exec(c, []string{"stat", "-c", "%u:%g %Y", "/rt/meta.so"})
	for _, p := range []string{"chown", "utimes", "setxattr"} {
		_, out := s.exec(c, []string{"/exploits/" + p, "/rt/meta.so"})
		s.Require().Containsf(out, eperm, "%s must stay denied in a lib_dir: %s", p, out)
	}
	_, after := s.exec(c, []string{"stat", "-c", "%u:%g %Y", "/rt/meta.so"})
	s.Require().Equal(before, after, "chown/utimes changed a lib_dir file")

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Vuln 3 (A): a read-only lib_dir is trusted for loading only by its own lib_binary writers. Its
// contents may predate the guard, so a library planted there must not load into an unrelated
// whitelisted binary.
func (s *IntegrationSuite) TestDaemon_LibDir_TrustedOnlyForItsWriters() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	// Owned by nobody: a root-owned file in a root-owned dir is an auto-trusted system library,
	// which would make both loads vacuous.
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /rt /tmp/rtw /exploits /etc/app-listener && echo s > /protected/secret" +
			" && cp /usr/bin/dash /tmp/rtw/dash"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/lib_probe.so"), "/rt/lib_probe.so", 0755),
		"copy lib_probe.so")
	s.exec(c, []string{"sh", "-c", "chown -R 65534 /rt"})
	s.startDaemon(c, `[libraries "Runtime"]
lib_dir /rt
lib_binary /tmp/rtw/dash

[watch /protected]
need_encryption: false
/usr/bin/bash`)

	loaded, out := s.preload(c, "/tmp/rtw/dash", "/rt/lib_probe.so")
	s.Require().Truef(loaded, "control: the lib_dir's writer must load its library: %s", out)

	loaded, out = s.preload(c, "/usr/bin/bash", "/rt/lib_probe.so")
	s.Require().Falsef(loaded, "a whitelisted non-writer must not load from a read-only lib_dir: %s", out)
	s.requireDenialLogged(c, "LIBLOAD", "lib_probe.so")

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

const (
	swapHome      = "/home/u"
	swapRoot      = swapHome + "/.config/gh"
	swapMarker    = "TOP-SECRET-GH-TOKEN-5C2E"
	reloadEnd     = "daemon: configuration reloaded from"
	reloadRefused = "reload failed, keeping previous configuration"
)

// reloadAndAwait sends SIGHUP and waits for the reload to commit (trust set included) or be refused.
func (s *IntegrationSuite) reloadAndAwait(c testcontainers.Container) {
	s.sigDaemon(c, "HUP")
	for dl := time.Now().Add(daemonShutdownTimeout); time.Now().Before(dl); {
		if l := s.readDaemonLog(c); strings.Contains(l, reloadEnd) || strings.Contains(l, reloadRefused) {
			return
		}
		time.Sleep(300 * time.Millisecond)
	}
	s.Failf("reload did not finish", "daemon log:\n%s", s.readDaemonLog(c))
}

// startSwappableRootDaemon guards swapRoot, whose ancestors belong to the unprivileged user, beside
// /protected (whitelisting bash); /tmp/w/cp is swapRoot's writer. Returns false when the user's
// ancestor swap was refused, i.e. the guard already blocks the attack's first step.
func (s *IntegrationSuite) startSwappableRootDaemon(c testcontainers.Container) bool {
	s.exec(c, []string{"sh", "-c",
		"mkdir -p " + swapRoot + " /protected /exploits /etc/app-listener /tmp/w && cp /usr/bin/cp /tmp/w/cp && " +
			"printf '" + swapMarker + "' > " + swapRoot +
			"/hosts.yml && echo s > /protected/secret && chown -R 65534:65534 " + swapHome +
			" && chmod 755 /protected " + swapHome + " " + swapHome + "/.config"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/lib_probe.so"), "/exploits/lib_probe.so", 0o755),
		"copy lib_probe.so")
	s.startDaemon(c, `[watch `+swapRoot+`]
need_encryption: false
/tmp/w/cp

[watch /protected]
need_encryption: false
/usr/bin/bash`)

	_, out := s.exec(c, []string{"sh", "-c", nobodyRun + "cat " + swapRoot + "/hosts.yml 2>&1"})
	s.Require().NotContainsf(out, swapMarker, "baseline: a non-whitelisted reader must be denied %s: %s", swapRoot, out)
	code, out := s.exec(c, []string{"sh", "-c", nobodyRun + "cp /exploits/lib_probe.so " + swapRoot + "/evil.so 2>&1"})
	s.Require().NotEqualf(0, code, "baseline: a non-whitelisted writer must be denied %s: %s", swapRoot, out)

	code, out = s.exec(c, []string{"sh", "-c", nobodyRun + "sh -c 'mv " + swapHome + "/.config " + swapHome +
		"/.config.real && mkdir -p " + swapRoot + " && cp /exploits/lib_probe.so " + swapRoot + "/evil.so' 2>&1"})
	if code != 0 {
		s.T().Logf("the user's swap of an ancestor of the watch root was refused: %s", out)
		return false
	}
	return true
}

// Bypass: a reload re-resolves a kept resource's root by path. After the user swaps an unguarded
// ancestor, the new guard anchors on the user's replacement dir, and the real tree (moved aside,
// still holding the secret) is off the new root's chain.
func (s *IntegrationSuite) TestDaemon_Bypass_ReloadReanchorsSwappedAncestor() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	if !s.startSwappableRootDaemon(c) {
		return
	}
	moved := swapHome + "/.config.real/gh/hosts.yml"
	_, out := s.exec(c, []string{"sh", "-c", nobodyRun + "cat " + moved + " 2>&1"})
	s.Require().NotContainsf(out, swapMarker, "control: the moved tree must stay guarded before any reload: %s", out)

	s.reloadAndAwait(c)
	_, out = s.exec(c, []string{"sh", "-c", nobodyRun + "cat " + moved + " 2>&1"})
	s.Require().NotContainsf(out, swapMarker,
		"after the reload the real tree became readable by a non-whitelisted process: %s\ndaemon log:\n%s",
		out, s.readDaemonLog(c))

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Bypass: the periodic sweep re-anchors a root whose path resolves to a new inode on the same
// filesystem, so the ancestor swap needs no reload at all, only one sweep tick (resyncSweepEvery).
func (s *IntegrationSuite) TestDaemon_Bypass_SweepReanchorsSwappedAncestor() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	if !s.startSwappableRootDaemon(c) {
		return
	}
	refused := s.awaitLog(c, "refusing to re-anchor the guard root", 45*time.Second)
	s.T().Logf("sweep refused the re-anchor: %v", refused)

	_, out := s.exec(c, []string{"sh", "-c", nobodyRun + "cat " + swapHome + "/.config.real/gh/hosts.yml 2>&1"})
	s.Require().NotContainsf(out, swapMarker,
		"after a sweep tick the real tree became readable by a non-whitelisted process: %s\ndaemon log:\n%s",
		out, s.readDaemonLog(c))

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Bypass: SetGuardedDirs re-stats each resource path at reload and marks it TRUSTED_DIR_ANY, which
// trusts every library below it, pre-existing ones included. After an ancestor swap that is the
// user's own dir, so a library the user left there preloads into any whitelisted binary.
func (s *IntegrationSuite) TestDaemon_Bypass_ReloadTrustsSwappedRootLibraries() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	if !s.startSwappableRootDaemon(c) {
		return
	}
	// Positive control: the moved real tree is still guarded, so its writer's library is trusted.
	good := swapHome + "/.config.real/gh/good.so"
	code, out := s.exec(c, []string{"/tmp/w/cp", "/exploits/lib_probe.so", good})
	s.Require().Equalf(0, code, "control: the root's whitelisted writer must plant a library: %s", out)
	loaded, out := s.preload(c, "/usr/bin/bash", good)
	s.T().Logf("control: library written by the root's writer inside the guarded root loaded=%v", loaded)

	lib := swapRoot + "/evil.so"
	loaded, out = s.preload(c, "/usr/bin/bash", lib)
	s.Require().Falsef(loaded, "control: the user's library must not load before the reload: %s", out)

	s.reloadAndAwait(c)
	loaded, out = s.preload(c, "/usr/bin/bash", lib)
	s.Require().Falsef(loaded,
		"after the reload a library the user placed under the swapped root loaded into a whitelisted binary: %s\n"+
			"daemon log:\n%s", out, s.readDaemonLog(c))

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Bypass: validateResources rejects nesting at config load, but NewGuard re-opens the root by path
// later. A read-only lib_dir re-pointed into another resource between the two becomes that
// subtree's innermost root and takes its guard_inodes rows, which read-only mode lets anyone read.
// The ci-only reload gate holds the window open so the test is deterministic.
func (s *IntegrationSuite) TestDaemon_Bypass_LibDirRepointedIntoResourceAtReload() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const (
		marker = "TOP-SECRET-LIBDIR-TAKEOVER-9E3B"
		libDir = swapHome + "/steam/linux64"
		gate   = "/tmp/reload-gate"
	)
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected/sub " + libDir + " /tmp/rtw /etc/app-listener && printf '" + marker +
			"' > /protected/sub/secret && chmod 755 /protected /protected/sub && chmod 644 /protected/sub/secret && " +
			"cp /usr/bin/dash /tmp/rtw/dash && chown -R 65534:65534 " + swapHome + " && chmod 755 " + swapHome})
	config := `[libraries "Runtime"]
lib_dir ` + libDir + `
lib_binary /tmp/rtw/dash

[watch /protected]
need_encryption: false
/usr/bin/true`
	s.exec(c, []string{"sh", "-c", fmt.Sprintf("cat > /etc/app-listener/daemon.conf <<'EOF'\n%s\nEOF", config)})
	code, out := s.exec(c, []string{"sh", "-c", "APPLISTENER_TEST_RELOAD_GATE=" + gate +
		" nohup /app-listener daemon --config /etc/app-listener/daemon.conf --headless > /tmp/daemon.log 2>&1 &"})
	s.Require().Equalf(0, code, "starting daemon: %s", out)
	s.awaitDaemonUp(c, config)

	_, out = s.exec(c, []string{"sh", "-c", nobodyRun + "cat /protected/sub/secret 2>&1"})
	s.Require().NotContainsf(out, marker, "baseline: a non-whitelisted reader must be denied /protected: %s", out)

	s.sigDaemon(c, "HUP")
	s.Require().Truef(s.awaitLog(c, "test: reload paused before guard build", daemonShutdownTimeout),
		"the reload never reached the gate (binary built without -tags ci?):\n%s", s.readDaemonLog(c))
	// Swap the unguarded ancestor (renaming the root itself is refused). Relative target:
	// path_symlink only matches absolute targets against guarded paths.
	code, out = s.exec(c, []string{"sh", "-c", nobodyRun + "sh -c 'cd " + swapHome + " && mv steam steam.real" +
		" && mkdir steam && ln -s ../../protected/sub steam/linux64' 2>&1"})
	s.exec(c, []string{"touch", gate})
	s.T().Logf("user's re-point of the lib_dir ancestor: rc=%d %s", code, out)
	for dl := time.Now().Add(daemonShutdownTimeout); time.Now().Before(dl); {
		if l := s.readDaemonLog(c); strings.Contains(l, reloadEnd) || strings.Contains(l, reloadRefused) {
			break
		}
		time.Sleep(300 * time.Millisecond)
	}

	_, out = s.exec(c, []string{"sh", "-c", nobodyRun + "cat /protected/sub/secret " + libDir + "/secret 2>&1"})
	s.Require().NotContainsf(out, marker,
		"after the reload /protected/sub became readable by a non-whitelisted process: %s\ndaemon log:\n%s",
		out, s.readDaemonLog(c))

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Vuln 3 (B): a globbed lib_dir (Steam's Proton*, SteamLinuxRuntime_*) becomes a trusted-library
// root for every future match at the next catalog refresh, so only Steam may create a match.
func (s *IntegrationSuite) TestDaemon_LibDirGlob_NonWriterCannotCreateMatch() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.startSteamGlobDaemon(c)

	for _, name := range []string{"Proton_x", "SteamLinuxRuntime_x"} {
		p := steamCommon + "/" + name
		s.assertNotPlanted(c, "mkdir "+name, p, []string{"mkdir", p})
		s.exec(c, []string{"sh", "-c", "mkdir -p /tmp/pre/" + name})
		s.assertNotPlanted(c, "move in "+name, p, []string{"mv", "/tmp/pre/" + name, p})
	}
	s.requireDenialLogged(c, "PLANT", "Proton_x")

	// dash has no mkdir builtin (a child mkdir runs as its own exe), so the control binds the
	// reserved name with a redirection, which the Steam client's own inode performs.
	for _, name := range []string{"Proton_ok", "SteamLinuxRuntime_ok"} {
		code, out := s.exec(c, []string{steamClient, "-c", ": > " + steamCommon + "/" + name})
		s.Require().Equalf(0, code, "Steam's own binary must create %s: %s", name, out)
	}

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Vuln 4: a lib_binary writer of a read-only tree may write code other whitelisted processes
// load, so a same-uid tracer must not steer it: ptrace ATTACH and a traced exec of the writer are
// refused, even though running it does not taint.
func (s *IntegrationSuite) TestDaemon_ReadOnlyWriter_NotPtraceable() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"sh", "-c",
		"mkdir -p /rt /exploits /etc/app-listener && echo 'TOP-SECRET-RT-WRITER' > /rt/marker && chmod 644 /rt/marker"})
	victim := absPath("./exploits/ptrace_race")
	s.Require().NoError(c.CopyFileToContainer(s.ctx, victim, "/exploits/race_victim", 0755), "copy race_victim")
	s.Require().NoError(c.CopyFileToContainer(s.ctx, victim, "/exploits/race_tracer", 0755), "copy race_tracer")
	s.startDaemon(c, `[libraries "Runtime"]
lib_dir /rt
lib_binary /exploits/race_victim`)

	code, out := s.exec(c, []string{"sh", "-c",
		"rm -f /tmp/race_addr* /tmp/race_go /tmp/race_done; " +
			"timeout 60 /exploits/race_tracer tracer /exploits/race_victim " +
			"/rt/marker /tmp/race_addr /tmp/race_go /tmp/race_done 2>&1"})
	s.T().Logf("tracer exit=%d out=%q", code, out)
	s.Require().NotEqualf(0, code, "tracing a read-only tree's writer must be refused: %s", out)
	s.Require().NotContains(out, "SECRET_DUMP|SUCCESS", "the tracer controlled the writer: %s", out)

	log := s.readDaemonLog(c)
	s.Require().Truef(regexp.MustCompile(`op=PTRACE .*mode=ATTACH|op=TRACED_EXEC`).MatchString(log),
		"the refused trace must be logged as PTRACE ATTACH or TRACED_EXEC, daemon log:\n%s", log)

	s.exec(c, []string{"sh", "-c", "pkill -f 'race_victim' ; pkill -f 'app-listener daemon' || true"})
}

// The per-resource guards judge a process by its exe inode; only the trust guard stops code injected
// into a whitelisted process or a whitelisted binary rewritten in place. Without it the daemon must
// refuse to start with EX_CONFIG (78, no systemd restart loop): before anything is unlocked when the
// guard can't load or attach, and re-sealing the vault when the post-unlock trusted set can't apply.
func (s *IntegrationSuite) TestDaemon_TrustGuardUnavailable_RefusesToStart() {
	for _, stage := range []string{"start", "after-unlock"} {
		s.Run(stage, func() {
			c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
			defer c.Terminate(s.ctx)
			s.copyFscryptHarness(c)

			const marker = "TRUST-GUARD-FAIL-SECRET-7B1C"
			const secretFile = "/protected/secret.txt"
			s.exec(c, []string{"sh", "-c",
				"mkdir -p /protected /etc/app-listener && head -c 32 /dev/zero > /etc/app-listener/fscrypt.key && chmod 600 /etc/app-listener/fscrypt.key"})
			s.resetFileVaultTarget(c, secretFile, marker)
			s.exec(c, []string{"chmod", "644", secretFile})
			s.harnessMigrate(c, secretFile)
			s.Require().True(s.harnessIsEncrypted(c, secretFile), "setup: the file must be sealed before the daemon starts")

			config := fmt.Sprintf("[watch %s]\nneed_encryption: true\n/usr/bin/grep", secretFile)
			s.exec(c, []string{"sh", "-c", fmt.Sprintf("cat > /etc/app-listener/daemon.conf <<'EOF'\n%s\nEOF", config)})

			_, out := s.exec(c, []string{"sh", "-c", "APPLISTENER_TEST_TRUST_FAIL=" + stage +
				" timeout 90 /app-listener daemon --config /etc/app-listener/daemon.conf --headless 2>&1; echo rc=$?"})
			s.Require().Containsf(out, "rc=78", "the daemon must refuse to start with EX_CONFIG: %s", out)
			s.Require().Containsf(out, "trust guard", "the refusal must name the trust guard: %s", out)
			if stage == "start" {
				// The self guards over /etc/app-listener attach first on purpose; the resource must not.
				for _, attached := range []string{"watching: /protected", "guarding: /protected"} {
					s.Require().NotContainsf(out, attached, "no resource guard may attach before the trust guard: %s", out)
				}
				s.Require().NotContainsf(out, "locked fscrypt directory", "nothing may be unlocked before the trust guard: %s", out)
			}

			s.Require().Truef(s.harnessIsEncrypted(c, secretFile), "the vault must be sealed after the refusal")
			s.assertFileVaultSealed(c, secretFile, marker, "trust guard refusal at "+stage)
		})
	}
}

const nobodyRun = "setpriv --reuid=65534 --regid=65534 --clear-groups "

// awaitLog polls the daemon log for needle.
func (s *IntegrationSuite) awaitLog(c testcontainers.Container, needle string, timeout time.Duration) bool {
	for dl := time.Now().Add(timeout); time.Now().Before(dl); {
		if strings.Contains(s.readDaemonLog(c), needle) {
			return true
		}
		time.Sleep(300 * time.Millisecond)
	}
	return false
}

// awaitFile polls until path exists and returns its content.
func (s *IntegrationSuite) awaitFile(c testcontainers.Container, path string, timeout time.Duration) (string, bool) {
	for dl := time.Now().Add(timeout); time.Now().Before(dl); {
		if code, out := s.exec(c, []string{"sh", "-c", "cat " + shQuote(path) + " 2>/dev/null"}); code == 0 {
			return out, true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return "", false
}

// Bypass: the catalog refresh writes each glob match into daemon.conf as `"<path>"` with no
// escaping, and a quoted path ends at its first `"`. A non-whitelisted process renames the
// wildcard directory of an existing install (Discord's ~/.config/discord/*/Discord) to a name
// holding `"` and newlines, and the unattended refresh writes the lines it carries into the
// section; a name that doesn't parse leaves a config the daemon refuses at its next start.
func (s *IntegrationSuite) TestDaemon_Bypass_CatalogRefreshInjectsThroughMatchedName() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.installFakeSystemctl(c)

	const (
		marker   = "TOP-SECRET-CONF-INJECT-6D0F"
		injected = "need_encryption: true"
	)
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /etc/app-listener /protected " + discordDir + "/sentry $(dirname " + discordClient + ")" +
			" && cp /usr/bin/dash " + discordClient + " && printf '" + marker + "' > /protected/secret" +
			" && chmod 755 /protected && chmod 644 /protected/secret" +
			" && head -c 32 /dev/zero > /etc/app-listener/fscrypt.key && chmod 600 /etc/app-listener/fscrypt.key"})
	s.startDaemon(c, `[watch `+discordDir+`/sentry]
need_encryption: false
`+discordClient+`

[watch /protected]
need_encryption: false
/usr/bin/true`)

	_, out := s.exec(c, []string{"sh", "-c", nobodyRun + "cat /protected/secret 2>&1"})
	s.Require().NotContainsf(out, marker, "baseline: a non-whitelisted reader must be denied /protected: %s", out)

	evil := discordDir + "/0.0.1\"\n" + injected + "\n#"
	code, attack := s.exec(c, []string{"sh", "-c", "mv " + shQuote(discordDir+"/0.0.1") + " " + shQuote(evil) + " 2>&1"})
	s.T().Logf("rename of the wildcard dir to a quote/newline name: rc=%d %s", code, attack)

	code, refresh := s.exec(c, []string{"/app-listener", "install", "--update-catalog-only", "--live", "--yes"})
	if code == 0 {
		s.awaitLog(c, "configuration reloaded", 60*time.Second)
	}
	_, conf := s.exec(c, []string{"cat", "/etc/app-listener/daemon.conf"})
	s.Require().NotContainsf(conf, "\n"+injected+"\n",
		"a directory name carried a directive into daemon.conf through the catalog refresh:\n"+
			"refresh (exit %d):\n%s\ndaemon.conf:\n%s", code, refresh, conf)

	// The config must still bring the daemon up: a refused start leaves every resource unguarded.
	s.sigDaemon(c, "TERM")
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not exit after SIGTERM")
	s.exec(c, []string{"rm", "-f", "/run/app-listener-daemon.pid"})
	s.launchDaemon(c)
	up := false
	for dl := time.Now().Add(daemonShutdownTimeout); time.Now().Before(dl) && !up; time.Sleep(500 * time.Millisecond) {
		code, _ = s.exec(c, []string{"test", "-f", "/run/app-listener-daemon.pid"})
		up = code == 0
	}
	_, out = s.exec(c, []string{"sh", "-c", nobodyRun + "cat /protected/secret 2>&1"})
	s.Require().NotContainsf(out, marker,
		"after a restart on the refreshed config /protected was readable (daemon up=%v): %s\ndaemon.conf:\n%s\n"+
			"daemon log:\n%s", up, out, conf, s.readDaemonLog(c))
	s.Require().Truef(up, "the daemon refused to start on the refreshed config:\n%s\ndaemon log:\n%s",
		conf, s.readDaemonLog(c))

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Bypass: reservation #3 checks only the name being bound below a glob root, and the catalog
// refresh globs THROUGH symlinks (filepath.Glob + os.Stat), then the daemon EvalSymlinks the match.
// A directory symlink with an unreserved name, pointing outside the root, gets an attacker binary
// whitelisted by the next unattended refresh. Two shapes: a wildcard component (Discord's
// ~/.config/discord/*/Discord) and a missing fixed component (Claude's ~/.local/bin/claude when
// ~/.local/bin does not exist yet).
func (s *IntegrationSuite) TestDaemon_Bypass_CatalogRefreshFollowsPlantedSymlink() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.installFakeSystemctl(c)

	const (
		discordMarker = "TOP-SECRET-DISCORD-TOKEN-3B8D"
		claudeMarker  = "TOP-SECRET-CLAUDE-TOKEN-8F21"
		discordSecret = discordDir + "/sentry/secret"
		claudeSecret  = "/root/.claude/secret"
	)
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/xres_move"), "/exploits/xres_move", 0o755),
		"copy xres_move")
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /exploits /etc/app-listener " + discordDir + "/sentry $(dirname " + discordClient + ") /root/.claude /root/.local" +
			" && rm -rf /root/.local/bin" +
			" && cp /usr/bin/dash " + discordClient + " && cp /usr/bin/true /usr/local/bin/claude" +
			" && printf '" + discordMarker + "' > " + discordSecret + " && printf '" + claudeMarker + "' > " + claudeSecret +
			" && chmod 644 " + discordSecret + " " + claudeSecret +
			" && head -c 32 /dev/zero > /etc/app-listener/fscrypt.key && chmod 600 /etc/app-listener/fscrypt.key"})

	s.startDaemon(c, `[watch `+discordDir+`/sentry]
need_encryption: false
`+discordClient+`

[watch /root/.claude]
need_encryption: false
/usr/local/bin/claude`)
	s.Require().Contains(s.readDaemonLog(c), "reserved glob name(s)", "trust guard #3 was not populated")

	// Controls: the reservation is live (a direct plant of the reserved name is refused).
	s.exec(c, []string{"mkdir", "-p", discordDir + "/direct"})
	s.assertNotPlanted(c, "direct Discord plant", discordDir+"/direct/Discord",
		[]string{"cp", "/exploits/xres_move", discordDir + "/direct/Discord"})
	s.exec(c, []string{"rmdir", discordDir + "/direct"})

	// Attack, as a non-whitelisted process: the binaries live outside every reserved root, and only
	// an unreserved directory name (a symlink) is created inside the user's tree. Not app-*: that is
	// Discord's lib_dir name, reserved.
	code, out := s.exec(c, []string{"sh", "-c",
		"mkdir -p /tmp/x /tmp/y && cp /exploits/xres_move /tmp/x/Discord && cp /exploits/xres_move /tmp/y/claude" +
			" && ln -s /tmp/x " + discordDir + "/0.0.999 && ln -s /tmp/y /root/.local/bin 2>&1"})
	s.Require().Equalf(0, code, "planting the symlinks: %s", out)

	for _, tc := range []struct{ bin, secret, marker string }{
		{"/tmp/x/Discord", discordSecret, discordMarker},
		{"/tmp/y/claude", claudeSecret, claudeMarker},
	} {
		_, out = s.exec(c, []string{tc.bin, "read", tc.secret})
		s.Require().NotContainsf(out, tc.marker, "baseline: %s must be denied before the refresh: %s", tc.bin, out)
	}

	// The unattended refresh the boot unit / package-manager hooks run.
	code, refresh := s.exec(c, []string{"/app-listener", "install", "--update-catalog-only", "--live", "--yes"})
	if code == 0 {
		s.awaitLog(c, "configuration reloaded", 60*time.Second)
	}
	_, conf := s.exec(c, []string{"cat", "/etc/app-listener/daemon.conf"})

	for _, tc := range []struct{ bin, secret, marker string }{
		{"/tmp/x/Discord", discordSecret, discordMarker},
		{"/tmp/y/claude", claudeSecret, claudeMarker},
	} {
		_, out = s.exec(c, []string{tc.bin, "read", tc.secret})
		s.Require().NotContainsf(out, tc.marker,
			"%s (outside every reserved root, reached through a planted directory symlink) was whitelisted "+
				"by the catalog refresh and read %s: %s\nrefresh (exit %d):\n%s\ndaemon.conf:\n%s",
			tc.bin, tc.secret, out, code, refresh, conf)
	}

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Bypass: is_system_trusted excludes only FUSE. A filesystem image the user builds keeps its
// on-disk ownership when mounted (udisks2 lets an active-session user loop-mount one, nosuid,nodev
// but exec allowed), so a "root-owned" library on it is auto-trusted and preloads into a whitelisted
// binary. The image is built unprivileged; only the mount runs as root, standing in for udisksd.
func (s *IntegrationSuite) TestDaemon_Bypass_UserImageRootOwnedLibNotAutoTrusted() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"sh", "-c", "mkdir -p /protected /exploits /etc/app-listener /realsys /tmp/nb && echo s > /protected/secret"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/lib_probe.so"), libProbe, 0o755), "copy lib_probe.so")

	build := `
set -e
rm -rf /tmp/u && mkdir -p /tmp/u/tree/lib
cp ` + libProbe + ` /tmp/u/tree/lib/evil.so
chmod 755 /tmp/u/tree/lib && chmod 644 /tmp/u/tree/lib/evil.so
chown -R 65534:65534 /tmp/u
cat > /tmp/u/build.sh <<'EOF'
set -e
cd /tmp/u
dd if=/dev/zero of=img bs=1M count=16 status=none
mkfs.ext4 -q -F -d tree img
for f in /lib /lib/evil.so; do
  debugfs -w -R "sif $f uid 0" img >/dev/null 2>&1
  debugfs -w -R "sif $f gid 0" img >/dev/null 2>&1
done
EOF
` + nobodyRun + `env PATH=/usr/sbin:/usr/bin:/sbin:/bin sh /tmp/u/build.sh
for i in $(seq 0 15); do [ -e /dev/loop$i ] || mknod /dev/loop$i b 7 "$i"; done
mkdir -p /media/u
LOOP=$(losetup -f --show /tmp/u/img)
mount -o nosuid,nodev "$LOOP" /media/u
printf 'LOOPDEV=%s\n' "$LOOP"
stat -f -c 'fstype=%T' /media/u
stat -c 'owner=%u:%g mode=%a' /media/u/lib /media/u/lib/evil.so
`
	code, out := s.exec(c, []string{"sh", "-c", build})
	if code != 0 || !strings.Contains(out, "LOOPDEV=/dev/loop") {
		s.T().Skipf("unprivileged image build / loop mount unavailable here (exit %d): %s", code, out)
	}
	defer s.exec(c, []string{"sh", "-c", "umount /media/u 2>/dev/null; losetup -D 2>/dev/null; true"})
	s.Require().NotContainsf(out, "fuse", "the image must be a real block filesystem, not FUSE: %s", out)
	s.Require().Containsf(out, "owner=0:0 mode=755", "the library's dir on the image must present as root-owned: %s", out)
	s.Require().Containsf(out, "owner=0:0 mode=644", "the library on the image must present as root-owned: %s", out)

	// Controls on the real fs: a genuine root-owned library loads, a user-owned one is refused.
	s.exec(c, []string{"sh", "-c",
		"cp " + libProbe + " /realsys/evil.so && chmod 644 /realsys/evil.so && " +
			"cp " + libProbe + " /tmp/nb/evil.so && chown -R 65534:65534 /tmp/nb && chmod 644 /tmp/nb/evil.so"})
	s.startDaemon(c, `[watch /protected]
need_encryption: false
/usr/bin/true`)

	_, out = s.exec(c, []string{"sh", "-c", "LD_PRELOAD=/realsys/evil.so /usr/bin/true 2>&1"})
	s.Require().Containsf(out, libProbeMarker, "control: a genuine root-owned system library must load: %s", out)
	_, out = s.exec(c, []string{"sh", "-c", "LD_PRELOAD=/tmp/nb/evil.so /usr/bin/true 2>&1"})
	s.Require().NotContainsf(out, libProbeMarker, "control: a user-owned library must be refused: %s", out)

	_, out = s.exec(c, []string{"sh", "-c", "LD_PRELOAD=/media/u/lib/evil.so /usr/bin/true 2>&1"})
	s.Require().NotContainsf(out, libProbeMarker,
		"a library on a user-built, loop-mounted ext4 image presenting as root-owned was auto-trusted and "+
			"loaded into a whitelisted binary: %s", out)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Bypass: trust_mmap only gates NEW mappings. A process of a not-yet-whitelisted binary preloads
// attacker code, waits, and reads the secret once a reload whitelists that exe.
func (s *IntegrationSuite) TestDaemon_Bypass_PreloadedBeforeWhitelistReadsAfterReload() {
	s.preloadedBeforeWhitelist(false)
}

// Same, but the preloaded process forks after the reload and the child, sharing the payload's
// mappings without an mmap of its own, does the read.
func (s *IntegrationSuite) TestDaemon_Bypass_PreloadedBeforeWhitelistForkedChildReads() {
	s.preloadedBeforeWhitelist(true)
}

func (s *IntegrationSuite) preloadedBeforeWhitelist(fork bool) {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const marker = "TOP-SECRET-PRELOAD-HELD-7D42"
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected && printf '" + marker + "' > /protected/secret && chmod 755 /protected && chmod 644 /protected/secret" +
			" && rm -f /tmp/go* /tmp/leak* /tmp/armed*"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/preload_wait.so"), "/tmp/wait.so", 0o755),
		"copy preload_wait.so")

	// Outside /etc/app-listener so the test can rewrite it (the daemon self-guards that dir).
	const cfgPath = "/tmp/poc-daemon.conf"
	writeCfg := func(body string) {
		code, out := s.exec(c, []string{"sh", "-c", fmt.Sprintf("cat > %s <<'EOF'\n%s\nEOF", cfgPath, body)})
		s.Require().Equalf(0, code, "writing %s: %s", cfgPath, out)
	}
	writeCfg("[watch /protected]\nneed_encryption: false\n/usr/bin/true")
	code, out := s.exec(c, []string{"sh", "-c",
		"nohup /app-listener daemon --config " + cfgPath + " --headless > /tmp/daemon.log 2>&1 &"})
	s.Require().Equalf(0, code, "starting daemon: %s", out)
	s.awaitDaemonUp(c, "[watch /protected]")

	forkEnv := ""
	if fork {
		forkEnv = "LEAK_FORK=1 "
	}
	preload := func(tag string) string {
		return forkEnv + "LD_PRELOAD=/tmp/wait.so LEAK_FILE=/protected/secret LEAK_GO=/tmp/go-" + tag +
			" LEAK_OUT=/tmp/leak-" + tag + " LEAK_ARMED=/tmp/armed-" + tag + " /usr/bin/grep x /dev/null"
	}

	// Control: grep is not whitelisted yet, so the preload maps but its read is denied.
	s.exec(c, []string{"sh", "-c", "touch /tmp/go-ctl && " + preload("ctl") + " >/dev/null 2>&1; true"})
	res, ok := s.awaitFile(c, "/tmp/leak-ctl", 15*time.Second)
	s.Require().Truef(ok, "control: the preload never ran in the non-whitelisted grep")
	s.Require().NotContainsf(res, marker, "control: non-whitelisted grep must be denied the secret: %s", res)

	// Attack: arm a grep process with the payload while grep is still untrusted and unwhitelisted.
	s.exec(c, []string{"sh", "-c", "nohup sh -c '" + preload("held") + "' >/dev/null 2>&1 &"})
	_, ok = s.awaitFile(c, "/tmp/armed-held", 15*time.Second)
	s.Require().Truef(ok, "the held preload never armed")

	// A reload (a catalog refresh picking up an installed app) whitelists grep.
	writeCfg("[watch /protected]\nneed_encryption: false\n/usr/bin/true\n/usr/bin/grep")
	s.sigDaemon(c, "HUP")
	s.Require().Truef(s.awaitLog(c, "configuration reloaded without dropping protection", daemonShutdownTimeout),
		"reload did not complete, daemon log:\n%s", s.readDaemonLog(c))
	s.Require().Truef(s.awaitLog(c, "trust guard: trusted set rebuilt after reload", 10*time.Second),
		"the trust guard did not rebuild after the reload, daemon log:\n%s", s.readDaemonLog(c))

	// Control: grep is now whitelisted and a NEW preload into it is refused by the trust guard.
	_, out = s.exec(c, []string{"sh", "-c", "grep -c " + marker + " /protected/secret"})
	s.Require().Containsf(out, "1", "control: grep must be whitelisted after the reload: %s", out)
	s.exec(c, []string{"sh", "-c", "touch /tmp/go-new && " + preload("new") + " >/dev/null 2>&1; true"})
	_, armedNew := s.awaitFile(c, "/tmp/armed-new", 3*time.Second)
	s.Require().Falsef(armedNew, "control: a fresh LD_PRELOAD into the now-trusted grep must be refused")

	// The held process, preloaded before grep was trusted, now runs as a whitelisted exe.
	s.exec(c, []string{"touch", "/tmp/go-held"})
	res, ok = s.awaitFile(c, "/tmp/leak-held", 15*time.Second)
	s.Require().Truef(ok, "the held preload never reported")
	s.Require().NotContainsf(res, marker,
		"code preloaded into grep BEFORE a reload whitelisted it read the secret afterwards (fork=%v): %s",
		fork, res)
	s.Require().Truef(s.awaitLog(c, "op=PRELOADED", 10*time.Second),
		"the refused read was not logged, daemon log:\n%s", s.readDaemonLog(c))

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// inNosuidNamespace runs cmd in a fresh mount namespace whose / is remounted nosuid, as systemd
// does for the daemon's unit and bwrap/Flatpak do for sandboxed apps. It prints ROOTOPTS= first so
// the caller can check the namespace really is nosuid.
func inNosuidNamespace(cmd string) string {
	return "unshare -m sh -c " + shQuote("mount -o remount,bind,nosuid / && "+
		"awk '$5==\"/\"{print \"ROOTOPTS=\" $6}' /proc/self/mountinfo && "+cmd)
}

// A nosuid mount says nothing about who mounted the superblock under it: the daemon's own unit
// namespace and every sandbox remount system mounts nosuid. The daemon's library closure and the
// kernel's auto-trust must both judge pid 1's mounts, or libc is refused to a sandboxed app.
func (s *IntegrationSuite) TestDaemon_NosuidNamespace_SystemLibsStillTrusted() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"sh", "-c", "mkdir -p /protected /exploits /etc/app-listener /realsys /tmp/nb && echo s > /protected/secret"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/lib_probe.so"), libProbe, 0o755), "copy lib_probe.so")
	s.exec(c, []string{"sh", "-c",
		"cp " + libProbe + " /realsys/evil.so && chmod 644 /realsys/evil.so && " +
			"cp " + libProbe + " /tmp/nb/evil.so && chown -R 65534:65534 /tmp/nb && chmod 644 /tmp/nb/evil.so"})

	const config = "[watch /protected]\nneed_encryption: false\n/usr/bin/true"
	s.exec(c, []string{"sh", "-c", "cat > /etc/app-listener/daemon.conf <<'EOF'\n" + config + "\nEOF"})
	code, out := s.exec(c, []string{"sh", "-c", "nohup " + inNosuidNamespace(
		"exec /app-listener daemon --config /etc/app-listener/daemon.conf --headless") + " > /tmp/daemon.log 2>&1 &"})
	s.Require().Equalf(0, code, "starting daemon: %s", out)
	s.awaitDaemonUp(c, config)
	defer s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})

	log := s.readDaemonLog(c)
	s.Require().Regexpf(`ROOTOPTS=\S*nosuid`, log, "the daemon's namespace must have a nosuid /: %s", log)
	s.Require().NotContainsf(log, "is on a nosuid mount",
		"the library closure judged the daemon's own nosuid namespace instead of pid 1's:\n%s", log)

	_, out = s.exec(c, []string{"sh", "-c", "LD_PRELOAD=/realsys/evil.so /usr/bin/true 2>&1"})
	s.Require().Containsf(out, libProbeMarker, "control: a root-owned system library must load in pid 1's namespace: %s", out)

	_, out = s.exec(c, []string{"sh", "-c", inNosuidNamespace("LD_PRELOAD=/realsys/evil.so /usr/bin/true 2>&1")})
	s.Require().Regexpf(`ROOTOPTS=\S*nosuid`, out, "the sandbox's / must be nosuid: %s", out)
	s.Require().Containsf(out, libProbeMarker,
		"a root-owned system library was refused to a whitelisted binary in a nosuid sandbox: %s", out)
	s.Require().NotContainsf(s.readDaemonLog(c), "path=/realsys/evil.so",
		"trust_mmap refused the system library in the sandbox")

	_, out = s.exec(c, []string{"sh", "-c", inNosuidNamespace("LD_PRELOAD=/tmp/nb/evil.so /usr/bin/true 2>&1")})
	s.Require().NotContainsf(out, libProbeMarker, "control: a user-owned library must still be refused in the sandbox: %s", out)
	s.requireDenialLogged(c, "LIBLOAD", "/tmp/nb/evil.so")
}

// Bypass: vouching is keyed by device number, and loop minors are reused. Once root's non-nosuid
// image is unmounted, a user image udisks2 attaches to the same loop device (nosuid) must not
// inherit its trust, however soon after the unmount it is used.
func (s *IntegrationSuite) TestDaemon_Bypass_ReusedLoopMinorNotVouched() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"sh", "-c", "mkdir -p /protected /exploits /etc/app-listener /mnt/r /media/u && echo s > /protected/secret"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/lib_probe.so"), libProbe, 0o755), "copy lib_probe.so")

	build := `
set -e
rm -rf /tmp/r /tmp/u && mkdir -p /tmp/r/tree/lib /tmp/u/tree/lib
cp ` + libProbe + ` /tmp/r/tree/lib/evil.so && chmod 755 /tmp/r/tree/lib && chmod 644 /tmp/r/tree/lib/evil.so
dd if=/dev/zero of=/tmp/r/img bs=1M count=16 status=none
mkfs.ext4 -q -F -d /tmp/r/tree /tmp/r/img
cp ` + libProbe + ` /tmp/u/tree/lib/evil.so && chmod 755 /tmp/u/tree/lib && chmod 644 /tmp/u/tree/lib/evil.so
chown -R 65534:65534 /tmp/u
cat > /tmp/u/build.sh <<'EOF'
set -e
cd /tmp/u
dd if=/dev/zero of=img bs=1M count=16 status=none
mkfs.ext4 -q -F -d tree img
for f in /lib /lib/evil.so; do
  debugfs -w -R "sif $f uid 0" img >/dev/null 2>&1
  debugfs -w -R "sif $f gid 0" img >/dev/null 2>&1
done
EOF
` + nobodyRun + `env PATH=/usr/sbin:/usr/bin:/sbin:/bin sh /tmp/u/build.sh
for i in $(seq 0 15); do [ -e /dev/loop$i ] || mknod /dev/loop$i b 7 "$i"; done
LOOP=$(losetup -f --show /tmp/r/img)
mount "$LOOP" /mnt/r
printf 'LOOPDEV=%s\n' "$LOOP"
printf 'ROOTDEV=%s\n' "$(mountpoint -d /mnt/r)"
`
	code, out := s.exec(c, []string{"sh", "-c", build})
	if code != 0 || !strings.Contains(out, "LOOPDEV=/dev/loop") {
		s.T().Skipf("image build / loop mount unavailable here (exit %d): %s", code, out)
	}
	defer s.exec(c, []string{"sh", "-c", "umount /mnt/r /media/u 2>/dev/null; losetup -D 2>/dev/null; true"})
	loop := regexp.MustCompile(`LOOPDEV=(\S+)`).FindStringSubmatch(out)[1]
	rootDev := regexp.MustCompile(`ROOTDEV=(\S+)`).FindStringSubmatch(out)[1]

	s.startDaemon(c, `[watch /protected]
need_encryption: false
/usr/bin/true`)
	defer s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})

	_, out = s.exec(c, []string{"sh", "-c", "LD_PRELOAD=/mnt/r/lib/evil.so /usr/bin/true 2>&1"})
	s.Require().Containsf(out, libProbeMarker, "control: root's non-nosuid loop image must be vouched: %s", out)

	// One shell, no pause: the user image must be judged before any mountinfo re-sync could run.
	swap := `
umount /mnt/r && losetup -d ` + loop + ` || exit 1
for i in 1 2 3 4 5 6 7 8 9 10; do losetup ` + loop + ` /tmp/u/img 2>/dev/null && break; sleep 0.05; done
mount -o nosuid,nodev ` + loop + ` /media/u || exit 1
printf 'USERDEV=%s\n' "$(mountpoint -d /media/u)"
stat -c 'owner=%u:%g mode=%a' /media/u/lib/evil.so
LD_PRELOAD=/media/u/lib/evil.so /usr/bin/true 2>&1
`
	code, out = s.exec(c, []string{"sh", "-c", swap})
	s.Require().Containsf(out, "USERDEV="+rootDev+"\n", "the user image must reuse root's loop device %s (exit %d): %s", rootDev, code, out)
	s.Require().Containsf(out, "owner=0:0 mode=644", "the user image's library must present as root-owned: %s", out)
	s.Require().NotContainsf(out, libProbeMarker,
		"a user image on a reused loop minor inherited the unmounted root image's trust: %s", out)
	// Trust events print the path below the superblock's root, not the mount point.
	s.requireDenialLogged(c, "LIBLOAD", "path=/lib/evil.so")
}

// Bypass: the .bun-* reservation turns on for ANY existing Bun tmp dir. A same-user process that
// creates it while the user has none (redirect declined) plants a .bun-* library that the next
// reload trusts in the Bun app: glob_lib_trusted judges name and root, not who wrote the file
// before the reservation existed. Owned by nobody like a user's cache (root-owned = system lib).
func (s *IntegrationSuite) TestDaemon_Bypass_BunTmpdirPlantedBeforeReservation() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /exploits /etc/app-listener /root/.config/opencode /root/bin /root/.cache" +
			" && cp /usr/bin/dash /root/bin/opencode"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/lib_probe.so"), libProbe, 0o755),
		"copy lib_probe.so")
	s.startDaemon(c, bunGlobConfig)

	plant := bunTmpDir + "/.bun-1000-evil.so"
	code, out := s.exec(c, []string{"sh", "-c", "mkdir -p " + bunTmpDir + "/sub && cp " + libProbe + " " + plant +
		" && cp " + libProbe + " " + bunTmpDir + "/sub/.bun-1000-deep.so && chown -R 65534 /root/.cache/app-listener"})
	s.Require().Equalf(0, code, "planting while nothing is reserved: %s", out)
	loaded, out := s.preload(c, "/root/bin/opencode", plant)
	s.Require().Falsef(loaded, "baseline: an unreserved plant must not load: %s", out)

	s.sigDaemon(c, "HUP")
	s.Require().True(s.awaitLog(c, "trusted set rebuilt after reload", 60*time.Second),
		"reload did not complete:\n%s", s.readDaemonLog(c))

	for _, p := range []string{plant, bunTmpDir + "/sub/.bun-1000-deep.so"} {
		loaded, out = s.preload(c, "/root/bin/opencode", p)
		s.Require().Falsef(loaded, "%s, planted before its dir was reserved, loaded into the whitelisted Bun app: %s",
			p, out)
	}

	// Control: the app's own extraction, written after the reservation, still loads.
	own := bunTmpDir + "/.bun-1000-own.so"
	code, out = s.exec(c, []string{"/root/bin/opencode", "-c", "cat " + libProbe + " > " + own})
	s.Require().Equalf(0, code, "the Bun app's own binary must extract %s: %s", own, out)
	loaded, out = s.preload(c, "/root/bin/opencode", own)
	s.Require().Truef(loaded, "control: the Bun app's own extraction must load: %s", out)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Bypass: the catalog refresh adopts every wildcard-free lib_dir that exists (Steam's
// compatibilitytools.d is absent until a custom Proton is installed), and only wildcard lib dirs
// are reserved. A same-user process creates the missing dir with a library in it, the unattended
// refresh turns it into a lib_dir, and the library preloads into Steam's whitelisted binaries.
func (s *IntegrationSuite) TestDaemon_Bypass_CatalogRefreshAdoptsPlantedLibDir() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.installFakeSystemctl(c)

	const (
		helper  = steamDir + "/ubuntu12_64/steamwebhelper"
		compat  = steamDir + "/compatibilitytools.d"
		plant   = compat + "/evil.so"
		shipped = steamDir + "/ubuntu12_64/libshipped.so"
	)
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/lib_probe.so"), libProbe, 0o755),
		"copy lib_probe.so")
	// ubuntu12_64 is Steam's shipped runtime (predates the daemon), handed to nobody so its library is
	// not an auto-trusted system file.
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /etc/app-listener " + steamDir + "/config " + steamDir + "/ubuntu12_64" +
			" && cp /usr/bin/dash " + helper + " && cp " + libProbe + " " + shipped +
			" && chown 65534 " + steamDir + "/ubuntu12_64 " + shipped +
			" && head -c 32 /dev/zero > /etc/app-listener/fscrypt.key && chmod 600 /etc/app-listener/fscrypt.key"})
	s.startDaemon(c, `[watch `+steamDir+`/config]
need_encryption: false
`+helper)

	// Attack: may be refused outright (the fix reserves the name); the outcome is what counts.
	_, attack := s.exec(c, []string{"sh", "-c",
		"mkdir " + compat + " && cp " + libProbe + " " + plant + " && chown -R 65534 " + compat + " 2>&1"})
	loaded, out := s.preload(c, helper, plant)
	s.Require().Falsef(loaded, "baseline: the plant must not load before the refresh: %s", out)

	code, refresh := s.exec(c, []string{"/app-listener", "install", "--update-catalog-only", "--live", "--yes"})
	if code == 0 {
		s.awaitLog(c, "configuration reloaded", 60*time.Second)
	}
	_, conf := s.exec(c, []string{"cat", "/etc/app-listener/daemon.conf"})

	loaded, out = s.preload(c, helper, plant)
	s.Require().Falsef(loaded,
		"a library in a lib dir a non-writer created was adopted by the unattended refresh and loaded into "+
			"steamwebhelper: %s\nattack: %s\nrefresh (exit %d):\n%s\ndaemon.conf:\n%s", out, attack, code, refresh, conf)

	// Controls: an existing Steam runtime dir is still adopted, and its library loads into Steam.
	s.Require().Containsf(conf, `lib_dir "`+steamDir+`/ubuntu12_64"`, "the shipped runtime must stay a lib_dir")
	loaded, out = s.preload(c, helper, shipped)
	s.Require().Truef(loaded, "control: Steam must load its shipped runtime library: %s", out)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Bypass: only the wildcard component of a lib dir glob is reserved, so inside a Proton* dir that
// lacks files/lib (the EasyAntiCheat runtime, old dist layouts) a non-writer creates the fixed tail,
// and the unattended refresh adopts it as a lib dir Steam's writers load from.
func (s *IntegrationSuite) TestDaemon_Bypass_WildcardLibDirTailPlanted() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.installFakeSystemctl(c)

	const (
		helper  = steamDir + "/ubuntu12_64/steamwebhelper"
		eac     = steamCommon + "/Proton_EasyAntiCheat_Runtime"
		plant   = eac + "/files/lib/evil.so"
		proton  = steamCommon + "/Proton_9.0/files/lib"
		shipped = proton + "/libwine.so"
	)
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/lib_probe.so"), libProbe, 0o755),
		"copy lib_probe.so")
	// Steam's Proton installs predate the daemon, handed to nobody so they aren't auto-trusted. No
	// spaces in the names: LD_PRELOAD splits on them, and the probe would never load either way.
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /etc/app-listener " + steamDir + "/config " + steamDir + "/ubuntu12_64 " +
			shQuote(eac) + " " + shQuote(proton) +
			" && cp /usr/bin/dash " + helper + " && cp " + libProbe + " " + shQuote(shipped) +
			" && chown -R 65534 " + steamCommon +
			" && head -c 32 /dev/zero > /etc/app-listener/fscrypt.key && chmod 600 /etc/app-listener/fscrypt.key"})
	s.startDaemon(c, `[watch `+steamDir+`/config]
need_encryption: false
`+helper)

	// Attack: may be refused outright (the fix reserves the tail); the outcome is what counts.
	_, attack := s.exec(c, []string{"sh", "-c",
		"mkdir -p " + shQuote(eac+"/files/lib") + " && cp " + libProbe + " " + shQuote(plant) +
			" && chown -R 65534 " + shQuote(eac+"/files") + " 2>&1; echo rc=$?"})
	loaded, out := s.preload(c, helper, plant)
	s.Require().Falsef(loaded, "baseline: the plant must not load before the refresh: %s", out)

	code, refresh := s.exec(c, []string{"/app-listener", "install", "--update-catalog-only", "--live", "--yes"})
	if code == 0 {
		s.awaitLog(c, "configuration reloaded", 60*time.Second)
	}
	_, conf := s.exec(c, []string{"cat", "/etc/app-listener/daemon.conf"})

	loaded, out = s.preload(c, helper, plant)
	s.Require().Falsef(loaded,
		"a lib dir tail a non-writer created under a Proton* dir was adopted by the unattended refresh and "+
			"its library loaded into steamwebhelper: %s\nattack: %s\nrefresh (exit %d):\n%s\ndaemon.conf:\n%s",
		out, attack, code, refresh, conf)

	// Controls: Steam's own Proton lib dir is adopted and its library loads into Steam.
	s.Require().Containsf(conf, `lib_dir "`+proton+`"`, "the genuine Proton lib dir must be adopted:\n%s", conf)
	loaded, out = s.preload(c, helper, shipped)
	s.Require().Truef(loaded, "control: Steam must load its Proton library: %s", out)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Bypass: daemon.conf stores expanded glob matches whose components below the wildcard are
// unreserved. Re-pointing one (files/bin -> /tmp/evil) is refused while the daemon runs, but a
// restart or reload re-parses the stored path and whitelists whatever it now resolves to.
func (s *IntegrationSuite) TestDaemon_Bypass_StoredMatchRepointedAcrossRestart() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const (
		marker  = "TOP-SECRET-STEAM-LOGIN-3C91"
		secret  = steamDir + "/config/loginusers.vdf"
		bin     = steamCommon + "/Proton_9.0/files/bin"
		ws      = bin + "/wineserver"
		control = steamCommon + "/Proton_8.0/files/bin/wineserver"
		evil    = "/tmp/evil/wineserver"
	)
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /etc/app-listener " + steamDir + "/config " + bin + " " + steamCommon + "/Proton_8.0/files/bin" +
			" && printf '" + marker + "' > " + secret + " && chmod 644 " + secret +
			" && cp /usr/bin/grep " + ws + " && cp /usr/bin/grep " + control})
	config := `[watch ` + steamDir + `/config]
need_encryption: false
` + ws + `
` + control
	s.startDaemon(c, config)

	readAs := func(exe string) string {
		_, out := s.exec(c, []string{"sh", "-c", exe + " -h . " + secret + " 2>&1; echo rc=$?"})
		return out
	}
	// grep, not cat: ubuntu:latest's cat is multi-call coreutils and refuses another basename.
	s.Require().NotContains(readAs("grep"), marker, "baseline: non-whitelisted grep must be denied")
	s.Require().Contains(readAs(ws), marker, "control: the stored wineserver must read the secret")

	_, attack := s.exec(c, []string{"sh", "-c",
		"mv " + bin + " " + bin + ".old && mkdir -p /tmp/evil && cp /usr/bin/grep " + evil +
			" && ln -s /tmp/evil " + bin + " 2>&1; echo rc=$?"})
	s.Require().NotContainsf(readAs(ws), marker, "baseline: the running daemon must refuse the re-pointed path "+
		"(attack: %s)", attack)

	s.sigDaemon(c, "TERM")
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not exit after SIGTERM")
	s.launchDaemon(c)
	s.awaitDaemonUp(c, config)
	s.Require().NotContainsf(readAs(evil), marker,
		"after a restart the stored path's re-pointed component got %s whitelisted\nattack: %s\ndaemon log:\n%s",
		evil, attack, s.readDaemonLog(c))
	s.Require().Contains(readAs(control), marker, "control: an intact stored match must still read after a restart")

	s.sigDaemon(c, "HUP")
	s.awaitLog(c, "configuration reloaded", 60*time.Second)
	s.Require().NotContainsf(readAs(evil), marker,
		"after a reload the stored path's re-pointed component got %s whitelisted\ndaemon log:\n%s",
		evil, s.readDaemonLog(c))
	s.Require().Contains(readAs(control), marker, "control: an intact stored match must still read after a reload")

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Bypass: every resource shares one guard_inodes map with one owner per inode. A read-only
// lib_dir's periodic sweep re-stats its root path through symlinks, so once a same-user process
// makes that path resolve to another resource's directory, the sweep re-anchors onto it and claims
// its inodes for the lib_dir, whose mode lets every process read.
// A file created after start (no inode row of its own) must stay the secret resource's too.
func (s *IntegrationSuite) TestDaemon_Bypass_LibDirSweepTakesOverResource() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const (
		marker = "TOP-SECRET-LIBDIR-TAKEOVER-5E07"
		secret = "/root/.ssh/id_secret"
		later  = "/root/.ssh/id_later"
		libF   = "/root/app/lib/libreal.txt"
	)
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /etc/app-listener /root/.ssh /root/app/lib /tmp/sshw /tmp/appw" +
			" && printf '" + marker + "' > " + secret + " && chmod 644 " + secret +
			" && echo LIB-CONTENT-OK > " + libF +
			" && cp /usr/bin/grep /tmp/sshw/grep && cp /usr/bin/dash /tmp/sshw/dash && cp /usr/bin/dash /tmp/appw/dash"})
	s.startDaemon(c, `[libraries "App"]
lib_dir /root/app/lib
lib_binary /tmp/appw/dash

[watch /root/.ssh]
need_encryption: false
/tmp/sshw/grep
/tmp/sshw/dash`)

	// grep, not cat: a copy of ubuntu:latest's multi-call cat outside root-placed dirs is opaque, so
	// the daemon drops its whitelist line.
	readAs := func(bin, path string) string {
		_, out := s.exec(c, []string{"sh", "-c", bin + " -h . " + path + " 2>&1; echo rc=$?"})
		return out
	}
	code, out := s.exec(c, []string{"/tmp/sshw/dash", "-c", "printf '" + marker + "' > " + later})
	s.Require().Equalf(0, code, "the whitelisted dash must create %s: %s", later, out)
	// Controls: the secrets are denied to a non-whitelisted reader and served to the whitelisted one;
	// the lib_dir's own content is world-readable.
	for _, f := range []string{secret, later} {
		s.Require().NotContainsf(readAs("grep", f), marker, "baseline: non-whitelisted grep must be denied %s", f)
		s.Require().Containsf(readAs("/tmp/sshw/grep", f), marker, "control: whitelisted grep must read %s", f)
	}
	s.Require().Contains(readAs("grep", libF), "LIB-CONTENT-OK", "control: the lib_dir is readable")

	// Attack, as a non-whitelisted process: only unguarded names change (the lib_dir's parent and a
	// new entry in its place). A relative target is not a watch-root path for path_symlink.
	_, out = s.exec(c, []string{"sh", "-c",
		"mv /root/app /root/app.old && mkdir /root/app && ln -s ../.ssh /root/app/lib; echo rc=$?; ls -la /root/app"})
	s.Require().Containsf(out, "rc=0", "the lib_dir's path must be re-pointable by a non-whitelisted process: %s", out)

	// Past at least two sweep ticks (resyncSweepEvery = 30s).
	for dl := time.Now().Add(75 * time.Second); time.Now().Before(dl); time.Sleep(3 * time.Second) {
		for _, f := range []string{secret, later} {
			out = readAs("grep", f)
			s.Require().NotContainsf(out, marker,
				"the lib_dir sweep re-anchored onto /root/.ssh and made %s readable to every process: %s\n"+
					"daemon log:\n%s", f, out, s.readDaemonLog(c))
		}
	}
	for _, f := range []string{secret, later} {
		s.Require().Containsf(readAs("/tmp/sshw/grep", f), marker, "control: whitelisted grep must still read %s", f)
	}
	log := s.readDaemonLog(c)
	s.Require().Containsf(log, "refusing to re-anchor",
		"the sweep must have seen and refused the re-pointed root, daemon log:\n%s", log)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}
