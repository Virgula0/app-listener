package binledger

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// Password guess budget shared by every verifier: MaxAuthFailures guesses, then AuthLockout.
const (
	MaxAuthFailures = 5
	AuthLockout     = time.Minute
)

// LockedError refuses a guess while the lockout lasts.
type LockedError struct{ Until time.Time }

func (e *LockedError) Error() string {
	return fmt.Sprintf("too many failed attempts — locked for another %s", time.Until(e.Until).Round(time.Second))
}

// EditAuthHash returns the stored edit-protected password record, "" when none is set.
func (l *Ledger) EditAuthHash() (string, error) {
	var h string
	err := l.db.QueryRowContext(context.Background(), `SELECT hash FROM edit_auth WHERE id = 1`).Scan(&h)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return h, err
}

// SetEditAuthHash stores encoded as the password record; there is no way to remove it.
func (l *Ledger) SetEditAuthHash(encoded string) error {
	if encoded == "" {
		return errors.New("empty password record")
	}
	_, err := l.db.ExecContext(context.Background(),
		`INSERT INTO edit_auth (id, hash, updated_at) VALUES (1, ?, ?)
		 ON CONFLICT (id) DO UPDATE SET hash = excluded.hash, updated_at = excluded.updated_at`,
		encoded, time.Now().Unix())
	return err
}

// BeginAuthAttempt reserves one password guess before it is verified, refusing while locked out.
// The reservation counts as a failure until AuthSucceeded, so concurrent guessers (several
// processes, the daemon and the CLI) share one budget instead of each checking a stale count.
func (l *Ledger) BeginAuthAttempt(now time.Time) error {
	tx, err := l.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck // no-op after Commit
	var failures, lockedUntil int64
	err = tx.QueryRowContext(context.Background(),
		`SELECT failures, locked_until FROM edit_auth_throttle WHERE id = 1`).Scan(&failures, &lockedUntil)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if until := time.Unix(lockedUntil, 0); now.Before(until) {
		return &LockedError{Until: until}
	}
	failures++
	lockedUntil = 0
	if failures >= MaxAuthFailures {
		failures, lockedUntil = 0, now.Add(AuthLockout).Unix()
	}
	if _, err := tx.ExecContext(context.Background(),
		`INSERT INTO edit_auth_throttle (id, failures, locked_until) VALUES (1, ?, ?)
		 ON CONFLICT (id) DO UPDATE SET failures = excluded.failures, locked_until = excluded.locked_until`,
		failures, lockedUntil); err != nil {
		return err
	}
	return tx.Commit()
}

// AuthSucceeded clears the guess budget after a correct password.
func (l *Ledger) AuthSucceeded() error {
	_, err := l.db.ExecContext(context.Background(), `DELETE FROM edit_auth_throttle`)
	return err
}
