// Package daemonconfig parses the daemon config (ssh-guard.conf grammar): each [watch <dir>]
// section protects one resource and lists allowed binaries with optional event types.
// chattr/exclude_chattr are deliberately unsupported (the LSM guard replaces them).
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
	// SharedAllowLibs are the allow_lib paths from [libraries] blocks. Library trust is
	// daemon-wide, so they belong to no resource; all go to the trust guard, which re-stats each
	// after the vaults unlock.
	SharedAllowLibs []string
}

// Resource is one guarded tree (directory or vault root) with its own guard, whitelist and, via
// EncryptionRoot, a shared fscrypt lifecycle. Symlinks, hard-linked files and special files are
// refused at parse time.
type Resource struct {
	Path string
	// NeedEncryption selects the fscrypt lifecycle (default true). An encrypted-marked resource
	// without an fscrypt policy aborts startup.
	NeedEncryption bool
	// EncryptionRoot is the vault root governing this resource, set by `watch:` directives in a
	// [watch <root>] group (section path = encryption root, watch paths = guarded trees). Empty =
	// the resource path is its own root.
	EncryptionRoot string
	Binaries       []BinaryRule
	// PendingBinaries parks whitelisted binaries unreadable at parse time (typically an
	// fscrypt-locked dir). Until the post-unlock pass moves them to Binaries they stay out of the
	// BPF whitelist: denied (fail-closed).
	PendingBinaries []BinaryRule
	// AllowLibs (`allow_lib <path>`) are extra libraries this section's whitelisted binaries may
	// load, beyond the static ELF closure (DT_NEEDED + interpreter): dlopen extras static analysis
	// can't see (NSS, gconv, GL drivers, plugins). Their inodes join the global trusted-library set
	// (internal/guard trust object): a whitelisted binary may exec-map only trusted libraries, and
	// no non-whitelisted process may overwrite one.
	AllowLibs []string
	// PendingLibs parks allow_lib paths unreadable at parse time (same post-unlock pass as
	// PendingBinaries).
	PendingLibs []string
	// ReadOnly guards the tree in guard.ModeReadOnly: every process may READ, modifications are
	// whitelist-gated. Backs `lib_dir`: a library dir holds no secrets (denying reads would break
	// unrelated processes), and the write monopoly is what makes its libraries trustworthy to load.
	// Never encrypted.
	ReadOnly bool
	// PathPending marks a grouped watch path unvalidated at parse time because its encryption root
	// is a locked vault (the sub-path doesn't resolve until unlock). After unlocking under an
	// ephemeral guard the daemon calls ResolvePendingPaths (same symlink/hard-link/type refusals as
	// addResource), which clears this flag. Still missing afterwards = hard error, never silently
	// dropped.
	PathPending bool
}

// EncryptionRootOrPath returns the vault root governing this resource: the `watch:` group's section
// path if set, else the resource path. Installer flows patching config text or driving fscrypt must
// address this path; grouped resources share ONE section ([watch <encryption root>]) and must be
// deduplicated by it.
func (r *Resource) EncryptionRootOrPath() string {
	if r.EncryptionRoot != "" {
		return r.EncryptionRoot
	}
	return r.Path
}

// EncryptionGroups returns resources deduplicated by encryption root, in config order, keeping the
// FIRST of each group (its EncryptionRootOrPath is the shared section header text-patching helpers
// address). Returned by pointer: patch grouped resources via their SECTION path, never a watch
// sub-path. Read-only (`lib_dir`) resources are excluded: no fscrypt question, migration or config
// patch may address them.
func (c *Config) EncryptionGroups() []*Resource {
	seen := make(map[string]bool)
	var out []*Resource
	for i := range c.Resources {
		if c.Resources[i].ReadOnly {
			continue
		}
		root := c.Resources[i].EncryptionRootOrPath()
		if !seen[root] {
			seen[root] = true
			out = append(out, &c.Resources[i])
		}
	}
	return out
}

