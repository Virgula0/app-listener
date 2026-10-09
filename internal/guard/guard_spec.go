package guard

import (
	"fmt"

	cilium "github.com/cilium/ebpf"
)

// ioctlCompatTarget is the LSM hook a compat (32-bit) ioctl calls since 6.8; before that it called
// file_ioctl, so the guard's file_ioctl program covers it there.
const ioctlCompatTarget = "bpf_lsm_file_ioctl_compat"

// guardSpec returns the guard object's spec and whether guard_file_ioctl_compat can attach. An LSM
// program's target is resolved at load, so on a pre-6.8 kernel it is swapped for an inert one.
func guardSpec() (*cilium.CollectionSpec, bool, error) {
	spec, err := LoadGuard()
	if err != nil {
		return nil, false, fmt.Errorf("reading embedded guard objects: %w", err)
	}
	if KernelHasFunc(ioctlCompatTarget) {
		return spec, true, nil
	}
	spec.Programs[GuardProgGuardFileIoctlCompat] = InertProgram(GuardProgGuardFileIoctlCompat)
	return spec, false, nil
}
