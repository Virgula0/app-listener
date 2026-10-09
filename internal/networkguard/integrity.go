package networkguard

import (
	"fmt"
	"time"

	cilium "github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/guard"
)

const memfdTarget = "memfd_alloc_file"

// integrityProgs only whitelist mode attaches. They need bpf_loop (5.17) and lsm/sb_delete (5.11):
// a blacklist swaps them for inert programs, so it still loads on the 5.10 BPF-LSM floor.
var integrityProgs = []string{
	GuardNetProgNetgMmapFile, GuardNetProgNetgFileMprotect, GuardNetProgNetgFileOpen,
	GuardNetProgNetgExecEnv, GuardNetProgNetgPtraceAccessCheck, GuardNetProgNetgBprmCheck,
	GuardNetProgNetgSbDelete, GuardNetProgNetgMemfdAlloc,
}

// netSpec returns the object's spec and whether netg_memfd_alloc can attach: an fexit's target is
// resolved at load, so on a kernel without memfd_alloc_file it is swapped for an inert program.
func netSpec(whitelist bool) (*cilium.CollectionSpec, bool, error) {
	spec, err := LoadGuardNet()
	if err != nil {
		return nil, false, fmt.Errorf("reading embedded network guard objects: %w", err)
	}
	if !whitelist {
		for _, name := range integrityProgs {
			spec.Programs[name] = guard.InertProgram(name)
		}
		return spec, false, nil
	}
	if guard.KernelHasFunc(memfdTarget) {
		return spec, true, nil
	}
	spec.Programs[GuardNetProgNetgMemfdAlloc] = guard.InertProgram(GuardNetProgNetgMemfdAlloc)
	return spec, false, nil
}

type hook struct {
	prog *cilium.Program
	name string
	lsm  bool
}

func (g *NetGuard) attach(h hook) (link.Link, error) {
	if h.lsm {
		return link.AttachLSM(link.LSMOptions{Program: h.prog})
	}
	return link.AttachTracing(link.TracingOptions{Program: h.prog})
}

// attachRequired attaches every hook or fails; what says what a missing one would leave open.
func (g *NetGuard) attachRequired(hooks []hook, what string) error {
	for _, h := range hooks {
		l, err := g.attach(h)
		if err != nil {
			return fmt.Errorf("required network guard hook %s failed to attach: %w — %s", h.name, err, what)
		}
		g.links = append(g.links, l)
	}
	return nil
}

// lifecycleHooks keep a key's identity over time (exe_supersede.h) and the per-process exec stamps
// and code-suspect marks: without unlink/rename/inode_free a reused inode number inherits a
// whitelist row, without fork/exec/exit a stamp or mark is lost or outlives its process.
func (g *NetGuard) lifecycleHooks() []hook {
	return []hook{
		{g.objs.NetgSchedFork, "sched_process_fork", false},
		{g.objs.NetgTaskFree, "task_free", true},
		{g.objs.NetgBprmCommitted, "bprm_committed_creds", true},
		{g.objs.NetgInodeUnlink, "inode_unlink", true},
		{g.objs.NetgInodeRename, "inode_rename", true},
		{g.objs.NetgInodeFree, "inode_free_security", true},
	}
}

// integrityHooks judge the code a protected process runs and guard its memory (networkguard.bpf.c,
// "Code integrity"). All required: each one missing is an injection path.
func (g *NetGuard) integrityHooks() []hook {
	return []hook{
		{g.objs.NetgMmapFile, "mmap_file", true},
		{g.objs.NetgFileMprotect, "file_mprotect", true},
		{g.objs.NetgFileOpen, "file_open", true},
		{g.objs.NetgExecEnv, "sched_process_exec", false},
		{g.objs.NetgPtraceAccessCheck, "ptrace_access_check", true},
		{g.objs.NetgBprmCheck, "bprm_check_security", true},
	}
}

// attachIntegrity attaches sb_delete and vouches the host's mounts before the code-integrity hooks,
// so no library is judged against an empty or stale vouched set. The memfd provenance program is
// best-effort: without it a memfd is never vouched for (stricter).
func (g *NetGuard) attachIntegrity() error {
	if err := g.attachRequired([]hook{{g.objs.NetgSbDelete, "sb_delete", true}},
		"a reused device number could inherit a dead superblock's system-library trust"); err != nil {
		return err
	}
	g.mounts = guard.NewMountVouch("network guard", "network guard", g.objs.NetgVouchedDevs,
		g.objs.NetgDeadDevs, g.objs.NetgDevSeq)
	if err := g.mounts.Open(); err != nil {
		return fmt.Errorf("vouching system-library mounts: %w", err)
	}
	if err := g.attachRequired(g.integrityHooks(),
		"code injected into a whitelisted process (LD_PRELOAD, ptrace) would get its network access"); err != nil {
		return err
	}
	if g.memfd {
		if l, err := g.attach(hook{g.objs.NetgMemfdAlloc, memfdTarget, false}); err == nil {
			g.links = append(g.links, l)
		} else {
			log.Warnf("network guard: memfd JIT provenance unavailable (%v) — a whitelisted process "+
				"mapping its own memfd code loses its network access", err)
		}
	}
	return g.mounts.StartWatch()
}

// net_guard_event.reason values — must match networkguard.bpf.c.
const (
	reasonCode       uint32 = 1
	reasonEnv        uint32 = 2
	reasonSuperseded uint32 = 3
)

// reasonText renders why a whitelisted binary's process was refused.
func reasonText(r uint32) string {
	switch r {
	case 0:
		return ""
	case reasonCode:
		return "code-suspect"
	case reasonEnv:
		return "ld-preload"
	case reasonSuperseded:
		return "superseded"
	default:
		return fmt.Sprintf("reason-%d", r)
	}
}

const degradedEvery = 30 * time.Second

// watchDegraded reports kernel bookkeeping the network guard could not record until done closes.
func (g *NetGuard) watchDegraded(done <-chan struct{}) {
	t := time.NewTicker(degradedEvery)
	defer t.Stop()
	for {
		select {
		case <-done:
			return
		case <-t.C:
			g.warnDegraded()
		}
	}
}

func (g *NetGuard) warnDegraded() {
	var lost uint64
	if g.objs.NetgSuspectLost.Lookup(uint32(0), &lost) == nil && lost > g.lostSeen[0] {
		g.lostSeen[0] = lost
		log.Errorf("network guard: CRITICAL: %d code-suspect mark(s) lost (more suspect processes than "+
			"the mark map holds) — no whitelisted binary is admitted until the network guard restarts", lost)
	}
	if g.objs.ExeSeq.Lookup(uint32(1), &lost) == nil && lost > g.lostSeen[1] {
		g.lostSeen[1] = lost
		log.Warnf("network guard: degraded: %d exec stamp(s) lost (process table larger than the stamp "+
			"map) — processes started before the guard lose access through binaries replaced since", lost)
	}
}
