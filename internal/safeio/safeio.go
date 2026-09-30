// Package safeio provides symlink-safe filesystem primitives for code that runs as root over paths
// an unprivileged user controls (the interactive editor, the fscrypt migration copy, the installer's
// per-user file drops). os.Open/os.Create/os.Chmod/os.Chown and os.MkdirAll all follow symlinks and
// resolve every path component, so a user who plants a symlink — at the target, at a temp-file name,
// or at a parent directory — can redirect a root write/chown to an arbitrary file. These helpers
// pin each step with O_NOFOLLOW and *at syscalls and validate the resulting descriptor, never the
// path, closing the check-then-use window.
package safeio

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// ErrSymlink reports that a component which had to be a real file or directory was a symlink.
var ErrSymlink = errors.New("refusing to operate through a symlink")

// DescendNoFollow opens each component in comps relative to dirFD with O_NOFOLLOW, optionally
// creating missing directories (mode dirPerm), and returns the deepest directory fd. A symlink (or
// non-directory) component is refused. The returned fd is always the deepest one held (or dirFD on
// early error), so a single deferred Close covers every path. dirFD is consumed: on success it is
// closed and the returned fd replaces it; on error dirFD is returned unchanged for the caller to
// close.
func DescendNoFollow(dirFD int, comps []string, createDirs bool, dirPerm os.FileMode) (int, error) {
	for _, comp := range comps {
		if comp == "" || comp == "." || comp == ".." {
			return dirFD, fmt.Errorf("illegal path component %q", comp)
		}
		if createDirs {
			if mkErr := unix.Mkdirat(dirFD, comp, uint32(dirPerm)); mkErr != nil && !errors.Is(mkErr, unix.EEXIST) {
				return dirFD, fmt.Errorf("creating %q: %w", comp, mkErr)
			}
		}
		next, err := unix.Openat(dirFD, comp, unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
		if err != nil {
			if errors.Is(err, unix.ELOOP) || errors.Is(err, unix.ENOTDIR) {
				return dirFD, fmt.Errorf("%w: component %q", ErrSymlink, comp)
			}
			return dirFD, fmt.Errorf("descending into %q: %w", comp, err)
		}
		_ = unix.Close(dirFD)
		dirFD = next
	}
	return dirFD, nil
}

// DescendCreateOwned is DescendNoFollow that also chowns each directory it CREATES to uid/gid
// (existing components are left untouched); uid/gid < 0 skips the chown. Missing directories are
// made 0755. Use it to build a path under a user's home as root without following a symlink the user
// planted at any parent component.
func DescendCreateOwned(dirFD int, comps []string, uid, gid int) (int, error) {
	for _, comp := range comps {
		if comp == "" || comp == "." || comp == ".." {
			return dirFD, fmt.Errorf("illegal path component %q", comp)
		}
		created := false
		if mkErr := unix.Mkdirat(dirFD, comp, 0o755); mkErr == nil {
			created = true
		} else if !errors.Is(mkErr, unix.EEXIST) {
			return dirFD, fmt.Errorf("creating %q: %w", comp, mkErr)
		}
		next, err := unix.Openat(dirFD, comp, unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
		if err != nil {
			if errors.Is(err, unix.ELOOP) || errors.Is(err, unix.ENOTDIR) {
				return dirFD, fmt.Errorf("%w: component %q", ErrSymlink, comp)
			}
			return dirFD, fmt.Errorf("descending into %q: %w", comp, err)
		}
		_ = unix.Close(dirFD)
		dirFD = next
		if created && uid >= 0 && gid >= 0 {
			if cErr := unix.Fchown(dirFD, uid, gid); cErr != nil {
				return dirFD, fmt.Errorf("chown %q: %w", comp, cErr)
			}
		}
	}
	return dirFD, nil
}

// RefuseExistingSymlink errors (ErrSymlink) when name under dirFD already exists as a symlink. A
// missing name is fine.
func RefuseExistingSymlink(dirFD int, name string) error {
	var st unix.Stat_t
	switch err := unix.Fstatat(dirFD, name, &st, unix.AT_SYMLINK_NOFOLLOW); {
	case err == nil && st.Mode&unix.S_IFMT == unix.S_IFLNK:
		return fmt.Errorf("%w: %s", ErrSymlink, name)
	case err != nil && !errors.Is(err, unix.ENOENT):
		return fmt.Errorf("checking %s: %w", name, err)
	}
	return nil
}

// CreateExclNoFollow creates name under dirFD with O_CREAT|O_EXCL|O_NOFOLLOW (O_RDWR, so the caller
// may read it back) and returns an *os.File. It fails if name already exists (including as a
// symlink). uid/gid < 0 skips the fchown.
func CreateExclNoFollow(dirFD int, name string, perm os.FileMode, uid, gid int) (*os.File, error) {
	fd, err := unix.Openat(dirFD, name,
		unix.O_RDWR|unix.O_CREAT|unix.O_EXCL|unix.O_NOFOLLOW|unix.O_CLOEXEC, uint32(perm))
	if err != nil {
		return nil, err
	}
	f := os.NewFile(uintptr(fd), name)
	if uid >= 0 && gid >= 0 {
		if cerr := f.Chown(uid, gid); cerr != nil {
			_ = f.Close()
			_ = unix.Unlinkat(dirFD, name, 0)
			return nil, cerr
		}
	}
	return f, nil
}

// AtomicWriteAt writes data to a temp file under dirFD (O_EXCL|O_NOFOLLOW), applies perm and, when
// uid/gid >= 0, ownership on the descriptor, then renames it over name. tmpName must differ from
// name and should be dot-prefixed. Never follows a symlink at either the temp or the final name.
func AtomicWriteAt(dirFD int, name, tmpName string, data []byte, perm os.FileMode, uid, gid int) error {
	_ = unix.Unlinkat(dirFD, tmpName, 0) // clear a stale temp from a crashed run
	fd, err := unix.Openat(dirFD, tmpName,
		unix.O_WRONLY|unix.O_CREAT|unix.O_EXCL|unix.O_NOFOLLOW|unix.O_CLOEXEC, uint32(perm))
	if err != nil {
		return fmt.Errorf("creating temp file %s: %w", tmpName, err)
	}
	f := os.NewFile(uintptr(fd), tmpName)
	cleanup := func() { _ = f.Close(); _ = unix.Unlinkat(dirFD, tmpName, 0) }
	if _, werr := f.Write(data); werr != nil {
		cleanup()
		return fmt.Errorf("writing %s: %w", tmpName, werr)
	}
	if serr := f.Sync(); serr != nil {
		cleanup()
		return fmt.Errorf("syncing %s: %w", tmpName, serr)
	}
	if cerr := f.Chmod(perm); cerr != nil {
		cleanup()
		return fmt.Errorf("chmod %s: %w", tmpName, cerr)
	}
	if uid >= 0 && gid >= 0 {
		if cerr := f.Chown(uid, gid); cerr != nil {
			cleanup()
			return fmt.Errorf("chown %s: %w", tmpName, cerr)
		}
	}
	if cerr := f.Close(); cerr != nil {
		_ = unix.Unlinkat(dirFD, tmpName, 0)
		return fmt.Errorf("closing %s: %w", tmpName, cerr)
	}
	if rerr := unix.Renameat(dirFD, tmpName, dirFD, name); rerr != nil {
		_ = unix.Unlinkat(dirFD, tmpName, 0)
		return fmt.Errorf("installing %s: %w", name, rerr)
	}
	return nil
}

// OpenRegularNoFollow opens path with O_NOFOLLOW added and verifies via fstat that the resulting
// descriptor is a regular file (never a symlink target, device, fifo, …). It closes the split-second
// gap that a path-based Lstat-then-open leaves: the descriptor, not the name, is validated.
func OpenRegularNoFollow(path string, flags int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(path, flags|unix.O_NOFOLLOW, perm)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	if !info.Mode().IsRegular() {
		_ = f.Close()
		return nil, fmt.Errorf("%s is not a regular file", path)
	}
	return f, nil
}
