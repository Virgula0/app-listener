package integrationtests

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

// netIntDir holds exploits/netinject twice: probe (whitelisted) and inject (not). evil.so is
// lib_probe.so in /tmp, which a root exe never vouches for (other-writable directory).
const (
	netIntDir   = "/netint"
	netIntProbe = netIntDir + "/probe"
	netIntInj   = netIntDir + "/inject"
	netIntEvil  = "/tmp/evil.so"
	// A closed local port: rc=EPERM is the guard's refusal, rc=ECONNREFUSED reached the stack.
	netIntDst = " 127.0.0.1 1"
)

// netIntegrityContainer is the pool for this file: drops probe processes, the fixtures and an ext4
// loop image a failed reuse test left mounted.
func (s *IntegrationSuite) netIntegrityContainer() testcontainers.Container {
	c := s.acquirePool("netintegrity")
	s.exec(c, []string{"sh", "-c", "pkill -f '" + netIntDir + "/[p]robe' ; pkill -f '/mnt/img/[b]in' ; " +
		"mountpoint -q /mnt/img && umount /mnt/img ; losetup -D 2>/dev/null ; " +
		"rm -rf " + netIntDir + " /mnt/img /tmp/img /tmp/go /tmp/pre-*.out " + netIntEvil + " ; true"})
	s.exec(c, []string{"sh", "-c", "mkdir -p " + netIntDir + " && chmod 755 " + netIntDir})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/netinject"), netIntProbe, 0o755),
		"copy netinject (run make -C integrationtests/exploits)")
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/netinject"), netIntInj, 0o755),
		"copy netinject")
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/lib_probe.so"), netIntEvil, 0o644),
		"copy lib_probe.so")
	return c
}

func (s *IntegrationSuite) netIntRun(c testcontainers.Container, cmd string) string {
	_, out := s.exec(c, []string{"sh", "-c", cmd + " 2>&1"})
	return out
}

// awaitNetGuardLine polls the guard log for a line containing every needle.
func (s *IntegrationSuite) awaitNetGuardLine(c testcontainers.Container, timeout time.Duration, needles ...string) {
	for dl := time.Now().Add(timeout); time.Now().Before(dl); time.Sleep(250 * time.Millisecond) {
		for _, line := range strings.Split(s.readNetGuardLog(c), "\n") {
			hit := true
			for _, n := range needles {
				hit = hit && strings.Contains(line, n)
			}
			if hit {
				return
			}
		}
	}
	s.T().Fatalf("no network-guard log line with %q; log tail:\n%s", needles, netGuardTail(s.readNetGuardLog(c)))
}

// TestNetworkGuard_Integrity_InjectedCode: a whitelisted binary running code its caller chose
// (LD_PRELOAD, LD_AUDIT, a dlopen of an unvouched library) loses its network access, as does a
// forked child sharing that code; a system library and a clean re-exec keep it.
func (s *IntegrationSuite) TestNetworkGuard_Integrity_InjectedCode() {
	c := s.netIntegrityContainer()
	s.startNetworkGuardStd(c, "-w", guardBinaryFlag(netIntProbe))
	defer s.stopNetGuard(c)

	cases := []struct {
		name, cmd string
		want      []string
	}{
		{"whitelisted, clean", netIntProbe + " connect" + netIntDst, []string{"rc=ECONNREFUSED"}},
		{"not whitelisted", netIntInj + " connect" + netIntDst, []string{"rc=EPERM"}},
		{"LD_PRELOAD", "LD_PRELOAD=" + netIntEvil + " " + netIntProbe + " connect" + netIntDst,
			[]string{"LIB_PROBE_LOADED", "rc=EPERM"}},
		{"LD_PRELOAD, forked child", "LD_PRELOAD=" + netIntEvil + " " + netIntProbe + " connect" + netIntDst + " fork",
			[]string{"LIB_PROBE_LOADED", "rc=EPERM"}},
		{"LD_AUDIT", "LD_AUDIT=" + netIntEvil + " " + netIntProbe + " connect" + netIntDst, []string{"rc=EPERM"}},
		{"dlopen of an unvouched library", netIntProbe + " dlopen " + netIntEvil + netIntDst,
			[]string{"LIB_PROBE_LOADED", "dlopen=OK", "rc=EPERM"}},
		{"dlopen of a system library", netIntProbe + " dlopen libm.so.6" + netIntDst,
			[]string{"dlopen=OK", "rc=ECONNREFUSED"}},
		// The exec drops the preloaded image and its environment: the new image is judged afresh.
		{"clean exec from a preloaded image", "LD_PRELOAD=" + netIntEvil + " " + netIntProbe + " exec " +
			netIntProbe + " connect" + netIntDst, []string{"LIB_PROBE_LOADED", "rc=ECONNREFUSED"}},
	}
	for _, tc := range cases {
		out := s.netIntRun(c, tc.cmd)
		for _, w := range tc.want {
			s.Require().Containsf(out, w, "%s (%s)", tc.name, tc.cmd)
		}
	}
	s.awaitNetGuardLine(c, 8*time.Second, "NETGUARD|CONNECT|probe|", "|true|ld-preload")
	s.awaitNetGuardLine(c, 8*time.Second, "NETGUARD|CONNECT|probe|", "|true|code-suspect")
}

