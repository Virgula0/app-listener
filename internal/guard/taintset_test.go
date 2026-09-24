package guard

import (
	"reflect"
	"testing"
)

func allowSet(exes ...GuardInodeKey) map[GuardInodeKey]uint8 {
	m := make(map[GuardInodeKey]uint8, len(exes))
	for _, e := range exes {
		m[e] = uint8(GUARD_ALLOW)
	}
	return m
}

func TestPlanTaintOwners(t *testing.T) {
	steam := GuardInodeKey{Dev: 1, Ino: 10}
	sshOnly := GuardInodeKey{Dev: 1, Ino: 11}
	writer := GuardInodeKey{Dev: 1, Ino: 12} // lib_binary: read-only trees only
	helper := GuardInodeKey{Dev: 1, Ino: 13} // one secret resource plus a read-only tree
	const (
		config, registry, ssh, libA, libB, dead = 1, 2, 3, 4, 5, 6
	)
	info := map[uint32]resInfo{
		config:   {path: "/s/config"},
		registry: {path: "/s/registry.vdf"},
		ssh:      {path: "/h/.ssh"},
		libA:     {path: "/s/linux64", readonly: true},
		libB:     {path: "/s/steamrt64", readonly: true},
	}
	allows := map[uint32]map[GuardInodeKey]uint8{
		config:   allowSet(steam, helper),
		registry: allowSet(steam),
		ssh:      allowSet(sshOnly),
		libA:     allowSet(steam, writer, helper),
		libB:     allowSet(writer),
		dead:     allowSet(sshOnly), // retired slot: must not count
	}
	got := planTaintOwners(allows, info)
	want := map[GuardInodeKey]taintOwner{
		steam:   {setKey: setKeyOf([]string{"/s/config", "/s/registry.vdf"})},
		sshOnly: {res: ssh},
		writer:  {res: libA}, // never taints: the exec hook skips a read-only owner
		helper:  {res: config},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("planTaintOwners:\n got %+v\nwant %+v", got, want)
	}
}

// Across a reload the old and new guard of a path both allow the exe: one set per path set, so
// the slot processes were tainted with stays valid.
func TestPlanTaintOwners_ReloadOverlapKeepsSetKey(t *testing.T) {
	steam := GuardInodeKey{Dev: 1, Ino: 10}
	info := map[uint32]resInfo{
		1: {path: "/s/config"}, 2: {path: "/s/registry.vdf"},
		7: {path: "/s/config"}, 8: {path: "/s/registry.vdf"},
	}
	allows := map[uint32]map[GuardInodeKey]uint8{1: allowSet(steam), 2: allowSet(steam), 7: allowSet(steam), 8: allowSet(steam)}
	got := planTaintOwners(allows, info)[steam]
	if want := setKeyOf([]string{"/s/config", "/s/registry.vdf"}); got.setKey != want {
		t.Fatalf("overlap set key = %q, want %q", got.setKey, want)
	}
}

func TestIntersectAllows(t *testing.T) {
	a, b, c := GuardInodeKey{Ino: 1}, GuardInodeKey{Ino: 2}, GuardInodeKey{Ino: 3}
	x := map[GuardInodeKey]uint8{a: uint8(GUARD_ALLOW), b: uint8(GUARD_ALLOW_ROOT), c: uint8(GUARD_ALLOW)}
	y := map[GuardInodeKey]uint8{a: uint8(GUARD_ALLOW), b: uint8(GUARD_ALLOW)}
	got := intersectAllows([]map[GuardInodeKey]uint8{x, y})
	want := map[GuardInodeKey]uint8{a: uint8(GUARD_ALLOW), b: uint8(GUARD_ALLOW_ROOT)}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("intersection = %v, want %v (root-only must stay root-only)", got, want)
	}
	if len(intersectAllows(nil)) != 0 {
		t.Fatal("a set with no live member must allow no one")
	}
}

func TestIsSubset(t *testing.T) {
	if !isSubset([]string{"/a"}, []string{"/a", "/b"}) || isSubset([]string{"/a", "/c"}, []string{"/a", "/b"}) {
		t.Fatal("isSubset")
	}
}

// Guards are added one at a time, each minting a larger set; only the sets an exe still maps to
// may carry rows, or guard_taint_members (4096) overflows past ~38 resources.
func TestPlanMemberRows_SupersededSetsHaveNoRows(t *testing.T) {
	const n = 48
	info := make(map[uint32]resInfo, n)
	var paths []string
	sets := make(map[string]uint32)
	for i := uint32(1); i <= n; i++ {
		info[i] = resInfo{path: string(rune('A'+i/26)) + string(rune('a'+i%26))}
		paths = append(paths, info[i].path)
		if i >= 2 {
			sets[setKeyOf(append([]string(nil), paths...))] = GuardMaxRes - i
		}
	}
	last := uint32(GuardMaxRes - n)
	rows := planMemberRows(sets, info, map[uint32]struct{}{last: {}})
	if len(rows) != n {
		t.Fatalf("got %d rows, want %d (one per member of the only used set)", len(rows), n)
	}
	for k := range rows {
		if k.Set != last {
			t.Fatalf("row for superseded set %d", k.Set)
		}
	}

	small := sets[setKeyOf(paths[:2])]
	rows = planMemberRows(sets, info, map[uint32]struct{}{last: {}, small: {}})
	if _, ok := rows[GuardTaintMemberKey{Set: last, Member: small}]; !ok {
		t.Fatal("a used set must list the used sets it covers")
	}
	if len(rows) != n+1+2 {
		t.Fatalf("got %d rows, want %d", len(rows), n+1+2)
	}
}
