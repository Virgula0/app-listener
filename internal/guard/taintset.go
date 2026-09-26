package guard

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"

	cilium "github.com/cilium/ebpf"
	log "github.com/sirupsen/logrus"
)

// Taint sets (guard.bpf.c, GUARD_RES_GLOBAL): a binary several non-read-only resources whitelist
// taints its process with the slot of exactly that set, whose whitelist is the intersection of
// only those resources'. A set is keyed by member PATHS, not slot ids: across a reload the old and
// new guard of a path are both members (every veto holds during the overlap), and the slot stays
// valid for processes tainted before it, so taint needs no remapping. Set slots are never reused
// while the engine lives; one with no live member allows no one.

type resInfo struct {
	path     string
	readonly bool
}

// taintOwner is where an exe's taint lands: a resource slot, or the set named by setKey.
type taintOwner struct {
	res    uint32
	setKey string
}

func setKeyOf(paths []string) string { return strings.Join(paths, "\x00") }

func setPaths(key string) []string { return strings.Split(key, "\x00") }

// planTaintOwners decides each allowed exe's taint owner. Read-only resources never taint
// (is_readonly_mode), so they are left out of a set; an exe only they allow keeps one of them,
// which the exec hook then skips.
func planTaintOwners(allows map[uint32]map[GuardInodeKey]uint8, info map[uint32]resInfo) map[GuardInodeKey]taintOwner {
	members := make(map[GuardInodeKey][]uint32)
	for res, set := range allows {
		if _, live := info[res]; !live {
			continue
		}
		for exe := range set {
			members[exe] = append(members[exe], res)
		}
	}
	out := make(map[GuardInodeKey]taintOwner, len(members))
	for exe, rs := range members {
		slices.Sort(rs)
		var taint []uint32
		for _, r := range rs {
			if !info[r].readonly {
				taint = append(taint, r)
			}
		}
		switch len(taint) {
		case 0:
			out[exe] = taintOwner{res: rs[0]}
		case 1:
			out[exe] = taintOwner{res: taint[0]}
		default:
			paths := make([]string, 0, len(taint))
			for _, r := range taint {
				paths = append(paths, info[r].path)
			}
			slices.Sort(paths)
			out[exe] = taintOwner{setKey: setKeyOf(slices.Compact(paths))}
		}
	}
	return out
}

// intersectAllows is the set of exe inodes every map allows, the stricter action winning:
// GUARD_ALLOW_ROOT (uid 0 only) over GUARD_ALLOW. No maps: allows nothing.
func intersectAllows(sets []map[GuardInodeKey]uint8) map[GuardInodeKey]uint8 {
	inter := make(map[GuardInodeKey]uint8)
	for i, set := range sets {
		if i == 0 {
			for k, a := range set {
				inter[k] = a
			}
			continue
		}
		for k, a := range inter {
			other, ok := set[k]
			switch {
			case !ok:
				delete(inter, k)
			case other == uint8(GUARD_ALLOW_ROOT) || a == uint8(GUARD_ALLOW_ROOT):
				inter[k] = uint8(GUARD_ALLOW_ROOT)
			}
		}
	}
	return inter
}

// isSubset reports whether every path of a is in b.
func isSubset(a, b []string) bool {
	for _, p := range a {
		if !slices.Contains(b, p) {
			return false
		}
	}
	return true
}

func (e *engine) liveResInfoLocked() map[uint32]resInfo {
	info := make(map[uint32]resInfo)
	for id, g := range &e.slots {
		if g != nil {
			info[uint32(id)] = resInfo{path: g.path, readonly: g.mode == ModeReadOnly} //nolint:gosec // id < GuardMaxRes
		}
	}
	return info
}

// allocSetSlotLocked reserves a slot for a new taint set, from the top of the table so resource
// slots (allocated from 1) rarely meet it. 0 when full.
func (e *engine) allocSetSlotLocked(key string) (uint32, error) {
	for id := GuardMaxRes - 1; id > int(resGlobal); id-- {
		slot := uint32(id) //nolint:gosec // id < GuardMaxRes
		if e.slots[id] != nil || e.setAt[slot] != "" {
			continue
		}
		if err := e.objs.GuardResConfig.Put(slot, GuardResConfig{
			Active:   1,
			Mode:     uint64(ModeWhitelist),
			TaintSet: 1,
		}); err != nil {
			return 0, fmt.Errorf("activating taint-set slot %d: %w", slot, err)
		}
		e.setAt[slot] = key
		e.sets[key] = slot
		return slot, nil
	}
	return 0, nil
}

