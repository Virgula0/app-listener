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
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/cmd/common"
	"github.com/Virgula0/app-listener/cmd/printers"
	"github.com/Virgula0/app-listener/internal/constants"
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
	// etcConfigPath is the system-wide config, preferred over the working-dir template.
	etcConfigPath = "/etc/app-listener/daemon.conf"
	// selfProtectDir holds the daemon's own state (daemon.conf, fscrypt.key); guarded whenever it
	// exists (selfguards.go).
	selfProtectDir   = "/etc/app-listener"
	sampleConfigPath = "daemon-samples/daemon.conf"
	// pidFilePath mirrors ssh-guard's /run/<name>.pid contract.
	pidFile = "/run/app-listener-daemon.pid"
	// bpffsMount: conventional bpffs mountpoint for LSM link pins (a SIGKILL leaves trees enforced
	// until ExecStopPost locks the vaults and the next start retires stale pins).
	bpffsMount = "/sys/fs/bpf"
)

// pinCfg: where and under which generation this instance pins LSM links. gen is fresh per start and
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

// daemonSelfBaselineEvents is the root-gated self mask every config/ephemeral guard registers its
// own binary under (guard.WithSelfAllowBinary): enough for directory-based fscrypt (kernel-keyring
// ops) and a file-vault's "already locked" read fast path, never write. Single definition so all
// construction sites and `daemon --lockdown` (which widens from exactly this baseline and restores
// it) can't drift.
var daemonSelfBaselineEvents = []ebpf.EventType{ebpf.EventOpen, ebpf.EventRead, ebpf.EventStat}

// groupUnlockSelfEvents is the baseline plus EventMknod, only for the ephemeral guard over a
// grouped encryption root. That guard stays attached through buildGuards, which pre-stages a
// file-vault's recovery sidecar with an O_CREAT|O_EXCL open
// (fscrypt.EnsureRecoverySidecarPlaceholder); creating a dentry is gated by EVENT_MKNOD alone
// (guard_path_mknod). Never use it for a permanent per-resource guard (the baseline is deliberately
// no-write).
var groupUnlockSelfEvents = []ebpf.EventType{ebpf.EventOpen, ebpf.EventRead, ebpf.EventStat, ebpf.EventMknod}

