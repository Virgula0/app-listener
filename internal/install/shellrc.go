package install

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/safeio"
)

const (
	rcBegin = "# >>> app-listener ssh-agent >>>"
	rcEnd   = "# <<< app-listener ssh-agent <<<"
)

// rcTarget is one user file that may carry our marked block. Create: made when missing (else only
// existing files are touched). Dedicated: deleted when its block is removed and nothing else is
// left. Setting: keyword whose presence in a non-comment line means the user configured it
// themselves (case-insensitive). Prepend: block goes first (ssh_config scopes a trailing
// directive to the last Host/Match, and the first value wins).
type rcTarget struct {
	Path      string
	Block     string
	Setting   string
	Mode      os.FileMode
	Create    bool
	Dedicated bool
	Prepend   bool
}

func rcBlock(lines ...string) string {
	return rcBegin + "\n" + strings.Join(lines, "\n") + "\n" + rcEnd + "\n"
}

// rcTargets lists every startup file the installer may edit for u. The socket path mirrors
// ssh-agent.service (%t/ssh-agent.socket = $XDG_RUNTIME_DIR); an already-set SSH_AUTH_SOCK
// (e.g. a forwarded agent) is never overridden.
func rcTargets(u User) []rcTarget {
	posix := rcBlock(fmt.Sprintf(
		`export SSH_AUTH_SOCK="${SSH_AUTH_SOCK:-${XDG_RUNTIME_DIR:-/run/user/%d}/ssh-agent.socket}"`, u.UID))
	fish := rcBlock(fmt.Sprintf(
		`set -q SSH_AUTH_SOCK; or set -gx SSH_AUTH_SOCK (set -q XDG_RUNTIME_DIR; and echo $XDG_RUNTIME_DIR; or echo /run/user/%d)/ssh-agent.socket`,
		u.UID))
	sh := filepath.Base(u.Shell)
	return []rcTarget{
		{Path: filepath.Join(u.Home, ".zshrc"), Block: posix, Setting: authSock, Mode: 0o644, Create: sh == "zsh"},
		{Path: filepath.Join(u.Home, ".bashrc"), Block: posix, Setting: authSock, Mode: 0o644, Create: sh == "bash"},
		{Path: filepath.Join(u.Home, ".config", "fish", "conf.d", "app-listener-ssh-agent.fish"),
			Block: fish, Setting: authSock, Mode: 0o644, Create: sh == "fish", Dedicated: true},
	}
}

const authSock = "SSH_AUTH_SOCK"

// sshConfigTarget is ~/.ssh/config with AddKeysToAgent, so ssh (whitelisted for ~/.ssh) loads a
// key into the agent on first use; the daemon itself never touches the agent.
func sshConfigTarget(u User) rcTarget {
	return rcTarget{
		Path:      filepath.Join(u.Home, ".ssh", "config"),
		Block:     rcBlock("AddKeysToAgent yes"),
		Setting:   "AddKeysToAgent",
		Mode:      0o600,
		Create:    true,
		Dedicated: true,
		Prepend:   true,
	}
}

// EnsureAddKeysToAgent puts `AddKeysToAgent yes` at the top of u's ~/.ssh/config (created when
// missing). A config that already sets AddKeysToAgent is left alone. Run before ~/.ssh is
// encrypted so the file is part of the migrated tree. Reports whether the file was modified.
func EnsureAddKeysToAgent(u User) (path string, changed bool, err error) {
	t := sshConfigTarget(u)
	if _, statErr := os.Stat(filepath.Dir(t.Path)); os.IsNotExist(statErr) {
		return t.Path, false, nil // never create ~/.ssh here (needs 0700 and its own guard)
	}
	changed, err = ensureRCBlock(t, u)
	if err != nil {
		return t.Path, false, fmt.Errorf("%s: %w", t.Path, err)
	}
	return t.Path, changed, nil
}

// EnsureSSHAgentEnv appends the SSH_AUTH_SOCK block to u's shell startup files (login shell's
// always, others only when present). Idempotent; a file that already sets SSH_AUTH_SOCK itself is
// left alone. Returns the files modified.
func EnsureSSHAgentEnv(u User) ([]string, error) {
	var changed []string
	for _, t := range rcTargets(u) {
		ok, err := ensureRCBlock(t, u)
		if err != nil {
			return changed, fmt.Errorf("%s: %w", t.Path, err)
		}
		if ok {
			changed = append(changed, t.Path)
		}
	}
	return changed, nil
}

// RemoveSSHAgentEnv strips the blocks EnsureSSHAgentEnv and EnsureAddKeysToAgent added from every
// startup file and ~/.ssh/config of u. Returns the files modified.
func RemoveSSHAgentEnv(u User) ([]string, error) {
	var changed []string
	for _, t := range append(rcTargets(u), sshConfigTarget(u)) {
		ok, err := removeRCBlock(t, u.UID)
		if err != nil {
			return changed, fmt.Errorf("%s: %w", t.Path, err)
		}
		if ok {
			changed = append(changed, t.Path)
		}
	}
	return changed, nil
}

