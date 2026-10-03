package integrationtests

import (
	"os"
	"strings"

	"github.com/Virgula0/app-listener/internal/guard"
)

// btrfsMountPrelude installs btrfs-progs and builds a 256M loop-mounted btrfs at /mnt/b. Callers
// append their own subvolume layout and a final `echo MOUNTED`, then skip when it is absent (many
// sandboxed kernels forbid the loop mount). Pair with btrfsTeardown.
const btrfsMountPrelude = `
set -e
command -v mkfs.btrfs >/dev/null 2>&1 || (timeout 120 apt-get update -qq && timeout 120 apt-get install -y -qq --no-install-recommends btrfs-progs >/dev/null)
mkdir -p /etc/app-listener /mnt/b
truncate -s 256M /tmp/bimg
mkfs.btrfs -q -f /tmp/bimg
for i in $(seq 0 15); do [ -e /dev/loop$i ] || mknod /dev/loop$i b 7 "$i"; done
mount -o loop /tmp/bimg /mnt/b
`

const btrfsTeardown = "umount /mnt/b 2>/dev/null; losetup -D 2>/dev/null; true"

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
	build := btrfsMountPrelude + `btrfs -q subvolume create /mnt/b/sys
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
	defer s.exec(c, []string{"sh", "-c", btrfsTeardown})

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

// A btrfs snapshot is a new subvolume (new st_dev, repeated inode numbers), so no guard key matches
// the copy; SEND and TREE_SEARCH copy a subvolume's contents — inline data of small files included
// — past the VFS. guard_file_ioctl denies all of them on any btrfs superblock hosting a guarded
// root, judged by GUARD_RES_GLOBAL like the raw block-device gate. Superblock-wide: an unguarded
// path on the same filesystem is denied too (intended; snapper/btrbk stop working on that btrfs).
func (s *IntegrationSuite) TestDaemon_Btrfs_CopyIoctlGate() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	s.copySwapFixtures(c)
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/btrfs_search"), "/exploits/btrfs_search", 0o755), "copy btrfs_search")
	// The 32-bit compat issuer is best-effort in the corpus build (needs multilib); skip its
	// sub-check when absent rather than fail the whole test.
	compat := false
	if _, err := os.Stat(absPath("./exploits/btrfs_search32")); err == nil {
		s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/btrfs_search32"), "/exploits/btrfs_search32", 0o755), "copy btrfs_search32")
		compat = true
	}

	const (
		marker = "TOP-SECRET-BTRFS-IOCTL-9A3F"
		vault  = "/mnt/b/home/vault"
	)
	// home is a subvolume so homero can be a read-only snapshot of it to feed `btrfs send`: a
	// snapshot taken BEFORE the guard (snapper's class), the one copy the gate cannot retract —
	// but on the same superblock, so SEND of it is still denied once the daemon is up.
	build := btrfsMountPrelude + `btrfs -q subvolume create /mnt/b/sys
btrfs -q subvolume create /mnt/b/home
mkdir -p /mnt/b/sys/bin ` + vault + `
cp ` + swapReaderPath + ` /mnt/b/sys/bin/app
printf '` + marker + `' > ` + vault + `/secret
chmod -R a+rX /mnt/b/home
btrfs -q subvolume snapshot -r /mnt/b/home /mnt/b/homero
echo MOUNTED`
	code, out := s.exec(c, []string{"sh", "-c", build})
	if code != 0 || !strings.Contains(out, "MOUNTED") {
		s.T().Skipf("btrfs loop mount unavailable here (exit %d): %s", code, out)
	}
	defer s.exec(c, []string{"sh", "-c", btrfsTeardown})

	// TREE_SEARCH reads the whole superblock's btree from ANY fd on it, so the vector opens an
	// UNGUARDED path on the same btrfs (opening the vault itself is denied by the ordinary
	// file_open guard, a different protection). Before the daemon it returns items: vector live.
	const unguarded = "/mnt/b/sys"
	code, out = s.exec(c, []string{"/exploits/btrfs_search", unguarded})
	s.Require().Equalf(0, code, "fixture: tree-search must work before the guard is up: %s", out)

	s.startDaemon(c, `[watch `+vault+`]
need_encryption: false
/mnt/b/sys/bin/app`)
	logs := func() string { return "\ndaemon log:\n" + s.readDaemonLog(c) }

	// TREE_SEARCH (nr 17) from unguarded paths of the guarded superblock: denied superblock-wide.
	for _, p := range []string{unguarded, "/mnt/b"} {
		code, out = s.exec(c, []string{"/exploits/btrfs_search", p})
		s.Require().Equalf(3, code, "btrfs tree-search of %s was not denied (exit %d): %s%s", p, code, out, logs())
	}

	// The compat (32-bit) hook, file_ioctl_compat on 6.8+.
	if compat {
		code, out = s.exec(c, []string{"/exploits/btrfs_search32", unguarded})
		s.Require().Equalf(3, code, "32-bit compat tree-search was not denied (exit %d): %s%s", code, out, logs())
	} else {
		s.T().Log("skipping the compat-path sub-check: btrfs_search32 not built " +
			"(run `make -C integrationtests/exploits btrfs_search32`, needs a 32-bit toolchain)")
	}

	// SNAP_CREATE_V2 (nr 23): a snapshot of the guarded tree would be an unguarded copy.
	code, out = s.exec(c, []string{"sh", "-c", "btrfs subvolume snapshot -r /mnt/b/home /mnt/b/snap 2>&1"})
	s.Require().NotEqualf(0, code, "a snapshot of a guarded btrfs tree was created: %s%s", out, logs())
	_, out = s.exec(c, []string{"sh", "-c", "cat /mnt/b/snap/vault/secret 2>&1 || true"})
	s.Require().NotContainsf(out, marker, "the snapshot exposed the secret: %s", out)

	// SEND (nr 38) of the pre-guard read-only snapshot, on the same superblock.
	code, out = s.exec(c, []string{"sh", "-c", "btrfs send -f /tmp/send.stream /mnt/b/homero 2>&1; echo EXIT=$?"})
	s.Require().Containsf(out, "EXIT=", "send did not run: %s%s", out, logs())
	s.Require().NotContainsf(out, "EXIT=0", "btrfs send of a snapshot on a guarded superblock succeeded: %s%s", out, logs())
	_, out = s.exec(c, []string{"sh", "-c", "tr -c '[:print:]' '\\n' < /tmp/send.stream 2>/dev/null | grep -q " + marker + " && echo LEAK || echo safe"})
	s.Require().Containsf(out, "safe", "btrfs send streamed the secret off the guarded superblock: %s", out)

	// The exploit's only guarded operation is the tree-search ioctl, so its denial is labelled
	// resource=btrfs-ioctl, never a watched path. (The btrfs CLI also stats other paths, so its
	// comm is not a clean probe for the label.)
	events := parseDaemonEvents(s.readDaemonLog(c))
	gated := false
	for _, ev := range events {
		if !ev.Denied || (ev.Comm != "btrfs_search" && ev.Comm != "btrfs_search32") {
			continue
		}
		gated = true
		s.Require().Equalf(guard.BtrfsIoctlResourceLabel, ev.Resource,
			"btrfs ioctl denial must be labelled %q, got resource=%q (comm=%s path=%s)",
			guard.BtrfsIoctlResourceLabel, ev.Resource, ev.Comm, ev.Path)
	}
	s.Require().Truef(gated, "expected a DENIED btrfs-ioctl event, daemon log:\n%s", s.readDaemonLog(c))

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}
