package integrationtests

import (
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"

	"github.com/Virgula0/app-listener/cmd/functions/editprotected"
)

const forwardPassword = "Sup3r-Secret-Fwd-42"

// Forward fixtures: grep is whitelisted, sed and perl are not. None are coreutils, whose applets are
// one multicall inode on ubuntu:latest (uutils): whitelisting head there admits cat and mkdir too.
const (
	fwdListed   = "/usr/bin/grep"
	fwdUnlisted = "/usr/bin/sed"
	fwdMaker    = "/usr/bin/perl"
	fwdConf     = "[watch /protected]\nneed_encryption: false\n" + fwdListed

	listedRead   = "grep -h . /protected/secret"
	unlistedRead = "sed -n p /protected/secret"
)

// forwardSetup starts a daemon guarding /protected (whitelist: fwdListed) with an edit-protected
// password, and waits for its control socket.
func (s *IntegrationSuite) forwardSetup() testcontainers.Container {
	hash, err := editprotected.Hash(forwardPassword, editprotected.OriginInstall)
	s.Require().NoError(err)

	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	s.requireDistinctExes(c, fwdListed, fwdUnlisted, fwdMaker)
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /protected /etc/app-listener && echo SECRET > /protected/secret && chmod 644 /protected/secret && " +
			"printf '%s\\n' " + shQuote(hash) + " > /etc/app-listener/edit-auth.hash && chmod 600 /etc/app-listener/edit-auth.hash"})
	s.startDaemon(c, fwdConf)
	s.Require().Truef(s.awaitLog(c, "edit-protected control socket ready", 20*time.Second),
		"control socket never came up, log:\n%s", s.readDaemonLog(c))
	return c
}

// startForward backgrounds `edit-protected --forward <args>`, its pid in /tmp/fwd.pid and output
// in /tmp/fwd.log, and waits for the n-th GRANTED line. The session idles out in 30s unless args
// set --timeout-session: a grant that should never have activated can't outlive the test by 30m.
func (s *IntegrationSuite) startForward(c testcontainers.Container, args string, n int) {
	if !strings.Contains(args, "--timeout-session") {
		args += " --timeout-session 30s"
	}
	cmd := "APP_LISTENER_EDIT_PASSWORD=" + forwardPassword +
		" nohup /app-listener edit-protected --forward --resource /protected " + args +
		" > /tmp/fwd.log 2>&1 & echo $! > /tmp/fwd.pid"
	code, out := s.exec(c, []string{"sh", "-c", cmd})
	s.Require().Equalf(0, code, "starting edit-protected --forward: %s", out)
	s.Require().Truef(s.awaitLogCount(c, "temporary access GRANTED", n, 15*time.Second),
		"grant #%d never became active, client:\n%s\ndaemon:\n%s", n, s.forwardLog(c), s.readDaemonLog(c))
}

func (s *IntegrationSuite) forwardLog(c testcontainers.Container) string {
	_, out := s.exec(c, []string{"sh", "-c", "cat /tmp/fwd.log 2>/dev/null"})
	return out
}

// stopForward signals the client and waits for the n-th REVOKED line.
func (s *IntegrationSuite) stopForward(c testcontainers.Container, sig string, n int) {
	s.exec(c, []string{"sh", "-c", "kill -" + sig + " $(cat /tmp/fwd.pid) 2>/dev/null; true"})
	s.Require().Truef(s.awaitLogCount(c, "temporary access REVOKED", n, 15*time.Second),
		"grant #%d was never revoked after SIG%s, daemon:\n%s", n, sig, s.readDaemonLog(c))
}

func (s *IntegrationSuite) awaitLogCount(c testcontainers.Container, needle string, n int, timeout time.Duration) bool {
	for dl := time.Now().Add(timeout); time.Now().Before(dl); {
		if strings.Count(s.readDaemonLog(c), needle) >= n {
			return true
		}
		time.Sleep(200 * time.Millisecond)
	}
	return false
}

// requireRead asserts cmd reads the secret (or, allowed=false, is refused by the guard).
func (s *IntegrationSuite) requireRead(c testcontainers.Container, cmd string, allowed bool, why string) {
	code, out := s.exec(c, []string{"sh", "-c", cmd + " 2>&1"})
	if allowed {
		s.Require().Equalf(0, code, "%s: %s", why, out)
		s.Require().Containsf(out, "SECRET", "%s: %s", why, out)
		return
	}
	s.Require().NotContainsf(out, "SECRET", "%s: %s", why, out)
	s.Require().Containsf(out, "Operation not permitted", "%s: expected a guard denial, got: %s", why, out)
}

