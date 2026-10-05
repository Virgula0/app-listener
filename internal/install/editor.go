package install

import (
	"errors"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Virgula0/app-listener/internal/codeedit"
)

// ErrEditCanceled is returned by EditText when the user aborts the
// embedded editor.
var ErrEditCanceled = errors.New("editing canceled by user")

const (
	editSaveKey   = "ctrl+s"
	editCancelKey = "esc"
	// editMinWidth and editMinHeight keep the editor usable on very small
	// terminals.
	editMinWidth  = 40
	editMinHeight = 8
	// editorChromeLines is the vertical space taken by the title, the key
	// hint and the blank lines around them; the textarea gets the rest.
	editorChromeLines = 6
	// editorSidePadding keeps the textarea off the terminal edges.
	editorSidePadding = 4
)

// EditText opens the embedded multiline editor pre-filled with initial, highlighted as the file
// name names (see codeedit.NewHighlighter). Ctrl+S saves and returns the text; Esc aborts with
// ErrEditCanceled.
func EditText(title, name, initial string) (string, error) {
	m := newEditorModel(title, name, initial)
	p := tea.NewProgram(&m, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		return "", err
	}
	em, ok := result.(*editorModel)
	if !ok {
		return "", errors.New("embedded editor returned an unexpected result")
	}
	if em.canceled {
		return "", ErrEditCanceled
	}
	return em.textarea.Value(), nil
}

type editorModel struct {
	textarea textarea.Model
	hl       *codeedit.Highlighter
	title    string
	canceled bool
}

func newEditorModel(title, name, initial string) editorModel {
	ta := codeedit.New()
	ta.SetHeight(24)
	ta.Placeholder = "Type your configuration here..."
	ta.Focus()
	codeedit.SetText(&ta, initial)
	return editorModel{textarea: ta, hl: codeedit.NewHighlighter(name, initial), title: title}
}

func (m *editorModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m *editorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.WindowSizeMsg); ok {
		m.textarea.SetWidth(max(msg.Width-editorSidePadding, editMinWidth))
		m.textarea.SetHeight(max(msg.Height-editorChromeLines, editMinHeight))
	}
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case editSaveKey:
			m.canceled = false
			return m, tea.Quit
		case editCancelKey:
			m.canceled = true
			return m, tea.Quit
		}
		if codeedit.Edit(&m.textarea, msg) {
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m *editorModel) View() string {
	var b strings.Builder
	b.WriteString(lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("212")).
		Render(m.title))
	b.WriteString("\n\n")
	b.WriteString(lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("Ctrl+S save  ·  Esc cancel (no changes are kept)  ·  Ctrl+←/→ word  ·  PgUp/PgDn page  ·  Tab indent"))
	b.WriteString("\n\n")
	b.WriteString(m.hl.View(&m.textarea))
	return b.String()
}
