package guard

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	cilium "github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// hostMountinfo is pid 1's mount table: the daemon's own namespace (PrivateTmp, ProtectSystem) is
// remounted nosuid by systemd and says nothing about who mounted what.
const hostMountinfo = "/proc/1/mountinfo"

// guard_dev_seq slots — must match guard_trust.bpf.c.
const (
	devSeqNow  uint32 = 0
	devSeqLost uint32 = 1
)

const (
	mountResyncInterval = time.Minute
	mountRetryInterval  = time.Second
)

// vouchedDevs returns the kernel dev (KernelDev) of every superblock in mountinfo with at least one
// mount that is neither nosuid nor FUSE. Any unparsable line fails the whole table.
func vouchedDevs(mountinfo []byte) (map[uint64]struct{}, error) {
	out := make(map[uint64]struct{})
	for _, line := range strings.Split(strings.TrimSpace(string(mountinfo)), "\n") {
		f := strings.Fields(line)
		sep := slices.Index(f, "-")
		if sep < 6 || sep+1 >= len(f) {
			return nil, fmt.Errorf("malformed mountinfo line %q", line)
		}
		dev, err := parseMajorMinor(f[2])
		if err != nil {
			return nil, fmt.Errorf("mountinfo line %q: %w", line, err)
		}
		if slices.Contains(strings.Split(f[5], ","), "nosuid") || isFuseType(f[sep+1]) {
			continue
		}
		out[dev] = struct{}{}
	}
	if len(out) == 0 {
		return nil, errors.New("no superblock mounted without nosuid (is / really nosuid?)")
	}
	return out, nil
}

func parseMajorMinor(s string) (uint64, error) {
	maj, mnr, ok := strings.Cut(s, ":")
	if !ok {
		return 0, fmt.Errorf("bad device %q", s)
	}
	major, err := strconv.ParseUint(maj, 10, 12)
	if err != nil {
		return 0, fmt.Errorf("bad major in %q: %w", s, err)
	}
	minor, err := strconv.ParseUint(mnr, 10, 20)
	if err != nil {
		return 0, fmt.Errorf("bad minor in %q: %w", s, err)
	}
	return ebpf.KernelDev(uint32(major), uint32(minor)), nil
}

func isFuseType(fstype string) bool {
	return fstype == "fuse" || fstype == "fuseblk" || strings.HasPrefix(fstype, "fuse.")
}

// MountVouch keeps an object's vouched-dev maps (guard_vouched_devs, guard_dead_devs,
// guard_dev_seq; the network guard's netg_* copies) in step with pid 1's mount table: only a
// superblock root mounted without nosuid vouches for root ownership bits. The object's sb_delete
// program must be attached before Open, so no superblock dies unseen after the first read.
type MountVouch struct {
	name                   string // log prefix
	vouchMap, deadMap, seq *cilium.Map
	restart                string // what the operator restarts once mount changes are untracked

	mu        sync.Mutex
	fd, stop  int
	done      chan struct{}
	count     int
	devs      map[uint64]struct{}
	untracked bool // the watch died, so Sync refuses to vouch anything again
}

// NewMountVouch tracks the given maps; name prefixes its log lines, restart names the process to
// restart once mount changes can no longer be tracked.
func NewMountVouch(name, restart string, vouched, dead, seq *cilium.Map) *MountVouch {
	return &MountVouch{name: name, restart: restart, vouchMap: vouched, deadMap: dead, seq: seq,
		fd: -1, stop: -1, count: -1}
}

// Open fills the vouched set from pid 1's mount table. The watch fd is opened before the first read
// so no mount change after it goes unseen.
func (v *MountVouch) Open() error {
	fd, err := unix.Open(hostMountinfo, unix.O_RDONLY|unix.O_CLOEXEC, 0)
	if err != nil {
		return fmt.Errorf("opening %s: %w", hostMountinfo, err)
	}
	v.fd = fd
	return v.Sync()
}

// Vouched reports whether the last successful sync vouched superblock dev.
func (v *MountVouch) Vouched(dev uint64) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	_, ok := v.devs[dev]
	return ok
}

// Sync re-syncs the vouched set to pid 1's current mount table. On failure it drops every vouched
// dev (system libraries become untrusted) rather than keep a possibly stale one.
func (v *MountVouch) Sync() error {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.untracked {
		return fmt.Errorf("mount changes are no longer tracked; restart the %s", v.restart)
	}
	// A btrfs mounted since start may have just loaded the module and its BTF.
	ensureBtrfsLayout()
	devs, err := v.syncVouchedDevs()
	if err != nil {
		v.devs = nil
		if cerr := syncMap(v.vouchMap, map[uint64]uint64{}); cerr != nil {
			err = fmt.Errorf("%w; clearing the vouched set also failed: %w", err, cerr)
		}
		return err
	}
	v.devs = devs
	// Container and snap churn remounts constantly: only the first count is worth an info line.
	if v.count < 0 {
		log.Infof("%s: %d superblock(s) vouched for system-library auto-trust "+
			"(mounted without nosuid in pid 1's namespace)", v.name, len(devs))
	} else if len(devs) != v.count {
		log.Debugf("%s: %d superblock(s) vouched after a mount change", v.name, len(devs))
	}
	v.count = len(devs)
	return nil
}

