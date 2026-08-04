package fscrypt

import (
	"os"
	"path/filepath"
	"testing"
)

// TestEncryptRollbackOnUnsupportedFilesystem exercises the migration on a
// filesystem without fscrypt support (t.TempDir usually lives on tmpfs).
// The policy application is expected to fail and the original directory
// must be restored with all contents — no data left in the backup-only
// staging state. When the filesystem happens to support fscrypt (and the
// environment has the privileges), the migration may actually succeed and
// the test skips.
func TestEncryptRollbackOnUnsupportedFilesystem(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "vault")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	inner := filepath.Join(dir, "nested")
	if err := os.MkdirAll(inner, 0o750); err != nil {
		t.Fatal(err)
	}
	secret := filepath.Join(inner, "secret.key")
	if err := os.WriteFile(secret, []byte("rollback-me"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := (&Vault{}).Encrypt(dir)
	if err == nil {
		t.Skip("filesystem supports fscrypt: migration succeeded, rollback not exercised")
	}

	backup := dir + BackupSuffix
	if _, statErr := os.Stat(backup); statErr == nil {
		t.Fatalf("backup %s left behind after failed migration (error: %v)", backup, err)
	}
	data, readErr := os.ReadFile(secret)
	if readErr != nil {
		t.Fatalf("original content lost after failed migration: %v (error: %v)", readErr, err)
	}
	if string(data) != "rollback-me" {
		t.Errorf("original content corrupted: %q", data)
	}
}
