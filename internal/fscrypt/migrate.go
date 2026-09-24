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
	"github.com/Virgula0/app-listener/internal/safeio"
	"github.com/google/fscrypt/actions"
	"github.com/google/fscrypt/metadata"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"google.golang.org/protobuf/proto"
)

// BackupSuffix forms a migration backup path, e.g. /home/alice/.ssh ->
// /home/alice/.ssh.app_listener.backup.
const BackupSuffix = ".app_listener.backup"

// DecryptSuffix forms the temporary plaintext copy made during permanent decryption
// (.ssh.app_listener.decrypt), renamed into place only after the copy completes.
const DecryptSuffix = ".app_listener.decrypt"

// encryptTmpSuffix is the empty sibling carrying the policy during single-file migration (policies
// apply only to EMPTY files; content is streamed in afterwards).
const encryptTmpSuffix = ".app_listener.encrypt"

// migrationTarget: which in-place migration applies to a path (classifyMigrationTarget).
type migrationTarget int

const (
	targetDirectory migrationTarget = iota
	targetRegularFile
)

// classifyMigrationTarget reports whether in-place migration may operate on path. Refuses symlinks
// (following one could migrate an unrelated target and break the rename swap), hard-linked files
// (one alias must never implicitly affect another path to the same inode), and anything not a
// directory/regular file.
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

// Bounded EBUSY-retry budget for forced deprovisions (as the daemon teardown loop): an inode still
// pinned by the key is retried, never fatal.
const (
	maxDeprovisionRetries = 100
	deprovisionRetryDelay = 10 * time.Millisecond
)

// chattr flag values, absent from x/sys/unix.
const (
	fsImmutableFlag = 0x10
	fsAppendFlag    = 0x20
)

// VerifyKey checks that the master key unlocks path's fscrypt policy without provisioning anything
// (directories), or decrypts the current file-vault ciphertext without mutating it (regular files).
// A locked/mismatched key errors fast (newBoundedKeyFn bounds the library's unwrap loop).
func (v *Vault) VerifyKey(path string) error {
	if isRegularFileTarget(path) {
		return v.verifyFileVaultKey(path)
	}
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

// EnsureSystemSetup creates /etc/fscrypt.conf like `fscrypt setup` (policy v2 on kernels >= 5.4,
// else v1); no-op if present.
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

// Encrypt migrates path in place to fscrypt (google/fscrypt actions API, raw-key protector),
// locking it on success for the daemon to unlock at startup. Backup-first: path is renamed to
// path+BackupSuffix BEFORE anything is touched and kept until installation fully succeeds; an
// existing backup aborts.
func (v *Vault) Encrypt(path string) error {
	return v.EncryptWithProgress(path, nil)
}

// EncryptWithProgress is Encrypt plus copy progress (copied/total bytes) via onBytes; unsupported
// targets refused (classifyMigrationTarget).
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

	// A system hardened by ssh-guard may carry immutable/append-only chattr flags that would block
	// the rename below.
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
			// Roll back so a failed migration leaves the original in place; RemoveAll runs while
			// the policy key is still provisioned (deleting encrypted entries requires it).
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
	// Keep the key provisioned: contents must be copied back while it is available, or writes fail
	// with "required key not available".
	if err := install.CopyTreeWithProgress(backup, path, onBytes); err != nil {
		return fmt.Errorf("copy contents back into %s: %w", path, err)
	}
	restore = false

	// Migration complete: lock and remove the key from the keyring (forced, bounded EBUSY retry) so
	// the directory stays locked until the daemon starts; failure only logs.
	if err := lockAndDeprovision(policy, nil); err != nil {
		log.Warnf("encrypted %s but could not remove its key from the kernel keyring (the directory stays unlocked until the daemon starts): %v", path, err)
	}
	return nil
}

// stampCopyMetadata applies the original's mode, owner (needs root) and xattrs to the copy. Xattrs
// are best-effort: failures are logged and skipped.
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

// testHookBeforeFileRead runs after a single-file migration classified path and before it reads it.
var testHookBeforeFileRead = func(string) {}

