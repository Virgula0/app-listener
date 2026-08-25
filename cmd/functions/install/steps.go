package install

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	inst "github.com/Virgula0/app-listener/internal/install"
)

// pickUsers asks which local users the installation should protect. Every
// user (root included) is preselected; root's /root is probed like any
// other user, but no ssh-agent unit is installed for root (see
// installSSHAgent).
func pickUsers() ([]inst.User, error) {
	users, err := inst.ListUsers()
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("no login users with a home directory found in /etc/passwd")
	}

	opts := make([]huh.Option[inst.User], 0, len(users))
	for i := range users {
		u := &users[i]
		label := fmt.Sprintf("%s (uid %d, home %s)", u.Name, u.UID, u.Home)
		if u.UID == 0 {
			label += " — no ssh-agent unit will be installed"
		}
		opts = append(opts, huh.NewOption(label, *u).Selected(true))
	}
	var picked []inst.User
	form := huh.NewForm(huh.NewGroup(
		huh.NewMultiSelect[inst.User]().
			Title("Which users should be protected?").
			Description("Every selected user's critical directories (SSH, keys, browsers, ...) are scanned in the next step.").
			Options(opts...).
			Height(10).
			Value(&picked),
	))
	if err := form.WithKeyMap(selectionKeymap()).Run(); err != nil {
		return nil, err
	}
	if len(picked) == 0 {
		return nil, fmt.Errorf("no users selected")
	}
	selectedUsers = picked
	for i := range picked {
		log.Infof("protecting user %s (home %s)", picked[i].Name, picked[i].Home)
	}
	return picked, nil
}

// pickDirectories probes the catalog for every selected user (and the
// system-level entries) and asks which of the found critical directories
// to protect. All are preselected.
func pickDirectories(users []inst.User) ([]inst.Candidate, error) {
	candidates := inst.DiscoverForUsers(users)
	if len(candidates) == 0 {
		log.Warn("no catalog directories found for the selected users — you can still add directories manually")
		return nil, nil
	}

	opts := make([]huh.Option[int], 0, len(candidates))
	for i := range candidates {
		c := &candidates[i]
		label := fmt.Sprintf("%s  %s", c.Entry.Name, c.Path)
		if c.User.Name != "" {
			label = fmt.Sprintf("%s (user %s)  %s", c.Entry.Name, c.User.Name, c.Path)
		}
		opts = append(opts, huh.NewOption(label, i).Selected(true))
	}
	var pickedIdx []int
	form := huh.NewForm(huh.NewGroup(
		huh.NewMultiSelect[int]().
			Title("Critical directories found — select the ones to protect").
			Description("Only existing paths are listed. All are preselected. Whitelisted binaries per directory are curated and minimal.").
			Options(opts...).
			Height(12).
			Value(&pickedIdx),
	))
	if err := form.WithKeyMap(selectionKeymap()).Run(); err != nil {
		return nil, err
	}
	var picked []inst.Candidate
	for _, i := range pickedIdx {
		picked = append(picked, candidates[i])
	}
	if len(picked) == 0 {
		log.Warn("no catalog directories selected — you can still add directories manually")
		return nil, nil
	}
	for i := range picked {
		c := &picked[i]
		allowed := c.FilterExistingWhitelist()
		log.Infof("will protect %s (%s) — %d whitelisted binaries", c.Path, c.Entry.Name, len(allowed))
	}
	return picked, nil
}

// selectionKeymap customizes the multi-select keys used by the user and
// directory pickers: Ctrl+K selects/deselects all entries, space/x selects
// one entry at a time. The legend at the bottom of the field shows both.
// It must be applied to the Form (not the field): NewForm overwrites
// every field's keymap with the form default.
func selectionKeymap() *huh.KeyMap {
	keys := huh.NewDefaultKeyMap()
	keys.MultiSelect.SelectAll = key.NewBinding(
		key.WithKeys("ctrl+k"), key.WithHelp("ctrl+k", "select/deselect all"))
	keys.MultiSelect.SelectNone = key.NewBinding(
		key.WithKeys("ctrl+k"), key.WithHelp("ctrl+k", "select/deselect all"))
	keys.MultiSelect.Toggle = key.NewBinding(
		key.WithKeys(" ", "x"), key.WithHelp("space/x", "select one"))
	return keys
}

