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

// launchDaemon backgrounds the daemon process against the already-written
// /etc/app-listener/daemon.conf, without waiting for it to come up. Split out
// of startDaemon so a caller that needs to race something against the exact
// startup window (see daemon_toctou_test.go) can start its racer immediately
// before this call and stop it once awaitDaemonUp returns.
func (s *IntegrationSuite) launchDaemon(c testcontainers.Container) {
	cmd := "nohup /app-listener daemon --config /etc/app-listener/daemon.conf --headless > /tmp/daemon.log 2>&1 &"
	code, out := s.exec(c, []string{"sh", "-c", cmd})
	s.Require().Equalf(0, code, "starting daemon: %s", out)
}

// awaitDaemonUp polls until the daemon's guards are attached and its event
// readers are running for EVERY resource in config: the pid file, then one
// "guard started — guarding: <path>" log line per "[watch <path>]" section.
// Waiting for only the first such line (as this used to) is a real race on a
// multi-resource config: each guard (re)attach is its own multi-second
// BPF-verifier pass on a slow host (see daemonShutdownTimeout), so a caller
// that immediately reads/writes a LATER resource can hit it before its guard
// — let alone its fscrypt unlock — is actually live. config may be "" (a
// caller re-launching against an already-written daemon.conf it doesn't have
// handy); that falls back to the old "at least one guard started" check.
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

// TestDaemon_SelfWhitelist_NoUniversalKey is the regression test for the
// self-key bypass: the daemon whitelists its own executable inode on every
// guarded resource so its fscrypt ioctls keep working. Identity used to be
// keyed by exe inode alone, with no uid binding, so ANY local user executing
// the same binary file ran with full allow rights on the whole tree — e.g.
// `app-listener network-monitor <guarded-file>` reads the file to hash it
// before any capability check. The self entry must only grant access to
// uid 0 (the daemon itself) and only for the events it needs (OPEN, READ).
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

// ---------------------------------------------------------------
// Bypass POC (finding #1): replacing a whitelisted binary at a
// user-writable path via rename gets the REPLACEMENT re-whitelisted.
//
// Whitelist identity is the exe inode. The daemon re-stats whitelisted
// binaries and re-admits whatever inode is now at the path
// (ReSyncBinaries, triggered on any denial and by the periodic sweep),
// with no check that the new file is the same binary, is root-owned, or
// sits in a non-user-writable directory. The in-place (same-inode)
// swap is caught by the hash-verify loop, but a RENAME replaces the
// path with a NEW inode — a different code path that only compares the
// stored inode, so the malicious replacement is admitted with full
// whitelist rights. Any binary whitelisted from a home directory
// (Claude Code in ~/.local/share/claude, Discord in ~/.config/discord)
// can be swapped by malware running as the same user.
//
// RED now: after the rename the replacement is re-whitelisted and reads
// the guarded secret. After the fix (re-admission must verify the
// binary — by hash, or refuse non-root-owned / user-writable paths),
// the replacement stays denied.
// ---------------------------------------------------------------
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

// TestDaemon_Bypass_LdPreloadWhitelistedBinary is the #2 regression: a
// whitelisted binary run with LD_PRELOAD pointing at an attacker .so under a
// user-writable path (/tmp) must NOT be able to map that code — the trust
// guard's library-load allowlist denies it (the .so is neither an allow_lib
// entry, an auto-trusted root-owned system library, nor inside a guarded
// tree), so the preload constructor never runs and the secret is never
// exfiltrated. A plain run of the same binary still works (its real libraries
// are root-owned system libs, auto-trusted). Enforcement is always on.
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

// ---------------------------------------------------------------
// Bypass POC (finding #4c): the memory-read taint is lost across a
// SIGHUP reload.
//
// A process that read a guarded file is tainted so process_vm_readv
// against it is denied. A reload rebuilds every guard with fresh BPF
// maps, so the tainted-pid set is emptied while the victim keeps running
// with the secret still in memory — the attacker can then dump it.
//
// RED now: after the reload the attacker dumps the secret from the
// still-live victim. After the fix (taint must survive a reload, e.g. by
// re-seeding from the surviving map or persisting it) the dump stays
// denied.
// ---------------------------------------------------------------
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

	// Reload rebuilds every guard with fresh, empty taint maps. Wait for the
	// reload to actually COMPLETE before the post-reload attempt: the
	// "guard started" markers awaitDaemonUp watches are already present from
	// the initial start, so only the reload-complete line proves the map
	// swap happened (otherwise the attacker races the still-tainted old
	// guard and the test passes for the wrong reason).
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

// daemonPinFiles lists the daemon's pin files under /sys/fs/bpf/app-listener.
func (s *IntegrationSuite) daemonPinFiles(c testcontainers.Container) []string {
	_, out := s.exec(c, []string{"sh", "-c", "ls -1 /sys/fs/bpf/app-listener 2>/dev/null || true"})
	return strings.Fields(out)
}

