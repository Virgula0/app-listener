package ebpf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/sys/unix"
)

var (
	homesMu   sync.RWMutex
	userHomes []string
)

// SetUserHomes sets the trees OpenConfined confines whatever their ownership: catalog globs expand
// there, and the guards must stop the owner's (or root's) malware too. "/" confines everything.
func SetUserHomes(homes []string) {
	var out []string
	for _, h := range homes {
		out = append(out, filepath.Clean(h))
		if r, err := filepath.EvalSymlinks(h); err == nil {
			out = append(out, r) // the physical path OpenConfined compares against (/home -> var/home)
		}
	}
	homesMu.Lock()
	userHomes = out
	homesMu.Unlock()
}

func inUserHome(dir string) bool {
	homesMu.RLock()
	defer homesMu.RUnlock()
	for _, h := range userHomes {
		if h == "/" || dir == h || strings.HasPrefix(dir, h+"/") {
			return true
		}
	}
	return false
}

// OpenConfined opens path O_PATH one component at a time, each relative to the previous fd. A
// directory component whose parent is in a user home (SetUserHomes) or writable by a non-root user
// resolves RESOLVE_BENEATH that parent: a symlink there can be re-pointed by malware, and a stored
// whitelist path must not follow it out of the tree it names (a catalog match's reserved glob
// root), the rule homeMatchConfined applies at refresh. The final component resolves freely, as
// there: its name is reserved or root-owned. Holding an fd per step, no rename between steps can
// redirect the walk.
func OpenConfined(path string) (*os.File, error) {
	if !filepath.IsAbs(path) {
		return nil, fmt.Errorf("%s is not an absolute path", path)
	}
	dir, err := unix.Open("/", unix.O_PATH|unix.O_DIRECTORY|unix.O_CLOEXEC, 0)
	if err != nil {
		return nil, &os.PathError{Op: "open", Path: "/", Err: err}
	}
	clean := filepath.Clean(path)
	if clean == "/" {
		return os.NewFile(uintptr(dir), clean), nil
	}
	comps := strings.Split(clean[1:], "/")
	for i, c := range comps {
		how := unix.OpenHow{Flags: unix.O_PATH | unix.O_CLOEXEC}
		if i < len(comps)-1 {
			how.Flags |= unix.O_DIRECTORY
			confine, cerr := confineBelow(dir)
			if cerr != nil {
				unix.Close(dir)
				return nil, fmt.Errorf("%s: %w", path, cerr)
			}
			if confine {
				how.Resolve = unix.RESOLVE_BENEATH | unix.RESOLVE_NO_MAGICLINKS
			}
		}
		next, oerr := unix.Openat2(dir, c, &how)
		unix.Close(dir)
		if oerr != nil {
			sub := "/" + strings.Join(comps[:i+1], "/")
			if how.Resolve != 0 && (errors.Is(oerr, unix.EXDEV) || errors.Is(oerr, unix.ELOOP)) {
				return nil, fmt.Errorf("%s: %s is a symlink leaving its parent directory, which malware "+
					"could re-point", path, sub)
			}
			return nil, &os.PathError{Op: "open", Path: sub, Err: oerr}
		}
		dir = next
	}
	return os.NewFile(uintptr(dir), clean), nil
}

// confineBelow reports whether entries of directory fd resolve RESOLVE_BENEATH it: fd is in a user
// home, or (mirroring rootOwnedInode) anyone but root may add, remove or rename entries in it.
func confineBelow(fd int) (bool, error) {
	var st unix.Stat_t
	if err := unix.Fstat(fd, &st); err != nil {
		return false, err
	}
	if st.Uid != 0 || st.Mode&unix.S_IWOTH != 0 || (st.Mode&unix.S_IWGRP != 0 && st.Gid != 0) {
		return true, nil
	}
	dir, err := os.Readlink(fmt.Sprintf("/proc/self/fd/%d", fd))
	if err != nil {
		return false, err
	}
	return inUserHome(dir), nil
}

// StatConfined is StatInode resolved by OpenConfined.
func StatConfined(path string) (dev, ino uint64, err error) {
	f, err := OpenConfined(path)
	if err != nil {
		return 0, 0, err
	}
	defer f.Close()
	var st unix.Stat_t
	if err := unix.Fstat(int(f.Fd()), &st); err != nil {
		return 0, 0, &os.PathError{Op: "fstat", Path: path, Err: err}
	}
	return KernelDev(unix.Major(st.Dev), unix.Minor(st.Dev)), st.Ino, nil
}

// ResolveConfined is filepath.EvalSymlinks resolved by OpenConfined.
func ResolveConfined(path string) (string, error) {
	f, err := OpenConfined(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	return os.Readlink(fmt.Sprintf("/proc/self/fd/%d", f.Fd()))
}
