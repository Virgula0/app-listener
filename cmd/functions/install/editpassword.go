package install

import (
	"fmt"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/cmd/functions/editprotected"
	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/systemd"
)

// promptEditPassword offers the optional edit-protected authentication
// password. It is collected here but persisted only as the very last install
// action (finalizeEditPassword), so an aborted install never leaves a
// dangling hash file. Returns "" when the user declines.
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

// finalizeEditPassword is the very last install step: persist the
// edit-protected password hash (0600 root, origin=install) and, when the
// daemon is already running, reload it so its self-protection guard covers
// the new file and the control socket comes up.
func finalizeEditPassword(password string) error {
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
	log.Infof("edit-protected password hash written to %s (0600, root) — this was the final install step", editprotected.HashFile)

	if systemd.IsDaemonActive() {
		log.Info("reloading the daemon so it self-guards the password file and opens the control socket ...")
		if err := systemd.ReloadDaemonIfActive(); err != nil {
			return fmt.Errorf("password saved but reloading the daemon failed (%w) — run: sudo systemctl reload %s",
				err, systemd.DaemonServiceName)
		}
	}
	return nil
}
