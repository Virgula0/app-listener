package install

import (
	"path/filepath"
	"sort"
	"strings"
)

// TrustGlob is a Whitelist/LibDirWriters glob under a home whose matches the catalog refresh turns
// into trust grants. Fixed are the wildcard-free components between Home and the first wildcard
// one; Name is the last component. The daemon reserves Name below Root (guard_trust.bpf.c #3).
// Lib marks a ReservedLibs pattern.
type TrustGlob struct {
	Home  string
	Fixed []string
	Name  string
	Lib   bool
}

// Root is the glob's leading wildcard-free directory.
func (g TrustGlob) Root() string {
	return filepath.Join(append([]string{g.Home}, g.Fixed...)...)
}

// TrustGlobs returns the entry's wildcarded %HOME% Whitelist and LibDirWriters patterns for one
// user, sorted and de-duplicated. Absolute patterns are skipped (their trees are root-owned), and
// so are fixed paths: one path is trusted on first use, then protection #1 guards it.
func (c *CandidateDir) TrustGlobs(user, home string) []TrustGlob {
	pats := make([]string, 0, len(c.Whitelist)+len(c.LibDirWriters))
	for p := range c.Whitelist {
		pats = append(pats, p)
	}
	pats = append(pats, c.LibDirWriters...)
	sort.Strings(pats)

	seen := make(map[string]bool)
	var out []TrustGlob
	for _, p := range pats {
		rel, ok := strings.CutPrefix(p, "%HOME%/")
		if !ok || !strings.ContainsAny(rel, "*?[") || seen[p] {
			continue
		}
		seen[p] = true
		out = append(out, splitTrustGlob(home, expandPlaceholders(rel, user, home)))
	}
	return out
}

// ReservedLibGlobs returns the entry's %HOME% ReservedLibs patterns for one user, with Lib set.
func (c *CandidateDir) ReservedLibGlobs(user, home string) []TrustGlob {
	out := make([]TrustGlob, 0, len(c.ReservedLibs))
	for _, p := range c.ReservedLibs {
		if rel, ok := strings.CutPrefix(p, "%HOME%/"); ok {
			g := splitTrustGlob(home, expandPlaceholders(rel, user, home))
			g.Lib = true
			out = append(out, g)
		}
	}
	return out
}

// LibDirGlobs returns, for each wildcarded %HOME%-relative LibDirRelPaths pattern, its first
// wildcard component as the reserved Name below the fixed prefix: the catalog refresh guards every
// new match, so only the entry's writers may create one (and nothing planted predates the guard).
func (c *CandidateDir) LibDirGlobs(user, home string) []TrustGlob {
	var out []TrustGlob
	for _, rel := range c.LibDirRelPaths {
		parts := strings.Split(expandPlaceholders(rel, user, home), "/")
		for i, part := range parts {
			if strings.ContainsAny(part, "*?[") {
				out = append(out, TrustGlob{Home: home, Fixed: parts[:i], Name: part})
				break
			}
		}
	}
	return out
}

func splitTrustGlob(home, rel string) TrustGlob {
	parts := strings.Split(rel, "/")
	g := TrustGlob{Home: home, Name: parts[len(parts)-1]}
	for _, part := range parts[:len(parts)-1] {
		if strings.ContainsAny(part, "*?[") {
			break
		}
		g.Fixed = append(g.Fixed, part)
	}
	return g
}

// OwnsPath reports whether path is one of the entry's guarded locations for this user: a RelPaths,
// WatchRelPaths or LibDirRelPaths entry (globs matched, not expanded against the filesystem).
func (c *CandidateDir) OwnsPath(user, home, path string) bool {
	pats := c.PathsFor(home, user)
	for _, rel := range c.WatchRelPaths {
		pats = append(pats, filepath.Join(home, expandPlaceholders(rel, user, home)))
	}
	for _, rel := range c.LibDirRelPaths {
		pats = append(pats, filepath.Join(home, expandPlaceholders(rel, user, home)))
	}
	for _, p := range pats {
		if ok, err := filepath.Match(p, path); err == nil && ok {
			return true
		}
	}
	return false
}
