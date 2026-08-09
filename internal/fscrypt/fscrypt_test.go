package fscrypt

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/fscrypt/actions"
	"github.com/google/fscrypt/filesystem"
)

// TestNewBoundedKeyFnFirstCall verifies that the initial callback call
// (retry=false) relays the master key untouched.
func TestNewBoundedKeyFnFirstCall(t *testing.T) {
	masterKey := make([]byte, FscryptKeySize)
	for i := range masterKey {
		masterKey[i] = byte(i)
	}

	key, err := newBoundedKeyFn(masterKey)(actions.ProtectorInfo{}, false)
	if err != nil {
		t.Fatalf("first call returned error: %v", err)
	}
	defer key.Wipe()
	if key.Len() != FscryptKeySize {
		t.Fatalf("key length = %d, want %d", key.Len(), FscryptKeySize)
	}
	if !bytes.Equal(key.Data(), masterKey) {
		t.Fatalf("key bytes = %v, want %v", key.Data(), masterKey)
	}
}

// TestNewBoundedKeyFnRetryAborts verifies the loop-breaking contract: once
// the unwrap loop signals retry=true, the callback must return an error
// (never the same wrong key again), otherwise google/fscrypt's
// unwrapProtectorKey spins forever logging "invalid wrapping key".
func TestNewBoundedKeyFnRetryAborts(t *testing.T) {
	masterKey := make([]byte, FscryptKeySize)
	for i := range masterKey {
		masterKey[i] = 0xAA
	}

	key, err := newBoundedKeyFn(masterKey)(actions.ProtectorInfo{}, true)
	if err == nil {
		key.Wipe()
		t.Fatal("retry call returned no error, want abort")
	}
	for _, want := range []string{"invalid wrapping key for protector", MasterKeyFile} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not mention %q", err, want)
		}
	}
}

// TestReadKeyFromMissingFile verifies that unlocking an already encrypted
// directory without a master key file fails fast instead of generating a
// fresh (useless) key.
func TestReadKeyFromMissingFile(t *testing.T) {
	_, err := readKeyFrom(filepath.Join(t.TempDir(), "fscrypt.key"))
	if err == nil {
		t.Fatal("expected error for missing key file")
	}
	if !strings.Contains(err.Error(), "read master key") {
		t.Errorf("error %q does not name the read step", err)
	}
}

// TestReadKeyFromBadLength verifies that a key file of the wrong size is
// rejected rather than fed into the unwrap logic.
func TestReadKeyFromBadLength(t *testing.T) {
	keyFile := filepath.Join(t.TempDir(), "fscrypt.key")
	if err := os.WriteFile(keyFile, make([]byte, 17), 0o400); err != nil {
		t.Fatal(err)
	}
	_, err := readKeyFrom(keyFile)
	if err == nil {
		t.Fatal("expected error for 17-byte key file")
	}
	if !strings.Contains(err.Error(), "must be exactly 32 bytes") {
		t.Errorf("error %q does not explain the size requirement", err)
	}
}

// TestReadKeyFromRoundTrip verifies that a well-formed key file is read
// back byte for byte.
func TestReadKeyFromRoundTrip(t *testing.T) {
	want := make([]byte, FscryptKeySize)
	for i := range want {
		want[i] = byte(255 - i)
	}
	keyFile := filepath.Join(t.TempDir(), "fscrypt.key")
	if err := os.WriteFile(keyFile, want, 0o400); err != nil {
		t.Fatal(err)
	}
	got, err := readKeyFrom(keyFile)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("read key = %v, want %v", got, want)
	}
}

// TestIsEncryptedRegularFile verifies that a plain file is reported as
// "not encrypted" (fscrypt policies exist only on directories) and that a
// missing path still surfaces an error.
func TestIsEncryptedRegularFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "secret")
	if err := os.WriteFile(file, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	encrypted, err := (&Vault{}).IsEncrypted(file)
	if err != nil {
		t.Fatalf("IsEncrypted on a file: %v", err)
	}
	if encrypted {
		t.Error("a regular file must never report as encrypted")
	}

	if _, err := (&Vault{}).IsEncrypted(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Error("IsEncrypted on a missing path must error")
	}
}

