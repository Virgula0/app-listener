// The `app-listener edit-protected` command tree. The bare command opens a protected directory in
// the embedded two-pane editor, in one of two modes chosen automatically:
//   - offline (daemon stopped, or no password configured): fatally refuses while the daemon runs,
//     re-scans the catalog like the uninstaller, unlocks ONE vault with the master key, edits,
//     re-locks; a vault never stays open after return.
//   - live (password set at install and daemon running): authenticates over the daemon's control
//     socket, which briefly widens that one resource's guard for the write, then narrows it. The
//     daemon keeps running; the vault lifecycle is untouched.
//
// Subcommands: `forward` (temporary whitelist/blacklist rows for other binaries), `config`
// (`config record`), `put` and `password`; all but `password` are live only.
package editprotected

import (
	"errors"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	resourceFlag     string
	forwardResources []string
	fromFlag         string
	whitelistFlags   []string
	blacklistFlags   []string
	eventsFlags      []string
	yesFlag          bool
	idleTimeout      time.Duration
	backendFlag      bool
)

// editPasswordEnv carries the password for the non-interactive live modes (`put`, `config
// --from`). It exists for automation; interactive runs prompt instead.
const editPasswordEnv = "APP_LISTENER_EDIT_PASSWORD" //nolint:gosec // env var name, not a credential

var EditProtectedCmd = &cobra.Command{
	Use:   "edit-protected",
	Short: "Edit protected directories and the daemon config, or grant temporary access",
	Long: `Open a protected directory in the embedded two-pane editor.

Live when the daemon runs with an edit-protected password: the daemon widens
only that directory's guard while you edit. Otherwise offline: the daemon must
be stopped, and the vault is unlocked with the master key, edited and locked
again. Both modes warn about edits that would sit outside the daemon's
protection.`,
	Example: `  sudo app-listener edit-protected
  sudo app-listener edit-protected --resource /home/alice/.ssh
  sudo app-listener edit-protected forward -w /usr/bin/python3 -e OPEN,READ
  sudo app-listener edit-protected config
  sudo app-listener edit-protected password`,
	Args: cobra.NoArgs,
	RunE: runEditProtected,
}

var forwardCmd = &cobra.Command{
	Use:   "forward (-w BINARY | -b BINARY)...",
	Short: "Give other binaries temporary access to guarded directories (live)",
	Long: `Hold temporary rules for other binaries on the guarded directories you pick.

-w admits binaries the directory doesn't whitelist; -b denies whitelisted ones.
-e narrows the rule to some event types (default: all). The terminal shows the
granted binaries' accesses; the grant is revoked when you quit (q or Ctrl+C),
the daemon reloads, a granted binary changes on disk, or --idle-timeout passes
with no activity. A -w grant applies to every user of that binary, and an
interpreter (python, node, bash) admits every script it runs.`,
	Example: `  sudo app-listener edit-protected forward -w ~/.local/bin/python3.12 -e OPEN,READ,STAT
  sudo app-listener edit-protected forward -b /usr/bin/git --resource /home/alice/.ssh`,
	Args: cobra.NoArgs,
	RunE: func(*cobra.Command, []string) error {
		if err := validateForwardEvents(); err != nil {
			return err
		}
		return runLive("forward", "", runForward)
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Edit the running daemon's daemon.conf, then reload it (live)",
	Long: `Edit the running daemon's daemon.conf. Saving shows the diff; the daemon then
writes the file and reloads. A configuration that fails to reload is not kept:
the daemon restores the file and you can edit again.

With the daemon stopped, edit /etc/app-listener/daemon.conf directly.`,
	Example: `  sudo app-listener edit-protected config
  sudo APP_LISTENER_EDIT_PASSWORD=... app-listener edit-protected config --from new.conf`,
	Args: cobra.NoArgs,
	RunE: func(*cobra.Command, []string) error {
		return runLive("config", "; with the daemon stopped, edit /etc/app-listener/daemon.conf directly", runEditConfig)
	},
}

var recordCmd = &cobra.Command{
	Use:   "record PATH...",
	Short: "Learn which binaries use PATH, then guard it with the ones you approve (live)",
	Long: `Record every access to the paths, keyed by the executable's inode, until you
stop it. Then review the binaries (and the libraries they load), optionally
encrypt the paths, and the daemon adds one section per path and reloads. A path
daemon.conf already guards keeps its whitelist and only gains new lines.

Everything that touched the paths is listed, malware included: approve only
what you ran on purpose.`,
	Example: `  sudo app-listener edit-protected config record /home/alice/.config/gh`,
	Args:    cobra.MinimumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		return runLive("config record", "", func() error { return runRecord(args) })
	},
}

