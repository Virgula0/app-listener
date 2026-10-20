package networkmonitor

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/cmd/common"
	"github.com/Virgula0/app-listener/cmd/printers"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/logging"
	"github.com/Virgula0/app-listener/internal/networkmonitor"
	"github.com/Virgula0/app-listener/internal/tui"
	"github.com/Virgula0/app-listener/internal/usecase"
)

var (
	eventsFlag []string
	headless   bool
)

var NetworkMonitorCmd = &cobra.Command{
	Use:   "monitor BINARY...",
	Short: "Show the network operations of binaries (observe only)",
	Long:  `Show the network operations (TCP, UDP, DNS, ...) of the listed binaries.`,
	Example: `  sudo app-listener network monitor /usr/bin/curl /usr/bin/wget
  sudo app-listener network monitor /usr/sbin/nginx -e CONNECT,ACCEPT,DNS`,
	Args: cobra.MinimumNArgs(1),
	RunE: runNetworkMonitor,
}

func init() {
	common.AddServeFlags(NetworkMonitorCmd)
	NetworkMonitorCmd.Flags().StringSliceVarP(&eventsFlag, "events", "e", nil,
		"Event types to show: "+common.NetEventList+" (default all)")
	common.AddHeadlessFlag(NetworkMonitorCmd, &headless)
}

func runNetworkMonitor(cmd *cobra.Command, args []string) error {
	serve, err := common.ParseServeFlags(cmd)
	if err != nil {
		return err
	}
	binaries := make([]networkmonitor.BinaryEntry, 0, len(args))
	for _, p := range args {
		entry, entryErr := networkmonitor.ComputeBinaryEntry(p)
		if entryErr != nil {
			return fmt.Errorf("processing binary %q: %w", p, entryErr)
		}
		binaries = append(binaries, entry)
	}

	parsedEvents, err := common.ParseNetEventsFlag(eventsFlag)
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
	if serve.Enabled {
		fan := tui.NewEventFanout(ucase.Events())
		defer fan.Stop()
		return tui.Serve(
			tui.NewNetModel(fan.Local(), binaries),
			tui.NewNetModel(fan.Browser(), binaries),
			tui.ServeOptions{
				Address: serve.Address, Username: serve.Username, Password: serve.Password,
			})
	}

	return runTUI(ucase, binaries)
}

func runHeadless(uc usecase.NetworkMonitorUseCase) {
	common.DrainUntilSignal(uc.Events(), func(ev *ebpf.NetEvent) {
		log.Infof("NETEVENT|%s|%s|%s|%s|%s|%d|%d|%d|%d|%d",
			ev.Type.String(), logging.SanitizeText(ev.Comm), ebpf.ProtocolString(ev.Protocol),
			logging.SanitizeText(ev.SrcAddr), logging.SanitizeText(ev.DstAddr), ev.Size, ev.PID, ev.TID, ev.NetNS, ev.CgroupID)
	})
}

func runTUI(uc usecase.NetworkMonitorUseCase, binaries []networkmonitor.BinaryEntry) error {
	p := tea.NewProgram(
		tui.NewNetModel(uc.Events(), binaries),
		tea.WithAltScreen(),
	)

	_, runErr := p.Run()
	return runErr
}