var (
	configFlag   string
	genKeyFlag   bool
	lockdownFlag bool
	checkFlag    bool
	headless     bool
	blockedOnly  bool
	pprofAddr    string
	// noLogMetadataBlocks drops metadata-only process-gate denials from the log (gatelog.go);
	// enforcement unchanged.
	noLogMetadataBlocks bool
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
With --no-log-metadata-blocks, denied metadata-only process inspections
(op=PTRACE mode=READ) are not printed at all; without it they are printed
once per caller/target pair and folded into one DAEMON DENIED-REPEAT
summary per minute.
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
	DaemonCmd.Flags().BoolVarP(&noLogMetadataBlocks, "no-log-metadata-blocks", "", false,
		"Do not log metadata-only process inspections that were denied (op=PTRACE mode=READ: another process "+
			"reading /proc/<pid> metadata of one holding guarded secrets — desktop services like the compositor, "+
			"audio server and portal do this constantly). They are still DENIED; only the log line is dropped. "+
			"Memory access (mode=ATTACH), /proc/<pid>/mem and every file denial are always logged (headless only)")
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

// startPprof serves the Go profiler on loopback only (opt-in --pprof), e.g. `go tool pprof
// http://127.0.0.1:6060/debug/pprof/heap`.
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
	// Own mux (not DefaultServeMux) so nothing else is exposed.
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

// catchLifecycleSignals registers termination and reload signals BEFORE the first fscrypt
// unlock/guard construction:
//   - SIGTERM/SIGINT: otherwise a signal during startup kills the process with no deferred Stop,
//     leaving unlocked vault keys provisioned until ExecStopPost (guards are pinned, so trees stay
//     enforced).
//   - SIGHUP: a reload (`edit-protected --set-password`, catalog-refresh hook) can arrive seconds
//     after start; SIGHUP's default is to terminate, so it would kill the daemon mid-startup.
//
// Both channels are buffered: an early signal is queued, and runDaemonUI drains hup once startup
// completes. Registered for the process lifetime; the caller defers the returned stop.
func catchLifecycleSignals() (termSig, hup chan os.Signal, stop func()) {
	termSig = make(chan os.Signal, 1)
	signal.Notify(termSig, syscall.SIGINT, syscall.SIGTERM)
	hup = make(chan os.Signal, 1)
	signal.Notify(hup, syscall.SIGHUP)
	return termSig, hup, func() {
		signal.Stop(termSig)
		signal.Stop(hup)
	}
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

	termSig, hup, stopSignals := catchLifecycleSignals()
	defer stopSignals()

	// Always-on guards over the daemon's own state (/etc/app-listener, fscrypt.key,
	// edit-auth.hash), attached synchronously HERE, before the config guards touch bpffs: otherwise
	// a killed predecessor's self-guard pins (swept by CleanupStalePins, same pin.base) would leave
	// a window with no replacements. pinstate.go's ensurePinStateFilePlaceholder keeps
	// writePinState working under the live RO self-guard. Self guards stay OUT of the usecase's
	// guard set so a SIGHUP reload's transient guard doubling stays under the kernel's per-LSM-hook
	// program cap (selfGuards.detach/attach around reload). A failed self guard is logged CRITICAL
	// but never aborts startup (config guards are the core function and compete for the same link
	// slots).
	sg := newSelfGuards()
	sg.attach(pin)
	defer sg.detach()

	d, err := startGuardedDaemonAbortable(termSig, cfg, vault, pin)
	if err != nil {
		return err
	}
	if d == nil {
		// Termination signal aborted startup; vaults were locked under the secure lockdown. Nothing
		// to unwind.
		return nil
	}
	defer d.Stop()

	// Trusted-binary/library protection: write-protect whitelisted binaries at user-writable paths
	// (#1) and enforce the library-load allowlist (#2). Never fatal; returns a no-op cleanup when
	// unavailable.
	defer startTrustGuard(cfg)()

	events := mergeDaemonEvents(d.Events(), sg.Events())

	// Record this run's config-guard pin generation for `daemon --lockdown` (pinstate.go).
	// Best-effort.
	if pinErr := writePinState(pin); pinErr != nil {
		log.Warnf("daemon: could not record pin state (%v) — `daemon --lockdown` will not be able to widen "+
			"self-access for a file-vault resource left unlocked by a crash of this run", pinErr)
	}

	notifySystemdReady()

	if err := writePidFile(); err != nil {
		return err
	}
	defer os.Remove(pidFile)

	// edit-protected control socket (only with a password configured). Best-effort: a bind failure
	// disables live editing, never the daemon.
	control := newControlManager(d)
	control.refresh()
	defer control.close()

	reload := makeReloadHandler(d, configPath, vault, pin, sg, control)
	return runDaemonUI(events, cfg, reload, termSig, hup, serve)
}

// runDaemonUI dispatches to the presentation chosen by flags: headless stderr, browser-mirrored
// TUI, or local TUI. hup is already registered (runDaemon), so a reload requested during startup is
// queued, never lost.
func runDaemonUI(events <-chan usecase.DaemonEvent, cfg *daemonconfig.Config, reload func(), termSig, hup <-chan os.Signal, serve common.ServeConfig) error {
	if headless {
		runHeadless(events, reload, termSig, hup)
		return nil
	}
	if serve.Enabled {
		// tui.Serve handles SIGINT/SIGTERM itself; leaving termSig registered is harmless
		// (tui.Serve drives teardown, then runDaemon's deferred d.Stop runs) and avoids an
		// unhandled window. Its reload is web-triggered only (tui.ServeOptions.Reload) and doesn't
		// consume hup.
		return runServedTUI(events, cfg, reload, serve)
	}
	return runTUI(events, cfg, reload, termSig, hup)
}

// startGuardedDaemonAbortable runs startup honoring a termination signal. Startup's blocking
// syscalls (fscrypt unlock, BPF attach) can't be interrupted, so on a signal it reaches the next
// consistent point, then runs the secure lockdown (Stop keeps guards attached until every vault is
// keyless). Returns (nil, nil) if aborted.
func startGuardedDaemonAbortable(termSig <-chan os.Signal, cfg *daemonconfig.Config, vault *fscrypt.Vault, pin pinCfg) (usecase.DaemonUseCase, error) {
	return awaitStartupOrSignal(termSig, func() (usecase.DaemonUseCase, error) {
		return startGuardedDaemon(cfg, vault, pin)
	})
}

// awaitStartupOrSignal races start (in a goroutine) against a termination signal; on a signal it
// waits for a consistent point, runs the secure lockdown on what start produced, and returns (nil,
// nil) = aborted.
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

// prepareDaemonStart runs the pre-startup checks common to every run mode: logging, eBPF/BPF-LSM
// availability, config load, bpffs preflight. Returns the config plus pin location/generation.
func prepareDaemonStart() (configPath string, cfg *daemonconfig.Config, pin pinCfg, err error) {
	printers.PrintLogo()
	configureDaemonLogging()

	if ebpfErr := common.CheckEBPF(); ebpfErr != nil {
		return "", nil, pinCfg{}, ebpfErr
	}
	// Enforcement uses BPF LSM hooks: refuse to start without an active bpf LSM.
	if lsmErr := common.CheckBPFLSM(); lsmErr != nil {
		return "", nil, pinCfg{}, lsmErr
	}

	configPath, cfg, err = loadDaemonConfig()
	if err != nil {
		return "", nil, pinCfg{}, err
	}
	log.Infof("daemon starting — config: %s, resources: %d", configPath, len(cfg.Resources))

	// ResolvePinBase mounts bpffs if missing and returns "" (CRITICAL log) when pinning is
	// unsupported; the daemon still runs, without SIGKILL survival.
	base := guard.ResolvePinBase(bpffsMount)
	if base == "" {
		log.Error("daemon: CRITICAL: LSM link pinning is UNAVAILABLE on this host — the guards enforce " +
			"while the daemon runs but will NOT survive a SIGKILL. Fix bpffs to restore the kill-safety guarantee.")
	}
	return configPath, cfg, pinCfg{base: base, gen: newPinGeneration()}, nil
}

// Plain log lines (no ANSI colors) for journald.
func configureDaemonLogging() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
	log.SetOutput(os.Stderr)
}

