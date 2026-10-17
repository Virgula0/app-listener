package workbench

import (
	"slices"
	"sync"
)

// activity buffers a workspace's events and aggregates them per executable.
type activity struct {
	rows   *ring
	ownPID uint32

	mu    sync.Mutex
	procs map[string]*Process
}

func newActivity(ownPID uint32) *activity {
	return &activity{rows: newRing(), ownPID: ownPID, procs: map[string]*Process{}}
}

func (a *activity) ingest(r *Row) {
	// The GUI's own accesses (its monitor stats inside a whitelisted tree) are noise, not activity.
	if r.PID == a.ownPID {
		return
	}
	a.rows.add(r)
	key := processKey(r.Exe, r.Comm)
	a.mu.Lock()
	p := a.procs[key]
	if p == nil {
		p = &Process{Exe: r.Exe, Comm: r.Comm}
		a.procs[key] = p
	}
	if p.Exe.Path == "" && r.Exe.Path != "" {
		p.Exe = r.Exe
	}
	p.Events++
	if r.Verdict == Denied {
		p.Denied++
	}
	p.LastSeen = r.Time
	a.mu.Unlock()
}

// Rows returns the buffered events matching f, newest first, and a version that changes with them.
func (a *activity) Rows(f Filter) (rows []Row, version uint64) { return a.rows.snapshot(f) }

// Counts returns the per-verdict totals and the buffer version.
func (a *activity) Counts() (c Counts, version uint64) { return a.rows.totals() }

// Clear empties the event buffer and the process list.
func (a *activity) Clear() {
	a.rows.clear()
	a.mu.Lock()
	a.procs = map[string]*Process{}
	a.mu.Unlock()
}

// Processes returns the executables seen, most recent first.
func (a *activity) Processes() []Process {
	a.mu.Lock()
	out := make([]Process, 0, len(a.procs))
	for _, p := range a.procs {
		out = append(out, *p)
	}
	a.mu.Unlock()
	slices.SortFunc(out, func(x, y Process) int { return y.LastSeen.Compare(x.LastSeen) })
	return out
}
