package guard

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"golang.org/x/sys/unix"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// guard_unix_roles values (guard.bpf.c UNIX_CONFINED / UNIX_AGENT).
const (
	unixConfined uint8 = 1
	unixAgent    uint8 = 2
)

// ConnectConfinedNames are whitelisted binaries that may unix-connect only to a socket bound inside
// a guarded tree, whose server end only an AgentNames binary may use: ssh-add sends the key it read
// to whatever agent SSH_AUTH_SOCK names. A name only adds denials, so matching the line's name (or
// its target's) weakens nothing.
var ConnectConfinedNames = []string{"ssh-add"}

// AgentNames are the whitelisted binaries that may serve a ConnectConfinedNames connection. A name
// match (never another file) only lets that inode read what the confined client sends.
var AgentNames = []string{"ssh-agent", "gpg-agent", "gcr-ssh-agent", "1password"}

// maxUnixRoles mirrors guard_unix_roles's max_entries in guard.bpf.c.
const maxUnixRoles = 64

// unixRoles holds every guard_unix_roles binary open (under sharedEngine.mu): its inode can't be
// freed and its number reused by another program while the row names it.
var unixRoles = map[GuardInodeKey]heldRole{}

type heldRole struct {
	f    *os.File
	role uint8
}

// confineConnect lists key in guard_unix_roles when path, or the file it resolves to, has a
// ConnectConfinedNames or AgentNames name. Called before the allow row is written, so no moment
// admits a confined binary unconfined.
func confineConnect(key GuardInodeKey, path string) error {
	role := unixRole(path) | unixRole(canonicalBinaryPath(path))
	if role == 0 {
		return nil
	}
	e := &sharedEngine
	e.mu.Lock()
	defer e.mu.Unlock()
	held, ok := unixRoles[key.Real()]
	if ok && held.role&role == role {
		return nil
	}
	if e.objs.GuardUnixRoles == nil {
		return fmt.Errorf("confining %s: guard engine not running", path)
	}
	if !ok && len(unixRoles) >= maxUnixRoles {
		return fmt.Errorf("confining %s: more than %d unix-role binaries", path, maxUnixRoles)
	}
	f, err := ebpf.OpenConfined(path)
	if err != nil {
		return fmt.Errorf("confining %s: %w", path, err)
	}
	var st unix.Stat_t
	if err := unix.Fstat(int(f.Fd()), &st); err != nil {
		f.Close()
		return fmt.Errorf("confining %s: %w", path, err)
	}
	// The allow row names key: listing another inode would leave that one unconfined.
	got := GuardInodeKey{Dev: ebpf.KernelDev(unix.Major(st.Dev), unix.Minor(st.Dev)), Ino: st.Ino}
	if got != key.Real() {
		f.Close()
		return fmt.Errorf("confining %s: replaced while it was admitted (inode %d, want %d)", path, st.Ino, key.Ino)
	}
	role |= held.role
	if err := e.objs.GuardUnixRoles.Put(got, role); err != nil {
		f.Close()
		return fmt.Errorf("confining %s: %w", path, err)
	}
	if ok {
		f.Close()
		f = held.f
	}
	unixRoles[got] = heldRole{f: f, role: role}
	return nil
}

func unixRole(path string) uint8 {
	var r uint8
	if slices.Contains(ConnectConfinedNames, filepath.Base(path)) {
		r |= unixConfined
	}
	if slices.Contains(AgentNames, filepath.Base(path)) {
		r |= unixAgent
	}
	return r
}

// syncUnixRolesLocked lists the binaries already held in a fresh engine's map: confineConnect skips
// a held one.
func (e *engine) syncUnixRolesLocked() error {
	for k, h := range unixRoles {
		if err := e.objs.GuardUnixRoles.Put(k, h.role); err != nil {
			return fmt.Errorf("confining exe inode %d: %w", k.Ino, err)
		}
	}
	return nil
}
