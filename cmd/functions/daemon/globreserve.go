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
// (an unconfigured app has no writers, and reserving would lock its own installer out).
func buildGlobReservations(cfg *daemonconfig.Config, users []install.User) guard.GlobReservations {
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
	return b.r
}

type globBuilder struct {
	cfg      *daemonconfig.Config
	r        guard.GlobReservations
	bitOf    map[string]int
	children map[[2]string]uint64
}

func (b *globBuilder) addEntry(e *install.CandidateDir, u install.User) {
	writers := entryWriters(b.cfg, e, u)
	if len(writers) == 0 {
		return
	}
	for _, g := range e.TrustGlobs(u.Name, u.Home) {
		if underResource(b.cfg, g.Root()) {
			continue // that resource's own guard already gates every write below it
		}
		if _, _, ok := guard.ParseGlobName(g.Name); !ok {
			log.Warnf("trust guard: %s/%s cannot be reserved (only exact, prefix* and *suffix names): "+
				"a same-user plant matching it would be trusted at the next catalog refresh", g.Root(), g.Name)
			continue
		}
		mask := b.bit(g.Name)
		reserveChain(g, mask, b.r.Roots, b.children)
		for _, w := range writers {
			b.r.Writers[w] |= mask
		}
	}
}

// bit returns name's mask, assigning the next bit on first use. Past 64 names the masks wrap, and
// SetGlobReservations refuses the set.
func (b *globBuilder) bit(name string) uint64 {
	i, ok := b.bitOf[name]
	if !ok {
		i = len(b.r.Patterns)
		b.bitOf[name] = i
		b.r.Patterns = append(b.r.Patterns, name)
	}
	return uint64(1) << (i % 64)
}

// reserveChain reserves g.Name below the deepest existing directory on the way to g.Root, so a root
// created later is covered from its parent until the next reload. Each existing fixed component is
// also reserved in its parent: roots are inode-keyed, and a renamed-away root recreated by a
// non-writer would otherwise be unregistered.
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
	roots[dir] |= mask
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