// BinaryRule is one whitelisted binary inside a resource section.
type BinaryRule struct {
	Path string
	// Events restricts the operations this binary may perform; empty means every type.
	Events []ebpf.EventType
}

// watchGroup is the parse state of one [watch <root>] section: the section path is the ENCRYPTION
// ROOT, `watch: <path>` directives are the GUARDED trees (each a Resource sharing the group's
// whitelist and lifecycle), binary directives are the shared whitelist. Without `watch:` the
// section path itself is the guarded tree.
type watchGroup struct {
	root           string
	needEncryption bool
	binaries       []BinaryRule
	pending        []BinaryRule
	libs           []string
	pendingLibs    []string
	watchPaths     []string
	libDirs        []string
	// libBinaries/libPending are writers of this group's lib_dirs ONLY, never of its own (secret,
	// whitelist-mode) resources: a runtime tree is maintained by vendor tools (Steam's
	// pressure-vessel rebuilds /usr each launch) that have no business reading the credential
	// directory in the same section.
	libBinaries []BinaryRule
	libPending  []BinaryRule
	// libraryBlock marks a [libraries <name>] block: no watch root, only library directives; its
	// lib_dirs are writable by its own lib_binary rules alone.
	libraryBlock bool
	blockName    string
	// skipped marks a section whose root is missing: its directives are warned and ignored, never
	// fatal (the group is dropped at finalize).
	skipped bool
	lineNo  int
}

// Load parses the config at path. Missing watch paths and directives outside any [watch] section
// are skipped with a warning; unreadable binaries go to PendingBinaries for post-unlock resolution.
// Malformed directives in valid sections fail fast.
func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &Config{}
	if err := parseConfig(cfg, file); err != nil {
		return nil, err
	}
	if err := validateResources(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// parseConfig reads the file into cfg, one Resource per guarded tree (a watch group emits one per
// `watch:` path).
func parseConfig(cfg *Config, file *os.File) error {
	var group *watchGroup

	scanner := bufio.NewScanner(file)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if err := applyConfigLine(cfg, &group, line, lineNo); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	finalizeGroup(cfg, group)
	return nil
}

// applyConfigLine dispatches one line: section headers close the previous group and open a new one;
// directives mutate the open group.
func applyConfigLine(cfg *Config, group **watchGroup, line string, lineNo int) error {
	if name, ok := parseLibrariesSection(line); ok {
		finalizeGroup(cfg, *group)
		*group = &watchGroup{libraryBlock: true, blockName: name, lineNo: lineNo}
		return nil
	}
	dirPath, isSection, headerErr := parseWatchSection(line)
	if headerErr != nil {
		return fmt.Errorf("daemon config line %d: %w", lineNo, headerErr)
	}
	if isSection {
		finalizeGroup(cfg, *group)
		*group = newWatchGroup(dirPath, lineNo)
		return nil
	}
	if *group == nil {
		// Tolerated like ssh-guard, but warned: a silently dropped security directive must be
		// spottable.
		log.Warnf("daemon config line %d: ignoring directive outside any [watch] section: %q", lineNo, line)
		return nil
	}
	g := *group
	if g.libraryBlock {
		return applyLibraryBlockDirective(g, line, lineNo)
	}
	if g.skipped {
		log.Warnf("daemon config line %d: skipped section: ignoring directive %q", lineNo, line)
		return nil
	}
	if isWatchDirective(line) {
		return g.addWatchPath(line, lineNo)
	}
	// TUI placeholder form: a bare [watch] header, then `path = <dir>` sets the section root.
	if g.root == "" {
		if p, ok := parsePathDirective(line); ok {
			g.root = p
			if _, statErr := os.Lstat(p); statErr != nil {
				g.skipped = true
			}
			return nil
		}
		log.Warnf("daemon config line %d: ignoring directive of placeholder section: %q", lineNo, line)
		return nil
	}
	return applyDirective(g, line, lineNo)
}

// newWatchGroup probes the section root: a missing root skips the whole group (directives warned,
// never fatal).
func newWatchGroup(dirPath string, lineNo int) *watchGroup {
	g := &watchGroup{root: dirPath, needEncryption: true, lineNo: lineNo}
	if dirPath != "" {
		if _, statErr := os.Lstat(dirPath); statErr != nil {
			g.skipped = true
		}
	}
	return g
}

// finalizeGroup materializes a closed group unless it was skipped.
func finalizeGroup(cfg *Config, group *watchGroup) {
	if group == nil || group.skipped {
		return
	}
	if group.libraryBlock {
		cfg.SharedAllowLibs = append(cfg.SharedAllowLibs, group.libs...)
		cfg.SharedAllowLibs = append(cfg.SharedAllowLibs, group.pendingLibs...)
		materializeLibDirs(cfg, group) // writers: the block's lib_binary rules only
		return
	}
	materializeWatchGroup(cfg, group)
}

// isWatchDirective recognizes `watch: <path>` (optional space after the colon), adding a guarded
// tree to the section's encryption group.
func isWatchDirective(line string) bool {
	return strings.HasPrefix(line, "watch:") || strings.HasPrefix(line, "watch ")
}

// addWatchPath validates and records one extra guarded tree: it must live INSIDE the section path
// (the encryption root; fscrypt is vault-wide) and not duplicate the section path or another watch
// path.
func (g *watchGroup) addWatchPath(line string, lineNo int) error {
	path := unquotePath(strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, "watch:"), "watch")))
	if path == "" {
		return fmt.Errorf("daemon config line %d: empty watch path", lineNo)
	}
	// Resolve lexical traversal before confinement: dir/../escape must not pass the inside-root
	// check.
	path = filepath.Clean(path)
	if !isInsidePath(path, g.root) {
		return fmt.Errorf("daemon config line %d: watch path %q must be inside the section directory %q", lineNo, path, g.root)
	}
	for _, existing := range g.watchPaths {
		if existing == path {
			return fmt.Errorf("daemon config line %d: duplicate watch path: %s", lineNo, path)
		}
	}
	g.watchPaths = append(g.watchPaths, path)
	return nil
}

