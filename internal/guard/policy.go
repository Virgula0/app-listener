package guard

import (
	"errors"
	"fmt"
	"slices"

	cilium "github.com/cilium/ebpf"
)

// SetBinaries replaces a running guard's whitelist or blacklist in place (the desktop GUI's live
// policy edit). No tree walk is needed, so it works while the guard denies its own process.
// A transition never admits a binary outside old ∪ new: whitelist removals land before additions,
// blacklist additions (blocks) before removals. An entry whose hash changed since it was added
// counts as removed and re-added. Only for plain guards: not pinned, nothing deferred, no self
// binary (the daemon's guards rebuild through Reload instead).
func (g *Guard) SetBinaries(entries []BinaryEntry) error {
	if g.mode != ModeWhitelist && g.mode != ModeBlacklist {
		return fmt.Errorf("guard %s: SetBinaries needs whitelist or blacklist mode", g.path)
	}
	g.policyMu.Lock()
	defer g.policyMu.Unlock()

	g.mu.Lock()
	if g.stopped {
		g.mu.Unlock()
		return fmt.Errorf("guard %s: stopped", g.path)
	}
	if g.pinPrefix != "" || len(g.deferred) > 0 || g.selfKeySet {
		g.mu.Unlock()
		return fmt.Errorf("guard %s: SetBinaries is for plain guards (no pinning, deferred or self binary)", g.path)
	}
	added, removed := diffBinaries(g.binaries, entries)
	g.mu.Unlock()

	if g.mode == ModeWhitelist {
		kept := slices.DeleteFunc(slices.Clone(entries), func(e BinaryEntry) bool { return slices.Contains(added, e) })
		if err := g.removeBinaries(removed, kept); err != nil {
			return err
		}
		return g.addBinaries(added)
	}
	if err := g.addBinaries(added); err != nil {
		return err
	}
	// After the adds: a re-added path or a second path to one inode owns the row just written.
	return g.removeBinaries(removed, entries)
}

// diffBinaries splits next against cur by path; a same-path entry with another hash is in both.
func diffBinaries(cur, next []BinaryEntry) (added, removed []BinaryEntry) {
	for _, n := range next {
		if !slices.Contains(cur, n) {
			added = append(added, n)
		}
	}
	for _, c := range cur {
		if !slices.Contains(next, c) {
			removed = append(removed, c)
		}
	}
	return added, removed
}

func (g *Guard) addBinaries(added []BinaryEntry) error {
	if len(added) == 0 {
		return nil
	}
	if err := checkBtrfsKeys(g.path, added); err != nil {
		return err
	}
	// Rows first, then bookkeeping: a failure part-way leaves extra rows only for entries the
	// caller asked to add (whitelist) or block (blacklist), never a removed one.
	if err := g.addBinaryActions(added, nil); err != nil {
		return err
	}
	g.mu.Lock()
	for _, b := range added {
		canonical := canonicalBinaryPath(b.Path)
		g.binaries = append(g.binaries, b)
		g.canonicalPaths[b.Path] = canonical
		if key, ok := g.deployed[canonical]; ok && g.verifyStop != nil {
			g.binaryVerifyStates[canonical] = newVerifyState(key, b.Hash)
		}
	}
	// A whitelist started empty has no verifier yet (Start); this one pins every entry.
	startVerifier := g.mode == ModeWhitelist && g.verifyStop == nil
	g.mu.Unlock()
	if startVerifier {
		g.startBinaryHashVerifier()
	}
	return nil
}

// removeBinaries drops every row this guard wrote for the entries' paths, superseded keys
// included (keyPaths), then forgets them. A key or canonical path one of keep still owns (same
// path re-added, symlink or hardlink to one inode) keeps its row and bookkeeping: deleting it
// would silently unblock (blacklist) or deny (whitelist) the kept entry.
func (g *Guard) removeBinaries(removed, keep []BinaryEntry) error {
	if len(removed) == 0 {
		return nil
	}
	g.mu.Lock()
	kept := g.keptLocked(keep)
	var keys []GuardInodeKey
	for _, b := range removed {
		for _, k := range g.entryKeysLocked(b.Path) {
			if _, owned := kept.owners[k]; !owned && !slices.Contains(keys, k) {
				keys = append(keys, k)
			}
		}
	}
	g.mu.Unlock()

	for _, k := range keys {
		if err := g.deleteExeRows(k); err != nil {
			return fmt.Errorf("guard %s: removing binary rows: %w", g.path, err)
		}
	}

	g.mu.Lock()
	defer g.mu.Unlock()
	g.forgetLocked(removed, keys, kept)
	return nil
}

type keptEntries struct {
	owners    map[GuardInodeKey]string // key -> a kept entry's path
	canonical map[string]bool
	paths     map[string]bool
}

func (g *Guard) keptLocked(keep []BinaryEntry) keptEntries {
	kept := keptEntries{
		owners:    make(map[GuardInodeKey]string),
		canonical: make(map[string]bool, len(keep)),
		paths:     make(map[string]bool, len(keep)),
	}
	for _, e := range keep {
		for _, k := range g.entryKeysLocked(e.Path) {
			kept.owners[k] = e.Path
		}
		kept.canonical[canonicalBinaryPath(e.Path)] = true
		kept.paths[e.Path] = true
	}
	return kept
}

// forgetLocked drops the bookkeeping of removed entries and their deleted keys, keeping what kept
// still uses; a kept key last recorded under a removed path moves to its owner's.
func (g *Guard) forgetLocked(removed []BinaryEntry, keys []GuardInodeKey, kept keptEntries) {
	for _, b := range removed {
		g.binaries = slices.DeleteFunc(g.binaries, func(e BinaryEntry) bool { return e == b })
		if kept.paths[b.Path] {
			continue
		}
		delete(g.canonicalPaths, b.Path)
		g.holdLocked(b.Path, nil)
		if canonical := canonicalBinaryPath(b.Path); !kept.canonical[canonical] {
			delete(g.deployed, canonical)
			delete(g.binaryVerifyStates, canonical)
		}
	}
	for _, k := range keys {
		delete(g.keyPaths, k)
	}
	for k, p := range kept.owners {
		if cur, ok := g.keyPaths[k]; ok && !kept.paths[cur] {
			g.keyPaths[k] = p
		}
	}
}

// entryKeysLocked is every key this guard wrote for path: the deployed one and superseded ones.
func (g *Guard) entryKeysLocked(path string) []GuardInodeKey {
	var keys []GuardInodeKey
	if k, ok := g.deployed[canonicalBinaryPath(path)]; ok {
		keys = append(keys, k)
	}
	for k, p := range g.keyPaths {
		if p == path && !slices.Contains(keys, k) {
			keys = append(keys, k)
		}
	}
	return keys
}

// deleteExeRows removes k's action and mask rows for this resource (and the unattested-applet row
// a blacklisted multicall key also wrote, putExeAction), then withdraws it from the shared views.
func (g *Guard) deleteExeRows(k GuardInodeKey) error {
	rows := []GuardInodeKey{k}
	if k.tagged() && g.mode == ModeBlacklist {
		rows = append(rows, k.withTag(mcTagNone))
	}
	for _, r := range rows {
		if err := g.objs().GuardExeActions.Delete(g.resKey(r)); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
			return err
		}
	}
	if err := g.objs().GuardExeEvents.Delete(g.resKey(k)); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
		return err
	}
	return sharedEngine.noteDeny(g.resID, k)
}
