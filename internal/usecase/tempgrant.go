package usecase

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/logging"
	"github.com/Virgula0/app-listener/internal/repository"
)

// tempBinaryCheckEvery is how often a temporary grant re-checks its binaries for an in-place
// rewrite (fstat only; a hash only when the fingerprint moved).
var tempBinaryCheckEvery = 2 * time.Second

// tempEventBuffer bounds the events a slow session client may lag behind; the rest are dropped
// (display only: enforcement already happened in the kernel).
const tempEventBuffer = 512

// TemporaryAccess is a live GrantTemporaryAccess grant.
type TemporaryAccess struct {
	// Report lists binaries the grant left untouched, one line each.
	Report []string
	// Events carries the granted resources' events whose process runs a granted binary, allowed
	// ones included. Never closed: stop reading at Done.
	Events <-chan guard.GuardEvent
	// Done is closed once the grant is revoked: by Revoke, or by the daemon when a granted binary
	// changed in place.
	Done   <-chan struct{}
	revoke func() error
}

// Revoke restores every row the grant changed; idempotent.
func (a *TemporaryAccess) Revoke() error { return a.revoke() }

func (d *daemonUseCase) GrantTemporaryAccess(resourcePaths, binaryPaths []string, rule guard.TempRule,
	journal func([]guard.GuardResInodeKey) error) (*TemporaryAccess, error) {
	if len(resourcePaths) == 0 || len(binaryPaths) == 0 {
		return nil, errors.New("a temporary grant needs at least one resource and one binary")
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.stopping {
		return nil, errors.New("daemon: shutdown in progress")
	}
	if d.editGrantActive {
		return nil, errors.New("another live edit session is already active")
	}
	guards, err := d.guardsForLocked(resourcePaths)
	if err != nil {
		return nil, err
	}
	bins, err := resolveTempBinaries(binaryPaths)
	if err != nil {
		return nil, err
	}
	grants, report, err := planTempGrants(guards, bins, rule)
	if err != nil {
		closeTempBinaries(bins)
		return nil, err
	}
	rows := tempJournalRows(grants)
	if err := applyTempGrants(grants, rows, journal); err != nil {
		closeTempBinaries(bins)
		return nil, err
	}
	d.editGrantActive = true
	logTempGrant(resourcePaths, bins, rule)
	events := make(chan guard.GuardEvent, tempEventBuffer)
	tap := newExeTap(bins, events)
	for _, g := range guards {
		g.SetEventTap(tap)
	}

	done := make(chan struct{})
	stop := make(chan struct{})
	var once sync.Once
	var rerr error
	revoke := func() error {
		once.Do(func() {
			close(stop)
			for _, g := range guards {
				g.SetEventTap(nil)
			}
			rerr = revokeTempGrants(grants)
			// A failed revoke keeps the journal: `daemon --lockdown` strips what is left.
			if rerr == nil && len(rows) > 0 {
				rerr = journal(nil)
			}
			closeTempBinaries(bins)
			d.mu.Lock()
			d.editGrantActive = false
			d.mu.Unlock()
			if rerr != nil {
				log.Errorf("daemon: revoking the temporary access grant: %v", rerr)
			} else {
				log.Infof("daemon: temporary access REVOKED on %s", strings.Join(sanitized(resourcePaths), ", "))
			}
			close(done)
		})
		return rerr
	}
	go watchTempBinaries(bins, tempBinaryCheckEvery, stop, revoke)
	return &TemporaryAccess{Report: report, Events: events, Done: done, revoke: revoke}, nil
}

func (d *daemonUseCase) guardsForLocked(resourcePaths []string) ([]repository.GuardRepository, error) {
	seen := make(map[string]bool, len(resourcePaths))
	var out []repository.GuardRepository
	for _, p := range resourcePaths {
		if seen[p] {
			continue
		}
		seen[p] = true
		idx := -1
		for i := range d.resources {
			if d.resources[i].Path == p {
				idx = i
				break
			}
		}
		if idx < 0 {
			return nil, fmt.Errorf("%s is not a guarded resource in the running configuration", p)
		}
		out = append(out, d.guards[idx])
	}
	return out, nil
}

// resolveTempBinaries resolves every path; two naming one inode collapse into one.
func resolveTempBinaries(paths []string) ([]*guard.TempBinary, error) {
	var bins []*guard.TempBinary
	seen := make(map[guard.GuardInodeKey]bool, len(paths))
	for _, p := range paths {
		b, err := guard.ResolveTempBinary(p)
		if err != nil {
			closeTempBinaries(bins)
			return nil, fmt.Errorf("resolving %s: %w", p, err)
		}
		if seen[b.Key] {
			b.Close()
			continue
		}
		seen[b.Key] = true
		bins = append(bins, b)
	}
	return bins, nil
}

func closeTempBinaries(bins []*guard.TempBinary) {
	for _, b := range bins {
		b.Close()
	}
}

func planTempGrants(guards []repository.GuardRepository, bins []*guard.TempBinary,
	rule guard.TempRule) ([]guard.TemporaryGrant, []string, error) {
	grants := make([]guard.TemporaryGrant, 0, len(guards))
	var report []string
	changes := 0
	for _, g := range guards {
		gr, err := g.PlanTemporaryGrant(bins, rule)
		if err != nil {
			return nil, nil, err
		}
		grants = append(grants, gr)
		report = append(report, gr.Report()...)
		changes += len(bins) - len(gr.Report())
	}
	if changes == 0 {
		return nil, nil, fmt.Errorf("nothing to change: %s", strings.Join(report, "; "))
	}
	return grants, report, nil
}

func tempJournalRows(grants []guard.TemporaryGrant) []guard.GuardResInodeKey {
	rows := make([]guard.GuardResInodeKey, 0, len(grants))
	for _, gr := range grants {
		rows = append(rows, gr.JournalRows()...)
	}
	return rows
}

// applyTempGrants journals the allow rows before the first write, then applies every grant,
// undoing all of them on any failure.
func applyTempGrants(grants []guard.TemporaryGrant, rows []guard.GuardResInodeKey,
	journal func([]guard.GuardResInodeKey) error) error {
	if len(rows) > 0 {
		if err := journal(rows); err != nil {
			return fmt.Errorf("recording the temporary grant before applying it: %w", err)
		}
	}
	for i, gr := range grants {
		if err := gr.Apply(); err != nil {
			rerr := revokeTempGrants(grants[:i])
			if rerr == nil && len(rows) > 0 {
				rerr = journal(nil)
			}
			if rerr != nil {
				log.Errorf("daemon: rolling back a partly applied temporary grant: %v", rerr)
			}
			return err
		}
	}
	return nil
}

func revokeTempGrants(grants []guard.TemporaryGrant) error {
	var errs []error
	for i := len(grants) - 1; i >= 0; i-- {
		if err := grants[i].Revoke(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// watchTempBinaries revokes the whole grant once a granted binary's content changes in place: its
// inode, and so its rows, would admit whatever was written there.
func watchTempBinaries(bins []*guard.TempBinary, every time.Duration, stop <-chan struct{}, revoke func() error) {
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		select {
		case <-stop:
			return
		case <-t.C:
			for _, b := range bins {
				if b.Changed() {
					log.Errorf("daemon: %s changed in place during its temporary grant — revoking the grant",
						logging.SanitizeText(b.Path))
					_ = revoke()
					return
				}
			}
		}
	}
}

// exeTapTTL bounds how long a pid's exe verdict is reused: an exec changes it under the same pid.
const exeTapTTL = time.Second

type exeVerdict struct {
	match bool
	at    time.Time
}

// newExeTap returns a guard event tap forwarding, without blocking, the events whose process
// runs one of bins. The pid's exe is read from /proc after the fact: this feeds a display, never a
// policy decision (those are the kernel's, by inode).
func newExeTap(bins []*guard.TempBinary, out chan<- guard.GuardEvent) func(guard.GuardEvent) {
	keys := make(map[guard.GuardInodeKey]bool, len(bins))
	for _, b := range bins {
		keys[b.Key] = true
	}
	var mu sync.Mutex
	seen := make(map[uint32]exeVerdict)
	return func(ev guard.GuardEvent) {
		now := time.Now()
		mu.Lock()
		v, ok := seen[ev.PID]
		if !ok || now.Sub(v.at) > exeTapTTL {
			dev, ino, err := ebpf.StatInode(fmt.Sprintf("/proc/%d/exe", ev.PID))
			v = exeVerdict{match: err == nil && keys[guard.GuardInodeKey{Dev: dev, Ino: ino}], at: now}
			if len(seen) > 4096 {
				clear(seen)
			}
			seen[ev.PID] = v
		}
		mu.Unlock()
		if !v.match {
			return
		}
		select {
		case out <- ev:
		default:
		}
	}
}

func logTempGrant(resourcePaths []string, bins []*guard.TempBinary, rule guard.TempRule) {
	verb := "deny"
	if rule.Allow {
		verb = "allow"
	}
	events := "ALL"
	if len(rule.Events) > 0 {
		names := make([]string, len(rule.Events))
		for i, e := range rule.Events {
			names[i] = e.String()
		}
		events = strings.Join(names, ",")
	}
	paths := make([]string, len(bins))
	for i, b := range bins {
		paths[i] = b.Path
	}
	log.Warnf("daemon: temporary access GRANTED: %s %s (events %s) on %s", verb,
		strings.Join(sanitized(paths), ", "), events, strings.Join(sanitized(resourcePaths), ", "))
}

func sanitized(in []string) []string {
	out := make([]string, len(in))
	for i, s := range in {
		out[i] = logging.SanitizeText(s)
	}
	return out
}
