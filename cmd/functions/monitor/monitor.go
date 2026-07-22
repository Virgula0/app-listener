package monitor

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/cmd/entity"
	"github.com/Virgula0/app-listener/cmd/printers"
	"github.com/Virgula0/app-listener/internal/gui"
	"github.com/Virgula0/app-listener/internal/infrastructure"
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

func init() {
	MonitorCmd.Flags().StringSliceVarP(&watchPaths, "watch", "w", nil,
		"Path to monitor (repeatable, required)")
	MonitorCmd.Flags().BoolVarP(&entity.Recursive, "recursive", "r", false,
		"Monitor directory recursively (default: false)")
	MonitorCmd.Flags().IntVarP(&entity.Depth, "depth", "d", 0,
		"Maximum directory depth (requires --recursive) (default: 0)")
	MonitorCmd.Flags().BoolVarP(&entity.Headless, "headless", "", false,
		"Run without TUI, print events to stderr (for testing/scripting)")
}

type rawTarget struct {
	absPath string
	isDir   bool
}

func runMonitor(cmd *cobra.Command, args []string) error {
	if len(watchPaths) == 0 {
		return errors.New("at least one -w/--watch path is required")
	}

	rawTargets, err := resolveTargets(watchPaths)
	if err != nil {
		return err
	}

	if err := validateTargets(rawTargets); err != nil {
		return err
	}

	recursive := entity.Recursive
	depth := entity.Depth

	if depth > 0 && !recursive {
		return errors.New("--depth requires --recursive flag")
	}

	ebpfTargets := make([]ebpf.Target, len(rawTargets))
	for i, rt := range rawTargets {
		t := ebpf.Target{Path: rt.absPath, IsDir: rt.isDir}
		if rt.isDir {
			t.Dir = rt.absPath
		} else {
			t.Dir = filepath.Dir(rt.absPath)
			t.File = filepath.Base(rt.absPath)
		}
		ebpfTargets[i] = t
	}

	for _, rt := range rawTargets {
		if !rt.isDir && recursive {
			log.Warnf("--recursive is ignored when monitoring a single file (%s)", rt.absPath)
		}
		if !rt.isDir && depth > 0 {
			log.Warnf("--depth is ignored when monitoring a single file (%s)", rt.absPath)
		}
	}

	printers.PrintLogo()

	log.Infof("Checking eBPF availability...")

	if checkErr := ebpf.Check(); checkErr != nil {
		log.Errorf("eBPF check failed: %v", checkErr)
		return checkErr
	}

	log.Infof("eBPF available — starting monitor")

	mon, err := ebpf.NewMonitor(ebpfTargets, recursive, depth)
	if err != nil {
		log.Errorf("Failed to create monitor: %v", err)
		return err
	}

	if err := mon.Start(); err != nil {
		log.Errorf("Failed to start monitor: %v", err)
		return err
	}

	defer mon.Stop()

	displayPaths := make([]string, len(rawTargets))
	for i, rt := range rawTargets {
		displayPaths[i] = rt.absPath
	}

	if entity.Headless {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		for {
			select {
			case ev, ok := <-mon.Events():
				if !ok {
					return nil
				}
				msg := fmt.Sprintf("EVENT|%s|%s|%s|%d|%s|%d",
					ev.Type.String(), ev.Comm, ev.Path, ev.PID, ev.Dest, ev.FD)
				log.Infof(msg)
			case <-sig:
				return nil
			}
		}
	} else if entity.GUI {
		gui.Run(mon.Events(), displayPaths, recursive, depth)
	} else {
		p := tea.NewProgram(
			tui.NewModel(mon.Events(), displayPaths, recursive, depth),
			tea.WithAltScreen(),
		)

		if _, err := p.Run(); err != nil {
			log.Errorf("TUI error: %v", err)
			return err
		}
	}

	return nil
}

func resolveTargets(paths []string) ([]rawTarget, error) {
	seen := make(map[string]bool)
	var targets []rawTarget

	for _, p := range paths {
		abs, err := filepath.Abs(p)
		if err != nil {
			log.Errorf("Invalid path %q: %v", p, err)
			return nil, err
		}

		if seen[abs] {
			return nil, fmt.Errorf("duplicate watch path: %s", abs)
		}
		seen[abs] = true

		info, err := os.Stat(abs)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				log.Errorf("Path does not exist: %s", abs)
			} else {
				log.Errorf("Cannot access path %s: %v", abs, err)
			}
			return nil, err
		}

		targets = append(targets, rawTarget{absPath: abs, isDir: info.IsDir()})
	}

	return targets, nil
}

func validateTargets(targets []rawTarget) error {
	for i, a := range targets {
		for j, b := range targets {
			if i == j {
				continue
			}

			if a.isDir && b.isDir {
				if isSubDir(a.absPath, b.absPath) {
					return fmt.Errorf(
						"%q is a subdirectory of %q — remove the more specific path",
						b.absPath, a.absPath)
				}
				if isSubDir(b.absPath, a.absPath) {
					return fmt.Errorf(
						"%q is a subdirectory of %q — remove the more specific path",
						a.absPath, b.absPath)
				}
			}

			if !a.isDir && b.isDir && pathWithinDir(a.absPath, b.absPath) {
				return fmt.Errorf(
					"%q is already covered by directory %q — remove the redundant file path",
					a.absPath, b.absPath)
			}
			if a.isDir && !b.isDir && pathWithinDir(b.absPath, a.absPath) {
				return fmt.Errorf(
					"%q is already covered by directory %q — remove the redundant file path",
					b.absPath, a.absPath)
			}
		}
	}
	return nil
}

func isSubDir(child, parent string) bool {
	return strings.HasPrefix(child+"/", parent+"/")
}

func pathWithinDir(path, dir string) bool {
	return strings.HasPrefix(path+"/", dir+"/")
}
