package integrationtests

import (
	"regexp"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

const (
	inspectorProbe = "/opt/insp/inspector_probe"
	plainProbe     = "/opt/plain/inspector_probe"
	tmpProbe       = "/tmp/inspector_probe"
	// r1Probe is whitelisted by /r1 only: not in GLOBAL's (intersection) whitelist.
	r1Probe = "/opt/r1/inspector_probe"
)

const inspectorConfig = `[watch /r1]
need_encryption: false
/opt/app/app
` + r1Probe + `

[watch /r3]
need_encryption: false
/opt/y/y

[inspectors]
` + inspectorProbe + `
` + tmpProbe

// startInspectorVictims starts a process tainted by /r1 alone and one tainted GLOBAL (/r3's content,
// then /r1's whitelisted image), each idling in a builtin so its image stays put.
func (s *IntegrationSuite) startInspectorVictims(c testcontainers.Container) (single, global string) {
	s.exec(c, []string{"sh", "-c",
		"(/opt/app/app -c 'read a < /r1/secret; read y < /tmp/f; exit 0' &);" +
			" (/opt/y/y -c 'read z < /r3/secret; exec /opt/app/app -c \"read y < /tmp/g; exit 0\"' &); sleep 1"})
	_, out := s.exec(c, []string{"sh", "-c", "pgrep -f 'app -c read a' | head -1"})
	single = strings.TrimSpace(out)
	_, out = s.exec(c, []string{"sh", "-c", "pgrep -f 'app -c read y < /tmp/g' | head -1"})
	global = strings.TrimSpace(out)
	s.Require().NotEmpty(single, "the /r1 victim did not start")
	s.Require().NotEmpty(global, "the GLOBAL victim did not start")
	return single, global
}

func (s *IntegrationSuite) probe(c testcontainers.Container, args ...string) string {
	_, out := s.exec(c, []string{"sh", "-c", strings.Join(args, " ") + " 2>&1"})
	return out
}

// An [inspectors] binary reads a tainted process's /proc metadata whatever slot holds it (one
// resource or GLOBAL) — what xdg-desktop-portal needs for screen sharing — and nothing more: memory
// (process_vm_readv, ptrace attach, /proc/<pid>/mem) stays refused to the same inode.
func (s *IntegrationSuite) TestDaemon_Inspector_ReadsMetadataNeverMemory() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.setupInspectorContainer(c)
	s.startDaemon(c, inspectorConfig)
	single, global := s.startInspectorVictims(c)

	for _, pid := range []string{single, global} {
		out := s.probe(c, inspectorProbe, "all", pid)
		for _, want := range []string{"EXE_OK", "ROOT_OK", "VMREAD_DENIED", "ATTACH_DENIED", "MEM_DENIED"} {
			s.Require().Containsf(out, want, "inspector against %s: %s", pid, out)
		}
		out = s.probe(c, plainProbe, "all", pid)
		for _, want := range []string{"EXE_DENIED", "ROOT_DENIED", "VMREAD_DENIED", "ATTACH_DENIED", "MEM_DENIED"} {
			s.Require().Containsf(out, want, "a copy outside [inspectors] against %s: %s", pid, out)
		}
		out = s.probe(c, tmpProbe, "exe", pid)
		s.Require().Containsf(out, "EXE_DENIED", "an inspector at a user-writable path must be refused: %s", out)
	}

	log := s.readDaemonLog(c)
	s.Require().Contains(log, "inspector "+tmpProbe+" refused", "the refused entry must be logged")
	s.Require().Regexpf(regexp.MustCompile(`op=PTRACE .*comm=inspector_probe .*resource=multiple `), log,
		"a GLOBAL-slot refusal must be labelled resource=multiple, not one resource:\n%s", log)
	s.exec(c, []string{"sh", "-c", "echo > /tmp/f; echo > /tmp/g; pkill -f 'app-listener daemon' || true"})
}

// An inspector process is a way into every protected process's metadata, so taking it over must
// satisfy every guard: attach, process_vm_readv and a traced exec of it are refused to a caller
// GLOBAL's (intersection) whitelist lacks, and it loads no untrusted library (LD_PRELOAD).
func (s *IntegrationSuite) TestDaemon_Inspector_CannotBeHijacked() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.setupInspectorContainer(c)
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/lib_probe.so"), "/tmp/evil.so", 0o755),
		"copy lib_probe.so")
	s.startDaemon(c, inspectorConfig)
	single, _ := s.startInspectorVictims(c)

	s.exec(c, []string{"sh", "-c", "(" + inspectorProbe + " wait &); sleep 0.5"})
	_, out := s.exec(c, []string{"sh", "-c", "pgrep -f '" + inspectorProbe + " wait' | head -1"})
	insp := strings.TrimSpace(out)
	s.Require().NotEmpty(insp, "the inspector did not start")
	s.Require().Contains(s.probe(c, r1Probe, "exe", single), "EXE_OK", "control: /r1's own caller inspects its process")
	for _, caller := range []string{plainProbe, r1Probe} {
		out = s.probe(c, caller, "vmread", insp)
		s.Require().Containsf(out, "VMREAD_DENIED", "%s must not read an inspector's memory: %s", caller, out)
		out = s.probe(c, caller, "attach", insp)
		s.Require().Containsf(out, "ATTACH_DENIED", "%s must not attach to an inspector: %s", caller, out)
	}
	out = s.probe(c, plainProbe, "tracedexec", inspectorProbe)
	s.Require().Containsf(out, "TRACED_EXEC_DENIED", "a traced inspector exec must be refused: %s", out)
	out = s.probe(c, plainProbe, "tracedexec", plainProbe)
	s.Require().Containsf(out, "TRACED_EXEC_OK", "control: tracing an ordinary exec works: %s", out)

	out = s.probe(c, "LD_PRELOAD=/tmp/evil.so", inspectorProbe, "exe", single)
	s.Require().NotContainsf(out, libProbeMarker, "an inspector must not load an untrusted library: %s", out)
	s.requireDenialLogged(c, "LIBLOAD", "evil.so")

	s.exec(c, []string{"sh", "-c", "echo > /tmp/f; echo > /tmp/g; pkill -f 'inspector_probe wait';" +
		" pkill -f 'app-listener daemon' || true"})
}

