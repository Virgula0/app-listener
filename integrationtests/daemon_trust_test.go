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
	// an unreserved directory name (a symlink) is created inside the user's tree.
	code, out := s.exec(c, []string{"sh", "-c",
		"mkdir -p /tmp/x /tmp/y && cp /exploits/xres_move /tmp/x/Discord && cp /exploits/xres_move /tmp/y/claude" +
			" && ln -s /tmp/x " + discordDir + "/app-0.0.999 && ln -s /tmp/y /root/.local/bin 2>&1"})
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
