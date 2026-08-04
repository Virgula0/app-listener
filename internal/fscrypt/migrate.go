package fscrypt

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/Virgula0/app-listener/internal/install"
	"github.com/google/fscrypt/actions"
	"github.com/google/fscrypt/metadata"
	"golang.org/x/sys/unix"
	"google.golang.org/protobuf/proto"
)

// BackupSuffix is appended to a directory to form its migration backup,
// e.g. /home/alice/.ssh -> /home/alice/.ssh.app_listener.backup.
const BackupSuffix = ".app_listener.backup"

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
		return fmt.Errorf("get policy for %s: %w", path, err)
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

// Encrypt migrates path to an fscrypt-encrypted directory in place, using
// the master key as a raw-key protector (no shelling out — the google/fscrypt
// actions API only). The original contents are preserved by renaming path to
// path+BackupSuffix BEFORE anything is touched; the backup is left in place
// for the installer to remove after the whole installation succeeded.
//
// If path+BackupSuffix already exists, Encrypt refuses to run: a previous
// migration may have failed midway, and overwriting an older backup would
// destroy data. The caller must resolve that situation manually.
//
// On success the directory is encrypted and locked: its key is removed from
// the kernel keyring, so the daemon must unlock it with the master key at
// startup.
func (v *Vault) Encrypt(path string) error {
	return v.EncryptWithProgress(path, nil)
}

// EncryptWithProgress behaves like Encrypt and additionally reports copy
// progress (bytes copied so far and total) through onBytes while the
// contents are moved into the encrypted directory.
func (v *Vault) EncryptWithProgress(path string, onBytes func(copied, total int64)) error {
	backup := path + BackupSuffix
	if _, err := os.Lstat(backup); err == nil {
		return fmt.Errorf("backup %s already exists: a previous migration was interrupted; restore or remove it manually before retrying", backup)
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("encrypt %s: only directories can be migrated", path)
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
	// the daemon unlocks it. Both are best-effort by design — the data is
	// already encrypted at rest, and a fatal error here would abort the
	// installer after a finished migration (and trip the backup-exists
	// refusal on retry).
	_ = policy.Lock()
	_ = policy.Deprovision(false)
	return nil
}

// RestoreBackup undoes a migration: the encrypted directory at path is
// deleted and path+BackupSuffix is moved back to path, restoring the
// original unencrypted content. The encrypted directory is unlocked with
// the master key first — deleting entries of a locked encrypted directory
// fails with ENOKEY because the kernel needs the key to touch the names.
// It refuses to delete a path that exists but is NOT encrypted (the user
// may have created new data there since the migration), and it is a no-op
// error when no backup exists.
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

// unlockForRemoval provisions the policy key of an encrypted directory so
// its contents can be deleted, then removes the key from the keyring again
// (leaving the directory locked) before returning.
func unlockForRemoval(path string) error {
	fsctx, err := actions.NewContextFromPath(path, nil)
	if err != nil {
		return err
	}
	policy, err := actions.GetPolicyFromPath(fsctx, path)
	if err != nil {
		return err
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
	defer lockAndDeprovision(policy, nil)
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
			lockAndDeprovision(policy, protector)
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
// directory locked. Best-effort: every failure is ignored.
func lockAndDeprovision(policy *actions.Policy, protector *actions.Protector) {
	_ = policy.Lock()
	if protector != nil {
		_ = protector.Lock()
	}
	_ = policy.Deprovision(false)
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
