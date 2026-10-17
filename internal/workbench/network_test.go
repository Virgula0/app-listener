package workbench

import (
	"errors"
	"fmt"
	"slices"
	"testing"
	"time"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/networkguard"
	"github.com/Virgula0/app-listener/internal/repository"
)

type fakeNetMonitor struct {
	j      *journal
	name   string
	events chan ebpf.NetEvent
	types  []ebpf.NetEventType
}

func (m *fakeNetMonitor) Start() error                        { m.j.add(m.name + " start"); return nil }
func (m *fakeNetMonitor) Stop()                               { m.j.add(m.name + " stop") }
func (m *fakeNetMonitor) Events() <-chan ebpf.NetEvent        { return m.events }
func (m *fakeNetMonitor) SetEventTypes(t []ebpf.NetEventType) { m.types = t }

type fakeNetGuard struct {
	j      *journal
	name   string
	cfg    NetGuardConfig
	prev   repository.NetworkGuardRepository
	events chan networkguard.NetGuardEvent
	err    error
}

func (g *fakeNetGuard) Start() error {
	g.j.add(g.name + " start")
	return g.err
}
func (g *fakeNetGuard) Stop()                                     { g.j.add(g.name + " stop") }
func (g *fakeNetGuard) Events() <-chan networkguard.NetGuardEvent { return g.events }

type fakeNet struct {
	j        *journal
	mons     []*fakeNetMonitor
	guards   []*fakeNetGuard
	hashes   map[string][32]byte
	infra    []string
	startErr error
}

func newFakeNet(t *testing.T) (*Network, *fakeNet) {
	t.Helper()
	f := &fakeNet{j: &journal{}, hashes: map[string][32]byte{"/a": {1}, "/b": {2}, "/infra": {3}, "/self": {4}}}
	n := NewNetwork(NetEngines{
		NewMonitor: func(bins []ebpf.BinaryEntry) (repository.NetworkMonitorRepository, error) {
			m := &fakeNetMonitor{j: f.j, name: fmt.Sprintf("mon%d(%d)", len(f.mons)+1, len(bins)),
				events: make(chan ebpf.NetEvent, 4)}
			f.mons = append(f.mons, m)
			return m, nil
		},
		NewGuard: func(cfg NetGuardConfig, prev repository.NetworkGuardRepository) (repository.NetworkGuardRepository, error) {
			g := &fakeNetGuard{j: f.j, name: fmt.Sprintf("guard%d", len(f.guards)+1), cfg: cfg, prev: prev,
				events: make(chan networkguard.NetGuardEvent), err: f.startErr}
			f.guards = append(f.guards, g)
			return g, nil
		},
		Infra: func() ([]string, error) { return f.infra, nil },
		Hash: func(p string) (ebpf.BinaryEntry, error) {
			h, ok := f.hashes[p]
			if !ok {
				return ebpf.BinaryEntry{}, errors.New("missing " + p)
			}
			return ebpf.BinaryEntry{Path: p, Hash: h}, nil
		},
		Self: func() (string, error) { return "/self", nil },
	})
	t.Cleanup(n.Close)
	return n, f
}

var (
	netA = ebpf.BinaryEntry{Path: "/a", Hash: [32]byte{1}}
	netB = ebpf.BinaryEntry{Path: "/b", Hash: [32]byte{2}}
)

func TestNetApplySwapsNewBeforeOld(t *testing.T) {
	n, f := newFakeNet(t)
	if err := n.Apply(NetPolicy{Mode: GuardWhitelist, Binaries: []ebpf.BinaryEntry{netA}}); err != nil {
		t.Fatal(err)
	}
	if err := n.Apply(NetPolicy{Mode: GuardBlacklist, Binaries: []ebpf.BinaryEntry{netB}}); err != nil {
		t.Fatal(err)
	}
	want := []string{"guard1 start", "guard2 start", "guard1 stop"}
	if got := f.j.entries(); !slices.Equal(got, want) {
		t.Fatalf("calls = %q, want %q", got, want)
	}
	if f.guards[1].prev != f.guards[0] {
		t.Fatal("the new guard must be handed the running one (code-suspect marks)")
	}
	if f.guards[1].cfg.Mode != networkguard.ModeBlacklist {
		t.Fatal("mode not passed")
	}
	if err := n.Apply(NetPolicy{Mode: GuardOff}); err != nil {
		t.Fatal(err)
	}
	if got := f.j.entries(); got[len(got)-1] != "guard2 stop" || n.Policy().Mode != GuardOff {
		t.Fatalf("off must stop the guard: %q", got)
	}
}

func TestNetApplyFailureKeepsRunningGuard(t *testing.T) {
	n, f := newFakeNet(t)
	if err := n.Apply(NetPolicy{Mode: GuardWhitelist, Binaries: []ebpf.BinaryEntry{netA}}); err != nil {
		t.Fatal(err)
	}
	f.startErr = errors.New("attach failed")
	if err := n.Apply(NetPolicy{Mode: GuardWhitelist}); err == nil {
		t.Fatal("start failure must surface")
	}
	if got := f.j.entries(); slices.Contains(got, "guard1 stop") {
		t.Fatalf("the running guard must survive a failed replacement: %q", got)
	}
	if p := n.Policy(); len(p.Binaries) != 1 || n.LastError() == "" {
		t.Fatalf("policy %+v, lastErr %q", p, n.LastError())
	}
}

