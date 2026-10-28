package integrationtests

import (
	"fmt"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

const (
	steamRoot     = "/home/ubuntu/.local/share/Steam"
	steamLibDir   = steamRoot + "/ubuntu12_32"
	steamRTClient = steamLibDir + "/steam"
	steamRT       = steamLibDir + "/steam-runtime"
	steamSh       = steamRoot + "/steam.sh"
	steamUser     = "setpriv --reuid=1000 --regid=1000 --clear-groups "
	// clientRan is written by the client when the sandbox fails to mask it: outside the overlay and
	// the sandbox's private home, /tmp and /run/user, so the write would survive the run.
	clientRan = "/var/tmp/client-ran"
)

// fakeSteamSh stands in for Valve's launcher: unpack_runtime's steps with the general tools the
// read-only lib_dir refuses, setup.sh's pins, then the client.
const fakeSteamSh = `#!/bin/bash
set -u
STEAMROOT="$(cd "$(dirname "$0")" && echo $PWD)"
srt="$STEAMROOT/ubuntu12_32/steam-runtime"
if [ -f "$srt.tar.xz.checksum" ] && ! cmp -s "$srt.tar.xz.checksum" "$srt/checksum"; then
	rm -rf "$srt.tmp" && mkdir "$srt.tmp" || exit 2
	(cd "$STEAMROOT/ubuntu12_32" && md5sum -c --status steam-runtime.tar.xz.checksum) || exit 3
	tar -xJf "$srt.tar.xz" -C "$srt.tmp" || exit 4
	rm -rf "$srt.old" && mv "$srt" "$srt.old" && mv "$srt.tmp"/* "$srt.tmp"/.. && rm -rf "$srt.tmp" || exit 5
	cp "$srt.tar.xz.checksum" "$srt/checksum" || exit 6
	"$srt/setup.sh" --force || exit 7
fi
exec "$STEAMROOT/ubuntu12_32/steam" -e 'open(F, ">", "` + clientRan + `") or die "$!\n"; print F "ran\n"'
`

// fakeSetupSh stands in for the runtime's setup.sh: what pin_newer_runtime_libs writes (the pin
// directories, a symlink pin, the markers, a system_ record). A runtime holding `evil` also drops a
// library into a pin directory, which the daemon must refuse to commit.
const fakeSetupSh = `#!/bin/bash
set -eu
top="$(cd "$(dirname "$0")" && pwd)"
rm -rf "$top/pinned_libs_32" "$top/pinned_libs_64"
mkdir "$top/pinned_libs_32" "$top/pinned_libs_64"
ln -s "$top/usr/lib/libdemo.so.1" "$top/pinned_libs_64/libdemo.so.1"
libc=/usr/lib/x86_64-linux-gnu/libc.so.6
printf '%s\n%s\n' "$libc" "$(readlink -f "$libc")" > "$top/pinned_libs_64/system_libc.so.6"
touch "$top/pinned_libs_64/has_pins" "$top/pinned_libs_32/done" "$top/pinned_libs_64/done"
if [ -e "$top/evil" ]; then cp /bin/true "$top/pinned_libs_64/libevil.so"; fi
`

// buildSteamArchive writes /tmp/out/steam-runtime.tar.xz{,.checksum} for version, as Valve ships
// them: md5sum's line for the archive in the checksum file.
func (s *IntegrationSuite) buildSteamArchive(c testcontainers.Container, version string, evil bool) {
	extra := ""
	if evil {
		extra = " && touch /tmp/rt/steam-runtime/evil"
	}
	code, out := s.exec(c, []string{"sh", "-c", "rm -rf /tmp/rt /tmp/out && mkdir -p /tmp/rt/steam-runtime/usr/lib /tmp/out" +
		" && printf '" + version + "' > /tmp/rt/steam-runtime/version.txt && printf ELF > /tmp/rt/steam-runtime/usr/lib/libdemo.so.1" +
		" && cat > /tmp/rt/steam-runtime/setup.sh <<'EOF'\n" + fakeSetupSh + "EOF\n" +
		"chmod 755 /tmp/rt/steam-runtime/setup.sh" + extra +
		" && tar -C /tmp/rt -cJf /tmp/out/steam-runtime.tar.xz --owner=0 --group=0 steam-runtime" +
		" && (cd /tmp/out && md5sum steam-runtime.tar.xz > steam-runtime.tar.xz.checksum)" +
		" && chmod -R a+rX /tmp/out && echo BUILT"})
	s.Require().Truef(code == 0 && strings.Contains(out, "BUILT"), "building the %s archive: %s", version, out)
}

// steamUpdates has the client (a whitelisted writer of the lib_dir) write the archive and its
// checksum as fresh files, as Steam's self-update does.
func (s *IntegrationSuite) steamUpdates(c testcontainers.Container) {
	_, out := s.exec(c, []string{"sh", "-c", steamUser + steamRTClient + ` -e '
for my $f ("steam-runtime.tar.xz", "steam-runtime.tar.xz.checksum") {
	open(I, "<", "/tmp/out/$f") or die "read $f: $!\n"; binmode I;
	unlink "` + steamLibDir + `/$f";
	open(O, ">", "` + steamLibDir + `/$f") or die "create $f: $!\n"; binmode O;
	{ local $/; print O <I>; } close O or die "close $f: $!\n";
}
print "WROTE\n"' 2>&1`})
	s.Require().Containsf(out, "WROTE", "the client must write its runtime archive: %s", out)
}

// steamFixture installs a fake Steam for uid 1000 running script as its steam.sh, and returns the
// daemon.conf guarding /r for the client with ubuntu12_32 as its lib_dir. The maintainer's stage
// gets a tmpfs: the container's overlayfs root can't hold an overlay's writable layer. xz is
// installed only for tests that build runtime archives.
func (s *IntegrationSuite) steamFixture(c testcontainers.Container, script string, withXz bool) string {
	setup := ""
	if withXz {
		setup = "(command -v xz >/dev/null 2>&1 || (timeout 120 apt-get update -qq && " +
			"timeout 120 apt-get install -y -qq --no-install-recommends xz-utils)) && "
	}
	code, out := s.exec(c, []string{"sh", "-c", setup + "(id -u ubuntu >/dev/null 2>&1 || useradd -m -u 1000 ubuntu)" +
		" && mkdir -p /etc/app-listener /var/lib/app-listener /r " + steamRT +
		" && mount -t tmpfs -o mode=0755 app-listener-stage /var/lib/app-listener && printf s > /r/secret" +
		" && printf old > " + steamRT + "/version.txt && printf 'old  steam-runtime.tar.xz\\n' > " + steamRT + "/checksum" +
		" && cp /usr/bin/perl " + steamRTClient + " && cat > " + steamSh + " <<'EOF' && chmod 755 " + steamSh +
		" && chown -R 1000:1000 /home/ubuntu/.local && echo READY\n" + script + "EOF\n"})
	s.Require().Truef(code == 0 && strings.Contains(out, "READY"), "setup: %s", out)
	return fmt.Sprintf("[watch /r]\nneed_encryption: false\n%s\nlib_dir %s", steamRTClient, steamLibDir)
}

func (s *IntegrationSuite) steamFile(c testcontainers.Container, path string) string {
	_, out := s.exec(c, []string{"sh", "-c", "cat " + shQuote(path) + " 2>/dev/null || echo MISSING"})
	return strings.TrimSpace(out)
}

// steam.sh unpacks and pins the scout runtime with mkdir/tar/ln, which the read-only lib_dir
// refuses. The daemon runs steam.sh itself in a sandbox, masks the client, and commits only content
// the lib_dir already holds (here: the client's archive and the old tree) plus pins and records.
// Pins holding a foreign library refuse the whole run; a steam.sh changed while the daemon was
// stopped is not run.
func (s *IntegrationSuite) TestDaemon_Maintainer_RunsSteamShSandboxed() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	conf := s.steamFixture(c, fakeSteamSh, true)
	s.startDaemon(c, conf)

	s.buildSteamArchive(c, "v1", false)
	s.steamUpdates(c)
	s.Require().Truef(s.awaitLog(c, "change(s) "+steamSh+" made to its lib_dirs", 120*time.Second),
		"the client's archive was not unpacked and pinned; daemon log:\n%s", s.readDaemonLog(c))
	s.Require().Equal("v1", s.steamFile(c, steamRT+"/version.txt"))
	s.Require().Equal("old", s.steamFile(c, steamRT+".old/version.txt"), "the previous tree is kept as .old")
	s.Require().Equal(s.steamFile(c, "/tmp/out/steam-runtime.tar.xz.checksum"), s.steamFile(c, steamRT+"/checksum"))
	_, out := s.exec(c, []string{"sh", "-c", "stat -c '%u' " + steamRT + " " + steamRT + "/version.txt " +
		steamRT + "/pinned_libs_64/system_libc.so.6 | sort -u; readlink " + steamRT + "/pinned_libs_64/libdemo.so.1; " +
		"test -e " + steamRT + "/pinned_libs_32/done && echo DONE; test -e " + steamRT + ".tmp || echo NO_TMP; " +
		"test -e " + clientRan + " || echo CLIENT_MASKED; ls -A /var/lib/app-listener/maintain 2>/dev/null | wc -l; " +
		"grep -c app-listener-maintain /proc/self/mounts || true"})
	s.Require().Equalf("1000\n"+steamRT+"/usr/lib/libdemo.so.1\nDONE\nNO_TMP\nCLIENT_MASKED\n0\n0", strings.TrimSpace(out),
		"the committed tree must be the user's, the client never started, no stage or overlay left")

	// The grant ended with the commit: the lib_dir is read-only again for everyone but its writers,
	// and steam.sh is a name only Steam may write.
	for _, who := range []string{steamUser, ""} {
		code, out := s.exec(c, []string{"sh", "-c", who + "touch " + steamLibDir + "/planted 2>&1"})
		s.Require().NotEqualf(0, code, "a non-writer (%q) wrote the lib_dir after the run: %s", who, out)
	}
	code, out := s.exec(c, []string{"sh", "-c", steamUser + "sh -c 'echo evil >> " + steamSh + "' 2>&1"})
	s.Require().NotEqualf(0, code, "a non-Steam process appended to steam.sh: %s", out)

	s.buildSteamArchive(c, "v2", true)
	s.steamUpdates(c)
	s.Require().Truef(s.awaitLog(c, "refusing the run's changes to "+steamLibDir, 120*time.Second),
		"pins holding a library must refuse the run; daemon log:\n%s", s.readDaemonLog(c))
	s.Require().Equal("v1", s.steamFile(c, steamRT+"/version.txt"), "nothing of a refused run is committed")
	s.Require().Equal("MISSING", s.steamFile(c, steamRT+"/pinned_libs_64/libevil.so"))

	// steam.sh changed while no daemon guarded it: the next daemon never runs it.
	s.exec(c, []string{"sh", "-c", "pkill -TERM -f 'app-listener daemon'; for i in $(seq 100); do " +
		"pgrep -f 'app-listener daemon' >/dev/null || break; sleep 0.1; done"})
	code, out = s.exec(c, []string{"sh", "-c", steamUser + "sh -c 'echo \"# edited\" >> " + steamSh + "' 2>&1"})
	s.Require().Equalf(0, code, "with no daemon running, nothing guards steam.sh: %s", out)
	s.launchDaemon(c)
	s.awaitDaemonUp(c, conf)
	s.Require().Truef(s.awaitLog(c, steamSh+" was not run: it is not confirmed", 60*time.Second),
		"a steam.sh changed while the daemon was stopped must not run; daemon log:\n%s", s.readDaemonLog(c))

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

const (
	evilJS     = "/home/ubuntu/.cache/evil.js"
	plantedPid = "/var/tmp/planted"
)

// overlayPlanter is same-uid malware outside the sandbox: the run's processes are dumpable after
// setpriv's exec, and their pid namespace hides nothing from the host's /proc, so /proc/<pid>/cwd
// (the Steam root, inside the run's overlay) is theirs to write. It plants the symlink
// ubuntu12_32/<name> -> target and records the pid it went through.
func overlayPlanter(name, target string) string {
	return `mkdir -p /home/ubuntu/.cache && echo 'steal()' > ` + evilJS + `
end=$(($(date +%s) + 150))
while [ "$(date +%s)" -lt "$end" ]; do
	for d in /proc/[0-9]*; do
		if ln -s ` + target + ` "$d/cwd/ubuntu12_32/` + name + `" 2>/dev/null; then
			echo "${d#/proc/}" > ` + plantedPid + `; exit 0
		fi
	done
done`
}

// userImage is the user's own ext4 image, every file root-owned as its author wrote it, mounted at
// /media/ubuntu/EVIL as udisks2 mounts one for an active-session user (nosuid,nodev). Built as
// nobody; only the mount runs as root, standing in for udisksd.
const userImage = `set -e
rm -rf /tmp/u && mkdir -p /tmp/u/tree && echo 'steal()' > '/tmp/u/tree/libraries~00299a408.js'
chmod 755 /tmp/u/tree && chmod 644 /tmp/u/tree/*.js && chown -R 65534:65534 /tmp/u
` + nobodyRun + `env PATH=/usr/sbin:/usr/bin:/sbin:/bin sh -ec '
cd /tmp/u && dd if=/dev/zero of=img bs=1M count=8 status=none
mkfs.ext4 -q -F -E root_owner=0:0 -d tree img
debugfs -w -R "sif /libraries~00299a408.js uid 0" img >/dev/null 2>&1
debugfs -w -R "sif /libraries~00299a408.js gid 0" img >/dev/null 2>&1'
for i in $(seq 0 15); do [ -e /dev/loop$i ] || mknod /dev/loop$i b 7 "$i"; done
mkdir -p /media/ubuntu/EVIL
mount -o nosuid,nodev "$(losetup -f --show /tmp/u/img)" /media/ubuntu/EVIL
stat -c 'owner=%u:%g' /media/ubuntu/EVIL '/media/ubuntu/EVIL/libraries~00299a408.js' | sort -u
echo READY`

// Whatever writes the run's overlay may be hostile: malware of the same uid reaches it through
// /proc/<pid>/cwd while steam.sh runs. A symlink it plants there must not be committed: the lib_dir
// name would hand the whitelisted client the user's file. A link to a directory a user fills is
// the same hole one level down (steamui/libraries/*.js is looked up through it, never judged),
// though root decides the directory's own name. So is a root-owned tree on the user's own mounted
// image, or the name where udisks2 would mount one later. steam.sh sleeps as long as a real runtime
// unpack takes, the window the racer needs.
func (s *IntegrationSuite) TestDaemon_Maintainer_RefusesPlantedSymlink() {
	for _, tc := range []struct{ name, link, target, setup string }{
		{"user file", "planted.js", evilJS, ""},
		{"user home", "libraries", "/home/ubuntu", ""},
		{"world-writable dir", "libraries", "/tmp", ""},
		{"user-mounted image", "libraries", "/media/ubuntu/EVIL", userImage},
		{"user mount point", "libraries", "/media/ubuntu/EVIL", ""},
	} {
		s.Run(tc.name, func() { s.plantedSymlinkRefused(tc.link, tc.target, tc.setup) })
	}
}

func (s *IntegrationSuite) plantedSymlinkRefused(link, target, setup string) {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	conf := s.steamFixture(c, strings.Replace(fakeSteamSh, "\nexec ", "\nsleep 6\nexec ", 1), false)
	if setup != "" {
		code, out := s.exec(c, []string{"sh", "-c", setup})
		if code != 0 || !strings.Contains(out, "READY") {
			s.T().Skipf("unprivileged image build / loop mount unavailable here (exit %d): %s", code, out)
		}
		defer s.exec(c, []string{"sh", "-c", "umount /media/ubuntu/EVIL 2>/dev/null; losetup -j /tmp/u/img | cut -d: -f1 | xargs -r losetup -d; true"})
		s.Require().Containsf(out, "owner=0:0\nREADY", "the image's tree must present as root-owned: %s", out)
	}

	code, out := s.exec(c, []string{"sh", "-c", "nohup " + steamUser + "sh -c " + shQuote(overlayPlanter(link, target)) +
		" >/var/tmp/planter.log 2>&1 &"})
	s.Require().Equalf(0, code, "starting the racer: %s", out)
	s.startDaemon(c, conf)

	_, reached := s.awaitFile(c, plantedPid, 90*time.Second)
	s.Require().Truef(reached, "the racer never reached a run's overlay (this vector is closed); daemon log:\n%s",
		s.readDaemonLog(c))
	var log string
	for dl := time.Now().Add(120 * time.Second); time.Now().Before(dl); time.Sleep(300 * time.Millisecond) {
		if log = s.readDaemonLog(c); strings.Contains(log, "made to its lib_dirs") ||
			strings.Contains(log, "refusing the run's changes") {
			break
		}
	}
	_, out = s.exec(c, []string{"sh", "-c", "readlink " + steamLibDir + "/" + link + " || echo MISSING"})
	s.Require().Equalf("MISSING", strings.TrimSpace(out),
		"EXPLOITED: the daemon committed the racer's symlink from the lib_dir to %s; daemon log:\n%s", target, log)
	s.Require().Containsf(log, "refusing the run's changes to "+steamLibDir+": "+link,
		"the run holding the planted symlink must be refused; daemon log:\n%s", log)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// swapAncestor, as the user: exchanges p with a fresh directory (renameat2 RENAME_EXCHANGE is
// syscall 316 on x86_64), moves it away, then tries to put a tree of its own at the name.
const swapAncestor = `my $p = shift; my $d = "/home/ubuntu/swap" . int(rand(1e9)); mkdir $d; mkdir "$d.2";
print syscall(316, -100, $p, -100, $d, 2) == 0 ? "EXCHANGED\n" : "exchange: $!\n";
print rename($p, "$p.moved") ? "moved away\n" : "rename: $!\n";
if (! -e $p) {
	print mkdir($p) ? "PLANTED\n" : "mkdir: $!\n";
	print rename("$d.2", $p) ? "RENAMED_ONTO\n" : "rename onto: $!\n";
}`

// Apply opens each lib_dir by path: a user who could put a tree of their own at an ancestor's name
// between the run and the commit would aim root's writes (and the content check) at it. Each
// ancestor below the home is a name reserved for Steam's writers (reserveFixedLibDirs): the trust
// guard refuses creating it or renaming onto it, exchanges included. Moving one away is allowed (it
// only breaks Steam) and leaves the lib_dir's path naming nothing. Deepest first, so each move
// leaves the shallower names in place.
func (s *IntegrationSuite) TestDaemon_Maintainer_LibDirAncestorsCannotBeSwapped() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.startDaemon(c, s.steamFixture(c, fakeSteamSh, false))
	_, lib := s.exec(c, []string{"stat", "-c", "%d:%i", steamLibDir})

	for _, p := range []string{steamLibDir, steamRoot, "/home/ubuntu/.local/share", "/home/ubuntu/.local"} {
		_, out := s.exec(c, []string{"sh", "-c", steamUser + "perl -e " + shQuote(swapAncestor) + " " + shQuote(p) + " 2>&1"})
		for _, bad := range []string{"EXCHANGED", "PLANTED", "RENAMED_ONTO"} {
			s.Require().NotContainsf(out, bad, "a non-writer put its own tree at the lib_dir ancestor %s: %s", p, out)
		}
		_, now := s.exec(c, []string{"sh", "-c", "stat -c %d:%i " + steamLibDir + " 2>/dev/null || echo MISSING"})
		s.Require().Truef(now == lib || strings.TrimSpace(now) == "MISSING",
			"after the attempts on %s the lib_dir path names inode %s, not %s: %s", p, now, lib, out)
	}

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}
