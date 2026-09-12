package integrationtests

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"
)

// ---------------------------------------------------------------
// File-target counterparts of the guard_test.go bypass corpus.
//
// Every regression above proves a bypass class is closed when the guard's
// watch root is a DIRECTORY. A single-file watch root (fscrypt's file vault,
// edit-protected's single-file resources) walks different BPF code paths —
// GUARD_ALLOW_ROOT is the file's own inode, there is no directory inode in
// guard_inodes, and the parent-directory rename-over-watchroot check in
// guard_path_rename special-cases exactly this shape (see
// TestGuard_Bypass_RenameOverGuardedFile). This file re-runs each bypass
// class from guard_test.go with the guard's watch root pointed directly at
// the file, instead of a directory containing it, so the same guarantee is
// proven for both watch-root shapes.
//
// TestGuard_BlocksAll_File, TestGuard_Whitelist_Binary_File,
// TestGuard_Blacklist_Binary_File, TestGuard_Bypass_OpenByHandleAt and
// TestGuard_Bypass_RenameOverGuardedFile already exercise a single-file
// watch root in guard_test.go and are not duplicated here.
// ---------------------------------------------------------------

// TestGuard_Exploits_File is the file-target counterpart of
// TestGuard_Exploits: every kprobe/LSM-adjacent read bypass (copy_file_range,
// io_uring, mmap, pread64, readv, sendfile, splice) plus execve, run against
// a watch root that IS the target file rather than a directory containing it.
func (s *IntegrationSuite) TestGuard_Exploits_File() {
	c := s.guardContainer()
	// pooled: terminated at suite end

	s.exec(c, []string{"mkdir", "-p", "/exploits"})

	// Copied up front: the loop below skips the shared table's "execve"
	// entry entirely (see the continue below), so its own CopyFileToContainer
	// never runs — the dedicated execve block further down needs the binary
	// staged independently of that loop.
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/execve"), "/exploits/execve", 0755))

	const targetPath = "/watch/exploit_target.txt"
	s.exec(c, []string{"sh", "-c", fmt.Sprintf("echo 'exploit target' > %s", targetPath)})

	// Single-file watch root: only exploit_target.txt's inode is guarded.
	s.startGuardStd(c, targetPath)
	logCursor := s.readGuardLog(c)

	for _, et := range guardExploitTests {
		if et.name == "execve" {
			// execve needs its own executable watch root (below): the
			// single-file guard here watches exploit_target.txt, a plain
			// data file with no +x bit, so running execve against it would
			// be denied by the ordinary MAY_EXEC permission check before
			// ever reaching a guard hook — a false-passing, event-less
			// "block" that isn't actually exercising the guard.
			continue
		}
		if et.name == "io_uring" && !s.checkKernelSupport(c) {
			s.T().Logf("skipping exploit %s (kernel may not support it)", et.name)
			continue
		}

		s.T().Run(et.name, func(t *testing.T) {
			t.Helper()

			exploitHostPath := absPath(fmt.Sprintf("./exploits/%s", et.binary))
			err := c.CopyFileToContainer(s.ctx, exploitHostPath, fmt.Sprintf("/exploits/%s", et.binary), 0755)
			s.Require().NoError(err, "copy exploit binary")

			args := []string{fmt.Sprintf("/exploits/%s", et.binary)}
			args = append(args, et.extraArgs...)
			args = append(args, targetPath)

			code, out := s.exec(c, args)
			if code == 2 && et.name == "io_uring_full" {
				t.Logf("skipping exploit %s (IORING_SETUP_SQPOLL not supported)", et.name)
				return
			}
			s.Require().NotEqualf(0, code,
				"exploit %s should be blocked by a single-file guard, got exit %d: %s", et.name, code, out)

			logAfter := s.readGuardLog(c)
			deltaEvents := guardDeltaEvents(logCursor, logAfter)
			logCursor = logAfter

			for _, expectedType := range et.events {
				s.requireBlockedEvent(deltaEvents, expectedType)
			}
		})
	}

	s.stopGuard(c)

	// execve needs its own watch root: the target must be executable and
	// distinct from the read-exploit target above.
	s.exec(c, []string{"sh", "-c",
		"echo '#!/bin/sh\necho executed' > /watch/.exec_target && chmod +x /watch/.exec_target"})
	s.startGuardStd(c, "/watch/.exec_target")
	logBefore := s.readGuardLog(c)

	code, out := s.exec(c, []string{"/exploits/execve", "/watch/.exec_target"})
	s.Require().NotEqualf(0, code, "execve of a single-file guarded target should be blocked: %s", out)

	logAfter := s.readGuardLog(c)
	s.requireBlockedEvent(guardDeltaEvents(logBefore, logAfter), "OPEN")

	s.stopGuard(c)
}

