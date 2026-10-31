package guard

import (
	"fmt"
	"slices"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/logging"
)

// trust_argv kinds, values and limits (guard_trust.bpf.c).
const (
	argvWord    uint32 = 1
	argvEnv     uint32 = 2
	argvEnvOK   uint32 = 3
	argvDeny    uint64 = 1
	argvWarn    uint64 = 2
	argvRules          = daemonconfig.MaxLaunchRules
	argvWords          = daemonconfig.MaxLaunchRuleWords
	argvSetsMax        = 1024
)

func deny(words ...string) daemonconfig.LaunchRule { return daemonconfig.LaunchRule{Words: words} }

// builtinLaunchRules: CLIs that print or spend their secret on request, keyed by the whitelist line's
// basename. A rule only adds denials, so picking it by name weakens nothing; identity stays the inode.
var builtinLaunchRules = func() map[string][]daemonconfig.LaunchRule {
	each := func(words ...string) []daemonconfig.LaunchRule {
		out := make([]daemonconfig.LaunchRule, 0, len(words))
		for _, w := range words {
			out = append(out, deny(w))
		}
		return out
	}
	// -stdin carries the method on stdin.
	coin := each("-stdin", "dumpprivkey", "dumpwallet", "walletpassphrase", "send", "sendall", "sendmany",
		"sendtoaddress")
	lightning := each("pay", "xpay", "renepay", "keysend", "withdraw", "sendpay", "sendonion", "fundchannel",
		"multifundchannel", "close", "exposesecret")
	return map[string][]daemonconfig.LaunchRule{
		"gh": {
			deny("auth", "token"), deny("auth", "status", "--show-token"), deny("auth", "status", "-t"),
			deny("alias", "set"), deny("alias", "import"),
			// git's credential helper: denying it breaks every git push, so it is only logged.
			{Words: []string{"auth", "git-credential"}, Warn: true},
		},
		"aws":           {deny("configure", "export-credentials"), deny("configure", "get")},
		"bitcoin-cli":   coin,
		"litecoin-cli":  coin,
		"lightning-cli": lightning,
	}
}()

// caEnvRules: Go's crypto/x509, OpenSSL and GnuTLS trust SSL_CERT_FILE / SSL_CERT_DIR, so a caller
// setting either makes a whitelisted client hand its token to a TLS server it chose. Every whitelisted
// binary gets them; a value naming a root-placed bundle (SetLaunchRules' allowEnv) passes.
var caEnvRules = []daemonconfig.LaunchRule{{Env: "SSL_CERT_FILE"}, {Env: "SSL_CERT_DIR"}}

// builtinAllowedEnv: the distributions' system CA bundles.
var builtinAllowedEnv = []string{
	"SSL_CERT_FILE=/etc/ssl/certs/ca-certificates.crt",
	"SSL_CERT_FILE=/etc/pki/tls/certs/ca-bundle.crt",
	"SSL_CERT_FILE=/etc/ssl/ca-bundle.pem",
	"SSL_CERT_FILE=/etc/ssl/cert.pem",
	"SSL_CERT_FILE=/etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem",
	"SSL_CERT_FILE=/etc/ca-certificates/extracted/tls-ca-bundle.pem",
	"SSL_CERT_DIR=/etc/ssl/certs",
	"SSL_CERT_DIR=/etc/pki/tls/certs",
}