// ---------------------------------------------------------------
// Test: the daemon protects its own on-disk state
//
// While the daemon runs it guards /etc/app-listener independent of any
// [watch] section:
//   - daemon.conf stays world-READABLE, but no process other than the
//     app-listener binary may write / rename-over / delete it, nor create
//     new files in the directory (ModeReadOnly);
//   - fscrypt.key is not even readable except by the app-listener binary
//     (ModeWhitelist, empty list — the key guard stacks on the RO dir guard).
//
// The app-listener binary itself (uid 0) keeps full access so install /
// --genkey / --update-catalog-only work.
// ---------------------------------------------------------------
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

	// 0. `systemctl reload` runs helper processes (ExecReload=/bin/kill …)
	// inside the unit's mount namespace, and setting that up bind-mounts the
	// guarded directory. A mount onto /etc/app-listener must NOT be blocked
	// in read-only mode or the helper dies with 226/NAMESPACE and reload
	// fails (regression: the RO guard used to deny sb_mount here).
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

// ---------------------------------------------------------------
// Test: raw block-device gate is daemon-wide and device-granular
//
// Two resources (/mnt/data/guardedA, /mnt/data/guardedB) share one backing
// block device (a loop-mounted ext4 at /mnt/data). The gate must:
//   - be stamped ONCE, not once per resource ("blocking raw access to
//     backing block device" appears a single time);
//   - block raw reads of the device even for an UNGUARDED path on it
//     (/mnt/data/open/*) — coarse by design, documented, accepted;
//   - log the denial as resource=raw-block-device, never as one of the
//     watched paths (the old bug attributed it to a random resource).
//
// Skips if the environment cannot provide a loop device.
// ---------------------------------------------------------------
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

// ---------------------------------------------------------------
// Test: edit-protected live mode (issue #40)
//
// With an edit-protected password configured, the daemon exposes a local
// control socket. It must:
//   - reject a peer that is not root (SO_PEERCRED);
//   - reject a wrong password (and, after enough failures, lock out — not
//     asserted here to keep the test quick);
//   - on a correct password from a root app-listener peer, briefly widen the
//     target resource's guard so `edit-protected --put` can write, then
//     narrow it again (GRANTED / REVOKED both logged);
//   - allow only one live session at a time.
//
// ---------------------------------------------------------------
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

// ---------------------------------------------------------------
// TOCTOU races across the daemon's lifecycle: start, SIGHUP reload,
// graceful stop, and SIGKILL — for both a real kernel-fscrypt directory and
// a real file-vault (single regular file) resource.
//
// The invariant under test throughout is the one internal/usecase/daemon.go
// documents everywhere: a resource is never readable in plaintext by an
// unauthorized (non-whitelisted) process at any observable instant,
// regardless of exactly when an attacker's read races the daemon's start,
// reload, stop, or a hard kill. The daemon's own ordering guarantee is
// attach -> unlock -> populate -> resolve -> re-sync (see daemon.go's doc
// comment); a SIGKILL cannot be caught, so the second half of the
// guarantee is the guard's LSM links staying pinned to /sys/fs/bpf, and
// `daemon --lockdown` (the systemd ExecStopPost safety net) being able to
// fully recover from ANY point in that ordering, including mid-transform
// of a file-vault resource (see fscrypt/filevault.go's recovery sidecar).
//
// These tests emulate an attacker racing the daemon rather than reasoning
// about the code: a background goroutine hammers the resource as an
// unprivileged user throughout the transition window and the test fails if
// it ever observes plaintext. Real Vault.Encrypt/IsEncrypted production
// code (never a reimplementation) sets resources up and inspects them,
// via a small root-capable exec harness compiled from internal/fscrypt —
// see harness_test.go and main_test.go's fscryptTestAmd64Bin, mirroring
// guardTestAmd64Bin's existing pattern in guard_test.go.
//
// The directory case needs a real fscrypt-capable filesystem (ext4 with
// the `encrypt` feature flag, kernel CONFIG_FS_ENCRYPTION support) that a
// container's overlayfs root cannot provide; setupFscryptDirFilesystem
// builds one on a loop device and skips the test when the environment
// cannot provide it, mirroring TestDaemon_RawBlockDevice_DeviceScope's own
// escape hatch. The file-vault case needs nothing beyond the master key,
// so it always runs.
// ---------------------------------------------------------------

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

