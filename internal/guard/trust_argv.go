package guard

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/logging"
)

// trust_argv kinds, values and limits (guard_trust.bpf.c).
const (
	argvWord     uint32 = 1
	argvEnv      uint32 = 2
	argvEnvOK    uint32 = 3
	argvProfile  uint32 = 4
	argvWordFold uint32 = 5
	argvDeny     uint64 = 1
	argvWarn     uint64 = 2
	argvRules           = daemonconfig.MaxLaunchRules
	argvWords           = daemonconfig.MaxLaunchRuleWords
	argvSetsMax         = 1024
	argvRootOnly        = 1
	argvAnchored        = 2
	argvFolded          = 4
)

// caEnvRules: Go's crypto/x509, OpenSSL and GnuTLS trust SSL_CERT_FILE / SSL_CERT_DIR, Node and Bun
// NODE_EXTRA_CA_CERTS, and NODE_TLS_REJECT_UNAUTHORIZED=0 trusts any certificate: a caller setting
// one makes a whitelisted client hand its token to a TLS server it chose. Every whitelisted binary
// gets them; a value naming a root-placed bundle (SetLaunchRules' allowEnv) passes.
var caEnvRules = []daemonconfig.LaunchRule{{Env: "SSL_CERT_FILE"}, {Env: "SSL_CERT_DIR"},
	{Env: "NODE_EXTRA_CA_CERTS"}, {Env: "NODE_TLS_REJECT_UNAUTHORIZED"}}

// builtinAllowedEnv: the distributions' system CA bundles (NODE_EXTRA_CA_CERTS only where the whole
// string fits the 63-byte key).
var builtinAllowedEnv = []string{
	"NODE_EXTRA_CA_CERTS=/etc/ssl/certs/ca-certificates.crt",
	"NODE_EXTRA_CA_CERTS=/etc/pki/tls/certs/ca-bundle.crt",
	"NODE_EXTRA_CA_CERTS=/etc/ssl/ca-bundle.pem",
	"NODE_EXTRA_CA_CERTS=/etc/ssl/cert.pem",
	"SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt",
	"SSL_CERT_FILE=/etc/pki/tls/certs/ca-bundle.crt",
	"SSL_CERT_FILE=/etc/ssl/ca-bundle.pem",
	"SSL_CERT_FILE=/etc/ssl/cert.pem",
	"SSL_CERT_FILE=/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem",
	"SSL_CERT_FILE=/etc/ca-certificates/extracted/tls-ca-bundle.pem",
	"SSL_CERT_DIR=/etc/ssl/certs",
	"SSL_CERT_DIR=/etc/pki/tls/certs",
}

// BuiltinAllowedEnv returns a copy of the built-in NAME=value strings no env rule matches; the
// caller vets each path as it does allow_env lines.
func BuiltinAllowedEnv() []string {
	return slices.Clone(builtinAllowedEnv)
}

type argvEntry struct {
	kind uint32
	name string
}

// encodeLaunchSet lays rules out as trust_argv entries and their trust_argv_sets row (see
// encodeArgvRule); an env name maps to argvDeny, or argvWarn when every rule naming it only warns;
// each allowed NAME=value is an argvEnvOK entry. rules must be canonical (canonLaunchRules).
func encodeLaunchSet(rules []daemonconfig.LaunchRule, allowed []string) (map[argvEntry]uint64,
	GuardTrustArgvSet, error) {
	words := map[argvEntry]uint64{}
	for _, a := range allowed {
		words[argvEntry{argvEnvOK, a}] = 1
	}
	var set GuardTrustArgvSet
	i := 0
	for _, r := range rules {
		var err error
		switch {
		case r.RootOnly || r.Short != nil:
			err = encodeSetFlags(&set, &r)
		case r.Profile != "":
			words[argvEntry{argvProfile, r.Profile}] = 1
			set.Flags |= argvAnchored
		case r.Env != "":
			encodeEnvRule(words, &r)
		default:
			err = encodeArgvRule(words, &set, &r, i)
			i++
		}
		if err != nil {
			return nil, set, err
		}
	}
	return words, set, nil
}

