package integrationtests

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/testcontainers/testcontainers-go"
	tcexec "github.com/testcontainers/testcontainers-go/exec"

	"github.com/Virgula0/app-listener/cmd/functions/editprotected"
	"github.com/Virgula0/app-listener/internal/guard"
)

// ---------------------------------------------------------------
// Daemon helpers
// ---------------------------------------------------------------

type daemonEvent struct {
	Denied   bool
	Op       string
	Comm     string
	UID      string
	Path     string
	Resource string
}

// parseDaemonEvents parses the daemon headless writer format:
//
//	<4>DAEMON DENIED  op=OPEN  comm=grep  pid=1  uid=nobody  resource=/p  path=/p/f
//	<6>DAEMON ALLOWED op=OPEN  comm=cat   pid=1  uid=root    resource=/p  path=/p/f
func parseDaemonEvents(logContent string) []daemonEvent {
	var events []daemonEvent
	for _, line := range strings.Split(logContent, "\n") {
		idx := strings.Index(line, "DAEMON ")
		if idx < 0 {
			continue
		}
		denied := strings.Contains(line, "DAEMON DENIED")
		ev := daemonEvent{Denied: denied}
		for _, field := range strings.Fields(line) {
			if !strings.Contains(field, "=") {
				continue
			}
			key, value, _ := strings.Cut(field, "=")
			switch key {
			case "op":
				ev.Op = value
			case "comm":
				ev.Comm = value
			case "uid":
				ev.UID = value
			case "path":
				ev.Path = value
			case "resource":
				ev.Resource = value
			}
		}
		events = append(events, ev)
	}
	return events
}

// startDaemon writes the config, starts the daemon headless and waits for
// the pid file (written after guards attach and readers run).
func (s *IntegrationSuite) startDaemon(c testcontainers.Container, config string) {
	writeConfig := fmt.Sprintf("cat > /etc/app-listener/daemon.conf <<'EOF'\n%s\nEOF", config)
	s.exec(c, []string{"sh", "-c", writeConfig})
	s.launchDaemon(c)
	s.awaitDaemonUp(c, config)
}

// watchPathsInConfig extracts every "[watch <path>]" section header's path
// from a daemon.conf, in order, for awaitDaemonUp to wait on each one's own
// "guard started" line individually rather than stopping at the first.
func watchPathsInConfig(config string) []string {
	var paths []string
	for _, line := range strings.Split(config, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "[watch ") || !strings.HasSuffix(line, "]") {
			continue
		}
		if path := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "[watch "), "]")); path != "" {
			paths = append(paths, path)
		}
	}
	return paths
}

// launchDaemon backgrounds the daemon on the already-written daemon.conf without waiting for it.
// Split from startDaemon so a caller can start a racer right before this call and stop it once
// awaitDaemonUp returns (daemon_toctou_test.go).
func (s *IntegrationSuite) launchDaemon(c testcontainers.Container) {
	cmd := "nohup /app-listener daemon --config /etc/app-listener/daemon.conf --headless > /tmp/daemon.log 2>&1 &"
	code, out := s.exec(c, []string{"sh", "-c", cmd})
	s.Require().Equalf(0, code, "starting daemon: %s", out)
}