// harnessProbe execs one subtest of the fscrypt harness binary against path
// without asserting success. The harness binary is never a whitelisted
// binary on any guarded resource in these tests, so a still-enforcing
// guard (its BPF LSM link can remain pinned and active for a window after
// a SIGKILL, by design — see the architecture note on ExecStopPost) will
// correctly deny it; callers that only want a best-effort diagnostic (not
// a setup-time guarantee) should use this instead of harnessRun, since
// that denial is itself proof the resource is NOT exposed, never a failure.
func (s *IntegrationSuite) harnessProbe(c testcontainers.Container, subtest, path string) (int, string) {
	return s.exec(c, []string{"sh", "-c", fmt.Sprintf(
		"APPLISTENER_HARNESS_PATH=%s /fscrypt.harness -test.run %s -test.v",
		shQuote(path), shQuote("^TestFscryptHarness/"+subtest+"$"))})
}

// harnessMigrate seals path into its encrypted-at-rest form — a real kernel
// fscrypt policy for a directory (the filesystem prerequisites are applied
// programmatically first), the userspace file-vault format for a regular
// file — via the real Vault.Encrypt production code path: exactly what
// `install` does before the daemon ever attaches a guard.
func (s *IntegrationSuite) harnessMigrate(c testcontainers.Container, path string) {
	s.harnessRun(c, "TestMigrate", path)
}

// assertFileVaultSealed confirms, via the already-whitelisted grep binary
// rather than the fscrypt harness, that a regular file-vault path's raw
// bytes no longer expose marker in the clear — the check `daemon
// --lockdown` completed its job after a SIGKILL. The harness binary is
// never whitelisted on the resource, so using it here would be wrong: its
// own guard's pinned BPF link can still be actively enforcing at this
// exact instant (lockdown widens self-access to relock the vault, it does
// not unpin the link — see runLockdown's doc comment), so a harness
// call would get denied regardless of whether the file is actually sealed,
// exactly the false-failure assertNoUnauthorizedPlaintext already guards
// against. Ciphertext reproducing an AEAD-sealed marker verbatim is not a
// realistic possibility, so grep finding nothing is real proof of sealing,
// and grep finding it is a genuine "lockdown failed to re-lock" failure —
// not a false positive from an unrelated denial.
func (s *IntegrationSuite) assertFileVaultSealed(c testcontainers.Container, path, marker, context string) {
	code, out := s.exec(c, []string{"grep", "-c", marker, path})
	s.Require().Falsef(code == 0 && strings.Contains(out, "1"),
		"%s: %s still exposes its plaintext marker after lockdown — not sealed", context, path)
}

// harnessIsEncrypted reports path's current on-disk encryption state via
// the real Vault.IsEncrypted, without mutating anything. Only safe to call
// when no guard could still be pinned-and-enforcing against the harness
// binary itself — i.e. before any daemon has started, or after a clean,
// unraced shutdown (which always unpins before the process exits). After a
// SIGKILL, use assertFileVaultSealed instead.
func (s *IntegrationSuite) harnessIsEncrypted(c testcontainers.Container, path string) bool {
	return strings.Contains(s.harnessRun(c, "TestIsEncrypted", path), "ENCRYPTED")
}

// resetFileVaultTarget clears any leftover backup/recovery sidecar from a
// previous round of a kill-loop test (safe: no guard is watching yet at
// this point — the "never delete, only rewrite in place" rule in
// filevault.go is what the running guard enforces, not a rule these test
// fixtures need to follow before any guard exists) and rewrites path as
// fresh plaintext containing marker, ready for a new harnessMigrate call.
func (s *IntegrationSuite) resetFileVaultTarget(c testcontainers.Container, path, marker string) {
	cmd := fmt.Sprintf("rm -f %s %s.app_listener.backup %s.app_listener.recover && printf '%%s' %s > %s && chmod 600 %s",
		path, path, path, shQuote(marker), path, path)
	code, out := s.exec(c, []string{"sh", "-c", cmd})
	s.Require().Equalf(0, code, "resetting file-vault target %s: %s", path, out)
}

// setupFscryptDirFilesystem creates a small ext4 loop-mounted filesystem
// with the fscrypt `encrypt` feature flag set at mountPoint, so a directory
// on it can carry a real kernel fscrypt policy (unlike the container's
// overlayfs root). Skips the calling test when the environment cannot
// provide a working loop device / encrypt-capable ext4 (older e2fsprogs, a
// kernel without CONFIG_FS_ENCRYPTION), mirroring
// TestDaemon_RawBlockDevice_DeviceScope's own escape hatch.
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

// daemonShutdownTimeout is the budget awaitDaemonDead callers give a
// graceful SIGTERM (Stop locks every vault back — see internal/usecase/daemon.go
// — before the process exits) or a kill racing one. Each guard (re)attach is a
// real BPF-verifier pass over 23 LSM hooks, observed taking 5-11s on a slower
// host/kernel (e.g. a hardened kernel's extra verifier work); a reload or a
// live-edit session can rebuild several guards (including the two
// self-protection ones, selfguards.go) back to back, so the wait comfortably
// covers a handful of those in sequence rather than the couple of seconds a
// fast host needs.
const daemonShutdownTimeout = 45 * time.Second

