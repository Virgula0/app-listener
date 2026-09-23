package guard

import (
	"errors"
	"fmt"
	"os"
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
// resource (the raw block-device gate, and a process tainted by more than one resource). Its
// whitelist is the intersection of every resource's, so it is never weaker than the separate
// per-resource guards it replaced. Real resources are allocated from 1.
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
	slots   [GuardMaxRes]*Guard
	refs    int
	rd      *ringbuf.Reader
	done    chan struct{}
	started bool

	// allows tracks each resource's whitelisted exe inodes and their action (GUARD_ALLOW or
	// GUARD_ALLOW_ROOT), the source for the union and intersection views (noteAllow).
	allows map[uint32]map[GuardInodeKey]uint8
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
	return id, nil
}

// startLocked loads the objects, attaches every LSM program once and starts the event reader.
func (e *engine) startLocked() error {
	var objs GuardObjects
	if err := LoadGuardObjects(&objs, nil); err != nil {
		return fmt.Errorf("loading guard BPF objects: %w", err)
	}
	e.objs = objs
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
	e.attachForkTaintLocked()

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

	log.Infof("guard engine — %d/%d LSM hooks attached once for every resource", lsmAttached, total)
	return nil
}

// attachLocked attaches every LSM program, pinning each at pinPrefix+<hook> when enabled. A pin
// failure degrades (pinDegraded, CRITICAL log) and drops the pins rather than aborting.
func (e *engine) attachLocked() (failedRequired []string, total int) {
	attachments := guardLSMHooks(&e.objs)
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

// attachForkTaintLocked attaches guard_sched_process_fork (copies a tainted parent's taint to
// forked children). Best-effort: failure only delays a child's tracking until its own guarded
// access.
func (e *engine) attachForkTaintLocked() {
	l, err := link.AttachTracing(link.TracingOptions{
		Program:    e.objs.GuardSchedProcessFork,
		AttachType: cilium.AttachTraceRawTp,
	})
	if err != nil {
		log.Warnf("guard: skipping fork taint propagation (%v) — a forked child of a process that "+
			"read guarded content is not taint-tracked until its own first guarded access; direct "+
			"enforcement is unaffected", err)
		return
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
}

// allocSlotLocked reserves a resource id for g. Slot 0 is reserved (resGlobal).
func (e *engine) allocSlotLocked(g *Guard) (uint32, error) {
	for id := 1; id < GuardMaxRes; id++ {
		if e.slots[id] == nil {
			e.slots[id] = g
			return uint32(id), nil //nolint:gosec // id < GuardMaxRes
		}
	}
	return 0, fmt.Errorf("no free guard resource slot (max %d)", GuardMaxRes-1)
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

// noteAllow records that res permits exe, and refreshes the two cross-resource views:
//
//   - guard_exe_union ("allowed by at least one resource") answers bprm_committed_creds, which
//     sees an exec with no guarded inode to resolve a resource from;
//   - the reserved resGlobal whitelist holds the INTERSECTION, used where a decision belongs to no
//     single resource. N separate guards each vetoed such an access independently, so anything
//     short of the intersection would be weaker than the model this replaced.
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

	owner := res
	var cur uint32
	switch err := e.objs.GuardExeUnion.Lookup(exe, &cur); {
	case err == nil:
		if cur != res {
			owner = resGlobal // allowed by several resources
		}
	case errors.Is(err, cilium.ErrKeyNotExist):
	default:
		return fmt.Errorf("reading exe union: %w", err)
	}
	if err := e.objs.GuardExeUnion.Put(exe, owner); err != nil {
		return fmt.Errorf("writing exe union: %w", err)
	}
	return e.refreshGlobalLocked()
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
	if err := e.rebuildUnionLocked(); err != nil {
		return err
	}
	return e.refreshGlobalLocked()
}

// forgetResource drops res from the cross-resource views when its guard goes away.
func (e *engine) forgetResource(res uint32) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.allows == nil {
		return
	}
	delete(e.allows, res)
	if !e.started {
		return
	}
	if err := e.rebuildUnionLocked(); err != nil {
		log.Warnf("guard: rebuilding exe union after dropping resource %d: %v", res, err)
	}
	if err := e.refreshGlobalLocked(); err != nil {
		log.Warnf("guard: refreshing shared whitelist after dropping resource %d: %v", res, err)
	}
}

// rebuildUnionLocked rewrites guard_exe_union from the live per-resource allow sets.
func (e *engine) rebuildUnionLocked() error {
	want := make(map[GuardInodeKey]uint32)
	for res, set := range e.allows {
		for exe := range set {
			if prev, seen := want[exe]; seen && prev != res {
				want[exe] = resGlobal
				continue
			}
			want[exe] = res
		}
	}
	return replaceInodeU32Map(e.objs.GuardExeUnion, want)
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

// allowIntersectionLocked is the set of exe inodes every live resource allows, with the stricter
// action where they differ: GUARD_ALLOW_ROOT (uid 0 only) wins over GUARD_ALLOW. Flattening to
// GUARD_ALLOW would both widen a root-only grant and break the checks that look for the root-gated
// self entry specifically (ptrace_access_check's metadata exception for the daemon's own process).
func (e *engine) allowIntersectionLocked() map[GuardInodeKey]uint8 {
	var inter map[GuardInodeKey]uint8
	for _, set := range e.allows {
		if inter == nil {
			inter = make(map[GuardInodeKey]uint8, len(set))
			for k, a := range set {
				inter[k] = a
			}
			continue
		}
		for k, a := range inter {
			other, ok := set[k]
			switch {
			case !ok:
				delete(inter, k)
			case other == uint8(GUARD_ALLOW_ROOT) || a == uint8(GUARD_ALLOW_ROOT):
				inter[k] = uint8(GUARD_ALLOW_ROOT)
			}
		}
	}
	return inter
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

// replaceInodeU32Map makes m hold exactly want.
func replaceInodeU32Map(m *cilium.Map, want map[GuardInodeKey]uint32) error {
	var stale []GuardInodeKey
	var k GuardInodeKey
	var v uint32
	it := m.Iterate()
	for it.Next(&k, &v) {
		if _, keep := want[k]; !keep {
			stale = append(stale, k)
		}
	}
	if err := it.Err(); err != nil {
		return err
	}
	for i := range stale {
		if err := m.Delete(stale[i]); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
			return err
		}
	}
	for key, val := range want {
		if err := m.Put(key, val); err != nil {
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