// isInsidePath reports whether path is a strict sub-path of dir.
func isInsidePath(path, dir string) bool {
	return path != dir && strings.HasPrefix(path+"/", dir+"/")
}

// unquotePath strips a surrounding pair of double quotes. Quoting is optional for paths without
// spaces/special characters (legacy form); otherwise required, since event types following a path
// on the same line would be ambiguous. The installer always quotes.
func unquotePath(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

// splitPathAndRest extracts a leading path: a double-quoted path (may contain whitespace) or, for
// hand-written legacy configs, the first whitespace-delimited token. Returns the unquoted path and
// the trimmed rest of the line.
func splitPathAndRest(line string) (path, rest string, err error) {
	if strings.HasPrefix(line, `"`) {
		closeIdx := strings.IndexByte(line[1:], '"')
		if closeIdx == -1 {
			return "", "", fmt.Errorf("unterminated quoted path: %s", line)
		}
		return line[1 : 1+closeIdx], strings.TrimSpace(line[1+closeIdx+1:]), nil
	}
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return "", "", nil
	}
	path = fields[0]
	return path, strings.TrimSpace(line[len(path):]), nil
}

// parsePathDirective recognizes the TUI placeholder `path = <dir>` (sets the root of a bare [watch]
// section).
func parsePathDirective(line string) (string, bool) {
	for _, prefix := range []string{"path = ", "path="} {
		if strings.HasPrefix(line, prefix) {
			if value := strings.TrimSpace(strings.TrimPrefix(line, prefix)); value != "" {
				return unquotePath(value), true
			}
		}
	}
	return "", false
}

