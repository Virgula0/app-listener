package daemonconfig

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"
)

// LaunchRule is one [launch_rules] line, or a catalog entry's (install.CandidateDir.LaunchRules).
// Words: an argv rule, its words matched in order anywhere in the launch's argv (a flag word compared
// up to '=', see LaunchWordForms). Env instead: a variable whose presence matches. A match marks the
// launch code-suspect (a tripwire, never a boundary); Warn only logs it.
type LaunchRule struct {
	Path  string
	Words []string
	Env   string
	Warn  bool
	// Catalog only (no config syntax). RootOnly: a launch by any uid but 0 matches.
	RootOnly bool
	Short    *ShortOpts
}

// ShortOpts (catalog only): getopt letters refused anywhere in an option cluster before one of
// Args, the letters taking the rest of it as their value.
type ShortOpts struct{ Deny, Args string }

// MaxLaunchRuleWords, MaxLaunchWordLen, MaxLaunchRules and MaxAllowEnvLen bound rules to the
// kernel's (guard_trust.bpf.c ARGV_WORDS, LAUNCH_KEY with its NUL, ARGV_RULES argv rules per binary,
// ARGV_KEY with its NUL).
const (
	MaxLaunchRuleWords = 4
	MaxLaunchWordLen   = 43
	MaxLaunchRules     = 16
	MaxAllowEnvLen     = 63
)

// Same reports whether r and o are the same rule.
func (r LaunchRule) Same(o LaunchRule) bool {
	return r.Path == o.Path && r.Env == o.Env && r.Warn == o.Warn && slices.Equal(r.Words, o.Words)
}

// applyLaunchRule records `deny_args|warn_args <path> <word>...`, `deny_env|warn_env <path> <NAME>`
// or `allow_env NAME=/path`.
func applyLaunchRule(cfg *Config, line string, lineNo int) error {
	if kw, rest, _ := strings.Cut(line, " "); kw == "allow_env" {
		a, err := ParseAllowEnv(strings.TrimSpace(rest))
		if err != nil {
			return fmt.Errorf("daemon config line %d: %w", lineNo, err)
		}
		if slices.Contains(cfg.AllowEnv, a) {
			return fmt.Errorf("daemon config line %d: duplicate allow_env", lineNo)
		}
		cfg.AllowEnv = append(cfg.AllowEnv, a)
		return nil
	}
	r, err := parseLaunchRule(line)
	if err != nil {
		return fmt.Errorf("daemon config line %d: %w", lineNo, err)
	}
	if slices.ContainsFunc(cfg.LaunchRules, r.Same) {
		return fmt.Errorf("daemon config line %d: duplicate launch rule", lineNo)
	}
	argvRules := 0
	for _, o := range cfg.LaunchRules {
		if o.Path == r.Path && o.Env == "" {
			argvRules++
		}
	}
	if r.Env == "" && argvRules >= MaxLaunchRules {
		return fmt.Errorf("daemon config line %d: more than %d argv rules for %s", lineNo, MaxLaunchRules, r.Path)
	}
	cfg.LaunchRules = append(cfg.LaunchRules, r)
	return nil
}

func parseLaunchRule(line string) (LaunchRule, error) {
	kw, rest, _ := strings.Cut(line, " ")
	var r LaunchRule
	switch kw {
	case "deny_args", "warn_args", "deny_env", "warn_env":
		r.Warn = strings.HasPrefix(kw, "warn_")
	default:
		return r, fmt.Errorf("unknown [launch_rules] directive %q (deny_args, warn_args, deny_env, warn_env, "+
			"allow_env)", kw)
	}
	p, words, err := splitPathAndRest(strings.TrimSpace(rest))
	if err != nil {
		return r, err
	}
	if !filepath.IsAbs(p) {
		return r, fmt.Errorf("launch rule binary %q is not an absolute path", p)
	}
	r.Path = filepath.Clean(p)
	fields := strings.Fields(words)
	if strings.HasSuffix(kw, "_env") {
		if len(fields) != 1 || !envName(fields[0]) {
			return r, fmt.Errorf("%s names one environment variable, got %q", kw, words)
		}
		r.Env = fields[0]
		return r, nil
	}
	if len(fields) == 0 || len(fields) > MaxLaunchRuleWords {
		return r, fmt.Errorf("%s takes 1 to %d words, got %q", kw, MaxLaunchRuleWords, words)
	}
	for _, w := range fields {
		if _, err := LaunchWordForms(w); err != nil {
			return r, err
		}
	}
	r.Words = fields
	return r, nil
}

// LaunchWordForms lists the argv words a rule word matches: itself, or for a flag written
// `--name[rest]` every prefix from --name to --namerest. gpg (like getopt_long) takes any
// unambiguous abbreviation of a long option, so `--export-secret-ke` runs --export-secret-keys.
func LaunchWordForms(w string) ([]string, error) {
	base, rest, abbrev := strings.Cut(w, "[")
	tail, closed := strings.CutSuffix(rest, "]")
	switch {
	case len(base)+len(tail) > MaxLaunchWordLen:
		return nil, fmt.Errorf("launch rule word %q is longer than %d bytes", w, MaxLaunchWordLen)
	case strings.HasPrefix(w, "-") && strings.Contains(w, "="):
		return nil, fmt.Errorf("launch rule flag %q: a flag is compared up to '=', so name it without a value", w)
	case !abbrev:
		return []string{w}, nil
	case !closed || tail == "" || strings.ContainsAny(tail, "[]") || !strings.HasPrefix(base, "-") ||
		strings.TrimLeft(base, "-") == "":
		return nil, fmt.Errorf("launch rule word %q: only a flag abbreviates, written --name[rest]", w)
	}
	out := make([]string, 0, len(tail)+1)
	for i := range len(tail) + 1 {
		out = append(out, base+tail[:i])
	}
	return out, nil
}

// ParseAllowEnv checks an allow_env NAME=/path: one clean absolute path, no ':' (SSL_CERT_DIR reads
// a list), the whole string matched exactly by the kernel.
func ParseAllowEnv(s string) (string, error) {
	name, p, ok := strings.Cut(s, "=")
	switch {
	case !ok || !envName(name):
		return "", fmt.Errorf("allow_env takes NAME=/path, got %q", s)
	case !filepath.IsAbs(p) || filepath.Clean(p) != p || strings.ContainsAny(p, ": \t"):
		return "", fmt.Errorf("allow_env %q: the value must be one clean absolute path", s)
	case len(s) > MaxAllowEnvLen:
		return "", fmt.Errorf("allow_env %q is longer than %d bytes", s, MaxAllowEnvLen)
	}
	return s, nil
}

func envName(s string) bool {
	if s == "" || len(s) > MaxLaunchWordLen || (s[0] >= '0' && s[0] <= '9') {
		return false
	}
	for _, c := range s {
		if c != '_' && (c < 'A' || c > 'Z') && (c < 'a' || c > 'z') && (c < '0' || c > '9') {
			return false
		}
	}
	return true
}
