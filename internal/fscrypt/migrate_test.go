// Package fscrypt migration tests. One test file per production source
// file is the convention here: EVERY test for migrate.go lives in this
// file — directory migrations, single-file migrations and rollback cases
// alike. Do not spawn variant files such as *_file_test.go or
// *_rollback_test.go; append new cases below so related coverage stays
// together.
package fscrypt

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"golang.org/x/sys/unix"
)

// ------------------------------------------------------------------
// Directory migration
// ------------------------------------------------------------------

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

// ------------------------------------------------------------------
// Single-file migration
// ------------------------------------------------------------------

// TestEncryptFileRefusesExistingBackup is the same safety-critical contract
// as the directory variant: an older backup must never be overwritten.
func TestEncryptFileRefusesExistingBackup(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "secret.env")
	if err := os.WriteFile(path, []byte("KEY=value"), 0o600); err != nil {
		t.Fatal(err)
	}
	backup := path + BackupSuffix
	if err := os.WriteFile(backup, []byte("older backup"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := (&Vault{}).Encrypt(path)
	if err == nil {
		t.Fatal("expected fatal error for existing backup")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("error does not explain the conflict: %v", err)
	}
	data, readErr := os.ReadFile(path)
	if readErr != nil || string(data) != "KEY=value" {
		t.Errorf("original was modified: %v %q", readErr, data)
	}
}

// TestEncryptFileRefusesSymlink: following a link could silently migrate an
// unrelated target and breaks the rename-based swap.
func TestEncryptFileRefusesSymlink(t *testing.T) {
	base := t.TempDir()
	target := filepath.Join(base, "real")
	if err := os.WriteFile(target, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(base, "link")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}

	err := (&Vault{}).Encrypt(link)
	if err == nil || !strings.Contains(err.Error(), "symbolic link") {
		t.Fatalf("err = %v, want symbolic-link refusal", err)
	}
}

// TestEncryptFileRefusesHardlink: migrating one alias must never implicitly
// affect another path to the same inode.
func TestEncryptFileRefusesHardlink(t *testing.T) {
	base := t.TempDir()
	a := filepath.Join(base, "a")
	if err := os.WriteFile(a, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	b := filepath.Join(base, "b")
	if err := os.Link(a, b); err != nil {
		t.Skipf("hardlinks unavailable here: %v", err)
	}

	for _, p := range []string{a, b} {
		err := (&Vault{}).Encrypt(p)
		if err == nil || !strings.Contains(err.Error(), "hard links") {
			t.Fatalf("Encrypt(%s) err = %v, want hard-link refusal", p, err)
		}
	}
}

// TestEncryptFileRefusesSpecialFile covers FIFOs; sockets and devices are
// refused by the same classification.
func TestEncryptFileRefusesSpecialFile(t *testing.T) {
	fifo := filepath.Join(t.TempDir(), "pipe")
	if err := syscall.Mkfifo(fifo, 0o600); err != nil {
		t.Skipf("mkfifo unavailable here: %v", err)
	}

	err := (&Vault{}).Encrypt(fifo)
	if err == nil || !strings.Contains(err.Error(), "regular files") {
		t.Fatalf("err = %v, want special-file refusal", err)
	}
}

// TestEncryptFileRollbackOnUnsupportedFilesystem mirrors the directory
// rollback test: on a filesystem without fscrypt support (t.TempDir is
// usually tmpfs) the migration fails before any rename and the original
// file plus its metadata must be untouched, with no staging leftovers.
func TestEncryptFileRollbackOnUnsupportedFilesystem(t *testing.T) {
	path := filepath.Join(t.TempDir(), "secret.env")
	content := "rollback-me"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}

	err := (&Vault{}).Encrypt(path)
	if err == nil {
		t.Skip("filesystem supports fscrypt: migration succeeded, rollback not exercised")
	}

	if _, statErr := os.Lstat(path + BackupSuffix); statErr == nil {
		t.Fatalf("backup left behind after failed migration (error: %v)", err)
	}
	if _, statErr := os.Lstat(path + encryptTmpSuffix); statErr == nil {
		t.Fatalf("temp file left behind after failed migration (error: %v)", err)
	}
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		t.Fatalf("original content lost after failed migration: %v (error: %v)", readErr, err)
	}
	if string(data) != content {
		t.Errorf("original content corrupted: %q", data)
	}
}

// TestDecryptFileRefusesNonEncrypted: a plain file has no policy to remove.
func TestDecryptFileRefusesNonEncrypted(t *testing.T) {
	path := filepath.Join(t.TempDir(), "plain.env")
	if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := (&Vault{}).Decrypt(path)
	if err == nil || !strings.Contains(err.Error(), "not encrypted") {
		t.Fatalf("err = %v, want not-encrypted refusal", err)
	}
}

// TestRestoreBackupMovesFileBackup covers the interrupted single-file
// migration case: the encrypted file never existed (or was already
// removed), so restoring only moves the backup back into place.
func TestRestoreBackupMovesFileBackup(t *testing.T) {
	base := t.TempDir()
	path := filepath.Join(base, "secret.env")
	backup := path + BackupSuffix
	if err := os.WriteFile(backup, []byte("original content"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := (&Vault{}).RestoreBackup(path); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Lstat(backup); !os.IsNotExist(err) {
		t.Errorf("backup still exists after restore: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "original content" {
		t.Errorf("restored file = %q (err %v), want original content", data, err)
	}
}

// TestCopyXattrsPreservesUserAttributes verifies that a user.* xattr set on
// the source reappears on the destination. Filesystems without user-namespace
// xattr support skip the test.
func TestCopyXattrsPreservesUserAttributes(t *testing.T) {
	base := t.TempDir()
	src := filepath.Join(base, "src")
	dst := filepath.Join(base, "dst")
	for _, p := range []string{src, dst} {
		if err := os.WriteFile(p, []byte("x"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	const name = "user.app_listener.test"
	if err := unix.Lsetxattr(src, name, []byte("hello"), 0); err != nil {
		t.Skipf("user xattrs unsupported here: %v", err)
	}

	if err := copyXattrs(src, dst); err != nil {
		t.Fatalf("copyXattrs: %v", err)
	}
	size, err := unix.Lgetxattr(dst, name, nil)
	if err != nil {
		t.Fatalf("xattr missing on dst: %v", err)
	}
	value := make([]byte, size)
	n, err := unix.Lgetxattr(dst, name, value)
	if err != nil {
		t.Fatal(err)
	}
	if string(value[:n]) != "hello" {
		t.Errorf("xattr value = %q, want hello", value[:n])
	}
}
