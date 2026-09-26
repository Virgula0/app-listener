package install

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/safeio"
)

// CopyTree recursively copies src into dst, preserving permissions, ownership and mtimes. Symlinks
// are recreated; sockets, FIFOs and device nodes are skipped.
//
// Symlink-safe: the fscrypt migration runs this as root, copying a user's tree back into a
// freshly-created vault while no guard is attached, so the destination's owner (the unprivileged
// user) could otherwise plant symlinks mid-copy to redirect root writes/chowns. Every destination
// step goes through a parent dir fd with *at syscalls: entries are created O_EXCL|O_NOFOLLOW, meta
// is applied on the descriptor, and a directory's own ownership is set only AFTER its subtree is
// fully copied (so the user never owns a destination directory while root is still writing into it).
func CopyTree(src, dst string) error {
	return CopyTreeWithProgress(src, dst, nil)
}

// CopyTreeWithProgress is CopyTree plus onBytes after every regular file (bytes copied so far,
// total to copy); a nil onBytes skips the pre-measurement walk and all reporting.
func CopyTreeWithProgress(src, dst string, onBytes func(copied, total int64)) error {
	var total int64
	if onBytes != nil {
		var err error
		if total, err = treeSize(src); err != nil {
			return fmt.Errorf("measuring %s: %w", src, err)
		}
	}
	var copied int64
	var report func(int64)
	if onBytes != nil {
		report = func(n int64) {
			copied += n
			onBytes(copied, total)
		}
	}
	return copyTreeRoot(src, dst, report)
}

// treeSize sums the sizes of every regular file under root (the pre-walk
// used to compute the copy total).
func treeSize(root string) (int64, error) {
	var total int64
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.Type().IsRegular() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			total += info.Size()
		}
		return nil
	})
	return total, err
}

// testHookBeforeSourceOpen runs after a source entry was classified and before it is opened.
var testHookBeforeSourceOpen = func(string) {}

// copyTreeRoot copies src (a directory the caller may have pre-created at dst, or a lone
// file/symlink) through parent dir fds on both sides. The parents themselves are the caller's
// choice and may be symlinks (a dotfiles-managed ~/.config); everything below them is pinned.
func copyTreeRoot(src, dst string, report func(int64)) error {
	sparent, err := unix.Open(filepath.Dir(src), unix.O_DIRECTORY|unix.O_CLOEXEC, 0)
	if err != nil {
		return fmt.Errorf("opening %s: %w", filepath.Dir(src), err)
	}
	defer unix.Close(sparent)
	dparent, err := unix.Open(filepath.Dir(dst), unix.O_DIRECTORY|unix.O_CLOEXEC, 0)
	if err != nil {
		return fmt.Errorf("opening %s: %w", filepath.Dir(dst), err)
	}
	defer unix.Close(dparent)
	return copyEntry(sparent, filepath.Dir(src), filepath.Base(src), dparent, filepath.Base(dst), 0, report)
}

// copyDirContents copies every child of the open source dir sfd into the open dst dir fd, then
// applies the directory's own mode/owner/mtime LAST — the destination directory is root-owned with
// the source's perms throughout the copy, so its (unprivileged) eventual owner cannot inject entries
// mid-copy.
func copyDirContents(sfd int, srcDir string, dfd int, info fs.FileInfo, depth int, report func(int64)) error {
	names, err := readDirNames(sfd)
	if err != nil {
		return fmt.Errorf("reading %s: %w", srcDir, err)
	}
	for _, name := range names {
		if cerr := copyEntry(sfd, srcDir, name, dfd, name, depth+1, report); cerr != nil {
			return cerr
		}
	}
	if cerr := unix.Fchmod(dfd, uint32(info.Mode().Perm())); cerr != nil {
		return fmt.Errorf("chmod %s: %w", srcDir, cerr)
	}
	if cerr := fchownFromInfo(dfd, info); cerr != nil {
		return cerr
	}
	return futimesFromInfo(dfd, info)
}