var netGateRe = regexp.MustCompile(`(attach|vmwrite|mem)=(\S+)`)

// TestNetworkGuard_Integrity_Ptrace: a process the whitelist doesn't admit can neither attach to,
// write into nor open the memory of a whitelisted process it spawned, nor exec one under its
// trace; the same probes against a non-whitelisted copy go through (the refusal is the guard's).
func (s *IntegrationSuite) TestNetworkGuard_Integrity_Ptrace() {
	c := s.netIntegrityContainer()
	s.startNetworkGuardStd(c, "-w", guardBinaryFlag(netIntProbe))
	defer s.stopNetGuard(c)

	gates := func(target string) map[string]string {
		out := s.netIntRun(c, netIntInj+" attachchild "+target+" wait /tmp/never"+netIntDst)
		got := map[string]string{}
		for _, m := range netGateRe.FindAllStringSubmatch(out, -1) {
			got[m[1]] = m[2]
		}
		s.Require().Lenf(got, 3, "attachchild %s: %s", target, out)
		return got
	}
	for gate, v := range gates(netIntProbe) {
		s.Require().Truef(v == "EPERM" || v == "EACCES", "%s on a whitelisted process: %s", gate, v)
	}
	ctl := gates(netIntInj)
	s.Require().Equal("OK", ctl["attach"], "control: attach to a non-whitelisted process")
	s.Require().Equal("EFAULT", ctl["vmwrite"], "control: process_vm_writev reached the target's memory")
	s.Require().Equal("OK", ctl["mem"], "control: /proc/<pid>/mem")

	s.Require().Contains(s.netIntRun(c, netIntInj+" traceme "+netIntProbe+" connect"+netIntDst), "exec=EPERM",
		"traced exec of a whitelisted binary")
	s.Require().Contains(s.netIntRun(c, netIntInj+" traceme "+netIntInj+" connect"+netIntDst), "exec=OK",
		"control: traced exec of a non-whitelisted binary")

	s.awaitNetGuardLine(c, 8*time.Second, "NETGUARD|PTRACE|inject|", "|true|target=")
	s.awaitNetGuardLine(c, 8*time.Second, "NETGUARD|TRACED_EXEC|inject|", "|true|tracer=")
}

