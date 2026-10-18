package protected

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func withCertPath(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "install-complete")
	prev := installCertPath
	installCertPath = path
	t.Cleanup(func() { installCertPath = prev })
	return path
}

func TestInstallCertRoundTrip(t *testing.T) {
	path := withCertPath(t)
	if c, err := ReadInstallCert(); err != nil || c != nil {
		t.Fatalf("no record: got %v, %v", c, err)
	}
	if err := RequireNoInstallCert(); err != nil {
		t.Fatalf("no record must allow the wizard: %v", err)
	}
	if err := WriteInstallCert("v1.2.3"); err != nil {
		t.Fatal(err)
	}
	c, err := ReadInstallCert()
	if err != nil || c == nil || c.Version != "v1.2.3" || c.CompletedAt.IsZero() {
		t.Fatalf("got %+v, %v", c, err)
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0o644 {
		t.Fatalf("mode: %v, %v", info, err)
	}
	err = RequireNoInstallCert()
	if err == nil || !strings.Contains(err.Error(), "already installed") || !strings.Contains(err.Error(), "v1.2.3") {
		t.Fatalf("a record must refuse the wizard: %v", err)
	}
	if err := RemoveInstallCert(); err != nil {
		t.Fatal(err)
	}
	if err := RemoveInstallCert(); err != nil {
		t.Fatalf("removing a missing record: %v", err)
	}
	if err := RequireNoInstallCert(); err != nil {
		t.Fatalf("after removal: %v", err)
	}
}

func TestInstallCertUnreadableRefuses(t *testing.T) {
	path := withCertPath(t)
	if err := os.WriteFile(path, []byte("{torn"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := RequireNoInstallCert(); err == nil {
		t.Fatal("a corrupt record must refuse the wizard")
	}
}

func TestInstallCertSymlinkRefuses(t *testing.T) {
	path := withCertPath(t)
	if err := os.Symlink("/etc/hostname", path); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadInstallCert(); err == nil {
		t.Fatal("a symlinked record must not be followed")
	}
}
