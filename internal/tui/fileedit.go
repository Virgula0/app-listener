package tui

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
// operates on a descriptor opened beneath the vault (vaultFS.withFile), never on a name a planted
// symlink could redirect.
func applyNewFileMeta(v *vaultFS, path string, isDir bool) error {
	flags := os.O_RDONLY
	if isDir {
		flags |= unix.O_DIRECTORY
	}
	return v.withFile(path, flags, func(f *os.File) error {
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
		return f.Chmod(defaultNewFileMode(isDir))
	})
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

// writeFileKeepMeta atomically replaces name (relative to vault) with data (temp sibling + rename
// after fsync), preserving mode and owner (needs root). The parent is opened beneath the vault and
// every later step works on that descriptor: the temp sibling is created O_EXCL|O_NOFOLLOW and
// metadata is applied on its fd, never on a name a planted symlink could redirect.
func writeFileKeepMeta(vault *os.Root, name string, data []byte) error {
	dir, err := vault.Open(filepath.Dir(name))
	if err != nil {
		return err
	}
	defer dir.Close()
	dirFD := int(dir.Fd())
	base := filepath.Base(name)
	var st unix.Stat_t
	if err := unix.Fstatat(dirFD, base, &st, unix.AT_SYMLINK_NOFOLLOW); err != nil {
		return fmt.Errorf("stating %s: %w", name, err)
	}
	if st.Mode&unix.S_IFMT != unix.S_IFREG {
		return fmt.Errorf("refusing to write %s: not a regular file", name)
	}
	uid, gid := -1, -1
	if os.Geteuid() == 0 {
		uid, gid = int(st.Uid), int(st.Gid)
	}
	tmp := "." + base + ".app_listener.edit"
	return safeio.AtomicWriteAt(dirFD, base, tmp, data, os.FileMode(st.Mode&0o777), uid, gid)
}

// writeFileInPlace rewrites f's inode (truncate, write, fsync), never a temp file or rename. Only
// for a single-file edit-protected resource (fileEditModel.singleFile): a single-file watch root's
// guard also protects its parent directory against anything created/renamed beside it
// (guard_path_rename's destination-parent check, the rename-over-watchroot defense), so
// writeFileKeepMeta's temp-sibling dance is denied. Mode and ownership stay as they are.
func writeFileInPlace(f *os.File, data []byte) error {
	if err := f.Truncate(0); err != nil {
		return err
	}
	if _, err := f.WriteAt(data, 0); err != nil {
		return err
	}
	return f.Sync()
}

// vaultFS pins the opened vault for the whole editor session. The editor runs as root over a
// user-owned tree, so re-walking an absolute path lets that user redirect a write, chmod, chown or
// delete by swapping any directory on it (inside the vault or above it) for a symlink. Every access
// resolves beneath the pinned directory (os.Root: no symlink or ".." leaves it) or goes through the
// pinned single file.
type vaultFS struct {
	base string
	dir  *os.Root
	file *os.File // single-file resource
}

// openVault pins path: a regular file O_RDWR|O_NOFOLLOW, a directory as an os.Root. A symlinked
// root is refused, never followed.
func openVault(path string) (*vaultFS, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	switch {
	case info.Mode().IsRegular():
		f, err := safeio.OpenRegularNoFollow(path, os.O_RDWR, 0)
		if err != nil {
			return nil, err
		}
		return &vaultFS{base: path, file: f}, nil
	case info.IsDir():
		r, err := os.OpenRoot(path)
		if err != nil {
			return nil, err
		}
		// OpenRoot follows a final symlink: reject a swap since the Lstat.
		if got, err := r.Stat("."); err != nil || !os.SameFile(got, info) {
			_ = r.Close()
			return nil, fmt.Errorf("%s changed while opening it", path)
		}
		return &vaultFS{base: path, dir: r}, nil
	default:
		return nil, fmt.Errorf("%s is neither a regular file nor a directory", path)
	}
}

func (v *vaultFS) close() {
	if v.dir != nil {
		_ = v.dir.Close()
	}
	if v.file != nil {
		_ = v.file.Close()
	}
}

// rel maps an absolute path the tree holds to its name beneath the vault.
func (v *vaultFS) rel(path string) (string, error) {
	if v.dir == nil && v.file == nil {
		return "", fmt.Errorf("the vault %s is not open", v.base)
	}
	r, err := filepath.Rel(v.base, path)
	if err != nil || r == ".." || strings.HasPrefix(r, "../") {
		return "", fmt.Errorf("%s is outside the vault", path)
	}
	if v.file != nil && r != "." {
		return "", fmt.Errorf("%s: a single-file resource has no entries", path)
	}
	return r, nil
}

func (v *vaultFS) lstat(path string) (os.FileInfo, error) {
	r, err := v.rel(path)
	if err != nil {
		return nil, err
	}
	if v.file != nil {
		return v.file.Stat()
	}
	return v.dir.Lstat(r)
}

func (v *vaultFS) readFile(path string) ([]byte, error) {
	r, err := v.rel(path)
	if err != nil {
		return nil, err
	}
	if v.file != nil {
		return io.ReadAll(io.NewSectionReader(v.file, 0, math.MaxInt64))
	}
	return v.dir.ReadFile(r)
}

func (v *vaultFS) readDir(path string) ([]os.DirEntry, error) {
	r, err := v.rel(path)
	if err != nil {
		return nil, err
	}
	if v.dir == nil {
		return nil, fmt.Errorf("%s is not a directory", path)
	}
	d, err := v.dir.Open(r)
	if err != nil {
		return nil, err
	}
	defer d.Close()
	return d.ReadDir(-1)
}

func (v *vaultFS) write(path string, data []byte) error {
	r, err := v.rel(path)
	if err != nil {
		return err
	}
	if v.file != nil {
		return writeFileInPlace(v.file, data)
	}
	return writeFileKeepMeta(v.dir, r, data)
}

// dirOp runs op on path's name beneath a directory vault; the vault root itself is refused.
func (v *vaultFS) dirOp(path string, op func(root *os.Root, name string) error) error {
	r, err := v.rel(path)
	if err != nil {
		return err
	}
	if v.dir == nil || r == "." {
		return fmt.Errorf("refusing to change the vault root %s", path)
	}
	return op(v.dir, r)
}

// withFile runs fn on path opened beneath the vault with O_NOFOLLOW added: fn acts on that
// descriptor, not on a name.
func (v *vaultFS) withFile(path string, flags int, fn func(*os.File) error) error {
	r, err := v.rel(path)
	if err != nil {
		return err
	}
	if v.file != nil {
		return fn(v.file)
	}
	f, err := v.dir.OpenFile(r, flags|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	return fn(f)
}

// dirEntry is one row of the right-pane directory listing.
type dirEntry struct {
	name  string
	isDir bool
	size  int64
	mtime time.Time
}

// listEntries stats up to limit children, skipping vanished entries.
func listEntries(v *vaultFS, n *fileNode, limit int) []dirEntry {
	out := make([]dirEntry, 0, min(limit, len(n.children)))
	for _, c := range n.children {
		if len(out) >= limit {
			break
		}
		info, err := v.lstat(c.path)
		if err != nil {
			continue
		}
		out = append(out, dirEntry{name: c.name, isDir: c.isDir, size: info.Size(), mtime: info.ModTime()})
	}
	return out
}
