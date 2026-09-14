package daemon

import (
	"slices"
	"testing"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// TestDaemonSelfBaselineEventsNeverWrites pins daemonSelfBaselineEvents to
// exactly {Open, Read, Stat} — the "never write" guarantee its own doc
// comment makes for every permanent config/ephemeral guard's self grant.
// Silently widening this package var would widen every construction site at
// once (buildGuards' per-resource guards included), so a change here must be
// deliberate, not a side effect of editing something else.
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

// TestGroupUnlockSelfEventsIsBaselinePlusMknod: the ephemeral guard
// unlockPendingGroupRoots attaches over a grouped encryption root stays live
// through buildGuards, which must be able to pre-create the file-vault
// recovery sidecar (an O_CREAT, gated by EVENT_MKNOD alone — see
// guard_path_mknod) for any regular-file WatchRelPaths sub-resource under
// that root. groupUnlockSelfEvents must therefore be exactly
// daemonSelfBaselineEvents plus EventMknod: no less (the sidecar create is
// denied again, reproducing the "could not pre-create the recovery sidecar"
// warning) and no more (this mask must stay as narrow as the baseline it
// widens, never drifting into write/delete/rename territory).
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
