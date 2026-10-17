package workbench

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// Exe is an executable as seen at event time: the path it ran from and the inode it ran. Path is
// display only; promoting it to a policy re-checks that the path still names Key (EntryFor).
type Exe struct {
	Path string
	Key  guard.GuardInodeKey // real key (KernelDev encoding); zero when unknown
}

// Known reports whether the inode was resolved.
func (e Exe) Known() bool { return e.Key != (guard.GuardInodeKey{}) }

// exeResolver maps a pid to its Exe through /proc/<pid>/exe, cached briefly: a burst of events from
// one process costs one readlink.
type exeResolver struct {
	procRoot string
	mu       sync.Mutex
	cache    map[uint32]exeCached
}

type exeCached struct {
	exe Exe
	at  time.Time
}

const exeCacheTTL = 2 * time.Second

func newExeResolver(procRoot string) *exeResolver {
	return &exeResolver{procRoot: procRoot, cache: map[uint32]exeCached{}}
}

func (r *exeResolver) resolve(pid uint32) Exe {
	if pid == 0 {
		return Exe{}
	}
	now := time.Now()
	r.mu.Lock()
	if c, ok := r.cache[pid]; ok && now.Sub(c.at) < exeCacheTTL {
		r.mu.Unlock()
		return c.exe
	}
	r.mu.Unlock()

	link := fmt.Sprintf("%s/%d/exe", r.procRoot, pid)
	var exe Exe
	if path, err := os.Readlink(link); err == nil {
		exe.Path = strings.TrimSuffix(path, " (deleted)")
		var st unix.Stat_t
		if unix.Stat(link, &st) == nil {
			exe.Key = guard.GuardInodeKey{Dev: ebpf.KernelDev(unix.Major(st.Dev), unix.Minor(st.Dev)), Ino: st.Ino}
		}
	}

	r.mu.Lock()
	if len(r.cache) > 4096 {
		r.cache = map[uint32]exeCached{}
	}
	r.cache[pid] = exeCached{exe: exe, at: now}
	r.mu.Unlock()
	return exe
}

// ErrExeUnknown: the event's executable could not be resolved (process gone before it was read).
var ErrExeUnknown = errors.New("the executable of this process is unknown: add it by path instead")

// EntryFor hashes the executable e names, refusing when its path no longer leads to the inode that
// ran: whatever is there now is not what was observed.
func EntryFor(e Exe) (ebpf.BinaryEntry, error) {
	if e.Path == "" || !e.Known() {
		return ebpf.BinaryEntry{}, ErrExeUnknown
	}
	f, err := ebpf.OpenConfined(e.Path)
	if err != nil {
		return ebpf.BinaryEntry{}, fmt.Errorf("opening %s: %w", e.Path, err)
	}
	defer f.Close()
	var st unix.Stat_t
	if err := unix.Fstat(int(f.Fd()), &st); err != nil {
		return ebpf.BinaryEntry{}, fmt.Errorf("stat %s: %w", e.Path, err)
	}
	now := guard.GuardInodeKey{Dev: ebpf.KernelDev(unix.Major(st.Dev), unix.Minor(st.Dev)), Ino: st.Ino}
	if now != e.Key.Real() {
		return ebpf.BinaryEntry{}, fmt.Errorf("%s was replaced since the process ran (inode %d, now %d): "+
			"add it by path if the new file is the one to trust", e.Path, e.Key.Ino, now.Ino)
	}
	if st.Mode&unix.S_IFMT != unix.S_IFREG {
		return ebpf.BinaryEntry{}, fmt.Errorf("%s is not a regular file", e.Path)
	}
	return ebpf.ComputeBinaryEntryFile(f, e.Path)
}

// Process aggregates one executable's activity on a target.
type Process struct {
	Exe      Exe
	Comm     string
	Events   uint64
	Denied   uint64
	LastSeen time.Time
}

func processKey(exe Exe, comm string) string {
	if exe.Known() {
		return fmt.Sprintf("%d:%d", exe.Key.Real().Dev, exe.Key.Real().Ino)
	}
	return "comm:" + comm
}
