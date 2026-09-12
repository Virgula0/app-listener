package daemon

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/cmd/common"
	"github.com/Virgula0/app-listener/cmd/printers"
	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/logging"
	"github.com/Virgula0/app-listener/internal/repository"
	"github.com/Virgula0/app-listener/internal/tui"
	"github.com/Virgula0/app-listener/internal/usecase"
)

const (
	// etcConfigPath is the system-wide configuration location, preferred
	// over the template in the working directory.
	etcConfigPath = "/etc/app-listener/daemon.conf"
	// selfProtectDir holds the daemon's own on-disk state (daemon.conf,
	// fscrypt.key). The daemon guards it whenever it exists — see
	// selfGuards in selfguards.go.
	selfProtectDir = "/etc/app-listener"
	// sampleConfigPath is the shipped template, relative to the working
	// directory.
	sampleConfigPath = "daemon-samples/daemon.conf"
	// pidFilePath mirrors ssh-guard's /run/<name>.pid contract.
	pidFile = "/run/app-listener-daemon.pid"
	// bpffsMount is the conventional bpffs mountpoint the guard pins its LSM
	// links under, so a SIGKILL leaves the guarded trees still enforced until
	// ExecStopPost locks the vaults and the next start retires the stale pins.
	bpffsMount = "/sys/fs/bpf"
)

// pinCfg is where and under which generation this daemon instance pins its
// LSM links. base is fixed for the process; gen is fresh per start and per
// reload so CleanupStalePins can tell a killed predecessor's pins apart.
type pinCfg struct {
	base string
	gen  string
}

func (p pinCfg) prefix(resourcePath string) string {
	return guard.PinPrefix(p.base, p.gen, resourcePath) // "" when p.base == "" (pinning unavailable)
}

// newPinGeneration returns a fresh generation tag for one batch of guards.
func newPinGeneration() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

// daemonSelfBaselineEvents is the root-gated self mask every config /
// ephemeral guard registers its own binary under (guard.WithSelfAllowBinary):
// enough for the daemon's normal directory-based fscrypt lifecycle (pure
// kernel-keyring operations) and for a file-vault resource's "already
// locked" read-only fast path, but never write. A single definition here so
// every construction site — and `daemon --lockdown`'s recovery path, which
// must widen from exactly this baseline and restore back to it — can never
// drift apart from one another.
var daemonSelfBaselineEvents = []ebpf.EventType{ebpf.EventOpen, ebpf.EventRead, ebpf.EventStat}

var (
	configFlag   string
	genKeyFlag   bool
	lockdownFlag bool
	checkFlag    bool
	headless     bool
	blockedOnly  bool
	pprofAddr    string
)

var DaemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Protect encrypted (or not) directories with a whitelist daemon (fscrypt + eBPF LSM)",
	Long: `Run the ssh-guard style daemon on top of the guard's eBPF LSM engine.

Each [watch <dir>] section in the config file protects one directory:
only the listed binaries may access it (whitelist mode, recursive by
default), optionally restricted to specific event types. Resources are
expected to be encrypted with fscrypt unless need_encryption: false is
set; encrypted resources are unlocked at startup and locked again on
shutdown while the guards remain attached, so there is never an
unprotected window. A SIGTERM/SIGINT that arrives mid-startup (before the
run loop is up) is caught: startup finishes its in-flight unlock/attach,
then the same secure lockdown runs. A hard SIGKILL cannot be caught, but
the guard's LSM links are pinned to /sys/fs/bpf, so the guarded trees stay
enforced after the process dies; the systemd unit's ExecStopPost
(app-listener daemon --lockdown) then removes the vault keys, and the next
start retires the stale pins once its own guards are attached.

The config file is resolved in this order:
  1. the --config flag, if given
  2. /etc/app-listener/daemon.conf, if present
  3. daemon-samples/daemon.conf in the working directory

With --headless, events are printed to stderr (captured by journald when
run as a systemd service), like ssh-guard's syslog output. Lines carry a
syslog priority marker (<4> warning, <6> info), so journald colors denied
events in yellow exactly like ssh-guard. With --blocked-only, only
denied (blocked) attempts are printed; allowed events are suppressed.
The guard itself never changes behavior — filtering is purely presentational.

SIGHUP reloads the configuration: every resource's binary whitelist is
recomputed (by inode, so pacman/system updates that replace binaries are
picked up) and applied atomically — new guards attach before the old
ones detach, so protection is never dropped. A malformed configuration
keeps the previous one running.`,
	Args: cobra.NoArgs,
	RunE: runDaemon,
}

func init() {
	common.AddServeFlags(DaemonCmd)
	DaemonCmd.Flags().StringVarP(&configFlag, "config", "", "",
		"Path to the daemon config file (default: /etc/app-listener/daemon.conf, then daemon-samples/daemon.conf)")
	DaemonCmd.Flags().BoolVarP(&headless, "headless", "", false,
		"Run without TUI, print events to stderr (for testing/scripting)")
	DaemonCmd.Flags().BoolVarP(&blockedOnly, "blocked-only", "", false,
		"Only print blocked (denied) events, skip allowed ones (headless only)")
	DaemonCmd.Flags().BoolVarP(&genKeyFlag, "genkey", "", false,
		"Generate the fscrypt master key file and exit")
	DaemonCmd.Flags().BoolVarP(&checkFlag, "check", "", false,
		"Preflight: verify the BPF-LSM prerequisites and load every guard eBPF program "+
			"into this kernel's verifier, then exit 0 (all accepted) or non-zero (any rejected). "+
			"Attaches nothing, changes nothing. The installer runs this against the deployed "+
			"binary before enabling the service, so a kernel whose verifier rejects a guard "+
			"program fails the install cleanly instead of crash-looping (or panicking) the daemon.")
	DaemonCmd.Flags().BoolVarP(&lockdownFlag, "lockdown", "", false,
		"Force-lock every encryption root in the config and exit. Wired into the systemd unit as ExecStopPost: "+
			"systemd runs it after every exit (clean stop, crash, SIGKILL, startup timeout), so a daemon that died "+
			"before its own lockdown finished never leaves a vault unlocked.")
	DaemonCmd.Flags().StringVarP(&pprofAddr, "pprof", "", "",
		"Serve net/http/pprof on this loopback address for profiling (e.g. 127.0.0.1:6060). Refused on non-loopback addresses. Off by default.")
}