// materializeWatchGroup appends one Resource per guarded tree, sharing the group's whitelist,
// need_encryption and encryption root. Unguardable paths (missing, symlinks, ...) are skipped with
// a warning by addResource, except a grouped watch path merely invisible because its encryption
// root is a locked vault: kept as PathPending for the daemon to re-validate after unlock
// (deferPendingWatchPath).
func materializeWatchGroup(cfg *Config, g *watchGroup) {
	paths := g.watchPaths
	if len(paths) == 0 {
		paths = []string{g.root}
	}
	grouped := len(g.watchPaths) > 0
	for _, watchPath := range paths {
		encRoot := ""
		if grouped {
			encRoot = g.root
		}
		res := addResource(cfg, watchPath, encRoot, g.lineNo)
		if res == nil {
			if grouped {
				res = deferPendingWatchPath(cfg, g, watchPath)
			}
			if res == nil {
				continue
			}
		}
		res.NeedEncryption = g.needEncryption
		if grouped {
			res.EncryptionRoot = g.root
		}
		res.Binaries = append(res.Binaries, g.binaries...)
		res.PendingBinaries = append(res.PendingBinaries, g.pending...)
		res.AllowLibs = append(res.AllowLibs, g.libs...)
		res.PendingLibs = append(res.PendingLibs, g.pendingLibs...)
	}
	materializeLibDirs(cfg, g)
}

// materializeLibDirs appends one read-only, unencrypted Resource per `lib_dir`, sharing the group's
// whitelist. The tree stays world-readable while only the section's binaries may
// create/replace/alter files; that write monopoly lets the trust object treat every .so under it as
// loadable without enumerating it, the only workable rule for per-launch runtime dirs (Steam's
// pressure-vessel `var/tmp-XXXXXX`).
func materializeLibDirs(cfg *Config, g *watchGroup) {
	// A lib_binary without any lib_dir grants nothing: say so rather than let the author believe it
	// was whitelisted.
	if len(g.libDirs) == 0 && len(g.libBinaries)+len(g.libPending) > 0 {
		log.Warnf("daemon config line %d: the section declares lib_binary but no lib_dir — "+
			"those binaries grant nothing (a lib_binary is a writer of this section's library directories only)", g.lineNo)
	}
	for _, dir := range g.libDirs {
		if _, statErr := os.Lstat(dir); statErr != nil {
			log.Warnf("daemon config line %d: lib_dir not present, ignoring: %s", g.lineNo, dir)
			continue
		}
		// A shared runtime tree may be named by several sections (one app, several config
		// locations; `install --diff-catalog` appending a second section). Guarding it once is
		// enough, so a repeat is a no-op, not a duplicate-watch-path error that would take the
		// daemon down. Only merge into another lib_dir: a lib_dir naming an already-guarded real
		// (whitelist-mode, possibly encrypted) resource must NOT fold into it (it would widen a
		// secret directory's whitelist); refuse.
		if existing := findResource(cfg, dir); existing != nil {
			if !existing.ReadOnly {
				log.Errorf("daemon config line %d: lib_dir %s is already a guarded resource — ignoring the directive "+
					"(a library directory must not share a path with a protected tree)", g.lineNo, dir)
				continue
			}
			existing.Binaries = append(existing.Binaries, libDirWriters(g.binaries)...)
			existing.PendingBinaries = append(existing.PendingBinaries, libDirWriters(g.pending)...)
			existing.Binaries = append(existing.Binaries, g.libBinaries...)
			existing.PendingBinaries = append(existing.PendingBinaries, g.libPending...)
			continue
		}
		res := addResource(cfg, dir, "", g.lineNo)
		if res == nil {
			continue
		}
		res.NeedEncryption = false
		res.ReadOnly = true
		res.Binaries = append(res.Binaries, libDirWriters(g.binaries)...)
		res.PendingBinaries = append(res.PendingBinaries, libDirWriters(g.pending)...)
		res.Binaries = append(res.Binaries, g.libBinaries...)
		res.PendingBinaries = append(res.PendingBinaries, g.libPending...)
	}
}

