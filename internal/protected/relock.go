package protected

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/fscrypt"
	"github.com/Virgula0/app-listener/internal/repository"
)

// Relock retry budget for the force-flush deprovision, mirroring the daemon
// teardown (internal/usecase) and the fscrypt migration lock
// (internal/fscrypt): a forced deprovision can fail with EBUSY while an
// inode is still pinned, and the key must be gone before the caller
// returns.
const (
	maxRelockRetries = 100
	relockRetryDelay = 10 * time.Millisecond
)

// RelockResources locks path with the same two-pass approach as the
// daemon's teardown: a plain (non-forced) lock first, then a force-flush
// with a bounded EBUSY retry. A key that is already gone (ErrKeyMissing)
// counts as fully locked.
func RelockResources(vault *fscrypt.Vault, path string) error {
	if err := vault.Lock(path, false); err != nil &&
		!errors.Is(err, repository.ErrKeyBusy) && !errors.Is(err, repository.ErrKeyMissing) {
		return err
	}
	for i := 0; i < maxRelockRetries; i++ {
		err := vault.Lock(path, true)
		switch {
		case err == nil || errors.Is(err, repository.ErrKeyMissing):
			log.Infof("%s is locked again", path)
			return nil
		case errors.Is(err, repository.ErrKeyBusy):
			time.Sleep(relockRetryDelay)
		default:
			return err
		}
	}
	return fmt.Errorf("%s is still busy after %d relock retries — a process may hold open files in it (investigate with: lsof +D %s, fuser -v %s)",
		path, maxRelockRetries, path, path)
}
