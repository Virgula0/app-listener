package repository

import "errors"

// Sentinel errors returned by a Vault implementation, used by callers to
// drive the fscrypt teardown retry loop (see the daemon usecase).
var (
	// ErrKeyBusy reports that some inodes using the key are still open;
	// the caller should retry shortly.
	ErrKeyBusy = errors.New("key still in use")

	// ErrKeyMissing reports that the kernel confirms the key is gone
	// (ENOKEY): the directory is fully locked.
	ErrKeyMissing = errors.New("key not present")
)

// Vault abstracts the fscrypt lifecycle of an encrypted directory.
type Vault interface {
	// IsEncrypted reports whether path carries an fscrypt policy (v1 or v2).
	IsEncrypted(path string) (bool, error)
	// IsProvisioned reports whether the policy key for path is currently
	// provisioned (i.e. Unlock would be a no-op).
	IsProvisioned(path string) (bool, error)
	// Unlock provisions the policy key for path so its contents are
	// readable. It is a no-op when the policy is already provisioned.
	Unlock(path string) error
	// Lock deprovisions the policy key for path. When forceFlush is true
	// the deprovision is attempted even if the policy does not appear
	// provisioned, and errors are translated to ErrKeyBusy/ErrKeyMissing.
	Lock(path string, forceFlush bool) error
}
