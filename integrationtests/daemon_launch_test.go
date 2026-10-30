package integrationtests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/testcontainers/testcontainers-go"
)

const (
	launchProbeSrc = "/exploits/launch_probe"
	launchProbe    = "/opt/rt/bin/launch_probe"
	launchJava     = "/opt/rt/bin/java"
	launchCode     = "/opt/rt/code"
	launchExt      = launchCode + "/ext.js"
	launchSecret   = "/protected/secret"
	launchMarker   = "TOP-SECRET-LAUNCH-4F18"
)

// startLaunchDaemon whitelists launch_probe (classed a Node/Electron runtime by its marker strings)
// and its `java` copy (a JVM by name) for /protected, and makes both writers of a read-only lib_dir.
// Root-placed, so the binary ledger admits them at every start.
func (s *IntegrationSuite) startLaunchDaemon(c testcontainers.Container) {
	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/launch_probe"), launchProbeSrc, 0o755),
		"copy launch_probe")
	code, out := s.exec(c, []string{"sh", "-c", "mkdir -p /protected /etc/app-listener /opt/rt/bin " + launchCode +
		" && printf '" + launchMarker + "' > " + launchSecret + " && echo base > " + launchExt +
		" && cp " + launchProbeSrc + " " + launchProbe + " && cp " + launchProbeSrc + " " + launchJava + " 2>&1"})
	s.Require().Equalf(0, code, "setup: %s", out)
	s.startDaemon(c, `[watch /protected]
need_encryption: false
`+launchProbe+`
`+launchJava+`
lib_dir `+launchCode)
}

type launchCase struct {
	what, env, bin string
	flags          []string
}

// launch runs bin with env prefixed (shell syntax) and flags after its own arguments.
func (s *IntegrationSuite) launch(c testcontainers.Container, lc launchCase, verb, file string) string {
	cmd := strings.TrimSpace(lc.env+" "+lc.bin+" "+verb+" "+file+" "+strings.Join(lc.flags, " ")) + " 2>&1"
	_, out := s.exec(c, []string{"sh", "-c", cmd})
	return out
}