// TestNetworkGuard_Integrity_PreexistingProcesses: a whitelisted process started before the guard,
// with LD_PRELOAD, is judged from /proc at start and refused; a clean one keeps its access.
func (s *IntegrationSuite) TestNetworkGuard_Integrity_PreexistingProcesses() {
	c := s.netIntegrityContainer()
	start := func(env, out string) {
		s.exec(c, []string{"sh", "-c", fmt.Sprintf("%s nohup %s wait /tmp/go%s > %s 2>&1 &",
			env, netIntProbe, netIntDst, out)})
	}
	start("LD_PRELOAD="+netIntEvil, "/tmp/pre-preload.out")
	start("", "/tmp/pre-clean.out")
	// Both must be running (the preload's constructor ran) before the guard starts.
	code := 1
	for dl := time.Now().Add(10 * time.Second); code != 0 && time.Now().Before(dl); time.Sleep(100 * time.Millisecond) {
		code, _ = s.exec(c, []string{"sh", "-c", "grep -q LIB_PROBE_LOADED /tmp/pre-preload.out && " +
			"[ \"$(pgrep -fc '" + netIntDir + "/[p]robe wait')\" = 2 ]"})
	}
	s.Require().Zero(code, "the pre-existing probes did not start")

	s.startNetworkGuardStd(c, "-w", guardBinaryFlag(netIntProbe))
	defer s.stopNetGuard(c)
	s.Require().Contains(s.readNetGuardLog(c), "started before the guard and is ld-preload")

	s.exec(c, []string{"touch", "/tmp/go"})
	for out, want := range map[string]string{"/tmp/pre-preload.out": "rc=EPERM", "/tmp/pre-clean.out": "rc=ECONNREFUSED"} {
		deadline := time.Now().Add(10 * time.Second)
		got := ""
		for ; time.Now().Before(deadline); time.Sleep(200 * time.Millisecond) {
			_, got = s.exec(c, []string{"cat", out})
			if strings.Contains(got, "rc=") {
				break
			}
		}
		s.Require().Containsf(got, want, "%s", out)
	}
}

// TestNetworkGuard_Bypass_ReusedInodeNumberInheritsWhitelist: ext4 hands a freed whitelisted inode's
// number to the next file created, any user's; that file must not inherit the network whitelist,
// whether the binary was deleted or renamed over.
func (s *IntegrationSuite) TestNetworkGuard_Bypass_ReusedInodeNumberInheritsWhitelist() {
	c := s.netIntegrityContainer()
	build := `
set -e
mkdir -p /mnt/img
dd if=/dev/zero of=/tmp/img bs=1M count=16 status=none
mkfs.ext4 -q -F /tmp/img
for i in $(seq 0 15); do [ -e /dev/loop$i ] || mknod /dev/loop$i b 7 "$i"; done
mount -o loop /tmp/img /mnt/img
mkdir -p /mnt/img/bin /mnt/img/drop && chmod 1777 /mnt/img/drop
echo MOUNTED`
	code, out := s.exec(c, []string{"sh", "-c", build})
	if code != 0 || !strings.Contains(out, "MOUNTED") {
		s.T().Skipf("ext4 loop mount unavailable here (exit %d): %s", code, out)
	}
	defer s.exec(c, []string{"sh", "-c", "umount /mnt/img 2>/dev/null; losetup -D 2>/dev/null; true"})

	for _, uninstall := range []string{
		"rm -f /mnt/img/bin/app",
		"cp " + netIntInj + " /mnt/img/bin/.app.new && mv -f /mnt/img/bin/.app.new /mnt/img/bin/app",
	} {
		s.exec(c, []string{"sh", "-c", "rm -f /mnt/img/bin/* /mnt/img/drop/* && cp " + netIntInj + " /mnt/img/bin/app"})
		s.startNetworkGuardStd(c, "-w", guardBinaryFlag("/mnt/img/bin/app"))
		s.Require().Contains(s.netIntRun(c, "/mnt/img/bin/app connect"+netIntDst), "rc=ECONNREFUSED",
			"baseline: the whitelisted binary connects")

		old := s.inodeOf(c, "/mnt/img/bin/app")
		code, out = s.exec(c, []string{"sh", "-c", uninstall + " && " + nobodyRun + "cp " + netIntInj +
			" /mnt/img/drop/evil && stat -c 'ino=%i' /mnt/img/drop/evil"})
		s.Require().Equalf(0, code, "%s + an unprivileged copy: %s", uninstall, out)
		m := regexp.MustCompile(`ino=(\d+)`).FindStringSubmatch(out)
		s.Require().NotNil(m, out)
		s.Require().Equalf(old, m[1], "fixture: ext4 must hand the freed number to the next file (%s)", uninstall)

		s.Require().Containsf(s.netIntRun(c, nobodyRun+"/mnt/img/drop/evil connect"+netIntDst), "rc=EPERM",
			"after %q, a file reusing the whitelisted inode number connected; log tail:\n%s",
			uninstall, netGuardTail(s.readNetGuardLog(c)))
		s.stopNetGuard(c)
	}
}
