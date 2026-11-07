package daemon

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/systemd"
	"github.com/Virgula0/app-listener/internal/usecase"
)

// Runtime sections (`runtime: true`) guard a directory in a user's /run/user/<uid>, a tmpfs that
// user-runtime-dir@<uid>.service mounts at the user's first login and unmounts at the last logout.
// Its drop-in (systemd.LoginDropIn) runs `daemon login <uid>` after the mount, before user@<uid>
// starts any of the user's processes, and `daemon logout <uid>` after the unmount:
//
//	client -> LOGIN <uid>\n | LOGOUT <uid>\n
//	server -> OK\n | ERR <reason>\n
//
// LOGIN makes each missing runtime directory as root 0700 (the user can't enter it), plants the
// sentinel, reloads so a guard anchors on it, checks that guard holds the very inode made here,
// and only then chowns it to the user. Meanwhile /run/user/<uid> is root's, so none of the user's
// processes can rename, replace or add an entry there; a directory this daemon didn't make (one of
// the user's, maybe holding listeners, even if a reload guarded it as found) is first moved into a
// root-only directory. A failure fails the login (the drop-in's ExecStartPost).
const (
	loginSocket = "/run/app-listener-daemon.login"
	loginLogin  = "LOGIN"
	loginLogout = "LOGOUT"
	loginOK     = "OK"
	loginErr    = "ERR"
	// loginTimeout bounds one request: a reload rebuilds every guard. Below user-runtime-dir@'s
	// default 90s start timeout, so the client reports the failure itself.
	loginTimeout = 75 * time.Second
	// runtimeSentinel is a root-owned, non-empty directory the daemon plants in every runtime
	// directory it makes: no whitelisted DELETE binary (rm -d, an agent) can then remove the root,
	// which would let the user recreate it unguarded until the next reload.
	runtimeSentinel = ".app-listener"
	sentinelNote    = "Keeps this directory non-empty: app-listener guards it, and an empty guarded root could be removed.\n"
)

var loginCmd = &cobra.Command{
	Use:   "login <uid>",
	Short: "Make and guard a user's runtime sections (run by user-runtime-dir@'s drop-in)",
	Long: `Ask the running daemon to make and guard the runtime sections (runtime: true)
of /run/user/<uid>, before that user's systemd manager starts. Exits 0 when no
daemon is running; non-zero when the daemon refused or timed out, which fails
the login. The installer's drop-in on user-runtime-dir@.service runs it.`,
	Args: cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error { return runLoginClient(loginLogin, args[0]) },
}

var logoutCmd = &cobra.Command{
	Use:   "logout <uid>",
	Short: "Drop a logged-out user's runtime sections (run by user-runtime-dir@'s drop-in)",
	Args:  cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		if err := runLoginClient(loginLogout, args[0]); err != nil {
			log.Warnf("daemon logout: %v", err) // never fail a stop
		}
	},
}

func init() {
	DaemonCmd.AddCommand(loginCmd, logoutCmd)
}

func runLoginClient(verb, uidText string) error {
	uid, err := strconv.ParseUint(uidText, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid uid %q", uidText)
	}
	conn, err := (&net.Dialer{Timeout: 5 * time.Second}).DialContext(context.Background(), "unix", loginSocket)
	if errors.Is(err, syscall.ENOENT) || errors.Is(err, syscall.ECONNREFUSED) {
		log.Infof("no app-listener daemon is running: nothing to guard in /run/user/%d", uid)
		return nil
	}
	if err != nil {
		return fmt.Errorf("connecting to the daemon: %w", err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(loginTimeout + 5*time.Second))
	if _, werr := fmt.Fprintf(conn, "%s %d\n", verb, uid); werr != nil {
		return fmt.Errorf("sending %s: %w", verb, werr)
	}
	line, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return fmt.Errorf("waiting for the daemon: %w", err)
	}
	line = strings.TrimRight(line, "\n")
	if line != loginOK {
		return fmt.Errorf("the daemon refused /run/user/%d: %s", uid, strings.TrimPrefix(line, loginErr+" "))
	}
	return nil
}

// loginServer answers the drop-in's requests, one at a time. Requests wait until the reload exists.
type loginServer struct {
	ln       net.Listener
	selfDev  uint64
	selfIno  uint64
	uc       usecase.DaemonUseCase
	mu       sync.Mutex
	ready    chan struct{}
	readyOne sync.Once
	reload   func() (*daemonconfig.Config, error)
	cfg      atomic.Pointer[daemonconfig.Config]
	secrets  *secretGates
	// made: the runtime directories this daemon made and handed to their uid (under mu). Only these
	// are admitted as found at a login; dropped at that uid's logout. A guard holds each one open, so
	// its tmpfs, and with it the device number in the key, can't be reused while it is guarded.
	made map[guard.GuardInodeKey]uint32
}