// loadDaemonConfig resolves, loads and sanity-checks the config. Failures are permanent
// (constants.ErrCriticalStartup), not worth a systemd restart.
func loadDaemonConfig() (string, *daemonconfig.Config, error) {
	configPath, err := resolveConfigPath()
	if err != nil {
		return "", nil, fmt.Errorf("%w: %w", constants.ErrCriticalStartup, err)
	}
	cfg, err := daemonconfig.Load(configPath)
	if err != nil {
		return "", nil, fmt.Errorf("%w: loading config %s: %w", constants.ErrCriticalStartup, configPath, err)
	}
	if len(cfg.Resources) == 0 {
		return "", nil, fmt.Errorf("%w: config %s contains no [watch] sections", constants.ErrCriticalStartup, configPath)
	}
	return configPath, cfg, nil
}

// makeReloadHandler returns the SIGHUP handler: re-parse the config, rebuild every guard
// (re-statting binaries so updated ones get new inodes) and hand the batch to the usecase, which
// swaps without dropping protection. Any failure keeps the previous config running.
func makeReloadHandler(d usecase.DaemonUseCase, configPath string, vault *fscrypt.Vault, pin pinCfg, sg *selfGuards, control *controlManager) func() {
	return func() {
		// A reload rebuilds every guard, so end any live edit-protected grant first (its client
		// gets EOF; the tree is read-only again before the swap).
		control.endActiveSession("configuration reload")
		// A password added/removed via `edit-protected --set-password` takes effect here:
		// start/stop the control socket.
		defer control.refresh()

		// Self guards step aside for the reload: config guards briefly run old+new together, which
		// must stay under the kernel's per-LSM-hook program cap (BPF_MAX_TRAMP_LINKS). Re-attached
		// (best effort) after the swap, on both success and keep-previous paths.
		sg.detach()

		liveGen, err := reloadOnce(d, configPath, vault, pin.base)
		if err != nil {
			log.Errorf("daemon: reload failed, keeping previous configuration: %v", err)
			sg.attach(pin)
			return
		}
		// The usecase's commit already unpinned the old generation's guards; sweep every other
		// generation, keeping the live batch. Self guards are detached here, so their pins are gone
		// and can't look stale.
		if _, cleanErr := guard.CleanupStalePins(pin.base, map[string]bool{liveGen: true}); cleanErr != nil {
			log.Warnf("daemon: could not sweep stale guard pins after reload: %v", cleanErr)
		}
		// The reload minted a fresh config-guard generation (not the self guards': pin.gen on this
		// outer variable never changes): record it, or `daemon --lockdown` would recover the
		// pre-reload generation whose pins CleanupStalePins just removed.
		if pinErr := writePinState(pinCfg{base: pin.base, gen: liveGen}); pinErr != nil {
			log.Warnf("daemon: could not record pin state after reload (%v) — `daemon --lockdown` will not be "+
				"able to widen self-access for a file-vault resource left unlocked by a crash of this run", pinErr)
		}
		sg.attach(pin)
		log.Infof("daemon: configuration reloaded from %s", configPath)
	}
}

