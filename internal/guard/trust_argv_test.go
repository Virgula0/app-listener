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

func deny(words ...string) daemonconfig.LaunchRule { return daemonconfig.LaunchRule{Words: words} }

func TestBuiltinAllowedEnvFits(t *testing.T) {
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
}

func TestEncodeLaunchSetSpellsOutAbbreviations(t *testing.T) {
	words, set, err := encodeLaunchSet([]daemonconfig.LaunchRule{deny("--sto[re]"), deny("x", "--ab[c]")}, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := map[argvEntry]uint64{
		{argvWord, "--sto"}: 1, {argvWord, "--stor"}: 1, {argvWord, "--store"}: 1,
		{argvWord, "x"}: 1 << 4, {argvWord, "--ab"}: 1 << 5, {argvWord, "--abc"}: 1 << 5,
	}
	if len(words) != len(want) || set.Lens != 1|2<<4 {
		t.Fatalf("words = %v, lens %#x", words, set.Lens)
	}
	for e, v := range want {
		if words[e] != v {
			t.Errorf("%v = %#x, want %#x", e, words[e], v)
		}
	}
	if _, _, err := encodeLaunchSet([]daemonconfig.LaunchRule{deny("sto[re]")}, nil); err == nil {
		t.Fatal("a non-flag abbreviation must be refused")
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

func TestShortMaskRefusesNonLetters(t *testing.T) {
	if _, err := shortMask("F-"); err == nil {
		t.Fatal("a '-' got a bit")
	}
	if m, _ := shortMask("Aa0"); m != 1|1<<26|1<<52 {
		t.Fatalf("mask %#x", m)
	}
}