// startPprof exposes the Go profiler on a loopback address only. It is opt-in
// (--pprof) and meant for diagnosing daemon CPU/memory: `go tool pprof
// http://127.0.0.1:6060/debug/pprof/{heap,profile,goroutine}`.
func startPprof(addr string) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		log.Errorf("pprof: ignoring --pprof %q: %v", addr, err)
		return
	}
	if ip := net.ParseIP(host); ip == nil || !ip.IsLoopback() {
		log.Errorf("pprof: refusing to expose the profiler on non-loopback address %q", addr)
		return
	}
	// Own mux, not DefaultServeMux: keeps the profiler entirely inside this
	// listener (nothing else is exposed) and off any global handler table.
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	go func() {
		srv := &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
		log.Warnf("pprof: profiler on http://%s/debug/pprof/ — diagnostic only, do not leave enabled", addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Errorf("pprof: %v", err)
		}
	}()
}

func runDaemon(cmd *cobra.Command, args []string) error {
	serve, err := common.ParseServeFlags(cmd)
	if err != nil {
		return err
	}
	if genKeyFlag {
		return runGenKey()
	}
	if checkFlag {
		return runBPFCheck()
	}
	if lockdownFlag {
		runLockdown()
		return nil
	}
	if pprofAddr != "" {
		startPprof(pprofAddr)
	}

	configPath, cfg, pin, err := prepareDaemonStart()
	if err != nil {
		return err
	}

	vault := fscrypt.New()

	// Catch termination signals BEFORE the first fscrypt unlock. Without this,
	// a SIGTERM/SIGINT during startup (systemctl stop, a startup timeout,
	// Ctrl+C) hits Go's default disposition and kills the process with no
	// deferred Stop — the just-unlocked vaults would then stay provisioned in
	// the kernel keyring until ExecStopPost runs. (The guards themselves are
	// pinned, so the trees stay enforced regardless; this is about locking the
	// vault key promptly, under the same secure lockdown a clean stop uses.)
	// Registered for the whole process lifetime; the run loop selects on it.
	termSig := make(chan os.Signal, 1)
	signal.Notify(termSig, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(termSig)

	d, err := startGuardedDaemonAbortable(termSig, cfg, vault, pin)
	if err != nil {
		return err
	}
	if d == nil {
		// A termination signal aborted startup; the vaults were locked back
		// under the secure lockdown. Nothing else to unwind.
		return nil
	}
	defer d.Stop()

	// Record this run's config-guard pin generation for `daemon --lockdown`
	// to recover later (see pinstate.go) — best-effort, and deliberately
	// before the self guards below ever attach, so the very first write on a
	// fresh host never has to fight its own not-yet-existent ReadOnly guard.
	if pinErr := writePinState(pin); pinErr != nil {
		log.Warnf("daemon: could not record pin state (%v) — `daemon --lockdown` will not be able to widen "+
			"self-access for a file-vault resource left unlocked by a crash of this run", pinErr)
	}

	notifySystemdReady()

	if err := writePidFile(); err != nil {
		return err
	}
	defer os.Remove(pidFile)

	// Always-on guards over the daemon's own state (/etc/app-listener,
	// fscrypt.key). Kept OUT of the usecase so a SIGHUP reload's transient
	// guard doubling stays under the kernel's per-LSM-hook program cap (see
	// selfGuards). Attached in the background — best-effort hardening, and
	// each LSM attach is slow on hardened kernels, so it must not delay the
	// event loop or "ready". Denials are merged into the event stream.
	sg := newSelfGuards()
	events := mergeDaemonEvents(d.Events(), sg.Events())
	go sg.attach(pin)
	defer sg.detach()

	// The edit-protected control socket (only when a password is configured).
	// Best-effort like the self guards: a socket that cannot bind disables
	// live editing but never blocks the daemon's core function.
	control := newControlManager(d)
	control.refresh()
	defer control.close()

	reload := makeReloadHandler(d, configPath, vault, pin, sg, control)
	return runDaemonUI(events, cfg, reload, termSig, serve)
}

// runDaemonUI dispatches to the presentation layer chosen by the flags:
// headless stderr stream, browser-mirrored TUI, or the local terminal TUI.
func runDaemonUI(events <-chan usecase.DaemonEvent, cfg *daemonconfig.Config, reload func(), termSig <-chan os.Signal, serve common.ServeConfig) error {
	if headless {
		runHeadless(events, reload, termSig)
		return nil
	}
	if serve.Enabled {
		// tui.Serve installs its own SIGINT/SIGTERM handling; leaving termSig
		// registered too is harmless (both channels get the signal, tui.Serve
		// drives the teardown, then runDaemon's defer d.Stop runs) and avoids
		// a brief unhandled window that signal.Stop here would open.
		return runServedTUI(events, cfg, reload, serve)
	}
	return runTUI(events, cfg, reload, termSig)
}

// startGuardedDaemonAbortable runs startup while honoring a termination
// signal. Startup is a chain of blocking syscalls (fscrypt unlock, BPF
// attach) that cannot be interrupted mid-call, so on a signal it lets startup
// reach its next consistent point and then runs the secure lockdown (Stop
// keeps the guards attached until every vault is keyless). Returns (nil, nil)
// when startup was aborted this way.
func startGuardedDaemonAbortable(termSig <-chan os.Signal, cfg *daemonconfig.Config, vault *fscrypt.Vault, pin pinCfg) (usecase.DaemonUseCase, error) {
	return awaitStartupOrSignal(termSig, func() (usecase.DaemonUseCase, error) {
		return startGuardedDaemon(cfg, vault, pin)
	})
}

// awaitStartupOrSignal runs start in a goroutine and races it against a
// termination signal. On a signal it waits for start to reach a consistent
// point (its blocking syscalls cannot be interrupted mid-call), then runs the
// secure lockdown on whatever it produced and returns (nil, nil) to signal
// "aborted — do not proceed".
func awaitStartupOrSignal(termSig <-chan os.Signal, start func() (usecase.DaemonUseCase, error)) (usecase.DaemonUseCase, error) {
	type result struct {
		d   usecase.DaemonUseCase
		err error
	}
	done := make(chan result, 1)
	go func() {
		d, err := start()
		done <- result{d, err}
	}()

	select {
	case r := <-done:
		return r.d, r.err
	case <-termSig:
		log.Warn("daemon: termination signal during startup — finishing the in-flight unlock/attach, then locking every vault back")
		r := <-done
		if r.d != nil {
			r.d.Stop() // secure lockdown: guards stay attached until keyless
		} else if r.err != nil {
			log.Infof("daemon: startup had already failed and cleaned up after itself: %v", r.err)
		}
		return nil, nil
	}
}

// prepareDaemonStart runs the pre-startup checks common to every run mode:
// logo + logging, eBPF / BPF-LSM availability, config load, and the bpffs
// preflight for guard link pinning. Returns the resolved config plus the pin
// location/generation for this daemon instance.
func prepareDaemonStart() (configPath string, cfg *daemonconfig.Config, pin pinCfg, err error) {
	printers.PrintLogo()
	configureDaemonLogging()

	if ebpfErr := common.CheckEBPF(); ebpfErr != nil {
		return "", nil, pinCfg{}, ebpfErr
	}
	// The daemon wraps fs guards (and the LSM network guard), which enforce
	// through BPF LSM hooks: refuse to start without an active bpf LSM.
	if lsmErr := common.CheckBPFLSM(); lsmErr != nil {
		return "", nil, pinCfg{}, lsmErr
	}

	configPath, cfg, err = loadDaemonConfig()
	if err != nil {
		return "", nil, pinCfg{}, err
	}
	log.Infof("daemon starting — config: %s, resources: %d", configPath, len(cfg.Resources))

	// The guard pins its LSM links to bpffs so they keep enforcing if the
	// daemon is SIGKILLed. ResolvePinBase mounts bpffs if it is missing and
	// returns "" (with a CRITICAL log) when this host cannot support pinning
	// at all — the daemon still runs, just without SIGKILL survival.
	base := guard.ResolvePinBase(bpffsMount)
	if base == "" {
		log.Error("daemon: CRITICAL: LSM link pinning is UNAVAILABLE on this host — the guards enforce " +
			"while the daemon runs but will NOT survive a SIGKILL. Fix bpffs to restore the kill-safety guarantee.")
	}
	return configPath, cfg, pinCfg{base: base, gen: newPinGeneration()}, nil
}

// configureDaemonLogging keeps daemon log lines plain (no ANSI colors) so
// journald renders them cleanly.
func configureDaemonLogging() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
	log.SetOutput(os.Stderr)
}

