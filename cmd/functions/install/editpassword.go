package install

import (
	"fmt"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/cmd/functions/editprotected"
	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

// promptEditPassword offers the optional edit-protected authentication
// password. It is collected here but persisted only by writeEditPasswordHash,
// called from deploy() right before the daemon's first start — after every
// step that could still abort the install (services, config, eBPF preflight)
// but before there is a running daemon to reload, so an aborted install never
// leaves a dangling hash file AND the very first start already self-guards it
// (no SIGHUP needed afterward). Returns "" when the user declines.
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

// writeEditPasswordHash persists the edit-protected password hash (0600
// root, origin=install), or is a no-op when no password was collected.
// Called from deploy(), after installServices/installConfig/
// preflightDeployedBinary have all succeeded but BEFORE systemd.EnableAndVerify
// brings up the daemon for the first time: the daemon isn't running yet, so
// there is no self-guard on /etc/app-listener to fight with, and that first
// Start() picks up the hash file directly — no SIGHUP reload required
// afterward (that used to re-attach every guard a second time on every
// fresh install).
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