// libDirWriters returns the section binaries allowed to WRITE a lib_dir: only unrestricted ones. An
// event list scopes a binary to the section's protected tree and is meaningless in read-only mode
// (no per-binary masks); promoting such a binary to an unrestricted writer would widen it, so it's
// left out (it can still read).
func libDirWriters(rules []BinaryRule) []BinaryRule {
	out := make([]BinaryRule, 0, len(rules))
	for _, r := range rules {
		if len(r.Events) == 0 {
			out = append(out, r)
		}
	}
	return out
}

// deferPendingWatchPath keeps a grouped watch path addResource rejected ONLY when it's "invisible
// while the vault is locked": the group declares need_encryption, the path fails to stat, and the
// encryption root is a real directory. Any other rejection (symlink, hard-linked file, missing path
// under an unencrypted group) stays dropped. Returns nil when the path must not be resurrected.
func deferPendingWatchPath(cfg *Config, g *watchGroup, watchPath string) *Resource {
	if !g.needEncryption {
		return nil
	}
	if _, statErr := os.Lstat(watchPath); statErr == nil {
		// The path resolves: addResource rejected it on merit (symlink/hard link/special file), not
		// visibility.
		return nil
	}
	rootInfo, rootErr := os.Lstat(g.root)
	if rootErr != nil || rootInfo.Mode()&os.ModeSymlink != 0 || !rootInfo.IsDir() {
		return nil
	}
	log.Warnf("daemon config line %d: grouped watch path not resolvable yet (encryption root %s locked?), deferring: %s",
		g.lineNo, g.root, watchPath)
	cfg.Resources = append(cfg.Resources, Resource{Path: watchPath, NeedEncryption: true, PathPending: true})
	return &cfg.Resources[len(cfg.Resources)-1]
}

// ResolvePendingPaths re-validates every PathPending resource after its encryption root is
// unlocked: the sub-path must now resolve to a directory or unique regular file (symlinks,
// hard-linked and special files refused as in addResource). A still-missing path is a hard error:
// dropping it would leave a declared-protected directory unguarded.
func ResolvePendingPaths(cfg *Config) error {
	for i := range cfg.Resources {
		r := &cfg.Resources[i]
		if !r.PathPending {
			continue
		}
		if err := validateWatchTarget(r.Path, r.EncryptionRoot); err != nil {
			return fmt.Errorf("grouped watch path %s (encryption root %s): %w", r.Path, r.EncryptionRoot, err)
		}
		r.PathPending = false
	}
	return nil
}

// findResource returns the already-materialized resource for path, or nil.
func findResource(cfg *Config, path string) *Resource {
	for i := range cfg.Resources {
		if cfg.Resources[i].Path == path {
			return &cfg.Resources[i]
		}
	}
	return nil
}

// validateResources rejects duplicate watch paths across the whole config.
func validateResources(cfg *Config) error {
	seen := make(map[string]bool, len(cfg.Resources))
	for i := range cfg.Resources {
		r := &cfg.Resources[i]
		if seen[r.Path] {
			return fmt.Errorf("daemon config: duplicate watch path: %s", r.Path)
		}
		seen[r.Path] = true
	}
	return nil
}

func parseWatchSection(line string) (path string, isSection bool, err error) {
	if !strings.HasPrefix(line, "[") {
		return "", false, nil
	}
	// Bare "[watch]" is the installer's placeholder header (path filled in later); accept it with
	// an empty path (skipped as unguardable by addResource).
	if line == "[watch]" {
		return "", true, nil
	}
	// Any other bracketed line is a section header: treating malformed ones as directives silently
	// merged following binaries into the PREVIOUS section (cross-resource whitelist contamination
	// on manual edits).
	if !strings.HasPrefix(line, "[watch ") || !strings.HasSuffix(line, "]") {
		return "", false, fmt.Errorf("malformed section header %q: expected \"[watch <path>]\"", line)
	}
	return unquotePath(strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "[watch "), "]"))), true, nil
}