// readDirNames lists the open directory sfd without consuming it.
func readDirNames(sfd int) ([]string, error) {
	dup, err := unix.Dup(sfd)
	if err != nil {
		return nil, err
	}
	d := os.NewFile(uintptr(dup), "")
	defer d.Close()
	return d.Readdirnames(-1)
}

// copyEntry copies sparent/sname into dparent/dname. The source is classified from an
// O_PATH|O_NOFOLLOW handle and reopened relative to sparent only if it is still that inode: the
// source tree belongs to the user and no guard is attached, so a name swapped for a symlink (or
// another file) between the two steps must not make root copy content from outside the tree.
func copyEntry(sparent int, srcDir, sname string, dparent int, dname string, depth int, report func(int64)) error {
	srcPath := filepath.Join(srcDir, sname)
	if depth > 64 {
		return fmt.Errorf("directory tree too deep at %s", srcPath)
	}
	pfd, err := unix.Openat(sparent, sname, unix.O_PATH|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if err != nil {
		return fmt.Errorf("opening %s: %w", srcPath, err)
	}
	pinned := os.NewFile(uintptr(pfd), srcPath)
	defer pinned.Close()
	info, err := pinned.Stat()
	if err != nil {
		return fmt.Errorf("stat %s: %w", srcPath, err)
	}

	switch {
	case info.IsDir():
		testHookBeforeSourceOpen(srcPath)
		sfd, oErr := reopenPinned(sparent, sname, info, unix.O_DIRECTORY)
		if oErr != nil {
			return fmt.Errorf("opening %s: %w", srcPath, oErr)
		}
		defer unix.Close(sfd)
		if mkErr := unix.Mkdirat(dparent, dname, uint32(info.Mode().Perm())); mkErr != nil && !errors.Is(mkErr, unix.EEXIST) {
			return fmt.Errorf("creating dir %s: %w", dname, mkErr)
		}
		cfd, oErr := unix.Openat(dparent, dname, unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
		if oErr != nil {
			if errors.Is(oErr, unix.ELOOP) || errors.Is(oErr, unix.ENOTDIR) {
				return fmt.Errorf("%w: destination %s", safeio.ErrSymlink, dname)
			}
			return fmt.Errorf("opening dir %s: %w", dname, oErr)
		}
		defer unix.Close(cfd)
		return copyDirContents(sfd, srcPath, cfd, info, depth, report)
	case info.Mode()&os.ModeSymlink != 0:
		return copySymlinkAt(pfd, srcPath, dparent, dname, info)
	case info.Mode().IsRegular():
		testHookBeforeSourceOpen(srcPath)
		sfd, oErr := reopenPinned(sparent, sname, info, unix.O_RDONLY|unix.O_NONBLOCK)
		if oErr != nil {
			return fmt.Errorf("opening %s: %w", srcPath, oErr)
		}
		in := os.NewFile(uintptr(sfd), srcPath)
		defer in.Close()
		return copyRegularAt(in, dparent, dname, info, report)
	default:
		// Sockets, FIFOs, device nodes: nothing meaningful to copy.
		return nil
	}
}

// reopenPinned opens sparent/name for reading (never through a symlink) and requires it to be the
// inode pinned already.
func reopenPinned(sparent int, name string, pinned fs.FileInfo, flags int) (int, error) {
	fd, err := unix.Openat(sparent, name, flags|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if err != nil {
		return -1, err
	}
	var st unix.Stat_t
	want, ok := pinned.Sys().(*syscall.Stat_t)
	if err := unix.Fstat(fd, &st); err != nil || !ok || st.Dev != want.Dev || st.Ino != want.Ino {
		_ = unix.Close(fd)
		return -1, errors.New("source entry changed while it was being copied")
	}
	return fd, nil
}

func copySymlinkAt(pfd int, srcPath string, dparent int, name string, info fs.FileInfo) error {
	target, err := readlinkFD(pfd)
	if err != nil {
		return fmt.Errorf("reading link %s: %w", srcPath, err)
	}
	if err := unix.Symlinkat(target, dparent, name); err != nil && !errors.Is(err, unix.EEXIST) {
		return fmt.Errorf("creating symlink %s: %w", name, err)
	}
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		if err := unix.Fchownat(dparent, name, int(stat.Uid), int(stat.Gid), unix.AT_SYMLINK_NOFOLLOW); err != nil &&
			!errors.Is(err, unix.EPERM) {
			return fmt.Errorf("chown symlink %s: %w", name, err)
		}
	}
	return nil
}

// readlinkFD reads the target of the symlink an O_PATH|O_NOFOLLOW fd refers to.
func readlinkFD(pfd int) (string, error) {
	buf := make([]byte, unix.PathMax)
	n, err := unix.Readlinkat(pfd, "", buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func copyRegularAt(in *os.File, dparent int, name string, info fs.FileInfo, report func(int64)) error {
	// Refuse to write through a symlink at name (a link planted by the destination's user), but
	// allow overwriting an existing regular file (binary reinstall). O_NOFOLLOW on every open
	// guards the check→open race: if name became a symlink after the Fstatat, the open fails.
	var st unix.Stat_t
	statErr := unix.Fstatat(dparent, name, &st, unix.AT_SYMLINK_NOFOLLOW)
	openFlags := unix.O_WRONLY | unix.O_NOFOLLOW | unix.O_CLOEXEC
	switch {
	case errors.Is(statErr, unix.ENOENT):
		openFlags |= unix.O_CREAT | unix.O_EXCL
	case statErr != nil:
		return fmt.Errorf("checking %s: %w", name, statErr)
	case st.Mode&unix.S_IFMT == unix.S_IFLNK:
		return fmt.Errorf("%w: destination %s", safeio.ErrSymlink, name)
	case st.Mode&unix.S_IFMT == unix.S_IFREG:
		openFlags |= unix.O_TRUNC
	default:
		return fmt.Errorf("refusing to overwrite non-regular destination %s", name)
	}
	fd, err := unix.Openat(dparent, name, openFlags, uint32(info.Mode().Perm()))
	if err != nil {
		return fmt.Errorf("creating %s: %w", name, err)
	}
	out := os.NewFile(uintptr(fd), name)
	if _, cerr := io.Copy(out, in); cerr != nil {
		_ = out.Close()
		return cerr
	}
	if cerr := unix.Fchmod(fd, uint32(info.Mode().Perm())); cerr != nil {
		_ = out.Close()
		return fmt.Errorf("chmod %s: %w", name, cerr)
	}
	if cerr := fchownFromInfo(fd, info); cerr != nil {
		_ = out.Close()
		return cerr
	}
	if cerr := futimesFromInfo(fd, info); cerr != nil {
		_ = out.Close()
		return cerr
	}
	if cerr := out.Close(); cerr != nil {
		return cerr
	}
	if report != nil {
		report(info.Size())
	}
	return nil
}

// fchownFromInfo restores ownership on the descriptor (best effort: EPERM is tolerated for a
// non-root caller, matching the old path-based behavior).
func fchownFromInfo(fd int, info fs.FileInfo) error {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return nil
	}
	if err := unix.Fchown(fd, int(stat.Uid), int(stat.Gid)); err != nil && !errors.Is(err, unix.EPERM) {
		return err
	}
	return nil
}

// futimesFromInfo sets atime=mtime=ModTime on the descriptor itself (AT_EMPTY_PATH).
func futimesFromInfo(fd int, info fs.FileInfo) error {
	ts := unix.NsecToTimespec(info.ModTime().UnixNano())
	return unix.UtimesNanoAt(fd, "", []unix.Timespec{ts, ts}, unix.AT_EMPTY_PATH)
}
