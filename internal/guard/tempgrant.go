package guard

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"sync"

	cilium "github.com/cilium/ebpf"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// Temporary grants back `edit-protected --forward`: an authenticated root session adds whitelist
// rows (-w) only where the resource has none, or narrows an existing GUARD_ALLOW (-b), for its
// lifetime. Revoke restores a row only while it still holds exactly what the grant wrote and the
// permanent whitelist hasn't claimed the key since, so a concurrent resync, system-pattern admission
// or tamper demotion always wins over the restore.

// TempBinary is a binary resolved for a temporary grant. Its inode stays open for the session, so
// the number can't be freed and reused by another file while the rows name it.
type TempBinary struct {
	Path string
	// Key is what the grant's rows name (ExeKey: a uutils applet's own key); Key.Real() the file.
	Key GuardInodeKey
	// SystemPlaced: only root could have placed it (ebpf.SystemPlacedInode). Otherwise its owner
	// can rewrite it in place during the window; Changed catches that.
	SystemPlaced bool
	// Runtime: the launch-injectable runtimes the image embeds (runtimeClass), set on the inode for
	// the session so trust_exec_launch judges its launches as it does a whitelisted one's.
	Runtime uint8

	// mu: the grant's watcher runs Changed while a revoke may Close.
	mu   sync.Mutex
	f    *os.File
	hash [32]byte
	stat ebpf.BinaryStat
}

// ResolveTempBinary opens path confined (ebpf.OpenConfined), requires an executable regular file,
// and pins the inode reached and its content hash.
func ResolveTempBinary(path string) (*TempBinary, error) {
	f, err := ebpf.OpenConfined(path)
	if err != nil {
		return nil, err
	}
	b := &TempBinary{Path: path, f: f}
	var st unix.Stat_t
	if serr := unix.Fstat(int(f.Fd()), &st); serr != nil {
		b.Close()
		return nil, &os.PathError{Op: "fstat", Path: path, Err: serr}
	}
	if st.Mode&unix.S_IFMT != unix.S_IFREG || st.Mode&0o111 == 0 {
		b.Close()
		return nil, fmt.Errorf("%s is not an executable regular file", path)
	}
	b.stat = binaryStatOf(&st)
	file := GuardInodeKey{Dev: b.stat.Dev, Ino: b.stat.Ino}
	entry, err := ebpf.ComputeBinaryEntryFile(f, path)
	if err != nil {
		b.Close()
		return nil, err
	}
	b.hash = entry.Hash
	b.SystemPlaced = ebpf.SystemPlacedInode(path, file.Dev, file.Ino)
	b.Runtime = runtimeClass(f, path)
	if b.Key, err = ExeKey(path, f, file); err != nil {
		b.Close()
		return nil, err
	}
	return b, nil
}

// RuntimeNames names the runtimes in b.Runtime, for a prompt.
func (b *TempBinary) RuntimeNames() []string {
	var out []string
	for _, c := range []struct {
		bit  uint8
		name string
	}{{trustedNode, "Node"}, {trustedElectron, "Electron"}, {trustedChromium, "Chromium"}, {trustedJVM, "JVM"}} {
		if b.Runtime&c.bit != 0 {
			out = append(out, c.name)
		}
	}
	return out
}

func binaryStatOf(st *unix.Stat_t) ebpf.BinaryStat {
	return ebpf.BinaryStat{
		Dev:     ebpf.KernelDev(unix.Major(st.Dev), unix.Minor(st.Dev)),
		Ino:     st.Ino,
		Size:    st.Size,
		MtimeNs: st.Mtim.Sec*1_000_000_000 + st.Mtim.Nsec,
		CtimeNs: st.Ctim.Sec*1_000_000_000 + st.Ctim.Nsec,
	}
}

// Changed reports whether the pinned inode's content changed since ResolveTempBinary (an in-place
// rewrite: the inode, and so the grant, stays the same). An unreadable inode counts as changed.
func (b *TempBinary) Changed() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.f == nil {
		return false // released with its grant
	}
	var st unix.Stat_t
	if err := unix.Fstat(int(b.f.Fd()), &st); err != nil {
		return true
	}
	fp := binaryStatOf(&st)
	if fp == b.stat {
		return false
	}
	entry, err := ebpf.ComputeBinaryEntryFile(b.f, b.Path)
	if err != nil || entry.Hash != b.hash {
		return true
	}
	b.stat = fp
	return false
}

