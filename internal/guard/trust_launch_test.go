package guard

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// testElectronApps are the LAUNCH_APP entries launchModel judges against.
var testElectronApps = map[string]bool{"/usr/lib/signal-desktop/app.asar": true}

// launchModel is trust_exec_launch's matching (guard_trust.bpf.c), so the rule table can be judged
// without loading BPF.
func launchModel(classes uint8, argv, env []string) bool {
	rule := func(kind uint32, name string) (launchRule, bool) {
		for _, r := range launchRules {
			if r.kind == kind && r.name == name {
				return r, true
			}
		}
		return launchRule{}, false
	}
	applies := func(r launchRule, ok bool) bool { return ok && r.classes&classes != 0 }
	name := func(s string, kind uint32, words int) string {
		b := []byte(s)
		dashes := 0
		for j, c := range b {
			if c == '=' {
				return string(b[:j])
			}
			if c == '_' && j >= 2 && kind != launchEnv {
				c = '-'
				b[j] = c
			}
			if c == '-' && j >= 2 && words > 0 {
				dashes++
			}
			if words > 0 && dashes == words {
				return string(b[:j])
			}
		}
		return string(b)
	}
	app := false
	for _, a := range argv[1:] {
		if classes&trustedElectron != 0 && !app && !strings.HasPrefix(a, "-") {
			app = true
			if len(a) >= launchKeyMax || !testElectronApps[a] {
				return true
			}
			continue
		}
		r, ok := rule(launchArg, name(a, launchArg, 0))
		if ok && r.mode == launchAllow {
			continue
		}
		for w := 1; w <= 2 && !applies(r, ok); w++ {
			r, ok = rule(launchArg, name(a, launchArg, w))
		}
		if applies(r, ok) {
			return true
		}
	}
	for _, e := range env {
		r, ok := rule(launchEnv, name(e, launchEnv, 0))
		if !applies(r, ok) {
			continue
		}
		if r.mode != launchTokens {
			return true
		}
		_, value, _ := strings.Cut(e, "=")
		for _, tok := range strings.Fields(value) {
			key := name(tok, launchOpt, 0)
			if len(key) > 4 && key[4] >= '0' && key[4] <= '9' {
				key = key[:4]
			}
			if !applies(rule(launchOpt, key)) {
				return true
			}
		}
	}
	return false
}

