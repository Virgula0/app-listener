package editprotected

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/protected"
	"github.com/Virgula0/app-listener/internal/systemd"
)

// authorizeFirstTimeSetPassword gates creating the FIRST edit-protected password (the transition
// from "no password" to "live mode armed"), for installs older than the mandatory install-time
// password. The daemon seals binaries.db so no non-daemon binary can write it, but that check is
// per-binary, not per-caller: `edit-protected --set-password` IS the daemon binary, so absent this gate any root process could arm live mode on a running
// production daemon (whose vaults are unlocked) with no interaction and no proof of ownership, then
// read every protected resource through the control socket. The audited bootstrap for the first
// password is `install`; refuse the first-time CLI set while the daemon is running and point there.
// Offline (daemon stopped, vaults locked) it is allowed: that flow already requires stopping the
// daemon, which a legitimate operator can do and a remote attacker cannot silently.
func authorizeFirstTimeSetPassword(daemonRunning bool) error {
	if daemonRunning {
		return errors.New("refusing to set the first edit-protected password while the daemon is running: " +
			"the first password is created by `sudo app-listener install` (the audited bootstrap). " +
			"Either run the installer, or stop the daemon first (systemctl stop app-listener-daemon) to set it offline")
	}
	return nil
}

// daemonRunning is anyDaemonRunning; tests replace it.
var daemonRunning = anyDaemonRunning

// anyDaemonRunning: the systemd unit, its control socket, or a daemon process outside systemd. An
// unreadable /proc counts as running, the stricter path.
func anyDaemonRunning() bool {
	if protected.DaemonRunning() || LiveModeAvailable() {
		return true
	}
	pids, err := protected.FindDaemonProcesses("/proc")
	return err != nil || len(pids) > 0
}

// runSetPassword changes the edit-protected password, or sets the first one with the daemon stopped.
func runSetPassword() error {
	return setPassword(func() (string, error) { return promptPassword("Current edit-protected password") },
		PromptNewPassword)
}

// setPassword is runSetPassword's flow, shared with the GUI backend: current is asked only when a
// password exists, next yields the new one. A change proves the current password first: to the
// daemon over its control socket while it runs (AUTH, throttled, then PASSWORD), else through
// CheckPassword's persistent guess budget, never an unthrottled local compare. There is no way to
// remove the password.
func setPassword(current, next func() (string, error)) error {
	configured, err := HashConfigured()
	if err != nil {
		return err
	}
	if !configured {
		// First-time set: arming live edit-protected on a resource that had no password.
		if aErr := authorizeFirstTimeSetPassword(daemonRunning()); aErr != nil {
			return aErr
		}
		return storeNewPassword(next)
	}
	cur, err := current()
	if err != nil {
		return err
	}
	if daemonRunning() {
		return changePasswordLive(cur, next)
	}
	if err := CheckPassword(cur); err != nil {
		if errors.Is(err, ErrWrongPassword) {
			return errors.New("current password is incorrect")
		}
		return err
	}
	return storeNewPassword(next)
}

func storeNewPassword(next func() (string, error)) error {
	newPw, err := next()
	if err != nil {
		return err
	}
	hashed, err := Hash(newPw, OriginCLI)
	if err != nil {
		return err
	}
	if err := StoreHash(hashed); err != nil {
		return fmt.Errorf("storing the edit-protected password: %w", err)
	}
	log.Info("edit-protected password saved")
	return nil
}

// changePasswordLive changes the password through the running daemon, which stores it itself.
func changePasswordLive(cur string, next func() (string, error)) error {
	if !LiveModeAvailable() {
		return errors.New("the daemon is running but its control socket " + ControlSocket + " is not present: " +
			"reload it (systemctl reload " + systemd.DaemonServiceName + ") or stop it, then retry")
	}
	session, err := dialLiveSession()
	if err != nil {
		return err
	}
	defer session.End()
	if _, authErr := session.Authenticate(cur); authErr != nil {
		return fmt.Errorf("the daemon refused the current password: %w", authErr)
	}
	newPw, err := next()
	if err != nil {
		return err
	}
	if err := session.ChangePassword(newPw); err != nil {
		return fmt.Errorf("the daemon refused the new password: %w", err)
	}
	log.Info("edit-protected password changed")
	return nil
}

// invalidPassword is a refused new password: too weak, or the two entries differ.
type invalidPassword struct{ error }

// passwordRules is shown with every new-password prompt.
const passwordRules = "At least 12 characters, no whitespace, three of: lowercase, uppercase, digits, symbols." //nolint:gosec // help text, not a credential

// PromptNewPassword reads the new password twice and enforces the strength floor.
func PromptNewPassword() (string, error) { return promptNewPassword("") }

// PromptNewPasswordUntilValid asks again, saying why, until an entry is accepted or the prompt is
// canceled: the installer's mandatory step.
func PromptNewPasswordUntilValid() (string, error) {
	problem := ""
	for {
		pw, err := promptNewPassword(problem)
		var inv invalidPassword
		if !errors.As(err, &inv) {
			return pw, err
		}
		problem = inv.Error() + ". Try again."
	}
}

func promptNewPassword(problem string) (string, error) {
	desc := passwordRules
	if problem != "" {
		desc = problem + "\n" + passwordRules
	}
	first, err := promptPasswordWith("New edit-protected password", desc)
	if err != nil {
		return "", err
	}
	if vErr := ValidatePassword(first); vErr != nil {
		return "", invalidPassword{vErr}
	}
	again, err := promptPassword("New edit-protected password (again)")
	if err != nil {
		return "", err
	}
	if first != again {
		return "", invalidPassword{errors.New("the two entries do not match")}
	}
	return first, nil
}
