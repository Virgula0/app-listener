package ebpf

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"golang.org/x/sys/unix"
)

// hashBufPool hands out reusable 128 KiB read buffers so hashing never allocates a file-sized
// slice: the binary verifier re-hashes large binaries (Electron/Chromium blobs, 100-200 MiB) on a
// timer across many guards, and os.ReadFile allocated and zeroed gigabytes per minute (heap/CPU
// profiles).
var hashBufPool = sync.Pool{New: func() any { b := make([]byte, 128*1024); return &b }}

// BinaryEntry describes an executable the engines key on: path, sha256 of contents, and comm (task
// name, truncated to the kernel's 16 bytes).
type BinaryEntry struct {
	Path string
	Hash [sha256.Size]byte
	Comm string
}

// ComputeBinaryEntry hashes the file at path (streamed through a pooled buffer, never read whole)
// and derives its comm. Shared by the guard, network guard and network monitor.
func ComputeBinaryEntry(path string) (BinaryEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return BinaryEntry{}, fmt.Errorf("reading binary %s: %w", path, err)
	}
	defer f.Close()
	return entryOf(f, path)
}

// entryOf hashes r into an entry named path.
func entryOf(r io.Reader, path string) (BinaryEntry, error) {
	h := sha256.New()
	bufp := hashBufPool.Get().(*[]byte)
	_, err := io.CopyBuffer(h, r, *bufp)
	hashBufPool.Put(bufp)
	if err != nil {
		return BinaryEntry{}, fmt.Errorf("reading binary %s: %w", path, err)
	}
	e := BinaryEntry{Path: path, Comm: filepath.Base(path)}
	h.Sum(e.Hash[:0])
	if len(e.Comm) > 15 {
		e.Comm = e.Comm[:15]
	}
	return e, nil
}

// KernelDev encodes major:minor as the kernel's internal dev_t (sb->s_dev), which the BPF programs
// key on; userspace's st_dev encoding differs.
func KernelDev(major, minor uint32) uint64 {
	return uint64(major)<<20 | uint64(minor)
}

// StatInode returns the dev/ino pair of path, with dev encoded by KernelDev.
func StatInode(path string) (dev, ino uint64, err error) {
	var s syscall.Stat_t
	if err := syscall.Stat(path, &s); err != nil {
		return 0, 0, err
	}
	return KernelDev(unix.Major(s.Dev), unix.Minor(s.Dev)), s.Ino, nil
}

// LstatInode is StatInode without following a final symlink: it identifies
// the directory entry itself rather than whatever it points at.
func LstatInode(path string) (dev, ino uint64, err error) {
	var s syscall.Stat_t
	if err := syscall.Lstat(path, &s); err != nil {
		return 0, 0, err
	}
	return KernelDev(unix.Major(s.Dev), unix.Minor(s.Dev)), s.Ino, nil
}

// BinaryStat is a cheap change-detection fingerprint: recompute the hash only when Size, MtimeNs or
// CtimeNs moved (an in-place overwrite always bumps mtime and ctime; ctime can't be restored
// without clock tampering).
type BinaryStat struct {
	Dev, Ino         uint64
	Size             int64
	MtimeNs, CtimeNs int64
}

// StatBinary returns the fingerprint of path (dev encoded as StatInode does).
func StatBinary(path string) (BinaryStat, error) {
	var s syscall.Stat_t
	if err := syscall.Stat(path, &s); err != nil {
		return BinaryStat{}, err
	}
	return BinaryStat{
		Dev:     KernelDev(unix.Major(s.Dev), unix.Minor(s.Dev)),
		Ino:     s.Ino,
		Size:    s.Size,
		MtimeNs: s.Mtim.Sec*1_000_000_000 + s.Mtim.Nsec,
		CtimeNs: s.Ctim.Sec*1_000_000_000 + s.Ctim.Nsec,
	}, nil
}

// BinariesSummary renders a whitelist as a compact log-friendly list.
func BinariesSummary(binaries []BinaryEntry) string {
	parts := make([]string, len(binaries))
	for i, b := range binaries {
		parts[i] = fmt.Sprintf("%s [sha256:%x..%x]", b.Path, b.Hash[:4], b.Hash[28:])
	}
	return strings.Join(parts, ", ")
}
