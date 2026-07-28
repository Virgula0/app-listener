package networkguard

import (
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
)

var NetworkGuardCmd = &cobra.Command{
	Use:   "network-guard",
	Short: "Block network operations of specific binaries using eBPF LSM",
	Long: `Block network operations (TCP, UDP, DNS, etc.) of one or more
binary executables using eBPF LSM programs.

Block network operations of listed binaries (blacklist) and optionally
allow others (whitelist) to bypass blocking. A binary cannot be in both lists.

Examples:
  app-listener network-guard -b /usr/bin/curl
  app-listener network-guard -b /usr/bin/curl /usr/bin/wget -e CONNECT,SEND
  app-listener network-guard -w /usr/bin/vim -e CONNECT,LISTEN
  app-listener network-guard -b /usr/bin/curl -w /usr/bin/ssh`,
	Args: cobra.NoArgs,
	RunE: runNetworkGuard,
}

func init() {
	NetworkGuardCmd.Flags().StringSliceVarP(&blacklistPaths, "blacklist", "b", nil,
		"Binary paths to blacklist (block network operations)")
	NetworkGuardCmd.Flags().StringSliceVarP(&whitelistPaths, "whitelist", "w", nil,
		"Binary paths to whitelist (bypass network blocking)")
	NetworkGuardCmd.Flags().StringSliceVarP(&eventsFlag, "events", "e", nil,
		"Event types to block (comma-separated: CONNECT,ACCEPT,SEND,RECV,CLOSE,DNS,BIND,LISTEN; default: all)")
	NetworkGuardCmd.Flags().BoolVarP(&entity.Headless, "headless", "", false,
		"Run without TUI, print events to stderr (for testing/scripting)")
}

func runNetworkGuard(cmd *cobra.Command, args []string) error {
	if len(blacklistPaths) == 0 && len(whitelistPaths) == 0 {
		return errors.New("specify at least one binary with --blacklist (-b) or --whitelist (-w)")
	}

	if err := checkDuplicatePaths(blacklistPaths, whitelistPaths); err != nil {
		return err
	}

	blockBinaries := make([]networkguard.BinaryEntry, 0, len(blacklistPaths))
	for _, p := range blacklistPaths {
		entry, err := networkguard.ComputeBinaryEntry(p)
		if err != nil {
			return fmt.Errorf("processing blacklist binary %q: %w", p, err)
		}
		blockBinaries = append(blockBinaries, entry)
	}

	allowBinaries := make([]networkguard.BinaryEntry, 0, len(whitelistPaths))
	for _, p := range whitelistPaths {
		entry, err := networkguard.ComputeBinaryEntry(p)
		if err != nil {
			return fmt.Errorf("processing whitelist binary %q: %w", p, err)
		}
		allowBinaries = append(allowBinaries, entry)
	}

	parsedEvents, err := parseNetGuardEventsFlag(eventsFlag)
	if err != nil {
		return err
	}

	printers.PrintLogo()

	if checkErr := common.CheckEBPF(); checkErr != nil {
		return checkErr
	}

	g, guardErr := networkguard.NewNetGuard(blockBinaries, allowBinaries, parsedEvents)
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

	return runTUI(g, blockBinaries, allowBinaries)
}

func checkDuplicatePaths(blacklist, whitelist []string) error {
	seen := make(map[string]bool, len(blacklist))
	for _, p := range blacklist {
		seen[p] = true
	}
	for _, p := range whitelist {
		if seen[p] {
			return fmt.Errorf("binary %q cannot be both blacklisted and whitelisted", p)
		}
	}
	return nil
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

func runTUI(g *networkguard.NetGuard, blockBinaries, allowBinaries []networkguard.BinaryEntry) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	p := tea.NewProgram(
		tui.NewNetGuardModel(g.Events(), blockBinaries, allowBinaries),
		tea.WithAltScreen(),
	)

	_, runErr := p.Run()
	return runErr
}
