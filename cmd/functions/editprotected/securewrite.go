package editprotected

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

// writeWithin writes data to dest (mode 0600) such that dest resolves to a
// location strictly inside root without traversing a single symlink below
// root. Missing intermediate directories are created (0700). The write is
// atomic (temp file in the target directory + rename).
//
// This is the symlink-safe writer for the non-interactive --put path:
// os.WriteFile / os.MkdirAll follow symlinks, so a symlink anywhere in a
// guarded tree (a dotfile manager's ~/.ssh/config, a planted link) would
// redirect the write to an arbitrary path outside every guard — where the
// daemon imposes no constraint at all. Every path component below root is
// opened O_NOFOLLOW; a symlink is refused, not followed. The interactive
// editor already enforces the same (internal/tui writeFileKeepMeta).
func writeWithin(root, dest string, data []byte) (err error) {
	rel, relErr := filepath.Rel(root, dest)
	if relErr != nil {
		return fmt.Errorf("resolving %s against %s: %w", dest, root, relErr)
	}
	parts := strings.Split(filepath.Clean(rel), string(os.PathSeparator))
	if len(parts) == 0 || parts[0] == "." || parts[0] == ".." {
		return fmt.Errorf("%s does not resolve to a path inside %s", dest, root)
	}
	name := parts[len(parts)-1]
	if name == "" || name == "." || name == ".." {
		return fmt.Errorf("illegal final path component in %s", dest)
	}

	// Anchor at root. root is the operator-supplied, daemon-validated watch
	// path, so following a symlink for root itself is fine (the guard is on
	// the real inode); every component BELOW it must not be a symlink.
	dirFD, openErr := unix.Open(root, unix.O_DIRECTORY|unix.O_CLOEXEC, 0)
	if openErr != nil {
		return fmt.Errorf("opening %s: %w", root, openErr)
	}
	defer func() { _ = unix.Close(dirFD) }()

	dirFD, err = descendNoFollow(dirFD, parts[:len(parts)-1], root, dest)
	if err != nil {
		return err
	}
	if symErr := refuseExistingSymlink(dirFD, name, dest); symErr != nil {
		return symErr
	}
	return atomicWriteAt(dirFD, name, dest, data)
}

// refuseExistingSymlink errors when name already exists as a symlink under
// dirFD (renameat would replace it, but silently swapping an operator's
// symlink for a real file is its own surprise — the interactive editor
// refuses too). A missing entry is fine.
func refuseExistingSymlink(dirFD int, name, dest string) error {
	var st unix.Stat_t
	switch err := unix.Fstatat(dirFD, name, &st, unix.AT_SYMLINK_NOFOLLOW); {
	case err == nil && st.Mode&unix.S_IFMT == unix.S_IFLNK:
		return fmt.Errorf("refusing to write through symlink %s", dest)
	case err != nil && !errors.Is(err, unix.ENOENT):
		return fmt.Errorf("checking %s: %w", dest, err)
	}
	return nil
}

// atomicWriteAt writes data to a temp file under dirFD (O_EXCL, O_NOFOLLOW)
// and renames it over name.
func atomicWriteAt(dirFD int, name, dest string, data []byte) error {
	tmp := "." + name + ".app_listener.put"
	_ = unix.Unlinkat(dirFD, tmp, 0) // clear a stale temp from a crashed run
	tmpFD, err := unix.Openat(dirFD, tmp,
		unix.O_WRONLY|unix.O_CREAT|unix.O_EXCL|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0o600)
	if err != nil {
		return fmt.Errorf("creating a temp file for %s: %w", dest, err)
	}
	if err := writeAndSync(tmpFD, data); err != nil {
		_ = unix.Unlinkat(dirFD, tmp, 0)
		return fmt.Errorf("writing %s: %w", dest, err)
	}
	if err := unix.Renameat(dirFD, tmp, dirFD, name); err != nil {
		_ = unix.Unlinkat(dirFD, tmp, 0)
		return fmt.Errorf("installing %s: %w", dest, err)
	}
	return nil
}

// descendNoFollow opens each component in comps relative to dirFD with
// O_NOFOLLOW, creating missing directories (0700), and returns the fd of the
// last directory. A symlink component is refused. The returned fd is always
// the deepest one successfully held (or the input on an early error), so the
// caller's single deferred Close covers both paths.
func descendNoFollow(dirFD int, comps []string, root, dest string) (int, error) {
	for _, comp := range comps {
		if comp == "" || comp == "." || comp == ".." {
			return dirFD, fmt.Errorf("illegal path component %q in %s", comp, dest)
		}
		if mkErr := unix.Mkdirat(dirFD, comp, 0o700); mkErr != nil && !errors.Is(mkErr, unix.EEXIST) {
			return dirFD, fmt.Errorf("creating %s under %s: %w", comp, root, mkErr)
		}
		next, descErr := unix.Openat(dirFD, comp, unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
		if descErr != nil {
			if errors.Is(descErr, unix.ELOOP) || errors.Is(descErr, unix.ENOTDIR) {
				return dirFD, fmt.Errorf("refusing to write %s: component %q is a symlink or not a directory", dest, comp)
			}
			return dirFD, fmt.Errorf("descending into %s: %w", comp, descErr)
		}
		_ = unix.Close(dirFD)
		dirFD = next
	}
	return dirFD, nil
}

func writeAndSync(fd int, data []byte) error {
	f := os.NewFile(uintptr(fd), "")
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}