// awaitDaemonUp polls until the guards are attached and event readers running for EVERY resource in
// config: the pid file, then one "guard started - guarding: <path>" line per "[watch <path>]".
// Waiting for only the first line races on multi-resource configs: each guard attach is a
// multi-second verifier pass on slow hosts (daemonShutdownTimeout), so a caller could touch a LATER
// resource before its guard or unlock is live. config == "" (relaunch on an existing daemon.conf)
// falls back to "at least one guard started".
func (s *IntegrationSuite) awaitDaemonUp(c testcontainers.Container, config string) {
	deadline := time.Now().Add(daemonShutdownTimeout)
	ready := false
	for time.Now().Before(deadline) {
		code, _ := s.exec(c, []string{"sh", "-c", "test -f /run/app-listener-daemon.pid && echo ready"})
		if code == 0 {
			ready = true
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	log := s.readDaemonLog(c)
	if !ready && strings.Contains(log, "OCI runtime exec failed") {
		// Transient host-level runc flake: the daemon never ran — retry once.
		s.launchDaemon(c)
		deadline = time.Now().Add(daemonShutdownTimeout)
		for time.Now().Before(deadline) {
			code, _ := s.exec(c, []string{"sh", "-c", "test -f /run/app-listener-daemon.pid && echo ready"})
			if code == 0 {
				ready = true
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
		log = s.readDaemonLog(c)
	}
	if !ready {
		s.Require().Failf("daemon did not attach its guards", "pid file missing after %s, daemon log:\n%s", daemonShutdownTimeout, log)
	}

	wantMarkers := []string{"guard started"}
	if paths := watchPathsInConfig(config); len(paths) > 0 {
		wantMarkers = wantMarkers[:0]
		for _, p := range paths {
			wantMarkers = append(wantMarkers, "guard started — guarding: "+p)
		}
	}
	allPresent := func(l string) bool {
		for _, m := range wantMarkers {
			if !strings.Contains(l, m) {
				return false
			}
		}
		return true
	}

	// Guards are attached and readers running before the pid file appears;
	// poll for every expected guard-started marker instead of a fixed settle
	// sleep or stopping at the first resource's.
	deadline2 := time.Now().Add(daemonShutdownTimeout)
	for time.Now().Before(deadline2) {
		if allPresent(log) {
			return
		}
		time.Sleep(200 * time.Millisecond)
		log = s.readDaemonLog(c)
	}
	s.Require().Truef(allPresent(log), "not every configured resource's guard started within %s, daemon log:\n%s", daemonShutdownTimeout, log)
}

func (s *IntegrationSuite) readDaemonLog(c testcontainers.Container) string {
	_, out := s.exec(c, []string{"sh", "-c", "cat /tmp/daemon.log 2>/dev/null || true"})
	return out
}

// ---------------------------------------------------------------
// Test: the daemon self-whitelist must not be a universal key
// ---------------------------------------------------------------

// Regression for the self-key bypass: the daemon whitelists its own exe inode on every resource for
// its fscrypt ioctls. Identity was keyed by exe inode alone, so ANY local user executing the same
// binary got full allow on the tree (e.g. `app-listener network-monitor <guarded-file>` hashes the
// file before any capability check). The self entry must grant only uid 0 and only the needed
// events (OPEN, READ).
func (s *IntegrationSuite) TestDaemon_SelfWhitelist_NoUniversalKey() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /etc/app-listener && echo TOP-SECRET-CONTENT > /protected/secret && chmod 755 /protected && chmod 644 /protected/secret"})
	s.startDaemon(c, `[watch /protected]
need_encryption: false
/usr/bin/sleep`)

	const nobody = "setpriv --reuid=65534 --regid=65534 --clear-groups"

	// Sanity: the whitelist guard enforces for ordinary binaries (grep has
	// its own inode on this image; the uutils applets hardlink together).
	code, out := s.exec(c, []string{"sh", "-c", nobody + " grep -c TOP-SECRET /protected/secret"})
	s.Require().NotEqualf(0, code, "baseline: non-whitelisted reader must be denied, got: %s", out)

	// Exploit: a non-root process executes the daemon's own binary inode and
	// makes it read the guarded file (network-monitor hashes its positional
	// arguments before any capability check). Pre-fix this read is allowed.
	code, out = s.exec(c, []string{"sh", "-c",
		nobody + " timeout 10 /app-listener network-monitor /protected/secret 2>&1"})
	s.T().Logf("exploit exit=%d out=%q", code, out)

	log := s.readDaemonLog(c)
	events := parseDaemonEvents(log)

	// The self-key must never let a NON-ROOT app-listener process touch the
	// guarded file. (The daemon itself — uid=root — legitimately stats every
	// file under /protected during its inode scan; those ALLOWED STAT events
	// are expected and are not the bypass.)
	for _, ev := range events {
		if ev.Comm == "app-listener" && ev.Path == "/protected/secret" && !ev.Denied && ev.UID != "root" {
			s.Require().Failf("self-key bypass",
				"daemon allowed uid=%s app-listener %s of %s (self-whitelist acted as a universal key)", ev.UID, ev.Op, ev.Path)
		}
	}

	// And the guarded access attempt by the non-root process must be visible
	// as denied enforcement, not silently dropped.
	denied := false
	for _, ev := range events {
		if ev.Comm == "app-listener" && ev.Path == "/protected/secret" && ev.Denied && ev.UID != "root" {
			denied = true
		}
	}
	s.Require().True(denied, "expected a DAEMON DENIED event for the non-root app-listener access of /protected/secret, got: %s", log)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Bypass (finding #1): replacing a whitelisted binary at a user-writable path via rename gets the
// REPLACEMENT re-whitelisted. Identity is the exe inode and the daemon re-admits whatever inode is
// now at the path (ReSyncBinaries, on any denial and the periodic sweep) with no check that it's
// the same binary, root-owned, or in a non-user-writable dir. The in-place swap is caught by the
// hash-verify loop, but a RENAME gives a NEW inode via a path that only compares the stored inode.
// Any home-directory whitelisted binary (Claude Code, Discord) could be swapped by same-user
// malware. Re-admission must verify the binary (hash, or refuse non-root-owned/user-writable
// paths): the replacement stays denied.
func (s *IntegrationSuite) TestDaemon_Bypass_BinaryRenameReplaceReWhitelisted() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const marker = "TOP-SECRET-RENAME-REWHITELIST-7B3D"
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /exploits /etc/app-listener && printf '" + marker + "' > /protected/secret && chmod 755 /protected && chmod 644 /protected/secret"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/swap_benign"), "/exploits/swap_benign", 0755), "copy swap_benign")
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/swap_reader"), "/exploits/swap_reader", 0755), "copy swap_reader")

	// The whitelisted binary lives at a user-writable path (mirrors a
	// home-directory app), initially the benign placeholder.
	s.exec(c, []string{"sh", "-c", "cp /exploits/swap_benign /tmp/app && chmod 755 /tmp/app"})
	s.startDaemon(c, `[watch /protected]
need_encryption: false
/tmp/app`)

	// Baseline: the whitelisted benign binary runs (proves the whitelist is
	// active for /tmp/app's current inode).
	code, out := s.exec(c, []string{"/tmp/app"})
	s.Require().Equalf(0, code, "baseline whitelisted binary should run: %s", out)
	s.Require().Contains(out, "BENIGN-APP-OK")

	// Attack: try to rename the malicious reader OVER the whitelisted path.
	// The binary lives at a user-writable path (/tmp), so the trust guard's
	// writer-attribution (#1) must deny the replacement by a non-whitelisted
	// process (mv), leaving the inode unchanged.
	_, inodeOut := s.exec(c, []string{"sh", "-c",
		"i1=$(stat -c %i /tmp/app); mv -f /exploits/swap_reader /tmp/app 2>/dev/null; i2=$(stat -c %i /tmp/app); echo $i1 $i2"})
	inodeRe := regexp.MustCompile(`(\d+)[^0-9]+(\d+)`)
	m := inodeRe.FindStringSubmatch(inodeOut)
	s.Require().NotNil(m, "inode probe output: %q", inodeOut)
	s.Require().Equalf(m[1], m[2],
		"protected whitelisted binary was replaced by a non-app process (before=%s after=%s)", m[1], m[2])

	// Belt-and-suspenders: even if the swap had somehow gone through, the
	// replacement must never be able to read the guarded secret. Poll the
	// whole re-sync window (denial-driven + the 30s periodic sweep).
	leaked := false
	var lastOut string
	deadline := time.Now().Add(50 * time.Second)
	for time.Now().Before(deadline) {
		code, out = s.exec(c, []string{"sh", "-c", "timeout 10 /tmp/app /protected/secret 2>&1"})
		lastOut = out
		if strings.Contains(out, "STOLEN|"+marker) {
			leaked = true
			break
		}
		time.Sleep(2 * time.Second)
	}

	// Secure expectation: the renamed-in replacement never gained the
	// whitelist entry, so the secret is never leaked.
	s.Require().Falsef(leaked,
		"renamed-in binary was re-whitelisted and read the guarded secret: %q", lastOut)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Regression #2: a whitelisted binary run with LD_PRELOAD pointing at an attacker .so under a
// user-writable path (/tmp) must NOT map it: the trust guard's library allowlist denies it (not
// allow_lib, not an auto-trusted root-owned system lib, not inside a guarded tree), so the
// constructor never runs. A plain run still works (real libraries are root-owned, auto-trusted).
// Always enforced.
func (s *IntegrationSuite) TestDaemon_Bypass_LdPreloadWhitelistedBinary() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const marker = "TOP-SECRET-LD-PRELOAD-5E7A"
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /etc/app-listener && printf '" + marker + "' > /protected/secret && chmod 755 /protected && chmod 644 /protected/secret"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/preload_leak.so"), "/tmp/leak.so", 0755), "copy preload_leak.so")

	s.startDaemon(c, `[watch /protected]
need_encryption: false
/usr/bin/true`)

	// Positive control: the whitelisted binary runs normally — its real
	// libraries are root-owned system libs and are auto-trusted, so
	// enforcement does not break it.
	code, out := s.exec(c, []string{"/usr/bin/true"})
	s.Require().Equalf(0, code, "whitelisted binary must still run under library enforcement: %s", out)

	// Attack: LD_PRELOAD an attacker .so from a user-writable path.
	_, out = s.exec(c, []string{"sh", "-c",
		"LD_PRELOAD=/tmp/leak.so LEAK_FILE=/protected/secret /usr/bin/true 2>&1"})
	s.Require().NotContainsf(out, marker,
		"guarded content exfiltrated through LD_PRELOAD into a whitelisted binary: %s", out)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Bypass (finding #5): is_system_trusted auto-trusts any library reported as root-owned in a
// root-owned directory. File ownership on a user-mountable filesystem (FUSE) is whatever the
// unprivileged mounter's server claims, so it proves nothing about root control: a normal user can
// mount a FUSE fs that reports an attacker library as root:root and LD_PRELOAD it into a whitelisted
// binary. The fix refuses to system-trust files on a FUSE superblock; a genuine root-owned library
// on the real filesystem must still load (positive control).
//
// Uses bindfs (a FUSE fs that can force root:root ownership) so no custom FUSE server is needed.
// Skips cleanly if FUSE/bindfs is unavailable in the container.
func (s *IntegrationSuite) TestDaemon_Bypass_FuseRootOwnedLibNotAutoTrusted() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	code, out := s.exec(c, []string{"sh", "-c",
		"command -v bindfs >/dev/null 2>&1 || (timeout 120 apt-get update -qq && " +
			"timeout 120 apt-get install -y -qq --no-install-recommends bindfs)"})
	if code2, _ := s.exec(c, []string{"sh", "-c", "command -v bindfs >/dev/null 2>&1"}); code2 != 0 {
		s.T().Skipf("bindfs (FUSE) could not be installed — this test needs it: %s", out)
	}
	// A privileged container usually has /dev/fuse; create the node and load the module if not.
	s.exec(c, []string{"sh", "-c",
		"test -e /dev/fuse || mknod -m 666 /dev/fuse c 10 229; modprobe fuse 2>/dev/null || true"})

	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /exploits /etc/app-listener /src /fusemnt /realsys && chmod 755 /protected /src /fusemnt /realsys"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/lib_probe.so"), "/exploits/lib_probe.so", 0o755), "copy lib_probe.so")
	// The library the attacker wants loaded: lib_probe announces LIB_PROBE_LOADED from its
	// constructor when ld.so maps and runs it.
	s.exec(c, []string{"sh", "-c", "cp /exploits/lib_probe.so /src/evil.so && chmod 644 /src/evil.so"})
	// Positive-control copy on the REAL fs, genuinely root:root in a root:root dir.
	s.exec(c, []string{"sh", "-c", "cp /exploits/lib_probe.so /realsys/evil.so && chown 0:0 /realsys/evil.so && chmod 644 /realsys/evil.so"})

	code, out = s.exec(c, []string{"sh", "-c",
		"bindfs --force-user=root --force-group=root --perms=0644 /src /fusemnt 2>&1"})
	if code != 0 {
		s.T().Skipf("bindfs mount unavailable in this container (no /dev/fuse or no SYS_ADMIN): %s", out)
	}
	defer s.exec(c, []string{"sh", "-c", "fusermount -u /fusemnt 2>/dev/null || umount /fusemnt 2>/dev/null || true"})

	// Confirm the mount really presents a FUSE fs reporting the library as root-owned; otherwise the
	// test would pass vacuously.
	_, ft := s.exec(c, []string{"sh", "-c", "stat -f -c %T /fusemnt 2>&1; stat -c '%U:%G' /fusemnt/evil.so 2>&1"})
	if !strings.Contains(ft, "fuse") || !strings.Contains(ft, "root:root") {
		s.T().Skipf("bindfs did not present a FUSE root-owned file (got %q) — cannot exercise the bypass", ft)
	}

	s.startDaemon(c, `[watch /protected]
need_encryption: false
/usr/bin/true`)

	// Positive control: a genuine root-owned library on the real filesystem is still auto-trusted
	// and loads into the whitelisted binary — the fix must not break system-library trust.
	_, out = s.exec(c, []string{"sh", "-c", "LD_PRELOAD=/realsys/evil.so /usr/bin/true 2>&1"})
	s.Require().Containsf(out, libProbeMarker,
		"a genuine root-owned system library on the real fs must still load (is_system_trusted regressed): %s", out)

	// Attack: the FUSE-hosted library reports root:root but is attacker-controlled. It must NOT be
	// auto-trusted, so ld.so cannot map it into the whitelisted binary.
	_, out = s.exec(c, []string{"sh", "-c", "LD_PRELOAD=/fusemnt/evil.so /usr/bin/true 2>&1"})
	s.Require().NotContainsf(out, libProbeMarker,
		"a FUSE-hosted 'root-owned' library was auto-trusted and loaded into a whitelisted binary: %s", out)
	s.requireDenialLogged(c, "LIBLOAD", "evil.so")

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// The trust guard's trusted set is built once at daemon start and is never rebuilt on SIGHUP, but
// SIGHUP is the normal way the whitelist changes: the pacman PostTransaction and apt
// DPkg::Post-Invoke catalog-refresh hooks reload rather than restart. A binary whitelisted by a
// reload is therefore attached to its per-resource guard while being ABSENT from
// guard_trusted_files, and trust_mmap returns early unless the mapping process's exe carries
// TRUSTED_BINARY — so it gets no library allowlist at all and LD_PRELOAD works against it, even
// though the identical binary whitelisted at startup is protected.
func (s *IntegrationSuite) TestDaemon_Bypass_LdPreloadBinaryWhitelistedByReload() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const marker = "TOP-SECRET-RELOAD-TRUST-9C3F"
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected && printf '" + marker + "' > /protected/secret && chmod 755 /protected && chmod 644 /protected/secret"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/preload_leak.so"), "/tmp/leak.so", 0755), "copy preload_leak.so")

	// The config lives OUTSIDE /etc/app-listener so this test can rewrite it: once the daemon runs,
	// /etc/app-listener is self-guarded read-only and only the app-listener binary may write it (the
	// real catalog refresh does so via `install --update-catalog-only`). /usr/bin/true is
	// whitelisted at startup; /usr/bin/cat is added only by the SIGHUP reload below, mirroring a
	// catalog refresh picking up a new binary.
	const cfgPath = "/tmp/poc-daemon.conf"
	writeCfg := func(body string) {
		code, out := s.exec(c, []string{"sh", "-c",
			fmt.Sprintf("cat > %s <<'EOF'\n%s\nEOF", cfgPath, body)})
		s.Require().Equalf(0, code, "writing %s: %s", cfgPath, out)
	}
	writeCfg("[watch /protected]\nneed_encryption: false\n/usr/bin/true")

	launch := "nohup /app-listener daemon --config " + cfgPath + " --headless > /tmp/daemon.log 2>&1 &"
	code, out := s.exec(c, []string{"sh", "-c", launch})
	s.Require().Equalf(0, code, "starting daemon: %s", out)
	s.awaitDaemonUp(c, "[watch /protected]")

	// Baseline: LD_PRELOAD into the startup-whitelisted binary is denied (trust set has it).
	_, out = s.exec(c, []string{"sh", "-c",
		"LD_PRELOAD=/tmp/leak.so LEAK_FILE=/protected/secret /usr/bin/true 2>&1"})
	s.Require().NotContainsf(out, marker, "baseline: startup-whitelisted binary must be preload-protected: %s", out)

	writeCfg("[watch /protected]\nneed_encryption: false\n/usr/bin/true\n/usr/bin/cat")
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
	s.Require().Truef(done, "reload did not complete within %s, daemon log:\n%s", daemonShutdownTimeout, s.readDaemonLog(c))

	// Positive control: the reload-added binary really is whitelisted now (per-resource access).
	code, out = s.exec(c, []string{"sh", "-c", "cat /protected/secret"})
	s.Require().Equalf(0, code, "reload-added binary should be whitelisted for the resource: %s", out)
	s.Require().Containsf(out, marker, "reload-added binary should read the secret normally: %s", out)

	// Attack: LD_PRELOAD an attacker .so into the reload-added whitelisted binary. Identical to the
	// baseline that was denied, only the binary was whitelisted by reload instead of at startup.
	_, out = s.exec(c, []string{"sh", "-c",
		"LD_PRELOAD=/tmp/leak.so LEAK_FILE=/protected/secret /usr/bin/cat /dev/null 2>&1"})
	s.Require().NotContainsf(out, marker,
		"guarded content exfiltrated through LD_PRELOAD into a binary whitelisted by SIGHUP: "+
			"the trust set was not rebuilt on reload: %s", out)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Bypass (finding #4c): memory-read taint lost across a SIGHUP reload. A reload rebuilds every
// guard with fresh BPF maps, emptying the tainted-pid set while the victim keeps the secret in
// memory. Taint must survive a reload (re-seed from the surviving map or persist it): the
// post-reload dump stays denied.
func (s *IntegrationSuite) TestDaemon_Bypass_TaintLostOnReload() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const marker = "TOP-SECRET-TAINT-RELOAD-6C4E"
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /exploits /etc/app-listener && printf '" + marker + "' > /protected/secret && chmod 755 /protected && chmod 644 /protected/secret"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/taint_victim"), "/exploits/taint_victim", 0755), "copy taint_victim")
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/process_vm_readv"), "/exploits/process_vm_readv", 0755), "copy process_vm_readv")
	s.exec(c, []string{"sh", "-c", "cp /exploits/taint_victim /tmp/victim && chmod 755 /tmp/victim"})

	config := `[watch /protected]
need_encryption: false
/tmp/victim`
	s.startDaemon(c, config)

	// The whitelisted victim reads the secret into its heap and stays
	// alive, tainted, holding the fd.
	s.exec(c, []string{"sh", "-c", "/tmp/victim hold /protected/secret & echo $! > /tmp/vpid; sleep 1"})

	// Before the reload: the attacker is blocked by taint tracking and
	// cannot see the secret. (Proves the victim is tainted.)
	_, out := s.exec(c, []string{"sh", "-c", "/exploits/process_vm_readv /protected/secret 2>&1"})
	s.Require().NotContainsf(out, marker,
		"pre-reload control: taint should have blocked the dump: %s", out)
	// The refusal must be VISIBLE: the process gates used to deny silently,
	// which once left a whole application failing with nothing in the log.
	_, logOut := s.exec(c, []string{"sh", "-c", "grep -c 'DAEMON DENIED  op=PTRACE' /tmp/daemon.log || true"})
	s.Require().NotEqualf("0", strings.TrimSpace(logOut),
		"the ptrace-class denial of process_vm_readv must be logged as op=PTRACE")

	// Wait for the reload to COMPLETE before the post-reload attempt: the "guard started" markers
	// awaitDaemonUp watches predate it, so only the reload-complete line proves the map swap (else
	// the attacker races the still-tainted old guard and the test passes for the wrong reason).
	s.sigDaemon(c, "HUP")
	const reloadDone = "configuration reloaded without dropping protection"
	reloaded := false
	deadline := time.Now().Add(daemonShutdownTimeout)
	for time.Now().Before(deadline) {
		if strings.Contains(s.readDaemonLog(c), reloadDone) {
			reloaded = true
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
	s.Require().Truef(reloaded, "reload did not complete within %s, daemon log:\n%s", daemonShutdownTimeout, s.readDaemonLog(c))

	// After the reload: the victim is still alive and holds the secret,
	// but the new guard's taint map is empty.
	_, out = s.exec(c, []string{"sh", "-c", "/exploits/process_vm_readv /protected/secret 2>&1"})
	s.Require().NotContainsf(out, marker,
		"taint was lost across the reload; attacker dumped the secret from the surviving victim: %s", out)

	s.exec(c, []string{"sh", "-c", "pkill -f taint_victim 2>/dev/null; pkill -f 'app-listener daemon' || true"})
}

// Bypass (finding #1): guard_path_rename and guard_path_link decide on the SOURCE side and return
// before checking the destination (guard.bpf.c ~1475/~1653). Since 84df28d moved every resource
// onto one hook set, a binary whitelisted for resource A can move/exchange/link into resource B
// (which does not whitelist it). The headline is confidentiality: a RENAME_EXCHANGE of a file in A
// with B's secret swaps their names, so B's content lands at a path in A that the SAME whitelisted
// identity may read. The two operations must be judged against BOTH sides' whitelists; a binary
// that is not a writer of B must be denied at B whatever its standing in A.
func (s *IntegrationSuite) TestDaemon_Bypass_CrossResourceRenameLinkOnlySourceChecked() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const marker = "TOP-SECRET-CROSS-RENAME-2D9B"
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /resourceA /resourceB /exploits /etc/app-listener && " +
			"printf '" + marker + "' > /resourceB/secret && chmod 755 /resourceA /resourceB && " +
			"chmod 644 /resourceB/secret && printf 'DECOY-CONTENT' > /resourceA/decoy && chmod 644 /resourceA/decoy"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/xres_move"), "/exploits/xres_move", 0o755),
		"copy xres_move")
	s.exec(c, []string{"sh", "-c", "cp /exploits/xres_move /tmp/mover && chmod 755 /tmp/mover"})

	// /tmp/mover is whitelisted for A only; /resourceB whitelists /usr/bin/true, never the mover.
	s.startDaemon(c, `[watch /resourceA]
need_encryption: false
/tmp/mover

[watch /resourceB]
need_encryption: false
/usr/bin/true`)

	// Baselines that make the exploit meaningful.
	code, out := s.exec(c, []string{"/tmp/mover", "read", "/resourceA/decoy"})
	s.Require().Equalf(0, code, "the mover is whitelisted for A and must read A: %s", out)
	code, out = s.exec(c, []string{"/tmp/mover", "read", "/resourceB/secret"})
	s.Require().NotEqualf(0, code, "the mover is NOT whitelisted for B and must be denied a direct read: %s", out)
	code, out = s.exec(c, []string{"sh", "-c", "setpriv --reuid=65534 --regid=65534 --clear-groups cat /resourceB/secret 2>&1"})
	s.Require().NotEqualf(0, code, "a non-whitelisted reader of B must be denied: %s", out)

	// Confidentiality: exchange A's decoy with B's secret. Must be denied at B's side.
	code, out = s.exec(c, []string{"/tmp/mover", "exchange", "/resourceA/decoy", "/resourceB/secret"})
	s.Require().NotEqualf(0, code,
		"a RENAME_EXCHANGE into resource B by a binary whitelisted only for A must be denied: %s", out)
	// Whether or not the exchange was refused, B's secret must never become readable through A.
	_, out = s.exec(c, []string{"/tmp/mover", "read", "/resourceA/decoy"})
	s.Require().NotContainsf(out, marker,
		"B's secret leaked into resource A via cross-resource exchange and was read by the A-whitelisted binary: %s", out)

	// Integrity: plant into B by rename and by hardlink from A. Both are source-side allowed today.
	code, out = s.exec(c, []string{"/tmp/mover", "rename", "/resourceA/decoy", "/resourceB/planted"})
	s.Require().NotEqualf(0, code, "renaming a file from A into B by an A-only binary must be denied: %s", out)
	code, _ = s.exec(c, []string{"sh", "-c", "test -e /resourceB/planted"})
	s.Require().NotEqualf(0, code, "the rename planted a file inside resource B")
	code, out = s.exec(c, []string{"/tmp/mover", "link", "/resourceA/decoy", "/resourceB/aliased"})
	s.Require().NotEqualf(0, code, "hardlinking a file from A into B by an A-only binary must be denied: %s", out)
	code, _ = s.exec(c, []string{"sh", "-c", "test -e /resourceB/aliased"})
	s.Require().NotEqualf(0, code, "the hardlink planted a file inside resource B")

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Bypass (finding #3): a reload builds a replacement guard for every resource in the NEW config
// before it validates the change. That replacement takes over the live guard's guard_inodes row.
// When the reload is then refused (Phase 0: a resource was dropped) the replacement is stopped and
// deletes the row it took, so a resource the daemon claims to be "keeping" is left readable while
// its old guard is still attached. Nothing re-adds the row (SweepInodes needs a changed root inode,
// ReconcileInodes only deletes).
//
// The config lives outside /etc/app-listener so the test can rewrite it (once running, the daemon
// self-guards its own config dir). Dropping /dropB makes the SIGHUP reload refuse; /keepA is a
// single-file root, so losing its one row fully unprotects it.
func (s *IntegrationSuite) TestDaemon_RefusedReload_KeepsRemainingResourceProtected() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const marker = "TOP-SECRET-REFUSED-RELOAD-6B4C"
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /dropB && printf '" + marker + "' > /keepA && chmod 644 /keepA && echo x > /dropB/f"})

	const cfgPath = "/tmp/poc-daemon.conf"
	writeCfg := func(body string) {
		code, out := s.exec(c, []string{"sh", "-c", fmt.Sprintf("cat > %s <<'EOF'\n%s\nEOF", cfgPath, body)})
		s.Require().Equalf(0, code, "writing %s: %s", cfgPath, out)
	}
	writeCfg("[watch /keepA]\nneed_encryption: false\n/usr/bin/true\n\n[watch /dropB]\nneed_encryption: false\n/usr/bin/true")

	launch := "nohup /app-listener daemon --config " + cfgPath + " --headless > /tmp/daemon.log 2>&1 &"
	code, out := s.exec(c, []string{"sh", "-c", launch})
	s.Require().Equalf(0, code, "starting daemon: %s", out)
	s.awaitDaemonUp(c, "[watch /keepA]\n[watch /dropB]")

	// Baseline: /keepA is guarded (a non-whitelisted reader is denied).
	code, out = s.exec(c, []string{"sh", "-c", "setpriv --reuid=65534 --regid=65534 --clear-groups cat /keepA 2>&1"})
	s.Require().NotEqualf(0, code, "baseline: /keepA must be guarded before the reload: %s", out)

	// Drop /dropB and reload: Phase 0 refuses (removing a resource would leave it unprotected).
	writeCfg("[watch /keepA]\nneed_encryption: false\n/usr/bin/true")
	s.sigDaemon(c, "HUP")

	refused := false
	for deadline := time.Now().Add(daemonShutdownTimeout); time.Now().Before(deadline); {
		if strings.Contains(s.readDaemonLog(c), "keeping previous configuration") {
			refused = true
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
	s.Require().Truef(refused, "the reload dropping /dropB was not refused, daemon log:\n%s", s.readDaemonLog(c))

	// The kept resource must still be guarded. Pre-fix the rolled-back replacement's Stop deleted
	// /keepA's only inode row, so the still-attached old guard fails open.
	code, out = s.exec(c, []string{"sh", "-c", "setpriv --reuid=65534 --regid=65534 --clear-groups cat /keepA 2>&1"})
	s.Require().NotEqualf(0, code,
		"a refused reload left the kept resource /keepA readable — its inode row was deleted by the rolled-back guard: %s", out)
	s.Require().NotContainsf(out, marker, "the kept resource's secret was disclosed after a refused reload: %s", out)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// daemonPinFiles lists the daemon's pin files under /sys/fs/bpf/app-listener.
func (s *IntegrationSuite) daemonPinFiles(c testcontainers.Container) []string {
	_, out := s.exec(c, []string{"sh", "-c", "ls -1 /sys/fs/bpf/app-listener 2>/dev/null || true"})
	return strings.Fields(out)
}

// The daemon protects its own on-disk state. While it runs, /etc/app-listener is guarded
// independent of any [watch]:
//   - daemon.conf stays world-READABLE, but nothing except the app-listener binary may
//     write/rename-over/delete it or create files in the directory (ModeReadOnly);
//   - fscrypt.key is unreadable except by the app-listener binary (ModeWhitelist, empty list;
//     stacks on the RO dir guard).
//
// The binary itself (uid 0) keeps full access so install / --genkey / --update-catalog-only work.
func (s *IntegrationSuite) TestDaemon_SelfProtection_ConfigAndKey() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /protected2 /etc/app-listener && echo SECRET > /protected/secret && " +
			"head -c 32 /dev/zero > /etc/app-listener/fscrypt.key && chmod 600 /etc/app-listener/fscrypt.key"})
	s.startDaemon(c, `[watch /protected]
need_encryption: false
/usr/bin/sleep

[watch /protected2]
need_encryption: false
/usr/bin/sleep`)

	// Self guards attach after the daemon signals ready (best-effort
	// hardening, slow LSM attach on hardened kernels) — wait for both.
	selfReady := false
	for dl := time.Now().Add(45 * time.Second); time.Now().Before(dl); {
		if strings.Contains(s.readDaemonLog(c), "self-protection: guarding /etc/app-listener/fscrypt.key") {
			selfReady = true
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
	log := s.readDaemonLog(c)
	s.Require().Truef(selfReady, "self-protection guards did not attach, log:\n%s", log)
	s.Require().Contains(log, "/etc/app-listener (readonly)")

	// 0. `systemctl reload` runs helpers (ExecReload=/bin/kill ...) in the unit's mount namespace,
	// bind-mounting the guarded dir. A mount onto /etc/app-listener must NOT be blocked in
	// read-only mode, or the helper dies with 226/NAMESPACE and reload fails (the RO guard used to
	// deny sb_mount).
	code, out := s.exec(c, []string{"sh", "-c",
		"mkdir -p /tmp/mnt-probe && mount --bind /tmp/mnt-probe /etc/app-listener && umount /etc/app-listener && echo OK"})
	s.Require().Equalf(0, code, "a bind mount over the RO-guarded dir must be allowed (systemctl reload namespacing): %s", out)

	// 1. daemon.conf stays world-readable.
	code, out = s.exec(c, []string{"cat", "/etc/app-listener/daemon.conf"})
	s.Require().Equalf(0, code, "daemon.conf must stay world-readable: %s", out)
	s.Require().Contains(out, "[watch /protected]")

	// 2. daemon.conf is not writable by a non-app-listener process.
	code, out = s.exec(c, []string{"sh", "-c", "echo pwned >> /etc/app-listener/daemon.conf 2>&1"})
	s.Require().NotEqualf(0, code, "append to daemon.conf must be denied: %s", out)

	// 3. daemon.conf cannot be replaced by renaming a file over it.
	code, out = s.exec(c, []string{"sh", "-c", "echo pwned > /tmp/evil && mv -f /tmp/evil /etc/app-listener/daemon.conf 2>&1"})
	s.Require().NotEqualf(0, code, "rename over daemon.conf must be denied: %s", out)

	// 4. no new files in the guarded directory.
	code, out = s.exec(c, []string{"sh", "-c", "touch /etc/app-listener/newfile 2>&1"})
	s.Require().NotEqualf(0, code, "creating a file in /etc/app-listener must be denied: %s", out)

	// 5. fscrypt.key is not readable at all.
	code, out = s.exec(c, []string{"sh", "-c", "cat /etc/app-listener/fscrypt.key 2>&1"})
	s.Require().NotEqualf(0, code, "reading fscrypt.key must be denied: %s", out)

	// 6. fscrypt.key is not writable.
	code, out = s.exec(c, []string{"sh", "-c", "echo x >> /etc/app-listener/fscrypt.key 2>&1"})
	s.Require().NotEqualf(0, code, "writing fscrypt.key must be denied: %s", out)

	// 7. SIGHUP reload: the self guards detach for the config swap (so the
	// transient old+new guard count stays under the kernel's per-LSM-hook
	// program cap) and re-attach afterwards. Self-protection must survive.
	_, pidOut := s.exec(c, []string{"sh", "-c", "cat /run/app-listener-daemon.pid"})
	s.exec(c, []string{"sh", "-c", "kill -HUP " + strings.TrimSpace(pidOut)})
	reloaded := false
	for rlDeadline := time.Now().Add(45 * time.Second); time.Now().Before(rlDeadline); {
		if strings.Contains(s.readDaemonLog(c), "configuration reloaded from") {
			reloaded = true
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
	s.Require().Truef(reloaded, "SIGHUP reload did not complete, log:\n%s", s.readDaemonLog(c))

	code, out = s.exec(c, []string{"cat", "/etc/app-listener/daemon.conf"})
	s.Require().Equalf(0, code, "daemon.conf must still be readable after reload: %s", out)
	code, out = s.exec(c, []string{"sh", "-c", "echo x >> /etc/app-listener/daemon.conf 2>&1"})
	s.Require().NotEqualf(0, code, "daemon.conf write must still be denied after reload: %s", out)
	code, out = s.exec(c, []string{"sh", "-c", "cat /etc/app-listener/fscrypt.key 2>&1"})
	s.Require().NotEqualf(0, code, "fscrypt.key read must still be denied after reload: %s", out)

	events := parseDaemonEvents(s.readDaemonLog(c))
	keyDenied, dirDenied := false, false
	for _, ev := range events {
		if !ev.Denied {
			continue
		}
		if ev.Resource == "/etc/app-listener/fscrypt.key" {
			keyDenied = true
		}
		if ev.Resource == "/etc/app-listener" {
			dirDenied = true
		}
	}
	s.Require().Truef(keyDenied, "expected a DENIED event for resource=/etc/app-listener/fscrypt.key, log:\n%s", s.readDaemonLog(c))
	s.Require().Truef(dirDenied, "expected a DENIED event for resource=/etc/app-listener, log:\n%s", s.readDaemonLog(c))

	// daemon.conf never lost its content (readable while the daemon runs).
	_, conf := s.exec(c, []string{"cat", "/etc/app-listener/daemon.conf"})
	s.Require().Contains(conf, "[watch /protected]")
	s.Require().NotContains(conf, "pwned")

	// Stop the daemon; its guards unpin on shutdown. Poll until fscrypt.key
	// is reachable again, then confirm it was never modified.
	s.exec(c, []string{"sh", "-c", "pkill -TERM -f 'app-listener daemon' || true"})
	var keyLen string
	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		code, out = s.exec(c, []string{"sh", "-c", "wc -c < /etc/app-listener/fscrypt.key 2>/dev/null"})
		if code == 0 {
			keyLen = out
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	s.Require().Equal("32", strings.TrimSpace(keyLen), "fscrypt.key was modified or the guard never detached")
}

// The raw block-device gate is daemon-wide and device-granular. Two resources (/mnt/data/guardedA,
// guardedB) share one backing device (loop-mounted ext4 at /mnt/data). The gate must:
//   - be stamped ONCE ("blocking raw access to backing block device" appears once);
//   - block raw reads even for an UNGUARDED path on the device (/mnt/data/open/*): coarse by
//     design, accepted;
//   - log the denial as resource=raw-block-device, never as a watched path.
//
// Skips if no loop device is available.
func (s *IntegrationSuite) TestDaemon_RawBlockDevice_DeviceScope() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const marker = "RAW-DEVICE-DAEMON-SECRET-9C1F"

	setup := `
set -e
mkdir -p /etc/app-listener /mnt/data
dd if=/dev/zero of=/img.ext4 bs=1M count=48 status=none
mkfs.ext4 -q -F /img.ext4
for i in $(seq 0 15); do [ -e /dev/loop$i ] || mknod /dev/loop$i b 7 "$i"; done
LOOP=$(losetup -f --show /img.ext4)
mount "$LOOP" /mnt/data
mkdir -p /mnt/data/guardedA /mnt/data/guardedB /mnt/data/open
echo "` + marker + `" > /mnt/data/guardedA/secret
echo "` + marker + `" > /mnt/data/open/plain
chmod -R 755 /mnt/data
sync
printf 'LOOPDEV=%s\n' "$LOOP"
`
	code, out := s.exec(c, []string{"sh", "-c", setup})
	loopDev := ""
	for _, ln := range strings.Split(out, "\n") {
		ln = strings.TrimSpace(ln)
		if strings.HasPrefix(ln, "LOOPDEV=/dev/loop") {
			loopDev = strings.TrimPrefix(ln, "LOOPDEV=")
		}
	}
	if code != 0 || loopDev == "" {
		s.T().Skipf("loop device / ext4 setup unavailable in this environment (exit %d): %s", code, out)
	}

	// Whitelist grep (a standalone inode on this image — the uutils coreutils
	// applets, dd/cat included, share one multi-call inode, so whitelisting
	// any of them would whitelist dd and defeat the dd assertion below).
	s.startDaemon(c, `[watch /mnt/data/guardedA]
need_encryption: false
/usr/bin/grep

[watch /mnt/data/guardedB]
need_encryption: false
/usr/bin/grep`)

	// The gate is stamped exactly once for the shared device.
	_, stamp := s.exec(c, []string{"sh", "-c",
		"grep -c 'blocking raw access to backing block device' /tmp/daemon.log || true"})
	s.Require().Equalf("1", strings.TrimSpace(stamp),
		"raw block-device gate must be stamped once for the shared device, daemon log:\n%s", s.readDaemonLog(c))

	// Raw read of an UNGUARDED path on the shared device is still blocked
	// (coarse by design) — debugfs and a plain dd.
	code, out = s.exec(c, []string{"sh", "-c", fmt.Sprintf(
		"debugfs -R 'cat /open/plain' %s 2>&1 | grep -q %s && echo LEAK || echo safe", loopDev, marker)})
	s.Require().Containsf(out, "safe", "debugfs raw read of the shared device must be blocked: %s", out)

	code, out = s.exec(c, []string{"sh", "-c", fmt.Sprintf(
		"dd if=%s bs=1M count=48 2>/dev/null | tr -c '[:print:]' '\\n' | grep -q %s && echo LEAK || echo safe", loopDev, marker)})
	s.Require().Containsf(out, "safe", "dd raw read of the shared device must be blocked: %s", out)

	// The denials are attributed to the raw-device gate, never to a watched path.
	events := parseDaemonEvents(s.readDaemonLog(c))
	rawDenied := false
	for _, ev := range events {
		if !ev.Denied || (ev.Comm != "debugfs" && ev.Comm != "dd") {
			continue
		}
		rawDenied = true
		s.Require().Equalf(guard.RawDeviceResourceLabel, ev.Resource,
			"raw block-device denial must be labelled %q, got resource=%q (comm=%s path=%s)",
			guard.RawDeviceResourceLabel, ev.Resource, ev.Comm, ev.Path)
	}
	s.Require().Truef(rawDenied, "expected a DENIED debugfs/dd raw block-device event, daemon log:\n%s", s.readDaemonLog(c))

	s.exec(c, []string{"sh", "-c",
		fmt.Sprintf("pkill -f 'app-listener daemon' || true; umount /mnt/data 2>/dev/null; losetup -d %s 2>/dev/null; true", loopDev)})
}

// Edit-protected live mode (issue #40). With a password configured the daemon exposes a control
// socket that must:
//   - reject a non-root peer (SO_PEERCRED);
//   - reject a wrong password (lockout after enough failures is not asserted, to keep the test
//     quick);
//   - on a correct password from a root app-listener peer, briefly widen the target resource's
//     guard so `edit-protected --put` can write, then narrow it (GRANTED / REVOKED both logged);
//   - allow only one live session at a time.
func (s *IntegrationSuite) TestDaemon_EditProtected_LiveMode() {
	const password = "Sup3r-Secret-99"

	hash, err := editprotected.Hash(password, editprotected.OriginInstall)
	s.Require().NoError(err)

	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /etc/app-listener /outside && echo SECRET > /protected/secret && echo ORIGINAL > /outside/victim && " +
			"chmod 755 /protected && chmod 600 /protected/secret && " +
			// a symlink placed inside the tree BEFORE the daemon guards it —
			// --put must refuse it, not follow it out of the guarded tree.
			"ln -s /outside/victim /protected/escape && " +
			"printf '%s\\n' " + shQuote(hash) + " > /etc/app-listener/edit-auth.hash && chmod 600 /etc/app-listener/edit-auth.hash"})

	s.startDaemon(c, `[watch /protected]
need_encryption: false
/usr/bin/sleep`)

	// The control socket comes up during startup.
	socketReady := false
	for dl := time.Now().Add(20 * time.Second); time.Now().Before(dl); {
		if strings.Contains(s.readDaemonLog(c), "edit-protected control socket ready") {
			socketReady = true
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
	s.Require().Truef(socketReady, "control socket never came up, log:\n%s", s.readDaemonLog(c))

	s.Require().NoError(c.CopyFileToContainer(s.ctx,
		absPath("./exploits/edit_auth_bypass"), "/exploits/edit_auth_bypass", 0o755))

	const sock = "/run/app-listener-daemon.control"
	const nobody = "setpriv --reuid=65534 --regid=65534 --clear-groups"
	const put = "/app-listener edit-protected --resource /protected --put live.txt"

	// 1. a peer that is not the app-listener binary is rejected before AUTH
	//    is even considered — it never learns a single protected path.
	code, out := s.exec(c, []string{"sh", "-c",
		"/exploits/edit_auth_bypass " + sock + " " + password + " 2>&1"})
	s.Require().NotEqualf(0, code, "a non-app-listener peer must be refused: %s", out)
	s.Require().Containsf(out, "not the installed app-listener binary", "expected the peer-exe rejection, got: %s", out)
	s.Require().NotContainsf(out, "/protected", "an unauthenticated peer must not see any protected path: %s", out)

	// 2. non-root is refused before it even connects.
	code, out = s.exec(c, []string{"sh", "-c",
		"echo x | " + nobody + " " + put + " 2>&1"})
	s.Require().NotEqualf(0, code, "non-root must be refused: %s", out)
	s.Require().Containsf(out, "must be run as root", "expected the root check, got: %s", out)

	// 3. wrong password -> refused, nothing written.
	code, out = s.exec(c, []string{"sh", "-c",
		"echo x | APP_LISTENER_EDIT_PASSWORD=wrong-Pass-1234 " + put + " 2>&1"})
	s.Require().NotEqualf(0, code, "a wrong password must be refused: %s", out)
	s.Require().Containsf(out, "authentication failed", "expected an auth failure, got: %s", out)
	code, _ = s.exec(c, []string{"sh", "-c", "test -e /protected/live.txt"})
	s.Require().NotEqualf(0, code, "no file may be written on a failed auth")

	// 4. correct password -> the grant works and the file is written.
	code, out = s.exec(c, []string{"sh", "-c",
		"echo LIVE-EDIT | APP_LISTENER_EDIT_PASSWORD=" + password + " " + put + " 2>&1"})
	s.Require().Equalf(0, code, "live --put must succeed: %s", out)
	_, content := s.exec(c, []string{"sh", "-c", "cat /protected/live.txt"})
	s.Require().Equal("LIVE-EDIT", strings.TrimSpace(content))

	// 5. the grant is transient: GRANTED and REVOKED both logged.
	log := s.readDaemonLog(c)
	s.Require().Containsf(log, "write access GRANTED", "expected a GRANTED log line, got:\n%s", log)
	s.Require().Containsf(log, "write access REVOKED", "expected a REVOKED log line, got:\n%s", log)

	// 6. between sessions a non-authenticated write to the tree is still denied.
	code, out = s.exec(c, []string{"sh", "-c", "echo pwned > /protected/pwned 2>&1"})
	s.Require().NotEqualf(0, code, "an unauthenticated write to the guarded tree must be denied: %s", out)

	// 7. --put must not follow a symlink out of the guarded tree.
	code, out = s.exec(c, []string{"sh", "-c",
		"echo PWNED | APP_LISTENER_EDIT_PASSWORD=" + password +
			" /app-listener edit-protected --resource /protected --put escape 2>&1"})
	s.Require().NotEqualf(0, code, "--put through a symlink must be refused: %s", out)
	s.Require().Containsf(out, "symlink", "expected a symlink refusal, got: %s", out)
	_, victim := s.exec(c, []string{"sh", "-c", "cat /outside/victim"})
	s.Require().Equal("ORIGINAL", strings.TrimSpace(victim), "the symlink target outside the tree was written")

	s.exec(c, []string{"sh", "-c", "pkill -TERM -f 'app-listener daemon' || true"})
}

// shQuote single-quotes s for safe embedding in a /bin/sh -c string.
func shQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// TOCTOU races across the daemon lifecycle (start, SIGHUP reload, graceful stop, SIGKILL) for both
// a kernel-fscrypt directory and a file-vault (single regular file) resource.
//
// Invariant (internal/usecase/daemon.go): a resource is never readable in plaintext by an
// unauthorized process at any observable instant, however an attacker's read races start, reload,
// stop or a hard kill. Ordering is attach -> unlock -> populate -> resolve -> re-sync; a SIGKILL
// can't be caught, so the rest rests on the guard's LSM links staying pinned to /sys/fs/bpf and
// `daemon --lockdown` (systemd ExecStopPost) recovering from ANY point, including mid-transform of
// a file-vault (recovery sidecar, fscrypt/filevault.go).
//
// The tests emulate an attacker rather than reason about code: a background goroutine hammers the
// resource as an unprivileged user through the transition and fails on any plaintext. Real
// Vault.Encrypt/IsEncrypted set up and inspect resources via a root-capable harness compiled from
// internal/fscrypt (harness_test.go, main_test.go's fscryptTestAmd64Bin, like guardTestAmd64Bin).
//
// The directory case needs a real fscrypt-capable filesystem (ext4 with `encrypt`,
// CONFIG_FS_ENCRYPTION), which overlayfs can't provide: setupFscryptDirFilesystem builds one on a
// loop device and skips if unavailable (like TestDaemon_RawBlockDevice_DeviceScope). The file-vault
// case needs only the master key and always runs.

// copyFscryptHarness copies the fscrypt root-capable test harness into the
// container at /fscrypt.harness (see internal/fscrypt/harness_test.go).
func (s *IntegrationSuite) copyFscryptHarness(c testcontainers.Container) {
	s.Require().NoError(c.CopyFileToContainer(s.ctx, fscryptTestAmd64Bin, "/fscrypt.harness", 0755))
}

// harnessRun execs one subtest of the fscrypt harness binary against path,
// asserting it exits 0, and returns its stdout+stderr.
func (s *IntegrationSuite) harnessRun(c testcontainers.Container, subtest, path string) string {
	code, out := s.harnessProbe(c, subtest, path)
	s.Require().Equalf(0, code, "fscrypt harness %s(%s) failed: %s", subtest, path, out)
	return out
}

// harnessProbe execs one harness subtest against path without asserting success. The harness is
// never whitelisted on a guarded resource here, so a still-enforcing guard (link may stay pinned
// for a window after SIGKILL, by design) correctly denies it: that denial proves the resource is
// NOT exposed, never a failure. Use this, not harnessRun, for best-effort diagnostics.
func (s *IntegrationSuite) harnessProbe(c testcontainers.Container, subtest, path string) (int, string) {
	return s.exec(c, []string{"sh", "-c", fmt.Sprintf(
		"APPLISTENER_HARNESS_PATH=%s /fscrypt.harness -test.run %s -test.v",
		shQuote(path), shQuote("^TestFscryptHarness/"+subtest+"$"))})
}

// harnessMigrate seals path into its encrypted-at-rest form (kernel fscrypt policy for a directory,
// with filesystem prerequisites applied first; file-vault format for a regular file) via the real
// Vault.Encrypt: exactly what `install` does before the daemon attaches a guard.
func (s *IntegrationSuite) harnessMigrate(c testcontainers.Container, path string) {
	s.harnessRun(c, "TestMigrate", path)
}

// assertFileVaultSealed confirms, via the whitelisted grep (not the harness), that a file-vault
// path's raw bytes no longer expose marker: the check that `daemon --lockdown` did its job after a
// SIGKILL. The harness is never whitelisted, and the guard's pinned link may still enforce right
// now (lockdown widens self-access to relock but doesn't unpin), so a harness call would be denied
// whether or not the file is sealed (the false failure assertNoUnauthorizedPlaintext also avoids).
// AEAD ciphertext won't reproduce the marker, so grep finding nothing is real proof and finding it
// is a genuine lockdown failure.
func (s *IntegrationSuite) assertFileVaultSealed(c testcontainers.Container, path, marker, context string) {
	code, out := s.exec(c, []string{"grep", "-c", marker, path})
	s.Require().Falsef(code == 0 && strings.Contains(out, "1"),
		"%s: %s still exposes its plaintext marker after lockdown — not sealed", context, path)
}

// harnessIsEncrypted reports path's on-disk encryption state via the real Vault.IsEncrypted,
// mutating nothing. Only safe when no guard can still be pinned-and-enforcing against the harness:
// before any daemon starts or after a clean unraced shutdown (which unpins). After a SIGKILL use
// assertFileVaultSealed.
func (s *IntegrationSuite) harnessIsEncrypted(c testcontainers.Container, path string) bool {
	return strings.Contains(s.harnessRun(c, "TestIsEncrypted", path), "ENCRYPTED")
}

// resetFileVaultTarget clears any leftover backup/recovery sidecar from a previous kill-loop round
// (safe: no guard watches yet; filevault.go's "only rewrite in place" rule is for a live guard) and
// rewrites path as fresh plaintext with marker, ready for harnessMigrate.
func (s *IntegrationSuite) resetFileVaultTarget(c testcontainers.Container, path, marker string) {
	cmd := fmt.Sprintf("rm -f %s %s.app_listener.backup %s.app_listener.recover && printf '%%s' %s > %s && chmod 600 %s",
		path, path, path, shQuote(marker), path, path)
	code, out := s.exec(c, []string{"sh", "-c", cmd})
	s.Require().Equalf(0, code, "resetting file-vault target %s: %s", path, out)
}

// setupFscryptDirFilesystem loop-mounts a small ext4 with the `encrypt` feature at mountPoint so a
// directory can carry a real kernel fscrypt policy (unlike the overlayfs root). Skips the test if
// no working loop device / encrypt-capable ext4 (older e2fsprogs, no CONFIG_FS_ENCRYPTION).
func (s *IntegrationSuite) setupFscryptDirFilesystem(c testcontainers.Container, mountPoint string) {
	setup := fmt.Sprintf(`
set -e
mkdir -p %s
dd if=/dev/zero of=/fscrypt.img bs=1M count=64 status=none
mkfs.ext4 -q -F -O encrypt /fscrypt.img
for i in $(seq 0 15); do [ -e /dev/loop$i ] || mknod /dev/loop$i b 7 "$i"; done
LOOP=$(losetup -f --show /fscrypt.img)
mount "$LOOP" %s
chmod 755 %s
printf 'LOOPDEV=%%s\n' "$LOOP"
`, mountPoint, mountPoint, mountPoint)
	code, out := s.exec(c, []string{"sh", "-c", setup})
	loopDev := ""
	for _, ln := range strings.Split(out, "\n") {
		ln = strings.TrimSpace(ln)
		if strings.HasPrefix(ln, "LOOPDEV=/dev/loop") {
			loopDev = strings.TrimPrefix(ln, "LOOPDEV=")
		}
	}
	if code != 0 || loopDev == "" {
		s.T().Skipf("encrypt-capable ext4 loop filesystem unavailable in this environment (exit %d): %s", code, out)
	}
}

// sigDaemon sends signal sig (e.g. "TERM", "KILL", "HUP") to the daemon
// process by comm match.
func (s *IntegrationSuite) sigDaemon(c testcontainers.Container, sig string) {
	s.exec(c, []string{"sh", "-c", "pkill -" + sig + " -f 'app-listener daemon' || true"})
}

// daemonShutdownTimeout is the budget for a graceful SIGTERM (Stop locks every vault back first) or
// a kill racing one. Each guard (re)attach is a verifier pass over 23 LSM hooks, 5-11s on slower
// hosts (hardened kernels), and a reload or live-edit can rebuild several guards (incl. the two
// self-protection ones) back to back, so it covers a handful in sequence.
const daemonShutdownTimeout = 45 * time.Second

// awaitDaemonDead polls until no app-listener process remains (false after timeout). Uses
// noLiveAppListenerProcs (pool_test.go), matching `comm`, not `pgrep -f 'app-listener daemon'`:
// under `sh -c "pgrep -f ..."` the pattern is in the wrapping shell's own cmdline, so `pgrep -f`
// matches it and reports "alive" forever.
func (s *IntegrationSuite) awaitDaemonDead(c testcontainers.Container, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		code, _ := s.exec(c, []string{"sh", "-c", noLiveAppListenerProcs})
		if code == 0 {
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}

// runLockdown execs `daemon --lockdown` (ExecStopPost safety net): force-locks every encryption
// root in daemon.conf and exits. It always exits 0 (best-effort per root), so callers check the
// real state via harnessIsEncrypted, not the exit code.
func (s *IntegrationSuite) runLockdown(c testcontainers.Container) (int, string) {
	return s.exec(c, []string{"/app-listener", "daemon", "--lockdown"})
}

// rawExec is s.exec without fatal assertions on transport errors: safe from a background goroutine
// (testify Require/FailNow off the test goroutine is undefined). A transport hiccup counts as
// "nothing observed", right for a tight polling racer.
func rawExec(ctx context.Context, c testcontainers.Container, cmd []string) (int, string) {
	code, reader, err := c.Exec(ctx, cmd, tcexec.Multiplexed())
	if err != nil {
		return -1, ""
	}
	b, err := io.ReadAll(reader)
	if err != nil {
		return -1, ""
	}
	return code, string(b)
}

// raceUnauthorizedReader repeatedly tries, as an unprivileged user, to read path until stop fires,
// setting leaked when an attempt succeeds (exit 0) and returns content containing marker: genuine
// plaintext exposure. A still-locked resource never matches, and a properly guarded unlocked one is
// denied outright, so it fires only on the failure these tests exist to catch.
func (s *IntegrationSuite) raceUnauthorizedReader(c testcontainers.Container, path, marker string, stop <-chan struct{}, leaked *atomic.Bool) {
	const nobody = "setpriv --reuid=65534 --regid=65534 --clear-groups"
	for {
		select {
		case <-stop:
			return
		default:
		}
		_, out := rawExec(s.ctx, c, []string{"sh", "-c", nobody + " cat " + path + " 2>/dev/null"})
		if strings.Contains(out, marker) {
			leaked.Store(true)
			return
		}
	}
}

// assertNoUnauthorizedPlaintext runs one direct check that an unprivileged reader can't currently
// see marker at path, even if the guard isn't attached right now (e.g. after SIGKILL): the resource
// must be safe by construction (ciphertext, or ciphertext-but-orphan-guarded), not "safe while the
// guard process is alive".
func (s *IntegrationSuite) assertNoUnauthorizedPlaintext(c testcontainers.Container, path, marker, context string) {
	// Best-effort only: the guard's pinned link may still enforce right now (a SIGKILL before
	// Stop() unpins leaves it orphaned-but-active until `daemon --lockdown`), in which case the
	// never-whitelisted harness is correctly denied. That is stronger evidence of "not exposed"
	// than any report, so never a hard failure; only the unprivileged-reader check below is the
	// real assertion.
	encStatus := "unknown (harness access denied by guard, itself a safe outcome)"
	if code, out := s.harnessProbe(c, "TestIsEncrypted", path); code == 0 {
		switch {
		case strings.Contains(out, "ENCRYPTED"):
			encStatus = "true"
		case strings.Contains(out, "PLAINTEXT"):
			encStatus = "false"
		}
	}
	code, out := s.exec(c, []string{"sh", "-c",
		"setpriv --reuid=65534 --regid=65534 --clear-groups cat " + path + " 2>/dev/null"})
	leaked := code == 0 && strings.Contains(out, marker)
	s.Require().Falsef(leaked,
		"%s: %s ended up plaintext (encrypted=%s) AND readable by an unauthorized user — TOCTOU window", context, path, encStatus)
}

// ---------------------------------------------------------------
// A. Startup race: an attacker hammers the resource from the moment the
// daemon process is launched until its guards are confirmed up.
// ---------------------------------------------------------------

// Races an unprivileged reader against startup on a real kernel-fscrypt-encrypted DIRECTORY: attach
// -> unlock -> populate must never leave the key provisioned with the guard not yet enforcing.
func (s *IntegrationSuite) TestDaemon_StartupRace_NoPlaintextWindow_Directory() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.copyFscryptHarness(c)

	const marker = "TOCTOU-DIR-STARTUP-SECRET-91AF"
	const mnt = "/mnt/fscrypt-dir"
	s.setupFscryptDirFilesystem(c, mnt) // skips the test if unsupported here

	secretDir := mnt + "/vault"
	secretFile := secretDir + "/secret.txt"
	s.exec(c, []string{"sh", "-c", fmt.Sprintf(
		"mkdir -p %s /etc/app-listener && printf '%s' > %s && chmod 755 %s && chmod 644 %s && "+
			"head -c 32 /dev/zero > /etc/app-listener/fscrypt.key && chmod 600 /etc/app-listener/fscrypt.key",
		secretDir, marker, secretFile, secretDir, secretFile)})

	s.harnessMigrate(c, secretDir)
	s.Require().True(s.harnessIsEncrypted(c, secretDir), "setup: the directory must be encrypted before the daemon ever starts")

	config := fmt.Sprintf("[watch %s]\nneed_encryption: true\n/usr/bin/grep", secretDir)
	s.exec(c, []string{"sh", "-c", fmt.Sprintf("cat > /etc/app-listener/daemon.conf <<'EOF'\n%s\nEOF", config)})

	stop := make(chan struct{})
	var leaked atomic.Bool
	go s.raceUnauthorizedReader(c, secretFile, marker, stop, &leaked)

	s.launchDaemon(c)
	s.awaitDaemonUp(c, config)
	close(stop)

	s.Require().Falsef(leaked.Load(),
		"an unauthorized reader observed plaintext content of %s before the guard was fully attached", secretFile)

	// Sanity: the daemon actually did its job — unlocked for its
	// whitelisted reader, still denying everyone else now that it is up.
	_, grepOut := s.exec(c, []string{"grep", "-c", marker, secretFile})
	s.Require().Contains(grepOut, "1", "whitelisted reader must see the unlocked content")
	code, _ := s.exec(c, []string{"sh", "-c", "setpriv --reuid=65534 --regid=65534 --clear-groups cat " + secretFile})
	s.Require().NotEqualf(0, code, "non-whitelisted reader must still be denied once the guard is fully up")

	s.sigDaemon(c, "TERM")
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not exit after SIGTERM")
	s.Require().Truef(s.harnessIsEncrypted(c, secretDir), "a graceful stop must re-lock the directory")

	s.exec(c, []string{"sh", "-c", fmt.Sprintf("umount %s 2>/dev/null; losetup -D 2>/dev/null; true", mnt)})
}

// TestDaemon_StartupRace_NoPlaintextWindow_File is
// TestDaemon_StartupRace_NoPlaintextWindow_Directory for a single-file
// watch root sealed with the userspace file-vault (filevault.go) instead of
// a kernel fscrypt policy — needs no special filesystem, so it always runs.
func (s *IntegrationSuite) TestDaemon_StartupRace_NoPlaintextWindow_File() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.copyFscryptHarness(c)

	const marker = "TOCTOU-FILE-STARTUP-SECRET-4D2E"
	const secretFile = "/protected/secret.txt"

	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /etc/app-listener && head -c 32 /dev/zero > /etc/app-listener/fscrypt.key && chmod 600 /etc/app-listener/fscrypt.key"})
	s.resetFileVaultTarget(c, secretFile, marker)
	s.exec(c, []string{"chmod", "644", secretFile})

	s.harnessMigrate(c, secretFile)
	s.Require().True(s.harnessIsEncrypted(c, secretFile), "setup: the file must be sealed before the daemon ever starts")

	config := fmt.Sprintf("[watch %s]\nneed_encryption: true\n/usr/bin/grep", secretFile)
	s.exec(c, []string{"sh", "-c", fmt.Sprintf("cat > /etc/app-listener/daemon.conf <<'EOF'\n%s\nEOF", config)})

	stop := make(chan struct{})
	var leaked atomic.Bool
	go s.raceUnauthorizedReader(c, secretFile, marker, stop, &leaked)

	s.launchDaemon(c)
	s.awaitDaemonUp(c, config)
	close(stop)

	s.Require().Falsef(leaked.Load(),
		"an unauthorized reader observed plaintext content of %s before the guard was fully attached", secretFile)

	_, grepOut := s.exec(c, []string{"grep", "-c", marker, secretFile})
	s.Require().Contains(grepOut, "1", "whitelisted reader must see the unlocked content")
	code, _ := s.exec(c, []string{"sh", "-c", "setpriv --reuid=65534 --regid=65534 --clear-groups cat " + secretFile})
	s.Require().NotEqualf(0, code, "non-whitelisted reader must still be denied once the guard is fully up")

	s.sigDaemon(c, "TERM")
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not exit after SIGTERM")
	s.Require().Truef(s.harnessIsEncrypted(c, secretFile), "a graceful stop must re-seal the file")
}

// Regression for the startup fan-out (buildGuards attaching concurrently, then startGuards
// unlocking roots and preparing guards concurrently; internal/usecase/daemon.go,
// cmd/functions/daemon/daemon.go): TWO `watch:` sub-paths share ONE encryption root
// (browser-profile shape: Local Storage + Cookies under one vault). The root unlocks once for the
// group (uniqueEncryptionRoots), making BOTH sub-paths' plaintext readable at once, so EACH
// sub-path's guard must already be attached or the lagging sibling is exposed. A naive per-resource
// pipeline would reintroduce this window; the concurrent version keeps the invariant because
// buildGuards attaches EVERY guard (all group members) as one completed phase before
// startGuards.Start() calls Unlock, and unlockRoots fans out only over deduplicated roots (never
// two concurrent unlocks of the shared root).
func (s *IntegrationSuite) TestDaemon_StartupRace_NoPlaintextWindow_GroupedResources() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.copyFscryptHarness(c)

	const markerA = "TOCTOU-GROUP-SUBA-SECRET-7F3C"
	const markerB = "TOCTOU-GROUP-SUBB-SECRET-9E1D"
	const mnt = "/mnt/fscrypt-group"
	s.setupFscryptDirFilesystem(c, mnt) // skips the test if unsupported here

	root := mnt + "/profile"
	// No space in these names, unlike the real catalog's "Local Storage":
	// these paths are interpolated unquoted into raw `sh -c` command strings
	// below, and a space would word-split into two separate shell arguments.
	subA := root + "/local-storage"
	subB := root + "/cookies"
	secretA := subA + "/secret.txt"
	secretB := subB + "/secret.txt"
	s.exec(c, []string{"sh", "-c", fmt.Sprintf(
		"mkdir -p %s %s /etc/app-listener && "+
			"printf '%s' > %s && printf '%s' > %s && "+
			"chmod 755 %s %s %s && chmod 644 %s %s && "+
			"head -c 32 /dev/zero > /etc/app-listener/fscrypt.key && chmod 600 /etc/app-listener/fscrypt.key",
		subA, subB, markerA, secretA, markerB, secretB, root, subA, subB, secretA, secretB)})

	s.harnessMigrate(c, root)
	s.Require().True(s.harnessIsEncrypted(c, root), "setup: the shared vault root must be encrypted before the daemon ever starts")

	config := fmt.Sprintf("[watch %s]\nwatch: %s\nwatch: %s\nneed_encryption: true\n/usr/bin/grep", root, subA, subB)
	s.exec(c, []string{"sh", "-c", fmt.Sprintf("cat > /etc/app-listener/daemon.conf <<'EOF'\n%s\nEOF", config)})

	stop := make(chan struct{})
	var leakedA, leakedB atomic.Bool
	go s.raceUnauthorizedReader(c, secretA, markerA, stop, &leakedA)
	go s.raceUnauthorizedReader(c, secretB, markerB, stop, &leakedB)

	s.launchDaemon(c)
	// watchPathsInConfig treats a `[watch <path>]` header as the guarded resource, wrong for a
	// group: the header names the ENCRYPTION ROOT, not itself guarded (TestLoadWatchGroup), so
	// awaitDaemonUp's per-path "guard started" markers never appear. Passing "" falls back to "pid
	// file exists", still exact: writePidFile runs only after startGuardedDaemon's d.Start()
	// returns, i.e. after EVERY guard (group members included) finished the full pipeline.
	s.awaitDaemonUp(c, "")
	close(stop)

	s.Require().Falsef(leakedA.Load(),
		"an unauthorized reader observed plaintext content of %s (shared vault root %s) before its guard was fully attached", secretA, root)
	s.Require().Falsef(leakedB.Load(),
		"an unauthorized reader observed plaintext content of %s (shared vault root %s) before its guard was fully attached", secretB, root)

	// Sanity: the daemon actually did its job on BOTH group members —
	// unlocked for the whitelisted reader, still denying everyone else.
	for _, pair := range []struct{ path, marker string }{{secretA, markerA}, {secretB, markerB}} {
		_, grepOut := s.exec(c, []string{"grep", "-c", pair.marker, pair.path})
		s.Require().Containsf(grepOut, "1", "whitelisted reader must see the unlocked content of %s", pair.path)
		code, _ := s.exec(c, []string{"sh", "-c", "setpriv --reuid=65534 --regid=65534 --clear-groups cat " + pair.path})
		s.Require().NotEqualf(0, code, "non-whitelisted reader must still be denied on %s once the guard is fully up", pair.path)
	}

	s.sigDaemon(c, "TERM")
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not exit after SIGTERM")
	s.Require().Truef(s.harnessIsEncrypted(c, root), "a graceful stop must re-lock the shared vault root")

	s.exec(c, []string{"sh", "-c", fmt.Sprintf("umount %s 2>/dev/null; losetup -D 2>/dev/null; true", mnt)})
}

// C. SIGKILL inside startup's unlock/populate window at several delays, for the file-vault
// resource: the newest, most crash-sensitive path (in-place AEAD transform + recovery sidecar).

// SIGKILLs the daemon at increasing delays after launch (early enough to plausibly land inside
// attach/unlock/populate) and asserts that whatever state the kill caught, the file is never both
// plaintext and reachable by an unauthorized reader; `daemon --lockdown` always recovers; and a
// clean restart round-trips the content uncorrupted (the recovery sidecar works on a real
// interrupted transform).
func (s *IntegrationSuite) TestDaemon_KillDuringUnlock_FileVault_NeverOrphansPlaintext() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.copyFscryptHarness(c)

	const marker = "TOCTOU-FILE-KILL-SECRET-2C6D"
	const secretFile = "/protected/secret.txt"
	config := fmt.Sprintf("[watch %s]\nneed_encryption: true\n/usr/bin/grep", secretFile)

	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /etc/app-listener && head -c 32 /dev/zero > /etc/app-listener/fscrypt.key && chmod 600 /etc/app-listener/fscrypt.key"})
	s.exec(c, []string{"sh", "-c", fmt.Sprintf("cat > /etc/app-listener/daemon.conf <<'EOF'\n%s\nEOF", config)})

	killDelays := []time.Duration{80 * time.Millisecond, 250 * time.Millisecond, 600 * time.Millisecond, 1200 * time.Millisecond}

	for i, delay := range killDelays {
		s.Run(fmt.Sprintf("kill_after_%s", delay), func() {
			s.resetFileVaultTarget(c, secretFile, marker)
			s.harnessMigrate(c, secretFile)
			s.Require().True(s.harnessIsEncrypted(c, secretFile), "round %d setup: must start from a sealed file", i)

			s.launchDaemon(c)
			time.Sleep(delay)
			s.sigDaemon(c, "KILL")
			s.Require().Truef(s.awaitDaemonDead(c, 10*time.Second), "round %d: daemon did not die after SIGKILL", i)

			s.assertNoUnauthorizedPlaintext(c, secretFile, marker, fmt.Sprintf("round %d (kill after %s)", i, delay))

			// ExecStopPost's safety net must always be able to fully
			// re-lock, regardless of exactly when the kill landed —
			// including mid in-place transform (recoverFileInPlace
			// consuming the staged recovery sidecar).
			code, out := s.runLockdown(c)
			s.Require().Equalf(0, code, "round %d: daemon --lockdown exited non-zero: %s", i, out)
			s.assertFileVaultSealed(c, secretFile, marker, fmt.Sprintf("round %d: lockdown", i))

			// A clean restart must round-trip the content without
			// corruption: AES-GCM authentication would hard-fail loudly on
			// any corruption, so a matching grep here is real proof, not a
			// coincidence.
			s.launchDaemon(c)
			s.awaitDaemonUp(c, config)
			_, grepOut := s.exec(c, []string{"grep", "-c", marker, secretFile})
			s.Require().Contains(grepOut, "1", "round %d: content must round-trip intact through the kill + lockdown + restart cycle", i)

			s.sigDaemon(c, "TERM")
			s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "round %d: daemon did not exit after SIGTERM", i)
			s.runLockdown(c)
		})
	}
}

