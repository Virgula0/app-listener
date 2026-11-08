package guard

import (
	"errors"
	"fmt"
	"os"
	"strings"

	cilium "github.com/cilium/ebpf"
	log "github.com/sirupsen/logrus"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// A handover keeps the WithHandOver resources (runtime sections) enforced across a clean stop: their
// slots, rows and res-id pins stay (Guard.HandOver), every other resource is retired as usual, and
// the links stay pinned. Unlike a SIGKILL's pins, a handover's guard only those resources, so the next start can
// run under them: AdmitSelfInHandover lets its own exe read them (an update changes its inode), and
// the start retires them once its own guards hold the trees. The marker pin is written last: its
// presence means no other resource is left enforced.
const handoverPinName = "handover"

// HandOver makes the next Stop keep a WithHandOver guard's resource enforced (see handoverPinName).
// Only for the final shutdown: a slot retained while the process goes on would leak, its rows
// judging whatever resource reuses the id. Without pinning nothing is kept.
func (g *Guard) HandOver() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.handingOver = g.handOver
}

// retainSlot keeps slot id enforced past its Guard's Stop, reporting whether it did. The kernel
// state stays as is; only event delivery ends.
func (e *engine) retainSlot(id uint32) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	if !e.started || e.pinPrefix == "" || e.pinDegraded || !e.mapsPinned ||
		id >= GuardMaxRes || e.slots[id] == nil {
		return false
	}
	e.slots[id] = nil
	e.retained = true
	e.refs--
	if e.refs <= 0 {
		e.stopLocked()
	}
	return true
}

// handOverLocked readies the last stop's kernel state for a handover: only retained slots judge
// anything once it returns true. False = the caller unpins everything as on any other stop.
func (e *engine) handOverLocked() bool {
	if !e.retained || e.pinPrefix == "" || e.pinDegraded || !e.mapsPinned {
		return false
	}
	// The raw block-device gate protects the retired trees' filesystems, not a runtime tmpfs.
	if err := deleteFsDevices(e.objs.GuardFsDevices); err != nil {
		log.Errorf("guard: handover: clearing the raw block-device gate: %v — not handing over", err)
		return false
	}
	m, err := cilium.NewMap(&cilium.MapSpec{Type: cilium.Array, KeySize: 4, ValueSize: 4, MaxEntries: 1})
	if err != nil {
		log.Errorf("guard: handover: creating the marker: %v — not handing over", err)
		return false
	}
	defer m.Close()
	if err := m.Pin(e.pinPrefix + handoverPinName); err != nil {
		log.Errorf("guard: handover: pinning the marker: %v — not handing over", err)
		return false
	}
	log.Infof("guard engine: runtime sections stay guarded by the pins at %s* until the next start", e.pinPrefix)
	return true
}

func deleteFsDevices(m *cilium.Map) error {
	var keys []uint32
	var k uint32
	var v uint8
	it := m.Iterate()
	for it.Next(&k, &v) {
		keys = append(keys, k)
	}
	if err := it.Err(); err != nil {
		return err
	}
	for _, key := range keys {
		if err := m.Delete(key); err != nil && !errors.Is(err, cilium.ErrKeyNotExist) {
			return err
		}
	}
	return nil
}

// handoverGen returns the generation of the one handover pinned under base, "" if none. Two
// handovers can't both be the last stop's, so that is an error.
func handoverGen(base string) (string, error) {
	if base == "" {
		return "", nil
	}
	entries, err := os.ReadDir(base)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", err
	}
	var gens []string
	for _, e := range entries {
		gen := genOfPinFile(e.Name())
		if gen != "" && e.Name() == sharedPinName(gen)+handoverPinName {
			gens = append(gens, gen)
		}
	}
	switch len(gens) {
	case 0:
		return "", nil
	case 1:
		return gens[0], nil
	default:
		return "", fmt.Errorf("%d handovers pinned under %s (%v)", len(gens), base, gens)
	}
}

// AdmitSelfInHandover finds the handover pinned under base and admits this process's exe on each of
// its resources, root-gated with events, as its own guards will (WithSelfAllowBinary). Returns the
// handover's generation, "" if there is none. The key is always /proc/self/exe's, never a parameter.
// On error the caller must not rely on the handover: RetireHandover drops it.
func AdmitSelfInHandover(base string, events []ebpf.EventType) (string, error) {
	gen, err := handoverGen(base)
	if err != nil || gen == "" {
		return gen, err
	}
	dev, ino, err := ebpf.StatInode("/proc/self/exe")
	if err != nil {
		return gen, fmt.Errorf("resolving own executable: %w", err)
	}
	mask, err := eventMask(events)
	if err != nil {
		return gen, err
	}
	shared := SharedPinPrefix(base, gen)
	actions, err := cilium.LoadPinnedMap(shared+ExeActionsPinName, nil)
	if err != nil {
		return gen, fmt.Errorf("loading the handover's %s: %w", ExeActionsPinName, err)
	}
	defer actions.Close()
	masks, err := cilium.LoadPinnedMap(shared+ExeEventsPinName, nil)
	if err != nil {
		return gen, fmt.Errorf("loading the handover's %s: %w", ExeEventsPinName, err)
	}
	defer masks.Close()

	prefixes, err := handoverResources(base, gen)
	if err != nil {
		return gen, err
	}
	for _, prefix := range prefixes {
		resID, idErr := readPinnedResID(prefix)
		if idErr != nil {
			return gen, idErr
		}
		key := GuardResInodeKey{ResId: resID, Ino: GuardInodeKey{Dev: dev, Ino: ino}}
		// The mask first: a row without one would allow every event.
		if putErr := masks.Put(key, mask); putErr != nil {
			return gen, fmt.Errorf("admitting this binary on handed-over slot %d: %w", resID, putErr)
		}
		if putErr := actions.Put(key, uint8(GUARD_ALLOW_ROOT)); putErr != nil {
			return gen, fmt.Errorf("admitting this binary on handed-over slot %d: %w", resID, putErr)
		}
	}
	log.Infof("guard: %d runtime section(s) still guarded by the stopped run's pins (generation %s)", len(prefixes), gen)
	return gen, nil
}

// handoverResources returns the PinPrefix of every resource gen left a res-id pin for.
func handoverResources(base, gen string) ([]string, error) {
	entries, err := os.ReadDir(base)
	if err != nil {
		return nil, err
	}
	shared := sharedPinName(gen)
	var out []string
	for _, e := range entries {
		name := e.Name()
		if genOfPinFile(name) != gen || strings.HasPrefix(name, shared) || !strings.HasSuffix(name, ResIDPinName) {
			continue
		}
		out = append(out, base+"/"+strings.TrimSuffix(name, ResIDPinName))
	}
	return out, nil
}

// RetireHandover unpins the handover under base, if any: its resources are unguarded from then on.
func RetireHandover(base string) error {
	gen, err := handoverGen(base)
	if err != nil {
		// Several: none can be told apart as the last, so all go.
		_, err = retirePins(base, func(g string) bool { return isHandoverGen(base, g) })
		return err
	}
	if gen == "" {
		return nil
	}
	_, err = retirePins(base, func(g string) bool { return g == gen })
	return err
}

func isHandoverGen(base, gen string) bool {
	_, err := os.Lstat(SharedPinPrefix(base, gen) + handoverPinName)
	return err == nil
}
