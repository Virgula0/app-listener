package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Virgula0/app-listener/cmd/functions/daemon"
	"github.com/Virgula0/app-listener/cmd/functions/desktop"
	"github.com/Virgula0/app-listener/cmd/functions/editprotected"
	"github.com/Virgula0/app-listener/cmd/functions/guard"
	"github.com/Virgula0/app-listener/cmd/functions/install"
	"github.com/Virgula0/app-listener/cmd/functions/monitor"
	"github.com/Virgula0/app-listener/cmd/functions/network"
	"github.com/Virgula0/app-listener/cmd/functions/trust"
	"github.com/Virgula0/app-listener/cmd/functions/uninstall"
	"github.com/Virgula0/app-listener/cmd/functions/update"
	"github.com/Virgula0/app-listener/internal/constants"
	"github.com/Virgula0/app-listener/internal/gui"
	"github.com/Virgula0/app-listener/internal/logging"
	"github.com/Virgula0/app-listener/internal/safeio"
	"github.com/Virgula0/app-listener/internal/wizard"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app-listener",
	Short: "Monitor and guard file system and network access with eBPF",
	Long: `app-listener monitors or denies file system and network operations with eBPF,
to keep credential directories (SSH, GPG, browsers, AI agents) away from
info-stealers. Every mode runs as root.

The desktop window is the separate app-listener-gui build.`,
	// Enables the built-in --version flag; value injected at build time
	// (see internal/constants.Version).
	Version:            constants.Version,
	SilenceErrors:      true, // main logs it
	PersistentPreRunE:  preRun,
	PersistentPostRunE: postRun,
}

// guiCmd is the root of the gui build: the window configures every operation itself, so it has no
// subcommands, and the daemon can never run the fyne build.
var guiCmd = &cobra.Command{
	Use:   "app-listener-gui",
	Short: "The app-listener desktop window",
	Long: `Open the app-listener desktop window: the file and network workspaces, the
protected directories, the installer and the updater. Run it as root from a
graphical session. The command-line modes live in the app-listener binary.`,
	Version:            constants.Version,
	SilenceErrors:      true,
	Args:               cobra.NoArgs,
	PersistentPreRunE:  preRun,
	PersistentPostRunE: postRun,
	RunE:               func(*cobra.Command, []string) error { return desktop.Run() },
}

func preRun(cmd *cobra.Command, _ []string) error {
	// Arguments and flags are valid by now: a later error is not a usage mistake.
	cmd.SilenceUsage = true
	initLogger()
	if err := validateVerbose(cmd); err != nil {
		return err
	}
	applyVerbosity(cmd)
	if err := setupDumpLog(cmd); err != nil {
		return err
	}
	restoreStdLog = logging.BridgeStdLog()
	return nil
}

func postRun(*cobra.Command, []string) error {
	if restoreStdLog != nil {
		restoreStdLog()
		restoreStdLog = nil
	}
	closeDumpFile()
	return nil
}

// verboseFlag backs the global --verbose (constants.VerboseLevel ladder) for the headless console
// and --dump-log threshold. Omitted keeps the current display; any explicit value, even 0, requires
// --headless. TUI/GUI is unaffected.
var verboseFlag int

// dumpLogFile backs the global --dump-log flag: every log line is mirrored
// into this file in plain text.
var dumpLogFile string

// dumpFileHandle is the open --dump-log file; closed after the command ran.
var dumpFileHandle *os.File

// restoreStdLog undoes logging.BridgeStdLog after the command ran.
var restoreStdLog func()

// overwritePrompt asks whether an existing dump file may be replaced; a package
// variable so tests can inject an answer without a terminal.
var overwritePrompt = func(path string) (bool, error) {
	return wizard.ConfirmOnce(fmt.Sprintf("%s already exists — overwrite it with the new dump?", path), "Overwrite")
}

func Execute() error {
	initLogger()
	if gui.Available {
		return guiCmd.Execute()
	}
	return rootCmd.Execute()
}