// D. Live catalog-refresh SIGHUP reload (whitelist change on an existing resource): the ordinary
// reload race window, and a kill landing mid-reload.

// installFakeSystemctl stubs /usr/local/bin/systemctl in the container. The images run no systemd
// (the daemon is launched via nohup) but install's LIVE refresh path
// (internal/systemd.IsDaemonActive / EnableAndVerify) shells out to systemctl. The stub covers what
// that path calls: is-active/is-enabled report the daemon up (live gate passes, no
// enable/start/restart branch), daemon-reload is a no-op, and reload delivers a real SIGHUP to the
// daemon via its pid file, making the reload real and worth racing.
func (s *IntegrationSuite) installFakeSystemctl(c testcontainers.Container) {
	script := "#!/bin/sh\n" +
		"case \"$1\" in\n" +
		"  is-active) echo active; exit 0;;\n" +
		"  is-enabled) echo enabled; exit 0;;\n" +
		"  daemon-reload) exit 0;;\n" +
		"  reload) kill -HUP \"$(cat /run/app-listener-daemon.pid 2>/dev/null)\" 2>/dev/null; exit 0;;\n" +
		"  *) exit 0;;\n" +
		"esac\n"
	code, out := s.exec(c, []string{"sh", "-c",
		fmt.Sprintf("cat > /usr/local/bin/systemctl <<'EOF'\n%sEOF\nchmod +x /usr/local/bin/systemctl", script)})
	s.Require().Equalf(0, code, "installing fake systemctl stub: %s", out)
}

