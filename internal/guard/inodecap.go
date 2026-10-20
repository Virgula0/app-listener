package guard

import (
	"errors"
	"fmt"

	cilium "github.com/cilium/ebpf"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

// inodeWarnPercent: occupancy past it is warned. A full guard_inodes drops rows, leaving those
// inodes to the ancestor walk alone (degradeInodesFull), so the operator gets room to act first.
const inodeWarnPercent = 75

// ConfigureInodeCapacity sets guard_inodes' max_entries (0 = compiled default) for the engine's next
// load, failing outside daemonconfig's bounds. It reports false when the engine is already loaded
// with another capacity: a live map can't be resized, so the change waits for a restart.
func ConfigureInodeCapacity(n uint32) (bool, error) {
	if n != 0 && (n < daemonconfig.MaxInodesFloor || n > daemonconfig.MaxInodesCeiling) {
		return false, fmt.Errorf("guard_inodes capacity %d is outside %d..%d", n,
			daemonconfig.MaxInodesFloor, daemonconfig.MaxInodesCeiling)
	}
	e := &sharedEngine
	e.mu.Lock()
	defer e.mu.Unlock()
	e.inodeCap = n
	if !e.started || e.objs.GuardInodes == nil {
		return true, nil
	}
	want := n
	if want == 0 {
		want = e.inodeCapDefault
	}
	return e.objs.GuardInodes.MaxEntries() == want, nil
}

// applyInodeCapacity sizes guard_inodes in spec before load and returns the object's default.
// Never below the floor (fail closed: refuse to load rather than run a map too small for the trees).
func applyInodeCapacity(spec *cilium.CollectionSpec, n uint32) (uint32, error) {
	m, ok := spec.Maps[GuardMapGuardInodes]
	if !ok {
		return 0, errors.New("guard object has no guard_inodes map")
	}
	def := m.MaxEntries
	if n == 0 {
		return def, nil
	}
	if n < daemonconfig.MaxInodesFloor {
		return 0, fmt.Errorf("guard_inodes capacity %d is below the floor %d", n, daemonconfig.MaxInodesFloor)
	}
	m.MaxEntries = n
	return def, nil
}

// InodeOccupancy counts guard_inodes rows against its capacity; ok is false while no engine runs.
// The walk runs on a dup of the map fd, outside e.mu: claimInode must not wait on it.
func InodeOccupancy() (rows, capacity int, ok bool, err error) {
	e := &sharedEngine
	e.mu.Lock()
	var m *cilium.Map
	if e.started && e.objs.GuardInodes != nil {
		m, err = e.objs.GuardInodes.Clone()
	}
	e.mu.Unlock()
	if m == nil || err != nil {
		return 0, 0, false, err
	}
	defer m.Close()
	rows, err = countInodeRows(func(c *cilium.MapBatchCursor, keys []GuardInodeKey, vals []uint32) (int, error) {
		return m.BatchLookup(c, keys, vals, nil)
	})
	if err != nil {
		rows = 0
		var key GuardInodeKey
		var val uint32
		it := m.Iterate()
		for it.Next(&key, &val) {
			rows++
		}
		err = it.Err()
	}
	return rows, int(m.MaxEntries()), true, err
}

func countInodeRows(lookup inodeBatchFn) (int, error) {
	rows := 0
	err := walkInodeBatches(lookup, func(GuardInodeKey, uint32) { rows++ })
	return rows, err
}

// WarnInodeOccupancy warns when guard_inodes is past inodeWarnPercent, once per crossing.
func WarnInodeOccupancy() {
	rows, capacity, ok, err := InodeOccupancy()
	if err != nil {
		log.Warnf("guard: counting guard_inodes rows: %v", err)
		return
	}
	if !ok || capacity == 0 {
		return
	}
	high := rows*100 >= capacity*inodeWarnPercent
	e := &sharedEngine
	e.mu.Lock()
	warned := e.inodeWarned
	e.inodeWarned = high
	e.mu.Unlock()
	switch {
	case high && !warned:
		log.Warnf("guard: guard_inodes holds %d of %d rows (%d%%): once full, new inodes are guarded by "+
			"the ancestor walk alone; raise max_inodes in daemon.conf's [settings] and restart",
			rows, capacity, rows*100/capacity)
	case !high:
		log.Debugf("guard: guard_inodes holds %d of %d rows", rows, capacity)
	}
}
