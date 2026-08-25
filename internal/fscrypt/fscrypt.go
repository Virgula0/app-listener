// Package fscrypt ports the fscrypt subsystem of the ssh-guard daemon: a 32-byte
// raw master key, policy detection, and the unlock/lock lifecycle for encrypted
// directories, preserving its tested EBUSY-retry/ENOKEY teardown semantics.
package fscrypt

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/google/fscrypt/actions"
	"github.com/google/fscrypt/crypto"
	"github.com/google/fscrypt/filesystem"
	"github.com/google/fscrypt/metadata"
	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/repository"
)

const (
	// MasterKeyFile is the default location of the 32-byte raw master key.
	MasterKeyFile = "/etc/app-listener/fscrypt.key"
	// FscryptKeySize is the AES-256-CTS master key length in bytes.
	FscryptKeySize = 32
)

// FS_IOC_GET_ENCRYPTION_POLICY_EX detects any fscrypt policy (v1 or v2).
const linux_FS_IOC_GET_ENCRYPTION_POLICY_EX = 0xc0096616

// Vault implements repository.Vault on top of the google/fscrypt actions API.
type Vault struct{}

// New returns a Vault rooted at the default master key file.
func New() *Vault {
	return &Vault{}
}

// hasEncryptionPolicy reports whether the file or directory at path carries an
// fscrypt policy (v1 or v2). A LOCKED encrypted regular file cannot be opened at
// all — ENOKEY there itself proves a policy exists. Anything neither directory
// nor regular file never carries one.
func hasEncryptionPolicy(path string) (bool, error) {
	info, statErr := os.Stat(path)
	if statErr != nil {
		return false, fmt.Errorf("stat %s: %w", path, statErr)
	}

	openFlags := unix.O_RDONLY
	if info.IsDir() {
		openFlags |= unix.O_DIRECTORY
	} else if !info.Mode().IsRegular() {
		return false, nil
	}

	fd, err := unix.Open(path, openFlags, 0)
	if err != nil {
		if errors.Is(err, unix.ENOKEY) {
			// Locked encrypted regular file: policy exists, key absent.
			return true, nil
		}
		if errors.Is(err, unix.ENOTDIR) {
			return false, nil
		}
		return false, fmt.Errorf("open %s: %w", path, err)
	}
	defer unix.Close(fd)

	var arg unix.FscryptGetPolicyExArg
	arg.Size = 24
	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(fd),
		uintptr(linux_FS_IOC_GET_ENCRYPTION_POLICY_EX),
		uintptr(unsafe.Pointer(&arg)),
	)
	if errno == 0 {
		return true, nil
	}
	// ENODATA: no policy set. EOPNOTSUPP/ENOTTY: the filesystem cannot carry
	// encryption policies at all.
	if errors.Is(errno, unix.ENODATA) ||
		errors.Is(errno, unix.EOPNOTSUPP) ||
		errors.Is(errno, unix.ENOTTY) {
		return false, nil
	}
	return false, fmt.Errorf("ioctl GET_ENCRYPTION_POLICY_EX on %s: %w", path, errno)
}

// IsEncrypted reports whether path carries an fscrypt policy.
func (v *Vault) IsEncrypted(path string) (bool, error) {
	return hasEncryptionPolicy(path)
}

// isLockedRegularFileErr reports whether err means "encrypted regular file whose
// key is not provisioned" (metadata.ErrLockedRegularFile): the kernel refuses to
// open such a file at all.
func isLockedRegularFileErr(err error) bool {
	var locked *metadata.ErrLockedRegularFile
	return errors.As(err, &locked)
}

// acceptFirstProtectorOption is the protector-picker shared by every policy
// unlock: there is exactly one protector per policy here.
var acceptFirstProtectorOption = func(_ string, _ []*actions.ProtectorOption) (int, error) {
	return 0, nil
}

// checkKeyLen validates the master key length against FscryptKeySize.
func checkKeyLen(keyFile string, data []byte) error {
	if len(data) != FscryptKeySize {
		return fmt.Errorf("master key %s must be exactly %d bytes, got %d",
			keyFile, FscryptKeySize, len(data))
	}
	return nil
}

// deprovisionKind classifies a failed forced deprovision for the retry
// loops shared with internal/usecase and internal/protected.
type deprovisionKind int

const (
	depOther   deprovisionKind = iota
	depBusy                    // inodes still pinned: retry after a delay
	depMissing                 // key already gone: counts as fully locked
)

