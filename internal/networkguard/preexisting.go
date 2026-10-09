package networkguard

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cilium/ebpf/link"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/logging"
)

// Processes older than the code-integrity hooks: the kernel saw neither their exec nor their
// mappings. markPreexisting judges each protected one once, by the kernel's rules, from /proc: traced,
// started with LD_PRELOAD/LD_AUDIT, or mapping executable code its exe doesn't vouch for
// (code_vouched) marks it code-suspect, keyed by the tgid netg_procs pairs with its /proc pid. Code written into its memory before the guard started is
// beyond any scan: whitelisted apps restarted after the guard get the full guarantee.

const mcTagMask = uint64(0xffff) << 32

// fileFacts is what code_vouched reads of an inode: owner, mode, superblock (mountinfo dev) and
// whether root vouches for that superblock's ownership bits.
type fileFacts struct {
	uid, gid, mode uint32
	sb             uint64
	vouched        bool
}

const noGID = ^uint32(0)

// roFor mirrors inode_ro_for: owned by uid, not other-writable, group-writable only for gwGid.
func roFor(f fileFacts, uid, gwGid uint32) bool {
	if f.uid != uid || f.mode&unix.S_IWOTH != 0 {
		return false
	}
	return f.mode&unix.S_IWGRP == 0 || f.gid == gwGid
}

// codeVouched mirrors code_vouched (networkguard.bpf.c).
func codeVouched(lib, dir, exe fileFacts) bool {
	if lib.vouched && roFor(lib, 0, 0) && roFor(dir, 0, 0) {
		return true
	}
	if exe.uid == 0 || (!lib.vouched && lib.sb != exe.sb) {
		return false
	}
	gw := noGID
	if exe.mode&unix.S_IWGRP != 0 {
		gw = exe.gid
	}
	return roFor(lib, exe.uid, gw) && (roFor(dir, exe.uid, gw) || roFor(dir, 0, 0))
}

// protectedFiles is the real key of every exe the kernel protects: an ALLOW row's file (an applet
// row's tag masked off), or a multicall with a whitelisted applet.
func (g *NetGuard) protectedFiles() (map[GuardNetInodeKey]bool, error) {
	out := make(map[GuardNetInodeKey]bool)
	var k GuardNetInodeKey
	var action uint8
	it := g.objs.GuardNetExeActions.Iterate()
	for it.Next(&k, &action) {
		if action == bpfAllow {
			out[GuardNetInodeKey{Dev: k.Dev &^ mcTagMask, Ino: k.Ino}] = true
		}
	}
	if err := it.Err(); err != nil {
		return nil, fmt.Errorf("reading whitelist rows: %w", err)
	}
	var mc uint32
	it = g.objs.McMulticall.Iterate()
	for it.Next(&k, &mc) {
		if mc&mcHasAllow != 0 {
			out[k] = true
		}
	}
	if err := it.Err(); err != nil {
		return nil, fmt.Errorf("reading multicall rows: %w", err)
	}
	return out, nil
}

// procRec is one netg_procs record: a thread group by its pid in this namespace and the kernel's tgid.
type procRec struct {
	nspid, tgid uint32
	start       uint64
}

const procRecSize = 16

func parseProcs(b []byte) ([]procRec, error) {
	if len(b)%procRecSize != 0 {
		return nil, fmt.Errorf("task iterator returned %d bytes, not whole %d-byte records", len(b), procRecSize)
	}
	out := make([]procRec, 0, len(b)/procRecSize)
	for ; len(b) > 0; b = b[procRecSize:] {
		out = append(out, procRec{nspid: binary.LittleEndian.Uint32(b), tgid: binary.LittleEndian.Uint32(b[4:]),
			start: binary.LittleEndian.Uint64(b[8:])})
	}
	return out, nil
}

// visibleProcs runs netg_procs once.
func (g *NetGuard) visibleProcs() ([]procRec, error) {
	it, err := link.AttachIter(link.IterOptions{Program: g.objs.NetgProcs})
	if err != nil {
		return nil, fmt.Errorf("attaching the task iterator: %w", err)
	}
	defer it.Close()
	rd, err := it.Open()
	if err != nil {
		return nil, fmt.Errorf("opening the task iterator: %w", err)
	}
	defer rd.Close()
	b, err := io.ReadAll(rd)
	if err != nil {
		return nil, fmt.Errorf("reading the task iterator: %w", err)
	}
	return parseProcs(b)
}