func ensureRCBlock(t rcTarget, u User) (bool, error) {
	f, err := openOwnedRegular(t.Path, u.UID)
	if errors.Is(err, fs.ErrNotExist) {
		if !t.Create {
			return false, nil
		}
		f, err = createOwned(u.Home, t.Path, u.UID, u.GID, t.Mode)
	}
	if err != nil {
		return false, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return false, err
	}
	text := string(data)
	if strings.Contains(text, rcBegin) {
		return false, nil
	}
	if setsKeyword(text, t.Setting) {
		log.Infof("%s already sets %s: leaving it alone", t.Path, t.Setting)
		return false, nil
	}
	if t.Prepend {
		// Rewrites the whole (small) file in place, longer than before, so no truncate is needed.
		merged := t.Block
		if len(data) > 0 {
			merged += "\n" + text
		}
		_, err = f.WriteAt([]byte(merged), 0)
		return err == nil, err
	}
	add := t.Block
	if len(data) > 0 {
		add = "\n" + add
		if !strings.HasSuffix(text, "\n") {
			add = "\n" + add
		}
	}
	_, err = f.WriteString(add)
	return err == nil, err
}

func removeRCBlock(t rcTarget, uid uint32) (bool, error) {
	f, err := openOwnedRegular(t.Path, uid)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return false, err
	}
	orig := string(data)
	text, err := stripRCBlocks(orig, t.Prepend)
	if err != nil {
		return false, err
	}
	if text == orig {
		return false, nil
	}
	if t.Dedicated && strings.TrimSpace(text) == "" {
		return true, os.Remove(t.Path)
	}
	// Shorter than before: overwrite in place first, truncate last, so a crash never leaves an
	// empty rc file.
	if _, err := f.WriteAt([]byte(text), 0); err != nil {
		return false, err
	}
	return true, f.Truncate(int64(len(text)))
}

// stripRCBlocks removes every marked block from text together with the blank separator that was
// added alongside it.
func stripRCBlocks(text string, prepend bool) (string, error) {
	for {
		start := strings.Index(text, rcBegin)
		if start < 0 {
			return text, nil
		}
		n := strings.Index(text[start:], rcEnd)
		if n < 0 {
			return "", errors.New("unterminated app-listener block: remove it by hand")
		}
		end := start + n + len(rcEnd)
		if strings.HasPrefix(text[end:], "\n") {
			end++
		}
		switch {
		case prepend && start == 0 && strings.HasPrefix(text[end:], "\n"):
			end++
		case !prepend && strings.HasSuffix(text[:start], "\n\n"):
			start--
		}
		text = text[:start] + text[end:]
	}
}

// setsKeyword reports whether any non-comment line mentions keyword (case-insensitive).
func setsKeyword(text, keyword string) bool {
	keyword = strings.ToLower(keyword)
	for _, line := range strings.Split(text, "\n") {
		line = strings.ToLower(strings.TrimSpace(line))
		if !strings.HasPrefix(line, "#") && strings.Contains(line, keyword) {
			return true
		}
	}
	return false
}

// openOwnedRegular opens path read-write and fstat-checks the descriptor: root edits a
// user-writable file, so a symlink swapped in to a root-owned target must be refused. Symlinks to
// the user's own files (dotfile managers) still work.
func openOwnedRegular(path string, uid uint32) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, err
	}
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !fi.Mode().IsRegular() || !ok || st.Uid != uid {
		f.Close()
		return nil, fmt.Errorf("refusing to edit: not a regular file owned by uid %d", uid)
	}
	return f, nil
}

// createOwned creates path (owned by uid/gid, missing parents made 0755) as root inside the user's
// home. The descent is symlink-safe: home is the trusted anchor, and every component below it is
// opened O_NOFOLLOW, so a parent directory the user has replaced with a symlink (e.g.
// ~/.config/fish/conf.d -> /etc/systemd/system) is refused instead of letting root create and chown
// a file outside the home. The final component is created O_EXCL|O_NOFOLLOW.
func createOwned(home, path string, uid, gid uint32, mode os.FileMode) (*os.File, error) {
	rel, err := filepath.Rel(home, filepath.Clean(path))
	if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
		return nil, fmt.Errorf("%s is not inside %s", path, home)
	}
	comps := strings.Split(rel, string(os.PathSeparator))

	dirfd, err := unix.Open(home, unix.O_DIRECTORY|unix.O_CLOEXEC, 0)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", home, err)
	}
	defer func() {
		if dirfd >= 0 {
			_ = unix.Close(dirfd)
		}
	}()

	dirfd, err = safeio.DescendCreateOwned(dirfd, comps[:len(comps)-1], int(uid), int(gid))
	if err != nil {
		return nil, err
	}
	return safeio.CreateExclNoFollow(dirfd, comps[len(comps)-1], mode, int(uid), int(gid))
}
