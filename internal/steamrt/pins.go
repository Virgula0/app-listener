package steamrt

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

// pinDirs are the only directories setup.sh's pin_newer_runtime_libs (re)creates in the runtime.
var pinDirs = map[string]bool{
	"pinned_libs_32": true, "pinned_libs_64": true, "libcurl_compat_32": true, "libcurl_compat_64": true,
}

const (
	maxPinEntries = 4096
	maxPinFile    = 4096
	pinTimeout    = 5 * time.Minute
	stageRoot     = "/run/app-listener/steamrt"
)

// PinsStale mirrors setup.sh's check_pins: pins never completed, or a pinned host library no longer
// resolves to the file setup.sh recorded (a distro upgrade).
func (r Runtime) PinsStale() bool {
	rt := filepath.Join(r.Dir, r.Name)
	if _, err := os.Lstat(filepath.Join(rt, "setup.sh")); err != nil {
		return false
	}
	for _, b := range []string{"32", "64"} {
		if _, err := os.Lstat(filepath.Join(rt, "pinned_libs_"+b, "done")); err != nil {
			return true
		}
	}
	pins, _ := filepath.Glob(filepath.Join(rt, "pinned_libs_*", "system_*"))
	for _, pin := range pins {
		data, err := os.ReadFile(pin)
		if err != nil {
			return true
		}
		lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
		actual, err := filepath.EvalSymlinks(lines[0])
		if err != nil {
			return true
		}
		if fi, err := os.Stat(actual); err != nil || !fi.Mode().IsRegular() || actual != lines[len(lines)-1] {
			return true
		}
	}
	return false
}

// Pin runs the runtime's `setup.sh --force` (steam.sh's pin_newer_runtime_libs) as the Steam user,
// in a private mount namespace over an overlay of the tree whose upper layer is a root-only tmpfs.
// What it wrote reaches the tree only if readPins accepts it: the pin directories, holding symlinks
// and setup.sh's marker files. A symlink can't make a library loadable that its target isn't:
// the trust guard judges the file mapped, never the name it was reached by.
func (r Runtime) Pin(grant Grant) error {
	stage, unmount, err := mountStage()
	if err != nil {
		return err
	}
	defer unmount()
	ops, err := r.pinIn(stage)
	if err != nil || len(ops) == 0 {
		return err
	}
	return grant(func() error { return r.applyPins(ops) })
}

func (r Runtime) pinIn(stage string) ([]pinOp, error) {
	upper, work := filepath.Join(stage, "upper"), filepath.Join(stage, "work")
	// The merged tree's root takes upper's owner and mode: setup.sh, as the user, creates in it.
	if err := os.Mkdir(upper, 0o755); err != nil {
		return nil, err
	}
	if err := os.Chown(upper, int(r.UID), int(r.GID)); err != nil {
		return nil, err
	}
	if err := os.Mkdir(work, 0o700); err != nil {
		return nil, err
	}
	if err := r.runSetup(upper, work); err != nil {
		return nil, err
	}
	ops, err := readPins(upper)
	if err != nil {
		return nil, fmt.Errorf("refusing setup.sh's output: %w", err)
	}
	return ops, nil
}

// mountStage mounts a private tmpfs for the overlay's upper and work dirs: a filesystem overlay
// accepts as upper (not the container's own overlay root), gone once unmounted.
func mountStage() (dir string, unmount func(), err error) {
	if mkErr := os.MkdirAll(stageRoot, 0o700); mkErr != nil {
		return "", nil, mkErr
	}
	dir, err = os.MkdirTemp(stageRoot, "pin-")
	if err != nil {
		return "", nil, err
	}
	mountErr := unix.Mount("tmpfs", dir, "tmpfs", unix.MS_NOSUID|unix.MS_NODEV|unix.MS_NOEXEC, "mode=0700,size=16m")
	if mountErr != nil {
		_ = os.Remove(dir)
		return "", nil, fmt.Errorf("mounting the pin stage: %w", mountErr)
	}
	return dir, func() {
		_ = unix.Unmount(dir, unix.MNT_DETACH)
		_ = os.Remove(dir)
	}, nil
}

