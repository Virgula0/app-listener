package fscrypt

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/Virgula0/app-listener/internal/install"
	"github.com/google/fscrypt/actions"
	"github.com/google/fscrypt/metadata"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"google.golang.org/protobuf/proto"
)

// BackupSuffix is appended to a directory to form its migration backup,
// e.g. /home/alice/.ssh -> /home/alice/.ssh.app_listener.backup.
const BackupSuffix = ".app_listener.backup"

// DecryptSuffix is appended to a directory to form the temporary plaintext
// copy created while a directory is being permanently decrypted, e.g.
// /home/alice/.ssh -> /home/alice/.ssh.app_listener.decrypt. The directory
// is renamed back to its original location only after the copy completed.
const DecryptSuffix = ".app_listener.decrypt"

// encryptTmpSuffix is the empty sibling file an encryption policy is applied
// to while a single regular file is migrated: fscrypt policies can only be
// set on EMPTY files, so the content is streamed in afterwards.
const encryptTmpSuffix = ".app_listener.encrypt"

// migrationTarget classifies what kind of in-place migration applies to a
// path (see classifyMigrationTarget).
type migrationTarget int

const (
	targetDirectory migrationTarget = iota
	targetRegularFile
)

// classifyMigrationTarget inspects path and reports whether an in-place
// encryption/decryption may operate on it. Symbolic links are refused
// outright: following them could silently migrate an unrelated target, and
// a link aliasing the watched content breaks the rename-based swap.
// Regular files with more than one link (hardlinks) are refused as well:
// migrating one alias must never implicitly affect another path to the same
// inode. Everything that is neither a directory nor a unique regular file
// (sockets, FIFOs, devices) is refused too.
func classifyMigrationTarget(path string) (migrationTarget, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return targetDirectory, fmt.Errorf("stat %s: %w", path, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return targetDirectory, fmt.Errorf("refusing to migrate %s: it is a symbolic link", path)
	}
	if info.IsDir() {
		return targetDirectory, nil
	}
	if !info.Mode().IsRegular() {
		return targetDirectory, fmt.Errorf("refusing to migrate %s: only directories and regular files can be migrated", path)
	}
	if stat, ok := info.Sys().(*syscall.Stat_t); ok && stat.Nlink > 1 {
		return targetDirectory, fmt.Errorf("refusing to migrate %s: the file has %d hard links; remove the extra links first", path, stat.Nlink)
	}
	return targetRegularFile, nil
}

// Deprovision retry budget, mirroring the daemon teardown loop: a forced
// deprovision can fail with EBUSY while an inode is still pinned, and the
// key must be gone before the daemon (or the gap before it starts) is
// considered safe.
const (
	maxDeprovisionRetries = 100
	deprovisionRetryDelay = 10 * time.Millisecond
)

// chattr flag values, absent from x/sys/unix.
const (
	fsImmutableFlag = 0x10
	fsAppendFlag    = 0x20
)

// VerifyKey checks that the master key in MasterKeyFile actually unlocks
// the fscrypt policy of path, without provisioning anything. It returns an
// error when the policy is locked and the key does not match — the loop in
// google/fscrypt's unwrapProtectorKey is bounded by newBoundedKeyFn, so a
// wrong key fails fast instead of spinning.
func (v *Vault) VerifyKey(path string) error {
	fsctx, err := actions.NewContextFromPath(path, nil)
	if err != nil {
		return fmt.Errorf("fscrypt context for %s: %w", path, err)
	}
	policy, err := actions.GetPolicyFromPath(fsctx, path)
	if err != nil {
		if !isLockedRegularFileErr(err) {
			return fmt.Errorf("get policy for %s: %w", path, err)
		}
		return v.verifyLockedFileKey(fsctx, path)
	}
	if policy.IsProvisionedByTargetUser() {
		return nil
	}

	keyBytes, err := readKey()
	if err != nil {
		return fmt.Errorf("verify key for %s: %w", path, err)
	}
	defer wipe(keyBytes)

	optionFn := func(_ string, _ []*actions.ProtectorOption) (int, error) {
		return 0, nil
	}
	if err := policy.Unlock(optionFn, newBoundedKeyFn(keyBytes)); err != nil {
		_ = policy.Lock()
		return fmt.Errorf("verify key for %s: %w", path, err)
	}
	return policy.Lock()
}

