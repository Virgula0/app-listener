package logging

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	stdlog "log"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/constants"
)

func TestVerboseToLogrus(t *testing.T) {
	cases := []struct {
		in   constants.VerboseLevel
		want log.Level
	}{
		{constants.NeededInfo, log.WarnLevel},
		{constants.InternalWarningsLevelTwo, log.InfoLevel},
		{constants.PrintAdditionalInfoLevelOne, log.DebugLevel},
		{constants.PrintAdditionalInfoLevelTwo, log.TraceLevel},
	}
	for _, tc := range cases {
		got, err := VerboseToLogrus(tc.in)
		if err != nil {
			t.Fatalf("VerboseToLogrus(%d): %v", int(tc.in), err)
		}
		if got != tc.want {
			t.Errorf("VerboseToLogrus(%d) = %v, want %v", int(tc.in), got, tc.want)
		}
	}
	for _, invalid := range []int{-1, 4, 100} {
		if _, err := VerboseToLogrus(constants.VerboseLevel(invalid)); err == nil {
			t.Errorf("VerboseToLogrus(%d) accepted an out-of-range level", invalid)
		}
	}
}

// TestDefaultVerboseMirrorsCurrentDisplay pins the contract that a dump
// taken without --verbose records exactly what is displayed today (the
// default console level is info).
func TestDefaultVerboseMirrorsCurrentDisplay(t *testing.T) {
	got, err := VerboseToLogrus(DefaultVerbose)
	if err != nil {
		t.Fatal(err)
	}
	if got != log.InfoLevel {
		t.Fatalf("default verbose maps to %v, want info", got)
	}
}

// TestDumpHookFiltersByVerbosity: with the "current display" threshold,
// info and more severe lines land in the dump while debug noise stays out.
func TestDumpHookFiltersByVerbosity(t *testing.T) {
	var buf bytes.Buffer
	hook, err := NewDumpHook(&buf, DefaultVerbose)
	if err != nil {
		t.Fatal(err)
	}

	logger := log.New()
	logger.SetOutput(&bytes.Buffer{}) // keep test output clean
	logger.AddHook(hook)
	logger.SetLevel(log.TraceLevel) // hook must do its own filtering

	logger.Debugln("debug-line")
	logger.Infoln("info-line")
	logger.Warnln("warn-line")

	out := buf.String()
	if strings.Contains(out, "debug-line") {
		t.Error("debug entry leaked past the verbosity filter:\n" + out)
	}
	for _, want := range []string{"info-line", "warn-line"} {
		if !strings.Contains(out, want) {
			t.Errorf("dump is missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "\x1b[") {
		t.Errorf("ANSI escape leaked into the dump:\n%q", out)
	}
}

// TestDumpHookFormat checks the plain-text line shape and field ordering.
func TestDumpHookFormat(t *testing.T) {
	var buf bytes.Buffer
	hook, err := NewDumpHook(&buf, constants.PrintAdditionalInfoLevelTwo)
	if err != nil {
		t.Fatal(err)
	}

	logger := log.New()
	logger.SetOutput(&bytes.Buffer{})
	logger.AddHook(hook)
	logger.WithField("path", "/tmp/x").Info("hello world")

	line := buf.String()
	if !strings.HasSuffix(line, "hello world path=/tmp/x\n") {
		t.Fatalf("unexpected line format: %q", line)
	}
	if !strings.Contains(line, "[INFO ]") || !strings.Contains(line, " [") {
		t.Fatalf("level not padded/bracketed: %q", line)
	}
	if len(line) < 19 || line[4] != '-' || line[7] != '-' || line[10] != ' ' {
		t.Fatalf("timestamp shape unexpected: %q", line[:20])
	}
}

// TestDumpHookFileIsSynced exercises the *os.File branch of Fire.
func TestDumpHookFileIsSynced(t *testing.T) {
	path := filepath.Join(t.TempDir(), "dump.log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	hook, err := NewDumpHook(f, constants.NeededInfo)
	if err != nil {
		t.Fatal(err)
	}
	logger := log.New()
	logger.SetOutput(&bytes.Buffer{})
	logger.AddHook(hook)

	logger.Info("dropped by threshold")
	logger.Error("kept-error")

	data, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatal(readErr)
	}
	out := string(data)
	if strings.Contains(out, "dropped") || !strings.Contains(out, "kept-error") {
		t.Fatalf("file filtering wrong: %q", out)
	}
}

// TestBridgeStdLog verifies that standard-library logging (fscrypt's
// channel) flows into logrus as info lines and that restore puts back the
// previous writer. Production installs the dump hook on the standard
// logrus logger, so the bridge must feed exactly that one.
func TestBridgeStdLog(t *testing.T) {
	var buf bytes.Buffer
	hook, err := NewDumpHook(&buf, DefaultVerbose)
	if err != nil {
		t.Fatal(err)
	}

	prevOut := log.StandardLogger().Out
	prevLevel := log.StandardLogger().GetLevel()
	t.Cleanup(func() {
		log.StandardLogger().SetOutput(prevOut)
		log.StandardLogger().SetLevel(prevLevel)
		log.StandardLogger().ReplaceHooks(log.LevelHooks{})
	})
	log.StandardLogger().SetOutput(&bytes.Buffer{}) // keep test output clean
	log.StandardLogger().AddHook(hook)

	restore := BridgeStdLog()
	stdlog.Printf("stdlib-bridge-line %d", 42)
	restore()

	if !strings.Contains(buf.String(), "stdlib-bridge-line 42") {
		t.Fatalf("stdlib line missing from logrus stream: %q", buf.String())
	}
	if stdlog.Writer() != os.Stderr {
		t.Fatalf("restore did not put back stderr, got %T", stdlog.Writer())
	}
}
