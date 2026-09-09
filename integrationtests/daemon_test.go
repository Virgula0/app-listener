package integrationtests

import (
	"fmt"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"

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
	cmd := "nohup /app-listener daemon --config /etc/app-listener/daemon.conf --headless > /tmp/daemon.log 2>&1 &"
	code, out := s.exec(c, []string{"sh", "-c", cmd})
	s.Require().Equalf(0, code, "starting daemon: %s", out)

	deadline := time.Now().Add(20 * time.Second)
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
		s.exec(c, []string{"sh", "-c", cmd})
		deadline = time.Now().Add(20 * time.Second)
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
		s.Require().Failf("daemon did not attach its guards", "pid file missing after 20s, daemon log:\n%s", log)
	}
	// Guards are attached and readers running before the pid file appears;
	// poll for the guard-started marker instead of a fixed settle sleep.
	deadline2 := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline2) {
		if strings.Contains(log, "guard started") || strings.Contains(s.readDaemonLog(c), "guard started") {
			break
		}
		time.Sleep(200 * time.Millisecond)
		log = s.readDaemonLog(c)
	}
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
		"mkdir -p /protected /etc/app-listener && echo SECRET > /protected/secret && chmod 755 /protected && chmod 600 /protected/secret && " +
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

	s.exec(c, []string{"sh", "-c", "pkill -TERM -f 'app-listener daemon' || true"})
}

// shQuote single-quotes s for safe embedding in a /bin/sh -c string.
func shQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
