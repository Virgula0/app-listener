package guard

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	cilium "github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"golang.org/x/sys/unix"
)

// trust_launch kinds, modes and key size (guard_trust.bpf.c).
const (
	launchEnv     uint32 = 1
	launchArg     uint32 = 2
	launchOpt     uint32 = 3
	launchPresent uint8  = 1
	launchTokens  uint8  = 2
	launchAllow   uint8  = 3
	launchKeyMax         = 44
)

// launchRule is one trust_launch entry: what in a runtime's launch makes it run code its caller
// chose (issue #79). Identity is the exe inode, so `ELECTRON_RUN_AS_NODE=1 code -e …` or
// `NODE_OPTIONS=--require=/tmp/x.js discord` is the whitelisted binary running the caller's code.
type launchRule struct {
	kind    uint32
	name    string
	classes uint8
	mode    uint8
}

// launchRules: an env var risky when set, one whose every token must be an allowed option, and
// argv flags, matched whole, by first word (--inspect covers VS Code's --inspect-extensions) or by
// first two (--remote-debugging-port); a whole-name launchAllow entry exempts a harmless one. Flag
// and option names are matched with '_' read as '-', so they are spelled with '-' only. Options are
// an allowlist (fail closed): an unknown token is risky. No short flags: app CLIs reuse them
// (`code -r` reuses a window); Chromium's single-dash long switches are listed.
var launchRules = func() []launchRule {
	rules := make([]launchRule, 0, 64)
	rules = append(rules,
		launchRule{launchEnv, "ELECTRON_RUN_AS_NODE", trustedNode, launchPresent},
		launchRule{launchEnv, "NODE_PATH", trustedNode, launchPresent},
		launchRule{launchEnv, "NODE_OPTIONS", trustedNode, launchTokens},
		launchRule{launchEnv, "JAVA_TOOL_OPTIONS", trustedJVM, launchTokens},
		launchRule{launchEnv, "_JAVA_OPTIONS", trustedJVM, launchTokens},
		launchRule{launchEnv, "JDK_JAVA_OPTIONS", trustedJVM, launchTokens},
	)
	for _, a := range []string{"--inspect", "--eval", "--print", "--require", "--import", "--loader",
		"--experimental-loader"} {
		rules = append(rules, launchRule{launchArg, a, trustedNode, launchPresent})
	}
	// The DevTools protocol runs script in the process and reads its cookies; an unpacked extension
	// reads them too. Chromium takes "-" as a switch prefix as well as "--".
	for _, a := range []string{"--remote-debugging", "--load-extension", "-remote-debugging", "-load-extension"} {
		rules = append(rules, launchRule{launchArg, a, trustedNode | trustedChromium, launchPresent})
	}
	// Only set where an inspector would listen, if one were started.
	for _, a := range []string{"--inspect-port", "--inspect-publish-uid"} {
		rules = append(rules, launchRule{launchArg, a, trustedNode, launchAllow})
	}
	for _, o := range []string{
		"--max-old-space-size", "--max-semi-space-size", "--stack-size", "--max-http-header-size",
		"--no-warnings", "--no-deprecation", "--disable-warning", "--trace-warnings",
		"--trace-deprecation", "--throw-deprecation", "--pending-deprecation", "--unhandled-rejections",
		"--enable-source-maps", "--use-openssl-ca", "--use-bundled-ca", "--use-system-ca",
		"--openssl-legacy-provider", "--dns-result-order",
		"--tls-min-v1.2", "--tls-min-v1.3", "--tls-max-v1.2", "--tls-max-v1.3",
	} {
		rules = append(rules, launchRule{launchOpt, o, trustedNode, 0})
	}
	for _, o := range []string{
		"-Xmx", "-Xms", "-Xss",
		"-Dawt.useSystemAAFontSettings", "-Dawt.toolkit.name", "-Dswing.aatext",
		"-Dswing.useSystemFontSettings", "-Dsun.java2d.opengl", "-Dsun.java2d.xrender",
		"-Dsun.java2d.uiScale", "-Dsun.java2d.uiScale.enabled", "-Djdk.gtk.version",
		"-Djava.awt.headless", "-Dfile.encoding", "-Duser.language", "-Duser.country",
	} {
		rules = append(rules, launchRule{launchOpt, o, trustedJVM, 0})
	}
	return rules
}()

