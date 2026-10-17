package networkguard

import (
	"errors"
	"fmt"
)

// Option configures NewNetGuard.
type Option func(*NetGuard)

// WithSuspectsFrom starts the guard with prev's code-suspect marks. A guard that replaces a running
// one must not clear marks /proc no longer shows (a library mapped, then unmapped). prev must still
// be running: a mark that can't be handed over refuses the start.
func WithSuspectsFrom(prev *NetGuard) Option {
	return func(g *NetGuard) { g.inheritFrom = prev }
}

// inheritSuspects copies prev's marks. It runs after this guard's lifecycle hooks attach, so a
// copied process that exits is dropped by task_free here, and a leftover for a reused tgid by fork.
func (g *NetGuard) inheritSuspects() error {
	prev := g.inheritFrom
	if prev == nil {
		return nil
	}
	prev.mu.Lock()
	defer prev.mu.Unlock()
	if prev.stop {
		return errors.New("the guard to take code-suspect marks from has stopped")
	}
	var (
		tgid uint32
		why  uint8
	)
	it := prev.objs.TrustCodeSuspect.Iterate()
	for it.Next(&tgid, &why) {
		if err := g.objs.TrustCodeSuspect.Put(tgid, why); err != nil {
			return fmt.Errorf("carrying the code-suspect mark of tgid %d: %w", tgid, err)
		}
	}
	if err := it.Err(); err != nil {
		return fmt.Errorf("reading the previous guard's code-suspect marks: %w", err)
	}
	return nil
}
