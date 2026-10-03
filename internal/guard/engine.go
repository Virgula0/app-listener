package guard

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	cilium "github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	log "github.com/sirupsen/logrus"
)

// GuardMaxRes mirrors GUARD_MAX_RES in guard.bpf.c: the size of the resource table.
const GuardMaxRes = 512

// resGlobal mirrors GUARD_RES_GLOBAL: the reserved slot for decisions that belong to no single
// resource (the raw block-device gate, and a process tainted by resources no taint set covers).
// Its whitelist is the intersection of every resource's, so it is never weaker than the separate
// per-resource guards it replaced. Real resources are allocated from 1, taint sets from the top
// (taintset.go).
const resGlobal uint32 = 0

// engine owns the single loaded copy of the guard BPF objects and the single set of LSM links that
// every Guard shares.
//
// Attaching per Guard instead costs one trampoline link per resource on each of the 26 attach
// points, and the kernel caps those at BPF_MAX_TRAMP_LINKS (38) — a limit shared by everything on
// the host, so the daemon both hit a ceiling on how many resources it could enforce and starved
// every other BPF-LSM user. Enforcement state is keyed by dev:ino, so one program set resolves the
// accessed inode to its resource in-kernel and attach usage stays O(1).
//
// The links outlive individual Guards: they are attached on the first Guard and detached only when
// the last one stops. A reload that replaces every Guard therefore never drops enforcement, which
// is what the attach-before-detach ordering guaranteed before.
type engine struct {
	mu   sync.Mutex
	objs GuardObjects
	// links and linkNames run in step: an optional hook that fails to attach is skipped, so the
	// pin basename cannot be recovered from the link's index alone.
	links     []link.Link
	linkNames []string
	// slots maps res_id to the owning Guard so ringbuf events reach the right reader. Index 0 is
	// resGlobal and never holds a Guard.
	slots [GuardMaxRes]*Guard
	// prevOwner records, per inode, the LIVE resource whose guard_inodes row a later claim
	// overwrote (a reload's replacement guard taking over the resource it replaces). On that
	// replacement's Stop the row is handed BACK to the displaced owner if it is still live, instead
	// of being deleted — otherwise a refused/rolled-back reload, or a failed buildGuards, would
	// leave the kept resource with no inode row while its old guard is still attached (fail open).
	prevOwner map[GuardInodeKey]uint32
	// keptLogged: per inode, the last resource told it may not take the row (log once).
	keptLogged map[GuardInodeKey]uint32
	refs       int
	rd         *ringbuf.Reader
	done       chan struct{}
	started    bool
	// ioctlCompat: this kernel has file_ioctl_compat (6.8+), so guard_file_ioctl_compat attaches.
	ioctlCompat bool

	// allows tracks each resource's whitelisted exe inodes and their action (GUARD_ALLOW or
	// GUARD_ALLOW_ROOT), the source for the union and intersection views (noteAllow).
	allows map[uint32]map[GuardInodeKey]uint8
	// Taint sets (taintset.go): key -> slot and back, and the rows last written to the kernel.
	sets       map[string]uint32
	setAt      map[uint32]string
	setRows    map[uint32]map[GuardInodeKey]uint8
	unionRows  map[GuardInodeKey]uint32
	memberRows map[GuardTaintMemberKey]struct{}
	// successor mirrors guard_taint_successor: old slot -> its reload replacement (same path).
	successor map[uint32]uint32
	// rawOwner receives the reserved slot's events (the raw block-device gate): it is the guard
	// that registered the backing devices, whose consumer labels them as such. Delivering them to
	// every guard duplicated each denial and let the self-guards label it with their own path.
	rawOwner *Guard

	// pinPrefix is where the shared links are pinned (ConfigureSharedPinning). Empty = no pinning.
	pinPrefix   string
	pinDegraded bool
	mapsPinned  bool
}

var sharedEngine engine

