package ebpf

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/unix"
)

// BinaryEntry describes an executable the engines key on: its path, a
// sha256 hash of its contents, and its comm (task name, truncated to the
// kernel's 16-byte limit).
type BinaryEntry struct {
	Path string
	Hash [sha256.Size]byte
	Comm string
}

// ComputeBinaryEntry hashes the file at path and derives its comm. It is
// shared by the guard, network guard and network monitor engines.
func ComputeBinaryEntry(path string) (BinaryEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return BinaryEntry{}, fmt.Errorf("reading binary %s: %w", path, err)
	}

	hash := sha256.Sum256(data)
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