// loadDaemonConfig resolves, loads and sanity-checks the daemon configuration.
func loadDaemonConfig() (string, *daemonconfig.Config, error) {
	configPath, err := resolveConfigPath()
	if err != nil {
		return "", nil, err
	}
	cfg, err := daemonconfig.Load(configPath)
	if err != nil {
		return "", nil, fmt.Errorf("loading config %s: %w", configPath, err)
	}
	if len(cfg.Resources) == 0 {
		return "", nil, fmt.Errorf("config %s contains no [watch] sections", configPath)
	}
	return configPath, cfg, nil
}

// makeReloadHandler returns the SIGHUP handler: re-parse the config file,
// rebuild every guard (re-statting the binaries, so updated binaries get
// their new inodes whitelisted), and hand the batch to the usecase, which
// applies it without ever dropping protection. Any failure keeps the
// previous configuration running, exactly like the original ssh-guard.
func makeReloadHandler(d usecase.DaemonUseCase, configPath string, vault *fscrypt.Vault, pin pinCfg, sg *selfGuards, control *controlManager) func() {
	return func() {
		// A reload rebuilds every guard: any live edit-protected write grant
		// is on a guard that is about to be stopped, so end the session first
		// (its client gets EOF and the tree is back to read-only before the
		// swap).
		control.endActiveSession("configuration reload")
		// A password added/removed via `edit-protected --set-password` on a
		// running daemon takes effect here: start or stop the control socket.
		defer control.refresh()

		// The self guards step aside for the reload: the config guards
		// briefly run old+new together and that peak must stay under the
		// kernel's per-LSM-hook program cap (BPF_MAX_TRAMP_LINKS). They are
		// re-attached (best effort) once the swap has settled — on both the
		// success and keep-previous-config paths.
		sg.detach()

		liveGen, err := reloadOnce(d, configPath, vault, pin.base)
		if err != nil {
			log.Errorf("daemon: reload failed, keeping previous configuration: %v", err)
			sg.attach(pin)
			return
		}
		// The pre-reload generation's guards were unpinned by the usecase's
		// commit (old guard Stop); sweep every other generation, keeping only
		// the batch that is now live. Self guards are detached here, so their
		// pins are gone and cannot be mistaken for stale.
		if _, cleanErr := guard.CleanupStalePins(pin.base, map[string]bool{liveGen: true}); cleanErr != nil {
			log.Warnf("daemon: could not sweep stale guard pins after reload: %v", cleanErr)
		}
		// The reload minted a fresh generation for the config guards (never
		// the self guards' — pin.gen on this outer variable never changes):
		// record it, or `daemon --lockdown` would keep recovering the
		// pre-reload generation, whose pins CleanupStalePins just removed.
		if pinErr := writePinState(pinCfg{base: pin.base, gen: liveGen}); pinErr != nil {
			log.Warnf("daemon: could not record pin state after reload (%v) — `daemon --lockdown` will not be "+
				"able to widen self-access for a file-vault resource left unlocked by a crash of this run", pinErr)
		}
		sg.attach(pin)
		log.Infof("daemon: configuration reloaded from %s", configPath)
	}
}

