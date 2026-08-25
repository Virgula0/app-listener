package fscrypt

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"golang.org/x/sys/unix"
)

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
