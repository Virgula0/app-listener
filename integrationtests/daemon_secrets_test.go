package integrationtests

import (
	"strings"
	"testing"

	"github.com/testcontainers/testcontainers-go"
)

const (
	ghToken     = "SECRET-gh-token"
	secretScope = "/run/app-listener/secret-scopes/gh"
	secretRelay = "/run/app-listener/secret-relay"
	keyringDir  = "/run/user/1000/keyring"
)

// A [secret gh] scope: gnome-keyring reaches the user bus only through the daemon's relay, and every
// call touching a gh item needs a caller the scope admits, attested by the pidfd the bus took at
// connect. The probes cover the item calls, a replace matching by a subset, owning the name,
// monitoring, exec after connecting or after sending, a send on a socket another process connected, a
// fake keyring on the relay, and the pkcs11 socket.
func (s *IntegrationSuite) TestDaemon_SecretService() {
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)

	code, out := s.exec(c, []string{"sh", "-c", "command -v gnome-keyring-daemon >/dev/null 2>&1 || " +
		"(timeout 120 apt-get update -qq && DEBIAN_FRONTEND=noninteractive timeout 300 apt-get install -y -qq " +
		"--no-install-recommends dbus-daemon gnome-keyring gnome-keyring-pkcs11 opensc >/dev/null)"})
	if code != 0 {
		s.T().Skipf("gnome-keyring unavailable: %s", out)
	}
	mount := "mkdir -p /run/user/1000 && mount -t tmpfs -o mode=0700,uid=1000,gid=1000 tmpfs /run/user/1000"
	code, out = s.exec(c, []string{"sh", "-c", "(id -u 1000 >/dev/null 2>&1 || useradd -u 1000 -m u1000) && " +
		"mkdir -p /watch /exploits /etc/app-listener && " + mount})
	s.Require().Equalf(0, code, "setup: %s", out)
	for _, name := range []string{"secret_probe", "secret_probe_ok"} {
		s.Require().NoError(c.CopyFileToContainer(s.ctx, secretProbeAmd64Bin, "/exploits/"+name, 0o755))
	}

	const always = "[watch /watch]\nneed_encryption: false\n/usr/bin/true\n"
	conf := always + "\n[secret_service]\nprovider /usr/bin/gnome-keyring-daemon\n\n" +
		"[secret gh]\nmatch service=gh:*\n/exploits/secret_probe_ok\n\n" +
		"[watch " + keyringDir + "]\nruntime: true\n/usr/bin/gnome-keyring-daemon MKNOD,DELETE,ATTR,CONNECT,STAT\n" +
		"need_encryption: false\n"
	s.exec(c, []string{"sh", "-c", "cat > /etc/app-listener/daemon.conf <<'EOF'\n" + conf + "EOF"})
	s.launchDaemon(c)
	// The runtime section guards nothing before the login.
	s.awaitDaemonUp(c, always+"\n[secret_service]\n[secret gh]\n")
	code, out = s.exec(c, []string{"ls", "/etc/dbus-1/session.d/app-listener-secrets.conf",
		"/etc/systemd/user/gnome-keyring-daemon.service.d/app-listener.conf"})
	s.Require().Equalf(0, code, "the daemon must install the bus policy and the keyring drop-in: %s", out)

	code, out = s.exec(c, []string{"/app-listener", "daemon", "login", "1000"})
	s.Require().Equalf(0, code, "daemon login: %s\n%s", out, s.readDaemonLog(c))
	code, out = s.exec(c, []string{"stat", "-c", "%u %a %F", secretRelay + "/1000.bus"})
	s.Require().Equalf(0, code, "relay socket: %s\n%s", out, s.readDaemonLog(c))
	s.Require().Equalf("1000 600 socket", out, "the relay socket is the user's, 0600")

	user := func(env, cmd string) (int, string) {
		return s.exec(c, []string{"setpriv", "--reuid=1000", "--regid=1000", "--clear-groups", "sh", "-c",
			"export HOME=$(getent passwd 1000 | cut -d: -f6) XDG_RUNTIME_DIR=/run/user/1000 " +
				"DBUS_SESSION_BUS_ADDRESS=unix:path=/run/user/1000/bus " + env + "; " + cmd})
	}
	code, out = user("", "(nohup dbus-daemon --session --address=$DBUS_SESSION_BUS_ADDRESS --nofork "+
		">/tmp/bus.log 2>&1 &) && for i in $(seq 50); do [ -S /run/user/1000/bus ] && exit 0; sleep 0.1; done; "+
		"cat /tmp/bus.log; exit 1")
	s.Require().Equalf(0, code, "session bus: %s", out)
	code, out = user("DBUS_SESSION_BUS_ADDRESS=unix:path="+secretRelay+"/1000.bus",
		"(echo -n hunter2 | nohup gnome-keyring-daemon --foreground --unlock --components=pkcs11,secrets "+
			"--control-directory="+keyringDir+" >/tmp/gk.log 2>&1 &); "+
			"for i in $(seq 100); do /exploits/secret_probe_ok store service gh:github.com "+ghToken+
			" >/tmp/store.log 2>&1 && exit 0; sleep 0.2; done; cat /tmp/store.log /tmp/gk.log; exit 1")
	s.Require().Equalf(0, code, "the whitelisted app could not store through the relay: %s\n%s", out, s.readDaemonLog(c))

	ok := func(what, cmd string) {
		code, out := user("", cmd)
		s.Require().Equalf(0, code, "%s: %s\n%s", what, out, s.readDaemonLog(c))
		s.Require().Containsf(out, "SECRET="+ghToken, "%s must read the original token", what)
	}
	refused := func(what, cmd, op, comm string) {
		code, out := user("", cmd)
		s.Require().NotContainsf(out, ghToken, "%s leaked the token", what)
		s.Require().Equalf(1, code, "%s must be refused (exit 1): %s\n%s", what, out, s.readDaemonLog(c))
		s.requireSecretDenial(c, op, comm, secretScope)
	}

	s.T().Run("items", func(*testing.T) {
		ok("whitelisted get", "/exploits/secret_probe_ok get service gh:github.com")
		refused("attacker get", "/exploits/secret_probe get service gh:github.com", "READ", "secret_probe")
		refused("attacker overwrite", "/exploits/secret_probe store service gh:github.com SECRET-evil",
			"WRITE", "secret_probe")
		ok("whitelisted get after the overwrite", "/exploits/secret_probe_ok get service gh:github.com")
		// No scope covers it: as open as the keyring keeps it.
		code, out := user("", "/exploits/secret_probe store service other SECRET-other && "+
			"/exploits/secret_probe get service other")
		s.Require().Equalf(0, code, "an unscoped item must stay usable: %s", out)
	})

	s.T().Run("replace", func(*testing.T) {
		// replace=true picks the item by a subset of its attributes: {user: alice} and {} reach gh's.
		refused("attacker replace by a subset", "/exploits/secret_probe replace ok user alice", "WRITE", "secret_probe")
		refused("attacker replace by no attributes", "/exploits/secret_probe replace ok", "WRITE", "secret_probe")
		refused("attacker replace with a failing secret", "/exploits/secret_probe replace fail user alice",
			"WRITE", "secret_probe")
		ok("whitelisted get after the replaces", "/exploits/secret_probe_ok get service gh:github.com")
	})

	s.T().Run("bus", func(*testing.T) {
		code, out := user("", "/exploits/secret_probe own")
		s.Require().Equalf(1, code, "the bus policy must refuse the Secret Service name: %s", out)
		code, out = user("", "(/exploits/secret_probe monitor SECRET-gh 4s >/tmp/mon.log 2>&1; echo rc=$? >>/tmp/mon.log) & "+
			"sleep 1; /exploits/secret_probe_ok get service gh:github.com >/dev/null; wait; cat /tmp/mon.log")
		s.Require().NotContainsf(out, "CAPTURED", "a monitor saw the token: %s", out)
		s.Require().NotContainsf(out, "rc=0", "monitoring must fail: %s", out)
	})

	s.T().Run("exec after connect", func(*testing.T) {
		code, out := user("", "/exploits/secret_probe connexec service gh:github.com /exploits/secret_probe "+
			"/exploits/secret_probe_ok sleep 3s")
		s.Require().NotContainsf(out, ghToken, "a connection made before exec'ing the app leaked: %s", out)
		s.Require().Containsf(out, "REFUSED", "want the gate's refusal (exit %d): %s\n%s", code, out, s.readDaemonLog(c))
		code, out = user("", "/exploits/secret_probe sendexec service gh:github.com /exploits/secret_probe "+
			"/exploits/secret_probe_ok sleep 6s")
		s.Require().NotContainsf(out, ghToken, "a call sent before exec'ing the app leaked (exit %d): %s", code, out)
	})

	s.T().Run("foreign sender", func(*testing.T) {
		// The app lends its bus socket: the bus credits the borrower's call to the app.
		code, out := user("", "/exploits/secret_probe_ok lend service gh:github.com /exploits/secret_probe")
		s.Require().NotContainsf(out, ghToken, "a call sent on the app's lent socket leaked: %s", out)
		s.Require().Equalf(1, code, "want the gate's refusal: %s\n%s", out, s.readDaemonLog(c))
		ok("the app on a fresh process after a lend", "/exploits/secret_probe_ok get service gh:github.com")
	})

	s.T().Run("fake keyring", func(*testing.T) {
		code, out := user("DBUS_SESSION_BUS_ADDRESS=unix:path="+secretRelay+"/1000.bus",
			"/exploits/secret_probe own")
		s.Require().Equalf(1, code, "only the provider may connect to the relay: %s", out)
		s.requireSecretDenial(c, "CONNECT", "secret_probe", secretRelay)
	})

	s.T().Run("pkcs11", func(t *testing.T) {
		const module = "/usr/lib/x86_64-linux-gnu/pkcs11/gnome-keyring-pkcs11.so"
		code, out := user("", "command -v pkcs11-tool && ls "+module)
		if code != 0 {
			t.Skipf("pkcs11-tool or the keyring module missing: %s", out)
		}
		_, out = user("", "pkcs11-tool --module "+module+" --slot-index 1 --read-object --type secrkey "+
			"--label 'probe gh:github.com' 2>&1")
		s.Require().NotContainsf(out, ghToken, "the pkcs11 socket handed out the token: %s", out)
		s.requireSecretDenial(c, "CONNECT", "pkcs11-tool", keyringDir)
	})

	s.reloadAndAwait(c)
	s.Require().NotContainsf(s.readDaemonLog(c), reloadRefused, "reload refused:\n%s", s.readDaemonLog(c))
	ok("whitelisted get after a reload", "/exploits/secret_probe_ok get service gh:github.com")
	refused("attacker get after a reload", "/exploits/secret_probe get service gh:github.com", "READ", "secret_probe")

	s.exec(c, []string{"sh", "-c", "pkill gnome-keyring; pkill dbus-daemon; pkill -f 'app-listener daemon' || true"})
	s.Require().True(s.awaitDaemonDead(c, daemonShutdownTimeout), "daemon did not exit")
}

// requireSecretDenial asserts a DAEMON DENIED line for op by comm under resource.
func (s *IntegrationSuite) requireSecretDenial(c testcontainers.Container, op, comm, resource string) {
	log := s.readDaemonLog(c)
	for _, e := range parseDaemonEvents(log) {
		if e.Denied && e.Op == op && e.Comm == comm && strings.HasPrefix(e.Resource, resource) {
			return
		}
	}
	s.Require().Failf("no denial logged", "no DAEMON DENIED op=%s comm=%s resource=%s:\n%s", op, comm, resource, log)
}