// reloadOnce performs one SIGHUP reload: re-parse, unlock any newly added
// locked grouped vault under an ephemeral guard (a manual edit \u2014 the install
// flow restarts the daemon; existing grouped resources are already unlocked
// so this is a no-op), rebuild the guards (pinned under a fresh generation)
// and hand the batch to the usecase. On any failure the freshly unlocked
// vaults are locked back and the previous configuration keeps running.
// Returns the pin generation that is live after a successful reload.
func reloadOnce(d usecase.DaemonUseCase, configPath string, vault *fscrypt.Vault, pinBase string) (string, error) {
	cfg, err := daemonconfig.Load(configPath)
	if err != nil {
		return "", err
	}
	if len(cfg.Resources) == 0 {
		return "", fmt.Errorf("config contains no [watch] sections")
	}

	pin := pinCfg{base: pinBase, gen: newPinGeneration()}

	pending, err := unlockPendingGroupRoots(cfg, vault, pin)
	if err != nil {
		return "", err
	}
	defer pending.stop()

	if resolveErr := daemonconfig.ResolvePendingPaths(cfg); resolveErr != nil {
		pending.lockRoots(vault)
		return "", resolveErr
	}
	newGuards, buildErr := buildGuards(cfg.Resources, pin)
	if buildErr != nil {
		pending.lockRoots(vault)
		return "", buildErr
	}
	if reloadErr := d.Reload(cfg.Resources, newGuards); reloadErr != nil {
		pending.lockRoots(vault)
		return "", reloadErr
	}
	// Reload committed: the new resources' roots are now owned by the
	// usecase (locked on Stop). The deferred pending.stop retires the
	// ephemeral unlock guards.
	return pin.gen, nil
}

// startGuardedDaemon brings the guard engine up in the fail-closed order:
// unlock any locked grouped vaults under ephemeral guards, re-validate the
// now-visible sub-paths, build the real per-resource guards, then start the
// usecase (attach -> unlock -> populate). The ephemeral vault-unlock guards
// are retired once the real guards are attached and their inodes populated.
// Every error path locks the freshly unlocked vaults back before returning.
func startGuardedDaemon(cfg *daemonconfig.Config, vault *fscrypt.Vault, pin pinCfg) (usecase.DaemonUseCase, error) {
	relockStaleVaults(cfg, vault, pin.base)

	// Retire pins a previous daemon left when it was killed. relockStaleVaults
	// (and the unit's ExecStopPost --lockdown) has locked any vault those pins
	// were the last guard for, so removing them now exposes only encrypted
	// data — and it must happen before this run touches fscrypt, because a
	// stale guard pinned by an OLDER daemon binary would deny this (new-inode)
	// process's unlock/lock ioctls on the guarded trees.
	if _, cleanErr := guard.CleanupStalePins(pin.base, map[string]bool{pin.gen: true}); cleanErr != nil {
		log.Warnf("daemon: could not sweep stale guard pins: %v", cleanErr)
	}

	pending, err := unlockPendingGroupRoots(cfg, vault, pin)
	if err != nil {
		return nil, err
	}
	defer pending.stop()

	if resolveErr := daemonconfig.ResolvePendingPaths(cfg); resolveErr != nil {
		pending.lockRoots(vault)
		return nil, resolveErr
	}

	guards, buildErr := buildGuards(cfg.Resources, pin)
	if buildErr != nil {
		pending.lockRoots(vault)
		return nil, buildErr
	}

	d, ucErr := usecase.NewDaemonUseCase(cfg.Resources, vault, guards)
	if ucErr != nil {
		for _, g := range guards {
			g.Stop()
		}
		pending.lockRoots(vault)
		return nil, ucErr
	}

	if startErr := d.Start(); startErr != nil {
		// d.Stop locks back every root it unlocked (the pending roots
		// included — uniqueEncryptionRoots covers them); the ephemeral
		// guards stay attached through that lockdown, then pending.stop
		// (deferred) retires them.
		d.Stop()
		return nil, startErr
	}
	// The real per-resource guards are attached and their inodes populated:
	// the ephemeral vault-unlock guards have done their job.
	pending.stop()
	return d, nil
}

