// filevault.go is a userspace fallback for encrypting a single regular file, which kernel fscrypt
// can't do: FS_IOC_SET_ENCRYPTION_POLICY is directory-scoped (a file only inherits its directory's
// policy; see the classifySetupError/classifySupportError text in fscrypt.go).
//
// Deliberately NOT fscrypt's on-disk format (AES-XTS, kernel-only per-file nonce): nothing else
// will read this ciphertext, and XTS has no integrity. Instead a subkey is derived from the master
// key (HKDF, domain-separated, never the raw key) and content is sealed with AES-256-GCM, which
// authenticates for free.
//
// The per-boot unlock/lock cycle (unlockFileInPlace/lockFileInPlace) transforms content on the SAME
// inode, never via rename: that is what prevents a TOCTOU regression on the guard (see their docs).
// The same rule covers the crash-recovery sidecar (fileVaultRecoverSuffix): only ever rewritten in
// place on its existing inode, never created/renamed/removed once a guard could be watching, and it
// holds a sealed record, not plaintext.
package fscrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// readMasterKey is the seam for reading the master key; a package var only so tests can supply a
// throwaway key instead of /etc/app-listener/fscrypt.key.
var readMasterKey = readKey

// fileVaultRecoverSuffix names the crash-recovery sidecar written before every in-place transform.
//
// It is inside the guard's scope: a single-file watch root's guard also denies create/rename/delete
// of anything beside it in the PARENT directory (the rename-over-watchroot defense;
// guard_path_rename's destination-parent check in guard.bpf.c), so a create-temp-then-rename for
// the sidecar is denied once the guard is live. So this package only rewrites the sidecar's
// EXISTING inode in place (EnsureRecoverySidecarPlaceholder, stageRecovery, clearRecovery). Empty
// sidecar = nothing to recover (like editprotected.HashFile's empty-means-unset).
//
// The staged content is sealed under the master key, never raw plaintext: a sidecar left by a crash
// survives a reboot (which drops every BPF-LSM pin), so it must be useless on its own.
const fileVaultRecoverSuffix = ".app_listener.recover"

func fileVaultRecoverPath(path string) string { return path + fileVaultRecoverSuffix }

// rejectSidecarSymlink Lstats sidecar and errors if it is a symlink. The sidecar is only ever a
// plain file this package created, and every entry point (EnsureRecoverySidecarPlaceholder,
// writeSidecarInPlace, recoverFileInPlace) calls this FIRST: a symlink swapped in by whoever
// controls the parent directory while no guard is attached (e.g. between install and first daemon
// start) must not be followed into an arbitrary root-writable file. Missing is not an error
// (callers handle it).
func rejectSidecarSymlink(sidecar string) error {
	info, err := os.Lstat(sidecar)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("checking %s: %w", sidecar, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("refusing recovery sidecar %s: it is a symbolic link", sidecar)
	}
	return nil
}

// EnsureRecoverySidecarPlaceholder creates path's (empty) recovery sidecar if missing. Must run
// BEFORE path's guard attaches (see fileVaultRecoverSuffix); same reason and timing as
// selfguards.go's ensureHashFilePlaceholder. No-op for anything not currently a regular file.
func EnsureRecoverySidecarPlaceholder(path string) error {
	if !isRegularFileTarget(path) {
		return nil
	}
	sidecar := fileVaultRecoverPath(path)
	if err := rejectSidecarSymlink(sidecar); err != nil {
		return err
	}
	if _, err := os.Lstat(sidecar); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("checking %s: %w", sidecar, err)
	}
	f, err := os.OpenFile(sidecar, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return fmt.Errorf("creating %s: %w", sidecar, err)
	}
	return f.Close()
}

// classifyRegularFileTarget is the safety gate every file-vault entry point runs first: only a
// plain, non-symlink, single-hard-link regular file (same refusals as classifyMigrationTarget).
func classifyRegularFileTarget(path string) error {
	target, err := classifyMigrationTarget(path)
	if err != nil {
		return err
	}
	if target != targetRegularFile {
		return fmt.Errorf("%s is not a regular file", path)
	}
	return nil
}