// EnsureSystemSetup creates the global /etc/fscrypt.conf the way `fscrypt
// setup` does (policy v2 on kernels >= 5.4, otherwise v1). It is a no-op
// when the config file already exists.
func EnsureSystemSetup() error {
	if _, err := os.Stat(actions.ConfigFileLocation); err == nil {
		return nil
	}
	policyVersion := int64(0)
	if isKernelAtLeast54() {
		policyVersion = 2
	}
	return actions.CreateConfigFile(5*time.Second, policyVersion)
}

// Encrypt migrates path to an fscrypt-encrypted directory or single file
// in place, using the master key as a raw-key protector (no shelling out —
// the google/fscrypt actions API only). The original contents are preserved
// by renaming path to path+BackupSuffix BEFORE anything is touched; the
// backup is left in place for the installer to remove after the whole
// installation succeeded.
//
// If path+BackupSuffix already exists, Encrypt refuses to run: a previous
// migration may have failed midway, and overwriting an older backup would
// destroy data. The caller must resolve that situation manually.
//
// On success the target is encrypted and locked: its key is removed from
// the kernel keyring, so the daemon must unlock it with the master key at
// startup.
func (v *Vault) Encrypt(path string) error {
	return v.EncryptWithProgress(path, nil)
}

// EncryptWithProgress behaves like Encrypt and additionally reports copy
// progress (bytes copied so far and total) through onBytes while the
// contents are moved into the encrypted directory. Directories and single
// regular files are both supported; symbolic links, hard-linked files and
// anything else is refused (see classifyMigrationTarget).
func (v *Vault) EncryptWithProgress(path string, onBytes func(copied, total int64)) error {
	backup := path + BackupSuffix
	if _, err := os.Lstat(backup); err == nil {
		return fmt.Errorf("backup %s already exists: a previous migration was interrupted; restore or remove it manually before retrying", backup)
	}

	target, err := classifyMigrationTarget(path)
	if err != nil {
		return err
	}
	if target == targetRegularFile {
		return v.encryptFileWithProgress(path, backup, onBytes)
	}
	return v.encryptDirWithProgress(path, backup, onBytes)
}

// encryptDirWithProgress migrates a directory in place (the original
// ssh-guard flow).
func (v *Vault) encryptDirWithProgress(path, backup string, onBytes func(copied, total int64)) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	// The daemon's LSM engine does not set chattr flags, but a system
	// previously hardened by the original ssh-guard may carry immutable or
	// append-only flags that would block the rename below.
	if err = stripImmutableFlags(path); err != nil {
		return fmt.Errorf("strip immutable flags on %s: %w", path, err)
	}

	if err = os.Rename(path, backup); err != nil {
		return fmt.Errorf("move %s to %s: %w", path, backup, err)
	}
	if err = os.Mkdir(path, info.Mode().Perm()); err != nil {
		return fmt.Errorf("recreate %s: %w", path, err)
	}

	var policy *actions.Policy
	restore := true
	defer func() {
		if restore {
			// Roll the directory back so a failed migration leaves the
			// original data in place. RemoveAll runs while the policy key
			// is still provisioned: deleting entries of an encrypted
			// directory requires its key, and dropping it first would
			// leave a partially copied directory unremovable.
			_ = os.RemoveAll(path)
			if policy != nil {
				_ = policy.Deprovision(false)
			}
			_ = os.Rename(backup, path)
		}
	}()

	policy, err = applyRawKeyPolicy(path)
	if err != nil {
		return fmt.Errorf("apply fscrypt policy to %s: %w", path, err)
	}
	// applyRawKeyPolicy leaves the policy key provisioned on purpose: the
	// contents must be copied back while the key is available, or the
	// kernel refuses every write with "required key not available".
	if err := install.CopyTreeWithProgress(backup, path, onBytes); err != nil {
		return fmt.Errorf("copy contents back into %s: %w", path, err)
	}
	restore = false

	// The migration is complete: wipe the in-memory keys and remove the
	// policy key from the kernel keyring so the directory is locked until
	// the daemon unlocks it. The deprovision is forced with a bounded
	// EBUSY retry (same budget as the daemon's teardown): a key left in
	// the keyring would leave the directory unlocked in the gap between
	// the installer stopping the daemon and it starting up again.
	if err := lockAndDeprovision(policy, nil); err != nil {
		log.Warnf("encrypted %s but could not remove its key from the kernel keyring (the directory stays unlocked until the daemon starts): %v", path, err)
	}
	return nil
}

