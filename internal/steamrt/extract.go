package steamrt

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/unix"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

const (
	maxEntries = 200_000
	maxBytes   = 4 << 30
	maxDepth   = 256
	// xzTimeout bounds one decompression (the scout runtime is ~65 MB packed).
	xzTimeout = 10 * time.Minute
	// nobody runs xz: a same-user process can't ptrace it into feeding other bytes to the tree.
	nobody = 65534
)

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

func openRegular(dir int, rel string) (*os.File, error) {
	fd, err := openBeneath(dir, rel, unix.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	f := os.NewFile(uintptr(fd), rel)
	if fi, err := f.Stat(); err != nil || !fi.Mode().IsRegular() {
		f.Close()
		return nil, fmt.Errorf("%s is not a regular file", rel)
	}
	return f, nil
}

func readSmall(dir int, rel string) ([]byte, error) {
	f, err := openRegular(dir, rel)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readAllBounded(f, maxChecksum)
}

func readAllBounded(r io.Reader, limit int64) ([]byte, error) {
	b, err := io.ReadAll(io.LimitReader(r, limit+1))
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > limit {
		return nil, fmt.Errorf("more than %d bytes", limit)
	}
	return b, nil
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

// extract unpacks the xz-compressed tar archive into root, which ends up holding r.Name.
func (r Runtime) extract(root int, archive *os.File) error {
	xz, err := systemTool("xz")
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), xzTimeout)
	defer cancel()
	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, xz, "-dc")
	cmd.Stdin = archive
	cmd.Stderr = &limitedWriter{w: &stderr, n: 4096}
	cmd.Env = []string{}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{Uid: nobody, Gid: nobody},
		Pdeathsig:  syscall.SIGKILL,
	}
	stream, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err = cmd.Start(); err != nil {
		return fmt.Errorf("starting xz: %w", err)
	}
	e := extractor{root: root, top: r.Name, uid: int(r.UID), gid: int(r.GID)}
	err = e.all(tar.NewReader(stream))
	if err == nil {
		_, err = io.Copy(io.Discard, io.LimitReader(stream, 1<<20)) // trailing padding: let xz finish
	}
	stream.Close()
	if werr := cmd.Wait(); err == nil && werr != nil {
		err = fmt.Errorf("xz: %w: %s", werr, strings.TrimSpace(stderr.String()))
	}
	if err == nil && !e.sawTop {
		err = fmt.Errorf("the archive holds no %s/", r.Name)
	}
	return err
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

// extractor writes an archive's entries below root: only below its top directory, never through a
// symlink (resolveBeneath), no device or fifo, no setuid/setgid/sticky bit, owned by uid:gid.
type extractor struct {
	root     int
	top      string
	uid, gid int
	entries  int
	bytes    int64
	sawTop   bool
}

func (e *extractor) all(tr *tar.Reader) error {
	for {
		h, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if err := e.entry(h, tr); err != nil {
			return fmt.Errorf("%q: %w", h.Name, err)
		}
	}
}

func (e *extractor) entry(h *tar.Header, body io.Reader) error {
	if h.Typeflag == tar.TypeXGlobalHeader {
		return nil
	}
	e.entries++
	if e.entries > maxEntries {
		return fmt.Errorf("more than %d entries", maxEntries)
	}
	name, err := e.clean(h.Name)
	if err != nil {
		return err
	}
	e.sawTop = true
	parent, base, done, err := e.parent(name)
	if err != nil {
		return err
	}
	defer done()
	mode := uint32(h.Mode) & 0o777 //nolint:gosec // permission bits only
	switch h.Typeflag {
	case tar.TypeDir:
		return e.dir(parent, base, mode)
	case tar.TypeReg:
		return e.file(parent, base, mode, h, body)
	case tar.TypeSymlink:
		if h.Linkname == "" || len(h.Linkname) >= unix.PathMax {
			return errors.New("bad symlink target")
		}
		if err := unix.Symlinkat(h.Linkname, parent, base); err != nil {
			return err
		}
		return e.finish(parent, base, h)
	case tar.TypeLink:
		return e.hardlink(parent, base, h)
	default:
		return fmt.Errorf("unsupported entry type %q", h.Typeflag)
	}
}

