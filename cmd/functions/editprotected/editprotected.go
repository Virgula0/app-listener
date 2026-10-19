// The `app-listener edit-protected` command opens fscrypt-encrypted (or plain) protected
// directories in the embedded two-pane editor. Two modes, chosen automatically:
//   - offline (daemon stopped, or no password configured): fatally refuses while the daemon runs,
//     re-scans the catalog like the uninstaller, unlocks ONE vault with the master key, edits,
//     re-locks; a vault never stays open after return.
//   - live (password set at install and daemon running): authenticates over the daemon's control
//     socket, which briefly widens that one resource's guard for the write, then narrows it. The
//     daemon keeps running; the vault lifecycle is untouched.
//
// `--set-password` changes the password (it can't be removed). `--forward -w/-b [-e]` runs no editor:
// it holds temporary whitelist/blacklist rows for other binaries on the picked resources (live only).
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
	setPasswordFlag bool
	resourceFlags   []string
	putFlag         string
	contentFileFlag string
	forwardFlag     bool
	whitelistFlags  []string
	blacklistFlags  []string
	eventsFlags     []string
	yesFlag         bool
	timeoutSession  time.Duration
	editConfigFlag  bool
	recordFlags     []string
	backendFlag     bool
)

// editPasswordEnv carries the password for the non-interactive live apply
// mode (--put). It exists for automation; interactive runs prompt instead.
const editPasswordEnv = "APP_LISTENER_EDIT_PASSWORD" //nolint:gosec // env var name, not a credential

var EditProtectedCmd = &cobra.Command{
	Use:   "edit-protected",
	Short: "Edit a protected directory in the embedded TUI editor (live when an edit-protected password is set)",
	Long: `Open one of the protected directories in the embedded two-pane editor.

When the daemon is running with the edit-protected password chosen during
` + "`app-listener install`" + `, the edit happens LIVE: the command authenticates to
the daemon over its local control socket, the daemon briefly grants write
access to that one resource, you edit, and the grant is dropped again. The
daemon keeps running and the fscrypt vaults are never touched. A live session
is revoked after --timeout-session (default 30m) without activity: every
editor keypress resets it.

Otherwise it runs OFFLINE: it fatally refuses while the daemon is running,
re-scans the catalog (NOT the installed daemon.conf) like the uninstaller,
unlocks ONE directory with the master key, runs the editor, and re-locks it
no matter how the editor exits.

Before exiting, both modes audit the edited tree against daemon.conf and
warn about anything that would sit outside the daemon's protection (a file
created outside every guarded watch path, a new symlink, a world-readable
new secret, a freshly dropped executable).

--edit-config edits the running daemon's daemon.conf (live mode only, same
password): saving shows the diff, then the daemon writes it and reloads. If the
reload fails, the daemon keeps the current configuration, restores the file,
and you can edit again. --content-file <file> replaces it without the editor.

--edit-config --record <path> (repeatable) adds new guarded paths by learning who
uses them: it records every access to the paths, keyed by the executable's
inode, until you stop it, lists the binaries (and the libraries they load) for
you to approve, can encrypt the paths first, then adds one section per path and
reloads the daemon. A path daemon.conf already guards is extended instead: its
current whitelist is shown and kept, and only new lines are added (the guard
keeps denying the rest while recording). Anything that touched the paths while
recording is listed, malware included: approve only what you ran on purpose.

Use --set-password to change the edit-protected password: it asks the
current one first (checked by the running daemon, or with the daemon stopped
against the same attempt budget). The password can't be removed. An install
older than the mandatory install-time password sets its first one here with
the daemon stopped.

For automation, --resource <path> --put <file-in-that-resource> performs one
live write non-interactively: the new content is read from --content-file
(or stdin) and the password from the ` + editPasswordEnv + ` environment
variable. It requires the daemon to be running with a configured password.

--forward opens no editor: it gives other binaries temporary access to the
guarded directories you pick (live mode only), with the guard command's
syntax. -w admits the listed binaries (every event, or only those given with
-e); -b denies already-whitelisted ones (every event, or only the -e ones).
The terminal shows the granted binaries' accesses in the guard view; the
grant lasts until you quit it (q or Ctrl+C), the daemon reloads, a granted
binary changes on disk, or --timeout-session passes without a granted binary
touching the directories or a keypress; then it is revoked.
A -w grant applies to every user running that binary, and admitting an
interpreter (python, node, bash) admits every script it runs.

  sudo app-listener edit-protected --forward \\
      -w ~/.local/share/uv/python/cpython-3.12/bin/python3.12 -e OPEN,READ,STAT,MKDIR,WRITE`,
	Args: cobra.NoArgs,
	RunE: runEditProtected,
}

