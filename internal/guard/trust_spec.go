package guard

import (
	"errors"
	"fmt"

	cilium "github.com/cilium/ebpf"
	"github.com/cilium/ebpf/asm"
	"github.com/cilium/ebpf/btf"
)

const memfdTarget = "memfd_alloc_file"

// trustSpec returns the trust object's spec and whether trust_memfd_alloc can attach. An fexit's
// target is resolved at load, so on a kernel without memfd_alloc_file the whole object would fail
// over that one best-effort program: it is swapped for an inert one, never attached.
func trustSpec() (*cilium.CollectionSpec, bool, error) {
	spec, err := LoadGuardTrust()
	if err != nil {
		return nil, false, fmt.Errorf("reading embedded trust objects: %w", err)
	}
	if KernelHasFunc(memfdTarget) {
		return spec, true, nil
	}
	spec.Programs[GuardTrustProgTrustMemfdAlloc] = InertProgram(GuardTrustProgTrustMemfdAlloc)
	return spec, false, nil
}

// InertProgram stands in for a program whose attach target this kernel lacks; never attached.
func InertProgram(name string) *cilium.ProgramSpec {
	return &cilium.ProgramSpec{
		Name:         name,
		Type:         cilium.SocketFilter,
		License:      "GPL",
		Instructions: asm.Instructions{asm.Mov.Imm(asm.R0, 0), asm.Return()},
	}
}

// KernelHasFunc reports whether vmlinux's BTF has function name. Unreadable BTF answers true: the
// load then reports the real error.
func KernelHasFunc(name string) bool {
	spec, err := btf.LoadKernelSpec()
	if err != nil {
		return true
	}
	var fn *btf.Func
	err = spec.TypeByName(name, &fn)
	return err == nil || !errors.Is(err, btf.ErrNotFound)
}
