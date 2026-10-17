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
		if err := g.removeBinaries(removed); err != nil {
			return err
		}
		return g.addBinaries(added)
	}
	if err := g.addBinaries(added); err != nil {
		return err
	}
	return g.removeBinaries(removed)
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
// included (keyPaths), then forgets them.
func (g *Guard) removeBinaries(removed []BinaryEntry) error {
	if len(removed) == 0 {
		return nil
	}
	g.mu.Lock()
	var keys []GuardInodeKey
	for _, b := range removed {
		if k, ok := g.deployed[canonicalBinaryPath(b.Path)]; ok {
			keys = append(keys, k)
		}
		for k, p := range g.keyPaths {
			if p == b.Path && !slices.Contains(keys, k) {
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
	for _, b := range removed {
		canonical := canonicalBinaryPath(b.Path)
		delete(g.deployed, canonical)
		delete(g.binaryVerifyStates, canonical)
		delete(g.canonicalPaths, b.Path)
		g.holdLocked(b.Path, nil)
		g.binaries = slices.DeleteFunc(g.binaries, func(e BinaryEntry) bool { return e == b })
	}
	for _, k := range keys {
		delete(g.keyPaths, k)
	}
	return nil
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
