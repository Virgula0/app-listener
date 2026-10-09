package install

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/repository"
	"github.com/Virgula0/app-listener/internal/systemd"
	"github.com/Virgula0/app-listener/internal/wizard"
)

// pickUsers asks which local users to protect. All (root included) are preselected. The per-user
// ssh-agent setup is asked afterwards, only for a user whose ~/.ssh ends up guarded
// (askSSHAgentUsers).
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
	if err := form.WithKeyMap(wizard.MultiSelectKeymap()).Run(); err != nil {
		return nil, err
	}
	if len(picked) == 0 {
		return nil, fmt.Errorf("no users selected")
	}
	for i := range picked {
		log.Infof("protecting user %s (home %s)", picked[i].Name, picked[i].Home)
	}
	return picked, nil
}

// catalogGroup bundles every discovered path of one resource (one catalog Name, one user) behind a
// single checkbox, so a multi-location resource (Steam's data dir + legacy home) is selected as a
// unit.
type catalogGroup struct {
	label      string
	candidates []inst.Candidate
}

// groupCandidates collapses per-path candidates into per-resource groups, in first-seen order.
func groupCandidates(cands []inst.Candidate) []catalogGroup {
	type key struct{ name, user string }
	var order []key
	byKey := map[key][]int{}
	for i := range cands {
		k := key{cands[i].Entry.Name, cands[i].User.Name}
		if _, ok := byKey[k]; !ok {
			order = append(order, k)
		}
		byKey[k] = append(byKey[k], i)
	}
	out := make([]catalogGroup, 0, len(order))
	for _, k := range order {
		idxs := byKey[k]
		cs := make([]inst.Candidate, 0, len(idxs))
		paths := make([]string, 0, len(idxs))
		for _, i := range idxs {
			cs = append(cs, cands[i])
			paths = append(paths, cands[i].Path)
		}
		label := k.name
		if k.user != "" {
			label = fmt.Sprintf("%s (user %s)", k.name, k.user)
		}
		if len(paths) == 1 {
			label += "  " + paths[0]
		} else {
			label += fmt.Sprintf("  [%d locations] %s", len(paths), strings.Join(paths, ", "))
		}
		out = append(out, catalogGroup{label: label, candidates: cs})
	}
	return out
}

// pickDirectories probes the catalog for every selected user (plus system entries) and asks which
// found directories to protect. All preselected.
func pickDirectories(users []inst.User) ([]inst.Candidate, error) {
	candidates := inst.DiscoverForUsers(users)
	if len(candidates) == 0 {
		log.Warn("no catalog directories found for the selected users — you can still add directories manually")
		return nil, nil
	}
	picked, err := pickFromCandidates(candidates,
		"Critical directories found — select the ones to protect",
		"Only existing paths are listed. All are preselected. A resource stored in several locations is one entry. Whitelisted binaries per directory are curated and minimal.")
	if err != nil {
		return nil, err
	}
	if len(picked) == 0 {
		log.Warn("no catalog directories selected — you can still add directories manually")
	}
	return picked, nil
}

// pickFromCandidates groups candidates per resource (one checkbox per catalog Name + user), runs
// the multi-select, and returns the flattened per-path candidates of picked groups.
func pickFromCandidates(candidates []inst.Candidate, title, description string) ([]inst.Candidate, error) {
	groups := groupCandidates(candidates)
	opts := make([]huh.Option[int], 0, len(groups))
	for i := range groups {
		opts = append(opts, huh.NewOption(groups[i].label, i).Selected(true))
	}
	var pickedIdx []int
	form := huh.NewForm(huh.NewGroup(
		huh.NewMultiSelect[int]().
			Title(title).
			Description(description).
			Options(opts...).
			Height(12).
			Value(&pickedIdx),
	))
	if err := form.WithKeyMap(wizard.MultiSelectKeymap()).Run(); err != nil {
		return nil, err
	}
	var picked []inst.Candidate
	for _, i := range pickedIdx {
		picked = append(picked, groups[i].candidates...)
	}
	for i := range picked {
		c := &picked[i]
		allowed := c.FilterExistingWhitelist()
		log.Infof("will protect %s (%s) — %d whitelisted binaries", c.Path, c.Entry.Name, len(allowed))
	}
	return picked, nil
}

// addManualDirectories asks for extra paths not found by the catalog. A manually entered path that
// doesn't exist is a fatal error (installer contract).
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

