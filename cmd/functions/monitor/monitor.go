package monitor

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/cmd/entity"
	"github.com/Virgula0/app-listener/cmd/printers"
	"github.com/Virgula0/app-listener/internal/gui"
	"github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/tui"
)

var MonitorCmd = &cobra.Command{
	Use:   "monitor <directory>",
	Short: "Monitor a directory for file system events",
	Long: `Monitor a directory for file system events using eBPF.

Supports recursive monitoring with configurable depth.

Examples:
  app-listener monitor /tmp
  app-listener monitor /var/log --recursive
  app-listener monitor /home --recursive --depth 3`,
	Args: cobra.ExactArgs(1),
	RunE: runMonitor,
}

func init() {
	MonitorCmd.Flags().BoolVarP(&entity.Recursive, "recursive", "r", false,
		"Monitor directory recursively (default: false)")
	MonitorCmd.Flags().IntVarP(&entity.Depth, "depth", "d", 0,
		"Maximum directory depth (requires --recursive) (default: 0)")
}

func runMonitor(cmd *cobra.Command, args []string) error {
	directory := args[0]

	absDir, err := filepath.Abs(directory)
	if err != nil {
		log.Errorf("Invalid directory path: %v", err)
		return err
	}

	info, err := os.Stat(absDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Errorf("Directory does not exist: %s", absDir)
		} else {
			log.Errorf("Cannot access directory %s: %v", absDir, err)
		}
		return err
	}
	if !info.IsDir() {
		log.Errorf("Not a directory: %s", absDir)
		return err
	}

	if entity.Depth > 0 && !entity.Recursive {
		return errors.New("--depth requires --recursive flag")
	}

	printers.PrintLogo()

	log.Infof("Checking eBPF availability...")

	if checkErr := ebpf.Check(); checkErr != nil {
		log.Errorf("eBPF check failed: %v", checkErr)
		return checkErr
	}

	log.Infof("eBPF available — starting monitor")
	log.Debugf("Monitoring directory: %s (recursive: %v, depth: %d)",
		absDir, entity.Recursive, entity.Depth)

	mon, err := ebpf.NewMonitor(absDir, entity.Recursive, entity.Depth)
	if err != nil {
		log.Errorf("Failed to create monitor: %v", err)
		return err
	}

	if err := mon.Start(); err != nil {
		log.Errorf("Failed to start monitor: %v", err)
		return err
	}

	defer mon.Stop()

	if entity.GUI {
		gui.Run(mon.Events(), absDir, entity.Recursive, entity.Depth)
	} else {
		p := tea.NewProgram(
			tui.NewModel(mon.Events(), absDir, entity.Recursive, entity.Depth),
			tea.WithAltScreen(),
		)

		if _, err := p.Run(); err != nil {
			log.Errorf("TUI error: %v", err)
			return err
		}
	}

	return nil
}