// progressWriter counts copied bytes and feeds the migration progress bar.
type progressWriter struct {
	w     io.Writer
	total int64
	n     int64
	cb    func(copied, total int64)
}

func (p *progressWriter) Write(b []byte) (int, error) {
	n, err := p.w.Write(b)
	p.n += int64(n)
	if p.cb != nil {
		p.cb(p.n, p.total)
	}
	return n, err
}

// stampCopyMetadata applies the original file's mode, owner and extended
// attributes to the freshly written copy. Ownership requires root (the
// installer runs as root); extended attributes are best-effort: one that
// cannot be carried over is logged and skipped, never fatal.
func stampCopyMetadata(out *os.File, src, dst string, info os.FileInfo) error {
	if err := out.Chmod(info.Mode().Perm()); err != nil {
		return fmt.Errorf("chmod %s: %w", dst, err)
	}
	if os.Geteuid() == 0 {
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			if chownErr := out.Chown(int(stat.Uid), int(stat.Gid)); chownErr != nil {
				return fmt.Errorf("chown %s: %w", dst, chownErr)
			}
		}
	}
	if xerr := copyXattrs(src, dst); xerr != nil {
		log.Warnf("%s: extended attributes were not fully preserved on %s: %v", src, dst, xerr)
	}
	return nil
}

// copyIntoEncryptedCopy streams srcPath's content through out — whose
// filesystem-level fscrypt encryption is already active — then syncs it to
// disk, stamps the original metadata onto it and closes it.
func copyIntoEncryptedCopy(out *os.File, srcPath, dst string, info os.FileInfo, onBytes func(copied, total int64)) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("open %s: %w", srcPath, err)
	}
	pw := &progressWriter{w: out, total: info.Size(), cb: onBytes}
	if _, copyErr := io.Copy(pw, src); copyErr != nil {
		_ = src.Close()
		return fmt.Errorf("copy %s to %s: %w", srcPath, dst, copyErr)
	}
	if err = src.Close(); err != nil {
		return fmt.Errorf("close %s: %w", srcPath, err)
	}
	if syncErr := out.Sync(); syncErr != nil {
		return fmt.Errorf("sync %s: %w", dst, syncErr)
	}
	if stampErr := stampCopyMetadata(out, srcPath, dst, info); stampErr != nil {
		return stampErr
	}
	if closeErr := out.Close(); closeErr != nil {
		return fmt.Errorf("close %s: %w", dst, closeErr)
	}
	return nil
}

// encryptFileWithProgress migrates a single regular file in place. The
// kernel only accepts an encryption policy on an EMPTY regular file, so a
// fresh empty sibling carries the policy and receives the content while its
// key is provisioned; after the key is removed again (ciphertext at rest)
// the two are swapped with the same crash-safety contract as the directory
// migration: original renamed to backup first, temp installed second, and a
// full rollback of both renames until then.
func (v *Vault) encryptFileWithProgress(path, backup string, onBytes func(copied, total int64)) error {
	tmp := path + encryptTmpSuffix
	if _, err := os.Lstat(tmp); err == nil {
		return fmt.Errorf("temporary file %s already exists: a previous migration was interrupted; remove it manually before retrying", tmp)
	}

	srcInfo, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, srcInfo.Mode().Perm())
	if err != nil {
		return fmt.Errorf("create %s: %w", tmp, err)
	}
	cleanupTmp := true
	defer func() {
		if cleanupTmp {
			_ = out.Close()
			_ = os.Remove(tmp)
		}
	}()

	// The policy must be applied to an empty file: create it first, apply
	// the policy, only then stream the content through the encrypted fd.
	policy, applyErr := applyRawKeyPolicy(tmp)
	if applyErr != nil {
		return fmt.Errorf("apply fscrypt policy to %s: %w", tmp, applyErr)
	}
	// applyRawKeyPolicy leaves the policy key provisioned on purpose: the
	// content can only be written while the key is available.
	restore := true
	defer func() {
		if restore {
			_ = lockAndDeprovision(policy, nil)
			_ = os.Remove(tmp)
			// Restores the original when it had already been moved to
			// the backup name; a harmless no-op error otherwise.
			_ = os.Rename(backup, path)
		}
	}()

	if copyErr := copyIntoEncryptedCopy(out, path, tmp, srcInfo, onBytes); copyErr != nil {
		return copyErr
	}
	cleanupTmp = false

	// Content is now encrypted at rest once the key is removed. Swap the
	// two names: original to backup first, encrypted copy into place
	// second; until both renames succeeded the rollback above undoes
	// everything (the backup rename-back is a no-op before the first one).
	if err := lockAndDeprovision(policy, nil); err != nil {
		log.Warnf("encrypted %s but could not remove its key from the kernel keyring (the file stays unlocked until the daemon starts): %v", path, err)
	}
	if err := os.Rename(path, backup); err != nil {
		return fmt.Errorf("move %s to %s: %w", path, backup, err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("move %s to %s: %w", tmp, path, err)
	}
	restore = false
	return nil
}

