// Package binledger is the daemon's record of every whitelisted binary it has confirmed, by content
// hash. Admission at start or reload otherwise trusts whatever sits at a whitelisted path, and
// nothing guards those paths while the daemon is not running (issue #80): a binary whose hash is
// not recorded is admitted only when the daemon can prove where it came from (cmd/functions/daemon
// binvet.go), else it waits for `app-listener trust-binaries`.
package binledger

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3" // database/sql driver

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

// DefaultPath lives in /etc/app-listener, which the daemon self-guards read-only: only the
// app-listener binary may modify it.
const DefaultPath = "/etc/app-listener/binaries.db"

// Source says how a hash entered the ledger.
type Source string

const (
	SourceInstall   Source = "install"
	SourceBootstrap Source = "bootstrap"
	SourceUpdater   Source = "updater"
	SourceConfirmed Source = "confirmed"
)

// Pending is a binary the daemon refused because its hash is not recorded for its line.
type Pending struct {
	Line      string
	Resolved  string
	SHA256    [32]byte
	Resources string
	SeenAt    time.Time
}

// Ledger is an open binaries.db.
type Ledger struct {
	db *sql.DB
}

const schema = `
CREATE TABLE IF NOT EXISTS meta (key TEXT PRIMARY KEY, value TEXT NOT NULL);
CREATE TABLE IF NOT EXISTS binaries (
	line TEXT PRIMARY KEY,
	sha256 BLOB NOT NULL,
	source TEXT NOT NULL,
	updated_at INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS grants (line TEXT PRIMARY KEY);
CREATE TABLE IF NOT EXISTS pending (
	line TEXT PRIMARY KEY,
	resolved TEXT NOT NULL,
	sha256 BLOB NOT NULL,
	resources TEXT NOT NULL,
	seen_at INTEGER NOT NULL
);`

// Line is the whitelist line the ledger keys r on: the link when r is a symlink line, so a link
// re-pointed at another confirmed binary is judged by that binary's content, not its path.
func Line(r daemonconfig.BinaryRule) string {
	if r.Link != "" {
		return r.Link
	}
	return r.Path
}

// Lines are the ledger lines of every binary cfg whitelists, pending ones included.
func Lines(cfg *daemonconfig.Config) []string {
	var out []string
	for i := range cfg.Resources {
		for _, list := range [][]daemonconfig.BinaryRule{cfg.Resources[i].Binaries, cfg.Resources[i].PendingBinaries} {
			for _, b := range list {
				out = append(out, Line(b))
			}
		}
	}
	return out
}

// JournalPath is the rollback journal beside path. journal_mode=PERSIST keeps it on disk between
// transactions: the self guard lets the daemon rewrite an existing file there, never create one.
func JournalPath(path string) string { return path + "-journal" }