// pendingGroupUnlock tracks the ephemeral guards and freshly unlocked
// encryption roots produced by unlockPendingGroupRoots, so startup / reload
// error paths can lock the vaults back (while the ephemeral guards still
// protect them) and every path can retire the ephemeral guards once the real
// per-resource guards take over.
type pendingGroupUnlock struct {
	guards []*guard.Guard
	roots  []string
}

// stop retires the ephemeral guards. Idempotent (guard.Stop is).
func (p *pendingGroupUnlock) stop() {
	for _, g := range p.guards {
		g.Stop()
	}
	p.guards = nil
}

// lockRoots force-flushes every root this unlock provisioned. The ephemeral
// guards stay attached until stop(), so the tree is never left unlocked and
// unguarded even if a lock-back fails.
func (p *pendingGroupUnlock) lockRoots(vault *fscrypt.Vault) {
	for _, root := range p.roots {
		if err := vault.Lock(root, true); err != nil && !errors.Is(err, repository.ErrKeyMissing) {
			log.Errorf("daemon: could not lock %s back after an aborted startup/reload "+
				"(ephemeral guard stays attached until it is retired): %v", root, err)
		}
	}
	p.roots = nil
}

// unlockPendingGroupRoots unlocks the fscrypt vault of every PathPending
// grouped resource so buildGuards can resolve the real sub-path inodes. Each
// root is unlocked UNDER an ephemeral recursive self-only guard attached
// first: the unlock window denies every reader except the root daemon, the
// same discipline the running daemon and the installer's --update-catalog-only
// use. Returns an empty tracker (no-op stop/lockRoots) when nothing is
// pending.
func unlockPendingGroupRoots(cfg *daemonconfig.Config, vault *fscrypt.Vault, pin pinCfg) (*pendingGroupUnlock, error) {
	seen := make(map[string]bool)
	var roots []string
	for i := range cfg.Resources {
		r := &cfg.Resources[i]
		if !r.PathPending {
			continue
		}
		root := r.EncryptionRootOrPath()
		if !seen[root] {
			seen[root] = true
			roots = append(roots, root)
		}
	}
	p := &pendingGroupUnlock{}
	if len(roots) == 0 {
		return p, nil
	}

	self, err := ebpf.ComputeBinaryEntry("/proc/self/exe")
	if err != nil {
		return nil, fmt.Errorf("resolving daemon executable for the pending-vault guard: %w", err)
	}

	for _, root := range roots {
		g, guardErr := guard.NewGuard(root, guard.ModeWhitelist, nil, true, 0,
			guard.WithSelfAllowBinary(self, daemonSelfBaselineEvents),
			// Pin the ephemeral guard too: a SIGKILL during the grouped-vault
			// unlock window must still leave the root enforced.
			guard.WithPinning(pin.prefix("ephemeral:"+root)))
		if guardErr != nil {
			p.lockRoots(vault)
			p.stop()
			return nil, fmt.Errorf("attaching ephemeral guard for locked vault %s: %w (vault left locked)", root, guardErr)
		}
		p.guards = append(p.guards, g)
		if unlockErr := vault.Unlock(root); unlockErr != nil {
			p.lockRoots(vault)
			p.stop()
			return nil, fmt.Errorf("unlocking vault %s for grouped watch paths: %w (vault left locked)", root, unlockErr)
		}
		p.roots = append(p.roots, root)
		log.Infof("daemon: unlocked encryption root %s for grouped watch paths (under an ephemeral guard)", root)
	}
	return p, nil
}

// runLockdown force-locks every encryption root in the resolved config and
// exits. It is the systemd ExecStopPost safety net: systemd runs ExecStopPost
// after EVERY exit — clean stop, crash, SIGTERM, SIGKILL, startup timeout — so
// a daemon that died before its own lockdown could finish still ends with its
// vaults keyless. Best-effort per root (a file another process holds open in
// the tree makes a force flush fail); failures are logged CRITICAL and the
// command still exits 0 so it never blocks the unit from settling.
//
// A directory root's lock is a pure kernel-keyring op that needs no
// self-access of any kind. A file-vault (regular file) root instead rewrites
// its own guarded bytes in place — if it is ALREADY locked (ciphertext on
// disk) that only needs a read, which every self binary gets unconditionally,
// but if the previous daemon died while it was still unlocked, actually
// re-sealing it needs the same write access `guard.WithSelfVaultAccess`
// grants a live daemon — and this process has no live *guard.Guard to call
// that on. recoverPinState + lockRootRecovering close that gap by widening
// the crashed daemon's own PINNED guard (see pinstate.go, guard.pinSelfMaps,
// guard.WithPinnedSelfVaultAccess) instead.
func runLockdown() {
	configureDaemonLogging()

	configPath, err := resolveConfigPath()
	if err != nil {
		log.Errorf("lockdown: %v — nothing to lock", err)
		return
	}
	cfg, err := daemonconfig.Load(configPath)
	if err != nil {
		log.Errorf("lockdown: cannot parse %s (%v) — lock the watched directories manually: fscrypt lock <dir>", configPath, err)
		return
	}

	pinBase := guard.ResolvePinBase(bpffsMount)
	rec := recoverPinState("lockdown")

	vault := fscrypt.New()
	seen := make(map[string]bool)
	locked, stillUnlocked := 0, 0
	for i := range cfg.Resources {
		r := &cfg.Resources[i]
		if !r.NeedEncryption {
			continue
		}
		root := r.EncryptionRootOrPath()
		if seen[root] {
			continue
		}
		seen[root] = true
		if lockRootRecovering(vault, root, r.Path, pinBase, rec) {
			locked++
		} else {
			stillUnlocked++
		}
	}
	log.Infof("lockdown: %d encryption root(s) locked, %d still unlocked", locked, stillUnlocked)
}