// startLoginServer opens the socket before READY=1: the drop-in orders after the daemon unit, so a
// login that finds no socket means no daemon.
func startLoginServer(uc usecase.DaemonUseCase, cfg *daemonconfig.Config) (*loginServer, error) {
	selfDev, selfIno, err := ebpf.StatInode("/proc/self/exe")
	if err != nil {
		return nil, fmt.Errorf("resolving own executable for the login socket: %w", err)
	}
	if rmErr := os.Remove(loginSocket); rmErr != nil && !os.IsNotExist(rmErr) {
		return nil, fmt.Errorf("removing stale login socket %s: %w", loginSocket, rmErr)
	}
	old := unix.Umask(0o177)
	ln, err := (&net.ListenConfig{}).Listen(context.Background(), "unix", loginSocket)
	unix.Umask(old)
	if err != nil {
		return nil, fmt.Errorf("listening on login socket %s: %w", loginSocket, err)
	}
	if len(cfg.RuntimePaths) > 0 && !systemd.LoginDropInInstalled() {
		log.Errorf("daemon: %s is missing: nothing guards the runtime sections at login — run `sudo app-listener install`",
			systemd.LoginDropInFile)
	}
	s := &loginServer{ln: ln, selfDev: selfDev, selfIno: selfIno, uc: uc, ready: make(chan struct{}),
		secrets: newSecretGates(cfg), made: map[guard.GuardInodeKey]uint32{}}
	s.cfg.Store(cfg)
	s.sentinelFound()
	go s.acceptLoop()
	return s, nil
}

// setReload hands over the reload; setConfig tracks the config it last committed (the reload
// handler's reloaded callback).
func (s *loginServer) setReload(reload func() (*daemonconfig.Config, error)) {
	s.reload = reload
	s.readyOne.Do(func() { close(s.ready) })
}

func (s *loginServer) setConfig(cfg *daemonconfig.Config) {
	if s != nil && cfg != nil {
		s.cfg.Store(cfg)
		s.secrets.setConfig(cfg)
		s.sentinelFound()
	}
}

// sentinelFound plants the sentinel in each guarded runtime directory without one: those found at a
// start or reload (made by the user before the daemon, as an agent unit's `mkdir -p` does) are
// guarded as found, and an agent's own cleanup could leave one empty for an rmdir.
func (s *loginServer) sentinelFound() {
	resources := s.uc.Resources()
	for i := range resources {
		r := &resources[i]
		uid, ok := daemonconfig.RuntimeDirFor(r.Path)
		if !r.Runtime || !ok {
			continue
		}
		if err := s.plantFound(uid, r.Path); err != nil {
			log.Errorf("daemon: %s: planting %s: %v — an emptied root could be removed and made anew unguarded",
				r.Path, runtimeSentinel, err)
		}
	}
}

func (s *loginServer) plantFound(uid uint32, path string) error {
	userFd, _, err := openUserRuntimeDir(uid)
	if err != nil {
		return nil // logged out: the reload drops the section
	}
	defer unix.Close(userFd)
	fd, err := unix.Openat(userFd, path[strings.LastIndexByte(path, '/')+1:],
		unix.O_RDONLY|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if errors.Is(err, unix.ENOENT) {
		return nil
	}
	if err != nil {
		return err
	}
	defer unix.Close(fd)
	if hasSentinel(fd) {
		return nil
	}
	var st unix.Stat_t
	if err := unix.Fstat(fd, &st); err != nil {
		return err
	}
	key := guard.GuardInodeKey{Dev: ebpf.KernelDev(unix.Major(st.Dev), unix.Minor(st.Dev)), Ino: st.Ino}
	return s.uc.WithRuntimeRoot(path, func(root guard.GuardInodeKey) error {
		if root != key {
			return errors.New("another directory is guarded at that path")
		}
		return plantSentinel(fd)
	})
}

func (s *loginServer) close() {
	if s == nil {
		return
	}
	_ = s.ln.Close()
	_ = os.Remove(loginSocket)
	s.secrets.close()
}

func (s *loginServer) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return
		}
		go s.handle(conn)
	}
}