// isRegularFileTarget is the cheap, lenient ROUTING check
// (IsEncrypted/IsProvisioned/VerifyKey/Unlock/Lock) choosing file-vault over the kernel-fscrypt
// directory path. Not the safety gate (classifyRegularFileTarget, used by every mutating function,
// is): on any doubt (symlink, stat error, special file) it returns false and the directory code
// runs unchanged.
func isRegularFileTarget(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

// sealFileVaultWithMasterKey derives the subkey from the current master key, seals plaintext, and
// wipes the subkey (NOT the caller's plaintext; callers wipe that). Shared by every mutating entry
// point.
func sealFileVaultWithMasterKey(plaintext []byte) ([]byte, error) {
	masterKey, err := readMasterKey()
	if err != nil {
		return nil, err
	}
	defer wipe(masterKey)
	subkey, err := deriveFileVaultSubkey(masterKey)
	if err != nil {
		return nil, err
	}
	defer wipe(subkey)
	return sealFileVault(subkey, plaintext)
}

// openFileVaultWithMasterKey derives the subkey from the current master key, opens record, and
// wipes the subkey.
func openFileVaultWithMasterKey(record []byte) ([]byte, error) {
	masterKey, err := readMasterKey()
	if err != nil {
		return nil, err
	}
	defer wipe(masterKey)
	subkey, err := deriveFileVaultSubkey(masterKey)
	if err != nil {
		return nil, err
	}
	defer wipe(subkey)
	return openFileVault(subkey, record)
}

// writeTempWithMetadata creates tmp (O_CREATE|O_EXCL, refusing to clobber an interrupted run's
// leftover), writes content, fsyncs, stamps src's mode/owner/xattrs onto it (src must still exist
// at its path; call before renaming src away), and closes it. Any failure removes tmp itself.
func writeTempWithMetadata(tmp, src string, content []byte, srcInfo os.FileInfo) error {
	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, srcInfo.Mode().Perm())
	if err != nil {
		return fmt.Errorf("create %s: %w", tmp, err)
	}
	if werr := func() error {
		if _, err := out.Write(content); err != nil {
			return fmt.Errorf("write %s: %w", tmp, err)
		}
		if err := out.Sync(); err != nil {
			return fmt.Errorf("sync %s: %w", tmp, err)
		}
		if err := stampCopyMetadata(out, src, tmp, srcInfo); err != nil {
			return err
		}
		return nil
	}(); werr != nil {
		_ = out.Close()
		_ = os.Remove(tmp)
		return werr
	}
	if err := out.Close(); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("close %s: %w", tmp, err)
	}
	return nil
}

// isFileVaultCiphertext reports whether path holds a valid file-vault header (magic + supported
// version): a cheap, key-free probe reading only the header. Used by IsEncrypted, deriving state
// from on-disk content rather than a side registry (as the kernel-fscrypt path does).
func isFileVaultCiphertext(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	buf := make([]byte, fileVaultHeaderSize)
	n, err := io.ReadFull(f, buf)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return false, fmt.Errorf("read %s: %w", path, err)
	}
	return looksLikeFileVaultRecord(buf[:n]), nil
}

// transformFileInPlace rewrites path with newContent on its EXISTING inode (open the existing fd,
// truncate, write, fsync; never unlink or rename the live path). The whole single-file design rests
// on this: guard_inodes' entry for path is keyed by (dev, ino) and stays valid, so new content is
// never reachable under an inode the guard doesn't know. A rename swap would present a new inode,
// and an unknown inode is default-ALLOW to the LSM hooks (they only enforce on known inodes).
func transformFileInPlace(path string, newContent []byte) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("truncate %s: %w", path, err)
	}
	if _, err := f.WriteAt(newContent, 0); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	if err := f.Sync(); err != nil {
		return fmt.Errorf("fsync %s: %w", path, err)
	}
	return nil
}

// writeSidecarInPlace rewrites sidecar with content (may be empty) on its EXISTING inode (open,
// truncate, write, fsync; never rename). The O_CREATE branch is a bootstrap fallback for when
// EnsureRecoverySidecarPlaceholder hasn't run (the install-time pass runs before any guard exists);
// once a guard watches path the sidecar always exists, so the open-existing branch is taken. Same
// two-branch shape and justification as editprotected.WriteHashFile.
func writeSidecarInPlace(sidecar string, content []byte) error {
	if err := rejectSidecarSymlink(sidecar); err != nil {
		return err
	}
	f, err := os.OpenFile(sidecar, os.O_RDWR, 0o600)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("opening %s: %w", sidecar, err)
		}
		f, err = os.OpenFile(sidecar, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err != nil {
			return fmt.Errorf("creating %s: %w", sidecar, err)
		}
	}
	defer f.Close()
	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("truncate %s: %w", sidecar, err)
	}
	if len(content) > 0 {
		if _, err := f.WriteAt(content, 0); err != nil {
			return fmt.Errorf("write %s: %w", sidecar, err)
		}
	}
	return f.Sync()
}

