package fscrypt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestEncryptRefusesExistingBackup is the safety-critical contract: when a
// previous migration left a backup behind, Encrypt must fail fatally
// BEFORE touching anything, so an older backup is never overwritten.
func TestEncryptRefusesExistingBackup(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "vault")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	backup := dir + BackupSuffix
	if err := os.MkdirAll(backup, 0o700); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(backup, "precious")
	if err := os.WriteFile(marker, []byte("do not touch"), 0o600); err != nil {
		t.Fatal(err)
	}

	vault := &Vault{}
	err := vault.Encrypt(dir)
	if err == nil {
		t.Fatal("expected fatal error for existing backup")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("error does not explain the conflict: %v", err)
	}

	// The directory must be untouched.
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("original directory was modified: %v", err)
	}
	if data, err := os.ReadFile(marker); err != nil || string(data) != "do not touch" {
		t.Errorf("older backup was modified: %v %q", err, data)
	}
}

// TestVerifyKeyOnPlainDirectory fails without side effects: the path is not
// encrypted, so there is no policy to verify against.
func TestVerifyKeyOnPlainDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "plain")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := (&Vault{}).VerifyKey(dir); err == nil {
		t.Fatal("expected error verifying a non-encrypted directory")
	}
}

// TestStripImmutableFlagsSkipsDanglingSymlink is the firefox regression:
// profiles contain a `lock` symlink pointing to a remote-debugging address
// that does not exist when the browser is closed. Opening it with
// O_RDONLY used to fail with ENOENT and abort the migration; the walk must
// skip symlinks and still reach the regular files.
func TestStripImmutableFlagsSkipsDanglingSymlink(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "profile")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "prefs.js"), []byte("// nothing"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("192.168.1.91:+15894", filepath.Join(dir, "lock")); err != nil {
		t.Fatal(err)
	}

	if err := stripImmutableFlags(dir); err != nil {
		t.Fatalf("stripImmutableFlags must tolerate a dangling symlink, got: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "prefs.js")); err != nil {
		t.Errorf("regular file was lost during the walk: %v", err)
	}
}

// TestRestoreBackupMovesBackup covers the interrupted-migration case: the
// encrypted directory never existed (or was already removed), so restoring
// only moves the backup back into place.
func TestRestoreBackupMovesBackup(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "vault")
	backup := path + BackupSuffix
	if err := os.MkdirAll(backup, 0o700); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(backup, "precious")
	if err := os.WriteFile(marker, []byte("original content"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := (&Vault{}).RestoreBackup(path); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(backup); !os.IsNotExist(err) {
		t.Errorf("backup still exists after restore: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(path, "precious"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "original content" {
		t.Errorf("restored content = %q", data)
	}
}

// TestRestoreBackupRefusesNonEncrypted is the safety contract: when the
// target exists but is NOT encrypted, RestoreBackup must refuse instead of
// deleting whatever the user has there now.
func TestRestoreBackupRefusesNonEncrypted(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "vault")
	if err := os.MkdirAll(path, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(path, "new-data"), []byte("created later"), 0o600); err != nil {
		t.Fatal(err)
	}
	backup := path + BackupSuffix
	if err := os.MkdirAll(backup, 0o700); err != nil {
		t.Fatal(err)
	}

	err := (&Vault{}).RestoreBackup(path)
	if err == nil {
		t.Fatal("expected refusal for a non-encrypted target")
	}
	if !strings.Contains(err.Error(), "not encrypted") {
		t.Errorf("error does not explain the refusal: %v", err)
	}
	if _, err := os.Stat(filepath.Join(path, "new-data")); err != nil {
		t.Errorf("non-encrypted directory was touched: %v", err)
	}
	if _, err := os.Lstat(backup); err != nil {
		t.Errorf("backup was touched: %v", err)
	}
}

// TestRestoreBackupMissingBackup fails without side effects when no backup
// exists.
func TestRestoreBackupMissingBackup(t *testing.T) {
	path := filepath.Join(t.TempDir(), "vault")
	if err := (&Vault{}).RestoreBackup(path); err == nil {
		t.Fatal("expected error for missing backup")
	}
}

// TestDecryptRefusesExistingTempDir is the safety-critical contract of the
// permanent decryption: when a previous decryption left its temporary
// plaintext directory behind, Decrypt must fail fatally BEFORE touching
// anything, so a half-finished copy is never overwritten or silently
// discarded.
func TestDecryptRefusesExistingTempDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "vault")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	tmp := dir + DecryptSuffix
	if err := os.MkdirAll(tmp, 0o700); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(tmp, "partial")
	if err := os.WriteFile(marker, []byte("do not touch"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := (&Vault{}).Decrypt(dir)
	if err == nil {
		t.Fatal("expected fatal error for existing temporary directory")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("error does not explain the conflict: %v", err)
	}

	// The encrypted directory must be untouched.
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("original directory was modified: %v", err)
	}
	if data, err := os.ReadFile(marker); err != nil || string(data) != "do not touch" {
		t.Errorf("older temporary copy was modified: %v %q", err, data)
	}
}

// TestDecryptRefusesPlainDirectory fails without side effects: the path
// carries no fscrypt policy, so there is nothing to decrypt.
func TestDecryptRefusesPlainDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "plain")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	marker := filepath.Join(dir, "data")
	if err := os.WriteFile(marker, []byte("content"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := (&Vault{}).Decrypt(dir)
	if err == nil {
		t.Fatal("expected error decrypting a non-encrypted directory")
	}
	if !strings.Contains(err.Error(), "not encrypted") {
		t.Errorf("error does not explain the refusal: %v", err)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Errorf("plain directory was touched: %v", err)
	}
	if _, err := os.Lstat(dir + DecryptSuffix); !os.IsNotExist(err) {
		t.Errorf("no temporary directory may be created: %v", err)
	}
}

// TestDecryptMissingPath fails fast with a stat error and leaves nothing
// behind.
func TestDecryptMissingPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing")
	if err := (&Vault{}).Decrypt(path); err == nil {
		t.Fatal("expected error for missing path")
	}
	if _, err := os.Lstat(path + DecryptSuffix); !os.IsNotExist(err) {
		t.Errorf("no temporary directory may be created: %v", err)
	}
}