// syncVouchedDevs: seq snapshot → mountinfo → delete stale → put → prune tombstones. Stale devs go
// first because a dev can be reused by another superblock; entries carry the snapshot so the
// kernel refuses one whose superblock died after it (the dead-devs map).
func (v *MountVouch) syncVouchedDevs() (map[uint64]struct{}, error) {
	var seq, lost uint64
	if err := v.seq.Lookup(devSeqNow, &seq); err != nil {
		return nil, fmt.Errorf("reading the shutdown sequence: %w", err)
	}
	if err := v.seq.Lookup(devSeqLost, &lost); err != nil {
		return nil, fmt.Errorf("reading the lost-tombstone count: %w", err)
	}
	data, err := os.ReadFile(hostMountinfo)
	if err != nil {
		return nil, err
	}
	devs, err := vouchedDevs(data)
	if err != nil {
		return nil, err
	}
	if err := deleteKeys(v.vouchMap, func(dev, _ uint64) bool { _, ok := devs[dev]; return !ok }); err != nil {
		return nil, fmt.Errorf("dropping unmounted devs: %w", err)
	}
	for dev := range devs {
		if err := v.vouchMap.Put(dev, seq); err != nil {
			return nil, fmt.Errorf("vouching dev %d:%d: %w", dev>>20, dev&0xfffff, err)
		}
	}
	if err := deleteKeys(v.deadMap, func(_, died uint64) bool { return died <= seq }); err != nil {
		return nil, fmt.Errorf("pruning tombstones: %w", err)
	}
	var lostNow uint64
	if err := v.seq.Lookup(devSeqLost, &lostNow); err != nil {
		return nil, fmt.Errorf("reading the lost-tombstone count: %w", err)
	}
	if lostNow != lost {
		return nil, errors.New("a superblock died during the sync and its tombstone could not be recorded")
	}
	return devs, nil
}

func deleteKeys(m *cilium.Map, drop func(k, v uint64) bool) error {
	var (
		k, v  uint64
		stale []uint64
	)
	it := m.Iterate()
	for it.Next(&k, &v) {
		if drop(k, v) {
			stale = append(stale, k)
		}
	}
	if err := it.Err(); err != nil {
		return err
	}
	for _, key := range stale {
		if err := m.Delete(key); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
			return err
		}
	}
	return nil
}

// watch re-syncs on every change to pid 1's mount table (POLLPRI on mountinfo) and every
// mountResyncInterval, which also prunes tombstones of superblocks that never had a host mount.
func (v *MountVouch) watch(stopFd int) {
	defer close(v.done)
	timeout := mountResyncInterval
	for {
		fds := []unix.PollFd{
			{Fd: int32(v.fd), Events: unix.POLLPRI},  //nolint:gosec // fds are < RLIMIT_NOFILE
			{Fd: int32(stopFd), Events: unix.POLLIN}, //nolint:gosec // same
		}
		if _, err := unix.Poll(fds, int(timeout.Milliseconds())); err != nil && !errors.Is(err, unix.EINTR) {
			v.fail(fmt.Errorf("polling %s: %w", hostMountinfo, err))
			return
		}
		if fds[1].Revents != 0 {
			return
		}
		timeout = mountResyncInterval
		if err := v.Sync(); err != nil {
			log.Errorf("%s: CRITICAL: mount re-sync failed (%v) — system libraries are refused "+
				"to whitelisted binaries until it succeeds", v.name, err)
			timeout = mountRetryInterval
		}
	}
}

// fail empties the vouched set for good: without the watch a later sync could not be trusted to
// stay current.
func (v *MountVouch) fail(err error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.untracked = true
	v.devs = nil
	_ = syncMap(v.vouchMap, map[uint64]uint64{})
	log.Errorf("%s: CRITICAL: %v — mount changes are no longer tracked; system libraries are "+
		"refused to whitelisted binaries until the %s restarts", v.name, err, v.restart)
}

// StartWatch runs the re-sync loop until Stop.
func (v *MountVouch) StartWatch() error {
	efd, err := unix.Eventfd(0, unix.EFD_CLOEXEC)
	if err != nil {
		return fmt.Errorf("creating the mount watch stop fd: %w", err)
	}
	v.stop = efd
	v.done = make(chan struct{})
	go v.watch(efd)
	return nil
}

// Stop ends the watch and closes its fds.
func (v *MountVouch) Stop() {
	if v.done != nil {
		one := []byte{1, 0, 0, 0, 0, 0, 0, 0}
		_, _ = unix.Write(v.stop, one)
		<-v.done
		v.done = nil
	}
	for _, fd := range []*int{&v.stop, &v.fd} {
		if *fd >= 0 {
			_ = unix.Close(*fd)
			*fd = -1
		}
	}
}

// startMountVouch attaches trust_sb_delete, then fills guard_vouched_devs from pid 1's mount table.
// Both precede trust_mmap's attach, so no library is judged against an empty or unguarded map.
func (t *TrustGuard) startMountVouch() error {
	l, err := link.AttachLSM(link.LSMOptions{Program: t.objs.TrustSbDelete})
	if err != nil {
		return fmt.Errorf("required trust hook sb_delete failed to attach: %w — a reused device "+
			"number could inherit a dead superblock's system-library trust", err)
	}
	t.links = append(t.links, l)
	if err := t.mounts.Open(); err != nil {
		return fmt.Errorf("vouching system-library mounts: %w", err)
	}
	return nil
}

// SyncMounts re-syncs guard_vouched_devs to pid 1's current mount table (MountVouch.Sync).
func (t *TrustGuard) SyncMounts() error {
	return t.mounts.Sync()
}
