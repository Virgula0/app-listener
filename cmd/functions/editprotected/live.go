package editprotected

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/tui"
)

// runLiveEdit is the daemon-running flow. The password is asked ONCE, up
// front — before any protected path is shown — then the daemon returns the
// list of guarded directories, the user picks one, the daemon widens that
// resource's guard for the session, the editor runs, and the grant is
// released. The daemon owns the vault lifecycle throughout.
//
// The password is read BEFORE connecting: the daemon arms a short handshake
// deadline the moment it accepts, so a connection that stays silent while the
// operator types would be dropped as a timeout.
func runLiveEdit() error {
	password, err := promptPassword("Edit-protected password")
	if err != nil {
		return err
	}

	session, err := dialLiveSession()
	if err != nil {
		return err
	}
	defer session.End()

	resources, err := session.Authenticate(password)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	if len(resources) == 0 {
		return errors.New("the daemon reports no guarded directories")
	}

	chosen, err := chooseResource(resources)
	if err != nil {
		return err
	}

	if err := session.Select(chosen); err != nil {
		return fmt.Errorf("activating live edit access to %s: %w", chosen, err)
	}

	log.Infof("editing %s LIVE — the daemon keeps guarding it; write access is granted only for this session", chosen)
	if err := tui.RunFileEditor(chosen); err != nil {
		return fmt.Errorf("edit %s: %w", chosen, err)
	}

	session.End() // narrow the guard back before the audit reads the tree
	auditAfterEdit(chosen)
	return nil
}

// chooseResource honors --resource when given (validating it against the
// daemon's list), otherwise prompts.
func chooseResource(resources []string) (string, error) {
	if resourceFlag != "" {
		if slices.Contains(resources, resourceFlag) {
			return resourceFlag, nil
		}
		return "", fmt.Errorf("%s is not a guarded directory in the running daemon", resourceFlag)
	}
	if len(resources) == 1 {
		log.Infof("only one guarded directory: %s", resources[0])
		return resources[0], nil
	}
	opts := make([]huh.Option[string], 0, len(resources))
	for _, r := range resources {
		opts = append(opts, huh.NewOption(r, r))
	}
	chosen := ""
	if err := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title("Which guarded directory do you want to edit?").
			Description("Write access is granted to this one directory for the edit session, then revoked.").
			Options(opts...).
			Value(&chosen),
	)).Run(); err != nil {
		return "", err
	}
	return chosen, nil
}

// runNonInteractiveLivePut performs one authenticated live write without any
// TUI: content from --content-file or stdin, password from the environment.
func runNonInteractiveLivePut() error {
	if resourceFlag == "" {
		return errors.New("--put requires --resource")
	}
	password := os.Getenv(editPasswordEnv)
	if password == "" {
		return fmt.Errorf("$%s is empty — set it to the edit-protected password for non-interactive --put", editPasswordEnv)
	}

	dest := putFlag
	if !filepath.IsAbs(dest) {
		dest = filepath.Join(resourceFlag, dest)
	}
	if !within(resourceFlag, dest) {
		return fmt.Errorf("--put target %s is outside the resource %s", dest, resourceFlag)
	}
	content, err := readPutContent()
	if err != nil {
		return err
	}

	session, err := dialLiveSession()
	if err != nil {
		return err
	}
	defer session.End()

	resources, err := session.Authenticate(password)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	if !slices.Contains(resources, resourceFlag) {
		return fmt.Errorf("%s is not a guarded directory in the running daemon", resourceFlag)
	}

	if err := session.Select(resourceFlag); err != nil {
		return fmt.Errorf("activating live edit access to %s: %w", resourceFlag, err)
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o700); err != nil {
		return fmt.Errorf("creating parent of %s: %w", dest, err)
	}
	if err := os.WriteFile(dest, content, 0o600); err != nil {
		return fmt.Errorf("writing %s: %w", dest, err)
	}
	log.Infof("wrote %d bytes to %s (live)", len(content), dest)

	session.End()
	auditAfterEdit(resourceFlag)
	return nil
}

func readPutContent() ([]byte, error) {
	if contentFileFlag != "" {
		return os.ReadFile(contentFileFlag)
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("reading --put content from stdin: %w", err)
	}
	return data, nil
}

// promptPassword reads a password without echoing it.
func promptPassword(title string) (string, error) {
	var pw string
	if err := huh.NewForm(huh.NewGroup(
		huh.NewInput().
			Title(title).
			EchoMode(huh.EchoModePassword).
			Value(&pw),
	)).Run(); err != nil {
		return "", err
	}
	if pw == "" {
		return "", fmt.Errorf("no password entered")
	}
	return pw, nil
}
