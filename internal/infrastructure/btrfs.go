package ebpf

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unsafe"

	"github.com/cilium/ebpf/btf"
	"golang.org/x/sys/unix"
)

// BtrfsSuperMagic is statfs's f_type for btrfs (BTRFS_SUPER_MAGIC).
const BtrfsSuperMagic = 0x9123683E

// BtrfsLayout mirrors inode_dev.h's struct btrfs_layout: the byte offsets the BPF objects need to
// read a btrfs inode's subvolume device. Valid is 0 until resolved.
type BtrfsLayout struct {
	Valid       uint32
	InodeRoot   uint32
	InodeVfs    uint32
	RootAnonDev uint32
}

// ResolveBtrfsLayout reads the offsets from the running kernel's BTF: vmlinux when btrfs is
// built in, else the btrfs module's (present only while it is loaded).
func ResolveBtrfsLayout() (BtrfsLayout, error) {
	base, err := btf.LoadKernelSpec()
	if err != nil {
		return BtrfsLayout{}, fmt.Errorf("loading kernel BTF: %w", err)
	}
	l, err := btrfsLayoutFrom(base)
	if err == nil {
		return l, nil
	}
	mod, modErr := btf.LoadKernelModuleSpec("btrfs")
	if modErr != nil {
		return BtrfsLayout{}, fmt.Errorf("%w; btrfs module BTF: %w", err, modErr)
	}
	return btrfsLayoutFrom(mod)
}

type typeSource interface {
	TypeByName(name string, typ any) error
}

func btrfsLayoutFrom(spec typeSource) (BtrfsLayout, error) {
	var ino, root *btf.Struct
	if err := spec.TypeByName("btrfs_inode", &ino); err != nil {
		return BtrfsLayout{}, fmt.Errorf("struct btrfs_inode: %w", err)
	}
	if err := spec.TypeByName("btrfs_root", &root); err != nil {
		return BtrfsLayout{}, fmt.Errorf("struct btrfs_root: %w", err)
	}
	rootOff, err := memberOffset(ino, "root", func(t btf.Type) bool {
		_, ok := t.(*btf.Pointer)
		return ok
	})
	if err != nil {
		return BtrfsLayout{}, err
	}
	vfsOff, err := memberOffset(ino, "vfs_inode", func(t btf.Type) bool {
		s, ok := t.(*btf.Struct)
		return ok && s.Name == "inode"
	})
	if err != nil {
		return BtrfsLayout{}, err
	}
	// inode_dev reads it as a dev_t (u32).
	devOff, err := memberOffset(root, "anon_dev", func(t btf.Type) bool {
		n, sizeErr := btf.Sizeof(t)
		return sizeErr == nil && n == 4
	})
	if err != nil {
		return BtrfsLayout{}, err
	}
	return BtrfsLayout{Valid: 1, InodeRoot: rootOff, InodeVfs: vfsOff, RootAnonDev: devOff}, nil
}

// memberOffset finds name in s, descending into anonymous struct/union members, and returns its
// byte offset. want checks the member's underlying type.
func memberOffset(s *btf.Struct, name string, want func(btf.Type) bool) (uint32, error) {
	off, t, ok := findMember(s.Members, name, 0)
	if !ok {
		return 0, fmt.Errorf("struct %s has no member %s", s.Name, name)
	}
	if off%8 != 0 {
		return 0, fmt.Errorf("struct %s: %s is not byte-aligned", s.Name, name)
	}
	if !want(btf.UnderlyingType(t)) {
		return 0, fmt.Errorf("struct %s: %s has unexpected type %v", s.Name, name, t)
	}
	return uint32(off / 8), nil
}

func findMember(members []btf.Member, name string, base btf.Bits) (btf.Bits, btf.Type, bool) {
	for _, m := range members {
		if m.Name == name {
			if m.BitfieldSize != 0 {
				return 0, nil, false
			}
			return base + m.Offset, m.Type, true
		}
		if m.Name != "" {
			continue
		}
		var inner []btf.Member
		switch c := btf.UnderlyingType(m.Type).(type) {
		case *btf.Struct:
			inner = c.Members
		case *btf.Union:
			inner = c.Members
		}
		if off, t, found := findMember(inner, name, base+m.Offset); found {
			return off, t, true
		}
	}
	return 0, nil, false
}

// BtrfsMounted reports whether any btrfs filesystem is mounted in this mount namespace.
func BtrfsMounted() (bool, error) {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return false, err
	}
	defer f.Close()
	return mountinfoHasFstype(f, "btrfs")
}

func mountinfoHasFstype(r io.Reader, fstype string) (bool, error) {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		_, after, ok := strings.Cut(sc.Text(), " - ")
		if !ok {
			continue
		}
		if f := strings.Fields(after); len(f) > 0 && f[0] == fstype {
			return true, nil
		}
	}
	return false, sc.Err()
}

// OnBtrfs reports whether path lies on btrfs (following symlinks).
func OnBtrfs(path string) (bool, error) {
	var fs unix.Statfs_t
	if err := unix.Statfs(path, &fs); err != nil {
		return false, err
	}
	return fs.Type == BtrfsSuperMagic, nil
}

