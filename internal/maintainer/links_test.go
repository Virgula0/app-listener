package maintainer

import (
	"os"
	"path/filepath"
	"testing"
)

// linkFixture: lib_dirs lib and ui below a user-owned root, the user's own files beside them, and
// a live symlink in lib out to them (as an older Steam or a pre-check commit could have left).
func linkFixture(t *testing.T) (root string, libs []string) {
	t.Helper()
	root = t.TempDir()
	write(t, filepath.Join(root, "lib", "rt", "usr", "lib", "libdemo.so.1"), "ELF", 0o755)
	write(t, filepath.Join(root, "ui", "index.js"), "ok", 0o644)
	write(t, filepath.Join(root, "cache", "evil.js"), "steal()", 0o644)
	if err := os.Symlink(filepath.Join(root, "cache"), filepath.Join(root, "lib", "out")); err != nil {
		t.Fatal(err)
	}
	return root, []string{filepath.Join(root, "lib"), filepath.Join(root, "ui")}
}

// checkLink judges only the link rel -> target, with extra in the plan beside it.
func checkLink(t *testing.T, libs []string, lib, rel, target string, extra ...op) error {
	t.Helper()
	ops := append([]op{{kind: opLink, rel: rel, target: target}}, extra...)
	lp := libPlan{dir: lib, ops: ops}
	only := libPlan{dir: lib, ops: ops[:1]}
	return newLinkView(libs, []libPlan{lp}).checkLinks(&only)
}

func TestCheckLinksAccepts(t *testing.T) {
	_, libs := linkFixture(t)
	lib := libs[0]
	for target, why := range map[string]string{
		filepath.Join(lib, "rt", "usr", "lib", "libdemo.so.1"): "a pin into the runtime",
		"../usr/lib/libdemo.so.1":                              "a relative pin",
		"../../../ui/index.js":                                 "another lib_dir",
		filepath.Join(lib, "rt", "not-yet"):                    "a name only the app may create",
		"/usr/lib/x86_64-linux-gnu/libc.so.6":                  "a system library",
		"/usr/lib/app-listener-missing/libc.so.6":              "a missing name in a root-only directory",
		"/run/host/usr/lib/libGL.so.1":                         "a name only the container runtime binds",
	} {
		if err := checkLink(t, libs, lib, "rt/pins/link", target); err != nil {
			t.Errorf("%s (%s) refused: %v", why, target, err)
		}
	}
	// The runtime's usr/lib/ssl/certs: a directory only root fills, all the way down.
	if _, err := os.Stat("/etc/ssl/certs"); err == nil {
		if err := checkLink(t, libs, lib, "rt/ssl/certs", "/etc/ssl/certs"); err != nil {
			t.Errorf("a root-only system directory refused: %v", err)
		}
	}
}

func TestCheckLinksRefusesUserFiles(t *testing.T) {
	root, libs := linkFixture(t)
	lib := libs[0]
	evil := filepath.Join(root, "cache", "evil.js")
	for target, why := range map[string]string{
		evil:                        "the user's file",
		"../cache/evil.js":          "a relative climb out of the lib_dir",
		"../cache":                  "the user's directory",
		"out/evil.js":               "through a live symlink out",
		"rt/../../cache/evil.js":    `a ".." after a name`,
		"/proc/self/cwd/evil.js":    "the reader's own cwd",
		"/tmp/evil.js":              "a world-writable directory",
		"/dev/shm/evil.js":          "a world-writable directory",
		filepath.Join(lib, "loop"):  "a loop",
		filepath.Join(lib, "hop/x"): "through a planned symlink out",
		// Directories a user fills, in a directory only root decides: what is looked up through
		// the link (steamui/libraries/*.js) is never judged again.
		"/tmp":     "a world-writable directory itself",
		"/var/tmp": "a world-writable directory itself",
		"/dev/shm": "a world-writable directory itself",
		"/home/app-listener-no-such-user/evil.js": "a home useradd would hand a user",
		"/run/user/4242424242/evil.js":            "a runtime dir logind would hand a user",
	} {
		hop := op{kind: opLink, rel: "hop", target: filepath.Join(root, "cache")}
		loop := op{kind: opLink, rel: "loop", target: "loop"}
		if err := checkLink(t, libs, lib, "x.js", target, hop, loop); err == nil {
			t.Errorf("%s (%s) accepted", why, target)
		}
	}
	if home, err := os.UserHomeDir(); err == nil && os.Geteuid() != 0 {
		if err := checkLink(t, libs, lib, "libraries", home); err == nil {
			t.Errorf("a link to the user's home (%s) accepted", home)
		}
	}
}

// Each target is judged by itself: a sibling name root decides vouches for none beside it.
func TestCheckLinksJudgesEachTarget(t *testing.T) {
	_, libs := linkFixture(t)
	lib := libs[0]
	ops := []op{{kind: opLink, rel: "ok", target: "/app-listener-no-such-dir"}, {kind: opLink, rel: "bad", target: "/tmp"}}
	if err := newLinkView(libs, []libPlan{{dir: lib, ops: ops}}).checkLinks(&libPlan{dir: lib, ops: ops}); err == nil {
		t.Error("a link to /tmp accepted after one to a missing name beside it")
	}
}

// A live link out is no longer followed once the plan deletes it or replaces its directory.
func TestCheckLinksSeesThePlannedTree(t *testing.T) {
	_, libs := linkFixture(t)
	lib := libs[0]
	if err := checkLink(t, libs, lib, "x", "out/evil.js", op{kind: opDelete, rel: "out"}); err != nil {
		t.Errorf("a deleted live link still judged: %v", err)
	}
	if err := checkLink(t, libs, lib, "x", "out/evil.js", op{kind: opDir, rel: "out"}); err != nil {
		t.Errorf("a live link replaced by a directory still judged: %v", err)
	}
	root := op{kind: opDir, rel: ".", opaque: true, names: map[string]bool{"x": true}}
	if err := checkLink(t, libs, lib, "x", "out/evil.js", root); err != nil {
		t.Errorf("a link gone with its replaced lib_dir still judged: %v", err)
	}
}
