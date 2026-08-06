package install

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	inst "github.com/Virgula0/app-listener/internal/install"
)

// cleanOrphanedFscrypt removes the fscrypt policy and protector metadata left
// behind by directories the installer encrypted and that no longer exist.
// The live set is derived from the catalog, exactly as the installer scopes
// itself: every watch directory discoverable for every local user, the
// system-level entries, and the existing resources of the final config. Only
// metadata of directories that are gone is deleted — an existing encrypted
// directory keeps its metadata even when it is not part of the config, and
// login protectors or foreign fscrypt setups are never touched. It runs after
// deployment as the second-to-last step.
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

	live, err := fscrypt.LivePolicyDescriptors(paths)
	if err != nil {
		return err
	}

	anchor := "/"
	for _, p := range paths {
		if _, statErr := os.Stat(p); statErr == nil {
			anchor = p
			break
		}
	}
	removedProtectors, removedPolicies, err := fscrypt.CleanOrphans(anchor, live)
	if err != nil {
		return fmt.Errorf("cleaning orphaned fscrypt metadata: %w", err)
	}
	if removedProtectors == 0 && removedPolicies == 0 {
		log.Info("no orphaned fscrypt metadata to clean")
		return nil
	}
	log.Infof("cleaned orphaned fscrypt metadata: removed %d protector(s) and %d policy(ies)", removedProtectors, removedPolicies)
	return nil
}
