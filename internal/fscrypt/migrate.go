package fscrypt

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Virgula0/app-listener/internal/install"
	"github.com/google/fscrypt/actions"
	"github.com/google/fscrypt/metadata"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"google.golang.org/protobuf/proto"
)

// BackupSuffix forms a migration backup path, e.g. /home/alice/.ssh ->
// /home/alice/.ssh.app_listener.backup.
const BackupSuffix = ".app_listener.backup"

// DecryptSuffix forms the temporary plaintext copy created during permanent
// decryption (/home/alice/.ssh.app_listener.decrypt), renamed into place only
// after the copy completed.
const DecryptSuffix = ".app_listener.decrypt"

// encryptTmpSuffix is the empty sibling file carrying the encryption policy
// during single-file migration: policies apply only to EMPTY files, so the
// content is streamed in afterwards.
const encryptTmpSuffix = ".app_listener.encrypt"

// migrationTarget classifies which in-place migration applies to a path
// (see classifyMigrationTarget).
type migrationTarget int

const (
	targetDirectory migrationTarget = iota
	targetRegularFile
)

// classifyMigrationTarget reports whether in-place migration may operate on path: symlinks are refused outright
// (following one could silently migrate an unrelated target and break the rename-based swap), hard-linked files too
// (one alias must never implicitly affect another path to the same inode), plus anything not a directory/regular file.
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

// Bounded EBUSY-retry budget for forced deprovisions, matching the daemon
// teardown loop: an inode still pinned by the key is retried, never fatal.
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
// path's fscrypt policy, without provisioning anything. A locked/mismatched
// key errors fast (newBoundedKeyFn bounds the library's unwrap loop).
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

	optionFn := acceptFirstProtectorOption
	if err := policy.Unlock(optionFn, newBoundedKeyFn(keyBytes)); err != nil {
		_ = policy.Lock()
		return fmt.Errorf("verify key for %s: %w", path, err)
	}
	return policy.Lock()
}

// EnsureSystemSetup creates the global /etc/fscrypt.conf like `fscrypt setup`
// (policy v2 on kernels >= 5.4, else v1); no-op when it already exists.
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

// Encrypt migrates path in place to fscrypt encryption via the google/fscrypt actions API (raw-key protector),
// locking it on success for the daemon to unlock at startup. Backup-first: path is renamed to path+BackupSuffix
// BEFORE anything is touched and kept until installation fully succeeds; an existing backup aborts the run.
func (v *Vault) Encrypt(path string) error {
	return v.EncryptWithProgress(path, nil)
}

// EncryptWithProgress behaves like Encrypt and additionally reports copy
// progress (copied/total bytes) through onBytes; unsupported targets are
// refused (see classifyMigrationTarget).
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

// encryptDirWithProgress migrates a directory in place.
func (v *Vault) encryptDirWithProgress(path, backup string, onBytes func(copied, total int64)) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}

	// A system previously hardened by ssh-guard may carry immutable or
	// append-only chattr flags that would block the rename below.
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
			// Roll back so a failed migration leaves the original data in
			// place; RemoveAll runs while the policy key is still
			// provisioned (deleting encrypted entries requires it).
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
	// The key stays provisioned on purpose: contents must be copied back
	// while it is available, or every write fails with "required key not
	// available".
	if err := install.CopyTreeWithProgress(backup, path, onBytes); err != nil {
		return fmt.Errorf("copy contents back into %s: %w", path, err)
	}
	restore = false

	// Migration complete: lock and remove the key from the kernel keyring
	// (forced, bounded EBUSY retry) so the directory stays locked until the
	// daemon starts; failure here only logs.
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
// attributes to the fresh copy (ownership requires root). Xattrs are
// best-effort: one that cannot be carried over is logged and skipped, never fatal.
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

// copyIntoEncryptedCopy streams srcPath's content through out (its fscrypt
// encryption is already active), then syncs, stamps the original metadata and closes it.
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

// encryptFileWithProgress migrates a single regular file in place: a fresh empty sibling carries the
// policy (policies apply only to EMPTY files) and receives the content while its key is provisioned;
// after key removal (ciphertext at rest) both are swapped under the directory migration's crash-safety contract.
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

	// The policy must be applied to an empty file: create it first, apply the
	// policy, then stream content through the encrypted fd.
	policy, applyErr := applyRawKeyPolicy(tmp)
	if applyErr != nil {
		return fmt.Errorf("apply fscrypt policy to %s: %w", tmp, applyErr)
	}
	// Key left provisioned on purpose: writing needs it.
	restore := true
	defer func() {
		if restore {
			_ = lockAndDeprovision(policy, nil)
			_ = os.Remove(tmp)
			// No-op unless the original was already moved to backup.
			_ = os.Rename(backup, path)
		}
	}()

	if copyErr := copyIntoEncryptedCopy(out, path, tmp, srcInfo, onBytes); copyErr != nil {
		return copyErr
	}
	cleanupTmp = false

	// Content is ciphertext at rest once the key is gone: swap the two names
	// (rollback above undoes everything until both renames succeeded).
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

// Decrypt permanently removes path's fscrypt policy, restoring plaintext in place: contents are copied
// into path+DecryptSuffix while the key is provisioned and the encrypted original removed only after
// the copy completes; any failure deletes the temp copy and leaves the encrypted target untouched.
func (v *Vault) Decrypt(path string) error {
	return v.DecryptWithProgress(path, nil)
}

// DecryptWithProgress behaves like Decrypt and additionally reports copy
// progress (copied/total bytes) through onBytes (see classifyMigrationTarget
// for supported targets).
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

