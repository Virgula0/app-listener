package tui

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/safeio"
)

// defaultNewFileMode models umask(022) defaults: 0644 files, 0755 dirs.
func defaultNewFileMode(isDir bool) os.FileMode {
	mask := os.FileMode(0o022)
	if isDir {
		return 0o777 &^ mask
	}
	return 0o666 &^ mask
}

// applyNewFileMeta gives fresh entries the ownership/mode the invoking user would have gotten:
// under sudo it chowns to SUDO_UID/SUDO_GID; non-root callers only get defaultNewFileMode. It
// operates on a symlink-safe descriptor (O_NOFOLLOW): os.Chmod follows a symlink, so a name planted
// as a link before the chmod would otherwise have root chmod the link target.
func applyNewFileMeta(path string, isDir bool) error {
	mode := defaultNewFileMode(isDir)
	flags := unix.O_NOFOLLOW | unix.O_CLOEXEC
	if isDir {
		flags |= unix.O_DIRECTORY
	} else {
		flags |= unix.O_RDONLY
	}
	fd, err := unix.Open(path, flags, 0)
	if err != nil {
		return fmt.Errorf("opening %s: %w", path, err)
	}
	f := os.NewFile(uintptr(fd), path)
	defer f.Close()
	if os.Geteuid() == 0 {
		uidStr, gidStr := os.Getenv("SUDO_UID"), os.Getenv("SUDO_GID")
		if uidStr != "" && gidStr != "" {
			uid, uidErr := strconv.ParseUint(uidStr, 10, 32)
			gid, gidErr := strconv.ParseUint(gidStr, 10, 32)
			if uidErr != nil || gidErr != nil {
				return fmt.Errorf("parsing SUDO_UID/SUDO_GID: %v / %v", uidErr, gidErr)
			}
			if err := f.Chown(int(uid), int(gid)); err != nil {
				return fmt.Errorf("chown %s: %w", path, err)
			}
		}
	}
	return f.Chmod(mode)
}

// isOctalMode reports a 1-4 digit octal mode (leading zeros allowed);
// callers must still reject values above 07777.
func isOctalMode(v string) bool {
	if len(v) < 1 || len(v) > 4 {
		return false
	}
	for _, r := range v {
		if r < '0' || r > '7' {
			return false
		}
	}
	return true
}

// unixModeToFileMode converts a unix mode number to os.FileMode; a raw cast
// would silently drop the setuid/setgid/sticky bits.
func unixModeToFileMode(mode uint32) os.FileMode {
	m := os.FileMode(mode & 0o777)
	if mode&0o4000 != 0 {
		m |= os.ModeSetuid
	}
	if mode&0o2000 != 0 {
		m |= os.ModeSetgid
	}
	if mode&0o1000 != 0 {
		m |= os.ModeSticky
	}
	return m
}

// fileModeToUnixMode is the inverse, restoring the special bits.
func fileModeToUnixMode(m os.FileMode) uint32 {
	mode := uint32(m.Perm())
	if m&os.ModeSetuid != 0 {
		mode |= 0o4000
	}
	if m&os.ModeSetgid != 0 {
		mode |= 0o2000
	}
	if m&os.ModeSticky != 0 {
		mode |= 0o1000
	}
	return mode
}

// humanSize renders a byte count with a binary unit suffix.
func humanSize(n int64) string {
	switch {
	case n < 1<<10:
		return fmt.Sprintf("%d B", n)
	case n < 1<<20:
		return fmt.Sprintf("%.1f KiB", float64(n)/(1<<10))
	case n < 1<<30:
		return fmt.Sprintf("%.1f MiB", float64(n)/(1<<20))
	default:
		return fmt.Sprintf("%.1f GiB", float64(n)/(1<<30))
	}
}

// clipLabel truncates s to w runes, appending an ellipsis when shortened.
func clipLabel(s string, w int) string {
	if w <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= w {
		return s
	}
	if w == 1 {
		return "…"
	}
	return string(r[:w-1]) + "…"
}

// isBinaryContent spots a NUL byte in the first binaryScanBytes bytes.
func isBinaryContent(data []byte) bool {
	n := min(len(data), binaryScanBytes)
	return bytes.IndexByte(data[:n], 0) >= 0
}

// writeFileKeepMeta atomically replaces path with data (temp sibling + rename after fsync),
// preserving mode and owner (needs root). The editor runs as root over a user-owned tree with no
// guard attached in the offline flow, so every step is symlink-safe: the temp sibling is created
// O_EXCL|O_NOFOLLOW inside the parent dir fd and metadata is applied on that descriptor, never on a
// name that a planted symlink could redirect (a plain O_CREATE|O_TRUNC on a predictable temp name
// let root truncate and chown the link target).
func writeFileKeepMeta(path string, data []byte) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("refusing to write through symlink %s", path)
	}
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	dirFD, err := unix.Open(dir, unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if err != nil {
		return fmt.Errorf("opening %s: %w", dir, err)
	}
	defer unix.Close(dirFD)

	uid, gid := -1, -1
	if os.Geteuid() == 0 {
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			uid, gid = int(stat.Uid), int(stat.Gid)
		}
	}
	tmp := "." + base + ".app_listener.edit"
	return safeio.AtomicWriteAt(dirFD, base, tmp, data, info.Mode().Perm(), uid, gid)
}

// writeFileInPlace rewrites path's EXISTING inode (open, truncate, write, fsync), never a temp file
// or rename. Only for a single-file edit-protected resource (fileEditModel.singleFile): a
// single-file watch root's guard also protects its parent directory against anything
// created/renamed beside it (guard_path_rename's destination-parent check, the
// rename-over-watchroot defense), so writeFileKeepMeta's temp-sibling dance is denied. Mode and
// ownership stay as they are (no temp file whose metadata needs copying).
func writeFileInPlace(path string, data []byte) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("refusing to write through symlink %s", path)
	}
	// O_NOFOLLOW + fstat: validate the descriptor, not the earlier Lstat, so a symlink swapped in
	// between the two cannot redirect the truncate+write.
	f, err := safeio.OpenRegularNoFollow(path, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := f.Truncate(0); err != nil {
		return err
	}
	if _, err := f.WriteAt(data, 0); err != nil {
		return err
	}
	return f.Sync()
}

// dirEntry is one row of the right-pane directory listing.
type dirEntry struct {
	name  string
	isDir bool
	size  int64
	mtime time.Time
}

// listEntries stats up to limit children, skipping vanished entries.
func listEntries(n *fileNode, limit int) []dirEntry {
	out := make([]dirEntry, 0, min(limit, len(n.children)))
	for _, c := range n.children {
		if len(out) >= limit {
			break
		}
		info, err := os.Lstat(c.path)
		if err != nil {
			continue
		}
		out = append(out, dirEntry{name: c.name, isDir: c.isDir, size: info.Size(), mtime: info.ModTime()})
	}
	return out
}