// ConfigureSharedPinning sets where the shared LSM links are pinned so enforcement survives a
// SIGKILL. Call it before building each batch of guards, passing that batch's generation prefix.
//
// The shared links outlive a reload (they are detached only when the last Guard stops), so a
// reload MOVES the existing pins to the new generation instead of leaving them under the retired
// one — otherwise `daemon --lockdown` would look for them under the generation pin-state records
// and find nothing, losing crash recovery for file vaults.
func ConfigureSharedPinning(prefix string) {
	sharedEngine.mu.Lock()
	defer sharedEngine.mu.Unlock()

	if !sharedEngine.started {
		sharedEngine.pinPrefix = prefix
		return
	}
	if prefix == "" || prefix == sharedEngine.pinPrefix || sharedEngine.pinDegraded {
		return
	}
	sharedEngine.repinLocked(prefix)
}

// repinLocked moves the shared pins to a new generation prefix. A failure mid-move leaves
// enforcement live but not kill-safe, which is reported like any other pin failure.
func (e *engine) repinLocked(prefix string) {
	old := e.pinPrefix
	for i, l := range e.links {
		if err := l.Pin(prefix + e.linkNames[i]); err != nil {
			log.Errorf("guard: CRITICAL: moving LSM link pin to %s failed (%v) — enforcement stays live "+
				"but will NOT survive a SIGKILL until the next restart", prefix, err)
			e.pinDegraded = true
			return
		}
	}
	if e.mapsPinned {
		for _, spec := range []struct {
			name string
			m    *cilium.Map
		}{
			{ExeActionsPinName, e.objs.GuardExeActions},
			{ExeEventsPinName, e.objs.GuardExeEvents},
		} {
			if err := spec.m.Pin(prefix + spec.name); err != nil {
				log.Errorf("guard: moving pinned %s to %s failed (%v) — 'daemon --lockdown' loses "+
					"self-access widening for this generation", spec.name, prefix, err)
			}
		}
	}
	e.pinPrefix = prefix
	log.Infof("guard engine: shared pins moved from %s* to %s*", old, prefix)
}

// SharedPinDegraded reports that pinning was requested but refused: enforcing while alive, but the
// links will not survive a SIGKILL.
func SharedPinDegraded() bool {
	sharedEngine.mu.Lock()
	defer sharedEngine.mu.Unlock()
	return sharedEngine.pinDegraded
}

// acquire loads and attaches the shared programs on first use and registers g in a resource slot.
// On any failure nothing is left attached for this caller.
func (e *engine) acquire(g *Guard) (uint32, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if !e.started {
		if err := e.startLocked(); err != nil {
			return 0, err
		}
	}

	id, err := e.allocSlotLocked(g)
	if err != nil {
		// Nothing else holds the engine: unwind the attach rather than leave hooks live with no
		// resource behind them.
		if e.refs == 0 {
			e.stopLocked()
		}
		return 0, err
	}
	e.refs++
	e.linkReloadTwinLocked(g, id)
	return id, nil
}

// linkReloadTwinLocked records id as the replacement of every live slot guarding g's path, before
// id claims any inode. Without it a process tainted under the old slot that reads the tree through
// the new one before carryTaintAcrossReload moves it is merged to GLOBAL for life. A failed write
// only leaves that merge (stricter), so it is logged, not fatal.
func (e *engine) linkReloadTwinLocked(g *Guard, id uint32) {
	for old := uint32(1); old < GuardMaxRes; old++ {
		o := e.slots[old]
		if old == id || o == nil || filepath.Clean(o.path) != filepath.Clean(g.path) {
			continue
		}
		if err := e.objs.GuardTaintSuccessor.Put(old, id); err != nil {
			log.Warnf("guard: linking reload slot %d to %d for %s: %v", old, id, g.path, err)
			continue
		}
		if e.successor == nil {
			e.successor = make(map[uint32]uint32)
		}
		e.successor[old] = id
	}
}

// unlinkReloadTwinLocked drops every successor row naming id before the slot can be reused: a
// stale row would collapse an unrelated resource's taint into id's.
func (e *engine) unlinkReloadTwinLocked(id uint32) {
	for old, next := range e.successor {
		if old != id && next != id {
			continue
		}
		if err := e.objs.GuardTaintSuccessor.Put(old, uint32(0)); err != nil {
			log.Warnf("guard: clearing reload link of slot %d: %v", old, err)
		}
		delete(e.successor, old)
	}
}

