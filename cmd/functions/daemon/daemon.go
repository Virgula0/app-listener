package daemon

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
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
	// sampleConfigPath is the shipped template, relative to the working
	// directory.
	sampleConfigPath = "daemon-samples/daemon.conf"
	// pidFilePath mirrors ssh-guard's /run/<name>.pid contract.
	pidFile = "/run/app-listener-daemon.pid"
)

var (
	configFlag   string
	genKeyFlag   bool
	lockdownFlag bool
	headless     bool
	blockedOnly  bool
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
then the same secure lockdown runs. A hard SIGKILL cannot be caught — the
systemd unit's ExecStopPost (app-listener daemon --lockdown) is the net
for that, and the next startup re-locks any vault a previous run left
provisioned.

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
	DaemonCmd.Flags().BoolVarP(&lockdownFlag, "lockdown", "", false,
		"Force-lock every encryption root in the config and exit. Wired into the systemd unit as ExecStopPost: "+
			"systemd runs it after every exit (clean stop, crash, SIGKILL, startup timeout), so a daemon that died "+
			"before its own lockdown finished never leaves a vault unlocked.")
}

func runDaemon(cmd *cobra.Command, args []string) error {
	serve, err := common.ParseServeFlags(cmd)
	if err != nil {
		return err
	}
	if genKeyFlag {
		return runGenKey()
	}
	if lockdownFlag {
		runLockdown()
		return nil
	}

	printers.PrintLogo()
	configureDaemonLogging()

	if ebpfErr := common.CheckEBPF(); ebpfErr != nil {
		return ebpfErr
	}

	// The daemon wraps fs guards (and the LSM network guard), which enforce
	// through BPF LSM hooks: refuse to start without an active bpf LSM.
	if lsmErr := common.CheckBPFLSM(); lsmErr != nil {
		return lsmErr
	}

	configPath, cfg, err := loadDaemonConfig()
	if err != nil {
		return err
	}
	log.Infof("daemon starting \u2014 config: %s, resources: %d", configPath, len(cfg.Resources))

	vault := fscrypt.New()

	// Catch termination signals BEFORE the first fscrypt unlock. Without this,
	// a SIGTERM/SIGINT during startup (systemctl stop, a startup timeout,
	// Ctrl+C) hits Go's default disposition and kills the process outright —
	// no deferred Stop runs, so the just-unlocked vaults stay provisioned in
	// the kernel keyring while the unpinned BPF guards detach on exit:
	// decrypted AND unguarded, with nothing to lock them back. Registered for
	// the whole process lifetime; the run loop selects on the same channel.
	termSig := make(chan os.Signal, 1)
	signal.Notify(termSig, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(termSig)

	d, err := startGuardedDaemonAbortable(termSig, cfg, vault)
	if err != nil {
		return err
	}
	if d == nil {
		// A termination signal aborted startup; the vaults were locked back
		// under the secure lockdown. Nothing else to unwind.
		return nil
	}
	defer d.Stop()

	notifySystemdReady()

	if err := writePidFile(); err != nil {
		return err
	}
	defer os.Remove(pidFile)

	reload := makeReloadHandler(d, configPath, vault)

	if headless {
		runHeadless(d, reload, termSig)
		return nil
	}
	if serve.Enabled {
		// tui.Serve installs its own SIGINT/SIGTERM handling; leaving termSig
		// registered too is harmless (both channels get the signal, tui.Serve
		// drives the teardown, then runDaemon's defer d.Stop runs) and avoids
		// a brief unhandled window that signal.Stop here would open.
		return runServedTUI(d, cfg, reload, serve)
	}

	return runTUI(d, cfg, reload, termSig)
}

// startGuardedDaemonAbortable runs startup while honoring a termination
// signal. Startup is a chain of blocking syscalls (fscrypt unlock, BPF
// attach) that cannot be interrupted mid-call, so on a signal it lets startup
// reach its next consistent point and then runs the secure lockdown (Stop
// keeps the guards attached until every vault is keyless). Returns (nil, nil)
// when startup was aborted this way.
func startGuardedDaemonAbortable(termSig <-chan os.Signal, cfg *daemonconfig.Config, vault *fscrypt.Vault) (usecase.DaemonUseCase, error) {
	return awaitStartupOrSignal(termSig, func() (usecase.DaemonUseCase, error) {
		return startGuardedDaemon(cfg, vault)
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
func makeReloadHandler(d usecase.DaemonUseCase, configPath string, vault *fscrypt.Vault) func() {
	return func() {
		if err := reloadOnce(d, configPath, vault); err != nil {
			log.Errorf("daemon: reload failed, keeping previous configuration: %v", err)
			return
		}
		log.Infof("daemon: configuration reloaded from %s", configPath)
	}
}

// reloadOnce performs one SIGHUP reload: re-parse, unlock any newly added
// locked grouped vault under an ephemeral guard (a manual edit \u2014 the install
// flow restarts the daemon; existing grouped resources are already unlocked
// so this is a no-op), rebuild the guards and hand the batch to the usecase.
// On any failure the freshly unlocked vaults are locked back and the previous
// configuration keeps running.
func reloadOnce(d usecase.DaemonUseCase, configPath string, vault *fscrypt.Vault) error {
	cfg, err := daemonconfig.Load(configPath)
	if err != nil {
		return err
	}
	if len(cfg.Resources) == 0 {
		return fmt.Errorf("config contains no [watch] sections")
	}

	pending, err := unlockPendingGroupRoots(cfg, vault)
	if err != nil {
		return err
	}
	defer pending.stop()

	if resolveErr := daemonconfig.ResolvePendingPaths(cfg); resolveErr != nil {
		pending.lockRoots(vault)
		return resolveErr
	}
	newGuards, buildErr := buildGuards(cfg.Resources)
	if buildErr != nil {
		pending.lockRoots(vault)
		return buildErr
	}
	if reloadErr := d.Reload(cfg.Resources, newGuards); reloadErr != nil {
		pending.lockRoots(vault)
		return reloadErr
	}
	// Reload committed: the new resources' roots are now owned by the
	// usecase (locked on Stop). The deferred pending.stop retires the
	// ephemeral unlock guards.
	return nil
}

// startGuardedDaemon brings the guard engine up in the fail-closed order:
// unlock any locked grouped vaults under ephemeral guards, re-validate the
// now-visible sub-paths, build the real per-resource guards, then start the
// usecase (attach -> unlock -> populate). The ephemeral vault-unlock guards
// are retired once the real guards are attached and their inodes populated.
// Every error path locks the freshly unlocked vaults back before returning.
func startGuardedDaemon(cfg *daemonconfig.Config, vault *fscrypt.Vault) (usecase.DaemonUseCase, error) {
	relockStaleVaults(cfg, vault)

	pending, err := unlockPendingGroupRoots(cfg, vault)
	if err != nil {
		return nil, err
	}
	defer pending.stop()

	if resolveErr := daemonconfig.ResolvePendingPaths(cfg); resolveErr != nil {
		pending.lockRoots(vault)
		return nil, resolveErr
	}

	guards, buildErr := buildGuards(cfg.Resources)
	if buildErr != nil {
		pending.lockRoots(vault)
		return nil, buildErr
	}

	d, ucErr := usecase.NewDaemonUseCase(cfg.Resources, vault, guards)
	if ucErr != nil {
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
func unlockPendingGroupRoots(cfg *daemonconfig.Config, vault *fscrypt.Vault) (*pendingGroupUnlock, error) {
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
			guard.WithSelfAllowBinary(self, []ebpf.EventType{ebpf.EventOpen, ebpf.EventRead}))
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
		if lockOneRoot(vault, root) {
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

// relockStaleVaults runs before the daemon unlocks anything: an encryption
// root already provisioned at this point means the PREVIOUS daemon exited
// without locking it (SIGKILL, OOM, power loss) — its plaintext has been
// exposed and, since that daemon's unpinned guards died with it, unguarded.
// Lock it back now; startGuards re-provisions it below with this run's guards
// already attached. Best-effort and loud; never fatal (aborting would only
// leave the stale vault unlocked for longer).
func relockStaleVaults(cfg *daemonconfig.Config, vault *fscrypt.Vault) {
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
			"(killed mid-run or mid-startup?). Its plaintext has been exposed and unguarded; locking it now before "+
			"re-provisioning it under this run's guards.", root)
		lockOneRoot(vault, root)
	}
}

// runGenKey implements --genkey: create the master key, but when one
// already exists ask the operator on the terminal first, because
// regenerating it invalidates every fscrypt directory using it.
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
func buildGuards(resources []daemonconfig.Resource) ([]repository.GuardRepository, error) {
	self, err := ebpf.ComputeBinaryEntry("/proc/self/exe")
	if err != nil {
		return nil, fmt.Errorf("resolving daemon executable: %w", err)
	}

	guards := make([]repository.GuardRepository, 0, len(resources))
	for _, r := range resources {
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

		g, err := guard.NewGuard(r.Path, guard.ModeWhitelist, binaries, true, 0,
			guard.WithBinaryEvents(events),
			guard.WithPendingBinaries(append(deferred, r.PendingBinaries...)),
			// Root-gated self access with the minimal event set the fscrypt
			// lifecycle needs; guarded content reads by non-root executors of
			// this binary stay denied (see the self-key bypass regression test).
			guard.WithSelfAllowBinary(self, []ebpf.EventType{ebpf.EventOpen, ebpf.EventRead}))
		if err != nil {
			for _, built := range guards {
				built.Stop()
			}
			return nil, fmt.Errorf("creating guard for %s: %w", r.Path, err)
		}
		guards = append(guards, g)
	}
	return guards, nil
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

func runHeadless(d usecase.DaemonUseCase, reload func(), termSig <-chan os.Signal) {
	uidr := common.NewUIDResolver()
	hup := make(chan os.Signal, 1)
	signal.Notify(hup, syscall.SIGHUP)
	defer signal.Stop(hup)

	for {
		select {
		case ev, ok := <-d.Events():
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

func runTUI(d usecase.DaemonUseCase, cfg *daemonconfig.Config, reload func(), termSig <-chan os.Signal) error {
	p := tea.NewProgram(newDaemonModel(d.Events(), cfg), tea.WithAltScreen())

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

func runServedTUI(d usecase.DaemonUseCase, cfg *daemonconfig.Config, reload func(), serve common.ServeConfig) error {
	// Two independent daemon models: one renders on the local terminal, one
	// is mirrored to the browser. Both consume their own fanout stream so
	// events, counters and dimensions never interfere.
	fan := tui.NewEventFanout(d.Events())
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