// Decrypt permanently removes the fscrypt policy of path, restoring a
// plaintext directory in place. The contents are copied into a temporary
// sibling directory (path+DecryptSuffix) while the policy key is
// provisioned, and the original encrypted directory is only removed after
// the copy completed; the plaintext copy is then renamed into place. On
// any failure the temporary directory is removed and the encrypted
// directory is left untouched, so a permanent decryption can never lose
// data. The policy and protector metadata left behind in /.fscrypt become
// orphaned and are cleaned up by the caller (see CleanOrphans).
func (v *Vault) Decrypt(path string) error {
	return v.DecryptWithProgress(path, nil)
}

// DecryptWithProgress behaves like Decrypt and additionally reports copy
// progress (bytes copied so far and total) through onBytes while the
// contents are moved into the plaintext copy. Directories and single
// regular files are both supported (see classifyMigrationTarget).
func (v *Vault) DecryptWithProgress(path string, onBytes func(copied, total int64)) error {
	target, err := classifyMigrationTarget(path)
	if err != nil {
		return err
	}
	if target == targetRegularFile {
		return v.decryptFileWithProgress(path, onBytes)
	}
	return v.decryptDirWithProgress(path, onBytes)
}

// decryptDirWithProgress permanently decrypts a directory in place (the
// original flow).
func (v *Vault) decryptDirWithProgress(path string, onBytes func(copied, total int64)) error {
	tmp := path + DecryptSuffix
	if _, err := os.Lstat(tmp); err == nil {
		return fmt.Errorf("temporary directory %s already exists: a previous decryption was interrupted; remove it manually before retrying", tmp)
	}
	if err := requireEncryptedPath(path); err != nil {
		return err
	}

	// Fetch the policy BEFORE the directory is removed: the policy object
	// is needed to drop the key from the kernel keyring after the copy.
	fsctx, err := actions.NewContextFromPath(path, nil)
	if err != nil {
		return fmt.Errorf("fscrypt context for %s: %w", path, err)
	}
	policy, err := actions.GetPolicyFromPath(fsctx, path)
	if err != nil {
		return fmt.Errorf("get policy for %s: %w", path, err)
	}

	// The contents must be readable while they are copied out: provision
	// the policy key first and always drop it again afterwards, whether
	// the copy succeeded or not.
	if err := v.Unlock(path); err != nil {
		return err
	}
	defer func() {
		if err := lockAndDeprovision(policy, nil); err != nil {
			log.Warnf("decrypted %s but could not remove its key from the kernel keyring: %v", path, err)
		}
	}()

	if err := install.CopyTreeWithProgress(path, tmp, onBytes); err != nil {
		_ = os.RemoveAll(tmp)
		return fmt.Errorf("copy %s to %s: %w — the encrypted directory was left untouched", path, tmp, err)
	}
	if err := os.RemoveAll(path); err != nil {
		_ = os.RemoveAll(tmp)
		return fmt.Errorf("remove encrypted directory %s: %w — the plaintext copy at %s was removed as well", path, err, tmp)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("move %s to %s: %w — the plaintext copy still holds the data", tmp, path, err)
	}
	return nil
}