// reloadOnce performs one SIGHUP reload: re-parse, unlock any newly added grouped vault under an
// ephemeral guard (only after a manual edit; install restarts the daemon), rebuild the guards under
// a fresh pin generation and hand them to the usecase. On failure the freshly unlocked vaults are
// locked back and the previous config keeps running. Returns the live pin generation.
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
	// Reload committed: the usecase now owns the new roots (locked on Stop); the deferred
	// pending.stop retires the ephemeral guards.
	return pin.gen, nil
}

// startGuardedDaemon brings the engine up in fail-closed order: unlock locked grouped vaults under
// ephemeral guards, re-validate the now-visible sub-paths, build the real guards, then start the
// usecase (attach -> unlock -> populate). Ephemeral guards retire once the real ones are attached
// and populated. Every error path locks freshly unlocked vaults back first.
func startGuardedDaemon(cfg *daemonconfig.Config, vault *fscrypt.Vault, pin pinCfg) (usecase.DaemonUseCase, error) {
	relockStaleVaults(cfg, vault, pin.base)

	// Retire pins a killed predecessor left. relockStaleVaults (and ExecStopPost --lockdown)
	// already locked any vault they were the last guard for, so removing them exposes only
	// encrypted data. Must precede this run's fscrypt use: a stale guard pinned by an OLDER binary
	// would deny this process's unlock/lock ioctls.
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
		// d.Stop locks back every root it unlocked (pending ones included); the ephemeral guards
		// stay attached through that, then deferred pending.stop retires them.
		d.Stop()
		return nil, startErr
	}
	// Real guards attached and populated: the ephemeral guards are done.
	pending.stop()
	return d, nil
}

// pendingGroupUnlock tracks the ephemeral guards and freshly unlocked roots from
// unlockPendingGroupRoots, so error paths can lock vaults back while those guards still protect
// them, and every path can retire them once the real guards take over.
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

