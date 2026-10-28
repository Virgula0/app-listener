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
// stamp must outlive any one engine, and both objects must judge a key by the same stamps (and key
// btrfs inodes by the same layout, and see the same code-suspect marks).
func supersedeMaps() (map[string]*cilium.Map, error) {
	supersedeMapsOnce.Do(func() {
		spec, err := LoadGuard()
		if err != nil {
			errSupersedeMaps = fmt.Errorf("reading embedded guard objects: %w", err)
			return
		}
		names := []string{GuardMapExeSuperseded, GuardMapExeStamps, GuardMapExeSeq, GuardMapBtrfsLayout,
			GuardMapTrustCodeSuspect}
		out := make(map[string]*cilium.Map, len(names))
		for _, name := range names {
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
	// hold: SetSupersedeHold. freedAt: when a prune first saw each freed key.
	hold    func() bool
	freedAt map[GuardInodeKey]time.Time
}

// SetSupersedeHold makes the prune keep freed keys while hold reports true (a catalog refresh
// pending or running), for at most supersedeHoldMax: a version directory an updater wrote is proven
// by its creator's guard_bin_updaters row (TrustGuard.CreatedByUpdater), and an app that deletes its
// old version frees that creator before the refresh judges the new one. nil clears it.
func SetSupersedeHold(hold func() bool) {
	supersede.mu.Lock()
	supersede.hold = hold
	supersede.mu.Unlock()
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
	return v, maps[GuardMapExeSuperseded].Lookup(k.Real(), &v) == nil
}

// exeGone reports a superseded key whose inode was freed.
func exeGone(k GuardInodeKey) bool {
	v, ok := supersededMark(k)
	return ok && v.Freed != 0
}

// PruneSuperseded drops every superseded key whose inode is gone, rows first: until the mark goes,
// a reused number is refused whatever rows it finds. A mark no resource or trusted file names any
// more goes too.
//
// Every guard re-syncs before a freed key's rows go, and a key freed after that re-sync waits for
// the next prune: a freed updater's guard_bin_updaters row is what proves the replacements it
// created (updaterCreated), so a guard judging one after the drop would refuse it. A self-updating
// app (Steam) frees its creator as soon as it restarts. A key also waits until an earlier prune saw
// it freed, so the catalog watch has seen the delete, and while SetSupersedeHold holds. Keeping
// rows longer admits nothing: the mark, which goes last, refuses a reused number.
func PruneSuperseded() {
	maps, err := supersedeMaps()
	if err != nil {
		return
	}
	marks := maps[GuardMapExeSuperseded]
	freed := freedSuperseded(marks, time.Now())
	if len(freed) == 0 {
		pruneSuperseded(marks, nil)
		return
	}
	log.Infof("guard: re-syncing every guard before forgetting %d freed binary inode(s)", len(freed))
	sharedEngine.resyncGuards()
	pruneSuperseded(marks, freed)
}

// freedSuperseded lists the freed keys the prune may drop at now (see PruneSuperseded).
func freedSuperseded(marks *cilium.Map, now time.Time) map[GuardInodeKey]bool {
	supersede.mu.Lock()
	hold := supersede.hold
	supersede.mu.Unlock()
	held := hold != nil && hold()

	supersede.mu.Lock()
	defer supersede.mu.Unlock()
	var keys []GuardInodeKey
	var k GuardInodeKey
	var v GuardExeSupersede
	it := marks.Iterate()
	for it.Next(&k, &v) {
		if v.Freed != 0 {
			keys = append(keys, k)
		}
	}
	if err := it.Err(); err != nil {
		log.Warnf("guard: reading superseded binaries: %v", err)
		return nil
	}
	freed, seen := droppableFreed(keys, supersede.freedAt, now, held)
	supersede.freedAt = seen
	return freed
}

// droppableFreed picks the freed keys a prune may drop: seen freed by an earlier prune (freedAt) and
// not held, or held past supersedeHoldMax. seen is freedAt for the next prune.
func droppableFreed(keys []GuardInodeKey, freedAt map[GuardInodeKey]time.Time, now time.Time,
	held bool) (drop map[GuardInodeKey]bool, seen map[GuardInodeKey]time.Time) {
	drop = map[GuardInodeKey]bool{}
	seen = make(map[GuardInodeKey]time.Time, len(keys))
	for _, k := range keys {
		at, known := freedAt[k]
		if !known {
			at = now
		}
		seen[k] = at
		if known && (!held || now.Sub(at) >= supersedeHoldMax) {
			drop[k] = true
		}
	}
	return drop, seen
}

// pruneSuperseded drops the freed keys in resynced (every guard re-synced after their free) and the
// unused marks.
func pruneSuperseded(marks *cilium.Map, resynced map[GuardInodeKey]bool) {
	supersede.mu.Lock()
	defer supersede.mu.Unlock()

	var gone, unused []GuardInodeKey
	var k GuardInodeKey
	var v GuardExeSupersede
	it := marks.Iterate()
	for it.Next(&k, &v) {
		switch {
		case v.Freed != 0:
			if resynced[k] {
				gone = append(gone, k)
			}
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
	k = k.Real()
	supersede.mu.Lock()
	defer supersede.mu.Unlock()
	if exeGone(k) {
		dropSupersededLocked(k)
	}
}

// supersededBinary is a superseded key a guard admits, carried to its replacement across a reload.
type supersededBinary struct {
	// path is the whitelist line the key was admitted for.
	path    string
	key     GuardInodeKey
	action  uint8
	mask    uint32
	hasMask bool
}

// inheritSuperseded admits in g, a reload's replacement whose slot is not active yet, the
// superseded keys its reload twins admit, so a process exec'd before its binary was replaced keeps
// its access from the first inode g claims. Read and written in one step under supersede.mu, which
// prune holds while dropping rows: a key freed meanwhile is never carried. A key is carried only
// while its inode exists, and only if g still whitelists its line or the line's file is gone (a
// vanished version): a line the operator dropped while its binary is installed takes its old
// versions with it. A failure carries less (that process is denied), never more.
func (g *Guard) inheritSuperseded() {
	supersede.mu.Lock()
	defer supersede.mu.Unlock()
	for _, twin := range sharedEngine.reloadTwins(g) {
		for _, b := range twin.supersededLocked() {
			if !g.carriesLine(b.path) {
				continue
			}
			if err := g.carrySupersededLocked(b); err != nil {
				log.Warnf("guard %s: superseded %s not carried (processes started before its update lose "+
					"access): %v", g.path, b.path, err)
				continue
			}
			log.Infof("guard %s: kept superseded %s (inode %d) for processes started before its update", g.path,
				b.path, b.key.Ino)
		}
	}
}

// supersededLocked lists g's superseded keys whose inode still exists: a process exec'd before the
// supersede may still run it. A row it can't read whole is left out.
func (g *Guard) supersededLocked() []supersededBinary {
	g.mu.Lock()
	paths := make(map[GuardInodeKey]string, len(g.keyPaths))
	for k, p := range g.keyPaths {
		paths[k] = p //nolint:modernize // this file's `maps` locals would shadow the package
	}
	g.mu.Unlock()

	var out []supersededBinary
	for k, p := range paths {
		if mark, ok := supersededMark(k); !ok || mark.Freed != 0 {
			continue
		}
		b := supersededBinary{path: p, key: k}
		if err := g.objs().GuardExeActions.Lookup(g.resKey(k), &b.action); err != nil {
			if !errors.Is(err, cilium.ErrKeyNotExist) {
				log.Warnf("guard %s: reading superseded %s: %v", g.path, p, err)
			}
			continue
		}
		// No mask row means every event: an unreadable one must not pass for that.
		switch err := g.objs().GuardExeEvents.Lookup(g.resKey(k), &b.mask); {
		case err == nil:
			b.hasMask = true
		case !errors.Is(err, cilium.ErrKeyNotExist):
			log.Warnf("guard %s: reading the event mask of superseded %s: %v", g.path, p, err)
			continue
		}
		out = append(out, b)
	}
	return out
}

// carrySupersededLocked writes b's mask before its action: an action row alone admits every event.
func (g *Guard) carrySupersededLocked(b supersededBinary) error {
	if b.hasMask {
		if err := g.objs().GuardExeEvents.Put(g.resKey(b.key), b.mask); err != nil {
			return err
		}
	}
	if err := g.putExeAction(b.key, b.action); err != nil {
		return err
	}
	g.mu.Lock()
	g.keyPaths[b.key] = b.path
	g.mu.Unlock()
	return nil
}

func (g *Guard) carriesLine(path string) bool {
	g.mu.Lock()
	if _, ok := g.links[path]; ok {
		g.mu.Unlock()
		return true
	}
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

// forgetExe drops what this guard remembers about k's rows (k a real key: its applet keys too), so a
// file reusing the number at a whitelisted path is judged as a replacement of k (retired), never as k.
func (g *Guard) forgetExe(k GuardInodeKey) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for p, key := range g.deployed {
		if key.Real() == k {
			delete(g.deployed, p)
			if g.retired == nil {
				g.retired = make(map[string]GuardInodeKey)
			}
			g.retired[p] = k
		}
	}
	for p, st := range g.binaryVerifyStates {
		if st.key.Real() == k {
			delete(g.binaryVerifyStates, p)
		}
	}
	for key := range g.keyPaths {
		if key.Real() == k {
			delete(g.keyPaths, key)
		}
	}
}

const (
	supersedePruneEvery = 30 * time.Second
	// supersedeHoldMax caps SetSupersedeHold: a held key keeps its exe_superseded entry, and a full
	// map refuses unlinking any whitelisted binary.
	supersedeHoldMax = 5 * time.Minute
)

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
