package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// sessionModel wraps a live edit-protected session's TUI: every key or mouse input is activity
// (the daemon's idle timeout), and the program quits when the daemon ends the session.
type sessionModel struct {
	inner      tea.Model
	onActivity func()
}

func (m *sessionModel) Init() tea.Cmd { return m.inner.Init() }

func (m *sessionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg, tea.MouseMsg:
		m.onActivity()
	}
	var cmd tea.Cmd
	m.inner, cmd = m.inner.Update(msg)
	return m, cmd
}

func (m *sessionModel) View() string { return m.inner.View() }

// RunSession runs m full-screen for a live session: onActivity runs on every input, and the
// program quits as soon as ended closes. ended reports whether that is why it returned.
func RunSession(m tea.Model, onActivity func(), ended <-chan struct{}) (endedByPeer bool, err error) {
	prog := tea.NewProgram(&sessionModel{inner: m, onActivity: onActivity}, tea.WithAltScreen())
	finished := make(chan struct{})
	peer := make(chan bool, 1)
	go func() {
		select {
		case <-ended:
			peer <- true
			prog.Quit()
		case <-finished:
			peer <- false
		}
	}()
	_, err = prog.Run()
	close(finished)
	return <-peer, err
}

// RunFileEditorSession is RunFileEditor under RunSession.
func RunFileEditorSession(root string, onActivity func(), ended <-chan struct{}) (endedByPeer bool, err error) {
	m := newFileEditModel(root)
	defer m.vault.close()
	return RunSession(m, onActivity, ended)
}
