package ebpf

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// DurableStat is a BinaryStat a later daemon run can still vouch for: plus the birth time (a new
// file reusing a freed inode number gets a new one) and the filesystem id.
type DurableStat struct {
	BinaryStat
	BtimeNs int64
	Fsid    uint64
}

// durableFS: local filesystems whose timestamps only the kernel sets. FUSE, network, overlay and
// tmpfs are left out: their stat comes from userspace or another host, or doesn't outlive a boot.
func durableFS(fsType int64) bool {
	switch fsType {
	case unix.EXT4_SUPER_MAGIC, unix.XFS_SUPER_MAGIC, unix.BTRFS_SUPER_MAGIC, unix.F2FS_SUPER_MAGIC:
		return true
	}
	return false
}

func fsidOf(s *unix.Statfs_t) uint64 {
	return uint64(uint32(s.Fsid.Val[0]))<<32 | uint64(uint32(s.Fsid.Val[1])) //nolint:gosec // bit reinterpretation
}

// StatDurable returns the DurableStat of the inode f holds when its fingerprint can vouch for its
// content across restarts: it still equals fp (taken on f), the filesystem is durableFS, the file
// has a birth time, and, seen from pid 1's root, the same inode sits on a mount that isn't nosuid
// (a user-mounted image carries whatever timestamps its author wrote; mountVouchesOwnership).
func StatDurable(f *os.File, fp BinaryStat) (DurableStat, bool) {
	var stx unix.Statx_t
	if err := unix.Statx(int(f.Fd()), "", unix.AT_EMPTY_PATH|unix.AT_STATX_SYNC_AS_STAT,
		unix.STATX_BASIC_STATS|unix.STATX_BTIME, &stx); err != nil || stx.Mask&unix.STATX_BTIME == 0 {
		return DurableStat{}, false
	}
	ds := DurableStat{
		BinaryStat: BinaryStat{
			Dev:     KernelDev(stx.Dev_major, stx.Dev_minor),
			Ino:     stx.Ino,
			Size:    int64(stx.Size), //nolint:gosec // st_size fits int64
			MtimeNs: stx.Mtime.Sec*1_000_000_000 + int64(stx.Mtime.Nsec),
			CtimeNs: stx.Ctime.Sec*1_000_000_000 + int64(stx.Ctime.Nsec),
		},
		BtimeNs: stx.Btime.Sec*1_000_000_000 + int64(stx.Btime.Nsec),
	}
	if ds.BinaryStat != fp {
		return DurableStat{}, false
	}
	var sfs unix.Statfs_t
	if err := unix.Fstatfs(int(f.Fd()), &sfs); err != nil || !durableFS(sfs.Type) {
		return DurableStat{}, false
	}
	ds.Fsid = fsidOf(&sfs)
	if !hostVouches(f, fp, ds.Fsid) {
		return DurableStat{}, false
	}
	return ds, true
}

// hostVouches: f's path, from pid 1's root, reaches the same inode on the same filesystem through
// a mount that isn't nosuid. The unit's own namespace may add nosuid everywhere, so it can't judge.
func hostVouches(f *os.File, fp BinaryStat, fsid uint64) bool {
	resolved, err := os.Readlink(fmt.Sprintf("/proc/self/fd/%d", f.Fd()))
	if err != nil {
		return false
	}
	host, err := hostView(resolved)
	if err != nil {
		return false
	}
	var st unix.Stat_t
	if unix.Stat(host, &st) != nil || KernelDev(unix.Major(st.Dev), unix.Minor(st.Dev)) != fp.Dev || st.Ino != fp.Ino {
		return false
	}
	var hfs unix.Statfs_t
	if unix.Statfs(host, &hfs) != nil {
		return false
	}
	return durableFS(hfs.Type) && hfs.Flags&unix.ST_NOSUID == 0 && fsidOf(&hfs) == fsid
}

// SyncFile writes back f's dirty pages, write-protecting its shared mappings: from here on any
// content change faults and bumps ctime. Without it, a writer holding an mmap can keep changing
// pages already dirty when the fingerprint is taken, with no ctime change, until writeback.
func SyncFile(f *os.File) error {
	r, err := os.Open(fmt.Sprintf("/proc/self/fd/%d", f.Fd()))
	if err != nil {
		return err
	}
	defer r.Close()
	return r.Sync()
}