// awaitDaemonDead polls until no live app-listener process remains,
// returning false if one is still alive after timeout. Uses
// noLiveAppListenerProcs (pool_test.go) — matching on `comm`, not the full
// cmdline — rather than `pgrep -f 'app-listener daemon'`: run via `sh -c
// "pgrep -f 'app-listener daemon' ..."`, that pattern string is itself
// embedded in the wrapping shell's own /proc/<pid>/cmdline, so `pgrep -f`
// (which searches the full cmdline of every process) matches the wrapping
// shell that is running the check and reports "alive" forever regardless of
// the real daemon's state.
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

// runLockdown execs `daemon --lockdown`, the systemd ExecStopPost safety
// net: force-locks every encryption root in /etc/app-listener/daemon.conf
// and exits. It is documented to always exit 0 (best-effort per root, so it
// never blocks the unit from settling) — callers must check the real
// on-disk state via harnessIsEncrypted, not this exit code, to know whether
// it actually succeeded.
func (s *IntegrationSuite) runLockdown(c testcontainers.Container) (int, string) {
	return s.exec(c, []string{"/app-listener", "daemon", "--lockdown"})
}

// rawExec is s.exec without the fatal assertions on transport errors: safe
// to call from a background goroutine (calling testify's Require/FailNow
// from a non-test goroutine is undefined per the testing package's own
// contract). A transport hiccup is treated as "nothing observed this
// iteration", which is exactly the right behavior for a racer polling in a
// tight loop.
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

// raceUnauthorizedReader repeatedly attempts, as an unprivileged user, to
// read path's content until stop fires, setting leaked when an attempt
// both succeeds (exit 0) and returns content containing marker — a genuine
// plaintext exposure to an unauthorized reader. A resource that is still
// locked (ciphertext, or unreadable before the daemon has even unlocked
// it) never matches marker, and one properly guarded once unlocked is
// denied outright — this only fires on the exact failure these tests exist
// to catch.
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

