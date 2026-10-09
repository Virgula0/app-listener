package networkguard

import (
	"testing"

	cilium "github.com/cilium/ebpf"
)

func TestNetSpecInertsIntegrityOnlyForBlacklist(t *testing.T) {
	black, memfd, err := netSpec(false)
	if err != nil {
		t.Fatal(err)
	}
	if memfd {
		t.Fatal("a blacklist must never attach the memfd provenance program")
	}
	white, _, err := netSpec(true)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range integrityProgs {
		if black.Programs[name].Type != cilium.SocketFilter {
			t.Errorf("blacklist: %s must be inert (it needs a newer kernel and is never attached)", name)
		}
		if name != GuardNetProgNetgMemfdAlloc && white.Programs[name].Type == cilium.SocketFilter {
			t.Errorf("whitelist: %s was made inert", name)
		}
	}
	if black.Programs[GuardNetProgNetgInodeFree].Type == cilium.SocketFilter {
		t.Error("blacklist: a lifecycle program was made inert")
	}
}