// TestGuard_Bypass_ForkExecFD_File is the file-target counterpart of
// TestGuard_Bypass_ForkExecFD.
func (s *IntegrationSuite) TestGuard_Bypass_ForkExecFD_File() {
	c := s.guardContainer()
	// pooled: terminated at suite end

	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	s.exec(c, []string{"sh", "-c", "echo 'fork+exec bypass target' > /watch/target.txt"})

	exploitHostPath := absPath("./exploits/fork_exec_fd")
	err := c.CopyFileToContainer(s.ctx, exploitHostPath, "/exploits/fork_exec_fd", 0755)
	s.Require().NoError(err, "copy fork_exec_fd binary")

	s.startGuardStd(c, "/watch/target.txt", "-b", "/usr/bin/cat")

	code, out := s.exec(c, []string{"/exploits/fork_exec_fd", "/watch/target.txt", "/usr/bin/cat"})
	s.Require().NotEqualf(0, code,
		"fork_exec_fd bypass against a single-file guard should be blocked (cat is blacklisted): %s", out)

	s.stopGuard(c)
}

// TestGuard_Bypass_SymlinkBinary_File is the file-target counterpart of
// TestGuard_Bypass_SymlinkBinary.
func (s *IntegrationSuite) TestGuard_Bypass_SymlinkBinary_File() {
	c := s.guardContainer()
	// pooled: terminated at suite end

	s.exec(c, []string{"sh", "-c", "echo 'secret' > /watch/target.txt"})
	s.exec(c, []string{"ln", "-sf", "/usr/bin/cat", "/tmp/myreader"})

	s.startGuardStd(c, "/watch/target.txt", "-b", "/usr/bin/cat")

	code, out := s.exec(c, []string{"/tmp/myreader", "/watch/target.txt"})
	s.Require().NotEqualf(0, code,
		"symlink to blacklisted binary should be blocked against a single-file guard: %s", out)

	s.stopGuard(c)
}

// TestGuard_Bypass_SCMRights_File is the file-target counterpart of
// TestGuard_Bypass_SCMRights.
func (s *IntegrationSuite) TestGuard_Bypass_SCMRights_File() {
	c := s.guardContainer()
	// pooled: terminated at suite end

	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	s.exec(c, []string{"sh", "-c", "echo 'SCM_RIGHTS bypass target' > /watch/target.txt"})

	exploitHostPath := absPath("./exploits/scm_rights_pass")
	err := c.CopyFileToContainer(s.ctx, exploitHostPath, "/exploits/scm_rights_pass", 0755)
	s.Require().NoError(err, "copy scm_rights_pass binary")

	s.startGuardStd(c, "/watch/target.txt", "-b", "/usr/bin/cat")

	code, out := s.exec(c, []string{"/exploits/scm_rights_pass", "/watch/target.txt"})
	s.Require().NotEqualf(0, code,
		"SCM_RIGHTS bypass against a single-file guard should be blocked (cat is blacklisted): %s", out)

	s.stopGuard(c)
}

