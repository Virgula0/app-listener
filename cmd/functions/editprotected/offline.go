package editprotected

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/fscrypt"
	"github.com/Virgula0/app-listener/internal/protected"
	"github.com/Virgula0/app-listener/internal/tui"
)

// runOfflineEdit is the daemon-stopped flow: re-scan encrypted catalog dirs, pick ONE, verify the
// master key, unlock, edit, re-lock (a vault never stays open after return).
func runOfflineEdit() error {
	if err := protected.RequireDaemonStopped(); err != nil {
		return err
	}

	vault := fscrypt.New()

	encrypted, err := protected.ScanEncryptedCatalogDirs(vault)
	if err != nil {
		return err
	}
	if len(encrypted) == 0 {
		return errors.New("no catalog directory is currently encrypted: nothing to edit (run `app-listener install` first)")
	}

	chosen, err := pickEncryptedDir(encrypted)
	if err != nil {
		return err
	}

	// Never edit a directory whose policy doesn't unlock with the master key: writes would land in
	// an unreadable policy and the re-lock couldn't complete.
	if err := protected.VerifyEncryptedKeys(vault, []string{chosen}); err != nil {
		return err
	}

	log.Infof("unlocking %s ...", chosen)
	if err := vault.Unlock(chosen); err != nil {
		return fmt.Errorf("unlock %s: %w", chosen, err)
	}

	// Re-lock however the editor exits: a protected directory must never stay unlocked after
	// return. On failure the operator gets the manual command.
	defer func() {
		if err := protected.RelockResources(vault, chosen); err != nil {
			log.Errorf("could not fully re-lock %s: %v — run `fscrypt lock %s` manually as soon as possible", chosen, err, chosen)
		}
	}()

	log.Infof("editing %s (it is unlocked and accessible while the editor is open)", chosen)
	if err := tui.RunFileEditor(chosen); err != nil {
		return fmt.Errorf("edit %s: %w", chosen, err)
	}

	auditAfterEdit(chosen, true)
	return nil
}

// pickEncryptedDir asks which encrypted catalog directory to open (only one vault is unlocked at a
// time); skipped with a single candidate.
func pickEncryptedDir(encrypted []string) (string, error) {
	if len(encrypted) == 1 {
		path := encrypted[0]
		log.Infof("only one encrypted catalog directory found: %s", path)
		return path, nil
	}
	opts := make([]huh.Option[string], 0, len(encrypted))
	for _, p := range encrypted {
		opts = append(opts, huh.NewOption(p, p))
	}
	chosen := ""
	if err := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title("Which protected directory do you want to open?").
			Description("Only one directory is unlocked at a time; it is re-locked when the editor closes.").
			Options(opts...).
			Value(&chosen),
	)).Run(); err != nil {
		return "", err
	}
	return chosen, nil
}
