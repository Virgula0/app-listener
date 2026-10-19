package recorder

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/Virgula0/app-listener/internal/binledger"
	"github.com/Virgula0/app-listener/internal/fscrypt"
)

// ErrStale: daemon.conf changed between Fetch and Put.
var ErrStale = errors.New("daemon.conf changed meanwhile")

// Store is the running daemon's configuration, reached over its control socket.
type Store interface {
	Fetch() ([]byte, error)
	// Put writes next only while daemon.conf still holds base (ErrStale otherwise); the daemon
	// reloads it and puts base back when the reload fails.
	Put(base, next []byte) error
	// Guarded lists the paths the running daemon guards.
	Guarded() ([]string, error)
}

// Outcome is a successful Apply.
type Outcome struct {
	Config []byte
	// Backups are the plaintext migration copies Apply left next to each encrypted target. They sit
	// outside every guarded tree: offer to delete them (DeleteBackup).
	Backups []string
	Notes   []string
}

// ProgressFunc reports a long step (an encryption's copy) as a fraction in [0, 1].
type ProgressFunc func(step string, fraction float64)

type vaultOps interface {
	FilesystemPrereqs(path string) ([]fscrypt.Prereq, error)
	EncryptWithProgress(path string, onBytes func(copied, total int64)) error
	RestoreBackup(path string) error
}

type ledgerOps interface {
	RecordIfAbsent(line string, sum [32]byte, src binledger.Source) (bool, error)
	Forget(line string, sum [32]byte, src binledger.Source) error
	Lookup(line string) ([32]byte, bool, error)
}

type applier struct {
	store       Store
	vault       vaultOps
	ledger      ledgerOps
	systemSetup func() error
	progress    ProgressFunc
}

// maxStaleRetries bounds re-rendering over a daemon.conf that keeps changing under the apply.
const maxStaleRetries = 3

// Apply performs p on the running daemon. Order keeps the daemon's invariant (never readable without
// a live guard) and fails closed:
//  1. render and validate against the current daemon.conf, before anything is touched;
//  2. migrate each Encrypt target in place (backup first): it ends locked, so it stays unreadable
//     until the daemon's reload has attached its guard and unlocks it;
//  3. pin each approved binary that isn't root-placed to the content the recording hashed (only a
//     line the ledger has no verdict for: a confirmed line keeps its hash), so the reload admits
//     exactly what was reviewed and refuses a file swapped since;
//  4. put the new daemon.conf (compare-and-swap, re-rendered on a concurrent change).
//
// Any failure undoes 2 and 3: the backups are restored and the pins forgotten.
func Apply(store Store, p *Plan, progress ProgressFunc) (*Outcome, error) {
	l, err := binledger.Open(binledger.DefaultPath)
	if err != nil {
		return nil, fmt.Errorf("opening the binary ledger: %w", err)
	}
	defer l.Close()
	a := &applier{store: store, vault: fscrypt.New(), ledger: l, systemSetup: fscrypt.EnsureSystemSetup,
		progress: progress}
	return a.apply(p)
}

func (a *applier) apply(p *Plan) (*Outcome, error) {
	base, err := a.store.Fetch()
	if err != nil {
		return nil, err
	}
	text, err := p.Render(base)
	if err != nil {
		return nil, err
	}
	if cerr := a.checkEncryptable(p.Encrypt); cerr != nil {
		return nil, cerr
	}
	encrypted, err := a.encrypt(p.Encrypt)
	if err != nil {
		return nil, a.rollback(err, encrypted, nil)
	}
	pinned, notes, err := a.pin(p.Binaries)
	if err != nil {
		return nil, a.rollback(err, encrypted, pinned)
	}
	applied, err := a.put(p, base, text)
	if err != nil {
		return nil, a.rollback(err, encrypted, pinned)
	}
	out := &Outcome{Config: applied, Notes: notes}
	for _, e := range encrypted {
		out.Backups = append(out.Backups, e+fscrypt.BackupSuffix)
	}
	out.Notes = append(out.Notes, a.unguarded(p)...)
	return out, nil
}

// checkEncryptable refuses before any migration when a target's filesystem isn't ready: the
// commands it needs (tune2fs, fscrypt setup) are the installer's to offer, not a live apply's.
func (a *applier) checkEncryptable(paths []string) error {
	for _, path := range paths {
		prereqs, err := a.vault.FilesystemPrereqs(path)
		if err != nil {
			return fmt.Errorf("%s can't be encrypted: %w", path, err)
		}
		if len(prereqs) > 0 {
			cmds := make([]string, len(prereqs))
			for i, pr := range prereqs {
				cmds[i] = pr.Command()
			}
			return fmt.Errorf("%s can't be encrypted until its filesystem is prepared (%s): run %s, or "+
				"record it without encryption", path, prereqs[0].Reason, strings.Join(cmds, " && "))
		}
	}
	if len(paths) > 0 {
		if err := a.systemSetup(); err != nil {
			return fmt.Errorf("setting up fscrypt: %w", err)
		}
	}
	return nil
}

