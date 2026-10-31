package guard

import (
	"errors"
	"fmt"
	"slices"

	cilium "github.com/cilium/ebpf"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

// WithPort makes this guard a `[port tcp/N]` resource: its root is the port's anchor file, and
// guard_ports routes the socket hooks on that port to it (guard.bpf.c, loopback TCP services).
func WithPort(port uint16) GuardOption {
	return func(g *Guard) {
		g.port = uint32(port)
	}
}

// portOwner is a live guard anchored on a port and the anchor key it published.
type portOwner struct {
	g   *Guard
	key GuardInodeKey
}

// publishPort points guard_ports at this guard's anchor. A reload's replacement anchors on the same
// file, so the row it writes equals the one it overwrites. g.mu is held across the write (Stop's
// order, g.mu then e.mu): a re-anchor racing Stop would re-add a row no cleanup ever removes.
func (g *Guard) publishPort() error {
	if g.port == 0 {
		return nil
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.stopped {
		return nil
	}
	return sharedEngine.publishPort(g, g.rootKey)
}

func (e *engine) publishPort(g *Guard, key GuardInodeKey) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if err := e.objs.GuardPorts.Put(g.port, key); err != nil {
		return fmt.Errorf("guarding port %d: %w", g.port, err)
	}
	log.Debugf("guard: port tcp/%d routed to resource %d", g.port, g.resID)
	owners := slices.DeleteFunc(e.ports[g.port], func(o portOwner) bool { return o.g == g })
	if e.ports == nil {
		e.ports = make(map[uint32][]portOwner)
	}
	e.ports[g.port] = append(owners, portOwner{g: g, key: key})
	return nil
}

// unpublishPort drops g from its port. The row stays while another guard is anchored there (the
// guard a reload replaces stops after its replacement published), pointing at that one's anchor.
func (e *engine) unpublishPort(g *Guard) {
	if g.port == 0 {
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	owners := slices.DeleteFunc(e.ports[g.port], func(o portOwner) bool { return o.g == g })
	if e.objs.GuardPorts == nil {
		return // engine already torn down
	}
	if len(owners) > 0 {
		e.ports[g.port] = owners
		live := owners[len(owners)-1]
		if err := e.objs.GuardPorts.Put(g.port, live.key); err != nil {
			log.Errorf("guard: CRITICAL: re-pointing port %d at its live resource: %v", g.port, err)
		}
		log.Infof("guard: port tcp/%d handed from resource %d to %d", g.port, g.resID, live.g.resID)
		return
	}
	delete(e.ports, g.port)
	if err := e.objs.GuardPorts.Delete(g.port); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
		log.Warnf("guard: clearing port %d: %v", g.port, err)
		return
	}
	log.Infof("guard: port tcp/%d no longer guarded", g.port)
}

// setPortNetNSLocked limits port decisions to the daemon's network namespace. Unknown leaves 0,
// which judges every namespace: stricter, never weaker.
func (e *engine) setPortNetNSLocked() {
	var st unix.Stat_t
	if err := unix.Stat("/proc/self/ns/net", &st); err != nil {
		log.Warnf("guard: own network namespace unknown (%v): guarded ports are judged in every namespace", err)
		return
	}
	if err := e.objs.GuardPortNetns.Put(uint32(0), st.Ino); err != nil {
		log.Warnf("guard: recording the network namespace (%v): guarded ports are judged in every namespace", err)
	}
}