// pinScript runs in the new mount namespace: $1 mount, $2 setpriv, $3 the runtime (overlay
// lower and mount point), $4 upper, $5 work, $6/$7 uid/gid, $8 bash.
const pinScript = `set -eu
"$1" -t overlay steamrt-pins -o "lowerdir=$3,upperdir=$4,workdir=$5,redirect_dir=off,metacopy=off" "$3"
cd "$3"
exec "$2" --reuid="$6" --regid="$7" --clear-groups --no-new-privs -- "$8" ./setup.sh --force`

func (r Runtime) runSetup(upper, work string) error {
	rt := filepath.Join(r.Dir, r.Name)
	for _, p := range []string{rt, upper, work} {
		if strings.ContainsAny(p, ",:\\\n") {
			return fmt.Errorf("%q can't be named in overlay mount options", p)
		}
	}
	tools := map[string]string{}
	for _, t := range []string{"unshare", "sh", "mount", "setpriv", "bash"} {
		p, err := systemTool(t)
		if err != nil {
			return err
		}
		tools[t] = p
	}
	ctx, cancel := context.WithTimeout(context.Background(), pinTimeout)
	defer cancel()
	//nolint:gosec // root-placed tools; the paths travel as positional args, never shell text
	cmd := exec.CommandContext(ctx, tools["unshare"], "--mount", "--propagation", "private", "--",
		tools["sh"], "-c", pinScript, "sh", tools["mount"], tools["setpriv"], rt, upper, work,
		strconv.FormatUint(uint64(r.UID), 10), strconv.FormatUint(uint64(r.GID), 10), tools["bash"])
	cmd.Dir = "/"
	// STEAM_ZENITY set and empty: no progress dialog (setup.sh would look for zenity).
	cmd.Env = []string{"PATH=/usr/bin:/bin", "HOME=" + r.Home, "LANG=C.UTF-8", "STEAM_ZENITY="}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pdeathsig: syscall.SIGKILL}
	cmd.Cancel = func() error { return unix.Kill(-cmd.Process.Pid, unix.SIGKILL) }
	var out bytes.Buffer
	cmd.Stdout = &limitedWriter{w: &out, n: 8192}
	cmd.Stderr = cmd.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("setup.sh --force: %w: %s", err, strings.TrimSpace(out.String()))
	}
	return nil
}

// pinOp is one pin directory as setup.sh left it in the overlay's upper layer.
type pinOp struct {
	dir string
	// removed: a whiteout, deleted and not recreated. replace: an opaque directory, recreated, so
	// its entries replace the tree's. Otherwise entries merge into the tree's directory.
	removed, replace bool
	entries          []pinEntry
}

type pinEntry struct {
	name     string
	whiteout bool
	link     string // a symlink's target; else a regular file holding data
	data     []byte
}

// readPins accepts only what setup.sh writes: the pin directories, each holding symlinks, the
// empty done / has_pins markers and system_* records (paths, no NUL). Anything else, a library
// copied in included, refuses the whole run.
func readPins(upper string) ([]pinOp, error) {
	top, err := os.ReadDir(upper)
	if err != nil {
		return nil, err
	}
	var ops []pinOp
	for _, de := range top {
		name := de.Name()
		if !pinDirs[name] {
			return nil, fmt.Errorf("%s is not a pin directory", name)
		}
		p := filepath.Join(upper, name)
		fi, err := os.Lstat(p)
		if err != nil {
			return nil, err
		}
		switch {
		case isWhiteout(fi):
			ops = append(ops, pinOp{dir: name, removed: true})
		case fi.IsDir():
			entries, err := readPinDir(p)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", name, err)
			}
			ops = append(ops, pinOp{dir: name, replace: opaque(p), entries: entries})
		default:
			return nil, fmt.Errorf("%s is not a directory", name)
		}
	}
	return ops, nil
}

