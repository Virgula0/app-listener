// Package trustbinaries implements `app-listener trust-binaries`: review and confirm the
// whitelisted binaries the daemon refused because their content is not the one recorded for their
// whitelist line (issue #80).
package trustbinaries

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/binledger"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/logging"
	"github.com/Virgula0/app-listener/internal/systemd"
	"github.com/Virgula0/app-listener/internal/wizard"
)

var (
	listOnly bool
	yesAll   bool
	noReload bool
)

// TrustBinariesCmd lists the binaries the daemon refused and records the confirmed ones.
var TrustBinariesCmd = &cobra.Command{
	Use:   "trust-binaries [path...]",
	Short: "Confirm whitelisted binaries the daemon refused (changed or new while it was not running)",
	Long: `The daemon records the content hash of every whitelisted binary it admits. A binary whose
content differs from the recorded one, or a new whitelist line, is refused unless only root could
have placed it or the app's own updater created it while the daemon ran: nothing guards those
paths while the daemon is stopped, booting or being reinstalled.

Without arguments, lists the refused binaries and asks to confirm each one. With paths, confirms
those whitelist lines as they are now (e.g. a line you added to daemon.conf by hand). A running
daemon is reloaded afterwards so confirmed binaries are admitted.`,
	RunE: run,
}

func init() {
	TrustBinariesCmd.Flags().BoolVar(&listOnly, "list", false, "Only list the refused binaries")
	TrustBinariesCmd.Flags().BoolVar(&yesAll, "yes", false, "Confirm every listed binary without asking")
	TrustBinariesCmd.Flags().BoolVar(&noReload, "no-reload", false,
		"Don't reload a running daemon afterwards (confirmed binaries are admitted at its next reload)")
}

type candidate struct {
	binledger.Pending
	current [32]byte
	readErr error
}

func run(cmd *cobra.Command, args []string) error {
	if os.Geteuid() != 0 {
		return errors.New("trust-binaries must run as root")
	}
	if _, err := os.Stat(filepath.Dir(binledger.DefaultPath)); err != nil {
		return fmt.Errorf("no %s: is the daemon installed? (%w)", filepath.Dir(binledger.DefaultPath), err)
	}
	l, err := binledger.Open(binledger.DefaultPath)
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
	confirmed, err := confirmAll(out, l, cands)
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
// --list), returning how many.
func confirmAll(out io.Writer, l *binledger.Ledger, cands []candidate) (int, error) {
	confirmed := 0
	for i := range cands {
		c := &cands[i]
		describe(out, c)
		if listOnly || !confirmable(out, c) {
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
		if err := l.Record(c.Line, c.current, binledger.SourceConfirmed); err != nil {
			return confirmed, fmt.Errorf("recording %s: %w", c.Line, err)
		}
		confirmed++
		fmt.Fprintf(out, "  confirmed.\n")
	}
	return confirmed, nil
}

// candidates are the ledger's pending rows, or the given lines as they are now.
func candidates(l *binledger.Ledger, args []string) ([]candidate, error) {
	var pend []binledger.Pending
	if len(args) == 0 {
		var err error
		if pend, err = l.ListPending(); err != nil {
			return nil, err
		}
	}
	for _, a := range args {
		pend = append(pend, binledger.Pending{Line: a, Resolved: a, Resources: "(given on the command line)",
			SeenAt: time.Now()})
	}
	out := make([]candidate, 0, len(pend))
	for _, p := range pend {
		c := candidate{Pending: p}
		c.current, c.readErr = hashNow(p.Line)
		if len(args) > 0 {
			c.SHA256 = c.current
		}
		out = append(out, c)
	}
	return out, nil
}

func hashNow(path string) ([32]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return [32]byte{}, err
	}
	defer f.Close()
	e, err := ebpf.ComputeBinaryEntryFile(f, path)
	return e.Hash, err
}

func describe(w io.Writer, c *candidate) {
	fmt.Fprintf(w, "\n%s\n", logging.SanitizeText(c.Line))
	if c.Resolved != c.Line {
		fmt.Fprintf(w, "  resolves to: %s\n", logging.SanitizeText(c.Resolved))
	}
	fmt.Fprintf(w, "  resource:    %s\n", logging.SanitizeText(c.Resources))
	fmt.Fprintf(w, "  refused at:  %s\n", c.SeenAt.Format(time.RFC3339))
	fmt.Fprintf(w, "  sha256:      %x\n", c.SHA256)
	var st unix.Stat_t
	if unix.Stat(c.Line, &st) == nil {
		fmt.Fprintf(w, "  owner uid:   %d   mode: %o   modified: %s\n", st.Uid, st.Mode&0o7777,
			time.Unix(st.Mtim.Sec, 0).Format(time.RFC3339))
	}
}

// confirmable: only the content the daemon refused is confirmed. A file that changed again since is
// judged afresh by the next reload instead.
func confirmable(w io.Writer, c *candidate) bool {
	switch {
	case c.readErr != nil:
		fmt.Fprintf(w, "  cannot read it now (%v): skipped\n", c.readErr)
		return false
	case c.current != c.SHA256:
		fmt.Fprintf(w, "  changed again since the daemon refused it (now %x): skipped — reload the daemon "+
			"to have it judged again\n", c.current)
		return false
	}
	return true
}
