package monitor

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/cmd/entity"
	"github.com/Virgula0/app-listener/cmd/printers"
	"github.com/Virgula0/app-listener/internal/infrastructure/ebpf"
	"github.com/Virgula0/app-listener/internal/tui"
)

var MonitorCmd = &cobra.Command{
	Use:   "monitor <directory>",
	Short: "Monitor a directory for file system events",
	Long:  `Monitor a directory for file system events using eBPF. Supports recursive monitoring.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runMonitor,
}

func init() {
	MonitorCmd.Flags().BoolVarP(&entity.Recursive, "recursive", "r", false, "Monitor directory recursively")
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
		if os.IsNotExist(err) {
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

	printers.PrintLogo()

	log.Infof("Checking eBPF availability...")

	if err := ebpf.Check(); err != nil {
		log.Errorf("eBPF check failed: %v", err)
		return err
	}

	log.Infof("eBPF available — starting TUI monitor")
	log.Debugf("Monitoring directory: %s (recursive: %v)", absDir, entity.Recursive)

	app := tui.NewModel(absDir, entity.Recursive)

	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Errorf("TUI error: %v", err)
		return err
	}

	return nil
}
