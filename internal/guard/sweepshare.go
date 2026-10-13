package guard

import (
	"sync"
	"time"
)

// sweepShareTTL bounds how old another guard's resolution may be when the periodic sweep reuses
// it. Sweeps fire on one shared pulse, so the guards whitelisting a binary resolve it within this
// window; reuse then looks like a sweep run up to this much earlier, well inside the 30s cadence.
const sweepShareTTL = 2 * time.Second

// sweepShared memoizes, per whitelisted path, the canonical path and inode key one confined open
// (resyncOne) resolved. It only lets the periodic sweep SKIP a binary whose resolution equals the
// key the guard already deploys: nothing is ever admitted from it. Any mismatch, a miss or an
// expired entry falls back to the guard's own confined open and replacement check.
var sweepShared = sweepMemo{ids: map[string]sweepID{}}

type sweepID struct {
	canonical string
	key       GuardInodeKey
	at        time.Time
}

type sweepMemo struct {
	mu  sync.Mutex
	ids map[string]sweepID
}

const sweepMemoMax = 8192 // pathological only; keeps the map bounded

func (m *sweepMemo) record(binPath, canonical string, key GuardInodeKey, now time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.ids) >= sweepMemoMax {
		m.ids = map[string]sweepID{}
	}
	m.ids[binPath] = sweepID{canonical: canonical, key: key, at: now}
}

func (m *sweepMemo) lookup(binPath string, now time.Time) (sweepID, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id, ok := m.ids[binPath]
	if !ok || now.Sub(id.at) > sweepShareTTL || now.Before(id.at) {
		return sweepID{}, false
	}
	return id, true
}

// unchangedPerShared reports that binPath still resolves, per a fresh shared resolution, to the
// key g deploys for it, so this sweep has nothing to re-admit.
func (g *Guard) unchangedPerShared(binPath string, now time.Time) bool {
	id, ok := sweepShared.lookup(binPath, now)
	if !ok {
		return false
	}
	g.mu.Lock()
	deployed, ok := g.deployed[id.canonical]
	g.mu.Unlock()
	return ok && deployed.Real() == id.key
}