// assertNoUnauthorizedPlaintext execs one direct (non-racing) check that an
// unprivileged reader cannot currently see marker at path, regardless of
// whether the guard is even attached right now (e.g. right after a
// SIGKILL) — the resource must be safe by construction (still ciphertext,
// or ciphertext-but-orphan-guarded), never "safe only while the guard
// process happens to still be alive".
func (s *IntegrationSuite) assertNoUnauthorizedPlaintext(c testcontainers.Container, path, marker, context string) {
	// Best-effort only: the guard's BPF LSM link can still be pinned and
	// actively enforcing at this exact instant (a SIGKILL landing before
	// Stop() reaches its own unpin step leaves it orphaned-but-active
	// until `daemon --lockdown` runs) — in which case the harness binary,
	// itself never whitelisted on this resource, is correctly denied by
	// the very protection this test exists to prove. That denial is
	// stronger evidence of "not exposed" than any encrypted/plaintext
	// report could be, so it must never be escalated into a hard failure
	// here — only the unprivileged-reader check below is the real assertion.
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

// TestDaemon_StartupRace_NoPlaintextWindow_Directory races an unprivileged
// reader against the daemon's startup on a real, kernel-fscrypt-encrypted
// DIRECTORY resource: attach -> unlock -> populate must never leave a
// window where the directory's key is provisioned but the guard is not yet
// enforcing.
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

// TestDaemon_StartupRace_NoPlaintextWindow_GroupedResources is the regression
// test for the daemon startup fan-out (buildGuards attaching guards
// concurrently, then startGuards unlocking encryption roots and preparing
// guards concurrently — see internal/usecase/daemon.go and
// cmd/functions/daemon/daemon.go): TWO `watch:` sub-paths share ONE
// encryption root, exactly the shape a browser-profile catalog entry uses
// (Local Storage + Cookies under one profile vault). The root's key is
// unlocked exactly once for the whole group (uniqueEncryptionRoots dedup),
// so unlocking it makes BOTH sub-paths' plaintext content readable at the
// kernel level at the same instant — the guard for EACH sub-path must
// already be attached by then, or the sibling whose own guard lags behind
// is exposed with nothing denying access to it. A naive per-resource
// pipeline (attach+unlock+populate+resolve+start for one resource, fully
// independent of its group siblings) would reintroduce exactly this window;
// the concurrency added to buildGuards/startGuards keeps the invariant by
// construction instead: buildGuards attaches EVERY resource's guard —
// including every member of this group — as one completed fan-out phase
// before startGuards.Start() ever calls Unlock, and unlockRoots only fans
// out over already-deduplicated unique roots, never issuing two concurrent
// unlocks for the one root this group shares.
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
	// watchPathsInConfig only understands a `[watch <path>]` section header as
	// the guarded resource, which is wrong for a group: the header names the
	// ENCRYPTION ROOT, never itself guarded (see TestLoadWatchGroup), so the
	// per-path "guard started — guarding: <path>" markers awaitDaemonUp would
	// derive from the real config never appear for it. Passing "" falls back
	// to "pid file exists" as the sole readiness signal, which is still
	// exact: writePidFile only runs after startGuardedDaemon's call to
	// d.Start() returns, which does not return until EVERY configured
	// guard — both group members included — has been through the full
	// attach/unlock/populate/resolve/start pipeline.
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

// ---------------------------------------------------------------
// C. SIGKILL landing somewhere inside startup's unlock/populate window,
// across several delays, for the file-vault resource — the newest, most
// crash-sensitive code path (in-place AEAD transform + recovery sidecar).
// ---------------------------------------------------------------

// TestDaemon_KillDuringUnlock_FileVault_NeverOrphansPlaintext repeatedly
// SIGKILLs the daemon at increasing delays after launch — early enough to
// plausibly land inside attach/unlock/populate at least once — and asserts
// that whatever state the kill caught it in, the file is never both
// plaintext and reachable by an unauthorized reader; that `daemon
// --lockdown` (the ExecStopPost safety net) can always fully recover from
// it; and that a subsequent clean restart round-trips the content without
// corruption (proving the recovery sidecar in filevault.go actually does
// its job on a real interrupted transform, not just in code review).
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

// ---------------------------------------------------------------
// D. Live catalog-refresh SIGHUP reload (whitelist change on an existing
// resource): both the ordinary (uninterrupted) race window and a kill
// landing mid-reload.
// ---------------------------------------------------------------

// installFakeSystemctl stubs a minimal /usr/local/bin/systemctl inside the
// container. These test images run no systemd at all — the daemon here is
// always launched directly via nohup, never as a unit — but install's LIVE
// catalog refresh path (internal/systemd.IsDaemonActive / EnableAndVerify)
// unconditionally shells out to the real systemctl to gate and deliver the
// change. The stub covers exactly the subcommands that path calls:
// is-active/is-enabled report the daemon as up (so the live-mode gate
// passes and no unwanted enable/start/restart branch fires), daemon-reload
// is a no-op, and reload is turned into what it would really do on a
// systemd host — deliver an actual SIGHUP to the running daemon via its
// pid file — which is what makes the resulting reload real and worth
// racing/killing against.
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

// TestDaemon_KillDuringReload_NoUnprotectedWindow exercises Reload's own
// documented ordering (the new guard attaches before the old one detaches)
// for a resource whose WHITELIST changes via a live SIGHUP reload: first
// with a live racer through a normal, uninterrupted reload, then with a
// SIGKILL landing mid-reload — in both cases the resource must never end up
// fully unguarded (open to any reader).
//
// The trigger is `install --update-catalog-only --live --yes` against the
// built-in WireGuard catalog entry (/etc/wireguard, whitelisting
// /usr/bin/nmcli), not a raw shell rewrite of daemon.conf: the self-guard on
// /etc/app-listener (see issue #53's follow-up hardening) now
// deterministically denies any non-daemon-binary write there, so a brand
// new [watch] section can never be introduced by hand-editing the config
// while the daemon runs — and per README.md, adding a genuinely NEW
// resource (`install` / `install --diff-catalog`) always stops the daemon
// first anyway, live or not. `--update-catalog-only --live` is the one
// documented live path: it re-expands an EXISTING catalog-matched section's
// whitelist from inside the daemon's own binary (GUARD_ALLOW_ROOT — same
// exe inode as the running daemon) and delivers the change via SIGHUP,
// without any external process ever writing to the guarded config.
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
	// installNmcli simulates the real-world trigger for this workflow: a
	// package manager installs a catalog-whitelisted binary after the
	// daemon already started guarding the resource with an empty (nothing
	// matched yet) whitelist. FilterExistingWhitelist only picks up binaries
	// that exist on disk at scan time. This must be a real, standalone ELF
	// binary, not a shebang script: guard identity is the CALLING PROCESS's
	// own exe inode, and a script's process runs under its interpreter's
	// inode (/bin/sh), never the script file's own — so a script here would
	// always be denied regardless of whitelisting. /bin/cat (and friends)
	// won't do either on modern Ubuntu: they're symlinks into one uutils
	// coreutils multi-call binary that dispatches on ITS OWN RESOLVED PATH's
	// basename (not argv[0] — `exec -a` doesn't help), so a copy landing at
	// a path named "nmcli" is rejected as an unknown applet regardless.
	// /usr/bin/grep is GNU grep, a genuine standalone binary indifferent to
	// its own path/name — copying it works as a stand-in reader.
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

// ---------------------------------------------------------------
// E. Shutdown interrupted mid-lockdown: SIGTERM (starts the secure
// lockdown — lock before detach) immediately followed by SIGKILL
// (simulating systemd's stop-timeout force-kill, or an impatient admin).
// ---------------------------------------------------------------

// TestDaemon_KillDuringShutdown_LockdownRecovers races a SIGKILL against
// the daemon's own graceful-shutdown lockdown sequence (Stop: lock vault
// before guard detach) and verifies `daemon --lockdown` can still always
// finish the job, without the content ever being both plaintext and
// unguarded, and without silent corruption on the next start.
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

// ---------------------------------------------------------------
// F. Live-edit-grant killed mid-session: the transient self-access
// widening (CLAUDE.md's "live-edit grant is a real escalation path") must
// never survive as a standing bypass usable by anyone but the exact
// trusted actor it was scoped to, and must never survive a restart at all.
// ---------------------------------------------------------------

// TestDaemon_LiveEditGrant_KilledMidSession_NoResidualEscalation opens a
// real live edit-protected session (AUTH+SELECT over the control socket,
// widening the daemon's own self-access on the target resource — see
// guard.Guard.GrantSelfEditAccess) and races a SIGKILL against it before it
// ever ends cleanly.
//
// The client must be the real app-listener binary (the control socket's
// authPeer check refuses any other exe inode outright — CLAUDE.md's
// "peer-exe check"), so the session's own write step cannot be slowed from
// outside; instead its --put content is made large (512 MiB) so the
// write+fsync between GRANTED and the client sending END takes long enough
// for a tight log-poll-then-kill loop to land inside that window on a
// reasonable fraction of runs. This is a best-effort race — like the
// existing TestGuard_Bypass_PtraceRace — but the assertions below hold
// regardless of whether the kill actually landed mid-grant or just after a
// clean revoke: the widened mask is scoped to exactly (uid 0, the daemon's
// own exe inode), so a non-root writer and a root shell that is NOT the
// app-listener binary must both stay denied either way (CLAUDE.md's
// "identity stays inode-based" invariant, exercised here against a live
// orphaned grant rather than a static whitelist entry), and the resource's
// other, unrelated content must never be touched. Finally, a fresh daemon
// instance must come back at the read-only baseline with no residual
// session — the grant lives only in the killed process's (and its orphaned
// pinned BPF map's) state, never on disk.
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

// ---------------------------------------------------------------
// F. Recovery-sidecar symlink-follow: the file-vault crash-recovery sidecar
// (internal/fscrypt/filevault.go's <resource>.app_listener.recover) must
// never be opened, read, or created through a symlink.
// ---------------------------------------------------------------

// TestDaemon_FileVault_RecoverySidecarSymlink_NeverFollowed proves a
// link-following bug in the file-vault crash-recovery sidecar: an attacker
// who fully controls the parent directory of a file-vault resource (a real
// precondition — the catalog's file-vault targets, e.g. a Steam
// registry.vdf, live directly inside a user's own home directory) can,
// while the daemon is not running, delete the not-yet-created sidecar and
// replace it with a symlink pointing at any existing file that only root
// can write. On the next daemon start, the unlock cycle's stageRecovery
// step opens that symlink with O_RDWR and writes sealed file-vault
// ciphertext straight into whatever it resolves to — a file completely
// outside the guarded resource and outside the attacker's own write
// permissions.
//
// To make the corruption deterministic instead of racing a timing window,
// the resource's own directory (/protected) is remounted read-only right
// after the malicious symlink is staged: this leaves the harmful write
// (through the symlink, into a file OUTSIDE /protected) unaffected, but
// makes the very next step — transformFileInPlace rewriting the real
// resource file, which lives under /protected — fail deterministically
// with EROFS. That failure aborts the unlock before clearRecovery (the
// step that would otherwise truncate the sidecar, and thus the symlinked
// victim file, back to empty) ever runs, so the corruption persists and can
// be asserted on directly.
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

	// The attack: the resource's parent directory is fully attacker-owned
	// (chmod 777 stands in for that — the real precondition is the
	// directory being the user's own, e.g. their home directory), so an
	// unprivileged user can freely create the not-yet-existing sidecar
	// themselves, as a symlink to a root-writable file elsewhere.
	// Remounting /protected read-only afterward only pins down the
	// deterministic-failure trick described above — it plays no part in
	// the vulnerability itself.
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
	// The daemon is expected to fail closed here (the resource's own
	// directory is read-only and its sidecar is hostile) rather than come
	// up — poll for either outcome without hard-failing, since the real
	// assertion is what happened to the symlink's target, not whether the
	// daemon itself started.
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

// ---------------------------------------------------------------
// G. Stale watch-root inode after in-place recreation (steady state, not
// startup): a guarded resource is deleted and recreated in place — an app
// rebuilding its own directory/file, an fscrypt migration, a backup restore
// — while the daemon keeps running. The kernel-side root-confinement anchor
// (guard_config[3..4], consulted by root_in_chain in guard.bpf.c) and
// g.rootKey must follow the new inode via the periodic SweepInodes, or two
// things go wrong at once: the real recreated resource silently stops being
// guarded (its ancestor chain no longer contains the stale anchor), and
// whatever unrelated path later ends up holding the freed OLD inode number
// gets denied purely by coincidence, misattributed to this resource. See
// internal/guard/guard.go's SweepInodes doc comment and
// TestSweepInodesRecreatedFileRoot / TestSweepInodesRecreatedDirRoot (unit
// tests, internal/guard/guard_test.go) for the same regression at the
// BPF-map level; these two exercise it through a real running daemon
// instead, waiting out the real periodic sweep.
// ---------------------------------------------------------------

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

// TestDaemon_StaleRootInode_Directory_RecreatedRootReguarded is the
// directory-root regression test for the bug fixed in SweepInodes: before
// the fix, only a single-file watch root's own inode change was detected and
// re-anchored — a directory root's own identity was never re-checked (only
// its mtime, to decide whether to re-walk top-level entries), so a
// wholesale directory recreation left guard_config/g.rootKey pointed at the
// freed old inode forever.
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

	// rename(2) preserves (dev, ino): moving the watch root itself out
	// leaves its old inode living at /old-root — exactly the state a
	// filesystem produces when it later reuses that freed inode number for
	// an unrelated path. A fresh mkdir at /watch then gets a genuinely new
	// inode (the parent, /, is not guarded, so this mkdir is unconditionally
	// allowed regardless of whitelist).
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

// TestDaemon_StaleRootInode_File_RecreatedRootReguarded is
// TestDaemon_StaleRootInode_Directory_RecreatedRootReguarded's single-file
// counterpart: the case SweepInodes already handled before this fix (see
// TestSweepInodesRecreatedFileRoot), exercised here end-to-end through a
// real running daemon instead of a direct unit-level SweepInodes call, so a
// future change to the shared plumbing (updateRootKey, the periodic sweep
// wiring) that breaks either case gets caught by both.
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

// The guard_path_unlink eviction fix itself (guard.bpf.c) is exercised at the
// guard level, not here: see TestGuard_PathUnlinkEvictsInodeImmediately and
// TestGuard_PathUnlinkDeniedDeleteKeepsGuardedInode in
// integrationtests/guard_test.go, alongside the rest of the
// TestGuard_InodeReuse_* family this fix complements.

// TestDaemon_JitProvenanceTrustsOnlySelfCreatedCode covers the trust guard's
// provenance rule for runtime-generated code (guard_jit_origin in
// guard_trust.bpf.c). GPU drivers JIT into a file they create themselves and
// then map executable — a memfd, an O_TMPFILE, or an mkstemp()'d file — and no
// path or ownership rule can ever cover those inodes, so the daemon trusts an
// exec mapping when the mapping process image created the inode and nothing
// else has written it.
//
// The rule is an ALLOWANCE, so most of this test is about what it must still
// refuse: an inode the attacker created (no provenance at all), a file that is
// merely new, an inode handed over through execve (same pid, different image),
// and one that a foreign writer touched after creation.
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

// TestDaemon_ReplacedSingleFileRootIsGuardedQuickly: an application's atomic
// save (write a temp, rename it over the guarded file — Steam's registry.vdf
// on every launch) gives the watch root a new inode that is neither the root
// nor in guard_inodes, i.e. unguarded until the daemon re-anchors. That used
// to wait for the 30 s sweep; single-file roots are now followed every second
// (fileRootFollowEvery). The kernel cannot follow the rename itself:
// guard_path_rename has no verifier budget left (issue #45).
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

// TestDaemon_ReadOnlyGuardDoesNotTaintItsWriters: running a binary that is a
// lib_binary writer of a read-only lib_dir must not taint the process.
// Taint shields a process holding SECRETS from ptrace-class inspection; a
// read-only tree is world-readable code, so it has nothing to shield. The exec
// hook tainted in every guard anyway (unlike the file-access taint sites,
// which skip read-only guards), which left each runtime tree's writer list as
// the only processes allowed to look at /proc/<pid> of half of Steam's
// process tree — and it failed silently.
func (s *IntegrationSuite) TestDaemon_ReadOnlyGuardDoesNotTaintItsWriters() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"sh", "-c",
		"mkdir -p /rt /etc/app-listener && printf 'code' > /rt/lib.so && cp /usr/bin/sleep /tmp/rtwriter && chmod 755 /tmp/rtwriter"})
	s.startDaemon(c, `[libraries "Runtime"]
lib_dir /rt
lib_binary /tmp/rtwriter`)

	s.exec(c, []string{"sh", "-c", "(/tmp/rtwriter 30 &) ; sleep 1"})
	_, pidOut := s.exec(c, []string{"sh", "-c", "pgrep -f '^/tmp/rtwriter 30' | head -1"})
	pid := strings.TrimSpace(pidOut)
	s.Require().NotEmptyf(pid, "the writer process did not start: %q", pidOut)

	// /proc/<pid>/fd listing is a ptrace-class access (PTRACE_MODE_READ).
	code, out := s.exec(c, []string{"sh", "-c", "ls /proc/" + pid + "/fd 2>&1"})
	s.Require().Equalf(0, code,
		"a process running a read-only tree's writer must stay inspectable by an unrelated process: %s", out)

	_, logOut := s.exec(c, []string{"sh", "-c", "grep -c 'op=PTRACE' /tmp/daemon.log || true"})
	s.Require().Equalf("0", strings.TrimSpace(logOut), "no ptrace-class denial may come from a read-only guard")

	s.exec(c, []string{"sh", "-c", "pkill -f '^/tmp/rtwriter' ; pkill -f 'app-listener daemon' || true"})
}