// classifyDeprovision maps the fscrypt library's deprovision error strings onto a
// deprovisionKind (substrings come from the kernel keyring/sysfs errors it surfaces).
func classifyDeprovision(err error) deprovisionKind {
	errStr := err.Error()
	switch {
	case strings.Contains(errStr, "some files using the key are still open"),
		strings.Contains(errStr, "in use"):
		return depBusy
	case strings.Contains(errStr, "Required key not available"),
		strings.Contains(errStr, "key not present"):
		return depMissing
	default:
		return depOther
	}
}

// provisionLockedFile unlocks a locked encrypted REGULAR file: unopenable before its key is provisioned,
// so its policy is unaddressable via the path — every policy on the filesystem is tried against the master
// key, tentatively provisioned until the file opens; wrong candidates are deprovisioned immediately.
func (v *Vault) provisionLockedFile(fsctx *actions.Context, path string) error {
	keyBytes, err := readKey()
	if err != nil {
		return err
	}
	defer wipe(keyBytes)

	descriptors, listErr := fsctx.Mount.ListPolicies(nil)
	if listErr != nil {
		return fmt.Errorf("listing policies on %s: %w", fsctx.Mount.Path, listErr)
	}

	optionFn := acceptFirstProtectorOption
	for _, descriptor := range descriptors {
		candidate, getErr := actions.GetPolicy(fsctx, descriptor)
		if getErr != nil {
			continue // unreadable metadata entry: skip
		}
		if candidate.IsProvisionedByTargetUser() {
			continue // already provisioned: cannot be what blocks the file
		}
		if unlockErr := candidate.Unlock(optionFn, newBoundedKeyFn(keyBytes)); unlockErr != nil {
			continue // wrapped by a different protector: not ours
		}
		if provErr := candidate.Provision(); provErr != nil {
			continue
		}
		fd, openErr := unix.Open(path, unix.O_RDONLY, 0)
		if openErr == nil {
			unix.Close(fd)
			return nil // this was the file's own policy: it opens now
		}
		// Wrong policy: undo the collateral provisioning immediately.
		_ = candidate.Deprovision(true)
	}
	return fmt.Errorf("no policy on %s could be unlocked with %s to open the locked file %s",
		fsctx.Mount.Path, MasterKeyFile, path)
}

// verifyLockedFileKey validates the master key against a locked encrypted regular file, provisioning
// nothing: its own policy is unaddressable while locked, so success means some policy on this
// filesystem unwraps with the master key (the single-key model).
func (v *Vault) verifyLockedFileKey(fsctx *actions.Context, path string) error {
	keyBytes, err := readKey()
	if err != nil {
		return fmt.Errorf("verify key for %s: %w", path, err)
	}
	defer wipe(keyBytes)

	descriptors, listErr := fsctx.Mount.ListPolicies(nil)
	if listErr != nil {
		return fmt.Errorf("verify key for %s: listing policies on %s: %w", path, fsctx.Mount.Path, listErr)
	}
	optionFn := acceptFirstProtectorOption
	for _, descriptor := range descriptors {
		candidate, getErr := actions.GetPolicy(fsctx, descriptor)
		if getErr != nil {
			continue
		}
		if candidate.IsProvisionedByTargetUser() {
			return nil // provisioned policies belong to this setup
		}
		if unlockErr := candidate.Unlock(optionFn, newBoundedKeyFn(keyBytes)); unlockErr == nil {
			_ = candidate.Lock()
			return nil
		}
	}
	return fmt.Errorf("no policy on %s unlocks with %s: the locked file %s would stay unreadable",
		fsctx.Mount.Path, MasterKeyFile, path)
}

// IsProvisioned reports whether the policy key for path is currently provisioned (unlocked).
func (v *Vault) IsProvisioned(path string) (bool, error) {
	fsctx, err := actions.NewContextFromPath(path, nil)
	if err != nil {
		return false, fmt.Errorf("fscrypt context for %s: %w", path, err)
	}
	policy, err := actions.GetPolicyFromPath(fsctx, path)
	if err != nil {
		if isLockedRegularFileErr(err) {
			return false, nil // locked by definition: not provisioned
		}
		return false, fmt.Errorf("get policy for %s: %w", path, err)
	}
	return policy.IsProvisionedByTargetUser(), nil
}

