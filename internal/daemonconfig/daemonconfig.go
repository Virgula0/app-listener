// Package daemonconfig parses the daemon mode configuration file, following ssh-guard.conf
// grammar: each [watch <dir>] section protects one resource and lists allowed binaries with
// optional event types; chattr/exclude_chattr are deliberately unsupported (the LSM guard engine replaces them).
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

// Resource is one [watch <path>] section: a directory tree or a single regular file with
// its own fscrypt policy. Symlinks, hard-linked files and special files are refused at parse time.
type Resource struct {
	Path string
	// NeedEncryption selects the fscrypt lifecycle; defaults to true. An encrypted-marked
	// resource without an fscrypt policy aborts the daemon at startup.
	NeedEncryption bool
	Binaries       []BinaryRule
	// PendingBinaries parks whitelisted binaries unreadable at parse time (typically an
	// fscrypt-locked directory). Until a post-unlock pass moves them into Binaries they stay
	// absent from the BPF whitelist — denied: fail-closed, never fail-open.
	PendingBinaries []BinaryRule
}

// BinaryRule is one whitelisted binary inside a resource section.
type BinaryRule struct {
	Path string
	// Events restricts the operations this binary may perform; empty means every type.
	Events []ebpf.EventType
}

// Load parses the configuration file at path. Missing watch paths are skipped with a warning,
// as are directives outside any [watch] section; unreadable binaries are parked in
// PendingBinaries for post-unlock resolution. Malformed directives in valid sections fail fast.
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

		dirPath, isSection, headerErr := parseWatchSection(line)
		if headerErr != nil {
			return nil, fmt.Errorf("daemon config line %d: %w", lineNo, headerErr)
		}
		if isSection {
			current = addResource(cfg, dirPath, lineNo)
			continue
		}

		if current == nil {
			// Tolerated like ssh-guard, but warned: a silently dropped
			// security directive must be spottable.
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

func parseWatchSection(line string) (path string, isSection bool, err error) {
	if !strings.HasPrefix(line, "[") {
		return "", false, nil
	}
	// Bare "[watch]" is the installer's placeholder header (path filled in
	// later); accept it with an empty path — the empty resource is skipped
	// as unguardable by addResource.
	if line == "[watch]" {
		return "", true, nil
	}
	// Any other bracketed line is a section header: accepting malformed
	// ones as "directives" silently merged the following binaries into the
	// PREVIOUS section (a cross-resource whitelist contamination on manual
	// edits).
	if !strings.HasPrefix(line, "[watch ") || !strings.HasSuffix(line, "]") {
		return "", false, fmt.Errorf("malformed section header %q: expected \"[watch <path>]\"", line)
	}
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "[watch "), "]")), true, nil
}

// addResource creates a resource for watchPath, returning nil for unguardable targets:
// missing paths, symlinks, hard-linked regular files (the inode-based guard would
// implicitly cover another path) and anything not a directory or unique regular file.
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

// applyDirective handles one [watch]-section directive: need_encryption, binary rules
// and (deliberately unsupported) chattr lines.
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
		// Binary unreadable (locked tree or genuinely gone): defer to the post-unlock
		// pass rather than drop — dropping would silently disable the entry; deferred
		// stays unlisted in BPF, i.e. denied, until resolvable (fail-closed).
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