// edit-protected --forward (issue #85): -w admits a binary only while the client holds the session,
// -e narrows what it may do, -b denies a whitelisted binary, and every exit path (END, client
// SIGKILL) restores the configured whitelist exactly.
func (s *IntegrationSuite) TestDaemon_EditProtected_Forward_AllowBlockRevoke() {
	c := s.forwardSetup()
	defer c.Terminate(s.ctx)

	s.requireRead(c, unlistedRead, false, "baseline: sed is not whitelisted")
	s.requireRead(c, listedRead, true, "baseline: grep is whitelisted")

	// -w with an event filter: sed reads, perl (granted only OPEN/READ/STAT) can't mkdir.
	s.startForward(c, "-w "+fwdUnlisted+" -w "+fwdMaker+" -e OPEN,READ,STAT", 1)
	s.requireRead(c, unlistedRead, true, "sed under -w")
	// The client shows the granted binaries' events, allowed ones included (guard's GUARD| format).
	s.Require().Eventuallyf(func() bool {
		return strings.Contains(s.forwardLog(c), "|sed|/protected/secret|")
	}, 10*time.Second, 200*time.Millisecond, "sed's access never reached the client:\n%s", s.forwardLog(c))
	code, out := s.exec(c, []string{"sh", "-c",
		`perl -e 'mkdir "/protected/newdir" or die "mkdir: $!\n"' 2>&1`})
	s.Require().NotEqualf(0, code, "MKDIR is outside the -e set: %s", out)
	s.Require().Containsf(out, "Operation not permitted", "perl mkdir: %s", out)
	s.stopForward(c, "TERM", 1)
	s.requireRead(c, unlistedRead, false, "sed after the client exited")

	// A killed client is an EOF: the daemon revokes on its own.
	s.startForward(c, "-w "+fwdUnlisted, 2)
	s.requireRead(c, unlistedRead, true, "sed under the second -w")
	s.stopForward(c, "KILL", 2)
	s.requireRead(c, unlistedRead, false, "sed after the client was killed")

	// -b denies a whitelisted binary for the session only.
	s.startForward(c, "-b "+fwdListed, 3)
	s.requireRead(c, listedRead, false, "grep under -b")
	s.stopForward(c, "TERM", 3)
	s.requireRead(c, listedRead, true, "grep after the -b session")

	s.exec(c, []string{"sh", "-c", "pkill -TERM -f 'app-listener daemon' || true"})
}

// Refusals: the daemon's own binary, a non-whitelisted -b, an unknown resource and a wrong
// password all leave nothing granted.
func (s *IntegrationSuite) TestDaemon_EditProtected_Forward_Refusals() {
	c := s.forwardSetup()
	defer c.Terminate(s.ctx)

	// Bounded twice: a request wrongly granted would otherwise hold the session for 30m.
	run := func(env, args string) (int, string) {
		return s.exec(c, []string{"sh", "-c", "APP_LISTENER_EDIT_PASSWORD=" + env +
			" timeout 30 /app-listener edit-protected --forward --timeout-session 10s " + args + " 2>&1"})
	}
	code, out := run(forwardPassword, "--resource /protected -w /app-listener --yes")
	s.Require().NotEqualf(0, code, "the app-listener binary must never be granted: %s", out)
	s.Require().Containsf(out, "app-listener binary itself", "%s", out)

	code, out = run(forwardPassword, "--resource /protected -b "+fwdUnlisted)
	s.Require().NotEqualf(0, code, "-b on a non-whitelisted binary changes nothing: %s", out)
	s.Require().Containsf(out, "already denied", "%s", out)

	code, out = run(forwardPassword, "--resource /etc/app-listener -w "+fwdUnlisted)
	s.Require().NotEqualf(0, code, "the self-guarded config dir is not a forwardable resource: %s", out)
	s.Require().Containsf(out, "not a guarded directory", "%s", out)

	code, out = run("wrong-Pass-1234", "--resource /protected -w "+fwdUnlisted)
	s.Require().NotEqualf(0, code, "a wrong password must be refused: %s", out)
	s.Require().Containsf(out, "authentication failed", "%s", out)
	s.Require().NotContainsf(out, "/protected ", "an unauthenticated caller must learn no resource: %s", out)

	s.Require().NotContainsf(s.readDaemonLog(c), "temporary access GRANTED", "a refused request granted something")
	s.requireRead(c, unlistedRead, false, "sed after the refusals")

	s.exec(c, []string{"sh", "-c", "pkill -TERM -f 'app-listener daemon' || true"})
}

