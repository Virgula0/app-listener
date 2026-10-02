package integrationtests

import (
	"strings"
)

// btrfs subvolumes share one superblock and repeat inode numbers, and stat reports each
// subvolume's own anonymous device, not the superblock's: keys must carry the subvolume.
func (s *IntegrationSuite) TestDaemon_Btrfs_SubvolumeIdentity() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	s.copySwapFixtures(c)

	const (
		marker       = "TOP-SECRET-BTRFS-5E1B"
		nestedMarker = "TOP-SECRET-BTRFS-NESTED-0C7D"
		vault        = "/mnt/b/home/vault"
	)
	// sys holds the whitelisted binary; nobody's own subvolume user gets a copy at the same inode
	// number (each subvolume numbers from 256). The vault is world-readable: only the guard denies.
	build := `
set -e
command -v mkfs.btrfs >/dev/null 2>&1 || (timeout 120 apt-get update -qq && ` +
		`timeout 120 apt-get install -y -qq --no-install-recommends btrfs-progs >/dev/null)
mkdir -p /etc/app-listener /mnt/b
truncate -s 256M /tmp/bimg
mkfs.btrfs -q -f /tmp/bimg
for i in $(seq 0 15); do [ -e /dev/loop$i ] || mknod /dev/loop$i b 7 "$i"; done
mount -o loop /tmp/bimg /mnt/b
btrfs -q subvolume create /mnt/b/sys
btrfs -q subvolume create /mnt/b/home
mkdir -p /mnt/b/sys/bin /mnt/b/drop ` + vault + ` && chmod 1777 /mnt/b/drop
cp ` + swapReaderPath + ` /mnt/b/sys/bin/app
printf '` + marker + `' > ` + vault + `/secret
btrfs -q subvolume create ` + vault + `/nested
printf '` + nestedMarker + `' > ` + vault + `/nested/secret
chmod -R a+rX ` + vault + `
` + nobodyRun + `btrfs -q subvolume create /mnt/b/drop/user
` + nobodyRun + `mkdir /mnt/b/drop/user/bin
` + nobodyRun + `cp ` + swapReaderPath + ` /mnt/b/drop/user/bin/evil
echo MOUNTED`
	code, out := s.exec(c, []string{"sh", "-c", build})
	if code != 0 || !strings.Contains(out, "MOUNTED") {
		s.T().Skipf("btrfs loop mount unavailable here (exit %d): %s", code, out)
	}
	defer s.exec(c, []string{"sh", "-c", "umount /mnt/b 2>/dev/null; losetup -D 2>/dev/null; true"})

	app, evil := s.inodeOf(c, "/mnt/b/sys/bin/app"), s.inodeOf(c, "/mnt/b/drop/user/bin/evil")
	s.Require().NotEmpty(app)
	s.Require().Equalf(app, evil, "fixture: both subvolumes must number the copies alike")

	s.startDaemon(c, `[watch `+vault+`]
need_encryption: false
/mnt/b/sys/bin/app`)
	logs := func() string { return "\ndaemon log:\n" + s.readDaemonLog(c) }

	for _, f := range []struct{ path, marker string }{
		{vault + "/secret", marker},
		{vault + "/nested/secret", nestedMarker},
	} {
		_, out = s.exec(c, []string{"/mnt/b/sys/bin/app", f.path})
		s.Require().Containsf(out, "STOLEN|"+f.marker, "the whitelisted binary reads %s: %s%s", f.path, out, logs())

		_, out = s.exec(c, []string{"sh", "-c", "cat " + f.path + " 2>&1"})
		s.Require().NotContainsf(out, f.marker, "a non-whitelisted binary read %s on btrfs: %s%s", f.path, out, logs())

		_, out = s.exec(c, []string{"sh", "-c", nobodyRun + "/mnt/b/drop/user/bin/evil " + f.path + " 2>&1"})
		s.Require().NotContainsf(out, f.marker, "a user's binary in another subvolume, at the whitelisted "+
			"binary's inode number, read %s: %s%s", f.path, out, logs())
	}

	// A subvolume is created by ioctl, not mkdir: still a write into the guarded directory.
	code, out = s.exec(c, []string{"sh", "-c", "btrfs subvolume create " + vault + "/planted 2>&1"})
	s.Require().NotEqualf(0, code, "a subvolume was created inside the guarded tree: %s%s", out, logs())

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}
