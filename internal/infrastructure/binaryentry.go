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

// hashBufPool hands out reusable 128 KiB read buffers so hashing a binary
// never allocates a file-sized slice. The guard's binary verifier re-hashes
// large binaries (Electron/Chromium blobs are 100-200 MiB) on a timer across
// many guards; os.ReadFile there was allocating — and zeroing — gigabytes per
// minute (confirmed by heap/CPU profiles).
var hashBufPool = sync.Pool{New: func() any { b := make([]byte, 128*1024); return &b }}

// BinaryEntry describes an executable the engines key on: its path, a
// sha256 hash of its contents, and its comm (task name, truncated to the
// kernel's 16-byte limit).
type BinaryEntry struct {
	Path string
	Hash [sha256.Size]byte
	Comm string
}

// ComputeBinaryEntry hashes the file at path and derives its comm. It is
// shared by the guard, network guard and network monitor engines. The file is
// streamed through a pooled buffer, never read whole into memory.
func ComputeBinaryEntry(path string) (BinaryEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return BinaryEntry{}, fmt.Errorf("reading binary %s: %w", path, err)
	}
	defer f.Close()

	h := sha256.New()
	bufp := hashBufPool.Get().(*[]byte)
	buf := *bufp
	for {
		n, readErr := f.Read(buf)
		if n > 0 {
			h.Write(buf[:n])
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			hashBufPool.Put(bufp)
			return BinaryEntry{}, fmt.Errorf("reading binary %s: %w", path, readErr)
		}
	}
	hashBufPool.Put(bufp)

	var hash [sha256.Size]byte
	h.Sum(hash[:0])

	comm := filepath.Base(path)
	if len(comm) > 15 {
		comm = comm[:15]
	}

	return BinaryEntry{
		Path: path,
		Hash: hash,
		Comm: comm,
	}, nil
}

// StatInode returns the dev/ino pair of path, with dev encoded the way the
// BPF programs expect ((major << 20) | minor).
func StatInode(path string) (dev, ino uint64, err error) {
	var s syscall.Stat_t
	if err := syscall.Stat(path, &s); err != nil {
		return 0, 0, err
	}
	return uint64((unix.Major(s.Dev) << 20) | unix.Minor(s.Dev)), s.Ino, nil
}

// BinaryStat is a cheap change-detection fingerprint of a file: a hash is
// only worth recomputing when Size, MtimeNs or CtimeNs moved (an in-place
// overwrite always bumps mtime and ctime; ctime cannot be restored without
// clock tampering).
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
		Dev:     uint64((unix.Major(s.Dev) << 20) | unix.Minor(s.Dev)),
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
