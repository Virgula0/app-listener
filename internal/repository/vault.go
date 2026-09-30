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

	// ErrNotEncrypted: path has no fscrypt policy at all, distinct from ErrKeyBusy/ErrKeyMissing
	// (which presuppose one). Permanent: retrying Lock can't make a policy-less path locked, so
	// callers must not retry it like ErrKeyBusy.
	ErrNotEncrypted = errors.New("path is not encrypted")
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
	// Lock deprovisions the policy key for path. With forceFlush it is attempted even if the policy
	// doesn't look provisioned; errors map to ErrKeyBusy/ErrKeyMissing.
	Lock(path string, forceFlush bool) error
}