// A whitelisted runtime started with env or flags that load its caller's code is that code running
// as the whitelisted inode: it must neither read, stat nor list the secret nor change the app's code
// tree, while the same binary started cleanly, or with harmless options, does all of it.
func (s *IntegrationSuite) TestDaemon_LaunchScan_RiskyLaunchRefusedSecretAndCode() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.startLaunchDaemon(c)

	clean := []launchCase{
		{"a clean launch", "", launchProbe, nil},
		{"an allowed NODE_OPTIONS option", "NODE_OPTIONS=--max-old-space-size=4096", launchProbe, nil},
		{"an inspector port with no inspector", "", launchProbe, []string{"--inspect-port=9229"}},
		{"an underscore-spelled inspector port", "", launchProbe, []string{"--inspect_port=9229"}},
		{"a clean JVM launch", "", launchJava, nil},
		{"an allowed JAVA_TOOL_OPTIONS option", "JAVA_TOOL_OPTIONS=-Xmx512m", launchJava, nil},
	}
	assertClean := func(when string) {
		for _, lc := range clean {
			out := s.launch(c, lc, "read", launchSecret)
			s.Require().Containsf(out, "STOLEN|"+launchMarker, "%s: %s must read its resource: %s\ndaemon log:\n%s",
				when, lc.what, out, s.readDaemonLog(c))
			out = s.launch(c, lc, "write", launchExt)
			s.Require().Containsf(out, "WROTE|", "%s: %s must write its code dir: %s", when, lc.what, out)
			out = s.launch(c, lc, "stat", launchSecret)
			s.Require().Containsf(out, "STAT|"+launchSecret, "%s: %s must stat its resource: %s", when, lc.what, out)
			out = s.launch(c, lc, "ls", "/protected")
			s.Require().Containsf(out, "ENTRY|secret", "%s: %s must list its resource: %s", when, lc.what, out)
		}
	}
	assertClean("control")
	_, before := s.exec(c, []string{"cat", launchExt})

	for _, lc := range []launchCase{
		{"ELECTRON_RUN_AS_NODE", "ELECTRON_RUN_AS_NODE=1", launchProbe, nil},
		{"NODE_OPTIONS --require", "NODE_OPTIONS='--require /tmp/x.js'", launchProbe, nil},
		{"a risky NODE_OPTIONS token after an allowed one", "NODE_OPTIONS='--max-old-space-size=4096 --import=/tmp/x.mjs'",
			launchProbe, nil},
		{"NODE_PATH", "NODE_PATH=/tmp/mods", launchProbe, nil},
		{"--inspect", "", launchProbe, []string{"--inspect=9229"}},
		{"--inspect-brk", "", launchProbe, []string{"--inspect-brk"}},
		{"--remote-debugging-port", "", launchProbe, []string{"--remote-debugging-port=9222"}},
		{"--load-extension", "", launchProbe, []string{"--load-extension=/tmp/ext"}},
		// Chromium takes single-dash switches; Node reads '_' in a flag name as '-'.
		{"-remote-debugging-port", "", launchProbe, []string{"-remote-debugging-port=9222"}},
		{"-load-extension", "", launchProbe, []string{"-load-extension=/tmp/ext"}},
		{"--renderer-cmd-prefix", "", launchProbe, []string{"--renderer-cmd-prefix=/tmp/w"}},
		{"--browser-subprocess-path", "", launchProbe, []string{"--browser-subprocess-path=/tmp/w"}},
		{"--extensionDevelopmentPath", "", launchProbe, []string{"--extensionDevelopmentPath=/tmp/e"}},
		{"--env-file", "", launchProbe, []string{"--env-file=/tmp/.env"}},
		{"NODE_REPL_EXTERNAL_MODULE", "NODE_REPL_EXTERNAL_MODULE=/tmp/m.js", launchProbe, nil},
		{"--inspect_brk", "", launchProbe, []string{"--inspect_brk"}},
		{"NODE_OPTIONS --experimental_loader", "NODE_OPTIONS=--experimental_loader=/tmp/x.mjs", launchProbe, nil},
		{"JAVA_TOOL_OPTIONS -javaagent", "JAVA_TOOL_OPTIONS=-javaagent:/tmp/a.jar", launchJava, nil},
		// The JVM splits on any whitespace: a tab must not hide the agent behind an allowed option.
		{"a tab-separated -javaagent after an allowed option", "JAVA_TOOL_OPTIONS=\"$(printf '%s\\t%s' -Xmx512m -javaagent:/tmp/a.jar)\"",
			launchJava, nil},
		{"JDK_JAVA_OPTIONS", "JDK_JAVA_OPTIONS=-Djava.class.path=/tmp", launchJava, nil},
	} {
		out := s.launch(c, lc, "read", launchSecret)
		s.Require().NotContainsf(out, launchMarker, "launched with %s, the runtime read its resource: %s\ndaemon log:\n%s",
			lc.what, out, s.readDaemonLog(c))
		// Names and sizes leak too: the trust object refuses only regular-file opens.
		out = s.launch(c, lc, "stat", launchSecret)
		s.Require().NotContainsf(out, "STAT|", "launched with %s, the runtime stat'ed its resource: %s", lc.what, out)
		out = s.launch(c, lc, "ls", "/protected")
		s.Require().NotContainsf(out, "ENTRY|", "launched with %s, the runtime listed its resource: %s", lc.what, out)
		for _, f := range []string{launchExt, launchCode + "/planted.js"} {
			out = s.launch(c, lc, "write", f)
			s.Require().NotContainsf(out, "WROTE|", "launched with %s, the runtime changed its code dir (%s): %s",
				lc.what, f, out)
		}
	}
	_, after := s.exec(c, []string{"cat", launchExt})
	s.Require().Equal(before, after, "a suspect launch modified "+launchExt)
	code, _ := s.exec(c, []string{"test", "-e", launchCode + "/planted.js"})
	s.Require().NotEqual(0, code, "a suspect launch created a file in the code dir")
	s.requireDenialLogged(c, "LAUNCH", "secret")

	// The mark is the process's, never the binary's: a clean exec after them is judged afresh.
	assertClean("after the risky launches")

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

const (
	launchBunProbe = "/opt/rt/bin/bunapp"
	launchBunReal  = "/opt/rt/bin/bunreal"
	launchBunEvil  = "/tmp/evil.js"
)

