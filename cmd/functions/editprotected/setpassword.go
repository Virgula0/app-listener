package editprotected

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/systemd"
)

// runSetPassword sets a first edit-protected password or rotates a
// cli-managed one. A password chosen during `app-listener install` cannot be
// rotated here (the installer owns it) — that is refused with a clear
// pointer.
func runSetPassword() error {
	encoded, err := LoadHashFile()
	switch {
	case errors.Is(err, ErrNoHashFile):
		// First-time set.
	case err != nil:
		return err
	default:
		origin, oErr := OriginOf(encoded)
		if oErr != nil {
			return fmt.Errorf("the existing password file is unreadable (%w) — remove %s and set a new one", oErr, HashFile)
		}
		if origin == OriginInstall {
			return errors.New("the edit-protected password was set during installation and cannot be rotated here — " +
				"re-run `sudo app-listener install` to change it (it re-prompts for the password)")
		}
		current, pErr := promptPassword("Current edit-protected password")
		if pErr != nil {
			return pErr
		}
		ok, vErr := Verify(encoded, current)
		if vErr != nil {
			return vErr
		}
		if !ok {
			return errors.New("current password is incorrect")
		}
	}

	newPw, err := PromptNewPassword()
	if err != nil {
		return err
	}
	hashed, err := Hash(newPw, OriginCLI)
	if err != nil {
		return err
	}
	if err := WriteHashFile(hashed); err != nil {
		return fmt.Errorf("writing %s: %w", HashFile, err)
	}
	log.Infof("edit-protected password saved to %s", HashFile)
	return reloadDaemonForPasswordChange()
}

// runClearPassword removes the password (of any origin) after proving
// knowledge of it, disabling live mode.
func runClearPassword() error {
	encoded, err := LoadHashFile()
	if errors.Is(err, ErrNoHashFile) {
		log.Info("no edit-protected password is configured — nothing to clear")
		return nil
	}
	if err != nil {
		return err
	}

	current, err := promptPassword("Current edit-protected password (to confirm removal)")
	if err != nil {
		return err
	}
	ok, err := Verify(encoded, current)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("password is incorrect — not clearing")
	}

	if err := RemoveHashFile(); err != nil {
		return fmt.Errorf("removing %s: %w", HashFile, err)
	}
	log.Infof("edit-protected password removed (%s) — live edit-protected is now disabled", HashFile)
	return reloadDaemonForPasswordChange()
}

// PromptNewPassword reads the new password twice and enforces the strength
// floor. Exported for the installer's optional edit-protected password step.
func PromptNewPassword() (string, error) {
	first, err := promptPassword("New edit-protected password")
	if err != nil {
		return "", err
	}
	if vErr := ValidatePassword(first); vErr != nil {
		return "", vErr
	}
	again, err := promptPassword("New edit-protected password (again)")
	if err != nil {
		return "", err
	}
	if first != again {
		return "", errors.New("the two entries do not match")
	}
	return first, nil
}

// reloadDaemonForPasswordChange asks a running daemon to reconcile its
// control socket with the new password state (SIGHUP; the reload handler
// calls control.refresh). A stopped daemon picks it up on next start.
func reloadDaemonForPasswordChange() error {
	if !systemd.IsDaemonActive() {
		return nil
	}
	log.Info("reloading the daemon so it picks up the password change ...")
	if err := systemd.ReloadDaemonIfActive(); err != nil {
		return fmt.Errorf("the password file was written but reloading the daemon failed (%w) — run: sudo systemctl reload %s",
			err, systemd.DaemonServiceName)
	}
	return nil
}