// Close releases the pinned inode.
func (b *TempBinary) Close() {
	if b == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.f != nil {
		_ = b.f.Close()
		b.f = nil
	}
}

// TempRule is what a temporary grant does to every listed binary: Allow admits it (-w), otherwise
// it is denied (-b). Events limits the rule to those types; empty means every event.
type TempRule struct {
	Allow  bool
	Events []ebpf.EventType
}

// TemporaryGrant is one resource's planned temporary rows (PlanTemporaryGrant).
type TemporaryGrant interface {
	// JournalRows are the allow rows Apply creates. Record them before Apply: after a crash
	// `daemon --lockdown` strips them from the pinned maps (StripPinnedTempAllows).
	JournalRows() []GuardResInodeKey
	// Report describes binaries the grant leaves untouched, one line each.
	Report() []string
	// Apply writes the rows, failing (and undoing its own writes) if any row moved since planning.
	Apply() error
	// Revoke restores what Apply changed; idempotent.
	Revoke() error
}

// exeRow is one key's whitelist state on a resource: its guard_exe_actions and guard_exe_events
// rows (absent mask = every event).
type exeRow struct {
	hasAction bool
	action    uint8
	hasMask   bool
	mask      uint32
}

type tempRow struct {
	bin   *TempBinary
	prior exeRow
	// What Apply writes: newAction 0 leaves the action row alone; maskOp says what happens to the
	// event-mask row.
	newAction uint8
	maskOp    tempMaskOp
	newMask   uint32
	applied   bool
	// The guard_trusted_files row as trustTempRuntime found it, and the bits it ORed in.
	trust      *TrustGuard
	trustHad   bool
	trustPrior uint8
	trustAdded uint8
}

type tempMaskOp int

const (
	maskKeep tempMaskOp = iota
	maskPut
	maskDelete
)

type tempGrant struct {
	g      *Guard
	rows   []*tempRow
	report []string

	mu      sync.Mutex
	revoked bool
}

// allEventBits is the base a -b event list is cut from when the binary has no mask row (= every
// event, including types the BPF side reports beyond ebpf.EventTypes).
const allEventBits = ^uint32(0)

// PlanTemporaryGrant validates rule against this resource's current rows and returns the grant
// without writing anything. Refused: a non-whitelist guard, the daemon's own binary, an inspector,
// and -w on a blocked (tamper-demoted) binary. Binaries already in the state asked for are left
// alone and reported.
func (g *Guard) PlanTemporaryGrant(bins []*TempBinary, rule TempRule) (TemporaryGrant, error) {
	if g.mode != ModeWhitelist {
		return nil, fmt.Errorf("guard %s: temporary grants need a whitelist guard", g.path)
	}
	g.mu.Lock()
	stopped, selfKey, selfSet := g.stopped, g.selfKey, g.selfKeySet
	g.mu.Unlock()
	if stopped {
		return nil, fmt.Errorf("guard %s is stopped", g.path)
	}

	t := &tempGrant{g: g}
	for _, b := range bins {
		if selfSet && b.Key == selfKey {
			return nil, fmt.Errorf("%s is the app-listener binary itself — never granted temporarily", b.Path)
		}
		if sharedEngine.isInspector(b.Key) {
			return nil, fmt.Errorf("%s is an inspector — its access is not a whitelist row", b.Path)
		}
		prior, err := g.readExeRow(b.Key)
		if err != nil {
			return nil, err
		}
		row := &tempRow{bin: b, prior: prior}
		var change bool
		if rule.Allow {
			change, err = g.planTempAllow(row, rule.Events)
		} else {
			change, err = planTempBlock(row, rule.Events)
		}
		if err != nil {
			return nil, err
		}
		if !change {
			t.report = append(t.report, g.tempSkipReason(b, rule))
			continue
		}
		t.rows = append(t.rows, row)
	}
	return t, nil
}

func (g *Guard) tempSkipReason(b *TempBinary, rule TempRule) string {
	if rule.Allow {
		return fmt.Sprintf("%s: %s is already whitelisted — left unchanged", g.path, b.Path)
	}
	return fmt.Sprintf("%s: %s is not whitelisted — already denied", g.path, b.Path)
}

