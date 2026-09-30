package guard

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
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

// startMountVouch attaches trust_sb_delete, then fills guard_vouched_devs from pid 1's mount table.
// Both precede trust_mmap's attach, so no library is judged against an empty or unguarded map. The
// watch fd is opened before the first read so no mount change after it goes unseen.
func (t *TrustGuard) startMountVouch() error {
	l, err := link.AttachLSM(link.LSMOptions{Program: t.objs.TrustSbDelete})
	if err != nil {
		return fmt.Errorf("required trust hook sb_delete failed to attach: %w — a reused device "+
			"number could inherit a dead superblock's system-library trust", err)
	}
	t.links = append(t.links, l)
	fd, err := unix.Open(hostMountinfo, unix.O_RDONLY|unix.O_CLOEXEC, 0)
	if err != nil {
		return fmt.Errorf("opening %s: %w", hostMountinfo, err)
	}
	t.mountFd = fd
	if err := t.SyncMounts(); err != nil {
		return fmt.Errorf("vouching system-library mounts: %w", err)
	}
	return nil
}

// SyncMounts re-syncs guard_vouched_devs to pid 1's current mount table. On failure it drops every
// vouched dev (system libraries become untrusted) rather than keep a possibly stale one.
func (t *TrustGuard) SyncMounts() error {
	t.mountMu.Lock()
	defer t.mountMu.Unlock()
	if t.mountUntracked {
		return errors.New("mount changes are no longer tracked; restart the daemon")
	}
	n, err := t.syncVouchedDevs()
	if err != nil {
		if cerr := syncMap(t.objs.GuardVouchedDevs, map[uint64]uint64{}); cerr != nil {
			err = fmt.Errorf("%w; clearing the vouched set also failed: %w", err, cerr)
		}
		return err
	}
	// Container and snap churn remounts constantly: only the first count is worth an info line.
	if t.vouched < 0 {
		log.Infof("trust guard: %d superblock(s) vouched for system-library auto-trust "+
			"(mounted without nosuid in pid 1's namespace)", n)
	} else if n != t.vouched {
		log.Debugf("trust guard: %d superblock(s) vouched after a mount change", n)
	}
	t.vouched = n
	return nil
}

// syncVouchedDevs: seq snapshot → mountinfo → delete stale → put → prune tombstones. Stale devs go
// first because a dev can be reused by another superblock; entries carry the snapshot so
// trust_mmap refuses one whose superblock died after it (guard_dead_devs).
func (t *TrustGuard) syncVouchedDevs() (int, error) {
	var seq, lost uint64
	if err := t.objs.GuardDevSeq.Lookup(devSeqNow, &seq); err != nil {
		return 0, fmt.Errorf("reading the shutdown sequence: %w", err)
	}
	if err := t.objs.GuardDevSeq.Lookup(devSeqLost, &lost); err != nil {
		return 0, fmt.Errorf("reading the lost-tombstone count: %w", err)
	}
	data, err := os.ReadFile(hostMountinfo)
	if err != nil {
		return 0, err
	}
	devs, err := vouchedDevs(data)
	if err != nil {
		return 0, err
	}
	if err := deleteKeys(t.objs.GuardVouchedDevs, func(dev, _ uint64) bool { _, ok := devs[dev]; return !ok }); err != nil {
		return 0, fmt.Errorf("dropping unmounted devs: %w", err)
	}
	for dev := range devs {
		if err := t.objs.GuardVouchedDevs.Put(dev, seq); err != nil {
			return 0, fmt.Errorf("vouching dev %d:%d: %w", dev>>20, dev&0xfffff, err)
		}
	}
	if err := deleteKeys(t.objs.GuardDeadDevs, func(_, died uint64) bool { return died <= seq }); err != nil {
		return 0, fmt.Errorf("pruning tombstones: %w", err)
	}
	var lostNow uint64
	if err := t.objs.GuardDevSeq.Lookup(devSeqLost, &lostNow); err != nil {
		return 0, fmt.Errorf("reading the lost-tombstone count: %w", err)
	}
	if lostNow != lost {
		return 0, errors.New("a superblock died during the sync and its tombstone could not be recorded")
	}
	return len(devs), nil
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

// watchMounts re-syncs on every change to pid 1's mount table (POLLPRI on mountinfo) and every
// mountResyncInterval, which also prunes tombstones of superblocks that never had a host mount.
func (t *TrustGuard) watchMounts(stopFd int) {
	defer close(t.mountDone)
	timeout := mountResyncInterval
	for {
		fds := []unix.PollFd{
			{Fd: int32(t.mountFd), Events: unix.POLLPRI}, //nolint:gosec // fds are < RLIMIT_NOFILE
			{Fd: int32(stopFd), Events: unix.POLLIN},     //nolint:gosec // same
		}
		if _, err := unix.Poll(fds, int(timeout.Milliseconds())); err != nil && !errors.Is(err, unix.EINTR) {
			t.failMounts(fmt.Errorf("polling %s: %w", hostMountinfo, err))
			return
		}
		if fds[1].Revents != 0 {
			return
		}
		timeout = mountResyncInterval
		if err := t.SyncMounts(); err != nil {
			log.Errorf("trust guard: CRITICAL: mount re-sync failed (%v) — system libraries are refused "+
				"to whitelisted binaries until it succeeds", err)
			timeout = mountRetryInterval
		}
	}
}

// failMounts empties guard_vouched_devs for good: without the watch a later sync could not be
// trusted to stay current.
func (t *TrustGuard) failMounts(err error) {
	t.mountMu.Lock()
	defer t.mountMu.Unlock()
	t.mountUntracked = true
	_ = syncMap(t.objs.GuardVouchedDevs, map[uint64]uint64{})
	log.Errorf("trust guard: CRITICAL: %v — mount changes are no longer tracked; system libraries are "+
		"refused to whitelisted binaries until the daemon restarts", err)
}

// startMountWatch runs watchMounts until stopMountWatch.
func (t *TrustGuard) startMountWatch() error {
	efd, err := unix.Eventfd(0, unix.EFD_CLOEXEC)
	if err != nil {
		return fmt.Errorf("creating the mount watch stop fd: %w", err)
	}
	t.mountStop = efd
	t.mountDone = make(chan struct{})
	go t.watchMounts(efd)
	return nil
}

func (t *TrustGuard) stopMountWatch() {
	if t.mountDone != nil {
		one := []byte{1, 0, 0, 0, 0, 0, 0, 0}
		_, _ = unix.Write(t.mountStop, one)
		<-t.mountDone
		t.mountDone = nil
	}
	for _, fd := range []*int{&t.mountStop, &t.mountFd} {
		if *fd >= 0 {
			_ = unix.Close(*fd)
			*fd = -1
		}
	}
}
