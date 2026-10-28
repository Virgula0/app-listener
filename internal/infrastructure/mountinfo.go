package ebpf

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

// HostMountinfo is pid 1's mount table: the daemon's own namespace (PrivateTmp, ProtectSystem) is
// remounted nosuid by systemd and says nothing about who mounted what.
const HostMountinfo = "/proc/1/mountinfo"

// MountInfo is one /proc/self/mountinfo line. BPF dentry walks stop at the superblock root, so
// their paths are relative to the filesystem, not to "/": MountPath maps one through a mount.
type MountInfo struct {
	ID         uint64
	Dev        uint32 // KernelDev encoding: the superblock's s_dev, as BPF reads it
	Root       string // path inside the filesystem that is mounted here
	Mountpoint string
}

// MountPath returns sbPath (relative to the superblock root) as an absolute path through m;
// false when sbPath is outside the part of the filesystem m mounts.
func (m MountInfo) MountPath(sbPath string) (string, bool) {
	rel, ok := PathUnder(sbPath, m.Root)
	if !ok {
		return "", false
	}
	switch {
	case m.Mountpoint == "/":
		return rel, true
	case rel == "/":
		return m.Mountpoint, true
	}
	return m.Mountpoint + rel, true
}

// PathUnder returns p relative to root (as "/..."), or false when p is outside root.
func PathUnder(p, root string) (string, bool) {
	switch {
	case root == "/":
		return p, strings.HasPrefix(p, "/")
	case p == root:
		return "/", true
	case strings.HasPrefix(p, root+"/"):
		return p[len(root):], true
	}
	return "", false
}

// ReadMountinfo parses /proc/self/mountinfo.
func ReadMountinfo() ([]MountInfo, error) {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseMountinfo(f)
}

// MountByID returns the mount with id (statx STATX_MNT_ID).
func MountByID(id uint64) (MountInfo, error) {
	mounts, err := ReadMountinfo()
	if err != nil {
		return MountInfo{}, err
	}
	for _, m := range mounts {
		if m.ID == id {
			return m, nil
		}
	}
	return MountInfo{}, fmt.Errorf("mount %d not in mountinfo", id)
}

// ParseMountinfo parses mountinfo lines in order (a later mount shadows an earlier one at the same
// mountpoint), skipping malformed ones.
func ParseMountinfo(r io.Reader) ([]MountInfo, error) {
	mounts := make([]MountInfo, 0, 64)
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		// id parent major:minor root mountpoint options ...
		f := strings.Fields(sc.Text())
		if len(f) < 5 {
			continue
		}
		id, err := strconv.ParseUint(f[0], 10, 64)
		if err != nil {
			continue
		}
		dev, ok := parseMountDev(f[2])
		if !ok {
			continue
		}
		mounts = append(mounts, MountInfo{
			ID:         id,
			Dev:        dev,
			Root:       unescapeMountinfo(f[3]),
			Mountpoint: unescapeMountinfo(f[4]),
		})
	}
	return mounts, sc.Err()
}

// VouchedDevs returns the kernel dev (KernelDev) of every superblock in mountinfo with at least one
// mount that is neither nosuid nor FUSE: only there does root vouch for the ownership bits a file
// presents (a user's udisks2 image or FUSE server claims any owner). Any unparsable line fails the
// whole table.
func VouchedDevs(mountinfo []byte) (map[uint64]struct{}, error) {
	out := make(map[uint64]struct{})
	for line := range strings.SplitSeq(strings.TrimSpace(string(mountinfo)), "\n") {
		f := strings.Fields(line)
		sep := slices.Index(f, "-")
		if sep < 6 || sep+1 >= len(f) {
			return nil, fmt.Errorf("malformed mountinfo line %q", line)
		}
		dev, ok := parseMountDev(f[2])
		if !ok {
			return nil, fmt.Errorf("mountinfo line %q: bad device", line)
		}
		if slices.Contains(strings.Split(f[5], ","), "nosuid") || isFuseType(f[sep+1]) {
			continue
		}
		out[uint64(dev)] = struct{}{}
	}
	if len(out) == 0 {
		return nil, errors.New("no superblock mounted without nosuid (is / really nosuid?)")
	}
	return out, nil
}

func isFuseType(fstype string) bool {
	return fstype == "fuse" || fstype == "fuseblk" || strings.HasPrefix(fstype, "fuse.")
}

// hostVouchedDevs is VouchedDevs of HostMountinfo.
func hostVouchedDevs() (map[uint64]struct{}, error) {
	data, err := os.ReadFile(HostMountinfo)
	if err != nil {
		return nil, err
	}
	return VouchedDevs(data)
}

// vouchedFd: fd's superblock (SuperblockDev: not st_dev on btrfs) is one of devs.
func vouchedFd(fd int, devs map[uint64]struct{}) error {
	dev, err := SuperblockDev(fd)
	if err != nil {
		return err
	}
	if _, ok := devs[dev]; !ok {
		return errors.New("is on a mount root does not vouch for (nosuid or FUSE: a user's image or server)")
	}
	return nil
}

// parseMountDev parses "major:minor" into the kernel dev_t (12-bit major, 20-bit minor).
func parseMountDev(s string) (uint32, bool) {
	major, minor, ok := strings.Cut(s, ":")
	if !ok {
		return 0, false
	}
	maj, err := strconv.ParseUint(major, 10, 12)
	if err != nil {
		return 0, false
	}
	mnr, err := strconv.ParseUint(minor, 10, 20)
	if err != nil {
		return 0, false
	}
	return uint32(maj)<<20 | uint32(mnr), true
}

// unescapeMountinfo decodes the kernel's \ooo octal escapes (space, tab, newline, backslash).
func unescapeMountinfo(s string) string {
	if !strings.Contains(s, `\`) {
		return s
	}
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+4 <= len(s) {
			if v, err := strconv.ParseUint(s[i+1:i+4], 8, 8); err == nil {
				b.WriteByte(byte(v))
				i += 3
				continue
			}
		}
		b.WriteByte(s[i])
	}
	return b.String()
}
