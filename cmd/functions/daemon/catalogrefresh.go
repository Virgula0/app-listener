package daemon

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
	"unsafe"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/logging"
)

const (
	// A refresh waits for refreshQuiet without events, at most refreshMaxWait after the first, and
	// runs at most every refreshMinGap: a flood of events costs one rescan per gap, never a stall.
	refreshQuiet   = 2 * time.Second
	refreshMaxWait = 10 * time.Second
	refreshMinGap  = 10 * time.Second
	// A system binary's re-sync is cheap (one stat per whitelisted path): shorter debounce.
	resyncQuiet   = 500 * time.Millisecond
	resyncMaxWait = 3 * time.Second
	resyncMinGap  = 2 * time.Second
)

// debouncer runs fire once events stop for quiet, at most maxWait after the first and never within
// minGap of the last run. Runs never overlap; a poke during a run schedules one more.
type debouncer struct {
	quiet, maxWait, minGap time.Duration
	fire                   func()

	mu      sync.Mutex
	first   time.Time
	lastRun time.Time
	timer   *time.Timer
	running bool
	again   bool
	stopped bool
}

func (d *debouncer) poke() { d.schedule(d.quiet) }

// now asks for a run as soon as minGap allows.
func (d *debouncer) now() { d.schedule(0) }

func (d *debouncer) schedule(quiet time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.stopped {
		return
	}
	if d.running {
		d.again = true
		return
	}
	now := time.Now()
	if d.first.IsZero() {
		d.first = now
	}
	at := now.Add(quiet)
	if limit := d.first.Add(d.maxWait); at.After(limit) {
		at = limit
	}
	if gap := d.lastRun.Add(d.minGap); at.Before(gap) {
		at = gap
	}
	if d.timer == nil {
		d.timer = time.AfterFunc(at.Sub(now), d.run)
	} else {
		d.timer.Reset(at.Sub(now))
	}
}

func (d *debouncer) run() {
	d.mu.Lock()
	if d.running || d.stopped {
		d.mu.Unlock()
		return
	}
	d.running = true
	d.first = time.Time{}
	d.mu.Unlock()

	d.fire()

	d.mu.Lock()
	d.running = false
	d.lastRun = time.Now()
	again := d.again
	d.again = false
	d.mu.Unlock()
	if again {
		d.poke()
	}
}

func (d *debouncer) stop() {
	d.mu.Lock()
	d.stopped = true
	if d.timer != nil {
		d.timer.Stop()
	}
	d.mu.Unlock()
}

// catalogRefresher runs the catalog refresh inside the daemon: at startup, and whenever the
// catalog watch sees a binary appear. Home files go through the refresh (rewrite daemon.conf, then
// the in-process reload); system files through the guards' re-sync, no reload.
type catalogRefresher struct {
	configPath string
	reload     func() *daemonconfig.Config // the serialized SIGHUP path; nil when it failed
	resync     func()
	editActive func() bool
	listUsers  func() ([]install.User, error)

	mu   sync.Mutex
	cfg  *daemonconfig.Config // the running config: Raw for the compare-and-swap, the plan's source
	plan []*watchPattern

	refreshes *debouncer
	resyncs   *debouncer
	wakeMu    sync.Mutex
	wake      int // eventfd: rebuild the watches, or stop; -1 once closed
	done      chan struct{}
	stopOnce  sync.Once
}

func newCatalogRefresher(configPath string, cfg *daemonconfig.Config, reload func() *daemonconfig.Config,
	resync func(), editActive func() bool) (*catalogRefresher, error) {
	wake, err := unix.Eventfd(0, unix.EFD_CLOEXEC|unix.EFD_NONBLOCK)
	if err != nil {
		return nil, fmt.Errorf("catalog watch: %w", err)
	}
	r := &catalogRefresher{configPath: configPath, reload: reload, resync: resync, editActive: editActive,
		listUsers: install.ListUsers, wake: wake, done: make(chan struct{})}
	r.refreshes = &debouncer{quiet: refreshQuiet, maxWait: refreshMaxWait, minGap: refreshMinGap, fire: r.refresh}
	r.resyncs = &debouncer{quiet: resyncQuiet, maxWait: resyncMaxWait, minGap: resyncMinGap, fire: resync}
	r.setConfig(cfg)
	return r, nil
}

