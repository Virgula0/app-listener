package ebpf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/sys/unix"
)

// maxPlacedHops bounds the symlinks OpenSystemPlaced follows, as the kernel's MAXSYMLINKS.
const maxPlacedHops = 40

// OpenSystemPlaced opens path O_PATH when only root could have placed what it names: every
// directory a name is looked up in, through every symlink hop, is root-owned, writable by no one
// else and outside every user home (SetUserHomes); the file it ends on is regular, root-owned, not
// writable by anyone else and not on FUSE. A symlink is judged by the directory holding it. Each name
// is opened O_NOFOLLOW relative to the checked directory's fd, so no rename can redirect the walk
// after a check, and the fd returned holds the judged inode. Whether root vouches for the
// superblock (nosuid, user images) is the caller's to check: guard_vouched_devs keys it by device.
func OpenSystemPlaced(path string) (*os.File, error) {
	fd, st, err := walkPlaced(path)
	if err != nil {
		return nil, err
	}
	f := os.NewFile(uintptr(fd), path)
	if err := systemFileStat(f, st); err != nil {
		f.Close()
		return nil, fmt.Errorf("%s %w", path, err)
	}
	return f, nil
}

func systemFileStat(f *os.File, st *unix.Stat_t) error {
	if st.Mode&unix.S_IFMT != unix.S_IFREG {
		return errors.New("is not a regular file")
	}
	if err := rootOwnedStat(st); err != nil {
		return err
	}
	var fs unix.Statfs_t
	if err := unix.Fstatfs(int(f.Fd()), &fs); err != nil {
		return err
	}
	if fs.Type == unix.FUSE_SUPER_MAGIC {
		return errors.New("is on a FUSE filesystem")
	}
	return nil
}

// walkPlaced resolves path for OpenSystemPlaced, expanding symlinks by hand; it returns the final
// fd and its stat.
func walkPlaced(path string) (int, *unix.Stat_t, error) {
	if !filepath.IsAbs(path) {
		return -1, nil, fmt.Errorf("%s is not an absolute path", path)
	}
	w := &placedWalk{path: path, dir: -1, comps: strings.Split(path, "/")}
	if err := w.restart(); err != nil {
		return -1, nil, err
	}
	for len(w.comps) > 0 {
		c := w.comps[0]
		w.comps = w.comps[1:]
		if c == "" || c == "." {
			continue
		}
		fd, done, err := w.step(c)
		if err != nil {
			unix.Close(w.dir)
			return -1, nil, err
		}
		if done {
			return fd, &w.st, nil
		}
	}
	return w.dir, &w.st, nil
}

// placedWalk is walkPlaced's state: dir is the placed directory names are looked up in, st the
// stat of the last inode reached, comps the names still to resolve.
type placedWalk struct {
	path  string
	dir   int
	st    unix.Stat_t
	comps []string
	hops  int
}

func (w *placedWalk) restart() error {
	if w.dir >= 0 {
		unix.Close(w.dir)
	}
	dir, err := unix.Open("/", unix.O_PATH|unix.O_DIRECTORY|unix.O_CLOEXEC, 0)
	if err != nil {
		w.dir = -1
		return &os.PathError{Op: "open", Path: "/", Err: err}
	}
	w.dir = dir
	return unix.Fstat(dir, &w.st)
}

// step looks c up in w.dir: a directory becomes w.dir, a symlink is expanded into w.comps, and any
// other file ends the walk (done, its fd).
func (w *placedWalk) step(c string) (fd int, done bool, err error) {
	if err = placedDir(w.dir); err != nil {
		return -1, false, fmt.Errorf("%s: %w", w.path, err)
	}
	how := unix.OpenHow{Flags: unix.O_PATH | unix.O_NOFOLLOW | unix.O_CLOEXEC,
		Resolve: unix.RESOLVE_NO_SYMLINKS | unix.RESOLVE_NO_MAGICLINKS}
	next, err := unix.Openat2(w.dir, c, &how)
	if err != nil {
		return -1, false, &os.PathError{Op: "open", Path: w.path, Err: err}
	}
	if err = unix.Fstat(next, &w.st); err != nil {
		unix.Close(next)
		return -1, false, &os.PathError{Op: "fstat", Path: w.path, Err: err}
	}
	switch w.st.Mode & unix.S_IFMT {
	case unix.S_IFLNK:
		err = w.follow(next)
		unix.Close(next)
		return -1, false, err
	case unix.S_IFDIR:
		unix.Close(w.dir)
		w.dir = next
		return -1, false, nil
	}
	if len(w.comps) > 0 {
		unix.Close(next)
		return -1, false, &os.PathError{Op: "open", Path: w.path, Err: unix.ENOTDIR}
	}
	unix.Close(w.dir)
	w.dir = -1
	return next, true, nil
}