// decryptFileWithProgress permanently decrypts a single encrypted regular
// file in place. The plaintext copy is written to a temporary sibling while
// the policy key is provisioned; the encrypted original is only removed
// after the copy completed and fsynced, then the plaintext copy is renamed
// into place. On any failure the temporary file is removed and the
// encrypted original is left untouched.
func (v *Vault) decryptFileWithProgress(path string, onBytes func(copied, total int64)) error {
	tmp := path + DecryptSuffix
	if _, err := os.Lstat(tmp); err == nil {
		return fmt.Errorf("temporary file %s already exists: a previous decryption was interrupted; remove it manually before retrying", tmp)
	}
	if err := requireEncryptedPath(path); err != nil {
		return err
	}

	// Fetch the policy BEFORE the original is removed: the policy object
	// is needed to drop the key from the kernel keyring after the copy.
	fsctx, err := actions.NewContextFromPath(path, nil)
	if err != nil {
		return fmt.Errorf("fscrypt context for %s: %w", path, err)
	}
	policy, err := actions.GetPolicyFromPath(fsctx, path)
	if err != nil {
		return fmt.Errorf("get policy for %s: %w", path, err)
	}

	// The contents must be readable while they are copied out: provision
	// the policy key first (a locked regular file needs the descriptor
	// fallback inside Unlock) and always drop it again afterwards.
	if unlockErr := v.Unlock(path); unlockErr != nil {
		return unlockErr
	}
	defer func() {
		if depErr := lockAndDeprovision(policy, nil); depErr != nil {
			log.Warnf("decrypted %s but could not remove its key from the kernel keyring: %v", path, depErr)
		}
	}()

	info, statErr := os.Stat(path)
	if statErr != nil {
		return fmt.Errorf("stat %s: %w", path, statErr)
	}
	out, openErr := os.OpenFile(tmp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, info.Mode().Perm())
	if openErr != nil {
		return fmt.Errorf("create %s: %w", tmp, openErr)
	}
	cleanupTmp := true
	defer func() {
		if cleanupTmp {
			_ = out.Close()
			_ = os.Remove(tmp)
		}
	}()

	if copyErr := copyIntoEncryptedCopy(out, path, tmp, info, onBytes); copyErr != nil {
		return fmt.Errorf("%w — the encrypted file was left untouched", copyErr)
	}
	cleanupTmp = false

	if err := os.Remove(path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("remove encrypted file %s: %w — the plaintext copy at %s was removed as well", path, err, tmp)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("move %s to %s: %w — the plaintext copy still holds the data", tmp, path, err)
	}
	return nil
}

// requireEncryptedPath reports an error unless path is a directory or a
// unique regular file carrying an fscrypt policy. Anything else — plain
// files, symlinks, hardlinks, special files — can never be decrypted and is
// refused, mirroring the safety contract of RestoreBackup.
func requireEncryptedPath(path string) error {
	if _, err := classifyMigrationTarget(path); err != nil {
		return err
	}
	encrypted, encErr := hasEncryptionPolicy(path)
	if encErr != nil {
		return fmt.Errorf("checking encryption of %s: %w", path, encErr)
	}
	if !encrypted {
		return fmt.Errorf("%s is not encrypted with fscrypt: refusing to decrypt it", path)
	}
	return nil
}

// RestoreBackup undoes a migration: the encrypted directory or file at
// path is deleted and path+BackupSuffix is moved back to path, restoring
// the original unencrypted content. The encrypted target is unlocked with
// the master key first — deleting a locked encrypted directory (or opening
// a locked encrypted file for removal) fails with ENOKEY because the
// kernel needs the key to touch it. It refuses to delete a path that
// exists but is NOT encrypted (the user may have created new data there
// since the migration), and it is an error when no backup exists.
func (v *Vault) RestoreBackup(path string) error {
	backup := path + BackupSuffix
	if _, err := os.Lstat(backup); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("backup %s does not exist", backup)
		}
		return fmt.Errorf("stat backup %s: %w", backup, err)
	}
	if _, err := os.Lstat(path); err == nil {
		encrypted, encErr := hasEncryptionPolicy(path)
		if encErr != nil {
			return fmt.Errorf("checking encryption of %s: %w", path, encErr)
		}
		if !encrypted {
			return fmt.Errorf("%s exists but is not encrypted: refusing to delete it — remove or restore it manually", path)
		}
		if err := unlockForRemoval(path); err != nil {
			return fmt.Errorf("unlock %s for removal: %w", path, err)
		}
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("remove encrypted directory %s: %w", path, err)
		}
	}
	if err := os.Rename(backup, path); err != nil {
		return fmt.Errorf("restore %s to %s: %w", backup, path, err)
	}
	return nil
}