// lockRoots force-flushes every root this unlock provisioned; the ephemeral guards stay attached
// until stop(), so the tree is never unlocked and unguarded even if a lock-back fails.
func (p *pendingGroupUnlock) lockRoots(vault *fscrypt.Vault) {
	for _, root := range p.roots {
		if err := vault.Lock(root, true); err != nil && !errors.Is(err, repository.ErrKeyMissing) {
			log.Errorf("daemon: could not lock %s back after an aborted startup/reload "+
				"(ephemeral guard stays attached until it is retired): %v", root, err)
		}
	}
	p.roots = nil
}

// unlockPendingGroupRoots unlocks the vault of every PathPending grouped resource so buildGuards
// can resolve real sub-path inodes. Each root is unlocked UNDER an ephemeral recursive self-only
// guard attached first (denies every reader but the root daemon; same discipline as the running
// daemon and --update-catalog-only). Returns an empty tracker when nothing is pending.
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
			// groupUnlockSelfEvents, not the baseline: buildGuards runs while this guard is
			// attached and must pre-create file-vault recovery sidecars.
			guard.WithSelfAllowBinary(self, groupUnlockSelfEvents),
			// Pin the ephemeral guard too: a SIGKILL during the unlock window must still leave the
			// root enforced.
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

// runLockdown force-locks every encryption root in the resolved config and exits: the systemd
// ExecStopPost safety net (runs after EVERY exit: stop, crash, SIGKILL, timeout), so a daemon that
// died early still ends with keyless vaults. Best-effort per root (an open file makes a force flush
// fail); failures are logged CRITICAL and it still exits 0 so it never blocks the unit.
//
// A directory root's lock is a pure keyring op needing no self-access. A file-vault root rewrites
// its own bytes in place: if already locked that needs only a read (every self binary has it), but
// if the daemon died while unlocked, re-sealing needs the write access guard.WithSelfVaultAccess
// grants a live daemon, and this process has no live *guard.Guard. recoverPinState +
// lockRootRecovering close that gap by widening the crashed daemon's own PINNED guard (pinstate.go,
// guard.pinSelfMaps, guard.WithPinnedSelfVaultAccess).
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