// movedFlags names the subcommand that replaced a flag, keyed by command path, for the unknown-flag
// error.
var movedFlags = map[string]map[string]string{
	"app-listener edit-protected": {
		"forward": "edit-protected forward", "edit-config": "edit-protected config",
		"record": "edit-protected config record PATH...", "put": "edit-protected put FILE --resource DIR",
		"set-password": "edit-protected password", "whitelist": "edit-protected forward",
		"blacklist": "edit-protected forward", "events": "edit-protected forward",
		"timeout-session": "--idle-timeout",
	},
	"app-listener install": {
		"binary-only": "install binary", "refresh-catalog-only": "install catalog refresh",
		"live": "install catalog refresh --live", "diff-catalog": "install catalog diff",
		"restore-backups": "install backups restore", "delete-post-backups": "install backups delete",
	},
	"app-listener monitor":                {"watch": "monitor PATH... (paths are arguments)"},
	"app-listener guard":                  {"events": "no -e: the guard covers every operation"},
	"app-listener trust":                  {"list": "trust list"},
	"app-listener edit-protected config":  {"content-file": "--from", "record": "edit-protected config record PATH..."},
	"app-listener edit-protected put":     {"content-file": "--from", "timeout-session": "--idle-timeout"},
	"app-listener edit-protected forward": {"timeout-session": "--idle-timeout"},
}

// movedFlagError adds the replacement to an unknown-flag error for a flag that became a
// subcommand.
func movedFlagError(cmd *cobra.Command, err error) error {
	msg := err.Error()
	var name string
	switch {
	case strings.HasPrefix(msg, "unknown flag: --"):
		name = strings.TrimPrefix(msg, "unknown flag: --")
	case strings.HasPrefix(msg, "unknown shorthand flag: 'w'") && cmd.Name() == "monitor":
		name = "watch"
	case strings.HasPrefix(msg, "unknown shorthand flag: "):
		if short, _, ok := strings.Cut(strings.TrimPrefix(msg, "unknown shorthand flag: '"), "'"); ok {
			for long, abbr := range map[string]string{"whitelist": "w", "blacklist": "b", "events": "e"} {
				if short == abbr {
					name = long
				}
			}
		}
	}
	if repl, ok := movedFlags[cmd.CommandPath()][name]; ok {
		return fmt.Errorf("%w\n--%s moved: use `app-listener %s`", err, name, repl)
	}
	return err
}

func initLogger() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:            true,
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
	log.SetOutput(os.Stderr)

	// Baseline is the historical default display; --verbose adjusts it in applyVerbosity.
	log.SetLevel(log.InfoLevel)
}

func init() {
	rootCmd.PersistentFlags().IntVarP(&verboseFlag, "verbose", "v", 0,
		fmt.Sprintf("With --headless: verbosity 0..%d (%d errors, %d essential, %d default, %d internal details, %d everything)",
			int(constants.PrintAdditionalInfoLevelTwo), int(constants.ErrorsOnly), int(constants.NeededInfo),
			int(constants.InternalWarningsLevelTwo), int(constants.PrintAdditionalInfoLevelOne),
			int(constants.PrintAdditionalInfoLevelTwo)))
	for _, c := range []*cobra.Command{rootCmd, guiCmd} {
		c.PersistentFlags().StringVar(&dumpLogFile, "dump-log", "",
			"Also write every log line to this file, in plain text (asks before overwriting)")
	}
	rootCmd.SetFlagErrorFunc(movedFlagError)
	// Accepted and ignored: launchers and docs of earlier releases start the window with --gui.
	guiCmd.Flags().Bool("gui", false, "Ignored: the window opens without it")

	rootCmd.AddGroup(
		&cobra.Group{ID: "modes", Title: "Monitor and enforce:"},
		&cobra.Group{ID: "daemon", Title: "Daemon:"},
		&cobra.Group{ID: "setup", Title: "Setup:"},
	)
	for group, cmds := range map[string][]*cobra.Command{
		"modes":  {monitor.MonitorCmd, guard.GuardCmd, network.NetworkCmd},
		"daemon": {daemon.DaemonCmd, editprotected.EditProtectedCmd, trust.TrustCmd},
		"setup":  {install.InstallCmd, uninstall.UninstallCmd, update.UpdateCmd},
	} {
		for _, c := range cmds {
			c.GroupID = group
			rootCmd.AddCommand(c)
		}
	}
	// An `update` of another version runs it on the new binary: never rename it.
	rootCmd.AddCommand(update.DaemonUnitCmd)
	network.NetworkCmd.SuggestFor = []string{"network-monitor", "network-guard"}
	trust.TrustCmd.SuggestFor = []string{"trust-binaries"}
}