// TestGuard_Bypass_Mount_File is the file-target counterpart of
// TestGuard_Bypass_Mount: mount_bypass.c binds a directory over a directory,
// which does not apply to a single-file watch root, so this uses a plain
// `mount --bind` of an attacker-controlled FILE over the guarded file's own
// path — the same "rename/mount-over-watchroot" bypass class CLAUDE.md calls
// out, via mount(2) instead of rename(2). The guard's sb_mount hook checks
// the mount point dentry against guard_config's watch-root inode.
func (s *IntegrationSuite) TestGuard_Bypass_Mount_File() {
	c := s.guardContainer()
	// pooled: terminated at suite end

	s.exec(c, []string{"sh", "-c",
		"printf 'ORIGINAL-GUARDED-SECRET' > /watch/secret.txt && printf 'ATTACKER-CONTROLLED-CONTENT' > /tmp/evilfile.txt"})

	s.startGuardStd(c, "/watch/secret.txt")
	logBefore := s.readGuardLog(c)

	code, out := s.exec(c, []string{"mount", "--bind", "/tmp/evilfile.txt", "/watch/secret.txt"})
	s.Require().NotEqualf(0, code, "bind mount over a single-file guarded target should be denied: %s", out)

	logAfter := s.readGuardLog(c)
	s.requireBlockedEvent(guardDeltaEvents(logBefore, logAfter), "OPEN")

	s.stopGuard(c)

	_, content := s.exec(c, []string{"cat", "/watch/secret.txt"})
	s.Require().Equalf("ORIGINAL-GUARDED-SECRET", content,
		"guarded file content was replaced through the mount-over bypass")
}

// TestGuard_Bypass_ProcessVmReadv_File is the file-target counterpart of
// TestGuard_Bypass_ProcessVmReadv.
func (s *IntegrationSuite) TestGuard_Bypass_ProcessVmReadv_File() {
	c := s.guardContainer()
	// pooled: terminated at suite end

	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	s.exec(c, []string{"sh", "-c", "echo 'process_vm_readv target' > /watch/target.txt"})

	exploitHostPath := absPath("./exploits/process_vm_readv")
	err := c.CopyFileToContainer(s.ctx, exploitHostPath, "/exploits/process_vm_readv", 0755)
	s.Require().NoError(err, "copy process_vm_readv binary")

	// bash is whitelisted so it may open the guarded file.
	s.startGuardStd(c, "/watch/target.txt", "-w", "/bin/bash")

	code, out := s.exec(c, []string{"sh", "-c",
		"bash -c 'exec 3< /watch/target.txt; IFS= read -r -u3 line; sleep 30' & " +
			"pid=$!; sleep 1; " +
			"/exploits/process_vm_readv /watch/target.txt; code=$?; " +
			"kill $pid 2>/dev/null; exit $code"})
	s.Require().NotEqualf(0, code,
		"process_vm_readv on a whitelisted victim should be blocked against a single-file guard: %s", out)

	s.stopGuard(c)
}

// TestGuard_Bypass_PtraceRace_File is the file-target counterpart of
// TestGuard_Bypass_PtraceRace.
func (s *IntegrationSuite) TestGuard_Bypass_PtraceRace_File() {
	c := s.guardContainer()
	// pooled: terminated at suite end

	s.exec(c, []string{"mkdir", "-p", "/protected", "/exploits"})
	s.exec(c, []string{"sh", "-c", "echo 'TOP-SECRET-CONTENT' > /protected/secret && chmod 755 /protected && chmod 644 /protected/secret"})

	victimHostPath := absPath("./exploits/ptrace_race")
	s.Require().NoError(c.CopyFileToContainer(s.ctx, victimHostPath, "/exploits/race_victim", 0755), "copy race_victim")
	s.Require().NoError(c.CopyFileToContainer(s.ctx, victimHostPath, "/exploits/race_tracer", 0755), "copy race_tracer")

	// Single-file watch root: only secret's own inode is guarded, no parent
	// directory inode.
	s.startGuardStd(c, "/protected/secret", "-w", "/exploits/race_victim")

	code, out := s.exec(c, []string{"sh", "-c",
		"rm -f /tmp/race_addr* /tmp/race_go /tmp/race_done; " +
			"timeout 60 /exploits/race_tracer tracer /exploits/race_victim " +
			"/protected/secret /tmp/race_addr /tmp/race_go /tmp/race_done"})
	s.T().Logf("tracer exit=%d out=%q", code, out)

	s.Require().NotEqualf(0, code, "ptrace race against a whitelisted victim of a single-file guard must be blocked")
	s.Require().NotContains(out, "SECRET_DUMP|SUCCESS",
		"tracer exfiltrated single-file guarded content from the victim heap: %s", out)

	s.stopGuard(c)
}