// bunRealApp compiles, on the host, a `bun build --compile` app that prints STOLEN|<content> of the
// file `read <path>` names. "" when bun is not on PATH or the build fails: the probe still covers the
// rules. The baseline target fetches Bun's official runtime: a distro bun (Arch's) compiling for its
// own target embeds itself, linked to the host's ICU, which the test image lacks.
func (s *IntegrationSuite) bunRealApp() string {
	bun, err := exec.LookPath("bun")
	if err != nil {
		s.T().Log("bun not on PATH: skipping the real Bun-compiled app")
		return ""
	}
	dir := s.T().TempDir()
	src := filepath.Join(dir, "app.js")
	s.Require().NoError(os.WriteFile(src, []byte(`const [verb, p] = process.argv.slice(2);
if (verb === "read") {
  try { console.log("STOLEN|" + require("fs").readFileSync(p, "utf8")); }
  catch (e) { console.log("DENIED|" + e.code); }
}
`), 0o644))
	out := filepath.Join(dir, "bunreal")
	b, err := exec.Command(bun, "build", "--compile", "--target=bun-linux-x64-baseline", src, //nolint:gosec // test fixture
		"--outfile", out).CombinedOutput()
	if err != nil {
		s.T().Logf("bun build --compile failed (offline?), skipping the real Bun-compiled app: %v\n%s", err, b)
		return ""
	}
	return out
}