func launchKey(kind uint32, name string) (GuardTrustLaunchKey, error) {
	k := GuardTrustLaunchKey{Kind: kind}
	if name == "" || len(name) >= launchKeyMax {
		return k, fmt.Errorf("launch rule %q does not fit the %d-byte key", name, launchKeyMax)
	}
	for i := range len(name) {
		k.Name[i] = int8(name[i]) //nolint:gosec // byte reinterpreted as a C char
	}
	return k, nil
}

// setLaunchRules fills trust_launch before trust_exec_launch attaches: an empty map judges nothing.
func (t *TrustGuard) setLaunchRules() error {
	for _, r := range launchRules {
		k, err := launchKey(r.kind, r.name)
		if err != nil {
			return err
		}
		if err := t.objs.TrustLaunch.Put(k, GuardTrustLaunchRule{Classes: r.classes, Mode: r.mode}); err != nil {
			return fmt.Errorf("launch rule %s: %w", r.name, err)
		}
	}
	return nil
}

// attachLaunchScan attaches trust_exec_launch. Required: without it a whitelisted runtime started
// with NODE_OPTIONS=--require or ELECTRON_RUN_AS_NODE runs its caller's code with full access.
func (t *TrustGuard) attachLaunchScan() error {
	if err := t.setLaunchRules(); err != nil {
		return fmt.Errorf("required trust launch rules: %w", err)
	}
	l, err := link.AttachTracing(link.TracingOptions{
		Program:    t.objs.TrustExecLaunch,
		AttachType: cilium.AttachTraceRawTp,
	})
	if err != nil {
		return fmt.Errorf("required trust hook sched_process_exec failed to attach: %w — a whitelisted "+
			"Node/Electron or JVM runtime could be launched to run its caller's code", err)
	}
	t.links = append(t.links, l)
	return nil
}

// Markers of a runtime whose launch env/argv can make it run other code. The java launcher's own
// options live in libjli, so a JVM is also known by its name.
var runtimeMarkers = []struct {
	marker []byte
	class  uint8
}{
	{[]byte("ELECTRON_RUN_AS_NODE"), trustedNode},
	{[]byte("NODE_OPTIONS"), trustedNode},
	{[]byte("remote-debugging-port"), trustedChromium},
	{[]byte("JDK_JAVA_OPTIONS"), trustedJVM},
	{[]byte("JAVA_TOOL_OPTIONS"), trustedJVM},
}

// trustedRuntimes: every class runtimeClass reports.
const trustedRuntimes = trustedNode | trustedJVM | trustedChromium

type classCacheKey struct {
	dev, ino     uint64
	size         int64
	mtime, ctime unix.Timespec
}

var (
	classMu    sync.Mutex
	classCache = map[classCacheKey]uint8{}
)

// runtimeClass returns the trustedRuntimes classes of the image f holds (path names it). Scanned once
// per inode and content fingerprint: an Electron app's image is ~200 MB. An unreadable image is
// classed as every runtime, so its launches are judged (fail closed).
func runtimeClass(f *os.File, path string) uint8 {
	var st unix.Stat_t
	if err := unix.Fstat(int(f.Fd()), &st); err != nil {
		return trustedRuntimes
	}
	key := classCacheKey{dev: st.Dev, ino: st.Ino, size: st.Size, mtime: st.Mtim, ctime: st.Ctim}
	classMu.Lock()
	c, ok := classCache[key]
	classMu.Unlock()
	if ok {
		return c
	}
	c, err := scanRuntimeClass(f)
	if err != nil {
		return trustedRuntimes
	}
	if base := filepath.Base(path); base == "java" || base == "javaw" {
		c |= trustedJVM
	}
	classMu.Lock()
	if len(classCache) > 4096 {
		classCache = map[classCacheKey]uint8{}
	}
	classCache[key] = c
	classMu.Unlock()
	return c
}

func scanRuntimeClass(f *os.File) (uint8, error) {
	r, err := os.Open(fmt.Sprintf("/proc/self/fd/%d", f.Fd()))
	if err != nil {
		return 0, err
	}
	defer r.Close()
	const chunk, overlap = 1 << 20, 32
	buf := make([]byte, chunk+overlap)
	var c uint8
	carry := 0
	for {
		n, err := io.ReadFull(r, buf[carry:])
		window := buf[:carry+n]
		for _, m := range runtimeMarkers {
			if bytes.Contains(window, m.marker) {
				c |= m.class
			}
		}
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			return c, nil
		}
		if err != nil {
			return 0, err
		}
		carry = copy(buf, window[len(window)-overlap:])
	}
}