// planTempAllow fills row for -w; false when the binary is whitelisted already.
func (g *Guard) planTempAllow(row *tempRow, events []ebpf.EventType) (bool, error) {
	switch {
	case row.prior.hasAction && row.prior.action == GUARD_ALLOW:
		return false, nil
	case row.prior.hasAction:
		return false, fmt.Errorf("%s is blocked on %s (action %d, e.g. demoted after an in-place change) — "+
			"refusing to allow it", row.bin.Path, g.path, row.prior.action)
	}
	row.newAction = GUARD_ALLOW
	if len(events) > 0 {
		mask, err := eventMask(events)
		if err != nil {
			return false, err
		}
		row.newMask, row.maskOp = mask, maskPut
	} else if row.prior.hasMask {
		row.maskOp = maskDelete
	}
	return true, nil
}

// planTempBlock fills row for -b; false when the binary isn't whitelisted (already denied).
func planTempBlock(row *tempRow, events []ebpf.EventType) (bool, error) {
	if !row.prior.hasAction || row.prior.action != GUARD_ALLOW {
		return false, nil
	}
	if len(events) == 0 {
		row.newAction = GUARD_BLOCK
		return true, nil
	}
	// Exact bits: eventMask's READ->OPEN|STAT implication would also strip what wasn't asked for.
	var deny uint32
	for _, t := range events {
		if t < 0 || t >= 32 {
			return false, fmt.Errorf("event type %d out of range", t)
		}
		deny |= 1 << uint(t) //nolint:gosec // range-checked above
	}
	base := allEventBits
	if row.prior.hasMask {
		base = row.prior.mask
	}
	row.newMask, row.maskOp = base&^deny, maskPut
	return true, nil
}

func (g *Guard) readExeRow(k GuardInodeKey) (row exeRow, err error) {
	err = g.objs().GuardExeActions.Lookup(g.resKey(k), &row.action)
	switch {
	case err == nil:
		row.hasAction = true
	case !errors.Is(err, cilium.ErrKeyNotExist):
		return row, fmt.Errorf("guard %s: reading whitelist row: %w", g.path, err)
	}
	err = g.objs().GuardExeEvents.Lookup(g.resKey(k), &row.mask)
	switch {
	case err == nil:
		row.hasMask = true
	case !errors.Is(err, cilium.ErrKeyNotExist):
		return row, fmt.Errorf("guard %s: reading event mask: %w", g.path, err)
	}
	return row, nil
}

func (t *tempGrant) JournalRows() []GuardResInodeKey {
	out := make([]GuardResInodeKey, 0, len(t.rows))
	for _, r := range t.rows {
		if r.newAction == GUARD_ALLOW {
			out = append(out, t.g.resKey(r.bin.Key))
		}
	}
	return out
}

func (t *tempGrant) Report() []string { return t.report }

// Apply runs under g.mu so Stop can't retire (and a reload reuse) the slot between the check and
// the writes. The trust guard is looked up first: PruneSuperseded takes supersede.mu, then g.mu.
func (t *tempGrant) Apply() error {
	g := t.g
	trust := liveTrust()
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.stopped {
		return fmt.Errorf("guard %s is stopped", g.path)
	}
	for _, r := range t.rows {
		if err := g.applyTempRow(r, trust); err != nil {
			if rerr := t.revokeLocked(trust); rerr != nil {
				log.Errorf("guard %s: undoing a partly applied temporary grant: %v", g.path, rerr)
			}
			return err
		}
	}
	return nil
}

func (g *Guard) applyTempRow(r *tempRow, trust *TrustGuard) error {
	cur, err := g.readExeRow(r.bin.Key)
	if err != nil {
		return err
	}
	if cur != r.prior {
		return fmt.Errorf("%s: the whitelist row of %s changed while the grant was prepared — retry", g.path, r.bin.Path)
	}
	r.applied = true
	if err := g.writeTempRow(r, trust); err != nil {
		return fmt.Errorf("%s: temporary row of %s: %w", g.path, r.bin.Path, err)
	}
	return nil
}

// writeTempRow puts an allow's mask in first and a block's action first: no moment is wider than
// the final state.
func (g *Guard) writeTempRow(r *tempRow, trust *TrustGuard) error {
	k := r.bin.Key
	if r.newAction == GUARD_BLOCK {
		if err := g.putExeAction(k, GUARD_BLOCK); err != nil {
			return err
		}
	}
	switch r.maskOp {
	case maskPut:
		if err := g.objs().GuardExeEvents.Put(g.resKey(k), r.newMask); err != nil {
			return err
		}
	case maskDelete:
		if err := g.restoreMask(k, false, 0); err != nil {
			return err
		}
	}
	if r.newAction == GUARD_ALLOW {
		if err := r.trustTempRuntime(trust); err != nil {
			return err
		}
		return g.putExeAction(k, GUARD_ALLOW)
	}
	return nil
}