// stageRecovery seals path's CURRENT bytes under the master key (same AEAD format as the vault
// record, never raw plaintext) and writes them to the sidecar in place (writeSidecarInPlace). A
// sidecar left by a crash, even across a reboot dropping every BPF-LSM pin, is a ciphertext blob
// useless without the master key.
func stageRecovery(path string, current []byte) error {
	sealed, err := sealFileVaultWithMasterKey(current)
	if err != nil {
		return fmt.Errorf("sealing recovery sidecar for %s: %w", path, err)
	}
	if err := writeSidecarInPlace(fileVaultRecoverPath(path), sealed); err != nil {
		return fmt.Errorf("staging recovery sidecar for %s: %w", path, err)
	}
	return nil
}

// clearRecovery empties the sidecar in place (truncate to zero on the same inode, never unlink)
// after a successful transform. Zero length = "nothing to recover" (recoverFileInPlace).
func clearRecovery(path string) error {
	if err := writeSidecarInPlace(fileVaultRecoverPath(path), nil); err != nil {
		return fmt.Errorf("clearing recovery sidecar for %s: %w", path, err)
	}
	return nil
}

// recoverFileInPlace restores path from a non-empty recovery sidecar, i.e. a previous unlock/lock
// was interrupted mid-transform. Both unlockFileInPlace and lockFileInPlace call it first. No-op
// when the sidecar is empty or missing (shouldn't happen once a guard exists). Needs the master key
// like a normal unlock (the sidecar is sealed); the restore is an in-place same-inode write
// (transformFileInPlace).
func recoverFileInPlace(path string) error {
	sidecar := fileVaultRecoverPath(path)
	if err := rejectSidecarSymlink(sidecar); err != nil {
		return err
	}
	sealed, err := os.ReadFile(sidecar)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("reading recovery sidecar for %s: %w", path, err)
	}
	if len(sealed) == 0 {
		return nil
	}
	backup, err := openFileVaultWithMasterKey(sealed)
	if err != nil {
		return fmt.Errorf("opening recovery sidecar for %s: %w", path, err)
	}
	defer wipe(backup)
	if err := transformFileInPlace(path, backup); err != nil {
		return fmt.Errorf("restoring %s from its recovery sidecar: %w", path, err)
	}
	return clearRecovery(path)
}

// unlockFileInPlace decrypts the ciphertext at path back to plaintext in place
// (transformFileInPlace). No-op if already plaintext (like the directory Unlock's "already
// provisioned"). A failed decrypt (wrong key, tampering) leaves the file untouched, never treats it
// as plaintext.
func (v *Vault) unlockFileInPlace(path string) error {
	if err := classifyRegularFileTarget(path); err != nil {
		return err
	}
	if err := recoverFileInPlace(path); err != nil {
		return err
	}
	current, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if !looksLikeFileVaultRecord(current) {
		return nil
	}

	plaintext, err := openFileVaultWithMasterKey(current)
	if err != nil {
		return fmt.Errorf("unlock %s: %w", path, err)
	}
	defer wipe(plaintext)

	if err := stageRecovery(path, current); err != nil {
		return err
	}
	if err := transformFileInPlace(path, plaintext); err != nil {
		return err
	}
	return clearRecovery(path)
}

// lockFileInPlace encrypts the plaintext at path in place with a fresh random nonce. No-op if
// already ciphertext (re-sealing would rotate the nonce for nothing; same idempotence as the
// directory Lock).
func (v *Vault) lockFileInPlace(path string) error {
	if err := classifyRegularFileTarget(path); err != nil {
		return err
	}
	if err := recoverFileInPlace(path); err != nil {
		return err
	}
	current, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if looksLikeFileVaultRecord(current) {
		return nil
	}

	record, err := sealFileVaultWithMasterKey(current)
	if err != nil {
		return fmt.Errorf("lock %s: %w", path, err)
	}

	if err := stageRecovery(path, current); err != nil {
		return err
	}
	if err := transformFileInPlace(path, record); err != nil {
		return err
	}
	return clearRecovery(path)
}

// verifyFileVaultKey confirms the master key decrypts path's ciphertext without mutating anything.
// nil if path isn't vault ciphertext.
func (v *Vault) verifyFileVaultKey(path string) error {
	current, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	if !looksLikeFileVaultRecord(current) {
		return nil
	}
	plaintext, err := openFileVaultWithMasterKey(current)
	if err != nil {
		return fmt.Errorf("verify key for %s: %w", path, err)
	}
	wipe(plaintext)
	return nil
}

