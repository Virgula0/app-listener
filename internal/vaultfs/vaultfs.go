// Package vaultfs is the edit-protected editor's file access, shared by the TUI and the GUI's
// backend (`edit-protected --backend`). It runs as root over a tree owned by an unprivileged user,
// often with no guard attached (offline flow), so no operation may follow a symlink below the
// opened root or re-walk a path that user could swap.
package vaultfs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/safeio"
)

// Limits for the embedded editors: binary content is never edited; files past the edit limit are
// read-only.
const (
	MaxEditableBytes = 2 << 20 // 2 MiB
	BinaryScanBytes  = 8 << 10 // 8 KiB of a file are scanned for NUL bytes
)

// DefaultNewFileMode models umask(022) defaults: 0644 files, 0755 dirs.
func DefaultNewFileMode(isDir bool) os.FileMode {
	mask := os.FileMode(0o022)
	if isDir {
		return 0o777 &^ mask
	}
	return 0o666 &^ mask
}

// IsOctalMode reports a 1-4 digit octal mode (leading zeros allowed);
// callers must still reject values above 07777.
func IsOctalMode(v string) bool {
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

// UnixModeToFileMode converts a unix mode number to os.FileMode; a raw cast
// would silently drop the setuid/setgid/sticky bits.
func UnixModeToFileMode(mode uint32) os.FileMode {
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

// FileModeToUnixMode is the inverse, restoring the special bits.
func FileModeToUnixMode(m os.FileMode) uint32 {
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

// IsBinaryContent spots a NUL byte in the first BinaryScanBytes bytes.
func IsBinaryContent(data []byte) bool {
	n := min(len(data), BinaryScanBytes)
	return bytes.IndexByte(data[:n], 0) >= 0
}

// WriteFileKeepMeta atomically replaces name (relative to vault) with data (temp sibling + rename
// after fsync), preserving mode and owner (needs root). The parent is opened beneath the vault and
// every later step works on that descriptor: the temp sibling is created O_EXCL|O_NOFOLLOW and
// metadata is applied on its fd, never on a name a planted symlink could redirect.
func WriteFileKeepMeta(vault *os.Root, name string, data []byte) error {
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

// WriteFileInPlace rewrites f's inode (truncate, write, fsync), never a temp file or rename. Only
// for a single-file edit-protected resource (Vault.SingleFile): a single-file watch root's
// guard also protects its parent directory against anything created/renamed beside it
// (guard_path_rename's destination-parent check, the rename-over-watchroot defense), so
// WriteFileKeepMeta's temp-sibling dance is denied. Mode and ownership stay as they are.
func WriteFileInPlace(f *os.File, data []byte) error {
	if err := f.Truncate(0); err != nil {
		return err
	}
	if _, err := f.WriteAt(data, 0); err != nil {
		return err
	}
	return f.Sync()
}

// Vault pins the opened vault for the whole editor session. The editor runs as root over a
// user-owned tree, so re-walking an absolute path lets that user redirect a write, chmod, chown or
// delete by swapping any directory on it (inside the vault or above it) for a symlink. Every access
// resolves beneath the pinned directory (os.Root: no symlink or ".." leaves it) or goes through the
// pinned single file.
type Vault struct {
	base string
	dir  *os.Root
	file *os.File // single-file resource
}

// Open pins path: a regular file O_RDWR|O_NOFOLLOW, a directory as an os.Root. A symlinked
// root is refused, never followed.
func Open(path string) (*Vault, error) {
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
		return &Vault{base: path, file: f}, nil
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
		return &Vault{base: path, dir: r}, nil
	default:
		return nil, fmt.Errorf("%s is neither a regular file nor a directory", path)
	}
}

// Unopened is a placeholder for a vault Open refused: every operation fails.
func Unopened(path string) *Vault { return &Vault{base: path} }

// Base is the path the vault was opened at.
func (v *Vault) Base() string { return v.base }

// SingleFile reports a single-file resource (no entries; writes go in place).
func (v *Vault) SingleFile() bool { return v.file != nil }

// Close releases the pinned root.
func (v *Vault) Close() {
	if v.dir != nil {
		_ = v.dir.Close()
	}
	if v.file != nil {
		_ = v.file.Close()
	}
}

// rel maps an absolute path the tree holds to its name beneath the vault.
func (v *Vault) rel(path string) (string, error) {
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

// Lstat stats path without following a final symlink.
func (v *Vault) Lstat(path string) (os.FileInfo, error) {
	r, err := v.rel(path)
	if err != nil {
		return nil, err
	}
	if v.file != nil {
		return v.file.Stat()
	}
	return v.dir.Lstat(r)
}

// ReadFile reads path whole.
func (v *Vault) ReadFile(path string) ([]byte, error) {
	r, err := v.rel(path)
	if err != nil {
		return nil, err
	}
	if v.file != nil {
		return io.ReadAll(io.NewSectionReader(v.file, 0, math.MaxInt64))
	}
	return v.dir.ReadFile(r)
}

// ReadDir lists path's entries.
func (v *Vault) ReadDir(path string) ([]os.DirEntry, error) {
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

// Write replaces an existing regular file's content, keeping its mode and owner.
func (v *Vault) Write(path string, data []byte) error {
	r, err := v.rel(path)
	if err != nil {
		return err
	}
	if v.file != nil {
		return WriteFileInPlace(v.file, data)
	}
	return WriteFileKeepMeta(v.dir, r, data)
}

// DirOp runs op on path's name beneath a directory vault; the vault root itself is refused.
func (v *Vault) DirOp(path string, op func(root *os.Root, name string) error) error {
	r, err := v.rel(path)
	if err != nil {
		return err
	}
	if v.dir == nil || r == "." {
		return fmt.Errorf("refusing to change the vault root %s", path)
	}
	return op(v.dir, r)
}

// WithFile runs fn on path opened beneath the vault with O_NOFOLLOW added: fn acts on that
// descriptor, not on a name.
func (v *Vault) WithFile(path string, flags int, fn func(*os.File) error) error {
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

// Create makes an empty file (0600, O_EXCL|O_NOFOLLOW: an existing symlink at target is refused,
// never followed) or a directory (0750) named name inside dir, and returns its path. ApplyNewMeta
// gives it its final owner and mode.
func (v *Vault) Create(dir, name string, isDir bool) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" || strings.Contains(name, "/") || name == "." || name == ".." {
		return "", errors.New("invalid name: must be non-empty and contain no '/'")
	}
	target := filepath.Join(dir, name)
	err := v.DirOp(target, func(r *os.Root, name string) error {
		if isDir {
			return r.Mkdir(name, 0o750)
		}
		f, openErr := r.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL|unix.O_NOFOLLOW, 0o600)
		if openErr != nil {
			return openErr
		}
		return f.Close()
	})
	if err != nil {
		return "", err
	}
	return target, nil
}

// ApplyNewMeta gives fresh entries the ownership/mode the invoking user would have gotten:
// under sudo it chowns to SUDO_UID/SUDO_GID; non-root callers only get DefaultNewFileMode. It
// operates on a descriptor opened beneath the vault (WithFile), never on a name a planted
// symlink could redirect.
func (v *Vault) ApplyNewMeta(path string, isDir bool) error {
	flags := os.O_RDONLY
	if isDir {
		flags |= unix.O_DIRECTORY
	}
	return v.WithFile(path, flags, func(f *os.File) error {
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
		return f.Chmod(DefaultNewFileMode(isDir))
	})
}

// ParseMode parses a 1-4 digit octal mode (e.g. 0644, 4755).
func ParseMode(v string) (uint32, error) {
	v = strings.TrimSpace(v)
	if !IsOctalMode(v) {
		return 0, errors.New("invalid mode: 1-4 octal digits (e.g. 0644)")
	}
	perm, _ := strconv.ParseUint(v, 8, 16)
	if perm > 0o7777 {
		return 0, errors.New("invalid mode: must not exceed 07777")
	}
	return uint32(perm), nil
}

// Chmod sets path's unix mode (setuid/setgid/sticky included) on an O_PATH|O_NOFOLLOW descriptor,
// re-validated at apply time: a symlink is refused.
func (v *Vault) Chmod(path string, mode uint32) error {
	return v.WithFile(path, unix.O_PATH, func(f *os.File) error {
		if info, statErr := f.Stat(); statErr != nil || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("%s is a symlink — chmod is refused", path)
		}
		// O_PATH descriptors can't fchmod directly; chmod via /proc/self/fd resolves the pinned
		// inode without re-walking the (swappable) path.
		return os.Chmod(fmt.Sprintf("/proc/self/fd/%d", f.Fd()), UnixModeToFileMode(mode))
	})
}

// Chown sets path's owner on an O_PATH|O_NOFOLLOW descriptor; a symlink is refused.
func (v *Vault) Chown(path string, uid, gid int) error {
	return v.WithFile(path, unix.O_PATH, func(f *os.File) error {
		if info, statErr := f.Stat(); statErr != nil || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("%s is a symlink — chown is refused", path)
		}
		return unix.Fchownat(int(f.Fd()), "", uid, gid, unix.AT_EMPTY_PATH)
	})
}

// Remove deletes path (recursively for a directory) beneath the vault; the root is refused.
func (v *Vault) Remove(path string) error {
	return v.DirOp(path, (*os.Root).RemoveAll)
}
