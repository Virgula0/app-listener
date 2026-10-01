package guard

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	cilium "github.com/cilium/ebpf"
	log "github.com/sirupsen/logrus"
)

// Superseded keys (exe_supersede.h): a whitelisted or trusted inode that lost its last link keeps
// admitting processes exec'd before, and refuses every later exec, a reused inode number's
// included. The kernel marks the key; userspace drops its rows once the inode is gone, and carries
// the rest across a reload.

var (
	supersedeMapsOnce sync.Once
	supersedeMapSet   map[string]*cilium.Map
	errSupersedeMaps  error
)

// supersedeMaps returns the maps the guard and trust objects share, created once per process: a
// stamp must outlive any one engine, and both objects must judge a key by the same stamps.
func supersedeMaps() (map[string]*cilium.Map, error) {
	supersedeMapsOnce.Do(func() {
		spec, err := LoadGuard()
		if err != nil {
			errSupersedeMaps = fmt.Errorf("reading embedded guard objects: %w", err)
			return
		}
		out := make(map[string]*cilium.Map, 3)
		for _, name := range []string{GuardMapExeSuperseded, GuardMapExeStamps, GuardMapExeSeq} {
			m, err := cilium.NewMap(spec.Maps[name])
			if err != nil {
				errSupersedeMaps = fmt.Errorf("creating shared map %s: %w", name, err)
				return
			}
			out[name] = m
		}
		supersedeMapSet = out
	})
	return supersedeMapSet, errSupersedeMaps
}

// supersede serializes prune, lift and carry: each reads a key's mark and acts on it.
var supersede struct {
	mu    sync.Mutex
	trust *TrustGuard
}

func setSupersedeTrust(t *TrustGuard) {
	supersede.mu.Lock()
	supersede.trust = t
	supersede.mu.Unlock()
}

func supersededMark(k GuardInodeKey) (GuardExeSupersede, bool) {
	var v GuardExeSupersede
	maps, err := supersedeMaps()
	if err != nil {
		return v, false
	}
	return v, maps[GuardMapExeSuperseded].Lookup(k, &v) == nil
}

// exeGone reports a superseded key whose inode was freed.
func exeGone(k GuardInodeKey) bool {
	v, ok := supersededMark(k)
	return ok && v.Freed != 0
}

// PruneSuperseded drops every superseded key whose inode is gone, rows first: until the mark goes,
// a reused number is refused whatever rows it finds. A mark no resource or trusted file names any
// more goes too.
func PruneSuperseded() {
	maps, err := supersedeMaps()
	if err != nil {
		return
	}
	supersede.mu.Lock()
	defer supersede.mu.Unlock()

	marks := maps[GuardMapExeSuperseded]
	var gone, unused []GuardInodeKey
	var k GuardInodeKey
	var v GuardExeSupersede
	it := marks.Iterate()
	for it.Next(&k, &v) {
		switch {
		case v.Freed != 0:
			gone = append(gone, k)
		case !sharedEngine.admitsExe(k) && (supersede.trust == nil || !supersede.trust.trusts(k)):
			unused = append(unused, k)
		}
	}
	if err := it.Err(); err != nil {
		log.Warnf("guard: reading superseded binaries: %v", err)
		return
	}
	for _, k := range gone {
		dropSupersededLocked(k)
	}
	for _, k := range unused {
		if err := marks.Delete(k); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
			log.Warnf("guard: forgetting superseded inode %d: %v", k.Ino, err)
		}
	}
}

// dropSupersededLocked removes every row naming k, then its mark.
func dropSupersededLocked(k GuardInodeKey) {
	if err := sharedEngine.forgetExe(k); err != nil {
		log.Warnf("guard: dropping rows of freed binary inode %d: %v", k.Ino, err)
		return
	}
	if t := supersede.trust; t != nil {
		if err := t.forgetExe(k); err != nil {
			log.Warnf("trust guard: dropping rows of freed inode %d: %v", k.Ino, err)
			return
		}
	}
	maps, _ := supersedeMaps()
	if err := maps[GuardMapExeSuperseded].Delete(k); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
		log.Warnf("guard: forgetting superseded inode %d: %v", k.Ino, err)
		return
	}
	log.Infof("guard: forgot freed binary inode %d (dev %d)", k.Ino, k.Dev)
}

// liftSuperseded runs before k is admitted from a vetted path. A mark on a freed inode belongs to a
// previous file whose number k reuses: its rows and mark go, so they neither refuse nor admit the
// new file. A mark on a live inode stays: the file just stat'ed may be the one being unlinked.
func liftSuperseded(k GuardInodeKey) {
	supersede.mu.Lock()
	defer supersede.mu.Unlock()
	if exeGone(k) {
		dropSupersededLocked(k)
	}
}