const (
	// fileVaultMagic identifies this package's own format, distinct from anything kernel fscrypt
	// writes (it never touches file content), so an unrelated file can't look like ours and
	// IsEncrypted is a pure content check.
	fileVaultMagic = "ALFV" // "App-Listener File Vault"
	// fileVaultVersion is the on-disk record format version.
	fileVaultVersion = 1
	// fileVaultNonceSize is the standard AES-GCM nonce size.
	fileVaultNonceSize = 12
	// fileVaultSubkeySize is the AES-256 key size.
	fileVaultSubkeySize = 32
	// fileVaultHeaderSize is magic + version + nonce, before ciphertext+tag.
	fileVaultHeaderSize = len(fileVaultMagic) + 1 + fileVaultNonceSize
	// fileVaultKeyInfo is the HKDF "info": fixed and versioned so this subkey can't collide with
	// any other derivation from the master key.
	fileVaultKeyInfo = "app-listener/file-vault/v1"
)

// deriveFileVaultSubkey derives the AES-256-GCM key from the fscrypt master key via HKDF-SHA256.
// The raw master key is never used as a cipher key for anything but the fscrypt policies it was
// generated for (domain separation rules out cross-protocol reuse).
func deriveFileVaultSubkey(masterKey []byte) ([]byte, error) {
	subkey, err := hkdf.Key(sha256.New, masterKey, nil, fileVaultKeyInfo, fileVaultSubkeySize)
	if err != nil {
		return nil, fmt.Errorf("derive file-vault subkey: %w", err)
	}
	return subkey, nil
}

// newFileVaultAEAD builds the AES-256-GCM instance for subkey.
func newFileVaultAEAD(subkey []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(subkey)
	if err != nil {
		return nil, fmt.Errorf("file-vault cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("file-vault AEAD: %w", err)
	}
	if gcm.NonceSize() != fileVaultNonceSize {
		// Guards against a stdlib behavior change silently breaking the on-disk format instead of
		// corrupting data.
		return nil, fmt.Errorf("file-vault AEAD: unexpected nonce size %d", gcm.NonceSize())
	}
	return gcm, nil
}

// sealFileVault encrypts plaintext under subkey with a fresh random nonce (crypto/rand, never
// reused) and returns the on-disk record: magic + version + nonce + ciphertext+tag.
func sealFileVault(subkey, plaintext []byte) ([]byte, error) {
	gcm, err := newFileVaultAEAD(subkey)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, fileVaultNonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("file-vault nonce: %w", err)
	}
	record := make([]byte, 0, fileVaultHeaderSize+len(plaintext)+gcm.Overhead())
	record = append(record, fileVaultMagic...)
	record = append(record, fileVaultVersion)
	record = append(record, nonce...)
	record = gcm.Seal(record, nonce, plaintext, nil)
	return record, nil
}

// openFileVault authenticates and decrypts a sealFileVault record. Tampering, wrong key, truncation
// or a malformed header is rejected outright (no partial-trust fallback).
func openFileVault(subkey, record []byte) ([]byte, error) {
	if len(record) < fileVaultHeaderSize {
		return nil, fmt.Errorf("file-vault record too short (%d bytes)", len(record))
	}
	if !bytes.Equal(record[:len(fileVaultMagic)], []byte(fileVaultMagic)) {
		return nil, fmt.Errorf("file-vault record has the wrong magic")
	}
	if record[len(fileVaultMagic)] != fileVaultVersion {
		return nil, fmt.Errorf("file-vault record has unsupported version %d", record[len(fileVaultMagic)])
	}
	gcm, err := newFileVaultAEAD(subkey)
	if err != nil {
		return nil, err
	}
	nonce := record[len(fileVaultMagic)+1 : fileVaultHeaderSize]
	ciphertext := record[fileVaultHeaderSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("file-vault: decryption failed (wrong key, or the file was tampered with)")
	}
	return plaintext, nil
}

// looksLikeFileVaultRecord reports whether content has this format's magic and a supported version:
// a cheap, key-free structural check used by IsEncrypted/IsProvisioned.
func looksLikeFileVaultRecord(content []byte) bool {
	return len(content) >= fileVaultHeaderSize &&
		bytes.Equal(content[:len(fileVaultMagic)], []byte(fileVaultMagic)) &&
		content[len(fileVaultMagic)] == fileVaultVersion
}
