package install

import (
	"bytes"
	"strings"
	"testing"

	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/systemd"
)

// TestSubcommandsRegistered pins the tree, and the hidden spelling an installer of another version
// runs on the installed binary.
func TestSubcommandsRegistered(t *testing.T) {
	for _, path := range [][]string{{"binary"}, {"catalog", "refresh"}, {"catalog", "diff"},
		{"backups", "restore"}, {"backups", "delete"}} {
		if c, _, err := InstallCmd.Find(path); err != nil || c.Name() != path[len(path)-1] {
			t.Errorf("install %s: %v", strings.Join(path, " "), err)
		}
	}
	if f := InstallCmd.Flags().Lookup("mark-complete"); f == nil || f.Hidden {
		t.Error("install --mark-complete must stay registered and listed")
	}
}

func TestDaemonUnitWithMetadataOutput(t *testing.T) {
	unit, err := inst.SampleContent("app-listener-daemon.service")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(unit, []byte("--no-log-metadata-blocks")) {
		t.Fatal("the bundled unit is expected to quiet metadata denials by default")
	}
	got := systemd.DaemonUnitWithMetadataOutput(unit)
	if bytes.Contains(got, []byte("--no-log-metadata-blocks")) {
		t.Fatalf("flag still present:\n%s", got)
	}
	if !bytes.Contains(got, []byte("ExecStart=/usr/local/sbin/app-listener daemon --headless --blocked-only\n")) {
		t.Fatalf("ExecStart mangled:\n%s", got)
	}
	if InstallCmd.Flags().Lookup("allow-metadata-output") == nil {
		t.Fatal("--allow-metadata-output not registered")
	}
}
