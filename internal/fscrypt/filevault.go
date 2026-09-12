// filevault.go implements a userspace fallback for encrypting a single
// regular file, for the class of catalog resources fscrypt itself cannot
// cover: FS_IOC_SET_ENCRYPTION_POLICY is directory-scoped in the Linux
// kernel — a standalone regular file can never be given its own policy
// directly, only inherit one from the directory it is created in. (The
// google/fscrypt library's own SetPolicy doc says "sets up the specified
// DIRECTORY", its CLI has no per-file command, and its own suggested
// workaround for "can't encrypt in place" is wrapping the file in a fresh
// encrypted directory — see the classifySetupError/classifySupportError
// remediation text in fscrypt.go for what actually happens if this is
// attempted through the kernel ioctl.)
//
// This is deliberately NOT an attempt to reproduce fscrypt's on-disk format
// (AES-XTS with a per-block tweak derived from the file's logical block
// number, keyed by a kernel-only per-file nonce). That format buys nothing
// here — nothing else (kernel, `fscrypt` CLI, other tooling) will ever read
// this ciphertext through the real fscrypt path — and XTS has no integrity
// protection by design. Instead: a subkey is derived from the same master
// key (domain-separated via HKDF, never the raw key bytes reused directly),
// and content is sealed with a standard AEAD (AES-256-GCM), which gives
// authentication for free and needs no filesystem-block-layout knowledge.
//
// The recurring per-boot unlock/lock cycle (unlockFileInPlace/
// lockFileInPlace) transforms content on the SAME inode, never via
// rename — see their doc comments for why that specific property is the
// one thing standing between this and a TOCTOU regression on the guard.
// The same rule applies to the crash-recovery sidecar the cycle stages
// before every transform (see fileVaultRecoverSuffix): it too is only ever
// rewritten in place on its existing inode, never created, renamed or
// removed once a guard could be watching, and it holds a sealed record, not
// raw plaintext.
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

// readMasterKey is the seam the file-vault functions read the master key
// through — a package var (not a direct readKey() call) purely so tests can
// supply a throwaway key instead of the real /etc/app-listener/fscrypt.key.
var readMasterKey = readKey

// fileVaultRecoverSuffix names the crash-recovery sidecar written before
// every in-place transform.
//
// It looks like it sits outside the guard's scope, but it does not: the
// guard for a single-file watch root also protects that file's PARENT
// DIRECTORY against create/rename/delete of anything beside it (the same
// "rename-over-watchroot" defense CLAUDE.md calls out — see
// guard_path_rename's destination-parent-directory check in guard.bpf.c),
// so a create-temp-then-rename dance for this sidecar is denied the moment
// the resource's guard is live. This package therefore only ever rewrites
// the sidecar's EXISTING inode in place (see EnsureRecoverySidecarPlaceholder,
// stageRecovery, clearRecovery) — never creates or removes it while a guard
// could be watching. An empty sidecar means "nothing to recover", exactly
// like editprotected.HashFile's empty-means-unset convention.
//
// The content staged there is also sealed under the master key (see
// stageRecovery), never raw plaintext: a sidecar left behind by a crash
// survives even a reboot (which drops every BPF-LSM pin), so it must be
// useless on its own, not a second unguarded plaintext copy of protected
// content.
const fileVaultRecoverSuffix = ".app_listener.recover"

func fileVaultRecoverPath(path string) string { return path + fileVaultRecoverSuffix }

// rejectSidecarSymlink Lstats sidecar and errors if it currently exists as a
// symbolic link. The sidecar is only ever a plain regular file this package
// created itself (see fileVaultRecoverSuffix): every entry point that opens,
// reads, or writes it (EnsureRecoverySidecarPlaceholder, writeSidecarInPlace,
// recoverFileInPlace) calls this FIRST, so a symlink swapped in by whoever
// controls the resource's parent directory — while no guard is attached yet,
// e.g. between an install and the daemon's first start — is refused outright
// instead of silently followed into an arbitrary root-writable file
// elsewhere on the system. Missing entirely is not an error here (callers
// each handle "does not exist" in their own next step).
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