func readPinDir(dir string) ([]pinEntry, error) {
	des, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	if len(des) > maxPinEntries {
		return nil, fmt.Errorf("more than %d entries", maxPinEntries)
	}
	entries := make([]pinEntry, 0, len(des))
	for _, de := range des {
		name := de.Name()
		p := filepath.Join(dir, name)
		fi, err := os.Lstat(p)
		if err != nil {
			return nil, err
		}
		switch {
		case isWhiteout(fi):
			entries = append(entries, pinEntry{name: name, whiteout: true})
		case fi.Mode()&os.ModeSymlink != 0:
			target, err := os.Readlink(p)
			if err != nil {
				return nil, err
			}
			entries = append(entries, pinEntry{name: name, link: target})
		case fi.Mode().IsRegular():
			data, err := pinRecord(p, name, fi.Size())
			if err != nil {
				return nil, err
			}
			entries = append(entries, pinEntry{name: name, data: data})
		default:
			return nil, fmt.Errorf("%s: unexpected file type %v", name, fi.Mode().Type())
		}
	}
	return entries, nil
}

func pinRecord(p, name string, size int64) ([]byte, error) {
	switch {
	case name == "done" || name == "has_pins":
		if size != 0 {
			return nil, fmt.Errorf("%s is not an empty marker", name)
		}
		return nil, nil
	case strings.HasPrefix(name, "system_") && size <= maxPinFile:
		data, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}
		if bytes.IndexByte(data, 0) >= 0 {
			return nil, fmt.Errorf("%s is not a text record", name)
		}
		return data, nil
	default:
		return nil, fmt.Errorf("%s is a regular file setup.sh does not write", name)
	}
}

// isWhiteout: overlayfs records a deletion as a 0:0 character device.
func isWhiteout(fi os.FileInfo) bool {
	st, ok := fi.Sys().(*syscall.Stat_t)
	return ok && fi.Mode()&os.ModeCharDevice != 0 && st.Rdev == 0
}

func opaque(dir string) bool {
	buf := make([]byte, 2)
	n, err := unix.Lgetxattr(dir, "trusted.overlay.opaque", buf)
	return err == nil && string(buf[:n]) == "y"
}

// applyPins writes ops into the runtime tree, owned by the Steam user (inside the grant).
func (r Runtime) applyPins(ops []pinOp) error {
	dir, err := openDir(r.Dir)
	if err != nil {
		return err
	}
	defer unix.Close(dir)
	rt, err := openBeneath(dir, r.Name, unix.O_RDONLY|unix.O_DIRECTORY, 0)
	if err != nil {
		return err
	}
	defer unix.Close(rt)
	for _, op := range ops {
		if op.removed || op.replace {
			if err := removeAll(rt, op.dir); err != nil {
				return err
			}
		}
		if op.removed {
			continue
		}
		if err := r.mkdirOwned(rt, op.dir, 0o755); err != nil {
			return err
		}
		if err := r.applyPinDir(rt, op); err != nil {
			return fmt.Errorf("%s: %w", op.dir, err)
		}
	}
	return nil
}

func (r Runtime) applyPinDir(rt int, op pinOp) error {
	d, err := openBeneath(rt, op.dir, unix.O_RDONLY|unix.O_DIRECTORY, 0)
	if err != nil {
		return err
	}
	defer unix.Close(d)
	for _, e := range op.entries {
		if err := removeAll(d, e.name); err != nil {
			return err
		}
		switch {
		case e.whiteout:
			continue
		case e.link != "":
			if err := unix.Symlinkat(e.link, d, e.name); err != nil {
				return err
			}
		default:
			if err := writeNew(d, e.name, e.data); err != nil {
				return err
			}
		}
		if err := unix.Fchownat(d, e.name, int(r.UID), int(r.GID), unix.AT_SYMLINK_NOFOLLOW); err != nil {
			return err
		}
	}
	return nil
}

func writeNew(dir int, name string, data []byte) error {
	fd, err := unix.Openat(dir, name, unix.O_WRONLY|unix.O_CREAT|unix.O_EXCL|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0o644)
	if err != nil {
		return err
	}
	f := os.NewFile(uintptr(fd), name)
	_, err = f.Write(data)
	return errors.Join(err, f.Close())
}
