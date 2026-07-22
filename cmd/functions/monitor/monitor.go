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
	"github.com/Virgula0/app-listener/cmd/entity"
	"github.com/Virgula0/app-listener/cmd/printers"
	"github.com/Virgula0/app-listener/internal/gui"
	"github.com/Virgula0/app-listener/internal/monitor"
	"github.com/Virgula0/app-listener/internal/tui"
)

var watchPaths []string

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

var eventsFlag []string

func init() {
	MonitorCmd.Flags().StringSliceVarP(&watchPaths, "watch", "w", nil,
		"Path to monitor (repeatable, required)")
	MonitorCmd.Flags().BoolVarP(&entity.Recursive, "recursive", "r", false,
		"Monitor directory recursively (default: false)")
	MonitorCmd.Flags().IntVarP(&entity.Depth, "depth", "d", 0,
		"Maximum directory depth (requires --recursive) (default: 0)")
	MonitorCmd.Flags().BoolVarP(&entity.Headless, "headless", "", false,
		"Run without TUI, print events to stderr (for testing/scripting)")
	MonitorCmd.Flags().StringSliceVarP(&eventsFlag, "events", "e", nil,
		"Event types to monitor (comma-separated: OPEN,READ,WRITE,DELETE,RENAME,SYMLINK,HARDLINK,MKDIR,MMAP; default: all)")
}

func runMonitor(cmd *cobra.Command, args []string) error {
	if len(watchPaths) == 0 {
		return errors.New("at least one -w/--watch path is required")
	}

	recursive := entity.Recursive
	depth := entity.Depth

	if depth > 0 && !recursive {
		return errors.New("--depth requires --recursive flag")
	}

	rawTargets, targetErr := common.ResolveTargets(watchPaths)
	if targetErr != nil {
		return targetErr
	}

	if err := common.ValidateTargets(rawTargets); err != nil {
		return err
	}

	ebpfTargets := common.BuildEBPFTargets(rawTargets)
	common.WarnIgnoredFlags(rawTargets, recursive, depth)

	parsedEvents, eventsErr := common.ParseEventsFlag(eventsFlag)
	if eventsErr != nil {
		return eventsErr
	}
	entity.EventTypes = parsedEvents

	printers.PrintLogo()

	if ebpfErr := common.CheckEBPF(); ebpfErr != nil {
		return ebpfErr
	}

	mon, monitorErr := monitor.NewMonitor(ebpfTargets, recursive, depth)
	if monitorErr != nil {
		return fmt.Errorf("creating monitor: %w", monitorErr)
	}

	mon.SetEventTypes(entity.EventTypes)

	if startErr := mon.Start(); startErr != nil {
		return fmt.Errorf("starting monitor: %w", startErr)
	}

	defer mon.Stop()

	displayPaths := common.MakeDisplayPaths(rawTargets)

	switch {
	case entity.Headless:
		runHeadless(mon)
	case entity.GUI:
		gui.Run(mon.Events(), displayPaths, recursive, depth)
	default:
		return runTUI(mon, displayPaths, recursive, depth)
	}

	return nil
}

func runHeadless(mon *monitor.Monitor) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case ev, ok := <-mon.Events():
			if !ok {
				return
			}
			log.Infof("EVENT|%s|%s|%s|%d|%s|%d",
				ev.Type.String(), ev.Comm, ev.Path, ev.PID, ev.Dest, ev.FD)
		case <-sig:
			return
		}
	}
}

func runTUI(mon *monitor.Monitor, paths []string, recursive bool, depth int) error {
	p := tea.NewProgram(
		tui.NewModel(mon.Events(), paths, recursive, depth),
		tea.WithAltScreen(),
	)

	_, runErr := p.Run()
	return runErr
}