// SupersededBinary is a superseded key a guard admits, carried to its replacement across a reload.
type SupersededBinary struct {
	// Path is the whitelist line the key was admitted for.
	Path    string
	Key     GuardInodeKey
	Action  uint8
	Mask    uint32
	HasMask bool
}

// SnapshotSuperseded lists this guard's superseded keys whose inode still exists: a process
// exec'd before the supersede may still run it.
func (g *Guard) SnapshotSuperseded() ([]SupersededBinary, error) {
	g.mu.Lock()
	paths := make(map[GuardInodeKey]string, len(g.keyPaths))
	for k, p := range g.keyPaths {
		paths[k] = p
	}
	g.mu.Unlock()

	var out []SupersededBinary
	for k, p := range paths {
		if mark, ok := supersededMark(k); !ok || mark.Freed != 0 {
			continue
		}
		var action uint8
		if err := g.objs().GuardExeActions.Lookup(g.resKey(k), &action); err != nil {
			if errors.Is(err, cilium.ErrKeyNotExist) {
				continue
			}
			return nil, fmt.Errorf("reading superseded binary %s for %s: %w", p, g.path, err)
		}
		b := SupersededBinary{Path: p, Key: k, Action: action}
		b.HasMask = g.objs().GuardExeEvents.Lookup(g.resKey(k), &b.Mask) == nil
		out = append(out, b)
	}
	return out, nil
}

// RestoreSuperseded admits bins in this (fresh) guard, so an old-version process keeps its access
// across a reload. A key is carried only while its inode exists, and only if this guard still
// whitelists its line or the line's file is gone (a vanished version): a line the operator dropped
// while its binary is still installed takes its old versions with it.
func (g *Guard) RestoreSuperseded(bins []SupersededBinary) error {
	supersede.mu.Lock()
	defer supersede.mu.Unlock()
	for _, b := range bins {
		if !g.carriesLine(b.Path) {
			continue
		}
		if mark, ok := supersededMark(b.Key); !ok || mark.Freed != 0 {
			continue
		}
		if err := g.putExeAction(b.Key, b.Action); err != nil {
			return fmt.Errorf("carrying superseded %s into %s: %w", b.Path, g.path, err)
		}
		if b.HasMask {
			if err := g.objs().GuardExeEvents.Put(g.resKey(b.Key), b.Mask); err != nil {
				return fmt.Errorf("carrying the event mask of superseded %s into %s: %w", b.Path, g.path, err)
			}
		}
		g.mu.Lock()
		g.keyPaths[b.Key] = b.Path
		g.mu.Unlock()
		log.Infof("guard %s: kept superseded %s (inode %d) for processes started before its update", g.path,
			b.Path, b.Key.Ino)
	}
	return nil
}

func (g *Guard) carriesLine(path string) bool {
	g.mu.Lock()
	for _, b := range g.binaries {
		if b.Path == path {
			g.mu.Unlock()
			return true
		}
	}
	for _, d := range g.deferred {
		if d.rule.Path == path {
			g.mu.Unlock()
			return true
		}
	}
	g.mu.Unlock()
	_, err := os.Lstat(path)
	return errors.Is(err, os.ErrNotExist)
}

// forgetExe drops what this guard remembers about k's rows, so a file reusing the number at a
// whitelisted path is judged as a new binary.
func (g *Guard) forgetExe(k GuardInodeKey) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for p, key := range g.deployed {
		if key == k {
			delete(g.deployed, p)
		}
	}
	for p, st := range g.binaryVerifyStates {
		if st.key == k {
			delete(g.binaryVerifyStates, p)
		}
	}
	delete(g.keyPaths, k)
}

const supersedePruneEvery = 30 * time.Second

// supersedePruneLoop runs PruneSuperseded while the engine lives.
func supersedePruneLoop(done <-chan struct{}) {
	t := time.NewTicker(supersedePruneEvery)
	defer t.Stop()
	var lostSeen uint64
	for {
		select {
		case <-done:
			return
		case <-t.C:
			PruneSuperseded()
			lostSeen = warnLostStamps(lostSeen)
		}
	}
}

// warnLostStamps reports exec stamps the kernel could not record (exe_stamps full): from then on a
// process with no stamp is refused by every superseded key, old-version apps started before the
// daemon included.
func warnLostStamps(seen uint64) uint64 {
	maps, err := supersedeMaps()
	if err != nil {
		return seen
	}
	var lost uint64
	if maps[GuardMapExeSeq].Lookup(uint32(1), &lost) != nil || lost <= seen {
		return seen
	}
	log.Warnf("guard: degraded: %d exec stamp(s) lost (process table larger than the stamp map) — "+
		"processes started before the daemon lose access through binaries replaced since", lost)
	return lost
}