// Exercises Reload's documented ordering (new guard attaches before the old detaches) for a
// resource whose WHITELIST changes via a live SIGHUP reload: with a live racer through a normal
// reload, then with a SIGKILL mid-reload. The resource must never end up fully unguarded.
//
// Trigger: `install --update-catalog-only --live --yes` on the built-in WireGuard entry
// (/etc/wireguard, whitelisting /usr/bin/nmcli), not a raw shell rewrite of daemon.conf: the
// /etc/app-listener self-guard (issue #53 follow-up) denies any non-daemon-binary write there, so a
// NEW [watch] section can't be hand-added while the daemon runs (and adding a genuinely new
// resource always stops the daemon first, per README.md). `--update-catalog-only --live` is the one
// documented live path: it re-expands an EXISTING section's whitelist from the daemon's own binary
// (GUARD_ALLOW_ROOT, same exe inode) and delivers via SIGHUP, with no external process writing the
// guarded config.
func (s *IntegrationSuite) TestDaemon_KillDuringReload_NoUnprotectedWindow() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.installFakeSystemctl(c)

	const marker = "TOCTOU-RELOAD-SECRET-7F31"
	const resourceB = "/etc/wireguard/secret.conf"
	const nmcli = "/usr/bin/nmcli"

	writeMarker := func() {
		s.exec(c, []string{"sh", "-c", fmt.Sprintf("mkdir -p /etc/wireguard && printf '%%s' '%s' > %s", marker, resourceB)})
	}
	// installNmcli simulates the real trigger: a package manager installs a catalog-whitelisted
	// binary after the daemon started guarding the resource with an empty whitelist
	// (FilterExistingWhitelist only picks up binaries on disk at scan time). It must be a real
	// standalone ELF, not a shebang script (identity is the CALLING PROCESS's exe inode, a script
	// runs under its interpreter's /bin/sh, so it would always be denied). /bin/cat and friends
	// won't do either on modern Ubuntu: they're symlinks into one uutils multicall binary that
	// dispatches on its own RESOLVED path's basename (not argv[0]; `exec -a` doesn't help), so a
	// copy at "nmcli" is an unknown applet. /usr/bin/grep is GNU grep, a standalone binary
	// indifferent to its name: a copy works as the stand-in reader.
	installNmcli := func() {
		s.exec(c, []string{"sh", "-c", fmt.Sprintf("cp /usr/bin/grep %s && chmod +x %s", nmcli, nmcli)})
	}
	removeNmcli := func() { s.exec(c, []string{"rm", "-f", nmcli}) }

	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /etc/app-listener && echo BASE > /protected/base.txt && chmod 755 /protected && " +
			"head -c 32 /dev/zero > /etc/app-listener/fscrypt.key && chmod 600 /etc/app-listener/fscrypt.key"})
	writeMarker()

	// /etc/wireguard starts with an empty whitelist (nmcli not installed
	// yet): the WireGuard catalog entry still matches the section by path,
	// so a later live refresh can re-expand it.
	baseConfig := "[watch /protected]\nneed_encryption: false\n/usr/bin/grep\n\n[watch /etc/wireguard]\nneed_encryption: false"

	s.startDaemon(c, baseConfig)

	// ---- Round 1: a normal, uninterrupted live catalog refresh must never
	// open a window where /etc/wireguard is fully unguarded. ----
	installNmcli()

	stop := make(chan struct{})
	var leaked atomic.Bool
	go s.raceUnauthorizedReader(c, resourceB, marker, stop, &leaked)

	code, out := s.exec(c, []string{"/app-listener", "install", "--update-catalog-only", "--live", "--yes"})
	s.Require().Equalf(0, code, "live catalog refresh failed: %s", out)

	reloaded := false
	for dl := time.Now().Add(30 * time.Second); time.Now().Before(dl); {
		if strings.Contains(s.readDaemonLog(c), "configuration reloaded") {
			reloaded = true
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	close(stop)
	s.Require().Truef(reloaded, "reload never completed, log:\n%s", s.readDaemonLog(c))
	s.Require().Falsef(leaked.Load(),
		"an unauthorized reader observed %s during the live catalog refresh", resourceB)

	code, out = s.exec(c, []string{nmcli, marker, resourceB})
	s.Require().Equalf(0, code, "the newly whitelisted binary must be able to read %s after the refresh: %s", resourceB, out)
	s.Require().Contains(out, marker)

	code, _ = s.exec(c, []string{"sh", "-c", "setpriv --reuid=65534 --regid=65534 --clear-groups cat " + resourceB})
	s.Require().NotEqualf(0, code, "a non-whitelisted reader must still be denied after the refresh")

	s.sigDaemon(c, "TERM")
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not exit after SIGTERM")
	s.runLockdown(c)

	// ---- Round 2: a live catalog refresh interrupted by SIGKILL must still
	// never leave /etc/wireguard fully unguarded. ----
	removeNmcli() // back to "not installed yet"
	writeMarker()
	s.startDaemon(c, baseConfig) // daemon.conf gets rewritten to baseConfig here, discarding round 1's patched whitelist
	installNmcli()

	code, out = s.exec(c, []string{"/app-listener", "install", "--update-catalog-only", "--live", "--yes"})
	s.Require().Equalf(0, code, "live catalog refresh (round 2 trigger) failed: %s", out)
	time.Sleep(150 * time.Millisecond) // aim somewhere inside the guard swap (new attaches, old detaches)
	s.sigDaemon(c, "KILL")
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not die after SIGKILL mid-reload")

	code, out = s.exec(c, []string{"sh", "-c", "setpriv --reuid=65534 --regid=65534 --clear-groups cat " + resourceB})
	s.Require().NotEqualf(0, code, "a non-whitelisted reader must still be denied after a kill mid-reload: %s", out)
}

// E. Shutdown interrupted mid-lockdown: SIGTERM (starts the secure lockdown: lock before detach)
// immediately followed by SIGKILL (systemd stop-timeout or an impatient admin).

// Races a SIGKILL against the graceful-shutdown lockdown (Stop: lock vault before guard detach) and
// verifies `daemon --lockdown` still finishes the job, with content never both plaintext and
// unguarded and no silent corruption on the next start.
func (s *IntegrationSuite) TestDaemon_KillDuringShutdown_LockdownRecovers() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.copyFscryptHarness(c)

	const marker = "TOCTOU-SHUTDOWN-SECRET-5A19"
	const secretFile = "/protected/secret.txt"
	config := fmt.Sprintf("[watch %s]\nneed_encryption: true\n/usr/bin/grep", secretFile)

	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /etc/app-listener && head -c 32 /dev/zero > /etc/app-listener/fscrypt.key && chmod 600 /etc/app-listener/fscrypt.key"})
	s.resetFileVaultTarget(c, secretFile, marker)
	s.harnessMigrate(c, secretFile)
	s.Require().True(s.harnessIsEncrypted(c, secretFile), "setup: must start from a sealed file")

	s.startDaemon(c, config)
	_, grepOut := s.exec(c, []string{"grep", "-c", marker, secretFile})
	s.Require().Contains(grepOut, "1", "setup: the daemon must actually be running with the file unlocked")

	// Race a SIGKILL right on top of the graceful SIGTERM lockdown.
	s.sigDaemon(c, "TERM")
	time.Sleep(30 * time.Millisecond)
	s.sigDaemon(c, "KILL")
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not die after the SIGTERM+SIGKILL race")

	s.assertNoUnauthorizedPlaintext(c, secretFile, marker, "SIGTERM raced by SIGKILL mid-shutdown")

	code, out := s.runLockdown(c)
	s.Require().Equalf(0, code, "lockdown after a raced shutdown exited non-zero: %s", out)
	s.assertFileVaultSealed(c, secretFile, marker, "lockdown after a raced shutdown")

	s.startDaemon(c, config)
	_, grepOut = s.exec(c, []string{"grep", "-c", marker, secretFile})
	s.Require().Contains(grepOut, "1", "content must round-trip intact through the raced shutdown + lockdown + restart cycle")

	s.sigDaemon(c, "TERM")
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not exit after the final SIGTERM")
	s.runLockdown(c)
}

