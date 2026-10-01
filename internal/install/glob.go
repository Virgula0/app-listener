package install

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
)

// maxGlobEntries bounds the directory entries one pattern's expansion reads: any process may fill a
// glob's wildcard directory, and the daemon expands catalog globs on every refresh.
var maxGlobEntries = 20000

// boundedGlob is filepath.Glob reading at most maxGlobEntries directory entries; past that it
// returns what it matched, with a warning. Matches are sorted.
func boundedGlob(pattern string) ([]string, error) {
	if _, err := filepath.Match(pattern, ""); err != nil {
		return nil, err
	}
	if !strings.ContainsAny(pattern, "*?[") {
		if _, err := os.Lstat(pattern); err != nil {
			return nil, nil //nolint:nilerr // like filepath.Glob: a missing path matches nothing
		}
		return []string{pattern}, nil
	}
	budget := maxGlobEntries
	cands := []string{"/"}
	if !filepath.IsAbs(pattern) {
		cands = []string{"."}
	}
	for _, comp := range strings.Split(strings.TrimPrefix(filepath.Clean(pattern), "/"), "/") {
		var next []string
		for _, dir := range cands {
			if !strings.ContainsAny(comp, "*?[") {
				if p := filepath.Join(dir, comp); exists(p) {
					next = append(next, p)
				}
				continue
			}
			next = append(next, matchDir(dir, comp, &budget)...)
		}
		cands = next
	}
	if budget <= 0 {
		log.Warnf("catalog: expanding %s stopped after %d directory entries; matches past them are not "+
			"admitted", pattern, maxGlobEntries)
	}
	sort.Strings(cands)
	return cands, nil
}

func exists(p string) bool {
	_, err := os.Lstat(p)
	return err == nil
}

// matchDir returns dir's entries matching comp, reading at most *budget entries.
func matchDir(dir, comp string, budget *int) []string {
	f, err := os.Open(dir)
	if err != nil {
		return nil
	}
	defer f.Close()
	var out []string
	for *budget > 0 {
		n := min(*budget, 512)
		entries, err := f.ReadDir(n)
		*budget -= len(entries)
		for _, e := range entries {
			if ok, _ := filepath.Match(comp, e.Name()); ok {
				out = append(out, filepath.Join(dir, e.Name()))
			}
		}
		if err != nil || len(entries) < n {
			if err != nil && !errors.Is(err, io.EOF) {
				return out
			}
			break
		}
	}
	return out
}
