package ebpf

import (
	"fmt"
	"os"
	"strings"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/asm"
	log "github.com/sirupsen/logrus"
)

const (
	unprivilegedBpfPath = "/proc/sys/kernel/unprivileged_bpf_disabled"
	bpfJitEnablePath    = "/proc/sys/net/core/bpf_jit_enable"
	bpfJitHardenPath    = "/proc/sys/net/core/bpf_jit_harden"
)

func Check() error {
	spec := &ebpf.ProgramSpec{
		Type: ebpf.SocketFilter,
		Instructions: asm.Instructions{
			asm.LoadImm(asm.R0, 0, asm.DWord),
			asm.Return(),
		},
		License: "GPL",
	}

	prog, err := ebpf.NewProgram(spec)
	if err == nil {
		prog.Close()
		return nil
	}

	diag := collectDiagnostics()
	return fmt.Errorf("eBPF program creation failed: %w%s", err, diag)
}

func collectDiagnostics() string {
	var hints []string

	if val, err := readSysctl(unprivilegedBpfPath); err == nil {
		if val == "1" {
			hints = append(hints,
				"kernel.unprivileged_bpf_disabled=1 (unprivileged users blocked)")
		}
	}

	if val, err := readSysctl(bpfJitEnablePath); err == nil {
		if val == "0" {
			hints = append(hints,
				"net.core.bpf_jit_enable=0 (JIT disabled, eBPF may not work)")
		}
	} else {
		log.Debugf("cannot read %s: %v", bpfJitEnablePath, err)
	}

	if val, err := readSysctl(bpfJitHardenPath); err == nil {
		if val == "2" {
			hints = append(hints,
				"net.core.bpf_jit_harden=2 (JIT hardening at maximum)")
		}
	}

	if len(hints) == 0 {
		return ""
	}

	hints = append(hints,
		"ensure the process has CAP_BPF and CAP_SYS_ADMIN",
		"or set kernel.unprivileged_bpf_disabled=0 via sysctl")

	return "\nsysctl hints:\n  - " + strings.Join(hints, "\n  - ")
}

func readSysctl(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("cannot read %s: %w", path, err)
	}
	return strings.TrimSpace(string(data)), nil
}