// TestClassifySetupErrorNotSetup verifies that the library's
// ErrNotSetup is translated into an error naming the exact remediation
// (run `fscrypt setup` once) so the installer can fail fast with guidance.
func TestClassifySetupErrorNotSetup(t *testing.T) {
	err := classifySetupError("/home/test/.ssh",
		&filesystem.ErrNotSetup{Mount: &filesystem.Mount{Path: "/"}})
	if err == nil {
		t.Fatal("expected an error")
	}
	for _, want := range []string{"/home/test/.ssh", "fscrypt setup --all-users", "not initialized"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not mention %q", err, want)
		}
	}
}

// TestClassifySetupErrorNotSupported verifies that unsupported filesystem
// types (tmpfs, FAT/exFAT) get their own distinct explanation.
func TestClassifySetupErrorNotSupported(t *testing.T) {
	err := classifySetupError("/tmp/foo",
		&filesystem.ErrSetupNotSupported{Mount: &filesystem.Mount{Path: "/tmp"}})
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "does not support fscrypt encryption") {
		t.Errorf("error %q does not explain the unsupported filesystem", err)
	}
}

// TestClassifySetupErrorGeneric verifies that unrelated failures pass
// through wrapped with context instead of being mislabeled.
func TestClassifySetupErrorGeneric(t *testing.T) {
	inner := errors.New("metadata corrupted")
	err := classifySetupError("/srv/data", inner)
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "fscrypt setup check for /srv/data") {
		t.Errorf("error %q does not carry the check context", err)
	}
	if !errors.Is(err, inner) {
		t.Errorf("error %q must wrap the original error", err)
	}
}

// TestClassifySupportErrorEncryptionNotEnabled verifies that a filesystem
// set up for fscrypt but lacking the ext4 `encrypt` feature flag gets the
// exact remediation (tune2fs -O encrypt) instead of a cryptic failure deep
// inside the library.
func TestClassifySupportErrorEncryptionNotEnabled(t *testing.T) {
	err := classifySupportError("/home/test/.ssh",
		&filesystem.ErrEncryptionNotEnabled{Mount: &filesystem.Mount{Path: "/", Device: "/dev/sda2", FilesystemType: "ext4"}})
	if err == nil {
		t.Fatal("expected an error")
	}
	for _, want := range []string{"/home/test/.ssh", "/dev/sda2", "tune2fs -O encrypt /dev/sda2", "unmounted"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not mention %q", err, want)
		}
	}
}

// TestClassifySupportErrorEncryptionNotEnabledF2fs verifies that f2fs
// filesystems get the f2fs-specific feature-enabling command.
func TestClassifySupportErrorEncryptionNotEnabledF2fs(t *testing.T) {
	err := classifySupportError("/srv/data",
		&filesystem.ErrEncryptionNotEnabled{Mount: &filesystem.Mount{Path: "/srv", Device: "/dev/sdb1", FilesystemType: "f2fs"}})
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "fsck.f2fs -O encrypt /dev/sdb1") {
		t.Errorf("error %q does not mention the f2fs enable command", err)
	}
	if strings.Contains(err.Error(), "tune2fs") {
		t.Errorf("error %q must not suggest tune2fs for f2fs", err)
	}
}

// TestClassifySupportErrorNotSupported verifies that a kernel without
// fscrypt support for the filesystem type is reported as such.
func TestClassifySupportErrorNotSupported(t *testing.T) {
	err := classifySupportError("/srv/data",
		&filesystem.ErrEncryptionNotSupported{Mount: &filesystem.Mount{Path: "/srv", FilesystemType: "xfs"}})
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "kernel") || !strings.Contains(err.Error(), "xfs") {
		t.Errorf("error %q does not explain the kernel limitation", err)
	}
}

// TestClassifySupportErrorGeneric verifies that unrelated failures pass
// through wrapped with context.
func TestClassifySupportErrorGeneric(t *testing.T) {
	inner := errors.New("probe failed")
	err := classifySupportError("/srv/data", inner)
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "fscrypt support check for /srv/data") {
		t.Errorf("error %q does not carry the check context", err)
	}
	if !errors.Is(err, inner) {
		t.Errorf("error %q must wrap the original error", err)
	}
}
