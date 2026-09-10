// Package backups handles the .app_listener.backup migration copies the
// installer leaves behind after an fscrypt migration: discovering them across
// the catalog and the installed daemon.conf, an interactive multi-select, and
// deleting or restoring them over the fscrypt vault. Shared by
// `install --restore-backups` / `--delete-post-backups`, the install
// post-step, and `uninstall`.
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

// Find lists every catalog or daemon.conf path that still carries a
// .app_listener.backup directory, deduplicated by path. The set mirrors the
// installer's own scope: the resources of /etc/app-listener/daemon.conf when
// it exists (covering manually added directories) plus every catalog path
// discovered for all local users and the system entries.
func Find() ([]Backup, error) {
	seen := make(map[string]bool)
	var paths []string
	add := func(p string) {
		if p != "" && !seen[p] {
			seen[p] = true
			paths = append(paths, p)
		}
	}

	if _, err := os.Stat(systemd.SystemConfigPath); err == nil {
		cfg, loadErr := daemonconfig.Load(systemd.SystemConfigPath)
		if loadErr != nil {
			return nil, fmt.Errorf("reading %s: %w", systemd.SystemConfigPath, loadErr)
		}
		for i := range cfg.Resources {
			// Grouped sections: the backup is at the encryption root (the whole
			// vault is renamed aside during migration), at the watch path
			// itself otherwise — keep both.
			add(cfg.Resources[i].Path)
			add(cfg.Resources[i].EncryptionRootOrPath())
		}
	}

	users, err := install.ListUsers()
	if err != nil {
		return nil, err
	}
	found := install.DiscoverForUsers(users)
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

// Restore deletes each entry's live (encrypted) directory and moves its
// backup back to the original path, behind a progress indicator. The caller
// MUST have stopped the daemon first — restoring under a live daemon would
// let it keep using the directories being deleted.
func Restore(entries []Backup) error {
	vault := fscrypt.New()
	return withProgress("Restoring", entries, func(b Backup) error {
		if err := vault.RestoreBackup(b.Path); err != nil {
			return fmt.Errorf("fatal: %v", err)
		}
		return nil
	})
}

// withProgress runs op for every entry while a single-line progress bar at
// the bottom of the terminal shows progress (one step per directory); the
// fscrypt logs scroll normally above it.
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
