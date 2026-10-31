package maintainer

import (
	"archive/tar"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	maxArchiveBytes   = 16 << 30
	maxArchiveMembers = 1_000_000
	decompressTimeout = 10 * time.Minute
	// nobody runs external decompressors: a parser bug then runs as nobody, and a same-user process
	// can't ptrace it into reporting other members.
	nobody = 65534
)

// sources is the content a lib_dir already holds: its regular files, and the regular members of
// the archives among them. A run may commit only that content. Trusting it trusts no more than the
// lib_dir itself: only its writers write it while the daemon runs.
type sources struct {
	dir     string
	dev     uint64 // the lib_dir's filesystem: content on any other (a mount below it) vouches nothing
	indexed bool
	bySize  map[int64][]string
	sums    map[string][32]byte
	members map[[32]byte]bool
}

func newSources(dir string) *sources {
	return &sources{dir: dir, sums: map[string][32]byte{}}
}

// has reports whether the lib_dir holds content of size and sum.
func (s *sources) has(size int64, sum [32]byte) (bool, error) {
	if err := s.index(); err != nil {
		return false, err
	}
	for _, p := range s.bySize[size] {
		h, err := s.sum(p)
		if err != nil {
			continue // gone or unreadable: not a source
		}
		if h == sum {
			return true, nil
		}
	}
	if s.members == nil {
		s.members = map[[32]byte]bool{}
		for _, list := range s.bySize {
			for _, p := range list {
				if err := s.addArchive(p); err != nil {
					return false, fmt.Errorf("reading archive %s: %w", p, err)
				}
			}
		}
	}
	return s.members[sum], nil
}

func (s *sources) index() error {
	if s.indexed {
		return nil
	}
	dev, err := libDirDev(s.dir)
	if err != nil {
		return err
	}
	s.dev = dev
	s.bySize = map[int64][]string{}
	err = filepath.WalkDir(s.dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && !d.Type().IsRegular() {
			return nil
		}
		fi, err := d.Info()
		if err != nil {
			return nil //nolint:nilerr // removed since the walk listed it: not a source
		}
		// A user's FUSE mount below a read-only lib_dir (user_allow_other) serves any bytes.
		if st, ok := fi.Sys().(*syscall.Stat_t); !ok || st.Dev != s.dev {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if d.Type().IsRegular() {
			s.bySize[fi.Size()] = append(s.bySize[fi.Size()], p)
		}
		return nil
	})
	s.indexed = err == nil
	return err
}

// libDirDev is dir's device; a FUSE one (a user's mount over an ancestor) is refused.
func libDirDev(dir string) (uint64, error) {
	f, err := os.Open(dir)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	var fsSt syscall.Statfs_t
	if err := syscall.Fstatfs(int(f.Fd()), &fsSt); err != nil {
		return 0, err
	}
	if fsSt.Type == fuseSuperMagic {
		return 0, fmt.Errorf("%s is on a FUSE filesystem: its content vouches for nothing", dir)
	}
	var st syscall.Stat_t
	if err := syscall.Fstat(int(f.Fd()), &st); err != nil {
		return 0, err
	}
	return st.Dev, nil
}

const fuseSuperMagic = 0x65735546

// open opens p, a source only while it is on the lib_dir's filesystem (checked on the fd: a mount
// made after the walk can't slip in).
func (s *sources) open(p string) (*os.File, *syscall.Stat_t, error) {
	f, err := os.OpenFile(p, os.O_RDONLY|syscall.O_NOFOLLOW, 0)
	if err != nil {
		return nil, nil, err
	}
	var st syscall.Stat_t
	if err := syscall.Fstat(int(f.Fd()), &st); err != nil {
		f.Close()
		return nil, nil, err
	}
	if st.Dev != s.dev {
		f.Close()
		return nil, nil, fmt.Errorf("%s is on another filesystem than the lib_dir", p)
	}
	return f, &st, nil
}

func (s *sources) sum(p string) ([32]byte, error) {
	if h, ok := s.sums[p]; ok {
		return h, nil
	}
	f, _, err := s.open(p)
	if err != nil {
		return [32]byte{}, err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return [32]byte{}, err
	}
	var out [32]byte
	copy(out[:], h.Sum(nil))
	s.sums[p] = out
	return out, nil
}

