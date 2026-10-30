package daemonconfig

import (
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// The Secret Service gatekeeper (internal/secretgate). [secret_service] names the keyring daemon
// (provider), the prompter it shows dialogs through and the managers that may administer the whole
// keyring; each [secret <name>] scope names the items it protects (`match` lines) and the binaries that
// may use them. A scope is a single-file resource rooted at a root-owned anchor the daemon makes, so
// its whitelist is resolved, re-synced and reloaded like any other; the gatekeeper judges a bus caller
// against it. The provider reaches the user bus only through the daemon's relay socket in
// SecretRelayDir, a resource whose only CONNECT binary is the provider.
const (
	SecretAnchorDir = "/run/app-listener/secret-scopes" //nolint:gosec // a path, not a credential
	SecretRelayDir  = "/run/app-listener/secret-relay"  //nolint:gosec // a path, not a credential
)

// Anchored reports whether r is a daemon-made root (a port's or secret scope's anchor, the secret relay
// dir) rather than a user's tree: never encrypted, and a reload may drop it.
func (r *Resource) Anchored() bool { return r.Port != 0 || r.Secret != nil || r.SecretRelay }

// SecretAnchor is scope name's anchor file.
func SecretAnchor(name string) string { return SecretAnchorDir + "/" + name }

// SecretRelaySocket is the socket uid's provider connects to (its DBUS_SESSION_BUS_ADDRESS).
func SecretRelaySocket(uid uint32) string { return fmt.Sprintf("%s/%d.bus", SecretRelayDir, uid) }

// SecretService is the [secret_service] block.
type SecretService struct {
	Providers []string
	Prompters []string
	Managers  []string
}

// SecretMatch is one `match` line: every attribute must equal its value, or start with it when the
// value ends in '*'.
type SecretMatch map[string]string

// Matches reports whether an item with attrs falls under m.
func (m SecretMatch) Matches(attrs map[string]string) bool {
	for k, want := range m {
		got, ok := attrs[k]
		if !ok {
			return false
		}
		if prefix, glob := strings.CutSuffix(want, "*"); glob {
			if !strings.HasPrefix(got, prefix) {
				return false
			}
		} else if got != want {
			return false
		}
	}
	return true
}

// SecretScope is a `[secret <name>]` section: an item belongs to it when any Matches entry does.
type SecretScope struct {
	Name    string
	Matches []SecretMatch
}

var secretNameRE = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,63}$`)

// ParseSecretSection recognizes a `[secret <name>]` header: ok false for any other line.
func ParseSecretSection(line string) (name string, ok bool, err error) {
	if !strings.HasPrefix(line, "[secret ") || !strings.HasSuffix(line, "]") {
		return "", false, nil
	}
	name = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "[secret "), "]"))
	if !secretNameRE.MatchString(name) {
		return "", true, fmt.Errorf("malformed section header %q: a secret scope name is lowercase letters, "+
			"digits, '.', '_' and '-'", line)
	}
	return name, true, nil
}

// applySecretService records one [secret_service] line: `provider|prompter|manager <absolute path>`.
func applySecretService(cfg *Config, line string, lineNo int) error {
	kw, rest, _ := strings.Cut(line, " ")
	p, extra, err := splitPathAndRest(strings.TrimSpace(rest))
	if err != nil {
		return fmt.Errorf("daemon config line %d: %w", lineNo, err)
	}
	if extra != "" || !filepath.IsAbs(p) {
		return fmt.Errorf("daemon config line %d: a [secret_service] line is `provider|prompter|manager "+
			"<absolute path>`, got %q", lineNo, line)
	}
	p = filepath.Clean(p)
	var list *[]string
	switch kw {
	case "provider":
		list = &cfg.SecretService.Providers
	case "prompter":
		list = &cfg.SecretService.Prompters
	case "manager":
		list = &cfg.SecretService.Managers
	default:
		return fmt.Errorf("daemon config line %d: unknown [secret_service] directive %q "+
			"(provider, prompter, manager)", lineNo, kw)
	}
	if slices.Contains(*list, p) {
		return fmt.Errorf("daemon config line %d: duplicate %s %s", lineNo, kw, p)
	}
	*list = append(*list, p)
	return nil
}

