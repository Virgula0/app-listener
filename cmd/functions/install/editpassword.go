package install

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/cmd/functions/editprotected"
	"github.com/Virgula0/app-listener/internal/wizard"
)

// promptEditPassword asks for the mandatory edit-protected password. An existing one (a resumed
// install) is kept unless the user replaces it: the installer runs as root with the daemon stopped,
// the audited bootstrap. Collected here, persisted only by writeEditPasswordHash
// from deploy() right before the daemon's first start, after every step that could abort. Returns
// "" when the existing password is kept.
func promptEditPassword() (string, error) {
	if err := editprotected.MigrateLegacyHash(); err != nil {
		log.Errorf("edit-protected password: %v", err)
	}
	configured, err := editprotected.HashConfigured()
	if err != nil {
		return "", err
	}
	if configured {
		keep, err := wizard.Confirm(wizard.Question{
			Title: "Keep the current edit-protected password?",
			Description: "It authenticates `app-listener edit-protected` (live edits, temporary access, " +
				"configuration edits). Change it any time with `edit-protected password`, which asks the " +
				"current one. Replace it here only if it is forgotten.",
			Affirmative: "Keep it",
			Negative:    "Replace it",
			Default:     true,
		})
		if err != nil {
			return "", err
		}
		if keep {
			return "", nil
		}
	}
	return editprotected.PromptNewPasswordUntilValid()
}

// writeEditPasswordHash persists the password hash (origin=install) in binaries.db, or keeps the
// existing one when password is "". Called from deploy() after
// installServices/installConfig/preflightDeployedBinary succeed but BEFORE systemd.EnableAndVerify
// first starts the daemon: no self guard on binaries.db yet, and the first Start() reads it.
func writeEditPasswordHash(password string) error {
	if password == "" {
		configured, err := editprotected.HashConfigured()
		if err != nil {
			return err
		}
		if !configured {
			return errors.New("no edit-protected password was set: it is mandatory")
		}
		log.Info("keeping the current edit-protected password")
		return nil
	}
	hashed, err := editprotected.Hash(password, editprotected.OriginInstall)
	if err != nil {
		return fmt.Errorf("hashing the edit-protected password: %w", err)
	}
	if err := editprotected.StoreHash(hashed); err != nil {
		return fmt.Errorf("storing the edit-protected password: %w", err)
	}
	log.Info("edit-protected password stored in the binary ledger")
	return nil
}
