package install

import (
	"errors"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

// EditText opens the embedded bubbletea multiline editor pre-filled with
// initial. Ctrl+S saves and returns the edited text; Esc aborts with
// ErrEditCanceled.
func EditText(title, initial string) (string, error) {
	m := newEditorModel(title, initial)
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
	title    string
	canceled bool
}

func newEditorModel(title, initial string) editorModel {
	ta := textarea.New()
	ta.SetValue(initial)
	ta.SetHeight(24)
	ta.ShowLineNumbers = true
	ta.Placeholder = "Type your configuration here..."
	ta.Focus()
	ta.KeyMap.LineNext.SetKeys("down", "ctrl+n")
	ta.KeyMap.LinePrevious.SetKeys("up", "ctrl+p")
	ta.KeyMap.CharacterBackward.SetKeys("left")
	ta.KeyMap.CharacterForward.SetKeys("right")
	return editorModel{textarea: ta, title: title}
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
		Render("Ctrl+S save  ·  Esc cancel (no changes are kept)"))
	b.WriteString("\n\n")
	b.WriteString(m.textarea.View())
	return b.String()
}
