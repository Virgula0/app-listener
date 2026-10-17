//go:build gui

package gui

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// pickEntry is one row of the picker.
type pickEntry struct {
	path string
	// rel is the name shown: the base name, or for a subfolder search the path below the searched
	// folder.
	rel string
	dir bool
}

type listOpts struct {
	hidden      bool
	foldersOnly bool
}

const (
	pickMaxResults = 1000
	pickMaxVisited = 200_000
)

// pseudoFS are never walked by a subfolder search: endless or huge, and nothing to pick there.
var pseudoFS = map[string]bool{"/proc": true, "/sys": true, "/dev": true}

// listDir lists dir, folders first then by name; a symlink counts as what it points to.
func listDir(dir string, opt listOpts) ([]pickEntry, error) {
	des, err := os.ReadDir(dir)
	out := make([]pickEntry, 0, len(des))
	for _, d := range des {
		name := d.Name()
		if !opt.hidden && strings.HasPrefix(name, ".") {
			continue
		}
		p := filepath.Join(dir, name)
		isDir := entryIsDir(p, d)
		if opt.foldersOnly && !isDir {
			continue
		}
		out = append(out, pickEntry{path: p, rel: name, dir: isDir})
	}
	sortEntries(out)
	return out, err
}

func entryIsDir(p string, d fs.DirEntry) bool {
	if d.Type()&fs.ModeSymlink == 0 {
		return d.IsDir()
	}
	fi, err := os.Stat(p)
	return err == nil && fi.IsDir()
}

// searchTerms splits a query into lower-case words; a name matches when it contains every one.
func searchTerms(query string) []string {
	return strings.Fields(strings.ToLower(query))
}

func matchesTerms(name string, terms []string) bool {
	name = strings.ToLower(name)
	for _, t := range terms {
		if !strings.Contains(name, t) {
			return false
		}
	}
	return true
}

func filterEntries(es []pickEntry, query string) []pickEntry {
	terms := searchTerms(query)
	if len(terms) == 0 {
		return es
	}
	out := make([]pickEntry, 0, len(es))
	for _, e := range es {
		if matchesTerms(e.rel, terms) {
			out = append(out, e)
		}
	}
	return out
}

// searchTree finds the names below root matching query, shallowest first. It follows no symlink
// and stops at pickMaxResults matches or pickMaxVisited entries (truncated), or once stop says so.
func searchTree(root, query string, opt listOpts, stop func() bool) (res []pickEntry, truncated bool) {
	t := &treeSearch{root: root, terms: searchTerms(query), opt: opt, stop: stop}
	_ = filepath.WalkDir(root, t.visit)
	sortEntries(t.res)
	return t.res, t.truncated
}

type treeSearch struct {
	root      string
	terms     []string
	opt       listOpts
	stop      func() bool
	visited   int
	res       []pickEntry
	truncated bool
}

func (t *treeSearch) visit(p string, d fs.DirEntry, err error) error {
	switch {
	case t.stop():
		return fs.SkipAll
	case err != nil: // unreadable (a guarded folder, say): left out, not fatal
		return skipEntry(d)
	case p == t.root:
		return nil
	}
	t.visited++
	if t.visited > pickMaxVisited {
		t.truncated = true
		return fs.SkipAll
	}
	name := d.Name()
	if (!t.opt.hidden && strings.HasPrefix(name, ".")) || (d.IsDir() && pseudoFS[p]) {
		return skipEntry(d)
	}
	if !matchesTerms(name, t.terms) {
		return nil
	}
	isDir := entryIsDir(p, d)
	if t.opt.foldersOnly && !isDir {
		return nil
	}
	rel, _ := filepath.Rel(t.root, p)
	t.res = append(t.res, pickEntry{path: p, rel: rel, dir: isDir})
	if len(t.res) >= pickMaxResults {
		t.truncated = true
		return fs.SkipAll
	}
	return nil
}

func skipEntry(d fs.DirEntry) error {
	if d != nil && d.IsDir() {
		return fs.SkipDir
	}
	return nil
}

func sortEntries(es []pickEntry) {
	sort.SliceStable(es, func(i, j int) bool {
		a, b := es[i], es[j]
		if da, db := strings.Count(a.rel, "/"), strings.Count(b.rel, "/"); da != db {
			return da < db
		}
		if a.dir != b.dir {
			return a.dir
		}
		return strings.ToLower(a.rel) < strings.ToLower(b.rel)
	})
}
