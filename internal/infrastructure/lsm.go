package ebpf

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	lsmPath = "/sys/kernel/security/lsm"

	// assumeBpfLsmEnv bypasses the BPF LSM preflight check entirely. It is
	// only meant for environments where securityfs is not mounted inside
	// the process's mount namespace (privileged containers used by the
	// integration suite) while the host kernel itself is verified.
	assumeBpfLsmEnv = "APPLISTENER_ASSUME_BPF_LSM"
)

// bpfLsmListed reports whether the kernel's active LSM stack includes the
// BPF LSM. The list is read from securityfs; entries are comma-separated.
func bpfLsmListed(file string) (bool, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return false, fmt.Errorf("cannot read %s (is securityfs mounted?): %w", file, err)
	}
	for _, lsm := range strings.Split(strings.TrimSpace(string(data)), ",") {
		if strings.TrimSpace(lsm) == "bpf" {
			return true, nil
		}
	}
	return false, nil
}

// CheckBPFLSM verifies that the kernel's active LSM stack includes the BPF
// LSM before any LSM-hook-based mode (guard, network guard, daemon) starts.
//
// A successful hook attach is NOT proof of enforcement: on kernels where
// bpf is absent from the LSM list the hooks attach but never fire, leaving
// the program silently inert. Refusing to start is safer than a guard that
// denies nothing.
func CheckBPFLSM() error {
	return CheckBPFLSMAt(lsmPath)
}

// CheckBPFLSMAt is CheckBPFLSM with an injectable lsm file path (for tests).
func CheckBPFLSMAt(file string) error {
	if os.Getenv(assumeBpfLsmEnv) == "1" {
		// The escape hatch is intentional (privileged containers cannot
		// read securityfs), but it silently disables the enforcement
		// preflight — make the degraded check visible in the logs.
		log.Warnf("%s=1: skipping the BPF LSM availability preflight — guards may silently not enforce if the host kernel lacks BPF LSM", assumeBpfLsmEnv)
		return nil
	}

	listed, err := bpfLsmListed(file)
	if err != nil {
		return err
	}
	if listed {
		return nil
	}

	return fmt.Errorf(
		"BPF LSM is not active: %s does not contain 'bpf'. "+
			"The guard relies on BPF LSM hooks (file_open/file_permission), which attach but never fire "+
			"when bpf is missing from the LSM stack. "+
			"Boot a kernel with CONFIG_BPF_LSM=y and 'bpf' in CONFIG_LSM (for example on the kernel "+
			"command line: lsm=landlock,lockdown,yama,integrity,apparmor,bpf) and reboot. "+
			"GitHub-hosted runners and many cloud kernels ship without bpf in CONFIG_LSM and cannot "+
			"run these hooks. Set %s (e.g. in containers without securityfs) only when the host "+
			"kernel does have the BPF LSM", file, assumeBpfLsmEnv)
}