// addResource creates a resource for watchPath, or nil for unguardable targets: missing paths,
// symlinks, hard-linked regular files (the inode-based guard would implicitly cover another path),
// anything not a directory or unique regular file. encRoot (grouped sub-path, "" otherwise) is used
// to reject a symlinked intermediate component.
func addResource(cfg *Config, watchPath, encRoot string, lineNo int) *Resource {
	if err := validateWatchTarget(watchPath, encRoot); err != nil {
		log.Warnf("daemon config line %d: skipping %s: %v", lineNo, watchPath, err)
		return nil
	}
	cfg.Resources = append(cfg.Resources, Resource{Path: watchPath, NeedEncryption: true})
	return &cfg.Resources[len(cfg.Resources)-1]
}

// validateWatchTarget refuses a target that would make the inode-based guard ambiguous or
// meaningless: unreadable/missing, symlink, hard-linked regular file, or not a directory/unique
// regular file. With encRoot set, an existing symlinked component between encRoot and watchPath is
// refused too (it could redirect the guard outside the vault while sharing the group's whitelist
// and lifecycle).
func validateWatchTarget(watchPath, encRoot string) error {
	if encRoot != "" {
		if err := rejectSymlinkedComponents(encRoot, watchPath); err != nil {
			return err
		}
	}
	info, err := os.Lstat(watchPath)
	if err != nil {
		return fmt.Errorf("missing or unreadable path")
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("symbolic links are refused as watch targets (guard identity is inode based)")
	}
	if !info.IsDir() {
		if !info.Mode().IsRegular() {
			return fmt.Errorf("only directories and regular files can be watched")
		}
		if stat, ok := info.Sys().(*syscall.Stat_t); ok && stat.Nlink > 1 {
			return fmt.Errorf("hard-linked files (%d links) are refused as watch targets", stat.Nlink)
		}
	}
	return nil
}

// rejectSymlinkedComponents refuses any existing component strictly between encRoot and target that
// is a symlink (validateWatchTarget Lstats only the final one, so `watch: <root>/link/sub` with
// link pointing outside the vault would redirect the guard). Components that don't exist yet are
// left to leaf validation and the deferred pass (a locked vault hides sub-path names).
func rejectSymlinkedComponents(encRoot, target string) error {
	rel, relErr := filepath.Rel(encRoot, target)
	if relErr != nil || rel == "." || rel == "" || strings.HasPrefix(rel, "..") {
		return nil //nolint:nilerr // a non-sub-path (relErr, or ".."/".") is not this check's concern — isInsidePath already rejected it
	}
	segs := strings.Split(rel, string(filepath.Separator))
	cur := encRoot
	for _, seg := range segs[:len(segs)-1] { // intermediates only; the leaf is validateWatchTarget's job
		if seg == "" || seg == "." {
			continue
		}
		cur = filepath.Join(cur, seg)
		info, statErr := os.Lstat(cur)
		if statErr != nil {
			return nil //nolint:nilerr // a component not present yet is deliberately deferred to the leaf validation / ResolvePendingPaths pass, not an error here
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("intermediate path component %q is a symbolic link (guard identity is inode based)", cur)
		}
	}
	return nil
}

// applyDirective handles one [watch]-group directive: need_encryption, binary rules and
// (unsupported) chattr lines. Directives are group-scoped: they apply to every tree materialized
// from the section.
func applyDirective(group *watchGroup, line string, lineNo int) error {
	if value, ok := parseNeedEncryption(line); ok {
		switch value {
		case "true":
			group.needEncryption = true
		case "false":
			group.needEncryption = false
		default:
			return fmt.Errorf("daemon config line %d: invalid need_encryption value %q (expected true or false)", lineNo, value)
		}
		return nil
	}

	if libPath, ok := parseAllowLib(line); ok {
		return applyAllowLib(group, libPath, lineNo)
	}

	if dirPath, ok := parseLibDir(line); ok {
		return applyLibDir(group, dirPath, lineNo)
	}

	if value, ok := parseLibBinary(line); ok {
		return applyLibBinary(group, value, lineNo)
	}

	binPath, rest, splitErr := splitPathAndRest(line)
	if splitErr != nil {
		return fmt.Errorf("daemon config line %d: %w", lineNo, splitErr)
	}
	if binPath == "" {
		return nil
	}
	rule := BinaryRule{Path: binPath}
	if rest != "" {
		events, parseErr := parseEvents(strings.Split(rest, ","))
		if parseErr != nil {
			return fmt.Errorf("daemon config line %d: %w", lineNo, parseErr)
		}
		rule.Events = events
	}
	recordBinaryRule(rule, &group.binaries, &group.pending, lineNo)
	return nil
}

