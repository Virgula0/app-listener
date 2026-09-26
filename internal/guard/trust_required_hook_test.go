package guard

import (
	"os"
	"testing"
)

// TrustGuard.Start attaches five LSM programs best-effort: each failure is logged and skipped, and
// Start only errors when ZERO attached. mmap_file is not merely one of five — it IS protection #2,
// the library allowlist that stops LD_PRELOAD of attacker code into a whitelisted process. When it
// alone fails to attach, Start returns nil and the daemon logs
// "enforcing write-protection + library allowlist", so the operator is told the protection is on
// while it is entirely absent.
//
// The failure is not hypothetical: each guarded resource stacks another program on the same
// attach point, and the per-attach-point trampoline cap (BPF_MAX_TRAMP_LINKS, 38) is reached on
// real installs — an exhausted hook returns E2BIG ("argument list too long") from AttachLSM. The
// per-resource guards treat their own required hooks as fatal (see requiredHooks in guard.go); the
// trust guard must do the same for mmap_file rather than degrade silently.
func TestTrustGuardStartFailsWhenLibraryAllowlistHookCannotAttach(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip("requires root to load and attach BPF programs")
	}

	tg, err := NewTrustGuard()
	if err != nil {
		t.Skipf("cannot load trust BPF objects on this kernel: %v", err)
	}
	defer tg.Stop()

	// Stand in for an exhausted attach point: closing the program makes its AttachLSM fail exactly
	// as E2BIG does, while every other hook still attaches.
	if cerr := tg.objs.TrustMmap.Close(); cerr != nil {
		t.Fatalf("closing the mmap_file program: %v", cerr)
	}

	if serr := tg.Start(); serr == nil {
		t.Fatal("Start() reported success with the library allowlist (mmap_file) unattached: " +
			"LD_PRELOAD into whitelisted binaries is unprotected while the daemon logs it as enforcing")
	}
}
