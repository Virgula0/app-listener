// Package daemonconfig parses the daemon mode configuration file.
//
// The grammar follows ssh-guard.conf: each [watch <dir>] section protects
// one resource and lists the binaries allowed to access it, optionally
// restricted to a set of event types. chattr/exclude_chattr directives are
// intentionally NOT supported — the LSM guard engine replaces that
// mechanism entirely.
package daemonconfig

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// Config is a parsed daemon configuration.
type Config struct {
	Resources []Resource
}

// Resource is one [watch <path>] section. The path is either a directory
// (the classic fscrypt-encrypted resource tree) or a single regular file
// carrying its own fscrypt policy. Symbolic links, hard-linked files and
// special files are refused at parse time.
type Resource struct {
	Path string
	// NeedEncryption selects the fscrypt lifecycle for this resource.
	// It defaults to true; a resource marked as encrypted that carries no
	// fscrypt policy aborts the daemon at startup.
	NeedEncryption bool
	Binaries       []BinaryRule
	// PendingBinaries lists whitelisted binaries that could not be read
	// when the config was parsed (typically because their directory is
	// still fscrypt-locked). One pass after the resource is unlocked
	// (see the daemon usecase) moves them into Binaries; until then they
	// are absent from the BPF whitelist, so in whitelist mode they are
	// simply denied — fail-closed, never fail-open.
	PendingBinaries []BinaryRule
}

// BinaryRule is one whitelisted binary inside a resource section.
type BinaryRule struct {
	Path string
	// Events restricts the operations this binary may perform; an empty
	// list means every event type is allowed.
	Events []ebpf.EventType
}

// Load parses the configuration file at path. Missing watch paths are
// skipped with a warning (matching ssh-guard's tolerance), and so are the
// directives of a skipped section — including directives that appear
// before any [watch] section. A binary that cannot be read yet (its
// directory is still fscrypt-locked) is parked in Resource.PendingBinaries
// and resolved by the daemon once the resource is unlocked. Malformed
// directives inside a valid section fail fast: a security configuration
// must not be silently misread.
func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &Config{}
	var current *Resource

	scanner := bufio.NewScanner(file)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if dirPath, ok := parseWatchSection(line); ok {
			current = addResource(cfg, dirPath, lineNo)
			continue
		}

		if current == nil {
			// ssh-guard ignores directives outside a watch section (its
			// currentIdx < 0 branch). Keep that tolerance, but say so: a
			// silently dropped security directive is a misconfiguration
			// the operator must be able to spot.
			log.Warnf("daemon config line %d: ignoring directive outside any [watch] section: %q", lineNo, line)
			continue
		}

		if err := applyDirective(current, line, lineNo); err != nil {
			return nil, err
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	seen := make(map[string]bool, len(cfg.Resources))
	for _, r := range cfg.Resources {
		if seen[r.Path] {
			return nil, fmt.Errorf("daemon config: duplicate watch path: %s", r.Path)
		}
		seen[r.Path] = true
	}
	return cfg, nil
}

func parseWatchSection(line string) (string, bool) {
	if !strings.HasPrefix(line, "[watch ") || !strings.HasSuffix(line, "]") {
		return "", false
	}
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "[watch "), "]")), true
}

// addResource creates a new resource for watchPath, skipping targets that
// cannot be guarded: missing paths, symbolic links, hard-linked regular
// files (more than one link would make the guard implicitly cover another
// path to the same inode) and anything that is neither a directory nor a
// unique regular file. It returns nil when the path is skipped.
func addResource(cfg *Config, watchPath string, lineNo int) *Resource {
	info, err := os.Lstat(watchPath)
	if err != nil {
		log.Warnf("daemon config line %d: skipping missing or unreadable path: %s", lineNo, watchPath)
		return nil
	}
	if info.Mode()&os.ModeSymlink != 0 {
		log.Warnf("daemon config line %d: skipping %s: symbolic links are refused as watch targets (guard identity is inode based)", lineNo, watchPath)
		return nil
	}
	if !info.IsDir() {
		if !info.Mode().IsRegular() {
			log.Warnf("daemon config line %d: skipping %s: only directories and regular files can be watched", lineNo, watchPath)
			return nil
		}
		if stat, ok := info.Sys().(*syscall.Stat_t); ok && stat.Nlink > 1 {
			log.Warnf("daemon config line %d: skipping %s: hard-linked files (%d links) are refused as watch targets", lineNo, watchPath, stat.Nlink)
			return nil
		}
	}
	cfg.Resources = append(cfg.Resources, Resource{Path: watchPath, NeedEncryption: true})
	return &cfg.Resources[len(cfg.Resources)-1]
}

// applyDirective handles one directive inside a [watch] section:
// need_encryption lines, binary lines and (deliberately unsupported)
// chattr lines.
func applyDirective(current *Resource, line string, lineNo int) error {
	if value, ok := parseNeedEncryption(line); ok {
		switch value {
		case "true":
			current.NeedEncryption = true
		case "false":
			current.NeedEncryption = false
		default:
			return fmt.Errorf("daemon config line %d: invalid need_encryption value %q (expected true or false)", lineNo, value)
		}
		return nil
	}

	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil
	}
	binPath := fields[0]
	rule := BinaryRule{Path: binPath}
	if len(fields) > 1 {
		events, parseErr := parseEvents(strings.Split(strings.Join(fields[1:], " "), ","))
		if parseErr != nil {
			return fmt.Errorf("daemon config line %d: %w", lineNo, parseErr)
		}
		rule.Events = events
	}
	if _, statErr := os.Stat(binPath); statErr != nil {
		// The binary lives inside a still-locked fscrypt tree (or is
		// genuinely gone). Keep the rule for the post-unlock resolution
		// pass instead of dropping it: dropping would silently disable
		// the whitelist entry, and keeping it is fail-closed (the binary
		// stays unlisted in the BPF whitelist, i.e. denied, until it
		// becomes resolvable).
		log.Warnf("daemon config line %d: binary not readable yet, deferring: %s", lineNo, binPath)
		current.PendingBinaries = append(current.PendingBinaries, rule)
		return nil //nolint:nilerr // returning nil despite statErr is the point: the unreadable rule is deliberately deferred, not dropped
	}
	if resolved, resolveErr := filepath.EvalSymlinks(binPath); resolveErr == nil && resolved != binPath {
		log.Infof("daemon config line %d: binary symlink resolved: %s -> %s", lineNo, binPath, resolved)
		rule.Path = resolved
	}
	current.Binaries = append(current.Binaries, rule)
	return nil
}

func parseNeedEncryption(line string) (string, bool) {
	if !strings.HasPrefix(line, "need_encryption:") {
		return "", false
	}
	return strings.TrimSpace(strings.TrimPrefix(line, "need_encryption:")), true
}

func parseEvents(tokens []string) ([]ebpf.EventType, error) {
	var types []ebpf.EventType
	seen := make(map[ebpf.EventType]bool)
	for _, token := range tokens {
		et, ok := ebpf.ParseEventType(strings.TrimSpace(token))
		if !ok {
			return nil, fmt.Errorf("unknown event type %q (valid: OPEN, READ, WRITE, DELETE, RENAME, SYMLINK, HARDLINK, MKDIR, MMAP, ATTR, STAT, MKNOD)", token)
		}
		if !seen[et] {
			seen[et] = true
			types = append(types, et)
		}
	}
	return types, nil
}