var putCmd = &cobra.Command{
	Use:   "put FILE --resource DIR",
	Short: "Write one file into a guarded directory without the editor (live)",
	Long: `Write FILE (absolute, or relative to --resource) with the content of --from,
or stdin. The password comes from $` + editPasswordEnv + `.`,
	Example: `  echo 'Host *' | sudo APP_LISTENER_EDIT_PASSWORD=... app-listener edit-protected put authorized_keys \
      --resource /home/alice/.ssh`,
	Args: cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		return runLive("put", "", func() error { return runPut(args[0]) })
	},
}

var passwordCmd = &cobra.Command{
	Use:   "password",
	Short: "Change the edit-protected password",
	Long: `Change the edit-protected password. The current one is asked first (checked by
the running daemon, or against the same attempt budget with it stopped). The
password can't be removed. An install older than the mandatory install-time
password sets its first one here, with the daemon stopped.`,
	Args: cobra.NoArgs,
	RunE: func(*cobra.Command, []string) error {
		if err := requireRoot(); err != nil {
			return err
		}
		return runSetPassword()
	},
}

func init() {
	addTimeoutFlag := func(c *cobra.Command, activity string) {
		c.Flags().DurationVar(&idleTimeout, "idle-timeout", SessionTimeoutDefault,
			"Revoke the access after this long without "+activity+", e.g. 10m or 45s")
	}

	f := EditProtectedCmd.Flags()
	f.StringVarP(&resourceFlag, "resource", "r", "", "Skip the picker and open this guarded directory")
	addTimeoutFlag(EditProtectedCmd, "editor input (live mode)")
	f.BoolVar(&backendFlag, "backend", false, "The desktop window's channel to the installed binary (JSON lines on stdin/stdout)")

	f = forwardCmd.Flags()
	f.StringSliceVarP(&whitelistFlags, "whitelist", "w", nil, "Binary to admit (repeatable)")
	f.StringSliceVarP(&blacklistFlags, "blacklist", "b", nil, "Whitelisted binary to deny (repeatable)")
	f.StringSliceVarP(&eventsFlags, "events", "e", nil,
		"Event types the rule covers: OPEN,READ,WRITE,DELETE,RENAME,SYMLINK,HARDLINK,MKDIR,MMAP,ATTR,STAT,MKNOD (default all)")
	f.StringArrayVarP(&forwardResources, "resource", "r", nil, "Guarded directory to grant, instead of the picker (repeatable)")
	f.BoolVarP(&yesFlag, "yes", "y", false, "Don't ask before admitting a binary that isn't root-placed (needed without a terminal)")
	addTimeoutFlag(forwardCmd, "a granted binary touching the directories or a keypress")
	forwardCmd.MarkFlagsMutuallyExclusive("whitelist", "blacklist")
	forwardCmd.MarkFlagsOneRequired("whitelist", "blacklist")

	configCmd.Flags().StringVar(&fromFlag, "from", "", "Replace daemon.conf with this file, without the editor")
	configCmd.AddCommand(recordCmd)

	f = putCmd.Flags()
	f.StringVarP(&resourceFlag, "resource", "r", "", "Guarded directory FILE belongs to (required)")
	f.StringVar(&fromFlag, "from", "", "File holding the new content (default: stdin)")
	addTimeoutFlag(putCmd, "activity")
	_ = putCmd.MarkFlagRequired("resource")

	EditProtectedCmd.AddCommand(forwardCmd, configCmd, putCmd, passwordCmd)
}

func requireRoot() error {
	if os.Geteuid() != 0 {
		return errors.New("edit-protected must be run as root: sudo app-listener edit-protected")
	}
	if idleTimeout < time.Second || idleTimeout > SessionTimeoutMax {
		return fmt.Errorf("--idle-timeout must be between 1s and %s", SessionTimeoutMax)
	}
	return nil
}

// runLive runs a live-only subcommand. The control socket exists only when a running daemon has a
// password configured: the exact precondition for live mode, independent of whether systemd
// supervises the daemon.
func runLive(name, alt string, run func() error) error {
	if err := requireRoot(); err != nil {
		return err
	}
	if !LiveModeAvailable() {
		return errors.New("edit-protected " + name + " needs the daemon running with an edit-protected password " +
			"configured (the control socket " + ControlSocket + " is not present)" + alt)
	}
	return run()
}

func runEditProtected(cmd *cobra.Command, _ []string) error {
	if err := requireRoot(); err != nil {
		return err
	}
	if backendFlag {
		if cmd.Flags().NFlag() != 1 {
			return errors.New("--backend takes no other flag")
		}
		return runBackend()
	}

	hasPassword, err := HashConfigured()
	if err != nil {
		return err
	}
	if LiveModeAvailable() {
		return runLiveEdit()
	}
	if hasPassword {
		log.Info("an edit-protected password is configured, but the daemon control socket is not present — using the offline flow")
	}
	return runOfflineEdit()
}