// startLocked loads the objects, attaches every LSM program once and starts the event reader.
func (e *engine) startLocked() error {
	shared, mapsErr := supersedeMaps()
	if mapsErr != nil {
		return mapsErr
	}
	spec, ioctlCompat, specErr := guardSpec()
	if specErr != nil {
		return specErr
	}
	var objs GuardObjects
	if err := spec.LoadAndAssign(&objs, &cilium.CollectionOptions{MapReplacements: shared}); err != nil {
		return fmt.Errorf("loading guard BPF objects: %w", err)
	}
	e.objs = objs
	e.ioctlCompat = ioctlCompat
	e.done = make(chan struct{})

	failedRequired, total := e.attachLocked()
	if len(failedRequired) > 0 {
		e.stopLocked()
		return fmt.Errorf(
			"required LSM hooks failed to attach: %v — read/write/open protection unavailable. "+
				"Ensure your kernel supports BPF LSM (CONFIG_BPF_LSM=y) and LSM=bpf is in the "+
				"boot command line (/sys/kernel/security/lsm). "+
				"The guard REQUIRES the file_open and file_permission hooks; if they cannot attach, "+
				"the guard cannot provide meaningful protection and will refuse to start",
			failedRequired)
	}

	lsmAttached := len(e.links)
	if err := e.attachForkLocked(); err != nil {
		e.stopLocked()
		return err
	}

	// The reserved slot must exist before any decision can land on it: res_cfg() treats an
	// inactive slot as unknown and denies, which would make the raw block-device gate refuse
	// everyone (and a process tainted by two resources un-ptraceable) rather than consult the
	// intersection whitelist refreshGlobalLocked maintains.
	if err := e.objs.GuardResConfig.Put(resGlobal, GuardResConfig{
		Active: 1,
		Mode:   uint64(ModeWhitelist),
	}); err != nil {
		e.stopLocked()
		return fmt.Errorf("initializing the shared resource slot: %w", err)
	}

	rd, err := ringbuf.NewReader(e.objs.Rb)
	if err != nil {
		e.stopLocked()
		return fmt.Errorf("opening guard ringbuf: %w", err)
	}
	e.rd = rd
	e.started = true
	go e.readLoop(rd)
	go supersedePruneLoop(e.done)

	log.Infof("guard engine — %d/%d LSM hooks attached once for every resource", lsmAttached, total)
	return nil
}

// attachLocked attaches every LSM program, pinning each at pinPrefix+<hook> when enabled. A pin
// failure degrades (pinDegraded, CRITICAL log) and drops the pins rather than aborting.
func (e *engine) attachLocked() (failedRequired []string, total int) {
	attachments := guardLSMHooks(&e.objs, e.ioctlCompat)
	for _, a := range attachments {
		l, attachErr := link.AttachLSM(link.LSMOptions{Program: a.prog})
		if attachErr != nil {
			if requiredHooks[a.hook] {
				failedRequired = append(failedRequired, a.hook)
				log.Errorf("CRITICAL: required LSM hook %s failed to attach: %v", a.hook, attachErr)
			} else {
				log.Warnf("skipping optional LSM hook %s: %v", a.hook, attachErr)
			}
			continue
		}
		e.links = append(e.links, l)
		e.linkNames = append(e.linkNames, strings.ReplaceAll(a.hook, "_", "-"))
		if e.pinPrefix != "" {
			// Some hardened bpffs only accept [a-z0-9-] in pin names.
			pinPath := e.pinPrefix + e.linkNames[len(e.linkNames)-1]
			if pinErr := l.Pin(pinPath); pinErr != nil {
				log.Errorf("guard: CRITICAL: LSM link pinning failed at %s (%v) — enforcement will NOT "+
					"survive a SIGKILL. The daemon keeps running with live enforcement; investigate bpffs "+
					"(kernel hardening, mount options).", pinPath, pinErr)
				for _, prev := range e.links {
					_ = prev.Unpin()
				}
				e.pinPrefix = ""
				e.pinDegraded = true
			}
		}
	}
	return failedRequired, len(attachments)
}