// CheckFilesystemReady verifies the filesystem backing path is initialized for fscrypt (`fscrypt setup`)
// and actually has encryption enabled (the ext4 `encrypt` feature flag) — the fscrypt CLI's non-destructive
// pre-check — so failures surface here instead of deep inside later policy creation.
func (v *Vault) CheckFilesystemReady(path string) error {
	mnt, err := filesystem.FindMount(path)
	if err != nil {
		return fmt.Errorf("resolve filesystem of %s: %w", path, err)
	}
	if err := mnt.CheckSetup(nil); err != nil {
		return classifySetupError(path, err)
	}
	if err := mnt.CheckSupport(); err != nil {
		return classifySupportError(path, err)
	}
	return nil
}

// classifySetupError translates fscrypt setup failures into actionable errors naming the remediation.
func classifySetupError(path string, err error) error {
	var notSetup *filesystem.ErrNotSetup
	var notSupported *filesystem.ErrSetupNotSupported
	switch {
	case errors.As(err, &notSetup):
		return fmt.Errorf(
			"filesystem containing %s is not initialized for fscrypt (%v). "+
				"Run once as root: 'fscrypt setup --all-users' (or 'fscrypt setup /'), then re-run the installer",
			path, err)
	case errors.As(err, &notSupported):
		return fmt.Errorf(
			"filesystem containing %s does not support fscrypt encryption (%v). "+
				"Only ext4/f2fs with encryption enabled can be protected",
			path, err)
	default:
		return fmt.Errorf("fscrypt setup check for %s: %w", path, err)
	}
}

// classifySupportError translates encryption-support probe failures into actionable errors, catching
// filesystems that are "set up" but cannot actually encrypt (typically ext4 created without the
// `encrypt` feature flag).
func classifySupportError(path string, err error) error {
	var notEnabled *filesystem.ErrEncryptionNotEnabled
	var notSupported *filesystem.ErrEncryptionNotSupported
	switch {
	case errors.As(err, &notEnabled):
		enable := "sudo tune2fs -O encrypt " + notEnabled.Mount.Device
		if notEnabled.Mount.FilesystemType == "f2fs" {
			enable = "sudo fsck.f2fs -O encrypt " + notEnabled.Mount.Device
		}
		return fmt.Errorf(
			"filesystem %s (%s) containing %s is not marked for fscrypt encryption. "+
				"Run %q as root, then re-run the installer",
			notEnabled.Mount.Device, notEnabled.Mount.FilesystemType, path, enable)
	case errors.As(err, &notSupported):
		return fmt.Errorf(
			"filesystem containing %s does not support fscrypt encryption (%v). "+
				"This kernel lacks support for encryption on %s filesystems",
			path, err, notSupported.Mount.FilesystemType)
	default:
		return fmt.Errorf("fscrypt support check for %s: %w", path, err)
	}
}

// readKey returns the master key, failing when the file is missing: a freshly
// generated key can never unlock an existing policy, so absence is always
// misconfiguration to surface, never paper over.
func readKey() ([]byte, error) {
	return readKeyFrom(MasterKeyFile)
}

func readKeyFrom(keyFile string) ([]byte, error) {
	data, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("read master key %s: %w", keyFile, err)
	}
	if err := checkKeyLen(keyFile, data); err != nil {
		return nil, err
	}
	return data, nil
}

// readOrCreateKey returns the master key, generating it on first use.
// Generation is atomic (O_EXCL); a lost race re-reads the winner's key.
func readOrCreateKey() ([]byte, error) {
	if data, err := os.ReadFile(MasterKeyFile); err == nil {
		if err := checkKeyLen(MasterKeyFile, data); err != nil {
			return nil, err
		}
		return data, nil
	}

	key := make([]byte, FscryptKeySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(MasterKeyFile), 0o700); err != nil {
		return nil, err
	}
	f, err := os.OpenFile(MasterKeyFile, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o400)
	if err != nil {
		if os.IsExist(err) {
			data, readErr := os.ReadFile(MasterKeyFile)
			if readErr != nil {
				return nil, readErr
			}
			if lenErr := checkKeyLen(MasterKeyFile, data); lenErr != nil {
				return nil, lenErr
			}
			return data, nil
		}
		return nil, err
	}
	defer f.Close()
	if _, err := f.Write(key); err != nil {
		os.Remove(MasterKeyFile)
		return nil, err
	}
	return key, nil
}