// TestGuard_Bypass_InPlaceBinarySwap_File is the file-target counterpart of
// TestGuard_Bypass_InPlaceBinarySwap.
func (s *IntegrationSuite) TestGuard_Bypass_InPlaceBinarySwap_File() {
	c := s.guardContainer()
	// pooled: terminated at suite end

	s.exec(c, []string{"mkdir", "-p", "/protected", "/exploits"})
	s.exec(c, []string{"sh", "-c",
		"echo 'TOP-SECRET-CONTENT' > /protected/secret && chmod 755 /protected && chmod 644 /protected/secret"})

	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/swap_benign"), "/exploits/swap_benign", 0755), "copy swap_benign")
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/swap_reader"), "/exploits/swap_reader", 0755), "copy swap_reader")

	s.exec(c, []string{"sh", "-c", "cp /exploits/swap_benign /tmp/app && chmod 755 /tmp/app"})
	s.startGuardStd(c, "/protected/secret", "-w", "/tmp/app")

	code, out := s.exec(c, []string{"/tmp/app"})
	s.Require().Equalf(0, code, "baseline whitelisted binary should run: %s", out)
	s.Require().Contains(out, "BENIGN-APP-OK")

	_, inodeOut := s.exec(c, []string{"sh", "-c",
		"i1=$(stat -c %i /tmp/app); cat /exploits/swap_reader > /tmp/app; chmod 755 /tmp/app; i2=$(stat -c %i /tmp/app); echo $i1 $i2"})
	inodeRe := regexp.MustCompile(`(\d+)[^0-9]+(\d+)`)
	m := inodeRe.FindStringSubmatch(inodeOut)
	s.Require().NotNil(m, "inode probe output: %q", inodeOut)
	s.Require().Equalf(m[1], m[2], "replacement must be same-inode for this vector (before=%s after=%s)", m[1], m[2])

	denied := false
	var lastOut string
	deadline := time.Now().Add(40 * time.Second)
	for time.Now().Before(deadline) {
		code, out = s.exec(c, []string{"sh", "-c", "timeout 10 /tmp/app /protected/secret 2>&1"})
		lastOut = out
		if code != 0 {
			denied = true
			break
		}
		time.Sleep(2 * time.Second)
	}
	s.Require().Truef(denied, "replaced binary was never denied within the verify window, last output: %q", lastOut)

	time.Sleep(2 * time.Second)
	code, out = s.exec(c, []string{"sh", "-c", "timeout 10 /tmp/app /protected/secret 2>&1"})
	s.Require().NotEqualf(0, code, "demoted whitelist entry must stay blocked: %s", out)
	s.Require().NotContains(out, "STOLEN|")

	log := s.readGuardLog(c)
	s.Require().Contains(log, "modified in place", "guard must log the in-place tampering detection")

	s.stopGuard(c)
}