// lockOneRoot force-flushes root's fscrypt key, retrying briefly on EBUSY
// (a pinned file). Reports whether the vault ended up keyless.
func lockOneRoot(vault *fscrypt.Vault, root string) bool {
	const attempts = 50
	for i := 0; i < attempts; i++ {
		err := vault.Lock(root, true)
		if err == nil || errors.Is(err, repository.ErrKeyMissing) {
			log.Infof("lockdown: %s is locked", root)
			return true
		}
		if !errors.Is(err, repository.ErrKeyBusy) {
			log.Errorf("lockdown: CRITICAL: could not lock %s: %v", root, err)
			return false
		}
		time.Sleep(100 * time.Millisecond)
	}
	log.Errorf("lockdown: CRITICAL: %s is STILL unlocked after %d attempts — a process holds files open in it "+
		"(lsof +D %q); its plaintext stays exposed until that process exits and the vault locks", root, attempts, root)
	return false
}

// lockRootRecovering is lockOneRoot, widening this process's self-access
// under root's pinned guard first when root is a file-vault (regular file)
// target and rec names a usable recovered pin generation (rec.Gen == ""
// means recoverPinState found nothing safe to act on — no record, or the
// record's PID still looks alive). Directories, and a file root that turns
// out to already be locked (see lockFileInPlace's read-only fast path — the
// widen is requested regardless, but only ever matters when a write is
// actually attempted), need no widening at all. A widen that fails
// structurally (pin missing/foreign/not-GUARD_ALLOW_ROOT — see
// guard.WithPinnedSelfVaultAccess) falls back to the plain attempt: no worse
// than before this recovery path existed.
func lockRootRecovering(vault *fscrypt.Vault, root, resourcePath, pinBase string, rec pinStateRecord) bool {
	if rec.Gen == "" {
		return lockOneRoot(vault, root)
	}
	info, statErr := os.Stat(root)
	if statErr != nil || !info.Mode().IsRegular() {
		return lockOneRoot(vault, root)
	}

	var ok bool
	pinPrefix := guard.PinPrefix(pinBase, rec.Gen, resourcePath)
	if widenErr := guard.WithPinnedSelfVaultAccess(pinPrefix, func() error {
		ok = lockOneRoot(vault, root)
		return nil
	}); widenErr != nil {
		log.Warnf("lockdown: could not widen self-access for %s (%v) — attempting a plain lock anyway "+
			"(succeeds if it was already locked)", root, widenErr)
		return lockOneRoot(vault, root)
	}
	return ok
}

// relockStaleVaults runs before the daemon unlocks anything: an encryption
// root already provisioned at this point means the PREVIOUS daemon exited
// without locking it (SIGKILL, OOM, power loss). Its pinned guards may still
// be enforcing (that is the point of pinning), but the vault key is stale —
// lock it back now, before this run re-provisions it under its own guards.
// Best-effort and loud; never fatal (aborting would only leave it unlocked
// for longer). ExecStopPost --lockdown normally does this first anyway, so
// this is usually a no-op; it reuses the exact same pin-state recovery path
// as --lockdown (recoverPinState/lockRootRecovering) for the case where it
// is not — e.g. --lockdown itself was killed mid-widen, or the unit's
// ExecStopPost never ran at all (a manual SIGKILL outside systemd).
func relockStaleVaults(cfg *daemonconfig.Config, vault *fscrypt.Vault, pinBase string) {
	rec := recoverPinState("daemon")

	seen := make(map[string]bool)
	for i := range cfg.Resources {
		r := &cfg.Resources[i]
		if !r.NeedEncryption {
			continue
		}
		root := r.EncryptionRootOrPath()
		if seen[root] {
			continue
		}
		seen[root] = true

		provisioned, err := vault.IsProvisioned(root)
		if err != nil {
			log.Warnf("daemon: could not check the startup lock state of %s: %v", root, err)
			continue
		}
		if !provisioned {
			continue
		}
		log.Warnf("daemon: SECURITY: %s was already unlocked at startup — the previous daemon did not lock it back "+
			"(killed mid-run or mid-startup?). Locking it now before re-provisioning it under this run's guards.", root)
		lockRootRecovering(vault, root, r.Path, pinBase, rec)
	}
}

// runGenKey implements --genkey: create the master key, but when one
// already exists ask the operator on the terminal first, because
// regenerating it invalidates every fscrypt directory using it.
// runBPFCheck is `daemon --check`: the enforcement preflight. It verifies the
// BPF-LSM stack is active and that every guard eBPF program is accepted by THIS
// kernel's verifier, without attaching anything. `install` runs it against the
// freshly deployed binary before enabling the service, so a kernel whose
// verifier rejects a guard program (e.g. a prebuilt object that is too complex
// for a newer kernel) fails the install cleanly here instead of at runtime.
func runBPFCheck() error {
	if lsmErr := common.CheckBPFLSM(); lsmErr != nil {
		return fmt.Errorf("BPF-LSM preflight failed: %w", lsmErr)
	}
	if loadErr := guard.VerifyLoad(); loadErr != nil {
		log.Error("this kernel cannot run the bundled guard programs — rebuild from source on " +
			"this host (`make build`) so the eBPF is compiled against this kernel's BTF; if the " +
			"rebuilt programs still fail, enforcement is not available on this kernel yet " +
			"(`monitor` mode, kprobes/observe-only, is unaffected)")
		return fmt.Errorf("guard eBPF preflight failed: %w", loadErr)
	}
	log.Info("BPF-LSM preflight OK: LSM stack active and every guard eBPF program is accepted by this kernel")
	return nil
}