// addManualDirectories asks for additional paths that were not discovered
// by the catalog. A manually entered path that does not exist is a fatal
// error, per the installer contract.
func addManualDirectories(candidates []inst.Candidate) ([]inst.Candidate, error) {
	for {
		var input string
		if err := huh.NewForm(huh.NewGroup(
			huh.NewInput().
				Title("Additional directory or file to protect").
				Description("Leave empty and confirm to continue. The path must exist. Whitelist its binaries later in the editor.").
				Prompt("> ").
				Placeholder("/path/to/extra-secret (empty = done)").
				Value(&input),
		)).Run(); err != nil {
			return nil, err
		}
		input = strings.TrimSpace(input)
		if input == "" {
			return candidates, nil
		}
		if !strings.HasPrefix(input, "/") {
			return nil, fmt.Errorf("fatal: %q is not an absolute path", input)
		}
		if _, err := os.Lstat(input); err != nil {
			return nil, fmt.Errorf("fatal: %q does not exist: %w", input, err)
		}
		duplicate := false
		for i := range candidates {
			if candidates[i].Path == input {
				duplicate = true
				break
			}
		}
		if duplicate {
			log.Warnf("%s was already selected, skipping", input)
			continue
		}
		log.Warnf("%s added manually: no binaries whitelisted yet — add them in the editor or every access will be denied", input)
		candidates = append(candidates, inst.Candidate{
			User:  inst.User{},
			Entry: inst.CandidateDir{Name: "manual"},
			Path:  input,
		})
	}
}

// editConfig renders the config from the selected candidates and opens the
// embedded editor. The user's version is validated with the real parser;
// an invalid config re-opens the editor until it parses or the user
// aborts.
func editConfig(candidates []inst.Candidate) (string, *daemonconfig.Config, error) {
	sections := make([]inst.Section, 0, len(candidates))
	for i := range candidates {
		c := &candidates[i]
		sections = append(sections, inst.Section{
			Path:    c.Path,
			Allow:   c.FilterExistingWhitelist(),
			Encrypt: true,
		})
	}
	confText := inst.GenerateConf(sections)

	for {
		edited, err := inst.EditText("app-listener daemon.conf — review and save (Ctrl+S)", confText)
		if err != nil {
			return "", nil, fmt.Errorf("config editing aborted: %w", err)
		}

		cfg, err := validateConfigText(edited)
		if err != nil {
			log.Errorf("config is invalid: %v — fix it and save again (or press Esc to abort)", err)
			confText = edited
			continue
		}
		return edited, cfg, nil
	}
}

// validateConfigText parses the edited configuration through the same
// strict parser the daemon uses.
func validateConfigText(text string) (*daemonconfig.Config, error) {
	tmp, err := os.CreateTemp("", "app-listener-conf-*.conf")
	if err != nil {
		return nil, err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if _, err := tmp.WriteString(text); err != nil {
		tmp.Close()
		return nil, err
	}
	if err := tmp.Close(); err != nil {
		return nil, err
	}
	return daemonconfig.Load(tmpPath)
}

// askEncryption asks, per resource that is NOT yet encrypted and that
// declares need_encryption: true, whether fscrypt encryption is required
// (default yes) and records the decision back into the config text.
// Resources declared need_encryption: false are never asked — the config
// already says no encryption and the installer honors it. Already-encrypted
// resources are never asked: they keep need_encryption: true and are not
// migrated (no backup is created for them). Both directories and single
// regular files are offered; symlinks, hardlinks and special files never
// reach this step (the daemon config parser refuses them). It returns the
// updated text and the list of resources that still need to be encrypted.
func askEncryption(vault *fscrypt.Vault, cfgText string, cfg *daemonconfig.Config) (text string, toEncrypt []string, err error) {
	for _, r := range cfg.Resources {
		if !r.NeedEncryption {
			log.Infof("%s declares need_encryption: false: skipping the encryption question", r.Path)
			continue
		}
		encrypted, encErr := vault.IsEncrypted(r.Path)
		if encErr != nil {
			return "", nil, fmt.Errorf("checking encryption of %s: %w", r.Path, encErr)
		}
		if encrypted {
			log.Infof("%s is already encrypted: keeping need_encryption: true (no question, no backup)", r.Path)
			continue
		}

		answer := true
		formErr := huh.NewForm(huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Require fscrypt encryption for %s?", r.Path)).
				Description("The resource is not encrypted yet: the installer will encrypt it in place, keeping a backup.").
				Affirmative("Yes, encrypt").
				Negative("No encryption").
				Value(&answer),
		)).Run()
		if formErr != nil {
			return "", nil, formErr
		}
		updated, updateErr := inst.SetNeedEncryption(cfgText, r.Path, answer)
		if updateErr != nil {
			return "", nil, updateErr
		}
		cfgText = updated
		if answer {
			toEncrypt = append(toEncrypt, r.Path)
		}
	}
	return cfgText, toEncrypt, nil
}