// F. Live-edit grant killed mid-session: the transient self-access widening (CLAUDE.md's "live-edit
// grant is a real escalation path") must never survive as a standing bypass for anyone but the
// exact trusted actor, and never survive a restart.

// Opens a real live edit-protected session (AUTH+SELECT over the control socket, widening the
// daemon's self-access on the target: guard.Guard.GrantSelfEditAccess) and races a SIGKILL against
// it before it ends cleanly.
//
// The client must be the real app-listener binary (authPeer refuses any other exe inode: the
// "peer-exe check"), so the write step can't be slowed from outside; instead --put content is large
// (512 MiB) so write+fsync between GRANTED and END is long enough for a tight log-poll-then-kill
// loop to land in the window on a fair fraction of runs. Best-effort race, but the assertions hold
// whether or not the kill landed mid-grant: the widened mask is scoped to exactly (uid 0, the
// daemon's exe inode), so a non-root writer and a root shell that is NOT the app-listener binary
// both stay denied (the "identity stays inode-based" invariant against a live orphaned grant), and
// unrelated content is never touched. A fresh daemon must return at the read-only baseline with no
// residual session: the grant lives only in the killed process's and its orphaned pinned map's
// state, never on disk.
func (s *IntegrationSuite) TestDaemon_LiveEditGrant_KilledMidSession_NoResidualEscalation() {
	const password = "Sup3r-Secret-KillRace-77"
	hash, err := editprotected.Hash(password, editprotected.OriginInstall)
	s.Require().NoError(err)

	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const secretFile = "/protected/secret.txt"
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /etc/app-listener && echo ORIGINAL-CONTENT > " + secretFile + " && chmod 600 " + secretFile + " && " +
			"printf '%s\\n' " + shQuote(hash) + " > /etc/app-listener/edit-auth.hash && chmod 600 /etc/app-listener/edit-auth.hash"})
	// Large --content-file: slows the write+fsync inside the grant window
	// (see the function doc comment) without touching secretFile itself.
	s.exec(c, []string{"sh", "-c", "dd if=/dev/zero of=/tmp/put_content.bin bs=1M count=512 status=none"})

	s.startDaemon(c, "[watch /protected]\nneed_encryption: false\n/usr/bin/sleep")

	socketReady := false
	for dl := time.Now().Add(20 * time.Second); time.Now().Before(dl); {
		if strings.Contains(s.readDaemonLog(c), "edit-protected control socket ready") {
			socketReady = true
			break
		}
		time.Sleep(300 * time.Millisecond)
	}
	s.Require().Truef(socketReady, "control socket never came up, log:\n%s", s.readDaemonLog(c))

	code, out := s.exec(c, []string{"sh", "-c", fmt.Sprintf(
		"APP_LISTENER_EDIT_PASSWORD=%s /app-listener edit-protected --resource /protected --put bigfile.bin --content-file /tmp/put_content.bin > /tmp/put.log 2>&1 &",
		password)})
	s.Require().Equalf(0, code, "backgrounding edit-protected --put: %s", out)

	// Tight, no-sleep poll: every millisecond spent here eats into the
	// (already narrow) window between GRANTED and the client's own END.
	granted := false
	for dl := time.Now().Add(10 * time.Second); time.Now().Before(dl); {
		if strings.Contains(s.readDaemonLog(c), "write access GRANTED") {
			granted = true
			break
		}
	}
	s.Require().Truef(granted, "edit-protected session never reached GRANTED, log:\n%s", s.readDaemonLog(c))

	s.sigDaemon(c, "KILL")
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not die after SIGKILL raced against the live-edit session")

	// Whether or not the kill actually beat the client's own END, the
	// widened self mask (orphaned in the still-pinned BPF program if the
	// kill won the race) is scoped to exactly (uid 0, the daemon's own exe
	// inode) — never anyone else.
	code, out = s.exec(c, []string{"sh", "-c",
		"setpriv --reuid=65534 --regid=65534 --clear-groups sh -c 'echo pwned > " + secretFile + "' 2>&1"})
	s.Require().NotEqualf(0, code, "a non-root writer must still be denied against a possibly-orphaned grant: %s", out)

	code, out = s.exec(c, []string{"sh", "-c", "echo pwned > " + secretFile + " 2>&1"})
	s.Require().NotEqualf(0, code, "a root shell that is NOT the app-listener binary must still be denied "+
		"(identity is inode-based, not uid-based): %s", out)

	_, content := s.exec(c, []string{"cat", secretFile})
	s.Require().Equal("ORIGINAL-CONTENT", strings.TrimSpace(content),
		"the session's grant covers the whole resource tree, but secret.txt itself — untouched by the legitimate --put — must stay exactly as it was")

	// A fresh instance comes back at the read-only baseline: no residual
	// session, and a plain write attempt outside any control-socket session
	// is denied exactly like it would be on a first-ever start.
	s.startDaemon(c, "[watch /protected]\nneed_encryption: false\n/usr/bin/sleep")
	code, out = s.exec(c, []string{"sh", "-c", "echo pwned > " + secretFile + " 2>&1"})
	s.Require().NotEqualf(0, code, "a restarted daemon must not carry over the killed session's grant: %s", out)

	s.sigDaemon(c, "TERM")
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not exit after the final SIGTERM")
}