// start watches and runs the startup refresh.
func (r *catalogRefresher) start() {
	go r.watchLoop()
	r.refreshes.now()
}

func (r *catalogRefresher) stop() {
	if r == nil {
		return
	}
	r.stopOnce.Do(func() {
		close(r.done)
		r.refreshes.stop()
		r.resyncs.stop()
		r.poke()
	})
}

// setConfig records the config the daemon now runs (start, every reload) and rebuilds the watches.
func (r *catalogRefresher) setConfig(cfg *daemonconfig.Config) {
	users, err := r.listUsers()
	if err != nil {
		log.Warnf("daemon: catalog watch: listing users (%v) — watching nothing until the next reload", err)
	}
	plan := catalogWatchPlan(cfg, users)
	r.mu.Lock()
	r.cfg, r.plan = cfg, plan
	r.mu.Unlock()
	r.poke()
}

func (r *catalogRefresher) poke() {
	var one uint64 = 1
	r.wakeMu.Lock()
	defer r.wakeMu.Unlock()
	if r.wake >= 0 {
		_, _ = unix.Write(r.wake, (*[8]byte)(unsafe.Pointer(&one))[:])
	}
}

// watchLoop owns the inotify instance, rebuilt on every config change and overflow.
func (r *catalogRefresher) watchLoop() {
	defer func() {
		r.wakeMu.Lock()
		_ = unix.Close(r.wake)
		r.wake = -1
		r.wakeMu.Unlock()
	}()
	for {
		select {
		case <-r.done:
			return
		default:
		}
		if !r.watchOnce() {
			// inotify unavailable: the 30s re-sync and SIGHUP remain.
			select {
			case <-r.done:
				return
			case <-time.After(time.Minute):
			}
		}
	}
}

// watchOnce runs one inotify instance until a rebuild or stop is asked for. false: it could not be
// created.
func (r *catalogRefresher) watchOnce() bool {
	fd, err := unix.InotifyInit1(unix.IN_CLOEXEC | unix.IN_NONBLOCK)
	if err != nil {
		log.Errorf("daemon: catalog watch unavailable (%v) — new app versions wait for a reload", err)
		return false
	}
	defer unix.Close(fd)
	drain(r.wake)
	set := newWatchSet(
		func(p string) (int32, error) {
			wd, err := unix.InotifyAddWatch(fd, p, watchMask)
			return int32(wd), err //nolint:gosec // inotify watch descriptors are int32 in the event ABI
		},
		func(wd int32) { _, _ = unix.InotifyRmWatch(fd, uint32(wd)) }, //nolint:gosec // see above
		time.Now)
	r.mu.Lock()
	plan := r.plan
	r.mu.Unlock()
	set.build(plan)
	log.Infof("daemon: catalog watch: %d pattern(s), %d director(ies) watched", len(plan), len(set.dirs))

	buf := make([]byte, 64*1024)
	fds := []unix.PollFd{{Fd: int32(fd), Events: unix.POLLIN}, {Fd: int32(r.wake), Events: unix.POLLIN}} //nolint:gosec // fds fit int32
	for {
		if _, err := unix.Poll(fds, int(time.Minute/time.Millisecond)); err != nil && !errors.Is(err, unix.EINTR) {
			log.Errorf("daemon: catalog watch: %v", err)
			return true
		}
		if fds[1].Revents != 0 {
			return true // rebuild or stop
		}
		if fds[0].Revents != 0 {
			readEvents(fd, buf, set.handle)
		}
		set.expire()
		refresh, resync, rebuild := set.take()
		if refresh {
			r.refreshes.poke()
		}
		if resync {
			r.resyncs.poke()
		}
		if rebuild {
			return true
		}
	}
}

func drain(fd int) {
	var b [8]byte
	_, _ = unix.Read(fd, b[:])
}

// readEvents parses every queued inotify event.
func readEvents(fd int, buf []byte, handle func(wd int32, mask uint32, name string)) {
	for {
		n, err := unix.Read(fd, buf)
		if n <= 0 || err != nil {
			return
		}
		for off := 0; off+unix.SizeofInotifyEvent <= n; {
			ev := (*unix.InotifyEvent)(unsafe.Pointer(&buf[off]))
			end := off + unix.SizeofInotifyEvent + int(ev.Len)
			if end > n {
				return
			}
			name := strings.TrimRight(string(buf[off+unix.SizeofInotifyEvent:end]), "\x00")
			handle(ev.Wd, ev.Mask, name)
			off = end
		}
	}
}

