// Package fscrypt file-vault tests. One test file per production source
// file is the convention here (see migrate_test.go's header): EVERY test
// for filevault.go lives in this one — crypto primitives, the recurring
// in-place unlock/lock cycle, and the Vault-method dispatch alike. Do not
// spawn variant files such as *_inplace_test.go or *_dispatch_test.go;
// append new cases below so related coverage stays together.
package fscrypt

import (
	"bytes"
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

////////////////////////////////////////////////////////////////////////////
// Shared helpers
////////////////////////////////////////////////////////////////////////////

// withMasterKey overrides the readMasterKey seam (MasterKeyFile is a real
// const pointing at /etc/app-listener/fscrypt.key — not writable/readable in
// a test) for the duration of the test and restores it after.
func withMasterKey(t *testing.T) []byte {
	t.Helper()
	key := bytes.Repeat([]byte{0x5A}, FscryptKeySize)
	withMasterKeyBytes(t, key)
	return key
}

// withMasterKeyBytes is withMasterKey but with a caller-chosen key.
func withMasterKeyBytes(t *testing.T, key []byte) {
	t.Helper()
	orig := readMasterKey
	readMasterKey = func() ([]byte, error) {
		return append([]byte(nil), key...), nil
	}
	t.Cleanup(func() { readMasterKey = orig })
}

// inode returns path's inode number — the load-bearing observable for the
// in-place tests below.
func inode(t *testing.T, path string) uint64 {
	t.Helper()
	info, err := os.Lstat(path)
	if err != nil {
		t.Fatalf("lstat %s: %v", path, err)
	}
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		t.Fatal("Stat_t unavailable on this platform")
	}
	return st.Ino
}

func mustSubkey(t *testing.T, master []byte) []byte {
	t.Helper()
	k, err := deriveFileVaultSubkey(master)
	if err != nil {
		t.Fatalf("deriveFileVaultSubkey: %v", err)
	}
	if len(k) != 32 {
		t.Fatalf("subkey length = %d, want 32", len(k))
	}
	return k
}

////////////////////////////////////////////////////////////////////////////
// Crypto primitives: subkey derivation, seal/open, tamper rejection
////////////////////////////////////////////////////////////////////////////

func TestDeriveFileVaultSubkeyDeterministicAndSeparated(t *testing.T) {
	masterA := bytes.Repeat([]byte{0xAA}, FscryptKeySize)
	masterB := bytes.Repeat([]byte{0xBB}, FscryptKeySize)

	k1 := mustSubkey(t, masterA)
	k2 := mustSubkey(t, masterA)
	if !bytes.Equal(k1, k2) {
		t.Error("subkey derivation is not deterministic for the same master key")
	}

	k3 := mustSubkey(t, masterB)
	if bytes.Equal(k1, k3) {
		t.Error("different master keys produced the same subkey")
	}

	// Domain separation: the subkey must never equal the raw master key
	// bytes it was derived from (the whole point of not reusing the fscrypt
	// key directly for a second, unrelated cipher construction).
	if bytes.Equal(k1, masterA) {
		t.Error("subkey must not equal the raw master key")
	}
}

func TestSealOpenFileVaultRoundTrip(t *testing.T) {
	subkey := mustSubkey(t, bytes.Repeat([]byte{0x42}, FscryptKeySize))
	plaintext := []byte("AutoLoginUser=someone\nRememberPassword=1\n")

	record, err := sealFileVault(subkey, plaintext)
	if err != nil {
		t.Fatalf("sealFileVault: %v", err)
	}
	if len(record) <= len(plaintext) {
		t.Fatalf("record (%d bytes) should carry header+tag overhead over plaintext (%d bytes)", len(record), len(plaintext))
	}

	got, err := openFileVault(subkey, record)
	if err != nil {
		t.Fatalf("openFileVault: %v", err)
	}
	if !bytes.Equal(got, plaintext) {
		t.Errorf("round trip mismatch: got %q, want %q", got, plaintext)
	}
}

