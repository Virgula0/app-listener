package daemon

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
	"github.com/Virgula0/app-listener/internal/install"
)

// buildGlobReservations derives trust protection #3 from the built-in catalog, the source of the
// globs the catalog refresh expands. An entry counts for a user only when one of its locations is
// a configured resource: that resource's binaries become the only writers of its reserved names
// (an unconfigured app has no writers, and reserving would lock its own installer out). It also
// returns the users whose Bun tmp dir is a reserved root (trustManager.reserveGlobs vets it).
func buildGlobReservations(cfg *daemonconfig.Config, users []install.User) (guard.GlobReservations, []install.User) {
	b := globBuilder{
		cfg:      cfg,
		r:        guard.GlobReservations{Roots: map[string]uint64{}, Writers: map[string]uint64{}},
		bitOf:    map[string]int{},
		children: map[[2]string]uint64{},
	}
	for _, u := range users {
		for i := range install.Catalog {
			if e := &install.Catalog[i]; !e.IsSystem() {
				b.addEntry(e, u)
			}
		}
		b.reserveBunTmp(u)
	}
	for k, bits := range b.children {
		b.r.Children = append(b.r.Children, guard.GlobChild{Parent: k[0], Name: k[1], Bits: bits})
	}
	sort.Slice(b.r.Children, func(i, j int) bool {
		if b.r.Children[i].Parent != b.r.Children[j].Parent {
			return b.r.Children[i].Parent < b.r.Children[j].Parent
		}
		return b.r.Children[i].Name < b.r.Children[j].Name
	})
	return b.r, b.bun
}

type globBuilder struct {
	cfg      *daemonconfig.Config
	r        guard.GlobReservations
	bitOf    map[string]int
	children map[[2]string]uint64
	bun      []install.User
}

func (b *globBuilder) addEntry(e *install.CandidateDir, u install.User) {
	writers := entryWriters(b.cfg, e, u)
	if len(writers) == 0 {
		return
	}
	for _, g := range e.TrustGlobs(u.Name, u.Home) {
		b.reserve(g, g.Name, writers)
	}
	b.reserveFixedLibDirs(e, u, writers)
	for _, g := range e.LibDirGlobs(u.Name, u.Home) {
		b.reserve(g.TrustGlob, g.Name, writers)
		b.reserveLibDirTail(e, &g, writers)
	}
	for _, g := range e.FixedBinaryGlobs(u.Name, u.Home) {
		b.reserve(g, g.Name, writers)
		if t, ok := symlinkTargetGlob(g); ok {
			b.reserve(t, t.Name, writers)
		}
	}
	// Library bits are per entry: the bit also grants loading (trust_mmap), and a "*.so" shared
	// across entries would let one app's writers plant, and load, below another app's root.
	for _, g := range e.ReservedLibGlobs(u.Name, u.Home) {
		b.reserve(g, e.Name+"\x00"+g.Name, writers)
	}
	// Per entry too: another app's binaries must not write this one's inputs.
	for _, g := range e.ReservedNameGlobs(u.Name, u.Home) {
		b.reserve(g, e.Name+"\x00name\x00"+g.Name, writers)
	}
}

// reserveBunTmp reserves the shared Bun extraction name (.bun-*) under ~/.cache/app-listener/bun for
// every configured Bun app's binaries, all under ONE shared bit: a launcher wrapper points $TMPDIR
// there (install), so a Bun app loads its own — and, harmlessly, another trusted Bun app's —
// per-launch native library, while a non-writer can neither plant nor load one. Unlike per-entry
// ReservedLibs the bit is shared: several apps write .bun-* into the same dir, so per-entry bits
// would make each app a non-writer of the others' extractions and deny creation. The dir must exist
// (the installer creates it only when the user opts in); a missing dir reserves nothing. An existing
// one may predate its reservation, so it is listed in b.bun for reserveGlobs to vet.
func (b *globBuilder) reserveBunTmp(u install.User) {
	seen := map[string]bool{}
	var writers []string
	for i := range install.Catalog {
		e := &install.Catalog[i]
		if !e.IsBun() {
			continue
		}
		for _, w := range entryWriters(b.cfg, e, u) {
			if !seen[w] {
				seen[w] = true
				writers = append(writers, w)
			}
		}
	}
	if len(writers) == 0 {
		return
	}
	sort.Strings(writers)
	g := install.TrustGlob{
		Home:  u.Home,
		Fixed: strings.Split(install.BunTmpRelDir, "/"),
		Name:  install.BunReservedName,
		Lib:   true,
	}
	b.reserve(g, "\x00bun-shared", writers)
	if b.r.Roots[g.Root()] != 0 {
		b.bun = append(b.bun, u)
	}
}

