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

// runOfflineEdit is the daemon-stopped flow: re-scan which catalog
// directories are encrypted, let the user pick ONE, verify the master key,
// unlock, edit, and re-lock — a vault never stays open after this returns.
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

	// A directory whose policy does not unlock with the master key must
	// never be edited: the writes would land in an unreadable policy and
	// the re-lock could not complete.
	if err := protected.VerifyEncryptedKeys(vault, []string{chosen}); err != nil {
		return err
	}

	log.Infof("unlocking %s ...", chosen)
	if err := vault.Unlock(chosen); err != nil {
		return fmt.Errorf("unlock %s: %w", chosen, err)
	}

	// Re-lock no matter how the editor exits: a protected directory must
	// never stay unlocked after this command returns. On failure the
	// operator gets the manual command.
	defer func() {
		if err := protected.RelockResources(vault, chosen); err != nil {
			log.Errorf("could not fully re-lock %s: %v — run `fscrypt lock %s` manually as soon as possible", chosen, err, chosen)
		}
	}()

	log.Infof("editing %s (it is unlocked and accessible while the editor is open)", chosen)
	if err := tui.RunFileEditor(chosen); err != nil {
		return fmt.Errorf("edit %s: %w", chosen, err)
	}

	auditAfterEdit(chosen)
	return nil
}

// pickEncryptedDir asks which of the encrypted catalog directories to open.
// Only one vault is ever unlocked at a time, so exactly one directory is
// returned. With a single candidate the picker is skipped.
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