// clean accepts a relative, already-clean name below e.top.
func (e *extractor) clean(raw string) (string, error) {
	name := strings.TrimSuffix(raw, "/")
	if name == "" || strings.HasPrefix(name, "/") || path.Clean(name) != name {
		return "", errors.New("not a clean relative name")
	}
	if name != e.top && !strings.HasPrefix(name, e.top+"/") {
		return "", fmt.Errorf("outside %s/", e.top)
	}
	for c := range strings.SplitSeq(name, "/") {
		if c == ".." || c == "." {
			return "", errors.New("not a clean relative name")
		}
	}
	return name, nil
}

// parent opens name's directory, creating missing ones as tar does (Valve's archive lists no entry
// for its top directory).
func (e *extractor) parent(name string) (dirfd int, base string, done func(), err error) {
	dir, base := path.Split(name)
	if dir == "" {
		return e.root, base, func() {}, nil
	}
	dir = strings.TrimSuffix(dir, "/")
	fd, err := openBeneath(e.root, dir, unix.O_RDONLY|unix.O_DIRECTORY, 0)
	if errors.Is(err, unix.ENOENT) {
		fd, err = e.mkdirAll(dir)
	}
	if err != nil {
		return -1, "", nil, err
	}
	return fd, base, func() { unix.Close(fd) }, nil
}

func (e *extractor) mkdirAll(dir string) (int, error) {
	fd, err := unix.Dup(e.root)
	if err != nil {
		return -1, err
	}
	for c := range strings.SplitSeq(dir, "/") {
		if err := e.dir(fd, c, 0o755); err != nil {
			unix.Close(fd)
			return -1, err
		}
		next, err := openBeneath(fd, c, unix.O_RDONLY|unix.O_DIRECTORY, 0)
		unix.Close(fd)
		if err != nil {
			return -1, err
		}
		fd = next
	}
	return fd, nil
}

func (e *extractor) dir(parent int, base string, mode uint32) error {
	// Owner rwx while extracting: a read-only directory in the archive must still receive entries.
	if err := unix.Mkdirat(parent, base, mode|0o700); err != nil {
		if !errors.Is(err, unix.EEXIST) {
			return err
		}
		var st unix.Stat_t
		if err := unix.Fstatat(parent, base, &st, unix.AT_SYMLINK_NOFOLLOW); err != nil {
			return err
		}
		if st.Mode&unix.S_IFMT != unix.S_IFDIR {
			return errors.New("exists and is not a directory")
		}
	}
	return unix.Fchownat(parent, base, e.uid, e.gid, unix.AT_SYMLINK_NOFOLLOW)
}

func (e *extractor) file(parent int, base string, mode uint32, h *tar.Header, body io.Reader) error {
	if h.Size < 0 || e.bytes+h.Size > maxBytes {
		return fmt.Errorf("archive larger than %d bytes", int64(maxBytes))
	}
	e.bytes += h.Size
	fd, err := unix.Openat(parent, base, unix.O_WRONLY|unix.O_CREAT|unix.O_EXCL|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0o600)
	if err != nil {
		return err
	}
	f := os.NewFile(uintptr(fd), base)
	defer f.Close()
	if _, err := io.CopyN(f, body, h.Size); err != nil {
		return err
	}
	if err := f.Chown(e.uid, e.gid); err != nil {
		return err
	}
	if err := f.Chmod(os.FileMode(mode)); err != nil {
		return err
	}
	return e.times(parent, base, h)
}

func (e *extractor) hardlink(parent int, base string, h *tar.Header) error {
	target, err := e.clean(h.Linkname)
	if err != nil {
		return fmt.Errorf("hard link target: %w", err)
	}
	tparent, tbase, done, err := e.parent(target)
	if err != nil {
		return err
	}
	defer done()
	var st unix.Stat_t
	if err := unix.Fstatat(tparent, tbase, &st, unix.AT_SYMLINK_NOFOLLOW); err != nil {
		return err
	}
	if st.Mode&unix.S_IFMT != unix.S_IFREG {
		return errors.New("hard link to a non-regular file")
	}
	return unix.Linkat(tparent, tbase, parent, base, 0)
}

func (e *extractor) finish(parent int, base string, h *tar.Header) error {
	if err := unix.Fchownat(parent, base, e.uid, e.gid, unix.AT_SYMLINK_NOFOLLOW); err != nil {
		return err
	}
	return e.times(parent, base, h)
}

func (e *extractor) times(parent int, base string, h *tar.Header) error {
	ts := []unix.Timespec{unix.NsecToTimespec(h.ModTime.UnixNano()), unix.NsecToTimespec(h.ModTime.UnixNano())}
	return unix.UtimesNanoAt(parent, base, ts, unix.AT_SYMLINK_NOFOLLOW)
}