// reserveFixedLibDirs reserves each wildcard-free lib dir of e, present or not, as exact names along
// its path for e's writers: the catalog refresh adopts such a dir once it exists and trusts what it
// holds for those writers, so a non-writer must not create it (or recreate it after an rmdir).
// Exact children, not a root: a root reserves at any depth, and would make every game dir named
// linux64 below steamapps writer-only.
func (b *globBuilder) reserveFixedLibDirs(e *install.CandidateDir, u install.User, writers []string) {
	dirs := e.FixedLibDirs(u.Name, u.Home)
	if len(dirs) == 0 {
		return
	}
	mask := b.libDirBit(e, dirs[0][len(dirs[0])-1], writers)
	for _, comps := range dirs {
		b.reserveExactPath(u.Home, comps, mask)
	}
}

// reserveLibDirTail reserves g's fixed tail (Proton*/files/lib) below every existing match of its
// wildcard, as for a fixed lib dir: the root reserves only the wildcard name, so in a match lacking
// the tail a non-writer could create it and the refresh would adopt the plant. A match made after
// this reload is covered from the next one.
func (b *globBuilder) reserveLibDirTail(e *install.CandidateDir, g *install.LibDirGlob, writers []string) {
	if len(g.Tail) == 0 {
		return
	}
	matches, err := filepath.Glob(filepath.Join(g.Root(), g.Name))
	if err != nil {
		log.Warnf("trust guard: lib dir pattern %s/%s is malformed (%v); its tail is not reserved", g.Root(),
			g.Name, err)
		return
	}
	mask := b.libDirBit(e, g.Tail[len(g.Tail)-1], writers)
	for _, m := range matches {
		if fi, err := os.Stat(m); err == nil && fi.IsDir() {
			b.reserveExactPath(m, g.Tail, mask)
		}
	}
}

// libDirBit is e's lib dir bit, granted to writers. One bit per entry; its pattern name matches no
// root, so it reserves only exact children.
func (b *globBuilder) libDirBit(e *install.CandidateDir, name string, writers []string) uint64 {
	mask := b.bit(e.Name+"\x00libdirs", name)
	for _, w := range writers {
		b.r.Writers[w] |= mask
	}
	return mask
}

// reserveExactPath reserves comps below dir as exact children, each in its parent, up to and
// including the first one missing.
func (b *globBuilder) reserveExactPath(dir string, comps []string, mask uint64) {
	for _, c := range comps {
		if underResource(b.cfg, dir) {
			return // that resource's guard already gates creating anything below it
		}
		if _, _, ok := guard.ParseGlobName(c); !ok || strings.ContainsAny(c, "*?[") {
			log.Warnf("trust guard: lib dir component %q below %s cannot be reserved; a same-user process "+
				"could create it and the catalog refresh would trust its libraries", c, dir)
			return
		}
		b.children[[2]string{dir, c}] |= mask
		next := filepath.Join(dir, c)
		if fi, err := os.Stat(next); err != nil || !fi.IsDir() {
			return
		}
		dir = next
	}
}