// readClassifiedFile reads path only if it is still the regular file info described: the migration
// runs as root on a file its user owns, so a name swapped for a symlink after classification must
// not seal (and hand back) content from elsewhere.
func readClassifiedFile(path string, info os.FileInfo) ([]byte, error) {
	f, err := safeio.OpenRegularNoFollow(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	got, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !os.SameFile(got, info) {
		return nil, errors.New("file changed while it was being migrated")
	}
	return io.ReadAll(f)
}

// encryptFileWithProgress migrates a single regular file in place: the plaintext is sealed under
// the file-vault subkey (filevault.go; the kernel ioctl can't target a standalone file) into a
// fresh sibling, swapped in under the directory migration's backup-first crash-safety contract.
// This one-time install/uninstall migration runs before any guard watches path, so the rename swap
// is safe here, unlike the recurring unlockFileInPlace/lockFileInPlace daemon cycle (see
// filevault.go).
func (v *Vault) encryptFileWithProgress(path, backup string, onBytes func(copied, total int64)) error {
	tmp := path + encryptTmpSuffix
	if _, err := os.Lstat(tmp); err == nil {
		return fmt.Errorf("temporary file %s already exists: a previous migration was interrupted; remove it manually before retrying", tmp)
	}

	srcInfo, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}
	testHookBeforeFileRead(path)
	plaintext, err := readClassifiedFile(path, srcInfo)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	record, err := sealFileVaultWithMasterKey(plaintext)
	wipe(plaintext)
	if err != nil {
		return fmt.Errorf("encrypt %s: %w", path, err)
	}

	if err := writeTempWithMetadata(tmp, path, record, srcInfo); err != nil {
		return err
	}

	if err := os.Rename(path, backup); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("move %s to %s: %w", path, backup, err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("move %s to %s: %w — the plaintext backup at %s still holds the data", tmp, path, err, backup)
	}
	if onBytes != nil {
		onBytes(int64(len(record)), int64(len(record)))
	}
	return nil
}

// Decrypt permanently removes path's fscrypt policy, restoring plaintext in place: contents are
// copied to path+DecryptSuffix while the key is provisioned and the encrypted original removed only
// after the copy completes; any failure deletes the temp copy and leaves the target untouched.
func (v *Vault) Decrypt(path string) error {
	return v.DecryptWithProgress(path, nil)
}

// DecryptWithProgress is Decrypt plus copy progress via onBytes (see classifyMigrationTarget for
// supported targets).
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
	if err := v.requireEncryptedPath(path); err != nil {
		return err
	}

	// Fetch the policy BEFORE the directory is removed: it is needed to drop the key after the
	// copy.
	fsctx, err := actions.NewContextFromPath(path, nil)
	if err != nil {
		return fmt.Errorf("fscrypt context for %s: %w", path, err)
	}
	policy, err := actions.GetPolicyFromPath(fsctx, path)
	if err != nil {
		return fmt.Errorf("get policy for %s: %w", path, err)
	}

	// Provision the key so contents are readable during copy-out; always dropped afterwards.
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

// decryptFileWithProgress permanently decrypts a file-vault regular file in place: plaintext goes
// to a temp sibling and the ciphertext original is removed only after the write completed and was
// fsynced; any failure removes the temp, leaving the original untouched. No keyring key to
// provision (the AEAD open yields plaintext directly).
func (v *Vault) decryptFileWithProgress(path string, onBytes func(copied, total int64)) error {
	tmp := path + DecryptSuffix
	if _, err := os.Lstat(tmp); err == nil {
		return fmt.Errorf("temporary file %s already exists: a previous decryption was interrupted; remove it manually before retrying", tmp)
	}
	if err := v.requireEncryptedPath(path); err != nil {
		return err
	}

	info, statErr := os.Lstat(path)
	if statErr != nil {
		return fmt.Errorf("stat %s: %w", path, statErr)
	}
	testHookBeforeFileRead(path)
	record, err := readClassifiedFile(path, info)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}

	plaintext, err := openFileVaultWithMasterKey(record)
	if err != nil {
		return fmt.Errorf("decrypt %s: %w — the encrypted file was left untouched", path, err)
	}
	defer wipe(plaintext)

	if err := writeTempWithMetadata(tmp, path, plaintext, info); err != nil {
		return fmt.Errorf("%w — the encrypted file was left untouched", err)
	}

	if err := os.Remove(path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("remove encrypted file %s: %w — the plaintext copy at %s was removed as well", path, err, tmp)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("move %s to %s: %w — the plaintext copy still holds the data", tmp, path, err)
	}
	if onBytes != nil {
		onBytes(int64(len(plaintext)), int64(len(plaintext)))
	}
	return nil
}

// requireEncryptedPath errors unless path is a directory with an fscrypt policy or a unique regular
// file with a file-vault header (IsEncrypted); anything else (plain/special files, symlinks,
// hardlinks) is refused, mirroring RestoreBackup's safety contract.
func (v *Vault) requireEncryptedPath(path string) error {
	if _, err := classifyMigrationTarget(path); err != nil {
		return err
	}
	encrypted, encErr := v.IsEncrypted(path)
	if encErr != nil {
		return fmt.Errorf("checking encryption of %s: %w", path, encErr)
	}
	if !encrypted {
		return fmt.Errorf("%s is not encrypted with fscrypt: refusing to decrypt it", path)
	}
	return nil
}

