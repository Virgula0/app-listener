package guard

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

func TestEncodeLaunchSet(t *testing.T) {
	rules, _ := canonLaunchRules([]daemonconfig.LaunchRule{
		deny("auth", "token"), deny("auth", "status", "-t"),
		{Words: []string{"auth", "git-credential"}, Warn: true},
		{Env: "SSL_CERT_FILE", Warn: true}, {Env: "SSL_CERT_FILE"}, {Env: "HTTPS_PROXY", Warn: true},
	})
	words, set, err := encodeLaunchSet(rules)
	if err != nil {
		t.Fatal(err)
	}
	// Canonical order: argv rules first (empty Env), denies before warns, then by words.
	if rules[0].Words[1] != "status" || rules[1].Words[1] != "token" || rules[2].Words[1] != "git-credential" {
		t.Fatalf("canonical order = %+v", rules)
	}
	want := map[argvEntry]uint64{
		{argvWord, "auth"}:           1<<0 | 1<<4 | 1<<8,
		{argvWord, "status"}:         1 << 1,
		{argvWord, "-t"}:             1 << 2,
		{argvWord, "token"}:          1 << 5,
		{argvWord, "git-credential"}: 1 << 9,
		{argvEnv, "SSL_CERT_FILE"}:   argvDeny, // a deny beats a warn naming the same variable
		{argvEnv, "HTTPS_PROXY"}:     argvWarn,
	}
	if len(words) != len(want) {
		t.Fatalf("words = %v", words)
	}
	for e, v := range want {
		if words[e] != v {
			t.Errorf("%v = %#x, want %#x", e, words[e], v)
		}
	}
	if set.Lens != 3|2<<4|2<<8 || set.Warn != 1<<2 {
		t.Errorf("set = lens %#x warn %#x", set.Lens, set.Warn)
	}
}

func TestEncodeLaunchSet_RefusesTooMany(t *testing.T) {
	var rules []daemonconfig.LaunchRule
	for i := range argvRules + 1 {
		rules = append(rules, deny(string(rune('a'+i))))
	}
	if _, _, err := encodeLaunchSet(rules); err == nil {
		t.Fatal("more than argvRules argv rules must be refused")
	}
	if _, _, err := encodeLaunchSet([]daemonconfig.LaunchRule{deny("a", "b", "c", "d", "e")}); err == nil {
		t.Fatal("a rule longer than argvWords must be refused")
	}
}

func TestBuiltinLaunchRulesFit(t *testing.T) {
	for name := range builtinLaunchRules {
		rules, _ := canonLaunchRules(append(BuiltinLaunchRules(name), goEnvRules...))
		if _, _, err := encodeLaunchSet(rules); err != nil {
			t.Errorf("%s: %v", name, err)
		}
		for _, r := range rules {
			for _, w := range r.Words {
				if len(w) > daemonconfig.MaxLaunchWordLen {
					t.Errorf("%s: word %q too long", name, w)
				}
			}
		}
	}
	a := BuiltinLaunchRules("bitcoin-cli")
	a = append(a, deny("mine"))
	if b := BuiltinLaunchRules("litecoin-cli"); len(b) != len(a)-1 || b[len(b)-1].Words[0] == "mine" {
		t.Fatal("BuiltinLaunchRules must return a copy")
	}
}

func TestIsGoBinary(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	self, err := ebpf.OpenConfined(exe) // O_PATH, as SetTrusted passes it
	if err != nil {
		t.Fatal(err)
	}
	defer self.Close()
	if !isGoBinary(self) {
		t.Error("the Go test binary is not detected")
	}
	p := filepath.Join(t.TempDir(), "script")
	if err := os.WriteFile(p, []byte("#!/bin/sh\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	f, err := ebpf.OpenConfined(p)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if isGoBinary(f) {
		t.Error("a script is detected as Go")
	}
}
