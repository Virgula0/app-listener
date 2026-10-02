package ebpf

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cilium/ebpf/btf"
)

type fakeTypes map[string]btf.Type

func (f fakeTypes) TypeByName(name string, typ any) error {
	t, ok := f[name]
	if !ok {
		return btf.ErrNotFound
	}
	p, ok := typ.(**btf.Struct)
	s, isStruct := t.(*btf.Struct)
	if !ok || !isStruct {
		return errors.New("not a struct")
	}
	*p = s
	return nil
}

func btrfsTypes(rootMember btf.Member) fakeTypes {
	u32 := &btf.Int{Name: "unsigned int", Size: 4}
	devT := &btf.Typedef{Name: "dev_t", Type: &btf.Typedef{Name: "__kernel_dev_t", Type: u32}}
	root := &btf.Struct{Name: "btrfs_root", Size: 64, Members: []btf.Member{
		{Name: "node", Type: &btf.Pointer{Target: &btf.Void{}}},
		{Name: "anon_dev", Type: devT, Offset: 40 * 8},
	}}
	rootMember.Type = &btf.Pointer{Target: root}
	inode := &btf.Struct{Name: "inode", Size: 600}
	// root sits in an anonymous union inside an anonymous struct, as kernel layouts often nest.
	anon := &btf.Struct{Size: 16, Members: []btf.Member{
		{Name: "", Offset: 8 * 8, Type: &btf.Union{Size: 8, Members: []btf.Member{rootMember}}},
	}}
	ino := &btf.Struct{Name: "btrfs_inode", Size: 1024, Members: []btf.Member{
		{Name: "", Type: anon, Offset: 16 * 8},
		{Name: "vfs_inode", Type: inode, Offset: 400 * 8},
	}}
	return fakeTypes{"btrfs_inode": ino, "btrfs_root": root}
}

func TestBtrfsLayoutFromNestedMembers(t *testing.T) {
	l, err := btrfsLayoutFrom(btrfsTypes(btf.Member{Name: "root"}))
	if err != nil {
		t.Fatal(err)
	}
	want := BtrfsLayout{Valid: 1, InodeRoot: 24, InodeVfs: 400, RootAnonDev: 40}
	if l != want {
		t.Fatalf("layout = %+v, want %+v", l, want)
	}
}

func TestBtrfsLayoutRejectsUnexpectedShapes(t *testing.T) {
	if _, err := btrfsLayoutFrom(btrfsTypes(btf.Member{Name: "root", BitfieldSize: 3})); err == nil {
		t.Fatal("a bitfield root member resolved")
	}
	if _, err := btrfsLayoutFrom(btrfsTypes(btf.Member{Name: "subvol"})); err == nil {
		t.Fatal("a missing root member resolved")
	}
	types := btrfsTypes(btf.Member{Name: "root"})
	r, ok := types["btrfs_root"].(*btf.Struct)
	if !ok {
		t.Fatal("fixture")
	}
	r.Members[1].Type = &btf.Int{Name: "u64", Size: 8}
	if _, err := btrfsLayoutFrom(types); err == nil {
		t.Fatal("an 8-byte anon_dev resolved: inode_dev reads a dev_t")
	}
	if _, err := btrfsLayoutFrom(fakeTypes{}); err == nil {
		t.Fatal("resolved without btrfs types")
	}
}

// The running kernel's own BTF, when btrfs is built in or loaded: the layout inode_dev.h relies on.
func TestResolveBtrfsLayoutHostBTF(t *testing.T) {
	if _, err := os.Stat("/sys/kernel/btf/vmlinux"); err != nil {
		t.Skip("no kernel BTF")
	}
	l, err := ResolveBtrfsLayout()
	if err != nil {
		if _, modErr := os.Stat("/sys/kernel/btf/btrfs"); modErr != nil {
			t.Skipf("btrfs BTF unavailable: %v", err)
		}
		t.Fatal(err)
	}
	if l.Valid != 1 || l.InodeVfs == 0 || l.RootAnonDev == 0 {
		t.Fatalf("implausible layout %+v", l)
	}
}

const mountinfoFixture = `22 1 259:2 / / rw,relatime shared:1 - ext4 /dev/nvme0n1p2 rw
35 22 0:31 /@home /home rw,relatime shared:2 - btrfs /dev/sda1 rw,subvol=/@home
36 35 0:31 /@home/snap /home/x\040y rw shared:3 - btrfs /dev/sda1 rw
`

func TestMountinfoDev(t *testing.T) {
	major, minor, err := mountinfoDev(strings.NewReader(mountinfoFixture), 35)
	if err != nil || major != 0 || minor != 31 {
		t.Fatalf("mount 35 = %d:%d, %v; want 0:31", major, minor, err)
	}
	if _, _, err := mountinfoDev(strings.NewReader(mountinfoFixture), 3); err == nil {
		t.Fatal("mount 3 found: ids must match whole fields")
	}
	if _, _, err := mountinfoDev(strings.NewReader("40 1 x / / rw - ext4 a rw\n"), 40); err == nil {
		t.Fatal("malformed device parsed")
	}
}

func TestMountinfoHasFstype(t *testing.T) {
	if ok, err := mountinfoHasFstype(strings.NewReader(mountinfoFixture), "btrfs"); err != nil || !ok {
		t.Fatalf("btrfs not found: %v", err)
	}
	ext4Only := strings.SplitAfter(mountinfoFixture, "\n")[0]
	if ok, _ := mountinfoHasFstype(strings.NewReader(ext4Only), "btrfs"); ok {
		t.Fatal("btrfs found on an ext4-only table")
	}
}

func TestSuperblockDevMatchesStatOffBtrfs(t *testing.T) {
	dir := t.TempDir()
	if on, err := OnBtrfs(dir); err != nil || on {
		t.Skip("temp dir on btrfs")
	}
	got, err := SuperblockDevPath(dir)
	if err != nil {
		t.Fatal(err)
	}
	dev, _, err := StatInode(dir)
	if err != nil || got != dev {
		t.Fatalf("SuperblockDev = %d, st_dev = %d (%v)", got, dev, err)
	}
}

func TestBtrfsSysfsDevices(t *testing.T) {
	dir := t.TempDir()
	for name, dev := range map[string]string{"1": "8:3\n", "2": "259:17\n"} {
		if err := os.MkdirAll(filepath.Join(dir, name), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name, "dev"), []byte(dev), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	got, err := btrfsSysfsDevices(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := map[uint32]bool{8<<20 | 3: true, 259<<20 | 17: true}
	if len(got) != 2 || !want[got[0]] || !want[got[1]] {
		t.Fatalf("devices = %v", got)
	}
	if _, err := btrfsSysfsDevices(t.TempDir()); err == nil {
		t.Fatal("an empty devices dir resolved")
	}
}
