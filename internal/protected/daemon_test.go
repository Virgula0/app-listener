package protected

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
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

	pids, err := FindDaemonProcesses(proc)
	if err != nil {
		t.Fatalf("FindDaemonProcesses: %v", err)
	}
	if !reflect.DeepEqual(pids, []int{100, 400}) {
		t.Errorf("pids = %v, want [100 400]", pids)
	}
}

// TestFindDaemonProcessesEmpty verifies an empty process tree yields no PIDs.
func TestFindDaemonProcessesEmpty(t *testing.T) {
	pids, err := FindDaemonProcesses(t.TempDir())
	if err != nil {
		t.Fatalf("FindDaemonProcesses: %v", err)
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