// sectionsFromCandidates turns candidates into config sections: filtered whitelist, encryption on
// by default, and the catalog entry's grouped extra watch sub-paths.
func sectionsFromCandidates(candidates []inst.Candidate) []inst.Section {
	sections := make([]inst.Section, 0, len(candidates))
	for i := range candidates {
		c := &candidates[i]
		sections = append(sections, inst.Section{
			Path:            c.Path,
			Allow:           c.FilterExistingWhitelist(),
			Encrypt:         true,
			ExtraWatchPaths: c.Entry.ExtraWatchPathsFor(c.User.Home, c.User.Name),
		})
	}
	return sections
}

// libraryBlocksFromCandidates returns one [libraries] block per catalog entry and user, in
// candidate order. An entry may yield several watch sections (Steam has three locations), but its
// library directives describe the APPLICATION, so they're declared once in their own block.
func libraryBlocksFromCandidates(candidates []inst.Candidate) []inst.LibraryBlock {
	seen := make(map[string]bool, len(candidates))
	var blocks []inst.LibraryBlock
	for i := range candidates {
		c := &candidates[i]
		block := c.Entry.LibraryBlockFor(c.User.Name, c.User.Home)
		if seen[block.Name] || block.Empty() {
			continue
		}
		seen[block.Name] = true
		blocks = append(blocks, block)
	}
	return blocks
}

// editConfig renders the config from the selected candidates and opens the embedded editor.
func editConfig(candidates []inst.Candidate) (string, *daemonconfig.Config, error) {
	return runConfigEditor(
		"app-listener daemon.conf — review and save (Ctrl+S)",
		inst.GenerateConf(sectionsFromCandidates(candidates), libraryBlocksFromCandidates(candidates)), nil)
}

// runConfigEditor opens initial in the embedded editor and validates the result with the daemon's
// strict parser; an invalid config re-opens the editor until it parses or the user aborts (Esc). A
// whitelist line not in prev (nil: any) that the daemon would drop as a multicall is invalid too.
func runConfigEditor(title, initial string, prev *daemonconfig.Config) (string, *daemonconfig.Config, error) {
	confText := initial
	for {
		edited, err := inst.EditText(title, "daemon.conf", confText)
		if err != nil {
			return "", nil, fmt.Errorf("config editing aborted: %w", err)
		}

		cfg, err := validateConfigText(edited)
		if err == nil {
			err = guard.RefuseMulticallLines(cfg, prev)
		}
		if err != nil {
			log.Errorf("config is invalid: %v — fix it and save again (or press Esc to abort)", err)
			confText = edited
			continue
		}
		return edited, cfg, nil
	}
}

// validateConfigText parses the edited config with the daemon's strict parser.
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

// askEncryption asks, per resource that is NOT yet encrypted and declares need_encryption: true,
// whether fscrypt encryption is required (default yes), recording the answer in the config text.
// Never asked: need_encryption: false resources (honored as-is) and already-encrypted ones (keep
// need_encryption: true, not migrated, no backup). Directories and single regular files are
// offered; symlinks, hardlinks and special files never reach here (the parser refuses them).
// Returns the updated text and the resources still to encrypt.
//
// Asked once per ENCRYPTION GROUP, not per watch path: grouped sections share the group's vault, so
// the answer and in-place encryption target the group root (addressing by watch sub-path failed
// with "section not found": a grouped config has one header per group).
func askEncryption(vault *fscrypt.Vault, cfgText string, cfg *daemonconfig.Config) (text string, toEncrypt []string, err error) {
	for _, r := range cfg.EncryptionGroups() {
		// Grouped sections share ONE [watch <root>] header: the fscrypt lifecycle and every
		// config-text patch address the encryption root, never a watch sub-path (no header of its
		// own; SetNeedEncryption would fail "section not found").
		sectionPath := r.EncryptionRootOrPath()
		if !r.NeedEncryption {
			log.Infof("%s declares need_encryption: false: skipping the encryption question", sectionPath)
			continue
		}
		encrypted, encErr := vault.IsEncrypted(sectionPath)
		if encErr != nil {
			return "", nil, fmt.Errorf("checking encryption of %s: %w", sectionPath, encErr)
		}
		if encrypted {
			log.Infof("%s is already encrypted: keeping need_encryption: true (no question, no backup)", sectionPath)
			continue
		}

		answer := true
		formErr := huh.NewForm(huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Require fscrypt encryption for %s?", sectionPath)).
				Description("The resource is not encrypted yet: the installer will encrypt it in place, keeping a backup.").
				Affirmative("Yes, encrypt").
				Negative("No encryption").
				Value(&answer),
		)).Run()
		if formErr != nil {
			return "", nil, formErr
		}
		updated, updateErr := inst.SetNeedEncryption(cfgText, sectionPath, answer)
		if updateErr != nil {
			return "", nil, updateErr
		}
		cfgText = updated
		if answer {
			toEncrypt = append(toEncrypt, sectionPath)
		}
	}
	return cfgText, toEncrypt, nil
}

