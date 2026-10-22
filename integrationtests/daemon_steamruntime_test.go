package integrationtests

import (
	"fmt"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

const (
	steamLibDir   = "/home/ubuntu/.local/share/Steam/ubuntu12_32"
	steamRTClient = steamLibDir + "/steam"
	steamRT       = steamLibDir + "/steam-runtime"
	steamUser     = "setpriv --reuid=1000 --regid=1000 --clear-groups "
)

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
		` && printf '%s  steam-runtime.tar.xz\n' "$(md5sum < /tmp/out/steam-runtime.tar.xz | cut -d' ' -f1)"` +
		" > /tmp/out/steam-runtime.tar.xz.checksum && chmod -R a+rX /tmp/out && echo BUILT"})
	s.Require().Truef(code == 0 && strings.Contains(out, "BUILT"), "building the %s archive: %s", version, out)
}

// steamUpdates has the client (a whitelisted writer of the lib_dir, so an updater) write the
// archive and its checksum as fresh files, as Steam's self-update does.
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

func (s *IntegrationSuite) steamFile(c testcontainers.Container, path string) string {
	_, out := s.exec(c, []string{"sh", "-c", "cat " + shQuote(path) + " 2>/dev/null || echo MISSING"})
	return strings.TrimSpace(out)
}

// steam.sh unpacks and pins the scout runtime with mkdir/tar/ln, which the read-only lib_dir
// refuses. The daemon does it instead, only for an archive the trust guard proves the client wrote:
// one placed before the daemon started stays packed; the client's own is unpacked as the user, the
// previous tree kept as .old, and setup.sh's pins committed. Pins holding a library are refused.
func (s *IntegrationSuite) TestDaemon_SteamRuntime_UnpacksOnlyTheClientsArchive() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	code, out := s.exec(c, []string{"sh", "-c", "(command -v xz >/dev/null 2>&1 || (timeout 120 apt-get update -qq && " +
		"timeout 120 apt-get install -y -qq --no-install-recommends xz-utils)) && (id -u ubuntu >/dev/null 2>&1 || " +
		"useradd -m -u 1000 ubuntu) && mkdir -p /etc/app-listener /r " + steamRT + " && printf s > /r/secret" +
		" && printf old > " + steamRT + "/version.txt && printf 'old  steam-runtime.tar.xz\\n' > " + steamRT + "/checksum" +
		" && cp /usr/bin/perl " + steamRTClient + " && chown -R 1000:1000 /home/ubuntu/.local && echo READY"})
	s.Require().Truef(code == 0 && strings.Contains(out, "READY"), "setup: %s", out)
	// Placed before the daemon started: no provenance, whoever wrote it.
	s.buildSteamArchive(c, "v0", false)
	s.exec(c, []string{"sh", "-c", "cp /tmp/out/steam-runtime.tar.xz* " + steamLibDir + "/"})

	s.startDaemon(c, fmt.Sprintf("[watch /r]\nneed_encryption: false\n%s\nlib_dir %s", steamRTClient, steamLibDir))
	s.Require().Truef(s.awaitLog(c, "not proven written by Steam's updater", 45*time.Second),
		"an unproven archive must be reported; daemon log:\n%s", s.readDaemonLog(c))
	s.Require().Equal("old", s.steamFile(c, steamRT+"/version.txt"), "an unproven archive must stay packed")

	s.buildSteamArchive(c, "v1", false)
	s.steamUpdates(c)
	s.Require().Truef(s.awaitLog(c, "pinned the Steam runtime", 90*time.Second),
		"the client's archive was not unpacked and pinned; daemon log:\n%s", s.readDaemonLog(c))
	s.Require().Equal("v1", s.steamFile(c, steamRT+"/version.txt"))
	s.Require().Equal("old", s.steamFile(c, steamRT+".old/version.txt"), "the previous tree is kept as .old")
	s.Require().Equal(s.steamFile(c, "/tmp/out/steam-runtime.tar.xz.checksum"), s.steamFile(c, steamRT+"/checksum"))
	_, out = s.exec(c, []string{"sh", "-c", "stat -c '%u' " + steamRT + " " + steamRT + "/version.txt " +
		steamRT + "/pinned_libs_64/libdemo.so.1 | sort -u; readlink " + steamRT + "/pinned_libs_64/libdemo.so.1; " +
		"test -e " + steamRT + "/pinned_libs_32/done && echo DONE; test -e " + steamRT + ".tmp || echo NO_TMP; " +
		"grep -c steam-runtime /proc/self/mounts || true"})
	s.Require().Equalf("1000\n"+steamRT+"/usr/lib/libdemo.so.1\nDONE\nNO_TMP\n0", strings.TrimSpace(out),
		"the unpacked, pinned tree must be the user's, with no stage or overlay left")

	// The grant ended with the run: the lib_dir is read-only again for everyone but its writers.
	for _, who := range []string{steamUser, ""} {
		code, out := s.exec(c, []string{"sh", "-c", who + "touch " + steamLibDir + "/planted 2>&1"})
		s.Require().NotEqualf(0, code, "a non-writer (%q) wrote the lib_dir after the run: %s", who, out)
	}

	s.buildSteamArchive(c, "v2", true)
	s.steamUpdates(c)
	s.Require().Truef(s.awaitLog(c, "refusing setup.sh's output", 90*time.Second),
		"pins holding a library must be refused; daemon log:\n%s", s.readDaemonLog(c))
	s.Require().Equal("v2", s.steamFile(c, steamRT+"/version.txt"), "the client's v2 archive is still unpacked")
	s.Require().Equal("MISSING", s.steamFile(c, steamRT+"/pinned_libs_64/libevil.so"),
		"a library in setup.sh's output must never reach the guarded runtime")

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}
