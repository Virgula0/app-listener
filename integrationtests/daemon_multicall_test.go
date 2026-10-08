package integrationtests

import (
	"fmt"
	"path"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

const (
	uutilsDir       = "/usr/lib/cargo/bin/coreutils"
	uutilsMulticall = "/usr/bin/coreutils"
	mcRoot          = "/mc"
	mcBypassRoot    = "/mc2"
	mcExploit       = "/exploits/multicall_exec"
	mcProbeTable    = "/mc-probes.tsv"
	// mcBits resources: each probe's identity is the set of resources admitting it, a 7-bit code.
	mcBits = 7
)

// uutilsVettedVersions are the rust-coreutils builds this test's applet table was vetted against.
// ubuntu:latest has shipped both (26.04 LTS moved 0.8.0 -> 0.10.0), and the rootful and rootless
// docker image caches can hold different ones, so either is accepted. Their `--list` is identical;
// only the version string and the stray hash links differ (0.10.0 adds b3sum), so the security-
// critical applet set below is one table and the strays are derived from the image.
var uutilsVettedVersions = []string{"0.8.0", "0.10.0"}

// uutilsApplets is `coreutils --list` of the vetted rust-coreutils builds: the build's applet set,
// which differs from the link names in uutilsDir (the strays below; hostname, kill, more and uptime
// have no link).
var uutilsApplets = strings.Fields(`[ arch b2sum base32 base64 basename basenc cat chcon chgrp chmod
	chown chroot cksum comm cp csplit cut date dd df dir dircolors dirname du echo env expand expr factor
	false fmt fold groups head hostid hostname id install join kill link ln logname ls md5sum mkdir mkfifo
	mknod mktemp more mv nice nl nohup nproc numfmt od paste pathchk pinky pr printenv printf ptx pwd
	readlink realpath rm rmdir runcon seq sha1sum sha224sum sha256sum sha384sum sha512sum shred shuf sleep
	sort split stat stdbuf stty sum sync tac tail tee test timeout touch tr true truncate tsort tty uname
	unexpand uniq unlink uptime users vdir wc who whoami yes`)

// uutilsStrayAllowed are the links in uutilsDir naming no applet, across the vetted builds: uutils
// runs the longest applet name the link ends with (b3sum, sha3sum, hashsum and shake128sum all end
// in "sum", so run sum) or exits "unknown program" (relpath). b3sum is 0.10.0-only. A uutilsDir link
// that is neither an applet nor one of these fails vetting.
var uutilsStrayAllowed = []string{"b3sum", "hashsum", "relpath", "sha3-224sum", "sha3-256sum",
	"sha3-384sum", "sha3-512sum", "sha3sum", "shake128sum", "shake256sum"}

// multiNameScan prints "<dev>:<ino> <path>" for every regular file reachable in the image's binary
// directories, symlinks followed: one inode under several basenames is the shape of a multicall.
const multiNameScan = `for d in /bin /sbin /usr/bin /usr/sbin /usr/libexec; do [ -d "$d" ] || continue
  for f in "$d"/*; do [ -f "$f" ] && stat -Lc "%d:%i %n" "$f"; done
done 2>/dev/null
true` // end on a zero exit: the trailing [ -f ] test would otherwise set the script's status

// vettedMultiNameGroups is every inode of ubuntu:latest's binary directories reachable under more
// than one name, besides the uutils multicall, vetted on 26.04.1 LTS. A name may be a glob (perl's
// versioned alias). Each group is one program under several names, or one file that picks its
// behaviour from argv[0] (bash/rbash, tune2fs/e2label): those identities do collapse, which is
// accepted here because each group is small, known, and ships that way upstream.
//
// A group appearing, growing or losing a name means the image's tooling changed. Re-vet it before
// editing this list: ClassifyMulticall tells apart only vetted families (uutils coreutils). A NEW
// family packing many tools into one file (a uutils findutils/diffutils switch, Ubuntu's next Rust
// swap) is refused when it prints "multi-call binary" (FamilyUnrecognized) and only warned about
// otherwise (warnMultiLinked); without one, every one of its tools shares one whitelist identity.
var vettedMultiNameGroups = [][]string{
	{"i386", "linux32", "linux64", "setarch", "x86_64"},                          // util-linux
	{"dnsdomainname", "domainname", "hostname", "nisdomainname", "ypdomainname"}, // hostname
	{"mke2fs", "mkfs.ext2", "mkfs.ext3", "mkfs.ext4"},                            // e2fsprogs
	{"e2fsck", "fsck.ext2", "fsck.ext3", "fsck.ext4"},                            // e2fsprogs
	{"captoinfo", "infotocap", "tic"},                                            // ncurses-bin
	{"awk", "mawk", "nawk"},                                                      // mawk
	{"agetty", "getty"},                                                          // util-linux
	{"bash", "rbash"},                                                            // bash
	{"cp", "gnucp"},                                                              // gnu-coreutils
	{"dash", "sh"},                                                               // dash
	{"df", "gnudf"},                                                              // gnu-coreutils
	{"dumpe2fs", "e2mmpstatus"},                                                  // e2fsprogs
	{"e2label", "tune2fs"},                                                       // e2fsprogs
	{"gnumv", "mv"},                                                              // gnu-coreutils
	{"gnurm", "rm"},                                                              // gnu-coreutils
	{"gnutrue", "true"},                                                          // gnu-coreutils
	{"gunzip", "uncompress"},                                                     // gzip
	{"killall5", "pidof"},                                                        // sysvinit-utils
	{"more", "pager"},                                                            // util-linux
	{"perl", "perl5.*"},                                                          // perl-base
	{"pgrep", "pkill"},                                                           // procps
	{"reset", "tset"},                                                            // ncurses-bin
	{"rmt", "rmt-tar"},                                                           // tar
	{"skill", "snice"},                                                           // procps
	{"vigr", "vipw"},                                                             // passwd
	{"which", "which.debianutils"},                                               // debianutils
}

// mcProbe is one exec of a binary whose stdout (stderr when stderr) is a file in every resource. code
// is the set of resources that must admit it; 0 = none. listed: path is the whitelist line.
type mcProbe struct {
	label  string
	path   string
	args   string
	stderr bool
	code   int
	listed bool
}

// uutilsUnlinked are applets of the build with no link in uutilsDir.
var uutilsUnlinked = []string{"hostname", "kill", "more", "uptime"}

// uutilsDirLinks is the sorted listing of uutilsDir in the container.
func (s *IntegrationSuite) uutilsDirLinks(c testcontainers.Container) []string {
	_, out := s.exec(c, []string{"sh", "-c", "ls -A " + uutilsDir})
	got := strings.Fields(out)
	slices.Sort(got)
	return got
}

// uutilsStrays are the container's uutilsDir links that name no applet (version-specific; derived,
// then vetted against uutilsStrayAllowed by requireVettedUutils).
func uutilsStrays(links []string) []string {
	var strays []string
	for _, l := range links {
		if !slices.Contains(uutilsApplets, l) {
			strays = append(strays, l)
		}
	}
	return strays
}

// requireVettedUutils pins the container's uutils build: the applet set is strict (a different one
// must be re-vetted, not silently tested against the wrong table), while the version and stray links
// may be any of the vetted builds (uutilsVettedVersions / uutilsStrayAllowed).
func (s *IntegrationSuite) requireVettedUutils(c testcontainers.Container) {
	_, ver := s.exec(c, []string{uutilsMulticall, "--version"})
	vetted := false
	for _, v := range uutilsVettedVersions {
		vetted = vetted || strings.Contains(ver, "coreutils "+v+" (multi-call binary)")
	}
	s.Require().Truef(vetted, "ubuntu:latest ships an unvetted uutils build %q (vetted: %v)", ver, uutilsVettedVersions)

	_, out := s.exec(c, []string{uutilsMulticall, "--list"})
	s.Require().Equal(uutilsApplets, strings.Fields(out), "uutils applet set drifted: re-vet the build")

	links := s.uutilsDirLinks(c)
	for _, a := range uutilsApplets {
		if slices.Contains(uutilsUnlinked, a) {
			s.Require().NotContainsf(links, a, "applet %s unexpectedly has a link now: re-vet the build", a)
			continue
		}
		s.Require().Containsf(links, a, "applet %s has no link in %s: re-vet the build", a, uutilsDir)
	}
	for _, l := range links {
		if slices.Contains(uutilsApplets, l) || slices.Contains(uutilsStrayAllowed, l) {
			continue
		}
		s.Require().Failf("unvetted uutils link",
			"%s/%s names neither an applet nor a known stray: re-vet the build", uutilsDir, l)
	}
	// Every APPLET link is the one multicall inode, or the test proves nothing about applet identity.
	// Stray links are excluded: some builds ship relpath (an "unknown program" to uutils) as its own
	// separate binary, which is harmless — strays are whitelisted nowhere regardless.
	appletPaths := []string{uutilsMulticall}
	for _, a := range uutilsApplets {
		if !slices.Contains(uutilsUnlinked, a) {
			appletPaths = append(appletPaths, uutilsDir+"/"+a)
		}
	}
	_, out = s.exec(c, append([]string{"stat", "-c", "%d:%i"}, appletPaths...))
	inodes := strings.Fields(out)
	slices.Sort(inodes)
	inodes = slices.Compact(inodes)
	s.Require().Lenf(inodes, 1, "uutils applet links are not one multicall inode: %v", inodes)
}

// uutilsToUpperLayer relinks every name of the uutils multicall to one copy in the container's
// writable layer. overlayfs copies an image-layer file up on link(2), under a new inode, which
// breaks its hard links: a test's `ln -f /usr/bin/coreutils x` would leave x and coreutils a
// different inode from every applet link (and a later test's one-inode check failing). Checked
// pristine first, so the relink can't hide an image whose links really differ. Once per container.
func (s *IntegrationSuite) uutilsToUpperLayer(c testcontainers.Container) {
	const done = "/.mc-upper-layer"
	if code, _ := s.exec(c, []string{"test", "-e", done}); code == 0 {
		return
	}
	s.requireVettedUutils(c)
	relink := "set -e; old=$(stat -c %i " + uutilsMulticall + "); cp -p " + uutilsMulticall + " " + uutilsDir + "/.mc-new; " +
		"for f in " + uutilsMulticall + " " + uutilsDir + "/*; do [ \"$(stat -c %i \"$f\")\" = \"$old\" ] || continue; " +
		"ln -f " + uutilsDir + "/.mc-new \"$f.mc-tmp\"; mv -f \"$f.mc-tmp\" \"$f\"; done; rm " + uutilsDir + "/.mc-new; touch " + done
	code, out := s.exec(c, []string{"bash", "-c", relink + " 2>&1"})
	s.Require().Equalf(0, code, "moving uutils to the writable layer: %s", out)
	s.requireVettedUutils(c)
}

// multicallProbes assigns applet i (1-based, uutilsApplets order) the code i, so every applet is
// admitted by a distinct resource set. Each applet runs from its uutilsDir link, and from /usr/bin
// when that name resolves to the multicall (Ubuntu routes cp, df, mv, rm, true to GNU). An applet
// with no link runs from a /usr/local/bin symlink to the multicall, which the daemon refuses to
// whitelist (the name exec'd must be the applet's own link): code 0, like the stray links. GNU
// binaries (own inodes) are the control: inode identity, unchanged.
func (s *IntegrationSuite) multicallProbes(c testcontainers.Container) []mcProbe {
	_, out := s.exec(c, []string{"sh", "-c", "for a in /usr/bin/*; do [ \"$(readlink -f \"$a\")\" = " +
		uutilsDir + "/\"${a##*/}\" ] && echo \"${a##*/}\"; done"})
	viaUsrBin := strings.Fields(out)
	links := s.uutilsDirLinks(c)

	var probes []mcProbe
	for i, a := range uutilsApplets {
		args, stderr := "--version", false
		if a == "test" {
			args, stderr = "1 -eq x", true // test prints nothing for --version unless invoked as [
		}
		p := mcProbe{label: a, path: uutilsDir + "/" + a, args: args, stderr: stderr, code: i + 1, listed: true}
		if !slices.Contains(links, a) {
			p = mcProbe{label: a + " (no link)", path: "/usr/local/bin/" + a, args: args}
		}
		probes = append(probes, p)
		if slices.Contains(viaUsrBin, a) {
			u := p
			u.label, u.path, u.listed = "/usr/bin/"+a, "/usr/bin/"+a, false
			probes = append(probes, u)
		}
	}
	for _, a := range uutilsStrays(links) {
		probes = append(probes, mcProbe{label: a + " (stray link)", path: uutilsDir + "/" + a, args: "--version"})
	}
	for i, g := range []string{"/usr/bin/gnucat", "/usr/bin/gnuhead"} {
		probes = append(probes, mcProbe{label: g, path: g, args: "--version", code: len(uutilsApplets) + 1 + i,
			listed: true})
	}
	s.Require().Less(len(uutilsApplets)+2, 1<<mcBits, "codes must fit in %d resources", mcBits)
	return probes
}

// multicallConfig whitelists, per resource j, every probe binary whose code has bit j, plus dash: the
// shell opens each probe's output file, so every write the probe itself makes is judged by its own
// identity alone. Resource 0 also lists the refused symlinked applets.
func multicallConfig(probes []mcProbe) string {
	var b strings.Builder
	for j := range mcBits {
		fmt.Fprintf(&b, "[watch %s/r%d]\nneed_encryption: false\n/usr/bin/dash\n", mcRoot, j)
		if j == 0 {
			for _, a := range uutilsUnlinked {
				fmt.Fprintf(&b, "/usr/local/bin/%s\n", a)
			}
		}
		for _, p := range probes {
			if p.listed && p.code&(1<<j) != 0 {
				fmt.Fprintf(&b, "%s\n", p.path)
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

// decodeMulticall names what an observed code means.
func decodeMulticall(probes []mcProbe, code int) string {
	switch code {
	case 0:
		return "admitted by no resource"
	case 1<<mcBits - 1:
		return "admitted by every resource (one identity for the whole inode)"
	}
	for _, p := range probes {
		if p.code == code {
			return "judged as " + p.label
		}
	}
	return "matches no binary's code"
}

// Policy keys on the exe inode, and every uutils applet is one inode: whitelisting `head` admitted
// `cat`, `dd`, `tee`. Each applet must be its own identity, whichever path runs it, while a name the
// build has no applet for is admitted nowhere and GNU binaries keep plain inode identity.
//
// Applet i is whitelisted by the resources of the bits of i, and writes its --version into a file in
// every resource; the bytes that landed, read back by the kernel's verdict and nothing else, give the
// code the guard judged it by.
func (s *IntegrationSuite) TestDaemon_Multicall_EveryUutilsAppletIsItsOwnIdentity() {
	c := s.multicallContainer()
	s.requireVettedUutils(c)
	s.requireDistinctExes(c, uutilsMulticall, "/usr/bin/gnucat", "/usr/bin/gnuhead", "/usr/bin/dash", "/usr/bin/gnutimeout")

	probes := s.multicallProbes(c)
	var table strings.Builder
	for i, p := range probes {
		stream := 1
		if p.stderr {
			stream = 2
		}
		fmt.Fprintf(&table, "%d\t%d\t%s\t%s\n", i, stream, p.path, p.args)
	}
	setup := "set -e; mkdir -p /etc/app-listener /usr/local/bin" +
		"; for a in " + strings.Join(uutilsUnlinked, " ") + "; do ln -sf " + uutilsMulticall + " /usr/local/bin/$a; done" +
		"; cat > " + mcProbeTable + " <<'EOF'\n" + table.String() + "EOF\n" +
		"for j in $(seq 0 " + strconv.Itoa(mcBits-1) + "); do mkdir -p " + mcRoot + "/r$j" +
		"; for i in $(seq 0 " + strconv.Itoa(len(probes)-1) + "); do : > " + mcRoot + "/r$j/p$i.out; done; done"
	code, out := s.exec(c, []string{"sh", "-c", setup + " 2>&1"})
	s.Require().Equalf(0, code, "setup: %s", out)

	s.startDaemon(c, multicallConfig(probes))
	defer s.exec(c, []string{"sh", "-c", "pkill -TERM -f 'app-listener daemon' || true"})
	log := s.readDaemonLog(c)
	for _, a := range uutilsUnlinked {
		s.Require().Regexpf(`CRITICAL: `+mcRoot+`/r0 whitelist line /usr/local/bin/`+regexp.QuoteMeta(a)+` dropped`, log,
			"a symlink of another name must not whitelist applet %s", a)
	}

	// dash builtins only: any other binary would be judged too. gnutimeout (own inode, never
	// whitelisted) bounds a probe without touching its output.
	run := `while IFS='	' read -r i s p a; do
  j=0
  while [ $j -lt ` + strconv.Itoa(mcBits) + ` ]; do
    f=` + mcRoot + `/r$j/p$i.out
    if [ "$s" = 2 ]; then gnutimeout 20 "$p" $a </dev/null 2>"$f" >/dev/null
    else gnutimeout 20 "$p" $a </dev/null >"$f" 2>/dev/null; fi
    j=$((j+1))
  done
done < ` + mcProbeTable + `
while IFS='	' read -r i s p a; do
  j=0
  while [ $j -lt ` + strconv.Itoa(mcBits) + ` ]; do
    if read -r l < ` + mcRoot + `/r$j/p$i.out; then echo "HIT $i $j $l"; fi
    j=$((j+1))
  done
done < ` + mcProbeTable + `
echo PROBES-DONE`
	_, out = s.exec(c, []string{"sh", "-c", run})
	s.Require().Containsf(out, "PROBES-DONE", "probe script did not finish: %s", out)

	got := make([]int, len(probes))
	for line := range strings.SplitSeq(out, "\n") {
		f := strings.Fields(line)
		if len(f) < 3 || f[0] != "HIT" {
			continue
		}
		i, ierr := strconv.Atoi(f[1])
		j, jerr := strconv.Atoi(f[2])
		s.Require().NoErrorf(ierr, "bad HIT line %q", line)
		s.Require().NoErrorf(jerr, "bad HIT line %q", line)
		got[i] |= 1 << j
	}

	var wrong []string
	for i, p := range probes {
		if got[i] != p.code {
			wrong = append(wrong, fmt.Sprintf("  %-22s want %07b (%s), got %07b: %s", p.label, p.code,
				decodeMulticall(probes, p.code), got[i], decodeMulticall(probes, got[i])))
		}
	}
	s.Require().Emptyf(wrong, "%d of %d probes judged by the wrong identity:\n%s\nprobe output:\n%s\ndaemon log tail:\n%s",
		len(wrong), len(probes), strings.Join(wrong, "\n"), tailLast(out, 40), tailLast(s.readDaemonLog(c), 40))
}

// mcNameGroup is one inode of the image reachable under several names.
type mcNameGroup struct {
	names []string
	path  string
}

// multiNameInodes groups the image's binary directories by inode, keeping the inodes reachable under
// more than one basename (symlink farms included: multiNameScan follows symlinks).
func (s *IntegrationSuite) multiNameInodes(c testcontainers.Container) map[string]mcNameGroup {
	code, out := s.exec(c, []string{"sh", "-c", multiNameScan})
	s.Require().Equalf(0, code, "scanning the image's binaries: %s", tailLast(out, 20))

	groups := map[string]mcNameGroup{}
	for line := range strings.SplitSeq(out, "\n") {
		ino, p, ok := strings.Cut(strings.TrimSpace(line), " ")
		if !ok || !strings.HasPrefix(p, "/") {
			continue
		}
		g := groups[ino]
		if name := path.Base(p); !slices.Contains(g.names, name) {
			g.names = append(g.names, name)
			g.path = p
			groups[ino] = g
		}
	}
	for ino, g := range groups {
		if len(g.names) < 2 {
			delete(groups, ino)
			continue
		}
		slices.Sort(g.names)
		groups[ino] = g
	}
	return groups
}

// vettedGroupIndex is the vettedMultiNameGroups entry names matches, or -1.
func vettedGroupIndex(names []string) int {
	for i, want := range vettedMultiNameGroups {
		if len(want) != len(names) {
			continue
		}
		matched := true
		for j, pat := range want {
			if ok, err := path.Match(pat, names[j]); err != nil || !ok {
				matched = false
				break
			}
		}
		if matched {
			return i
		}
	}
	return -1
}

// A multicall family the classifier doesn't recognize shares one whitelist identity across its
// tools: refused if it calls itself a multi-call binary, else only a warning. ebpf.ClassifyMulticall
// works by content markers, so this pins the image's whole inventory of "one inode, several names". Ubuntu swapping another package for a Rust
// multicall build (findutils, diffutils, whatever follows) lands here as a new group, naming the
// package that needs a classifier and an applet table before its tools can be told apart.
func (s *IntegrationSuite) TestDaemon_Multicall_NoUnvettedMulticallFamilies() {
	c := s.multicallContainer()
	s.requireVettedUutils(c)

	code, out := s.exec(c, []string{"stat", "-Lc", "%d:%i", uutilsMulticall})
	s.Require().Equalf(0, code, "stat %s: %s", uutilsMulticall, out)
	mcInode := strings.TrimSpace(out)

	groups := s.multiNameInodes(c)
	s.Require().Containsf(groups, mcInode, "the uutils multicall is not reachable under several names in the image's bin dirs")

	seen := make([]bool, len(vettedMultiNameGroups))
	var unexpected []string
	for ino, g := range groups {
		if ino == mcInode {
			continue // requireVettedUutils pins this one, applet by applet
		}
		if i := vettedGroupIndex(g.names); i >= 0 {
			seen[i] = true
			continue
		}
		_, owner := s.exec(c, []string{"sh", "-c", "dpkg -S " + g.path + " 2>/dev/null | head -1 || true"})
		unexpected = append(unexpected, fmt.Sprintf("  %s runs as %d names: %s  [%s]",
			g.path, len(g.names), strings.Join(g.names, " "), owner))
	}
	var missing []string
	for i, want := range vettedMultiNameGroups {
		if !seen[i] {
			missing = append(missing, "  "+strings.Join(want, " "))
		}
	}
	slices.Sort(unexpected)

	s.Require().Emptyf(unexpected, "%d inode(s) of this image run under several names and are not vetted:\n%s\n\n"+
		"If one file packs many tools and picks between them by argv[0], whitelisting any of its names admits all of "+
		"them: teach ebpf.ClassifyMulticall that family (content markers + its applet list) and give it a case in the "+
		"applet matrix test. If it is only an alias of one program, add it to vettedMultiNameGroups.",
		len(unexpected), strings.Join(unexpected, "\n"))
	s.Require().Emptyf(missing, "%d vetted multi-name group(s) are gone from the image:\n%s\n\n"+
		"The image's tooling changed; re-vet and update vettedMultiNameGroups.", len(missing), strings.Join(missing, "\n"))
}

// mcBypass is one attempt to exec the uutils multicall as the whitelisted applet `head` without the
// kernel attesting head: each writes --version to its own file in the guarded resource. admit=false
// means the exec is MC_TAG_NONE (or another applet) and no row covers it, so file_permission must
// deny the write and the file stays empty; admit=true is a positive control that must land bytes, so
// a run where everything is denied cannot pass vacuously. run is executed by a whitelisted dash
// (`dash -c 'run >f'`): dash opens the file as itself before exec'ing run, so the open is always
// allowed and the bytes that land are judged by the exec'd image's own identity alone.
type mcBypass struct {
	label string
	setup string // optional: plant a fixture before the daemon classifies the whitelist
	run   string // exec'd as: bash -c "exec <run> >$f 2>/dev/null"
	admit bool
}

// multicallBypasses whitelists only the head link; every negative shape dispatches to head (or cat)
// but under an identity the kernel will not attest as head.
func multicallBypasses() []mcBypass {
	headLink := uutilsDir + "/head"
	catLink := uutilsDir + "/cat"
	return []mcBypass{
		// Positive controls.
		{label: "head link (control)", run: headLink + " --version", admit: true},
		{label: "symlink head->cat (control)", setup: "ln -sf " + catLink + " /usr/local/bin/head",
			run: "/usr/local/bin/head --version", admit: true}, // same inode, attested head: identity is the name, not the path
		{label: "gnuhead own inode (control)", run: "/usr/bin/gnuhead --version", admit: true},

		// exec -a: argv[0] forced (via a nested bash, whose exec replaces it), so argv[0] disagrees
		// with the kernel filename. The nested bash never writes the file; it execs the applet.
		{label: "exec -a head on cat link", run: `bash -c "exec -a head ` + catLink + ` --version"`, admit: false},
		{label: "exec -a cat on head link", run: `bash -c "exec -a cat ` + headLink + ` --version"`, admit: false},

		// argv[1] dispatch: the exec'd name ("coreutils"/"*utils") is not an applet.
		{label: "coreutils head", run: uutilsMulticall + " head --version", admit: false},
		{label: "catutils argv1 head", setup: "ln -f " + uutilsMulticall + " /usr/local/bin/catutils",
			run: "/usr/local/bin/catutils head --version", admit: false},

		// Suffix dispatch: gcat->cat, xhead.sh->head; the name is not an applet.
		{label: "suffix gcat", setup: "ln -f " + uutilsMulticall + " /usr/local/bin/gcat",
			run: "/usr/local/bin/gcat --version", admit: false},
		{label: "suffix xhead.sh", setup: "ln -f " + uutilsMulticall + " /usr/local/bin/xhead.sh",
			run: "/usr/local/bin/xhead.sh --version", admit: false},

		// Hidden filename: fexecve/execveat/proc-fd exec, argv0 head. uutils falls back to argv0 for
		// these paths and runs head; the kernel filename basename is a number, never "head".
		{label: "fexecve", run: mcExploit + " fexecve " + headLink + " head", admit: false},
		{label: "execveat AT_EMPTY_PATH", run: mcExploit + " execveat " + headLink + " head", admit: false},
		{label: "/proc/self/fd exec", run: mcExploit + " pathfd " + headLink + " head", admit: false},

		// A #! script whose interpreter is the head link: bprm->interp != bprm->filename, so no tag.
		{label: "script with head interpreter", setup: "printf '#!" + headLink + "\\n' > /usr/local/bin/hdscript; chmod +x /usr/local/bin/hdscript",
			run: "/usr/local/bin/hdscript", admit: false},
	}
}

// The guard keys a uutils multicall on (inode, attested applet). Whitelisting the `head` applet must
// admit head and nothing else out of that one inode — not cat, not head reached under any name the
// kernel cannot attest as head (exec -a, argv[1] dispatch, suffix names, fexecve/execveat/proc-fd,
// a script interpreter). Bytes land only for the three controls; every bypass file stays empty.
func (s *IntegrationSuite) TestDaemon_Multicall_BypassCorpus() {
	c := s.multicallContainer()
	s.requireVettedUutils(c)
	s.requireDistinctExes(c, uutilsMulticall, "/usr/bin/gnuhead", "/usr/bin/dash")
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/multicall_exec"), mcExploit, 0o755),
		"copy multicall_exec")

	cases := multicallBypasses()
	var setups []string
	for _, b := range cases {
		if b.setup != "" {
			setups = append(setups, b.setup)
		}
	}
	// Pre-create every output file so the daemon populates its inode into the guard at start, as the
	// matrix test does: a file first created at probe time might escape the subtree enforcement and
	// let a denied write land vacuously.
	precreate := fmt.Sprintf("for i in $(seq 0 %d); do : > %s/b$i.out; done", len(cases)-1, mcBypassRoot)
	setup := "set -e; mkdir -p /etc/app-listener /usr/local/bin " + mcBypassRoot + "; " +
		strings.Join(append(setups, precreate), "; ")
	code, out := s.exec(c, []string{"bash", "-c", setup + " 2>&1"})
	s.Require().Equalf(0, code, "bypass setup: %s", out)

	// Only dash (the opener/reader), gnuhead and the head link are whitelisted. bash's absence is
	// deliberate: it links the exec -a chain but must never be trusted to open or write the file.
	config := "[watch " + mcBypassRoot + "]\nneed_encryption: false\n/usr/bin/dash\n/usr/bin/gnuhead\n" + uutilsDir + "/head\n"
	s.startDaemon(c, config)
	defer s.exec(c, []string{"sh", "-c", "pkill -TERM -f 'app-listener daemon' || true"})

	var run strings.Builder
	for i, b := range cases {
		f := fmt.Sprintf("%s/b%d.out", mcBypassRoot, i)
		// dash opens (and truncates) the file as itself — whitelisted — then exec hands fd 1 to the
		// image under test, whose own write is what the guard judges. dash emits nothing of its own.
		fmt.Fprintf(&run, "dash -c '%s >%s 2>/dev/null' </dev/null || true\n", b.run, f)
	}
	run.WriteString("echo BYPASS-DONE")
	_, out = s.exec(c, []string{"sh", "-c", run.String()})
	s.Require().Containsf(out, "BYPASS-DONE", "bypass script did not finish: %s", out)

	// Read each verdict with a whitelisted dash using only the `[ -s ]` builtin (a stat, not a read of
	// a non-whitelisted binary): HIT i means file i has bytes, i.e. its write was admitted.
	var rb strings.Builder
	rb.WriteString("dash -c '")
	for i := range cases {
		fmt.Fprintf(&rb, "[ -s %s/b%d.out ] && echo HIT%d; ", mcBypassRoot, i, i)
	}
	rb.WriteString("true'")
	_, hits := s.exec(c, []string{"sh", "-c", rb.String()})
	hit := map[string]bool{}
	for line := range strings.SplitSeq(hits, "\n") {
		hit[strings.TrimSpace(line)] = true
	}

	var wrong []string
	for i, b := range cases {
		landed := hit[fmt.Sprintf("HIT%d", i)]
		if landed != b.admit {
			verdict := "admitted (bytes landed)"
			if !landed {
				verdict = "denied (empty)"
			}
			want := "deny"
			if b.admit {
				want = "admit"
			}
			wrong = append(wrong, fmt.Sprintf("  %-32s want %s, got %s", b.label, want, verdict))
		}
	}
	s.Require().Emptyf(wrong, "%d of %d bypass cases judged wrong:\n%s\nreadback: %q\ndaemon log tail:\n%s",
		len(wrong), len(cases), strings.Join(wrong, "\n"), hits, tailLast(s.readDaemonLog(c), 40))
}

// mcLanded runs each command with fd 1 on its own file under dir (opened by a whitelisted dash, so
// only the exec'd image's write is judged) and reports which files got bytes: the kernel's verdict.
func (s *IntegrationSuite) mcLanded(c testcontainers.Container, dir string, runs ...string) []bool {
	var b strings.Builder
	for i, r := range runs {
		fmt.Fprintf(&b, "dash -c '%s >%s/v%d.out 2>/dev/null' </dev/null || true\n", r, dir, i)
	}
	b.WriteString("dash -c '")
	for i := range runs {
		fmt.Fprintf(&b, "[ -s %s/v%d.out ] && echo HIT%d; ", dir, i, i)
	}
	b.WriteString("echo LANDED-DONE'")
	_, out := s.exec(c, []string{"sh", "-c", b.String()})
	s.Require().Containsf(out, "LANDED-DONE", "verdict script did not finish: %s", out)
	got := make([]bool, len(runs))
	for line := range strings.SplitSeq(out, "\n") {
		if i, err := strconv.Atoi(strings.TrimPrefix(strings.TrimSpace(line), "HIT")); err == nil && i < len(got) {
			got[i] = true
		}
	}
	return got
}

// mcPrecreate creates dir and n empty verdict files before the daemon starts, so their inodes are
// populated into the guard (see TestDaemon_Multicall_BypassCorpus).
func (s *IntegrationSuite) mcPrecreate(c testcontainers.Container, dir string, n int) {
	code, out := s.exec(c, []string{"sh", "-c", fmt.Sprintf(
		"mkdir -p /etc/app-listener %s && for i in $(seq 0 %d); do : > %s/v$i.out; done", dir, n-1, dir)})
	s.Require().Equalf(0, code, "precreate: %s", out)
}

// A uutils upgrade replaces the multicall by a new inode at every applet link (dpkg: temp file +
// rename). The whitelisted applet is re-admitted on the new inode under the same applet identity;
// its siblings still are not.
func (s *IntegrationSuite) TestDaemon_Multicall_UpgradeKeepsAppletIdentity() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.requireVettedUutils(c)
	s.requireDistinctExes(c, uutilsMulticall, "/usr/bin/dash")

	const dir = "/mcup"
	s.mcPrecreate(c, dir, 2)
	s.startDaemon(c, "[watch "+dir+"]\nneed_encryption: false\n/usr/bin/dash\n"+uutilsDir+"/head\n")
	defer s.exec(c, []string{"sh", "-c", "pkill -TERM -f 'app-listener daemon' || true"})
	head, cat := uutilsDir+"/head --version", uutilsDir+"/cat --version"
	s.Require().Equal([]bool{true, false}, s.mcLanded(c, dir, head, cat), "baseline: head admitted, cat not")

	old := s.inodeOf(c, uutilsDir+"/head")
	upgrade := "set -e; cd " + uutilsDir + "; cp -p head .new; for f in *; do " +
		`[ "$(stat -c %i "$f")" = "` + old + `" ] || continue; ln -f .new ".tmp-$f"; mv -f ".tmp-$f" "$f"; done; rm .new`
	code, out := s.exec(c, []string{"bash", "-c", upgrade + " 2>&1"})
	s.Require().Equalf(0, code, "the upgrade: %s", out)
	s.Require().NotEqual(old, s.inodeOf(c, uutilsDir+"/head"), "fixture: the upgrade must make a new inode")
	s.Require().Equal(s.inodeOf(c, uutilsDir+"/head"), s.inodeOf(c, uutilsDir+"/cat"),
		"fixture: the applets must share the new inode")

	// Each denial schedules a re-sync; the re-admission needs no reload.
	var got []bool
	for dl := time.Now().Add(40 * time.Second); time.Now().Before(dl); time.Sleep(time.Second) {
		if got = s.mcLanded(c, dir, head, cat); got[0] {
			break
		}
	}
	s.Require().Equalf([]bool{true, false}, got, "after the upgrade: head re-admitted, cat not\ndaemon log:\n%s",
		tailLast(s.readDaemonLog(c), 40))
	s.Require().NotContains(s.readDaemonLog(c), "configuration reloaded", "re-admission must need no reload")
}

// A daemon.conf line naming a multicall whose applets can't be told apart is dropped (CRITICAL)
// and stays denied; the daemon still starts and enforces every other line. A uutils build a user
// could have placed is one: its `--list` is never run, so its applets are unknown.
func (s *IntegrationSuite) TestDaemon_Multicall_OpaqueLineDropped() {
	c := s.multicallContainer()
	s.requireDistinctExes(c, "/usr/bin/dash", "/usr/bin/gnuhead", "/usr/bin/gnucat")
	const dir, opq, userUu = "/mcopq", "/mc-fixtures/opq/wget", "/mc-fixtures/user/head"
	s.exec(c, []string{"mkdir", "-p", "/mc-fixtures/opq"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/netmc_opaque"), opq, 0o755),
		"copy netmc_opaque (run make -C integrationtests/exploits)")
	code, out := s.exec(c, []string{"sh", "-c", "set -e; mkdir -p /mc-fixtures/user; cp " + uutilsMulticall + " " +
		userUu + "; chown -R 65534:65534 /mc-fixtures/user"})
	s.Require().Equalf(0, code, "user-owned uutils copy: %s", out)
	s.mcPrecreate(c, dir, 2)

	config := "[watch " + dir + "]\nneed_encryption: false\n/usr/bin/dash\n/usr/bin/gnuhead\n" + opq + "\n" + userUu + "\n"
	s.exec(c, []string{"sh", "-c", "cat > /etc/app-listener/daemon.conf <<'EOF'\n" + config + "\nEOF"})
	code, out = s.exec(c, []string{"/app-listener", "trust-binaries", "--yes", "--no-reload",
		"/usr/bin/dash", "/usr/bin/gnuhead"})
	s.Require().Equalf(0, code, "confirming the single binaries: %s", out)
	s.launchDaemonUnconfirmed(c)
	s.awaitDaemonUp(c, config)
	defer s.exec(c, []string{"sh", "-c", "pkill -TERM -f 'app-listener daemon' || true"})

	log := s.readDaemonLog(c)
	s.Require().Containsf(log, "CRITICAL: "+dir+" whitelist line "+opq+" dropped", "daemon log:\n%s", log)
	s.Require().Contains(log, "cannot be told apart")
	s.Require().Containsf(log, "CRITICAL: "+dir+" whitelist line "+userUu+" dropped",
		"a user-owned uutils build must be opaque:\n%s", log)
	s.Require().Equal([]bool{true, false}, s.mcLanded(c, dir, "/usr/bin/gnuhead --version", "/usr/bin/gnucat --version"),
		"the other lines stay enforced: gnuhead admitted, gnucat not")
}

// GNU coreutils (debian:stable) is one binary per tool: identity stays the plain inode, no applet
// tagging engages, and whitelisting head admits head only.
func (s *IntegrationSuite) TestDaemon_Multicall_GNUControl() {
	c := s.startContainer("debian:stable", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.requireDistinctExes(c, "/usr/bin/head", "/usr/bin/cat", "/usr/bin/dash")

	const dir = "/gnu"
	s.mcPrecreate(c, dir, 2)
	s.startDaemon(c, "[watch "+dir+"]\nneed_encryption: false\n/usr/bin/dash\n/usr/bin/head\n")
	defer s.exec(c, []string{"sh", "-c", "pkill -TERM -f 'app-listener daemon' || true"})

	s.Require().Equal([]bool{true, false}, s.mcLanded(c, dir, "/usr/bin/head --version", "/usr/bin/cat --version"),
		"GNU head admitted, cat not")
	log := s.readDaemonLog(c)
	s.Require().NotContainsf(log, "multicall binary", "no multicall refusal may engage on GNU coreutils:\n%s", log)
	s.Require().NotContainsf(log, "CRITICAL", "daemon log:\n%s", log)
}

// BusyBox runs its applets in-process, so they can't be told apart: guard -w/-b on an applet refuse
// to start, naming the siblings.
func (s *IntegrationSuite) TestGuard_Multicall_BusyboxRefused() {
	c := s.startContainer("busybox:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.exec(c, []string{"mkdir", "-p", "/protected"})
	for _, flag := range []string{"-w", "-b"} {
		_, out := s.exec(c, []string{"sh", "-c",
			"timeout 20 /app-listener guard /protected " + flag + " /bin/cat --headless 2>&1; echo rc=$?"})
		s.Require().NotRegexpf(`rc=(0|124|143)$`, out, "guard %s busybox started: %s", flag, out)
		s.Require().Containsf(out, "busybox multicall binary", "guard %s: %s", flag, out)
		s.Require().Containsf(out, "the same file also runs as:", "guard %s: siblings not named: %s", flag, out)
	}
}

// A per-binary event list binds a uutils applet too: its mask row is written under the applet's
// tagged key, the one the hooks look up (a mask under the plain inode is never read, leaving the
// applet unrestricted). The same applet unrestricted in a second resource proves it is admitted.
func (s *IntegrationSuite) TestDaemon_Multicall_AppletEventMaskEnforced() {
	c := s.multicallContainer()
	s.requireDistinctExes(c, uutilsMulticall, "/usr/bin/dash", "/usr/bin/gnuhead")
	const masked, open, dd = "/mcmask", "/mcmask2", uutilsDir + "/dd"
	code, out := s.exec(c, []string{"sh", "-c", "set -e; for d in " + masked + " " + open + "; do mkdir -p $d; " +
		"echo MASK-SECRET > $d/secret; : > $d/out; done; rm -f /tmp/leak*"})
	s.Require().Equalf(0, code, "fixture: %s", out)
	s.startDaemon(c, "[watch "+masked+"]\nneed_encryption: false\n/usr/bin/gnuhead\n"+dd+" OPEN,WRITE,STAT\n\n"+
		"[watch "+open+"]\nneed_encryption: false\n"+dd+"\n")
	defer s.exec(c, []string{"sh", "-c", "pkill -TERM -f 'app-listener daemon' || true"})

	_, out = s.exec(c, []string{"sh", "-c", dd + " if=" + open + "/secret of=/tmp/leak-open 2>&1; cat /tmp/leak-open"})
	s.Require().Containsf(out, "MASK-SECRET", "baseline: unrestricted dd must read in %s: %s", open, out)

	_, out = s.exec(c, []string{"sh", "-c", dd + " if=" + masked + "/secret of=/tmp/leak 2>&1; cat /tmp/leak 2>/dev/null"})
	s.Require().NotContainsf(out, "MASK-SECRET", "dd granted OPEN,WRITE,STAT read %s/secret", masked)
	s.Require().Regexpf("Permission denied|not permitted", out, "the read must be refused by the kernel: %s", out)

	_, out = s.exec(c, []string{"sh", "-c", "printf WROTE | " + dd + " of=" + masked + "/out conv=notrunc 2>&1; " +
		"/usr/bin/gnuhead -c5 " + masked + "/out"})
	s.Require().Containsf(out, "WROTE", "dd granted WRITE must write %s/out: %s", masked, out)

	denied := false
	for _, ev := range parseDaemonEvents(s.readDaemonLog(c)) {
		denied = denied || (ev.Denied && ev.Comm == "dd" && ev.Path == masked+"/secret" && ev.Op == "READ")
	}
	s.Require().Truef(denied, "no DAEMON DENIED op=READ comm=dd for %s/secret:\n%s", masked,
		tailLast(s.readDaemonLog(c), 40))
}
