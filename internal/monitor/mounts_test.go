package monitor

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/stretchr/testify/suite"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

const fixtureMountinfo = `22 1 259:2 / / rw,relatime shared:1 - ext4 /dev/nvme0n1p2 rw
30 22 0:40 / /tmp rw,nosuid shared:2 - tmpfs tmpfs rw
31 22 0:40 /sub /mnt/bind\040dir rw shared:2 - tmpfs tmpfs rw
32 22 0:50 /@home /home rw - btrfs /dev/sda2 rw
`

var (
	rootDev  = uint32(ebpf.KernelDev(259, 2))
	tmpfsDev = uint32(ebpf.KernelDev(0, 40))
	btrfsDev = uint32(ebpf.KernelDev(0, 50))
)

type mountsUnitTest struct {
	suite.Suite
}

func TestMountsUnitTest(t *testing.T) {
	suite.Run(t, new(mountsUnitTest))
}

func (s *mountsUnitTest) fixtureTable() *mountTable {
	f := filepath.Join(s.T().TempDir(), "mountinfo")
	s.Require().NoError(os.WriteFile(f, []byte(fixtureMountinfo), 0o600))
	return &mountTable{file: f}
}

func (s *mountsUnitTest) TestResolve() {
	t := s.fixtureTable()
	underTmp := func(p string) bool { _, ok := underRoot(p, "/tmp"); return ok }

	tests := []struct {
		name   string
		dev    uint32
		path   string
		prefer func(string) bool
		want   string
		ok     bool
	}{
		{name: "tmpfs file", dev: tmpfsDev, path: "/foo", want: "/tmp/foo", ok: true},
		{name: "tmpfs root", dev: tmpfsDev, path: "/", want: "/tmp", ok: true},
		{name: "bind mount wins on longer root", dev: tmpfsDev, path: "/sub/x", want: "/mnt/bind dir/x", ok: true},
		{name: "preferred candidate", dev: tmpfsDev, path: "/sub/x", prefer: underTmp, want: "/tmp/sub/x", ok: true},
		{name: "btrfs subvolume", dev: btrfsDev, path: "/@home/u/f", want: "/home/u/f", ok: true},
		{name: "unmounted subvolume", dev: btrfsDev, path: "/@other/f", ok: false},
		{name: "root filesystem", dev: rootDev, path: "/etc/x", want: "/etc/x", ok: true},
		{name: "unknown superblock", dev: 12345, path: "/foo", ok: false},
		{name: "empty path", dev: tmpfsDev, path: "", ok: false},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, ok := t.resolve(tt.dev, tt.path, tt.prefer)
			s.Require().Equal(tt.ok, ok)
			if ok {
				s.Require().Equal(tt.want, got)
			}
		})
	}
}

func (s *mountsUnitTest) TestDevOf() {
	t := s.fixtureTable()
	s.Require().Equal(tmpfsDev, t.devOf("/tmp/a"))
	s.Require().Equal(tmpfsDev, t.devOf("/tmp"))
	s.Require().Equal(rootDev, t.devOf("/tmpx"))
	s.Require().Equal(tmpfsDev, t.devOf("/mnt/bind dir/q"))
	s.Require().Equal(btrfsDev, t.devOf("/home/u"))
}

func (s *mountsUnitTest) TestUnescapeMountinfo() {
	s.Require().Equal("/a b\tc\\d", unescapeMountinfo(`/a\040b\011c\134d`))
	s.Require().Equal(`/trailing\04`, unescapeMountinfo(`/trailing\04`))
}

func (s *mountsUnitTest) newMonitor() *Monitor {
	return &Monitor{
		targets:     dirTarget("/tmp/watch"),
		recursive:   true,
		mounts:      s.fixtureTable(),
		paths:       newPathCache(),
		watchInodes: make(map[uint64][]watchedInode),
	}
}

func (s *mountsUnitTest) TestFilterEventMapsMountRelativePath() {
	m := s.newMonitor()
	ev := &ebpf.FileEvent{Type: ebpf.EventWrite, Path: "/watch/f"}
	s.Require().True(m.filterEvent(ev, &ebpf.BpfEvent{Dev: tmpfsDev, Mode: syscall.S_IFREG, Ino: 7}))
	s.Require().Equal("/tmp/watch/f", ev.Path)
}

// Events outside every target must not touch their file: the path cache stays empty because
// EvalSymlinks never ran, and the inode is not a watched one, so nothing was lstat'd.
func (s *mountsUnitTest) TestFilterEventDoesNotResolveForeignFiles() {
	m := s.newMonitor()
	ev := &ebpf.FileEvent{Type: ebpf.EventStat, Path: "/etc/app-listener/fscrypt.key"}
	s.Require().False(m.filterEvent(ev, &ebpf.BpfEvent{Dev: rootDev, Mode: syscall.S_IFREG, Ino: 9}))
	s.Require().Empty(m.paths.m)

	link := &ebpf.FileEvent{Type: ebpf.EventStat, Path: "/home/u/link"}
	s.Require().False(m.filterEvent(link, &ebpf.BpfEvent{Dev: rootDev, Mode: syscall.S_IFLNK, Ino: 10}))
	s.Require().Len(m.paths.m, 1, "a symlink is still resolved: it may lead into a target")
}

func (s *mountsUnitTest) TestFilterEventUserStringPath() {
	m := s.newMonitor()
	ev := &ebpf.FileEvent{Type: ebpf.EventStat, Path: "/tmp/watch/g"}
	s.Require().True(m.filterEvent(ev, &ebpf.BpfEvent{}))
	s.Require().Equal("/tmp/watch/g", ev.Path)
}

func (s *mountsUnitTest) TestFilterEventRenameDest() {
	m := s.newMonitor()
	ev := &ebpf.FileEvent{Type: ebpf.EventRename, Path: "/watch/old", Dest: "/watch/new"}
	s.Require().True(m.filterEvent(ev, &ebpf.BpfEvent{Dev: tmpfsDev, Mode: syscall.S_IFREG, Ino: 3}))
	s.Require().Equal("/tmp/watch/old", ev.Path)
	s.Require().Equal("/tmp/watch/new", ev.Dest)
}

func (s *mountsUnitTest) TestFilterEventSymlinkRelativeBody() {
	m := s.newMonitor()
	ev := &ebpf.FileEvent{Type: ebpf.EventSymlink, Path: "/watch/lnk", Dest: "../target"}
	s.Require().True(m.filterEvent(ev, &ebpf.BpfEvent{Dev: tmpfsDev}))
	s.Require().Equal("/tmp/watch/lnk", ev.Path)
	s.Require().Equal("/tmp/target", ev.Dest)
}
