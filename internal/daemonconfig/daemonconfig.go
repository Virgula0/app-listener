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
	"strings"

	log "github.com/sirupsen/logrus"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// Config is a parsed daemon configuration.
type Config struct {
	Resources []Resource
}

// Resource is one [watch <dir>] section.
type Resource struct {
	Path string
	// NeedEncryption selects the fscrypt lifecycle for this resource.
	// It defaults to true; a resource marked as encrypted that carries no
	// fscrypt policy aborts the daemon at startup.
	NeedEncryption bool
	Binaries       []BinaryRule
}

// BinaryRule is one whitelisted binary inside a resource section.
type BinaryRule struct {
	Path string
	// Events restricts the operations this binary may perform; an empty
	// list means every event type is allowed.
	Events []ebpf.EventType
}

// Load parses the configuration file at path. Missing watch paths and
// missing binaries are skipped with a warning (matching ssh-guard's
// tolerance), and so are the directives of a skipped section — including
// directives that appear before any [watch] section. Malformed directives
// inside a valid section fail fast: a security configuration must not be
// silently misread.
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

// addResource creates a new resource for dirPath, skipping paths that do
// not exist or are not directories (matching ssh-guard's tolerance). It
// returns nil when the path is skipped.
func addResource(cfg *Config, dirPath string, lineNo int) *Resource {
	info, err := os.Stat(dirPath)
	if err != nil {
		log.Warnf("daemon config line %d: skipping missing or unreadable path: %s", lineNo, dirPath)
		return nil
	}
	if !info.IsDir() {
		log.Warnf("daemon config line %d: skipping non-directory watch path: %s", lineNo, dirPath)
		return nil
	}
	cfg.Resources = append(cfg.Resources, Resource{Path: dirPath, NeedEncryption: true})
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
	if info, statErr := os.Stat(binPath); statErr != nil {
		log.Warnf("daemon config line %d: skipping missing binary: %s", lineNo, binPath)
	} else {
		_ = info
		rule := BinaryRule{Path: binPath}
		if len(fields) > 1 {
			events, parseErr := parseEvents(strings.Split(strings.Join(fields[1:], " "), ","))
			if parseErr != nil {
				return fmt.Errorf("daemon config line %d: %w", lineNo, parseErr)
			}
			rule.Events = events
		}
		current.Binaries = append(current.Binaries, rule)
	}
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
			return nil, fmt.Errorf("unknown event type %q (valid: OPEN, READ, WRITE, DELETE, RENAME, SYMLINK, HARDLINK, MKDIR, MMAP)", token)
		}
		if !seen[et] {
			seen[et] = true
			types = append(types, et)
		}
	}
	return types, nil
}
