package editprotected

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/x/term"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/cmd/common"
	guardcmd "github.com/Virgula0/app-listener/cmd/functions/guard"
	"github.com/Virgula0/app-listener/internal/guard"
	"github.com/Virgula0/app-listener/internal/tui"
)

// validateForwardFlags enforces the --forward flag surface: -w/-b/-e/--yes only with it, -w xor -b
// required with it, and no editor or password flags alongside.
func validateForwardFlags() error {
	if !forwardFlag {
		if len(whitelistFlags)+len(blacklistFlags)+len(eventsFlags) > 0 || yesFlag {
			return errors.New("-w, -b, -e and --yes work only with --forward")
		}
		return nil
	}
	if err := validateForwardRule(); err != nil {
		return err
	}
	if putFlag != "" || contentFileFlag != "" || setPasswordFlag || clearPasswordFlag {
		return errors.New("--forward can't be combined with --put, --content-file, --set-password or --clear-password")
	}
	return nil
}

func validateForwardRule() error {
	switch {
	case len(whitelistFlags) > 0 && len(blacklistFlags) > 0:
		return errors.New("--whitelist and --blacklist are mutually exclusive")
	case len(whitelistFlags) == 0 && len(blacklistFlags) == 0:
		return errors.New("--forward needs binaries: -w <binary> to admit or -b <binary> to deny")
	}
	_, err := common.ParseEventsFlag(eventsFlags)
	return err
}

// forwardBinaries is the -w or -b list made absolute (the daemon resolves them in its own cwd).
func forwardBinaries() ([]string, error) {
	paths := whitelistFlags
	if len(blacklistFlags) > 0 {
		paths = blacklistFlags
	}
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		abs, err := filepath.Abs(p)
		if err != nil {
			return nil, fmt.Errorf("resolving %s: %w", p, err)
		}
		out = append(out, abs)
	}
	return out, nil
}

// runForward is the --forward flow: confirm risky binaries, authenticate (before any protected
// path is shown), pick resources, activate the grant, and hold it until the operator or the daemon
// ends it.
func runForward() error {
	allow := len(whitelistFlags) > 0
	bins, err := forwardBinaries()
	if err != nil {
		return err
	}
	if allow {
		if cerr := confirmUserPlaced(bins); cerr != nil {
			return cerr
		}
	}

	password := os.Getenv(editPasswordEnv)
	if password == "" {
		if password, err = promptPassword("Edit-protected password"); err != nil {
			return err
		}
	}
	session, err := dialLiveSession()
	if err != nil {
		return err
	}
	defer session.End()
	session.timeout = timeoutSession

	resources, err := session.Authenticate(password)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	chosen, err := chooseResources(resources)
	if err != nil {
		return err
	}

	notes, err := session.Forward(chosen, bins, allow, normalizedEvents())
	if err != nil {
		return fmt.Errorf("activating temporary access: %w", err)
	}
	for _, n := range notes {
		log.Warn(n)
	}
	verb := "DENIED"
	if allow {
		verb = "ALLOWED"
	}
	log.Warnf("temporary access active: %s %s on %s (events %s) — revoked at exit, or after %s "+
		"without activity", strings.Join(bins, ", "), verb, strings.Join(chosen, ", "), eventsLabel(),
		timeoutSession)
	return holdForward(session, chosen, bins, allow)
}

func normalizedEvents() []string {
	out := make([]string, 0, len(eventsFlags))
	for _, e := range eventsFlags {
		out = append(out, strings.ToUpper(strings.TrimSpace(e)))
	}
	return out
}

func eventsLabel() string {
	if len(eventsFlags) == 0 {
		return "ALL"
	}
	return strings.Join(normalizedEvents(), ",")
}

