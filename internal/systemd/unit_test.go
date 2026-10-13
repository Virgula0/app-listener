package systemd

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const testUnit = "[Unit]\nDescription=x\n\n[Service]\nExecStart=/usr/local/sbin/app-listener daemon --headless --no-log-metadata-blocks\n"

// fakeBinary writes a script that answers DaemonUnitCommand by printing body (exit code exit).
func fakeBinary(t *testing.T, body string, exit int) string {
	t.Helper()
	dir := t.TempDir()
	payload := filepath.Join(dir, "unit.out")
	if err := os.WriteFile(payload, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	bin := filepath.Join(dir, "app-listener")
	script := "#!/bin/sh\n[ \"$1\" = " + DaemonUnitCommand + " ] || exit 2\ncat '" + payload + "'\nexit " +
		strconv.Itoa(exit) + "\n"
	if err := os.WriteFile(bin, []byte(script), 0o700); err != nil {
		t.Fatal(err)
	}
	return bin
}

func withUnitFile(t *testing.T, content string) {
	t.Helper()
	orig := daemonUnitFile
	t.Cleanup(func() { daemonUnitFile = orig })
	daemonUnitFile = filepath.Join(t.TempDir(), DaemonServiceName+".service")
	if content != "" {
		if err := os.WriteFile(daemonUnitFile, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestShippedDaemonUnit(t *testing.T) {
	got, err := ShippedDaemonUnit(fakeBinary(t, testUnit, 0))
	if err != nil || string(got) != testUnit {
		t.Fatalf("ShippedDaemonUnit = %q, %v", got, err)
	}
	for name, bin := range map[string]string{
		"fails":   fakeBinary(t, testUnit, 1),
		"no unit": fakeBinary(t, "Error: unknown command\n", 0),
		"too big": fakeBinary(t, testUnit+strings.Repeat("#\n", maxUnitSize), 0),
	} {
		if _, err := ShippedDaemonUnit(bin); err == nil {
			t.Errorf("%s: ShippedDaemonUnit accepted it", name)
		}
	}
}

func TestDesiredDaemonUnitKeepsMetadataChoice(t *testing.T) {
	if got := DesiredDaemonUnit([]byte(testUnit), []byte(testUnit)); string(got) != testUnit {
		t.Fatalf("default install: got %q", got)
	}
	allowed := DaemonUnitWithMetadataOutput([]byte(testUnit))
	got := DesiredDaemonUnit(allowed, []byte(testUnit))
	if strings.Contains(string(got), "--no-log-metadata-blocks") {
		t.Fatalf("--allow-metadata-output install lost its choice: %q", got)
	}
}

func TestPlanDaemonUnitUpdate(t *testing.T) {
	newUnit := strings.Replace(testUnit, "Description=x", "Description=y", 1)
	bin := fakeBinary(t, newUnit, 0)
	accept := func(_ string, _, _ []byte) (bool, error) { return true, nil }
	decline := func(_ string, _, _ []byte) (bool, error) { return false, nil }
	mustNotAsk := func(_ string, _, _ []byte) (bool, error) {
		t.Fatal("asked although nothing differs")
		return false, nil
	}

	withUnitFile(t, "")
	if got, err := PlanDaemonUnitUpdate(bin, mustNotAsk); got != nil || err != nil {
		t.Fatalf("no installed unit: %q, %v", got, err)
	}

	withUnitFile(t, newUnit)
	if got, err := PlanDaemonUnitUpdate(bin, mustNotAsk); got != nil || err != nil {
		t.Fatalf("identical unit: %q, %v", got, err)
	}

	withUnitFile(t, testUnit)
	if got, err := PlanDaemonUnitUpdate(bin, accept); err != nil || string(got) != newUnit {
		t.Fatalf("accepted: %q, %v", got, err)
	}
	if got, err := PlanDaemonUnitUpdate(bin, decline); got != nil || err != nil {
		t.Fatalf("declined: %q, %v", got, err)
	}
	abort := errors.New("ctrl-c")
	if _, err := PlanDaemonUnitUpdate(bin, func(_ string, _, _ []byte) (bool, error) { return false, abort }); !errors.Is(err, abort) {
		t.Fatalf("aborted prompt: err = %v", err)
	}
	if got, err := PlanDaemonUnitUpdate(fakeBinary(t, "", 1), mustNotAsk); got != nil || err != nil {
		t.Fatalf("binary without the command: %q, %v", got, err)
	}
}

func TestWriteDaemonUnit(t *testing.T) {
	withUnitFile(t, "old")
	if err := WriteDaemonUnit([]byte(testUnit)); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(daemonUnitFile)
	if err != nil || string(got) != testUnit {
		t.Fatalf("unit = %q, %v", got, err)
	}
	info, err := os.Stat(daemonUnitFile)
	if err != nil || info.Mode().Perm() != 0o644 {
		t.Fatalf("mode = %v, %v", info.Mode(), err)
	}
	entries, _ := os.ReadDir(filepath.Dir(daemonUnitFile))
	if len(entries) != 1 {
		t.Fatalf("staging file left behind: %v", entries)
	}
}
