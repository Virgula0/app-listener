package editprotected

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/fscrypt"
	inst "github.com/Virgula0/app-listener/internal/install"
	"github.com/Virgula0/app-listener/internal/logging"
	"github.com/Virgula0/app-listener/internal/recorder"
	"github.com/Virgula0/app-listener/internal/wizard"
)

// runRecord is `--edit-config --record <path>...`: record which binaries use the paths, review
// them, optionally encrypt the paths, then add the sections to the running daemon's daemon.conf.
func runRecord() error {
	if !term.IsTerminal(os.Stdin.Fd()) {
		return errors.New("--record needs a terminal: it records until you stop it, then asks you to review")
	}
	password := os.Getenv(editPasswordEnv)
	if password == "" {
		var err error
		if password, err = promptPassword("Edit-protected password"); err != nil {
			return err
		}
	}
	flow, err := newRecordFlow(password, recordFlags)
	if err != nil {
		return err
	}
	defer flow.close()
	err = flow.start()
	if err == nil {
		err = runRecordView(flow)
	}
	if err != nil {
		return err
	}
	rec := flow.stop()
	printRecording(os.Stdout, rec)
	sel, ok, err := reviewRecording(flow, rec)
	if err != nil || !ok {
		if errors.Is(err, huh.ErrUserAborted) || (err == nil && !ok) {
			log.Info("nothing applied — daemon.conf is unchanged")
			return nil
		}
		return err
	}
	return confirmAndApply(flow, &sel)
}

type recordTick struct{}

// recordModel shows a running recording until the operator stops it.
type recordModel struct {
	flow  *recordFlow
	prog  recorder.Progress
	width int
}

func runRecordView(flow *recordFlow) error {
	m := &recordModel{flow: flow, prog: flow.rec.Progress()}
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	return err
}

func recordTickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(time.Time) tea.Msg { return recordTick{} })
}

func (m *recordModel) Init() tea.Cmd { return recordTickCmd() }

func (m *recordModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "enter", "esc":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case recordTick:
		m.prog = m.flow.rec.Progress()
		if m.prog.Elapsed >= recorder.MaxDuration {
			return m, tea.Quit
		}
		return m, recordTickCmd()
	}
	return m, nil
}

var (
	recordTitle = lipgloss.NewStyle().Bold(true)
	recordMuted = lipgloss.NewStyle().Faint(true)
	recordLive  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9"))
)

func (m *recordModel) View() string {
	var b strings.Builder
	b.WriteString(recordLive.Render("● REC") + "  " + recordTitle.Render(formatElapsed(m.prog.Elapsed)) + "\n\n")
	b.WriteString("Recording every access to:\n")
	for _, t := range m.flow.targets {
		b.WriteString("  " + logging.SanitizeText(t.Path) + "\n")
	}
	fmt.Fprintf(&b, "\n%d accesses by %d executables so far\n\n", m.prog.Events, m.prog.Binaries)
	for _, s := range m.prog.Recent {
		line := fmt.Sprintf("  %-6s ago  %6d×  %-28s %s", formatElapsed(time.Since(s.Last)), s.Count,
			strings.Join(s.Ops, ","), logging.SanitizeText(s.Exe))
		if m.width > 4 && len(line) > m.width-1 {
			line = line[:m.width-1]
		}
		b.WriteString(line + "\n")
	}
	b.WriteString("\n" + recordMuted.Render("Use the applications that need these paths now (open, save, sync …). "+
		"Anything that touches them is listed for your review, malware included.") + "\n")
	b.WriteString(recordMuted.Render("Enter, q or Ctrl+C stops the recording.") + "\n")
	return b.String()
}

func formatElapsed(d time.Duration) string {
	d = d.Truncate(time.Second)
	h, mnt, s := int(d.Hours()), int(d.Minutes())%60, int(d.Seconds())%60
	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, mnt, s)
	}
	return fmt.Sprintf("%02d:%02d", mnt, s)
}

// printRecording lists everything recorded. Paths and comms come from other processes: sanitized.
func printRecording(w io.Writer, rec *recorder.Recording) {
	fmt.Fprintf(w, "\nRecorded %s, %d accesses.\n", formatElapsed(rec.Stopped.Sub(rec.Started)), rec.Events)
	if rec.Dropped > 0 {
		fmt.Fprintf(w, "  %d events were dropped (too many distinct executables or libraries)\n", rec.Dropped)
	}
	if len(rec.Binaries) == 0 {
		fmt.Fprintln(w, "  Nothing touched the recorded paths.")
		return
	}
	for i := range rec.Binaries {
		printBinary(w, &rec.Binaries[i])
	}
	fmt.Fprintln(w)
}

