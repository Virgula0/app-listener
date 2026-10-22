package daemon

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/steamrt"
	"github.com/Virgula0/app-listener/internal/usecase"
)

const (
	runtimeCheckEvery = 15 * time.Second
	// pinRetryAfter spaces setup.sh runs that failed or left the pins stale.
	pinRetryAfter = 30 * time.Minute
)

// keptRuntime is one catalog Runtimes tree whose parent is a lib_dir of the running config.
type keptRuntime struct {
	rt steamrt.Runtime
	// owners are the lib_dir's lines: the archive's creator must be one of their updaters.
	owners []string
}

// runtimeKeeper unpacks a pending Steam runtime archive the trust guard proves the client wrote, and
// re-pins a runtime whose pins went stale (internal/steamrt). Archives are judged on a timer and,
// deterministically, right before the prune forgets a freed updater (guard.OnBeforeForgetFreed):
// the client frees its creating binary as soon as it restarts.
type runtimeKeeper struct {
	uc    usecase.DaemonUseCase
	trust *guard.TrustGuard
	// mu serializes runs: the timer's and the prune hook's.
	mu       sync.Mutex
	cfgMu    sync.Mutex
	runtimes []keptRuntime
	reported map[string]string
	pinRetry map[string]time.Time
	stop     chan struct{}
	unhook   func()
}

func startRuntimeKeeper(uc usecase.DaemonUseCase, trust *guard.TrustGuard, cfg *daemonconfig.Config) *runtimeKeeper {
	k := &runtimeKeeper{uc: uc, trust: trust, reported: map[string]string{}, pinRetry: map[string]time.Time{},
		stop: make(chan struct{})}
	k.setConfig(cfg)
	k.unhook = guard.OnBeforeForgetFreed(func() { k.check(false) })
	go k.loop()
	return k
}

func (k *runtimeKeeper) close() {
	k.unhook()
	close(k.stop)
}

func (k *runtimeKeeper) setConfig(cfg *daemonconfig.Config) {
	users, err := install.ListUsers()
	if err != nil {
		log.Warnf("daemon: listing users for the Steam runtime keeper: %v", err)
	}
	kept := keptRuntimes(cfg, users)
	k.cfgMu.Lock()
	k.runtimes = kept
	k.cfgMu.Unlock()
}

// keptRuntimes pairs every user's catalog Runtimes with the read-only lib_dir holding it, when
// that lib_dir is the user's own (steam.sh writes the runtime as them).
func keptRuntimes(cfg *daemonconfig.Config, users []install.User) []keptRuntime {
	libDirs := map[string]*daemonconfig.Resource{}
	for i := range cfg.Resources {
		if r := &cfg.Resources[i]; r.ReadOnly && !r.PathPending {
			libDirs[filepath.Clean(r.Path)] = r
		}
	}
	var out []keptRuntime
	for _, u := range users {
		for i := range install.Catalog {
			for _, p := range install.Catalog[i].ExpandRuntimes(u.Name, u.Home) {
				res := libDirs[filepath.Dir(p)]
				if res == nil || !ownedBy(res.Path, u.UID) {
					continue
				}
				out = append(out, keptRuntime{
					rt:     steamrt.Runtime{Dir: res.Path, Name: filepath.Base(p), UID: u.UID, GID: u.GID, Home: u.Home},
					owners: ownerPaths(res),
				})
			}
		}
	}
	return out
}

func ownedBy(path string, uid uint32) bool {
	fi, err := os.Lstat(path)
	if err != nil || !fi.IsDir() {
		return false
	}
	st, ok := fi.Sys().(*syscall.Stat_t)
	return ok && st.Uid == uid
}

func (k *runtimeKeeper) loop() {
	t := time.NewTicker(runtimeCheckEvery)
	defer t.Stop()
	for {
		select {
		case <-k.stop:
			return
		case <-t.C:
			k.check(true)
		}
	}
}

// check unpacks every pending runtime and, with pins, re-pins the stale ones.
func (k *runtimeKeeper) check(pins bool) {
	k.cfgMu.Lock()
	kept := k.runtimes
	k.cfgMu.Unlock()
	k.mu.Lock()
	defer k.mu.Unlock()
	for i := range kept {
		kr := &kept[i]
		pending, err := kr.rt.Pending()
		if err != nil {
			k.report(kr, "check", "daemon: checking the Steam runtime "+k.path(kr)+": "+err.Error())
			continue
		}
		if pending {
			if k.unpack(kr) && pins {
				k.pin(kr)
			}
			continue
		}
		if pins && kr.rt.PinsStale() && time.Now().After(k.pinRetry[k.path(kr)]) {
			k.pin(kr)
		}
	}
}

func (k *runtimeKeeper) path(kr *keptRuntime) string { return filepath.Join(kr.rt.Dir, kr.rt.Name) }

func (k *runtimeKeeper) grant(kr *keptRuntime) steamrt.Grant {
	return func(fn func() error) error { return k.uc.WithRuntimeAccess(kr.rt.Dir, fn) }
}

func (k *runtimeKeeper) unpack(kr *keptRuntime) bool {
	prove := func(f *os.File) bool {
		dev, ino, err := ebpf.StatFile(f)
		return err == nil && k.trust.CreatedByUpdater(f, guard.GuardInodeKey{Dev: dev, Ino: ino}, kr.owners)
	}
	err := kr.rt.Unpack(prove, k.grant(kr))
	switch {
	case err == nil:
		delete(k.reported, k.path(kr)+"|unpack")
		log.Infof("daemon: unpacked the Steam runtime update into %s — start Steam again", k.path(kr))
		return true
	case errors.Is(err, steamrt.ErrProvenance):
		k.report(kr, "unpack", "daemon: a Steam runtime update is pending in "+k.path(kr)+", but its archive was "+
			"not proven written by Steam's updater since this daemon started — not unpacked "+
			"(see docs/daemon.md, Steam runtime updates)")
	default:
		k.report(kr, "unpack", "daemon: unpacking the Steam runtime update into "+k.path(kr)+": "+err.Error())
	}
	return false
}

func (k *runtimeKeeper) pin(kr *keptRuntime) {
	err := kr.rt.Pin(k.grant(kr))
	if err == nil && kr.rt.PinsStale() {
		err = errors.New("setup.sh left the pins stale")
	}
	if err != nil {
		k.pinRetry[k.path(kr)] = time.Now().Add(pinRetryAfter)
		k.report(kr, "pin", "daemon: pinning the Steam runtime "+k.path(kr)+" libraries: "+err.Error())
		return
	}
	delete(k.reported, k.path(kr)+"|pin")
	log.Infof("daemon: pinned the Steam runtime %s libraries (setup.sh)", k.path(kr))
}

// report warns once per runtime, step and message: the timer retries every runtimeCheckEvery.
func (k *runtimeKeeper) report(kr *keptRuntime, step, msg string) {
	key := k.path(kr) + "|" + step
	if k.reported[key] == msg {
		return
	}
	k.reported[key] = msg
	log.Warn(msg)
}