// FirstOnBtrfs returns the first of paths on btrfs, "" if none. Unstat-able paths are skipped.
func FirstOnBtrfs(paths []string) string {
	for _, p := range paths {
		if on, err := OnBtrfs(p); err == nil && on {
			return p
		}
	}
	return ""
}

// SuperblockDev returns the KernelDev of the superblock holding fd's file: what mountinfo reports
// and the BPF objects' sb_dev reads. That is st_dev everywhere but on btrfs, where st_dev is the
// subvolume's anonymous device and the superblock has its own (inode_dev.h).
func SuperblockDev(fd int) (uint64, error) {
	var fs unix.Statfs_t
	if err := unix.Fstatfs(fd, &fs); err != nil {
		return 0, fmt.Errorf("statfs: %w", err)
	}
	var sx unix.Statx_t
	if err := unix.Statx(fd, "", unix.AT_EMPTY_PATH|unix.AT_STATX_SYNC_AS_STAT, unix.STATX_MNT_ID, &sx); err != nil {
		return 0, fmt.Errorf("statx: %w", err)
	}
	if fs.Type != BtrfsSuperMagic {
		return KernelDev(sx.Dev_major, sx.Dev_minor), nil
	}
	if sx.Mask&unix.STATX_MNT_ID == 0 {
		return 0, errors.New("statx: no mount id")
	}
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return 0, err
	}
	defer f.Close()
	major, minor, err := mountinfoDev(f, sx.Mnt_id)
	if err != nil {
		return 0, err
	}
	return KernelDev(major, minor), nil
}

// SuperblockDevPath is SuperblockDev for path, following symlinks.
func SuperblockDevPath(path string) (uint64, error) {
	fd, err := unix.Open(path, unix.O_PATH|unix.O_CLOEXEC, 0)
	if err != nil {
		return 0, &os.PathError{Op: "open", Path: path, Err: err}
	}
	defer func() { _ = unix.Close(fd) }()
	return SuperblockDev(fd)
}

// mountinfoDev returns the major:minor field of mount mnt's mountinfo line.
func mountinfoDev(r io.Reader, mnt uint64) (major, minor uint32, err error) {
	id := strconv.FormatUint(mnt, 10)
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		f := strings.Fields(sc.Text())
		if len(f) < 3 || f[0] != id {
			continue
		}
		maj, mi, ok := strings.Cut(f[2], ":")
		if !ok {
			return 0, 0, fmt.Errorf("mount %d: malformed device %q", mnt, f[2])
		}
		a, errA := strconv.ParseUint(maj, 10, 32)
		b, errB := strconv.ParseUint(mi, 10, 32)
		if errA != nil || errB != nil {
			return 0, 0, fmt.Errorf("mount %d: malformed device %q", mnt, f[2])
		}
		return uint32(a), uint32(b), nil
	}
	if err := sc.Err(); err != nil {
		return 0, 0, err
	}
	return 0, 0, fmt.Errorf("mount %d not in mountinfo", mnt)
}

// btrfsIoctlFsInfoArgs is struct btrfs_ioctl_fs_info_args (1024 bytes).
type btrfsIoctlFsInfoArgs struct {
	MaxID          uint64
	NumDevices     uint64
	Fsid           [16]byte
	Nodesize       uint32
	Sectorsize     uint32
	CloneAlignment uint32
	CsumType       uint16
	CsumSize       uint16
	Flags          uint64
	Generation     uint64
	MetadataUUID   [16]byte
	Reserved       [944]byte
}

// btrfsIocFsInfo is BTRFS_IOC_FS_INFO: _IOR(0x94, 31, struct btrfs_ioctl_fs_info_args).
const btrfsIocFsInfo = 2<<30 | 1024<<16 | 0x94<<8 | 31

// BtrfsDevices returns the block devices (KernelDev form) of the btrfs filesystem holding path,
// from /sys/fs/btrfs/<fsid>/devices: st_dev on btrfs names no block device.
func BtrfsDevices(path string) ([]uint32, error) {
	fd, err := unix.Open(path, unix.O_RDONLY|unix.O_CLOEXEC, 0)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: path, Err: err}
	}
	defer func() { _ = unix.Close(fd) }()
	var args btrfsIoctlFsInfoArgs
	if _, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), btrfsIocFsInfo,
		uintptr(unsafe.Pointer(&args))); errno != 0 {
		return nil, fmt.Errorf("BTRFS_IOC_FS_INFO on %s: %w", path, errno)
	}
	u := args.Fsid
	fsid := fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:16])
	return btrfsSysfsDevices("/sys/fs/btrfs/" + fsid + "/devices")
}

func btrfsSysfsDevices(dir string) ([]uint32, error) {
	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var out []uint32
	for _, e := range ents {
		raw, err := os.ReadFile(filepath.Join(dir, e.Name(), "dev"))
		if err != nil {
			return nil, err
		}
		maj, mi, ok := strings.Cut(strings.TrimSpace(string(raw)), ":")
		a, errA := strconv.ParseUint(maj, 10, 12)
		b, errB := strconv.ParseUint(mi, 10, 20)
		if !ok || errA != nil || errB != nil {
			return nil, fmt.Errorf("%s/%s/dev: malformed %q", dir, e.Name(), raw)
		}
		out = append(out, uint32(KernelDev(uint32(a), uint32(b)))) //nolint:gosec // 12+20 bits
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("%s lists no device", dir)
	}
	return out, nil
}