// encodeArgvRule records argv rule i: its word j sets bit i*argvWords+j of each of that word's
// forms' entries. An abbreviation's every spelling is an exact key, so the kernel needs no prefix
// match; a FoldCase rule's are lowercase argvWordFold keys, which the kernel looks up lowercased.
func encodeArgvRule(words map[argvEntry]uint64, set *GuardTrustArgvSet, r *daemonconfig.LaunchRule, i int) error {
	if i >= argvRules {
		return fmt.Errorf("more than %d argv rules", argvRules)
	}
	if len(r.Words) == 0 || len(r.Words) > argvWords {
		return fmt.Errorf("argv rule %q: 1 to %d words", strings.Join(r.Words, " "), argvWords)
	}
	for j, w := range r.Words {
		forms, err := daemonconfig.LaunchWordForms(w)
		if err != nil {
			return err
		}
		kind := argvWord
		if r.FoldCase {
			kind = argvWordFold
			set.Flags |= argvFolded
		}
		for _, f := range ruleForms(r, forms) {
			words[argvEntry{kind, f}] |= 1 << (i*argvWords + j)
		}
	}
	set.Lens |= uint64(len(r.Words)) << (i * 4)
	if r.Warn {
		set.Warn |= 1 << i
	}
	return nil
}

// ruleForms adds forms' other-dash spellings for an EitherDash rule and lowercases a FoldCase one's.
func ruleForms(r *daemonconfig.LaunchRule, forms []string) []string {
	out := slices.Clone(forms)
	if r.EitherDash {
		for _, f := range forms {
			switch {
			case strings.HasPrefix(f, "--"):
				out = append(out, f[1:])
			case strings.HasPrefix(f, "-") && len(f) > 1:
				out = append(out, "-"+f)
			}
		}
	}
	if r.FoldCase {
		for k, f := range out {
			out[k] = asciiLower(f)
		}
	}
	return out
}

// asciiLower folds A-Z only, as the kernel does: no Unicode case mapping.
func asciiLower(s string) string {
	b := []byte(s)
	for k, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[k] = c + 'a' - 'A'
		}
	}
	return string(b)
}

// encodeEnvRule maps an env rule's name to argvDeny, or argvWarn while every rule naming it warns.
func encodeEnvRule(words map[argvEntry]uint64, r *daemonconfig.LaunchRule) {
	switch e := (argvEntry{argvEnv, r.Env}); {
	case !r.Warn:
		words[e] = argvDeny
	case words[e] == 0:
		words[e] = argvWarn
	}
}

// encodeSetFlags records a RootOnly or Short rule in set.
func encodeSetFlags(set *GuardTrustArgvSet, r *daemonconfig.LaunchRule) error {
	if r.RootOnly {
		set.Flags |= argvRootOnly
	}
	if r.Short == nil {
		return nil
	}
	deny, err := shortMask(r.Short.Deny)
	if err != nil {
		return err
	}
	args, err := shortMask(r.Short.Args)
	if err != nil {
		return err
	}
	set.ShortDeny |= deny
	set.ShortArg |= args
	return nil
}

// shortMask maps getopt letters to argv_set's short_* bits: A-Z, a-z, 0-9 = bits 0-61.
func shortMask(letters string) (uint64, error) {
	var m uint64
	for _, c := range letters {
		switch {
		case c >= 'A' && c <= 'Z':
			m |= 1 << (c - 'A')
		case c >= 'a' && c <= 'z':
			m |= 1 << (26 + c - 'a')
		case c >= '0' && c <= '9':
			m |= 1 << (52 + c - '0')
		default:
			return 0, fmt.Errorf("short option %q: not a letter or digit", c)
		}
	}
	return m, nil
}