// unlockForRemoval provisions the policy key of an encrypted path so its
// contents can be deleted, then removes the key from the keyring again
// (leaving it locked) before returning. A locked encrypted regular file is
// handled through the descriptor fallback: provision first, then re-read
// the policy now that the file opens.
func unlockForRemoval(path string) error {
	fsctx, err := actions.NewContextFromPath(path, nil)
	if err != nil {
		return err
	}
	policy, err := actions.GetPolicyFromPath(fsctx, path)
	if err != nil {
		if !isLockedRegularFileErr(err) {
			return err
		}
		if provErr := (&Vault{}).provisionLockedFile(fsctx, path); provErr != nil {
			return provErr
		}
		policy, err = actions.GetPolicyFromPath(fsctx, path)
		if err != nil {
			return err
		}
	}
	if policy.IsProvisionedByTargetUser() {
		return nil
	}
	keyBytes, err := readKey()
	if err != nil {
		return err
	}
	defer wipe(keyBytes)
	optionFn := func(_ string, _ []*actions.ProtectorOption) (int, error) {
		return 0, nil
	}
	if err := policy.Unlock(optionFn, newBoundedKeyFn(keyBytes)); err != nil {
		return err
	}
	defer func() { _ = lockAndDeprovision(policy, nil) }()
	if err := policy.Provision(); err != nil {
		return err
	}
	return nil
}

// applyRawKeyPolicy protects the (empty) directory at path with a fresh
// policy wrapped by a raw-key protector holding the master key. On success
// the policy key is left provisioned in the kernel keyring: the caller
// must copy the directory contents back while the key is available and
// deprovision afterwards (see Encrypt). On any failure the in-memory keys
// are wiped and the keyring is cleaned up.
func applyRawKeyPolicy(path string) (*actions.Policy, error) {
	fsctx, ctxErr := actions.NewContextFromPath(path, nil)
	if ctxErr != nil {
		return nil, ctxErr
	}
	rawCtx := modifiedContextWithSource(fsctx, metadata.SourceType_raw_key)

	keyBytes, readErr := readKey()
	if readErr != nil {
		return nil, readErr
	}
	defer wipe(keyBytes)

	name := fmt.Sprintf("app-listener-key-%d", time.Now().UnixNano())
	protector, createErr := actions.CreateProtector(rawCtx, name, newBoundedKeyFn(keyBytes), nil)
	if createErr != nil {
		return nil, createErr
	}
	created := true
	defer func() {
		if created {
			_ = protector.Revert()
		}
	}()

	if err := protector.Unlock(newBoundedKeyFn(keyBytes)); err != nil {
		return nil, err
	}
	policy, err := actions.CreatePolicy(rawCtx, protector)
	if err != nil {
		return nil, err
	}
	defer func() {
		if created {
			_ = policy.Revert()
		}
	}()
	// On failure after the policy exists, wipe the in-memory keys and drop
	// the policy key from the keyring (a no-op if Provision never ran).
	// On success both are skipped: the key must stay provisioned for the
	// copy-back in Encrypt.
	ok := false
	defer func() {
		if !ok {
			_ = lockAndDeprovision(policy, protector)
		}
	}()

	optionFn := func(_ string, _ []*actions.ProtectorOption) (int, error) {
		return 0, nil
	}
	if err := policy.Unlock(optionFn, newBoundedKeyFn(keyBytes)); err != nil {
		return nil, err
	}
	if err := policy.Provision(); err != nil {
		return nil, err
	}
	if err := policy.Apply(path); err != nil {
		return nil, err
	}
	created = false
	ok = true
	return policy, nil
}