// BuiltinLaunchRules returns a copy of the built-in rules for a whitelisted binary named name.
func BuiltinLaunchRules(name string) []daemonconfig.LaunchRule {
	return slices.Clone(builtinLaunchRules[name])
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

// encodeLaunchSet lays rules out as trust_argv entries and their trust_argv_sets row: argv rule i's
// word j sets bit i*argvWords+j of that word's entry; an env name maps to argvDeny, or argvWarn when
// every rule naming it only warns; each allowed NAME=value is an argvEnvOK entry. rules must be
// canonical (canonLaunchRules).
func encodeLaunchSet(rules []daemonconfig.LaunchRule, allowed []string) (map[argvEntry]uint64,
	GuardTrustArgvSet, error) {
	words := map[argvEntry]uint64{}
	for _, a := range allowed {
		words[argvEntry{argvEnvOK, a}] = 1
	}
	var set GuardTrustArgvSet
	i := 0
	for _, r := range rules {
		if r.Env != "" {
			switch e := (argvEntry{argvEnv, r.Env}); {
			case !r.Warn:
				words[e] = argvDeny
			case words[e] == 0:
				words[e] = argvWarn
			}
			continue
		}
		if i >= argvRules {
			return nil, set, fmt.Errorf("more than %d argv rules", argvRules)
		}
		if len(r.Words) == 0 || len(r.Words) > argvWords {
			return nil, set, fmt.Errorf("argv rule %q: 1 to %d words", strings.Join(r.Words, " "), argvWords)
		}
		for j, w := range r.Words {
			words[argvEntry{argvWord, w}] |= 1 << (i*argvWords + j)
		}
		set.Lens |= uint64(len(r.Words)) << (i * 4)
		if r.Warn {
			set.Warn |= 1 << i
		}
		i++
	}
	return words, set, nil
}

// canonLaunchRules dedupes and orders rules and the allowed NAME=value strings of the variables they
// name, returning both with the key naming their set.
func canonLaunchRules(rules []daemonconfig.LaunchRule, allowEnv []string) (canon []daemonconfig.LaunchRule,
	allowed []string, key string) {
	enc := func(r daemonconfig.LaunchRule) string {
		return fmt.Sprintf("%s\x00%t\x00%s", r.Env, r.Warn, strings.Join(r.Words, "\x00"))
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

// launchSet returns the trust_argv set of the binary at path, writing a new set first; ok false: its
// rules don't fit. Sets are never rewritten nor reused, so a row always names the rules it was given.
func (t *TrustGuard) launchSet(path string) (id uint64, ok bool, err error) {
	t.argvMu.Lock()
	defer t.argvMu.Unlock()
	rules, allowed, canon := canonLaunchRules(append(slices.Clone(t.argvByPath[path]), caEnvRules...), t.allowEnv)
	if known, ok := t.argvSets[canon]; ok {
		return known, true, nil
	}
	words, set, err := encodeLaunchSet(rules, allowed)
	if err != nil {
		log.Errorf("trust guard: launch rules of %s not applied: %v", logging.SanitizeText(path), err)
		return 0, false, nil
	}
	if len(t.argvSets) >= argvSetsMax {
		log.Errorf("trust guard: launch rules of %s not applied: more than %d distinct rule sets since the "+
			"daemon started; restart it", logging.SanitizeText(path), argvSetsMax)
		return 0, false, nil
	}
	id = uint64(len(t.argvSets))
	for e, v := range words {
		k := GuardTrustArgvKey{Set: uint32(id), Kind: e.kind} //nolint:gosec // id < argvSetsMax
		if len(e.name) >= len(k.Name) {
			return 0, false, fmt.Errorf("launch rule word %q does not fit the %d-byte key", e.name, len(k.Name))
		}
		for i := range len(e.name) {
			k.Name[i] = int8(e.name[i]) //nolint:gosec // byte reinterpreted as a C char
		}
		if err := t.objs.TrustArgv.Put(k, v); err != nil {
			return 0, false, fmt.Errorf("launch rule word %s: %w", e.name, err)
		}
	}
	if err := t.objs.TrustArgvSets.Put(uint32(id), set); err != nil { //nolint:gosec // id < argvSetsMax
		return 0, false, fmt.Errorf("launch rule set: %w", err)
	}
	if t.argvSets == nil {
		t.argvSets = map[string]uint64{}
	}
	t.argvSets[canon] = id
	return id, true, nil
}