func runGenKey() error {
	exists, err := fscrypt.MasterKeyExists()
	if err != nil {
		return fmt.Errorf("checking master key: %w", err)
	}
	if exists && !confirmMasterKeyOverwrite() {
		log.Info("master key already exists — aborting, no changes made")
		return nil
	}
	if err := fscrypt.GenerateMasterKey(exists); err != nil {
		return fmt.Errorf("generating master key: %w", err)
	}
	log.Infof("master key ready at %s", fscrypt.MasterKeyFile)
	return nil
}

// confirmMasterKeyOverwrite asks the operator on the terminal whether the
// existing master key may be replaced. Regenerating a key invalidates
// every fscrypt directory provisioned with it, so the prompt is
// required before any destructive --genkey action.
func confirmMasterKeyOverwrite() bool {
	fmt.Fprintf(os.Stderr, "fscrypt master key %s already exists.\n", fscrypt.MasterKeyFile)
	fmt.Fprintf(os.Stderr, "Regenerating it invalidates every fscrypt-encrypted directory using it.\n")
	fmt.Fprintf(os.Stderr, "Regenerate? [y/N]: ")

	answer, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes"
}

// resolveConfigPath applies the documented precedence: --config flag,
// /etc/app-listener/daemon.conf, then the daemon-samples template.
func resolveConfigPath() (string, error) {
	if configFlag != "" {
		if _, err := os.Stat(configFlag); err != nil {
			return "", fmt.Errorf("config file %s: %w", configFlag, err)
		}
		return configFlag, nil
	}
	if _, err := os.Stat(etcConfigPath); err == nil {
		return etcConfigPath, nil
	}
	if _, err := os.Stat(sampleConfigPath); err == nil {
		return sampleConfigPath, nil
	}
	return "", errors.New("no config file found: pass --config, install " + etcConfigPath +
		", or keep " + sampleConfigPath + " in the working directory")
}

// buildGuards creates one whitelist guard engine per resource. The daemon's
// own executable is registered with a root-gated allow action and a minimal
// event mask (OPEN, READ) so the fscrypt ioctls (which open the watched
// directories by path) keep working while the guards are attached — without
// turning the binary into a universal key that any local user could execute
// to inherit full access to every guarded tree. On failure the partially
// built guards (which are already attached to the kernel) are detached
// before returning.
func buildGuards(resources []daemonconfig.Resource, pin pinCfg) ([]repository.GuardRepository, error) {
	self, err := ebpf.ComputeBinaryEntry("/proc/self/exe")
	if err != nil {
		return nil, fmt.Errorf("resolving daemon executable: %w", err)
	}

	// The raw block-device gate is device-granular: a block device is a whole
	// filesystem, so debugfs/dd/fsck reading it bypasses every per-path check.
	// It is enforced ONCE, daemon-wide — one guard carries the union of every
	// resource's backing device, the rest opt out — instead of each
	// per-resource guard redundantly stamping (and, in its denial logs,
	// mis-attributing) the same shared device.
	rawDevices := backingDeviceUnion(resources)

	guards := make([]repository.GuardRepository, 0, len(resources))
	for i, r := range resources {
		binaries := make([]guard.BinaryEntry, 0, len(r.Binaries)+1)
		events := make(map[string][]ebpf.EventType, len(r.Binaries)+1)
		var deferred []daemonconfig.BinaryRule
		for _, b := range r.Binaries {
			entry, err := ebpf.ComputeBinaryEntry(b.Path)
			if err != nil {
				// The binary lives in a tree that is still locked (or is
				// genuinely gone). Defer it so the guard can resolve it
				// after its resource is unlocked; until then it is not
				// whitelisted and therefore denied — fail-closed.
				log.Warnf("binary %q for %s not readable yet, deferring: %v", b.Path, r.Path, err)
				deferred = append(deferred, b)
				continue
			}
			binaries = append(binaries, entry)
			events[b.Path] = b.Events
		}

		// Exactly one guard carries the raw block-device union; the others
		// pass an empty set so the device is not stamped N times.
		deviceSet := []uint32(nil)
		if i == 0 {
			deviceSet = rawDevices
		}

		g, err := guard.NewGuard(r.Path, guard.ModeWhitelist, binaries, true, 0,
			guard.WithBinaryEvents(events),
			guard.WithPendingBinaries(append(deferred, r.PendingBinaries...)),
			// Root-gated self access with the minimal event set the fscrypt
			// lifecycle needs; guarded content reads by non-root executors of
			// this binary stay denied (see the self-key bypass regression test).
			guard.WithSelfAllowBinary(self, daemonSelfBaselineEvents),
			// Pin the LSM links to bpffs so a SIGKILL leaves this tree still
			// enforced until ExecStopPost locks the vault.
			guard.WithPinning(pin.prefix(r.Path)),
			guard.WithBackingDevices(deviceSet))
		if err != nil {
			for _, built := range guards {
				built.Stop()
			}
			return nil, fmt.Errorf("creating guard for %s: %w", r.Path, err)
		}
		if pin.base != "" && g.PinDegraded() {
			log.Errorf("daemon: CRITICAL: guard for %s could not pin its LSM links — it will not survive a SIGKILL", r.Path)
		}
		guards = append(guards, g)
	}
	return guards, nil
}