// follow queues link's target; an absolute one restarts at "/". It is judged by the directory
// holding the link (already placedDir), so the link's own owner does not matter.
func (w *placedWalk) follow(link int) error {
	w.hops++
	if w.hops > maxPlacedHops {
		return &os.PathError{Op: "open", Path: w.path, Err: unix.ELOOP}
	}
	target, err := readlinkFd(link)
	if err != nil {
		return &os.PathError{Op: "readlink", Path: w.path, Err: err}
	}
	w.comps = append(strings.Split(target, "/"), w.comps...)
	if filepath.IsAbs(target) {
		return w.restart()
	}
	return unix.Fstat(w.dir, &w.st) // a walk ending here ends on w.dir
}

// placedDir: entries of directory fd only root can add, remove or rename (confineBelow's rule),
// and its owner is not a user's FUSE server's say-so.
func placedDir(fd int) error {
	var st unix.Stat_t
	if err := unix.Fstat(fd, &st); err != nil {
		return err
	}
	name, err := os.Readlink(fmt.Sprintf("/proc/self/fd/%d", fd))
	if err != nil {
		return err
	}
	if err := rootOwnedStat(&st); err != nil {
		return fmt.Errorf("%s %w", name, err)
	}
	if inUserHome(name) {
		return fmt.Errorf("%s is in a user home", name)
	}
	// A user's FUSE server reports whatever owner it likes for the directories it serves.
	var fs unix.Statfs_t
	if err := unix.Fstatfs(fd, &fs); err != nil {
		return err
	}
	if fs.Type == unix.FUSE_SUPER_MAGIC {
		return fmt.Errorf("%s is on a FUSE filesystem", name)
	}
	return nil
}

func readlinkFd(fd int) (string, error) {
	buf := make([]byte, unix.PathMax)
	n, err := unix.Readlinkat(fd, "", buf)
	if err != nil {
		return "", err
	}
	if n == len(buf) {
		return "", unix.ENAMETOOLONG
	}
	return string(buf[:n]), nil
}

// maxPlacedDirEntries bounds the entries GlobSystemPlaced reads from one directory.
const maxPlacedDirEntries = 4096

// GlobSystemPlaced is filepath.Glob expanding wildcards only in directories OpenSystemPlaced would
// walk, so no user can add a match or crowd root's out of limit; existing matches, sorted, at most
// limit. Only a hint: each match must still pass OpenSystemPlaced.
func GlobSystemPlaced(pattern string, limit int) []string {
	if _, err := filepath.Match(pattern, ""); err != nil || !filepath.IsAbs(pattern) {
		return nil
	}
	cands := []string{"/"}
	for _, comp := range strings.Split(strings.TrimPrefix(filepath.Clean(pattern), "/"), "/") {
		var next []string
		for _, dir := range cands {
			if strings.ContainsAny(comp, "*?[") {
				next = append(next, placedMatches(dir, comp)...)
			} else {
				next = append(next, filepath.Join(dir, comp))
			}
		}
		cands = next
	}
	var out []string
	for _, m := range cands {
		if _, err := os.Lstat(m); err == nil {
			out = append(out, m)
		}
	}
	sort.Strings(out)
	if len(out) > limit {
		out = out[:limit]
	}
	return out
}

// placedMatches lists dir's entries matching comp when dir is a placed directory.
func placedMatches(dir, comp string) []string {
	fd, st, err := walkPlaced(dir)
	if err != nil {
		return nil
	}
	defer unix.Close(fd)
	if st.Mode&unix.S_IFMT != unix.S_IFDIR || placedDir(fd) != nil {
		return nil
	}
	d, err := os.Open(fmt.Sprintf("/proc/self/fd/%d", fd))
	if err != nil {
		return nil
	}
	defer d.Close()
	names, _ := d.Readdirnames(maxPlacedDirEntries)
	var out []string
	for _, n := range names {
		if ok, _ := filepath.Match(comp, n); ok {
			out = append(out, filepath.Join(dir, n))
		}
	}
	return out
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

// SystemPlacedInode reports whether path is root-placed (OpenSystemPlaced), ends on inode dev:ino
// and lies on a mount root vouches for (mountVouchesOwnership): a file no user could have put
// there, whatever happened while no guard ran.
func SystemPlacedInode(path string, dev, ino uint64) bool {
	f, err := OpenSystemPlaced(path)
	if err != nil {
		return false
	}
	defer f.Close()
	if d, i, serr := StatFile(f); serr != nil || d != dev || i != ino {
		return false
	}
	resolved, err := os.Readlink(fmt.Sprintf("/proc/self/fd/%d", f.Fd()))
	return err == nil && mountVouchesOwnership(resolved) == nil
}
