package daemon

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/guard"
	"github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/logging"
	"github.com/Virgula0/app-listener/internal/maintainer"
	"github.com/Virgula0/app-listener/internal/usecase"
)

const (
	maintainCheckEvery = 15 * time.Second
	// maintainEvery re-runs each maintainer this often: steam.sh re-pins after a host library
	// upgrade, which no lib_dir change announces. A failed run waits as long unless its lib_dirs
	// change.
	maintainEvery = 30 * time.Minute
)

// maintainedApp is one install.Maintainer the daemon runs.
type maintainedApp struct {
	install.Maintainer
	// owner is the lib_dir whose writers are the script's updaters (vetMaintainer).
	owner *daemonconfig.Resource
	masks []string
	// sig is the lib_dirs' directory times at the last run, failedSig those of a failed one.
	sig, failedSig string
	due            time.Time
}

// maintainerKeeper runs apps' maintainer scripts (internal/maintainer): on start and reload, when
// a lib_dir's entries change (the app wrote an update), and every maintainEvery. Never while the
// app runs: its writes would race the commit, and it would load a half-written runtime.
type maintainerKeeper struct {
	uc    usecase.DaemonUseCase
	trust *guard.TrustGuard
	vet   *binaryVetter
	stop  chan struct{}

	// mu guards apps and binaries, swapped by reloads.
	mu   sync.Mutex
	apps []*maintainedApp
	// binaries are every configured whitelist path: a run's commit must touch none (trust #1 lets
	// only a binary's updaters replace it, and the daemon is none).
	binaries []string
	// reported: the loop goroutine's last warning per script.
	reported map[string]string
}

func startMaintainers(uc usecase.DaemonUseCase, trust *trustManager, cfg *daemonconfig.Config) *maintainerKeeper {
	if err := maintainer.CleanStages(); err != nil {
		log.Warnf("daemon: clearing stale maintainer stages: %v", err)
	}
	k := &maintainerKeeper{uc: uc, trust: trust.tg, vet: trust.vet, stop: make(chan struct{}),
		reported: map[string]string{}}
	k.setConfig(cfg)
	go k.loop()
	return k
}

func (k *maintainerKeeper) close() { close(k.stop) }

func (k *maintainerKeeper) setConfig(cfg *daemonconfig.Config) {
	users, err := install.ListUsers()
	if err != nil {
		log.Warnf("daemon: listing users for app maintainers: %v", err)
	}
	apps := maintainedApps(cfg, users)
	var bins []string
	for i := range cfg.Resources {
		for _, list := range [][]daemonconfig.BinaryRule{cfg.Resources[i].Binaries, cfg.Resources[i].PendingBinaries} {
			for _, b := range list {
				bins = append(bins, b.Path)
				if b.Link != "" {
					bins = append(bins, b.Link)
				}
			}
		}
	}
	k.mu.Lock()
	k.apps, k.binaries = apps, bins
	k.mu.Unlock()
}

