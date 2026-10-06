package integrationtests

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/testcontainers/testcontainers-go"
)

const (
	uutilsDir       = "/usr/lib/cargo/bin/coreutils"
	uutilsMulticall = "/usr/bin/coreutils"
	uutilsVersion   = "coreutils 0.8.0 (multi-call binary)"
	mcRoot          = "/mc"
	mcProbeTable    = "/mc-probes.tsv"
	// mcBits resources: each probe's identity is the set of resources admitting it, a 7-bit code.
	mcBits = 7
)

// uutilsApplets is `coreutils --list` of Ubuntu 26.04's rust-coreutils 0.8.0: the build's applet set,
// which differs from the link names in uutilsDir (uutilsStrayLinks; hostname, kill, more and uptime
// have no link).
var uutilsApplets = strings.Fields(`[ arch b2sum base32 base64 basename basenc cat chcon chgrp chmod
	chown chroot cksum comm cp csplit cut date dd df dir dircolors dirname du echo env expand expr factor
	false fmt fold groups head hostid hostname id install join kill link ln logname ls md5sum mkdir mkfifo
	mknod mktemp more mv nice nl nohup nproc numfmt od paste pathchk pinky pr printenv printf ptx pwd
	readlink realpath rm rmdir runcon seq sha1sum sha224sum sha256sum sha384sum sha512sum shred shuf sleep
	sort split stat stdbuf stty sum sync tac tail tee test timeout touch tr true truncate tsort tty uname
	unexpand uniq unlink uptime users vdir wc who whoami yes`)

// uutilsStrayLinks are links in uutilsDir naming no applet: uutils runs the longest applet name the
// link ends with (sha3sum, hashsum and shake128sum run sum) or exits "unknown program" (relpath).
var uutilsStrayLinks = []string{"hashsum", "relpath", "sha3-224sum", "sha3-256sum", "sha3-384sum",
	"sha3-512sum", "sha3sum", "shake128sum", "shake256sum"}

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

// uutilsLinks is the pinned listing of uutilsDir.
func uutilsLinks() []string {
	var links []string
	for _, a := range uutilsApplets {
		if !slices.Contains(uutilsUnlinked, a) {
			links = append(links, a)
		}
	}
	links = append(links, uutilsStrayLinks...)
	slices.Sort(links)
	return links
}

// requireVettedUutils pins the container's uutils build: a different applet set must be re-vetted,
// not silently tested against the wrong table.
func (s *IntegrationSuite) requireVettedUutils(c testcontainers.Container) {
	_, out := s.exec(c, []string{uutilsMulticall, "--version"})
	s.Require().Containsf(out, uutilsVersion, "ubuntu:latest no longer ships the vetted uutils build")
	_, out = s.exec(c, []string{uutilsMulticall, "--list"})
	s.Require().Equal(uutilsApplets, strings.Fields(out), "uutils applet set drifted: re-vet the build")
	_, out = s.exec(c, []string{"sh", "-c", "ls -A " + uutilsDir})
	got := strings.Fields(out)
	slices.Sort(got)
	s.Require().Equal(uutilsLinks(), got, "uutils link names drifted: re-vet the build")
	// One inode behind every name, or the test proves nothing about multicall identity.
	_, out = s.exec(c, []string{"sh", "-c", "stat -c %d:%i " + uutilsMulticall + " " + uutilsDir + "/*  | sort -u"})
	s.Require().Lenf(strings.Fields(out), 1, "uutils links are not one multicall inode: %s", out)
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
	links := uutilsLinks()

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
	for _, a := range uutilsStrayLinks {
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
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
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