// F. Recovery-sidecar symlink-follow: the file-vault sidecar (filevault.go's
// <resource>.app_listener.recover) must never be opened, read or created through a symlink.

// Proves a link-following bug in the file-vault recovery sidecar: an attacker controlling a
// file-vault resource's parent directory (real precondition: catalog file-vault targets like
// Steam's registry.vdf sit in the user's home) can, while the daemon is down, replace the
// not-yet-created sidecar with a symlink to any root-writable file. On the next start stageRecovery
// opens it O_RDWR and writes sealed ciphertext into that file, outside the resource and the
// attacker's own permissions.
//
// To make the corruption deterministic, /protected is remounted read-only right after staging the
// symlink: the harmful write (through the symlink, outside /protected) is unaffected, but
// transformFileInPlace on the real file fails with EROFS. That aborts the unlock before
// clearRecovery (which would truncate the sidecar, hence the symlinked victim, to empty), so the
// corruption persists and can be asserted.
func (s *IntegrationSuite) TestDaemon_FileVault_RecoverySidecarSymlink_NeverFollowed() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.copyFscryptHarness(c)

	const secretFile = "/protected/secret.txt"
	const victim = "/etc/victim.txt"

	setup := fmt.Sprintf(`set -e
mkdir -p /protected /etc/app-listener
head -c 32 /dev/zero > /etc/app-listener/fscrypt.key
chmod 600 /etc/app-listener/fscrypt.key
printf 'MARKER-CONTENT' > %s
chmod 644 %s
: > %s
chown 0:0 %s
chmod 600 %s
`, secretFile, secretFile, victim, victim, victim)
	code, out := s.exec(c, []string{"sh", "-c", setup})
	s.Require().Equalf(0, code, "setup: %s", out)

	s.harnessMigrate(c, secretFile)
	s.Require().True(s.harnessIsEncrypted(c, secretFile), "setup: the file must be sealed before the daemon ever starts")

	// The attack: the resource's parent dir is attacker-owned (chmod 777 stands in; the real
	// precondition is a user's own home directory), so an unprivileged user can create the missing
	// sidecar as a symlink to a root-writable file. The read-only remount only pins the
	// deterministic-failure trick; it plays no part in the vulnerability.
	attack := fmt.Sprintf(`set -e
chmod 777 /protected
setpriv --reuid=65534 --regid=65534 --clear-groups ln -s %s %s.app_listener.recover
mount --bind /protected /protected
mount -o remount,ro,bind /protected
`, victim, secretFile)
	code, out = s.exec(c, []string{"sh", "-c", attack})
	s.Require().Equalf(0, code, "staging the symlink attack: %s", out)
	defer s.exec(c, []string{"sh", "-c", "mount -o remount,rw,bind /protected 2>/dev/null; umount /protected 2>/dev/null; true"})

	config := fmt.Sprintf("[watch %s]\nneed_encryption: true\n/usr/bin/grep", secretFile)
	s.exec(c, []string{"sh", "-c", fmt.Sprintf("cat > /etc/app-listener/daemon.conf <<'EOF'\n%s\nEOF", config)})

	s.launchDaemon(c)
	// The daemon is expected to fail closed (read-only resource dir, hostile sidecar): poll for
	// either outcome without hard-failing, since the real assertion is what happened to the
	// symlink's target.
	deadline := time.Now().Add(20 * time.Second)
	for time.Now().Before(deadline) {
		if code, _ := s.exec(c, []string{"sh", "-c", "test -f /run/app-listener-daemon.pid && echo ready"}); code == 0 {
			break
		}
		if code, _ := s.exec(c, []string{"sh", "-c", noLiveAppListenerProcs}); code == 0 {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}

	_, victimContent := s.exec(c, []string{"cat", victim})
	s.Require().Emptyf(strings.TrimSpace(victimContent),
		"%s ended up non-empty (%q): an unprivileged user's symlink swap of the file-vault recovery sidecar made "+
			"the daemon write sealed file-vault bytes into a file it does not own, entirely outside the guarded "+
			"resource — daemon log:\n%s", victim, victimContent, s.readDaemonLog(c))

	s.sigDaemon(c, "KILL")
}

// G. Stale watch-root inode after in-place recreation (steady state): a guarded resource is deleted
// and recreated in place (app rebuild, fscrypt migration, backup restore) while the daemon runs.
// The kernel-side root anchor (guard_config[3..4], root_in_chain) and g.rootKey must follow the new
// inode via the periodic SweepInodes, or the recreated resource silently stops being guarded and
// whatever unrelated path gets the freed OLD inode number is denied by coincidence, misattributed
// to this resource. See SweepInodes' doc and TestSweepInodesRecreatedFileRoot/DirRoot
// (internal/guard/guard_test.go) for the BPF-map-level regression; these two go through a real
// daemon and the real periodic sweep.

// staleRootSweepPollTimeout comfortably exceeds the daemon's hardcoded
// periodic sweep interval (resyncSweepEvery, internal/usecase/daemon.go,
// currently 30s) so these tests poll for the sweep's effect rather than
// sleeping a fixed, easily-stale duration.
const staleRootSweepPollTimeout = 75 * time.Second

// awaitDenied polls until execing cmd in c is denied (nonzero exit),
// returning false if it is still allowed after timeout — used here to wait
// out the daemon's periodic SweepInodes tick rather than sleeping a fixed
// duration tied to its interval constant.
func (s *IntegrationSuite) awaitDenied(c testcontainers.Container, cmd []string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if code, _ := s.exec(c, cmd); code != 0 {
			return true
		}
		time.Sleep(2 * time.Second)
	}
	return false
}

// Directory-root regression for SweepInodes: only a single-file root's inode change was detected
// and re-anchored; a directory root's identity was never re-checked (only its mtime, to decide on a
// top-level re-walk), so wholesale recreation left guard_config/g.rootKey on the freed inode
// forever.
func (s *IntegrationSuite) TestDaemon_StaleRootInode_Directory_RecreatedRootReguarded() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"sh", "-c",
		"mkdir -p /watch /etc/app-listener && echo 'inside content' > /watch/inside.txt"})

	// mv is whitelisted so the moves below (simulating the app deleting and
	// recreating its own directory) succeed while every other process stays
	// denied on the guarded tree, exactly like the existing
	// TestGuard_InodeReuse_StaleEntryOutsideTree family.
	config := "[watch /watch]\nneed_encryption: false\n/usr/bin/mv"
	s.startDaemon(c, config)

	// Positive control: the original tree is guarded.
	code, out := s.exec(c, []string{"sh", "-c", "cat /watch/inside.txt"})
	s.Require().NotEqualf(0, code, "cat on the original in-tree file should be blocked: %s", out)

	// rename(2) preserves (dev, ino): moving the watch root out leaves its old inode at /old-root,
	// as when a filesystem later reuses that freed number. A fresh mkdir at /watch gets a genuinely
	// new inode (the parent / isn't guarded, so it's allowed regardless of whitelist).
	code, out = s.exec(c, []string{"mv", "/watch", "/old-root"})
	s.Require().Equalf(0, code, "whitelisted mv should move the watch root itself out: %s", out)
	code, out = s.exec(c, []string{"mkdir", "/watch"})
	s.Require().Equalf(0, code, "recreating /watch should not be denied (its parent is unguarded): %s", out)
	s.exec(c, []string{"sh", "-c", "echo 'recreated content' > /watch/new.txt"})

	// Wait out the periodic sweep: it must re-anchor the root to the new
	// inode, at which point the RECREATED tree becomes guarded again.
	s.Require().Truef(
		s.awaitDenied(c, []string{"sh", "-c", "cat /watch/new.txt"}, staleRootSweepPollTimeout),
		"the recreated directory root must become guarded again once the periodic sweep re-anchors it — daemon log:\n%s",
		s.readDaemonLog(c))

	// Regression capture: the OLD path (now holding the stale inode) must
	// no longer be denied by a leftover root-confinement match — it is
	// genuinely outside the tree once the anchor has moved on.
	code, out = s.exec(c, []string{"cat", "/old-root/inside.txt"})
	s.Require().Equalf(0, code, "the old path holding the stale root inode must not be denied once the anchor has moved on: %s", out)
	s.Require().Contains(out, "inside content")

	s.sigDaemon(c, "TERM")
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not exit after SIGTERM")
}

// Single-file counterpart of TestDaemon_StaleRootInode_Directory_RecreatedRootReguarded: the case
// SweepInodes already handled (TestSweepInodesRecreatedFileRoot), now end-to-end through a real
// daemon so a change to the shared plumbing (updateRootKey, sweep wiring) breaking either case is
// caught by both.
func (s *IntegrationSuite) TestDaemon_StaleRootInode_File_RecreatedRootReguarded() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"sh", "-c",
		"mkdir -p /etc/app-listener && echo 'v1' > /registry.vdf"})

	config := "[watch /registry.vdf]\nneed_encryption: false\n/usr/bin/mv"
	s.startDaemon(c, config)

	code, out := s.exec(c, []string{"cat", "/registry.vdf"})
	s.Require().NotEqualf(0, code, "cat on the original guarded file should be blocked: %s", out)

	// Same rename-preserves-inode trick as the directory case: moving the
	// single-file watch root out leaves its old inode at /old-registry.vdf,
	// then a fresh file created at the original path gets a new inode.
	code, out = s.exec(c, []string{"mv", "/registry.vdf", "/old-registry.vdf"})
	s.Require().Equalf(0, code, "whitelisted mv should move the watch root file itself out: %s", out)
	code, out = s.exec(c, []string{"sh", "-c", "echo 'v2-recreated' > /registry.vdf"})
	s.Require().Equalf(0, code, "recreating /registry.vdf should not be denied (its parent is unguarded): %s", out)

	s.Require().Truef(
		s.awaitDenied(c, []string{"cat", "/registry.vdf"}, staleRootSweepPollTimeout),
		"the recreated file root must become guarded again once the periodic sweep re-anchors it — daemon log:\n%s",
		s.readDaemonLog(c))

	code, out = s.exec(c, []string{"cat", "/old-registry.vdf"})
	s.Require().Equalf(0, code, "the old path holding the stale root inode must not be denied once the anchor has moved on: %s", out)
	s.Require().Contains(out, "v1")

	s.sigDaemon(c, "TERM")
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not exit after SIGTERM")
}

// The guard_path_unlink eviction fix is tested at guard level
// (TestGuard_PathUnlinkEvictsInodeImmediately, TestGuard_PathUnlinkDeniedDeleteKeepsGuardedInode in
// integrationtests/guard_test.go, with the TestGuard_InodeReuse_* family).

// Trust guard provenance rule for runtime-generated code (guard_jit_origin, guard_trust.bpf.c). GPU
// drivers JIT into a file they create (memfd, O_TMPFILE, mkstemp) and map it executable; no path or
// ownership rule covers those inodes, so the daemon trusts an exec mapping when the mapping process
// image created the inode and nothing else wrote it.
//
// The rule is an ALLOWANCE, so most of the test is what it must still refuse: an inode the attacker
// created (no provenance), a merely new file, an inode handed over through execve (same pid,
// different image), and one a foreign writer touched after creation.
func (s *IntegrationSuite) TestDaemon_JitProvenanceTrustsOnlySelfCreatedCode() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const marker = "TOP-SECRET-JIT-PROVENANCE-9B1D"
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /exploits /etc/app-listener && printf '" + marker + "' > /protected/secret && chmod 755 /protected && chmod 644 /protected/secret"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/jit_provenance"), "/exploits/jit_provenance", 0755), "copy jit_provenance")
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/preload_leak.so"), "/exploits/leak.so", 0755), "copy preload_leak.so")
	// The same program at two paths: only /tmp/app is whitelisted, so
	// /tmp/attacker stands in for same-user malware.
	s.exec(c, []string{"sh", "-c", "cp /exploits/jit_provenance /tmp/app && cp /exploits/jit_provenance /tmp/attacker && chmod 755 /tmp/app /tmp/attacker"})

	// A benign library to stand in for JIT output: what matters is that the
	// bytes are a loadable ELF the process produced itself.
	const benignLib = "/lib/x86_64-linux-gnu/libz.so.1"
	code, out := s.exec(c, []string{"sh", "-c", "test -f " + benignLib})
	s.Require().Equalf(0, code, "test fixture expects %s in the image: %s", benignLib, out)

	s.startDaemon(c, `[watch /protected]
need_encryption: false
/tmp/app
/usr/bin/true`)

	// 1. FUNCTIONAL: the whitelisted app may map code it created itself with
	//    mkstemp — the driver's file-backed JIT path, and the reason this rule
	//    exists at all.
	code, out = s.exec(c, []string{"/tmp/app", "self", benignLib})
	s.Require().Containsf(out, "MAPPED",
		"a whitelisted process must be able to exec-map code it created itself (exit %d): %s", code, out)

	// 2. FUNCTIONAL: the memfd variant, the first step of the fallback chain.
	//    It depends on the fexit attach (memfd_create never reaches
	//    security_file_open), so tolerate a kernel where that program could
	//    not attach — the daemon says so, and the denial is fail-closed.
	code, out = s.exec(c, []string{"/tmp/app", "memfd", benignLib})
	if !strings.Contains(out, "MAPPED") {
		_, logOut := s.exec(c, []string{"sh", "-c", "grep -c 'skipping memfd provenance' /tmp/daemon.log || true"})
		s.Require().NotEqualf("0", strings.TrimSpace(logOut),
			"memfd-backed self-created code was denied although the memfd provenance program attached (exit %d): %s", code, out)
	}

	// 3. ATTACK: the attacker allocates the memfd and hands it to a
	//    whitelisted victim through an inherited fd plus LD_PRELOAD. The
	//    victim maps an inode somebody else made: no provenance, no mapping.
	_, out = s.exec(c, []string{"/tmp/attacker", "passfd", "/exploits/leak.so", "/usr/bin/true", "/protected/secret"})
	s.Require().NotContainsf(out, marker,
		"guarded content exfiltrated: a whitelisted binary mapped a memfd created by an unwhitelisted process: %s", out)

	// 4. ATTACK: the same handover performed by a WHITELISTED creator. execve
	//    keeps the thread group id, so a pid-only provenance check would let
	//    the victim inherit the attacker's memfd as if it had made it itself.
	//    The mm and exe identity are what refuse it.
	_, out = s.exec(c, []string{"/tmp/app", "passfd", "/exploits/leak.so", "/usr/bin/true", "/protected/secret"})
	s.Require().NotContainsf(out, marker,
		"guarded content exfiltrated: provenance followed a pid across execve into a different image: %s", out)

	// 5. ATTACK: a file is not trustworthy merely for being new — the
	//    attacker creates it and the whitelisted victim preloads it.
	_, out = s.exec(c, []string{"sh", "-c",
		"/tmp/attacker self /exploits/leak.so >/tmp/fresh.txt 2>&1; p=$(sed -n 's/^\\(MAPPED\\|DENIED\\) \\([^ :]*\\).*/\\2/p' /tmp/fresh.txt | head -1); " +
			"LD_PRELOAD=$p LEAK_FILE=/protected/secret /usr/bin/true 2>&1"})
	s.Require().NotContainsf(out, marker,
		"guarded content exfiltrated through a freshly created attacker .so: %s", out)

	// 6. ATTACK: a foreign write after creation poisons the inode for good.
	//    The whitelisted app creates its JIT file, an unrelated process
	//    write-opens it (the /proc/<pid>/fd shape of the attack), and the
	//    app's own mapping must then be refused.
	s.exec(c, []string{"sh", "-c", "(/tmp/app hold " + benignLib + " 6 >/tmp/hold.log 2>&1 &) ; sleep 2"})
	_, pathOut := s.exec(c, []string{"sh", "-c", "sed -n 's/^PATH=//p' /tmp/hold.log | head -1"})
	jitPath := strings.TrimSpace(pathOut)
	s.Require().NotEmptyf(jitPath, "hold mode did not report its JIT file: %s", pathOut)
	// A write-open by another process is all it takes; the bytes need not change.
	s.exec(c, []string{"sh", "-c", "exec 3<>" + jitPath + "; exec 3>&-"})
	s.exec(c, []string{"sh", "-c", "sleep 6"})
	_, holdOut := s.exec(c, []string{"sh", "-c", "cat /tmp/hold.log"})
	s.Require().Containsf(holdOut, "DENIED",
		"a JIT file that a foreign process write-opened must no longer be mappable: %s", holdOut)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// An application's atomic save (temp write + rename over the guarded file; Steam's registry.vdf
// every launch) gives the watch root a new inode that is neither root nor in guard_inodes:
// unguarded until the daemon re-anchors. This used to wait for the 30 s sweep; single-file roots
// are now followed every second (fileRootFollowEvery). The kernel can't follow the rename itself
// (guard_path_rename has no verifier budget left, issue #45).
func (s *IntegrationSuite) TestDaemon_ReplacedSingleFileRootIsGuardedQuickly() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const secret = "REPLACED-ROOT-SECRET-4B7E"
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /data /etc/app-listener && printf 'OLD' > /data/registry.vdf && chmod 755 /data"})
	s.startDaemon(c, `[watch /data/registry.vdf]
need_encryption: false
/usr/bin/mv`)

	code, out := s.exec(c, []string{"sh", "-c", "cat /data/registry.vdf 2>&1"})
	s.Require().NotEqualf(0, code, "control: the file must be guarded before the save: %s", out)

	code, out = s.exec(c, []string{"sh", "-c",
		"printf '" + secret + "' > /data/reg.tmp && /usr/bin/mv /data/reg.tmp /data/registry.vdf"})
	s.Require().Equalf(0, code, "the whitelisted atomic save must succeed: %s", out)

	// Re-anchored within a few follow ticks, not the 30 s sweep.
	deadline := time.Now().Add(4 * time.Second)
	for {
		code, out = s.exec(c, []string{"sh", "-c", "cat /data/registry.vdf 2>&1"})
		if code != 0 && !strings.Contains(out, secret) {
			break
		}
		s.Require().Truef(time.Now().Before(deadline),
			"the replaced file was still readable by a non-whitelisted process after 4s: %s", out)
		time.Sleep(200 * time.Millisecond)
	}

	// And it stays guarded.
	_, out = s.exec(c, []string{"sh", "-c", "sleep 2; cat /data/registry.vdf 2>&1"})
	s.Require().NotContainsf(out, secret, "the re-anchored root must stay guarded: %s", out)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// procReadProbe asks whether a PTRACE_MODE_READ-gated /proc access against pid is refused,
