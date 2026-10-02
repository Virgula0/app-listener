package install

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

// ErrNoCatalogMatch: the config has no section mapping to a catalog entry (all-manual).
// Interactive callers surface it; automated ones treat it as "nothing to refresh".
var ErrNoCatalogMatch = errors.New("no catalog entries match any configured directory")

// SectionScan re-expands one catalog section's whitelist through expand and reports whether its
// vault is encrypted. The CLI wraps it to unlock a locked vault around the scan.
type SectionScan func(r *daemonconfig.Resource, expand func() []BinaryRule) (fresh []BinaryRule, encrypted bool, err error)

// SectionChange is one refreshed section's whitelist delta.
type SectionChange struct {
	Section  string
	Admitted []string
	Dropped  []string
}

// RefreshOptions tune RefreshCatalog.
type RefreshOptions struct {
	// Live: the vaults are unlocked and guarded (the running daemon's view).
	Live bool
	// KeepExisting keeps a section's whitelist lines whose file still exists, dropping only
	// vanished ones (the daemon's own refresh). Off, the section is rebuilt from the catalog alone,
	// so a line the catalog doesn't produce is dropped (`install --update-catalog-only`).
	KeepExisting bool
	Scan         SectionScan
	// Admit, when set, vets each line the refresh would add to a section (one it doesn't list yet).
	Admit func(u User, b BinaryRule) bool
}

// RefreshCatalog re-expands the whitelist of every catalog section of confText (the text cfg was
// parsed from). What it admits is what `install --update-catalog-only` admits: catalog globs,
// confSafeMatch, homeMatchConfined, and library directives only ever added. User sections are
// kept verbatim.
func RefreshCatalog(confText string, cfg *daemonconfig.Config, users []User, opts RefreshOptions) (string, []SectionChange, error) {
	patched := 0
	var changes []SectionChange
	// Encryption groups, not resources: a grouped section is one [watch <root>] header shared by its
	// watch paths, so its whitelist is re-expanded once.
	for _, r := range cfg.EncryptionGroups() {
		text, change, ok, err := RefreshSection(confText, r, users, opts)
		if err != nil {
			return "", nil, err
		}
		if !ok {
			continue
		}
		confText = text
		patched++
		if len(change.Admitted)+len(change.Dropped) > 0 {
			changes = append(changes, change)
		}
	}
	if patched == 0 {
		return "", nil, ErrNoCatalogMatch
	}
	return confText, changes, nil
}

// RefreshSection re-expands one resource's section if it matches a catalog entry; false for a
// user-added section. A live re-scan that comes back EMPTY for an encrypted section that listed
// binaries is refused (vault locked, or this binary isn't the running daemon's): persisting it
// would silently shrink the whitelist.
func RefreshSection(confText string, r *daemonconfig.Resource, users []User,
	opts RefreshOptions) (text string, change SectionChange, patched bool, err error) {
	// The SECTION path is what the text helpers address: the group root for grouped sections.
	sectionPath := r.EncryptionRootOrPath()
	change = SectionChange{Section: sectionPath}
	entry, user := ResolveCatalogEntry(sectionPath, users)
	if entry == nil {
		log.Infof("keeping user section as-is: %s", sectionPath)
		return confText, change, false, nil
	}

	candidate := Candidate{User: *user, Entry: *entry, Path: sectionPath}
	old := ParseSectionWhitelist(confText, sectionPath)
	fresh, encrypted, err := opts.Scan(r, candidate.FilterExistingWhitelist)
	if err != nil {
		return "", change, false, err
	}
	if opts.Admit != nil {
		fresh = admitNew(fresh, old, *user, opts.Admit)
	}
	if opts.KeepExisting {
		fresh = keepExisting(fresh, r)
	}
	log.Infof("re-scanned %s (%s) — %d whitelisted binaries", sectionPath, entry.Name, len(fresh))
	if LiveEmptyWhitelistRejected(encrypted, opts.Live, len(fresh), len(old)) {
		return "", change, false, fmt.Errorf("live re-scan of %s produced an empty whitelist while the config "+
			"lists binaries for it: the vault appears locked or this installer is not the running daemon's "+
			"binary — refusing to shrink the whitelist", sectionPath)
	}

	updated, err := SetSectionWhitelist(confText, sectionPath, fresh)
	if err != nil {
		return "", change, false, fmt.Errorf("patching section %s: %w", sectionPath, err)
	}
	// Library directives are only ever ADDED, in the app's own [libraries] block. Catalog-generated
	// ones an older installer nested under this section move there first (hand-written ones stay).
	// When the catalog moves a binary from the whitelist to the lib writers, only this re-adds it.
	block := entry.LibraryBlockFor(user.Name, user.Home)
	migrated, err := RemoveSectionLibDirectives(updated, sectionPath, block.DirectiveLines())
	if err != nil {
		return "", change, false, fmt.Errorf("moving the library directives of %s: %w", sectionPath, err)
	}
	freshPaths := make([]string, 0, len(fresh))
	for _, b := range fresh {
		freshPaths = append(freshPaths, b.Path)
	}
	change.Admitted = missingFrom(freshPaths, old)
	change.Dropped = missingFrom(old, freshPaths)
	return EnsureLibraryBlock(migrated, &block), change, true, nil
}