func TestLaunchRulesJudgeRuntimes(t *testing.T) {
	electron := trustedNode | trustedChromium
	generic := electron | trustedElectron
	for _, c := range []struct {
		name    string
		classes uint8
		argv    []string
		env     []string
		risky   bool
	}{
		{"plain launch", electron, []string{"code", "--no-sandbox", "/tmp/project"}, []string{"HOME=/h"}, false},
		{"run as node", electron, []string{"code", "-e", "x"}, []string{"ELECTRON_RUN_AS_NODE=1"}, true},
		{"node options heap only", electron, []string{"code"}, []string{"NODE_OPTIONS=--max-old-space-size=8192 --no-warnings"}, false},
		{"node options require", electron, []string{"discord"}, []string{"NODE_OPTIONS=--max-old-space-size=8192 --require /tmp/x.js"}, true},
		{"node options short require", electron, []string{"discord"}, []string{"NODE_OPTIONS=-r /tmp/x.js"}, true},
		{"node options inspect", electron, []string{"discord"}, []string{"NODE_OPTIONS=--inspect=9229"}, true},
		{"node path", electron, []string{"discord"}, []string{"NODE_PATH=/tmp/mods"}, true},
		{"inspect flag", electron, []string{"code", "--inspect=9229"}, nil, true},
		{"vscode extension inspector", electron, []string{"code", "--inspect-brk-extensions=9333"}, nil, true},
		{"inspector port alone", electron, []string{"code", "--type=utility", "--inspect-port=0"}, nil, false},
		{"devtools protocol", electron, []string{"discord", "--remote-debugging-port=9222"}, nil, true},
		{"devtools single dash", electron, []string{"discord", "-remote-debugging-port=9222"}, nil, true},
		{"devtools pipe single dash", trustedChromium, []string{"chrome", "-remote-debugging-pipe"}, nil, true},
		{"extension single dash", electron, []string{"discord", "-load-extension=/tmp/ext"}, nil, true},
		{"inspect underscore", electron, []string{"discord", "--inspect_brk"}, nil, true},
		{"inspect wait underscore", electron, []string{"discord", "--inspect_wait=9229"}, nil, true},
		{"inspect brk node underscores", electron, []string{"discord", "--inspect_brk_node"}, nil, true},
		{"loader underscore", electron, []string{"discord", "--experimental_loader=/tmp/x.mjs"}, nil, true},
		{"inspector port underscore", electron, []string{"code", "--inspect_port=0"}, nil, false},
		{"node options underscore heap", electron, []string{"code"}, []string{"NODE_OPTIONS=--max_old_space_size=4096"}, false},
		{"node options underscore require", electron, []string{"code"}, []string{"NODE_OPTIONS=--experimental_loader=/tmp/x.mjs"}, true},
		{"code reuse window", electron, []string{"code", "-r", "file.go"}, nil, false},
		{"code remote", electron, []string{"code", "--remote", "ssh-remote+host"}, nil, false},
		{"browser devtools", trustedChromium, []string{"chrome", "--remote-debugging-pipe"}, nil, true},
		{"browser extension", trustedChromium, []string{"chrome", "--load-extension=/tmp/ext"}, nil, true},
		{"browser ignores node env", trustedChromium, []string{"chrome"}, []string{"NODE_OPTIONS=--require /tmp/x.js"}, false},
		{"jvm fonts", trustedJVM, []string{"java"}, []string{"_JAVA_OPTIONS=-Dawt.useSystemAAFontSettings=on -Dswing.aatext=true -Xmx2g"}, false},
		{"jvm agent", trustedJVM, []string{"java"}, []string{"JAVA_TOOL_OPTIONS=-javaagent:/tmp/a.jar"}, true},
		{"jvm plugin path", trustedJVM, []string{"java"}, []string{"JDK_JAVA_OPTIONS=-Didea.plugins.path=/tmp/p"}, true},
		{"not a runtime", 0, []string{"ssh", "--inspect"}, []string{"ELECTRON_RUN_AS_NODE=1", "NODE_OPTIONS=-r x"}, false},
		{"generic electron, its app", generic, []string{"electron37", "/usr/lib/signal-desktop/app.asar", "/tmp/f"}, nil, false},
		{"generic electron, switch before its app", generic,
			[]string{"electron37", "--ozone-platform=wayland", "/usr/lib/signal-desktop/app.asar"}, nil, false},
		{"generic electron, no app", generic, []string{"electron37", "--type=renderer"}, nil, false},
		{"generic electron, another app", generic, []string{"electron37", "/tmp/evil"}, nil, true},
		{"generic electron, app after --", generic, []string{"electron37", "--", "/tmp/evil"}, nil, true},
		{"generic electron, relative app", generic, []string{"electron37", "app.asar"}, nil, true},
		{"generic electron, '=' past its app", generic,
			[]string{"electron37", "/usr/lib/signal-desktop/app.asar=/../../../tmp/x"}, nil, true},
		{"generic electron, '_' in its app", generic, []string{"electron37", "/usr/lib/signal_desktop/app.asar"}, nil, true},
		{"generic electron, require", generic, []string{"electron37", "-r", "/tmp/x.js", "/usr/lib/signal-desktop/app.asar"}, nil, true},
		{"generic electron, repl", generic, []string{"electron37", "-i"}, nil, true},
		{"generic electron, --app= file", generic, []string{"electron37", "--app=/tmp/evil.js"}, nil, true},
		{"generic electron, --app= before its app", generic,
			[]string{"electron37", "--app=/tmp/evil.js", "/usr/lib/signal-desktop/app.asar"}, nil, true},
		{"generic electron, single-dash repl", generic, []string{"electron37", "-repl"}, nil, true},
		{"bundled electron ignores its positional", electron, []string{"code", "/tmp/evil"}, nil, false},
		{"renderer launcher", electron, []string{"discord", "--renderer-cmd-prefix=/tmp/w"}, nil, true},
		{"utility launcher single dash", trustedChromium, []string{"chrome", "-utility-cmd-prefix=/tmp/w"}, nil, true},
		{"subprocess path", trustedChromium, []string{"chrome", "--browser-subprocess-path=/tmp/w"}, nil, true},
		{"gpu launcher", electron, []string{"code", "--gpu-launcher=/tmp/w"}, nil, true},
		{"vscode agent host", electron, []string{"code", "--type=utility", "--logsPath", "/h/.config/Code/logs/1",
			"--user-data-dir", "/h/.config/Code"}, nil, false},
		{"vscode extensions dir", electron, []string{"code", "--extensions-dir", "/tmp/e"}, nil, true},
		{"vscode extension development", electron, []string{"code", "--extensionDevelopmentPath=/tmp/e"}, nil, true},
		{"vscode install extension", electron, []string{"code", "--install-extension", "/tmp/e.vsix"}, nil, true},
		{"node env file", electron, []string{"node", "--env-file=/tmp/.env"}, nil, true},
		{"node repl module", electron, []string{"node"}, []string{"NODE_REPL_EXTERNAL_MODULE=/tmp/m.js"}, true},
		{"logging to a file", electron, []string{"code", "--enable-logging=file", "--log-file=/tmp/l"}, nil, false},
	} {
		if got := launchModel(c.classes, c.argv, c.env); got != c.risky {
			t.Errorf("%s: risky = %v, want %v", c.name, got, c.risky)
		}
	}
}

