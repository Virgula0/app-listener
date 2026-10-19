package monitor

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/cmd/common"
	"github.com/Virgula0/app-listener/cmd/printers"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/logging"
	"github.com/Virgula0/app-listener/internal/monitor"
	"github.com/Virgula0/app-listener/internal/tui"
	"github.com/Virgula0/app-listener/internal/usecase"
)

var (
	eventsFlag []string

	recursive bool
	depth     int
	headless  bool
)

var MonitorCmd = &cobra.Command{
	Use:   "monitor PATH...",
	Short: "Show file system events on files or directories (observe only)",
	Long: `Show the file system events on each PATH (a directory, or a single file), with
the process behind them. Nothing is blocked: use guard for that.`,
	Example: `  sudo app-listener monitor /tmp
  sudo app-listener monitor /home/alice --recursive --depth 3
  sudo app-listener monitor /etc/passwd /var/log -e OPEN,WRITE`,
	Args: cobra.MinimumNArgs(1),
	RunE: runMonitor,
}

func init() {
	common.AddServeFlags(MonitorCmd)
	f := MonitorCmd.Flags()
	f.BoolVarP(&recursive, "recursive", "r", false, "Include subdirectories")
	f.IntVarP(&depth, "depth", "d", 0, "With --recursive: maximum directory depth (0 = unlimited)")
	common.AddHeadlessFlag(MonitorCmd, &headless)
	f.StringSliceVarP(&eventsFlag, "events", "e", nil, "Event types to show: "+common.FileEventList+" (default all)")
}

func runMonitor(cmd *cobra.Command, args []string) error {
	serve, err := common.ParseServeFlags(cmd)
	if err != nil {
		return err
	}

	rawTargets, ebpfTargets, err := prepareMonitorTargets(args)
	if err != nil {
		return err
	}

	parsedEvents, err := common.ParseEventsFlag(eventsFlag)
	if err != nil {
		return err
	}

	printers.PrintLogo()

	if ebpfErr := common.CheckEBPF(); ebpfErr != nil {
		return ebpfErr
	}

	mon, monitorErr := monitor.NewMonitor(ebpfTargets, recursive, depth)
	if monitorErr != nil {
		return fmt.Errorf("creating monitor: %w", monitorErr)
	}
	mon.SetEventTypes(parsedEvents)

	ucase := usecase.NewMonitorUseCase(mon)

	if startErr := ucase.Start(); startErr != nil {
		return fmt.Errorf("starting monitor: %w", startErr)
	}

	defer ucase.Stop()

	displayPaths := common.MakeDisplayPaths(rawTargets)

	switch {
	case headless:
		runHeadless(ucase)
	case serve.Enabled:
		fan := tui.NewEventFanout(ucase.Events())
		defer fan.Stop()
		return tui.Serve(
			tui.NewModel(fan.Local(), displayPaths, recursive, depth),
			tui.NewModel(fan.Browser(), displayPaths, recursive, depth),
			tui.ServeOptions{
				Address: serve.Address, Username: serve.Username, Password: serve.Password,
			})
	default:
		return runTUI(ucase, displayPaths, recursive, depth)
	}

	return nil
}

// prepareMonitorTargets validates and resolves the paths into eBPF targets, warning about flags
// that single files silently ignore.
func prepareMonitorTargets(watchPaths []string) ([]common.RawTarget, []ebpf.Target, error) {
	if depth > 0 && !recursive {
		return nil, nil, errors.New("--depth requires --recursive")
	}

	rawTargets, err := common.ResolveTargets(watchPaths)
	if err != nil {
		return nil, nil, err
	}

	if err := common.ValidateTargets(rawTargets); err != nil {
		return nil, nil, err
	}

	common.WarnIgnoredFlags(rawTargets, recursive, depth)
	return rawTargets, common.BuildEBPFTargets(rawTargets), nil
}

func runHeadless(uc usecase.MonitorUseCase) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case ev, ok := <-uc.Events():
			if !ok {
				return
			}
			log.Infof("EVENT|%s|%s|%s|%d|%s|%d",
				ev.Type.String(), logging.SanitizeText(ev.Comm), logging.SanitizeText(ev.Path), ev.PID, logging.SanitizeText(ev.Dest), ev.FD)
		case <-sig:
			return
		}
	}
}

func runTUI(uc usecase.MonitorUseCase, paths []string, recursive bool, depth int) error {
	p := tea.NewProgram(
		tui.NewModel(uc.Events(), paths, recursive, depth),
		tea.WithAltScreen(),
	)

	_, runErr := p.Run()
	return runErr
}
