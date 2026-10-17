package ebpf

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type mountinfoUnitTest struct {
	suite.Suite
}

func TestMountinfoUnitTest(t *testing.T) {
	suite.Run(t, new(mountinfoUnitTest))
}

func (s *mountinfoUnitTest) TestParseMountinfo() {
	mounts, err := ParseMountinfo(strings.NewReader(`22 1 259:2 / / rw shared:1 - ext4 /dev/nvme0n1p2 rw
garbage
31 22 0:40 /sub /mnt/bind\040dir rw - tmpfs tmpfs rw
32 22 4096:1 / /bad rw - x x rw
`))
	s.Require().NoError(err)
	s.Require().Equal([]MountInfo{
		{ID: 22, Dev: uint32(KernelDev(259, 2)), Root: "/", Mountpoint: "/"},
		{ID: 31, Dev: uint32(KernelDev(0, 40)), Root: "/sub", Mountpoint: "/mnt/bind dir"},
	}, mounts, "a malformed line and a major over 12 bits are skipped")
}

func (s *mountinfoUnitTest) TestMountPath() {
	tests := []struct {
		name   string
		m      MountInfo
		sbPath string
		want   string
		ok     bool
	}{
		{name: "root filesystem", m: MountInfo{Root: "/", Mountpoint: "/"}, sbPath: "/etc/x", want: "/etc/x", ok: true},
		{name: "tmpfs file", m: MountInfo{Root: "/", Mountpoint: "/tmp"}, sbPath: "/a/b", want: "/tmp/a/b", ok: true},
		{name: "tmpfs root", m: MountInfo{Root: "/", Mountpoint: "/tmp"}, sbPath: "/", want: "/tmp", ok: true},
		{name: "subvolume", m: MountInfo{Root: "/@home", Mountpoint: "/home"}, sbPath: "/@home/u/.ssh", want: "/home/u/.ssh", ok: true},
		{name: "subvolume root", m: MountInfo{Root: "/@home", Mountpoint: "/home"}, sbPath: "/@home", want: "/home", ok: true},
		{name: "outside the mounted part", m: MountInfo{Root: "/@home", Mountpoint: "/home"}, sbPath: "/@homeX/f", ok: false},
		{name: "relative path", m: MountInfo{Root: "/", Mountpoint: "/tmp"}, sbPath: "a", ok: false},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			got, ok := tt.m.MountPath(tt.sbPath)
			s.Require().Equal(tt.ok, ok)
			s.Require().Equal(tt.want, got)
		})
	}
}

func (s *mountinfoUnitTest) TestUnescapeMountinfo() {
	s.Require().Equal("/a b\tc\\d", unescapeMountinfo(`/a\040b\011c\134d`))
	s.Require().Equal(`/trailing\04`, unescapeMountinfo(`/trailing\04`))
}
