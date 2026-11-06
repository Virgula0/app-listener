package install

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

// NarrowSection is one `[narrow <path>]` section: a subtree of its parent section judged by Allow,
// a subset of the parent's whitelist (daemonconfig.ParseNarrowSection).
type NarrowSection struct {
	Path  string
	Allow []BinaryRule
}

func (n NarrowSection) text() string {
	var b strings.Builder
	b.WriteString("\n[narrow ")
	b.WriteString(quotePath(n.Path))
	b.WriteString("]\n\n")
	writeAllow(&b, n.Allow)
	return b.String()
}

// NarrowSectionsFor is c's Narrowed subtrees of location that exist, each granted the lines of allow
// (the location's whitelist) its spec names. A missing one is skipped: daemon.conf refuses a
// [narrow] target that doesn't exist.
func (c *CandidateDir) NarrowSectionsFor(location string, allow []BinaryRule) []NarrowSection {
	var out []NarrowSection
	for _, rel := range slices.Sorted(maps.Keys(c.Narrowed)) {
		p := filepath.Join(location, rel)
		if fi, err := os.Lstat(p); err != nil || !fi.IsDir() || !ConfSafePath(p) {
			continue
		}
		out = append(out, NarrowSection{Path: p, Allow: narrowAllow(allow, c.Narrowed[rel])})
	}
	return out
}

// SeedNarrowed makes c's missing Narrowed directories of location (inside u's home) as u. One made
// later, by the app, would sit under the location's wider whitelist until a catalog refresh.
func (c *CandidateDir) SeedNarrowed(u User, location string) error {
	for _, rel := range slices.Sorted(maps.Keys(c.Narrowed)) {
		homeRel, err := filepath.Rel(u.Home, filepath.Join(location, rel))
		if err != nil || homeRel == ".." || strings.HasPrefix(homeRel, "../") {
			continue // a system location: nothing to make as the user
		}
		if err := seedPath(u, homeRel, ""); err != nil {
			return fmt.Errorf("seeding %s: %w", filepath.Join(location, rel), err)
		}
	}
	return nil
}

// narrowAllow keeps the lines of allow whose basename spec names, with spec's events where it
// lists some (intersected with the line's own: a narrow only takes access away).
func narrowAllow(allow []BinaryRule, spec map[string][]string) []BinaryRule {
	var out []BinaryRule
	for _, r := range allow {
		events, ok := spec[filepath.Base(r.Path)]
		if !ok {
			continue
		}
		if ev, keep := intersectEvents(r.Events, events); keep {
			out = append(out, BinaryRule{Path: r.Path, Events: ev})
		}
	}
	return out
}

// intersectEvents is the narrower of two event lists (nil = every event); keep false when they
// share none.
func intersectEvents(parent, child []string) ([]string, bool) {
	switch {
	case len(child) == 0:
		return slices.Clone(parent), true
	case len(parent) == 0:
		return slices.Clone(child), true
	}
	var out []string
	for _, e := range child {
		if slices.ContainsFunc(parent, func(p string) bool { return strings.EqualFold(p, e) }) {
			out = append(out, e)
		}
	}
	return out, len(out) > 0
}

// SyncNarrows brings the [narrow] sections inside sectionPath in line with its new whitelist fresh,
// so a refresh never leaves daemon.conf refusing to load: entry's own narrows are regenerated from
// fresh (added after the section when their directory exists), any other keeps only the lines fresh
// still grants, events narrowed to fresh's.
func SyncNarrows(confText, sectionPath string, entry *CandidateDir, fresh []BinaryRule) string {
	own := map[string]NarrowSection{}
	if entry != nil {
		for _, n := range entry.NarrowSectionsFor(sectionPath, fresh) {
			own[n.Path] = n
		}
	}
	lines := strings.Split(confText, "\n")
	var out []string
	for i := 0; i < len(lines); {
		p, ok, err := daemonconfig.ParseNarrowSection(strings.TrimSpace(lines[i]))
		if !ok || err != nil || !isInsidePath(p, sectionPath) {
			out = append(out, lines[i])
			i++
			continue
		}
		end := findSectionEnd(lines, i)
		if n, isOwn := own[p]; isOwn {
			delete(own, p)
			out = append(out, strings.Split(strings.TrimPrefix(n.text(), "\n"), "\n")...)
			out = trimTrailingBlank(out)
			out = append(out, "")
		} else {
			out = append(out, lines[i])
			for _, l := range lines[i+1 : end] {
				if kept, keep := pruneNarrowLine(l, fresh); keep {
					out = append(out, kept)
				}
			}
		}
		i = end
	}
	text := strings.Join(out, "\n")
	if len(own) == 0 {
		return text
	}
	return insertAfterSection(text, sectionPath, own)
}

// pruneNarrowLine keeps a [narrow] body line: comments and blanks as they are, a binary line only
// while fresh lists it, with events narrowed to fresh's.
func pruneNarrowLine(line string, fresh []BinaryRule) (string, bool) {
	t := strings.TrimSpace(line)
	if t == "" || strings.HasPrefix(t, "#") {
		return line, true
	}
	path, rest := splitAllowLine(t)
	i := slices.IndexFunc(fresh, func(r BinaryRule) bool { return r.Path == path })
	if i < 0 {
		return "", false
	}
	var events []string
	if rest != "" {
		events = strings.Split(rest, ",")
		for j := range events {
			events[j] = strings.TrimSpace(events[j])
		}
	}
	ev, keep := intersectEvents(fresh[i].Events, events)
	if !keep {
		return "", false
	}
	var b strings.Builder
	writeAllow(&b, []BinaryRule{{Path: path, Events: ev}})
	return keepIndent(line, strings.TrimSuffix(b.String(), "\n")), true
}

func splitAllowLine(t string) (path, rest string) {
	if strings.HasPrefix(t, `"`) {
		if end := strings.IndexByte(t[1:], '"'); end != -1 {
			return t[1 : 1+end], strings.TrimSpace(t[2+end:])
		}
	}
	fields := strings.Fields(t)
	return fields[0], strings.TrimSpace(t[len(fields[0]):])
}

// insertAfterSection puts the narrows right after the [watch sectionPath] body: a [narrow] must
// follow its parent.
func insertAfterSection(text, sectionPath string, narrows map[string]NarrowSection) string {
	lines := strings.Split(text, "\n")
	start := findSectionStart(lines, sectionPath)
	if start < 0 {
		return text
	}
	end := findSectionEnd(lines, start)
	var add []string
	for _, p := range slices.Sorted(maps.Keys(narrows)) {
		add = append(add, strings.Split(strings.TrimSuffix(narrows[p].text(), "\n"), "\n")...)
	}
	head := trimTrailingBlank(slices.Clone(lines[:end]))
	merged := slices.Concat(head, add, []string{""}, lines[end:])
	return strings.Join(merged, "\n")
}

func trimTrailingBlank(lines []string) []string {
	for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}
