package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// newFlagTestCommand registers a throwaway subcommand carrying a headless
// flag, mirroring the daemon/guard/monitor family, and returns a cleanup.
func newFlagTestCommand(t *testing.T) *cobra.Command {
	t.Helper()
	testCmd := &cobra.Command{
		Use: "flagtest",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("flagtest executed")
			return nil
		},
	}
	testCmd.Flags().BoolP("headless", "", false, "")
	rootCmd.AddCommand(testCmd)
	t.Cleanup(func() {
		rootCmd.RemoveCommand(testCmd)
	})
	return testCmd
}

// resetLoggingState restores global logger state mutated by the root hooks
// and resets the persistent flag set: cobra keeps "Changed" markers and
// parsed values across in-process Execute calls, which would otherwise leak
// between tests.
func resetLoggingState(t *testing.T) {
	t.Helper()
	prevLevel := log.StandardLogger().GetLevel()
	prevOut := log.StandardLogger().Out
	t.Cleanup(func() {
		closeDumpFile()
		if restoreStdLog != nil {
			restoreStdLog()
			restoreStdLog = nil
		}
		log.StandardLogger().SetOutput(prevOut)
		log.SetLevel(prevLevel)
	})

	logLevel, guiFlag, verboseFlag, dumpLogFile = "info", false, 0, ""
	dumpFileHandle = nil
	rootCmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		f.Changed = false
	})
}

// executeRoot runs the root command with the given arguments, silencing
// cobra's usage/error printing.
func executeRoot(t *testing.T, args ...string) error {
	t.Helper()
	newFlagTestCommand(t)
	resetLoggingState(t)

	prevArgs, prevSilenceUsage, prevSilenceErrors := os.Args, rootCmd.SilenceUsage, rootCmd.SilenceErrors
	rootCmd.SetArgs(args)
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	defer func() {
		rootCmd.SetArgs(nil)
		rootCmd.SilenceUsage = prevSilenceUsage
		rootCmd.SilenceErrors = prevSilenceErrors
		_ = prevArgs
	}()
	return rootCmd.Execute()
}

func TestVerboseWithoutHeadlessRejected(t *testing.T) {
	err := executeRoot(t, "flagtest", "--verbose", "2")
	if err == nil || !strings.Contains(err.Error(), "--headless") {
		t.Fatalf("err = %v, want --headless requirement", err)
	}
}

func TestVerboseOutOfRangeRejected(t *testing.T) {
	err := executeRoot(t, "flagtest", "--headless", "--verbose", "9")
	if err == nil || !strings.Contains(err.Error(), "out of range") {
		t.Fatalf("err = %v, want out-of-range rejection", err)
	}
}

func TestVerboseNegativeRejected(t *testing.T) {
	err := executeRoot(t, "flagtest", "--headless", "--verbose", "-1")
	if err == nil || !strings.Contains(err.Error(), "out of range") {
		t.Fatalf("err = %v, want out-of-range rejection", err)
	}
}

func TestVerboseWithHeadlessRaisesLevel(t *testing.T) {
	if err := executeRoot(t, "flagtest", "--headless", "--verbose", "2"); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if got := log.GetLevel(); got != log.DebugLevel {
		t.Fatalf("level after --verbose 2 = %v, want debug", got)
	}
}

func TestVerboseZeroIsNoOpAnywhere(t *testing.T) {
	if err := executeRoot(t, "flagtest", "--verbose", "0"); err != nil {
		t.Fatalf("--verbose 0 must not demand headless: %v", err)
	}
	if got := log.GetLevel(); got != log.InfoLevel {
		t.Fatalf("--verbose 0 changed the level to %v, want default info", got)
	}
}

func TestExplicitLogLevelWinsOverVerbose(t *testing.T) {
	if err := executeRoot(t, "flagtest", "--headless", "--verbose", "2", "--log-level", "error"); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if got := log.GetLevel(); got != log.ErrorLevel {
		t.Fatalf("explicit --log-level ignored: %v, want error", got)
	}
}

func TestDumpLogCreatesFileAndCapturesLines(t *testing.T) {
	dump := filepath.Join(t.TempDir(), "session.log")
	if err := executeRoot(t, "flagtest", "--dump-log", dump); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	data, readErr := os.ReadFile(dump)
	if readErr != nil {
		t.Fatalf("dump file not created: %v", readErr)
	}
	out := string(data)
	if !strings.Contains(out, "[INFO") {
		t.Fatalf("expected at least one info line in the dump, got: %q", out)
	}
	if strings.Contains(out, "\x1b[") {
		t.Fatalf("ANSI colors leaked into the dump: %q", out)
	}
}

func TestDumpLogExistingRefusedUnderHeadless(t *testing.T) {
	dump := filepath.Join(t.TempDir(), "session.log")
	if err := os.WriteFile(dump, []byte("previous session"), 0o600); err != nil {
		t.Fatal(err)
	}
	err := executeRoot(t, "flagtest", "--headless", "--dump-log", dump)
	if err == nil || !strings.Contains(err.Error(), "refusing to overwrite") {
		t.Fatalf("err = %v, want overwrite refusal under headless", err)
	}
	data, readErr := os.ReadFile(dump)
	if readErr != nil || string(data) != "previous session" {
		t.Fatalf("existing dump was modified: %q (%v)", data, readErr)
	}
}

func TestDumpLogExistingInteractivePromptHonored(t *testing.T) {
	dump := filepath.Join(t.TempDir(), "session.log")
	if err := os.WriteFile(dump, []byte("previous session"), 0o600); err != nil {
		t.Fatal(err)
	}

	// Decline: run aborts, file untouched.
	overwritePrompt = func(path string) (bool, error) { return false, nil }
	err := executeRoot(t, "flagtest", "--dump-log", dump)
	if err == nil || !strings.Contains(err.Error(), "aborted") {
		t.Fatalf("err = %v, want aborted-by-user", err)
	}
	data, _ := os.ReadFile(dump)
	if string(data) != "previous session" {
		t.Fatalf("file modified despite declining: %q", data)
	}

	// Accept: file truncated and used for the new dump.
	overwritePrompt = func(path string) (bool, error) { return true, nil }
	if err := executeRoot(t, "flagtest", "--dump-log", dump); err != nil {
		t.Fatalf("Execute with accepted overwrite: %v", err)
	}
	data, _ = os.ReadFile(dump)
	if strings.Contains(string(data), "previous session") {
		t.Fatal("accepted overwrite did not truncate the old content")
	}
	if !strings.Contains(string(data), "[INFO") {
		t.Fatalf("new dump has no content: %q", data)
	}
}

func TestDumpLogWithVerboseUsesVerboseThresholdInFile(t *testing.T) {
	dump := filepath.Join(t.TempDir(), "session.log")
	if err := executeRoot(t, "flagtest", "--headless", "--verbose", "0", "--dump-log", dump); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	// verbose 0 is a no-op: the dump uses DefaultVerbose (info).
	data, _ := os.ReadFile(dump)
	if !strings.Contains(string(data), "[INFO") {
		t.Fatalf("default-threshold dump missing info lines: %q", data)
	}
}
