package ebpf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

// CheckSystemTrusted is rootOwnedSafe judged on the inode f (from OpenConfined) holds, not on a
// path: the file is regular, root-owned and not writable by anyone else, and is right now an entry
// of a directory that is likewise, so only root could have put it there. A rename after the check
// cannot change which inode was judged. Not FUSE. Whether root vouches for the superblock (nosuid,
// user images) is the caller's to check: guard_vouched_devs keys it by device.
func CheckSystemTrusted(f *os.File) error {
	var st unix.Stat_t
	if err := unix.Fstat(int(f.Fd()), &st); err != nil {
		return err
	}
	if st.Mode&unix.S_IFMT != unix.S_IFREG {
		return errors.New("not a regular file")
	}
	if err := rootOwnedStat(&st); err != nil {
		return err
	}
	var fs unix.Statfs_t
	if err := unix.Fstatfs(int(f.Fd()), &fs); err != nil {
		return err
	}
	if fs.Type == unix.FUSE_SUPER_MAGIC {
		return errors.New("on a FUSE filesystem")
	}
	resolved, err := os.Readlink(fmt.Sprintf("/proc/self/fd/%d", f.Fd()))
	if err != nil {
		return err
	}
	return checkSystemParent(resolved, &st)
}

// checkSystemParent: resolved's directory is root-owned and not writable by anyone else, and holds
// the inode st describes right now.
func checkSystemParent(resolved string, st *unix.Stat_t) error {
	if !filepath.IsAbs(resolved) || strings.HasSuffix(resolved, " (deleted)") {
		return fmt.Errorf("%s has no name", resolved)
	}
	dir, err := OpenConfined(filepath.Dir(resolved))
	if err != nil {
		return err
	}
	defer dir.Close()
	var dst unix.Stat_t
	if err := unix.Fstat(int(dir.Fd()), &dst); err != nil {
		return err
	}
	if dst.Mode&unix.S_IFMT != unix.S_IFDIR {
		return errors.New("parent is not a directory")
	}
	if err := rootOwnedStat(&dst); err != nil {
		return fmt.Errorf("parent %w", err)
	}
	var entry unix.Stat_t
	if err := unix.Fstatat(int(dir.Fd()), filepath.Base(resolved), &entry, unix.AT_SYMLINK_NOFOLLOW); err != nil {
		return err
	}
	if entry.Dev != st.Dev || entry.Ino != st.Ino {
		return errors.New("moved while being checked")
	}
	return nil
}

func rootOwnedStat(st *unix.Stat_t) error {
	switch {
	case st.Uid != 0:
		return fmt.Errorf("is not root-owned (uid %d)", st.Uid)
	case st.Mode&unix.S_IWOTH != 0:
		return errors.New("is world-writable")
	case st.Mode&unix.S_IWGRP != 0 && st.Gid != 0:
		return fmt.Errorf("is group-writable by non-root gid %d", st.Gid)
	}
	return nil
}

// ComputeBinaryEntryFile is ComputeBinaryEntry for the inode f (an O_PATH fd) holds, reopened
// through its magic link: the path is never resolved again. path only names the entry.
func ComputeBinaryEntryFile(f *os.File, path string) (BinaryEntry, error) {
	r, err := os.Open(fmt.Sprintf("/proc/self/fd/%d", f.Fd()))
	if err != nil {
		return BinaryEntry{}, fmt.Errorf("reading binary %s: %w", path, err)
	}
	defer r.Close()
	return entryOf(r, path)
}

// StatFile is the (dev, ino) of the inode f holds, dev encoded by KernelDev.
func StatFile(f *os.File) (dev, ino uint64, err error) {
	var st unix.Stat_t
	if err := unix.Fstat(int(f.Fd()), &st); err != nil {
		return 0, 0, &os.PathError{Op: "fstat", Path: f.Name(), Err: err}
	}
	return KernelDev(unix.Major(st.Dev), unix.Minor(st.Dev)), st.Ino, nil
}
