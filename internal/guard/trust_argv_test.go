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
		if j == len(word)-1 {
			return false
		}
	}
	return true
}

func TestShortMaskRefusesNonLetters(t *testing.T) {
	if _, err := shortMask("F-"); err == nil {
		t.Fatal("a '-' got a bit")
	}
	if m, _ := shortMask("Aa0"); m != 1|1<<26|1<<52 {
		t.Fatalf("mask %#x", m)
	}
}

// Anchors are set flags and names, not argv rules: a binary with every argv rule slot used keeps them.
func TestEncodeLaunchSetProfileAnchors(t *testing.T) {
	rules := []daemonconfig.LaunchRule{{Profile: "SingletonLock"}, {Profile: "Local State"}}
	for i := range argvRules {
		rules = append(rules, deny(string(rune('a'+i))))
	}
	canon, allowed, key := canonLaunchRules(rules, nil)
	words, set, err := encodeLaunchSet(canon, allowed)
	if err != nil {
		t.Fatal(err)
	}
	if set.Flags&argvAnchored == 0 || words[argvEntry{argvProfile, "SingletonLock"}] != 1 ||
		words[argvEntry{argvProfile, "Local State"}] != 1 {
		t.Fatalf("anchors not encoded: flags %#x words %v", set.Flags, words)
	}
	if _, _, plain := canonLaunchRules(rules[2:], nil); plain == key {
		t.Fatal("a set's key must name its anchors")
	}
	if _, set, _ := encodeLaunchSet(rules[2:], nil); set.Flags&argvAnchored != 0 {
		t.Fatal("a set without anchors must not be flagged")
	}
}

// wordMask mirrors guard_trust.bpf.c argv_word's lookups for one argv word.
func wordMask(words map[argvEntry]uint64, set GuardTrustArgvSet, w string) uint64 {
	if strings.HasPrefix(w, "-") {
		w, _, _ = strings.Cut(w, "=")
	}
	m := words[argvEntry{argvWord, w}]
	if set.Flags&argvFolded != 0 {
		m |= words[argvEntry{argvWordFold, asciiLower(w)}]
	}
	return m
}

// One binary may mix a parser that ignores case (Firefox's flags) with one that doesn't (its
// -MOZ_LOG): each rule says how its words compare, and either dash count where it says so.
func TestEncodeLaunchSetFoldsPerRule(t *testing.T) {
	words, set, err := encodeLaunchSet([]daemonconfig.LaunchRule{
		{Words: []string{"-headless"}, FoldCase: true, EitherDash: true},
		{Words: []string{"-MOZ_LOG"}, EitherDash: true},
		deny("--exact"),
		{Words: []string{"--ab[c]"}, FoldCase: true},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if set.Flags&argvFolded == 0 {
		t.Fatal("a set with a folded rule must be flagged for the folded lookup")
	}
	for w, want := range map[string]uint64{
		"-headless": 1, "--headless": 1, "-HEADLESS": 1, "--HeadLess": 1, "-headless=1": 1, "---headless": 0,
		"-MOZ_LOG": 1 << 4, "--MOZ_LOG=x": 1 << 4, "-moz_log": 0, "--Moz_Log": 0,
		"--exact": 1 << 8, "-exact": 0, "--EXACT": 0,
		"--AB": 1 << 12, "--aBc": 1 << 12, "-ab": 0,
	} {
		if got := wordMask(words, set, w); got != want {
			t.Errorf("%q matches %#x, want %#x", w, got, want)
		}
	}
	for e := range words {
		if e.kind == argvWordFold && e.name != asciiLower(e.name) {
			t.Errorf("folded key %q is not lowercase: the kernel looks it up lowercased", e.name)
		}
	}
	_, plain, _ := encodeLaunchSet([]daemonconfig.LaunchRule{deny("-headless")}, nil)
	if plain.Flags&argvFolded != 0 {
		t.Error("a set without folded rules must not pay for the folded lookup")
	}
	a := []daemonconfig.LaunchRule{{Words: []string{"-x"}, FoldCase: true}}
	b := []daemonconfig.LaunchRule{{Words: []string{"-x"}}}
	if _, _, ka := canonLaunchRules(a, nil); func() bool { _, _, kb := canonLaunchRules(b, nil); return ka == kb }() {
		t.Error("a folded and an exact rule must not share a set")
	}
}
