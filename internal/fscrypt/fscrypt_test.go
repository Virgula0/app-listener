package fscrypt

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/fscrypt/actions"
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
