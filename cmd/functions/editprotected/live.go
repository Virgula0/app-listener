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

// runLiveEdit is the daemon-running flow: password asked ONCE up front (before any protected path
// is shown), the daemon returns the guarded directories, the user picks one, the daemon widens that
// resource's guard for the session, the editor runs, the grant is released. The daemon owns the
// vault lifecycle.
//
// The password is read BEFORE connecting: the daemon arms a short handshake deadline on accept, so
// a connection silent while the operator types would time out.
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
	session.timeout = timeoutSession

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

	if serr := session.Select(chosen); serr != nil {
		return fmt.Errorf("activating live edit access to %s: %w", chosen, serr)
	}

	log.Infof("editing %s LIVE — the daemon keeps guarding it; write access is granted only for this session "+
		"and revoked after %s without input", chosen, timeoutSession)
	_, ended := session.Watch()
	byDaemon, err := tui.RunFileEditorSession(chosen, session.Pinger(pingInterval()), ended)
	if err != nil {
		return fmt.Errorf("edit %s: %w", chosen, err)
	}

	session.End() // narrow the guard back before the audit reads the tree
	auditAfterEdit(chosen, true)
	if byDaemon {
		return errors.New("the daemon ended the live edit session (idle timeout or reload) — " +
			"saves made before it are kept, later ones were not written")
	}
	return nil
}

// chooseResource honors --resource when given (validating it against the
// daemon's list), otherwise prompts.
func chooseResource(resources []string) (string, error) {
	if len(resourceFlags) > 0 {
		if slices.Contains(resources, resourceFlags[0]) {
			return resourceFlags[0], nil
		}
		return "", fmt.Errorf("%s is not a guarded directory in the running daemon", resourceFlags[0])
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
	if len(resourceFlags) == 0 {
		return errors.New("--put requires --resource")
	}
	resource := resourceFlags[0]
	password := os.Getenv(editPasswordEnv)
	if password == "" {
		return fmt.Errorf("$%s is empty — set it to the edit-protected password for non-interactive --put", editPasswordEnv)
	}

	dest := putFlag
	if !filepath.IsAbs(dest) {
		dest = filepath.Join(resource, dest)
	}
	// Lexical containment is only a first gate; writeWithin re-checks every component with
	// O_NOFOLLOW so a symlink in the tree can't redirect the write outside the guarded resource.
	if !within(resource, dest) {
		return fmt.Errorf("--put target %s is outside the resource %s", dest, resource)
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
	session.timeout = timeoutSession

	resources, err := session.Authenticate(password)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	if !slices.Contains(resources, resource) {
		return fmt.Errorf("%s is not a guarded directory in the running daemon", resource)
	}

	if err := session.Select(resource); err != nil {
		return fmt.Errorf("activating live edit access to %s: %w", resource, err)
	}

	if err := writeWithin(resource, dest, content); err != nil {
		return err
	}
	log.Infof("wrote %d bytes to %s (live)", len(content), dest)

	session.End()
	auditAfterEdit(resource, false)
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
