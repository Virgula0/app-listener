package daemon

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/secretgate"
	"github.com/Virgula0/app-listener/internal/systemd"
)

// secretRelaySelfEvents is the relay dir guard's self mask: the daemon binds (MKNOD), hands over
// (ATTR) and unlinks (DELETE) the per-user relay sockets there. Never any other resource's.
var secretRelaySelfEvents = append(append([]ebpf.EventType(nil), daemonSelfBaselineEvents...),
	ebpf.EventMknod, ebpf.EventAttr, ebpf.EventDelete)

// selfEventsFor is r's guard's self mask.
func selfEventsFor(r *daemonconfig.Resource) []ebpf.EventType {
	if r.SecretRelay {
		return secretRelaySelfEvents
	}
	return daemonSelfBaselineEvents
}

// ensureSecretAnchors makes each [secret] scope's anchor and the relay dir before their guards claim
// them. /run/app-listener is 0755 so a provider can reach the relay dir; the anchor dirs stay 0700.
func ensureSecretAnchors(cfg *daemonconfig.Config) error {
	if len(cfg.SecretService.Providers) == 0 {
		return nil
	}
	if err := ensureRootDir(filepath.Dir(daemonconfig.SecretRelayDir), 0o755); err != nil {
		return err
	}
	if err := ensureRootDir(daemonconfig.SecretRelayDir, 0o755); err != nil {
		return err
	}
	if err := ensureRootDir(daemonconfig.SecretAnchorDir, 0o700); err != nil {
		return err
	}
	for i := range cfg.Resources {
		if s := cfg.Resources[i].Secret; s != nil {
			if err := ensureAnchor(daemonconfig.SecretAnchor(s.Name)); err != nil {
				return err
			}
		}
	}
	return nil
}

// secretJudge decides for the gates on the live guards: a scope's whitelist is its anchor's guard.
type secretJudge struct {
	svc    atomic.Pointer[daemonconfig.SecretService]
	scopes atomic.Pointer[[]secretgate.Scope]
}

func (j *secretJudge) setConfig(cfg *daemonconfig.Config) {
	svc := cfg.SecretService
	var scopes []secretgate.Scope
	for i := range cfg.Resources {
		if s := cfg.Resources[i].Secret; s != nil {
			scopes = append(scopes, secretgate.Scope{Name: s.Name, Matches: s.Matches})
		}
	}
	j.svc.Store(&svc)
	j.scopes.Store(&scopes)
}

func (j *secretJudge) Attest(pidfd int) (secretgate.Caller, error) {
	return guard.AttestBusClient(pidfd)
}

func (j *secretJudge) Scopes() []secretgate.Scope { return *j.scopes.Load() }

// RelayGuarded: a relay socket outside a live relay guard would let any of the user's processes
// connect as the keyring (a root bus connection that may own org.freedesktop.secrets).
func (j *secretJudge) RelayGuarded() bool {
	return guard.LiveGuardAt(daemonconfig.SecretRelayDir) != nil
}

func (j *secretJudge) Admits(scope string, c secretgate.Caller, op ebpf.EventType, what string) bool {
	g := guard.LiveGuardAt(daemonconfig.SecretAnchor(scope))
	if g == nil {
		log.Warnf("secret gate: refused %s: scope %q has no running guard", what, scope)
		return false
	}
	return g.AdmitsBusClient(c.(*guard.BusClient), op, what)
}

func (j *secretJudge) InRole(r secretgate.Role, c secretgate.Caller, what string) bool {
	bc := c.(*guard.BusClient)
	svc := j.svc.Load()
	paths, role := svc.Managers, "manager"
	if r == secretgate.RolePrompter {
		paths, role = svc.Prompters, "prompter"
	}
	for _, p := range paths {
		if key, err := placedKey(p); err != nil {
			log.Debugf("secret gate: %s %s: %v", role, p, err)
		} else if bc.RunsExe(key) {
			return true
		}
	}
	log.Warnf("secret gate: refused %s from pid %d (%s): not a [secret_service] %s", what, bc.PID, bc.Comm, role)
	return false
}

// placedKey is the inode a root-placed path names now (a role line names a system binary).
func placedKey(path string) (guard.GuardInodeKey, error) {
	f, err := ebpf.OpenSystemPlaced(path)
	if err != nil {
		return guard.GuardInodeKey{}, err
	}
	defer f.Close()
	var st unix.Stat_t
	if err := unix.Fstat(int(f.Fd()), &st); err != nil {
		return guard.GuardInodeKey{}, err
	}
	return guard.GuardInodeKey{Dev: ebpf.KernelDev(unix.Major(st.Dev), unix.Minor(st.Dev)), Ino: st.Ino}, nil
}

