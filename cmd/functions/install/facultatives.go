package install

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/systemd"
	"github.com/Virgula0/app-listener/internal/wizard"
)

type facultativeKind int

const (
	facultativeSSHAgent facultativeKind = iota
	facultativeBunTmpdir
)

// facultative is one optional per-user setup the wizard asks about. `install facultatives` offers
// it whether or not daemon.conf guards the directory it serves.
type facultative struct {
	kind facultativeKind
	user inst.User
	// launchers are the Bun commands to wrap (facultativeBunTmpdir).
	launchers []string
	// sshGuarded: daemon.conf covers ~/.ssh, so its config is left alone: the daemon denies the
	// edit while it runs, and a stopped daemon leaves the vault locked (ENOKEY).
	sshGuarded bool
	installed  bool
}

func (f *facultative) label() string {
	var s string
	switch f.kind {
	case facultativeSSHAgent:
		s = fmt.Sprintf("ssh-agent for %s: user unit %s, SSH_AUTH_SOCK in the shell startup file",
			f.user.Name, inst.SSHAgentUnitPath(f.user.Home))
		if f.sshGuarded {
			s += " (~/.ssh is guarded: its config is left as is)"
		} else {
			s += ", AddKeysToAgent in ~/.ssh/config"
		}
	case facultativeBunTmpdir:
		s = fmt.Sprintf("Bun $TMPDIR redirect for %s: %s -> %s",
			f.user.Name, strings.Join(f.launchers, ", "), inst.BunTmpDir(f.user.Home))
	}
	if f.installed {
		s += "  [installed]"
	}
	return s
}

var facultativesCmd = &cobra.Command{
	Use:   "facultatives",
	Short: "Install the optional per-user setup you pick (ssh-agent, Bun redirect)",
	Long: `List the optional per-user setup the wizard offers (the ssh-agent unit, the Bun
$TMPDIR redirect) for every user, whether or not daemon.conf guards what it
serves, and install only the ones you pick. Nothing else changes and the daemon
keeps running: a new Bun redirect is applied with a reload.`,
	Args: cobra.NoArgs,
	RunE: func(*cobra.Command, []string) error { return asRoot(runFacultatives) },
}

// runFacultatives is `install facultatives`.
func runFacultatives() error {
	cfg, err := inst.LoadConfigIfPresent(systemd.SystemConfigPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", systemd.SystemConfigPath, err)
	}
	users, err := inst.ListUsers()
	if err != nil {
		return err
	}
	items, err := listFacultatives(cfg, users)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		log.Info("no facultative applies to this host's users")
		return nil
	}
	choices := make([]wizard.Choice, len(items))
	for i := range items {
		choices[i] = wizard.Choice{Label: items[i].label()}
	}
	idx, err := wizard.MultiSelect("Facultatives: select the ones to install",
		"Optional per-user setup, offered whether or not daemon.conf guards the directory it serves. "+
			"Only the selected ones are installed; the daemon keeps running.", choices, 12)
	if err != nil {
		return err
	}
	if len(idx) == 0 {
		log.Info("no facultative selected: nothing installed")
		return nil
	}
	picked := make([]facultative, 0, len(idx))
	for _, i := range idx {
		picked = append(picked, items[i])
	}
	return installFacultatives(cfg, picked)
}

// listFacultatives returns the ssh-agent setup for every non-root user and the Bun redirect for
// every user with a Bun catalog app configured or present.
func listFacultatives(cfg *daemonconfig.Config, users []inst.User) ([]facultative, error) {
	sample, err := inst.SampleContent("ssh-agent.service")
	if err != nil {
		return nil, err
	}
	covered := inst.CoveredPaths(cfg)
	var out []facultative
	for i := range users {
		u := users[i]
		if u.UID != 0 {
			unit, _ := os.ReadFile(inst.SSHAgentUnitPath(u.Home))
			out = append(out, facultative{kind: facultativeSSHAgent, user: u,
				sshGuarded: inst.PathCovered(filepath.Join(u.Home, ".ssh"), covered),
				installed:  bytes.Equal(unit, sample) && inst.HasSSHAgentEnv(u)})
		}
		if launchers := bunLaunchersOffered(cfg, u); len(launchers) > 0 {
			_, dirErr := os.Stat(inst.BunTmpDir(u.Home))
			out = append(out, facultative{kind: facultativeBunTmpdir, user: u, launchers: launchers,
				installed: dirErr == nil && inst.HasBunLauncherEnv(u)})
		}
	}
	return out, nil
}

// bunLaunchersOffered is bunLaunchersFor widened to the Bun catalog apps present in u's home.
func bunLaunchersOffered(cfg *daemonconfig.Config, u inst.User) []string {
	found := map[string]bool{}
	cands := inst.Discover(u)
	for i := range cands {
		found[cands[i].Entry.Name] = true
	}
	return bunLaunchers(func(e *inst.CandidateDir) bool { return found[e.Name] || bunEntryConfigured(cfg, e, u) })
}

// installFacultatives sets up each picked facultative like the wizard does. The daemon reserves a
// Bun tmp dir only when it builds its trust set, so a new dir a configured Bun app uses is applied
// with a reload (guards stay attached, no unlock churn); an unreserved one only keeps the load denied.
func installFacultatives(cfg *daemonconfig.Config, picked []facultative) error {
	var bun []bunUser
	reload := false
	for i := range picked {
		f := &picked[i]
		switch f.kind {
		case facultativeSSHAgent:
			if err := installSSHAgentFacultative(f); err != nil {
				return err
			}
		case facultativeBunTmpdir:
			bun = append(bun, bunUser{User: f.user, Launchers: f.launchers})
			reload = reload || len(bunLaunchersFor(cfg, f.user)) > 0
		}
	}
	if err := setupBunTmpdir(bun); err != nil {
		return err
	}
	if reload && systemd.IsDaemonActive() {
		log.Info("reloading the daemon so it reserves the new Bun tmp dir (no restart)")
		if err := systemd.ReloadDaemonIfActive(); err != nil {
			return fmt.Errorf("the Bun redirect is set up, but the daemon reload failed (the Bun app's "+
				"native library load stays denied until the next reload): %w", err)
		}
	}
	log.Infof("installed %d facultative(s)", len(picked))
	return nil
}

func installSSHAgentFacultative(f *facultative) error {
	if f.sshGuarded {
		log.Warnf("%s is guarded: AddKeysToAgent was not added to its config — add `AddKeysToAgent yes` "+
			"with `app-listener edit-protected` for ssh to load keys into the agent on first use",
			filepath.Join(f.user.Home, ".ssh"))
	} else if err := addKeysToAgent([]inst.User{f.user}); err != nil {
		return err
	}
	return installSSHAgent(f.user)
}