func printBinary(w io.Writer, b *recorder.Binary) {
	name := b.Line
	if name == "" {
		name = b.Exe
	}
	if name == "" {
		name = "(unnamed executable)"
	}
	fmt.Fprintf(w, "\n%s\n", logging.SanitizeText(name))
	fmt.Fprintf(w, "  comm %s, uid %v, %d accesses (%s)\n", logging.SanitizeText(strings.Join(b.Comms, ",")),
		b.UIDs, b.Count, strings.Join(b.Ops, ","))
	fmt.Fprintf(w, "  paths: %s\n", logging.SanitizeText(strings.Join(b.Targets, ", ")))
	if !b.Admissible() {
		fmt.Fprintf(w, "  NOT WHITELISTABLE: %s\n", logging.SanitizeText(b.Problem))
		return
	}
	fmt.Fprintf(w, "  sha256 %s\n", b.SHA256)
	for _, n := range b.Notes {
		fmt.Fprintf(w, "  ! %s\n", n)
	}
	for _, l := range b.Libs {
		if l.Problem != "" {
			fmt.Fprintf(w, "  library %s: %s\n", logging.SanitizeText(l.Path), l.Problem)
		} else {
			fmt.Fprintf(w, "  library %s needs allow_lib\n", logging.SanitizeText(l.Path))
		}
	}
	if b.TrustedLibs > 0 {
		fmt.Fprintf(w, "  %d more libraries are already trusted (system or guarded)\n", b.TrustedLibs)
	}
}

// reviewRecording asks which binaries, libraries and encryptions to apply. ok=false: nothing to do.
func reviewRecording(flow *recordFlow, rec *recorder.Recording) (recorder.Selection, bool, error) {
	var sel recorder.Selection
	bins, err := pickBinaries(rec)
	if err != nil {
		return sel, false, err
	}
	if len(bins) == 0 {
		empty, cerr := wizard.Confirm(wizard.Question{
			Title: "No binary selected: guard the paths with an empty whitelist?",
			Description: "Every access to them would be denied until you edit the section " +
				"(edit-protected --edit-config).",
			Affirmative: "Guard, deny all", Negative: "Cancel",
		})
		if cerr != nil || !empty {
			return sel, false, cerr
		}
	}
	sel.Binaries = lines(bins)
	if sel.Libs, err = pickLibs(bins); err != nil {
		return sel, false, err
	}
	if len(bins) > 0 {
		if sel.Masks, err = wizard.Confirm(wizard.Question{
			Title: "Restrict each binary to the events it was seen doing?",
			Description: "Least privilege, but an operation the recording missed is then denied. Without it, " +
				"each binary gets every event, like most catalog entries.",
			Affirmative: "Restrict", Negative: "All events",
		}); err != nil {
			return sel, false, err
		}
	}
	sel.Encrypt, err = pickEncryption(flow)
	return sel, err == nil, err
}

func pickBinaries(rec *recorder.Recording) ([]*recorder.Binary, error) {
	var cands []*recorder.Binary
	var choices []wizard.Choice
	for i := range rec.Binaries {
		b := &rec.Binaries[i]
		if !b.Admissible() {
			continue
		}
		label := logging.SanitizeText(b.Line) + "  (" + strings.Join(b.Ops, ",") + ")"
		if !b.SystemPlaced {
			label += "  [not root-placed]"
		}
		if b.GeneralTool {
			label += "  [general tool]"
		}
		if b.ShellHelper {
			label += "  [shell prompt/hook]"
		}
		cands = append(cands, b)
		choices = append(choices, wizard.Choice{Label: label, Selected: b.Suggested})
	}
	if len(cands) == 0 {
		return nil, nil
	}
	picked, err := wizard.MultiSelect("Whitelist which binaries?",
		"Only root-placed binaries that aren't general tools start selected. Select only what you ran on "+
			"purpose: anything else that touched the paths during the recording is listed too.",
		choices, min(len(choices)+2, 16))
	if err != nil {
		return nil, err
	}
	out := make([]*recorder.Binary, len(picked))
	for i, p := range picked {
		out[i] = cands[p]
	}
	return out, nil
}

