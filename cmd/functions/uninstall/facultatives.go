package uninstall

import (
	"fmt"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/systemd"
	"github.com/Virgula0/app-listener/internal/wizard"
)

var facultativesCmd = &cobra.Command{
	Use:   "facultatives",
	Short: "Remove the optional per-user setup you pick (ssh-agent, Bun redirect)",
	Long: `List the optional per-user setup the installer added (the ssh-agent unit and its
SSH_AUTH_SOCK block, the Bun launcher wrappers) and remove only the ones you
pick. Nothing else changes and the daemon may keep running.`,
	Args: cobra.NoArgs,
	RunE: func(*cobra.Command, []string) error { return run(Request{Facultatives: true}) },
}

func init() { UninstallCmd.AddCommand(facultativesCmd) }

// installedFacultative is one facultative `install facultatives` (or the wizard) set up for user.
type installedFacultative struct {
	user inst.User
	bun  bool
	// unit: the ssh-agent unit matches the bundled sample, so it is ours to remove.
	unit bool
	// sshGuarded: daemon.conf covers ~/.ssh, whose AddKeysToAgent block is then left in place.
	sshGuarded bool
}

func (f *installedFacultative) label() string {
	if f.bun {
		return fmt.Sprintf("Bun launcher wrappers for %s (%s is kept)", f.user.Name, inst.BunTmpDir(f.user.Home))
	}
	s := fmt.Sprintf("ssh-agent for %s: SSH_AUTH_SOCK block", f.user.Name)
	if f.unit {
		s += ", user unit " + inst.SSHAgentUnitPath(f.user.Home)
	}
	if !f.sshGuarded {
		s += ", AddKeysToAgent block in ~/.ssh/config"
	}
	return s
}

// removeFacultatives is `uninstall facultatives`.
func removeFacultatives() error {
	cfg, err := inst.LoadConfigIfPresent(systemd.SystemConfigPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", systemd.SystemConfigPath, err)
	}
	users, err := inst.ListUsers()
	if err != nil {
		return err
	}
	sample, err := inst.SampleContent("ssh-agent.service")
	if err != nil {
		return err
	}
	covered := inst.CoveredPaths(cfg)
	var items []installedFacultative
	for i := range users {
		u := users[i]
		unit := u.UID != 0 && isInstallerSSHAgentUnit(sample, inst.SSHAgentUnitPath(u.Home))
		if unit || (u.UID != 0 && inst.HasSSHAgentEnv(u)) {
			items = append(items, installedFacultative{user: u, unit: unit,
				sshGuarded: inst.PathCovered(filepath.Join(u.Home, ".ssh"), covered)})
		}
		if inst.HasBunLauncherEnv(u) {
			items = append(items, installedFacultative{user: u, bun: true})
		}
	}
	if len(items) == 0 {
		log.Info("no installed facultative found")
		return nil
	}
	choices := make([]wizard.Choice, len(items))
	for i := range items {
		choices[i] = wizard.Choice{Label: items[i].label()}
	}
	idx, err := wizard.MultiSelect("Facultatives: select the ones to remove",
		"Only the selected ones are removed; the daemon keeps running.", choices, 12)
	if err != nil {
		return err
	}
	for _, i := range idx {
		if err := removeFacultative(&items[i]); err != nil {
			return err
		}
	}
	log.Infof("removed %d facultative(s)", len(idx))
	return nil
}

func removeFacultative(f *installedFacultative) error {
	if f.bun {
		return removeBunLauncherEnv(f.user)
	}
	if f.unit {
		if err := removeSSHAgentUnit(userSSHAgentUnit{User: f.user, Path: inst.SSHAgentUnitPath(f.user.Home)}); err != nil {
			return err
		}
		log.Infof("an ssh-agent already running for %s keeps running until logout (systemctl --user stop ssh-agent)",
			f.user.Name)
	}
	return removeSSHAgentEnv(f.user, f.sshGuarded)
}
