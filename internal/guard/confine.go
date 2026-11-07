package guard

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"golang.org/x/sys/unix"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// ConnectConfinedNames are whitelisted binaries that may unix-connect only to a socket bound inside
// a guarded tree: ssh-add sends the key it read to whatever agent SSH_AUTH_SOCK names. A name only
// adds denials, so matching the line's name (or its target's) weakens nothing.
var ConnectConfinedNames = []string{"ssh-add"}

// maxConnectConfined mirrors guard_connect_confined's max_entries in guard.bpf.c.
const maxConnectConfined = 64

// connectConfined holds every guard_connect_confined binary open (under sharedEngine.mu): its inode
// can't be freed and its number reused by another program while the row confines it.
var connectConfined = map[GuardInodeKey]*os.File{}

// confineConnect lists key in guard_connect_confined when path, or the file it resolves to, is a
// ConnectConfinedNames binary. Called before the allow row is written, so no moment admits the
// binary unconfined.
func confineConnect(key GuardInodeKey, path string) error {
	if !isConnectConfined(path) && !isConnectConfined(canonicalBinaryPath(path)) {
		return nil
	}
	e := &sharedEngine
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, held := connectConfined[key.Real()]; held {
		return nil
	}
	if e.objs.GuardConnectConfined == nil {
		return fmt.Errorf("confining %s: guard engine not running", path)
	}
	if len(connectConfined) >= maxConnectConfined {
		return fmt.Errorf("confining %s: more than %d confined binaries", path, maxConnectConfined)
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
	// The allow row names key: confining another inode would leave that one free.
	got := GuardInodeKey{Dev: ebpf.KernelDev(unix.Major(st.Dev), unix.Minor(st.Dev)), Ino: st.Ino}
	if got != key.Real() {
		f.Close()
		return fmt.Errorf("confining %s: replaced while it was admitted (inode %d, want %d)", path, st.Ino, key.Ino)
	}
	if err := e.objs.GuardConnectConfined.Put(got, uint8(1)); err != nil {
		f.Close()
		return fmt.Errorf("confining %s: %w", path, err)
	}
	connectConfined[got] = f
	return nil
}

func isConnectConfined(path string) bool {
	return slices.Contains(ConnectConfinedNames, filepath.Base(path))
}

// syncConnectConfinedLocked lists the binaries already held in a fresh engine's map: confineConnect
// skips a held one.
func (e *engine) syncConnectConfinedLocked() error {
	for k := range connectConfined {
		if err := e.objs.GuardConnectConfined.Put(k, uint8(1)); err != nil {
			return fmt.Errorf("confining exe inode %d: %w", k.Ino, err)
		}
	}
	return nil
}