// lockOneRoot force-flushes root's fscrypt key, retrying briefly on EBUSY. Reports whether the
// vault ended up keyless.
func lockOneRoot(vault *fscrypt.Vault, root string) bool {
	const attempts = 50
	for i := 0; i < attempts; i++ {
		err := vault.Lock(root, true)
		if err == nil || errors.Is(err, repository.ErrKeyMissing) {
			log.Infof("lockdown: %s is locked", root)
			return true
		}
		if errors.Is(err, repository.ErrNotEncrypted) {
			log.Errorf("lockdown: CRITICAL: %s is not encrypted and daemon.conf needs to be re-checked — need_encryption no longer matches this resource's actual on-disk state (was it decrypted, restored from a backup, or otherwise modified outside the daemon?): %v", root, err)
			return false
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

// lockRootRecovering is lockOneRoot, first widening this process's self-access under root's pinned
// guard when root is a file-vault and rec names a usable recovered generation (rec.Gen == "" =
// nothing safe to act on: no record, or its PID looks alive). Directories and already-locked file
// roots need no widening (checked first: a graceful shutdown already locked and unpinned it, so the
// pinned map is usually gone and widening would only warn confusingly). A widen that runs and fails
// structurally (pin missing/foreign/not GUARD_ALLOW_ROOT) falls back to the plain attempt.
func lockRootRecovering(vault *fscrypt.Vault, root, resourcePath, pinBase string, rec pinStateRecord) bool {
	if rec.Gen == "" {
		return lockOneRoot(vault, root)
	}
	info, statErr := os.Stat(root)
	if statErr != nil || !info.Mode().IsRegular() {
		return lockOneRoot(vault, root)
	}
	if alreadyLocked, encErr := vault.IsEncrypted(root); encErr == nil && alreadyLocked {
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

// relockStaleVaults runs before the daemon unlocks anything: a root already provisioned means the
// PREVIOUS daemon exited without locking it (SIGKILL, OOM, power loss). Its pinned guards may still
// enforce but the key is stale, so lock it back before re-provisioning. Best-effort and loud, never
// fatal (aborting would leave it unlocked longer). ExecStopPost --lockdown usually did this
// already; this covers --lockdown being killed mid-widen or never run, via the same
// recoverPinState/lockRootRecovering path.
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

// runBPFCheck is `daemon --check`: verifies the BPF-LSM stack is active and every guard eBPF
// program passes THIS kernel's verifier, attaching nothing. `install` runs it on the deployed
// binary before enabling the service, so a verifier rejection fails the install cleanly instead of
// at runtime.
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

// runGenKey implements --genkey: creates the master key, asking the operator first if one exists
// (regenerating invalidates every fscrypt directory using it).
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

// confirmMasterKeyOverwrite asks on the terminal before replacing the master key (invalidates every
// fscrypt directory provisioned with it).
func confirmMasterKeyOverwrite() bool {
	fmt.Fprintf(os.Stderr, "fscrypt master key %s already exists.\n", fscrypt.MasterKeyFile)
	fmt.Fprintf(os.Stderr, "Regenerating it invalidates every fscrypt-encrypted directory using it.\n")
	fmt.Fprintf(os.Stderr, "Regenerate? [y/N]: ")

	answer, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes"
}

// resolveConfigPath: --config flag, /etc/app-listener/daemon.conf, then the daemon-samples
// template.
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

// buildOneGuard builds and attaches one resource's guard. Inputs (self, deviceSet, pin) are
// read-only and precomputed, and it touches nothing beyond r and its result, so buildGuards runs it
// concurrently.
func buildOneGuard(r *daemonconfig.Resource, self guard.BinaryEntry, deviceSet []uint32, pin pinCfg) (*guard.Guard, error) {
	binaries := make([]guard.BinaryEntry, 0, len(r.Binaries)+1)
	events := make(map[string][]ebpf.EventType, len(r.Binaries)+1)
	var deferred []daemonconfig.BinaryRule
	for _, b := range r.Binaries {
		entry, err := ebpf.ComputeBinaryEntry(b.Path)
		if err != nil {
			// Binary is in a still-locked tree (or gone): defer it; denied until resolved after
			// unlock (fail closed).
			log.Warnf("binary %q for %s not readable yet, deferring: %v", b.Path, r.Path, err)
			deferred = append(deferred, b)
			continue
		}
		binaries = append(binaries, entry)
		events[b.Path] = b.Events
	}

	// File-vault resources stage a crash-recovery sidecar before every unlock/lock transform, which
	// must exist before THIS guard attaches: once live, creating a dentry beside a single-file
	// watch root is denied (fscrypt.EnsureRecoverySidecarPlaceholder, guard_path_rename
	// destination-parent check). No-op for directories.
	if err := fscrypt.EnsureRecoverySidecarPlaceholder(r.Path); err != nil {
		log.Warnf("daemon: could not pre-create the recovery sidecar for %s (%v) — "+
			"a lock/unlock interrupted by a crash may not be recoverable", r.Path, err)
	}

	// lib_dir resource: every process must keep READING the tree (denying would break every binary
	// loading the .so) while modifications stay whitelist-gated; the write monopoly is what makes
	// the libraries trustworthy.
	mode := guard.ModeWhitelist
	if r.ReadOnly {
		mode = guard.ModeReadOnly
	}

	g, err := guard.NewGuard(r.Path, mode, binaries, true, 0,
		guard.WithBinaryEvents(events),
		guard.WithPendingBinaries(append(deferred, r.PendingBinaries...)),
		// Root-gated self access with the minimal event set the fscrypt lifecycle needs; non-root
		// executors of this binary stay denied (self-key bypass regression test).
		guard.WithSelfAllowBinary(self, daemonSelfBaselineEvents),
		// Pin LSM links so a SIGKILL leaves the tree enforced until ExecStopPost locks the vault.
		guard.WithPinning(pin.prefix(r.Path)),
		guard.WithBackingDevices(deviceSet))
	if err != nil {
		return nil, err
	}
	if pin.base != "" && g.PinDegraded() {
		log.Errorf("daemon: CRITICAL: guard for %s could not pin its LSM links — it will not survive a SIGKILL", r.Path)
	}
	return g, nil
}

// buildConcurrency bounds concurrent guard attaches (verifier work is CPU-heavy per guard): core
// count, never more than the work.
func buildConcurrency(n int) int {
	if n < 1 {
		return 1
	}
	if lim := runtime.NumCPU(); lim < n {
		return lim
	}
	return n
}

// buildGuards creates one whitelist guard per resource. The daemon's own executable is registered
// root-gated with a minimal mask (OPEN, READ) so fscrypt ioctls keep working without making it a
// universal key.
//
// Guards are built concurrently (bounded by buildConcurrency); nothing here unlocks, so the
// invariant lives in the caller: it must not unlock ANY resource until buildGuards returns (every
// guard, incl. every member of a shared-vault `watch:` group, attached). On failure every guard
// built so far (a later resource can finish before an earlier one fails) is detached, so none stays
// attached across an error return.
func buildGuards(resources []daemonconfig.Resource, pin pinCfg) ([]repository.GuardRepository, error) {
	self, err := ebpf.ComputeBinaryEntry("/proc/self/exe")
	if err != nil {
		return nil, fmt.Errorf("resolving daemon executable: %w", err)
	}

	// The raw block-device gate is device-granular (debugfs/dd/fsck on a block device bypass
	// per-path checks), so it is enforced ONCE daemon-wide: one guard carries the union of every
	// resource's backing device, the rest opt out, avoiding N duplicate, misattributed stamps.
	// Computed up front (pure function of resources); workers only read it.
	rawDevices := backingDeviceUnion(resources)

	built := make([]*guard.Guard, len(resources))
	errs := make([]error, len(resources))

	var wg sync.WaitGroup
	sem := make(chan struct{}, buildConcurrency(len(resources)))
	for i := range resources {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }()
			// Exactly one guard carries the union; the rest pass an empty set. Index-keyed (not
			// completion order) for determinism under concurrency.
			deviceSet := []uint32(nil)
			if i == 0 {
				deviceSet = rawDevices
			}
			built[i], errs[i] = buildOneGuard(&resources[i], self, deviceSet, pin)
		}(i)
	}
	wg.Wait()

	for i, buildErr := range errs {
		if buildErr == nil {
			continue
		}
		for _, g := range built {
			if g != nil {
				g.Stop()
			}
		}
		// A resource daemon.conf declares but guard.NewGuard refuses (symlink, multi-hardlink file,
		// special file, invalid path) fails identically on every run until fixed: permanent, not
		// transient.
		return nil, fmt.Errorf("%w: creating guard for %s: %w", constants.ErrCriticalStartup, resources[i].Path, buildErr)
	}

	guards := make([]repository.GuardRepository, len(resources))
	for i, g := range built {
		guards[i] = g
	}
	return guards, nil
}

// backingDeviceUnion is the deduplicated set of backing block devices in guard_fs_devices key form.
// Major-0 filesystems (tmpfs/overlay) have no device to raw-read and are skipped; un-stat-able
// paths are reported and left uncovered rather than aborting.
func backingDeviceUnion(resources []daemonconfig.Resource) []uint32 {
	seen := make(map[uint32]bool, len(resources))
	out := make([]uint32, 0, len(resources))
	for i := range resources {
		r := &resources[i]
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

// notifySystemdReady tells systemd the daemon is up (NOTIFY_SOCKET handshake).
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

// Syslog priority markers journald maps to its priority field: <4> = warning (yellow), <6> = info.
// Gives level coloring without touching the parseable DAEMON|... layout.
const (
	syslogWarning = "<4>"
	syslogInfo    = "<6>"
)

// writeEvent prints one event to w with a syslog priority marker (denials yellow in journald).
// blockedOnly drops allowed events. Reports whether it wrote.
func writeEvent(w io.Writer, blockedOnly bool, uidr *common.UIDResolver, ev *usecase.DaemonEvent) bool {
	if !ev.Event.Blocked && blockedOnly {
		return false
	}
	who := uidr.Resolve(ev.Event.UID)
	// Sanitize event-controlled fields (path, comm, resource): newlines or terminal escapes in a
	// filename must not forge journald audit lines.
	resource := logging.SanitizeText(ev.Resource)
	path := logging.SanitizeText(ev.Event.Path)
	comm := logging.SanitizeText(ev.Event.Comm)
	op := ev.Event.Type.String()
	if ev.Event.Process != "" {
		op = ev.Event.Process // PTRACE / TRACED_EXEC: no file involved
	}
	if ev.Event.Blocked {
		// commFullPath is best-effort telemetry, not identity: comm/pid come off the kernel event
		// and can be spoofed or recycled before this readlink (see checkCommSpoof); enforcement
		// keys on exe inode. "~" = unresolved (exited, unreadable, or pid reused).
		commFullPath := "~"
		if target, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", ev.Event.PID)); err == nil {
			commFullPath = logging.SanitizeText(target)
		}
		fmt.Fprintf(w, "%sDAEMON DENIED  op=%s  comm=%s  commFullPath=%s  pid=%d  uid=%s  resource=%s  path=%s\n",
			syslogWarning, op, comm, commFullPath, ev.Event.PID, who, resource, path)
		return true
	}
	fmt.Fprintf(w, "%sDAEMON ALLOWED  op=%s  comm=%s  pid=%d  uid=%s  resource=%s  path=%s\n",
		syslogInfo, op, comm, ev.Event.PID, who, resource, path)
	return true
}

func runHeadless(events <-chan usecase.DaemonEvent, reload func(), termSig, hup <-chan os.Signal) {
	uidr := common.NewUIDResolver()
	// Folds repeated metadata-only process-gate denials into one summary per window (gatelog.go);
	// enforcement unaffected.
	limiter := newGateLogLimiter()
	flush := time.NewTicker(gateLogWindow / 4)
	defer flush.Stop()
	defer limiter.Flush(os.Stderr, true)

	for {
		select {
		case <-flush.C:
			limiter.Flush(os.Stderr, false)
		case ev, ok := <-events:
			if !ok {
				return
			}
			if !admitEvent(limiter, &ev, noLogMetadataBlocks) {
				continue
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

func runTUI(events <-chan usecase.DaemonEvent, cfg *daemonconfig.Config, reload func(), termSig, hup <-chan os.Signal) error {
	p := tea.NewProgram(newDaemonModel(events, cfg), tea.WithAltScreen())

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
	return runErr
}

func runServedTUI(events <-chan usecase.DaemonEvent, cfg *daemonconfig.Config, reload func(), serve common.ServeConfig) error {
	// Two independent models (local terminal, browser mirror), each on its own fanout stream so
	// events and counters never interfere.
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
	for i := range cfg.Resources {
		r := &cfg.Resources[i]
		resources = append(resources, tui.DaemonResourceInfo{
			Path:           r.Path,
			NeedEncryption: r.NeedEncryption,
			Binaries:       len(r.Binaries),
		})
	}
	return tui.NewDaemonModel(events, resources)
}
