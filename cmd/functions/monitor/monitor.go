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
	"github.com/Virgula0/app-listener/internal/gui"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/logging"
	"github.com/Virgula0/app-listener/internal/monitor"
	"github.com/Virgula0/app-listener/internal/tui"
	"github.com/Virgula0/app-listener/internal/usecase"
)

var (
	watchPaths []string
	eventsFlag []string

	recursive bool
	depth     int
	headless  bool
)

var MonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor files/directories for file system events",
	Long: `Monitor one or more files or directories for file system events using eBPF.

Use -w/--watch (repeatable) to specify paths to monitor.
Accepts directories (to monitor all files within) or single files.
When monitoring a file, only events for that specific file are shown.

Examples:
  app-listener monitor -w /tmp
  app-listener monitor -w /var/log --recursive
  app-listener monitor -w /home --recursive --depth 3
  app-listener monitor -w /path/to/file.txt
  app-listener monitor -w /dir1 -w /dir2 -w /path/to/file`,
	Args: cobra.NoArgs,
	RunE: runMonitor,
}

func init() {
	common.AddServeFlags(MonitorCmd)
	MonitorCmd.Flags().StringSliceVarP(&watchPaths, "watch", "w", nil,
		"Path to monitor (repeatable, required)")
	MonitorCmd.Flags().BoolVarP(&recursive, "recursive", "r", false,
		"Monitor directory recursively (default: false)")
	MonitorCmd.Flags().IntVarP(&depth, "depth", "d", 0,
		"Maximum directory depth (requires --recursive) (default: 0)")
	MonitorCmd.Flags().BoolVarP(&headless, "headless", "", false,
		"Run without TUI, print events to stderr (for testing/scripting)")
	MonitorCmd.Flags().StringSliceVarP(&eventsFlag, "events", "e", nil,
		"Event types to monitor (comma-separated: OPEN,READ,WRITE,DELETE,RENAME,SYMLINK,HARDLINK,MKDIR,MMAP,ATTR,STAT,MKNOD; default: all)")
}

func runMonitor(cmd *cobra.Command, args []string) error {
	serve, err := common.ParseServeFlags(cmd)
	if err != nil {
		return err
	}

	rawTargets, ebpfTargets, err := prepareMonitorTargets()
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

	// --gui is a root persistent flag, inherited by this command.
	launchGUI, _ := cmd.Flags().GetBool("gui")

	switch {
	case headless:
		runHeadless(ucase)
	case launchGUI:
		gui.Run(ucase.Events(), displayPaths, recursive, depth)
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

// prepareMonitorTargets validates and resolves the --watch paths into eBPF
// targets, warning about flags that single files silently ignore.
func prepareMonitorTargets() ([]common.RawTarget, []ebpf.Target, error) {
	if len(watchPaths) == 0 {
		return nil, nil, errors.New("at least one -w/--watch path is required")
	}

	if depth > 0 && !recursive {
		return nil, nil, errors.New("--depth requires --recursive flag")
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