func init() {
	EditProtectedCmd.Flags().BoolVar(&setPasswordFlag, "set-password", false,
		"Change the edit-protected authentication password (asks the current one first)")
	EditProtectedCmd.Flags().StringArrayVar(&resourceFlags, "resource", nil,
		"Skip the picker and act on this configured watch path (repeatable with --forward)")
	EditProtectedCmd.Flags().StringVar(&putFlag, "put", "",
		"Non-interactive live mode: write this file (must be inside --resource) from --content-file/stdin, authenticating with $"+editPasswordEnv)
	EditProtectedCmd.Flags().StringVar(&contentFileFlag, "content-file", "",
		"Source of the --put content (default: stdin), or the new daemon.conf for --edit-config")
	EditProtectedCmd.Flags().BoolVar(&forwardFlag, "forward", false,
		"Run no editor: grant the -w/-b binaries temporary access to the picked guarded directories (live mode)")
	EditProtectedCmd.Flags().StringSliceVarP(&whitelistFlags, "whitelist", "w", nil,
		"With --forward: binaries to admit temporarily (repeatable, mutually exclusive with -b)")
	EditProtectedCmd.Flags().StringSliceVarP(&blacklistFlags, "blacklist", "b", nil,
		"With --forward: whitelisted binaries to deny temporarily (repeatable, mutually exclusive with -w)")
	EditProtectedCmd.Flags().StringSliceVarP(&eventsFlags, "events", "e", nil,
		"With --forward: event types the rule covers (comma-separated: OPEN,READ,WRITE,DELETE,RENAME,SYMLINK,HARDLINK,MKDIR,MMAP,ATTR,STAT,MKNOD; default: all)")
	EditProtectedCmd.Flags().DurationVar(&timeoutSession, "timeout-session", SessionTimeoutDefault,
		"Live sessions: revoke the access after this long without activity (editor input, or with "+
			"--forward a granted binary touching the directory), e.g. 10m or 45s")
	EditProtectedCmd.Flags().BoolVar(&editConfigFlag, "edit-config", false,
		"Edit the running daemon's daemon.conf (live mode): saving reloads it; a configuration that fails to "+
			"reload is not kept (with --content-file: replace it non-interactively)")
	EditProtectedCmd.Flags().StringArrayVar(&recordFlags, "record", nil,
		"With --edit-config: record which binaries use this new path (repeatable), then add it to daemon.conf "+
			"with the binaries you approve")
	EditProtectedCmd.Flags().BoolVar(&yesFlag, "yes", false,
		"With --forward: don't ask before admitting a binary that isn't root-placed (required without a terminal)")
	EditProtectedCmd.Flags().BoolVar(&backendFlag, "backend", false,
		"The desktop window's channel to the installed binary (JSON lines on stdin/stdout)")
	_ = EditProtectedCmd.Flags().MarkHidden("backend")
}

func validateFlags() error {
	if err := validateForwardFlags(); err != nil {
		return err
	}
	if !forwardFlag && len(resourceFlags) > 1 {
		return errors.New("--resource is repeatable only with --forward")
	}
	if timeoutSession < time.Second || timeoutSession > SessionTimeoutMax {
		return fmt.Errorf("--timeout-session must be between 1s and %s", SessionTimeoutMax)
	}
	if editConfigFlag && (forwardFlag || putFlag != "" || setPasswordFlag || len(resourceFlags) > 0) {
		return errors.New("--edit-config can't be combined with --forward, --put, --resource or --set-password")
	}
	if len(recordFlags) > 0 && (!editConfigFlag || contentFileFlag != "") {
		return errors.New("--record needs --edit-config, without --content-file")
	}
	return nil
}

func runEditProtected(cmd *cobra.Command, args []string) error {
	if os.Geteuid() != 0 {
		return errors.New("edit-protected must be run as root: sudo app-listener edit-protected")
	}
	if backendFlag {
		if cmd.Flags().NFlag() != 1 {
			return errors.New("--backend takes no other flag")
		}
		return runBackend()
	}

	if err := validateFlags(); err != nil {
		return err
	}
	if setPasswordFlag {
		return runSetPassword()
	}

	hasPassword, err := HashConfigured()
	if err != nil {
		return err
	}
	// The control socket exists only when a running daemon has a password configured: the exact
	// precondition for live mode, independent of whether systemd supervises the daemon.
	liveReady := LiveModeAvailable()

	for _, m := range []struct {
		on        bool
		flag, alt string
		run       func() error
	}{
		{editConfigFlag, "--edit-config", "; with the daemon stopped, edit /etc/app-listener/daemon.conf directly", runEditConfig},
		{forwardFlag, "--forward", "", runForward},
		{putFlag != "", "--put", "", runNonInteractiveLivePut},
	} {
		if !m.on {
			continue
		}
		if !liveReady {
			return errors.New(m.flag + " needs the daemon running with an edit-protected password configured " +
				"(the control socket " + ControlSocket + " is not present)" + m.alt)
		}
		return m.run()
	}

	if liveReady {
		return runLiveEdit()
	}
	if hasPassword {
		log.Info("an edit-protected password is configured, but the daemon control socket is not present — using the offline flow")
	}
	return runOfflineEdit()
}