// attachForkLocked attaches guard_sched_process_fork, which copies a parent's taint and exec
// stamp to its child. Required: a child without its parent's stamp would pass for a process
// started before the guard, so a superseded image re-exec'd by its parent would admit the child.
func (e *engine) attachForkLocked() error {
	l, err := link.AttachTracing(link.TracingOptions{
		Program:    e.objs.GuardSchedProcessFork,
		AttachType: cilium.AttachTraceRawTp,
	})
	if err != nil {
		return fmt.Errorf("required hook sched_process_fork failed to attach: %w — a forked child "+
			"would lose its parent's exec stamp", err)
	}
	if e.pinPrefix != "" {
		if pinErr := l.Pin(e.pinPrefix + "sched-process-fork"); pinErr != nil {
			log.Warnf("guard: pinning fork taint propagation failed (%v) — it will not survive a "+
				"SIGKILL; enforcement links are unaffected", pinErr)
			_ = l.Unpin()
		}
	}
	e.links = append(e.links, l)
	e.linkNames = append(e.linkNames, "sched-process-fork")
	return nil
}

// FreeResourceSlots is how many resource slots a new Guard can still claim.
func FreeResourceSlots() int {
	sharedEngine.mu.Lock()
	defer sharedEngine.mu.Unlock()
	free := 0
	for id := 1; id < GuardMaxRes; id++ {
		if sharedEngine.slots[id] == nil && sharedEngine.setAt[uint32(id)] == "" { //nolint:gosec // id < GuardMaxRes
			free++
		}
	}
	return free
}

// allocSlotLocked reserves a resource id for g. Slot 0 is reserved (resGlobal), and so is every
// taint-set slot.
func (e *engine) allocSlotLocked(g *Guard) (uint32, error) {
	for id := 1; id < GuardMaxRes; id++ {
		if e.slots[id] == nil && e.setAt[uint32(id)] == "" { //nolint:gosec // id < GuardMaxRes
			e.slots[id] = g
			return uint32(id), nil //nolint:gosec // id < GuardMaxRes
		}
	}
	return 0, fmt.Errorf("no free guard resource slot (max %d)", GuardMaxRes-1)
}

// claimInode maps key (found at statPath, logically path) to g's resource unless path lies in a
// live resource rooted strictly inside g's root. guard_inodes holds one owner per inode, so the
// innermost resource must win whatever order the guards scan in (the sealed fscrypt.key inside the
// read-only /etc/app-listener guard). Checked and written under e.mu: an inner guard registers its
// slot before scanning, so an outer scan can't overwrite its rows on a stale check.
//
// A row of another live resource, or an in-tree symlink's target, is taken only when the inode's
// physical chain says g is its nearest root (mayClaimLocked): a path that came to resolve into
// another tree must not hand that tree's inodes to g's mode.
func (e *engine) claimInode(g *Guard, path, statPath string, key GuardInodeKey, viaLink bool) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.innerResourceOwnsLocked(g, filepath.Clean(path)) {
		return nil
	}
	var cur uint32
	owner := uint32(0)
	if e.objs.GuardInodes.Lookup(key, &cur) == nil && cur != g.resID &&
		cur < GuardMaxRes && e.slots[cur] != nil {
		owner = cur
	}
	// A reload's replacement over the same root takes the tree it replaces.
	replacing := owner != 0 && filepath.Clean(e.slots[owner].path) == filepath.Clean(g.path)
	if !replacing && (owner != 0 || viaLink) && !e.mayClaimLocked(g, path, statPath, key, owner, viaLink) {
		return nil
	}
	// Remember a displaced live owner so releaseInodes can hand the row back rather than delete it.
	if owner != 0 {
		if e.prevOwner == nil {
			e.prevOwner = make(map[GuardInodeKey]uint32)
		}
		e.prevOwner[key] = owner
	}
	return e.objs.GuardInodes.Put(key, g.resID)
}