func (s *loginServer) handle(conn net.Conn) {
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(loginTimeout))
	if err := authSelfPeer(conn, s.selfDev, s.selfIno); err != nil {
		log.Warnf("daemon: login: rejected connection: %v", err)
		_, _ = fmt.Fprintf(conn, "%s %s\n", loginErr, err)
		return
	}
	line, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return
	}
	verb, uidText, _ := strings.Cut(strings.TrimRight(line, "\n"), " ")
	uid, parseErr := strconv.ParseUint(uidText, 10, 32)
	if parseErr != nil || (verb != loginLogin && verb != loginLogout) {
		_, _ = fmt.Fprintf(conn, "%s malformed request\n", loginErr)
		return
	}
	select {
	case <-s.ready:
	case <-time.After(loginTimeout / 2):
		_, _ = fmt.Fprintf(conn, "%s the daemon is still starting\n", loginErr)
		return
	}
	s.mu.Lock()
	if verb == loginLogin {
		err = s.login(uint32(uid))
	} else {
		err = s.logout(uint32(uid))
	}
	s.mu.Unlock()
	if err != nil {
		log.Errorf("daemon: %s %d: %v", strings.ToLower(verb), uid, err)
		_, _ = fmt.Fprintf(conn, "%s %s\n", loginErr, strings.ReplaceAll(err.Error(), "\n", " "))
		return
	}
	_, _ = fmt.Fprintf(conn, "%s\n", loginOK)
}

// runtimePathsFor returns cfg's runtime sections in /run/user/<uid>.
func runtimePathsFor(cfg *daemonconfig.Config, uid uint32) []string {
	var out []string
	for _, p := range cfg.RuntimePaths {
		if owner, ok := daemonconfig.RuntimeDirFor(p); ok && owner == uid {
			out = append(out, p)
		}
	}
	return out
}

// madeDir is a runtime directory this login created, held open: the chown goes through this fd,
// after the guard is seen anchored on its inode.
type madeDir struct {
	path string
	fd   int
	key  guard.GuardInodeKey
}

func (s *loginServer) login(uid uint32) error {
	// Before user@<uid> starts: the keyring's first connect, and every bus connect, find it ready.
	s.secrets.login(uid)
	paths := runtimePathsFor(s.cfg.Load(), uid)
	if len(paths) == 0 {
		return nil
	}
	userFd, gid, err := openUserRuntimeDir(uid)
	if err != nil {
		return err
	}
	defer unix.Close(userFd)
	if s.allOurs(userFd, uid, paths) {
		return nil
	}
	unlock, err := lockParent(userFd, uid)
	if err != nil {
		return err
	}
	defer unlock()
	made, err := s.prepareDirs(userFd, uid, paths)
	defer func() {
		for _, m := range made {
			unix.Close(m.fd)
		}
	}()
	if err != nil {
		return err
	}
	if len(made) == 0 && s.allRunning(paths) {
		return nil
	}
	if _, err := s.reload(); err != nil {
		return fmt.Errorf("reload: %w", err)
	}
	if missing := s.notRunning(paths); len(missing) > 0 {
		return fmt.Errorf("%v not guarded after the reload", missing)
	}
	for _, m := range made {
		if chErr := s.uc.WithRuntimeRoot(m.path, func(root guard.GuardInodeKey) error {
			if root != m.key {
				return fmt.Errorf("%s was replaced before its guard attached", m.path)
			}
			return unix.Fchown(m.fd, int(uid), int(gid))
		}); chErr != nil {
			return fmt.Errorf("handing %s to uid %d: %w", m.path, uid, chErr)
		}
		s.made[m.key] = uid
		log.Infof("daemon: %s made and guarded for uid %d", m.path, uid)
	}
	return nil
}

// lockParent makes /run/user/<uid> root's until unlock (its group, and mode 0700, stay): none of
// the user's processes can then add, rename or remove an entry in it.
func lockParent(userFd int, uid uint32) (unlock func(), err error) {
	if err := unix.Fchownat(userFd, "", 0, -1, unix.AT_EMPTY_PATH); err != nil {
		return nil, fmt.Errorf("locking /run/user/%d: %w", uid, err)
	}
	return func() {
		if err := unix.Fchownat(userFd, "", int(uid), -1, unix.AT_EMPTY_PATH); err != nil {
			log.Errorf("daemon: handing /run/user/%d back to its user: %v — the next login retries", uid, err)
		}
	}, nil
}

// logout reloads when one of uid's runtime directories went with its mount, so the reload drops it.
func (s *loginServer) logout(uid uint32) error {
	s.secrets.logout(uid)
	for k, owner := range s.made {
		if owner == uid {
			delete(s.made, k)
		}
	}
	resources := s.uc.Resources()
	for i := range resources {
		r := &resources[i]
		if owner, ok := daemonconfig.RuntimeDirFor(r.Path); !r.Runtime || !ok || owner != uid {
			continue
		}
		if _, err := os.Lstat(r.Path); errors.Is(err, os.ErrNotExist) {
			if _, reloadErr := s.reload(); reloadErr != nil {
				return fmt.Errorf("reload: %w", reloadErr)
			}
			return nil
		}
	}
	return nil
}