// EnsureRecoverySidecarPlaceholder creates path's (empty) recovery sidecar
// if it does not already exist yet. Callers must run this BEFORE path's
// guard ever attaches (see fileVaultRecoverSuffix) — mirrors
// cmd/functions/daemon/selfguards.go's ensureHashFilePlaceholder, same
// reason and same timing requirement. A no-op for anything that is not
// currently a regular file (directories never take the file-vault path).
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

// classifyRegularFileTarget is the safety gate every file-vault entry point
// runs first: only a plain, non-symlink, single-hard-link regular file may
// be touched (same refusals as the directory-migration path, for the same
// reasons — see classifyMigrationTarget).
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

// isRegularFileTarget is the cheap, lenient ROUTING check used by
// IsEncrypted/IsProvisioned/VerifyKey/Unlock/Lock to pick the file-vault
// path over the kernel-fscrypt directory path. It is deliberately not the
// safety gate (classifyRegularFileTarget, called by every mutating file-
// vault function, is) — on any doubt (symlink, stat error, special file) it
// answers false, falling through to today's existing directory-oriented
// code, unchanged behavior for every shape this design does not target.
func isRegularFileTarget(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

// sealFileVaultWithMasterKey derives the subkey from the current master key
// and seals plaintext under it, wiping the subkey (and NOT the caller's
// plaintext slice — callers wipe that themselves once done with it)
// afterward. Convenience wrapper shared by every mutating entry point below,
// so the read-derive-wipe dance lives in exactly one place.
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

// openFileVaultWithMasterKey derives the subkey from the current master key
// and opens record under it, wiping the subkey afterward.
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

// writeTempWithMetadata creates tmp (O_CREATE|O_EXCL, refusing to clobber a
// leftover from an interrupted run), writes content, fsyncs it, stamps src's
// mode/owner/xattrs onto it (src must still exist at its original path for
// the xattr copy — call this before renaming src away), and closes it. Any
// failure removes tmp itself, so callers need no cleanup of their own.
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

// isFileVaultCiphertext reports whether path currently holds a valid
// file-vault header (magic + supported version) — a cheap, key-free probe:
// only the header is read, not the whole file. Used by IsEncrypted for a
// regular-file target, mirroring how the kernel-fscrypt path derives state
// from on-disk metadata rather than a side registry.
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

// transformFileInPlace rewrites path's content with newContent on path's
// EXISTING inode — open the already-there fd, truncate, write, fsync; never
// unlink, never rename the live path. This is the property the whole
// single-file design rests on: the guard's guard_inodes entry for path
// (populated once, before any unlock/lock ever runs) is keyed by (dev,ino)
// and is never invalidated by this call, so there is no window where new
// content at this path is reachable under an inode the guard does not yet
// recognize — unlike a rename-based swap, which would present a brand-new
// inode the guard has no entry for until the next populate/sweep (and an
// unrecognized inode is default-ALLOW to the LSM hooks, not default-deny:
// they only enforce on inodes they know about).
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

// writeSidecarInPlace rewrites sidecar with content (which may be empty),
// on the sidecar's EXISTING inode when there is one — open, truncate,
// write, fsync — never a rename. The O_CREATE branch is only a bootstrap
// fallback for a sidecar EnsureRecoverySidecarPlaceholder never got to run
// for yet (the initial `install`-time encryption pass runs before any guard
// exists at all, so nothing denies it there); by the time a guard is
// watching path, EnsureRecoverySidecarPlaceholder has always already made
// the sidecar exist, so this always takes the open-existing branch — this
// mirrors editprotected.WriteHashFile's identical two-branch shape and
// justification.
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

// stageRecovery seals path's CURRENT bytes under the master key — the same
// AEAD format as the vault's own on-disk record, never raw plaintext — and
// writes the result to path's recovery sidecar in place (see
// fileVaultRecoverSuffix / writeSidecarInPlace). Sealing means a sidecar
// left behind by a crash, even across a reboot that drops every BPF-LSM
// pin, is just another vault-format ciphertext blob: useless without the
// master key, never a second unguarded plaintext copy of protected content.
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

// clearRecovery empties path's recovery sidecar in place — truncate to zero
// length on the same inode, never unlink (see fileVaultRecoverSuffix) — once
// an in-place transform has completed successfully. A zero-length sidecar is
// the "nothing to recover" state recoverFileInPlace reads.
func clearRecovery(path string) error {
	if err := writeSidecarInPlace(fileVaultRecoverPath(path), nil); err != nil {
		return fmt.Errorf("clearing recovery sidecar for %s: %w", path, err)
	}
	return nil
}

// recoverFileInPlace restores path from a leftover recovery sidecar, if one
// is staged (non-empty) — meaning a previous unlock/lock was interrupted
// mid-transform. Both unlockFileInPlace and lockFileInPlace call this first,
// unconditionally, before doing anything else. No-op when the sidecar is
// empty or (tolerated, though EnsureRecoverySidecarPlaceholder means it
// should not normally happen once a guard exists) missing. The sidecar
// holds a sealed record (see stageRecovery), so recovering needs the master
// key exactly like a normal unlock/lock. The restore itself is an in-place,
// same-inode write (transformFileInPlace): recovery never touches the live
// path's identity either.
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

// unlockFileInPlace decrypts the file-vault ciphertext at path back to
// plaintext, in place (see transformFileInPlace). No-op when the file is
// already plaintext (mirrors the directory Unlock's "already provisioned"
// no-op) — a failed decrypt (wrong key, tampered content) leaves the file
// untouched, never falls back to treating it as plaintext.
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

// lockFileInPlace encrypts the plaintext at path back to file-vault
// ciphertext, in place, with a fresh random nonce. No-op when already
// ciphertext (never re-seals — that would rotate the nonce for nothing and
// discard the "already locked" idempotence the directory Lock also has).
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

// verifyFileVaultKey confirms the master key actually decrypts path's
// current file-vault ciphertext, without mutating anything. A no-op (nil)
// when path is not currently vault-ciphertext.
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
	// fileVaultMagic identifies a file as this package's own ciphertext
	// format, distinct from anything the kernel fscrypt ioctl would ever
	// write (which never touches the file's own content, only its policy
	// metadata) — so a plain, unrelated file can never accidentally look
	// like one of ours, and IsEncrypted for a file is a pure content check.
	fileVaultMagic = "ALFV" // "App-Listener File Vault"
	// fileVaultVersion is the on-disk record format version.
	fileVaultVersion = 1
	// fileVaultNonceSize is the standard AES-GCM nonce size.
	fileVaultNonceSize = 12
	// fileVaultSubkeySize is the AES-256 key size.
	fileVaultSubkeySize = 32
	// fileVaultHeaderSize is magic + version + nonce, before the
	// ciphertext+tag.
	fileVaultHeaderSize = len(fileVaultMagic) + 1 + fileVaultNonceSize
	// fileVaultKeyInfo is the HKDF "info" parameter: fixed and versioned so
	// this subkey can never collide with any other derivation from the same
	// master key, present or future.
	fileVaultKeyInfo = "app-listener/file-vault/v1"
)

