package integrationtests

import (
	"fmt"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

// awaitStolen runs bin against secret until it prints marker (true) or timeout. Each denial makes
// the daemon re-sync its whitelist (resyncMinInterval).
func (s *IntegrationSuite) awaitStolen(c testcontainers.Container, bin, secret, marker string,
	timeout time.Duration) (bool, string) {
	var out string
	for dl := time.Now().Add(timeout); time.Now().Before(dl); time.Sleep(time.Second) {
		_, out = s.exec(c, []string{"sh", "-c", bin + " " + secret + " 2>&1"})
		if strings.Contains(out, "STOLEN|"+marker) {
			return true, out
		}
	}
	return false, out
}

// A package upgrade replaces a root-owned whitelisted binary by temp file + rename: the new inode
// is admitted by the re-sync alone, with no hook and no reload.
func (s *IntegrationSuite) TestDaemon_LiveAdmission_SystemBinaryReplacedNoReload() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	s.copySwapFixtures(c)

	const marker = "TOP-SECRET-LIVE-SYSTEM-0B7E"
	s.exec(c, []string{"sh", "-c", "mkdir -p /protected /etc/app-listener /opt/app/bin && printf '" + marker +
		"' > /protected/secret && cp " + swapReaderPath + " /opt/app/bin/app"})
	s.startDaemon(c, `[watch /protected]
need_encryption: false
/opt/app/bin/app`)
	_, out := s.exec(c, []string{"/opt/app/bin/app", "/protected/secret"})
	s.Require().Containsf(out, "STOLEN|"+marker, "baseline: the whitelisted binary reads its resource: %s", out)

	old := s.inodeOf(c, "/opt/app/bin/app")
	code, out := s.exec(c, []string{"sh", "-c", "cp " + swapReaderPath + " /opt/app/bin/.app.new && " +
		"printf x >> /opt/app/bin/.app.new && mv -f /opt/app/bin/.app.new /opt/app/bin/app 2>&1"})
	s.Require().Equalf(0, code, "the upgrade: %s", out)
	s.Require().NotEqual(old, s.inodeOf(c, "/opt/app/bin/app"), "fixture: the upgrade must make a new inode")

	ok, out := s.awaitStolen(c, "/opt/app/bin/app", "/protected/secret", marker, 40*time.Second)
	s.Require().Truef(ok, "the upgraded system binary was not admitted: %s\ndaemon log:\n%s", out, s.readDaemonLog(c))
	s.Require().NotContains(s.readDaemonLog(c), "configuration reloaded", "admission must need no reload")

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// A new match of an absolute catalog glob (JetBrains' /opt/*/jbr/bin/java) is admitted live when
// root installed it; the same name in a directory a user owns is not.
func (s *IntegrationSuite) TestDaemon_LiveAdmission_NewCatalogGlobMatch() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	s.copySwapFixtures(c)

	const (
		marker = "TOP-SECRET-JBR-GLOB-5A61"
		secret = "/root/.config/JetBrains/secret"
	)
	s.exec(c, []string{"sh", "-c", "mkdir -p /etc/app-listener /root/.config/JetBrains /opt/idea/jbr/bin && printf '" +
		marker + "' > " + secret + " && cp " + swapReaderPath + " /opt/idea/jbr/bin/java"})
	s.startDaemon(c, `[watch /root/.config/JetBrains]
need_encryption: false
/opt/idea/jbr/bin/java`)

	s.exec(c, []string{"sh", "-c", "mkdir -p /opt/clion/jbr/bin /opt/evil/jbr/bin && chown 65534 /opt/evil/jbr/bin && " +
		"cp " + swapReaderPath + " /opt/clion/jbr/bin/java && " + nobodyRun + "cp " + swapReaderPath +
		" /opt/evil/jbr/bin/java"})

	ok, out := s.awaitStolen(c, "/opt/clion/jbr/bin/java", secret, marker, 40*time.Second)
	s.Require().Truef(ok, "a root-installed match of /opt/*/jbr/bin/java was not admitted: %s\ndaemon log:\n%s",
		out, s.readDaemonLog(c))
	s.Require().NotContains(s.readDaemonLog(c), "configuration reloaded", "admission must need no reload")

	// A reload rebuilds the trusted set from config paths, which do not name the match: it must
	// keep its library allowlist (TRUSTED_BINARY), or a user's LD_PRELOAD runs inside it.
	s.exec(c, []string{"mkdir", "-p", "/tmp/u"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/lib_probe.so"), "/tmp/u/evil.so", 0o755),
		"copy lib_probe.so")
	s.exec(c, []string{"chown", "-R", "65534:65534", "/tmp/u"})
	s.reloadAndAwait(c)
	ok, out = s.awaitStolen(c, "/opt/clion/jbr/bin/java", secret, marker, 10*time.Second)
	s.Require().Truef(ok, "after a reload the live-admitted match lost its resource: %s", out)
	_, out = s.exec(c, []string{"sh", "-c", "LD_PRELOAD=/tmp/u/evil.so /opt/clion/jbr/bin/java " + secret + " 2>&1"})
	s.Require().NotContainsf(out, libProbeMarker, "after a reload a user library preloaded into the live-admitted "+
		"match: %s", out)

	s.assertNeverStolen(c, "/opt/evil/jbr/bin/java", secret, marker)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

const jetbrainsSecret = "/root/.config/JetBrains/secret"

// startJetBrainsDaemon guards jetbrainsSecret (a catalog section, so /opt/*/jbr/bin/java is a
// system pattern) with a root-installed /opt/idea/jbr/bin/java plus extra lines. The secret holds
// "STOLEN|marker" so a plain reader (grep, sed) prints what the probes look for.
func (s *IntegrationSuite) startJetBrainsDaemon(c testcontainers.Container, marker, setup string, extra ...string) {
	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	s.copySwapFixtures(c)
	cmd := "mkdir -p /etc/app-listener /root/.config/JetBrains /opt/idea/jbr/bin && printf 'STOLEN|" + marker +
		"\\n' > " + jetbrainsSecret + " && cp " + swapReaderPath + " /opt/idea/jbr/bin/java"
	if setup != "" {
		cmd += " && " + setup
	}
	code, out := s.exec(c, []string{"sh", "-c", cmd + " 2>&1"})
	s.Require().Equalf(0, code, "setup: %s", out)
	s.startDaemon(c, "[watch /root/.config/JetBrains]\nneed_encryption: false\n/opt/idea/jbr/bin/java\n"+
		strings.Join(extra, "\n"))
}

// asNobody runs cmd as the unprivileged user and requires it to succeed: it is the attack's setup,
// not what the guard must refuse.
func (s *IntegrationSuite) asNobody(c testcontainers.Container, what, cmd string) {
	code, out := s.exec(c, []string{"sh", "-c", nobodyRun + "sh -c " + shQuote(cmd) + " 2>&1"})
	s.Require().Equalf(0, code, "fixture: %s as the user: %s", what, out)
}

// A whitelisted name in a directory the user owns is the user's to re-point. The trust guard keeps
// the whitelisted inode itself, so the user swaps the parent directory and leaves a symlink to
// root's grep; a glob match in a user-owned /opt dir is a symlink to root's sed. Both land on
// root-owned files in root-owned dirs, but no root action placed the names: neither may be
// admitted. A root package-style upgrade of the root-placed line is still admitted, live.
func (s *IntegrationSuite) TestDaemon_LiveAdmission_UserPlacedNameRefused() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const marker = "TOP-SECRET-USER-PLACED-3D9A"
	s.startJetBrainsDaemon(c, marker, "mkdir -p /opt/uapp /opt/evil && chown 65534:65534 /opt/uapp /opt/evil && "+
		nobodyRun+"sh -c 'mkdir /opt/uapp/bin && cp "+swapReaderPath+" /opt/uapp/bin/tool'", "/opt/uapp/bin/tool")
	_, out := s.exec(c, []string{"sh", "-c", "/opt/uapp/bin/tool " + jetbrainsSecret + " 2>&1"})
	s.Require().Containsf(out, "STOLEN|"+marker, "baseline: the user's whitelisted tool reads the resource: %s", out)

	s.asNobody(c, "swapping the whitelisted tool's directory", "mv /opt/uapp/bin /opt/uapp/old && "+
		"mkdir /opt/uapp/bin && ln -s /usr/bin/grep /opt/uapp/bin/tool")
	s.asNobody(c, "planting a glob match", "mkdir -p /opt/evil/jbr/bin && ln -s /usr/bin/sed /opt/evil/jbr/bin/java")

	old := s.inodeOf(c, "/opt/idea/jbr/bin/java")
	code, out := s.exec(c, []string{"sh", "-c", "cp " + swapReaderPath + " /opt/idea/jbr/bin/.java.new && " +
		"printf x >> /opt/idea/jbr/bin/.java.new && mv -f /opt/idea/jbr/bin/.java.new /opt/idea/jbr/bin/java 2>&1"})
	s.Require().Equalf(0, code, "the upgrade: %s", out)
	s.Require().NotEqual(old, s.inodeOf(c, "/opt/idea/jbr/bin/java"), "fixture: the upgrade must make a new inode")
	ok, out := s.awaitStolen(c, "/opt/idea/jbr/bin/java", jetbrainsSecret, marker, 40*time.Second)
	s.Require().Truef(ok, "control: the root-placed upgrade was not admitted live: %s\ndaemon log:\n%s", out,
		s.readDaemonLog(c))

	// The control's admission proves re-syncs judged both planted names.
	s.assertNeverStolen(c, "/opt/uapp/bin/tool ''", jetbrainsSecret, marker)
	s.assertNeverStolen(c, "/opt/evil/jbr/bin/java -n p", jetbrainsSecret, marker)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// The daemon's own catalog refresh (here the startup one) expands /opt/*/jbr/bin/java. A match the
// user placed, a symlink to root's sed, must not be written to daemon.conf: the reload that follows
// resolves it by path and would admit sed.
func (s *IntegrationSuite) TestDaemon_CatalogRefresh_UserPlacedSystemMatchRefused() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const marker = "TOP-SECRET-REFRESH-PLANT-71C4"
	s.startJetBrainsDaemon(c, marker, "mkdir -p /opt/evil && chown 65534:65534 /opt/evil && "+nobodyRun+
		"sh -c 'mkdir -p /opt/evil/jbr/bin && ln -s /usr/bin/sed /opt/evil/jbr/bin/java'")
	s.Require().Truef(s.awaitLog(c, "re-scanned /root/.config/JetBrains", 30*time.Second),
		"the startup catalog refresh did not run: %s", s.readDaemonLog(c))

	s.assertNeverStolen(c, "/opt/evil/jbr/bin/java -n p", jetbrainsSecret, marker)
	_, conf := s.exec(c, []string{"cat", "/etc/app-listener/daemon.conf"})
	s.Require().NotContainsf(conf, "/opt/evil", "the refresh wrote the user-placed match:\n%s", conf)
	s.Require().NotContainsf(conf, "/usr/bin/sed", "the refresh wrote the match's target:\n%s", conf)

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// An updater replacing a binary several resources whitelist (Proton's wineserver: Steam's config,
// registry.vdf, ...) is re-admitted by every one of them, including guards that re-sync only after
// the prune forgot the freed old inode: they still judge the new inode a replacement of it. perl
// stands in for the app: it is its own updater and swaps itself in by rename, as Steam does.
func (s *IntegrationSuite) TestDaemon_LiveAdmission_UpdaterReplacementSurvivesPrune() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const marker = "TOP-SECRET-PRUNE-RACE-71B3"
	s.exec(c, []string{"sh", "-c", "mkdir -p /etc/app-listener /tmp/apps && for i in 1 2 3 4; do mkdir -p /r$i &&" +
		" printf '" + marker + "' > /r$i/secret; done && cp /usr/bin/perl /tmp/apps/app && chmod 755 /tmp/apps/app"})
	var conf strings.Builder
	for i := 1; i <= 4; i++ {
		fmt.Fprintf(&conf, "[watch /r%d]\nneed_encryption: false\n/tmp/apps/app\n\n", i)
	}
	s.startDaemon(c, conf.String())

	read := func(i int) string {
		_, out := s.exec(c, []string{"sh", "-c", fmt.Sprintf(
			`/tmp/apps/app -e 'open(F, "<", "/r%d/secret") or die "DENIED $!\n"; print "STOLEN|", <F>' 2>&1`, i)})
		return out
	}
	for i := 1; i <= 4; i++ {
		s.Require().Containsf(read(i), "STOLEN|"+marker, "baseline /r%d", i)
	}

	before := s.inodeOf(c, "/tmp/apps/app")
	_, out := s.exec(c, []string{"sh", "-c", `/tmp/apps/app -e '
open(I, "<", "/usr/bin/perl") or die "read $!\n"; binmode I;
open(O, ">", "/tmp/apps/app.new") or die "create $!\n"; binmode O;
{ local $/; print O <I>; } close O or die "close $!\n";
chmod 0755, "/tmp/apps/app.new" or die "chmod $!\n";
rename("/tmp/apps/app.new", "/tmp/apps/app") or die "rename $!\n"; print "UPDATED\n"' 2>&1`})
	s.Require().Containsf(out, "UPDATED", "the app must update itself: %s", out)
	s.Require().NotEqual(before, s.inodeOf(c, "/tmp/apps/app"), "the update must be a new inode")

	// No read until the prune has forgotten the old inode: a denial would re-sync its guard first.
	s.Require().Truef(s.awaitLog(c, "forgot freed binary inode", 70*time.Second),
		"the old inode was never pruned; daemon log:\n%s", s.readDaemonLog(c))
	for i := 1; i <= 4; i++ {
		deadline := time.Now().Add(45 * time.Second)
		for time.Now().Before(deadline) && !strings.Contains(read(i), "STOLEN|"+marker) {
			time.Sleep(2 * time.Second)
		}
		s.Require().Containsf(read(i), "STOLEN|"+marker,
			"/r%d never re-admitted its updater's replacement; daemon log:\n%s", i, s.readDaemonLog(c))
	}
	s.Require().NotContains(s.readDaemonLog(c), "not re-admitted", "no guard may refuse the updater's inode")

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

const (
	relinkMarker = "TOP-SECRET-RELINK-6E2D"
	relinkLink   = "/tmp/apps/bin/app"
	relinkRead   = relinkLink + " /tmp/read.pl"
)

// relinkScripts: read.pl prints the file it is given, hold.pl answers probeAsk's "read" as an app
// left running across an update, relink.pl src dst is the updater: it writes dst from src unless
// src is empty, then points the link at dst by symlink + rename, as Claude Code's installer does.
const relinkScripts = `cat > /tmp/read.pl <<'PL'
open(F, "<", $ARGV[0]) or die "DENIED $!\n"; print "STOLEN|", <F>;
PL
cat > /tmp/hold.pl <<'PL'
while (1) {
	if (unlink "/tmp/ctl/read") {
		my $r = open(F, "<", "/r/secret") ? "STOLEN|" . join("", <F>) : "DENIED $!";
		close F;
		open(O, ">", "/tmp/ctl/read.tmp") or die; print O $r; close O;
		rename "/tmp/ctl/read.tmp", "/tmp/ctl/read.out";
	}
	select(undef, undef, undef, 0.2);
}
PL
cat > /tmp/relink.pl <<'PL'
my ($src, $dst) = @ARGV;
if ($src ne "") {
	open(I, "<", $src) or die "read $!\n"; binmode I;
	open(O, ">", $dst) or die "create $!\n"; binmode O;
	{ local $/; print O <I>; } close O or die "close $!\n";
	chmod 0755, $dst or die "chmod $!\n";
}
unlink "/tmp/apps/bin/.app.tmp";
symlink($dst, "/tmp/apps/bin/.app.tmp") or die "symlink $!\n";
rename("/tmp/apps/bin/.app.tmp", "/tmp/apps/bin/app") or die "rename $!\n";
print "RELINKED\n";
PL`

// relink runs the updater through the link: the running version points it at dst.
func (s *IntegrationSuite) relink(c testcontainers.Container, src, dst string) {
	_, out := s.exec(c, []string{"sh", "-c", relinkLink + " /tmp/relink.pl " + shQuote(src) + " " + dst + " 2>&1"})
	s.Require().Containsf(out, "RELINKED", "the updater must re-point its link: %s", out)
}

// An updater re-points its whitelisted link at a version it wrote (Claude Code's
// ~/.local/bin/claude -> versions/<v>): the new target is admitted live, no reload, and the old one
// keeps only the process started before, across a reload too. A target the updater did not create
// stays refused even when the updater itself points the link at it.
func (s *IntegrationSuite) TestDaemon_LiveAdmission_UpdaterRepointsLink() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	code, out := s.exec(c, []string{"sh", "-c", "mkdir -p /etc/app-listener /r /tmp/apps/versions /tmp/apps/bin " +
		supersedeCtl + " && printf '" + relinkMarker + "' > /r/secret && cp /usr/bin/perl /tmp/apps/versions/1 && " +
		"ln -s /tmp/apps/versions/1 " + relinkLink + " && " + relinkScripts + "\necho READY"})
	s.Require().Truef(code == 0 && strings.Contains(out, "READY"), "setup: %s", out)
	s.startDaemon(c, "[watch /r]\nneed_encryption: false\n"+relinkLink)
	_, out = s.exec(c, []string{"sh", "-c", relinkRead + " /r/secret 2>&1"})
	s.Require().Containsf(out, "STOLEN|"+relinkMarker, "baseline: the linked binary reads its resource: %s", out)
	s.exec(c, []string{"sh", "-c", "nohup " + relinkLink + " /tmp/hold.pl >/dev/null 2>&1 &"})
	s.Require().Contains(s.probeAsk(c, "read"), "STOLEN|"+relinkMarker, "baseline: the running app reads")

	s.relink(c, "/usr/bin/perl", "/tmp/apps/versions/2")
	ok, out := s.awaitStolen(c, relinkRead, "/r/secret", relinkMarker, 40*time.Second)
	s.Require().Truef(ok, "the version the updater linked was not admitted: %s\ndaemon log:\n%s", out, s.readDaemonLog(c))
	s.Require().Contains(s.readDaemonLog(c), "re-pointed to /tmp/apps/versions/2", "admitted by following the link")
	s.Require().NotContains(s.readDaemonLog(c), "configuration reloaded", "admission must need no reload")

	assertOldSplit := func(when string) {
		s.Require().Containsf(s.probeAsk(c, "read"), "STOLEN|"+relinkMarker,
			"%s: the app started before the update lost its resource", when)
		_, out := s.exec(c, []string{"sh", "-c", "/tmp/apps/versions/1 /tmp/read.pl /r/secret 2>&1"})
		s.Require().NotContainsf(out, relinkMarker, "%s: a new exec of the previous target read the resource: %s",
			when, out)
	}
	assertOldSplit("after the update")
	s.reloadAndAwait(c)
	s.Require().Contains(s.readDaemonLog(c), reloadEnd, "the reload must commit")
	assertOldSplit("after a reload")
	_, out = s.exec(c, []string{"sh", "-c", relinkRead + " /r/secret 2>&1"})
	s.Require().Containsf(out, "STOLEN|"+relinkMarker, "after a reload the linked version lost its resource: %s", out)

	// A file root's cp wrote has no updater provenance: the updater linking it admits nothing.
	s.exec(c, []string{"cp", "/usr/bin/perl", "/tmp/apps/versions/3"})
	s.relink(c, "", "/tmp/apps/versions/3")
	s.assertNeverStolen(c, relinkRead, "/r/secret", relinkMarker)
	s.Require().Contains(s.readDaemonLog(c), "which was neither created by its updater", "the planted target is refused")

	s.exec(c, []string{"sh", "-c", "pkill -f hold.pl; pkill -f 'app-listener daemon' || true"})
}
