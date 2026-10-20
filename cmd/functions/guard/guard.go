package guard

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/cmd/common"
	"github.com/Virgula0/app-listener/cmd/printers"
	"github.com/Virgula0/app-listener/internal/guard"
	"github.com/Virgula0/app-listener/internal/logging"
	"github.com/Virgula0/app-listener/internal/tui"
	"github.com/Virgula0/app-listener/internal/usecase"
)

var (
	blacklistPaths []string
	whitelistPaths []string
	guardRecursive bool
	guardDepth     int
	headless       bool
)

var GuardCmd = &cobra.Command{
	Use:   "guard PATH [-w BINARY... | -b BINARY...]",
	Short: "Deny access to a file or directory with eBPF LSM",
	Long: `Deny access to PATH. With -w only the listed binaries may access it (with
neither -w nor -b, nothing may); with -b everything but the listed binaries may.
Binaries are identified by inode, so renaming one doesn't move it past the rule.`,
	Example: `  sudo app-listener guard /home/alice/.ssh -w /usr/bin/ssh -w /usr/bin/ssh-add
  sudo app-listener guard /secret -b /usr/bin/cat`,
	Args: cobra.ExactArgs(1),
	RunE: runGuard,
}

func init() {
	common.AddServeFlags(GuardCmd)
	f := GuardCmd.Flags()
	f.StringSliceVarP(&whitelistPaths, "whitelist", "w", nil, "Binary allowed to access PATH (repeatable)")
	f.StringSliceVarP(&blacklistPaths, "blacklist", "b", nil, "Binary denied access to PATH (repeatable)")
	common.AddHeadlessFlag(GuardCmd, &headless)
	f.BoolVarP(&guardRecursive, "recursive", "r", true, "Include subdirectories (--recursive=false: PATH's own entries only)")
	f.IntVarP(&guardDepth, "depth", "d", 0, "Maximum directory depth (0 = unlimited, 1 = PATH's own entries)")
	GuardCmd.MarkFlagsMutuallyExclusive("whitelist", "blacklist")
}

func runGuard(cmd *cobra.Command, args []string) error {
	serve, err := common.ParseServeFlags(cmd)
	if err != nil {
		return err
	}
	guardPath := args[0]

	mode, binaries, err := resolveGuardConfig()
	if err != nil {
		return err
	}

	printers.PrintLogo()

	if checkErr := common.CheckEBPF(); checkErr != nil {
		return checkErr
	}

	// The fs guard's enforcement is the LSM hooks (file_open, file_permission): refuse kernels
	// whose LSM stack lacks bpf (attach can succeed while the hooks never fire).
	if checkErr := common.CheckBPFLSM(); checkErr != nil {
		return checkErr
	}

	g, guardErr := guard.NewGuard(guardPath, mode, binaries, guardRecursive, guardDepth, guard.WithEagerPopulate())
	if guardErr != nil {
		return fmt.Errorf("creating guard: %w", guardErr)
	}

	ucase := usecase.NewGuardUseCase(g)

	// Every file/dir under the guarded path was already registered in guard_inodes inside NewGuard
	// (WithEagerPopulate) before the hooks attached: the tree is protected from the first live
	// instant and the guard's own scan can't be blocked by its file_open hook.
	if startErr := ucase.Start(); startErr != nil {
		return fmt.Errorf("starting guard: %w", startErr)
	}

	defer ucase.Stop()

	if headless {
		runGuardHeadless(ucase)
		return nil
	}
	if serve.Enabled {
		fan := tui.NewEventFanout(ucase.Events())
		defer fan.Stop()
		return tui.Serve(
			tui.NewGuardModel(fan.Local(), guardPath, mode, binaries),
			tui.NewGuardModel(fan.Browser(), guardPath, mode, binaries),
			tui.ServeOptions{
				Address: serve.Address, Username: serve.Username, Password: serve.Password,
			})
	}

	return runGuardTUI(ucase, guardPath, mode, binaries)
}

func resolveGuardConfig() (guard.Mode, []guard.BinaryEntry, error) {
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

func runGuardHeadless(uc usecase.GuardUseCase) {
	common.DrainUntilSignal(uc.Events(), func(ev *guard.GuardEvent) { log.Info(HeadlessLine(ev)) })
}

// HeadlessLine is ev in the stable GUARD| format --headless prints (the integration tests parse it).
func HeadlessLine(ev *guard.GuardEvent) string {
	return fmt.Sprintf("GUARD|%s|%s|%s|%d|%s|%t|%d",
		ev.Type.String(), logging.SanitizeText(ev.Comm), logging.SanitizeText(ev.Path), ev.PID,
		logging.SanitizeText(ev.Dest), ev.Blocked, ev.UID)
}

func runGuardTUI(uc usecase.GuardUseCase, guardPath string, mode guard.Mode, binaries []guard.BinaryEntry) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	p := tea.NewProgram(
		tui.NewGuardModel(uc.Events(), guardPath, mode, binaries),
		tea.WithAltScreen(),
	)

	_, runErr := p.Run()
	return runErr
}