func TestLaunchRulesFitTheKey(t *testing.T) {
	seen := map[[2]string]bool{}
	for _, r := range launchRules {
		if _, err := launchKey(r.kind, r.name); err != nil {
			t.Error(err)
		}
		if r.kind != launchEnv && strings.Contains(r.name[min(2, len(r.name)):], "_") {
			t.Errorf("launch rule %q: flag and option names are matched with '_' read as '-'", r.name)
		}
		k := [2]string{string(rune('0' + r.kind)), r.name}
		if seen[k] {
			t.Errorf("duplicate launch rule %q", r.name)
		}
		seen[k] = true
	}
}

func TestRuntimeClassMarkers(t *testing.T) {
	dir := t.TempDir()
	write := func(name string, content []byte) *os.File {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, content, 0o755); err != nil {
			t.Fatal(err)
		}
		f, err := os.Open(p)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { f.Close() })
		return f
	}
	// A marker split across the 1 MiB read boundary still counts.
	big := make([]byte, 1<<20+64)
	copy(big[1<<20-8:], "ELECTRON_RUN_AS_NODE")
	if c := runtimeClass(write("electron", big), filepath.Join(dir, "electron")); c&trustedNode == 0 {
		t.Errorf("Electron marker across chunks not found: class %d", c)
	}
	if c := runtimeClass(write("chrome", []byte("x\x00remote-debugging-port\x00")), "chrome"); c != trustedChromium {
		t.Errorf("chrome class = %d, want Chromium only", c)
	}
	if c := runtimeClass(write("java", []byte("ELF")), filepath.Join(dir, "java")); c != trustedJVM {
		t.Errorf("java launcher class = %d, want JVM by name", c)
	}
	if c := runtimeClass(write("ssh", []byte("ELF ssh")), "ssh"); c != 0 {
		t.Errorf("ssh class = %d, want none", c)
	}
}

func TestRuntimeClassGenericElectron(t *testing.T) {
	exe := func(dir string, resources ...string) (*os.File, string) {
		for _, r := range resources {
			if err := os.MkdirAll(filepath.Join(dir, "resources", filepath.Dir(r)), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(dir, "resources", r), nil, 0o644); err != nil {
				t.Fatal(err)
			}
		}
		p := filepath.Join(dir, "electron")
		if err := os.WriteFile(p, []byte("ELF\x00ELECTRON_RUN_AS_NODE\x00"), 0o755); err != nil {
			t.Fatal(err)
		}
		f, err := os.Open(p)
		if err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { f.Close() })
		return f, p
	}
	f, p := exe(t.TempDir(), "default_app.asar")
	if c := runtimeClass(f, p); c&trustedElectron == 0 || c&trustedNode == 0 {
		t.Errorf("an Electron with default_app.asar beside it runs its argv's app: class %d", c)
	}
	// Electron falls back to default_app.asar when app.asar doesn't load: still generic.
	f, p = exe(t.TempDir(), "default_app.asar", "app.asar")
	if c := runtimeClass(f, p); c&trustedElectron == 0 {
		t.Errorf("default_app.asar beside a bundled app is still reachable: class %d", c)
	}
	f, p = exe(t.TempDir(), "app.asar")
	if c := runtimeClass(f, p); c&trustedElectron != 0 || c&trustedNode == 0 {
		t.Errorf("a bundled app's Electron chooses no app from argv: class %d", c)
	}
	// The class follows the real exe, not the configured link (/usr/bin/electron37).
	dir := t.TempDir()
	f, p = exe(dir, "default_app.asar")
	link := filepath.Join(t.TempDir(), "electron37")
	if err := os.Symlink(p, link); err != nil {
		t.Fatal(err)
	}
	if c := runtimeClass(f, link); c&trustedElectron == 0 {
		t.Errorf("generic Electron through a link: class %d", c)
	}
}

// A temporary grant classes the inode it pins: trustTempRuntime puts these bits, never TRUSTED_BINARY.
func TestResolveTempBinaryClassesRuntime(t *testing.T) {
	dir := t.TempDir()
	for _, tc := range []struct {
		name    string
		content string
		want    uint8
		names   []string
	}{
		{"code", "ELF\x00NODE_OPTIONS\x00", trustedNode, []string{"Node"}},
		{"chrome", "ELF\x00remote-debugging-port\x00", trustedChromium, []string{"Chromium"}},
		{"tool", "ELF plain", 0, nil},
	} {
		p := filepath.Join(dir, tc.name)
		if err := os.WriteFile(p, []byte(tc.content), 0o755); err != nil {
			t.Fatal(err)
		}
		b, err := ResolveTempBinary(p)
		if err != nil {
			t.Fatalf("%s: %v", tc.name, err)
		}
		if b.Runtime != tc.want || b.Runtime&trustedBinary != 0 {
			t.Errorf("%s: Runtime = %d, want %d", tc.name, b.Runtime, tc.want)
		}
		if got := b.RuntimeNames(); !slices.Equal(got, tc.names) {
			t.Errorf("%s: RuntimeNames = %v, want %v", tc.name, got, tc.names)
		}
		b.Close()
	}
}
