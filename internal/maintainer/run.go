// Package maintainer runs an app's own maintenance script (Steam's steam.sh) for the read-only
// lib_dirs it can't write itself, and commits what it wrote only after validating it. The script
// runs as its user in private mount, pid, ipc and network namespaces, over an overlay of its own
// directory whose writable layer is root-only, with the app's binaries masked by true so the app
// itself never starts. Only changes below the given lib_dirs are kept, and of those only what
// plan.go accepts: content the lib_dir already holds, as a file or inside an archive, plus
// symlinks, directories, deletions and small path records.
package maintainer

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

const (
	runTimeout = 15 * time.Minute
	// maxScript bounds the copy of the maintainer script (steam.sh is ~40 KiB); a longer one fails
	// the hash check.
	maxScript = 16 << 20
	// stageRoot holds each run's overlay layers: on disk, not tmpfs, as an unpacked runtime and the
	// copy of the tree it replaces take hundreds of MB.
	stageRoot = "/var/lib/app-listener/maintain"
)

// Spec is one maintainer run.
type Spec struct {
	// Script is the vetted script (any fd to it, O_PATH included) and ScriptSum the content the
	// ledger confirmed. The sandbox runs a copy holding exactly that at ScriptPath, whatever the
	// path names by then.
	Script     *os.File
	ScriptSum  [32]byte
	ScriptPath string
	UID, GID   uint32
	User, Home string
	// Masks are binaries below the script's directory that the sandbox replaces with true.
	Masks []string
	// LibDirs are the directories below the script's directory whose changes may be committed.
	LibDirs []string
}

// CleanStages removes what a run interrupted by a crash left in stageRoot.
func CleanStages() error { return os.RemoveAll(stageRoot) }

// Run runs spec's script in the sandbox and returns its validated changes. The caller closes the
// plan.
func Run(spec *Spec) (*Plan, error) {
	root := filepath.Dir(spec.ScriptPath)
	for _, lib := range spec.LibDirs {
		if !strings.HasPrefix(lib, root+"/") {
			return nil, fmt.Errorf("%s is not below %s", lib, root)
		}
	}
	st, err := newStage(root)
	if err != nil {
		return nil, err
	}
	if err = st.copyScript(spec); err == nil {
		err = runSandbox(spec, root, st)
	}
	if err != nil {
		st.close()
		return nil, err
	}
	p, err := planFrom(spec, root, st)
	if err != nil {
		st.close()
		return nil, err
	}
	return p, nil
}

type stage struct{ dir, upper, work, merged, script string }

func (s *stage) close() { _ = os.RemoveAll(s.dir) }

// newStage makes the run's root-only stage. The overlay's root takes upper's owner and mode, so
// upper copies root's.
func newStage(root string) (*stage, error) {
	if err := os.MkdirAll(stageRoot, 0o700); err != nil {
		return nil, err
	}
	var st unix.Stat_t
	if err := unix.Lstat(stageRoot, &st); err != nil {
		return nil, err
	}
	if st.Mode&unix.S_IFMT != unix.S_IFDIR || st.Uid != 0 || st.Mode&0o077 != 0 {
		return nil, fmt.Errorf("%s must be a root-only directory", stageRoot)
	}
	var sfs unix.Statfs_t
	if err := unix.Statfs(stageRoot, &sfs); err != nil {
		return nil, err
	}
	if sfs.Type == unix.OVERLAYFS_SUPER_MAGIC {
		return nil, fmt.Errorf("%s is on overlayfs, which can't hold an overlay's writable layer", stageRoot)
	}
	var rst unix.Stat_t
	if err := unix.Stat(root, &rst); err != nil {
		return nil, err
	}
	dir, err := os.MkdirTemp(stageRoot, "run-")
	if err != nil {
		return nil, err
	}
	s := &stage{dir: dir, upper: filepath.Join(dir, "upper"), work: filepath.Join(dir, "work"),
		merged: filepath.Join(dir, "merged"), script: filepath.Join(dir, "script")}
	err = os.Mkdir(s.upper, 0o700)
	if err == nil {
		err = os.Chown(s.upper, int(rst.Uid), int(rst.Gid))
	}
	if err == nil {
		err = os.Chmod(s.upper, os.FileMode(rst.Mode&0o777))
	}
	if err == nil {
		err = os.Mkdir(s.work, 0o700)
	}
	if err == nil {
		err = os.Mkdir(s.merged, 0o700)
	}
	if err != nil {
		s.close()
		return nil, err
	}
	return s, nil
}

