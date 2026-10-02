package daemon

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/install"
)

// The catalog watch tells the daemon when an install may have changed what the catalog admits. It
// is only ever a hint: an event schedules the same vetting a refresh or re-sync does on its own,
// and nothing about the event (name, path, who wrote it) is trusted.

// watchPattern is one catalog pattern (or whitelisted system binary) split into path components.
type watchPattern struct {
	comps []string
	// home: a file in a user's home, vetted by the catalog refresh. Otherwise root's: the re-sync
	// admits it (system files only).
	home bool
	// dir: a lib dir, whose appearance is the trigger. Otherwise a binary: only a write that
	// completes (IN_CLOSE_WRITE) or a rename onto the name (IN_MOVED_TO) triggers, never a mkdir.
	dir bool
}

func newWatchPattern(path string, home, dir bool) *watchPattern {
	return &watchPattern{comps: strings.Split(strings.TrimPrefix(filepath.Clean(path), "/"), "/"), home: home, dir: dir}
}

func (p *watchPattern) key() string {
	return strings.Join(p.comps, "/") + map[bool]string{true: "|h", false: "|r"}[p.home] +
		map[bool]string{true: "d", false: "b"}[p.dir]
}

// root is the deepest existing directory along p's fixed leading components, and the index of the
// component its entries must match.
func (p *watchPattern) root() (dir string, next int) {
	dir = "/"
	for i, c := range p.comps[:len(p.comps)-1] {
		if strings.ContainsAny(c, "*?[") {
			return dir, i
		}
		next := filepath.Join(dir, c)
		if fi, err := os.Stat(next); err != nil || !fi.IsDir() {
			return dir, i
		}
		dir = next
	}
	return dir, len(p.comps) - 1
}

// catalogWatchPlan lists what the daemon watches for cfg: every catalog pattern of a configured
// resource (its user's home patterns and the absolute ones), and every whitelisted root-owned binary
// with the file its path resolves to (package managers install by temp file + rename).
func catalogWatchPlan(cfg *daemonconfig.Config, users []install.User) []*watchPattern {
	seen := map[string]bool{}
	var out []*watchPattern
	add := func(p *watchPattern) {
		if k := p.key(); !seen[k] {
			seen[k] = true
			out = append(out, p)
		}
	}
	for i := range cfg.Resources {
		r := &cfg.Resources[i]
		catalogPatterns(r, users, add)
		systemBinaryPatterns(r, add)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].key() < out[j].key() })
	return out
}

func catalogPatterns(r *daemonconfig.Resource, users []install.User, add func(*watchPattern)) {
	entry, user := install.ResolveCatalogEntry(r.EncryptionRootOrPath(), users)
	if entry == nil {
		return
	}
	bins, dirs := entry.WatchPatterns(user.Name, user.Home)
	for _, b := range bins {
		add(newWatchPattern(b, inHome(b, user.Home), false))
	}
	for _, d := range dirs {
		add(newWatchPattern(d, true, true))
	}
}

// systemBinaryPatterns: each whitelisted root-owned binary, and the file its path resolves to.
func systemBinaryPatterns(r *daemonconfig.Resource, add func(*watchPattern)) {
	for _, list := range [][]daemonconfig.BinaryRule{r.Binaries, r.PendingBinaries} {
		for _, b := range list {
			if !ebpf.SystemTrusted(b.Path) {
				continue
			}
			add(newWatchPattern(b.Path, false, false))
			if resolved, err := filepath.EvalSymlinks(b.Path); err == nil && resolved != b.Path {
				add(newWatchPattern(resolved, false, false))
			}
		}
	}
}

func inHome(path, home string) bool {
	home = strings.TrimSuffix(home, "/")
	return home != "" && strings.HasPrefix(path, home+"/")
}

const (
	// tempWatchTTL bounds how long a new directory below a glob root is followed: an updater
	// writes its binary within minutes of creating the version dir.
	tempWatchTTL = 10 * time.Minute
	// maxTempWatches caps followed directories: a process filling a glob root with directories
	// gets a rescan, not a watch each.
	maxTempWatches = 256
	// maxDescend caps the entries read from one new directory.
	maxDescend = 64

	watchMask = unix.IN_CREATE | unix.IN_MOVED_TO | unix.IN_CLOSE_WRITE | unix.IN_ONLYDIR | unix.IN_EXCL_UNLINK
)

type watchPos struct {
	p    *watchPattern
	next int // index of the component this directory's entries must match
}

type watchDir struct {
	path    string
	pos     []watchPos
	expires time.Time // zero: permanent
}