// prepareDirs makes each missing runtime directory and moves aside each one found that this daemon
// didn't make; made holds the ones to hand over, open, even on error. The parent is locked.
func (s *loginServer) prepareDirs(userFd int, uid uint32, paths []string) (made []madeDir, err error) {
	for _, p := range paths {
		m, created, mkErr := makeRuntimeDir(userFd, p)
		if mkErr != nil {
			return made, mkErr
		}
		if created {
			made = append(made, m)
			continue
		}
		if s.ours(m, uid) {
			unix.Close(m.fd)
			continue
		}
		if err := s.moveAside(userFd, m); err != nil {
			return made, err
		}
		if m, created, mkErr = makeRuntimeDir(userFd, p); mkErr != nil {
			return made, mkErr
		}
		if !created {
			unix.Close(m.fd)
			return made, fmt.Errorf("%s: made anew while /run/user/%d was locked", p, uid)
		}
		made = append(made, m)
	}
	return made, nil
}

// allOurs: every path is a directory this daemon made for uid, anchored by its running guard
// (user-runtime-dir@ restarted over a live mount). Unchanged since: the user can't rename a root.
func (s *loginServer) allOurs(userFd int, uid uint32, paths []string) bool {
	for _, p := range paths {
		var st unix.Stat_t
		if unix.Fstatat(userFd, p[strings.LastIndexByte(p, '/')+1:], &st, unix.AT_SYMLINK_NOFOLLOW) != nil {
			return false
		}
		key := guard.GuardInodeKey{Dev: ebpf.KernelDev(unix.Major(st.Dev), unix.Minor(st.Dev)), Ino: st.Ino}
		if !s.ours(madeDir{path: p, key: key}, uid) {
			return false
		}
	}
	return true
}

// ours: m is a directory this daemon made for uid and its running guard still anchors it.
func (s *loginServer) ours(m madeDir, uid uint32) bool {
	if owner, ok := s.made[m.key]; !ok || owner != uid {
		return false
	}
	// WithRuntimeRoot only reads the anchor here; nothing is changed.
	return s.uc.WithRuntimeRoot(m.path, func(root guard.GuardInodeKey) error {
		if root != m.key {
			return errors.New("another directory is guarded at that path")
		}
		return nil
	}) == nil
}