// maintainedApps keeps each user's maintainer with the lib_dirs it may change: the user's own, on
// the script's filesystem (the overlay sees no other mount below its root).
func maintainedApps(cfg *daemonconfig.Config, users []install.User) []*maintainedApp {
	var out []*maintainedApp
	for _, m := range install.Maintainers(cfg, users) {
		root := filepath.Dir(m.Script)
		a := &maintainedApp{Maintainer: m}
		a.LibDirs = nil
		for _, r := range m.LibDirs {
			if ownedBy(r.Path, m.User.UID) && sameDevice(root, r.Path) {
				a.LibDirs = append(a.LibDirs, r)
				if a.owner == nil && len(r.Binaries) > 0 {
					a.owner = r
				}
			}
		}
		if a.owner == nil {
			continue
		}
		for _, w := range entryWriters(cfg, m.Entry, m.User) {
			if strings.HasPrefix(w, root+"/") {
				a.masks = append(a.masks, w)
			}
		}
		out = append(out, a)
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

func sameDevice(a, b string) bool {
	var sa, sb syscall.Stat_t
	return syscall.Stat(a, &sa) == nil && syscall.Stat(b, &sb) == nil && sa.Dev == sb.Dev
}

func (k *maintainerKeeper) loop() {
	k.tick()
	t := time.NewTicker(maintainCheckEvery)
	defer t.Stop()
	for {
		select {
		case <-k.stop:
			return
		case <-t.C:
			k.tick()
		}
	}
}

// tick runs every maintainer that is due, or whose lib_dirs changed since its last run (unless that
// run failed on the same lib_dirs), while its app is not running. Only the loop goroutine runs it
// and touches an app's state; a reload swaps in new apps without waiting for a run.
func (k *maintainerKeeper) tick() {
	k.mu.Lock()
	apps, bins := k.apps, k.binaries
	k.mu.Unlock()
	for _, a := range apps {
		now := time.Now()
		sig := libDirSig(a)
		if (now.Before(a.due) && (sig == a.sig || sig == a.failedSig)) || appRunning(a.masks) {
			continue
		}
		a.due = now.Add(maintainEvery)
		if err := k.run(a, bins); err != nil {
			a.failedSig = libDirSig(a)
			k.report(a, err.Error())
			continue
		}
		a.sig, a.failedSig = libDirSig(a), ""
		delete(k.reported, a.Script)
	}
}

// libDirSig is the lib_dirs' directory change and modification times: an entry created, removed
// or renamed in one moves them.
func libDirSig(a *maintainedApp) string {
	var b strings.Builder
	for _, r := range a.LibDirs {
		var st syscall.Stat_t
		if syscall.Stat(r.Path, &st) == nil {
			fmt.Fprintf(&b, "%d.%d.%d.%d;", st.Mtim.Sec, st.Mtim.Nsec, st.Ctim.Sec, st.Ctim.Nsec)
		} else {
			b.WriteString("-;")
		}
	}
	return b.String()
}

func (k *maintainerKeeper) run(a *maintainedApp, binaries []string) error {
	script, sum, err := k.vet.vetMaintainer(a.owner, a.Script)
	if errors.Is(err, errScriptUnconfirmed) {
		return fmt.Errorf("daemon: %s was not run: it is not confirmed in the binary ledger; confirm it "+
			"with `sudo app-listener trust` if %s updated it", logging.SanitizeText(a.Script), a.Entry.Name)
	}
	if err != nil {
		return fmt.Errorf("daemon: opening the maintainer %s: %w", logging.SanitizeText(a.Script), err)
	}
	defer script.Close()
	libs := make([]string, len(a.LibDirs))
	for i, r := range a.LibDirs {
		libs[i] = r.Path
	}
	plan, err := maintainer.Run(&maintainer.Spec{Script: script, ScriptSum: sum, ScriptPath: a.Script, UID: a.User.UID,
		GID: a.User.GID, User: a.User.Name, Home: a.User.Home, Masks: a.masks, LibDirs: libs})
	if err != nil {
		return fmt.Errorf("daemon: running the maintainer %s: %w", logging.SanitizeText(a.Script), err)
	}
	defer plan.Close()
	n := plan.Changes()
	if n == 0 {
		return nil
	}
	for _, b := range binaries {
		if plan.Touches(b) {
			return fmt.Errorf("daemon: refusing the maintainer %s's changes: they replace the whitelisted %s",
				logging.SanitizeText(a.Script), logging.SanitizeText(b))
		}
	}
	if appRunning(a.masks) {
		return fmt.Errorf("daemon: %s started during the maintainer run: its changes wait for the next one",
			a.Entry.Name)
	}
	err = plan.Apply(func(dir string, fn func() error) error {
		return k.trust.WithRootWriter(func() error { return k.uc.WithRuntimeAccess(dir, fn) })
	})
	if err != nil {
		return fmt.Errorf("daemon: committing the maintainer %s's changes: %w", logging.SanitizeText(a.Script), err)
	}
	log.Infof("daemon: committed %d change(s) %s made to its lib_dirs; start %s again if it failed to start",
		n, logging.SanitizeText(a.Script), a.Entry.Name)
	return nil
}

// appRunning reports whether a process runs one of bins (by inode).
func appRunning(bins []string) bool {
	ids := map[[2]uint64]bool{}
	for _, b := range bins {
		var st syscall.Stat_t
		if syscall.Stat(b, &st) == nil {
			ids[[2]uint64{st.Dev, st.Ino}] = true
		}
	}
	if len(ids) == 0 {
		return false
	}
	ents, err := os.ReadDir("/proc")
	if err != nil {
		return true
	}
	for _, e := range ents {
		if _, err := strconv.Atoi(e.Name()); err != nil {
			continue
		}
		var st syscall.Stat_t
		if syscall.Stat("/proc/"+e.Name()+"/exe", &st) == nil && ids[[2]uint64{st.Dev, st.Ino}] {
			return true
		}
	}
	return false
}

// report warns once per maintainer and message: a failed run is retried every maintainEvery.
func (k *maintainerKeeper) report(a *maintainedApp, msg string) {
	if k.reported[a.Script] == msg {
		return
	}
	k.reported[a.Script] = msg
	log.Warn(msg)
}
