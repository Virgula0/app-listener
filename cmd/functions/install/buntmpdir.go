package install

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	inst "github.com/Virgula0/app-listener/internal/install"
)

// bunUser is one user who opted into the Bun $TMPDIR redirect, with the launcher command names to
// wrap (union of every configured Bun catalog app's launchers).
type bunUser struct {
	User      inst.User
	Launchers []string
}

// gatherPerUserSetup runs the interactive per-user prompts before encryption: the ssh-agent setup
// (with its ~/.ssh/config edit, which must land in the tree that gets migrated) and the Bun $TMPDIR
// redirect.
func gatherPerUserSetup(cfg *daemonconfig.Config) ([]inst.User, []bunUser, error) {
	sshUsers, err := askSSHAgentUsers(cfg)
	if err != nil {
		return nil, nil, err
	}
	if keysErr := addKeysToAgent(sshUsers); keysErr != nil {
		return nil, nil, keysErr
	}
	bunUsers, err := askBunTmpdirUsers(cfg)
	if err != nil {
		return nil, nil, err
	}
	return sshUsers, bunUsers, nil
}

// askBunTmpdirUsers returns the users to set the Bun launcher redirect up for: those with a
// configured Bun catalog app, after one question each. The redirect points the app's $TMPDIR at a
// private guarded dir so its per-launch native library (.bun-<uid>-<hash>.so) loads under trust
// instead of being denied on world-writable /tmp.
func askBunTmpdirUsers(cfg *daemonconfig.Config) ([]bunUser, error) {
	users, err := inst.ListUsers()
	if err != nil {
		return nil, err
	}
	var accepted []bunUser
	for i := range users {
		u := users[i]
		launchers := bunLaunchersFor(cfg, u)
		if len(launchers) == 0 {
			continue
		}
		set := true
		if ferr := huh.NewForm(huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Redirect Bun app extraction for user %s?", u.Name)).
				Description(fmt.Sprintf(
					"%s bundles a native library Bun extracts to $TMPDIR each launch and loads;\n"+
						"on world-writable /tmp the trust guard denies that load. This will:\n"+
						"  - create %s (0700, owned by that user)\n"+
						"  - add a shell-function wrapper for %s to the user's startup file\n"+
						"    (.zshrc / .bashrc / fish conf.d) running them with $TMPDIR there\n"+
						"The daemon reserves .bun-* below that dir for the app, so it loads its own\n"+
						"extraction while nothing else can plant or load one. Interactive shells only —\n"+
						"a Bun app launched from a desktop menu still uses /tmp.",
					strings.Join(launchers, ", "), inst.BunTmpDir(u.Home), strings.Join(launchers, ", "))).
				Affirmative("Set it up").
				Negative("Skip").
				Value(&set),
		)).Run(); ferr != nil {
			return nil, ferr
		}
		if !set {
			log.Infof("skipping the Bun $TMPDIR redirect for %s — %s keeps extracting to /tmp and its native library load stays denied",
				u.Name, strings.Join(launchers, ", "))
			continue
		}
		accepted = append(accepted, bunUser{User: u, Launchers: launchers})
	}
	return accepted, nil
}

// bunLaunchersFor returns the launcher command names of every Bun catalog app configured for u.
func bunLaunchersFor(cfg *daemonconfig.Config, u inst.User) []string {
	seen := map[string]bool{}
	var out []string
	for i := range inst.Catalog {
		e := &inst.Catalog[i]
		if !e.IsBun() || !bunEntryConfigured(cfg, e, u) {
			continue
		}
		for _, l := range e.BunLaunchers {
			if !seen[l] {
				seen[l] = true
				out = append(out, l)
			}
		}
	}
	sort.Strings(out)
	return out
}

// bunEntryConfigured reports whether any of cfg's resources is one of e's guarded locations for u.
func bunEntryConfigured(cfg *daemonconfig.Config, e *inst.CandidateDir, u inst.User) bool {
	for i := range cfg.Resources {
		if e.OwnsPath(u.Name, u.Home, cfg.Resources[i].Path) {
			return true
		}
	}
	return false
}

// setupBunTmpdir creates each opted-in user's private Bun tmp dir and installs the launcher
// wrappers. Run before the daemon (re)starts so the reserved dir exists when it builds its
// reservations (a missing dir reserves nothing).
func setupBunTmpdir(users []bunUser) error {
	for _, bu := range users {
		dir, err := inst.EnsureBunTmpDir(bu.User)
		if err != nil {
			return fmt.Errorf("creating the Bun tmp dir for %s: %w", bu.User.Name, err)
		}
		changed, err := inst.EnsureBunLauncherEnv(bu.User, bu.Launchers)
		if err != nil {
			return fmt.Errorf("wrapping Bun launchers for %s: %w", bu.User.Name, err)
		}
		log.Infof("Bun $TMPDIR redirect for %s: %s; wrapped %s in %v",
			bu.User.Name, dir, strings.Join(bu.Launchers, ", "), changed)
	}
	return nil
}