// secretGates runs one relay per logged-in user while [secret_service] names a provider. Each user's
// /run/user/<uid> is registered with the guard first, so its bus connects are attested from then on;
// processes connected earlier (before the daemon started) are refused scoped items until restarted.
type secretGates struct {
	judge secretJudge

	mu    sync.Mutex
	on    bool
	users map[uint32]*userGate
}

type userGate struct {
	gate *secretgate.Gate
	dir  guard.GuardInodeKey
}

func newSecretGates(cfg *daemonconfig.Config) *secretGates {
	s := &secretGates{users: map[uint32]*userGate{}}
	s.setConfig(cfg)
	return s
}

// setConfig follows the committed config: the relay guard is attached (or gone) by then.
func (s *secretGates) setConfig(cfg *daemonconfig.Config) {
	if s == nil || cfg == nil {
		return
	}
	s.judge.setConfig(cfg)
	on := len(cfg.SecretService.Providers) > 0
	if changed, err := systemd.SyncSecretServiceFiles(on); err != nil {
		log.Errorf("daemon: CRITICAL: %v — the session bus policy may not match [secret_service]", err)
	} else if changed && on {
		log.Warnf("daemon: installed %s and %s: users logged in now must log in again before the "+
			"Secret Service gate covers their keyring", systemd.SecretBusPolicyFile, systemd.SecretKeyringDropIn)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if on == s.on {
		return
	}
	s.on = on
	if !on {
		for uid := range s.users {
			s.stopLocked(uid)
		}
		return
	}
	for _, uid := range loggedInUsers() {
		s.startLocked(uid)
	}
}

func (s *secretGates) login(uid uint32) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.on {
		s.stopLocked(uid) // a previous session's tmpfs, if logout was missed
		s.startLocked(uid)
	}
}

func (s *secretGates) logout(uid uint32) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stopLocked(uid)
}

func (s *secretGates) close() {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.on = false
	for uid := range s.users {
		s.stopLocked(uid)
	}
}

// startLocked fails closed either way: without a gate the keyring can't reach the bus, without the
// bus dir registered no caller attests.
func (s *secretGates) startLocked(uid uint32) {
	runDir := fmt.Sprintf("/run/user/%d", uid)
	fd, gid, err := openUserRuntimeDir(uid)
	if err != nil {
		log.Errorf("secret gate %d: %v", uid, err)
		return
	}
	unix.Close(fd)
	dir, err := guard.RegisterBusDir(runDir)
	if err != nil {
		log.Errorf("secret gate %d: registering %s: %v", uid, runDir, err)
		return
	}
	sock := daemonconfig.SecretRelaySocket(uid)
	if rmErr := os.Remove(sock); rmErr != nil && !errors.Is(rmErr, os.ErrNotExist) {
		log.Errorf("secret gate %d: removing a stale %s: %v", uid, sock, rmErr)
	}
	g, err := secretgate.Listen(uid, gid, sock, userBus{uid: uid}, &s.judge)
	if err != nil {
		guard.UnregisterBusDir(dir)
		log.Errorf("secret gate %d: %v — the keyring can't reach the session bus", uid, err)
		return
	}
	s.users[uid] = &userGate{gate: g, dir: dir}
	go g.Serve()
}

// userBus is uid's session bus socket, opened without following a symlink and tagged so only a bus
// broker answers on it (guard.TagRelayBus): the user can put any listener at that name.
type userBus struct{ uid uint32 }

func (b userBus) String() string { return fmt.Sprintf("/run/user/%d/bus", b.uid) }

func (b userBus) Open() (path string, release func(), err error) {
	dir, _, err := openUserRuntimeDir(b.uid)
	if err != nil {
		return "", nil, err
	}
	defer unix.Close(dir)
	fd, err := unix.Openat(dir, "bus", unix.O_PATH|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0)
	if err != nil {
		return "", nil, err
	}
	untag, err := guard.TagRelayBus(fd)
	if err != nil {
		unix.Close(fd)
		return "", nil, err
	}
	return fmt.Sprintf("/proc/self/fd/%d", fd), func() {
		untag()
		unix.Close(fd)
	}, nil
}

func (s *secretGates) stopLocked(uid uint32) {
	u, ok := s.users[uid]
	if !ok {
		return
	}
	delete(s.users, uid)
	u.gate.Close()
	guard.UnregisterBusDir(u.dir)
}

// loggedInUsers lists the uids with a /run/user/<uid> mounted now.
func loggedInUsers() []uint32 {
	entries, err := os.ReadDir("/run/user")
	if err != nil {
		return nil
	}
	var out []uint32
	for _, e := range entries {
		if uid, err := strconv.ParseUint(e.Name(), 10, 32); err == nil && e.IsDir() {
			out = append(out, uint32(uid))
		}
	}
	return out
}
