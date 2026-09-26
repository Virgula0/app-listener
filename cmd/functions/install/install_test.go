package install

import (
	"bytes"
	"strings"
	"testing"

	inst "github.com/Virgula0/app-listener/internal/install"
)

func TestValidateMaintenanceFlags(t *testing.T) {
	cases := []struct {
		name    string
		f       maintenanceFlags
		wantErr string // substring; "" = no error
	}{
		{"none", maintenanceFlags{}, ""},
		{"diff alone", maintenanceFlags{diffCatalog: true}, ""},
		{"update alone", maintenanceFlags{updateCatalog: true}, ""},
		{"diff + update", maintenanceFlags{diffCatalog: true, updateCatalog: true}, "mutually exclusive"},
		{"diff + binary-only", maintenanceFlags{diffCatalog: true, binaryOnly: true}, "mutually exclusive"},
		{"diff + yes", maintenanceFlags{diffCatalog: true, autoConfirm: true}, "--yes can only be used"},
		{"diff + live", maintenanceFlags{diffCatalog: true, live: true}, "--live can only be used"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := validateMaintenanceFlags(c.f)
			switch {
			case c.wantErr == "" && err != nil:
				t.Errorf("unexpected error: %v", err)
			case c.wantErr != "" && (err == nil || !strings.Contains(err.Error(), c.wantErr)):
				t.Errorf("error = %v, want substring %q", err, c.wantErr)
			}
		})
	}
}

// TestDiffCatalogFlagRegistered guards against the flag being dropped from
// the command wiring.
func TestDiffCatalogFlagRegistered(t *testing.T) {
	if InstallCmd.Flags().Lookup("diff-catalog") == nil {
		t.Fatal("install --diff-catalog flag is not registered")
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
	got := daemonUnitWithMetadataOutput(unit)
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
