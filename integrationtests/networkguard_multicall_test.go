package integrationtests

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

// netMcDir holds exploits/netmc (a synthetic uutils multicall, applets neta/netb) under three hard
// links (neta, netb, netz = a non-applet name) and netmc_opaque as opq/neta. Each probe connects to
// a closed port: "rc=EPERM" is the network guard's refusal, "rc=ECONNREFUSED" means it let the
// connect reach the stack.
const netMcDir = "/mcnet"

func (s *IntegrationSuite) netMulticallContainer() testcontainers.Container {
	c := s.newNetGuardContainer(nil)
	s.installNetMc(c)
	return c
}

func (s *IntegrationSuite) installNetMc(c testcontainers.Container) {
	s.exec(c, []string{"sh", "-c", "rm -rf " + netMcDir + " && mkdir -p " + netMcDir + "/opq && chmod 755 " +
		netMcDir + " " + netMcDir + "/opq"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/netmc"), netMcDir+"/neta", 0o755),
		"copy netmc (run make -C integrationtests/exploits)")
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/netmc_opaque"), netMcDir+"/opq/neta", 0o755),
		"copy netmc_opaque")
	code, out := s.exec(c, []string{"sh", "-c", fmt.Sprintf(
		"cd %s && ln -f neta netb && ln -f neta netz && ln -f opq/neta opq/netb && ./neta --list", netMcDir)})
	s.Require().Equalf(0, code, "netmc fixture: %s", out)
	s.Require().Equal("neta\nnetb", out, "netmc --list")
}

// netMcProbe runs one connect probe and returns its rc= verdict line.
func (s *IntegrationSuite) netMcProbe(c testcontainers.Container, cmd string) string {
	_, out := s.exec(c, []string{"sh", "-c", cmd + " 2>&1"})
	return out
}

// TestNetworkGuard_Multicall_PerApplet: whitelisting one applet of a uutils multicall admits that
// applet only, as attested by the kernel at exec, not every program sharing the inode.
func (s *IntegrationSuite) TestNetworkGuard_Multicall_PerApplet() {
	c := s.netMulticallContainer()
	s.startNetworkGuardStd(c, "-w", guardBinaryFlag(netMcDir+"/neta"))
	defer s.stopNetGuard(c)

	cases := []struct {
		name, cmd, want string
	}{
		{"whitelisted applet", netMcDir + "/neta 127.0.0.1 1", "rc=ECONNREFUSED"},
		{"whitelisted applet, forked child", netMcDir + "/neta 127.0.0.1 1 fork", "rc=ECONNREFUSED"},
		{"sibling applet, same inode", netMcDir + "/netb 127.0.0.1 1", "rc=EPERM"},
		{"sibling forked child", netMcDir + "/netb 127.0.0.1 1 fork", "rc=EPERM"},
		// argv[0] names the whitelisted applet, but the kernel's filename does not: no tag.
		{"exec -a neta on netb", "bash -c 'exec -a neta " + netMcDir + "/netb 127.0.0.1 1'", "rc=EPERM"},
		{"exec -a netb on neta", "bash -c 'exec -a netb " + netMcDir + "/neta 127.0.0.1 1'", "rc=EPERM"},
	}
	for _, tc := range cases {
		s.Require().Containsf(s.netMcProbe(c, tc.cmd), tc.want, "%s (%s)", tc.name, tc.cmd)
	}
	s.waitForNetGuardBlockedEvent(c, "netb", "CONNECT", 8*time.Second)
}

// TestNetworkGuard_Multicall_Blacklist: a blacklisted applet blocks only itself.
func (s *IntegrationSuite) TestNetworkGuard_Multicall_Blacklist() {
	c := s.netMulticallContainer()
	s.startNetworkGuardStd(c, "-b", netMcDir+"/neta")
	defer s.stopNetGuard(c)

	s.Require().Contains(s.netMcProbe(c, netMcDir+"/neta 127.0.0.1 1"), "rc=EPERM", "blacklisted applet")
	s.Require().Contains(s.netMcProbe(c, netMcDir+"/neta 127.0.0.1 1 fork"), "rc=EPERM", "blacklisted forked child")
	s.Require().Contains(s.netMcProbe(c, netMcDir+"/netb 127.0.0.1 1"), "rc=ECONNREFUSED", "sibling applet")
	// Unattested (argv[0] != filename): the multicall still runs neta from argv[0], so it must block.
	s.Require().Contains(s.netMcProbe(c, "bash -c 'exec -a neta "+netMcDir+"/netb 127.0.0.1 1'"), "rc=EPERM",
		"unattested exec running the blacklisted applet")
	s.waitForNetGuardBlockedEvent(c, "neta", "CONNECT", 8*time.Second)
}

