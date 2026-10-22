// Package steamrt keeps Steam's scout runtime (ubuntu12_32/steam-runtime) current inside a
// read-only lib_dir: the job steam.sh's unpack_runtime and the runtime's setup.sh do with general
// tools (mkdir, tar, ln) that the guard rightly refuses, since malware runs the same tools. The
// daemon unpacks only an archive the trust guard proves the Steam client wrote, and commits
// setup.sh's library pins only after validating them (pins.go).
package steamrt

import (
	"bytes"
	"crypto/md5" //nolint:gosec // steam.sh's own integrity check; provenance is what admits the archive
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

const (
	archiveExt  = ".tar.xz"
	checksumExt = archiveExt + ".checksum"
	// maxChecksum bounds the checksum file: one md5sum line.
	maxChecksum = 4096
)

// Runtime is one Steam runtime: Dir/Name the unpacked tree, Dir/Name.tar.xz the archive the client
// downloads beside it. Dir is a read-only lib_dir, so every write happens inside a Grant.
type Runtime struct {
	// Dir is the lib_dir as the running config names it (canonical: no symlink is followed in it).
	Dir  string
	Name string
	// UID/GID own what the daemon writes: steam.sh would have written it as the Steam user.
	UID, GID uint32
	// Home is the Steam user's home, setup.sh's HOME.
	Home string
}

// Grant runs fn while the daemon may write r.Dir (usecase.DaemonUseCase.WithRuntimeAccess).
type Grant func(fn func() error) error

// Provenance reports whether the file f holds was created by one of the lib_dir's updaters and
// never written by another exe since (guard.TrustGuard.CreatedByUpdater).
type Provenance func(f *os.File) bool

// ErrProvenance: the archive or its checksum was not proven to be the Steam client's.
var ErrProvenance = errors.New("not proven to be written by the runtime's updater")

// Pending reports an archive whose checksum file differs from the unpacked tree's: what
// steam.sh's unpack_runtime would extract.
func (r Runtime) Pending() (bool, error) {
	dir, err := openDir(r.Dir)
	if err != nil {
		return false, err
	}
	defer unix.Close(dir)
	want, err := readSmall(dir, r.Name+checksumExt)
	if errors.Is(err, unix.ENOENT) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	var st unix.Stat_t
	if err = unix.Fstatat(dir, r.Name+archiveExt, &st, unix.AT_SYMLINK_NOFOLLOW); err != nil {
		if errors.Is(err, unix.ENOENT) {
			return false, nil
		}
		return false, err
	}
	have, err := readSmall(dir, r.Name+"/checksum")
	if err != nil && !errors.Is(err, unix.ENOENT) {
		return false, err
	}
	return !bytes.Equal(want, have), nil
}

// Unpack does what unpack_runtime does once the checksums differ, provided prove vouches for both
// files the client wrote: verify the md5, extract into Name.tmp, move the old tree to Name.old and
// the new one in place, copy the checksum. The archive is judged, hashed and decompressed from one
// open file, so a rename after the check can't slip another in. Every write runs inside grant; a
// failure before the swap leaves the old tree in place.
func (r Runtime) Unpack(prove Provenance, grant Grant) error {
	dir, err := openDir(r.Dir)
	if err != nil {
		return err
	}
	defer unix.Close(dir)
	archive, err := openRegular(dir, r.Name+archiveExt)
	if err != nil {
		return err
	}
	defer archive.Close()
	sumFile, err := openRegular(dir, r.Name+checksumExt)
	if err != nil {
		return err
	}
	defer sumFile.Close()
	if !prove(archive) || !prove(sumFile) {
		return ErrProvenance
	}
	want, err := readAllBounded(sumFile, maxChecksum)
	if err != nil {
		return err
	}
	if err := verifyMD5(archive, r.Name+archiveExt, want); err != nil {
		return err
	}
	return grant(func() error { return r.install(dir, archive, want) })
}

// verifyMD5 is unpack_runtime's check: the checksum file holds md5sum's output for the archive.
func verifyMD5(archive *os.File, name string, want []byte) error {
	if _, err := archive.Seek(0, io.SeekStart); err != nil {
		return err
	}
	h := md5.New() //nolint:gosec // see the import
	if _, err := io.Copy(h, archive); err != nil {
		return fmt.Errorf("hashing %s: %w", name, err)
	}
	if _, err := archive.Seek(0, io.SeekStart); err != nil {
		return err
	}
	got := fmt.Sprintf("%x  %s", h.Sum(nil), name)
	if got != strings.TrimRight(string(want), "\n") {
		return fmt.Errorf("%s does not match its checksum file (got %q)", name, got)
	}
	return nil
}

func (r Runtime) install(dir int, archive *os.File, checksum []byte) error {
	tmp := r.Name + ".tmp"
	if err := removeAll(dir, tmp); err != nil {
		return fmt.Errorf("clearing %s: %w", tmp, err)
	}
	if err := r.mkdirOwned(dir, tmp, 0o755); err != nil {
		return err
	}
	tmpfd, err := openBeneath(dir, tmp, unix.O_RDONLY|unix.O_DIRECTORY, 0)
	if err != nil {
		return err
	}
	defer unix.Close(tmpfd)
	if err := r.extract(tmpfd, archive); err != nil {
		if cerr := removeAll(dir, tmp); cerr != nil {
			err = errors.Join(err, cerr)
		}
		return fmt.Errorf("unpacking %s%s: %w", r.Name, archiveExt, err)
	}
	return r.swap(dir, tmpfd, checksum)
}

func (r Runtime) swap(dir, tmpfd int, checksum []byte) error {
	old := r.Name + ".old"
	var st unix.Stat_t
	hadTree := unix.Fstatat(dir, r.Name, &st, unix.AT_SYMLINK_NOFOLLOW) == nil
	if hadTree {
		if err := removeAll(dir, old); err != nil {
			return fmt.Errorf("clearing %s: %w", old, err)
		}
		if err := unix.Renameat(dir, r.Name, dir, old); err != nil {
			return fmt.Errorf("moving %s aside: %w", r.Name, err)
		}
	}
	if err := unix.Renameat(tmpfd, r.Name, dir, r.Name); err != nil {
		if hadTree {
			if rerr := unix.Renameat(dir, old, dir, r.Name); rerr != nil {
				err = errors.Join(err, fmt.Errorf("restoring %s: %w", r.Name, rerr))
			}
		}
		return fmt.Errorf("moving the new %s in place: %w", r.Name, err)
	}
	if err := removeAll(dir, r.Name+".tmp"); err != nil {
		return err
	}
	fd, err := openBeneath(dir, r.Name+"/checksum", unix.O_WRONLY|unix.O_CREAT|unix.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("writing %s/checksum: %w", r.Name, err)
	}
	f := os.NewFile(uintptr(fd), "checksum")
	defer f.Close()
	if _, err := f.Write(checksum); err != nil {
		return err
	}
	return f.Chown(int(r.UID), int(r.GID))
}

func (r Runtime) mkdirOwned(dir int, name string, mode uint32) error {
	if err := unix.Mkdirat(dir, name, mode); err != nil && !errors.Is(err, unix.EEXIST) {
		return fmt.Errorf("creating %s: %w", name, err)
	}
	var st unix.Stat_t
	if err := unix.Fstatat(dir, name, &st, unix.AT_SYMLINK_NOFOLLOW); err != nil {
		return err
	}
	if st.Mode&unix.S_IFMT != unix.S_IFDIR {
		return fmt.Errorf("%s is not a directory", name)
	}
	return unix.Fchownat(dir, name, int(r.UID), int(r.GID), unix.AT_SYMLINK_NOFOLLOW)
}
