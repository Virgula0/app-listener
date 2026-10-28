// Package bintrust is `app-listener trust`'s core, shared by the CLI and the GUI's edit-protected
// backend: the binaries the daemon refused (binledger pending rows) as they are now, and recording
// the confirmed ones.
package bintrust

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"

	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/binledger"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// Row is one binary to confirm: Pending as refused (SHA256 = the refused content) or as given by
// hand (SHA256 = Current), plus what the file holds now.
type Row struct {
	binledger.Pending
	Current [32]byte
	ReadErr error
	// Stat is the file's metadata as hashed; nil when it couldn't be read.
	Stat *unix.Stat_t
}

// OpenLedger opens the installed binary ledger.
func OpenLedger() (*binledger.Ledger, error) {
	if _, err := os.Stat(filepath.Dir(binledger.DefaultPath)); err != nil {
		return nil, fmt.Errorf("no %s: is the daemon installed? (%w)", filepath.Dir(binledger.DefaultPath), err)
	}
	return binledger.Open(binledger.DefaultPath)
}

// Pending lists the ledger's refused binaries as they are now.
func Pending(l *binledger.Ledger) ([]Row, error) {
	pend, err := l.ListPending()
	if err != nil {
		return nil, err
	}
	out := make([]Row, 0, len(pend))
	for _, p := range pend {
		r := Row{Pending: p}
		r.Current, r.Stat, r.ReadErr = Inspect(p.Line)
		out = append(out, r)
	}
	return out, nil
}

// Given describes lines to confirm as they are now (`trust add`), resources naming their origin.
func Given(lines []string, resources string) []Row {
	out := make([]Row, 0, len(lines))
	for _, line := range lines {
		r := Row{Pending: binledger.Pending{Line: line, Resolved: line, Resources: resources, SeenAt: time.Now()}}
		r.Current, r.Stat, r.ReadErr = Inspect(line)
		r.SHA256 = r.Current
		out = append(out, r)
	}
	return out
}

// Inspect hashes the regular file at path. O_NONBLOCK: a FIFO planted at a whitelist line must not
// hang the reader.
func Inspect(path string) ([32]byte, *unix.Stat_t, error) {
	fd, err := unix.Open(path, unix.O_RDONLY|unix.O_NONBLOCK|unix.O_CLOEXEC, 0)
	if err != nil {
		return [32]byte{}, nil, &os.PathError{Op: "open", Path: path, Err: err}
	}
	f := os.NewFile(uintptr(fd), path)
	defer f.Close()
	var st unix.Stat_t
	if err = unix.Fstat(fd, &st); err != nil {
		return [32]byte{}, nil, &os.PathError{Op: "fstat", Path: path, Err: err}
	}
	if st.Mode&unix.S_IFMT != unix.S_IFREG {
		return [32]byte{}, nil, fmt.Errorf("%s is not a regular file", path)
	}
	e, err := ebpf.ComputeBinaryEntryFile(f, path)
	if err != nil {
		return [32]byte{}, nil, err
	}
	return e.Hash, &st, nil
}

// Skip is why r can't be confirmed now, "" when it can: only the content that was refused (or
// reviewed) is confirmed; a file that changed again since is judged afresh by the next reload.
func (r *Row) Skip() string {
	switch {
	case r.ReadErr != nil:
		return fmt.Sprintf("cannot read it now (%v)", r.ReadErr)
	case r.Current != r.SHA256:
		return fmt.Sprintf("changed again since the daemon refused it (now %x): reload the daemon to have it "+
			"judged again", r.Current)
	}
	return ""
}

// ErrChanged: the line no longer holds the content that was reviewed.
var ErrChanged = errors.New("its content changed since it was reviewed")

var errNotRefused = errors.New("the daemon no longer refuses this content (confirmed meanwhile, or refused " +
	"again with other content)")

// ConfirmPending records line as sum only while the daemon's refusal still names sum and the file
// still holds it. Nothing is opened unless the ledger refuses exactly that: line and sum may come
// from a caller, and hashing a path it names would be a content oracle on guarded files.
func ConfirmPending(l *binledger.Ledger, line string, sum [32]byte) error {
	if err := refuses(l, line, sum); err != nil {
		return err
	}
	if err := stillHolds(line, sum); err != nil {
		return err
	}
	ok, err := l.ConfirmPending(line, sum)
	switch {
	case err != nil:
		return err
	case !ok:
		return errNotRefused
	}
	return nil
}

func refuses(l *binledger.Ledger, line string, sum [32]byte) error {
	pend, err := l.ListPending()
	if err != nil {
		return err
	}
	if !slices.ContainsFunc(pend, func(p binledger.Pending) bool { return p.Line == line && p.SHA256 == sum }) {
		return errNotRefused
	}
	return nil
}

// ConfirmGiven records line as sum (`trust add`) while the file still holds sum.
func ConfirmGiven(l *binledger.Ledger, line string, sum [32]byte) error {
	if err := stillHolds(line, sum); err != nil {
		return err
	}
	return l.Record(line, sum, binledger.SourceConfirmed)
}

func stillHolds(line string, sum [32]byte) error {
	cur, _, err := Inspect(line)
	switch {
	case err != nil:
		return err
	case cur != sum:
		return ErrChanged
	}
	return nil
}