// RestoreBackup undoes a migration: the encrypted target is deleted and path+BackupSuffix moved
// back. The master key unlocks it first (removing a locked target fails with ENOKEY); an existing
// unencrypted path is never deleted (may hold new user data); a missing backup is an error.
func (v *Vault) RestoreBackup(path string) error {
	backup := path + BackupSuffix
	if _, err := os.Lstat(backup); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("backup %s does not exist", backup)
		}
		return fmt.Errorf("stat backup %s: %w", backup, err)
	}
	if _, err := os.Lstat(path); err == nil {
		encrypted, encErr := v.IsEncrypted(path)
		if encErr != nil {
			return fmt.Errorf("checking encryption of %s: %w", path, encErr)
		}
		if !encrypted {
			return fmt.Errorf("%s exists but is not encrypted: refusing to delete it — remove or restore it manually", path)
		}
		// A file-vault regular file needs no unlock before removal (no kernel key holds it busy;
		// it's just ciphertext). Only the kernel-fscrypt directory case needs its key provisioned
		// first.
		if !isRegularFileTarget(path) {
			if err := unlockForRemoval(path); err != nil {
				return fmt.Errorf("unlock %s for removal: %w", path, err)
			}
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

// unlockForRemoval provisions path's policy key so its contents can be deleted, then removes the
// key again (leaving it locked). A locked encrypted regular file takes the descriptor fallback:
// provision first, then re-read the policy now that the file opens.
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

// applyRawKeyPolicy protects the (empty) path with a fresh policy wrapped by a raw-key protector
// holding the master key, returning with the policy key provisioned: the caller copies contents
// back while the key is available, then deprovisions (see Encrypt). On failure keys are wiped and
// the keyring cleaned.
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
	// On failure wipe in-memory keys and drop the policy key (no-op if Provision never ran); on
	// success keep the key for Encrypt's copy-back.
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

// lockAndDeprovision wipes in-memory keys and removes the policy key from the keyring, leaving the
// target locked until the daemon unlocks it. The forced deprovision retries EBUSY (inodes pinned by
// the key) within the shared budget; ENOKEY counts as success.
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

// isDeprovisionBusy: the error means inodes are still pinned by the key (retry).
func isDeprovisionBusy(err error) bool {
	return classifyDeprovision(err) == depBusy
}

// isDeprovisionMissing: the error means the key is already gone from the keyring (counts as
// locked).
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

// stripImmutableFlags clears immutable and append-only flags on every file and directory under root
// (like `chattr -R -i -a`). Symlinks are skipped (not followed); entries vanishing mid-walk
// (ENOENT) are tolerated, only the root itself is fatal; other errors name the path.
func stripImmutableFlags(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			if errors.Is(walkErr, fs.ErrNotExist) && path != root {
				// Entry vanished mid-walk (live app churning its profile): nothing to strip.
				return nil
			}
			return fmt.Errorf("walk %s: %w", path, walkErr)
		}
		if d.Type()&fs.ModeSymlink != 0 {
			// Never follow links: a dangling symlink (e.g. a firefox lock file) would fail the open
			// below.
			return nil
		}
		if !d.IsDir() && !d.Type().IsRegular() {
			// These flags exist only on regular files and directories. Opening a FIFO O_RDONLY
			// blocks forever with no writer (a real hang on ~/.gnupg and browser profiles) and a
			// device node would be opened for real; skip everything else, as copyTree does.
			return nil
		}
		return stripEntryFlags(path, path == root)
	})
}

// stripEntryFlags clears the flags on one regular file or directory. isRoot suppresses the ENOENT
// tolerance (a vanished watch root is fatal). Opens use O_NONBLOCK so a race turning `path` into a
// FIFO after the WalkDir stat can't wedge the installer.
func stripEntryFlags(path string, isRoot bool) error {
	fd, openErr := unix.Open(path, unix.O_RDONLY|unix.O_NONBLOCK, 0)
	if openErr != nil {
		if errors.Is(openErr, fs.ErrNotExist) && !isRoot {
			return nil
		}
		return fmt.Errorf("open %s: %w", path, openErr)
	}
	// On error (e.g. tmpfs without chattr support) flags is zeroed, equal to the cleared value
	// below.
	flags, _, _ := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), unix.FS_IOC_GETFLAGS, 0)
	unix.Close(fd)
	cleared := flags &^ (fsImmutableFlag | fsAppendFlag)
	if cleared == flags {
		return nil
	}
	fd, openErr = unix.Open(path, unix.O_RDONLY|unix.O_NONBLOCK, 0)
	if openErr != nil {
		return fmt.Errorf("open %s: %w", path, openErr)
	}
	_, _, _ = unix.Syscall(unix.SYS_IOCTL, uintptr(fd), unix.FS_IOC_SETFLAGS, cleared)
	unix.Close(fd)
	return nil
}

func wipe(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// isKernelAtLeast54: only used to pick the default fscrypt policy version (as `fscrypt setup`
// does).
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
