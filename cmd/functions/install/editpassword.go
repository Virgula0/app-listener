package install

import (
	"fmt"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/cmd/functions/editprotected"
	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

// promptEditPassword offers the optional edit-protected password. Collected here, persisted only by
// writeEditPasswordHash from deploy() right before the daemon's first start: after every step that
// could abort (services, config, eBPF preflight), so an aborted install leaves no dangling hash
// file, and the first start already self-guards it (no SIGHUP). Returns "" if declined.
func promptEditPassword(cfg *daemonconfig.Config) (string, error) {
	anyEncrypted := false
	for i := range cfg.Resources {
		if cfg.Resources[i].NeedEncryption {
			anyEncrypted = true
			break
		}
	}

	desc := "Separate from the fscrypt master key: it only authenticates `app-listener edit-protected`. " +
		"With it set, edit-protected can modify a protected directory LIVE, without stopping the daemon."
	if anyEncrypted {
		desc += "\nRecommended: at least one directory will be encrypted."
	}

	want := anyEncrypted
	if err := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().
			Title("Set an edit-protected authentication password? (optional)").
			Description(desc).
			Affirmative("Set a password").
			Negative("Skip").
			Value(&want),
	)).Run(); err != nil {
		return "", err
	}
	if !want {
		return "", nil
	}
	return editprotected.PromptNewPassword()
}

// writeEditPasswordHash persists the password hash (0600 root, origin=install), or is a no-op if
// none was collected. Called from deploy() after
// installServices/installConfig/preflightDeployedBinary succeed but BEFORE systemd.EnableAndVerify
// first starts the daemon: no self-guard on /etc/app-listener yet to fight with, and the first
// Start() picks the hash up directly (no reload).
func writeEditPasswordHash(password string) error {
	if password == "" {
		log.Infof("no edit-protected password set — live edit-protected is disabled "+
			"(add one later with `sudo app-listener edit-protected --set-password`); %s not created", editprotected.HashFile)
		return nil
	}
	hashed, err := editprotected.Hash(password, editprotected.OriginInstall)
	if err != nil {
		return fmt.Errorf("hashing the edit-protected password: %w", err)
	}
	if err := editprotected.WriteHashFile(hashed); err != nil {
		return fmt.Errorf("writing %s: %w", editprotected.HashFile, err)
	}
	log.Infof("edit-protected password hash written to %s (0600, root)", editprotected.HashFile)
	return nil
}
