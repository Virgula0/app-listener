package maintainer

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

const maxDepth = 256

// resolveBeneath confines a walk to its dirfd: no symlink, magic link, ".." escape or mount
// crossing on the way. Everything below a lib_dir is written through it.
const resolveBeneath = unix.RESOLVE_BENEATH | unix.RESOLVE_NO_SYMLINKS | unix.RESOLVE_NO_MAGICLINKS |
	unix.RESOLVE_NO_XDEV

func openDir(dir string) (int, error) {
	fd, err := unix.Openat2(unix.AT_FDCWD, dir, &unix.OpenHow{
		Flags:   unix.O_RDONLY | unix.O_DIRECTORY | unix.O_CLOEXEC,
		Resolve: unix.RESOLVE_NO_SYMLINKS | unix.RESOLVE_NO_MAGICLINKS,
	})
	if err != nil {
		return -1, &os.PathError{Op: "open", Path: dir, Err: err}
	}
	return fd, nil
}

func openBeneath(dir int, rel string, flags int, mode uint32) (int, error) {
	fd, err := unix.Openat2(dir, rel, &unix.OpenHow{
		Flags:   uint64(flags | unix.O_CLOEXEC | unix.O_NOFOLLOW), //nolint:gosec // non-negative open flags
		Mode:    uint64(mode),
		Resolve: resolveBeneath,
	})
	if err != nil {
		return -1, &os.PathError{Op: "open", Path: rel, Err: err}
	}
	return fd, nil
}

// removeAll is rm -rf of name in dir without following a symlink anywhere below it.
func removeAll(dir int, name string) error { return removeAllDepth(dir, name, 0) }

func removeAllDepth(dir int, name string, depth int) error {
	if depth > maxDepth {
		return fmt.Errorf("%s: tree deeper than %d", name, maxDepth)
	}
	var st unix.Stat_t
	if err := unix.Fstatat(dir, name, &st, unix.AT_SYMLINK_NOFOLLOW); err != nil {
		if errors.Is(err, unix.ENOENT) {
			return nil
		}
		return err
	}
	if st.Mode&unix.S_IFMT != unix.S_IFDIR {
		return ignoreNotExist(unix.Unlinkat(dir, name, 0))
	}
	fd, err := unix.Openat(dir, name, unix.O_RDONLY|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if err != nil {
		return err
	}
	f := os.NewFile(uintptr(fd), name)
	names, err := f.Readdirnames(-1)
	if err != nil {
		f.Close()
		return err
	}
	for _, n := range names {
		if err := removeAllDepth(fd, n, depth+1); err != nil {
			f.Close()
			return err
		}
	}
	f.Close()
	return ignoreNotExist(unix.Unlinkat(dir, name, unix.AT_REMOVEDIR))
}

func ignoreNotExist(err error) error {
	if errors.Is(err, unix.ENOENT) {
		return nil
	}
	return err
}

// systemTool finds name in the system bin dirs, root-placed only (ebpf.SystemTrusted).
func systemTool(name string) (string, error) {
	for _, dir := range []string{"/usr/bin", "/bin", "/usr/sbin", "/sbin"} {
		if p := filepath.Join(dir, name); ebpf.SystemTrusted(p) {
			return p, nil
		}
	}
	return "", fmt.Errorf("no root-owned %s in the system bin directories", name)
}

type limitedWriter struct {
	w io.Writer
	n int
}

func (l *limitedWriter) Write(p []byte) (int, error) {
	if l.n > 0 {
		k := min(len(p), l.n)
		l.n -= k
		_, _ = l.w.Write(p[:k])
	}
	return len(p), nil
}
