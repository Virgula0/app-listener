package install

import (
	"path/filepath"
	"strings"
)

// ResolveCatalogEntry reverse-looks-up which Catalog entry an absolute path from the existing
// daemon.conf came from. AbsPaths entries match exactly; RelPaths entries via PathsFor for each
// known user (any location matches the shared whitelist); WatchRelPaths entries also match their
// watch sub-paths (~/.config/discord/Local Storage -> Discord), so a refresh re-expands every
// grouped section. Returns nil for a user-added section.
func ResolveCatalogEntry(resourcePath string, users []User) (*CandidateDir, *User) {
	if entry, user := findCatalogRoot(resourcePath, users); entry != nil {
		return entry, user
	}
	return findCatalogWatchSubPath(resourcePath, users)
}

// catalogEntryMatch returns the user whose expansion of entry contains a path satisfying match, or
// nil. System entries yield the zero user.
func catalogEntryMatch(entry *CandidateDir, users []User, match func(catalogPath string) bool) *User {
	if entry.IsSystem() {
		for _, p := range entry.PathsFor("", "") {
			if match(p) {
				return &User{}
			}
		}
		return nil
	}
	for j := range users {
		u := &users[j]
		for _, p := range entry.PathsFor(u.Home, u.Name) {
			if match(p) {
				return u
			}
		}
	}
	return nil
}

// findCatalogRoot matches one of the entry's own watch roots exactly.
func findCatalogRoot(resourcePath string, users []User) (*CandidateDir, *User) {
	for i := range Catalog {
		entry := &Catalog[i]
		if u := catalogEntryMatch(entry, users, func(p string) bool { return p == resourcePath }); u != nil {
			return entry, u
		}
	}
	return nil, nil
}

// findCatalogWatchSubPath matches grouped watch sub-paths: a resource path strictly inside a
// catalog root resolves to that entry, so refreshes re-expand every grouped section.
func findCatalogWatchSubPath(resourcePath string, users []User) (*CandidateDir, *User) {
	for i := range Catalog {
		entry := &Catalog[i]
		if u := catalogEntryMatch(entry, users, func(p string) bool { return isInsidePath(resourcePath, p) }); u != nil {
			return entry, u
		}
	}
	return nil, nil
}

func isInsidePath(path, dir string) bool {
	return path != dir && strings.HasPrefix(path+"/", dir+"/")
}

// SystemWhitelist returns the entry's absolute Whitelist patterns (fixed paths or globs, never
// %HOME%-relative), sorted, for the daemon's live admission of system binaries.
func (c *CandidateDir) SystemWhitelist() []BinaryRule {
	var out []BinaryRule
	for _, r := range c.ExpandWhitelist("", "") {
		if _, raw := c.Whitelist[r.Path]; raw && filepath.IsAbs(r.Path) {
			out = append(out, r) // unchanged by expansion: no %HOME%/%USER%
		}
	}
	return out
}