// resolveOwnPID sets ownPID to this process's tgid, which events carry, from its namespace pid.
func (g *NetGuard) resolveOwnPID() error {
	procs, err := g.visibleProcs()
	if err != nil {
		return err
	}
	for _, p := range procs {
		if int(p.nspid) == os.Getpid() {
			g.ownPID = int(p.tgid)
			return nil
		}
	}
	return errors.New("own process not in the task iterator's output")
}

// maxJudgePasses bounds judgePasses; not converging refuses to start.
const maxJudgePasses = 16

// judgePasses calls mark once per process (tgid, start) of each pass of list, until a pass marks
// nothing new: a process forked after a pass listed its parent but before the parent was marked
// inherited no mark, and only the next pass sees it.
func judgePasses(list func() ([]procRec, error), mark func(procRec) (bool, error)) error {
	type id struct {
		tgid  uint32
		start uint64
	}
	judged := make(map[id]bool)
	for range maxJudgePasses {
		procs, err := list()
		if err != nil {
			return err
		}
		marked := false
		for _, p := range procs {
			if judged[id{p.tgid, p.start}] {
				continue
			}
			judged[id{p.tgid, p.start}] = true
			m, err := mark(p)
			if err != nil {
				return err
			}
			marked = marked || m
		}
		if !marked {
			return nil
		}
	}
	return fmt.Errorf("suspect processes kept forking while being judged (%d passes)", maxJudgePasses)
}

func (g *NetGuard) markPreexisting() error {
	protected, err := g.protectedFiles()
	if err != nil || len(protected) == 0 {
		return err
	}
	return judgePasses(g.visibleProcs, func(p procRec) (bool, error) { return g.markProcess(p, protected) })
}

// markProcess judges p by its /proc pid and marks its tgid. A pid reused meanwhile marks a dead
// tgid; the new process is in the next pass.
func (g *NetGuard) markProcess(p procRec, protected map[GuardNetInodeKey]bool) (bool, error) {
	pid := int(p.nspid)
	if pid <= 0 || pid == os.Getpid() {
		return false, nil
	}
	why, jerr := g.judgeProcess(pid, protected)
	if errors.Is(jerr, os.ErrNotExist) || errors.Is(jerr, unix.ESRCH) {
		return false, nil // exited meanwhile
	}
	if jerr != nil {
		why = reasonCode // unreadable: fail closed
	}
	if why == 0 {
		return false, nil
	}
	if err := g.objs.TrustCodeSuspect.Put(p.tgid, uint8(why)); err != nil { //nolint:gosec // a reason code
		return false, fmt.Errorf("marking pid %d (tgid %d) code-suspect: %w", pid, p.tgid, err)
	}
	comm, _ := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	log.Warnf("network guard: %s (pid %d) started before the guard and is %s (%v): refused the "+
		"network until it restarts", logging.SanitizeText(strings.TrimSpace(string(comm))), pid,
		reasonText(why), jerr)
	return true, nil
}

// judgeProcess returns pid's code-suspect reason, 0 if it is clean or not protected.
func (g *NetGuard) judgeProcess(pid int, protected map[GuardNetInodeKey]bool) (uint32, error) {
	proc := fmt.Sprintf("/proc/%d", pid)
	var st unix.Stat_t
	if err := unix.Stat(proc+"/exe", &st); err != nil {
		if errors.Is(err, unix.ENOENT) {
			return 0, nil // a kernel thread, or gone
		}
		return 0, err
	}
	exeKey := GuardNetInodeKey{Dev: ebpf.KernelDev(unix.Major(st.Dev), unix.Minor(st.Dev)), Ino: st.Ino}
	if !protected[exeKey] {
		return 0, nil
	}
	if traced, err := isTraced(proc); err != nil || traced {
		return reasonCode, err
	}
	env, err := os.ReadFile(proc + "/environ")
	if err != nil {
		return 0, err
	}
	if loaderEnv(env) {
		return reasonEnv, nil
	}
	mounts, err := readMountDevs(proc + "/mountinfo")
	if err != nil {
		return 0, err
	}
	exe, err := g.facts(proc+"/exe", mounts)
	if err != nil {
		return 0, err
	}
	return g.judgeMaps(proc, exeKey, exe, mounts)
}

func isTraced(proc string) (bool, error) {
	b, err := os.ReadFile(proc + "/status")
	if err != nil {
		return false, err
	}
	for _, line := range strings.Split(string(b), "\n") {
		if v, ok := strings.CutPrefix(line, "TracerPid:"); ok {
			return strings.TrimSpace(v) != "0", nil
		}
	}
	return false, errors.New("no TracerPid in status")
}