// reporting in-band. It uses readlink(/proc/<pid>/exe): proc_pid_get_link() consults
// proc_fd_access_allowed() -> ptrace_may_access(PTRACE_MODE_READ) and FAILS the syscall when the
// guard refuses (seen in the daemon log). Two look-alike probes are NOT usable:
//   - `ls /proc/<pid>/fd`: proc_fd_permission() returns once generic_permission() passes, and
//     container root passes via CAP_DAC_OVERRIDE, so the LSM hook is never consulted;
//   - `cat /proc/<pid>/maps`: was observed reading a guarded process's maps in full WHILE the
//     daemon logged a matching op=PTRACE mode=READ denial. UNEXPLAINED: proc_maps_open() takes its
//     mm from proc_mem_open(PTRACE_MODE_READ) and returns the error, and the read path has no
//     advisory ptrace check, so the open should fail. Until the raw ptrace mode bits are logged
//     (0x04 = NOAUDIT, an advisory call site) the denial can't be attributed to that open, so maps
//     is unusable either way. ATTACH-class access (/proc/<pid>/mem, process_vm_readv) IS refused:
//     see TestDaemon_OwnMetadataReadableMemoryNot and TestGuard_Bypass_ProcessVmReadv.
//
// rc is readlink's own status, not the exec's (a docker exec exit code has been 0 for a refused
// command). state proves the target is alive: a dead pid also fails readlink and would pass a
// denial assertion vacuously.
func procReadProbe(pid string) string {
	return "exe=$(readlink /proc/" + pid + "/exe 2>&1); rc=$?; " +
		"st=$(awk '{print $3}' /proc/" + pid + "/stat 2>/dev/null); " +
		"echo \"rc=$rc state=$st exe=$exe\""
}

// procGateMatrix reports every /proc access class for one pid, so a failure
// says which ptrace-gated reads this kernel actually refuses.
func procGateMatrix(pid string) string {
	return "for f in exe maps environ mem stat status cmdline; do " +
		"if [ \"$f\" = exe ]; then o=$(readlink /proc/" + pid + "/exe 2>&1); " +
		"else o=$(head -c 40 /proc/" + pid + "/$f 2>&1 | tr -d \"\\0\" | head -1); fi; " +
		"echo \"  $f rc=$? out=$o\"; done"
}

// Running a lib_binary writer of a read-only lib_dir must not taint the process. Taint shields a
// process holding SECRETS from ptrace-class inspection; a read-only tree is world-readable code
// with nothing to shield. The exec hook tainted in every guard anyway (unlike the file-access taint
// sites, which skip read-only guards), leaving each runtime tree's writer list as the only
// processes allowed to look at /proc/<pid> of half of Steam's process tree, silently.
func (s *IntegrationSuite) TestDaemon_ReadOnlyGuardDoesNotTaintItsWriters() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	// The writer copy must keep the basename `sleep`: ubuntu:latest ships
	// coreutils as one multi-call binary that dispatches on argv[0], so a
	// copy named /tmp/rtwriter dies instantly with "coreutils: unknown
	// program 'rtwriter'" and the test found no process to inspect.
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /rt /tmp/rtw /etc/app-listener && printf 'code' > /rt/lib.so && cp /usr/bin/sleep /tmp/rtw/sleep && chmod 755 /tmp/rtw/sleep"})
	s.startDaemon(c, `[libraries "Runtime"]
lib_dir /rt
lib_binary /tmp/rtw/sleep`)

	s.exec(c, []string{"sh", "-c", "(/tmp/rtw/sleep 30 &) ; sleep 1"})
	_, pidOut := s.exec(c, []string{"sh", "-c", "pgrep -f '^/tmp/rtw/sleep 30' | head -1"})
	pid := strings.TrimSpace(pidOut)
	s.Require().NotEmptyf(pid, "the writer process did not start: %q", pidOut)

	_, out := s.exec(c, []string{"sh", "-c", procReadProbe(pid)})
	s.Require().Containsf(out, "rc=0",
		"a process running a read-only tree's writer must stay inspectable by an unrelated process: %s", out)
	s.Require().Containsf(out, "/tmp/rtw/sleep",
		"the probe did not actually resolve the writer's exe, so it proves nothing: %s", out)

	_, logOut := s.exec(c, []string{"sh", "-c", "grep -c 'op=PTRACE' /tmp/daemon.log || true"})
	s.Require().Equalf("0", strings.TrimSpace(logOut), "no ptrace-class denial may come from a read-only guard")

	s.exec(c, []string{"sh", "-c", "pkill -f '^/tmp/rtw/sleep' ; pkill -f 'app-listener daemon' || true"})
}

// Pins the taint lifecycle:
//   - the process that read the secret is tainted (not inspectable);
//   - a forked child keeping the parent's memory stays tainted;
//   - a child exec'ing a NON-whitelisted image is cleared (exec discards the address space, so it
//     holds nothing from the vault).
//
// Before the exec rule every descendant of a tainted process (the whole Steam tree) stayed tainted
// for life, and GameMode, PipeWire, the portal and Wine's wineserver were refused every look. The
// refusal must also be logged with its access mode.
func (s *IntegrationSuite) TestDaemon_TaintFollowsForkButNotExec() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /etc/app-listener && printf 'TAINT-LIFECYCLE-SECRET' > /protected/secret && " +
			"chmod 755 /protected && chmod 644 /protected/secret && cp /bin/bash /tmp/wsh && chmod 755 /tmp/wsh && mkfifo /tmp/idle_fifo"})
	s.startDaemon(c, `[watch /protected]
need_encryption: false
/tmp/wsh`)

	// The whitelisted shell reads the secret with a builtin (no exec), then
	// starts a fork-only subshell and a fork+exec child, recording pids.
	s.exec(c, []string{"sh", "-c", `nohup /tmp/wsh -c '
read -r x < /protected/secret
printf "%s" "$x" > /tmp/p_value
echo $$ > /tmp/p_reader
( echo $BASHPID > /tmp/p_fork; while :; do read -t 1 -r _ <> /tmp/idle_fifo || true; done ) &
/usr/bin/sleep 60 & echo $! > /tmp/p_exec
wait' >/dev/null 2>&1 &
sleep 2`})

	pidOf := func(f string) string {
		_, out := s.exec(c, []string{"sh", "-c", "cat " + f})
		p := strings.TrimSpace(out)
		s.Require().NotEmptyf(p, "missing %s", f)
		return p
	}
	inspect := func(pid string) (int, string) {
		return s.exec(c, []string{"sh", "-c", procReadProbe(pid)})
	}
	// diagnose is printed when a denial assertion fails: the whole /proc
	// access matrix for the target (which classes this kernel actually
	// refuses), its identity and liveness, a control read by an
	// unwhitelisted process, and the daemon's own log.
	diagnose := func(pid string) string {
		_, diag := s.exec(c, []string{"sh", "-c",
			"echo 'proc gate matrix:'; " + procGateMatrix(pid) + "; " +
				"echo \"comm=$(cat /proc/" + pid + "/comm)\"; " +
				"echo \"state=$(awk '{print $3}' /proc/" + pid + "/stat) threads=$(awk '/^Threads:/{print $2}' /proc/" + pid + "/status)\"; " +
				"echo \"whitelist-inode=$(stat -c %d:%i /tmp/wsh)\"; " +
				"echo \"control-read: $(cat /protected/secret 2>&1)\"; " +
				"echo '--- daemon.log ---'; tail -30 /tmp/daemon.log"})
		return diag
	}
	requireDenied := func(label, pid string) {
		_, out := inspect(pid)
		s.Require().Regexpf(`state=[RSD]`, out,
			"%s: the target is not a live process, so a refusal proves nothing: %s", label, out)
		if strings.Contains(out, "rc=0") {
			s.Require().Failf(label, "the probe was not refused: pid=%s %s\n%s", pid, out, diagnose(pid))
		}
	}
	requireAllowed := func(label, pid string) {
		_, out := inspect(pid)
		s.Require().Containsf(out, "rc=0", "%s: %s", label, out)
		s.Require().Containsf(out, "exe=/", "%s: the probe resolved no exe, so it proves nothing: %s", label, out)
	}

	// Control first: taint is stamped when the whitelisted image execs and
	// when it reads, so if the read itself was refused the whole test is
	// vacuous — and the failure looks identical to a missing gate.
	_, readValue := s.exec(c, []string{"sh", "-c", "cat /tmp/p_value 2>&1"})
	s.Require().Equalf("TAINT-LIFECYCLE-SECRET", strings.TrimSpace(readValue),
		"the whitelisted shell did not read the secret, so nothing tainted it: %q", readValue)

	requireDenied("the process that read the secret must not be inspectable by an unwhitelisted process",
		pidOf("/tmp/p_reader"))
	requireDenied("a forked child sharing the reader's memory image must stay tainted", pidOf("/tmp/p_fork"))
	requireAllowed("a child that exec'd an unwhitelisted image holds none of the vault's memory: its taint must be cleared",
		pidOf("/tmp/p_exec"))

	_, logOut := s.exec(c, []string{"sh", "-c", "grep -c 'op=PTRACE.*mode=READ' /tmp/daemon.log || true"})
	s.Require().NotEqualf("0", strings.TrimSpace(logOut),
		"a /proc/<pid>/maps denial must be logged as op=PTRACE with mode=READ")

	s.exec(c, []string{"sh", "-c", "pkill -f /tmp/wsh; pkill -f 'sleep 60'; pkill -f 'app-listener daemon' || true"})
}

// The daemon is tainted (it reads what it guards) and holds secrets (fscrypt key, edit-auth hash)
// in memory. A ptrace ATTACH-class access (/proc/<pid>/mem, process_vm_readv) must stay denied; a
// READ-class one (/proc/<pid>/maps, environ, exe; what journald does to label log lines) is
// allowed, else each logged denial makes journald look it up and fail again.
//
// In a pid namespace the denial can only come from ptrace_access_check, which reads child->tgid
// (init-namespace). The second gate, is_proc_mem_of_tainted() in file_open/file_permission, parses
// the /proc directory NAME (the container's pid) and compares it to guard_tainted_pids, keyed on
// the init-namespace tgid by mark_tainted(), so it can't match in a namespace. Both agree on the
// host (the target); the namespace asymmetry only costs the secondary catch for reads on a
// /proc/<pid>/mem fd opened before tainting.
func (s *IntegrationSuite) TestDaemon_OwnMetadataReadableMemoryNot() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected/sub /etc/app-listener && printf 'S' > /protected/secret && chmod 755 /protected"})
	s.startDaemon(c, `[watch /protected]
need_encryption: false
/usr/bin/true`)

	_, pidOut := s.exec(c, []string{"sh", "-c", "cat /run/app-listener-daemon.pid"})
	pid := strings.TrimSpace(pidOut)
	s.Require().NotEmpty(pid, "daemon pid")

	// ATTACH-class: still refused; this also proves the daemon IS tainted, so the READ assertion
	// below isn't vacuous (an untainted process would open it and give an I/O error at offset 0).
	// The refusal is EACCES, not the gate's EPERM: opening /proc/<pid>/mem calls
	// mm_access(PTRACE_MODE_ATTACH), which maps any LSM denial to -EACCES before
	// security_file_open. Keyed on that message, not the exec exit code (docker exec has reported 0
	// for a refused command).
	_, out := s.exec(c, []string{"sh", "-c", "head -c1 /proc/" + pid + "/mem 2>&1"})
	s.Require().Regexpf("Permission denied|not permitted", out,
		"memory of the daemon must stay unreadable to an unwhitelisted process: %s", out)

	// READ-class: allowed.
	_, out = s.exec(c, []string{"sh", "-c", procReadProbe(pid)})
	s.Require().Containsf(out, "rc=0",
		"metadata of the daemon's own process must be readable — journald resolves /proc/<pid>/exe for every line it labels: %s", out)
	s.Require().Containsf(out, "exe=/", "the probe resolved no exe, so it proves nothing: %s", out)
	s.Require().NotEmptyf(out, "empty maps: the daemon's address space was not really read, so the probe proves nothing")

	_, logOut := s.exec(c, []string{"sh", "-c", "grep -c 'op=PTRACE.*comm=app-listener mode=READ' /tmp/daemon.log || true"})
	s.Require().Equalf("0", strings.TrimSpace(logOut), "no READ-mode denial may be logged for the daemon itself")

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// daemonScaleResources exceeds BPF_MAX_TRAMP_LINKS (38), the kernel's cap on trampoline links at a
// single attach point. A real install guards ~34 catalog resources plus the daemon's two
// self-guards, so production already sits just under it.
const daemonScaleResources = 48

// Every guarded resource used to attach its own copy of all 26 LSM programs, so each one added a
// link to each attach point and usage on security_file_open was O(resources). Past the cap the
// kernel returns E2BIG and the daemon refuses to start, which made the number of [watch] sections a
// hard ceiling — and, because BPF-LSM links are global to the kernel rather than namespaced, a
// daemon near the cap also starved every other BPF-LSM user on the host, containers included.
//
// Enforcement state is keyed by dev:ino, so one attached program set serves every resource by
// resolving the accessed inode to its resource in-kernel, keeping attach usage O(1).
func (s *IntegrationSuite) TestDaemon_ManyResources_AttachStaysWithinTrampolineCap() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	var mkdirs, config strings.Builder
	mkdirs.WriteString("mkdir -p /etc/app-listener")
	for i := range daemonScaleResources {
		dir := fmt.Sprintf("/protected%02d", i)
		fmt.Fprintf(&mkdirs, " %s", dir)
		fmt.Fprintf(&config, "[watch %s]\nneed_encryption: false\n/usr/bin/sleep\n\n", dir)
	}
	s.exec(c, []string{"sh", "-c", mkdirs.String()})
	s.exec(c, []string{"sh", "-c",
		"echo TOP-SECRET-CONTENT > /protected00/secret && chmod 755 /protected00 && chmod 644 /protected00/secret"})

	// startDaemon fails the test if any resource's guard never reports "guard started", which is
	// what an exhausted attach point (E2BIG) causes.
	s.startDaemon(c, config.String())

	log := s.readDaemonLog(c)
	s.Require().NotContainsf(log, "argument list too long",
		"the daemon hit the per-attach-point trampoline cap with %d resources: %s", daemonScaleResources, log)

	// All 48 attached — now prove enforcement still works, so the shared attach did not trade the
	// cap for a guard that denies nothing.
	const nobody = "setpriv --reuid=65534 --regid=65534 --clear-groups"
	code, out := s.exec(c, []string{"sh", "-c", nobody + " grep -c TOP-SECRET /protected00/secret"})
	s.Require().NotEqualf(0, code, "non-whitelisted reader must still be denied with %d resources, got: %s",
		daemonScaleResources, out)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// steamGlobConfig guards a Steam location so the catalog's Steam globs are reserved for root's
// home, with a dash copy standing in for the Steam client (its builtins write as its own inode).
// /usr/bin/bash is whitelisted for an unrelated resource: the confused deputy.
const (
	steamDir        = "/root/.local/share/Steam"
	steamCommon     = steamDir + "/steamapps/common"
	steamClient     = steamDir + "/ubuntu12_32/steam"
	steamGlobConfig = `[watch ` + steamDir + `/config]
need_encryption: false
` + steamClient + `

[watch /protected]
need_encryption: false
/usr/bin/bash`
)

// steamTools holds pressure-vessel's pv-*/srt-*/*-capsule-capture-libs helpers. Created before the
// daemon starts: pv-runtime itself matches pv-*, so only a Steam binary may create it afterwards.
const steamTools = steamDir + "/steamrt64/pv-runtime/x/pressure-vessel/libexec/steam-runtime-tools-0"

func (s *IntegrationSuite) startSteamGlobDaemon(c testcontainers.Container) {
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /exploits /etc/app-listener " + steamDir + "/config " + steamCommon + " " + steamTools +
			" $(dirname " + steamClient + ") && cp /usr/bin/dash " + steamClient +
			" && cp /usr/bin/true /tmp/stealer"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/glob_plant"), "/exploits/glob_plant", 0755),
		"copy glob_plant")
	s.startDaemon(c, steamGlobConfig)
	s.Require().Contains(s.readDaemonLog(c), "reserved glob name(s)", "trust guard #3 was not populated")
}

// assertNotPlanted runs cmd (a plant attempt) and requires it to fail with path left absent.
func (s *IntegrationSuite) assertNotPlanted(c testcontainers.Container, what, path string, cmd []string) {
	code, out := s.exec(c, cmd)
	s.Require().NotEqualf(0, code, "%s: planting %s must be denied: %s", what, path, out)
	code, _ = s.exec(c, []string{"sh", "-c", "test -e " + shQuote(path) + " || test -L " + shQuote(path)})
	s.Require().NotEqualf(0, code, "%s: %s exists after a denied plant", what, path)
}

