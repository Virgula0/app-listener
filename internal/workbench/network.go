package workbench

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"sync"
	"time"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/networkguard"
	"github.com/Virgula0/app-listener/internal/repository"
)

// NetPolicy is the network guard's configuration: network-guard's -w/-b, --auto-infra, --unsafe,
// --no-throttle (inverted) and -e.
type NetPolicy struct {
	Mode      GuardMode
	Binaries  []ebpf.BinaryEntry
	AutoInfra bool
	Unsafe    bool
	Throttle  bool
	Events    []ebpf.NetEventType // reported; empty = all
}

// NetExtra is a binary the guard admits beyond the user's list, and why.
type NetExtra struct {
	Entry ebpf.BinaryEntry
	Why   string
}

const (
	WhyInfra = "network daemon (auto-allow)"
	WhySelf  = "this window (local sockets)"
)

// NetGuardConfig is what a network guard is built from.
type NetGuardConfig struct {
	Mode     networkguard.Mode
	Binaries []ebpf.BinaryEntry
	Events   []ebpf.NetEventType
	Unsafe   bool
	Throttle bool
}

// NetEngines builds the engines; tests swap in fakes.
type NetEngines struct {
	NewMonitor func(bins []ebpf.BinaryEntry) (repository.NetworkMonitorRepository, error)
	// NewGuard builds a guard that replaces prev (nil: none), taking over its code-suspect marks.
	NewGuard func(cfg NetGuardConfig, prev repository.NetworkGuardRepository) (repository.NetworkGuardRepository, error)
	Infra    func() ([]string, error)
	Hash     func(path string) (ebpf.BinaryEntry, error)
	Self     func() (string, error)
}

// Network is the Network workspace: one network monitor over the watched binaries and the host-wide
// network guard.
type Network struct {
	*activity
	engines NetEngines
	exes    *exeResolver

	mu        sync.Mutex
	watched   []ebpf.BinaryEntry
	monEvents []ebpf.NetEventType
	mon       repository.NetworkMonitorRepository
	monStop   chan struct{}
	grd       repository.NetworkGuardRepository
	grdStop   chan struct{}
	policy    NetPolicy
	extras    []NetExtra
	lastErr   string
	closed    bool

	// opMu serializes Watch, Apply and Close.
	opMu sync.Mutex
}

// NewNetwork returns the workspace with nothing watched and the guard off.
func NewNetwork(engines NetEngines) *Network {
	return &Network{
		activity: newActivity(uint32(os.Getpid())), //nolint:gosec // pid fits
		engines:  engines, exes: newExeResolver("/proc"),
		policy: NetPolicy{Throttle: true},
	}
}

// Watched returns the binaries the monitor follows.
func (n *Network) Watched() []ebpf.BinaryEntry {
	n.mu.Lock()
	defer n.mu.Unlock()
	return slices.Clone(n.watched)
}

// Monitoring reports whether a monitor runs.
func (n *Network) Monitoring() bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.mon != nil
}

// Watch makes bins the monitored set. The new monitor starts before the old one stops so no event
// is lost; an empty set stops monitoring.
func (n *Network) Watch(bins []ebpf.BinaryEntry) error {
	n.opMu.Lock()
	defer n.opMu.Unlock()
	if n.isClosed() {
		return errClosed
	}
	bins = slices.Clone(bins)
	var next repository.NetworkMonitorRepository
	var stop chan struct{}
	if len(bins) > 0 {
		m, err := n.engines.NewMonitor(bins)
		if err != nil {
			return fmt.Errorf("creating network monitor: %w", err)
		}
		n.mu.Lock()
		events := slices.Clone(n.monEvents)
		n.mu.Unlock()
		if len(events) > 0 {
			m.SetEventTypes(events)
		}
		if err := m.Start(); err != nil {
			m.Stop()
			return fmt.Errorf("starting network monitor: %w", err)
		}
		next, stop = m, make(chan struct{})
		go n.pumpMonitor(m.Events(), stop)
	}
	n.mu.Lock()
	prev, prevStop := n.mon, n.monStop
	n.mon, n.monStop, n.watched = next, stop, bins
	n.mu.Unlock()
	if prev != nil {
		close(prevStop)
		prev.Stop()
	}
	return nil
}

