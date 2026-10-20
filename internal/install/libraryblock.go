package install

import (
	"fmt"
	"strings"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

// LibraryBlock is one `[libraries "<name>"]` block: the library-trust directives of ONE
// application, declared once rather than under one of its watch sections.
//
// Library trust is daemon-wide (every allow_lib and lib_dir is trusted for every whitelisted
// binary, whichever section named it), so nesting under a watch path implied a per-resource scoping
// that never existed. What IS scoped is who may WRITE a lib_dir, hence one block per application: a
// block's lib_binary rules write only that block's lib_dirs, so one app's updater can't write
// another's runtime tree.
type LibraryBlock struct {
	// Name labels the application ("Steam (alice)"); the daemon scopes the
	// block by position, the installer finds it again by this name.
	Name       string
	Libs       []string // allow_lib
	LibDirs    []string // lib_dir
	LibWriters []string // lib_binary
}

// Empty reports whether the block would render no directive at all.
func (b *LibraryBlock) Empty() bool {
	return len(b.Libs)+len(b.LibDirs)+len(b.LibWriters) == 0
}

// DirectiveLines renders the block's directives in documented order (allowed library files, library
// dirs, then their writers), exactly as GenerateLibraryBlocks writes them (the installer matches
// config lines against this text).
func (b *LibraryBlock) DirectiveLines() []string {
	out := make([]string, 0, len(b.Libs)+len(b.LibDirs)+len(b.LibWriters))
	for _, d := range [...]struct {
		keyword string
		values  []string
	}{
		{"allow_lib", b.Libs},
		{"lib_dir", b.LibDirs},
		{"lib_binary", b.LibWriters},
	} {
		for _, v := range d.values {
			out = append(out, d.keyword+" "+quotePath(v))
		}
	}
	return out
}

func libraryBlockHeader(name string) string {
	if name == "" {
		return "[libraries]"
	}
	return "[libraries " + quotePath(name) + "]"
}

// GenerateLibraryBlocks renders the blocks that carry at least one directive,
// each preceded by a blank line, for appending after the watch sections.
func GenerateLibraryBlocks(blocks []LibraryBlock) string {
	var b strings.Builder
	for i := range blocks {
		if blocks[i].Empty() {
			continue
		}
		b.WriteString("\n")
		b.WriteString(libraryBlockHeader(blocks[i].Name))
		b.WriteString("\n\n")
		for _, line := range blocks[i].DirectiveLines() {
			b.WriteString(line)
			b.WriteString("\n")
		}
	}
	return b.String()
}

// parseLibrariesHeaderName returns the name of a `[libraries ...]` header line, as the daemon reads it.
func parseLibrariesHeaderName(line string) (string, bool) {
	return daemonconfig.ParseLibrariesSection(strings.TrimSpace(line))
}

// EnsureLibraryBlock merges block into confText: appended if no `[libraries <name>]` header exists,
// otherwise only its missing directive lines are added. Never removes or reorders a line, so
// operator hand edits survive catalog refreshes.
//
// Presence is checked within the block only: a lib_binary writes its OWN block's lib_dirs, so a
// lib_dir declared elsewhere (legacy nesting under a watch section) is still added here or this
// block's writers wouldn't reach it. The daemon guards a lib_dir named twice once.
func EnsureLibraryBlock(confText string, block *LibraryBlock) string {
	if block.Empty() {
		return confText
	}
	lines := strings.Split(confText, "\n")
	start := -1
	for i, l := range lines {
		if name, ok := parseLibrariesHeaderName(l); ok && name == block.Name {
			start = i
			break
		}
	}
	if start == -1 {
		return strings.TrimRight(confText, "\n") + "\n" + GenerateLibraryBlocks([]LibraryBlock{*block})
	}

	end := findSectionEnd(lines, start)
	present := make(map[string]bool, end-start)
	insertAt := start + 1
	for i := start + 1; i < end; i++ {
		t := strings.TrimSpace(lines[i])
		if t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		present[t] = true
		insertAt = i + 1
	}
	var add []string
	for _, line := range block.DirectiveLines() {
		if !present[line] {
			present[line] = true
			add = append(add, line)
		}
	}
	if len(add) == 0 {
		return confText
	}
	result := make([]string, 0, len(lines)+len(add))
	result = append(result, lines[:insertAt]...)
	result = append(result, add...)
	result = append(result, lines[insertAt:]...)
	return strings.Join(result, "\n")
}

// PruneLibraryBlock deletes from the `[libraries <name>]` block each lib_dir / lib_binary line whose
// path stale reports, returning those paths. allow_lib lines and other blocks are untouched.
func PruneLibraryBlock(confText, name string, stale func(keyword, path string) bool) (text string, dropped []string) {
	lines := strings.Split(confText, "\n")
	start := -1
	for i, l := range lines {
		if n, ok := parseLibrariesHeaderName(l); ok && n == name {
			start = i
			break
		}
	}
	if start == -1 {
		return confText, nil
	}
	end := findSectionEnd(lines, start)
	result := make([]string, 0, len(lines))
	for i, l := range lines {
		if i > start && i < end {
			if keyword, path, ok := parseLibDirLine(strings.TrimSpace(l)); ok && stale(keyword, path) {
				dropped = append(dropped, path)
				continue
			}
		}
		result = append(result, l)
	}
	if len(dropped) == 0 {
		return confText, nil
	}
	return strings.Join(result, "\n"), dropped
}

// parseLibDirLine splits a trimmed lib_dir / lib_binary line as the daemon parser reads it.
func parseLibDirLine(t string) (keyword, path string, ok bool) {
	for _, d := range [...]string{"lib_dir", "lib_binary"} {
		rest, found := daemonconfig.CutDirective(t, d)
		if !found {
			continue
		}
		path = daemonconfig.UnquotePath(rest)
		if path == "" || strings.Contains(path, `"`) {
			return "", "", false
		}
		return d, path, true
	}
	return "", "", false
}

// RemoveSectionLibDirectives deletes the library directive lines in drop (trimmed comparison) from
// the [watch path] section body, migrating legacy configs (directives nested under a watch section)
// to the [libraries] layout: the installer drops the catalog-generated lines and EnsureLibraryBlock
// re-adds them in the app's block. Hand-written lines stay.
func RemoveSectionLibDirectives(confText, path string, drop []string) (string, error) {
	lines := strings.Split(confText, "\n")
	start := findSectionStart(lines, path)
	if start == -1 {
		return "", fmt.Errorf("section for %q not found in configuration", path)
	}
	end := findSectionEnd(lines, start)
	dropSet := make(map[string]bool, len(drop))
	for _, d := range drop {
		dropSet[d] = true
	}
	result := make([]string, 0, len(lines))
	for i, l := range lines {
		if i > start && i < end && dropSet[strings.TrimSpace(l)] {
			continue
		}
		result = append(result, l)
	}
	return strings.Join(result, "\n"), nil
}

// InsertSectionsBeforeLibraries adds rendered watch sections before the first [libraries] block (or
// at the end), keeping the layout: all watch sections first, library blocks after.
func InsertSectionsBeforeLibraries(confText, sections string) string {
	lines := strings.Split(confText, "\n")
	for i, l := range lines {
		if _, ok := parseLibrariesHeaderName(l); ok {
			head := strings.TrimRight(strings.Join(lines[:i], "\n"), "\n")
			return head + "\n" + sections + "\n" + strings.Join(lines[i:], "\n")
		}
	}
	return strings.TrimRight(confText, "\n") + "\n" + sections + "\n"
}
