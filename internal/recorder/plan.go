package recorder

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"syscall"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	"github.com/Virgula0/app-listener/internal/guard"
	"github.com/Virgula0/app-listener/internal/install"
)

// LibraryBlockName names the [libraries] block recorded allow_lib entries go to.
const LibraryBlockName = "recorded access"

// MaxTargets bounds the paths of one recording.
const MaxTargets = 32

// refusedTrees are never targets, nor anything inside them: guarding them breaks the system or the
// daemon itself (its config, key and ledger live under /etc/app-listener; /sys/fs/bpf holds its
// pinned links).
var refusedTrees = []string{"/proc", "/sys", "/dev", "/run", "/boot", "/usr", "/bin", "/sbin", "/lib",
	"/lib32", "/lib64", "/libx32", "/etc/app-listener", "/var/lib/app-listener", "/var/run", "/var/lock"}

// refusedExact are never targets themselves (their subtrees may be): whole system or home roots.
var refusedExact = []string{"/", "/etc", "/home", "/var", "/opt", "/srv", "/mnt", "/media", "/tmp", "/var/tmp",
	"/root", "/snap", "/nix"}

// ValidateTargets checks the paths a recording may cover before any of them is watched: absolute,
// existing, a directory or a singly-linked regular file (never a symlink), outside system trees,
// not a user's whole home, and neither inside nor around a guarded path of base or another target.
// homes are the login users' home directories.
func ValidateTargets(paths []string, base *daemonconfig.Config, homes []string) ([]Target, error) {
	if len(paths) == 0 {
		return nil, errors.New("name at least one directory or file to record")
	}
	if len(paths) > MaxTargets {
		return nil, fmt.Errorf("at most %d paths per recording", MaxTargets)
	}
	var out []Target
	for _, raw := range paths {
		t, err := validateTarget(raw, base, homes)
		if err != nil {
			return nil, err
		}
		for _, o := range out {
			if overlaps(t.Path, o.Path) {
				return nil, fmt.Errorf("%s and %s overlap: a file follows exactly one whitelist, record sibling "+
					"paths instead", o.Path, t.Path)
			}
		}
		out = append(out, t)
	}
	return out, nil
}

func validateTarget(raw string, base *daemonconfig.Config, homes []string) (Target, error) {
	if !filepath.IsAbs(raw) || strings.ContainsRune(raw, 0) {
		return Target{}, fmt.Errorf("%q: give an absolute path", raw)
	}
	p := filepath.Clean(raw)
	if err := refusedPath(p, homes); err != nil {
		return Target{}, err
	}
	for _, suffix := range []string{fscrypt.BackupSuffix, fscrypt.DecryptSuffix} {
		if strings.HasSuffix(p, suffix) {
			return Target{}, fmt.Errorf("%s is a migration copy, not a resource", p)
		}
	}
	info, err := os.Lstat(p)
	if err != nil {
		return Target{}, fmt.Errorf("%s: %w", p, err)
	}
	if resolved, rerr := filepath.EvalSymlinks(p); rerr != nil || resolved != p {
		return Target{}, fmt.Errorf("%s is (or goes through) a symbolic link: name the real path", p)
	}
	switch {
	case info.IsDir():
	case info.Mode().IsRegular():
		if st, ok := info.Sys().(*syscall.Stat_t); ok && st.Nlink > 1 {
			return Target{}, fmt.Errorf("%s has %d hard links: an alias would extend the rule to another path", p, st.Nlink)
		}
	default:
		return Target{}, fmt.Errorf("%s is not a directory or a regular file", p)
	}
	t := Target{Path: p, IsDir: info.IsDir()}
	if base != nil {
		ex, err := existingSection(p, base)
		if err != nil {
			return Target{}, err
		}
		t.Existing = ex
	}
	return t, nil
}