// A Bun exe (a runtime or a `bun build --compile` app: opencode, the native claude) ignores
// NODE_OPTIONS and hands argv to its script, but BUN_BE_BUN=1 turns it into the bun CLI and
// BUN_OPTIONS/the inspector variables load its caller's code into it. Whitelisted, it must read
// nothing below the secret when launched with them, and keep reading when launched cleanly.
func (s *IntegrationSuite) TestDaemon_LaunchScan_BunRuntimeEnv() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	s.exec(c, []string{"mkdir", "-p", "/exploits", "/opt/rt/bin"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/launch_probe_bun"), launchBunProbe, 0o755),
		"copy launch_probe_bun (run make -C integrationtests/exploits)")
	whitelist := launchBunProbe
	real := s.bunRealApp()
	if real != "" {
		s.Require().NoError(c.CopyFileToContainer(s.ctx, real, launchBunReal, 0o755), "copy the Bun-compiled app")
		whitelist += "\n" + launchBunReal
	}
	code, out := s.exec(c, []string{"sh", "-c", "mkdir -p /protected /etc/app-listener && printf '" + launchMarker +
		"' > " + launchSecret + " && printf '%s' 'try { console.log(\"STOLEN|\" + " +
		"require(\"fs\").readFileSync(\"" + launchSecret + "\", \"utf8\")) } catch (e) { console.log(\"DENIED|\" + e.code) }' > " +
		launchBunEvil + " 2>&1"})
	s.Require().Equalf(0, code, "setup: %s", out)
	s.startDaemon(c, "[watch /protected]\nneed_encryption: false\n"+whitelist)

	bins := []string{launchBunProbe}
	if real != "" {
		bins = append(bins, launchBunReal)
	}
	for _, bin := range bins {
		clean := []launchCase{{"a clean launch", "", bin, nil}, {"a Bun tmpdir", "BUN_TMPDIR=/tmp", bin, nil}}
		if bin == launchBunProbe {
			// The Bun class alone judges no Node knob. A real Bun image also carries Node's markers
			// (NODE_OPTIONS), so it is judged as Node too: stricter, never looser.
			clean = append(clean, launchCase{"NODE_OPTIONS", "NODE_OPTIONS='--require /tmp/x.js'", bin, nil},
				launchCase{"Node flags", "", bin, []string{"--print", "--inspect"}})
		}
		for _, lc := range clean {
			out := s.launch(c, lc, "read", launchSecret)
			s.Require().Containsf(out, "STOLEN|"+launchMarker, "%s with %s must read its resource: %s\ndaemon log:\n%s",
				bin, lc.what, out, s.readDaemonLog(c))
		}
		for _, lc := range []launchCase{
			{"BUN_BE_BUN", "BUN_BE_BUN=1", bin, nil},
			{"BUN_OPTIONS --preload", "BUN_OPTIONS='--preload " + launchBunEvil + "'", bin, nil},
			{"BUN_INSPECT", "BUN_INSPECT=ws://127.0.0.1:1/x", bin, nil},
			{"BUN_INSPECT_PRELOAD", "BUN_INSPECT_PRELOAD=" + launchBunEvil, bin, nil},
		} {
			out := s.launch(c, lc, "read", launchSecret)
			s.Require().NotContainsf(out, launchMarker, "%s launched with %s read its resource: %s\ndaemon log:\n%s",
				bin, lc.what, out, s.readDaemonLog(c))
		}
	}
	if real != "" {
		// The attack itself: BUN_BE_BUN=1 makes the whitelisted app evaluate its caller's script.
		_, out := s.exec(c, []string{"sh", "-c", "BUN_BE_BUN=1 " + launchBunReal + " " + launchBunEvil + " 2>&1"})
		s.T().Logf("BUN_BE_BUN=1 %s %s: %s", launchBunReal, launchBunEvil, out)
		s.Require().NotContainsf(out, launchMarker, "BUN_BE_BUN ran the caller's script as the whitelisted app: %s", out)
		s.Require().Containsf(out, "DENIED|", "the caller's script must have run and been refused: %s", out)
	}
	s.requireDenialLogged(c, "LAUNCH", "secret")

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// A generic Electron (default_app.asar beside it, Arch's /usr/lib/electronNN) runs whatever app its
// first positional argv names: whitelisted, it may be started only with an [electron_apps] entry root
// placed. Any other app is its caller's code and gets nothing below the secret.
func (s *IntegrationSuite) TestDaemon_LaunchScan_GenericElectronRunsOnlyListedApps() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const (
		electron = "/usr/lib/electron99/electron"
		app      = "/usr/lib/probe-app/app.asar"
		userApp  = "/tmp/u/app.asar"
	)
	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/launch_probe"), launchProbeSrc, 0o755),
		"copy launch_probe")
	code, out := s.exec(c, []string{"sh", "-c", "mkdir -p /protected /etc/app-listener /usr/lib/electron99/resources " +
		"/usr/lib/probe-app /usr/lib/other-app /tmp/u && printf '" + launchMarker + "' > " + launchSecret +
		" && cp " + launchProbeSrc + " " + electron + " && : > /usr/lib/electron99/resources/default_app.asar" +
		" && : > " + app + " && : > /usr/lib/other-app/app.asar && : > " + userApp + " 2>&1"})
	s.Require().Equalf(0, code, "setup: %s", out)
	s.startDaemon(c, `[electron_apps]
`+app+`
`+userApp+`

[watch /protected]
need_encryption: false
`+electron)
	s.Require().Contains(s.readDaemonLog(c), "electron app "+userApp+" refused",
		"an app a user could write must not be admitted")

	run := func(args, verb, file string) string {
		_, out := s.exec(c, []string{"sh", "-c", electron + " " + args + " " + verb + " " + file + " 2>&1"})
		return out
	}
	assertClean := func(when string) {
		out := run(app, "read", launchSecret)
		s.Require().Containsf(out, "STOLEN|"+launchMarker, "%s: the listed app must read its resource: %s\ndaemon log:\n%s",
			when, out, s.readDaemonLog(c))
		s.Require().Containsf(run(app, "stat", launchSecret), "STAT|", "%s: the listed app must stat its resource", when)
		s.Require().Containsf(run(app, "ls", "/protected"), "ENTRY|secret", "%s: the listed app must list its resource", when)
	}
	assertClean("control")

	// A switch is followed by the listed app, or the probe's verb would be the positional (risky anyway).
	for what, args := range map[string]string{
		"an unlisted app":                   "/tmp/evil.asar",
		"an unlisted root-placed app":       "/usr/lib/other-app/app.asar",
		"a listed app a user could write":   userApp,
		"a path that resolves past its app": app + "=/../../../tmp/evil",
		"a relative app":                    "app.asar",
		"default_app's --require":           "-r /tmp/x.js " + app,
		"default_app's REPL":                "-i " + app,
		"default_app's --app= file":         "--app=/tmp/evil.js " + app,
		"default_app's -repl":               "-repl " + app,
		"an app after a switch and a --":    "--enable-logging -- /tmp/evil.asar",
	} {
		out := run(args, "read", launchSecret)
		s.Require().NotContainsf(out, launchMarker, "started with %s, the generic Electron read its resource: %s\ndaemon log:\n%s",
			what, out, s.readDaemonLog(c))
		s.Require().NotContainsf(run(args, "stat", launchSecret), "STAT|", "started with %s, it stat'ed its resource", what)
	}
	s.requireDenialLogged(c, "LAUNCH", "secret")
	assertClean("after the risky launches")

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// Discord's updater trusts settings.json (update endpoint) and installer.db (what is installed), and
// runs the code in its app-* dirs: only Discord's binaries may write any of them or create an app-*
// dir, so neither an unrelated process nor a binary whitelisted elsewhere can steer the update.
func (s *IntegrationSuite) TestDaemon_DiscordUpdaterInputs_OnlyDiscordWrites() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const (
		settings  = discordDir + "/settings.json"
		installer = discordDir + "/installer.db"
		appDir    = discordDir + "/app-1.0.9"
		asar      = appDir + "/resources/app.asar"
		evil      = "EVIL-UPDATER-INPUT-71D0"
	)
	s.exec(c, []string{"mkdir", "-p", "/exploits"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/launch_probe"), launchProbeSrc, 0o755),
		"copy launch_probe")
	code, out := s.exec(c, []string{"sh", "-c", "mkdir -p /protected /etc/app-listener " + discordDir + "/sentry " +
		discordDir + "/0.0.1 " + appDir + "/resources && cp " + launchProbeSrc + " " + discordClient +
		" && printf '{}' > " + settings + " && printf db > " + installer + " && printf asar > " + asar + " 2>&1"})
	s.Require().Equalf(0, code, "setup: %s", out)
	s.startDaemon(c, `[watch `+discordDir+`/sentry]
need_encryption: false
`+discordClient+`
lib_dir `+appDir+`

[watch /protected]
need_encryption: false
/usr/bin/bash`)
	s.Require().Contains(s.readDaemonLog(c), "reserved glob name(s)", "trust guard #3 was not populated")

	// sh: an unrelated process; bash: whitelisted for another resource, not one of Discord's writers.
	for _, shell := range []string{"sh", "/usr/bin/bash"} {
		for _, attack := range []string{
			"printf " + evil + " >> " + settings,
			"printf " + evil + " > " + installer,
			"printf " + evil + " > /tmp/swap && mv -f /tmp/swap " + settings,
			"printf " + evil + " >> " + asar,
			"printf " + evil + " > " + appDir + "/evil.js",
			"mkdir " + discordDir + "/app-2.0.0",
		} {
			s.exec(c, []string{shell, "-c", attack + " 2>&1"})
		}
	}
	// Only what exists is printed: sentry/ is a guarded tree whose denials would name it too.
	_, out = s.exec(c, []string{"sh", "-c", "grep -l " + evil + " " + settings + " " + installer + " " + asar +
		" 2>/dev/null; ls -d " + appDir + "/evil.js " + discordDir + "/app-2.0.0 2>/dev/null; true"})
	s.Require().NotContainsf(out, discordDir+"/", "a non-Discord process wrote Discord's updater inputs or code: %s\n"+
		"daemon log:\n%s", out, s.readDaemonLog(c))
	s.requireDenialLogged(c, "PLANT", "settings.json")

	// Control: Discord itself (launch_probe as the client: dash has no mkdir builtin) updates all of
	// them.
	for _, step := range [][]string{
		{"write", settings}, {"write", installer}, {"write", asar},
		{"mkdir", discordDir + "/app-1.0.10"}, {"write", discordDir + "/app-1.0.10/x"},
	} {
		code, out = s.exec(c, append([]string{discordClient}, step...))
		s.Require().Equalf(0, code, "Discord must update its own inputs and code (%v): %s\ndaemon log:\n%s", step,
			out, s.readDaemonLog(c))
	}

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}

