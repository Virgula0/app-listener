package install

import (
	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	inst "github.com/Virgula0/app-listener/internal/install"
)

// cleanOrphanedFscrypt removes fscrypt policy/protector metadata left by directories that no longer
// exist. The live set mirrors the installer's scope: every discoverable catalog dir, system entries
// and the final config resources. Runs after deployment, second-to-last.
func cleanOrphanedFscrypt(cfg *daemonconfig.Config) error {
	users, err := inst.ListUsers()
	if err != nil {
		return err
	}
	var paths []string
	catalog := inst.DiscoverForUsers(users)
	for i := range catalog {
		paths = append(paths, catalog[i].Path)
	}
	system := inst.DiscoverSystem()
	for i := range system {
		paths = append(paths, system[i].Path)
	}
	for i := range cfg.Resources {
		r := &cfg.Resources[i]
		// Keep both the guarded watch path and its encryption root: grouped sections carry the
		// policy on the root, which must not be classified as orphaned.
		paths = append(paths, r.Path, r.EncryptionRootOrPath())
	}

	if err := fscrypt.CleanOrphanedMetadata(paths); err != nil {
		return err
	}
	return nil
}