// canonLaunchRules dedupes and orders rules and the allowed NAME=value strings of the variables they
// name, returning both with the key naming their set.
func canonLaunchRules(rules []daemonconfig.LaunchRule, allowEnv []string) (canon []daemonconfig.LaunchRule,
	allowed []string, key string) {
	enc := func(r daemonconfig.LaunchRule) string {
		short := ""
		if r.Short != nil {
			short = r.Short.Deny + "/" + r.Short.Args
		}
		return fmt.Sprintf("%s\x00%t\x00%t\x00%s\x00%s\x00%t\x00%t\x00%s", r.Env, r.Warn, r.RootOnly, short,
			r.Profile, r.FoldCase, r.EitherDash, strings.Join(r.Words, "\x00"))
	}
	canon = slices.Clone(rules)
	slices.SortFunc(canon, func(a, b daemonconfig.LaunchRule) int { return strings.Compare(enc(a), enc(b)) })
	canon = slices.CompactFunc(canon, func(a, b daemonconfig.LaunchRule) bool { return enc(a) == enc(b) })
	keys := make([]string, 0, len(canon)+len(allowEnv))
	for _, r := range canon {
		keys = append(keys, enc(r))
	}
	for _, a := range allowEnv {
		name, _, _ := strings.Cut(a, "=")
		if slices.ContainsFunc(canon, func(r daemonconfig.LaunchRule) bool { return r.Env == name }) {
			allowed = append(allowed, a)
		}
	}
	slices.Sort(allowed)
	allowed = slices.Compact(allowed)
	for _, a := range allowed {
		keys = append(keys, "\x02"+a)
	}
	return canon, allowed, strings.Join(keys, "\x01")
}

// SetLaunchRules sets each whitelisted binary's argv/env rules by resolved path, applied by the next
// SetTrusted and by live admissions; every binary also gets caEnvRules. allowEnv: the vetted
// NAME=value strings no env rule matches.
func (t *TrustGuard) SetLaunchRules(byPath map[string][]daemonconfig.LaunchRule, allowEnv []string) {
	t.argvMu.Lock()
	defer t.argvMu.Unlock()
	t.argvByPath = byPath
	t.allowEnv = allowEnv
}

// CheckLaunchRules reports whether rules (plus caEnvRules) fit one trust_argv set.
func CheckLaunchRules(rules []daemonconfig.LaunchRule) error {
	canon, _, _ := canonLaunchRules(append(slices.Clone(rules), caEnvRules...), nil)
	words, _, err := encodeLaunchSet(canon, nil)
	if err != nil {
		return err
	}
	var k GuardTrustArgvKey
	for e := range words {
		if len(e.name) >= len(k.Name) {
			return fmt.Errorf("launch rule word %q does not fit the %d-byte key", e.name, len(k.Name))
		}
	}
	return nil
}

// launchSet returns the trust_argv set of the binary at path, writing a new set first. Rules that
// don't fit are an error, never a binary admitted without them. Sets are never rewritten nor
// reused, so a row always names the rules it was given.
func (t *TrustGuard) launchSet(path string) (uint64, error) {
	t.argvMu.Lock()
	defer t.argvMu.Unlock()
	rules, allowed, canon := canonLaunchRules(append(slices.Clone(t.argvByPath[path]), caEnvRules...), t.allowEnv)
	if known, ok := t.argvSets[canon]; ok {
		return known, nil
	}
	words, set, err := encodeLaunchSet(rules, allowed)
	if err != nil {
		return 0, fmt.Errorf("launch rules of %s: %w", logging.SanitizeText(path), err)
	}
	if len(t.argvSets) >= argvSetsMax {
		return 0, fmt.Errorf("launch rules of %s: more than %d distinct rule sets since the daemon "+
			"started; restart it", logging.SanitizeText(path), argvSetsMax)
	}
	id := uint64(len(t.argvSets))
	for e, v := range words {
		k := GuardTrustArgvKey{Set: uint32(id), Kind: e.kind} //nolint:gosec // id < argvSetsMax
		if len(e.name) >= len(k.Name) {
			return 0, fmt.Errorf("launch rule word %q does not fit the %d-byte key", e.name, len(k.Name))
		}
		for i := range len(e.name) {
			k.Name[i] = int8(e.name[i]) //nolint:gosec // byte reinterpreted as a C char
		}
		if err := t.objs.TrustArgv.Put(k, v); err != nil {
			return 0, fmt.Errorf("launch rule word %s: %w", e.name, err)
		}
	}
	if err := t.objs.TrustArgvSets.Put(uint32(id), set); err != nil { //nolint:gosec // id < argvSetsMax
		return 0, fmt.Errorf("launch rule set: %w", err)
	}
	if t.argvSets == nil {
		t.argvSets = map[string]uint64{}
	}
	t.argvSets[canon] = id
	return id, nil
}
