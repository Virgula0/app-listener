package workbench

import (
	"os"

	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/monitor"
	"github.com/Virgula0/app-listener/internal/networkguard"
	"github.com/Virgula0/app-listener/internal/networkmonitor"
	"github.com/Virgula0/app-listener/internal/repository"
)

// RealFileEngines builds the eBPF monitor and guard, as `monitor` and `guard` do.
func RealFileEngines() FileEngines {
	return FileEngines{
		NewMonitor: func(t ebpf.Target, recursive bool, depth int) (repository.MonitorRepository, error) {
			return monitor.NewMonitor([]ebpf.Target{t}, recursive, depth)
		},
		NewGuard: func(path string, mode guard.Mode, bins []ebpf.BinaryEntry, recursive bool, depth int) (FileGuard, error) {
			// Eager: the tree is claimed before the hooks attach, as `guard` does (a post-attach
			// walk would be denied by the whitelist itself).
			return guard.NewGuard(path, mode, bins, recursive, depth, guard.WithEagerPopulate())
		},
		Hash: ebpf.ComputeBinaryEntry,
	}
}

// RealNetEngines builds the eBPF network monitor and guard, as `network-monitor` and
// `network-guard` do.
func RealNetEngines() NetEngines {
	return NetEngines{
		NewMonitor: func(bins []ebpf.BinaryEntry) (repository.NetworkMonitorRepository, error) {
			return networkmonitor.NewNetworkMonitor(bins)
		},
		NewGuard: func(cfg NetGuardConfig, prev repository.NetworkGuardRepository) (repository.NetworkGuardRepository, error) {
			var opts []networkguard.Option
			if p, ok := prev.(*networkguard.NetGuard); ok {
				opts = append(opts, networkguard.WithSuspectsFrom(p))
			}
			return networkguard.NewNetGuard(cfg.Mode, cfg.Binaries, cfg.Events, cfg.Unsafe, cfg.Throttle, opts...)
		},
		Infra: networkguard.DiscoverInfraBinaries,
		Hash:  ebpf.ComputeBinaryEntry,
		Self:  os.Executable,
	}
}