// A whitelisted CLI started with words one of its launch rules names (built in by basename, or a
// [launch_rules] line) is refused its secret; started otherwise it keeps reading. A warn rule only
// logs. Every binary is judged on SSL_CERT_FILE/SSL_CERT_DIR unless the value is a root-placed bundle
// (built in or allow_env); an allow_env a user can fill is refused. Rule sets must survive a reload.
func (s *IntegrationSuite) TestDaemon_LaunchRules() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	const gh, tool, goread = "/opt/rt/bin/gh", "/opt/rt/bin/tool", "/opt/rt/bin/goread"
	const corpCA = "/etc/corp-ca.pem"
	s.exec(c, []string{"mkdir", "-p", "/exploits", "/opt/rt/bin"})
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/launch_probe_plain"), gh, 0o755),
		"copy launch_probe_plain (run make -C integrationtests/exploits)")
	s.Require().NoError(c.CopyFileToContainer(s.ctx, netTesterAmd64Bin, goread, 0o755))
	code, out := s.exec(c, []string{"sh", "-c", "mkdir -p /protected /etc/app-listener && printf '" + launchMarker +
		"' > " + launchSecret + " && cp " + gh + " " + tool + " && echo ca > " + corpCA + " 2>&1"})
	s.Require().Equalf(0, code, "setup: %s", out)
	s.startDaemon(c, "[launch_rules]\ndeny_args "+tool+" export --all\nwarn_env "+tool+" NOTE_ME\n"+
		"allow_env SSL_CERT_FILE="+corpCA+"\nallow_env SSL_CERT_FILE=/tmp/user-ca.pem\n\n"+
		"[watch /protected]\nneed_encryption: false\n"+gh+"\n"+tool+"\n"+goread)

	run := func(cmd string) string {
		_, out := s.exec(c, []string{"sh", "-c", cmd + " 2>&1"})
		return out
	}
	check := func(when string) {
		for what, cmd := range map[string]string{
			"a clean launch":            gh + " read " + launchSecret,
			"a partial rule":            gh + " auth status read " + launchSecret,
			"a warn rule":               gh + " auth git-credential read " + launchSecret,
			"rule words out of order":   gh + " token auth read " + launchSecret,
			"another binary's built-in": tool + " auth token read " + launchSecret,
			"a warn env rule":           "NOTE_ME=1 " + tool + " read " + launchSecret,
			"a clean Go launch":         goread + " read " + launchSecret,
			"the system CA bundle":      "SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt " + gh + " read " + launchSecret,
			"the system CA dir on Go":   "SSL_CERT_DIR=/etc/ssl/certs " + goread + " read " + launchSecret,
			"a root-placed allow_env":   "SSL_CERT_FILE=" + corpCA + " " + goread + " read " + launchSecret,
		} {
			s.Require().Containsf(run(cmd), "STOLEN|"+launchMarker, "%s: %s must read its resource\ndaemon log:\n%s",
				when, what, s.readDaemonLog(c))
		}
		for what, cmd := range map[string]string{
			"words around a flag's value":  gh + " auth --hostname x token read " + launchSecret,
			"a short flag after operands":  gh + " auth status read " + launchSecret + " -t",
			"a long flag with a value":     gh + " auth status --show-token=1 read " + launchSecret,
			"a [launch_rules] argv rule":   tool + " export read " + launchSecret + " --all",
			"SSL_CERT_FILE on a Go binary": "SSL_CERT_FILE=/tmp/ca.pem " + goread + " read " + launchSecret,
			"SSL_CERT_FILE on any CLI":     "SSL_CERT_FILE=/tmp/ca.pem " + gh + " read " + launchSecret,
			"SSL_CERT_DIR of a user":       "SSL_CERT_DIR=/tmp " + goread + " read " + launchSecret,
			"a user-fillable allow_env":    "SSL_CERT_FILE=/tmp/user-ca.pem " + gh + " read " + launchSecret,
			"an allowed path extended":     "SSL_CERT_FILE=" + corpCA + ".x " + gh + " read " + launchSecret,
			"an allowed bundle in a list":  "SSL_CERT_DIR=/etc/ssl/certs:/tmp " + goread + " read " + launchSecret,
		} {
			s.Require().NotContainsf(run(cmd), launchMarker, "%s: %s read the secret\ndaemon log:\n%s",
				when, what, s.readDaemonLog(c))
		}
	}
	check("at start")
	s.reloadAndAwait(c)
	check("after a reload")

	log := s.readDaemonLog(c)
	for _, want := range []string{"word=token", "word=-t", "word=--show-token", "word=--all", "word=SSL_CERT_FILE",
		"ARGV comm=gh", "word=git-credential", "word=NOTE_ME", "word=SSL_CERT_DIR",
		"allow_env SSL_CERT_FILE=/tmp/user-ca.pem refused"} {
		s.Require().Containsf(log, want, "daemon log lacks %q:\n%s", want, log)
	}
	s.Require().NotContainsf(log, "word=SSL_CERT_FILE=", "an env match must log the name alone:\n%s", log)
	s.requireDenialLogged(c, "ARGV", "secret")

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}