func TestSealFileVaultNoncesNeverRepeat(t *testing.T) {
	subkey := mustSubkey(t, bytes.Repeat([]byte{0x07}, FscryptKeySize))
	plaintext := []byte("same plaintext every time")

	seen := map[string]bool{}
	for i := 0; i < 200; i++ {
		record, err := sealFileVault(subkey, plaintext)
		if err != nil {
			t.Fatalf("sealFileVault: %v", err)
		}
		nonce := string(record[len(fileVaultMagic)+1 : len(fileVaultMagic)+1+fileVaultNonceSize])
		if seen[nonce] {
			t.Fatalf("nonce reuse detected on iteration %d — catastrophic for GCM", i)
		}
		seen[nonce] = true
	}
}

func TestOpenFileVaultRejectsTampering(t *testing.T) {
	subkey := mustSubkey(t, bytes.Repeat([]byte{0x11}, FscryptKeySize))
	otherSubkey := mustSubkey(t, bytes.Repeat([]byte{0x99}, FscryptKeySize))
	record, err := sealFileVault(subkey, []byte("top secret"))
	if err != nil {
		t.Fatalf("sealFileVault: %v", err)
	}

	cases := map[string][]byte{
		"flipped ciphertext byte": func() []byte {
			r := append([]byte(nil), record...)
			r[len(r)-1] ^= 0xFF
			return r
		}(),
		"wrong subkey": record, // opened below with otherSubkey
		"truncated":    record[:len(record)-1],
		"empty":        nil,
		"too short":    record[:3],
		"bad magic":    append([]byte("XXXX"), record[4:]...),
		"bad version":  append(append([]byte(fileVaultMagic), 0xFF), record[5:]...),
	}
	for name, r := range cases {
		key := subkey
		if name == "wrong subkey" {
			key = otherSubkey
		}
		if _, err := openFileVault(key, r); err == nil {
			t.Errorf("%s: openFileVault should have failed", name)
		}
	}
}

func TestIsFileVaultRecordShape(t *testing.T) {
	subkey := mustSubkey(t, bytes.Repeat([]byte{0x22}, FscryptKeySize))
	record, err := sealFileVault(subkey, []byte("x"))
	if err != nil {
		t.Fatalf("sealFileVault: %v", err)
	}
	if !looksLikeFileVaultRecord(record) {
		t.Error("a genuine sealed record must look like one")
	}
	if looksLikeFileVaultRecord([]byte("plain old file content")) {
		t.Error("ordinary plaintext must not look like a vault record")
	}
	if looksLikeFileVaultRecord(nil) {
		t.Error("empty content must not look like a vault record")
	}
}

////////////////////////////////////////////////////////////////////////////
// Recurring in-place unlock/lock (unlockFileInPlace/lockFileInPlace)
////////////////////////////////////////////////////////////////////////////

// TestUnlockLockFileInPlaceRoundTrip: lock then unlock restores the original
// bytes, and the plaintext is never visible in the locked content — this
// one just proves the crypto+I/O round trip is correct (the inode-preserving
// property is the next test's job).
func TestUnlockLockFileInPlaceRoundTrip(t *testing.T) {
	withMasterKey(t)
	v := New()
	path := filepath.Join(t.TempDir(), "registry.vdf")
	original := []byte("AutoLoginUser=steamuser\n")
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatal(err)
	}

	if err := v.lockFileInPlace(path); err != nil {
		t.Fatalf("lockFileInPlace: %v", err)
	}
	locked, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !looksLikeFileVaultRecord(locked) {
		t.Fatalf("file does not look encrypted after lock: %q", locked)
	}
	if bytes.Contains(locked, original) {
		t.Error("plaintext is still visible in the locked file")
	}

	if err := v.unlockFileInPlace(path); err != nil {
		t.Fatalf("unlockFileInPlace: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, original) {
		t.Errorf("unlock did not restore the original content: got %q, want %q", got, original)
	}
}