// syncTaintLocked rewrites guard_exe_union, the taint sets' whitelists and guard_taint_members
// from the live per-resource allow sets. Writes are diffed against what was last written.
func (e *engine) syncTaintLocked() error {
	if !e.started {
		return nil
	}
	if e.sets == nil {
		e.sets = make(map[string]uint32)
		e.setAt = make(map[uint32]string)
		e.setRows = make(map[uint32]map[GuardInodeKey]uint8)
	}
	info := e.liveResInfoLocked()
	union := make(map[GuardInodeKey]uint32)
	for exe, o := range planTaintOwners(e.allows, info) {
		if o.setKey == "" {
			union[exe] = o.res
			continue
		}
		slot, ok := e.sets[o.setKey]
		if !ok {
			var err error
			if slot, err = e.allocSetSlotLocked(o.setKey); err != nil {
				return err
			}
			if slot == 0 {
				log.Warnf("guard: no free slot for a taint set, falling back to the shared "+
					"intersection for %s (stricter: its own processes may not inspect each other)",
					strings.ReplaceAll(o.setKey, "\x00", ", "))
				slot = resGlobal
			}
		}
		union[exe] = slot
	}
	// Sets before the union: an exe must never point at a set whose whitelist isn't written yet.
	if err := e.syncSetsLocked(info, union); err != nil {
		return err
	}
	if err := syncU32Map(e.objs.GuardExeUnion, e.unionRows, union); err != nil {
		return fmt.Errorf("writing exe union: %w", err)
	}
	e.unionRows = union
	return nil
}

// syncSetsLocked recomputes every taint set ever allocated: its whitelist (the intersection of its
// live members'), and membership rows for the sets union still points at (planMemberRows).
func (e *engine) syncSetsLocked(info map[uint32]resInfo, union map[GuardInodeKey]uint32) error {
	for key, slot := range e.sets {
		var allowSets []map[GuardInodeKey]uint8
		for res, ri := range info {
			if slices.Contains(setPaths(key), ri.path) {
				allowSets = append(allowSets, e.allows[res])
			}
		}
		if err := e.syncSetRowsLocked(slot, intersectAllows(allowSets)); err != nil {
			return err
		}
	}
	used := make(map[uint32]struct{})
	for _, slot := range union {
		if _, isSet := e.setAt[slot]; isSet {
			used[slot] = struct{}{}
		}
	}
	members := planMemberRows(e.sets, info, used)
	if err := syncSetMap(e.objs.GuardTaintMembers, e.memberRows, members); err != nil {
		return fmt.Errorf("writing taint-set members: %w", err)
	}
	e.memberRows = members
	return nil
}

// planMemberRows lists each used set's members: its live resources and the used sets it covers.
// Superseded sets (every guard added one at a time mints a larger one) get no rows, or the table
// grows quadratically with the resource count; a process still tainted by one merges to
// GUARD_RES_GLOBAL, the stricter intersection.
func planMemberRows(sets map[string]uint32, info map[uint32]resInfo,
	used map[uint32]struct{}) map[GuardTaintMemberKey]struct{} {
	members := make(map[GuardTaintMemberKey]struct{})
	for key, slot := range sets {
		if _, ok := used[slot]; !ok {
			continue
		}
		paths := setPaths(key)
		for res, ri := range info {
			if slices.Contains(paths, ri.path) {
				members[GuardTaintMemberKey{Set: slot, Member: res}] = struct{}{}
			}
		}
		for otherKey, other := range sets {
			if _, ok := used[other]; ok && other != slot && isSubset(setPaths(otherKey), paths) {
				members[GuardTaintMemberKey{Set: slot, Member: other}] = struct{}{}
			}
		}
	}
	return members
}

func (e *engine) syncSetRowsLocked(slot uint32, want map[GuardInodeKey]uint8) error {
	have := e.setRows[slot]
	for exe := range have {
		if _, keep := want[exe]; !keep {
			err := e.objs.GuardExeActions.Delete(GuardResInodeKey{ResId: slot, Ino: exe})
			if err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
				return fmt.Errorf("taint-set %d whitelist: %w", slot, err)
			}
		}
	}
	for exe, a := range want {
		if cur, ok := have[exe]; ok && cur == a {
			continue
		}
		if err := e.objs.GuardExeActions.Put(GuardResInodeKey{ResId: slot, Ino: exe}, a); err != nil {
			return fmt.Errorf("taint-set %d whitelist: %w", slot, err)
		}
	}
	e.setRows[slot] = want
	return nil
}

// setOwnerLocked is the guard that reports a taint set's denials: its live member with the lowest
// path, so they read as coming from one of the resources involved.
func (e *engine) setOwnerLocked(slot uint32) *Guard {
	paths := setPaths(e.setAt[slot])
	sort.Strings(paths)
	for _, p := range paths {
		for _, g := range &e.slots {
			if g != nil && g.path == p {
				return g
			}
		}
	}
	return nil
}

// syncU32Map writes want over have: puts first, then deletes, so no key passes through absent.
func syncU32Map(m *cilium.Map, have, want map[GuardInodeKey]uint32) error {
	for k, v := range want {
		if cur, ok := have[k]; ok && cur == v {
			continue
		}
		if err := m.Put(k, v); err != nil {
			return err
		}
	}
	for k := range have {
		if _, keep := want[k]; !keep {
			if err := m.Delete(k); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
				return err
			}
		}
	}
	return nil
}

func syncSetMap(m *cilium.Map, have, want map[GuardTaintMemberKey]struct{}) error {
	for k := range want {
		if _, ok := have[k]; !ok {
			if err := m.Put(k, uint8(1)); err != nil {
				return err
			}
		}
	}
	for k := range have {
		if _, keep := want[k]; !keep {
			if err := m.Delete(k); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
				return err
			}
		}
	}
	return nil
}