// recordBinaryRule stats and symlink-resolves one rule into resolved or, if unreadable (locked tree
// or gone), deferred. Deferring rather than dropping: a dropped rule silently disables the entry,
// while a deferred one stays out of the BPF whitelist (denied) until it resolves (fail-closed).
func recordBinaryRule(rule BinaryRule, resolved, deferred *[]BinaryRule, lineNo int) {
	if _, statErr := os.Stat(rule.Path); statErr != nil {
		log.Warnf("daemon config line %d: binary not readable yet, deferring: %s", lineNo, rule.Path)
		*deferred = append(*deferred, rule)
		return
	}
	if target, resolveErr := filepath.EvalSymlinks(rule.Path); resolveErr == nil && target != rule.Path {
		log.Infof("daemon config line %d: binary symlink resolved: %s -> %s", lineNo, rule.Path, target)
		rule.Path = target
	}
	*resolved = append(*resolved, rule)
}

// parseAllowLib recognizes `allow_lib <path>` / `allow_lib: <path>` and returns the unquoted path.
func parseAllowLib(line string) (string, bool) {
	rest, ok := strings.CutPrefix(line, "allow_lib")
	if !ok {
		return "", false
	}
	// Require a separator so "/usr/bin/allow_libfoo" isn't mistaken for the directive.
	if rest == "" || (rest[0] != ' ' && rest[0] != '\t' && rest[0] != ':') {
		return "", false
	}
	rest = strings.TrimPrefix(rest, ":")
	return unquotePath(strings.TrimSpace(rest)), true
}

// applyAllowLib records one allow_lib path on the open group, deferring it (like a binary) if not
// yet readable.
func applyAllowLib(group *watchGroup, libPath string, lineNo int) error {
	if libPath == "" {
		return fmt.Errorf("daemon config line %d: empty allow_lib path", lineNo)
	}
	if _, statErr := os.Stat(libPath); statErr != nil {
		log.Warnf("daemon config line %d: allow_lib not readable yet, deferring: %s", lineNo, libPath)
		group.pendingLibs = append(group.pendingLibs, libPath)
		return nil //nolint:nilerr // deferred, not dropped — same fail-closed handling as a binary
	}
	if resolved, resolveErr := filepath.EvalSymlinks(libPath); resolveErr == nil && resolved != libPath {
		libPath = resolved
	}
	group.libs = append(group.libs, libPath)
	return nil
}

// parseLibDir recognizes `lib_dir <path>` / `lib_dir: <path>` and returns the unquoted path.
func parseLibDir(line string) (string, bool) {
	rest, ok := strings.CutPrefix(line, "lib_dir")
	if !ok {
		return "", false
	}
	// Require a separator so "/usr/bin/lib_dirfoo" isn't mistaken for the directive.
	if rest == "" || (rest[0] != ' ' && rest[0] != '\t' && rest[0] != ':') {
		return "", false
	}
	rest = strings.TrimPrefix(rest, ":")
	return unquotePath(strings.TrimSpace(rest)), true
}