// copyScript writes the vetted script into the stage, refusing content other than what the ledger
// confirmed. The sandbox binds this copy: a bind's source must be a mount of the caller's own
// namespace, and the vetted fd's is the daemon's.
func (s *stage) copyScript(spec *Spec) error {
	src, err := os.Open(fmt.Sprintf("/proc/self/fd/%d", spec.Script.Fd()))
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.OpenFile(s.script, os.O_WRONLY|os.O_CREATE|os.O_EXCL|syscall.O_NOFOLLOW, 0o755)
	if err != nil {
		return err
	}
	defer dst.Close()
	h := sha256.New()
	if _, err := io.Copy(io.MultiWriter(dst, h), io.LimitReader(src, maxScript)); err != nil {
		return err
	}
	if !bytes.Equal(h.Sum(nil), spec.ScriptSum[:]) {
		return fmt.Errorf("%s changed after the binary ledger confirmed it", spec.ScriptPath)
	}
	return dst.Close()
}

// sandboxScript runs as root in the new namespaces, then becomes the user. The script sees a
// private, empty home (and /run/user/<uid>, /tmp) holding only its own directory, as the overlay:
// whatever it writes elsewhere is discarded, and it reaches none of the user's files or sockets.
// The script is the stage's checked copy (copyScript), bound read-only over its path.
const sandboxScript = `set -eu
mount=$1 setpriv=$2 bash=$3 true=$4 chown=$5 mkdir=$6 root=$7 upper=$8 work=$9
shift 9
merged=$1 script=$2 uid=$3 gid=$4 home=$5 copy=$6
shift 6
"$mount" -t overlay app-listener-maintain -o "lowerdir=$root,upperdir=$upper,workdir=$work,redirect_dir=off,metacopy=off,index=off,xino=off" "$merged"
"$mount" -t tmpfs -o mode=1777,nosuid,nodev app-listener-tmp /tmp
if [ -d "/run/user/$uid" ]; then
	"$mount" -t tmpfs -o "mode=0700,uid=$uid,gid=$gid,nosuid,nodev" app-listener-run "/run/user/$uid"
fi
if [ -n "$home" ]; then
	"$mount" -t tmpfs -o "mode=0700,uid=$uid,gid=$gid,nosuid,nodev" app-listener-home "$home"
	"$mkdir" -p "$root"
	"$chown" -R "$uid:$gid" "$home"
fi
"$mount" --move "$merged" "$root"
"$mount" -c --bind "$copy" "$root/$script"
"$mount" -c -o remount,bind,ro "$root/$script"
for m in "$@"; do
	if [ -f "$m" ]; then "$mount" --bind "$true" "$m"; fi
done
cd "$root"
exec "$setpriv" --reuid="$uid" --regid="$gid" --clear-groups --no-new-privs -- "$bash" "$root/$script"`

func runSandbox(spec *Spec, root string, st *stage) error {
	for _, p := range []string{root, st.upper, st.work, st.merged} {
		if strings.ContainsAny(p, ",:\\\n") {
			return fmt.Errorf("%q can't be named in overlay mount options", p)
		}
	}
	tools := map[string]string{}
	for _, t := range []string{"unshare", "sh", "mount", "setpriv", "bash", "true", "chown", "mkdir"} {
		p, err := systemTool(t)
		if err != nil {
			return err
		}
		tools[t] = p
	}
	args := []string{"--mount", "--pid", "--ipc", "--net", "--fork", "--kill-child", "--mount-proc",
		"--propagation", "private", "--", tools["sh"], "-c", sandboxScript, "sh",
		tools["mount"], tools["setpriv"], tools["bash"], tools["true"], tools["chown"], tools["mkdir"],
		root, st.upper, st.work, st.merged, filepath.Base(spec.ScriptPath),
		strconv.FormatUint(uint64(spec.UID), 10), strconv.FormatUint(uint64(spec.GID), 10), "", st.script}
	if home := filepath.Clean(spec.Home); home != "/" && strings.HasPrefix(root, home+"/") {
		args[len(args)-2] = home
	}
	for _, m := range spec.Masks {
		if strings.HasPrefix(m, root+"/") {
			args = append(args, m)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), runTimeout)
	defer cancel()
	//nolint:gosec // root-placed tools; every path travels as a positional arg, never shell text
	cmd := exec.CommandContext(ctx, tools["unshare"], args...)
	cmd.Dir = "/"
	cmd.Env = []string{"PATH=/usr/bin:/bin", "HOME=" + spec.Home, "USER=" + spec.User, "LOGNAME=" + spec.User,
		"LANG=C.UTF-8"}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pdeathsig: syscall.SIGKILL}
	cmd.Cancel = func() error { return unix.Kill(-cmd.Process.Pid, unix.SIGKILL) }
	var out bytes.Buffer
	cmd.Stdout = &limitedWriter{w: &out, n: 8192}
	cmd.Stderr = cmd.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s: %w: %s", spec.ScriptPath, err, tail(out.String()))
	}
	return nil
}

func tail(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 2048 {
		s = s[len(s)-2048:]
	}
	return s
}