// mayClaimLocked applies the kernel's nearest-root rule to the inode's physical chain. owner's row
// only guards where its root is on the chain (root_in_chain), so a row whose owner root is absent
// is stale and g may take it; one whose owner root is nearer than g's is not g's. An unknowable
// chain keeps the current owner.
func (e *engine) mayClaimLocked(g *Guard, path, statPath string, key GuardInodeKey, owner uint32, viaLink bool) bool {
	chain, ok := inodeChain(statPath, key, viaLink)
	if !ok {
		if owner != 0 {
			e.warnKeptLocked(g, path, key, owner, "its physical location could not be verified")
		}
		return false
	}
	// g's reload twin shares its root: not a competing resource.
	self := func(id uint32) bool { return id == g.resID || e.slots[id].path == g.path }
	outranks := func(id uint32) bool { return tieRank(g) > tieRank(e.slots[id]) }
	if nearestRootClaims(chain, e.rootsLocked(), self, owner, outranks) {
		return true
	}
	if owner != 0 {
		e.warnKeptLocked(g, path, key, owner, "it lies in that resource's tree")
	}
	return false
}

// tieRank orders resources anchored on the same root inode under different paths, i.e. one path
// re-pointed onto another resource: a path that IS its root beats one reaching it through a
// symlink, then a whitelist beats a read-only (world-readable) mode.
func tieRank(g *Guard) int {
	r := 0
	if g.canonicalRoot.Load() {
		r += 2
	}
	if g.mode != ModeReadOnly {
		r++
	}
	return r
}

// outrankedOnRoot names a live resource (not g's reload twin) anchored on root that g must yield
// to: a higher tieRank, or an equal one that holds the root row already.
func (e *engine) outrankedOnRoot(g *Guard, root GuardInodeKey) string {
	e.mu.Lock()
	defer e.mu.Unlock()
	var rowOwner uint32
	_ = e.objs.GuardInodes.Lookup(root, &rowOwner)
	for id, r := range e.rootsLocked() {
		o := e.slots[id]
		if r != root || id == g.resID || o.path == g.path {
			continue
		}
		if tieRank(o) > tieRank(g) || (tieRank(o) == tieRank(g) && rowOwner == id) {
			return o.path
		}
	}
	return ""
}

// nearestRootClaims applies root_in_chain's view to chain (nearest first). A live owner's row is
// claimable when its root is off the chain (stale) or farther than self's, or tied and outranked.
// A symlink target (owner 0) only inside self's tree with no other root nearer, or tied and not
// outranked.
func nearestRootClaims(chain []GuardInodeKey, roots map[uint32]GuardInodeKey, self func(uint32) bool,
	owner uint32, outranks func(uint32) bool) bool {
	nearestOf := func(match func(uint32) bool) int {
		for i, k := range chain {
			for id, r := range roots {
				if r == k && match(id) {
					return i
				}
			}
		}
		return -1
	}
	sIdx := nearestOf(self)
	if owner != 0 {
		oIdx := nearestOf(func(id uint32) bool { return id == owner })
		return oIdx < 0 || (sIdx >= 0 && (sIdx < oIdx || (sIdx == oIdx && outranks(owner))))
	}
	blocking := nearestOf(func(id uint32) bool {
		if self(id) {
			return false
		}
		i := slices.Index(chain, roots[id])
		return i < sIdx || !outranks(id)
	})
	return sIdx >= 0 && (blocking < 0 || blocking > sIdx)
}

// rootsLocked returns every live resource's root anchor as the kernel sees it (guard_res_config),
// read there rather than from Guard.rootKey: Guard.mu may be held across a call into the engine.
func (e *engine) rootsLocked() map[uint32]GuardInodeKey {
	roots := make(map[uint32]GuardInodeKey)
	for id, o := range &e.slots {
		if o == nil {
			continue
		}
		var cfg GuardResConfig
		if e.objs.GuardResConfig.Lookup(uint32(id), &cfg) == nil { //nolint:gosec // id < GuardMaxRes
			roots[uint32(id)] = GuardInodeKey{Dev: cfg.RootDev, Ino: cfg.RootIno} //nolint:gosec // id < GuardMaxRes
		}
	}
	return roots
}

// liveOwner returns the live resource other than self whose row or root anchor key is, 0 if none.
func (e *engine) liveOwner(self uint32, key GuardInodeKey) uint32 {
	e.mu.Lock()
	defer e.mu.Unlock()
	var cur uint32
	if e.objs.GuardInodes.Lookup(key, &cur) == nil && cur != self && cur < GuardMaxRes && e.slots[cur] != nil {
		return cur
	}
	for id, r := range e.rootsLocked() {
		if id != self && r == key {
			return id
		}
	}
	return 0
}

