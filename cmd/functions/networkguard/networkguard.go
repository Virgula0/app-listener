package networkguard

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/cmd/common"
	"github.com/Virgula0/app-listener/cmd/entity"
	"github.com/Virgula0/app-listener/cmd/printers"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/networkguard"
	"github.com/Virgula0/app-listener/internal/tui"
)

var (
	blacklistPaths []string
	whitelistPaths []string
	eventsFlag     []string
	unsafeFlag     bool
)

var NetworkGuardCmd = &cobra.Command{
	Use:   "network-guard",
	Short: "Block network operations of specific binaries using eBPF LSM",
	Long: `Block or allow network operations (TCP, UDP, DNS, etc.) of specific
binary executables using eBPF LSM programs.

Blacklist mode (-b): block AF_INET/AF_INET6 network operations for the
listed binaries. Everything else is allowed. This mode is safe for desktop
environments since only the listed binaries are affected.

Whitelist mode (-w): block AF_INET/AF_INET6 network operations for all
binaries except the listed ones. Add --unsafe to also block AF_UNIX
sockets (used by X11, D-Bus, etc.) — this may break desktop applications.
--unsafe is only available with -w.

The two modes are mutually exclusive.

Examples:
  app-listener network-guard -b /usr/bin/curl
  app-listener network-guard -b /usr/bin/curl /usr/bin/wget -e CONNECT,SEND
  app-listener network-guard -w /usr/bin/vim -e CONNECT,LISTEN
  app-listener network-guard -w /usr/bin/python3 --unsafe`,
	Args: cobra.NoArgs,
	RunE: runNetworkGuard,
}

func init() {
	NetworkGuardCmd.Flags().StringSliceVarP(&blacklistPaths, "blacklist", "b", nil,
		"Binary paths to blacklist (block AF_INET/AF_INET6 network operations)")
	NetworkGuardCmd.Flags().StringSliceVarP(&whitelistPaths, "whitelist", "w", nil,
		"Binary paths to whitelist (allow AF_INET/AF_INET6, block everything else)")
	NetworkGuardCmd.Flags().StringSliceVarP(&eventsFlag, "events", "e", nil,
		"Event types to intercept (comma-separated: CONNECT,ACCEPT,SEND,RECV,CLOSE,DNS,BIND,LISTEN; default: all)")
	NetworkGuardCmd.Flags().BoolVarP(&entity.Headless, "headless", "", false,
		"Run without TUI, print events to stderr (for testing/scripting)")
	NetworkGuardCmd.Flags().BoolVarP(&unsafeFlag, "unsafe", "", false,
		"Also block AF_UNIX sockets (only valid with -w). May break desktop applications that use X11/D-Bus communication")
}

func runNetworkGuard(cmd *cobra.Command, args []string) error {
	hasBlacklist := len(blacklistPaths) > 0
	hasWhitelist := len(whitelistPaths) > 0

	if hasBlacklist && hasWhitelist {
		return errors.New("--blacklist (-b) and --whitelist (-w) are mutually exclusive")
	}
	if !hasBlacklist && !hasWhitelist {
		return errors.New("specify at least one binary with --blacklist (-b) or --whitelist (-w)")
	}

	if unsafeFlag && hasBlacklist {
		return errors.New("--unsafe is only available with --whitelist (-w)")
	}

	mode := networkguard.ModeBlacklist
	binPaths := blacklistPaths
	if hasWhitelist {
		mode = networkguard.ModeWhitelist
		binPaths = whitelistPaths
	}

	binaries := make([]networkguard.BinaryEntry, 0, len(binPaths))
	for _, p := range binPaths {
		entry, err := networkguard.ComputeBinaryEntry(p)
		if err != nil {
			return fmt.Errorf("processing binary %q: %w", p, err)
		}
		binaries = append(binaries, entry)
	}

	parsedEvents, err := parseNetGuardEventsFlag(eventsFlag)
	if err != nil {
		return err
	}

	if unsafeFlag {
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
	}

	printers.PrintLogo()

	if checkErr := common.CheckEBPF(); checkErr != nil {
		return checkErr
	}

	g, guardErr := networkguard.NewNetGuard(mode, binaries, parsedEvents, unsafeFlag)
	if guardErr != nil {
		return fmt.Errorf("creating network guard: %w", guardErr)
	}

	if startErr := g.Start(); startErr != nil {
		return fmt.Errorf("starting network guard: %w", startErr)
	}

	defer g.Stop()

	if entity.Headless {
		runHeadless(g)
		return nil
	}

	return runTUI(g, mode, binaries)
}

func parseNetGuardEventsFlag(flag []string) ([]ebpf.NetEventType, error) {
	if len(flag) == 0 {
		return ebpf.NetEventTypes(), nil
	}

	var parsed []ebpf.NetEventType
	for _, s := range flag {
		et, ok := ebpf.ParseNetEventType(strings.TrimSpace(s))
		if !ok {
			return nil, fmt.Errorf("unknown network event type %q (valid: CONNECT, ACCEPT, SEND, RECV, CLOSE, DNS, BIND, LISTEN)", s)
		}
		parsed = append(parsed, et)
	}
	return parsed, nil
}

func runHeadless(g *networkguard.NetGuard) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case ev, ok := <-g.Events():
			if !ok {
				return
			}
			log.Infof("NETGUARD|%s|%s|%s|%s|%s|%d|%d|%d|%d|%t",
				ev.Type.String(), ev.Comm, ebpf.ProtocolString(ev.Protocol),
				ev.SrcAddr, ev.DstAddr, ev.Size, ev.PID, ev.TID, ev.NetNS, ev.Blocked)
		case <-sig:
			return
		}
	}
}

func runTUI(g *networkguard.NetGuard, mode networkguard.Mode, binaries []networkguard.BinaryEntry) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	p := tea.NewProgram(
		tui.NewNetGuardModel(g.Events(), mode, binaries),
		tea.WithAltScreen(),
	)

	_, runErr := p.Run()
	return runErr
}
