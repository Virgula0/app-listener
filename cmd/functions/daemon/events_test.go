package daemon

import (
	"slices"
	"testing"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// Pins daemonSelfBaselineEvents to exactly {Open, Read, Stat}: its "never write" guarantee for
// every permanent config/ephemeral guard's self grant. Widening this var would widen every
// construction site (buildGuards included) at once, so a change must be deliberate.
func TestDaemonSelfBaselineEventsNeverWrites(t *testing.T) {
	want := []ebpf.EventType{ebpf.EventOpen, ebpf.EventRead, ebpf.EventStat}
	if !slices.Equal(daemonSelfBaselineEvents, want) {
		t.Errorf("daemonSelfBaselineEvents = %v, want %v (never write — see its doc comment)",
			daemonSelfBaselineEvents, want)
	}
	for _, forbidden := range []ebpf.EventType{ebpf.EventWrite, ebpf.EventDelete, ebpf.EventRename, ebpf.EventMknod, ebpf.EventMkdir} {
		if slices.Contains(daemonSelfBaselineEvents, forbidden) {
			t.Errorf("daemonSelfBaselineEvents must never contain %s", forbidden)
		}
	}
}

// The ephemeral guard unlockPendingGroupRoots puts over a grouped root stays live through
// buildGuards, which must pre-create the file-vault recovery sidecar (an O_CREAT gated by
// EVENT_MKNOD alone; guard_path_mknod). groupUnlockSelfEvents must be exactly baseline +
// EventMknod: no less (the create is denied, giving the "could not pre-create the recovery sidecar"
// warning) and no more (never write/delete/rename).
func TestGroupUnlockSelfEventsIsBaselinePlusMknod(t *testing.T) {
	if !slices.Contains(groupUnlockSelfEvents, ebpf.EventMknod) {
		t.Fatalf("groupUnlockSelfEvents = %v must contain EventMknod, or the ephemeral "+
			"grouped-root guard denies its own file-vault sidecar pre-create", groupUnlockSelfEvents)
	}
	for _, e := range daemonSelfBaselineEvents {
		if !slices.Contains(groupUnlockSelfEvents, e) {
			t.Errorf("groupUnlockSelfEvents = %v is missing baseline event %s", groupUnlockSelfEvents, e)
		}
	}
	want := append(slices.Clone(daemonSelfBaselineEvents), ebpf.EventMknod)
	if !slices.Equal(groupUnlockSelfEvents, want) {
		t.Errorf("groupUnlockSelfEvents = %v, want exactly daemonSelfBaselineEvents+EventMknod = %v",
			groupUnlockSelfEvents, want)
	}
}