// keepExisting adds to fresh every whitelist line of r whose file still exists, events kept.
func keepExisting(fresh []BinaryRule, r *daemonconfig.Resource) []BinaryRule {
	have := make(map[string]bool, len(fresh))
	for _, b := range fresh {
		have[b.Path] = true
	}
	for _, list := range [][]daemonconfig.BinaryRule{r.Binaries, r.PendingBinaries} {
		for _, b := range list {
			if b.LibBinary || have[b.Path] {
				continue
			}
			if _, err := os.Stat(b.Path); err != nil {
				continue
			}
			have[b.Path] = true
			rule := BinaryRule{Path: b.Path}
			for _, e := range b.Events {
				rule.Events = append(rule.Events, e.String())
			}
			fresh = append(fresh, rule)
		}
	}
	sort.Slice(fresh, func(i, j int) bool { return fresh[i].Path < fresh[j].Path })
	return fresh
}

// admitNew drops each line of fresh that old doesn't list and admit refuses.
func admitNew(fresh []BinaryRule, old []string, u User, admit func(User, BinaryRule) bool) []BinaryRule {
	var out []BinaryRule
	for _, b := range fresh {
		if slices.Contains(old, b.Path) || admit(u, b) {
			out = append(out, b)
		}
	}
	return out
}

func missingFrom(a, b []string) []string {
	var out []string
	for _, p := range a {
		if !slices.Contains(b, p) {
			out = append(out, p)
		}
	}
	return out
}

// ParseSectionWhitelist extracts the binary paths listed in the [watch <path>] section of
// confText. nil when the section is absent.
func ParseSectionWhitelist(confText, resourcePath string) []string {
	var out []string
	inSection := false
	for line := range strings.SplitSeq(confText, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") {
			// Any header ends the section — a [libraries] block too.
			headerPath, ok := ParseSectionHeaderPath(trimmed)
			inSection = ok && headerPath == resourcePath
			continue
		}
		if !inSection || trimmed == "" || strings.HasPrefix(trimmed, "#") || IsLibraryDirective(trimmed) {
			continue
		}
		if trimmed == "need_encryption: true" || trimmed == "need_encryption: false" ||
			strings.HasPrefix(trimmed, "watch:") {
			continue
		}
		// A binary line's path may be double-quoted (spaces included, as the installer emits): take
		// it whole, like the daemon's parser.
		if strings.HasPrefix(trimmed, `"`) {
			if end := strings.IndexByte(trimmed[1:], '"'); end != -1 {
				out = append(out, trimmed[1:1+end])
				continue
			}
		}
		if fields := strings.Fields(trimmed); len(fields) > 0 {
			out = append(out, fields[0])
		}
	}
	return out
}

// LiveEmptyWhitelistRejected is the live refresh's fail-closed contract (RefreshSection).
// Non-encrypted resources are exempt: empty there is a legitimately uninstalled binary.
func LiveEmptyWhitelistRejected(encrypted, live bool, fresh, old int) bool {
	return live && encrypted && fresh == 0 && old > 0
}
