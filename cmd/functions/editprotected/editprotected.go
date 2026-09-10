// The `app-listener edit-protected` command opens fscrypt-encrypted
// (or plain) protected directories in the embedded two-pane editor.
//
// Two modes, chosen automatically:
//
//   - offline (no edit-protected password configured, or the daemon is
//     stopped): fatally refuses while the daemon runs, re-scans the catalog
//     like the uninstaller, unlocks ONE vault with the master key, edits, and
//     re-locks it — a vault never stays open after the command returns.
//
//   - live (an edit-protected password was set at install time and the
//     daemon is running): authenticates to the daemon over its local control
//     socket, which briefly widens that one resource's guard so the edit can
//     be written, then narrows it again. The daemon keeps running and the
//     vault lifecycle is never touched.
//
// `--set-password` / `--clear-password` manage the password itself.
package editprotected

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	setPasswordFlag   bool
	clearPasswordFlag bool
	resourceFlag      string
	putFlag           string
	contentFileFlag   string
)

// editPasswordEnv carries the password for the non-interactive live apply
// mode (--put). It exists for automation; interactive runs prompt instead.
const editPasswordEnv = "APP_LISTENER_EDIT_PASSWORD" //nolint:gosec // env var name, not a credential

var EditProtectedCmd = &cobra.Command{
	Use:   "edit-protected",
	Short: "Edit a protected directory in the embedded TUI editor (live when an edit-protected password is set)",
	Long: `Open one of the protected directories in the embedded two-pane editor.

If an edit-protected password was chosen during ` + "`app-listener install`" + ` and
the daemon is running, the edit happens LIVE: the command authenticates to
the daemon over its local control socket, the daemon briefly grants write
access to that one resource, you edit, and the grant is dropped again. The
daemon keeps running and the fscrypt vaults are never touched.

Otherwise it runs OFFLINE: it fatally refuses while the daemon is running,
re-scans the catalog (NOT the installed daemon.conf) like the uninstaller,
unlocks ONE directory with the master key, runs the editor, and re-locks it
no matter how the editor exits.

Before exiting, both modes audit the edited tree against daemon.conf and
warn about anything that would sit outside the daemon's protection (a file
created outside every guarded watch path, a new symlink, a world-readable
new secret, a freshly dropped executable).

Use --set-password to set or rotate the edit-protected password (only when
it was not chosen during installation — rotating that one requires
re-running the installer), and --clear-password to remove it.

For automation, --resource <path> --put <file-in-that-resource> performs one
live write non-interactively: the new content is read from --content-file
(or stdin) and the password from the ` + editPasswordEnv + ` environment
variable. It requires the daemon to be running with a configured password.`,
	Args: cobra.NoArgs,
	RunE: runEditProtected,
}

func init() {
	EditProtectedCmd.Flags().BoolVar(&setPasswordFlag, "set-password", false,
		"Set or rotate the edit-protected authentication password (refused if it was set during installation)")
	EditProtectedCmd.Flags().BoolVar(&clearPasswordFlag, "clear-password", false,
		"Remove the edit-protected authentication password (disables live mode)")
	EditProtectedCmd.Flags().StringVar(&resourceFlag, "resource", "",
		"Skip the picker and act on this configured watch path")
	EditProtectedCmd.Flags().StringVar(&putFlag, "put", "",
		"Non-interactive live mode: write this file (must be inside --resource) from --content-file/stdin, authenticating with $"+editPasswordEnv)
	EditProtectedCmd.Flags().StringVar(&contentFileFlag, "content-file", "",
		"Source of the --put content (default: stdin)")
}

func runEditProtected(cmd *cobra.Command, args []string) error {
	if os.Geteuid() != 0 {
		return errors.New("edit-protected must be run as root: sudo app-listener edit-protected")
	}

	if setPasswordFlag && clearPasswordFlag {
		return errors.New("--set-password and --clear-password are mutually exclusive")
	}
	if setPasswordFlag {
		return runSetPassword()
	}
	if clearPasswordFlag {
		return runClearPassword()
	}

	hasPassword, err := HashFileExists()
	if err != nil {
		return err
	}
	// The control socket exists only when a running daemon has a password
	// configured — the exact precondition for live mode, and independent of
	// whether systemd supervises the daemon.
	liveReady := LiveModeAvailable()

	if putFlag != "" {
		if !liveReady {
			return errors.New("--put needs the daemon running with an edit-protected password configured " +
				"(the control socket " + ControlSocket + " is not present)")
		}
		return runNonInteractiveLivePut()
	}

	if liveReady {
		return runLiveEdit()
	}
	if hasPassword {
		log.Info("an edit-protected password is configured, but the daemon control socket is not present — using the offline flow")
	}
	return runOfflineEdit()
}
