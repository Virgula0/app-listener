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
	if base != nil {
		if g := guardedOverlap(p, base); g != "" {
			return Target{}, fmt.Errorf("%s overlaps %s, which the daemon already guards: edit that section "+
				"instead (or record a sibling path)", p, g)
		}
	}
	return Target{Path: p, IsDir: info.IsDir()}, nil
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

// guardedOverlap is the guarded path (watch path or encryption root) of base p overlaps, or "".
func guardedOverlap(p string, base *daemonconfig.Config) string {
	for i := range base.Resources {
		r := &base.Resources[i]
		for _, g := range []string{r.Path, r.EncryptionRoot} {
			if g != "" && overlaps(p, g) {
				return g
			}
		}
	}
	return ""
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
	Sections []install.Section
	Libs     []string
	Encrypt  []string
	Binaries []PlannedBinary
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
		if !slices.ContainsFunc(rec.Targets, func(t Target) bool { return t.Path == e }) {
			return nil, fmt.Errorf("%s is not a path of this recording", e)
		}
		p.Encrypt = append(p.Encrypt, e)
	}
	for _, t := range rec.Targets {
		p.Sections = append(p.Sections, sectionFor(t.Path, chosen, slices.Contains(p.Encrypt, t.Path), sel.Masks))
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

// sectionFor whitelists, for path, the chosen binaries that touched it.
func sectionFor(path string, chosen map[string]*Binary, encrypt, masks bool) install.Section {
	s := install.Section{Path: path, Encrypt: encrypt}
	for _, line := range sortedLines(chosen) {
		b := chosen[line]
		if !slices.Contains(b.Targets, path) {
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

// Render adds the plan to base and checks the result as the daemon will: it parses, every new
// section materializes as a guarded resource with its approved lines, and no line names a multicall
// the daemon would drop. A section the parser skips (a path that vanished) would otherwise apply
// cleanly and guard nothing.
func (p *Plan) Render(base []byte) (string, error) {
	text := install.InsertSectionsBeforeLibraries(string(base), install.GenerateSections(p.Sections))
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
		if err := sectionMaterialized(cfg, &p.Sections[i]); err != nil {
			return "", err
		}
	}
	return text, nil
}

func sectionMaterialized(cfg *daemonconfig.Config, s *install.Section) error {
	for i := range cfg.Resources {
		r := &cfg.Resources[i]
		if r.Path != s.Path {
			continue
		}
		rules := slices.Concat(r.Binaries, r.PendingBinaries)
		for _, want := range s.Allow {
			if !slices.ContainsFunc(rules, func(b daemonconfig.BinaryRule) bool {
				return b.Path == want.Path || b.Link == want.Path
			}) {
				return fmt.Errorf("the daemon would not whitelist %s for %s (check its log)", want.Path, s.Path)
			}
		}
		if r.NeedEncryption != s.Encrypt {
			return fmt.Errorf("%s: need_encryption did not parse as written", s.Path)
		}
		return nil
	}
	return fmt.Errorf("the daemon would not guard %s: its section is skipped (the path vanished or is "+
		"no longer a plain directory or file)", s.Path)
}
