// Package pulse shares one ticker between many periodic loops, so N guards wake the process once
// per period instead of N times at drifting phases.
package pulse

import (
	"sync"
	"sync/atomic"
	"time"
)

// Pulse fires every period from its first subscription until the process exits.
type Pulse struct {
	period time.Duration
	once   sync.Once
	gen    atomic.Uint64
	mu     sync.Mutex
	ch     chan struct{} // closed and replaced on every tick
}

var closed = func() chan struct{} {
	c := make(chan struct{})
	close(c)
	return c
}()

// New returns a Pulse ticking every period (> 0).
func New(period time.Duration) *Pulse {
	return &Pulse{period: period, ch: make(chan struct{})}
}

func (p *Pulse) run() {
	t := time.NewTicker(p.period)
	defer t.Stop()
	for range t.C {
		p.mu.Lock()
		p.gen.Add(1)
		close(p.ch)
		p.ch = make(chan struct{})
		p.mu.Unlock()
	}
}

// Subscribe returns a subscriber whose first tick is the Pulse's next one.
func (p *Pulse) Subscribe() *Sub {
	p.once.Do(func() { go p.run() })
	return &Sub{p: p, seen: p.gen.Load()}
}

// Sub is one loop's view of a Pulse, with time.Ticker's delivery rule: a tick that fires while
// the loop is busy is delivered late, not lost, and several missed ticks collapse into one. A
// periodic security sweep must never silently skip a period.
type Sub struct {
	p    *Pulse
	seen uint64
}

// C is ready when a tick happened since the last Ack. Call it on every select, then Ack on
// receive. Not safe for concurrent use: one Sub per loop.
func (s *Sub) C() <-chan struct{} {
	if s.p.gen.Load() != s.seen {
		return closed
	}
	s.p.mu.Lock()
	defer s.p.mu.Unlock()
	if s.p.gen.Load() != s.seen {
		return closed
	}
	return s.p.ch
}

// Ack consumes every tick so far.
func (s *Sub) Ack() {
	s.seen = s.p.gen.Load()
}