func (a *applier) encrypt(paths []string) ([]string, error) {
	var done []string
	for _, path := range paths {
		err := a.vault.EncryptWithProgress(path, func(copied, total int64) {
			if a.progress != nil && total > 0 {
				a.progress("Encrypting "+path, float64(copied)/float64(total))
			}
		})
		if err != nil {
			return done, fmt.Errorf("encrypting %s: %w", path, err)
		}
		done = append(done, path)
	}
	return done, nil
}

func (a *applier) pin(bins []PlannedBinary) ([]PlannedBinary, []string, error) {
	var pinned []PlannedBinary
	var notes []string
	for _, b := range bins {
		if b.SystemPlaced {
			continue
		}
		added, err := a.ledger.RecordIfAbsent(b.Line, b.SHA256, binledger.SourceRecorded)
		if err != nil {
			return pinned, nil, fmt.Errorf("recording %s in the binary ledger: %w", b.Line, err)
		}
		if added {
			pinned = append(pinned, b)
			continue
		}
		if sum, ok, lerr := a.ledger.Lookup(b.Line); lerr == nil && ok && sum != b.SHA256 {
			notes = append(notes, fmt.Sprintf("%s: the ledger confirms other content for this line, so the "+
				"daemon admits it only once `app-listener trust` confirms the current file", b.Line))
		}
	}
	return pinned, notes, nil
}

func (a *applier) put(p *Plan, base []byte, text string) ([]byte, error) {
	for attempt := 0; ; attempt++ {
		err := a.store.Put(base, []byte(text))
		if err == nil {
			return []byte(text), nil
		}
		if !errors.Is(err, ErrStale) {
			// A reload outlasting the reply still applied: the file then holds exactly text.
			if cur, ferr := a.store.Fetch(); ferr == nil && bytes.Equal(cur, []byte(text)) {
				return cur, nil
			}
			return nil, err
		}
		if attempt+1 >= maxStaleRetries {
			return nil, fmt.Errorf("daemon.conf keeps changing while it is applied: %w", err)
		}
		if base, err = a.store.Fetch(); err != nil {
			return nil, err
		}
		if text, err = p.Render(base); err != nil {
			return nil, err
		}
	}
}

// rollback undoes encrypted migrations and ledger pins after cause, newest first.
func (a *applier) rollback(cause error, encrypted []string, pinned []PlannedBinary) error {
	var errs []error
	for i := len(encrypted) - 1; i >= 0; i-- {
		if err := a.vault.RestoreBackup(encrypted[i]); err != nil {
			errs = append(errs, fmt.Errorf("restoring %s from %s%s: %w — restore it with `app-listener install "+
				"install backups restore`", encrypted[i], encrypted[i], fscrypt.BackupSuffix, err))
		}
	}
	for _, b := range pinned {
		if err := a.ledger.Forget(b.Line, b.SHA256, binledger.SourceRecorded); err != nil {
			errs = append(errs, fmt.Errorf("forgetting the ledger entry for %s: %w", b.Line, err))
		}
	}
	if len(errs) == 0 {
		return fmt.Errorf("%w — nothing was changed", cause)
	}
	return errors.Join(append([]error{cause}, errs...)...)
}

// unguarded notes every new or extended section the running daemon doesn't report as guarded after the reload.
func (a *applier) unguarded(p *Plan) []string {
	guarded, err := a.store.Guarded()
	if err != nil {
		return []string{"could not confirm the new guards with the daemon: " + err.Error()}
	}
	paths := make([]string, 0, len(p.Sections)+len(p.Extensions))
	for _, s := range p.Sections {
		paths = append(paths, s.Path)
	}
	for _, e := range p.Extensions {
		paths = append(paths, e.Path)
	}
	var out []string
	for _, path := range paths {
		if !slices.Contains(guarded, path) {
			out = append(out, "WARNING: the daemon does not report "+path+" as guarded: check `journalctl -u "+
				"app-listener-daemon`")
		}
	}
	return out
}

// DeleteBackup removes one of Outcome.Backups (the plaintext copy of a now-encrypted target).
func DeleteBackup(o *Outcome, backup string) error {
	if !slices.Contains(o.Backups, backup) || !strings.HasSuffix(backup, fscrypt.BackupSuffix) {
		return fmt.Errorf("%s is not a backup this apply made", backup)
	}
	if err := os.RemoveAll(backup); err != nil {
		return fmt.Errorf("removing %s: %w", backup, err)
	}
	return nil
}
