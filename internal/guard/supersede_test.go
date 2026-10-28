package guard

import (
	"testing"
	"time"
)

func TestDroppableFreed_WaitsOnePruneAndHold(t *testing.T) {
	k := GuardInodeKey{Dev: 1, Ino: 7}
	t0 := time.Unix(1000, 0)

	drop, seen := droppableFreed([]GuardInodeKey{k}, nil, t0, false)
	if drop[k] {
		t.Fatal("a key first seen freed must wait for the next prune")
	}
	drop, seen = droppableFreed([]GuardInodeKey{k}, seen, t0.Add(supersedePruneEvery), true)
	if drop[k] {
		t.Fatal("a held key must stay")
	}
	if seen[k] != t0 {
		t.Fatalf("the first sighting must be kept across prunes, got %v", seen[k])
	}
	if drop, _ = droppableFreed([]GuardInodeKey{k}, seen, t0.Add(2*supersedePruneEvery), false); !drop[k] {
		t.Fatal("once released, a key seen by an earlier prune must go")
	}
	if drop, _ = droppableFreed([]GuardInodeKey{k}, seen, t0.Add(supersedeHoldMax), true); !drop[k] {
		t.Fatal("a hold must not outlast supersedeHoldMax")
	}
}

func TestDroppableFreed_ForgetsKeysNoLongerFreed(t *testing.T) {
	gone := GuardInodeKey{Dev: 1, Ino: 9}
	_, seen := droppableFreed(nil, map[GuardInodeKey]time.Time{gone: time.Unix(1, 0)}, time.Unix(2, 0), false)
	if _, ok := seen[gone]; ok {
		t.Fatal("a key no longer marked freed must leave freedAt")
	}
}
