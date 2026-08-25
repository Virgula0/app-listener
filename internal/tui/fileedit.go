package tui

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"
)

// defaultNewFileMode models umask(022) defaults: 0644 files, 0755 dirs.
func defaultNewFileMode(isDir bool) os.FileMode {
	mask := os.FileMode(0o022)
	if isDir {
		return 0o777 &^ mask
	}
	return 0o666 &^ mask
}

// applyNewFileMeta gives fresh entries the ownership/mode the invoking user
// would have gotten: under sudo it chowns to SUDO_UID/SUDO_GID natively;
// non-root callers only get the defaultNewFileMode mode.
func applyNewFileMeta(path string, isDir bool) error {
	mode := defaultNewFileMode(isDir)
	if os.Geteuid() == 0 {
		uidStr, gidStr := os.Getenv("SUDO_UID"), os.Getenv("SUDO_GID")
		if uidStr != "" && gidStr != "" {
			uid, uidErr := strconv.ParseUint(uidStr, 10, 32)
			gid, gidErr := strconv.ParseUint(gidStr, 10, 32)
			if uidErr != nil || gidErr != nil {
				return fmt.Errorf("parsing SUDO_UID/SUDO_GID: %v / %v", uidErr, gidErr)
			}
			if err := os.Lchown(path, int(uid), int(gid)); err != nil {
				return fmt.Errorf("chown %s: %w", path, err)
			}
		}
	}
	return os.Chmod(path, mode)
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

// writeFileKeepMeta atomically replaces path with data (temp sibling +
// rename after fsync), preserving mode and owner; symlinks are refused so
// writes never follow links. Owner restoration requires root.
func writeFileKeepMeta(path string, data []byte) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("refusing to write through symlink %s", path)
	}
	tmp := path + ".app_listener.edit"
	f, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode().Perm())
	if err != nil {
		return err
	}
	cleanup := true
	defer func() {
		_ = f.Close()
		if cleanup {
			_ = os.Remove(tmp)
		}
	}()
	if _, werr := f.Write(data); werr != nil {
		return werr
	}
	if serr := f.Sync(); serr != nil {
		return serr
	}
	if cerr := f.Chmod(info.Mode()); cerr != nil {
		return cerr
	}
	// Chown needs root (the editor runs as root via --edit-protected);
	// non-root temp files already carry the caller's ownership.
	if os.Geteuid() == 0 {
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			if cerr := f.Chown(int(stat.Uid), int(stat.Gid)); cerr != nil {
				return cerr
			}
		}
	}
	if cerr := f.Close(); cerr != nil {
		return cerr
	}
	cleanup = false
	return os.Rename(tmp, path)
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
