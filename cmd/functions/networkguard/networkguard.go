package networkguard

import (
	"bufio"
	"errors"
	"fmt"
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
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/logging"
	"github.com/Virgula0/app-listener/internal/networkguard"
	"github.com/Virgula0/app-listener/internal/tui"
	"github.com/Virgula0/app-listener/internal/usecase"
)

var (
	blacklistPaths []string
	whitelistPaths []string
	eventsFlag     []string
	unsafeFlag     bool
	autoInfraFlag  bool
	noThrottleFlag bool
	headless       bool
)

var NetworkGuardCmd = &cobra.Command{
	Use:   "guard [-w BINARY... | -b BINARY...]",
	Short: "Deny network access (AF_INET/AF_INET6) with eBPF LSM",
	Long: `Deny IPv4/IPv6 network operations. With -b only the listed binaries are
blocked (safe on a desktop). With -w every binary but the listed ones is
blocked; with neither, everything is.

Whitelist mode only: --auto-infra also admits the running system network
daemons (DNS resolver, network manager), and --unsafe also blocks AF_UNIX
sockets, which breaks X11 and D-Bus clients.`,
	Example: `  sudo app-listener network guard -b /usr/bin/curl -b /usr/bin/wget
  sudo app-listener network guard -w /usr/bin/firefox --auto-infra`,
	Args: cobra.NoArgs,
	RunE: runNetworkGuard,
}

func init() {
	common.AddServeFlags(NetworkGuardCmd)
	f := NetworkGuardCmd.Flags()
	f.StringSliceVarP(&whitelistPaths, "whitelist", "w", nil, "Binary allowed on the network (repeatable)")
	f.StringSliceVarP(&blacklistPaths, "blacklist", "b", nil, "Binary blocked from the network (repeatable)")
	f.StringSliceVarP(&eventsFlag, "events", "e", nil,
		"Event types to show: "+common.NetEventList+" (default all; display only, enforcement covers every operation)")
	common.AddHeadlessFlag(NetworkGuardCmd, &headless)
	f.BoolVar(&unsafeFlag, "unsafe", false, "Also block AF_UNIX sockets (whitelist mode; breaks X11/D-Bus clients)")
	f.BoolVar(&autoInfraFlag, "auto-infra", false, "Admit the running system network daemons, e.g. the DNS resolver (whitelist mode)")
	f.BoolVar(&noThrottleFlag, "no-throttle", false,
		"Show every event (default: at most 1 per event type and process per 250ms)")
	NetworkGuardCmd.MarkFlagsMutuallyExclusive("whitelist", "blacklist")
}

func runNetworkGuard(cmd *cobra.Command, args []string) error {
	serve, err := common.ParseServeFlags(cmd)
	if err != nil {
		return err
	}
	mode, binPaths, err := resolveGuardMode()
	if err != nil {
		return err
	}

	binaries, err := computeGuardBinaries(binPaths)
	if err != nil {
		return err
	}

	infraPaths, binaries, err := discoverInfra(mode, binaries)
	if err != nil {
		return err
	}

	parsedEvents, err := common.ParseNetEventsFlag(eventsFlag)
	if err != nil {
		return err
	}

	if unsafeFlag {
		if err := confirmUnsafe(); err != nil {
			return err
		}
	}

	printers.PrintLogo()

	for _, p := range infraPaths {
		fmt.Printf("infra allowlisted: %s\n", p)
	}

	if checkErr := common.CheckEBPF(); checkErr != nil {
		return checkErr
	}

	if checkErr := common.CheckBPFLSM(); checkErr != nil {
		return checkErr
	}

	g, guardErr := networkguard.NewNetGuard(mode, binaries, parsedEvents, unsafeFlag, !noThrottleFlag)
	if guardErr != nil {
		return fmt.Errorf("creating network guard: %w", guardErr)
	}

	ucase := usecase.NewNetworkGuardUseCase(g)

	if startErr := ucase.Start(); startErr != nil {
		return fmt.Errorf("starting network guard: %w", startErr)
	}

	defer ucase.Stop()

	if headless {
		runHeadless(ucase, !noThrottleFlag)
		return nil
	}
	if serve.Enabled {
		fan := tui.NewEventFanout(ucase.Events())
		defer fan.Stop()
		return tui.Serve(
			tui.NewNetGuardModel(fan.Local(), mode, binaries),
			tui.NewNetGuardModel(fan.Browser(), mode, binaries),
			tui.ServeOptions{
				Address: serve.Address, Username: serve.Username, Password: serve.Password,
			})
	}

	return runTUI(ucase, mode, binaries)
}