// liveTrust is the attached trust guard; nil before its Start and after its Stop.
func liveTrust() *TrustGuard {
	supersede.mu.Lock()
	defer supersede.mu.Unlock()
	return supersede.trust
}

// trustTempRuntime ORs the binary's runtime classes into guard_trusted_files before its allow row
// lands: trust_exec_launch judges only exes carrying them, so ELECTRON_RUN_AS_NODE, NODE_OPTIONS or
// --inspect would otherwise run the caller's code with the grant. Never TRUSTED_BINARY, and never
// over a row that has it: that is the permanent whitelist's. No trust guard, no allow (trust_mmap's
// code-suspect marking and the launch checks both live in it).
func (r *tempRow) trustTempRuntime(t *TrustGuard) error {
	if t == nil {
		return errors.New("the trust guard is not running — refusing to admit a binary without its launch checks")
	}
	if r.bin.Runtime == 0 {
		return nil
	}
	file := r.bin.Key.Real()
	k := GuardTrustInodeKey{Dev: file.Dev, Ino: file.Ino}
	var flags uint8
	err := t.objs.GuardTrustedFiles.Lookup(k, &flags)
	had := err == nil
	switch {
	case had:
		if flags&trustedBinary != 0 || flags&r.bin.Runtime == r.bin.Runtime {
			return nil
		}
		err = t.objs.GuardTrustedFiles.Update(k, flags|r.bin.Runtime, cilium.UpdateExist)
	case errors.Is(err, cilium.ErrKeyNotExist):
		err = t.objs.GuardTrustedFiles.Update(k, r.bin.Runtime, cilium.UpdateNoExist)
	}
	if err != nil {
		return fmt.Errorf("marking %s's runtime launches for judging: %w", r.bin.Path, err)
	}
	r.trust, r.trustHad, r.trustPrior, r.trustAdded = t, had, flags, r.bin.Runtime&^flags
	return nil
}

// untrustTempRuntime takes back the bits trustTempRuntime added, after the allow row is gone, and
// only while the row holds exactly what it wrote and no guard admits the key: a row SetTrusted,
// AllowReplacement or another admission rewrote since is theirs to keep (extra bits only judge more).
func (r *tempRow) untrustTempRuntime(live *TrustGuard) error {
	t, added := r.trust, r.trustAdded
	r.trust, r.trustAdded = nil, 0
	if added == 0 || t != live || sharedEngine.admitsExe(r.bin.Key.Real()) {
		return nil
	}
	file := r.bin.Key.Real()
	k := GuardTrustInodeKey{Dev: file.Dev, Ino: file.Ino}
	var cur uint8
	if err := t.objs.GuardTrustedFiles.Lookup(k, &cur); err != nil {
		if errors.Is(err, cilium.ErrKeyNotExist) {
			return nil
		}
		return err
	}
	if cur != r.trustPrior|added {
		return nil
	}
	if r.trustHad {
		return t.objs.GuardTrustedFiles.Update(k, r.trustPrior, cilium.UpdateExist)
	}
	if err := t.objs.GuardTrustedFiles.Delete(k); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
		return err
	}
	return nil
}

func (t *tempGrant) Revoke() error {
	trust := liveTrust() // before g.mu: see Apply
	t.g.mu.Lock()
	defer t.g.mu.Unlock()
	return t.revokeLocked(trust)
}

// revokeLocked runs under g.mu. A stopped guard's rows were dropped with its slot, which a later
// guard may own: none is written then. trust: the live trust guard (a stopped one's map is gone).
func (t *tempGrant) revokeLocked(trust *TrustGuard) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.revoked {
		return nil
	}
	t.revoked = true
	var errs []error
	for _, r := range slices.Backward(t.rows) {
		if !r.applied {
			continue
		}
		// A stopped guard's whitelist rows went with its slot; the trust row is daemon-wide.
		if !t.g.stopped {
			if err := t.g.restoreTempRowLocked(r); err != nil {
				errs = append(errs, fmt.Errorf("%s: restoring the whitelist row of %s: %w", t.g.path, r.bin.Path, err))
			}
		}
		if err := r.untrustTempRuntime(trust); err != nil {
			errs = append(errs, fmt.Errorf("%s: restoring the trust row of %s: %w", t.g.path, r.bin.Path, err))
		}
	}
	return errors.Join(errs...)
}

