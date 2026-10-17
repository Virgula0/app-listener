package update

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Virgula0/app-listener/internal/codeedit"
	"github.com/Virgula0/app-listener/internal/wizard"
)

const (
	changelogTitleStyle = "212"
	changelogHintStyle  = "240"
	changelogMinWidth   = 40
	changelogMinHeight  = 5
	changelogChrome     = 4
)

// changelogText returns the release notes of r, falling back to a short
// placeholder when GitHub did not publish a body for the release.
func changelogText(r *githubRelease) string {
	if notes := strings.TrimSpace(r.Body); notes != "" {
		return notes
	}
	return fmt.Sprintf("No changelog available for this release.\n\nTag: %s\nPublished: %s", r.TagName, r.PublishedAt)
}

// cumulativeChangelog joins the notes of every channel release after installed up to latest,
// newest first: each release's notes cover only the step from its predecessor.
func cumulativeChangelog(installed string, latest *githubRelease, all []githubRelease, channel string) string {
	span, found := releasesSince(installed, latest, all, channel)
	if len(span) == 1 && found {
		return changelogText(&span[0])
	}
	var b strings.Builder
	fmt.Fprintf(&b, "Changes from %s to %s (%d releases)\n", installed, latest.TagName, len(span))
	if found {
		fmt.Fprintf(&b, "Full diff: https://github.com/%s/compare/%s...%s\n", defaultRepo, installed, latest.TagName)
	} else {
		fmt.Fprintf(&b, "The installed version %s is not among the fetched releases: showing every fetched %s release up to %s.\n",
			installed, channel, latest.TagName)
	}
	for i := range span {
		fmt.Fprintf(&b, "\n# %s (published %s)\n\n%s\n", span[i].TagName, span[i].PublishedAt, changelogText(&span[i]))
	}
	return b.String()
}

// releasesSince returns the channel releases newer than installed and not newer than latest, newest
// first; latest is always included. found is false when installed's own release is unknown, so the
// lower bound is missing and every fetched release up to latest is returned. Stable tags compare by
// semver; anything else by published_at, located in all (any channel) so a channel switch works.
func releasesSince(installed string, latest *githubRelease, all []githubRelease, channel string) (span []githubRelease, found bool) {
	w := newReleaseWindow(installed, latest, all, channel)
	for _, r := range filterChannel(all, channel) {
		if r.TagName != latest.TagName && w.contains(&r) {
			span = append(span, r)
		}
	}
	span = append(span, *latest)
	slices.SortStableFunc(span, func(a, b githubRelease) int { return w.newestFirst(&a, &b) })
	return span, w.found
}

// releaseWindow is the (installed, latest] interval releasesSince selects.
type releaseWindow struct {
	installedVer, latestVer stableVersion
	bySemver, found         bool
	installedAt, latestAt   time.Time
}

func newReleaseWindow(installed string, latest *githubRelease, all []githubRelease, channel string) releaseWindow {
	w := releaseWindow{latestAt: latestTime(latest)}
	var installedOK, latestOK bool
	w.installedVer, installedOK = parseStableVersion(installed)
	w.latestVer, latestOK = parseStableVersion(latest.TagName)
	w.bySemver = channel == channelStable && installedOK && latestOK
	for i := range all {
		if all[i].TagName == installed {
			w.installedAt, w.found = latestTime(&all[i]), true
		}
	}
	w.found = w.found || w.bySemver
	return w
}

func (w *releaseWindow) contains(r *githubRelease) bool {
	if v, ok := parseStableVersion(r.TagName); w.bySemver && ok {
		return compareStableVersions(v, w.installedVer) > 0 && compareStableVersions(v, w.latestVer) <= 0
	}
	at := latestTime(r)
	return !at.IsZero() && !at.After(w.latestAt) && (!w.found || at.After(w.installedAt))
}

func (w *releaseWindow) newestFirst(a, b *githubRelease) int {
	va, oka := parseStableVersion(a.TagName)
	vb, okb := parseStableVersion(b.TagName)
	if w.bySemver && oka && okb {
		return compareStableVersions(vb, va)
	}
	return latestTime(b).Compare(latestTime(a))
}

// showChangelog displays release notes in a full-screen read-only TUI
// viewer (or the active wizard UI's). Close it with q, Esc or Enter.
func showChangelog(title, notes string) error {
	if ui := wizard.Active(); ui != nil {
		return ui.ShowText("Changelog up to "+title, notes)
	}
	m := newChangelogModel(title, notes)
	p := tea.NewProgram(&m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}

type changelogModel struct {
	textarea textarea.Model
	hl       *codeedit.Highlighter
	find     codeedit.Finder
	sel      codeedit.Selection
	title    string
	// top is the viewer's first screen row in the last View (mouse coordinates).
	top int
}

func newChangelogModel(title, notes string) changelogModel {
	ta := codeedit.New()
	ta.SetHeight(24)
	ta.Focus()
	codeedit.SetText(&ta, notes)
	return changelogModel{
		textarea: ta,
		hl:       codeedit.NewHighlighter("CHANGELOG.md", notes),
		find:     codeedit.NewFinder(false),
		title:    title,
	}
}

func (m *changelogModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m *changelogModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.textarea.SetWidth(max(msg.Width-2, changelogMinWidth))
		m.textarea.SetHeight(max(msg.Height-changelogChrome, changelogMinHeight))
	}
	if msg, ok := msg.(tea.MouseMsg); ok {
		m.sel.Mouse(&m.textarea, m.hl, msg, msg.X, msg.Y-m.top)
		return m, nil
	}
	if msg, ok := msg.(tea.KeyMsg); ok {
		if m.find.Update(&m.textarea, msg) {
			return m, nil
		}
		if ok, cmd := m.sel.Key(&m.textarea, msg, false); ok {
			return m, cmd
		}
		switch msg.String() {
		case "q", "esc", "enter", "ctrl+q":
			return m, tea.Quit
		}
		if codeedit.Navigate(&m.textarea, msg) {
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m *changelogModel) View() string {
	var b strings.Builder
	b.WriteString(lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(changelogTitleStyle)).
		Render(m.title))
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color(changelogHintStyle)).
		Render("q / Esc / Enter: continue  ·  Ctrl+F find  ·  Shift+arrows/mouse select  ·  Ctrl+C copy"))
	b.WriteString("\n\n")
	m.top = lipgloss.Height(b.String()) - 1
	b.WriteString(m.find.View(m.hl, &m.textarea, &m.sel))
	return b.String()
}