// EnsurePlaceholders creates path and its journal empty (0600) when missing, so writes keep working
// once the read-only self guard over their directory is attached. No-op when the directory is
// missing (nothing installed).
func EnsurePlaceholders(path string) error {
	if _, err := os.Stat(filepath.Dir(path)); err != nil {
		return nil //nolint:nilerr // no install yet: nothing to pre-create
	}
	for _, p := range []string{path, JournalPath(path)} {
		f, err := os.OpenFile(p, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if errors.Is(err, os.ErrExist) {
			continue
		}
		if err != nil {
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Open opens (creating the schema in) the ledger at path.
func Open(path string) (*Ledger, error) {
	if err := EnsurePlaceholders(path); err != nil {
		return nil, err
	}
	dsn := "file:" + path + "?_journal_mode=PERSIST&_busy_timeout=10000&_txlock=immediate&_sync=FULL"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	if _, err := db.ExecContext(context.Background(), schema); err != nil {
		db.Close()
		return nil, fmt.Errorf("binary ledger %s: %w", path, err)
	}
	return &Ledger{db: db}, nil
}

// Close closes the ledger.
func (l *Ledger) Close() error { return l.db.Close() }

const initializedKey = "initialized"

// Initialized reports whether a run already recorded the whitelist: until then (an upgrade from a
// version without the ledger) the daemon records what it finds once.
func (l *Ledger) Initialized() (bool, error) {
	var v string
	err := l.db.QueryRowContext(context.Background(), `SELECT value FROM meta WHERE key = ?`, initializedKey).Scan(&v)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return err == nil, err
}

// MarkInitialized ends the one-time bootstrap.
func (l *Ledger) MarkInitialized() error {
	_, err := l.db.ExecContext(context.Background(), `INSERT OR REPLACE INTO meta (key, value) VALUES (?, ?)`, initializedKey,
		time.Now().UTC().Format(time.RFC3339))
	return err
}

// Lookup returns the hash recorded for line.
func (l *Ledger) Lookup(line string) (sum [32]byte, ok bool, err error) {
	var b []byte
	err = l.db.QueryRowContext(context.Background(), `SELECT sha256 FROM binaries WHERE line = ?`, line).Scan(&b)
	if errors.Is(err, sql.ErrNoRows) {
		return sum, false, nil
	}
	if err != nil {
		return sum, false, err
	}
	if len(b) != len(sum) {
		return sum, false, fmt.Errorf("binary ledger: malformed hash for %s", line)
	}
	copy(sum[:], b)
	return sum, true, nil
}

// Record stores sum as line's confirmed content and drops any pending row for it.
func (l *Ledger) Record(line string, sum [32]byte, src Source) error {
	tx, err := l.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck // no-op after Commit
	if _, err := tx.ExecContext(context.Background(), `INSERT OR REPLACE INTO binaries (line, sha256, source, updated_at) VALUES (?, ?, ?, ?)`,
		line, sum[:], string(src), time.Now().Unix()); err != nil {
		return err
	}
	if _, err := tx.ExecContext(context.Background(), `DELETE FROM pending WHERE line = ?`, line); err != nil {
		return err
	}
	return tx.Commit()
}

// AddPending records a refused binary for `trust-binaries` to show.
func (l *Ledger) AddPending(p *Pending) error {
	_, err := l.db.ExecContext(context.Background(), `INSERT OR REPLACE INTO pending (line, resolved, sha256, resources, seen_at) VALUES (?, ?, ?, ?, ?)`,
		p.Line, p.Resolved, p.SHA256[:], p.Resources, p.SeenAt.Unix())
	return err
}

// ListPending returns every refused binary, oldest first.
func (l *Ledger) ListPending() ([]Pending, error) {
	rows, err := l.db.QueryContext(context.Background(), `SELECT line, resolved, sha256, resources, seen_at FROM pending ORDER BY seen_at, line`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Pending
	for rows.Next() {
		var p Pending
		var sum []byte
		var seen int64
		if err := rows.Scan(&p.Line, &p.Resolved, &sum, &p.Resources, &seen); err != nil {
			return nil, err
		}
		copy(p.SHA256[:], sum)
		p.SeenAt = time.Unix(seen, 0)
		out = append(out, p)
	}
	return out, rows.Err()
}

// DropPending removes line's pending row (declined, or no longer configured).
func (l *Ledger) DropPending(line string) error {
	_, err := l.db.ExecContext(context.Background(), `DELETE FROM pending WHERE line = ?`, line)
	return err
}

// Grant lets the next daemon start record line on first sight, whatever its content: the install
// wizard's confirmation of a section it just added, whose binaries may sit in a vault it can't read.
func (l *Ledger) Grant(line string) error {
	_, err := l.db.ExecContext(context.Background(), `INSERT OR IGNORE INTO grants (line) VALUES (?)`, line)
	return err
}

// TakeGrant consumes line's grant, reporting whether there was one.
func (l *Ledger) TakeGrant(line string) (bool, error) {
	res, err := l.db.ExecContext(context.Background(), `DELETE FROM grants WHERE line = ?`, line)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n > 0, err
}

// DropGrants removes every unused grant: a grant is valid for the start that follows the install.
func (l *Ledger) DropGrants() error {
	_, err := l.db.ExecContext(context.Background(), `DELETE FROM grants`)
	return err
}
