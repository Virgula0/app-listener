package integrationtests

import (
	"context"
	"fmt"
	"io"
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
	s.awaitDaemonUp(c)
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
// readers are running (pid file, then the "guard started" log marker).
func (s *IntegrationSuite) awaitDaemonUp(c testcontainers.Container) {
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
		s.launchDaemon(c)
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
	s.awaitDaemonUp(c)
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
	s.Require().True(s.awaitDaemonDead(c, 15*time.Second), "daemon did not exit after SIGTERM")
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
	s.awaitDaemonUp(c)
	close(stop)

	s.Require().Falsef(leaked.Load(),
		"an unauthorized reader observed plaintext content of %s before the guard was fully attached", secretFile)

	_, grepOut := s.exec(c, []string{"grep", "-c", marker, secretFile})
	s.Require().Contains(grepOut, "1", "whitelisted reader must see the unlocked content")
	code, _ := s.exec(c, []string{"sh", "-c", "setpriv --reuid=65534 --regid=65534 --clear-groups cat " + secretFile})
	s.Require().NotEqualf(0, code, "non-whitelisted reader must still be denied once the guard is fully up")

	s.sigDaemon(c, "TERM")
	s.Require().True(s.awaitDaemonDead(c, 15*time.Second), "daemon did not exit after SIGTERM")
	s.Require().Truef(s.harnessIsEncrypted(c, secretFile), "a graceful stop must re-seal the file")
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
			s.awaitDaemonUp(c)
			_, grepOut := s.exec(c, []string{"grep", "-c", marker, secretFile})
			s.Require().Contains(grepOut, "1", "round %d: content must round-trip intact through the kill + lockdown + restart cycle", i)

			s.sigDaemon(c, "TERM")
			s.Require().True(s.awaitDaemonDead(c, 15*time.Second), "round %d: daemon did not exit after SIGTERM", i)
			s.runLockdown(c)
		})
	}
}

// ---------------------------------------------------------------
// D. SIGHUP reload adding a new encrypted resource: both the ordinary
// (uninterrupted) race window and a kill landing mid-reload.
// ---------------------------------------------------------------