// MonitorEvents returns the monitor's event filter (empty = all).
func (n *Network) MonitorEvents() []ebpf.NetEventType {
	n.mu.Lock()
	defer n.mu.Unlock()
	return slices.Clone(n.monEvents)
}

// SetMonitorEvents changes the monitor's event filter live (network-monitor's -e).
func (n *Network) SetMonitorEvents(types []ebpf.NetEventType) {
	n.mu.Lock()
	n.monEvents = slices.Clone(types)
	mon := n.mon
	n.mu.Unlock()
	if mon != nil {
		if len(types) == 0 {
			types = ebpf.NetEventTypes()
		}
		mon.SetEventTypes(types)
	}
}

func (n *Network) pumpMonitor(events <-chan ebpf.NetEvent, stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			return
		case ev, ok := <-events:
			if !ok {
				return
			}
			row := n.netRow(&ev, "monitor", Observed)
			n.ingest(&row)
		}
	}
}

func (n *Network) pumpGuard(events <-chan networkguard.NetGuardEvent, stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			return
		case ev, ok := <-events:
			if !ok {
				return
			}
			v := Allowed
			if ev.Blocked {
				v = Denied
			}
			row := n.netRow(&ev.NetEvent, "guard", v)
			row.Detail = ev.Reason
			n.ingest(&row)
		}
	}
}

// netRow converts an event. The exe comes from /proc after the fact: display only, and promoting it
// re-checks the inode (EntryFor).
func (n *Network) netRow(ev *ebpf.NetEvent, source string, v Verdict) Row {
	return Row{
		Time: time.Unix(0, ev.Timestamp), Source: source, Op: ev.Type.String(), Comm: ev.Comm, PID: ev.PID,
		UID: ev.UID, Proto: ebpf.ProtocolString(ev.Protocol), Local: ev.SrcAddr, Remote: ev.DstAddr,
		Size: ev.Size, Verdict: v, Exe: n.exes.resolve(ev.PID),
	}
}

// Policy returns the policy the guard enforces now.
func (n *Network) Policy() NetPolicy {
	n.mu.Lock()
	defer n.mu.Unlock()
	return clonePolicy(n.policy)
}

// Extras returns what the enforced guard admits beyond Policy().Binaries.
func (n *Network) Extras() []NetExtra {
	n.mu.Lock()
	defer n.mu.Unlock()
	return slices.Clone(n.extras)
}

// LastError is the most recent guard failure, for the view.
func (n *Network) LastError() string {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.lastErr
}

func clonePolicy(p NetPolicy) NetPolicy {
	p.Binaries = slices.Clone(p.Binaries)
	p.Events = slices.Clone(p.Events)
	return p
}

var errClosed = errors.New("the workspace is closed")

// validate mirrors network-guard's flag checks.
func (p *NetPolicy) validate() error {
	if p.Mode != GuardWhitelist && (p.Unsafe || p.AutoInfra) {
		return errors.New("blocking local sockets and auto-allowing network daemons need whitelist mode")
	}
	return nil
}

// Apply makes p the enforced policy. A new guard always comes up before the old one stops: while
// both run, an operation needs both verdicts, so the swap never admits what either policy denies.
// On error the previous guard keeps running.
func (n *Network) Apply(p NetPolicy) error {
	n.opMu.Lock()
	defer n.opMu.Unlock()
	if n.isClosed() {
		return errClosed
	}
	p = clonePolicy(p)
	err := p.validate()
	if err == nil {
		err = n.apply(p)
	}
	n.mu.Lock()
	n.lastErr = ""
	if err != nil {
		n.lastErr = err.Error()
	}
	n.mu.Unlock()
	return err
}