// TestDaemon_TaintFollowsForkButNotExec pins the taint lifecycle:
//   - the process that read the secret is tainted (not inspectable);
//   - a forked child that keeps the parent's memory stays tainted;
//   - a child that execs a NON-whitelisted image is cleared: exec discards
//     the address space, so it holds nothing read from the vault.
// Before the exec rule every descendant of a tainted process stayed tainted
// for life — the whole Steam process tree (Proton, audio tools, the game) —
// and GameMode, PipeWire, the portal and Wine's own wineserver were refused
// every look at it. The refusal must also be logged with its access mode.
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
	inspect := func(pid string) int {
		code, _ := s.exec(c, []string{"sh", "-c", "ls /proc/" + pid + "/fd >/dev/null 2>&1"})
		return code
	}

	s.Require().NotEqualf(0, inspect(pidOf("/tmp/p_reader")),
		"the process that read the secret must not be inspectable by an unwhitelisted process")
	s.Require().NotEqualf(0, inspect(pidOf("/tmp/p_fork")),
		"a forked child sharing the reader's memory image must stay tainted")
	s.Require().Equalf(0, inspect(pidOf("/tmp/p_exec")),
		"a child that exec'd an unwhitelisted image holds none of the vault's memory: its taint must be cleared")

	_, logOut := s.exec(c, []string{"sh", "-c", "grep -c 'op=PTRACE.*mode=READ' /tmp/daemon.log || true"})
	s.Require().NotEqualf("0", strings.TrimSpace(logOut),
		"a /proc/<pid>/fd denial must be logged as op=PTRACE with mode=READ")

	s.exec(c, []string{"sh", "-c", "pkill -f /tmp/wsh; pkill -f 'sleep 60'; pkill -f 'app-listener daemon' || true"})
}

