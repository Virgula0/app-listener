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
