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

type rawTarget struct {
	absPath string
	isDir   bool
}

func runMonitor(cmd *cobra.Command, args []string) error {
	rawTargets, targetErr := prepareTargets()
	if targetErr != nil {
		return targetErr
	}

	recursive := entity.Recursive
	depth := entity.Depth

	if depth > 0 && !recursive {
		return errors.New("--depth requires --recursive flag")
	}

	ebpfTargets := buildEBPFTargets(rawTargets)
	warnIgnoredFlags(rawTargets, recursive, depth)

	if eventsErr := parseEventsFlag(); eventsErr != nil {
		return eventsErr
	}

	printers.PrintLogo()

	if ebpfErr := checkEBPF(); ebpfErr != nil {
		return ebpfErr
	}

	mon, monitorErr := ebpf.NewMonitor(ebpfTargets, recursive, depth)
	if monitorErr != nil {
		return fmt.Errorf("creating monitor: %w", monitorErr)
	}

	mon.SetEventTypes(entity.EventTypes)

	if startErr := mon.Start(); startErr != nil {
		return fmt.Errorf("starting monitor: %w", startErr)
	}

	defer mon.Stop()

	displayPaths := makeDisplayPaths(rawTargets)

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

func prepareTargets() ([]rawTarget, error) {
	if len(watchPaths) == 0 {
		return nil, errors.New("at least one -w/--watch path is required")
	}

	rawTargets, err := resolveTargets(watchPaths)
	if err != nil {
		return nil, err
	}
	if err := validateTargets(rawTargets); err != nil {
		return nil, err
	}
	return rawTargets, nil
}

func checkEBPF() error {
	log.Info("Checking eBPF availability...")
	if err := ebpf.Check(); err != nil {
		log.Errorf("eBPF check failed: %v", err)
		return err
	}
	log.Info("eBPF available — starting monitor")
	return nil
}

func makeDisplayPaths(rawTargets []rawTarget) []string {
	paths := make([]string, len(rawTargets))
	for i, rt := range rawTargets {
		paths[i] = rt.absPath
	}
	return paths
}

func runHeadless(mon *ebpf.Monitor) {
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

func runTUI(mon *ebpf.Monitor, paths []string, recursive bool, depth int) error {
	p := tea.NewProgram(
		tui.NewModel(mon.Events(), paths, recursive, depth),
		tea.WithAltScreen(),
	)

	_, runErr := p.Run()
	return runErr
}

func buildEBPFTargets(rawTargets []rawTarget) []ebpf.Target {
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
	return ebpfTargets
}

func warnIgnoredFlags(rawTargets []rawTarget, recursive bool, depth int) {
	for _, rt := range rawTargets {
		if !rt.isDir && recursive {
			log.Warnf("--recursive is ignored when monitoring a single file (%s)", rt.absPath)
		}
		if !rt.isDir && depth > 0 {
			log.Warnf("--depth is ignored when monitoring a single file (%s)", rt.absPath)
		}
	}
}

func parseEventsFlag() error {
	if len(eventsFlag) == 0 {
		entity.EventTypes = ebpf.EventTypes()
		return nil
	}

	var parsed []ebpf.EventType
	for _, s := range eventsFlag {
		et, ok := ebpf.ParseEventType(strings.TrimSpace(s))
		if !ok {
			return fmt.Errorf("unknown event type %q (valid: OPEN, READ, WRITE, DELETE, RENAME, SYMLINK, HARDLINK, MKDIR, MMAP)", s)
		}
		parsed = append(parsed, et)
	}
	entity.EventTypes = parsed
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
			if err := validatePair(a, b); err != nil {
				return err
			}
		}
	}
	return nil
}

func validatePair(a, b rawTarget) error {
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
	return nil
}

func isSubDir(child, parent string) bool {
	return strings.HasPrefix(child+"/", parent+"/")
}

func pathWithinDir(path, dir string) bool {
	return strings.HasPrefix(path+"/", dir+"/")
}