func refusedPath(p string, homes []string) error {
	if slices.Contains(refusedExact, p) {
		return fmt.Errorf("%s is a system root: record a directory inside it", p)
	}
	for _, t := range refusedTrees {
		if p == t || strings.HasPrefix(p, t+"/") {
			return fmt.Errorf("%s is inside %s, which can't be guarded", p, t)
		}
		if strings.HasPrefix(t, p+"/") {
			return fmt.Errorf("%s contains %s, which can't be guarded", p, t)
		}
	}
	for _, h := range homes {
		if h != "" && filepath.Clean(h) == p {
			return fmt.Errorf("%s is a whole home directory: record the sensitive directory inside it", p)
		}
	}
	return nil
}

func overlaps(a, b string) bool {
	return a == b || strings.HasPrefix(a, b+"/") || strings.HasPrefix(b, a+"/")
}

// Existing is the daemon.conf section already guarding a target: the recording extends it.
type Existing struct {
	// Section is the [watch] header new lines go to: the target, or for a `watch:` member its
	// group's encryption root, whose whitelist every tree of the group shares (SharedWith).
	Section    string
	SharedWith []string
	Allow      []ExistingRule
	Encrypted  bool
}

// ExistingRule is one whitelist line of an Existing section; no Events = every event.
type ExistingRule struct {
	Line   string
	Events []string
}

// Has reports whether line is already whitelisted there.
func (e *Existing) Has(line string) bool {
	return e != nil && slices.ContainsFunc(e.Allow, func(r ExistingRule) bool { return r.Line == line })
}

// Warning tells the operator that recording path, which daemon.conf already guards, extends it.
func (e *Existing) Warning(path string) string {
	w := path + " is already in daemon.conf ([watch " + e.Section + "]): the recording only adds whitelist " +
		"lines to that section, every line it has stays. The guard keeps denying what it doesn't whitelist " +
		"while you record, so such a program is seen only up to its first refused operation: its events are " +
		"incomplete, prefer every event for it."
	if len(e.SharedWith) > 0 {
		w += " The whitelist is shared with " + strings.Join(e.SharedWith, ", ") + ": a line added for " + path +
			" reaches them too."
	}
	return w
}

// Lines are the section's whitelist lines as daemon.conf reads them.
func (e *Existing) Lines() []string {
	out := make([]string, len(e.Allow))
	for i, r := range e.Allow {
		out[i] = r.Line
		if len(r.Events) > 0 {
			out[i] += " " + strings.Join(r.Events, ",")
		}
	}
	return out
}

// existingSection is the section base guards exactly p with, or nil. Any other overlap with a
// guarded path is refused: a file follows exactly one whitelist, so lines recorded for p would
// neither reach a guarded tree nested in it nor the one around it.
func existingSection(p string, base *daemonconfig.Config) (*Existing, error) {
	match, err := guardedMatch(p, base)
	if err != nil || match == nil {
		return nil, err
	}
	ex := &Existing{Section: match.EncryptionRootOrPath(), Encrypted: match.NeedEncryption}
	for _, b := range slices.Concat(match.Binaries, match.PendingBinaries) {
		if b.LibBinary {
			continue
		}
		rule := ExistingRule{Line: b.Path}
		if b.Link != "" {
			rule.Line = b.Link
		}
		for _, e := range b.Events {
			rule.Events = append(rule.Events, e.String())
		}
		ex.Allow = append(ex.Allow, rule)
	}
	for i := range base.Resources {
		r := &base.Resources[i]
		if r != match && match.EncryptionRoot != "" && r.EncryptionRoot == match.EncryptionRoot {
			ex.SharedWith = append(ex.SharedWith, r.Path)
		}
	}
	return ex, nil
}

