package guard

import (
	"fmt"
	"testing"

	"golang.org/x/sys/unix"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

func TestVouchedDevs(t *testing.T) {
	const mountinfo = `22 1 259:2 / / rw,relatime shared:1 - ext4 /dev/nvme0n1p2 rw
23 22 0:21 / /proc rw,nosuid,nodev,noexec,relatime shared:5 - proc proc rw
24 22 0:40 / /tmp rw,nosuid,nodev shared:20 - tmpfs tmpfs rw
25 22 7:3 / /snap/core/1 ro,nodev,relatime shared:30 - squashfs /dev/loop3 ro
26 22 7:4 / /media/u/img rw,nosuid,nodev,relatime shared:31 - ext4 /dev/loop4 rw
27 22 0:50 / /home/u/fuse rw,relatime shared:32 - fuse.bindfs bindfs rw,user_id=0
28 22 8:1 / /mnt/usb rw,relatime shared:33 - fuseblk /dev/sda1 rw
29 22 0:51 / /home/u/x rw,relatime shared:34 - fuse sshfs rw
30 22 259:3 / /boot rw,nosuid shared:35 - vfat /dev/nvme0n1p1 rw
31 22 259:3 / /efi rw,relatime shared:36 - vfat /dev/nvme0n1p1 rw
32 22 0:40 /x /var/tmp rw,nosuid shared:21 - tmpfs tmpfs rw
`
	got, err := vouchedDevs([]byte(mountinfo))
	if err != nil {
		t.Fatal(err)
	}
	want := map[uint64]bool{
		ebpf.KernelDev(259, 2): true,  // /
		ebpf.KernelDev(7, 3):   true,  // snap squashfs: nodev only
		ebpf.KernelDev(259, 3): true,  // one suid mount is enough (/efi)
		ebpf.KernelDev(0, 21):  false, // every mount nosuid
		ebpf.KernelDev(0, 40):  false,
		ebpf.KernelDev(7, 4):   false, // udisks2 user image
		ebpf.KernelDev(0, 50):  false, // FUSE, whatever the mount flags
		ebpf.KernelDev(8, 1):   false,
		ebpf.KernelDev(0, 51):  false,
	}
	for dev, vouched := range want {
		if _, ok := got[dev]; ok != vouched {
			t.Errorf("dev %d:%d vouched=%v, want %v", dev>>20, dev&0xfffff, ok, vouched)
		}
	}
	if len(got) != 3 {
		t.Errorf("vouched %d devs, want 3: %v", len(got), got)
	}
}

func TestVouchedDevsFailsClosed(t *testing.T) {
	for name, mi := range map[string]string{
		"empty":          "",
		"all nosuid":     "22 1 259:2 / / rw,nosuid shared:1 - ext4 /dev/root rw\n",
		"no separator":   "22 1 259:2 / / rw,relatime shared:1 ext4 /dev/root rw\n",
		"bad device":     "22 1 259-2 / / rw - ext4 /dev/root rw\n",
		"minor overflow": "22 1 7:1048576 / / rw - ext4 /dev/root rw\n",
		"major overflow": "22 1 4096:1 / / rw - ext4 /dev/root rw\n",
		"no fstype":      "22 1 259:2 / / rw -\n",
		"one bad line":   "22 1 259:2 / / rw - ext4 /dev/root rw\ngarbage\n",
	} {
		if devs, err := vouchedDevs([]byte(mi)); err == nil {
			t.Errorf("%s: vouched %v, want an error", name, devs)
		}
	}
}

// mountinfo prints sb->s_dev, which trust_mmap keys on; it must encode like StatInode's dev.
func TestParseMajorMinorMatchesStatEncoding(t *testing.T) {
	for _, path := range []string{"/", "/proc"} {
		var st unix.Stat_t
		if err := unix.Stat(path, &st); err != nil {
			t.Fatal(err)
		}
		dev, err := parseMajorMinor(fmt.Sprintf("%d:%d", unix.Major(st.Dev), unix.Minor(st.Dev)))
		want, _, serr := ebpf.StatInode(path)
		if err != nil || serr != nil || dev != want {
			t.Errorf("%s: parseMajorMinor = %#x (%v), StatInode dev = %#x (%v)", path, dev, err, want, serr)
		}
	}
}