// deriveFileVaultSubkey derives the AES-256-GCM key used for single-file
// encryption from the raw fscrypt master key via HKDF-SHA256. The raw
// master key bytes are never used directly as a cipher key for anything
// other than the fscrypt directory policies they were generated for —
// reusing key material across unrelated constructions is the kind of
// cross-protocol risk domain separation exists to rule out.
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
		// Guards a stdlib behavior change silently breaking the on-disk
		// format rather than corrupting data.
		return nil, fmt.Errorf("file-vault AEAD: unexpected nonce size %d", gcm.NonceSize())
	}
	return gcm, nil
}

// sealFileVault encrypts plaintext under subkey with a fresh, random nonce
// (crypto/rand — never reused, never derived) and returns the full on-disk
// record: magic + version + nonce + ciphertext+tag.
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

// openFileVault authenticates and decrypts a record produced by
// sealFileVault. Any tampering, wrong key, truncation, or malformed header
// is rejected outright — there is no partial-trust fallback.
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

// looksLikeFileVaultRecord reports whether content has this format's magic
// and a supported version — a cheap, key-free structural check used by
// IsEncrypted/IsProvisioned for files, mirroring how the kernel-fscrypt path
// derives state from on-disk metadata rather than a side registry.
func looksLikeFileVaultRecord(content []byte) bool {
	return len(content) >= fileVaultHeaderSize &&
		bytes.Equal(content[:len(fileVaultMagic)], []byte(fileVaultMagic)) &&
		content[len(fileVaultMagic)] == fileVaultVersion
}
