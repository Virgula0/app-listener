package guard

import (
	"errors"
	"testing"

	cilium "github.com/cilium/ebpf"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

func TestApplyInodeCapacity(t *testing.T) {
	spec, _, err := guardSpec()
	if err != nil {
		t.Fatal(err)
	}
	def := spec.Maps[GuardMapGuardInodes].MaxEntries
	if got, err := applyInodeCapacity(spec, 0); err != nil || got != def || spec.Maps[GuardMapGuardInodes].MaxEntries != def {
		t.Fatalf("0 must keep the object's default %d: got %d, %v", def, got, err)
	}
	if _, err := applyInodeCapacity(spec, 65536); err != nil || spec.Maps[GuardMapGuardInodes].MaxEntries != 65536 {
		t.Fatalf("65536 not applied: %v", err)
	}
	spec.Maps[GuardMapGuardInodes].MaxEntries = def
	if _, err := applyInodeCapacity(spec, daemonconfig.MaxInodesFloor-1); err == nil ||
		spec.Maps[GuardMapGuardInodes].MaxEntries != def {
		t.Fatal("a capacity below the floor must refuse to load, leaving the spec alone")
	}
	delete(spec.Maps, GuardMapGuardInodes)
	if _, err := applyInodeCapacity(spec, 65536); err == nil {
		t.Fatal("a spec without guard_inodes must fail")
	}
}

func TestConfigureInodeCapacityBounds(t *testing.T) {
	prev := sharedEngine.inodeCap
	t.Cleanup(func() { sharedEngine.inodeCap = prev })
	for _, n := range []uint32{1, daemonconfig.MaxInodesFloor - 1, daemonconfig.MaxInodesCeiling + 1} {
		if _, err := ConfigureInodeCapacity(n); err == nil {
			t.Fatalf("%d must be refused", n)
		}
	}
	if applied, err := ConfigureInodeCapacity(65536); err != nil || !applied || sharedEngine.inodeCap != 65536 {
		t.Fatalf("65536 before the engine loads: applied %v, %v", applied, err)
	}
}

func TestCountInodeRows(t *testing.T) {
	batches := [][]uint64{{1, 2, 3}, {4, 5}}
	vals := [][]uint32{{7, 8, 7}, {7, 9}}
	if n, err := countInodeRows(fakeBatches(batches, vals, cilium.ErrKeyNotExist)); err != nil || n != 5 {
		t.Fatalf("rows %d err %v: want every resource's rows, last batch included", n, err)
	}
	boom := errors.New("ENOSPC")
	if _, err := countInodeRows(fakeBatches(batches, vals, boom)); !errors.Is(err, boom) {
		t.Fatalf("err %v, want the batch error so the caller falls back", err)
	}
	if _, err := countInodeRows(fakeBatches(batches, vals, nil)); err == nil {
		t.Fatal("a stalled walk must error")
	}
}
