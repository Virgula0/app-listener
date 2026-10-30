package guard

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	cilium "github.com/cilium/ebpf"
	"golang.org/x/sys/unix"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// BusClient is a session-bus client resolved from the pidfd its bus took at connect, for the secret
// gatekeeper's decisions (internal/secretgate). It replays check_and_emit_args on the guard's own
// maps: the caller isn't the current task, so no LSM hook can judge it.
type BusClient struct {
	PID     uint32
	UID     uint32
	Comm    string // display only
	key     GuardInodeKey
	refused bool  // exe_refused: a superseded image exec'd at or after the supersede
	suspect uint8 // trust_code_suspect mark, 0 = none
}

// ErrUnattested refuses a bus client whose connection can't be tied to its current image.
var ErrUnattested = errors.New("bus client not attested")

// RegisterBusDir lets the guard record connects to dir/bus (a /run/user/<uid> dir) for AttestBusClient.
// A user's connects made before it are unattested.
func RegisterBusDir(dir string) (GuardInodeKey, error) {
	dev, ino, err := ebpf.StatInode(dir)
	if err != nil {
		return GuardInodeKey{}, err
	}
	key := GuardInodeKey{Dev: dev, Ino: ino}
	sharedEngine.mu.Lock()
	defer sharedEngine.mu.Unlock()
	if sharedEngine.objs.GuardBusDirs == nil {
		return GuardInodeKey{}, errors.New("guard engine not running")
	}
	return key, sharedEngine.objs.GuardBusDirs.Put(key, uint8(1))
}

// UnregisterBusDir drops a dir RegisterBusDir added (its tmpfs is gone at logout).
func UnregisterBusDir(key GuardInodeKey) {
	sharedEngine.mu.Lock()
	defer sharedEngine.mu.Unlock()
	if sharedEngine.objs.GuardBusDirs != nil {
		_ = sharedEngine.objs.GuardBusDirs.Delete(key)
	}
}

// AttestBusClient resolves pidfd (GetConnectionCredentials' ProcessFD). ErrUnattested when the process
// exec'd since its first connect to a registered bus (the connection may be another image's), made it
// before the guard recorded it, or is gone. The exe is read between two equal exec stamps, so it is
// the image the bus row vouches for.
func AttestBusClient(pidfd int) (*BusClient, error) {
	pid, err := pidfdPid(pidfd)
	if err != nil {
		return nil, err
	}
	sharedEngine.mu.Lock()
	objs := sharedEngine.objs
	sharedEngine.mu.Unlock()
	if objs.GuardBusConns == nil {
		return nil, errors.New("guard engine not running")
	}
	var bus GuardBusConn
	if objs.GuardBusConns.Lookup(pid, &bus) != nil {
		return nil, fmt.Errorf("%w: pid %d has no recorded bus connect", ErrUnattested, pid)
	}
	stamp, err := execStamp(&objs, pid)
	if err != nil {
		return nil, err
	}
	if stamp.Start != bus.Start || stamp.Seq != bus.Seq || bus.Seq == 0 {
		return nil, fmt.Errorf("%w: pid %d exec'd after connecting to the bus", ErrUnattested, pid)
	}
	c, err := imageOf(&objs, pid, stamp)
	if err != nil {
		return nil, err
	}
	c.UID, c.Comm = procIdentity(pid)
	again, err := execStamp(&objs, pid)
	if err != nil || again != stamp {
		return nil, fmt.Errorf("%w: pid %d exec'd while being attested", ErrUnattested, pid)
	}
	// The pidfd pins the process, not its pid number: alive now = /proc/<pid> was this process.
	if err := unix.PidfdSendSignal(pidfd, 0, nil, 0); err != nil {
		return nil, fmt.Errorf("%w: pid %d gone: %v", ErrUnattested, pid, err)
	}
	return c, nil
}

// busClientKey is pid's exe key as the guard keys it: a multicall's carries the applet stamped at exec.
func busClientKey(objs *GuardObjects, pid uint32, stamp GuardExecStamp) (GuardInodeKey, error) {
	var st unix.Stat_t
	if err := unix.Stat(fmt.Sprintf("/proc/%d/exe", pid), &st); err != nil {
		return GuardInodeKey{}, fmt.Errorf("%w: pid %d exe: %v", ErrUnattested, pid, err)
	}
	key := GuardInodeKey{Dev: ebpf.KernelDev(unix.Major(st.Dev), unix.Minor(st.Dev)), Ino: st.Ino}
	err := objs.GuardMulticall.Lookup(key, new(uint8))
	if errors.Is(err, cilium.ErrKeyNotExist) {
		return key, nil
	}
	if err != nil {
		return GuardInodeKey{}, err
	}
	tag := mcTagNone
	if stamp.Tag != 0 {
		tag = stamp.Tag
	}
	return key.withTag(tag), nil
}

