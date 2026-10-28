package maintainer

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"syscall"

	"golang.org/x/sys/unix"
)

// Apply writes the plan into each lib_dir inside grant(dir, fn), owned by the run's user. Each file
// is written from what was validated: a record from memory, any other copied from the stage and
// re-hashed on the way, so nothing that reached the stage after validation is committed. A failure
// stops there; what was applied before it stays, all of it validated.
func (p *Plan) Apply(grant func(dir string, fn func() error) error) error {
	for i := range p.libs {
		lp := &p.libs[i]
		if err := grant(lp.dir, func() error { return p.applyLib(lp) }); err != nil {
			return fmt.Errorf("%s: %w", lp.dir, err)
		}
	}
	return nil
}

func (p *Plan) applyLib(lp *libPlan) error {
	root, err := openDir(lp.dir)
	if err != nil {
		return err
	}
	defer unix.Close(root)
	for i := range lp.ops {
		if err := p.apply(root, &lp.ops[i]); err != nil {
			return fmt.Errorf("%s: %w", lp.ops[i].rel, err)
		}
	}
	return nil
}

func (p *Plan) apply(root int, o *op) error {
	if o.rel == "." {
		return clearExcept(root, o.names)
	}
	dir, name := path.Split(o.rel)
	parent := root
	if dir != "" {
		fd, err := openBeneath(root, dir[:len(dir)-1], unix.O_RDONLY|unix.O_DIRECTORY, 0)
		if err != nil {
			return err
		}
		defer unix.Close(fd)
		parent = fd
	}
	switch o.kind {
	case opDelete:
		return removeAll(parent, name)
	case opDir:
		return p.applyDir(parent, name, o)
	case opLink:
		return p.applyLink(parent, name, o)
	default:
		return p.applyFile(parent, name, o)
	}
}

func (p *Plan) applyDir(parent int, name string, o *op) error {
	var st unix.Stat_t
	statErr := unix.Fstatat(parent, name, &st, unix.AT_SYMLINK_NOFOLLOW)
	switch {
	case statErr == nil && st.Mode&unix.S_IFMT == unix.S_IFDIR:
	case statErr == nil || errors.Is(statErr, unix.ENOENT):
		if err := removeAll(parent, name); err != nil {
			return err
		}
		if err := unix.Mkdirat(parent, name, 0o700); err != nil {
			return err
		}
	default:
		return statErr
	}
	fd, err := openBeneath(parent, name, unix.O_RDONLY|unix.O_DIRECTORY, 0)
	if err != nil {
		return err
	}
	defer unix.Close(fd)
	if o.opaque {
		if err := clearExcept(fd, o.names); err != nil {
			return err
		}
	}
	return p.own(fd, o.mode)
}

func (p *Plan) applyLink(parent int, name string, o *op) error {
	if cur, err := readlinkat(parent, name); err == nil && cur == o.target {
		return unix.Fchownat(parent, name, int(p.uid), int(p.gid), unix.AT_SYMLINK_NOFOLLOW)
	}
	if err := removeAll(parent, name); err != nil {
		return err
	}
	if err := unix.Symlinkat(o.target, parent, name); err != nil {
		return err
	}
	return unix.Fchownat(parent, name, int(p.uid), int(p.gid), unix.AT_SYMLINK_NOFOLLOW)
}

func (p *Plan) applyFile(parent int, name string, o *op) error {
	if fd, err := openBeneath(parent, name, unix.O_RDONLY|unix.O_NONBLOCK, 0); err == nil {
		f := os.NewFile(uintptr(fd), name)
		same := holds(f, o)
		if same {
			err = p.own(fd, o.mode)
		}
		f.Close()
		if same {
			return err
		}
	}
	tmp, f, err := createTemp(parent)
	if err != nil {
		return err
	}
	committed := false
	defer func() {
		f.Close()
		if !committed {
			_ = unix.Unlinkat(parent, tmp, 0)
		}
	}()
	if err := p.write(f, o); err != nil {
		return err
	}
	if err := p.own(int(f.Fd()), o.mode); err != nil {
		return err
	}
	if err := removeDir(parent, name); err != nil {
		return err
	}
	if err := unix.Renameat(parent, tmp, parent, name); err != nil {
		return err
	}
	committed = true
	return nil
}

// write copies o's validated content into f: a record from memory, else from the stage, re-hashed.
func (p *Plan) write(f *os.File, o *op) error {
	if o.data != nil || o.size == 0 {
		_, err := f.Write(o.data)
		return err
	}
	src, err := os.OpenFile(o.src, os.O_RDONLY|syscall.O_NOFOLLOW, 0)
	if err != nil {
		return err
	}
	defer src.Close()
	h := sha256.New()
	n, err := io.Copy(io.MultiWriter(f, h), io.LimitReader(src, o.size+1))
	if err != nil {
		return err
	}
	if n != o.size || !bytes.Equal(h.Sum(nil), o.sum[:]) {
		return errors.New("its staged copy changed after validation")
	}
	return nil
}

func (p *Plan) own(fd int, mode uint32) error {
	var st unix.Stat_t
	if err := unix.Fstat(fd, &st); err != nil {
		return err
	}
	if st.Uid != p.uid || st.Gid != p.gid {
		if err := unix.Fchown(fd, int(p.uid), int(p.gid)); err != nil {
			return err
		}
	}
	if st.Mode&0o7777 != mode {
		return unix.Fchmod(fd, mode)
	}
	return nil
}

// holds reports whether f is a regular file with o's validated content.
func holds(f *os.File, o *op) bool {
	fi, err := f.Stat()
	if err != nil || !fi.Mode().IsRegular() || fi.Size() != o.size {
		return false
	}
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return false
	}
	return bytes.Equal(h.Sum(nil), o.sum[:])
}

func createTemp(parent int) (string, *os.File, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", nil, err
	}
	name := ".app-listener-" + hex.EncodeToString(b[:])
	fd, err := unix.Openat(parent, name, unix.O_WRONLY|unix.O_CREAT|unix.O_EXCL|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0o600)
	if err != nil {
		return "", nil, err
	}
	return name, os.NewFile(uintptr(fd), name), nil
}

// removeDir removes name when it is a directory: a rename can't replace one with a file.
func removeDir(parent int, name string) error {
	var st unix.Stat_t
	if err := unix.Fstatat(parent, name, &st, unix.AT_SYMLINK_NOFOLLOW); err != nil || st.Mode&unix.S_IFMT != unix.S_IFDIR {
		return nil //nolint:nilerr // missing or not a directory: the rename handles it
	}
	return removeAll(parent, name)
}

// clearExcept removes every entry of dir not in keep.
func clearExcept(dir int, keep map[string]bool) error {
	fd, err := unix.Openat(dir, ".", unix.O_RDONLY|unix.O_DIRECTORY|unix.O_CLOEXEC, 0)
	if err != nil {
		return err
	}
	f := os.NewFile(uintptr(fd), ".")
	names, err := f.Readdirnames(-1)
	f.Close()
	if err != nil {
		return err
	}
	for _, n := range names {
		if !keep[n] {
			if err := removeAll(dir, n); err != nil {
				return err
			}
		}
	}
	return nil
}

func readlinkat(dir int, name string) (string, error) {
	buf := make([]byte, unix.PathMax)
	n, err := unix.Readlinkat(dir, name, buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}