// The grant follows the config: an empty [inspectors] block on reload revokes it, and ps-style
// NOAUDIT probes (/proc/<pid>/stat) are refused without a log line.
func (s *IntegrationSuite) TestDaemon_Inspector_RevokedOnReload() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.setupInspectorContainer(c)
	// Outside /etc/app-listener, which the running daemon guards read-only: the test rewrites it.
	const cfgPath = "/tmp/inspector-daemon.conf"
	writeCfg := func(body string) {
		code, out := s.exec(c, []string{"sh", "-c", "cat > " + cfgPath + " <<'EOF'\n" + body + "\nEOF"})
		s.Require().Equalf(0, code, "writing %s: %s", cfgPath, out)
	}
	writeCfg(inspectorConfig)
	code, out := s.exec(c, []string{"sh", "-c",
		"nohup /app-listener daemon --config " + cfgPath + " --headless > /tmp/daemon.log 2>&1 &"})
	s.Require().Equalf(0, code, "starting daemon: %s", out)
	s.awaitDaemonUp(c, inspectorConfig)
	single, _ := s.startInspectorVictims(c)

	s.Require().Contains(s.probe(c, inspectorProbe, "exe", single), "EXE_OK")

	// The plain probe's refusal is logged after cat's would be (one ringbuf, one reader).
	s.exec(c, []string{"sh", "-c", "cat /proc/" + single + "/stat > /dev/null"})
	s.Require().Contains(s.probe(c, plainProbe, "exe", single), "EXE_DENIED")
	logged := false
	for dl := time.Now().Add(10 * time.Second); time.Now().Before(dl) && !logged; {
		logged = strings.Contains(s.readDaemonLog(c), "op=PTRACE  comm=inspector_probe ")
		time.Sleep(200 * time.Millisecond)
	}
	s.Require().Truef(logged, "control: the plain probe's refusal must be logged:\n%s", s.readDaemonLog(c))
	s.Require().NotContainsf(s.readDaemonLog(c), "op=PTRACE  comm=cat ",
		"a NOAUDIT probe must not be logged:\n%s", s.readDaemonLog(c))

	writeCfg(strings.Split(inspectorConfig, "[inspectors]")[0] + "[inspectors]")
	// Logged after the trust guard and [inspectors] are re-applied, unlike the usecase's commit line.
	const reloadDone = "daemon: configuration reloaded from"
	s.sigDaemon(c, "HUP")
	done := false
	for dl := time.Now().Add(daemonShutdownTimeout); time.Now().Before(dl) && !done; {
		done = strings.Contains(s.readDaemonLog(c), reloadDone)
		time.Sleep(300 * time.Millisecond)
	}
	s.Require().Truef(done, "reload did not complete:\n%s", s.readDaemonLog(c))
	out = s.probe(c, inspectorProbe, "exe", single)
	s.Require().Containsf(out, "EXE_DENIED", "an empty [inspectors] block must revoke the grant: %s", out)

	s.exec(c, []string{"sh", "-c", "echo > /tmp/f; echo > /tmp/g; pkill -f 'app-listener daemon' || true"})
}

func (s *IntegrationSuite) setupInspectorContainer(c testcontainers.Container) {
	s.Require().NoError(c.CopyFileToContainer(s.ctx, absPath("./exploits/inspector_probe"), "/exploits/inspector_probe", 0o755),
		"copy inspector_probe")
	s.exec(c, []string{"sh", "-c",
		"mkdir -p /r1 /r3 /etc/app-listener /opt/app /opt/y /opt/insp /opt/plain /opt/r1" +
			" && echo s1 > /r1/secret && echo s3 > /r3/secret" +
			" && cp /usr/bin/dash /opt/app/app && cp /usr/bin/dash /opt/y/y" +
			" && cp /exploits/inspector_probe " + inspectorProbe +
			" && cp /exploits/inspector_probe " + plainProbe +
			" && cp /exploits/inspector_probe " + tmpProbe +
			" && cp /exploits/inspector_probe " + r1Probe +
			" && mkfifo /tmp/f /tmp/g"})
}