func (e *engine) warnKeptLocked(g *Guard, path string, key GuardInodeKey, owner uint32, why string) {
	if e.keptLogged == nil {
		e.keptLogged = make(map[GuardInodeKey]uint32)
	}
	if e.keptLogged[key] == g.resID {
		return
	}
	e.keptLogged[key] = g.resID
	log.Errorf("guard %s: CRITICAL: %s (inode %d) is guarded by resource %s and %s — not taking it over; "+
		"a path in this resource resolves into another resource's tree", g.path, path, key.Ino,
		e.slots[owner].path, why)
}

// releaseInodes drops resID's guard_inodes rows on Stop. A row whose previous owner (displaced by a
// reload claim) is still live is handed BACK to it instead of deleted, so a resource kept across a
// refused/rolled-back reload never loses its inode rows while its old guard is still attached. Done
// under e.mu (like claimInode) so the check-and-restore is atomic against concurrent claims.
func (e *engine) releaseInodes(resID uint32) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.objs.GuardInodes == nil {
		return nil // engine already torn down
	}

	var owned []GuardInodeKey
	var key GuardInodeKey
	var val uint32
	it := e.objs.GuardInodes.Iterate()
	for it.Next(&key, &val) {
		if val == resID {
			owned = append(owned, key)
		}
	}
	if err := it.Err(); err != nil {
		return err
	}

	for _, k := range owned {
		prev, hasPrev := e.prevOwner[k]
		if hasPrev && prev != resID && prev < GuardMaxRes && e.slots[prev] != nil {
			if err := e.objs.GuardInodes.Put(k, prev); err != nil {
				return err
			}
		} else if err := e.objs.GuardInodes.Delete(k); err != nil &&
			!errors.Is(err, cilium.ErrKeyNotExist) {
			return err
		}
		delete(e.prevOwner, k)
	}
	// Drop any remaining hand-back records that point AT this (now-gone) resource.
	for k, v := range e.prevOwner {
		if v == resID {
			delete(e.prevOwner, k)
		}
	}
	return nil
}

// innerResourceOwnsLocked: a replacement guard over the same root (reload) is not "inner".
func (e *engine) innerResourceOwnsLocked(g *Guard, path string) bool {
	outer := filepath.Clean(g.path)
	for _, o := range &e.slots {
		if o == nil || o == g {
			continue
		}
		root := filepath.Clean(o.path)
		if root != outer && pathWithin(root, outer) && pathWithin(path, root) {
			return true
		}
	}
	return false
}

func pathWithin(path, root string) bool {
	return path == root || strings.HasPrefix(path, root+"/")
}

// release retires g's slot and, once the last Guard is gone, detaches the shared programs.
func (e *engine) release(id uint32) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if id < GuardMaxRes && e.slots[id] != nil {
		e.slots[id] = nil
		// Retire the kernel-side slot so an inode entry that outlives the Guard cannot be judged
		// by a later resource that reuses the id.
		if e.started {
			if err := e.objs.GuardResConfig.Put(id, GuardResConfig{}); err != nil {
				log.Warnf("guard: clearing resource slot %d: %v", id, err)
			}
			e.unlinkReloadTwinLocked(id)
			// Drop the id from every taint set now: a later resource reusing it must not be
			// judged as a member of a set it never joined.
			if err := e.syncTaintLocked(); err != nil {
				log.Warnf("guard: refreshing taint sets after retiring slot %d: %v", id, err)
			}
		}
		e.refs--
	}
	if e.refs <= 0 {
		e.stopLocked()
	}
}

// stopLocked detaches everything and resets the engine so a later Guard starts clean.
func (e *engine) stopLocked() {
	e.unpinSharedMaps()
	if e.done != nil {
		select {
		case <-e.done:
		default:
			close(e.done)
		}
	}
	if e.rd != nil {
		_ = e.rd.Close()
		e.rd = nil
	}
	for _, l := range e.links {
		if e.pinPrefix != "" {
			_ = l.Unpin()
		}
		_ = l.Close()
	}
	e.links = nil
	e.linkNames = nil
	_ = e.objs.Close()
	e.objs = GuardObjects{}
	e.started = false
	e.refs = 0
	e.done = nil
	e.sets, e.setAt, e.setRows, e.unionRows, e.memberRows = nil, nil, nil, nil, nil
	e.prevOwner = nil
	e.successor = nil
}