// loaderEnv mirrors env_loader: LD_PRELOAD or LD_AUDIT set to a non-empty value.
func loaderEnv(environ []byte) bool {
	for _, kv := range bytes.Split(environ, []byte{0}) {
		for _, name := range []string{"LD_PRELOAD=", "LD_AUDIT="} {
			if v, ok := bytes.CutPrefix(kv, []byte(name)); ok && len(v) > 0 {
				return true
			}
		}
	}
	return false
}

type mountDev struct {
	dev  uint64
	fuse bool
}

// readMountDevs maps each mount id of a mountinfo to its superblock dev and FUSE-ness.
func readMountDevs(path string) (map[uint64]mountDev, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	out := make(map[uint64]mountDev)
	for _, line := range strings.Split(strings.TrimSpace(string(b)), "\n") {
		f := strings.Fields(line)
		sep := -1
		for i, x := range f {
			if x == "-" {
				sep = i
				break
			}
		}
		if sep < 6 || sep+1 >= len(f) {
			return nil, fmt.Errorf("malformed mountinfo line %q", line)
		}
		id, err := strconv.ParseUint(f[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("mountinfo line %q: %w", line, err)
		}
		maj, mnr, ok := strings.Cut(f[2], ":")
		major, err1 := strconv.ParseUint(maj, 10, 12)
		minor, err2 := strconv.ParseUint(mnr, 10, 20)
		if !ok || err1 != nil || err2 != nil {
			return nil, fmt.Errorf("mountinfo line %q: bad device", line)
		}
		fstype := f[sep+1]
		out[id] = mountDev{dev: ebpf.KernelDev(uint32(major), uint32(minor)),
			fuse: fstype == "fuse" || fstype == "fuseblk" || strings.HasPrefix(fstype, "fuse.")}
	}
	return out, nil
}

// facts stats path (following it) for code_vouched, its superblock found by mount id in mounts. A
// mount the namespace doesn't list (an internal one: memfd, anonymous files) vouches for nothing.
func (g *NetGuard) facts(path string, mounts map[uint64]mountDev) (fileFacts, error) {
	var stx unix.Statx_t
	if err := unix.Statx(unix.AT_FDCWD, path, 0, unix.STATX_BASIC_STATS|unix.STATX_MNT_ID, &stx); err != nil {
		return fileFacts{}, err
	}
	f := fileFacts{uid: stx.Uid, gid: stx.Gid, mode: uint32(stx.Mode)}
	if stx.Mask&unix.STATX_MNT_ID == 0 {
		return f, nil
	}
	if m, ok := mounts[stx.Mnt_id]; ok {
		f.sb = m.dev
		f.vouched = !m.fuse && g.mounts != nil && g.mounts.Vouched(m.dev)
	} else {
		f.sb = ^uint64(0) - stx.Mnt_id // unlisted: no other mount's dev
	}
	return f, nil
}

// judgeMaps judges pid's executable file mappings (/proc/<pid>/maps), each through map_files so the
// inode is the mapped one, whatever its path names now.
func (g *NetGuard) judgeMaps(proc string, exeKey GuardNetInodeKey, exe fileFacts,
	mounts map[uint64]mountDev) (uint32, error) {
	b, err := os.ReadFile(proc + "/maps")
	if err != nil {
		return 0, err
	}
	sc := bufio.NewScanner(bytes.NewReader(b))
	sc.Buffer(make([]byte, 0, 64<<10), 1<<20)
	for sc.Scan() {
		f := strings.Fields(sc.Text())
		if len(f) < 6 || len(f[1]) < 3 || f[1][2] != 'x' || f[4] == "0" {
			continue // not executable, or anonymous
		}
		var st unix.Stat_t
		mapped := proc + "/map_files/" + f[0]
		if err := unix.Stat(mapped, &st); err != nil {
			return 0, err
		}
		if (GuardNetInodeKey{Dev: ebpf.KernelDev(unix.Major(st.Dev), unix.Minor(st.Dev)), Ino: st.Ino}) == exeKey {
			continue
		}
		lib, err := g.facts(mapped, mounts)
		if err != nil {
			return 0, err
		}
		name := strings.TrimSuffix(strings.Join(f[5:], " "), " (deleted)")
		dir, err := g.facts(proc+"/root"+filepath.Dir(name), mounts)
		if err != nil || !codeVouched(lib, dir, exe) {
			return reasonCode, err
		}
	}
	return 0, sc.Err()
}
