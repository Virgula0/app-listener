//go:build gui

package gui

import (
	"os"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/protected"
)

// sysStatus is what the host offers the workspaces: guards need eBPF AND the BPF LSM (a hook that
// attaches without bpf in the LSM list denies nothing), monitors only eBPF.
type sysStatus struct {
	root    bool
	ebpfErr error
	lsmErr  error
	daemon  bool
}

func probeSystem() sysStatus {
	return sysStatus{
		root:    os.Geteuid() == 0,
		ebpfErr: ebpf.Check(),
		lsmErr:  ebpf.CheckBPFLSM(),
		daemon:  protected.DaemonRunning(),
	}
}

func (s sysStatus) canMonitor() bool { return s.root && s.ebpfErr == nil }
func (s sysStatus) canGuard() bool   { return s.canMonitor() && s.lsmErr == nil }