// imageOf is pid's image as check_and_emit_args judges it: its key, whether it is a superseded key
// exec'd at or after the supersede, its code-suspect mark.
func imageOf(objs *GuardObjects, pid uint32, stamp GuardExecStamp) (*BusClient, error) {
	key, err := busClientKey(objs, pid, stamp)
	if err != nil {
		return nil, err
	}
	c := &BusClient{PID: pid, key: key}
	var sup GuardExeSupersede
	if err := objs.ExeSuperseded.Lookup(key.Real(), &sup); err == nil {
		c.refused = stamp.Seq >= sup.Seq
	} else if !errors.Is(err, cilium.ErrKeyNotExist) {
		return nil, err
	}
	if err := objs.TrustCodeSuspect.Lookup(pid, &c.suspect); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
		return nil, err
	}
	return c, nil
}

func execStamp(objs *GuardObjects, pid uint32) (GuardExecStamp, error) {
	var stamp GuardExecStamp
	if err := objs.ExeStamps.Lookup(pid, &stamp); err != nil {
		return stamp, fmt.Errorf("%w: pid %d has no exec stamp", ErrUnattested, pid)
	}
	return stamp, nil
}

// AdmitsBusClient reports whether this guard's whitelist admits c for op, as check_and_emit_args
// would for a file below it: an allow row, not code-suspect, op in the row's mask. GUARD_ALLOW_ROOT
// (the daemon's own key) never admits a bus client. It logs the decision like a kernel event.
func (g *Guard) AdmitsBusClient(c *BusClient, op ebpf.EventType, what string) bool {
	ok, suspect := g.busDecision(c, op)
	ev := &GuardEvent{
		FileEvent: ebpf.FileEvent{PID: c.PID, UID: c.UID, Type: op, Comm: c.Comm, Path: what,
			ExeDev: c.key.Dev, ExeIno: c.key.Ino},
		Blocked: !ok,
		Exe:     c.key,
	}
	if suspect {
		ev.Suspect = suspectLabel(uint32(c.suspect) << guardReasonSuspectShift)
	}
	if !ok || !g.quietAllow() {
		g.dispatch(ev)
	}
	return ok
}

func (g *Guard) busDecision(c *BusClient, op ebpf.EventType) (ok, suspect bool) {
	if c.refused || g.mode != ModeWhitelist {
		return false, false
	}
	k := g.resKey(c.key)
	var action uint8
	if err := g.objs().GuardExeActions.Lookup(k, &action); err != nil || action != GUARD_ALLOW {
		return false, false
	}
	if c.suspect != 0 {
		return false, true
	}
	var mask uint32
	if g.objs().GuardExeEvents.Lookup(k, &mask) == nil && op < 32 && mask&(1<<uint(op)) == 0 { //nolint:gosec // op < 32
		return false, false
	}
	return true, false
}

func (g *Guard) quietAllow() bool {
	var quiet uint8
	return g.objs().GuardQuietAllow.Lookup(g.resID, &quiet) == nil && quiet != 0
}

// pidfdPid reads the pid a pidfd refers to; -1 there means the process exited.
func pidfdPid(pidfd int) (uint32, error) {
	f, err := os.Open(fmt.Sprintf("/proc/self/fdinfo/%d", pidfd))
	if err != nil {
		return 0, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		v, ok := strings.CutPrefix(s.Text(), "Pid:")
		if !ok {
			continue
		}
		pid, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil || pid <= 0 {
			return 0, fmt.Errorf("%w: pidfd %d: process gone", ErrUnattested, pidfd)
		}
		return uint32(pid), nil //nolint:gosec // pid > 0
	}
	return 0, fmt.Errorf("pidfd %d: no Pid in fdinfo (not a pidfd?)", pidfd)
}

// procIdentity is display data: the real uid and comm of pid, read racily.
func procIdentity(pid uint32) (uid uint32, comm string) {
	raw, _ := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
	if f, err := os.Open(fmt.Sprintf("/proc/%d/status", pid)); err == nil {
		s := bufio.NewScanner(f)
		for s.Scan() {
			if v, ok := strings.CutPrefix(s.Text(), "Uid:"); ok {
				if fs := strings.Fields(v); len(fs) > 0 {
					if n, err := strconv.ParseUint(fs[0], 10, 32); err == nil {
						uid = uint32(n)
					}
				}
				break
			}
		}
		f.Close()
	}
	return uid, strings.TrimSpace(string(raw))
}

// RunsExe reports whether c runs the file keyed file (a role binary from [secret_service]) under a
// clean image: not a refused superseded one, not code-suspect.
func (c *BusClient) RunsExe(file GuardInodeKey) bool {
	return c.key == file && !c.refused && c.suspect == 0
}

// LiveGuardAt is a running guard rooted at path (a reload briefly runs two; either holds a whitelist
// of a config that passed), nil if none.
func LiveGuardAt(path string) *Guard {
	sharedEngine.mu.Lock()
	var cands []*Guard
	for _, g := range &sharedEngine.slots {
		if g != nil && g.path == path {
			cands = append(cands, g)
		}
	}
	sharedEngine.mu.Unlock()
	// g.mu outside e.mu: publishPort takes them in the other order.
	for _, g := range cands {
		g.mu.Lock()
		stopped := g.stopped
		g.mu.Unlock()
		if !stopped {
			return g
		}
	}
	return nil
}
