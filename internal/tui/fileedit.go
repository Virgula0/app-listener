package tui

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"
)

// defaultNewFileMode returns the mode a freshly created entry gets when no
// explicit octal mode was given: the conventional 0666/0777 masked by the
// caller's umask (022), i.e. 0644 files and 0755 directories.
func defaultNewFileMode(isDir bool) os.FileMode {
	mask := os.FileMode(0o022)
	if isDir {
		return 0o777 &^ mask
	}
	return 0o666 &^ mask
}

// applyNewFileMeta gives a freshly created entry the ownership it would
// have had if the invoking user had created it directly. When running as
// root under sudo the real caller (SUDO_UID/SUDO_GID) becomes the owner —
// natively, with a single chown, no subprocess — and the mode follows the
// standard umask default modeled by defaultNewFileMode. Non-root callers
// cannot chown and are only given the default mode.
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

// isOctalMode reports whether v is a 1-4 digit octal permission mode like
// 0644. Leading zeros are allowed since 3- and 4-digit modes (0644 vs 644)
// mean the same permission bits; a parser of the value must still reject
// anything longer than 07777.
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

// unixModeToFileMode converts a conventional unix mode number (0o0000 to
// 0o7777, as typed in the chmod dialog) into Go's os.FileMode. The two are
// not interchangeable: os.FileMode keeps only the rwx bits in its low bits
// and carries setuid, setgid and sticky as separate high flags, so a raw
// cast drops the special bits (a raw-cast chmod 4755 silently applied 0755).
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

// fileModeToUnixMode is the inverse: it renders an os.FileMode as the
// conventional octal mode number including the special bits.
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

// clipLabel truncates s to at most w runes, appending an ellipsis when it
// was shortened.
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

// isBinaryContent reports whether data looks like a binary file: any NUL
// byte inside the first binaryScanBytes bytes classifies it as binary, and
// such files are never handed to the editor.
func isBinaryContent(data []byte) bool {
	n := min(len(data), binaryScanBytes)
	return bytes.IndexByte(data[:n], 0) >= 0
}

// writeFileKeepMeta atomically replaces the file at path with data while
// preserving its mode and owner. The new content is written to a temporary
// sibling and renamed over the original, so a crash or an error can never
// leave the file half-written. Only regular files may be replaced: writing
// through a symlink is refused.
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
	// Preserving the owner is only possible as root (the editor is used by
	// `install --edit-protected`, which runs as root); a non-root user's
	// temporary file already carries the caller's ownership.
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

// listEntries stats up to limit children of n for the right-pane listing.
// Entries that vanish while listing are skipped.
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