// TestGuard_Bypass_RawBlockDevice_File is the file-target counterpart of
// TestGuard_Bypass_RawBlockDevice: the watch root is a single file living on
// the loop-mounted ext4 filesystem, instead of the mount directory itself.
func (s *IntegrationSuite) TestGuard_Bypass_RawBlockDevice_File() {
	c := s.guardContainer()
	// pooled: terminated at suite end

	const marker = "TOP-SECRET-RAW-BLOCK-DEVICE-FILE-91BE"

	setup := `
set -e
dd if=/dev/zero of=/img.ext4 bs=1M count=32 status=none
mkfs.ext4 -q -F /img.ext4
for i in $(seq 0 15); do [ -e /dev/loop$i ] || mknod /dev/loop$i b 7 "$i"; done
mkdir -p /watch
LOOP=$(losetup -f --show /img.ext4)
mount "$LOOP" /watch
echo "` + marker + `" > /watch/secret.txt
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

	defer s.exec(c, []string{"sh", "-c", fmt.Sprintf(
		"umount /watch 2>/dev/null; losetup -d %s 2>/dev/null; rm -f /img.ext4; true", loopDev)})

	s.Require().NoError(
		c.CopyFileToContainer(s.ctx, absPath("./exploits/raw_block_device"), "/exploits/raw_block_device", 0755),
		"copy raw_block_device exploit")

	// Positive control: the bypass works BEFORE the guard starts.
	code, out = s.exec(c, []string{"/exploits/raw_block_device", loopDev, "/secret.txt"})
	s.Require().Equalf(0, code, "pre-guard: raw block device read should succeed (bypass): %s", out)
	s.Require().Containsf(out, marker, "pre-guard: exploit must recover the secret: %s", out)

	// Single-file watch root — BackingDevice() stats the file directly.
	s.startGuardStd(c, "/watch/secret.txt")
	startLog := s.readGuardLog(c)
	s.Require().Containsf(startLog, "blocking raw access to backing block device",
		"guard must log backing-block-device detection for a single-file watch root:\n%s", startLog)

	code, out = s.exec(c, []string{"sh", "-c", "cat /watch/secret.txt 2>&1"})
	s.Require().NotEqualf(0, code, "direct read of the guarded file must be blocked: %s", out)

	logBefore := s.readGuardLog(c)

	code, out = s.exec(c, []string{"/exploits/raw_block_device", loopDev, "/secret.txt"})
	s.Require().NotEqualf(0, code, "raw_block_device must be blocked by a single-file guard: %s", out)
	s.Require().NotContainsf(out, marker, "guard must prevent the exploit reading the secret: %s", out)

	code, out = s.exec(c, []string{"sh", "-c", fmt.Sprintf(
		"dd if=%s bs=1M count=32 2>/dev/null | tr -c '[:print:]' '\\n' | grep -q %s && echo LEAK || echo safe",
		loopDev, marker)})
	s.Require().Containsf(out, "safe", "dd raw read of the backing device must not surface the secret: %s", out)

	logAfter := s.readGuardLog(c)
	deltaEvents := guardDeltaEvents(logBefore, logAfter)
	blocked := s.requireBlockedEventSilent(deltaEvents, "OPEN",
		"raw_block_devic", "raw_block_device", "debugfs", "dd")
	if !blocked {
		s.T().Logf("raw guard log:\n%s", logAfter)
	}
	s.Require().Truef(blocked, "expected a blocked OPEN on the backing block device, got %v", deltaEvents)

	s.stopGuard(c)
}

// ---------------------------------------------------------------
// File-target counterpart of TestGuard_BypassVectors: ATTR-class ops that
// never create a struct file (truncate, chmod, chown, utimes, setxattr),
// run directly against a single-file watch root. Mknod and Rmdir are
// directory-creation operations with no equivalent on a single-file watch
// root (there is no parent guarded directory to create/remove an entry in)
// and are intentionally not duplicated here.
// ---------------------------------------------------------------

func (s *IntegrationSuite) runAttrOpFile(binary, expect string) {
	exploitHostPath := absPath(fmt.Sprintf("./exploits/%s", binary))

	c := s.guardContainer()
	// pooled: terminated at suite end

	const target = "/watch/file.txt"
	s.exec(c, []string{"sh", "-c", "echo secret > " + target})
	err := c.CopyFileToContainer(s.ctx, exploitHostPath, fmt.Sprintf("/exploits/%s", binary), 0755)
	s.Require().NoError(err, "copy exploit binary")

	// ---- Phase 1: single-file guard, block-all ----
	s.startGuardStd(c, target)
	logBefore := s.readGuardLog(c)

	code, out := s.exec(c, []string{"sh", "-c", "cat " + target + " > /dev/null 2>&1"})
	s.Require().NotEqualf(0, code, "control cat should be blocked: %s", out)

	code, out = s.exec(c, []string{fmt.Sprintf("/exploits/%s", binary), target})
	s.Require().NotEqualf(0, code, "%s must be blocked against a single-file guard in block-all mode: %s", binary, out)

	logAfter := s.readGuardLog(c)
	deltaEvents := guardDeltaEvents(logBefore, logAfter)
	if !s.requireBlockedEventSilent(deltaEvents, expect, binary) {
		s.T().Logf("phase1 raw guard log:\n%s", logAfter)
		s.requireBlockedEvent(deltaEvents, expect, binary)
	}

	// ---- Phase 2: single-file guard, exploit binary whitelisted ----
	s.stopGuard(c)
	s.startGuardStd(c, target, "-w", "/exploits/"+binary)

	code, out = s.exec(c, []string{"sh", "-c", "cat " + target + " > /dev/null 2>&1"})
	s.Require().NotEqualf(0, code, "control cat should stay blocked under whitelist: %s", out)

	code, out = s.exec(c, []string{fmt.Sprintf("/exploits/%s", binary), target})
	s.Require().Equalf(0, code, "whitelisted %s must succeed against a single-file guard: %s", binary, out)

	logAfter2 := s.readGuardLog(c)
	s.requireNoBlockedEvent(guardDeltaEvents(logAfter, logAfter2), expect)

	s.stopGuard(c)
}

func (s *IntegrationSuite) TestGuard_BypassVectors_File() {
	cases := []attrOpCase{
		{name: "Truncate", binary: "truncate", expect: "ATTR"},
		{name: "Chmod", binary: "chmod", expect: "ATTR"},
		{name: "Chown", binary: "chown", expect: "ATTR"},
		{name: "Utimes", binary: "utimes", expect: "ATTR"},
		{name: "Setxattr", binary: "setxattr", expect: "ATTR"},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			s.runAttrOpFile(tc.binary, tc.expect)
		})
	}
}

// TestGuard_Bypass_StatMetadata_File is the file-target counterpart of
// TestGuard_Bypass_StatMetadata. Unlike the directory case — where the
// watch-root DIRECTORY node itself is a deliberate stat exception — a
// single-file watch root's own inode IS the protected resource, so stat of
// the root file itself must be denied; there is no "root node" exception
// to carve out here.
func (s *IntegrationSuite) TestGuard_Bypass_StatMetadata_File() {
	c := s.guardContainer()
	// pooled: terminated at suite end

	const target = "/watch/secret"
	s.exec(c, []string{"sh", "-c", "printf 'TOP-SECRET-SIZE-LEAK' > " + target})
	exploitHostPath := absPath("./exploits/statonly")
	s.Require().NoError(c.CopyFileToContainer(s.ctx, exploitHostPath, "/exploits/statonly", 0755))

	// ---- Phase 1: block-all mode ----
	s.startGuardStd(c, target)
	logBefore := s.readGuardLog(c)

	code, out := s.exec(c, []string{"/exploits/statonly", target})
	s.Require().NotEqualf(0, code, "stat of a single-file guarded target must be denied: %s", out)
	s.Require().NotContainsf(out, "SIZE=", "stat leaked metadata of a single-file guarded target: %s", out)

	code, out = s.exec(c, []string{"sh", "-c", "ls -la " + target + " 2>&1"})
	s.Require().NotEqualf(0, code, "ls -la of a single-file guarded target must be denied: %s", out)

	logAfter := s.readGuardLog(c)
	s.requireBlockedEvent(guardDeltaEvents(logBefore, logAfter), "STAT", "statonly")

	// ---- Phase 2: whitelist mode (statonly whitelisted) ----
	s.stopGuard(c)
	s.startGuardStd(c, target, "-w", "/exploits/statonly")

	code, out = s.exec(c, []string{"/exploits/statonly", target})
	s.Require().Equalf(0, code, "whitelisted stat of a single-file guarded target must succeed: %s", out)
	s.Require().Contains(out, "SIZE=20")

	s.stopGuard(c)

	// ---- Phase 3: a READ-restricted whitelist entry still gets STAT ----
	s.startGuardStd(c, target, "-w", "/exploits/statonly", "-e", "READ")
	code, out = s.exec(c, []string{"/exploits/statonly", target})
	s.Require().Equalf(0, code, "a READ-masked whitelist entry must still be allowed to stat a single-file target: %s", out)
	s.Require().Contains(out, "SIZE=20")
	s.stopGuard(c)
}

// ---------------------------------------------------------------
// Hardlink / symlink / rename-out coverage for a single-file watch root.
//
// TestGuard_BlocksAll_File already renames/hardlinks the watch root to a
// sibling path with plain mv/ln, but those coreutils lstat() their source
// first (see rawfsop.c's header comment) and are denied at the STAT hook
// before path_rename/path_link is ever reached — the same masking
// TestGuard_Bypass_RenameOverGuardedFile calls out for the destination
// side. The tests below use rawfsop's bare syscalls to drive straight into
// path_rename's and path_link's SOURCE-side checks (guard.bpf.c's
// read_inode_guard/guarded_map_hit on old_dentry), proving the hooks
// themselves — not just the STAT hook — deny moving or aliasing the
// watch root's inode out to an unguarded name. A dedicated symlink test
// covers the third bypass class this hook set defends: path_symlink's
// target-content match, which denies creating a symlink ANYWHERE (even
// outside the guarded path) whose link content names the watch root.
// ---------------------------------------------------------------

// TestGuard_Bypass_RenameOutOfWatchRoot_File proves path_rename's
// source-side check (guard.bpf.c ~line 1292: read_inode_guard(inode) &&
// root_in_chain) denies renaming the single-file watch root itself to an
// unguarded destination — the inverse direction of
// TestGuard_Bypass_RenameOverGuardedFile's destination-victim check.
func (s *IntegrationSuite) TestGuard_Bypass_RenameOutOfWatchRoot_File() {
	c := s.guardContainer()
	// pooled: terminated at suite end

	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	s.exec(c, []string{"sh", "-c", "rm -f /tmp/escaped.txt; printf 'ORIGINAL-GUARDED-SECRET' > /watch/secret.txt"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/rawfsop"), "/exploits/rawfsop", 0755))

	s.startGuardStd(c, "/watch/secret.txt")
	logBefore := s.readGuardLog(c)

	// Bare rename(2): coreutils mv would be stopped at STAT before ever
	// reaching path_rename, masking whether the source-side check works.
	code, out := s.exec(c, []string{"/exploits/rawfsop", "rename", "/watch/secret.txt", "/tmp/escaped.txt"})
	s.Require().NotEqualf(0, code, "renaming the single-file watch root itself out must be denied: %s", out)

	logAfter := s.readGuardLog(c)
	s.requireBlockedEvent(guardDeltaEvents(logBefore, logAfter), "RENAME")

	s.stopGuard(c)

	code, _ = s.exec(c, []string{"sh", "-c", "test -e /tmp/escaped.txt"})
	s.Require().NotEqualf(0, code, "the watch root must not have moved to the unguarded destination")
	_, content := s.exec(c, []string{"cat", "/watch/secret.txt"})
	s.Require().Equal("ORIGINAL-GUARDED-SECRET", content, "the watch root's content must be intact after the denied move")
}

// TestGuard_Bypass_HardlinkOutOfWatchRoot_File proves path_link's
// source-side check (guard.bpf.c ~line 1430: guarded_map_hit(old_dentry, ...))
// denies creating a hardlink to the single-file watch root's inode at an
// unguarded name — aliasing the guarded inode under a name the guard never
// sees would otherwise be a path-based bypass of an inode-keyed guard.
func (s *IntegrationSuite) TestGuard_Bypass_HardlinkOutOfWatchRoot_File() {
	c := s.guardContainer()
	// pooled: terminated at suite end

	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	s.exec(c, []string{"sh", "-c", "rm -f /tmp/hardlink_escape.txt; printf 'ORIGINAL-GUARDED-SECRET' > /watch/secret.txt"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/rawfsop"), "/exploits/rawfsop", 0755))

	s.startGuardStd(c, "/watch/secret.txt")
	logBefore := s.readGuardLog(c)

	// Bare link(2): coreutils ln would be stopped at STAT first, same
	// masking concern as the rename case above.
	code, out := s.exec(c, []string{"/exploits/rawfsop", "link", "/watch/secret.txt", "/tmp/hardlink_escape.txt"})
	s.Require().NotEqualf(0, code,
		"hardlinking the single-file watch root out to an unguarded name must be denied: %s", out)

	logAfter := s.readGuardLog(c)
	s.requireBlockedEvent(guardDeltaEvents(logBefore, logAfter), "HARDLINK")

	s.stopGuard(c)

	code, _ = s.exec(c, []string{"sh", "-c", "test -e /tmp/hardlink_escape.txt"})
	s.Require().NotEqualf(0, code, "no alias of the watch root's inode may exist at the unguarded destination")
}

// TestGuard_Bypass_SymlinkTargetMatch_File proves path_symlink's second
// check (guard.bpf.c: the prefix match against the stored guard_path) denies
// creating a symlink ANYWHERE — even entirely outside the guarded area —
// whose link content names the single-file watch root. Unlike the
// destination-side symlink-creation check (creating a link INSIDE a guarded
// directory), this one has nothing to do with where the new dentry lives;
// it is keyed on what the symlink POINTS TO.
func (s *IntegrationSuite) TestGuard_Bypass_SymlinkTargetMatch_File() {
	c := s.guardContainer()
	// pooled: terminated at suite end

	s.exec(c, []string{"sh", "-c", "rm -f /tmp/alias /tmp/unrelated; echo secret > /watch/secret.txt"})

	s.startGuardStd(c, "/watch/secret.txt")
	logBefore := s.readGuardLog(c)

	code, out := s.exec(c, []string{"ln", "-s", "/watch/secret.txt", "/tmp/alias"})
	s.Require().NotEqualf(0, code,
		"a symlink outside the guarded area whose target names the single-file watch root must be denied: %s", out)

	logAfter := s.readGuardLog(c)
	s.requireBlockedEvent(guardDeltaEvents(logBefore, logAfter), "SYMLINK")

	code, _ = s.exec(c, []string{"sh", "-c", "test -e /tmp/alias -o -L /tmp/alias"})
	s.Require().NotEqualf(0, code, "the decoy symlink must not have been created")

	// Negative control: an unrelated symlink target is unaffected — the
	// check is content-specific, not a blanket symlink ban outside the tree.
	logCursor := s.readGuardLog(c)
	code, out = s.exec(c, []string{"ln", "-s", "/etc/passwd", "/tmp/unrelated"})
	s.Require().Equalf(0, code, "a symlink unrelated to the guarded path must succeed: %s", out)
	s.requireNoBlockedEvent(guardDeltaEvents(logCursor, s.readGuardLog(c)), "SYMLINK")

	s.stopGuard(c)
}

// TestGuard_Bypass_SymlinkAliasAccess_File proves that a symlink pointing at
// the single-file watch root, created BEFORE the guard ever attached (so
// path_symlink's target-content check never saw it), still cannot be used to
// read the guarded content: file_open resolves the symlink and checks the
// FINAL inode, so enforcement is inode-keyed and independent of the path
// used to reach it.
func (s *IntegrationSuite) TestGuard_Bypass_SymlinkAliasAccess_File() {
	c := s.guardContainer()
	// pooled: terminated at suite end

	s.exec(c, []string{"sh", "-c", "rm -f /tmp/alias; echo secret > /watch/secret.txt"})
	s.exec(c, []string{"ln", "-s", "/watch/secret.txt", "/tmp/alias"})

	s.startGuardStd(c, "/watch/secret.txt")
	logBefore := s.readGuardLog(c)

	code, out := s.exec(c, []string{"sh", "-c", "cat /tmp/alias > /dev/null 2>&1"})
	s.Require().NotEqualf(0, code, "reading the watch root through a pre-existing symlink alias must be denied: %s", out)

	logAfter := s.readGuardLog(c)
	s.requireCatBlocked(guardDeltaEvents(logBefore, logAfter))

	s.stopGuard(c)
}
