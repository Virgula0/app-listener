package editprotected

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/x/term"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	inst "github.com/Virgula0/app-listener/internal/install"
)

// errConfigChanged: daemon.conf moved between fetching it and applying the edit.
var errConfigChanged = errors.New("daemon.conf changed while you were editing it")

// runEditConfig edits the running daemon's daemon.conf: fetch it over the control socket, edit it
// locally (no connection is held meanwhile: nothing is widened), then send it back. The daemon
// writes it and reloads; a failed reload keeps the running configuration and restores the file,
// and the operator may edit again. --content-file replaces it non-interactively.
func runEditConfig() error {
	interactive := contentFileFlag == ""
	if interactive && !term.IsTerminal(os.Stdin.Fd()) {
		return errors.New("--edit-config needs a terminal for its editor; use --content-file <file> to replace " +
			"the configuration non-interactively")
	}
	password := os.Getenv(editPasswordEnv)
	if password == "" {
		var err error
		if password, err = promptPassword("Edit-protected password"); err != nil {
			return err
		}
	}

	base, err := fetchConfig(password)
	if err != nil {
		return err
	}
	text := string(base)
	for {
		edited, err := nextConfig(text, interactive)
		if errors.Is(err, errAborted) {
			return nil
		}
		if err != nil {
			return err
		}
		if edited == string(base) {
			log.Info("configuration unchanged — nothing to apply")
			return nil
		}
		text = edited
		if done, err := applyEdit(password, &base, edited, interactive); done {
			return err
		}
	}
}

// applyEdit sends one edit; done=false means back to the editor. A stale base is refreshed so the
// next diff is against what the daemon now holds.
func applyEdit(password string, base *[]byte, edited string, interactive bool) (done bool, err error) {
	err = checkAndApply(password, *base, edited, interactive)
	switch {
	case err == nil:
		log.Info("configuration saved and reloaded by the daemon")
		return true, nil
	case errors.Is(err, errBackToEditor):
		return false, nil
	case errors.Is(err, errConfigChanged):
		fresh, ferr := fetchConfig(password)
		if ferr != nil {
			return true, ferr
		}
		*base = fresh
	}
	if !interactive || !editAgain(err) {
		return true, err
	}
	return false, nil
}

var (
	errAborted      = errors.New("aborted")
	errBackToEditor = errors.New("back to the editor")
)

func nextConfig(text string, interactive bool) (string, error) {
	if !interactive {
		data, err := os.ReadFile(contentFileFlag)
		return string(data), err
	}
	edited, err := inst.EditText("daemon.conf (live — Ctrl+S reviews, applies and reloads)", "daemon.conf", text)
	if errors.Is(err, inst.ErrEditCanceled) {
		log.Info("edit canceled — the configuration is unchanged")
		return "", errAborted
	}
	return edited, err
}

// checkAndApply parses edited locally, shows the diff for confirmation, then sends it.
func checkAndApply(password string, base []byte, edited string, interactive bool) error {
	if _, err := daemonconfig.Parse([]byte(edited)); err != nil {
		return fmt.Errorf("the configuration does not parse: %w", err)
	}
	if interactive {
		if err := inst.ShowDiff("daemon.conf — changes to apply", inst.UnifiedDiff(string(base), edited)); err != nil {
			return err
		}
		apply := false
		if err := huh.NewForm(huh.NewGroup(huh.NewConfirm().
			Title("Apply these changes and reload the daemon?").
			Description("If the reload fails, the daemon keeps the current configuration and restores daemon.conf.").
			Affirmative("Apply").Negative("Back to the editor").
			Value(&apply))).Run(); err != nil {
			return err
		}
		if !apply {
			return errBackToEditor
		}
	}
	return putConfig(password, base, []byte(edited))
}

func editAgain(cause error) bool {
	again := false
	if err := huh.NewForm(huh.NewGroup(huh.NewConfirm().
		Title("The configuration was not applied").
		Description(cause.Error()).
		Affirmative("Edit again").Negative("Discard my changes").
		Value(&again))).Run(); err != nil {
		return false
	}
	return again
}

func fetchConfig(password string) ([]byte, error) {
	session, err := dialLiveSession()
	if err != nil {
		return nil, err
	}
	defer session.End()
	if _, err := session.Authenticate(password); err != nil {
		return nil, fmt.Errorf("authentication failed: %w", err)
	}
	return session.Config()
}

// putConfig sends edited on a fresh connection, only while the daemon's file still holds base
// (the daemon re-checks it before writing).
func putConfig(password string, base, edited []byte) error {
	session, err := dialLiveSession()
	if err != nil {
		return err
	}
	defer session.End()
	if _, err := session.Authenticate(password); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	cur, cerr := session.Config()
	if cerr != nil {
		return cerr
	}
	if !bytes.Equal(cur, base) {
		return errConfigChanged
	}
	return session.PutConfig(edited)
}
