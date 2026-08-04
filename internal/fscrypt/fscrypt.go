// Package fscrypt ports the fscrypt subsystem of the ssh-guard daemon: a
// 32-byte raw master key, policy detection, and the unlock/lock lifecycle
// for encrypted directories. The teardown semantics (EBUSY retry, ENOKEY
// detection) were designed and tested against TOC-TOU races when tearing
// up and down encrypted directories, and are preserved verbatim.
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

// hasEncryptionPolicy reports whether the directory at dirPath carries an
// fscrypt encryption policy (v1 or v2).
func hasEncryptionPolicy(dirPath string) (bool, error) {
	dirFd, err := unix.Open(dirPath, unix.O_RDONLY|unix.O_DIRECTORY, 0)
	if err != nil {
		return false, fmt.Errorf("open %s: %w", dirPath, err)
	}
	defer unix.Close(dirFd)

	var arg unix.FscryptGetPolicyExArg
	arg.Size = 24
	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		uintptr(dirFd),
		uintptr(linux_FS_IOC_GET_ENCRYPTION_POLICY_EX),
		uintptr(unsafe.Pointer(&arg)),
	)
	if errno == 0 {
		return true, nil
	}
	if errors.Is(errno, unix.ENODATA) || errors.Is(errno, unix.EOPNOTSUPP) {
		return false, nil
	}
	return false, fmt.Errorf("ioctl GET_ENCRYPTION_POLICY_EX on %s: %w", dirPath, errno)
}

// IsEncrypted reports whether path carries an fscrypt policy.
func (v *Vault) IsEncrypted(path string) (bool, error) {
	return hasEncryptionPolicy(path)
}

// IsProvisioned reports whether the policy key for path is currently
// provisioned, i.e. the directory is unlocked.
func (v *Vault) IsProvisioned(path string) (bool, error) {
	fsctx, err := actions.NewContextFromPath(path, nil)
	if err != nil {
		return false, fmt.Errorf("fscrypt context for %s: %w", path, err)
	}
	policy, err := actions.GetPolicyFromPath(fsctx, path)
	if err != nil {
		return false, fmt.Errorf("get policy for %s: %w", path, err)
	}
	return policy.IsProvisionedByTargetUser(), nil
}

// readKey returns the 32-byte master key, failing when the file is
// missing. The key is never generated implicitly: for an already encrypted
// directory a freshly generated key can never unlock the policy, so a
// missing file is always a misconfiguration to surface, not to paper over.
func readKey() ([]byte, error) {
	return readKeyFrom(MasterKeyFile)
}

func readKeyFrom(keyFile string) ([]byte, error) {
	data, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("read master key %s: %w", keyFile, err)
	}
	if len(data) != FscryptKeySize {
		return nil, fmt.Errorf("master key %s must be exactly %d bytes, got %d",
			keyFile, FscryptKeySize, len(data))
	}
	return data, nil
}

// readOrCreateKey returns the 32-byte master key, generating it on first
// use. Generation is atomic (O_EXCL); a lost race re-reads the winner's key.
func readOrCreateKey() ([]byte, error) {
	if data, err := os.ReadFile(MasterKeyFile); err == nil {
		if len(data) != FscryptKeySize {
			return nil, fmt.Errorf("master key %s must be exactly %d bytes, got %d",
				MasterKeyFile, FscryptKeySize, len(data))
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
			if len(data) != FscryptKeySize {
				return nil, fmt.Errorf("master key %s must be exactly %d bytes, got %d",
					MasterKeyFile, FscryptKeySize, len(data))
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

// GenerateMasterKey creates the 32-byte master key file. When force is
// false it is a no-op if the key already exists; when force is true the
// existing key is atomically replaced (write to a temp file in the same
// directory, then rename). Replacing the key invalidates every fscrypt
// directory provisioned with the old one — callers must confirm first.
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

// newBoundedKeyFn returns the key callback passed to policy.Unlock. The
// google/fscrypt unwrap loop (actions/callback.go, unwrapProtectorKey)
// re-invokes the callback with retry=true whenever the returned key fails
// to unwrap a protector, and only a callback error aborts the loop — a
// callback that keeps returning the same wrong key spins forever, logging
// "invalid wrapping key for protector ..." at full speed. Because the
// daemon has exactly one master key, a retry can never succeed: the first
// failed unwrap is answered with an error that terminates the loop and
// names the key file in question.
func newBoundedKeyFn(masterKey []byte) actions.KeyFunc {
	return func(info actions.ProtectorInfo, retry bool) (*crypto.Key, error) {
		if retry {
			return nil, fmt.Errorf("invalid wrapping key for protector %s: master key %s does not match the policy",
				info.Descriptor(), MasterKeyFile)
		}
		return crypto.NewFixedLengthKeyFromReader(bytes.NewReader(masterKey), FscryptKeySize)
	}
}

// Unlock provisions the policy key of path so its contents become
// readable. It is a no-op when the policy is already provisioned by the
// target user.
func (v *Vault) Unlock(path string) error {
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
		return fmt.Errorf("unlock %s: %w", path, err)
	}
	defer func() {
		for i := range keyBytes {
			keyBytes[i] = 0
		}
	}()

	optionFn := func(_ string, _ []*actions.ProtectorOption) (int, error) {
		return 0, nil
	}
	if err := policy.Unlock(optionFn, newBoundedKeyFn(keyBytes)); err != nil {
		return fmt.Errorf("unlock policy for %s: %w", path, err)
	}
	defer func() { _ = policy.Lock() }()
	if err := policy.Provision(); err != nil {
		return fmt.Errorf("provision policy for %s: %w", path, err)
	}
	return nil
}

// Lock deprovisions the policy key of path. When forceFlush is false and
// the policy is not provisioned by the target user, the deprovision is
// skipped (already locked). Errors are translated into repository.ErrKeyBusy
// (retry) and repository.ErrKeyMissing (fully locked) sentinels.
func (v *Vault) Lock(path string, forceFlush bool) error {
	fsctx, err := actions.NewContextFromPath(path, nil)
	if err != nil {
		return fmt.Errorf("fscrypt context for %s: %w", path, err)
	}
	policy, err := actions.GetPolicyFromPath(fsctx, path)
	if err != nil {
		return fmt.Errorf("get policy for %s: %w", path, err)
	}

	if !forceFlush && !policy.IsProvisionedByTargetUser() {
		return nil
	}

	if err := policy.Deprovision(true); err != nil {
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "some files using the key are still open"),
			strings.Contains(errStr, "in use"):
			return fmt.Errorf("%w: %v", repository.ErrKeyBusy, err)
		case strings.Contains(errStr, "Required key not available"),
			strings.Contains(errStr, "key not present"):
			return fmt.Errorf("%w: %v", repository.ErrKeyMissing, err)
		default:
			return fmt.Errorf("deprovision policy for %s: %w", path, err)
		}
	}
	return nil
}
