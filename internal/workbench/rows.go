// Package workbench is the desktop GUI's engine layer, free of any toolkit: it owns the monitors
// and guards a GUI session runs, applies policy edits in a fail-closed order, and buffers events for
// the view to poll. internal/gui renders it.
package workbench

import (
	"strings"
	"sync"
	"time"
)

// Verdict is what happened to an access: only seen (monitor), or judged by a guard.
type Verdict int

const (
	Observed Verdict = iota
	Allowed
	Denied
)

func (v Verdict) String() string {
	switch v {
	case Allowed:
		return "allowed"
	case Denied:
		return "denied"
	default:
		return "observed"
	}
}

// Row is one event as the GUI shows it.
type Row struct {
	Seq    uint64
	Time   time.Time
	Source string // "monitor" or "guard"
	Op     string
	Comm   string
	PID    uint32
	UID    uint32
	Path   string
	Dest   string
	// Proto, Local, Remote and Size are network events' (Size: bytes for SEND/RECV, backlog for LISTEN).
	Proto   string
	Local   string
	Remote  string
	Size    uint32
	Verdict Verdict
	// Detail names a gate or a suspect mark (guard only).
	Detail string
	// Exe identifies the acting executable, when known (see Exe).
	Exe Exe
}

// matches reports whether r contains needle (lower-case) in any text column.
func (r *Row) matches(needle string) bool {
	for _, s := range []string{r.Op, r.Comm, r.Path, r.Dest, r.Proto, r.Local, r.Remote, r.Exe.Path, r.Detail, r.Source} {
		if strings.Contains(strings.ToLower(s), needle) {
			return true
		}
	}
	return false
}

// ring keeps the newest rows up to a cap; version moves on every change so a poller can skip
// unchanged snapshots.
type ring struct {
	mu      sync.Mutex
	rows    []Row
	start   int
	n       int
	seq     uint64
	version uint64
	counts  [3]uint64
}

const ringCap = 20000

func newRing() *ring { return &ring{rows: make([]Row, ringCap)} }

func (r *ring) add(row *Row) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seq++
	row.Seq = r.seq
	r.counts[row.Verdict]++
	r.rows[(r.start+r.n)%ringCap] = *row
	if r.n < ringCap {
		r.n++
	} else {
		r.start = (r.start + 1) % ringCap
	}
	r.version++
}

func (r *ring) clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.start, r.n = 0, 0
	r.counts = [3]uint64{}
	r.version++
}

// Filter narrows a snapshot.
type Filter struct {
	Text     string
	Verdicts []Verdict // empty = all
}

func (f Filter) admits(r *Row, needle string) bool {
	if len(f.Verdicts) > 0 {
		ok := false
		for _, v := range f.Verdicts {
			ok = ok || r.Verdict == v
		}
		if !ok {
			return false
		}
	}
	return needle == "" || r.matches(needle)
}

// snapshot returns the matching rows, newest first, and the ring's version.
func (r *ring) snapshot(f Filter) (rows []Row, version uint64) {
	needle := strings.ToLower(strings.TrimSpace(f.Text))
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]Row, 0, r.n)
	for i := r.n - 1; i >= 0; i-- {
		row := &r.rows[(r.start+i)%ringCap]
		if f.admits(row, needle) {
			out = append(out, *row)
		}
	}
	return out, r.version
}

// Counts is the number of rows ever added per verdict since the last clear.
type Counts struct{ Observed, Allowed, Denied uint64 }

func (r *ring) totals() (c Counts, version uint64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return Counts{r.counts[Observed], r.counts[Allowed], r.counts[Denied]}, r.version
}
