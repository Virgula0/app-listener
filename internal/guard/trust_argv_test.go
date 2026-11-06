package guard

import (
	"slices"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

func TestEncodeLaunchSet(t *testing.T) {
	rules, allowed, _ := canonLaunchRules([]daemonconfig.LaunchRule{
		deny("auth", "token"), deny("auth", "status", "-t"),
		{Words: []string{"auth", "git-credential"}, Warn: true},
		{Env: "SSL_CERT_FILE", Warn: true}, {Env: "SSL_CERT_FILE"}, {Env: "HTTPS_PROXY", Warn: true},
	}, []string{"SSL_CERT_FILE=/etc/b.pem", "SSL_CERT_DIR=/etc/c", "SSL_CERT_FILE=/etc/a.pem", "SSL_CERT_FILE=/etc/a.pem"})
	if !slices.Equal(allowed, []string{"SSL_CERT_FILE=/etc/a.pem", "SSL_CERT_FILE=/etc/b.pem"}) {
		t.Fatalf("allowed = %v: only the variables a rule names, sorted, once", allowed)
	}
	words, set, err := encodeLaunchSet(rules, allowed)
	if err != nil {
		t.Fatal(err)
	}
	// Canonical order: argv rules first (empty Env), denies before warns, then by words.
	if rules[0].Words[1] != "status" || rules[1].Words[1] != "token" || rules[2].Words[1] != "git-credential" {
		t.Fatalf("canonical order = %+v", rules)
	}
	want := map[argvEntry]uint64{
		{argvWord, "auth"}:                      1<<0 | 1<<4 | 1<<8,
		{argvWord, "status"}:                    1 << 1,
		{argvWord, "-t"}:                        1 << 2,
		{argvWord, "token"}:                     1 << 5,
		{argvWord, "git-credential"}:            1 << 9,
		{argvEnv, "SSL_CERT_FILE"}:              argvDeny, // a deny beats a warn naming the same variable
		{argvEnv, "HTTPS_PROXY"}:                argvWarn,
		{argvEnvOK, "SSL_CERT_FILE=/etc/a.pem"}: 1,
		{argvEnvOK, "SSL_CERT_FILE=/etc/b.pem"}: 1,
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
	if _, _, err := encodeLaunchSet(rules, nil); err == nil {
		t.Fatal("more than argvRules argv rules must be refused")
	}
	if _, _, err := encodeLaunchSet([]daemonconfig.LaunchRule{deny("a", "b", "c", "d", "e")}, nil); err == nil {
		t.Fatal("a rule longer than argvWords must be refused")
	}
}

func TestCanonLaunchRulesKeyCoversAllowed(t *testing.T) {
	_, _, a := canonLaunchRules(caEnvRules, []string{"SSL_CERT_FILE=/etc/a.pem"})
	_, _, b := canonLaunchRules(caEnvRules, []string{"SSL_CERT_FILE=/etc/b.pem"})
	_, _, c := canonLaunchRules(caEnvRules, []string{"SSL_CERT_FILE=/etc/a.pem", "OTHER=/x"})
	if a == b || a != c {
		t.Fatalf("a set's key must name the allowed values of its variables, only those: %q %q %q", a, b, c)
	}
}

func TestBuiltinLaunchRulesFit(t *testing.T) {
	for name := range builtinLaunchRules {
		rules, allowed, _ := canonLaunchRules(append(BuiltinLaunchRules(name), caEnvRules...), BuiltinAllowedEnv())
		if _, _, err := encodeLaunchSet(rules, allowed); err != nil {
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
	var k GuardTrustArgvKey
	for _, a := range BuiltinAllowedEnv() {
		if _, err := daemonconfig.ParseAllowEnv(a); err != nil || len(a) >= len(k.Name) {
			t.Errorf("built-in %q: %v (len %d)", a, err, len(a))
		}
		if name, _, _ := strings.Cut(a, "="); !slices.ContainsFunc(caEnvRules,
			func(r daemonconfig.LaunchRule) bool { return r.Env == name }) {
			t.Errorf("built-in %q allows a variable no CA rule names", a)
		}
	}
	a := BuiltinLaunchRules("bitcoin-cli")
	a = append(a, deny("mine"))
	if b := BuiltinLaunchRules("litecoin-cli"); len(b) != len(a)-1 || b[len(b)-1].Words[0] == "mine" {
		t.Fatal("BuiltinLaunchRules must return a copy")
	}
}

func TestPiLaunchRulesRefuseExtensionLoaders(t *testing.T) {
	rules := BuiltinLaunchRules("pi")
	for _, want := range []daemonconfig.LaunchRule{deny("-e"), deny("--extension"), {Env: "PI_CODING_AGENT_DIR"}} {
		if !slices.ContainsFunc(rules, func(r daemonconfig.LaunchRule) bool {
			return r.Env == want.Env && slices.Equal(r.Words, want.Words) && !r.Warn
		}) {
			t.Errorf("pi has no deny rule %+v: %+v", want, rules)
		}
	}
}

// Node and Bun trust NODE_EXTRA_CA_CERTS and skip verification with NODE_TLS_REJECT_UNAUTHORIZED=0.
func TestCAEnvRulesCoverNode(t *testing.T) {
	for _, env := range []string{"SSL_CERT_FILE", "SSL_CERT_DIR", "NODE_EXTRA_CA_CERTS", "NODE_TLS_REJECT_UNAUTHORIZED"} {
		if !slices.ContainsFunc(caEnvRules, func(r daemonconfig.LaunchRule) bool { return r.Env == env && !r.Warn }) {
			t.Errorf("no CA env deny rule for %s", env)
		}
	}
}

// Env and config flags that redirect an agent's token or load a config of its caller's choosing.
func TestAgentRedirectRules(t *testing.T) {
	has := func(name, env string) bool {
		return slices.ContainsFunc(BuiltinLaunchRules(name), func(r daemonconfig.LaunchRule) bool {
			return r.Env == env && !r.Warn
		})
	}
	for _, c := range [][2]string{
		{"opencode", "OPENCODE_CONFIG"}, {"opencode", "OPENCODE_CONFIG_DIR"}, {"opencode", "OPENCODE_CONFIG_CONTENT"},
		{"codex", "OPENAI_BASE_URL"},
	} {
		if !has(c[0], c[1]) {
			t.Errorf("%s: no deny_env %s", c[0], c[1])
		}
	}
}

// shortHit mirrors guard_trust.bpf.c argv_short over set's masks.
func shortHit(set GuardTrustArgvSet, word string) bool {
	if set.ShortDeny == 0 || len(word) < 2 || word[0] != '-' || word[1] == '-' {
		return false
	}
	for j := 1; j < len(word) && j < 32; j++ {
		m, err := shortMask(word[j : j+1])
		if err != nil {
			return false
		}
		if set.ShortDeny&m != 0 {
			return true
		}
		if set.ShortArg&m != 0 {
			return false
		}
	}
	return false
}

func TestSSHConfigFileRefusedInEveryCluster(t *testing.T) {
	rules, allowed, _ := canonLaunchRules(BuiltinLaunchRules("ssh"), nil)
	_, set, err := encodeLaunchSet(rules, allowed)
	if err != nil {
		t.Fatal(err)
	}
	for word, want := range map[string]bool{
		"-F": true, "-F/tmp/c": true, "-4F/tmp/c": true, "-vvqF": true, "-TF": true,
		"-oF": false, "-oForwardAgent=yes": false, "-L8080:F": false, "-iF": false, "--F": false,
		"-": false, "host": false, "-v": false,
	} {
		if got := shortHit(set, word); got != want {
			t.Errorf("ssh %q: refused = %v, want %v", word, got, want)
		}
	}
}

func TestSSHServerBinariesAreRootOnly(t *testing.T) {
	for _, n := range []string{"sshd", "sshd-session", "sshd-auth"} {
		rules, allowed, _ := canonLaunchRules(BuiltinLaunchRules(n), nil)
		_, set, err := encodeLaunchSet(rules, allowed)
		if err != nil {
			t.Fatal(err)
		}
		if set.Flags&argvRootOnly == 0 || set.Lens != 0 {
			t.Errorf("%s: flags %#x lens %#x, want root-only and no argv rule", n, set.Flags, set.Lens)
		}
	}
}

func TestShortMaskRefusesNonLetters(t *testing.T) {
	if _, err := shortMask("F-"); err == nil {
		t.Fatal("a '-' got a bit")
	}
	if m, _ := shortMask("Aa0"); m != 1|1<<26|1<<52 {
		t.Fatalf("mask %#x", m)
	}
}
