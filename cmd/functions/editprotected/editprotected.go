// The `app-listener edit-protected` command opens the fscrypt-encrypted
// catalog directories in the embedded two-pane editor. It fatally refuses
// while the daemon is running, re-scans the catalog (NOT the installed
// daemon.conf) exactly like the uninstaller, verifies the master key, lets
// the user pick ONE directory to open (only one vault is ever unlocked at
// a time), unlocks it, runs the TUI editor, and re-locks it again no
// matter how the editor exits: a vault never stays open after this
// command returns.
package editprotected

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/internal/fscrypt"
	"github.com/Virgula0/app-listener/internal/protected"
	"github.com/Virgula0/app-listener/internal/tui"
)

var EditProtectedCmd = &cobra.Command{
	Use:   "edit-protected",
	Short: "Edit one of the fscrypt-encrypted catalog directories in the embedded TUI editor",
	Long: `Open one of the fscrypt-encrypted catalog directories in the embedded
two-pane editor instead of configuring the daemon or installing anything.

The command re-scans the catalog (NOT the installed daemon.conf) like the
uninstaller does, verifies that every found directory unlocks with the
master key, lets you pick ONE directory to open (only one is ever unlocked
at a time), unlocks it with the master key, runs the editor, and
re-locks it again no matter how the editor exits — a protected directory
never stays unlocked after this command returns.

It fatally refuses while the daemon is running: the directory being edited
is the very one the daemon guards and unlocks.`,
	Args: cobra.NoArgs,
	RunE: runEditProtected,
}

// runEditProtected drives the whole edit-protected flow: require the daemon
// to be stopped, re-scan which catalog directories are encrypted, let the
// user pick ONE, verify the master key, unlock, edit, and re-lock.
func runEditProtected(cmd *cobra.Command, args []string) error {
	if os.Geteuid() != 0 {
		return errors.New("edit-protected must be run as root: sudo app-listener edit-protected")
	}

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