// TestDaemon_KillDuringReload_NoUnprotectedWindow exercises Reload's own
// documented ordering (attach the new guard -> unlock -> populate ->
// resolve -> start, old guards never drop protection in the meantime) for
// a resource ADDED via SIGHUP: first with a live racer through a normal,
// uninterrupted reload, then with a SIGKILL landing mid-reload — in both
// cases the newly added resource must never end up unlocked and
// unprotected.
func (s *IntegrationSuite) TestDaemon_KillDuringReload_NoUnprotectedWindow() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.copyFscryptHarness(c)

	const marker = "TOCTOU-RELOAD-SECRET-7F31"
	const resourceB = "/vault/added.txt"

	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /vault /etc/app-listener && echo BASE > /protected/base.txt && chmod 755 /protected && " +
			"head -c 32 /dev/zero > /etc/app-listener/fscrypt.key && chmod 600 /etc/app-listener/fscrypt.key"})

	baseConfig := "[watch /protected]\nneed_encryption: false\n/usr/bin/grep"
	reloadedConfig := fmt.Sprintf("%s\n\n[watch %s]\nneed_encryption: true\n/usr/bin/grep", baseConfig, resourceB)

	// resourceB must be sealed BEFORE any daemon here ever starts: once one
	// is running, its self-guard denies every non-daemon-binary read of
	// /etc/app-listener/fscrypt.key (including this harness), exactly like
	// a real `install` migration has to run before the daemon is up.
	s.resetFileVaultTarget(c, resourceB, marker)
	s.harnessMigrate(c, resourceB)
	s.Require().True(s.harnessIsEncrypted(c, resourceB), "round 1 setup: must start from a sealed file")

	s.startDaemon(c, baseConfig)

	// ---- Round 1: a normal, uninterrupted reload must never expose
	// plaintext early. ----
	s.exec(c, []string{"sh", "-c", fmt.Sprintf("cat > /etc/app-listener/daemon.conf <<'EOF'\n%s\nEOF", reloadedConfig)})

	_, pidOut := s.exec(c, []string{"sh", "-c", "cat /run/app-listener-daemon.pid"})
	pid := strings.TrimSpace(pidOut)

	stop := make(chan struct{})
	var leaked atomic.Bool
	go s.raceUnauthorizedReader(c, resourceB, marker, stop, &leaked)

	s.exec(c, []string{"sh", "-c", "kill -HUP " + pid})
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
		"an unauthorized reader observed plaintext content of %s during the SIGHUP reload window", resourceB)

	_, grepOut := s.exec(c, []string{"grep", "-c", marker, resourceB})
	s.Require().Contains(grepOut, "1", "the newly added resource must be unlocked and readable by its whitelisted binary after reload")
	code, _ := s.exec(c, []string{"sh", "-c", "setpriv --reuid=65534 --regid=65534 --clear-groups cat " + resourceB})
	s.Require().NotEqualf(0, code, "non-whitelisted reader must be denied once the new guard is live")

	s.sigDaemon(c, "TERM")
	s.Require().True(s.awaitDaemonDead(c, 15*time.Second), "daemon did not exit after SIGTERM")
	s.runLockdown(c)

	// ---- Round 2: a reload interrupted by SIGKILL must still never leave
	// the newly added resource unlocked and unguarded. ----
	s.resetFileVaultTarget(c, resourceB, marker)
	s.harnessMigrate(c, resourceB)
	s.startDaemon(c, baseConfig) // back to just the base resource

	s.exec(c, []string{"sh", "-c", fmt.Sprintf("cat > /etc/app-listener/daemon.conf <<'EOF'\n%s\nEOF", reloadedConfig)})
	_, pidOut = s.exec(c, []string{"sh", "-c", "cat /run/app-listener-daemon.pid"})
	pid = strings.TrimSpace(pidOut)

	s.exec(c, []string{"sh", "-c", "kill -HUP " + pid})
	time.Sleep(150 * time.Millisecond) // aim somewhere inside prepareNewResources/prepareGuards/startNewGuards
	s.sigDaemon(c, "KILL")
	s.Require().True(s.awaitDaemonDead(c, 15*time.Second), "daemon did not die after SIGKILL mid-reload")

	s.assertNoUnauthorizedPlaintext(c, resourceB, marker, "kill mid-reload")

	code, out := s.runLockdown(c)
	s.Require().Equalf(0, code, "lockdown after a kill-mid-reload exited non-zero: %s", out)
	s.assertFileVaultSealed(c, resourceB, marker, "lockdown after kill-mid-reload")

	s.startDaemon(c, reloadedConfig)
	_, grepOut = s.exec(c, []string{"grep", "-c", marker, resourceB})
	s.Require().Contains(grepOut, "1", "content must round-trip intact through the kill-mid-reload + lockdown + restart cycle")

	s.sigDaemon(c, "TERM")
	s.Require().True(s.awaitDaemonDead(c, 15*time.Second), "daemon did not exit after SIGTERM")
	s.runLockdown(c)
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
	s.Require().True(s.awaitDaemonDead(c, 15*time.Second), "daemon did not die after the SIGTERM+SIGKILL race")

	s.assertNoUnauthorizedPlaintext(c, secretFile, marker, "SIGTERM raced by SIGKILL mid-shutdown")

	code, out := s.runLockdown(c)
	s.Require().Equalf(0, code, "lockdown after a raced shutdown exited non-zero: %s", out)
	s.assertFileVaultSealed(c, secretFile, marker, "lockdown after a raced shutdown")

	s.startDaemon(c, config)
	_, grepOut = s.exec(c, []string{"grep", "-c", marker, secretFile})
	s.Require().Contains(grepOut, "1", "content must round-trip intact through the raced shutdown + lockdown + restart cycle")

	s.sigDaemon(c, "TERM")
	s.Require().True(s.awaitDaemonDead(c, 15*time.Second), "daemon did not exit after the final SIGTERM")
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
	s.Require().True(s.awaitDaemonDead(c, 15*time.Second), "daemon did not die after SIGKILL raced against the live-edit session")

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
	s.Require().True(s.awaitDaemonDead(c, 15*time.Second), "daemon did not exit after the final SIGTERM")
}