func (n *Network) apply(p NetPolicy) error {
	if p.Mode == GuardOff {
		n.swapGuard(nil, p, nil)
		return nil
	}
	if err := n.verify(p.Binaries); err != nil {
		return err
	}
	extras, err := n.extrasFor(p)
	if err != nil {
		return err
	}
	bins := slices.Clone(p.Binaries)
	for _, x := range extras {
		bins = append(bins, x.Entry)
	}
	mode := networkguard.ModeBlacklist
	if p.Mode == GuardWhitelist {
		mode = networkguard.ModeWhitelist
	}
	n.mu.Lock()
	prev := n.grd
	n.mu.Unlock()
	g, err := n.engines.NewGuard(NetGuardConfig{
		Mode: mode, Binaries: bins, Events: p.Events, Unsafe: p.Unsafe, Throttle: p.Throttle,
	}, prev)
	if err != nil {
		return fmt.Errorf("creating network guard: %w", err)
	}
	if err := g.Start(); err != nil {
		g.Stop()
		return fmt.Errorf("starting network guard: %w", err)
	}
	n.swapGuard(g, p, extras)
	return nil
}

func (n *Network) verify(bins []ebpf.BinaryEntry) error { return verifyHashes(n.engines.Hash, bins) }

// verifyHashes re-hashes every listed binary: a guard keys on whatever inode the path names when
// it applies, so a file replaced since it was chosen must not be admitted in its place.
func verifyHashes(hash func(string) (ebpf.BinaryEntry, error), bins []ebpf.BinaryEntry) error {
	for _, b := range bins {
		now, err := hash(b.Path)
		if err != nil {
			return err
		}
		if now.Hash != b.Hash {
			return fmt.Errorf("%s changed since it was added: remove it and add it again to trust the new file", b.Path)
		}
	}
	return nil
}

// extrasFor resolves --auto-infra's daemons and, with --unsafe, this process's own binary: a
// whitelist that blocks AF_UNIX would otherwise cut the window off its display server.
func (n *Network) extrasFor(p NetPolicy) ([]NetExtra, error) {
	have := map[string]bool{}
	for _, b := range p.Binaries {
		have[b.Path] = true
	}
	var out []NetExtra
	add := func(path, why string) error {
		if have[path] {
			return nil
		}
		have[path] = true
		e, err := n.engines.Hash(path)
		if err != nil {
			return err
		}
		out = append(out, NetExtra{Entry: e, Why: why})
		return nil
	}
	if p.AutoInfra {
		paths, err := n.engines.Infra()
		if err != nil {
			return nil, fmt.Errorf("discovering network daemons: %w", err)
		}
		for _, path := range paths {
			if err := add(path, WhyInfra); err != nil {
				return nil, err
			}
		}
	}
	if p.Unsafe {
		self, err := n.engines.Self()
		if err != nil {
			return nil, fmt.Errorf("resolving this program's executable: %w", err)
		}
		if err := add(self, WhySelf); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// swapGuard installs next (nil = none) and only then stops the previous guard.
func (n *Network) swapGuard(next repository.NetworkGuardRepository, p NetPolicy, extras []NetExtra) {
	var stop chan struct{}
	if next != nil {
		stop = make(chan struct{})
		go n.pumpGuard(next.Events(), stop)
	}
	n.mu.Lock()
	prev, prevStop := n.grd, n.grdStop
	n.grd, n.grdStop, n.policy, n.extras = next, stop, p, extras
	n.mu.Unlock()
	if prev != nil {
		close(prevStop)
		prev.Stop()
	}
}

func (n *Network) isClosed() bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.closed
}

// Close stops the guard, then the monitor (window closed).
func (n *Network) Close() {
	n.opMu.Lock()
	defer n.opMu.Unlock()
	n.mu.Lock()
	if n.closed {
		n.mu.Unlock()
		return
	}
	n.closed = true
	mon, monStop, grd, grdStop := n.mon, n.monStop, n.grd, n.grdStop
	n.mon, n.grd = nil, nil
	n.mu.Unlock()
	if grd != nil {
		close(grdStop)
		grd.Stop()
	}
	if mon != nil {
		close(monStop)
		mon.Stop()
	}
}