// archiveCache keeps each archive's member sums across runs, by file identity.
var archiveCache = struct {
	sync.Mutex
	m map[fileID]map[[32]byte]bool
}{m: map[fileID]map[[32]byte]bool{}}

type fileID struct {
	dev, ino     uint64
	size         int64
	mtime, ctime syscall.Timespec
}

// addArchive adds p's member sums when p is a tar archive, plain or compressed; anything else adds
// nothing.
func (s *sources) addArchive(p string) error {
	f, st, openErr := s.open(p)
	if openErr != nil {
		return nil // gone, unreadable or mounted over: not a source
	}
	defer f.Close()
	id := fileID{dev: st.Dev, ino: st.Ino, size: st.Size, mtime: st.Mtim, ctime: st.Ctim}
	archiveCache.Lock()
	cached, ok := archiveCache.m[id]
	archiveCache.Unlock()
	if !ok {
		var err error
		if cached, err = archiveMembers(f); err != nil {
			return err
		}
		archiveCache.Lock()
		archiveCache.m[id] = cached
		archiveCache.Unlock()
	}
	for sum := range cached {
		s.members[sum] = true
	}
	return nil
}

// archiveMembers returns the sums of f's regular members, or none when f is no tar archive.
func archiveMembers(f *os.File) (map[[32]byte]bool, error) {
	head := make([]byte, 512)
	n, err := io.ReadFull(f, head)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return nil, err
	}
	head = head[:n]
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	switch {
	case bytes.HasPrefix(head, []byte("\xfd7zXZ\x00")):
		return decompressed(f, "xz")
	case bytes.HasPrefix(head, []byte("\x28\xb5\x2f\xfd")):
		return decompressed(f, "zstd")
	case bytes.HasPrefix(head, []byte("\x1f\x8b")):
		zr, err := gzip.NewReader(f)
		if err != nil {
			return map[[32]byte]bool{}, nil // not gzip after all: not an archive
		}
		defer zr.Close()
		return tarMembers(zr), nil
	case bytes.HasPrefix(head, []byte("BZh")):
		return tarMembers(bzip2.NewReader(f)), nil
	case len(head) >= 262 && string(head[257:262]) == "ustar":
		return tarMembers(f), nil
	}
	return map[[32]byte]bool{}, nil
}

// decompressed runs tool -dc over f as nobody and reads the tar stream it prints.
func decompressed(f *os.File, tool string) (map[[32]byte]bool, error) {
	bin, err := systemTool(tool)
	if err != nil {
		return map[[32]byte]bool{}, nil // no decompressor: its archives vouch for nothing
	}
	ctx, cancel := context.WithTimeout(context.Background(), decompressTimeout)
	defer cancel()
	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, bin, "-dc")
	cmd.Stdin = f
	cmd.Stderr = &limitedWriter{w: &stderr, n: 4096}
	cmd.Env = []string{}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{Uid: nobody, Gid: nobody},
		Pdeathsig:  syscall.SIGKILL,
	}
	stream, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("starting %s: %w", tool, err)
	}
	sums := tarMembers(stream)
	_, _ = io.Copy(io.Discard, stream)
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("%s: %w: %s", tool, err, strings.TrimSpace(stderr.String()))
	}
	return sums, nil
}

// tarMembers sums r's regular members up to the bounds. A stream that is no tar, or is cut short,
// vouches for what it held so far.
func tarMembers(r io.Reader) map[[32]byte]bool {
	sums := map[[32]byte]bool{}
	tr := tar.NewReader(io.LimitReader(r, maxArchiveBytes))
	for range maxArchiveMembers {
		h, err := tr.Next()
		if err != nil {
			break
		}
		if h.Typeflag != tar.TypeReg {
			continue
		}
		d := sha256.New()
		if _, err := io.Copy(d, tr); err != nil { //nolint:gosec // bounded: maxArchiveBytes caps the stream
			break
		}
		var sum [32]byte
		copy(sum[:], d.Sum(nil))
		sums[sum] = true
	}
	return sums
}
