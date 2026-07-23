package guard

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
	"github.com/Virgula0/app-listener/internal/guard"
	"github.com/Virgula0/app-listener/internal/tui"
)

var (
	blacklistPaths []string
	whitelistPaths []string
	eventsFlag     []string
)

var GuardCmd = &cobra.Command{
	Use:   "guard <path>",
	Short: "Block file system access using eBPF LSM",
	Long: `Guard a file or directory using eBPF LSM programs to block access
by unauthorized processes.

Use -b/--blacklist or -w/--whitelist (mutually exclusive) to specify
process binaries to block or allow. Each binary is identified by its
SHA256 hash.

Examples:
  app-listener guard /secret -b /usr/bin/cat
  app-listener guard /etc/shadow -w /usr/bin/vim
  app-listener guard /data -w /usr/bin/cat -e OPEN,READ`,
	Args: cobra.ExactArgs(1),
	RunE: runGuard,
}

func init() {
	GuardCmd.Flags().StringSliceVarP(&blacklistPaths, "blacklist", "b", nil,
		"Binary paths to blacklist (repeatable, mutually exclusive with -w)")
	GuardCmd.Flags().StringSliceVarP(&whitelistPaths, "whitelist", "w", nil,
		"Binary paths to whitelist (repeatable, mutually exclusive with -b)")
	GuardCmd.Flags().StringSliceVarP(&eventsFlag, "events", "e", nil,
		"Event types to monitor (comma-separated: OPEN,READ,WRITE,DELETE,RENAME,SYMLINK,HARDLINK,MKDIR,MMAP; default: all)")
	GuardCmd.Flags().BoolVarP(&entity.Headless, "headless", "", false,
		"Run without TUI, print events to stderr (for testing/scripting)")
}

func runGuard(cmd *cobra.Command, args []string) error {
	guardPath := args[0]

	mode, binaries, err := resolveGuardConfig()
	if err != nil {
		return err
	}

	parsedEvents, err := common.ParseEventsFlag(eventsFlag)
	if err != nil {
		return err
	}
	entity.EventTypes = parsedEvents

	printers.PrintLogo()

	if checkErr := common.CheckEBPF(); checkErr != nil {
		return checkErr
	}

	g, guardErr := guard.NewGuard(guardPath, mode, binaries)
	if guardErr != nil {
		return fmt.Errorf("creating guard: %w", guardErr)
	}

	if startErr := g.Start(); startErr != nil {
		return fmt.Errorf("starting guard: %w", startErr)
	}

	defer g.Stop()

	if entity.Headless {
		runGuardHeadless(g)
		return nil
	}

	return runGuardTUI(g, guardPath, mode, binaries)
}

func resolveGuardConfig() (guard.Mode, []guard.BinaryEntry, error) {
	if len(blacklistPaths) > 0 && len(whitelistPaths) > 0 {
		return 0, nil, errors.New("--blacklist and --whitelist are mutually exclusive")
	}

	var mode guard.Mode
	var paths []string

	switch {
	case len(blacklistPaths) > 0:
		mode = guard.ModeBlacklist
		paths = blacklistPaths
	case len(whitelistPaths) > 0:
		mode = guard.ModeWhitelist
		paths = whitelistPaths
	default:
		mode = guard.ModeWhitelist
		paths = nil
	}

	binaries := make([]guard.BinaryEntry, 0, len(paths))
	for _, p := range paths {
		entry, err := guard.ComputeBinaryEntry(p)
		if err != nil {
			return 0, nil, fmt.Errorf("processing binary %q: %w", p, err)
		}
		binaries = append(binaries, entry)
	}

	return mode, binaries, nil
}

func runGuardHeadless(g *guard.Guard) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case ev, ok := <-g.Events():
			if !ok {
				return
			}
			log.Infof("GUARD|%s|%s|%s|%d|%s|%t|%d",
				ev.Type.String(), ev.Comm, ev.Path, ev.PID, ev.Dest, ev.Blocked, ev.UID)
		case <-sig:
			return
		}
	}
}

func runGuardTUI(g *guard.Guard, guardPath string, mode guard.Mode, binaries []guard.BinaryEntry) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	p := tea.NewProgram(
		tui.NewGuardModel(g.Events(), guardPath, mode, binaries),
		tea.WithAltScreen(),
	)

	_, runErr := p.Run()
	return runErr
}