// readLoop drains the shared ringbuf and routes each event to the Guard that owns its resource.
// Takes the reader as an argument: stopLocked clears e.rd, which this goroutine outlives.
func (e *engine) readLoop(rd *ringbuf.Reader) {
	for {
		rec, err := rd.Read()
		if err != nil {
			if errors.Is(err, ringbuf.ErrClosed) {
				return
			}
			log.Errorf("guard: ringbuf read error: %v", err)
			continue
		}
		ev, resID, ok := parseGuardEvent(rec.RawSample)
		if !ok {
			continue
		}
		e.deliver(resID, ev)
	}
}

// deliver hands an event to its owning Guard. An event from the reserved slot belongs to no
// resource; it goes to exactly one guard (see rawOwner), as the gate used to live in one guard.
func (e *engine) deliver(resID uint32, ev *GuardEvent) {
	e.mu.Lock()
	var target *Guard
	switch {
	case resID == resGlobal:
		target = e.globalOwnerLocked()
	case e.setAt[resID] != "":
		target = e.setOwnerLocked(resID)
	case resID < GuardMaxRes:
		target = e.slots[resID]
	}
	e.mu.Unlock()

	if target != nil {
		target.dispatch(ev)
	}
}

// globalOwnerLocked is the guard that reports reserved-slot events: the raw-device owner while it
// lives, else the lowest live slot so a denial is never silently dropped.
func (e *engine) globalOwnerLocked() *Guard {
	if e.rawOwner != nil {
		for _, g := range &e.slots {
			if g == e.rawOwner {
				return g
			}
		}
		e.rawOwner = nil
	}
	for _, g := range &e.slots {
		if g != nil {
			return g
		}
	}
	return nil
}

// claimRawDevices makes g the reporter of raw block-device denials (it registered the devices).
func (e *engine) claimRawDevices(g *Guard) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rawOwner = g
}

// noteAllow records that res permits exe, and refreshes the cross-resource views:
//
//   - guard_exe_union (the slot judging an exe's taint: its resource, or its taint set) answers
//     bprm_check/bprm_committed_creds, which see an exec with no guarded inode to resolve;
//   - the taint sets' whitelists and the reserved resGlobal whitelist, which holds the
//     INTERSECTION of every resource's for decisions that belong to no single resource.
func (e *engine) noteAllow(res uint32, exe GuardInodeKey, action uint8) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.allows == nil {
		e.allows = make(map[uint32]map[GuardInodeKey]uint8)
	}
	if e.allows[res] == nil {
		e.allows[res] = make(map[GuardInodeKey]uint8)
	}
	e.allows[res][exe] = action
	return e.syncSharedLocked()
}

// noteDeny withdraws an allow (a whitelisted binary demoted after in-place tampering, or an entry
// dropped on re-sync) from the cross-resource views, so the tampered image loses the shared
// whitelist too.
func (e *engine) noteDeny(res uint32, exe GuardInodeKey) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.allows[res] == nil {
		return nil
	}
	if _, had := e.allows[res][exe]; !had {
		return nil
	}
	delete(e.allows[res], exe)
	return e.syncSharedLocked()
}

// admitsExe reports whether any live resource allows exe.
func (e *engine) admitsExe(exe GuardInodeKey) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, set := range e.allows {
		if _, ok := set[exe]; ok {
			return true
		}
	}
	return false
}

