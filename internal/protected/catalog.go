package protected

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/fscrypt"
)

// vaultKeyCheck is the subset of *fscrypt.Vault needed for the key
// verification step; an interface keeps this package decoupled and
// headlessly testable.
type vaultKeyCheck interface {
	IsEncrypted(path string) (bool, error)
	VerifyKey(path string) error
}

// VerifyEncryptedKeys checks the master key against every encrypted
// directory. A directory that does not unlock with the master key is a
// fatal error: removing the daemon (and possibly the key with
// --delete-key) would lock that directory forever since no usable key
// remains.
func VerifyEncryptedKeys(vault vaultKeyCheck, encrypted []string) error {
	for _, path := range encrypted {
		log.Infof("verifying master key against %s ...", path)
		if err := vault.VerifyKey(path); err != nil {
			return fmt.Errorf("fatal: %v — the master key %s does not match the policy of %s; restore the correct key before continuing",
				err, fscrypt.MasterKeyFile, path)
		}
	}
	return nil
}
