package daemonconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// ParseNarrowSection recognizes `[narrow <path>]`: a guarded subtree (directory or file) inside a
// tree an earlier section guards, judged by its own whitelist, a subset of that tree's. guard_inodes
// keeps one owner per inode, the innermost resource's, so the subtree's files follow only its list.
func ParseNarrowSection(line string) (path string, ok bool, err error) {
	if !strings.HasPrefix(line, "[narrow ") || !strings.HasSuffix(line, "]") {
		return "", false, nil
	}
	path = UnquotePath(strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "[narrow "), "]")))
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return "", true, fmt.Errorf("malformed section header %q: a [narrow] path is one clean absolute path", line)
	}
	return path, true, nil
}

// applyNarrowLine records one binary line of a [narrow] section; nothing else belongs there (its
// encryption and libraries are its parent's).
func applyNarrowLine(g *watchGroup, line string, lineNo int) error {
	_, isEnc := parseNeedEncryption(line)
	_, isLib := parseAllowLib(line)
	_, isLibDir := parseLibDir(line)
	_, isLibBin := parseLibBinary(line)
	_, isRuntime := CutDirective(line, "runtime")
	if isEnc || isLib || isLibDir || isLibBin || isRuntime || isWatchDirective(line) {
		return fmt.Errorf("daemon config line %d: a [narrow] section takes only binary lines (its "+
			"encryption and libraries are its parent tree's)", lineNo)
	}
	return applyBinaryLine(g, line, lineNo)
}

// materializeNarrowGroup adds a [narrow] section's resource below its parent, the innermost
// whitelist-mode tree declared before it. A target missing (or invisible in a locked vault the
// parent encrypts: deferred as PathPending) is an error, never a silent widening to the parent's list.
func materializeNarrowGroup(cfg *Config, g *watchGroup) error {
	parent, err := narrowParent(cfg, g.root)
	if err != nil {
		return fmt.Errorf("daemon config line %d: %w", g.lineNo, err)
	}
	if err := narrowSubsetErr(g, parent); err != nil {
		return fmt.Errorf("daemon config line %d: [narrow %s]: %w", g.lineNo, g.root, err)
	}
	// One section group with its parent (EncryptionGroups keeps the parent): never its own vault.
	res := Resource{Path: g.root, NeedEncryption: parent.NeedEncryption,
		EncryptionRoot: parent.EncryptionRootOrPath(), NarrowOf: parent.Path}
	switch vErr := validateWatchTarget(g.root, parent.Path); {
	case vErr == nil:
	case parent.PathPending || narrowInLockedVault(g.root, parent):
		res.PathPending = true
	default:
		return fmt.Errorf("daemon config line %d: [narrow %s]: %w (create it, or remove the section: its "+
			"files would follow %s's wider whitelist)", g.lineNo, g.root, vErr, parent.Path)
	}
	res.Binaries = append(res.Binaries, g.binaries...)
	res.PendingBinaries = append(res.PendingBinaries, g.pending...)
	cfg.Resources = append(cfg.Resources, res)
	return nil
}

// narrowInLockedVault: target doesn't stat, and the parent's vault root is an encrypted directory
// whose names are unreadable while locked. ResolvePendingPaths re-validates it after unlock.
func narrowInLockedVault(target string, parent *Resource) bool {
	if !parent.NeedEncryption {
		return false
	}
	if _, err := os.Lstat(target); err == nil {
		return false
	}
	fi, err := os.Lstat(parent.EncryptionRootOrPath())
	return err == nil && fi.IsDir()
}

// narrowParent finds the resource path lies strictly inside: a [watch] tree (or watch: path) in
// whitelist mode, never a lib_dir, runtime, port, secret or another [narrow].
func narrowParent(cfg *Config, path string) (*Resource, error) {
	var parent *Resource
	for i := range cfg.Resources {
		r := &cfg.Resources[i]
		if isInsidePath(path, r.Path) && (parent == nil || isInsidePath(r.Path, parent.Path)) {
			parent = r
		}
	}
	switch {
	case parent == nil:
		return nil, fmt.Errorf("[narrow %s] is not inside a tree a [watch] section declared before it", path)
	case parent.NarrowOf != "":
		return nil, fmt.Errorf("[narrow %s] is inside [narrow %s]: only one level", path, parent.Path)
	case parent.ReadOnly || parent.Runtime || parent.Port != 0 || parent.Secret != nil || parent.SecretRelay:
		return nil, fmt.Errorf("[narrow %s]: its parent %s is not a whitelist tree", path, parent.Path)
	}
	return parent, nil
}

// narrowSubsetErr refuses a line the parent doesn't whitelist, or one granting an event the
// parent's line doesn't: a [narrow] only takes access away.
func narrowSubsetErr(g *watchGroup, parent *Resource) error {
	lines := slices.Concat(parent.Binaries, parent.PendingBinaries)
	for _, c := range slices.Concat(g.binaries, g.pending) {
		i := slices.IndexFunc(lines, func(p BinaryRule) bool { return p.configured() == c.configured() })
		if i < 0 {
			return fmt.Errorf("%s is not whitelisted in %s", c.configured(), parent.Path)
		}
		p := lines[i]
		if len(p.Events) == 0 {
			continue
		}
		if len(c.Events) == 0 {
			return fmt.Errorf("%s may do every event here but not in %s: list its events", c.configured(),
				parent.Path)
		}
		for _, e := range c.Events {
			if !slices.Contains(p.Events, e) {
				return fmt.Errorf("%s may not %s in %s", c.configured(), e, parent.Path)
			}
		}
	}
	return nil
}

// configured is the line as written: the symlink when Path is its resolved target.
func (b BinaryRule) configured() string {
	if b.Link != "" {
		return b.Link
	}
	return b.Path
}