// TestNetworkGuard_Multicall_Refusals: a multicall entry with no applet identity (a non-applet link,
// or a build whose applets can't be listed) refuses to start in either mode, naming its siblings.
func (s *IntegrationSuite) TestNetworkGuard_Multicall_Refusals() {
	c := s.netMulticallContainer()
	cases := []struct {
		name, flags, want string
	}{
		{"non-applet name, whitelist", "-w " + netMcDir + "/netz", `"netz" is not one of its applets`},
		{"non-applet name, blacklist", "-b " + netMcDir + "/netz", `"netz" is not one of its applets`},
		{"opaque build, whitelist", "-w " + netMcDir + "/opq/neta", "cannot be told apart"},
		{"opaque build, blacklist", "-b " + netMcDir + "/opq/neta", "cannot be told apart"},
	}
	for _, tc := range cases {
		_, out := s.exec(c, []string{"sh", "-c",
			"timeout 20 /app-listener network-guard " + tc.flags + " --headless 2>&1; echo rc=$?"})
		s.Require().NotContainsf(out, "rc=0", "%s: network-guard started: %s", tc.name, out)
		s.Require().NotContainsf(out, "rc=124", "%s: network-guard ran until the timeout: %s", tc.name, out)
		s.Require().Containsf(out, tc.want, "%s: %s", tc.name, out)
		s.Require().Containsf(out, "the same file also runs as:", "%s: siblings not named: %s", tc.name, out)
	}
	s.stopNetGuard(c)
}

// TestNetworkMonitor_Multicall_PerApplet: watching one applet of a uutils multicall reports that
// applet only; a sibling applet that renames its comm to the watched one (the monitor's first filter)
// is not reported as it.
func (s *IntegrationSuite) TestNetworkMonitor_Multicall_PerApplet() {
	c := s.netmonContainer()
	s.installNetMc(c)
	s.startNetworkMonitorStd(c, netMcDir+"/neta")
	defer s.stopNetMonitor(c)
	logBefore := s.readNetMonitorLog(c)

	s.Require().Contains(s.netMcProbe(c, netMcDir+"/netb 127.0.0.1 40002 comm=neta"), "rc=ECONNREFUSED")
	s.Require().Contains(s.netMcProbe(c, netMcDir+"/neta 127.0.0.1 40003"), "rc=ECONNREFUSED")

	// Ring-buffer order: once the watched applet's connect is in, the sibling's would be too.
	deadline := time.Now().Add(8 * time.Second)
	var events []netMonitorEvent
	for {
		events = netMonitorDeltaEvents(logBefore, s.readNetMonitorLog(c))
		if slices.ContainsFunc(events, func(e netMonitorEvent) bool {
			return e.Type == "CONNECT" && strings.HasSuffix(e.DstAddr, ":40003")
		}) {
			break
		}
		if time.Now().After(deadline) {
			s.T().Fatalf("watched applet's connect never reported; captured: %+v", events)
		}
		time.Sleep(250 * time.Millisecond)
	}
	for _, e := range events {
		s.Require().Falsef(strings.HasSuffix(e.DstAddr, ":40002"),
			"sibling applet posing as the watched one was reported: %+v", e)
	}
}

// TestNetworkGuard_Multicall_GNUCoreutilsSingle: Fedora's coreutils-single runs every applet as one
// inode and prints no "multi-call"; whitelisting or blacklisting it, or an applet of it, is refused.
// The image ships split coreutils: the fixture swaps in coreutils-single, whose applets are
// `#!/usr/bin/coreutils --coreutils-prog-shebang=<applet>` scripts.
func (s *IntegrationSuite) TestNetworkGuard_Multicall_GNUCoreutilsSingle() {
	c := s.startContainer("fedora:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	code, out := s.exec(c, []string{"sh", "-c", "rpm -q coreutils-single >/dev/null || " +
		"timeout 300 dnf -y -q swap coreutils coreutils-single >/dev/null 2>&1; " +
		"rpm -q coreutils-single && /usr/bin/coreutils --coreutils-prog=true && " +
		"head -n1 /usr/bin/cat | grep -qx '#!/usr/bin/coreutils --coreutils-prog-shebang=cat'"})
	s.Require().Equalf(0, code, "fixture: coreutils-single with shebang applets: %s", out)

	for _, flags := range []string{"-w /usr/bin/coreutils", "-w /usr/bin/cat", "-b /usr/bin/cat"} {
		_, out := s.exec(c, []string{"sh", "-c",
			"timeout 20 /app-listener network-guard " + flags + " --headless 2>&1; echo rc=$?"})
		s.Require().NotContainsf(out, "rc=0", "%s: network-guard started: %s", flags, out)
		s.Require().NotContainsf(out, "rc=124", "%s: network-guard ran until the timeout: %s", flags, out)
		s.Require().Containsf(out, "is a gnu-coreutils multicall binary", "%s: %s", flags, out)
	}
}