// moveAside renames the directory m found into a new root-only directory beside it: whatever the
// user put there (listeners bound before any guard) stays out of reach of the path's clients. A
// running guard anchored on it lets the daemon move its root (ownerAccessEvents). Closes m.fd.
func (s *loginServer) moveAside(userFd int, m madeDir) error {
	defer unix.Close(m.fd)
	var rnd [8]byte
	if _, err := rand.Read(rnd[:]); err != nil {
		return err
	}
	name := m.path[strings.LastIndexByte(m.path, '/')+1:]
	aside := ".app-listener-moved-" + hex.EncodeToString(rnd[:])
	if err := unix.Mkdirat(userFd, aside, 0o700); err != nil {
		return fmt.Errorf("moving %s aside: %w", m.path, err)
	}
	asideFd, err := unix.Openat(userFd, aside, unix.O_PATH|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if err != nil {
		return fmt.Errorf("moving %s aside: %w", m.path, err)
	}
	defer unix.Close(asideFd)
	move := func() error { return unix.Renameat2(userFd, name, asideFd, name, unix.RENAME_NOREPLACE) }
	if s.allRunning([]string{m.path}) {
		err = s.uc.WithRuntimeRoot(m.path, func(guard.GuardInodeKey) error { return move() })
	} else {
		err = move()
	}
	if err != nil {
		return fmt.Errorf("%s was not made by the daemon and can't be moved aside: %w", m.path, err)
	}
	log.Warnf("daemon: %s was not made by the daemon: moved to %s/%s, a fresh one is guarded", m.path,
		m.path[:strings.LastIndexByte(m.path, '/')], aside)
	return nil
}

func (s *loginServer) allRunning(paths []string) bool { return len(s.notRunning(paths)) == 0 }

func (s *loginServer) notRunning(paths []string) []string {
	running := make(map[string]bool)
	resources := s.uc.Resources()
	for i := range resources {
		running[resources[i].Path] = resources[i].Runtime
	}
	var out []string
	for _, p := range paths {
		if !running[p] {
			out = append(out, p)
		}
	}
	return out
}

// openUserRuntimeDir opens /run/user/<uid> without following links and checks it is the user's (or
// root's: a login's lock its daemon didn't undo), returning its fd and group.
func openUserRuntimeDir(uid uint32) (fd int, gid uint32, err error) {
	base, err := unix.Open("/run/user", unix.O_PATH|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if err != nil {
		return -1, 0, fmt.Errorf("opening /run/user: %w", err)
	}
	defer unix.Close(base)
	fd, err = unix.Openat(base, strconv.FormatUint(uint64(uid), 10), unix.O_PATH|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if err != nil {
		return -1, 0, fmt.Errorf("opening /run/user/%d: %w", uid, err)
	}
	var st unix.Stat_t
	var fs unix.Statfs_t
	err = unix.Fstat(fd, &st)
	switch {
	case err != nil:
	case st.Uid != uid && st.Uid != 0:
		err = fmt.Errorf("owned by uid %d", st.Uid)
	case unix.Fstatfs(fd, &fs) != nil || fs.Type != unix.TMPFS_MAGIC:
		// user-runtime-dir@'s tmpfs, not a FUSE mount the user put over it (setuid fusermount3).
		err = errors.New("not a tmpfs")
	}
	if err != nil {
		unix.Close(fd)
		return -1, 0, fmt.Errorf("/run/user/%d: %w", uid, err)
	}
	return fd, st.Gid, nil
}

// makeRuntimeDir creates path's directory in userFd as root 0700 with the sentinel and returns it
// open. created is false for a directory already there and not this daemon's unfinished one (root's,
// sentinel planted: a login that failed after the mkdir); the caller then vets it.
func makeRuntimeDir(userFd int, path string) (m madeDir, created bool, err error) {
	name := path[strings.LastIndexByte(path, '/')+1:]
	mkErr := unix.Mkdirat(userFd, name, 0o700)
	existed := errors.Is(mkErr, unix.EEXIST)
	if mkErr != nil && !existed {
		return madeDir{}, false, fmt.Errorf("creating %s: %w", path, mkErr)
	}
	fd, err := unix.Openat(userFd, name, unix.O_RDONLY|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if err != nil {
		return madeDir{}, false, fmt.Errorf("opening %s: %w", path, err)
	}
	var st unix.Stat_t
	if err := unix.Fstat(fd, &st); err != nil {
		unix.Close(fd)
		return madeDir{}, false, fmt.Errorf("%s: %w", path, err)
	}
	m = madeDir{path: path, fd: fd, key: guard.GuardInodeKey{Dev: ebpf.KernelDev(unix.Major(st.Dev), unix.Minor(st.Dev)), Ino: st.Ino}}
	switch {
	case existed && st.Uid == 0 && hasSentinel(fd):
		return m, true, nil
	case existed:
		return m, false, nil
	case st.Uid != 0:
		unix.Close(fd)
		return madeDir{}, false, fmt.Errorf("%s: not the directory just made", path)
	}
	if err := plantSentinel(fd); err != nil {
		unix.Close(fd)
		return madeDir{}, false, fmt.Errorf("%s: %w", path, err)
	}
	return m, true, nil
}

func hasSentinel(dirFd int) bool {
	var st unix.Stat_t
	return unix.Fstatat(dirFd, runtimeSentinel+"/README", &st, unix.AT_SYMLINK_NOFOLLOW) == nil && st.Uid == 0
}

func plantSentinel(dirFd int) error {
	if err := unix.Mkdirat(dirFd, runtimeSentinel, 0o700); err != nil {
		return fmt.Errorf("creating %s: %w", runtimeSentinel, err)
	}
	sfd, err := unix.Openat(dirFd, runtimeSentinel, unix.O_RDONLY|unix.O_DIRECTORY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if err != nil {
		return fmt.Errorf("opening %s: %w", runtimeSentinel, err)
	}
	defer unix.Close(sfd)
	f, err := unix.Openat(sfd, "README", unix.O_WRONLY|unix.O_CREAT|unix.O_EXCL|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0o444)
	if err != nil {
		return fmt.Errorf("creating %s/README: %w", runtimeSentinel, err)
	}
	_, werr := unix.Write(f, []byte(sentinelNote))
	if cerr := unix.Close(f); werr == nil {
		werr = cerr
	}
	if werr != nil {
		return fmt.Errorf("writing %s/README: %w", runtimeSentinel, werr)
	}
	return unix.Fchmod(sfd, 0o555)
}