// TestUnlockLockFileInPlacePreservesInode is the load-bearing regression
// test for the whole design: the guard's guard_inodes map is keyed by
// (dev,ino) and is populated once, before this ever runs. If a future
// change swaps in-place transform for a rename-based one, the inode WILL
// change and this test WILL catch it — that change would silently reopen
// the exact TOCTOU window this design exists to close (a freshly renamed-in
// inode is unrecognized by the guard until the next populate/sweep, and the
// LSM hooks default to allow what they do not recognize).
func TestUnlockLockFileInPlacePreservesInode(t *testing.T) {
	withMasterKey(t)
	v := New()
	path := filepath.Join(t.TempDir(), "localconfig.vdf")
	if err := os.WriteFile(path, []byte("token=abc123"), 0o600); err != nil {
		t.Fatal(err)
	}

	before := inode(t, path)

	if err := v.lockFileInPlace(path); err != nil {
		t.Fatalf("lockFileInPlace: %v", err)
	}
	if got := inode(t, path); got != before {
		t.Fatalf("inode changed across lockFileInPlace: %d -> %d", before, got)
	}

	if err := v.unlockFileInPlace(path); err != nil {
		t.Fatalf("unlockFileInPlace: %v", err)
	}
	if got := inode(t, path); got != before {
		t.Fatalf("inode changed across unlockFileInPlace: %d -> %d", before, got)
	}
}

