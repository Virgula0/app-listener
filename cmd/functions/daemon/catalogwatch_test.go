package daemon

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/install"
)

// fakeInotify hands out one watch descriptor per path, like inotify does per inode.
type fakeInotify struct {
	wds     map[string]int32
	removed []int32
	now     time.Time
}

func newFakeInotify() *fakeInotify {
	return &fakeInotify{wds: map[string]int32{}, now: time.Unix(1000, 0)}
}

func (f *fakeInotify) set() *watchSet {
	return newWatchSet(func(p string) (int32, error) {
		if _, err := os.Stat(p); err != nil {
			return 0, err
		}
		if wd, ok := f.wds[p]; ok {
			return wd, nil
		}
		wd := int32(len(f.wds) + 1) //nolint:gosec // test
		f.wds[p] = wd
		return wd, nil
	}, func(wd int32) { f.removed = append(f.removed, wd) }, func() time.Time { return f.now })
}

func TestWatchSetBinaryWriteTriggersNotMkdir(t *testing.T) {
	root := filepath.Join(t.TempDir(), ".config", "discord")
	mkdirs(t, root)
	fi := newFakeInotify()
	s := fi.set()
	s.build([]*watchPattern{newWatchPattern(filepath.Join(root, "*", "Discord"), true, false)})

	mkdirs(t, filepath.Join(root, "0.0.2"))
	s.handle(fi.wds[root], unix.IN_CREATE|unix.IN_ISDIR, "0.0.2")
	if r, _, _ := s.take(); r {
		t.Fatal("a new version dir alone must not trigger a refresh")
	}
	vdir := fi.wds[filepath.Join(root, "0.0.2")]
	if vdir == 0 {
		t.Fatal("the new version dir must be followed")
	}
	s.handle(vdir, unix.IN_CREATE, "Discord")
	if r, _, _ := s.take(); r {
		t.Fatal("creating the binary must not trigger before its write completes")
	}
	s.handle(vdir, unix.IN_CLOSE_WRITE, "Other")
	if r, _, _ := s.take(); r {
		t.Fatal("a name the pattern does not match must not trigger")
	}
	s.handle(vdir, unix.IN_CLOSE_WRITE, "Discord")
	// The re-sync too: an in-place update (same path) changes no config line for the refresh.
	if r, rs, _ := s.take(); !r || !rs {
		t.Fatalf("a completed write of the binary in a home pattern: refresh=%v resync=%v, want both", r, rs)
	}
}