// Vuln 3: a same-user process must not be able to create any path a catalog whitelist glob
// matches (it would be whitelisted at the next refresh), by any route, while the wildcard dirs stay
// otherwise writable.
func (s *IntegrationSuite) TestDaemon_GlobPlant_NonWriterDenied() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.startSteamGlobDaemon(c)

	bin := steamCommon + "/exploit/files/bin"
	ws := bin + "/wineserver"
	code, out := s.exec(c, []string{"mkdir", "-p", bin})
	s.Require().Equalf(0, code, "the wildcard dirs must stay writable: %s", out)

	s.assertNotPlanted(c, "create", ws, []string{"sh", "-c", "echo x > " + ws})
	s.assertNotPlanted(c, "symlink", ws, []string{"ln", "-s", "/tmp/stealer", ws})
	s.assertNotPlanted(c, "hardlink", ws, []string{"ln", "/tmp/stealer", ws})
	s.assertNotPlanted(c, "rename", ws, []string{"mv", "/tmp/stealer", ws})
	s.assertNotPlanted(c, "mkdir", ws, []string{"mkdir", ws})
	s.assertNotPlanted(c, "O_TMPFILE+linkat", ws, []string{"/exploits/glob_plant", "tmpfile", bin, ws})

	// A directory assembled outside and moved in carries its content past the per-name check.
	s.exec(c, []string{"sh", "-c", "mkdir -p /tmp/evil/files/bin && cp /usr/bin/true /tmp/evil/files/bin/wineserver"})
	s.assertNotPlanted(c, "directory move", steamCommon+"/evil", []string{"mv", "/tmp/evil", steamCommon + "/evil"})

	// Wildcard names: prefix (pv-*) and suffix (*-capsule-capture-libs).
	s.assertNotPlanted(c, "prefix name", steamTools+"/pv-adverb", []string{"sh", "-c", "echo x > " + steamTools + "/pv-adverb"})
	s.assertNotPlanted(c, "suffix name", steamTools+"/x86_64-linux-gnu-capsule-capture-libs",
		[]string{"sh", "-c", "echo x > " + steamTools + "/x86_64-linux-gnu-capsule-capture-libs"})

	// Negative controls: games keep writing their dirs, and only the fixed name is reserved.
	for _, p := range []string{steamCommon + "/exploit/save.dat", bin + "/wineserver2", "/tmp/wineserver"} {
		code, out = s.exec(c, []string{"sh", "-c", "echo x > " + p})
		s.Require().Equalf(0, code, "%s is not a glob match and must stay writable: %s", p, out)
	}

	s.Require().Contains(s.readDaemonLog(c), "TRUST DENIED  op=PLANT", "denials must be logged")

	// Renaming the root away is allowed, but recreating it would give an unregistered root.
	code, out = s.exec(c, []string{"mv", steamCommon, steamCommon + ".old"})
	s.Require().Equalf(0, code, "renaming the root away is harmless: %s", out)
	s.assertNotPlanted(c, "root recreation", steamCommon, []string{"mkdir", steamCommon})

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Vuln 3: the owning app still installs new versions (the wildcard dir is new, the reserved name
// is written by its own binary), and until the refresh whitelists that fresh match, nothing else,
// not even a binary whitelisted for another resource, may rewrite or alias it.
func (s *IntegrationSuite) TestDaemon_GlobPlant_OwnWriterAllowed() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.startSteamGlobDaemon(c)

	bin := steamCommon + "/newVersion/files/bin"
	ws := bin + "/wineserver"
	s.exec(c, []string{"mkdir", "-p", bin})
	code, out := s.exec(c, []string{steamClient, "-c", "printf STEAM > " + ws})
	s.Require().Equalf(0, code, "Steam's own binary must create %s: %s", ws, out)

	for _, tc := range []struct {
		what string
		cmd  []string
	}{
		{"append", []string{"sh", "-c", "echo evil >> " + ws}},
		{"truncate", []string{"truncate", "-s", "0", ws}},
		{"hardlink alias", []string{"ln", ws, "/tmp/alias"}},
		{"rename over", []string{"sh", "-c", "cp /usr/bin/true /tmp/over && mv -f /tmp/over " + ws}},
		{"exchange onto it", []string{"sh", "-c", "cp /usr/bin/true /tmp/swap && /exploits/glob_plant exchange /tmp/swap " + ws}},
		{"exchange from it", []string{"sh", "-c", "cp /usr/bin/true /tmp/swap && /exploits/glob_plant exchange " + ws + " /tmp/swap"}},
	} {
		code, out = s.exec(c, tc.cmd)
		s.Require().NotEqualf(0, code, "%s of the fresh match by a non-writer must be denied: %s", tc.what, out)
		_, content := s.exec(c, []string{"cat", ws})
		s.Require().Equalf("STEAM", content, "%s altered the fresh match", tc.what)
	}

	// Confused deputy: whitelisted, but for another resource.
	deputy := steamCommon + "/deputy/files/bin"
	s.exec(c, []string{"mkdir", "-p", deputy})
	code, out = s.exec(c, []string{"/usr/bin/bash", "-c", "echo x > " + deputy + "/wineserver"})
	s.Require().NotEqualf(0, code, "a binary whitelisted for another resource must not plant: %s", out)

	code, out = s.exec(c, []string{steamClient, "-c", "printf UPDATED >> " + ws})
	s.Require().Equalf(0, code, "Steam's own binary must keep updating it: %s", out)
	_, out = s.exec(c, []string{"cat", ws})
	s.Require().Equal("STEAMUPDATED", out)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// discordLibConfig guards a Discord watch path so the catalog's ReservedLibs are reserved below
// root's ~/.config/discord, with a dash copy standing in for the Discord client (a writer).
// /usr/bin/bash is whitelisted for an unrelated resource: whitelisted, but not a writer.
const (
	discordDir       = "/root/.config/discord"
	discordModules   = discordDir + "/0.0.1/modules"
	discordClient    = discordDir + "/0.0.1/Discord"
	libProbe         = "/exploits/lib_probe.so"
	libProbeMarker   = "LIB_PROBE_LOADED"
	discordLibConfig = `[watch ` + discordDir + `/sentry]
need_encryption: false
` + discordClient + `

[watch /protected]
need_encryption: false
/usr/bin/bash`
)

// startDiscordLibDaemon hands modules/ to nobody: a root-owned dir would make every file in it an
// auto-trusted system library (is_system_trusted) and the load tests vacuous.
func (s *IntegrationSuite) startDiscordLibDaemon(c testcontainers.Container) {
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /exploits /etc/app-listener " + discordDir + "/sentry " + discordModules +
			"/discord_voice && cp /usr/bin/dash " + discordClient + " && chown -R 65534 " + discordModules})
	for _, f := range []string{"lib_probe.so", "glob_plant"} {
		s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/"+f), "/exploits/"+f, 0755), "copy "+f)
	}
	s.startDaemon(c, discordLibConfig)
	s.Require().Contains(s.readDaemonLog(c), "reserved glob name(s)", "trust guard #3 was not populated")
}

// writeAsDiscord creates path with lib_probe's bytes; dash's redirection opens it as the client.
func (s *IntegrationSuite) writeAsDiscord(c testcontainers.Container, path string) {
	code, out := s.exec(c, []string{discordClient, "-c", "cat " + libProbe + " > " + shQuote(path)})
	s.Require().Equalf(0, code, "Discord's own binary must write %s: %s", path, out)
}

// preload runs bin with lib LD_PRELOADed and reports whether lib's constructor ran.
func (s *IntegrationSuite) preload(c testcontainers.Container, bin, lib string) (bool, string) {
	_, out := s.exec(c, []string{"sh", "-c", "LD_PRELOAD=" + shQuote(lib) + " " + bin + " -c true 2>&1"})
	return strings.Contains(out, libProbeMarker), out
}

func (s *IntegrationSuite) requireDenialLogged(c testcontainers.Container, op, name string) {
	for _, line := range strings.Split(s.readDaemonLog(c), "\n") {
		if strings.Contains(line, "TRUST DENIED  op="+op) && strings.Contains(line, name) {
			return
		}
	}
	s.Failf("denial not logged", "no op=%s line for %s in:\n%s", op, name, s.readDaemonLog(c))
}

// Discord's self-updated native code under ~/.config/discord is outside every guarded tree and
// not root-owned; a reserved name its own binary wrote there must load, at any depth.
func (s *IntegrationSuite) TestDaemon_ReservedLib_WriterLibraryLoads() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.startDiscordLibDaemon(c)

	for _, lib := range []string{
		discordModules + "/probe.so",
		discordModules + "/libprobe.so.1",
		discordModules + "/discord_voice/discord_voice.node",
	} {
		s.writeAsDiscord(c, lib)
		loaded, out := s.preload(c, discordClient, lib)
		s.Require().Truef(loaded, "Discord must load its own reserved library %s: %s", lib, out)
	}
	s.Require().NotContains(s.readDaemonLog(c), "op=LIBLOAD", "no load may have been refused")

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Loading a reserved name is sound only if nothing but Discord can bind or rewrite one: every
// other process, including a binary whitelisted for another resource, is refused with op=PLANT.
func (s *IntegrationSuite) TestDaemon_ReservedLib_NonWriterPlantDenied() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.startDiscordLibDaemon(c)

	deep := discordModules + "/a/b"
	code, out := s.exec(c, []string{"mkdir", "-p", deep})
	s.Require().Equalf(0, code, "unreserved directories must stay writable: %s", out)
	s.exec(c, []string{"cp", libProbe, "/tmp/evil.so"})

	s.assertNotPlanted(c, "create *.so", discordModules+"/evil.so",
		[]string{"cp", libProbe, discordModules + "/evil.so"})
	s.assertNotPlanted(c, "create lib*", discordModules+"/libevil.so.1",
		[]string{"sh", "-c", "cat " + libProbe + " > " + discordModules + "/libevil.so.1"})
	s.assertNotPlanted(c, "create at depth", deep+"/evil.node",
		[]string{"sh", "-c", "cat " + libProbe + " > " + deep + "/evil.node"})
	s.assertNotPlanted(c, "rename", deep+"/evil.so", []string{"mv", "/tmp/evil.so", deep + "/evil.so"})
	s.assertNotPlanted(c, "hardlink", deep+"/evil.so", []string{"ln", "/tmp/evil.so", deep + "/evil.so"})
	s.assertNotPlanted(c, "symlink", deep+"/evil.so", []string{"ln", "-s", "/tmp/evil.so", deep + "/evil.so"})
	s.assertNotPlanted(c, "O_TMPFILE+linkat", deep+"/evil.so",
		[]string{"/exploits/glob_plant", "tmpfile", deep, deep + "/evil.so"})
	s.assertNotPlanted(c, "confused deputy", deep+"/evil.node",
		[]string{"/usr/bin/bash", "-c", "cat " + libProbe + " > " + deep + "/evil.node"})
	s.exec(c, []string{"sh", "-c", "mkdir -p /tmp/pkg && cp " + libProbe + " /tmp/pkg/evil.so"})
	s.assertNotPlanted(c, "directory move", discordModules+"/pkg", []string{"mv", "/tmp/pkg", discordModules + "/pkg"})

	lib := discordModules + "/libgood.so"
	s.writeAsDiscord(c, lib)
	for _, tc := range []struct {
		what string
		cmd  []string
	}{
		{"overwrite", []string{"sh", "-c", "cat /usr/bin/true > " + lib}},
		{"append", []string{"sh", "-c", "echo evil >> " + lib}},
		{"truncate", []string{"truncate", "-s", "0", lib}},
		{"rename over", []string{"sh", "-c", "cp " + libProbe + " /tmp/over && mv -f /tmp/over " + lib}},
		{"exchange", []string{"sh", "-c", "cp " + libProbe + " /tmp/swap && /exploits/glob_plant exchange /tmp/swap " + lib}},
		{"hardlink alias", []string{"ln", lib, "/tmp/alias.so"}},
	} {
		code, out = s.exec(c, tc.cmd)
		s.Require().NotEqualf(0, code, "%s of Discord's library by a non-writer must be denied: %s", tc.what, out)
		code, _ = s.exec(c, []string{"cmp", "-s", libProbe, lib})
		s.Require().Equalf(0, code, "%s altered Discord's library", tc.what)
	}
	s.requireDenialLogged(c, "PLANT", "evil.so")
	s.requireDenialLogged(c, "PLANT", "libgood.so")

	for _, p := range []string{deep + "/settings.json", "/tmp/free.so"} {
		code, out = s.exec(c, []string{"sh", "-c", "echo x > " + p})
		s.Require().Equalf(0, code, "%s is not reserved and must stay writable: %s", p, out)
	}

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// The reservation trusts names, not the tree: an unreserved name in the same dir, a reserved name
// outside the root, and Discord's own library mapped into a whitelisted non-writer stay refused.
func (s *IntegrationSuite) TestDaemon_ReservedLib_UnreservedNameRefused() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.startDiscordLibDaemon(c)

	good := discordModules + "/probe.so"
	s.writeAsDiscord(c, good)
	loaded, out := s.preload(c, discordClient, good)
	s.Require().Truef(loaded, "control: the reserved library must load: %s", out)

	unreserved := discordModules + "/probe.dat"
	s.writeAsDiscord(c, unreserved)
	loaded, out = s.preload(c, discordClient, unreserved)
	s.Require().Falsef(loaded, "an unreserved name below the root must not load: %s", out)
	s.requireDenialLogged(c, "LIBLOAD", "probe.dat")

	s.exec(c, []string{"sh", "-c", "mkdir -p /root/.config/other && chown 65534 /root/.config/other"})
	outside := "/root/.config/other/probe.so"
	s.writeAsDiscord(c, outside)
	loaded, out = s.preload(c, discordClient, outside)
	s.Require().Falsef(loaded, "a reserved name outside the root must not load: %s", out)

	loaded, out = s.preload(c, "/usr/bin/bash", good)
	s.Require().Falsef(loaded, "a whitelisted non-writer must not load Discord's library: %s", out)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// probeExeAs reads /proc/<pid>/exe with a copy of readlink at bin (its inode is the caller's
// identity), reporting the command's own status: docker-exec exit codes lie about denials.
func probeExeAs(bin, pid string) string {
	return "exe=$(" + bin + " /proc/" + pid + "/exe 2>&1); echo \"rc=$? exe=$exe\""
}

// A binary two resources whitelist (Steam's client: config + registry.vdf) taints its process with
// the set of exactly those two: a caller both whitelist may inspect it, a caller only one does may
// not. Judging it by every resource's intersection instead let no Steam process inspect another
// and Steam's UI never came up.
func (s *IntegrationSuite) TestDaemon_TaintSet_JudgedByItsOwnResources() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	// readlink copies keep their basename: ubuntu's coreutils is one multi-call binary.
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /r1 /r2 /r3 /etc/app-listener /opt/app /opt/y /opt/pboth /opt/pone /opt/pnone" +
			" && echo s1 > /r1/secret && echo s2 > /r2/secret && echo s3 > /r3/secret" +
			" && cp /usr/bin/dash /opt/app/app && cp /usr/bin/dash /opt/y/y" +
			" && for p in pboth pone pnone; do cp /usr/bin/readlink /opt/$p/readlink; done" +
			" && mkfifo /tmp/f /tmp/g"})
	s.startDaemon(c, `[watch /r1]
need_encryption: false
/opt/app/app
/opt/pboth/readlink
/opt/pone/readlink

[watch /r2]
need_encryption: false
/opt/app/app
/opt/pboth/readlink

[watch /r3]
need_encryption: false
/opt/y/y`)

	// The victim reads both resources (each read must keep the set, not widen it to GLOBAL), then
	// idles in a builtin so its image stays /opt/app/app.
	s.exec(c, []string{"sh", "-c",
		"(/opt/app/app -c 'read a < /r1/secret; read b < /r2/secret; read y < /tmp/f; exit 0' &); sleep 1"})
	_, pidOut := s.exec(c, []string{"sh", "-c", "pgrep -f 'app -c read a' | head -1"})
	pid := strings.TrimSpace(pidOut)
	s.Require().NotEmptyf(pid, "the victim did not start: %q", pidOut)

	_, out := s.exec(c, []string{"sh", "-c", probeExeAs("/opt/pboth/readlink", pid)})
	s.Require().Containsf(out, "rc=0 exe=/opt/app/app",
		"a caller every tainting resource whitelists must inspect the process: %s", out)
	for _, p := range []string{"/opt/pone/readlink", "/opt/pnone/readlink"} {
		_, out = s.exec(c, []string{"sh", "-c", probeExeAs(p, pid)})
		s.Require().NotContainsf(out, "exe=/opt/app/app",
			"%s is not whitelisted by every resource the process holds: %s", p, out)
		s.Require().Containsf(out, "rc=1", "%s: the probe must report its own refusal: %s", p, out)
	}

	// Content from a resource outside the set (/r3, read before exec'ing the set's binary) can't
	// be covered by it: the merge must fall back to GLOBAL, where pboth is not enough.
	s.exec(c, []string{"sh", "-c",
		"(/opt/y/y -c 'read z < /r3/secret; exec /opt/app/app -c \"read y < /tmp/g; exit 0\"' &); sleep 1"})
	_, pidOut = s.exec(c, []string{"sh", "-c", "pgrep -f 'app -c read y < /tmp/g' | head -1"})
	mixed := strings.TrimSpace(pidOut)
	s.Require().NotEmptyf(mixed, "the mixed victim did not start: %q", pidOut)
	_, out = s.exec(c, []string{"sh", "-c", probeExeAs("/opt/pboth/readlink", mixed)})
	s.Require().NotContainsf(out, "exe=/opt/app/app",
		"content from a resource outside the set must not be inspectable by the set's callers: %s", out)

	s.Require().Contains(s.readDaemonLog(c), "op=PTRACE", "the refusals must be logged")
	s.exec(c, []string{"sh", "-c", "echo > /tmp/f; echo > /tmp/g; pkill -f 'app-listener daemon' || true"})
}