// MasterKeyExists reports whether the master key file is present.
func MasterKeyExists() (bool, error) {
	_, err := os.Stat(MasterKeyFile)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GenerateMasterKey creates the master key file: a no-op when it exists and force is false, otherwise
// an atomic replacement (temp file + rename). Replacing the key invalidates every directory provisioned
// with the old one — callers must confirm first.
func GenerateMasterKey(force bool) error {
	if !force {
		if _, err := os.Stat(MasterKeyFile); err == nil {
			return nil
		}
		_, err := readOrCreateKey()
		return err
	}

	key := make([]byte, FscryptKeySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(MasterKeyFile), 0o700); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(MasterKeyFile), ".fscrypt.key.*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if err := tmp.Chmod(0o400); err != nil {
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(key); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpPath, MasterKeyFile)
}

// newBoundedKeyFn bounds the library's unwrap loop: the callback is re-invoked with retry=true on each
// failed protector unwrap and only a callback error aborts, so a constant wrong key would spin forever.
// With one master key retries never succeed — failure ends the loop with an error naming the key file.
func newBoundedKeyFn(masterKey []byte) actions.KeyFunc {
	return func(info actions.ProtectorInfo, retry bool) (*crypto.Key, error) {
		if retry {
			return nil, fmt.Errorf("invalid wrapping key for protector %s: master key %s does not match the policy",
				info.Descriptor(), MasterKeyFile)
		}
		return crypto.NewFixedLengthKeyFromReader(bytes.NewReader(masterKey), FscryptKeySize)
	}
}

// Unlock provisions the policy key of path so its contents become readable;
// no-op when already provisioned by the target user.
func (v *Vault) Unlock(path string) error {
	fsctx, err := actions.NewContextFromPath(path, nil)
	if err != nil {
		return fmt.Errorf("fscrypt context for %s: %w", path, err)
	}
	policy, err := actions.GetPolicyFromPath(fsctx, path)
	if err != nil {
		if !isLockedRegularFileErr(err) {
			return fmt.Errorf("get policy for %s: %w", path, err)
		}
		// Locked file: unopenable, so fall back to provisionLockedFile.
		if provErr := v.provisionLockedFile(fsctx, path); provErr != nil {
			return fmt.Errorf("unlock %s: %w", path, provErr)
		}
		return nil
	}
	if policy.IsProvisionedByTargetUser() {
		return nil
	}

	keyBytes, err := readKey()
	if err != nil {
		return fmt.Errorf("unlock %s: %w", path, err)
	}
	defer wipe(keyBytes)

	optionFn := acceptFirstProtectorOption
	if err := policy.Unlock(optionFn, newBoundedKeyFn(keyBytes)); err != nil {
		return fmt.Errorf("unlock policy for %s: %w", path, err)
	}
	defer func() { _ = policy.Lock() }()
	if err := policy.Provision(); err != nil {
		return fmt.Errorf("provision policy for %s: %w", path, err)
	}
	return nil
}

// Lock deprovisions the policy key of path; skipped when not provisioned by
// the target user and forceFlush is false (already locked). Errors map onto
// repository.ErrKeyBusy (retry) / repository.ErrKeyMissing (fully locked).
func (v *Vault) Lock(path string, forceFlush bool) error {
	fsctx, err := actions.NewContextFromPath(path, nil)
	if err != nil {
		return fmt.Errorf("fscrypt context for %s: %w", path, err)
	}
	policy, err := actions.GetPolicyFromPath(fsctx, path)
	if err != nil {
		if isLockedRegularFileErr(err) {
			// The file cannot even be opened: its key is already gone.
			return fmt.Errorf("%w: %s is a locked encrypted regular file", repository.ErrKeyMissing, path)
		}
		return fmt.Errorf("get policy for %s: %w", path, err)
	}

	if !forceFlush && !policy.IsProvisionedByTargetUser() {
		return nil
	}

	if err := policy.Deprovision(true); err != nil {
		return classifyDeprovisionErr(path, err)
	}
	return nil
}

// classifyDeprovisionErr maps deprovision errors onto the repository sentinels:
// EBUSY-style (inodes still pinned) → ErrKeyBusy, ENOKEY (key already gone) →
// ErrKeyMissing, anything else wrapped with the path.
func classifyDeprovisionErr(path string, err error) error {
	switch classifyDeprovision(err) {
	case depBusy:
		return fmt.Errorf("%w: %v", repository.ErrKeyBusy, err)
	case depMissing:
		return fmt.Errorf("%w: %v", repository.ErrKeyMissing, err)
	default:
		return fmt.Errorf("deprovision policy for %s: %w", path, err)
	}
}
