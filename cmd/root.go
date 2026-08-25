package cmd

import (
	"fmt"
	"os"

	"github.com/Virgula0/app-listener/cmd/functions/daemon"
	"github.com/Virgula0/app-listener/cmd/functions/editprotected"
	"github.com/Virgula0/app-listener/cmd/functions/guard"
	"github.com/Virgula0/app-listener/cmd/functions/install"
	"github.com/Virgula0/app-listener/cmd/functions/monitor"
	"github.com/Virgula0/app-listener/cmd/functions/networkguard"
	"github.com/Virgula0/app-listener/cmd/functions/networkmonitor"
	"github.com/Virgula0/app-listener/cmd/functions/uninstall"
	"github.com/Virgula0/app-listener/cmd/functions/update"
	"github.com/Virgula0/app-listener/internal/constants"
	"github.com/Virgula0/app-listener/internal/logging"
	"github.com/Virgula0/app-listener/internal/wizard"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app-listener",
	Short: "Monitor and guard file system events using eBPF",
	Long: `app-listener is a TUI/GUI application that monitors and guards
file system operations using eBPF.

Use --gui to launch the graphical interface instead of the terminal TUI
(only available for the monitor subcommand).`,
	// Enables the built-in --version flag; the value is injected at build
	// time by the release workflow (see internal/constants.Version).
	Version: constants.Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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
	},
	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		if restoreStdLog != nil {
			restoreStdLog()
			restoreStdLog = nil
		}
		closeDumpFile()
		return nil
	},
}

var logLevel string

// guiFlag is a persistent root flag consumed by the monitor subcommand (the
// only one that currently supports a GUI backend).
var guiFlag bool

// verboseFlag backs the global --verbose flag. Zero (or omitted) preserves
// today's output exactly; values above zero raise verbosity and are only
// valid together with --headless, since they apply to the headless console
// stream — TUI/GUI presentation never changes.
var verboseFlag int

// dumpLogFile backs the global --dump-log flag: when set, every log line is
// mirrored into this file in plain text.
var dumpLogFile string

// dumpFileHandle is the open --dump-log file; closed after the command ran.
var dumpFileHandle *os.File

// restoreStdLog undoes logging.BridgeStdLog after the command ran.
var restoreStdLog func()

// overwritePrompt asks whether an existing dump file may be replaced.
// A package variable so tests can inject an answer without a terminal.
var overwritePrompt = func(path string) (bool, error) {
	return wizard.ConfirmOnce(fmt.Sprintf("%s already exists — overwrite it with the new dump?", path), "Overwrite")
}

func Execute() error {
	return rootCmd.Execute()
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

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Warnf("invalid log level %q, using info", logLevel)
		level = log.InfoLevel
	}
	log.SetLevel(level)
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Set log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().BoolVarP(&guiFlag, "gui", "", false, "Launch GUI instead of TUI (default: false)")
	rootCmd.PersistentFlags().IntVarP(&verboseFlag, "verbose", "v", 0,
		fmt.Sprintf("Verbosity level 0..%d for the headless console stream and --dump-log "+
			"(0 = default display; %d=essential, %d=current display, %d=+internal details, %d=everything). "+
			"Requires --headless", int(constants.PrintAdditionalInfoLevelTwo),
			int(constants.NeededInfo), int(constants.InternalWarningsLevelTwo),
			int(constants.PrintAdditionalInfoLevelOne), int(constants.PrintAdditionalInfoLevelTwo)))
	rootCmd.PersistentFlags().StringVar(&dumpLogFile, "dump-log", "",
		"Mirror all log lines into this file in plain text (asks before overwriting an existing file)")

	rootCmd.AddCommand(monitor.MonitorCmd)
	rootCmd.AddCommand(guard.GuardCmd)
	rootCmd.AddCommand(networkmonitor.NetworkMonitorCmd)
	rootCmd.AddCommand(networkguard.NetworkGuardCmd)
	rootCmd.AddCommand(daemon.DaemonCmd)
	rootCmd.AddCommand(install.InstallCmd)
	rootCmd.AddCommand(uninstall.UninstallCmd)
	rootCmd.AddCommand(update.UpdateCmd)
	rootCmd.AddCommand(editprotected.EditProtectedCmd)
}

// validateVerbose enforces the --verbose contract: values above zero demand
// --headless on the executed subcommand (verbosity applies to the headless
// console stream only), and the value must be a defined VerboseLevel.
func validateVerbose(cmd *cobra.Command) error {
	flag := cmd.Flags().Lookup("verbose")
	if flag == nil || !flag.Changed || verboseFlag == 0 {
		return nil // omitted or explicit no-op: nothing to enforce
	}
	maxLevel := int(constants.PrintAdditionalInfoLevelTwo)
	if verboseFlag < 0 || verboseFlag > maxLevel {
		return fmt.Errorf("--verbose %d is out of range: valid levels are 0..%d", verboseFlag, maxLevel)
	}
	if !headlessRequested(cmd) {
		return fmt.Errorf("--verbose requires --headless: verbosity applies to the headless console stream " +
			"(supported by daemon, guard, monitor, network-monitor, network-guard)")
	}
	return nil
}

// headlessRequested reports whether the executed command carries a headless
// flag and runs with it enabled.
func headlessRequested(cmd *cobra.Command) bool {
	flag := cmd.Flags().Lookup("headless")
	return flag != nil && flag.Value.String() == "true"
}

// applyVerbosity raises the console log level from --verbose. Omitted or
// zero keeps today's display untouched; an explicit --log-level always wins
// for the console stream.
func applyVerbosity(cmd *cobra.Command) {
	flag := cmd.Flags().Lookup("verbose")
	if flag == nil || !flag.Changed || verboseFlag == 0 {
		return
	}
	if logFlag := cmd.Flags().Lookup("log-level"); logFlag != nil && logFlag.Changed {
		log.Debugf("--log-level overrides the console side of --verbose")
		return
	}
	level, err := logging.VerboseToLogrus(constants.VerboseLevel(verboseFlag))
	if err != nil {
		log.Warnf("ignoring --verbose: %v", err) // unreachable: validateVerbose range-checked
		return
	}
	log.SetLevel(level)
}

// effectiveVerbose returns the verbosity governing the dump file: an
// explicit positive --verbose when given, otherwise the level that mirrors
// what is displayed today (--verbose 0 is a full no-op by contract).
func effectiveVerbose(cmd *cobra.Command) constants.VerboseLevel {
	flag := cmd.Flags().Lookup("verbose")
	if flag != nil && flag.Changed && verboseFlag > 0 {
		return constants.VerboseLevel(verboseFlag)
	}
	return logging.DefaultVerbose
}

// setupDumpLog opens the --dump-log file and installs the mirroring hook.
// An existing file is never destroyed silently: interactively the user is
// asked, under --headless (systemd, no terminal) the run refuses so the
// operator can clear the file deliberately.
func setupDumpLog(cmd *cobra.Command) error {
	if dumpLogFile == "" {
		return nil
	}

	if _, err := os.Stat(dumpLogFile); err == nil {
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

	f, err := os.OpenFile(dumpLogFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
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

// closeDumpFile flushes state after the command completed. The hook itself
// syncs on every write, so abrupt exits never lose emitted lines anyway.
// The dump hook is the only hook this CLI installs, so replacing the whole
// hook set is a precise removal.
func closeDumpFile() {
	if dumpFileHandle == nil {
		return
	}
	_ = dumpFileHandle.Sync()
	_ = dumpFileHandle.Close()
	dumpFileHandle = nil
	log.StandardLogger().ReplaceHooks(log.LevelHooks{})
}
