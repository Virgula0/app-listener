package guard

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// rootHandle pins a guard root: scans go through proc (the fd's magic link), so swapping an
// unguarded ancestor mid-scan can't redirect them into another tree.
type rootHandle struct {
	fd       int
	key      GuardInodeKey
	resolved string // the kernel's d_path for fd: where the configured path resolved
	dir      bool
	proc     string
}

// openRoot resolves path once, following symlinks like StatInode (the anchor's historical flavor).
func openRoot(path string) (*rootHandle, error) {
	fd, err := unix.Open(path, unix.O_PATH|unix.O_CLOEXEC, 0)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: path, Err: err}
	}
	st, err := statFD(fd)
	if err != nil {
		_ = unix.Close(fd)
		return nil, fmt.Errorf("stating guard root %s: %w", path, err)
	}
	proc := procFDDir + strconv.Itoa(fd)
	resolved, err := os.Readlink(proc)
	if err != nil {
		_ = unix.Close(fd)
		return nil, fmt.Errorf("resolving guard root %s: %w", path, err)
	}
	return &rootHandle{fd: fd, key: st.key, resolved: resolved, dir: st.dir, proc: proc}, nil
}

const procFDDir = "/proc/self/fd/"

func (r *rootHandle) close() { _ = unix.Close(r.fd) }

// logical maps a path below r.proc back to the configured root path.
func (r *rootHandle) logical(root, p string) string {
	return root + strings.TrimPrefix(p, r.proc)
}

type fdStat struct {
	key   GuardInodeKey
	dir   bool
	nlink uint32
	mnt   uint64
}

func statFD(fd int) (fdStat, error) {
	var sx unix.Statx_t
	mask := unix.STATX_TYPE | unix.STATX_INO | unix.STATX_NLINK | unix.STATX_MNT_ID
	if err := unix.Statx(fd, "", unix.AT_EMPTY_PATH|unix.AT_STATX_SYNC_AS_STAT, mask, &sx); err != nil {
		return fdStat{}, err
	}
	if sx.Mask&unix.STATX_MNT_ID == 0 {
		return fdStat{}, errors.New("statx: no mount id")
	}
	return fdStat{
		key:   GuardInodeKey{Dev: ebpf.KernelDev(sx.Dev_major, sx.Dev_minor), Ino: sx.Ino},
		dir:   sx.Mode&unix.S_IFMT == unix.S_IFDIR,
		nlink: sx.Nlink,
		mnt:   sx.Mnt_id,
	}, nil
}

const maxChainDepth = 4096

// inodeChain returns key's physical ancestor chain, nearest first and key itself included: what
// the kernel's root_in_chain walks (d_parent), not what path's symlinks spell. Every step is
// fd-relative and key is re-verified at path, so a symlink swapped mid-walk can't redirect it.
// ok=false when the chain can't be known exactly (raced, hard-linked non-dir, bind mount of a
// subdirectory whose hidden parents d_parent would still walk).
func inodeChain(path string, key GuardInodeKey, follow bool) ([]GuardInodeKey, bool) {
	flags := unix.O_PATH | unix.O_CLOEXEC
	// A pinned root's own proc path is a magic link to that fd: NOFOLLOW would open the link.
	if !follow && !isProcFD(path) {
		flags |= unix.O_NOFOLLOW
	}
	fd, err := unix.Open(path, flags, 0)
	if err != nil {
		return nil, false
	}
	st, err := statFD(fd)
	if err != nil || st.key != key {
		_ = unix.Close(fd)
		return nil, false
	}
	var chain []GuardInodeKey
	if !st.dir {
		pfd, pst, ok := physicalParent(fd, st)
		_ = unix.Close(fd)
		if !ok {
			return nil, false
		}
		chain = append(chain, key)
		fd, st = pfd, pst
	}
	defer func() { _ = unix.Close(fd) }()

	for range maxChainDepth {
		chain = append(chain, st.key)
		up, err := unix.Openat(fd, "..", unix.O_PATH|unix.O_DIRECTORY|unix.O_CLOEXEC, 0)
		if err != nil {
			return nil, false
		}
		ust, err := statFD(up)
		if err != nil {
			_ = unix.Close(up)
			return nil, false
		}
		if ust.mnt != st.mnt || ust.key == st.key {
			// fd is a mount root: d_parent ends here only if the mount shows its superblock root.
			_ = unix.Close(up)
			whole, err := mountShowsSbRoot(st.mnt)
			return chain, err == nil && whole
		}
		_ = unix.Close(fd)
		fd, st = up, ust
	}
	return nil, false
}

// isProcFD: path is exactly /proc/self/fd/<n>.
func isProcFD(path string) bool {
	n, ok := strings.CutPrefix(path, procFDDir)
	_, err := strconv.Atoi(n)
	return ok && err == nil
}

// physicalParent opens the directory holding the non-directory fd, verified to contain it.
func physicalParent(fd int, st fdStat) (int, fdStat, bool) {
	if st.nlink != 1 {
		return -1, fdStat{}, false
	}
	p, err := os.Readlink(procFDDir + strconv.Itoa(fd))
	if err != nil || !filepath.IsAbs(p) || strings.HasSuffix(p, " (deleted)") {
		return -1, fdStat{}, false
	}
	pfd, err := unix.Openat2(unix.AT_FDCWD, filepath.Dir(p), &unix.OpenHow{
		Flags:   unix.O_PATH | unix.O_DIRECTORY | unix.O_CLOEXEC,
		Resolve: unix.RESOLVE_NO_SYMLINKS,
	})
	if err != nil {
		return -1, fdStat{}, false
	}
	pst, err := statFD(pfd)
	var ent unix.Stat_t
	if err == nil {
		err = unix.Fstatat(pfd, filepath.Base(p), &ent, unix.AT_SYMLINK_NOFOLLOW)
	}
	entKey := GuardInodeKey{Dev: ebpf.KernelDev(unix.Major(ent.Dev), unix.Minor(ent.Dev)), Ino: ent.Ino}
	if err != nil || entKey != st.key || pst.mnt != st.mnt {
		_ = unix.Close(pfd)
		return -1, fdStat{}, false
	}
	return pfd, pst, true
}

// mountShowsSbRoot reports whether mount mnt (this namespace) is rooted at its superblock's root.
func mountShowsSbRoot(mnt uint64) (bool, error) {
	data, err := os.ReadFile("/proc/self/mountinfo")
	if err != nil {
		return false, err
	}
	return mountinfoShowsSbRoot(data, mnt)
}

func mountinfoShowsSbRoot(mountinfo []byte, mnt uint64) (bool, error) {
	id := strconv.FormatUint(mnt, 10)
	for _, line := range strings.Split(string(mountinfo), "\n") {
		f := strings.Fields(line)
		if len(f) > 3 && f[0] == id {
			return f[3] == "/", nil
		}
	}
	return false, fmt.Errorf("mount %d not in mountinfo", mnt)
}
