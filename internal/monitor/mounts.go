package monitor

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// mountTable maps what the BPF dentry walk yields (a superblock's kernel s_dev plus a path from
// that superblock's root) to absolute paths in this process's mount namespace. Without it a file
// on a separate mount (a tmpfs /tmp) reaches userspace as "/name" and matches no watch target.
type mountTable struct {
	mu       sync.Mutex
	file     string
	entries  []mountEntry // mountinfo order: a later mount shadows an earlier one at the same point
	byDev    map[uint32][]int
	loadedAt time.Time
}

type mountEntry struct {
	dev        uint32 // ebpf.KernelDev encoding, as BPF reads sb->s_dev
	root       string // path inside the filesystem that is mounted here
	mountpoint string
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
		rel, ok := underRoot(fsPath, e.root)
		if !ok {
			continue
		}
		cand := joinMount(e.mountpoint, rel)
		if prefer != nil && prefer(cand) {
			return cand, true
		}
		if len(e.root) >= bestRoot {
			best, bestRoot = cand, len(e.root)
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
		if _, ok := underRoot(p, e.mountpoint); ok && len(e.mountpoint) >= bestLen {
			dev, bestLen = e.dev, len(e.mountpoint)
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
	t.entries, t.byDev = parseMountinfo(bufio.NewScanner(f))
}

func parseMountinfo(sc *bufio.Scanner) (entries []mountEntry, byDev map[uint32][]int) {
	entries = make([]mountEntry, 0, 64)
	byDev = make(map[uint32][]int)
	for sc.Scan() {
		// id parent major:minor root mountpoint options ...
		fields := strings.Fields(sc.Text())
		if len(fields) < 5 {
			continue
		}
		dev, ok := parseKernelDev(fields[2])
		if !ok {
			continue
		}
		byDev[dev] = append(byDev[dev], len(entries))
		entries = append(entries, mountEntry{
			dev:        dev,
			root:       unescapeMountinfo(fields[3]),
			mountpoint: unescapeMountinfo(fields[4]),
		})
	}
	return entries, byDev
}

// parseKernelDev parses mountinfo's "major:minor" into the kernel dev_t (12-bit major, 20-bit minor).
func parseKernelDev(s string) (uint32, bool) {
	major, minor, ok := strings.Cut(s, ":")
	if !ok {
		return 0, false
	}
	maj, err := strconv.ParseUint(major, 10, 12)
	if err != nil {
		return 0, false
	}
	mnr, err := strconv.ParseUint(minor, 10, 20)
	if err != nil {
		return 0, false
	}
	return uint32(maj)<<20 | uint32(mnr), true
}

// unescapeMountinfo decodes the kernel's \ooo octal escapes (space, tab, newline, backslash).
func unescapeMountinfo(s string) string {
	if !strings.Contains(s, `\`) {
		return s
	}
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+4 <= len(s) {
			if v, err := strconv.ParseUint(s[i+1:i+4], 8, 8); err == nil {
				b.WriteByte(byte(v))
				i += 3
				continue
			}
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

// underRoot returns p relative to root (as "/..."), or false when p is outside root.
func underRoot(p, root string) (string, bool) {
	switch {
	case root == "/":
		return p, strings.HasPrefix(p, "/")
	case p == root:
		return "/", true
	case strings.HasPrefix(p, root+"/"):
		return p[len(root):], true
	}
	return "", false
}

func joinMount(mountpoint, rel string) string {
	switch {
	case mountpoint == "/":
		return rel
	case rel == "/":
		return mountpoint
	}
	return mountpoint + rel
}