// isInsidePath reports whether path is a strict sub-path of dir.
func isInsidePath(path, dir string) bool {
	return path != dir && strings.HasPrefix(path+"/", dir+"/")
}

// updateCatalogConfig reads the installed daemon.conf, re-expands the catalog whitelist for every
// matched section (inst.RefreshCatalog) and shows a diff for confirmation; unmatched (user-added)
// sections are kept verbatim. autoConfirm logs the diff and overwrites without prompting. In live
// mode the daemon keeps running (vaults already unlocked and guarded), so sections are re-scanned
// without touching lock state and the caller applies the config via SIGHUP. Returns whether the
// file was written (deliver to the daemon only then).
func updateCatalogConfig(vault *fscrypt.Vault, autoConfirm, live bool) (changed bool, err error) {
	oldText, err := os.ReadFile(systemd.SystemConfigPath)
	if err != nil {
		return false, fmt.Errorf("no previous installation found — run `app-listener install` first: %w", err)
	}

	cfg, err := daemonconfig.Load(systemd.SystemConfigPath)
	if err != nil {
		return false, fmt.Errorf("parsing existing configuration: %w", err)
	}
	if len(cfg.Resources) == 0 {
		return false, fmt.Errorf("existing configuration contains no [watch] sections")
	}

	users, err := inst.ListUsers()
	if err != nil {
		return false, fmt.Errorf("listing users: %w", err)
	}

	confText, changes, err := inst.RefreshCatalog(string(oldText), cfg, users,
		inst.RefreshOptions{Live: live, Scan: vaultScan(vault, live)})
	if err != nil {
		return false, err
	}

	newBytes := []byte(confText)
	if bytes.Equal(oldText, newBytes) {
		log.Info("configuration is already up to date")
		return false, nil
	}
	if !autoConfirm {
		overwrite, overErr := inst.ConfirmOverwrite(systemd.SystemConfigPath, oldText, newBytes)
		if overErr != nil {
			return false, overErr
		}
		if !overwrite {
			log.Info("configuration update aborted by user")
			return false, nil
		}
	} else {
		log.Infof("catalog refresh: %d section(s) changed, diff below", len(changes))
		log.Info(inst.UnifiedDiff(string(oldText), string(newBytes)))
	}

	if err := os.WriteFile(systemd.SystemConfigPath, newBytes, 0o600); err != nil {
		return false, fmt.Errorf("writing updated configuration: %w", err)
	}
	log.Infof("daemon.conf updated (%d section(s) changed)", len(changes))
	return true, nil
}

// vaultScan wraps a section's re-scan in its vault lifecycle. Encrypted resources are unlocked for
// the re-scan under an EPHEMERAL self-only guard attached BEFORE the key is provisioned, so the
// unlock window denies every reader except the root installer (same discipline as the daemon). The
// guard stays attached through the re-lock and drops only once the vault is keyless; a vault that
// can't be re-locked blocks rather than leaving it unlocked.
//
// Live (the daemon runs): the vaults are already unlocked and guarded, so the re-scan touches no
// lock state and needs no ephemeral guard.
func vaultScan(vault *fscrypt.Vault, live bool) inst.SectionScan {
	return func(r *daemonconfig.Resource, expand func() []inst.BinaryRule) ([]inst.BinaryRule, bool, error) {
		// Grouped sections share one encryption root: the vault lifecycle addresses the root.
		root := r.EncryptionRootOrPath()
		if !r.NeedEncryption {
			return expand(), false, nil
		}
		encrypted, err := vault.IsEncrypted(root)
		if err != nil {
			return nil, false, fmt.Errorf("checking encryption of %s: %w", root, err)
		}
		if !encrypted {
			return expand(), false, nil
		}
		if live {
			log.Infof("re-scanning %s (daemon running: vault already unlocked and guarded) ...", root)
			return expand(), true, nil
		}
		log.Infof("unlocking %s for whitelist re-expansion (under an ephemeral guard) ...", root)
		ephemeral, err := unlockUnderGuard(vault, root)
		if err != nil {
			return nil, true, err
		}
		defer ephemeral.Stop()
		fresh := expand()
		log.Infof("re-locking %s ...", root)
		// Ephemeral guard still attached: until the key is gone the tree denies every non-root
		// reader. Retry unbounded like the daemon's lockdown: a pinned fd must not downgrade this
		// to a warning that leaves the vault unlocked after exit.
		lockVaultFully(vault, root, ephemeral)
		return fresh, true, nil
	}
}