// An updater re-pointing its link by unlink + symlink() emits IN_CREATE only; a regular file's
// IN_CREATE still waits for its write.
func TestWatchSetSymlinkCreateTriggers(t *testing.T) {
	bin := filepath.Join(t.TempDir(), ".local", "bin")
	mkdirs(t, bin)
	fi := newFakeInotify()
	s := fi.set()
	s.build([]*watchPattern{newWatchPattern(filepath.Join(bin, "claude"), true, false)})

	if err := os.WriteFile(filepath.Join(bin, "claude"), nil, 0o755); err != nil {
		t.Fatal(err)
	}
	s.handle(fi.wds[bin], unix.IN_CREATE, "claude")
	if r, rs, _ := s.take(); r || rs {
		t.Fatal("a regular file's IN_CREATE must not trigger before its write completes")
	}
	if err := os.Remove(filepath.Join(bin, "claude")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("../share/claude/versions/2", filepath.Join(bin, "claude")); err != nil {
		t.Fatal(err)
	}
	s.handle(fi.wds[bin], unix.IN_CREATE, "claude")
	if r, rs, _ := s.take(); !r || !rs {
		t.Fatalf("a symlink created at the binary's name: refresh=%v resync=%v, want both", r, rs)
	}
}

func TestWatchSetVersionDirMovedInWhole(t *testing.T) {
	root := filepath.Join(t.TempDir(), "discord")
	mkdirs(t, root, filepath.Join(root, "0.0.3"))
	if err := os.WriteFile(filepath.Join(root, "0.0.3", "Discord"), nil, 0o755); err != nil {
		t.Fatal(err)
	}
	fi := newFakeInotify()
	s := fi.set()
	s.build([]*watchPattern{newWatchPattern(filepath.Join(root, "*", "Discord"), true, false)})
	s.handle(fi.wds[root], unix.IN_MOVED_TO|unix.IN_ISDIR, "0.0.3")
	if r, _, _ := s.take(); !r {
		t.Fatal("a version dir moved in with its binary must trigger a refresh")
	}
}

func TestWatchSetSystemBinaryRenameTriggersResync(t *testing.T) {
	bin := filepath.Join(t.TempDir(), "usr", "bin")
	mkdirs(t, bin)
	fi := newFakeInotify()
	s := fi.set()
	s.build([]*watchPattern{newWatchPattern(filepath.Join(bin, "ssh"), false, false)})
	s.handle(fi.wds[bin], unix.IN_CLOSE_WRITE, "ssh.pacnew")
	if r, rs, _ := s.take(); r || rs {
		t.Fatal("the package manager's temp file must not trigger")
	}
	s.handle(fi.wds[bin], unix.IN_MOVED_TO, "ssh")
	if r, rs, _ := s.take(); r || !rs {
		t.Fatalf("a rename onto a whitelisted system binary: refresh=%v resync=%v, want resync only", r, rs)
	}
}

func TestWatchSetMissingRootWatchesDeepestAncestor(t *testing.T) {
	home := t.TempDir()
	fi := newFakeInotify()
	s := fi.set()
	s.build([]*watchPattern{newWatchPattern(filepath.Join(home, ".local", "bin", "claude"), true, false)})
	if fi.wds[home] == 0 {
		t.Fatalf("the deepest existing ancestor must be watched, got %v", fi.wds)
	}
	mkdirs(t, filepath.Join(home, ".local", "bin"))
	s.handle(fi.wds[home], unix.IN_CREATE|unix.IN_ISDIR, ".local")
	if fi.wds[filepath.Join(home, ".local", "bin")] == 0 {
		t.Fatal("a fixed component created after the watch must be followed down (mkdir -p)")
	}
	s.handle(fi.wds[filepath.Join(home, ".local", "bin")], unix.IN_MOVED_TO, "claude")
	if r, _, _ := s.take(); !r {
		t.Fatal("the binary landing in the created chain must trigger")
	}
}

func TestWatchSetFloodIsCapped(t *testing.T) {
	root := t.TempDir()
	fi := newFakeInotify()
	s := fi.set()
	s.build([]*watchPattern{newWatchPattern(filepath.Join(root, "*", "Discord"), true, false)})
	for i := range maxTempWatches + 50 {
		name := fmt.Sprintf("junk%d", i)
		mkdirs(t, filepath.Join(root, name))
		s.handle(fi.wds[root], unix.IN_CREATE|unix.IN_ISDIR, name)
	}
	if s.temps > maxTempWatches {
		t.Fatalf("%d temporary watches, cap %d", s.temps, maxTempWatches)
	}
	if r, _, _ := s.take(); !r {
		t.Fatal("past the cap a new directory must fall back to a rescan")
	}
}

func TestWatchSetOverflowRescansAndRebuilds(t *testing.T) {
	s := newFakeInotify().set()
	s.handle(-1, unix.IN_Q_OVERFLOW, "")
	if r, rs, b := s.take(); !r || !rs || !b {
		t.Fatalf("overflow: refresh=%v resync=%v rebuild=%v, want all", r, rs, b)
	}
}

func TestWatchSetTemporaryWatchesExpire(t *testing.T) {
	root := t.TempDir()
	mkdirs(t, filepath.Join(root, "v1"))
	fi := newFakeInotify()
	s := fi.set()
	s.build([]*watchPattern{newWatchPattern(filepath.Join(root, "*", "Discord"), true, false)})
	s.handle(fi.wds[root], unix.IN_CREATE|unix.IN_ISDIR, "v1")
	fi.now = fi.now.Add(tempWatchTTL + time.Second)
	s.expire()
	if s.temps != 0 || len(fi.removed) != 1 || len(s.dirs) != 1 {
		t.Fatalf("after the TTL: temps=%d removed=%v dirs=%d, want only the root", s.temps, fi.removed, len(s.dirs))
	}
}

func TestWatchSetLibDirAppearanceTriggers(t *testing.T) {
	root := t.TempDir()
	mkdirs(t, filepath.Join(root, "Proton 9", "files"))
	fi := newFakeInotify()
	s := fi.set()
	s.build([]*watchPattern{newWatchPattern(filepath.Join(root, "Proton*", "files", "lib"), true, true)})
	s.handle(fi.wds[root], unix.IN_CREATE|unix.IN_ISDIR, "Proton 9")
	mkdirs(t, filepath.Join(root, "Proton 9", "files", "lib"))
	s.handle(fi.wds[filepath.Join(root, "Proton 9", "files")], unix.IN_CREATE|unix.IN_ISDIR, "lib")
	if r, _, _ := s.take(); !r {
		t.Fatal("a lib dir appearing must trigger the refresh that adopts it")
	}
}

func TestDebouncerCoalescesAndBoundsWait(t *testing.T) {
	var runs atomic.Int32
	d := &debouncer{quiet: 50 * time.Millisecond, maxWait: 200 * time.Millisecond, minGap: 0,
		fire: func() { runs.Add(1) }}
	defer d.stop()
	start := time.Now()
	for time.Since(start) < 400*time.Millisecond {
		d.poke()
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(100 * time.Millisecond)
	if n := runs.Load(); n < 2 || n > 4 {
		t.Fatalf("400ms of constant pokes ran %d time(s): want a run per maxWait, not one per poke or none", n)
	}
}

func TestDebouncerMinGap(t *testing.T) {
	var runs atomic.Int32
	d := &debouncer{quiet: time.Millisecond, maxWait: time.Millisecond, minGap: 300 * time.Millisecond,
		fire: func() { runs.Add(1) }}
	defer d.stop()
	d.now()
	time.Sleep(50 * time.Millisecond)
	d.now()
	time.Sleep(100 * time.Millisecond)
	if n := runs.Load(); n != 1 {
		t.Fatalf("a second run within minGap: %d runs, want 1", n)
	}
	time.Sleep(300 * time.Millisecond)
	if n := runs.Load(); n != 2 {
		t.Fatalf("the deferred run must happen after minGap: %d runs, want 2", n)
	}
}

func TestLogRefreshChange(t *testing.T) {
	var b bytes.Buffer
	logRefreshChange(&b, install.SectionChange{Section: "/h/.config/discord", Admitted: []string{"/a", "/b"}})
	want := "<6>DAEMON catalog-refresh  resource=/h/.config/discord  admitted=/a,/b  dropped=-\n"
	if b.String() != want {
		t.Fatalf("got %q, want %q", b.String(), want)
	}
}