// TestDaemon_OwnMetadataReadableMemoryNot: the daemon is tainted (it reads
// what it guards) and keeps its secrets — the fscrypt key, the edit-auth hash
// — in memory. A ptrace ATTACH-class access to it (/proc/<pid>/mem,
// process_vm_readv) must stay denied; a READ-class one (/proc/<pid>/fd,
// environ — what systemd-journald does to label every log line) is allowed,
// or each denial the daemon logs makes journald look it up and fail again.
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

	// ATTACH-class: still refused — and this also proves the daemon IS
	// tainted, so the READ assertion below is not vacuous (an untainted
	// process would give an I/O error at offset 0, not EPERM).
	_, out := s.exec(c, []string{"sh", "-c", "head -c1 /proc/" + pid + "/mem 2>&1"})
	s.Require().Containsf(out, "not permitted",
		"memory of the daemon must stay unreadable to an unwhitelisted process: %s", out)

	// READ-class: allowed.
	code, out := s.exec(c, []string{"sh", "-c", "ls /proc/" + pid + "/fd >/dev/null 2>&1 && cat /proc/" + pid + "/environ >/dev/null 2>&1"})
	s.Require().Equalf(0, code, "metadata of the daemon's own process must be readable: %s", out)

	_, logOut := s.exec(c, []string{"sh", "-c", "grep -c 'op=PTRACE.*comm=app-listener mode=READ' /tmp/daemon.log || true"})
	s.Require().Equalf("0", strings.TrimSpace(logOut), "no READ-mode denial may be logged for the daemon itself")

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}