// forgetExe drops every row naming exe, in every resource and the views derived from them, and has
// each guard forget it (Guard.forgetExe).
func (e *engine) forgetExe(exe GuardInodeKey) error {
	e.mu.Lock()
	if !e.started {
		e.mu.Unlock()
		return nil
	}
	for _, set := range e.allows {
		delete(set, exe)
	}
	err := e.syncSharedLocked()
	if err == nil {
		err = deleteInoKeys(e.objs.GuardExeActions, exe)
	}
	if err == nil {
		err = deleteInoKeys(e.objs.GuardExeEvents, exe)
	}
	var guards []*Guard
	for _, g := range &e.slots {
		if g != nil {
			guards = append(guards, g)
		}
	}
	e.mu.Unlock()
	for _, g := range guards {
		g.forgetExe(exe)
	}
	return err
}

// forgetResource drops res from the cross-resource views when its guard goes away.
func (e *engine) forgetResource(res uint32) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.allows == nil {
		return
	}
	delete(e.allows, res)
	if err := e.syncSharedLocked(); err != nil {
		log.Warnf("guard: refreshing shared views after dropping resource %d: %v", res, err)
	}
}

func (e *engine) syncSharedLocked() error {
	if err := e.syncTaintLocked(); err != nil {
		return err
	}
	return e.refreshGlobalLocked()
}

// refreshGlobalLocked rewrites the reserved slot's whitelist with the intersection of every live
// resource's, so a cross-resource decision needs an allow from all of them.
func (e *engine) refreshGlobalLocked() error {
	if !e.started {
		return nil
	}
	if err := e.clearGlobalAllowsLocked(); err != nil {
		return err
	}
	for exe, action := range e.allowIntersectionLocked() {
		k := GuardResInodeKey{ResId: resGlobal, Ino: exe}
		if err := e.objs.GuardExeActions.Put(k, action); err != nil {
			return err
		}
	}
	return nil
}

// allowIntersectionLocked is the set of exe inodes every live resource allows (intersectAllows).
// Flattening GUARD_ALLOW_ROOT to GUARD_ALLOW would both widen a root-only grant and break the
// checks that look for the root-gated self entry (ptrace_access_check's metadata exception).
func (e *engine) allowIntersectionLocked() map[GuardInodeKey]uint8 {
	sets := make([]map[GuardInodeKey]uint8, 0, len(e.allows))
	for _, set := range e.allows {
		sets = append(sets, set)
	}
	return intersectAllows(sets)
}

// clearGlobalAllowsLocked drops every row of the reserved slot's whitelist.
func (e *engine) clearGlobalAllowsLocked() error {
	var stale []GuardResInodeKey
	var key GuardResInodeKey
	var val uint8
	it := e.objs.GuardExeActions.Iterate()
	for it.Next(&key, &val) {
		if key.ResId == resGlobal {
			stale = append(stale, key)
		}
	}
	if err := it.Err(); err != nil {
		return err
	}
	for i := range stale {
		if err := e.objs.GuardExeActions.Delete(stale[i]); err != nil &&
			!errors.Is(err, cilium.ErrKeyNotExist) {
			return err
		}
	}
	return nil
}

// pinSharedMaps pins the shared whitelist maps once, so `daemon --lockdown` can widen a file
// vault's self-access after the daemon died. Idempotent: the maps are pinned on first request and
// re-pinning would MOVE the pin, orphaning the earlier path.
func (e *engine) pinSharedMaps() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.pinPrefix == "" || e.mapsPinned || !e.started {
		return
	}
	for _, spec := range []struct {
		name string
		m    *cilium.Map
	}{
		{ExeActionsPinName, e.objs.GuardExeActions},
		{ExeEventsPinName, e.objs.GuardExeEvents},
	} {
		if err := spec.m.Pin(e.pinPrefix + spec.name); err != nil {
			log.Errorf("guard: pinning %s failed (%v) — 'daemon --lockdown' will not be able to widen "+
				"self-access if the daemon crashes while a vault is unlocked", spec.name, err)
			return
		}
	}
	e.mapsPinned = true
}

// unpinSharedMaps drops the shared map pins on clean shutdown.
func (e *engine) unpinSharedMaps() {
	if !e.mapsPinned {
		return
	}
	for _, name := range []string{ExeActionsPinName, ExeEventsPinName} {
		if err := os.Remove(e.pinPrefix + name); err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Errorf("guard: removing pinned %s failed (%v) — remove it manually under %s*",
				name, err, e.pinPrefix)
		}
	}
	e.mapsPinned = false
}
