package networkmonitor

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/cmd/common"
	"github.com/Virgula0/app-listener/cmd/printers"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/networkmonitor"
	"github.com/Virgula0/app-listener/internal/tui"
	"github.com/Virgula0/app-listener/internal/usecase"
)

var (
	eventsFlag []string
	headless   bool
)

var NetworkMonitorCmd = &cobra.Command{
	Use:   "network-monitor <binary1> [binary2 ...]",
	Short: "Monitor network operations of specific binaries using eBPF",
	Long: `Monitor network operations (TCP, UDP, ICMP, DNS, etc.) of one or more
binary executables using eBPF tracepoints.

Positional arguments are paths to binary executables to watch.
Only network events originating from those binaries are shown.

Examples:
  app-listener network-monitor /usr/bin/bash
  app-listener network-monitor /usr/bin/curl /usr/bin/wget
  app-listener network-monitor /usr/bin/tcpdump --headless
  app-listener network-monitor /usr/sbin/nginx -e CONNECT,ACCEPT,DNS`,
	Args: cobra.MinimumNArgs(1),
	RunE: runNetworkMonitor,
}

func init() {
	NetworkMonitorCmd.Flags().StringSliceVarP(&eventsFlag, "events", "e", nil,
		"Event types to monitor (comma-separated: CONNECT,ACCEPT,SEND,RECV,CLOSE,DNS,BIND,LISTEN; default: all)")
	NetworkMonitorCmd.Flags().BoolVarP(&headless, "headless", "", false,
		"Run without TUI, print events to stderr (for testing/scripting)")
}

func runNetworkMonitor(cmd *cobra.Command, args []string) error {
	binaries := make([]networkmonitor.BinaryEntry, 0, len(args))
	for _, p := range args {
		entry, err := networkmonitor.ComputeBinaryEntry(p)
		if err != nil {
			return fmt.Errorf("processing binary %q: %w", p, err)
		}
		binaries = append(binaries, entry)
	}

	parsedEvents, err := parseNetEventsFlag(eventsFlag)
	if err != nil {
		return err
	}

	printers.PrintLogo()

	if checkErr := common.CheckEBPF(); checkErr != nil {
		return checkErr
	}

	nm, netErr := networkmonitor.NewNetworkMonitor(binaries)
	if netErr != nil {
		return fmt.Errorf("creating network monitor: %w", netErr)
	}

	ucase := usecase.NewNetworkMonitorUseCase(nm)
	ucase.SetEventTypes(parsedEvents)

	if startErr := ucase.Start(); startErr != nil {
		return fmt.Errorf("starting network monitor: %w", startErr)
	}

	defer ucase.Stop()

	if headless {
		runHeadless(ucase)
		return nil
	}

	return runTUI(ucase, binaries)
}

func parseNetEventsFlag(flag []string) ([]ebpf.NetEventType, error) {
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

func runHeadless(uc usecase.NetworkMonitorUseCase) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case ev, ok := <-uc.Events():
			if !ok {
				return
			}
			log.Infof("NETEVENT|%s|%s|%s|%s|%s|%d|%d|%d|%d|%d",
				ev.Type.String(), ev.Comm, ebpf.ProtocolString(ev.Protocol),
				ev.SrcAddr, ev.DstAddr, ev.Size, ev.PID, ev.TID, ev.NetNS, ev.CgroupID)
		case <-sig:
			return
		}
	}
}

func runTUI(uc usecase.NetworkMonitorUseCase, binaries []networkmonitor.BinaryEntry) error {
	p := tea.NewProgram(
		tui.NewNetModel(uc.Events(), binaries),
		tea.WithAltScreen(),
	)

	_, runErr := p.Run()
	return runErr
}
