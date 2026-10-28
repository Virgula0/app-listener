package ebpf

import (
	"errors"
	"fmt"
	"io"
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
	return placedFileStat(int(f.Fd()), st)
}

func placedFileStat(fd int, st *unix.Stat_t) error {
	if st.Mode&unix.S_IFMT != unix.S_IFREG {
		return errors.New("is not a regular file")
	}
	if err := rootOwnedStat(st); err != nil {
		return err
	}
	var fs unix.Statfs_t
	if err := unix.Fstatfs(fd, &fs); err != nil {
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
	w := &placedWalk{path: path, dir: -1}
	fd, err := w.run()
	if err != nil {
		return -1, nil, err
	}
	return fd, &w.st, nil
}

// placedWalk is walkPlaced's state: dir is the placed directory names are looked up in, st the
// stat of the last inode reached, comps the names still to resolve. missing is the physical path
// of the name a failed lookup did not find, every symlink hop before it resolved.
type placedWalk struct {
	path    string
	dir     int
	st      unix.Stat_t
	comps   []string
	hops    int
	missing string
}

func (w *placedWalk) run() (int, error) {
	if !filepath.IsAbs(w.path) {
		return -1, fmt.Errorf("%s is not an absolute path", w.path)
	}
	w.comps = strings.Split(w.path, "/")
	if err := w.restart(); err != nil {
		return -1, err
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
			return -1, err
		}
		if done {
			return fd, nil
		}
	}
	return w.dir, nil
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
		if errors.Is(err, unix.ENOENT) {
			if dir, e := os.Readlink(fmt.Sprintf("/proc/self/fd/%d", w.dir)); e == nil {
				w.missing = filepath.Join(dir, c)
			}
		}
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
	switch fs.Type {
	case unix.FUSE_SUPER_MAGIC:
		return fmt.Errorf("%s is on a FUSE filesystem", name)
	case unix.PROC_SUPER_MAGIC:
		// /proc/self (and /dev/fd through it) names the reader's own directories, not ours.
		return fmt.Errorf("%s is on procfs", name)
	}
	return nil
}

// RootPlacedName reports whether only root decides what the clean absolute path names, now and
// later, as walkPlaced walks it (every symlink hop judged). An existing target is a root-owned file
// no one else may write, or a directory only root fills all the way down (placedTree), on a mount
// root vouches for (VouchedDevs): what a symlink to it reaches below is never judged again. A
// missing one is looked up in a placed directory and is not a name root creates for a user or
// where a user's media is mounted (userAssignable, userMountable).
func RootPlacedName(path string) error {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return fmt.Errorf("%q is not a clean absolute path", path)
	}
	if userMountable(path) {
		return fmt.Errorf("%s is where a user's own filesystem may be mounted", path)
	}
	w := &placedWalk{path: path, dir: -1}
	fd, err := w.run()
	switch {
	case err == nil:
		defer unix.Close(fd)
		return placedTarget(fd, &w.st, path)
	case !errors.Is(err, unix.ENOENT):
		return err
	case w.missing == "":
		return fmt.Errorf("%s: where the lookup failed is unknown: %w", path, err)
	case userAssignable(w.missing):
		return fmt.Errorf("%s: %s is a name root creates for a user", path, w.missing)
	case userMountable(w.missing):
		return fmt.Errorf("%s: %s is where a user's own filesystem may be mounted", path, w.missing)
	}
	return nil
}

// placedTarget judges RootPlacedName's existing target fd. A tree is vouched at its top only:
// placedEntry never crosses into a mount below.
func placedTarget(fd int, st *unix.Stat_t, path string) error {
	j := &treeJudge{}
	var err error
	if j.devs, err = hostVouchedDevs(); err != nil {
		return err
	}
	if err := vouchedFd(fd, j.devs); err != nil {
		return fmt.Errorf("%s %w", path, err)
	}
	if st.Mode&unix.S_IFMT == unix.S_IFDIR {
		return placedTree(fd, j)
	}
	if err := placedFileStat(fd, st); err != nil {
		return fmt.Errorf("%s %w", path, err)
	}
	return nil
}

// userAssignable: a missing name root creates for a user later: a home beside the users' homes
// (useradd) or a /run/user/<uid> (logind).
func userAssignable(p string) bool {
	dir := filepath.Dir(p)
	if dir == "/home" || dir == "/run/user" {
		return true
	}
	homesMu.RLock()
	defer homesMu.RUnlock()
	for _, h := range userHomes {
		if filepath.Dir(h) == dir {
			return true
		}
	}
	return false
}

// userMountable: p is at or below where udisks2 mounts a user's media (/media, /run/media) or fstab
// "user" mounts conventionally go (/mnt). A filesystem mounted there later presents any owner its
// author wrote, over whatever the name holds now.
func userMountable(p string) bool {
	for _, root := range []string{"/media", "/run/media", "/mnt"} {
		if _, ok := PathUnder(p, root); ok {
			return true
		}
	}
	return false
}

// treeJudge is placedTree's state: n counts the entries judged, at most maxPlacedDirEntries (a tree
// too large to judge is refused); devs are the superblocks root vouches for.
type treeJudge struct {
	n    int
	devs map[uint64]struct{}
}

// placedTree: directory fd and everything below it only root decides. A symlink below must end, as
// OpenSystemPlaced walks it, on a placed regular file on a vouched mount; mounts below are not
// crossed.
func placedTree(fd int, j *treeJudge) error {
	if err := placedDir(fd); err != nil {
		return err
	}
	name, err := os.Readlink(fmt.Sprintf("/proc/self/fd/%d", fd))
	if err != nil {
		return err
	}
	d, err := os.Open(fmt.Sprintf("/proc/self/fd/%d", fd))
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(maxPlacedDirEntries + 1 - j.n)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	for _, c := range names {
		j.n++
		if j.n > maxPlacedDirEntries {
			return fmt.Errorf("%s holds over %d entries: too many to judge", name, maxPlacedDirEntries)
		}
		if err := placedEntry(fd, name, c, j); err != nil {
			return err
		}
	}
	return nil
}

func placedEntry(dir int, dirName, c string, j *treeJudge) error {
	how := unix.OpenHow{Flags: unix.O_PATH | unix.O_NOFOLLOW | unix.O_CLOEXEC,
		Resolve: unix.RESOLVE_NO_SYMLINKS | unix.RESOLVE_NO_MAGICLINKS | unix.RESOLVE_NO_XDEV}
	fd, err := unix.Openat2(dir, c, &how)
	if err != nil {
		return &os.PathError{Op: "open", Path: filepath.Join(dirName, c), Err: err}
	}
	defer unix.Close(fd)
	var st unix.Stat_t
	if err := unix.Fstat(fd, &st); err != nil {
		return err
	}
	p := filepath.Join(dirName, c)
	switch st.Mode & unix.S_IFMT {
	case unix.S_IFDIR:
		return placedTree(fd, j)
	case unix.S_IFLNK:
		target, err := readlinkFd(fd)
		if err != nil {
			return err
		}
		if !filepath.IsAbs(target) {
			target = filepath.Join(dirName, target) // dirName is physical: the lexical ".." is the kernel's
		}
		f, err := OpenSystemPlaced(target)
		if err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}
		defer f.Close()
		if err := vouchedFd(int(f.Fd()), j.devs); err != nil {
			return fmt.Errorf("%s: %s %w", p, target, err)
		}
		return nil
	}
	if err := placedFileStat(fd, &st); err != nil {
		return fmt.Errorf("%s %w", p, err)
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
	for comp := range strings.SplitSeq(strings.TrimPrefix(filepath.Clean(pattern), "/"), "/") {
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