// guardedMatch is the resource of base guarding exactly p (nil when none), refusing every other
// overlap but its own group's encryption root.
func guardedMatch(p string, base *daemonconfig.Config) (*daemonconfig.Resource, error) {
	var match *daemonconfig.Resource
	for i := range base.Resources {
		if base.Resources[i].Path == p {
			match = &base.Resources[i]
			break
		}
	}
	if match != nil && match.ReadOnly {
		return nil, fmt.Errorf("%s is a lib_dir of daemon.conf, not a guarded resource: edit its section "+
			"instead", p)
	}
	for i := range base.Resources {
		r := &base.Resources[i]
		if r == match {
			continue
		}
		for _, g := range []string{r.Path, r.EncryptionRoot} {
			if g == "" || !overlaps(p, g) || (match != nil && g == match.EncryptionRoot) {
				continue
			}
			return nil, fmt.Errorf("%s overlaps %s, which the daemon already guards: record exactly a "+
				"guarded path to extend its whitelist, or a sibling path", p, g)
		}
	}
	return match, nil
}

// Selection is what the operator approved after reviewing a Recording.
type Selection struct {
	// Binaries are the approved Binary.Line values.
	Binaries []string `json:"binaries"`
	// Libs are the approved allow_lib paths (each a Lib of an approved binary).
	Libs []string `json:"libs"`
	// Encrypt are the targets to migrate to fscrypt before they are guarded.
	Encrypt []string `json:"encrypt"`
	// Masks restricts each binary to the events it was seen doing (Binary.Mask).
	Masks bool `json:"masks"`
}

// Plan is a validated Selection: what Render writes and Apply performs.
type Plan struct {
	Sections   []install.Section
	Extensions []Extension
	Libs       []string
	Encrypt    []string
	Binaries   []PlannedBinary
}

// Extension adds whitelist lines to a target's existing section.
type Extension struct {
	Path    string
	Section string
	Allow   []install.BinaryRule
}

// PlannedBinary is an approved whitelist line and the content the recording saw.
type PlannedBinary struct {
	Line         string
	SHA256       [32]byte
	SystemPlaced bool
}

// Plan validates sel against the recording: only what was recorded can be approved.
func (rec *Recording) Plan(sel *Selection) (*Plan, error) {
	p := &Plan{}
	chosen, err := rec.planBinaries(p, sel.Binaries)
	if err != nil {
		return nil, err
	}
	for _, l := range dedupe(sel.Libs) {
		if !libOfChosen(l, chosen) {
			return nil, fmt.Errorf("%s is not a library an approved binary needs", l)
		}
		p.Libs = append(p.Libs, l)
	}
	for _, e := range dedupe(sel.Encrypt) {
		i := slices.IndexFunc(rec.Targets, func(t Target) bool { return t.Path == e })
		switch {
		case i < 0:
			return nil, fmt.Errorf("%s is not a path of this recording", e)
		case rec.Targets[i].Existing != nil:
			return nil, fmt.Errorf("%s is already in daemon.conf: change its need_encryption with --edit-config", e)
		}
		p.Encrypt = append(p.Encrypt, e)
	}
	for _, t := range rec.Targets {
		s := sectionFor(&t, chosen, slices.Contains(p.Encrypt, t.Path), sel.Masks)
		switch {
		case t.Existing == nil:
			p.Sections = append(p.Sections, s)
		case len(s.Allow) > 0:
			p.Extensions = append(p.Extensions, Extension{Path: t.Path, Section: t.Existing.Section, Allow: s.Allow})
		}
	}
	if len(p.Sections) == 0 && len(p.Extensions) == 0 {
		return nil, errors.New("nothing to add: every recorded path is already in daemon.conf and no binary " +
			"it doesn't whitelist yet is selected")
	}
	return p, nil
}

func (rec *Recording) planBinaries(p *Plan, lines []string) (map[string]*Binary, error) {
	byLine := map[string]*Binary{}
	for i := range rec.Binaries {
		if b := &rec.Binaries[i]; b.Admissible() {
			byLine[b.Line] = b
		}
	}
	chosen := map[string]*Binary{}
	for _, line := range dedupe(lines) {
		b := byLine[line]
		if b == nil {
			return nil, fmt.Errorf("%s is not a whitelistable binary of this recording", line)
		}
		if b.Covered() {
			return nil, fmt.Errorf("%s is already whitelisted for every path it touched", line)
		}
		sum, err := hex.DecodeString(b.SHA256)
		if err != nil || len(sum) != 32 {
			return nil, fmt.Errorf("%s: no content hash recorded", line)
		}
		pb := PlannedBinary{Line: line, SystemPlaced: b.SystemPlaced}
		copy(pb.SHA256[:], sum)
		p.Binaries = append(p.Binaries, pb)
		chosen[line] = b
	}
	return chosen, nil
}

