package install

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

// writeCmdline creates a fake /proc/<pid>/cmdline entry with the given
// NUL-separated arguments.
func writeCmdline(t *testing.T, procRoot, pid string, args ...string) {
	t.Helper()
	dir := filepath.Join(procRoot, pid)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	data := make([]byte, 0, 256)
	for i, a := range args {
		if i > 0 {
			data = append(data, 0)
		}
		data = append(data, a...)
	}
	if err := os.WriteFile(filepath.Join(dir, "cmdline"), data, 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestFindDaemonProcesses exercises the manual-daemon detection used to
// refuse installation when a daemon is running outside systemd.
func TestFindDaemonProcesses(t *testing.T) {
	proc := t.TempDir()
	writeCmdline(t, proc, "100", "/usr/local/sbin/app-listener", "daemon", "--headless", "--blocked-only")
	writeCmdline(t, proc, "200", "/usr/local/sbin/app-listener", "install")
	writeCmdline(t, proc, "300", "/usr/bin/something", "daemon")
	writeCmdline(t, proc, "400", "/usr/local/sbin/app-listener", "daemon")

	pids, err := findDaemonProcesses(proc)
	if err != nil {
		t.Fatalf("findDaemonProcesses: %v", err)
	}
	if !reflect.DeepEqual(pids, []int{100, 400}) {
		t.Errorf("pids = %v, want [100 400]", pids)
	}
}

// TestFindDaemonProcessesEmpty verifies an empty process tree yields no PIDs.
func TestFindDaemonProcessesEmpty(t *testing.T) {
	pids, err := findDaemonProcesses(t.TempDir())
	if err != nil {
		t.Fatalf("findDaemonProcesses: %v", err)
	}
	if len(pids) != 0 {
		t.Errorf("pids = %v, want none", pids)
	}
}

// TestIsDaemonCmdline covers the classifier directly.
func TestIsDaemonCmdline(t *testing.T) {
	cases := []struct {
		cmdline string
		want    bool
	}{
		{"app-listener\x00daemon\x00--headless", true},
		{"/usr/local/sbin/app-listener\x00daemon", true},
		{"app-listener\x00install", false},
		{"app-listener", false},
		{"", false},
		{"bash\x00daemon", false},
	}
	for _, c := range cases {
		if got := isDaemonCmdline(c.cmdline); got != c.want {
			t.Errorf("isDaemonCmdline(%q) = %v, want %v", c.cmdline, got, c.want)
		}
	}
}

// TestVerifyResourcesLocked exercises the post-stop lockdown verification:
// every encrypted resource must be keyless, unencrypted ones are skipped,
// and a still-unlocked directory is a fatal error.
func TestVerifyResourcesLocked(t *testing.T) {
	locked := func(path string) (bool, error) { return false, nil }
	unlocked := func(path string) (bool, error) { return true, nil }
	checkErr := func(path string) (bool, error) { return false, errors.New("boom") }

	// Everything locked: pass.
	if err := verifyResourcesLocked([]daemonconfig.Resource{
		{Path: "/enc/a", NeedEncryption: true},
		{Path: "/enc/b", NeedEncryption: true},
	}, locked); err != nil {
		t.Errorf("all locked: unexpected error: %v", err)
	}

	// One still unlocked: fatal error naming the directory.
	err := verifyResourcesLocked([]daemonconfig.Resource{
		{Path: "/enc/a", NeedEncryption: true},
		{Path: "/enc/b", NeedEncryption: true},
	}, unlocked)
	if err == nil {
		t.Fatal("expected error when a resource is still unlocked")
	}
	if !strings.Contains(err.Error(), "/enc/b") {
		t.Errorf("error must name the unlocked resource, got: %v", err)
	}

	// Unencrypted resources are never checked, even when "unlocked".
	if err := verifyResourcesLocked([]daemonconfig.Resource{
		{Path: "/plain", NeedEncryption: false},
	}, unlocked); err != nil {
		t.Errorf("unencrypted resource must be skipped: %v", err)
	}

	// Empty configuration: nothing to verify.
	if err := verifyResourcesLocked(nil, unlocked); err != nil {
		t.Errorf("empty configuration: unexpected error: %v", err)
	}

	// A failing lock-state check is propagated, not swallowed.
	if err := verifyResourcesLocked([]daemonconfig.Resource{
		{Path: "/enc/a", NeedEncryption: true},
	}, checkErr); err == nil || !strings.Contains(err.Error(), "/enc/a") {
		t.Errorf("check failure must be wrapped with the path, got: %v", err)
	}
}
