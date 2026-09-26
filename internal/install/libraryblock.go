package install

import (
	"fmt"
	"strings"
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

// parseLibrariesHeaderName returns the name of a `[libraries ...]` header
// line, mirroring the daemon parser (bare, unquoted and quoted forms).
func parseLibrariesHeaderName(line string) (string, bool) {
	t := strings.TrimSpace(line)
	if !strings.HasPrefix(t, "[libraries") || !strings.HasSuffix(t, "]") {
		return "", false
	}
	rest := strings.TrimSuffix(strings.TrimPrefix(t, "[libraries"), "]")
	if rest != "" && rest[0] != ' ' && rest[0] != '\t' {
		return "", false
	}
	rest = strings.TrimSpace(rest)
	if len(rest) >= 2 && rest[0] == '"' && rest[len(rest)-1] == '"' {
		rest = rest[1 : len(rest)-1]
	}
	return rest, true
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

// IsLibraryDirective: a trimmed config line is a library-trust directive (allow_lib, lib_dir,
// lib_binary): section structure, never a whitelist entry.
func IsLibraryDirective(trimmed string) bool {
	return isLibDirectiveText(trimmed)
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
