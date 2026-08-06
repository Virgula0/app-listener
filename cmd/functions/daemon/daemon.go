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

	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/cmd/common"
	"github.com/Virgula0/app-listener/cmd/printers"
	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/fscrypt"
	"github.com/Virgula0/app-listener/internal/guard"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
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
	configFlag  string
	genKeyFlag  bool
	headless    bool
	blockedOnly bool
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
unprotected window.

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
	DaemonCmd.Flags().StringVarP(&configFlag, "config", "", "",
		"Path to the daemon config file (default: /etc/app-listener/daemon.conf, then daemon-samples/daemon.conf)")
	DaemonCmd.Flags().BoolVarP(&headless, "headless", "", false,
		"Run without TUI, print events to stderr (for testing/scripting)")
	DaemonCmd.Flags().BoolVarP(&blockedOnly, "blocked-only", "", false,
		"Only print blocked (denied) events, skip allowed ones (headless only)")
	DaemonCmd.Flags().BoolVarP(&genKeyFlag, "genkey", "", false,
		"Generate the fscrypt master key file and exit")
}

func runDaemon(cmd *cobra.Command, args []string) error {
	if genKeyFlag {
		return runGenKey()
	}

	printers.PrintLogo()

	// Daemon logs land in journald; keep them plain (no ANSI colors).
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		DisableLevelTruncation: true,
		PadLevelText:           true,
	})
	log.SetOutput(os.Stderr)

	if err := common.CheckEBPF(); err != nil {
		return err
	}

	configPath, err := resolveConfigPath()
	if err != nil {
		return err
	}
	cfg, err := daemonconfig.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config %s: %w", configPath, err)
	}
	if len(cfg.Resources) == 0 {
		return fmt.Errorf("config %s contains no [watch] sections", configPath)
	}
	log.Infof("daemon starting \u2014 config: %s, resources: %d", configPath, len(cfg.Resources))

	vault := fscrypt.New()

	guards, err := buildGuards(cfg.Resources)
	if err != nil {
		return err
	}

	d, err := usecase.NewDaemonUseCase(cfg.Resources, vault, guards)
	if err != nil {
		return err
	}

	if startErr := d.Start(); startErr != nil {
		// Lock back any resources that were already unlocked.
		d.Stop()
		return startErr
	}
	defer d.Stop()

	notifySystemdReady()

	if err := writePidFile(); err != nil {
		return err
	}
	defer os.Remove(pidFile)

	reload := makeReloadHandler(d, configPath)

	if headless {
		runHeadless(d, reload)
		return nil
	}

	return runTUI(d, cfg, reload)
}

// makeReloadHandler returns the SIGHUP handler: re-parse the config file,
// rebuild every guard (re-statting the binaries, so updated binaries get
// their new inodes whitelisted), and hand the batch to the usecase, which
// applies it without ever dropping protection. Any failure keeps the
// previous configuration running, exactly like the original ssh-guard.
func makeReloadHandler(d usecase.DaemonUseCase, configPath string) func() {
	return func() {
		cfg, err := daemonconfig.Load(configPath)
		if err != nil {
			log.Errorf("daemon: reload failed, keeping previous configuration: %v", err)
			return
		}
		if len(cfg.Resources) == 0 {
			log.Error("daemon: reload failed, keeping previous configuration: config contains no [watch] sections")
			return
		}

		newGuards, err := buildGuards(cfg.Resources)
		if err != nil {
			log.Errorf("daemon: reload failed, keeping previous configuration: %v", err)
			return
		}
		if err := d.Reload(cfg.Resources, newGuards); err != nil {
			log.Errorf("daemon: reload failed, keeping previous configuration: %v", err)
			return
		}
		log.Infof("daemon: configuration reloaded from %s \u2014 resources: %d", configPath, len(cfg.Resources))
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
// own executable is appended to every whitelist so the fscrypt ioctls (which
// open the watched directories by path) keep working while the guards are
// attached during shutdown locking. On failure the partially built guards
// (which are already attached to the kernel) are detached before returning.
func buildGuards(resources []daemonconfig.Resource) ([]repository.GuardRepository, error) {
	self, err := ebpf.ComputeBinaryEntry("/proc/self/exe")
	if err != nil {
		return nil, fmt.Errorf("resolving daemon executable: %w", err)
	}

	guards := make([]repository.GuardRepository, 0, len(resources))
	for _, r := range resources {
		binaries := make([]guard.BinaryEntry, 0, len(r.Binaries)+1)
		events := make(map[string][]ebpf.EventType, len(r.Binaries)+1)
		for _, b := range r.Binaries {
			entry, err := ebpf.ComputeBinaryEntry(b.Path)
			if err != nil {
				for _, built := range guards {
					built.Stop()
				}
				return nil, fmt.Errorf("processing binary %q for %s: %w", b.Path, r.Path, err)
			}
			binaries = append(binaries, entry)
			events[b.Path] = b.Events
		}
		binaries = append(binaries, self)
		events[self.Path] = nil // the daemon itself: all events

		g, err := guard.NewGuard(r.Path, guard.ModeWhitelist, binaries, true, 0,
			guard.WithBinaryEvents(events))
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
func writeEvent(w io.Writer, blockedOnly bool, ev *usecase.DaemonEvent) bool {
	if !ev.Event.Blocked && blockedOnly {
		return false
	}
	if ev.Event.Blocked {
		fmt.Fprintf(w, "%sDAEMON DENIED|%s|%s|%s|%s|%d|%d\n",
			syslogWarning, ev.Resource, ev.Event.Type.String(), ev.Event.Comm, ev.Event.Path, ev.Event.PID, ev.Event.UID)
		return true
	}
	fmt.Fprintf(w, "%sDAEMON|%s|%s|%s|%s|%d|%d\n",
		syslogInfo, ev.Resource, ev.Event.Type.String(), ev.Event.Comm, ev.Event.Path, ev.Event.PID, ev.Event.UID)
	return true
}

func runHeadless(d usecase.DaemonUseCase, reload func()) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case ev, ok := <-d.Events():
			if !ok {
				return
			}
			if !writeEvent(os.Stderr, blockedOnly, &ev) {
				continue
			}
		case s := <-sig:
			switch s {
			case syscall.SIGHUP:
				log.Info("daemon: SIGHUP received, reloading configuration")
				reload()
			case syscall.SIGINT, syscall.SIGTERM:
				log.Info("daemon: caught termination signal, shutting down")
				return
			}
		}
	}
}

func runTUI(d usecase.DaemonUseCase, cfg *daemonconfig.Config, reload func()) error {
	resources := make([]tui.DaemonResourceInfo, 0, len(cfg.Resources))
	for _, r := range cfg.Resources {
		resources = append(resources, tui.DaemonResourceInfo{
			Path:           r.Path,
			NeedEncryption: r.NeedEncryption,
			Binaries:       len(r.Binaries),
		})
	}
	p := tea.NewProgram(
		tui.NewDaemonModel(d.Events(), resources),
		tea.WithAltScreen(),
	)

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case s := <-sig:
				switch s {
				case syscall.SIGHUP:
					log.Info("daemon: SIGHUP received, reloading configuration")
					reload()
				case syscall.SIGINT, syscall.SIGTERM:
					p.Quit()
					return
				}
			case <-quit:
				return
			}
		}
	}()

	_, runErr := p.Run()
	close(quit)
	signal.Stop(sig)
	return runErr
}