// confirmUserPlaced asks before admitting a binary that isn't root-placed: its owner can rewrite it
// in place during the window (the daemon revokes the grant when it notices). --yes skips the
// question; without a terminal it is required.
func confirmUserPlaced(bins []string) error {
	var risky []string
	for _, p := range bins {
		b, err := guard.ResolveTempBinary(p)
		if err != nil {
			return fmt.Errorf("%s: %w", p, err)
		}
		if !b.SystemPlaced {
			risky = append(risky, p)
		}
		if names := b.RuntimeNames(); len(names) > 0 {
			log.Infof("%s is a %s runtime: started with code-loading env or flags (NODE_OPTIONS, "+
				"ELECTRON_RUN_AS_NODE, --inspect, JAVA_TOOL_OPTIONS, …) it is refused the directories", p,
				strings.Join(names, "/"))
		}
		b.Close()
	}
	if len(risky) == 0 || yesFlag {
		return nil
	}
	msg := "Not root-placed (a user can rewrite them while the grant is active; it is revoked if their " +
		"content changes):\n  " + strings.Join(risky, "\n  ") +
		"\nThe grant covers every user running them, and an interpreter (python, bash, node) admits " +
		"every script it is given."
	if !term.IsTerminal(os.Stdin.Fd()) {
		return errors.New(msg + "\nRe-run with --yes to admit them anyway.")
	}
	ok := false
	if err := huh.NewForm(huh.NewGroup(
		huh.NewConfirm().Title("Admit binaries that are not root-placed?").Description(msg).Value(&ok),
	)).Run(); err != nil {
		return err
	}
	if !ok {
		return errors.New("aborted")
	}
	return nil
}

// chooseResources honors --resource (each must be in the daemon's list), otherwise asks.
func chooseResources(resources []string) ([]string, error) {
	if len(resourceFlags) > 0 {
		for _, r := range resourceFlags {
			if !slices.Contains(resources, r) {
				return nil, fmt.Errorf("%s is not a guarded directory in the running daemon", r)
			}
		}
		return resourceFlags, nil
	}
	if len(resources) == 1 {
		log.Infof("only one guarded directory: %s", resources[0])
		return resources, nil
	}
	opts := make([]huh.Option[string], 0, len(resources))
	for _, r := range resources {
		opts = append(opts, huh.NewOption(r, r))
	}
	var chosen []string
	if err := huh.NewForm(huh.NewGroup(
		huh.NewMultiSelect[string]().
			Title("Which guarded directories should the binaries reach?").
			Description("Space selects; the access is revoked when this command exits.").
			Options(opts...).
			Validate(func(s []string) error {
				if len(s) == 0 {
					return errors.New("select at least one directory")
				}
				return nil
			}).
			Value(&chosen),
	)).Run(); err != nil {
		return nil, err
	}
	return chosen, nil
}

// errForwardEnded: the daemon revoked the grant on its own.
var errForwardEnded = errors.New("the daemon ended the temporary access (idle timeout, reload, or a granted " +
	"binary changed) — see its log")

// holdForward shows the granted binaries' events until the operator quits (the guard TUI's q /
// ctrl+c, or a signal without a terminal) or the daemon ends the grant.
func holdForward(session *liveSession, resources, bins []string, allow bool) error {
	events, ended := session.Watch()
	if term.IsTerminal(os.Stdin.Fd()) {
		mode := guard.ModeBlacklist
		if allow {
			mode = guard.ModeWhitelist
		}
		model := tui.NewGuardModel(events, strings.Join(resources, ", "), mode, displayEntries(bins))
		byDaemon, err := tui.RunSession(model, session.Pinger(pingInterval()), ended)
		if err != nil {
			return err
		}
		if byDaemon {
			return errForwardEnded
		}
	} else if holdHeadless(events, ended) {
		return errForwardEnded
	}
	session.End()
	log.Info("temporary access revoked")
	return nil
}

// holdHeadless prints events in guard's GUARD| format until a signal (false) or the daemon ending
// the grant (true).
func holdHeadless(events <-chan guard.GuardEvent, ended <-chan struct{}) bool {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer signal.Stop(sig)
	for {
		select {
		case ev := <-events:
			log.Info(guardcmd.HeadlessLine(&ev))
		case <-ended:
			return true
		case <-sig:
			return false
		}
	}
}

// displayEntries names bins for the guard TUI header; identity stays the daemon's inode check.
func displayEntries(bins []string) []guard.BinaryEntry {
	out := make([]guard.BinaryEntry, 0, len(bins))
	for _, p := range bins {
		e, err := guard.ComputeBinaryEntry(p)
		if err != nil {
			e = guard.BinaryEntry{Path: p, Comm: filepath.Base(p)}
		}
		out = append(out, e)
	}
	return out
}

// pingInterval throttles activity PINGs well inside the idle timeout.
func pingInterval() time.Duration {
	return min(5*time.Second, timeoutSession/4)
}