// validateVerbose enforces the --verbose contract: any explicit value (even 0, errors-only) demands
// --headless, and the value must be a defined VerboseLevel.
func validateVerbose(cmd *cobra.Command) error {
	flag := cmd.Flags().Lookup("verbose")
	if flag == nil || !flag.Changed {
		return nil // omitted: default display everywhere
	}
	maxLevel := int(constants.PrintAdditionalInfoLevelTwo)
	if verboseFlag < 0 || verboseFlag > maxLevel {
		return fmt.Errorf("--verbose %d is out of range: valid levels are 0..%d", verboseFlag, maxLevel)
	}
	if !headlessRequested(cmd) {
		return fmt.Errorf("--verbose requires --headless: verbosity applies to the headless console stream " +
			"(supported by daemon, monitor, guard, network monitor, network guard)")
	}
	return nil
}

// headlessRequested reports whether the executed command runs with --headless enabled.
func headlessRequested(cmd *cobra.Command) bool {
	flag := cmd.Flags().Lookup("headless")
	return flag != nil && flag.Value.String() == "true"
}

// applyVerbosity sets the console log level from an explicit --verbose; omitted keeps
// the initLogger info baseline (today's display).
func applyVerbosity(cmd *cobra.Command) {
	flag := cmd.Flags().Lookup("verbose")
	if flag == nil || !flag.Changed {
		return
	}
	level, err := logging.VerboseToLogrus(constants.VerboseLevel(verboseFlag))
	if err != nil {
		log.Warnf("ignoring --verbose: %v", err) // unreachable: validateVerbose range-checked
		return
	}
	log.SetLevel(level)
}

// effectiveVerbose returns the dump-file verbosity: an explicit --verbose when given,
// otherwise the DefaultVerbose mirror of today's display.
func effectiveVerbose(cmd *cobra.Command) constants.VerboseLevel {
	flag := cmd.Flags().Lookup("verbose")
	if flag != nil && flag.Changed {
		return constants.VerboseLevel(verboseFlag)
	}
	return logging.DefaultVerbose
}

// setupDumpLog opens the --dump-log file and installs the mirroring hook. An existing file is never
// silently destroyed: interactive runs ask; --headless refuses so the operator clears it
// deliberately.
func setupDumpLog(cmd *cobra.Command) error {
	if dumpLogFile == "" {
		return nil
	}

	// Lstat, not Stat: setupDumpLog runs as root in PersistentPreRunE (before the daemon's
	// self-guards attach), so a symlink planted at an operator-chosen path in a user-writable dir
	// must never be followed — a symlink whose target does not exist yet would slip past a Stat-based
	// existence check and have root create the target, and one pointing at an existing file would
	// truncate it. Refuse any symlink outright; the open below adds O_NOFOLLOW for the race.
	if lst, err := os.Lstat(dumpLogFile); err == nil {
		if lst.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("refusing to use a symlink as the --dump-log target: %s", dumpLogFile)
		}
		if headlessRequested(cmd) {
			return fmt.Errorf("refusing to overwrite existing dump file %s under --headless: remove it first", dumpLogFile)
		}
		overwrite, promptErr := overwritePrompt(dumpLogFile)
		if promptErr != nil {
			return fmt.Errorf("dump-log overwrite prompt: %w", promptErr)
		}
		if !overwrite {
			return fmt.Errorf("dump-log aborted: %s already exists", dumpLogFile)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("checking dump file %s: %w", dumpLogFile, err)
	}

	f, err := safeio.OpenRegularNoFollow(dumpLogFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("opening dump file %s: %w", dumpLogFile, err)
	}
	hook, err := logging.NewDumpHook(f, effectiveVerbose(cmd))
	if err != nil {
		_ = f.Close()
		return err
	}
	log.AddHook(hook)
	dumpFileHandle = f
	return nil
}

// closeDumpFile flushes after the command completes (the hook syncs every write, so abrupt exits
// lose nothing). Replacing the whole hook set is a precise removal: it's the only hook this CLI
// installs.
func closeDumpFile() {
	if dumpFileHandle == nil {
		return
	}
	_ = dumpFileHandle.Sync()
	_ = dumpFileHandle.Close()
	dumpFileHandle = nil
	log.StandardLogger().ReplaceHooks(log.LevelHooks{})
}