// A granted binary its owner can rewrite in place keeps its inode, and so its rows: the daemon
// must revoke the whole grant once the content changes.
func (s *IntegrationSuite) TestDaemon_EditProtected_Forward_InPlaceRewriteRevokes() {
	c := s.forwardSetup()
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"sh", "-c", "cp " + fwdUnlisted + " /tmp/mysed && chmod 755 /tmp/mysed"})
	code, out := s.exec(c, []string{"sh", "-c", "APP_LISTENER_EDIT_PASSWORD=" + forwardPassword +
		" timeout 30 /app-listener edit-protected --forward --resource /protected -w /tmp/mysed < /dev/null 2>&1"})
	s.Require().NotEqualf(0, code, "a user-placed binary needs --yes without a terminal: %s", out)
	s.Require().Containsf(out, "--yes", "%s", out)

	s.startForward(c, "-w /tmp/mysed --yes", 1)
	s.requireRead(c, "/tmp/mysed -n p /protected/secret", true, "mysed under -w")

	s.exec(c, []string{"sh", "-c", "printf 'x' >> /tmp/mysed"})
	s.Require().Truef(s.awaitLogCount(c, "temporary access REVOKED", 1, 15*time.Second),
		"an in-place rewrite did not revoke the grant, daemon:\n%s", s.readDaemonLog(c))
	s.Require().Contains(s.readDaemonLog(c), "changed in place during its temporary grant")
	s.requireRead(c, "/tmp/mysed -n p /protected/secret", false, "mysed after its rewrite")

	// The client notices the daemon closing the session.
	s.Require().Eventuallyf(func() bool {
		return strings.Contains(s.forwardLog(c), "the daemon ended the temporary access")
	}, 10*time.Second, 200*time.Millisecond, "client log:\n%s", s.forwardLog(c))

	s.exec(c, []string{"sh", "-c", "pkill -TERM -f 'app-listener daemon' || true"})
}

// Pinned links keep enforcing after a SIGKILL, with the temporary allow still in the pinned
// whitelist: `daemon --lockdown` (ExecStopPost) must strip it from the journal, and a restart must
// come back without it.
func (s *IntegrationSuite) TestDaemon_EditProtected_Forward_DaemonKilledMidGrant() {
	c := s.forwardSetup()
	defer c.Terminate(s.ctx)

	s.startForward(c, "-w "+fwdUnlisted, 1)
	s.requireRead(c, unlistedRead, true, "sed under -w")

	// The daemon dies first: a client exiting first would be an EOF, revoked cleanly.
	s.sigDaemon(c, "KILL")
	s.exec(c, []string{"sh", "-c", "kill -KILL $(cat /tmp/fwd.pid) 2>/dev/null; true"})
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not die after SIGKILL")
	s.requireRead(c, unlistedRead, true, "the pinned whitelist still holds the allow (what lockdown closes)")

	code, out := s.runLockdown(c)
	s.Require().Equalf(0, code, "daemon --lockdown: %s", out)
	s.Require().Containsf(out, "stripped 1 temporary allow row", "lockdown did not strip the journaled allow: %s", out)
	s.requireRead(c, unlistedRead, false, "sed after lockdown, daemon dead")
	_, journal := s.exec(c, []string{"sh", "-c", "cat /etc/app-listener/temp-grants.json"})
	s.Require().Emptyf(strings.TrimSpace(journal), "lockdown left the journal behind: %s", journal)

	s.startDaemon(c, fwdConf)
	s.requireRead(c, unlistedRead, false, "sed after a restart")
	s.requireRead(c, listedRead, true, "grep after a restart")

	s.sigDaemon(c, "TERM")
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not exit after the final SIGTERM")
}

// The session ends after --timeout-session without activity; a granted binary touching the
// resource is activity and keeps it alive.
func (s *IntegrationSuite) TestDaemon_EditProtected_Forward_IdleTimeout() {
	c := s.forwardSetup()
	defer c.Terminate(s.ctx)

	s.startForward(c, "-w "+fwdUnlisted+" --timeout-session 4s", 1)
	for range 6 { // ~6s of reads, one per second: past the 4s timeout only if activity is ignored
		s.requireRead(c, unlistedRead, true, "sed while the session is active")
		time.Sleep(time.Second)
	}
	s.Require().NotContainsf(s.readDaemonLog(c), "temporary access REVOKED",
		"an active session hit the idle timeout:\n%s", s.readDaemonLog(c))

	s.Require().Truef(s.awaitLogCount(c, "temporary access REVOKED", 1, 15*time.Second),
		"an idle session was never revoked, daemon:\n%s", s.readDaemonLog(c))
	s.Require().Contains(s.readDaemonLog(c), "idle for 4s")
	s.requireRead(c, unlistedRead, false, "sed after the idle timeout")
	s.Require().Eventuallyf(func() bool {
		return strings.Contains(s.forwardLog(c), "the daemon ended the temporary access")
	}, 10*time.Second, 200*time.Millisecond, "client log:\n%s", s.forwardLog(c))

	s.exec(c, []string{"sh", "-c", "pkill -TERM -f 'app-listener daemon' || true"})
}