func lines(bins []*recorder.Binary) []string {
	out := make([]string, len(bins))
	for i, b := range bins {
		out[i] = b.Line
	}
	return out
}

func pickLibs(bins []*recorder.Binary) ([]string, error) {
	var libs []string
	for _, b := range bins {
		for _, l := range b.Libs {
			if l.Problem == "" && !slices.Contains(libs, l.Path) {
				libs = append(libs, l.Path)
			}
		}
	}
	if len(libs) == 0 {
		return nil, nil
	}
	choices := make([]wizard.Choice, len(libs))
	for i, l := range libs {
		choices[i] = wizard.Choice{Label: logging.SanitizeText(l)}
	}
	picked, err := wizard.MultiSelect("Trust which libraries (allow_lib)?",
		"A whitelisted binary may load only trusted code, which is what stops LD_PRELOAD and library injection. "+
			"Trust is daemon-wide and write-protects the file: approve only libraries the application itself "+
			"ships. None start selected.", choices, min(len(choices)+2, 16))
	if err != nil {
		return nil, err
	}
	out := make([]string, len(picked))
	for i, p := range picked {
		out[i] = libs[p]
	}
	return out, nil
}

func pickEncryption(flow *recordFlow) ([]string, error) {
	blockers := flow.encryptBlockers()
	var out []string
	for _, t := range flow.targets {
		if why := blockers[t.Path]; why != "" {
			log.Warnf("%s will be guarded unencrypted: %s", t.Path, why)
			continue
		}
		yes, err := wizard.Confirm(wizard.Question{
			Title: "Encrypt " + t.Path + " at rest?",
			Description: "It is moved to " + t.Path + fscrypt.BackupSuffix + ", re-created encrypted with the " +
				"daemon's key and locked; the reload attaches its guard, then unlocks it. Without encryption " +
				"it is guarded only while the daemon runs.",
			Affirmative: "Encrypt", Negative: "Guard only", Default: true,
		})
		if err != nil {
			return nil, err
		}
		if yes {
			out = append(out, t.Path)
		}
	}
	return out, nil
}

// confirmAndApply shows the diff, applies on confirmation, then offers to delete the backups.
func confirmAndApply(flow *recordFlow, sel *recorder.Selection) error {
	base, text, err := flow.preview(sel)
	if err != nil {
		return err
	}
	if derr := inst.ShowDiff("daemon.conf — sections to add", inst.UnifiedDiff(string(base), text)); derr != nil {
		return derr
	}
	ok, err := wizard.Confirm(wizard.Question{
		Title: "Apply and reload the daemon?",
		Description: "If anything fails, the encrypted paths are restored from their backups and daemon.conf " +
			"keeps its current content.",
		Affirmative: "Apply", Negative: "Cancel",
	})
	if err != nil || !ok {
		log.Info("nothing applied — daemon.conf is unchanged")
		return err
	}
	var out *recorder.Outcome
	err = wizard.WithBottomBar(func(bar *wizard.BottomBar) error {
		var aerr error
		out, aerr = flow.apply(func(step string, fraction float64) { bar.Set(step, fraction) })
		return aerr
	})
	if err != nil {
		return err
	}
	log.Info("the recorded paths are added to daemon.conf and guarded by the daemon")
	for _, n := range out.Notes {
		log.Warn(n)
	}
	return offerBackupDeletion(flow, out)
}

func offerBackupDeletion(flow *recordFlow, out *recorder.Outcome) error {
	for _, b := range out.Backups {
		remove, err := wizard.Confirm(wizard.Question{
			Title: "Remove the unencrypted backup " + b + "?",
			Description: "It is a plaintext copy of what is now encrypted, outside every guarded tree: anything " +
				"running as its owner can read it. Keep it only until you checked the encrypted copy.",
			Affirmative: "Remove", Negative: "Keep", Default: true,
		})
		if err != nil {
			return err
		}
		if !remove {
			log.Warnf("kept %s: delete it later with `app-listener install --delete-post-backups`", b)
			continue
		}
		if err := flow.deleteBackup(b); err != nil {
			return err
		}
		log.Infof("removed %s", b)
	}
	return nil
}