func (b *globBuilder) reserve(g install.TrustGlob, key string, writers []string) {
	if underResource(b.cfg, g.Root()) {
		return // that resource's own guard already gates every write below it
	}
	if _, _, ok := guard.ParseGlobName(g.Name); !ok {
		if g.Lib {
			log.Warnf("trust guard: library pattern %s/%s cannot be reserved (only exact, prefix* and "+
				"*suffix names): those libraries stay untrusted", g.Root(), g.Name)
		} else {
			log.Warnf("trust guard: %s/%s cannot be reserved (only exact, prefix* and *suffix names): "+
				"a same-user plant matching it would be trusted at the next catalog refresh", g.Root(), g.Name)
		}
		return
	}
	mask := b.bit(key, g.Name)
	reserveChain(g, mask, b.r.Roots, b.children)
	for _, w := range writers {
		b.r.Writers[w] |= mask
	}
}

// bit returns key's mask for pattern name, assigning the next bit on first use. Past 64 keys the
// masks wrap, and SetGlobReservations refuses the set.
func (b *globBuilder) bit(key, name string) uint64 {
	i, ok := b.bitOf[key]
	if !ok {
		i = len(b.r.Patterns)
		b.bitOf[key] = i
		b.r.Patterns = append(b.r.Patterns, name)
	}
	return uint64(1) << (i % 64)
}

// reserveChain reserves g.Name below the deepest existing directory on the way to g.Root, so a root
// created later is covered from its parent until the next reload. Each existing fixed component is
// also reserved in its parent: roots are inode-keyed, and a renamed-away root recreated by a
// non-writer would otherwise be unregistered. A library or input root never falls back: "lib*" or
// "settings.json" below all of ~/.config would lock every other app out of those names.
func reserveChain(g install.TrustGlob, mask uint64, roots map[string]uint64, children map[[2]string]uint64) {
	dir := g.Home
	for _, c := range g.Fixed {
		next := filepath.Join(dir, c)
		if fi, err := os.Stat(next); err != nil || !fi.IsDir() {
			break
		}
		if _, _, ok := guard.ParseGlobName(c); ok { // Fixed holds no wildcards: an exact name
			children[[2]string{dir, c}] |= mask
		} else {
			log.Warnf("trust guard: path component %q of %s cannot be reserved; renaming it away "+
				"would unregister the glob root until the next reload", c, g.Root())
		}
		dir = next
	}
	if (g.Lib || g.Input) && dir != g.Root() {
		log.Infof("trust guard: %s is missing, its %s files are not reserved", g.Root(), g.Name)
		return
	}
	roots[dir] |= mask
}

// symlinkTargetGlob is the in-home file a fixed binary's symlink resolves to (~/.local/bin/claude
// -> versions/X): the daemon whitelists the target, so its directories need the same reservation.
func symlinkTargetGlob(g install.TrustGlob) (install.TrustGlob, bool) {
	link := filepath.Join(g.Root(), g.Name)
	target, err := filepath.EvalSymlinks(link)
	if err != nil || target == link {
		return install.TrustGlob{}, false
	}
	rel, ok := strings.CutPrefix(target, g.Home+"/")
	if !ok {
		return install.TrustGlob{}, false // outside the home: root-owned, protection #1 skips it
	}
	parts := strings.Split(rel, "/")
	return install.TrustGlob{Home: g.Home, Fixed: parts[:len(parts)-1], Name: parts[len(parts)-1]}, true
}

// entryWriters returns the whitelisted binaries of every configured resource belonging to e for u.
func entryWriters(cfg *daemonconfig.Config, e *install.CandidateDir, u install.User) []string {
	seen := make(map[string]bool)
	var out []string
	for i := range cfg.Resources {
		res := &cfg.Resources[i]
		if !e.OwnsPath(u.Name, u.Home, res.Path) && !e.OwnsPath(u.Name, u.Home, res.EncryptionRootOrPath()) {
			continue
		}
		for _, list := range [][]daemonconfig.BinaryRule{res.Binaries, res.PendingBinaries} {
			for _, b := range list {
				if !seen[b.Path] {
					seen[b.Path] = true
					out = append(out, b.Path)
				}
			}
		}
	}
	sort.Strings(out)
	return out
}

func underResource(cfg *daemonconfig.Config, dir string) bool {
	for i := range cfg.Resources {
		p := cfg.Resources[i].Path
		if dir == p || strings.HasPrefix(dir, p+"/") {
			return true
		}
	}
	return false
}
