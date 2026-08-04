package install

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
)

// CopyTree recursively copies the tree rooted at src into dst, preserving
// directory and file permissions, ownership and modification times. Symbolic
// links are recreated; sockets, FIFOs and device nodes are skipped (they
// cannot be copied meaningfully). The destination must not exist yet.
func CopyTree(src, dst string) error {
	return CopyTreeWithProgress(src, dst, nil)
}

// CopyTreeWithProgress behaves like CopyTree and additionally invokes
// onBytes after every regular file with the bytes copied so far and the
// total bytes to copy (a nil onBytes skips the pre-measurement walk and
// all reporting).
func CopyTreeWithProgress(src, dst string, onBytes func(copied, total int64)) error {
	var total int64
	if onBytes != nil {
		var err error
		if total, err = treeSize(src); err != nil {
			return fmt.Errorf("measuring %s: %w", src, err)
		}
	}
	var copied int64
	var report func(int64)
	if onBytes != nil {
		report = func(n int64) {
			copied += n
			onBytes(copied, total)
		}
	}
	return copyTree(src, dst, 0, report)
}

// treeSize sums the sizes of every regular file under root (the pre-walk
// used to compute the copy total).
func treeSize(root string) (int64, error) {
	var total int64
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.Type().IsRegular() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			total += info.Size()
		}
		return nil
	})
	return total, err
}

func copyTree(src, dst string, depth int, report func(int64)) error {
	if depth > 64 {
		return fmt.Errorf("directory tree too deep at %s", src)
	}
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	switch {
	case info.IsDir():
		return copyDir(src, dst, info, depth, report)
	case info.Mode()&os.ModeSymlink != 0:
		return copySymlink(src, dst)
	case info.Mode().IsRegular():
		return copyFile(src, dst, info, report)
	default:
		// Sockets, FIFOs, device nodes: nothing meaningful to copy.
		return nil
	}
}

func copyDir(src, dst string, info fs.FileInfo, depth int, report func(int64)) error {
	if err := os.Mkdir(dst, info.Mode().Perm()); err != nil && !os.IsExist(err) {
		return err
	}
	if err := os.Chmod(dst, info.Mode().Perm()); err != nil {
		return err
	}
	if err := preserveMeta(dst, info); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if err := copyTree(filepath.Join(src, entry.Name()), filepath.Join(dst, entry.Name()), depth+1, report); err != nil {
			return err
		}
	}
	return nil
}

func copySymlink(src, dst string) error {
	target, err := os.Readlink(src)
	if err != nil {
		return err
	}
	if err := os.Symlink(target, dst); err != nil && !os.IsExist(err) {
		return err
	}
	return nil
}

func copyFile(src, dst string, info fs.FileInfo, report func(int64)) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode().Perm())
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	if err := preserveMeta(dst, info); err != nil {
		return err
	}
	if report != nil {
		report(info.Size())
	}
	return nil
}

// preserveMeta restores ownership (best effort, requires root) and
// modification time on a freshly created path.
func preserveMeta(dst string, info fs.FileInfo) error {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		if err := os.Chown(dst, int(stat.Uid), int(stat.Gid)); err != nil && !os.IsPermission(err) {
			return err
		}
	}
	return os.Chtimes(dst, info.ModTime(), info.ModTime())
}