// refresh re-expands the catalog sections of the running config and applies the result through
// the in-process reload. Admission is exactly `install --update-catalog-only`'s (RefreshCatalog).
// Any failure leaves the running config, and daemon.conf, as they were.
func (r *catalogRefresher) refresh() {
	defer func() {
		if p := recover(); p != nil {
			log.Errorf("daemon: catalog refresh failed (%v) — keeping the running configuration", p)
		}
	}()
	// A reload revokes a live edit-protected grant: wait for the session to end instead.
	for r.editActive() {
		select {
		case <-r.done:
			return
		case <-time.After(time.Second):
		}
	}
	r.mu.Lock()
	cfg := r.cfg
	r.mu.Unlock()
	if err := r.refreshOnce(cfg); err != nil {
		log.Errorf("daemon: catalog refresh: %v — keeping the running configuration", err)
	}
}

func (r *catalogRefresher) refreshOnce(cfg *daemonconfig.Config) error {
	cur, err := os.ReadFile(r.configPath)
	if err != nil {
		return err
	}
	// Compare-and-swap: only a file that is still the running config is rewritten. An operator's
	// edit not yet reloaded is theirs to apply.
	if !bytes.Equal(cur, cfg.Raw) {
		log.Warnf("daemon: catalog refresh skipped: %s changed since the daemon loaded it — reload "+
			"(systemctl reload app-listener-daemon) to apply it", r.configPath)
		return nil
	}
	users, err := r.listUsers()
	if err != nil {
		return fmt.Errorf("listing users: %w", err)
	}
	live := func(res *daemonconfig.Resource, expand func() []install.BinaryRule) ([]install.BinaryRule, bool, error) {
		return expand(), res.NeedEncryption, nil
	}
	// A line the operator added that the catalog doesn't produce stays while its file exists.
	text, changes, err := install.RefreshCatalog(string(cur), cfg, users,
		install.RefreshOptions{Live: true, KeepExisting: true, Scan: live})
	if errors.Is(err, install.ErrNoCatalogMatch) || (err == nil && text == string(cur)) {
		return nil
	}
	if err != nil {
		return err
	}
	next, err := daemonconfig.Parse([]byte(text))
	if err != nil {
		return fmt.Errorf("the refreshed configuration does not parse: %w", err)
	}
	if sameConfig(next, cfg) {
		return nil // only formatting moved: no reload for nothing
	}
	if err := r.swap(cur, []byte(text)); err != nil {
		return err
	}
	for _, c := range changes {
		logRefreshChange(os.Stderr, c)
	}
	if r.reload() == nil {
		// The running config stays: put the file back so it keeps describing it.
		if err := r.swap([]byte(text), cur); err != nil {
			log.Errorf("daemon: catalog refresh: restoring %s after a failed reload: %v", r.configPath, err)
		}
		return errors.New("the reload of the refreshed configuration failed")
	}
	return nil
}

func sameConfig(a, b *daemonconfig.Config) bool {
	x, y := *a, *b
	x.Raw, y.Raw = nil, nil
	return reflect.DeepEqual(x, y)
}

// swap writes next over the config only while it still holds want.
func (r *catalogRefresher) swap(want, next []byte) error {
	cur, err := os.ReadFile(r.configPath)
	if err != nil {
		return err
	}
	if !bytes.Equal(cur, want) {
		return fmt.Errorf("%s changed while being refreshed — not overwriting it", r.configPath)
	}
	return os.WriteFile(r.configPath, next, 0o600)
}

// logRefreshChange prints one refreshed section as a journal line.
func logRefreshChange(w io.Writer, c install.SectionChange) {
	list := func(paths []string) string {
		if len(paths) == 0 {
			return "-"
		}
		return logging.SanitizeText(strings.Join(paths, ","))
	}
	fmt.Fprintf(w, "%sDAEMON catalog-refresh  resource=%s  admitted=%s  dropped=%s\n", syslogInfo,
		logging.SanitizeText(c.Section), list(c.Admitted), list(c.Dropped))
}
