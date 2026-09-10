package install

import (
	"strings"
	"testing"
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