// sectionFor whitelists, for t, the chosen binaries that touched it and its section doesn't list.
func sectionFor(t *Target, chosen map[string]*Binary, encrypt, masks bool) install.Section {
	s := install.Section{Path: t.Path, Encrypt: encrypt}
	for _, line := range sortedLines(chosen) {
		b := chosen[line]
		if !slices.Contains(b.Targets, t.Path) || t.Existing.Has(line) {
			continue
		}
		rule := install.BinaryRule{Path: line}
		if masks {
			rule.Events = b.Mask
		}
		s.Allow = append(s.Allow, rule)
	}
	return s
}

func libOfChosen(lib string, chosen map[string]*Binary) bool {
	for _, b := range chosen {
		if slices.ContainsFunc(b.Libs, func(l Lib) bool { return l.Path == lib && l.Problem == "" }) {
			return true
		}
	}
	return false
}

func sortedLines(m map[string]*Binary) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func dedupe(in []string) []string {
	var out []string
	for _, s := range in {
		if !slices.Contains(out, s) {
			out = append(out, s)
		}
	}
	return out
}

// Render adds the plan to base and checks the result as the daemon will: it parses, every new or
// extended section materializes as a guarded resource with its approved lines, and no line names a
// multicall the daemon would drop. A section the parser skips (a path that vanished) would otherwise
// apply cleanly and guard nothing.
func (p *Plan) Render(base []byte) (string, error) {
	text := install.InsertSectionsBeforeLibraries(string(base), install.GenerateSections(p.Sections))
	for _, e := range p.Extensions {
		var err error
		if text, err = install.AppendSectionWhitelist(text, e.Section, e.Allow); err != nil {
			return "", fmt.Errorf("extending %s: %w", e.Path, err)
		}
	}
	if len(p.Libs) > 0 {
		text = install.EnsureLibraryBlock(text, &install.LibraryBlock{Name: LibraryBlockName, Libs: p.Libs})
	}
	cfg, err := daemonconfig.Parse([]byte(text))
	if err != nil {
		return "", fmt.Errorf("the new configuration does not parse: %w", err)
	}
	prev, err := daemonconfig.Parse(base)
	if err != nil {
		prev = nil
	}
	if err := guard.RefuseMulticallLines(cfg, prev); err != nil {
		return "", err
	}
	for i := range p.Sections {
		s := &p.Sections[i]
		r, err := materialized(cfg, s.Path, s.Allow)
		if err != nil {
			return "", err
		}
		if r.NeedEncryption != s.Encrypt {
			return "", fmt.Errorf("%s: need_encryption did not parse as written", s.Path)
		}
	}
	for _, e := range p.Extensions {
		if _, err := materialized(cfg, e.Path, e.Allow); err != nil {
			return "", err
		}
	}
	return text, nil
}

// materialized is cfg's resource guarding path, once it whitelists every line of allow.
func materialized(cfg *daemonconfig.Config, path string, allow []install.BinaryRule) (*daemonconfig.Resource, error) {
	for i := range cfg.Resources {
		r := &cfg.Resources[i]
		if r.Path != path {
			continue
		}
		rules := slices.Concat(r.Binaries, r.PendingBinaries)
		for _, want := range allow {
			if !slices.ContainsFunc(rules, func(b daemonconfig.BinaryRule) bool {
				return !b.LibBinary && (b.Path == want.Path || b.Link == want.Path)
			}) {
				return nil, fmt.Errorf("the daemon would not whitelist %s for %s (check its log)", want.Path, path)
			}
		}
		return r, nil
	}
	return nil, fmt.Errorf("the daemon would not guard %s: its section is skipped (the path vanished or is "+
		"no longer a plain directory or file)", path)
}
