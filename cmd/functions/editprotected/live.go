package editprotected

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/systemd"
	"github.com/Virgula0/app-listener/internal/tui"
)

// runLiveEdit is the daemon-running flow: pick a guarded resource from the
// installed daemon.conf, authenticate to the daemon over the control socket
// (which widens that resource's guard for the session), edit, then release
// the grant. The daemon owns the vault lifecycle throughout — nothing is
// unlocked or locked here.
func runLiveEdit() error {
	cfg, err := daemonconfig.Load(systemd.SystemConfigPath)
	if err != nil {
		return fmt.Errorf("reading the installed daemon config %s: %w", systemd.SystemConfigPath, err)
	}
	if len(cfg.Resources) == 0 {
		return fmt.Errorf("%s has no [watch] sections — nothing to edit", systemd.SystemConfigPath)
	}

	chosen, err := chooseResource(cfg)
	if err != nil {
		return err
	}

	password, err := promptPassword("Edit-protected password")
	if err != nil {
		return err
	}

	log.Infof("authenticating to the daemon for live edit access to %s ...", chosen)
	session, err := beginLiveSession(chosen, password)
	if err != nil {
		return fmt.Errorf("live edit request refused: %w", err)
	}
	defer session.End()

	log.Infof("editing %s LIVE — the daemon keeps guarding it; write access is granted only for this session", chosen)
	if err := tui.RunFileEditor(chosen); err != nil {
		return fmt.Errorf("edit %s: %w", chosen, err)
	}

	session.End() // narrow the guard back before the audit reads the tree
	auditAfterEditWithConfig(cfg, chosen)
	return nil
}

// chooseResource honors --resource when given (validating it is configured),
// otherwise prompts.
func chooseResource(cfg *daemonconfig.Config) (string, error) {
	if resourceFlag != "" {
		for i := range cfg.Resources {
			if cfg.Resources[i].Path == resourceFlag {
				return resourceFlag, nil
			}
		}
		return "", fmt.Errorf("%s is not a [watch] path in %s", resourceFlag, systemd.SystemConfigPath)
	}
	return pickGuardedResource(cfg)
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

	cfg, err := daemonconfig.Load(systemd.SystemConfigPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", systemd.SystemConfigPath, err)
	}
	if findResource(cfg, resourceFlag) == nil {
		return fmt.Errorf("%s is not a [watch] path in %s", resourceFlag, systemd.SystemConfigPath)
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

	session, err := beginLiveSession(resourceFlag, password)
	if err != nil {
		return fmt.Errorf("live edit request refused: %w", err)
	}
	defer session.End()

	if err := os.MkdirAll(filepath.Dir(dest), 0o700); err != nil {
		return fmt.Errorf("creating parent of %s: %w", dest, err)
	}
	if err := os.WriteFile(dest, content, 0o600); err != nil {
		return fmt.Errorf("writing %s: %w", dest, err)
	}
	log.Infof("wrote %d bytes to %s (live)", len(content), dest)

	session.End()
	auditAfterEditWithConfig(cfg, resourceFlag)
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

// pickGuardedResource asks which configured watch path to open. A single
// resource skips the prompt.
func pickGuardedResource(cfg *daemonconfig.Config) (string, error) {
	if len(cfg.Resources) == 1 {
		log.Infof("only one guarded resource: %s", cfg.Resources[0].Path)
		return cfg.Resources[0].Path, nil
	}
	opts := make([]huh.Option[string], 0, len(cfg.Resources))
	for i := range cfg.Resources {
		r := &cfg.Resources[i]
		label := r.Path
		if r.NeedEncryption {
			label += "  (encrypted)"
		}
		opts = append(opts, huh.NewOption(label, r.Path))
	}
	chosen := ""
	if err := huh.NewForm(huh.NewGroup(
		huh.NewSelect[string]().
			Title("Which guarded directory do you want to edit?").
			Description("The daemon grants write access to this one resource for the edit session, then revokes it.").
			Options(opts...).
			Value(&chosen),
	)).Run(); err != nil {
		return "", err
	}
	return chosen, nil
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