// unlockUnderGuard attaches an ephemeral self-only whitelist guard on path BEFORE provisioning the
// key, so the unlock window denies every reader except the root installer (GUARD_ALLOW_ROOT, same
// uid-gated mechanism as the daemon). The caller must Stop the guard only after the vault is
// confirmed locked (lockVaultFully).
//
// A single-file resource (userspace file vault) is transformed IN PLACE (the unlock truncates and
// rewrites its own bytes), which the read-only self mask denies. As in the daemon
// (vaultOpForGuard), write is granted only for the vault call via WithSelfVaultAccess, and the
// recovery sidecar is pre-staged before the guard attaches (creating it beside a guarded file is
// denied once live). Without both, `install --update-catalog-only` failed with "truncate ...:
// operation not permitted" on file-vault sections (e.g. Steam's registry.vdf).
func unlockUnderGuard(vault *fscrypt.Vault, path string) (*guard.Guard, error) {
	self, err := ebpf.ComputeBinaryEntry("/proc/self/exe")
	if err != nil {
		return nil, fmt.Errorf("resolving installer executable: %w", err)
	}
	if err := fscrypt.EnsureRecoverySidecarPlaceholder(path); err != nil {
		log.Warnf("install: could not pre-create the recovery sidecar for %s (%v) — "+
			"a lock/unlock interrupted by a crash may not be recoverable", path, err)
	}
	ephemeral, guardErr := guard.NewGuard(path, guard.ModeWhitelist, nil, true, 0,
		guard.WithSelfAllowBinary(self, []ebpf.EventType{ebpf.EventOpen, ebpf.EventRead}))
	if guardErr != nil {
		// Fail-closed: without the guard there is no unlock at all.
		return nil, fmt.Errorf("attaching ephemeral guard for %s: %w (vault left locked)", path, guardErr)
	}
	if unlockErr := vaultOpUnderGuard(ephemeral, path, func() error { return vault.Unlock(path) }); unlockErr != nil {
		ephemeral.Stop()
		return nil, fmt.Errorf("unlocking %s: %w", path, unlockErr)
	}
	return ephemeral, nil
}

// vaultOpUnderGuard runs a vault Unlock/Lock on path, widening g's self mask to write only for a
// file vault (a directory's lifecycle is pure keyring work). g may be nil: op still runs so a
// denial surfaces as an error, not a silent skip.
func vaultOpUnderGuard(g *guard.Guard, path string, op func() error) error {
	info, statErr := os.Stat(path)
	if statErr != nil || !info.Mode().IsRegular() || g == nil {
		return op()
	}
	return g.WithSelfVaultAccess(op)
}

// lockVaultFully force-flushes the vault key until it is gone, never giving up (like the daemon's
// lockdown): the caller keeps the ephemeral guard attached, so the tree stays guarded while this
// retries. A persistent pin blocks the refresh with a loud log rather than leaving the vault
// unlocked.
func lockVaultFully(vault *fscrypt.Vault, path string, g *guard.Guard) {
	for {
		// Widened per attempt, never across the retry sleep.
		err := vaultOpUnderGuard(g, path, func() error { return vault.Lock(path, true) })
		if err == nil {
			return
		}
		if errors.Is(err, repository.ErrKeyMissing) {
			return // fully locked
		}
		log.Errorf("install: %s is still unlocked: a process holds open files in it "+
			"(investigate with: lsof +D %s, fuser -v %s): %v — retrying while the ephemeral guard keeps it protected",
			path, path, path, err)
		time.Sleep(time.Second)
	}
}
