package install

import (
	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	inst "github.com/Virgula0/app-listener/internal/install"
)

// cleanOrphanedFscrypt removes the fscrypt policy and protector metadata
// left behind by directories that no longer exist. The live set mirrors the
// installer's own scope: every discoverable catalog directory, the system
// entries, and the final config resources. Runs after deployment as the
// second-to-last step.
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
	for _, r := range cfg.Resources {
		paths = append(paths, r.Path)
	}

	if err := fscrypt.CleanOrphanedMetadata(paths); err != nil {
		return err
	}
	return nil
}