func TestNetApplyRefusesReplacedBinary(t *testing.T) {
	n, f := newFakeNet(t)
	f.hashes["/a"] = [32]byte{9}
	if err := n.Apply(NetPolicy{Mode: GuardWhitelist, Binaries: []ebpf.BinaryEntry{netA}}); err == nil {
		t.Fatal("a binary whose content changed since it was chosen must be refused")
	}
	if len(f.guards) != 0 {
		t.Fatal("no guard may be built from a changed binary")
	}
}

func TestNetApplyFlagRules(t *testing.T) {
	n, _ := newFakeNet(t)
	for _, p := range []NetPolicy{{Mode: GuardBlacklist, Unsafe: true}, {Mode: GuardBlacklist, AutoInfra: true}} {
		if err := n.Apply(p); err == nil {
			t.Fatalf("%+v must be refused, as network-guard does", p)
		}
	}
}

func TestNetApplyExtras(t *testing.T) {
	n, f := newFakeNet(t)
	f.infra = []string{"/infra", "/a"}
	p := NetPolicy{Mode: GuardWhitelist, Binaries: []ebpf.BinaryEntry{netA}, AutoInfra: true, Unsafe: true}
	if err := n.Apply(p); err != nil {
		t.Fatal(err)
	}
	var paths []string
	for _, b := range f.guards[0].cfg.Binaries {
		paths = append(paths, b.Path)
	}
	if want := []string{"/a", "/infra", "/self"}; !slices.Equal(paths, want) {
		t.Fatalf("guard binaries = %q, want %q (listed once each)", paths, want)
	}
	extras := n.Extras()
	if len(extras) != 2 || extras[0].Why != WhyInfra || extras[1].Why != WhySelf {
		t.Fatalf("extras = %+v", extras)
	}
	if !f.guards[0].cfg.Unsafe {
		t.Fatal("unsafe not passed")
	}
}

func TestNetWatchRestartsMonitor(t *testing.T) {
	n, f := newFakeNet(t)
	n.SetMonitorEvents([]ebpf.NetEventType{ebpf.NetDNS})
	if err := n.Watch([]ebpf.BinaryEntry{netA}); err != nil {
		t.Fatal(err)
	}
	if err := n.Watch([]ebpf.BinaryEntry{netA, netB}); err != nil {
		t.Fatal(err)
	}
	if err := n.Watch(nil); err != nil {
		t.Fatal(err)
	}
	want := []string{"mon1(1) start", "mon2(2) start", "mon1(1) stop", "mon2(2) stop"}
	if got := f.j.entries(); !slices.Equal(got, want) {
		t.Fatalf("calls = %q, want %q", got, want)
	}
	if !slices.Equal(f.mons[1].types, []ebpf.NetEventType{ebpf.NetDNS}) {
		t.Fatal("a new monitor must keep the event filter")
	}
	if n.Monitoring() || len(n.Watched()) != 0 {
		t.Fatal("an empty set stops monitoring")
	}
}

func TestNetCloseStopsGuardThenMonitor(t *testing.T) {
	n, f := newFakeNet(t)
	if err := n.Watch([]ebpf.BinaryEntry{netA}); err != nil {
		t.Fatal(err)
	}
	if err := n.Apply(NetPolicy{Mode: GuardBlacklist, Binaries: []ebpf.BinaryEntry{netB}}); err != nil {
		t.Fatal(err)
	}
	n.Close()
	got := f.j.entries()
	if want := []string{"guard1 stop", "mon1(1) stop"}; !slices.Equal(got[len(got)-2:], want) {
		t.Fatalf("calls = %q", got)
	}
	if err := n.Apply(NetPolicy{Mode: GuardBlacklist}); err == nil {
		t.Fatal("apply after close must fail")
	}
}

func TestNetGuardEventsBecomeRows(t *testing.T) {
	n, f := newFakeNet(t)
	if err := n.Apply(NetPolicy{Mode: GuardBlacklist, Binaries: []ebpf.BinaryEntry{netB}}); err != nil {
		t.Fatal(err)
	}
	f.guards[0].events <- networkguard.NetGuardEvent{
		NetEvent: ebpf.NetEvent{PID: 4242, Comm: "curl", Type: ebpf.NetConnect, Protocol: 6, DstAddr: "1.2.3.4:443"},
		Blocked:  true,
	}
	for range 100 {
		if rows, _ := n.Rows(Filter{}); len(rows) == 1 {
			r := rows[0]
			if r.Verdict != Denied || r.Remote != "1.2.3.4:443" || r.Op != "CONNECT" || r.Proto != "TCP" {
				t.Fatalf("row = %+v", r)
			}
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("event never ingested")
}