// watchSet maps inotify watches to pattern positions and turns events into triggers.
type watchSet struct {
	add   func(path string) (int32, error)
	rm    func(wd int32)
	now   func() time.Time
	dirs  map[int32]*watchDir
	temps int
	// set by handle: what the event asks for
	refresh, resync, rebuild bool
}

func newWatchSet(add func(string) (int32, error), rm func(int32), now func() time.Time) *watchSet {
	return &watchSet{add: add, rm: rm, now: now, dirs: map[int32]*watchDir{}}
}

// build watches every pattern's root.
func (s *watchSet) build(plan []*watchPattern) {
	for _, p := range plan {
		dir, next := p.root()
		s.watch(dir, watchPos{p, next}, false)
	}
}

func (s *watchSet) trigger(p *watchPattern) {
	if p.home {
		s.refresh = true
	} else {
		s.resync = true
	}
}

// watch adds pos at dir. A temporary watch past maxTempWatches is not added: the caller rescans.
func (s *watchSet) watch(dir string, pos watchPos, temp bool) bool {
	if temp && s.temps >= maxTempWatches {
		return false
	}
	wd, err := s.add(dir)
	if err != nil {
		return false
	}
	d := s.dirs[wd]
	if d == nil {
		d = &watchDir{path: dir}
		s.dirs[wd] = d
		if temp {
			d.expires = s.now().Add(tempWatchTTL)
			s.temps++
		}
	} else if !temp && !d.expires.IsZero() {
		d.expires = time.Time{}
		s.temps--
	}
	for _, have := range d.pos {
		if have == pos {
			return true
		}
	}
	d.pos = append(d.pos, pos)
	return true
}

// follow watches a directory that just appeared below a pattern position, and catches a match
// already inside it (a version dir moved in whole, or entries created before the watch).
func (s *watchSet) follow(dir string, pos watchPos) {
	if !s.watch(dir, pos, true) {
		s.trigger(pos.p)
		return
	}
	f, err := os.Open(dir)
	if err != nil {
		return
	}
	entries, _ := f.ReadDir(maxDescend)
	f.Close()
	for _, e := range entries {
		var mask uint32 = unix.IN_MOVED_TO
		if e.IsDir() {
			mask |= unix.IN_ISDIR
		}
		s.match(dir, pos, mask, e.Name())
	}
}

// match applies one entry event in dir to one pattern position.
func (s *watchSet) match(dir string, pos watchPos, mask uint32, name string) {
	comps := pos.p.comps
	if ok, err := filepath.Match(comps[pos.next], name); err != nil || !ok {
		return
	}
	isDir := mask&unix.IN_ISDIR != 0
	arrived := mask&(unix.IN_CREATE|unix.IN_MOVED_TO) != 0
	switch last := pos.next == len(comps)-1; {
	case last && pos.p.dir:
		if isDir && arrived {
			s.trigger(pos.p)
		}
	case last:
		if !isDir && mask&(unix.IN_CLOSE_WRITE|unix.IN_MOVED_TO) != 0 {
			s.trigger(pos.p)
		}
	case isDir && arrived:
		s.follow(filepath.Join(dir, name), watchPos{pos.p, pos.next + 1})
	}
}

// handle applies one inotify event.
func (s *watchSet) handle(wd int32, mask uint32, name string) {
	if mask&unix.IN_Q_OVERFLOW != 0 {
		s.refresh, s.resync, s.rebuild = true, true, true
		return
	}
	d := s.dirs[wd]
	if d == nil {
		return
	}
	if mask&unix.IN_IGNORED != 0 {
		delete(s.dirs, wd)
		if d.expires.IsZero() {
			s.rebuild = true // a root went away: find its deepest existing ancestor again
		} else {
			s.temps--
		}
		return
	}
	for _, pos := range append([]watchPos(nil), d.pos...) {
		s.match(d.path, pos, mask, name)
	}
}

// expire drops temporary watches past their TTL.
func (s *watchSet) expire() {
	now := s.now()
	for wd, d := range s.dirs {
		if !d.expires.IsZero() && now.After(d.expires) {
			s.rm(wd)
			delete(s.dirs, wd)
			s.temps--
		}
	}
}

// take returns and clears the pending triggers.
func (s *watchSet) take() (refresh, resync, rebuild bool) {
	refresh, resync, rebuild = s.refresh, s.resync, s.rebuild
	s.refresh, s.resync, s.rebuild = false, false, false
	return refresh, resync, rebuild
}
