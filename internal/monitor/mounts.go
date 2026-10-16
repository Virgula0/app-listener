package monitor

import (
	"os"
	"sync"
	"time"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// mountTable maps what the BPF dentry walk yields (a superblock's kernel s_dev plus a path from
// that superblock's root) to absolute paths in this process's mount namespace. Without it a file
// on a separate mount (a tmpfs /tmp) reaches userspace as "/name" and matches no watch target.
type mountTable struct {
	mu       sync.Mutex
	file     string
	entries  []ebpf.MountInfo
	byDev    map[uint32][]int
	loadedAt time.Time
}

const (
	mountTTL       = 10 * time.Second
	mountReloadMin = time.Second
)

func newMountTable() *mountTable {
	return &mountTable{file: "/proc/self/mountinfo"}
}

// resolve returns the absolute path of fsPath on superblock dev, preferring a candidate prefer
// accepts (a bind-mounted filesystem has several). False when no mount of dev covers fsPath.
func (t *mountTable) resolve(dev uint32, fsPath string, prefer func(string) bool) (string, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	idx := t.lookupLocked(dev)

	best, bestRoot := "", -1
	for _, i := range idx {
		e := t.entries[i]
		cand, ok := e.MountPath(fsPath)
		if !ok {
			continue
		}
		if prefer != nil && prefer(cand) {
			return cand, true
		}
		if len(e.Root) >= bestRoot {
			best, bestRoot = cand, len(e.Root)
		}
	}
	return best, bestRoot >= 0
}

// devOf returns the s_dev of the mount holding the absolute path p, or 0.
func (t *mountTable) devOf(p string) uint32 {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.entries == nil || time.Since(t.loadedAt) > mountTTL {
		t.loadLocked()
	}
	var dev uint32
	bestLen := -1
	for _, e := range t.entries {
		if _, ok := ebpf.PathUnder(p, e.Mountpoint); ok && len(e.Mountpoint) >= bestLen {
			dev, bestLen = e.Dev, len(e.Mountpoint)
		}
	}
	return dev
}

func (t *mountTable) lookupLocked(dev uint32) []int {
	if t.entries == nil || time.Since(t.loadedAt) > mountTTL {
		t.loadLocked()
	}
	idx, ok := t.byDev[dev]
	if !ok && time.Since(t.loadedAt) >= mountReloadMin {
		t.loadLocked()
		idx = t.byDev[dev]
	}
	return idx
}

func (t *mountTable) loadLocked() {
	t.loadedAt = time.Now()
	f, err := os.Open(t.file)
	if err != nil {
		return
	}
	defer f.Close()
	entries, err := ebpf.ParseMountinfo(f)
	if err != nil && len(entries) == 0 {
		return
	}
	t.entries = entries
	t.byDev = make(map[uint32][]int, len(entries))
	for i, e := range entries {
		t.byDev[e.Dev] = append(t.byDev[e.Dev], i)
	}
}