// decryptDirWithProgress permanently decrypts a directory in place.
func (v *Vault) decryptDirWithProgress(path string, onBytes func(copied, total int64)) error {
	tmp := path + DecryptSuffix
	if _, err := os.Lstat(tmp); err == nil {
		return fmt.Errorf("temporary directory %s already exists: a previous decryption was interrupted; remove it manually before retrying", tmp)
	}
	if err := requireEncryptedPath(path); err != nil {
		return err
	}

	// Fetch the policy BEFORE the directory is removed: it is needed to
	// drop the key after the copy.
	fsctx, err := actions.NewContextFromPath(path, nil)
	if err != nil {
		return fmt.Errorf("fscrypt context for %s: %w", path, err)
	}
	policy, err := actions.GetPolicyFromPath(fsctx, path)
	if err != nil {
		return fmt.Errorf("get policy for %s: %w", path, err)
	}

	// Provision the key so contents are readable during copy-out; always
	// dropped again afterwards.
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

// decryptFileWithProgress permanently decrypts a single encrypted regular file in place: the plaintext
// copy is written to a temporary sibling while the key is provisioned and the fsynced encrypted original
// removed only after the copy completed; any failure removes the temp file, leaving the original untouched.
func (v *Vault) decryptFileWithProgress(path string, onBytes func(copied, total int64)) error {
	tmp := path + DecryptSuffix
	if _, err := os.Lstat(tmp); err == nil {
		return fmt.Errorf("temporary file %s already exists: a previous decryption was interrupted; remove it manually before retrying", tmp)
	}
	if err := requireEncryptedPath(path); err != nil {
		return err
	}

	// Fetch the policy BEFORE the original is removed: it is needed to
	// drop the key after the copy.
	fsctx, err := actions.NewContextFromPath(path, nil)
	if err != nil {
		return fmt.Errorf("fscrypt context for %s: %w", path, err)
	}
	policy, err := actions.GetPolicyFromPath(fsctx, path)
	if err != nil {
		return fmt.Errorf("get policy for %s: %w", path, err)
	}

	// Provision the key for the copy-out; always dropped again afterwards
	// (a locked regular file needs the descriptor fallback inside Unlock).
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

// requireEncryptedPath errors unless path is a directory or unique regular file carrying an fscrypt
// policy; anything else (plain/special files, symlinks, hardlinks) can never be decrypted and is
// refused, mirroring RestoreBackup's safety contract.
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

// RestoreBackup undoes a migration: the encrypted target at path is deleted and path+BackupSuffix moved back.
// The master key unlocks it first (removing a locked encrypted target fails with ENOKEY); an existing but
// unencrypted path is never deleted (may hold new user data), and a missing backup is an error.
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

// unlockForRemoval provisions path's policy key so its contents can be deleted, then removes the key
// again (leaving it locked) before returning. A locked encrypted regular file takes the descriptor
// fallback: provision first, then re-read the policy now that the file opens.
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
	optionFn := acceptFirstProtectorOption
	if err := policy.Unlock(optionFn, newBoundedKeyFn(keyBytes)); err != nil {
		return err
	}
	defer func() { _ = lockAndDeprovision(policy, nil) }()
	if err := policy.Provision(); err != nil {
		return err
	}
	return nil
}

// applyRawKeyPolicy protects the (empty) path with a fresh policy wrapped by a raw-key protector holding
// the master key, returning with the policy key provisioned: the caller must copy contents back while the
// key is available and deprovision afterwards (see Encrypt). On failure keys are wiped and the keyring cleaned.
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
	// On failure wipe in-memory keys and drop the policy key (no-op if
	// Provision never ran); on success keep the key for Encrypt's copy-back.
	ok := false
	defer func() {
		if !ok {
			_ = lockAndDeprovision(policy, protector)
		}
	}()

	optionFn := acceptFirstProtectorOption
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

// lockAndDeprovision wipes in-memory keys and removes the policy key from the kernel keyring, leaving
// the encrypted target locked until the daemon unlocks it at startup. The forced deprovision retries
// EBUSY (inodes still pinned by the key) within the shared budget; ENOKEY counts as success.
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
// still pinned by the key (retry).
func isDeprovisionBusy(err error) bool {
	return classifyDeprovision(err) == depBusy
}

// isDeprovisionMissing reports whether a deprovision error means the key is
// already gone from the kernel keyring (counts as locked).
func isDeprovisionMissing(err error) bool {
	return classifyDeprovision(err) == depMissing
}

// modifiedContextWithSource returns a copy of ctx whose protector source is source.
func modifiedContextWithSource(ctx *actions.Context, source metadata.SourceType) *actions.Context {
	modified := *ctx
	modified.Config = proto.Clone(ctx.Config).(*metadata.Config)
	modified.Config.Source = source
	return &modified
}

// stripImmutableFlags clears the immutable and append-only flags on every file and directory under root
// (the equivalent of `chattr -R -i -a`). Symlinks are skipped (targets not followed), entries vanishing
// mid-walk (ENOENT) are tolerated — only the root itself is fatal; other errors always name the path.
func stripImmutableFlags(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			if errors.Is(walkErr, fs.ErrNotExist) && path != root {
				// Entry vanished mid-walk (a live app churning
				// its profile): nothing to strip.
				return nil
			}
			return fmt.Errorf("walk %s: %w", path, walkErr)
		}
		if d.Type()&fs.ModeSymlink != 0 {
			// Never follow links: a dangling symlink (e.g. a
			// firefox lock file) would fail the open below.
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
		// comparing equal to the cleared value below.
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

// isKernelAtLeast54 reports whether the running kernel is 5.4 or newer; used
// only to pick the default fscrypt policy version (as `fscrypt setup` does).
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