// backingDeviceUnion is the deduplicated set of block devices backing the
// configured resources, in guard_fs_devices key form. Major-0 filesystems
// (tmpfs/overlay) and paths that cannot be stat'd are skipped — the former
// have no device to raw-read, the latter are reported and left uncovered
// rather than aborting guard construction.
func backingDeviceUnion(resources []daemonconfig.Resource) []uint32 {
	seen := make(map[uint32]bool, len(resources))
	out := make([]uint32, 0, len(resources))
	for _, r := range resources {
		rdev, hasDevice, err := guard.BackingDevice(r.Path)
		if err != nil {
			log.Warnf("daemon: raw block-device gate: %v (raw access to this device is not blocked)", err)
			continue
		}
		if !hasDevice || seen[rdev] {
			continue
		}
		seen[rdev] = true
		out = append(out, rdev)
	}
	return out
}

// notifySystemdReady tells systemd (when started as a service) that the
// daemon is fully up, mirroring ssh-guard's NOTIFY_SOCKET handshake.
func notifySystemdReady() {
	socket := os.Getenv("NOTIFY_SOCKET")
	if socket == "" {
		return
	}
	conn, err := (&net.Dialer{}).DialContext(context.Background(), "unixgram", socket)
	if err != nil {
		return
	}
	defer conn.Close()
	_, _ = conn.Write([]byte("READY=1\n"))
}

func writePidFile() error {
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d\n", os.Getpid())), 0o600); err != nil {
		return fmt.Errorf("writing pid file %s: %w", pidFile, err)
	}
	return nil
}

// Syslog priority markers journald parses at the start of a stdout/stderr
// line and maps to its own priority field: <4> = warning (rendered yellow
// by journalctl), <6> = info (plain). This reproduces the level-based
// coloring ssh-guard gets from real syslog logging without touching the
// guard or the parseable DAEMON|... field layout.
const (
	syslogWarning = "<4>"
	syslogInfo    = "<6>"
)

// writeEvent prints one guard event to w, prefixed with the syslog
// priority marker so journald colors denied attempts yellow. When
// blockedOnly is set, allowed events are dropped. It reports whether the
// event was written.
func writeEvent(w io.Writer, blockedOnly bool, uidr *common.UIDResolver, ev *usecase.DaemonEvent) bool {
	if !ev.Event.Blocked && blockedOnly {
		return false
	}
	who := uidr.Resolve(ev.Event.UID)
	// Event-controlled fields (path, comm, resource) are sanitized: a
	// filename containing newlines or terminal escapes must not forge
	// audit lines in journald.
	resource := logging.SanitizeText(ev.Resource)
	path := logging.SanitizeText(ev.Event.Path)
	comm := logging.SanitizeText(ev.Event.Comm)
	if ev.Event.Blocked {
		fmt.Fprintf(w, "%sDAEMON DENIED  op=%s  comm=%s  pid=%d  uid=%s  resource=%s  path=%s\n",
			syslogWarning, ev.Event.Type.String(), comm, ev.Event.PID, who, resource, path)
		return true
	}
	fmt.Fprintf(w, "%sDAEMON ALLOWED  op=%s  comm=%s  pid=%d  uid=%s  resource=%s  path=%s\n",
		syslogInfo, ev.Event.Type.String(), comm, ev.Event.PID, who, resource, path)
	return true
}

func runHeadless(events <-chan usecase.DaemonEvent, reload func(), termSig <-chan os.Signal) {
	uidr := common.NewUIDResolver()
	hup := make(chan os.Signal, 1)
	signal.Notify(hup, syscall.SIGHUP)
	defer signal.Stop(hup)

	for {
		select {
		case ev, ok := <-events:
			if !ok {
				return
			}
			if !writeEvent(os.Stderr, blockedOnly, uidr, &ev) {
				continue
			}
		case <-hup:
			log.Info("daemon: SIGHUP received, reloading configuration")
			reload()
		case <-termSig:
			log.Info("daemon: caught termination signal, shutting down")
			return
		}
	}
}

func runTUI(events <-chan usecase.DaemonEvent, cfg *daemonconfig.Config, reload func(), termSig <-chan os.Signal) error {
	p := tea.NewProgram(newDaemonModel(events, cfg), tea.WithAltScreen())

	hup := make(chan os.Signal, 1)
	signal.Notify(hup, syscall.SIGHUP)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-hup:
				log.Info("daemon: SIGHUP received, reloading configuration")
				reload()
			case <-termSig:
				p.Quit()
				return
			case <-quit:
				return
			}
		}
	}()

	_, runErr := p.Run()
	close(quit)
	signal.Stop(hup)
	return runErr
}

func runServedTUI(events <-chan usecase.DaemonEvent, cfg *daemonconfig.Config, reload func(), serve common.ServeConfig) error {
	// Two independent daemon models: one renders on the local terminal, one
	// is mirrored to the browser. Both consume their own fanout stream so
	// events, counters and dimensions never interfere.
	fan := tui.NewEventFanout(events)
	defer fan.Stop()
	return tui.Serve(
		newDaemonModel(fan.Local(), cfg),
		newDaemonModel(fan.Browser(), cfg),
		tui.ServeOptions{
			Address:  serve.Address,
			Username: serve.Username,
			Password: serve.Password,
			Reload:   reload,
		})
}

func newDaemonModel(events <-chan usecase.DaemonEvent, cfg *daemonconfig.Config) tea.Model {
	resources := make([]tui.DaemonResourceInfo, 0, len(cfg.Resources))
	for _, r := range cfg.Resources {
		resources = append(resources, tui.DaemonResourceInfo{
			Path:           r.Path,
			NeedEncryption: r.NeedEncryption,
			Binaries:       len(r.Binaries),
		})
	}
	return tui.NewDaemonModel(events, resources)
}