// restoreTempRowLocked undoes r only where the row still holds what Apply wrote: a key the
// permanent whitelist claimed meanwhile, or a mask someone else rewrote, is theirs to keep.
func (g *Guard) restoreTempRowLocked(r *tempRow) error {
	k := r.bin.Key
	cur, err := g.readExeRow(k)
	if err != nil {
		return err
	}
	permanent := g.deployedKeyLocked(k)
	if r.newAction == GUARD_ALLOW && permanent {
		return nil
	}
	if err := g.restoreTempAction(r, cur, permanent); err != nil {
		return err
	}
	if (r.maskOp == maskPut && cur.hasMask && cur.mask == r.newMask) || (r.maskOp == maskDelete && !cur.hasMask) {
		return g.restoreMask(k, r.prior.hasMask, r.prior.mask)
	}
	return nil
}

func (g *Guard) restoreTempAction(r *tempRow, cur exeRow, permanent bool) error {
	k := r.bin.Key
	switch {
	case r.newAction == GUARD_ALLOW && cur.hasAction && cur.action == GUARD_ALLOW:
		if err := g.objs().GuardExeActions.Delete(g.resKey(k)); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
			return err
		}
		return sharedEngine.noteDeny(g.resID, k)
	// Never re-allow a key the whitelist dropped or the tamper detector demoted.
	case r.newAction == GUARD_BLOCK && cur.hasAction && cur.action == GUARD_BLOCK && permanent &&
		!g.demotedKeyLocked(k):
		return g.putExeAction(k, GUARD_ALLOW)
	}
	return nil
}

func (g *Guard) restoreMask(k GuardInodeKey, had bool, mask uint32) error {
	if had {
		return g.objs().GuardExeEvents.Put(g.resKey(k), mask)
	}
	if err := g.objs().GuardExeEvents.Delete(g.resKey(k)); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
		return err
	}
	return nil
}

func (g *Guard) deployedKeyLocked(k GuardInodeKey) bool {
	for _, key := range g.deployed {
		if key == k {
			return true
		}
	}
	return false
}

func (g *Guard) demotedKeyLocked(k GuardInodeKey) bool {
	for _, st := range g.binaryVerifyStates {
		if st.key == k && st.demoted {
			return true
		}
	}
	return false
}

// StripPinnedTempAllows deletes, from a dead daemon's pinned shared maps, every journaled row that
// is still GUARD_ALLOW, with its event mask. Strip-only: a forged journal can deny, never admit.
func StripPinnedTempAllows(sharedPrefix string, rows []GuardResInodeKey) (int, error) {
	if sharedPrefix == "" || len(rows) == 0 {
		return 0, nil
	}
	actions, events, err := loadPinnedExeMaps(sharedPrefix)
	if err != nil || actions == nil {
		return 0, err
	}
	defer actions.Close()
	if events != nil {
		defer events.Close()
	}

	stripped := 0
	var errs []error
	for _, k := range rows {
		ok, err := stripPinnedAllow(actions, events, k)
		if ok {
			stripped++
		}
		if err != nil {
			errs = append(errs, err)
		}
	}
	if stripped > 0 {
		log.Warnf("guard: stripped %d temporary allow row(s) a dead daemon left in its pinned whitelist", stripped)
	}
	return stripped, errors.Join(errs...)
}

// loadPinnedExeMaps returns nil maps when a clean shutdown already unpinned them.
func loadPinnedExeMaps(sharedPrefix string) (actions, events *cilium.Map, err error) {
	actions, err = cilium.LoadPinnedMap(sharedPrefix+ExeActionsPinName, nil)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, fmt.Errorf("loading pinned %s: %w", ExeActionsPinName, err)
	}
	events, err = cilium.LoadPinnedMap(sharedPrefix+ExeEventsPinName, nil)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		actions.Close()
		return nil, nil, fmt.Errorf("loading pinned %s: %w", ExeEventsPinName, err)
	}
	return actions, events, nil
}

func stripPinnedAllow(actions, events *cilium.Map, k GuardResInodeKey) (bool, error) {
	var action uint8
	if err := actions.Lookup(k, &action); err != nil {
		if errors.Is(err, cilium.ErrKeyNotExist) {
			return false, nil
		}
		return false, err
	}
	if action != GUARD_ALLOW {
		return false, nil
	}
	if err := actions.Delete(k); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
		return false, err
	}
	if events != nil {
		if err := events.Delete(k); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
			return true, err
		}
	}
	return true, nil
}
