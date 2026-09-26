// Package backups handles the .app_listener.backup copies left after an fscrypt migration:
// discovery across the catalog and installed daemon.conf, interactive multi-select, delete/restore
// over the vault. Shared by `install --restore-backups`/`--delete-post-backups`, the install
// post-step and `uninstall`.
package backups

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	"github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/systemd"
	"github.com/Virgula0/app-listener/internal/wizard"
)

// Backup is one discovered migration backup and the directory it belongs to.
type Backup struct {
	// Path is the original directory (its encrypted replacement is live).
	Path string
	// BackupPath is Path + fscrypt.BackupSuffix — the backup directory.
	BackupPath string
}

// Find lists every catalog or daemon.conf path still carrying a .app_listener.backup, deduplicated.
// Mirrors the installer's scope: resources of /etc/app-listener/daemon.conf (incl. manual dirs)
// plus every catalog path for all local users and system entries.
func Find() ([]Backup, error) {
	return find(systemd.SystemConfigPath, install.ListUsers, install.DiscoverForUsers)
}

// find is Find's injectable core (systemConfigPath, listUsers, discoverForUsers as parameters) so
// tests use fake data: on a host with the daemon installed, some catalog paths (e.g. Steam's
// registry.vdf) are live guarded resources, and os.Lstat from outside their whitelist (the test
// binary) trips the guard.
func find(systemConfigPath string, listUsers func() ([]install.User, error), discoverForUsers func([]install.User) []install.Candidate) ([]Backup, error) {
	seen := make(map[string]bool)
	var paths []string
	add := func(p string) {
		if p != "" && !seen[p] {
			seen[p] = true
			paths = append(paths, p)
		}
	}

	if _, err := os.Stat(systemConfigPath); err == nil {
		cfg, loadErr := daemonconfig.Load(systemConfigPath)
		if loadErr != nil {
			return nil, fmt.Errorf("reading %s: %w", systemConfigPath, loadErr)
		}
		for i := range cfg.Resources {
			// Grouped sections: the backup is at the encryption root (the whole vault is renamed
			// aside), else at the watch path itself; keep both.
			add(cfg.Resources[i].Path)
			add(cfg.Resources[i].EncryptionRootOrPath())
		}
	}

	users, err := listUsers()
	if err != nil {
		return nil, err
	}
	found := discoverForUsers(users)
	for i := range found {
		add(found[i].Path)
	}

	var out []Backup
	for _, p := range paths {
		bp := p + fscrypt.BackupSuffix
		if _, err := os.Lstat(bp); err == nil {
			out = append(out, Backup{Path: p, BackupPath: bp})
		}
	}
	return out, nil
}

// Select shows the backups in a preselected TUI multi-select and returns the
// chosen ones.
func Select(entries []Backup, title, description string) ([]Backup, error) {
	opts := make([]huh.Option[int], 0, len(entries))
	for i := range entries {
		opts = append(opts, huh.NewOption(entries[i].BackupPath, i).Selected(true))
	}
	var pickedIdx []int
	form := huh.NewForm(huh.NewGroup(
		huh.NewMultiSelect[int]().
			Title(title).
			Description(description).
			Options(opts...).
			Height(10).
			Value(&pickedIdx),
	))
	if err := form.WithKeyMap(wizard.MultiSelectKeymap()).Run(); err != nil {
		return nil, err
	}
	picked := make([]Backup, 0, len(pickedIdx))
	for _, i := range pickedIdx {
		picked = append(picked, entries[i])
	}
	return picked, nil
}

// Delete removes the backup directories of the given entries behind a
// bottom-bar progress indicator.
func Delete(entries []Backup) error {
	return withProgress("Deleting", entries, func(b Backup) error {
		if err := os.RemoveAll(b.BackupPath); err != nil {
			return fmt.Errorf("removing backup %s: %w", b.BackupPath, err)
		}
		return nil
	})
}

// Restore deletes each entry's live (encrypted) directory and moves its backup back, behind a
// progress indicator. The caller MUST have stopped the daemon first (it would keep using the
// directories being deleted).
func Restore(entries []Backup) error {
	vault := fscrypt.New()
	return withProgress("Restoring", entries, func(b Backup) error {
		if err := vault.RestoreBackup(b.Path); err != nil {
			return fmt.Errorf("fatal: %v", err)
		}
		return nil
	})
}

// withProgress runs op per entry under a one-line progress bar at the terminal bottom (one step per
// directory); fscrypt logs scroll above it.
func withProgress(verb string, entries []Backup, op func(Backup) error) error {
	total := len(entries)
	return wizard.WithBottomBar(func(bar *wizard.BottomBar) error {
		for i := range entries {
			label := fmt.Sprintf("%s %d/%d: %s", verb, i+1, total, entries[i].BackupPath)
			bar.Set(label, float64(i)/float64(total))
			if err := op(entries[i]); err != nil {
				return err
			}
			bar.Set(label+" (done)", float64(i+1)/float64(total))
		}
		return nil
	})
}
