// Package trust implements `app-listener trust`: review and confirm the whitelisted binaries the
// daemon refused because their content is not the one recorded for their whitelist line (issue #80).
package trust

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/Virgula0/app-listener/internal/binledger"
	"github.com/Virgula0/app-listener/internal/bintrust"
	"github.com/Virgula0/app-listener/internal/logging"
	"github.com/Virgula0/app-listener/internal/systemd"
	"github.com/Virgula0/app-listener/internal/wizard"
)

var (
	listOnly bool
	yesAll   bool
	noReload bool
)

// TrustCmd reviews the binaries the daemon refused and records the confirmed ones.
var TrustCmd = &cobra.Command{
	Use:   "trust",
	Short: "Confirm whitelisted binaries the daemon refused (changed or new while it was stopped)",
	Long: `The daemon records the content hash of every whitelisted binary it admits. One
whose content differs, or a new whitelist line, is refused unless only root
could have placed it or the app's own updater wrote it while the daemon ran.

This lists the refused binaries and asks to confirm each; a running daemon is
then reloaded so the confirmed ones are admitted.`,
	Example: `  sudo app-listener trust
  sudo app-listener trust list
  sudo app-listener trust add /home/alice/.local/bin/tool`,
	Args: cobra.NoArgs,
	RunE: run,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the refused binaries without confirming any",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		listOnly = true
		return run(cmd, args)
	},
}

var addCmd = &cobra.Command{
	Use:   "add PATH...",
	Short: "Confirm whitelist lines as they are now (e.g. one you added to daemon.conf by hand)",
	Args:  cobra.MinimumNArgs(1),
	RunE:  run,
}

func init() {
	for _, c := range []*cobra.Command{TrustCmd, addCmd} {
		c.Flags().BoolVarP(&yesAll, "yes", "y", false, "Confirm every binary without asking")
		c.Flags().BoolVar(&noReload, "no-reload", false,
			"Don't reload a running daemon (confirmed binaries are admitted at its next reload)")
	}
	TrustCmd.AddCommand(listCmd, addCmd)
}

func run(cmd *cobra.Command, args []string) error {
	if os.Geteuid() != 0 {
		return errors.New("trust must be run as root: sudo app-listener trust")
	}
	l, err := bintrust.OpenLedger()
	if err != nil {
		return err
	}
	defer l.Close()

	cands, err := candidates(l, args)
	if err != nil {
		return err
	}
	out := cmd.OutOrStdout()
	if len(cands) == 0 {
		fmt.Fprintln(out, "No refused binaries: nothing to confirm.")
		return nil
	}
	confirmed, err := confirmAll(out, l, cands, len(args) > 0)
	if err != nil {
		return err
	}
	if confirmed == 0 || noReload {
		return nil
	}
	if err := systemd.ReloadDaemonIfActive(); err != nil {
		return fmt.Errorf("%d binary(ies) confirmed, but reloading the daemon failed (%w): run `systemctl reload %s`",
			confirmed, err, systemd.DaemonServiceName)
	}
	fmt.Fprintf(out, "%d binary(ies) confirmed; a running daemon was reloaded to admit them.\n", confirmed)
	return nil
}

// confirmAll describes each candidate and records the ones confirmed (all with --yes, none with
// `trust list`), returning how many.
func confirmAll(out io.Writer, l *binledger.Ledger, cands []bintrust.Row, given bool) (int, error) {
	confirmed := 0
	for i := range cands {
		c := &cands[i]
		describe(out, c)
		if listOnly {
			continue
		}
		if why := c.Skip(); why != "" {
			fmt.Fprintf(out, "  %s: skipped\n", why)
			continue
		}
		ok := yesAll
		if !ok {
			var err error
			if ok, err = wizard.ConfirmOnce(fmt.Sprintf("Trust %s for %s?", c.Line, c.Resources), "Trust"); err != nil {
				return confirmed, err
			}
		}
		if !ok {
			continue
		}
		if err := record(l, c, given); err != nil {
			return confirmed, fmt.Errorf("recording %s: %w", c.Line, err)
		}
		confirmed++
		fmt.Fprintf(out, "  confirmed.\n")
	}
	return confirmed, nil
}

// record stores the hash shown; a refused row only while the daemon's refusal still names it.
func record(l *binledger.Ledger, c *bintrust.Row, given bool) error {
	if given {
		return l.Record(c.Line, c.Current, binledger.SourceConfirmed)
	}
	ok, err := l.ConfirmPending(c.Line, c.Current)
	if err == nil && !ok {
		err = errors.New("the daemon refused other content for it meanwhile: run trust again")
	}
	return err
}

// candidates are the ledger's pending rows, or the given lines as they are now.
func candidates(l *binledger.Ledger, args []string) ([]bintrust.Row, error) {
	if len(args) > 0 {
		return bintrust.Given(args, "(given on the command line)"), nil
	}
	return bintrust.Pending(l)
}

func describe(w io.Writer, c *bintrust.Row) {
	fmt.Fprintf(w, "\n%s\n", logging.SanitizeText(c.Line))
	if c.Resolved != c.Line {
		fmt.Fprintf(w, "  resolves to: %s\n", logging.SanitizeText(c.Resolved))
	}
	fmt.Fprintf(w, "  resource:    %s\n", logging.SanitizeText(c.Resources))
	fmt.Fprintf(w, "  refused at:  %s\n", c.SeenAt.Format(time.RFC3339))
	fmt.Fprintf(w, "  sha256:      %x\n", c.SHA256)
	if st := c.Stat; st != nil {
		fmt.Fprintf(w, "  owner uid:   %d   mode: %o   modified: %s\n", st.Uid, st.Mode&0o7777,
			time.Unix(st.Mtim.Sec, 0).Format(time.RFC3339))
	}
}