// lockAndDeprovision wipes the in-memory keys of policy and protector and
// removes the policy key from the kernel keyring, leaving the encrypted
// directory locked. The deprovision is forced (all users) with a bounded
// EBUSY retry loop mirroring the daemon's teardown; a key that is already
// gone (ENOKEY) counts as success. The returned error is non-nil only when
// the key could not be removed after the full retry budget.
func lockAndDeprovision(policy *actions.Policy, protector *actions.Protector) error {
	_ = policy.Lock()
	if protector != nil {
		_ = protector.Lock()
	}
	for i := 0; i < maxDeprovisionRetries; i++ {
		err := policy.Deprovision(true)
		switch {
		case err == nil:
			return nil
		case isDeprovisionBusy(err):
			time.Sleep(deprovisionRetryDelay)
		case isDeprovisionMissing(err):
			return nil
		default:
			return err
		}
	}
	return fmt.Errorf("deprovision after %d retries", maxDeprovisionRetries)
}

// isDeprovisionBusy reports whether a deprovision error means inodes are
// still pinned by the key (retry), mirroring the daemon's classification.
func isDeprovisionBusy(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "some files using the key are still open") ||
		strings.Contains(errStr, "in use")
}

// isDeprovisionMissing reports whether a deprovision error means the key
// is already gone from the kernel keyring (counts as locked).
func isDeprovisionMissing(err error) bool {
	errStr := err.Error()
	return strings.Contains(errStr, "Required key not available") ||
		strings.Contains(errStr, "key not present")
}

// modifiedContextWithSource returns a copy of ctx whose protector source is
// source (mirrors the unexported helper in google/fscrypt's actions).
func modifiedContextWithSource(ctx *actions.Context, source metadata.SourceType) *actions.Context {
	modified := *ctx
	modified.Config = proto.Clone(ctx.Config).(*metadata.Config)
	modified.Config.Source = source
	return &modified
}

// stripImmutableFlags clears the immutable and append-only flags on every
// file and directory under root, mirroring the original migration script's
// `chattr -R -i -a`. Symbolic links are skipped (they cannot block the
// migration and their targets are not followed), and entries that vanish
// mid-walk (ENOENT) are tolerated: only the root itself is fatal. Errors
// on other paths always include the failing path.
func stripImmutableFlags(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			if errors.Is(walkErr, fs.ErrNotExist) && path != root {
				// The entry disappeared while walking (a live
				// application churning its profile): nothing to
				// strip, nothing to do.
				return nil
			}
			return fmt.Errorf("walk %s: %w", path, walkErr)
		}
		if d.Type()&fs.ModeSymlink != 0 {
			// A dangling symlink (e.g. firefox's lock file) would
			// fail the open below: never follow links.
			return nil
		}
		fd, openErr := unix.Open(path, unix.O_RDONLY, 0)
		if openErr != nil {
			if errors.Is(openErr, fs.ErrNotExist) && path != root {
				return nil
			}
			return fmt.Errorf("open %s: %w", path, openErr)
		}
		// On error (e.g. tmpfs without chattr support) flags is zeroed,
		// which compares equal to the cleared value below.
		flags, _, _ := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), unix.FS_IOC_GETFLAGS, 0)
		unix.Close(fd)
		cleared := flags &^ (fsImmutableFlag | fsAppendFlag)
		if cleared == flags {
			return nil
		}
		fd, openErr = unix.Open(path, unix.O_RDONLY, 0)
		if openErr != nil {
			return fmt.Errorf("open %s: %w", path, openErr)
		}
		_, _, _ = unix.Syscall(unix.SYS_IOCTL, uintptr(fd), unix.FS_IOC_SETFLAGS, cleared)
		unix.Close(fd)
		return nil
	})
}

func wipe(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// isKernelAtLeast54 reports whether the running kernel is 5.4 or newer
// (used only to decide the default fscrypt policy version, as `fscrypt
// setup` does).
func isKernelAtLeast54() bool {
	var uts unix.Utsname
	if err := unix.Uname(&uts); err != nil {
		return false
	}
	release := string(uts.Release[:])
	major := 0
	if _, err := fmt.Sscanf(release, "%d", &major); err != nil {
		return false
	}
	return major >= 5
}
