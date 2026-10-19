package ebpf

import (
	"fmt"
	_ "unsafe" // go:linkname

	_ "github.com/cilium/ebpf" // initializes internal/linux before init() below overrides it
	"golang.org/x/sys/unix"
)

// ciliumVersion mirrors cilium/ebpf's internal.Version; TestCiliumKernelVersionOverride catches
// a layout change on upgrade.
type ciliumVersion [3]uint16

// ciliumKernelVersion is cilium/ebpf's cached kernel-version lookup, used for every Kprobe program
// load, its own feature probes included (link.Kprobe's bpf_link_perf_event one). It reads the
// vDSO through /proc/self/mem, which the guard's PROC_MEM gate denies to a tainted process: every
// app-listener process once the daemon runs (it execs a whitelisted image). uname replaces it.
//
//go:linkname ciliumKernelVersion github.com/cilium/ebpf/internal/linux.KernelVersion
var ciliumKernelVersion func() (ciliumVersion, error)

func init() {
	ciliumKernelVersion = unameKernelVersion
}

func unameKernelVersion() (ciliumVersion, error) {
	var uts unix.Utsname
	if err := unix.Uname(&uts); err != nil {
		return ciliumVersion{}, fmt.Errorf("uname: %w", err)
	}
	return parseKernelRelease(unix.ByteSliceToString(uts.Release[:]))
}

// parseKernelRelease reads "7.0.11-arch1-1" as 7.0.11, each part clamped to 255 like
// KERNEL_VERSION's SUBLEVEL.
func parseKernelRelease(release string) (ciliumVersion, error) {
	var v [3]uint16
	if n, _ := fmt.Sscanf(release, "%d.%d.%d", &v[0], &v[1], &v[2]); n < 2 {
		return ciliumVersion{}, fmt.Errorf("unparsable kernel release %q", release)
	}
	return ciliumVersion{min(v[0], 255), min(v[1], 255), min(v[2], 255)}, nil
}