// edit-protected --edit-config: a valid edit is written and reloaded live; one the reload refuses
// leaves the running configuration and is rolled back on disk.
func (s *IntegrationSuite) TestDaemon_EditProtected_EditConfig() {
	c := s.forwardSetup()
	defer c.Terminate(s.ctx)

	const added = fwdConf + "\n" + fwdUnlisted + "\n"
	editConfig := func(content string) (int, string) {
		s.exec(c, []string{"sh", "-c", "printf '%s' " + shQuote(content) + " > /tmp/new.conf"})
		return s.exec(c, []string{"sh", "-c", "APP_LISTENER_EDIT_PASSWORD=" + forwardPassword +
			" timeout 60 /app-listener edit-protected --edit-config --content-file /tmp/new.conf 2>&1"})
	}

	s.requireRead(c, unlistedRead, false, "baseline: sed is not whitelisted")
	code, out := editConfig(added)
	s.Require().Equalf(0, code, "--edit-config: %s", out)
	s.Require().Contains(out, "saved and reloaded")
	_, conf := s.exec(c, []string{"cat", "/etc/app-listener/daemon.conf"})
	s.Require().Equal(strings.TrimSpace(added), conf)
	s.requireRead(c, unlistedRead, true, "sed after the live config edit")

	// Parses, but the daemon can't run it: no [watch] section.
	code, out = editConfig("# nothing to guard\n")
	s.Require().NotEqualf(0, code, "an unloadable configuration must be refused: %s", out)
	s.Require().Containsf(out, "restored", "%s", out)
	_, conf = s.exec(c, []string{"cat", "/etc/app-listener/daemon.conf"})
	s.Require().Equal(strings.TrimSpace(added), conf, "a refused edit must leave daemon.conf as it was")
	s.requireRead(c, unlistedRead, true, "sed: the running configuration is unchanged")

	code, out = s.exec(c, []string{"sh", "-c", "APP_LISTENER_EDIT_PASSWORD=wrong-Pass-1234" +
		" timeout 60 /app-listener edit-protected --edit-config --content-file /tmp/new.conf 2>&1"})
	s.Require().NotEqualf(0, code, "a wrong password must be refused: %s", out)
	s.Require().Contains(out, "authentication failed")

	s.exec(c, []string{"sh", "-c", "pkill -TERM -f 'app-listener daemon' || true"})
}

// A runtime admitted by -w is judged at launch like a whitelisted one (trust_exec_launch keys on
// guard_trusted_files' runtime bits): started with env or flags that load its caller's code it gets
// nothing, started cleanly it reads. launch_probe's marker strings class it a Node runtime.
func (s *IntegrationSuite) TestDaemon_EditProtected_Forward_RuntimeLaunchJudged() {
	c := s.forwardSetup()
	defer c.Terminate(s.ctx)

	const probe = "/opt/rt/bin/launch_probe"
	s.exec(c, []string{"mkdir", "-p", "/opt/rt/bin"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/launch_probe"), probe, 0o755),
		"copy launch_probe")
	read := func(env, flags string) string {
		_, out := s.exec(c, []string{"sh", "-c", env + " " + probe + " read /protected/secret " + flags + " 2>&1"})
		return out
	}
	s.Require().NotContains(read("", ""), "STOLEN|", "baseline: launch_probe is not whitelisted")

	s.startForward(c, "-w "+probe, 1)
	s.Require().Containsf(read("", ""), "STOLEN|SECRET", "a clean launch under -w must read:\n%s", s.readDaemonLog(c))
	for _, lc := range []struct{ what, env, flags string }{
		{"ELECTRON_RUN_AS_NODE", "ELECTRON_RUN_AS_NODE=1", ""},
		{"NODE_OPTIONS --require", "NODE_OPTIONS='--require /tmp/x.js'", ""},
		{"--inspect", "", "--inspect=9229"},
	} {
		out := read(lc.env, lc.flags)
		s.Require().NotContainsf(out, "STOLEN|", "launched with %s, the granted runtime read the secret: %s\ndaemon log:\n%s",
			lc.what, out, s.readDaemonLog(c))
		s.Require().Containsf(out, "open denied: Operation not permitted", "%s: expected the kernel's refusal: %s", lc.what, out)
	}
	s.requireDenialLogged(c, "LAUNCH", "secret")
	s.Require().Containsf(read("", ""), "STOLEN|SECRET", "the mark is the process's: a clean launch reads again")

	s.stopForward(c, "TERM", 1)
	s.Require().NotContains(read("", ""), "STOLEN|", "launch_probe after the client exited")
	s.Require().NotContains(s.readDaemonLog(c), "restoring the trust row", "the revoke left a trust row behind")

	s.exec(c, []string{"sh", "-c", "pkill -TERM -f 'app-listener daemon' || true"})
}