// resolveGuardMode picks the guard mode from the flags. Neither -b nor -w = whitelist with no
// allowed binaries (all network operations blocked).
func resolveGuardMode() (networkguard.Mode, []string, error) {
	switch {
	case len(blacklistPaths) > 0 && len(whitelistPaths) > 0:
		return 0, nil, errors.New("--blacklist (-b) and --whitelist (-w) are mutually exclusive")
	case len(blacklistPaths) > 0:
		if unsafeFlag {
			return 0, nil, errors.New("--unsafe is only available in whitelist mode")
		}
		if autoInfraFlag {
			return 0, nil, errors.New("--auto-infra is only available in whitelist mode")
		}
		return networkguard.ModeBlacklist, blacklistPaths, nil
	default:
		return networkguard.ModeWhitelist, whitelistPaths, nil
	}
}

func computeGuardBinaries(paths []string) ([]networkguard.BinaryEntry, error) {
	binaries := make([]networkguard.BinaryEntry, 0, len(paths))
	for _, p := range paths {
		entry, err := networkguard.ComputeBinaryEntry(p)
		if err != nil {
			return nil, fmt.Errorf("processing binary %q: %w", p, err)
		}
		binaries = append(binaries, entry)
	}
	return binaries, nil
}

func discoverInfra(mode networkguard.Mode, binaries []networkguard.BinaryEntry) (infraPaths []string, resolved []networkguard.BinaryEntry, err error) {
	if mode != networkguard.ModeWhitelist || !autoInfraFlag {
		return nil, binaries, nil
	}
	paths, err := networkguard.DiscoverInfraBinaries()
	if err != nil {
		return nil, binaries, fmt.Errorf("discovering infrastructure binaries: %w", err)
	}
	for _, p := range paths {
		entry, err := networkguard.ComputeBinaryEntry(p)
		if err != nil {
			return nil, binaries, fmt.Errorf("processing infrastructure binary %q: %w", p, err)
		}
		binaries = append(binaries, entry)
	}
	return paths, binaries, nil
}

func confirmUnsafe() error {
	fmt.Fprintf(os.Stderr, `
WARNING: --unsafe also blocks AF_UNIX sockets (used by X11, D-Bus, etc.)
for all binaries. This may break desktop applications and cause system
instability. Only proceed if you understand these risks.

Continue? [y/N] `)

	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		return errors.New("aborted by user")
	}
	return nil
}

// runHeadless prints guard events. With throttle on, the BPF layer already rate-limits per (type,
// verdict, comm) and this printer also dedupes repeats within 1s; --no-throttle prints every event.
func runHeadless(uc usecase.NetworkGuardUseCase, throttle bool) {
	// Rate-limit per (type, verdict, comm): the LSM hooks are global, so in whitelist mode every
	// blocked socket op of noisy host processes (mDNS, IDE workers) would flood the log. The
	// verdict is in the key so an allowed event never hides a refusal of the same comm.
	type eventKey struct {
		typ     string
		comm    string
		blocked bool
		reason  string
	}
	lastLog := make(map[eventKey]time.Time)

	common.DrainUntilSignal(uc.Events(), func(ev *networkguard.NetGuardEvent) {
		if throttle {
			k := eventKey{typ: ev.Type.String(), comm: ev.Comm, blocked: ev.Blocked, reason: ev.Reason}
			now := time.Now()
			if prev, seen := lastLog[k]; seen && now.Sub(prev) < time.Second {
				return
			}
			lastLog[k] = now
		}
		// A refused whitelisted process, or a process gate, appends |<reason>.
		reason := ""
		if ev.Reason != "" {
			reason = "|" + ev.Reason
		}
		log.Infof("NETGUARD|%s|%s|%s|%s|%s|%d|%d|%d|%d|%t%s",
			ev.Type.String(), logging.SanitizeText(ev.Comm), ebpf.ProtocolString(ev.Protocol),
			logging.SanitizeText(ev.SrcAddr), logging.SanitizeText(ev.DstAddr), ev.Size, ev.PID, ev.TID, ev.NetNS,
			ev.Blocked, reason)
	})
}

func runTUI(uc usecase.NetworkGuardUseCase, mode networkguard.Mode, binaries []networkguard.BinaryEntry) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	p := tea.NewProgram(
		tui.NewNetGuardModel(uc.Events(), mode, binaries),
		tea.WithAltScreen(),
	)

	_, runErr := p.Run()
	return runErr
}