// TestLockUnlockFileInPlaceIdempotent: locking an already-locked file (and
// unlocking an already-plaintext one) is a no-op, matching the directory
// Vault.Unlock/Lock "already provisioned"/"already locked" contract.
func TestLockUnlockFileInPlaceIdempotent(t *testing.T) {
	withMasterKey(t)
	v := New()
	path := filepath.Join(t.TempDir(), "config.vdf")
	if err := os.WriteFile(path, []byte("plain"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := v.unlockFileInPlace(path); err != nil {
		t.Fatalf("unlock on already-plaintext file should no-op, got: %v", err)
	}
	if got, _ := os.ReadFile(path); string(got) != "plain" {
		t.Errorf("no-op unlock must not touch content: %q", got)
	}

	if err := v.lockFileInPlace(path); err != nil {
		t.Fatalf("lockFileInPlace: %v", err)
	}
	locked, _ := os.ReadFile(path)
	if err := v.lockFileInPlace(path); err != nil {
		t.Fatalf("lock on already-locked file should no-op, got: %v", err)
	}
	again, _ := os.ReadFile(path)
	if !bytes.Equal(locked, again) {
		t.Error("no-op lock must not re-encrypt (would rotate the nonce for nothing)")
	}
}

// TestUnlockFileInPlaceWrongKeyFails: a different master key must not be
// able to unlock the file — fail closed, never fall back to treating it as
// plaintext.
func TestUnlockFileInPlaceWrongKeyFails(t *testing.T) {
	withMasterKey(t)
	v := New()
	path := filepath.Join(t.TempDir(), "id.json")
	original := []byte(`{"secret":"key"}`)
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := v.lockFileInPlace(path); err != nil {
		t.Fatal(err)
	}

	// Swap in a different master key.
	withMasterKeyBytes(t, bytes.Repeat([]byte{0x01}, FscryptKeySize))

	if err := v.unlockFileInPlace(path); err == nil {
		t.Fatal("unlock with the wrong master key must fail")
	}
	// The file must be left exactly as it was (still valid ciphertext) — a
	// failed unlock is not license to touch the content.
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !looksLikeFileVaultRecord(got) {
		t.Error("file must remain valid ciphertext after a failed unlock attempt")
	}
}

// TestUnlockFileInPlaceRecoversFromInterruptedTransform simulates a crash
// between staging the recovery sidecar and completing the in-place write:
// the next call must restore the pre-crash content and clean up, rather
// than leaving the file (or the sidecar) behind.
func TestUnlockFileInPlaceRecoversFromInterruptedTransform(t *testing.T) {
	withMasterKey(t)
	v := New()
	path := filepath.Join(t.TempDir(), "localconfig.vdf")
	original := []byte("pre-crash content")
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatal(err)
	}
	before := inode(t, path)

	// Simulate the crash: a recovery sidecar exists (staged before the
	// interrupted write) and the live file was left half-written (garbage).
	if err := stageRecovery(path, original); err != nil {
		t.Fatalf("stageRecovery: %v", err)
	}
	if err := transformFileInPlace(path, []byte("garbage-mid-write")); err != nil {
		t.Fatalf("transformFileInPlace: %v", err)
	}

	// The next unlock call (whatever the caller intended) must self-heal
	// before doing anything else.
	if err := v.unlockFileInPlace(path); err != nil {
		t.Fatalf("unlockFileInPlace after simulated crash: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, original) {
		t.Errorf("recovery did not restore the pre-crash content: got %q, want %q", got, original)
	}
	if _, err := os.Lstat(fileVaultRecoverPath(path)); !os.IsNotExist(err) {
		t.Error("recovery sidecar must be removed once healed")
	}
	if got := inode(t, path); got != before {
		t.Errorf("recovery must not change the live path's inode: %d -> %d", before, got)
	}
}

////////////////////////////////////////////////////////////////////////////
// Vault method dispatch: IsEncrypted/IsProvisioned/VerifyKey/Unlock/Lock
// route a regular-file target to the file vault instead of the kernel-
// fscrypt (directory) path, which would error out on a file with no
// kernel policy.
////////////////////////////////////////////////////////////////////////////

func TestIsEncryptedDispatchesToFileVaultForRegularFiles(t *testing.T) {
	withMasterKey(t)
	v := New()
	path := filepath.Join(t.TempDir(), "config.vdf")
	if err := os.WriteFile(path, []byte("plain content"), 0o600); err != nil {
		t.Fatal(err)
	}

	encrypted, err := v.IsEncrypted(path)
	if err != nil {
		t.Fatalf("IsEncrypted on plaintext file: %v", err)
	}
	if encrypted {
		t.Error("a plain file must not report as encrypted")
	}

	if err := v.lockFileInPlace(path); err != nil {
		t.Fatal(err)
	}
	encrypted, err = v.IsEncrypted(path)
	if err != nil {
		t.Fatalf("IsEncrypted on locked file: %v", err)
	}
	if !encrypted {
		t.Error("a file-vault-locked file must report as encrypted")
	}
}

func TestIsProvisionedForFile(t *testing.T) {
	withMasterKey(t)
	v := New()
	path := filepath.Join(t.TempDir(), "registry.vdf")
	if err := os.WriteFile(path, []byte("plain"), 0o600); err != nil {
		t.Fatal(err)
	}

	provisioned, err := v.IsProvisioned(path)
	if err != nil {
		t.Fatalf("IsProvisioned on plaintext file: %v", err)
	}
	if !provisioned {
		t.Error("a plaintext file is, by definition, currently readable/provisioned")
	}

	if err := v.lockFileInPlace(path); err != nil {
		t.Fatal(err)
	}
	provisioned, err = v.IsProvisioned(path)
	if err != nil {
		t.Fatalf("IsProvisioned on locked file: %v", err)
	}
	if provisioned {
		t.Error("a locked (ciphertext) file must report as not provisioned")
	}
}

func TestVerifyKeyForFile(t *testing.T) {
	withMasterKey(t)
	v := New()
	path := filepath.Join(t.TempDir(), "localconfig.vdf")
	if err := os.WriteFile(path, []byte("plain"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := v.VerifyKey(path); err != nil {
		t.Errorf("VerifyKey on an unencrypted file should be a no-op, got: %v", err)
	}

	if err := v.lockFileInPlace(path); err != nil {
		t.Fatal(err)
	}
	if err := v.VerifyKey(path); err != nil {
		t.Errorf("VerifyKey with the correct master key should succeed, got: %v", err)
	}

	withMasterKeyBytes(t, make([]byte, FscryptKeySize)) // wrong key
	if err := v.VerifyKey(path); err == nil {
		t.Error("VerifyKey with the wrong master key must fail")
	}
}

func TestVaultUnlockLockDispatchForRegularFile(t *testing.T) {
	withMasterKey(t)
	v := New()
	path := filepath.Join(t.TempDir(), "id.json")
	original := []byte(`{"secret":"key"}`)
	if err := os.WriteFile(path, original, 0o600); err != nil {
		t.Fatal(err)
	}
	before := inode(t, path)

	if err := v.Lock(path, false); err != nil {
		t.Fatalf("Vault.Lock on a regular file: %v", err)
	}
	locked, _ := os.ReadFile(path)
	if !looksLikeFileVaultRecord(locked) {
		t.Fatal("Vault.Lock did not route to the file vault")
	}

	if err := v.Unlock(path); err != nil {
		t.Fatalf("Vault.Unlock on a regular file: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(original) {
		t.Errorf("round trip via Vault.Unlock/Lock mismatch: got %q, want %q", got, original)
	}
	if after := inode(t, path); after != before {
		t.Errorf("Vault.Unlock/Lock changed the inode: %d -> %d", before, after)
	}
}
