package install

import (
	"fmt"
	"os"
	"slices"
	"syscall"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	"github.com/Virgula0/app-listener/internal/systemd"
)

// resolveFilesystemPrereqs makes every encryption root's filesystem ready for fscrypt. For each
// fixable prerequisite it asks whether to run the command now (the installer is root) or abort
// (declining aborts the install). A terminal condition (unsupported filesystem or kernel) always
// aborts. Offerable commands live in fscrypt.(*Vault).FilesystemPrereqs: audit them there.
func resolveFilesystemPrereqs(vault *fscrypt.Vault, cfg *daemonconfig.Config) error {
	// Per-filesystem `fscrypt setup` needs the global /etc/fscrypt.conf: create it first (no-op if
	// present) so a first-ever run doesn't fail.
	if err := fscrypt.EnsureSystemSetup(); err != nil {
		return fmt.Errorf("initializing fscrypt (/etc/fscrypt.conf): %w", err)
	}

	ran := map[string]bool{}
	for {
		prereqs, err := collectFilesystemPrereqs(vault, cfg)
		if err != nil {
			return err // terminal — unsupported filesystem/kernel
		}
		if len(prereqs) == 0 {
			return nil
		}
		p := prereqs[0]
		if ran[p.Command()] {
			return fmt.Errorf("the filesystem is still not ready after running `%s` — "+
				"run it manually and check its output, then re-run the installer", p.Command())
		}

		proceed, err := confirmRunPrereq(p)
		if err != nil {
			return err
		}
		if !proceed {
			return fmt.Errorf("aborted: a prerequisite is not met — run `%s` as root, then re-run the installer (%s)",
				p.Command(), p.Reason)
		}

		log.Infof("running prerequisite command: %s", p.Command())
		if err := systemd.RunCmd(p.Argv[0], p.Argv[1:]...); err != nil {
			return fmt.Errorf("prerequisite command `%s` failed: %w", p.Command(), err)
		}
		ran[p.Command()] = true
		log.Infof("prerequisite satisfied: %s", p.Command())
	}
}

// collectFilesystemPrereqs returns every fixable prerequisite across the config's encryption roots
// (each backing filesystem once), or a terminal error. Only stat + read-only fscrypt probes, so
// safe to loop between remediation commands.
//
// Regular-file resources are skipped: single files use the userspace AEAD vault
// (internal/fscrypt/filevault.go), never kernel fscrypt (FS_IOC_SET_ENCRYPTION_POLICY can't target
// a standalone file), so neither the `encrypt` feature flag nor `fscrypt setup` applies.
func collectFilesystemPrereqs(vault *fscrypt.Vault, cfg *daemonconfig.Config) ([]fscrypt.Prereq, error) {
	var checkedDevs []uint64
	var out []fscrypt.Prereq
	for _, r := range cfg.EncryptionGroups() {
		if !r.NeedEncryption {
			continue
		}
		root := r.EncryptionRootOrPath()
		info, statErr := os.Stat(root)
		if statErr != nil {
			return nil, fmt.Errorf("stat %s: %w", root, statErr)
		}
		if info.Mode().IsRegular() {
			continue
		}
		dev := info.Sys().(*syscall.Stat_t).Dev
		if slices.Contains(checkedDevs, dev) {
			continue
		}
		checkedDevs = append(checkedDevs, dev)

		ps, err := vault.FilesystemPrereqs(root)
		if err != nil {
			return nil, err
		}
		out = append(out, ps...)
	}
	return out, nil
}

// confirmRunPrereq asks whether the installer should run one remediation
// command now. Negative = abort the install.
func confirmRunPrereq(p fscrypt.Prereq) (bool, error) {
	proceed := false
	err := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("Prerequisite not met: " + p.Title).
			Description(p.Reason +
				"\n\nThe installer can run this now, as root:\n    " + p.Command() +
				"\n\nChoose “Abort” to stop and run it yourself instead.").
			Affirmative("Run it now").
			Negative("Abort").
			Value(&proceed),
	)).Run()
	return proceed, err
}
