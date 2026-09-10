package install

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/repository"
	"github.com/Virgula0/app-listener/internal/systemd"
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

// catalogGroup bundles every discovered path of one resource (one catalog
// Name, one user) behind a single TUI checkbox — a resource stored in more
// than one location (Steam's data dir + legacy home) is selected as a unit.
type catalogGroup struct {
	label      string
	candidates []inst.Candidate
}

// groupCandidates collapses the per-path candidates into per-resource
// groups, preserving first-seen order.
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

// pickDirectories probes the catalog for every selected user (and the
// system-level entries) and asks which of the found critical directories
// to protect. All are preselected.
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

// pickFromCandidates groups the candidates per resource (one checkbox per
// catalog Name + user) and runs the preselected multi-select, returning the
// flattened per-path candidates of every picked group.
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
	if err := form.WithKeyMap(selectionKeymap()).Run(); err != nil {
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

// sectionsFromCandidates turns discovered/added candidates into config
// sections: the filtered whitelist, encryption on by default, and any
// grouped extra watch sub-paths from the catalog entry.
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

// editConfig renders the config from the selected candidates and opens the
// embedded editor.
func editConfig(candidates []inst.Candidate) (string, *daemonconfig.Config, error) {
	return runConfigEditor(
		"app-listener daemon.conf — review and save (Ctrl+S)",
		inst.GenerateConf(sectionsFromCandidates(candidates)))
}

// runConfigEditor opens initial in the embedded editor and validates the
// result through the same strict parser the daemon uses; an invalid config
// re-opens the editor until it parses or the user aborts (Esc).
func runConfigEditor(title, initial string) (string, *daemonconfig.Config, error) {
	confText := initial
	for {
		edited, err := inst.EditText(title, confText)
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
// askEncryption asks the fscrypt question once per ENCRYPTION GROUP (not
// per watch path): grouped sections share the group's vault, so the answer
// and the in-place encryption target the group root. Regression guard for
// the grouped-watch config: addressing sections by watch sub-path failed
// with "section not found" (a grouped config has one section header per
// group, not one per watch path).
func askEncryption(vault *fscrypt.Vault, cfgText string, cfg *daemonconfig.Config) (text string, toEncrypt []string, err error) {
	for _, r := range cfg.EncryptionGroups() {
		// Grouped sections share ONE [watch <root>] header: the fscrypt
		// lifecycle and every config-text patch address the encryption
		// root, never a watch sub-path (which has no section header of its
		// own — SetNeedEncryption would fail "section not found").
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

// resolveCatalogEntry performs a reverse lookup: given an absolute path from
// the existing daemon.conf, finds which Catalog entry it originated from.
// System-level entries (AbsPaths) are matched by exact path. User-level
// entries (RelPaths) are matched by computing PathsFor for each known user —
// any of an entry's locations matches its shared whitelist. Grouped entries
// (WatchRelPaths) also match their watch sub-paths — e.g. ~/.config/discord/
// Local Storage resolves to the Discord entry — so a refresh re-expands the
// whitelist of every grouped section. Returns nil when the section was
// user-added and has no catalog origin.
func resolveCatalogEntry(resourcePath string, users []inst.User) (*inst.CandidateDir, *inst.User) {
	if entry, user := findCatalogRoot(resourcePath, users); entry != nil {
		return entry, user
	}
	return findCatalogWatchSubPath(resourcePath, users)
}

// catalogEntryMatch returns the user whose expansion of entry contains a
// path satisfying match, or nil. System entries yield the zero user.
func catalogEntryMatch(entry *inst.CandidateDir, users []inst.User, match func(catalogPath string) bool) *inst.User {
	if entry.IsSystem() {
		for _, p := range entry.PathsFor("", "") {
			if match(p) {
				return &inst.User{}
			}
		}
		return nil
	}
	for j := range users {
		u := &users[j]
		for _, p := range entry.PathsFor(u.Home, u.Name) {
			if match(p) {
				return u
			}
		}
	}
	return nil
}

// findCatalogRoot matches one of the entry's own watch roots exactly.
func findCatalogRoot(resourcePath string, users []inst.User) (*inst.CandidateDir, *inst.User) {
	for i := range inst.Catalog {
		entry := &inst.Catalog[i]
		if u := catalogEntryMatch(entry, users, func(p string) bool { return p == resourcePath }); u != nil {
			return entry, u
		}
	}
	return nil, nil
}

// findCatalogWatchSubPath matches grouped watch sub-paths: a resource path
// strictly inside a catalog root resolves to that entry (the `watch:` group),
// so refreshes re-expand the whitelist of every grouped section.
func findCatalogWatchSubPath(resourcePath string, users []inst.User) (*inst.CandidateDir, *inst.User) {
	for i := range inst.Catalog {
		entry := &inst.Catalog[i]
		if u := catalogEntryMatch(entry, users, func(p string) bool { return isInsidePath(resourcePath, p) }); u != nil {
			return entry, u
		}
	}
	return nil, nil
}

// isInsidePath reports whether path is a strict sub-path of dir.
func isInsidePath(path, dir string) bool {
	return path != dir && strings.HasPrefix(path+"/", dir+"/")
}

// errNoCatalogMatch means the installed config has no section that maps to a
// catalog entry (an all-manual config). Interactive callers surface it as an
// error; automated ones (`--yes`, the package hooks and the boot unit) treat
// it as "nothing to refresh".
var errNoCatalogMatch = errors.New("no catalog entries match any configured directory")

// updateCatalogConfig reads the installed daemon.conf, re-expands the
// catalog whitelist for every matched section and shows a diff for
// confirmation. Unmatched (user-added) sections are preserved verbatim.
// When autoConfirm is true the diff is logged and the config is overwritten
// without prompting. In live mode the daemon keeps running: its vaults are
// already unlocked and guarded, so sections are re-scanned without touching
// lock state, and the caller applies the patched config via SIGHUP reload.
// Returns whether the config file was actually written (the caller must
// deliver the change to the daemon only then).
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

	confText := string(oldText)
	patched := 0
	// Grouped configs: iterate ENCRYPTION GROUPS, not resources — a grouped
	// section is one [watch <root>] header shared by all its watch paths,
	// so the whitelist is re-expanded and written exactly once per group.
	// Regression guard for the first grouped implementation: patching by
	// watch sub-path failed with "section not found in configuration".
	for _, r := range cfg.EncryptionGroups() {
		text, ok, patchErr := patchCatalogSection(vault, confText, r, users, live)
		if patchErr != nil {
			return false, patchErr
		}
		if ok {
			confText = text
			patched++
		}
	}

	if patched == 0 {
		return false, errNoCatalogMatch
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
		log.Infof("catalog refresh: %d sections re-scanned, diff below", patched)
		log.Info(inst.UnifiedDiff(string(oldText), string(newBytes)))
	}

	if err := os.WriteFile(systemd.SystemConfigPath, newBytes, 0o600); err != nil {
		return false, fmt.Errorf("writing updated configuration: %w", err)
	}
	log.Infof("daemon.conf updated (%d sections re-scanned)", patched)
	return true, nil
}

// patchCatalogSection re-expands the whitelist for one config resource
// if it matches a catalog entry. Returns the updated text and true when
// patched, or the original text and false when the section is user-added.
//
// Encrypted resources are unlocked for the re-scan under an EPHEMERAL
// self-only guard attached BEFORE the key is provisioned: the unlock window
// then denies every reader except the root installer process — the same
// protection discipline as the running daemon — instead of leaving the tree
// readable with no LSM attached. The guard stays attached through the
// re-lock and is dropped only once the vault is keyless; a vault that
// cannot be re-locked is a hard error, so the hook fails visibly and the
// config write / daemon restart never proceed on an unresolved vault.
//
// In live mode (live=true) the daemon keeps running: its vaults are already
// unlocked and guarded, so the re-scan touches no lock state at all and no
// ephemeral guard is needed (the running guards + the installer's
// GUARD_ALLOW_ROOT identity already protect the tree). A resource whose
// fresh whitelist comes back EMPTY while the section previously had entries
// is a hard error: that state means the vault is locked or this installer
// is not the running daemon's binary — persisting it would silently shrink
// the whitelist.
func patchCatalogSection(vault *fscrypt.Vault, confText string, r *daemonconfig.Resource, users []inst.User, live bool) (updated string, patched bool, err error) {
	// The SECTION path is what the text helpers address: the group root for
	// grouped sections (`watch:` directives — one [watch <root>] header
	// shared by all watch paths), the resource path otherwise.
	sectionPath := r.EncryptionRootOrPath()

	entry, user := resolveCatalogEntry(sectionPath, users)
	if entry == nil {
		log.Infof("keeping user section as-is: %s", sectionPath)
		return confText, false, nil
	}

	// Grouped sections (watch: directives) share one encryption root AND one
	// [watch <root>] header carrying one shared whitelist: the vault-level
	// lifecycle (ephemeral-guarded unlock, re-lock) and the whitelist text
	// patch both address that root, never a watch sub-path.
	root := sectionPath

	wasEncrypted := false
	if r.NeedEncryption {
		encrypted, encErr := vault.IsEncrypted(root)
		if encErr != nil {
			return "", false, fmt.Errorf("checking encryption of %s: %w", root, encErr)
		}
		if encrypted {
			wasEncrypted = true
			if live {
				log.Infof("re-scanning %s (daemon running: vault already unlocked and guarded) ...", sectionPath)
			} else {
				log.Infof("unlocking %s for whitelist re-expansion (under an ephemeral guard) ...", root)
				release, unlockErr := unlockUnderGuard(vault, root)
				if unlockErr != nil {
					return "", false, unlockErr
				}
				defer release()
			}
		}
	}

	candidate := inst.Candidate{User: *user, Entry: *entry, Path: sectionPath}
	freshWhitelist := candidate.FilterExistingWhitelist()
	log.Infof("re-scanned %s (%s) — %d whitelisted binaries", sectionPath, entry.Name, len(freshWhitelist))

	if liveEmptyWhitelistRejected(wasEncrypted, live, len(freshWhitelist), len(parseSectionWhitelist(confText, sectionPath))) {
		return "", false, fmt.Errorf("live re-scan of %s produced an empty whitelist while the config lists binaries for it: "+
			"the vault appears locked or this installer is not the running daemon's binary — refusing to shrink the whitelist", sectionPath)
	}

	if wasEncrypted && !live {
		log.Infof("re-locking %s ...", root)
		// The ephemeral guard is still attached: until the key is gone the
		// tree keeps denying every non-root reader. Retry unbounded, like
		// the daemon's lockdown — a pinned fd must not downgrade this to a
		// warning that leaves the vault unlocked once the process exits.
		lockVaultFully(vault, root)
	}

	updated, patchErr := inst.SetSectionWhitelist(confText, sectionPath, freshWhitelist)
	if patchErr != nil {
		return "", false, fmt.Errorf("patching section %s: %w", sectionPath, patchErr)
	}
	return updated, true, nil
}

// parseSectionWhitelist extracts the binary paths currently listed in the
// [watch <path>] section of confText, for the live empty-whitelist safety
// check. Returns nil when the section is absent.
func parseSectionWhitelist(confText, resourcePath string) []string {
	needle := "[watch " + resourcePath + "]"
	var out []string
	inSection := false
	for _, line := range strings.Split(confText, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[watch") {
			inSection = trimmed == needle
			continue
		}
		if !inSection || trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if trimmed == "need_encryption: true" || trimmed == "need_encryption: false" {
			continue
		}
		if fields := strings.Fields(trimmed); len(fields) > 0 {
			out = append(out, fields[0])
		}
	}
	return out
}

// unlockUnderGuard attaches an ephemeral self-only whitelist guard on path
// BEFORE provisioning the key, so the unlock window denies every reader
// except the root installer process (GUARD_ALLOW_ROOT — the same uid-gated
// mechanism the daemon uses for itself). The returned release func stops the
// guard and must run only after the vault is confirmed locked back.
func unlockUnderGuard(vault *fscrypt.Vault, path string) (release func(), err error) {
	self, err := ebpf.ComputeBinaryEntry("/proc/self/exe")
	if err != nil {
		return nil, fmt.Errorf("resolving installer executable: %w", err)
	}
	ephemeral, guardErr := guard.NewGuard(path, guard.ModeWhitelist, nil, true, 0,
		guard.WithSelfAllowBinary(self, []ebpf.EventType{ebpf.EventOpen, ebpf.EventRead}))
	if guardErr != nil {
		// Fail-closed: without the guard there is no unlock at all.
		return nil, fmt.Errorf("attaching ephemeral guard for %s: %w (vault left locked)", path, guardErr)
	}
	if unlockErr := vault.Unlock(path); unlockErr != nil {
		ephemeral.Stop()
		return nil, fmt.Errorf("unlocking %s: %w", path, unlockErr)
	}
	return ephemeral.Stop, nil
}

// lockVaultFully force-flushes the vault key until it is gone, with the same
// never-give-up discipline as the daemon's lockdown: the caller keeps the
// ephemeral guard attached, so the tree stays guarded for as long as this
// retry loop runs. A persistent pin blocks the pacman hook with a loud log
// instead of silently leaving the vault unlocked.
// liveEmptyWhitelistRejected is the live refresh's fail-closed contract: a
// re-scan that comes back empty for a previously-populated encrypted
// resource must be refused. That state means the vault is locked (stat on
// plaintext names fails) or this installer binary is not the running
// daemon's (reads denied) — persisting it would silently shrink the
// whitelist. Non-encrypted resources are exempt: an empty re-scan there is
// a legitimately uninstalled binary ("empty whitelists still deny
// everything").
func liveEmptyWhitelistRejected(encrypted, live bool, fresh, old int) bool {
	return live && encrypted && fresh == 0 && old > 0
}

func lockVaultFully(vault *fscrypt.Vault, path string) {
	for {
		err := vault.Lock(path, true)
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