// applyLibDir records one library directory on the open group. Unlike `watch:` it is NOT confined
// to the section root and never joins its encryption group: it becomes its own read-only resource
// (materializeLibDirs). A missing directory is warned and dropped, not deferred: it holds no
// secret, and failing the config over an optional runtime tree (a removed Proton version) would
// take the daemon down.
func applyLibDir(group *watchGroup, dirPath string, lineNo int) error {
	if dirPath == "" {
		return fmt.Errorf("daemon config line %d: empty lib_dir path", lineNo)
	}
	dirPath = filepath.Clean(dirPath)
	for _, existing := range group.libDirs {
		if existing == dirPath {
			return fmt.Errorf("daemon config line %d: duplicate lib_dir path: %s", lineNo, dirPath)
		}
	}
	group.libDirs = append(group.libDirs, dirPath)
	return nil
}

// parseLibBinary recognizes `lib_binary <path>` / `lib_binary: <path>` and returns the raw value
// (path plus rest).
func parseLibBinary(line string) (string, bool) {
	rest, ok := strings.CutPrefix(line, "lib_binary")
	if !ok {
		return "", false
	}
	// Require a separator so "/usr/bin/lib_binaryfoo" isn't mistaken for the directive.
	if rest == "" || (rest[0] != ' ' && rest[0] != '\t' && rest[0] != ':') {
		return "", false
	}
	rest = strings.TrimPrefix(rest, ":")
	return strings.TrimSpace(rest), true
}

// parseLibrariesSection recognizes a `[libraries]`, `[libraries <name>]` or `[libraries "<name>"]`
// header. The name is only a label for the application; it scopes nothing but the block itself (a
// lib_binary writes only its own block's lib_dirs).
func parseLibrariesSection(line string) (string, bool) {
	if !strings.HasPrefix(line, "[libraries") || !strings.HasSuffix(line, "]") {
		return "", false
	}
	rest := strings.TrimSuffix(strings.TrimPrefix(line, "[libraries"), "]")
	if rest != "" && rest[0] != ' ' && rest[0] != '\t' {
		return "", false // "[librariesfoo]": let parseWatchSection reject it
	}
	return unquotePath(strings.TrimSpace(rest)), true
}

// applyLibraryBlockDirective handles one line of a [libraries] block. Only the three library
// directives are meaningful; anything else, a binary path above all, is refused: in a watch section
// a bare path is a whitelist entry, and treating it as one here would look like a grant while
// granting nothing.
func applyLibraryBlockDirective(g *watchGroup, line string, lineNo int) error {
	if libPath, ok := parseAllowLib(line); ok {
		return applyAllowLib(g, libPath, lineNo)
	}
	if dirPath, ok := parseLibDir(line); ok {
		return applyLibDir(g, dirPath, lineNo)
	}
	if value, ok := parseLibBinary(line); ok {
		return applyLibBinary(g, value, lineNo)
	}
	return fmt.Errorf("daemon config line %d: %q is not allowed in a [libraries] block "+
		"(only allow_lib, lib_dir and lib_binary are)", lineNo, line)
}

// applyLibBinary records one writer of the group's `lib_dir` trees. Deliberately NOT a whitelist
// entry: it reaches the read-only library resources only, never the section's own protected tree
// (vendor helpers such as Steam's pressure-vessel that own a runtime tree must not also get the
// credential directory listed in the same section).
//
// Event restrictions are rejected, not ignored: a lib_dir is read-only guarded, which has no
// per-binary masks, so accepting a list would promise a narrowing that can't be enforced.
func applyLibBinary(group *watchGroup, value string, lineNo int) error {
	binPath, rest, splitErr := splitPathAndRest(value)
	if splitErr != nil {
		return fmt.Errorf("daemon config line %d: %w", lineNo, splitErr)
	}
	if binPath == "" {
		return fmt.Errorf("daemon config line %d: empty lib_binary path", lineNo)
	}
	if rest != "" {
		return fmt.Errorf("daemon config line %d: lib_binary does not accept event restrictions "+
			"(%q): a library directory is guarded read-only, which carries no per-binary event mask", lineNo, rest)
	}
	recordBinaryRule(BinaryRule{Path: binPath}, &group.libBinaries, &group.libPending, lineNo)
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