// parseSecretMatch reads `match key=value...`; ok false for any other line.
func parseSecretMatch(line string) (SecretMatch, bool, error) {
	rest, ok := strings.CutPrefix(line, "match ")
	if !ok {
		return nil, false, nil
	}
	m := SecretMatch{}
	for pair := range strings.FieldsSeq(rest) {
		k, v, found := strings.Cut(pair, "=")
		if !found || k == "" || v == "" || v == "*" {
			return nil, true, fmt.Errorf("match: %q is not key=value (a bare * would match every item)", pair)
		}
		if strings.Contains(strings.TrimSuffix(v, "*"), "*") {
			return nil, true, fmt.Errorf("match: %q: '*' only ends a value", pair)
		}
		if _, dup := m[k]; dup {
			return nil, true, fmt.Errorf("match: duplicate attribute %q", k)
		}
		m[k] = v
	}
	if len(m) == 0 {
		return nil, true, fmt.Errorf("match: no attributes")
	}
	return m, true, nil
}

// materializeSecretGroup appends the scope's Resource; its anchor is made by the daemon, like a port's.
func materializeSecretGroup(cfg *Config, g *watchGroup) error {
	if g.runtime || len(g.watchPaths) > 0 || len(g.libs) > 0 || len(g.pendingLibs) > 0 ||
		len(g.libDirs) > 0 || len(g.libBinaries) > 0 || len(g.libPending) > 0 || g.needEncryption {
		return fmt.Errorf("daemon config line %d: a secret section takes only `match` and binary lines", g.lineNo)
	}
	if len(g.secretMatches) == 0 {
		return fmt.Errorf("daemon config line %d: secret section %q has no `match` line", g.lineNo, g.secret)
	}
	if err := checkSecretEvents(g); err != nil {
		return err
	}
	for i := range cfg.Resources {
		if s := cfg.Resources[i].Secret; s != nil && s.Name == g.secret {
			return fmt.Errorf("daemon config line %d: duplicate secret section %q", g.lineNo, g.secret)
		}
	}
	cfg.Resources = append(cfg.Resources, Resource{
		Path:            SecretAnchor(g.secret),
		Secret:          &SecretScope{Name: g.secret, Matches: g.secretMatches},
		Binaries:        g.binaries,
		PendingBinaries: g.pending,
	})
	return nil
}

func checkSecretEvents(g *watchGroup) error {
	for _, b := range slices.Concat(g.binaries, g.pending) {
		for _, e := range b.Events {
			if e != ebpf.EventRead && e != ebpf.EventWrite && e != ebpf.EventDelete {
				return fmt.Errorf("daemon config line %d: %s %s: a secret section's events are READ, "+
					"WRITE and DELETE", g.lineNo, b.Path, e)
			}
		}
	}
	return nil
}

// finishSecretService checks the scopes have a provider and adds the relay resource: only a provider
// may connect to its sockets.
func finishSecretService(cfg *Config) error {
	scopes := slices.ContainsFunc(cfg.Resources, func(r Resource) bool { return r.Secret != nil })
	svc := &cfg.SecretService
	if len(svc.Providers) == 0 {
		if scopes || len(svc.Prompters) > 0 || len(svc.Managers) > 0 {
			return fmt.Errorf("daemon config: secret sections need a [secret_service] block with a provider line")
		}
		return nil
	}
	rules := make([]BinaryRule, 0, len(svc.Providers))
	pending := []BinaryRule{}
	for _, p := range svc.Providers {
		recordBinaryRule(BinaryRule{Path: p, Events: []ebpf.EventType{ebpf.EventConnect, ebpf.EventStat}},
			&rules, &pending, 0)
	}
	cfg.Resources = append(cfg.Resources, Resource{
		Path:            SecretRelayDir,
		SecretRelay:     true,
		Binaries:        rules,
		PendingBinaries: pending,
	})
	return nil
}
