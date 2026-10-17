package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/unix"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Virgula0/app-listener/internal/codeedit"
	inst "github.com/Virgula0/app-listener/internal/install"
)

// Limits for the embedded editor: binary content is never edited; files past the edit limit are
// read-only.
const (
	maxEditableBytes = 2 << 20 // 2 MiB
	binaryScanBytes  = 8 << 10 // 8 KiB of a file are scanned for NUL bytes

	keyEnter = "enter"
	keyEsc   = "esc"
)

// fileEditMode is the model's current interaction mode.
type fileEditMode int

const (
	modeNav fileEditMode = iota
	modeEdit
	modeInput
	modeConfirm
	modeChmod
	modeChown
)

// fileKind is the type of entry being created (modeInput).
type fileKind int

const (
	createFile fileKind = iota
	createDir
)

// chownUser is one chown-picker candidate (uid -1 keeps current ownership).
type chownUser struct {
	label string
	uid   int
	gid   int
}

// hint pairs a key label with its action for the persistent legend.
type hint struct {
	keys  string
	label string
}

var (
	fileTreeCursor = lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("255"))

	fileEditHintKey = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("141"))

	paneSepStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	// The focused pane's title is underlined across the pane, like tmux's active border.
	paneTitleActive = lipgloss.NewStyle().Bold(true).Underline(true).Foreground(lipgloss.Color("141"))
	paneTitleIdle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

// fileNode is one row of the lazily-expanding tree of the opened vault.
type fileNode struct {
	path     string
	name     string
	isDir    bool
	open     bool
	loaded   bool
	depth    int
	parent   *fileNode
	children []*fileNode
}

// fileEditModel is the two-pane protected-file editor (file tree left; directory listing / editor /
// chmod / chown right). The vault must be unlocked by the caller for the whole session and is
// re-locked when RunFileEditor returns.
type fileEditModel struct {
	root  *fileNode
	vault *vaultFS
	rows  []*fileNode
	cur   int
	top   int // tree scroll offset

	// singleFile: root itself is a guarded regular file (file-vault resource), not a directory: the
	// tree is one childless row and save() writes in place, not temp-file+rename
	// (writeFileInPlace).
	singleFile bool

	mode fileEditMode

	editor    textarea.Model
	hl        *codeedit.Highlighter
	find      codeedit.Finder
	sel       codeedit.Selection
	editPath  string
	original  string
	dirty     bool
	previewOK bool
	// editorOpen keeps the file in the right pane while the tree has the focus (Alt+←); Esc or
	// Ctrl+S in the editor closes it.
	editorOpen bool
	treeHidden bool

	input     textinput.Model
	inputKind fileKind
	inputDir  string

	chmodInput textinput.Model
	chmodPath  string

	chownPath  string
	chownUsers []chownUser
	chownCur   int
	chownTop   int

	confirmPrompt string
	confirmQuit   bool
	pendingDelete string
	pendingOpen   *fileNode

	width  int
	height int
	leftW  int
	rightW int

	status string
}

// RunFileEditor opens a full-screen editor rooted at root, which must already be
// unlocked/plaintext. A directory opens the two-pane tree; a single guarded regular file opens
// straight on it. Returns when closed; fscrypt state is untouched, so the caller must re-lock.
func RunFileEditor(root string) error {
	m := newFileEditModel(root)
	defer m.vault.close()
	// Mouse capture keeps a drag inside the editor (the terminal's own selection spans both
	// panes); RunFileEditorSession captures it too.
	prog := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := prog.Run()
	return err
}

func newFileEditModel(root string) *fileEditModel {
	root = filepath.Clean(root)
	name := filepath.Base(root)
	if name == "" || name == "." || name == "/" {
		name = root
	}
	m := &fileEditModel{
		width:  80,
		height: 24,
		leftW:  40,
		rightW: 39,
	}
	m.editor = codeedit.New()
	m.editor.Blur()
	m.find = codeedit.NewFinder(true)

	vault, vaultErr := openVault(root)
	if vaultErr != nil {
		vault = &vaultFS{base: root}
	}
	m.vault = vault
	if vault.file != nil {
		m.singleFile = true
		m.root = &fileNode{path: root, name: name, loaded: true}
	} else {
		m.root = &fileNode{path: root, name: name, isDir: true, open: true}
		m.status = "loading " + root
		if err := m.loadChildren(m.root); err != nil {
			m.status = "error: " + err.Error()
		}
		if vaultErr != nil {
			m.status = "error: " + vaultErr.Error()
		}
	}
	m.rebuild()
	if m.singleFile {
		// Nothing to navigate to first: open straight into the editor.
		if err := m.enterEditor(); err != nil {
			m.status = "error: " + err.Error()
		}
	}
	m.refit() // sane geometry for the first frame, before the WindowSizeMsg
	return m
}

func (m *fileEditModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m *fileEditModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
	case tea.KeyMsg:
		before := m.layoutKey()
		cmd, quit := m.handleKey(msg)
		if quit {
			return m, tea.Quit
		}
		if m.layoutKey() != before {
			m.refit()
		}
		return m, cmd
	case tea.MouseMsg:
		before := m.layoutKey()
		m.handleMouse(msg)
		if m.layoutKey() != before {
			m.refit()
		}
	}
	return m, nil
}

// layoutKey is what the pane split depends on: a change re-fits before the next render.
func (m *fileEditModel) layoutKey() [3]int {
	b := func(v bool) int {
		if v {
			return 1
		}
		return 0
	}
	return [3]int{int(m.mode), b(m.editorOpen), b(m.treeHidden)}
}

// treeSidebarWidth: tree pane width while editing (the editor is the focus; the tree shrinks to a
// sidebar).
const treeSidebarWidth = 32

// paneLayout returns the tree and right-pane widths for the current mode from the usable content
// width minus one separator column: an open editor gets a sidebar tree, a hidden tree none, other
// modes split evenly.
func (m *fileEditModel) paneLayout() (left, right int) {
	usable := max(m.width-appStyle.GetHorizontalMargins(), 4)
	if m.treeHidden {
		return 0, usable
	}
	if m.mode == modeEdit || m.editorOpen {
		left = min(treeSidebarWidth, usable/3)
	} else {
		left = usable / 2
	}
	left = max(left, 16)
	right = max(usable-left-1, 10)
	return left, right
}

// resize records the new terminal size and re-fits the panes/editor.
func (m *fileEditModel) resize(width, height int) {
	m.width = width
	m.height = height
	m.refit()
}

// refit recomputes the pane geometry and resizes the editor; call on terminal-size or mode change.
func (m *fileEditModel) refit() {
	m.leftW, m.rightW = m.paneLayout()
	m.editor.SetWidth(max(m.rightW, 10))
	m.editor.SetHeight(max(m.contentHeight(), 3))
	m.ensureCursorVisible()
}

func (m *fileEditModel) handleKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	if (m.mode == modeNav || m.mode == modeEdit && !m.find.Active()) && m.paneKey(msg.String()) {
		return nil, false
	}
	switch m.mode {
	case modeNav:
		return m.updateNavKey(msg)
	case modeEdit:
		var cmd tea.Cmd
		m.editor, cmd = m.updateEditKey(msg)
		return cmd, false
	case modeInput:
		var cmd tea.Cmd
		m.input, cmd = m.updateInputKey(msg)
		return cmd, false
	case modeChmod:
		var cmd tea.Cmd
		m.chmodInput, cmd = m.updateChmodKey(msg)
		return cmd, false
	case modeChown:
		return m.updateChownKey(msg)
	case modeConfirm:
		return m.updateConfirmKey(msg)
	}
	return nil, false
}

func (m *fileEditModel) updateNavKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch msg.String() {
	case "j", "down":
		m.move(1)
	case "k", "up":
		m.move(-1)
	case "h", "left":
		m.collapse()
	case "l", "right", keyEnter, "e":
		m.openSelected()
	default:
		return nil, m.navAction(msg)
	}
	return nil, false
}

// navAction handles non-movement nav keys; reports whether the editor must quit.
func (m *fileEditModel) navAction(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "g":
		m.cur = 0
		m.top = 0
	case "G":
		m.cur = len(m.rows) - 1
		m.ensureCursorVisible()
	case "a":
		m.beginCreate(createFile)
	case "A":
		m.beginCreate(createDir)
	case "d":
		m.beginDelete()
	case "m":
		m.beginChmod()
	case "c":
		m.beginChown()
	case "ctrl+s":
		m.savePending()
	case keyEsc, "q", "ctrl+c":
		if msg.String() == keyEsc && m.editorOpen {
			m.closeEditor()
			return false
		}
		if m.editPath != "" && m.dirty {
			m.confirmPrompt = "Discard unsaved changes to " + filepath.Base(m.editPath) + "?"
			m.confirmQuit = true
			m.mode = modeConfirm
			return false
		}
		return true
	}
	return false
}

func (m *fileEditModel) updateEditKey(msg tea.KeyMsg) (textarea.Model, tea.Cmd) {
	if m.find.Update(&m.editor, msg) {
		m.dirty = m.editor.Value() != m.original
		return m.editor, nil
	}
	if ok, cmd := m.sel.Key(&m.editor, msg, true); ok {
		m.dirty = m.editor.Value() != m.original
		if note := m.sel.TakeNote(); note != "" {
			m.status = note
		}
		return m.editor, cmd
	}
	switch msg.String() {
	case "ctrl+s":
		if err := m.save(); err != nil {
			m.status = err.Error()
		}
		m.closeEditor()
	case keyEsc:
		m.closeEditor()
	default:
		var cmd tea.Cmd
		if !codeedit.Edit(&m.editor, msg) {
			m.editor, cmd = m.editor.Update(msg)
		}
		if m.editor.Value() != m.original {
			m.dirty = true
		}
		return m.editor, cmd
	}
	return m.editor, nil
}

func (m *fileEditModel) updateInputKey(msg tea.KeyMsg) (textinput.Model, tea.Cmd) {
	switch msg.String() {
	case keyEnter:
		m.createEntry()
	case keyEsc:
		m.input.Blur()
		m.mode = modeNav
	default:
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m.input, cmd
	}
	return m.input, nil
}

func (m *fileEditModel) updateChmodKey(msg tea.KeyMsg) (textinput.Model, tea.Cmd) {
	switch msg.String() {
	case keyEnter:
		m.applyChmod()
	case keyEsc:
		m.chmodInput.Blur()
		m.mode = modeNav
	default:
		var cmd tea.Cmd
		m.chmodInput, cmd = m.chmodInput.Update(msg)
		return m.chmodInput, cmd
	}
	return m.chmodInput, nil
}

func (m *fileEditModel) updateChownKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch msg.String() {
	case "j", "down":
		m.chownCur = clamp(m.chownCur+1, len(m.chownUsers)-1)
		rows := max(m.contentHeight()-2, 1)
		if m.chownCur >= m.chownTop+rows {
			m.chownTop = m.chownCur - rows + 1
		}
	case "k", "up":
		m.chownCur = clamp(m.chownCur-1, len(m.chownUsers)-1)
		if m.chownCur < m.chownTop {
			m.chownTop = m.chownCur
		}
	case keyEnter:
		m.applyChown()
	case keyEsc:
		m.mode = modeNav
	}
	return nil, false
}

func (m *fileEditModel) updateConfirmKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch msg.String() {
	case "y", keyEnter:
		if m.confirmQuit {
			return nil, true
		}
		if n := m.pendingOpen; n != nil {
			m.pendingOpen, m.dirty, m.mode = nil, false, modeNav
			if i := m.indexOf(n); i >= 0 {
				m.cur = i
			}
			if err := m.enterEditor(); err != nil {
				m.status = err.Error()
			}
			return nil, false
		}
		m.confirmDelete()
	case "n", keyEsc, "q":
		m.pendingOpen = nil
		m.mode = modeNav
	}
	return nil, false
}

// paneKey moves between the panes (Alt+←/→) and hides the tree (Ctrl+B); reports whether it
// used k. Ctrl+←/→ stay the editor's word jumps.
func (m *fileEditModel) paneKey(k string) bool {
	switch k {
	case "alt+left":
		m.focusTree()
	case "alt+right":
		m.focusEditor()
	case "ctrl+b":
		m.toggleTree()
	default:
		return false
	}
	return true
}

// focusTree gives the tree the keyboard; an open file stays on the right.
func (m *fileEditModel) focusTree() {
	m.treeHidden = false
	if m.mode == modeEdit {
		m.find.Close()
		m.sel.Clear()
		m.editor.Blur()
		m.mode = modeNav
	}
}

// focusEditor returns to the open file, else opens the selected file.
func (m *fileEditModel) focusEditor() {
	switch {
	case m.mode == modeEdit:
	case m.editorOpen:
		m.mode = modeEdit
		m.editor.Focus()
		m.status = "editing " + m.editPath
	default:
		if n := m.selected(); n != nil && !n.isDir {
			m.openSelected()
		}
	}
}

func (m *fileEditModel) toggleTree() {
	switch {
	case m.treeHidden:
		m.treeHidden = false
	case !m.editorOpen:
		m.status = "open a file first: the file tree is the only pane"
	default:
		m.treeHidden = true
		m.focusEditor()
	}
}

// closeEditor leaves the editor for the tree and its details pane; unsaved changes stay in the
// buffer until another file is opened.
func (m *fileEditModel) closeEditor() {
	m.find.Close()
	m.sel.Clear()
	m.editor.Blur()
	m.editorOpen, m.treeHidden = false, false
	m.mode = modeNav
}

// handleMouse: a click selects a tree row (a second click opens it) or places the editor's cursor
// and starts a selection; dragging extends it; the wheel scrolls the pane under the pointer.
func (m *fileEditModel) handleMouse(msg tea.MouseMsg) {
	if m.mode != modeNav && m.mode != modeEdit {
		return
	}
	x, y := msg.X-appStyle.GetMarginLeft(), msg.Y-m.contentTop()
	editorX := m.leftW + 1
	if m.treeHidden {
		editorX = 0
	}
	if m.sel.Dragging() || m.editorOpen && x >= editorX {
		m.editorMouse(msg, x-editorX, y)
		return
	}
	if !m.treeHidden && x < m.leftW {
		m.treeMouse(msg, y)
	}
}

// editorMouse: input on the editor focuses it first (a blurred textarea neither scrolls nor
// takes the cursor).
func (m *fileEditModel) editorMouse(msg tea.MouseMsg, x, y int) {
	wheel := msg.Button == tea.MouseButtonWheelUp || msg.Button == tea.MouseButtonWheelDown
	if m.mode == modeNav && (wheel || msg.Action == tea.MouseActionPress) {
		m.focusEditor()
	}
	if m.mode == modeEdit {
		m.sel.Mouse(&m.editor, m.hl, msg, x, y)
	}
}

func (m *fileEditModel) treeMouse(msg tea.MouseMsg, y int) {
	switch {
	case msg.Button == tea.MouseButtonWheelUp:
		m.move(-3)
	case msg.Button == tea.MouseButtonWheelDown:
		m.move(3)
	case msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft:
		i := m.top + y
		if y < 0 || i >= len(m.rows) {
			return
		}
		m.focusTree()
		if i == m.cur {
			m.openSelected()
			return
		}
		m.cur = i
	}
}

// contentTop is the screen row of the panes' first content row: margin, header, pane titles.
func (m *fileEditModel) contentTop() int {
	return appStyle.GetMarginTop() + 2
}

// loadChildren fills a directory node with sorted, lazily loaded children (directories first, then
// files, case-insensitive by name).
func (m *fileEditModel) loadChildren(n *fileNode) error {
	if !n.isDir || n.loaded {
		return nil
	}
	entries, err := m.vault.readDir(n.path)
	if err != nil {
		return err
	}
	n.children = make([]*fileNode, 0, len(entries))
	for _, e := range entries {
		n.children = append(n.children, &fileNode{
			path:   filepath.Join(n.path, e.Name()),
			name:   e.Name(),
			isDir:  e.IsDir(),
			depth:  n.depth + 1,
			parent: n,
		})
	}
	sortNodes(n.children)
	n.loaded = true
	return nil
}

func sortNodes(nodes []*fileNode) {
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].isDir != nodes[j].isDir {
			return nodes[i].isDir
		}
		return strings.ToLower(nodes[i].name) < strings.ToLower(nodes[j].name)
	})
}

// rebuild flattens the visible tree (root plus loaded, open directories) into rows, keeping the
// cursor on the same node when possible.
func (m *fileEditModel) rebuild() {
	selected := m.selected()
	m.rows = m.rows[:0]
	var walk func(n *fileNode)
	walk = func(n *fileNode) {
		m.rows = append(m.rows, n)
		if n.isDir && n.open && n.loaded {
			for _, c := range n.children {
				walk(c)
			}
		}
	}
	walk(m.root)
	m.cur = min(m.cur, len(m.rows)-1)
	if idx := m.indexOf(selected); idx >= 0 {
		m.cur = idx
	}
	m.ensureCursorVisible()
}

func (m *fileEditModel) selected() *fileNode {
	if len(m.rows) == 0 {
		return nil
	}
	return m.rows[m.cur]
}

func (m *fileEditModel) indexOf(n *fileNode) int {
	for i, row := range m.rows {
		if row == n {
			return i
		}
	}
	return -1
}

func (m *fileEditModel) ensureCursorVisible() {
	if m.cur < m.top {
		m.top = m.cur
	}
	rows := m.contentHeight()
	if m.cur >= m.top+rows {
		m.top = m.cur - rows + 1
	}
}

func (m *fileEditModel) move(delta int) {
	m.cur = clamp(m.cur+delta, len(m.rows)-1)
	m.ensureCursorVisible()
}

func (m *fileEditModel) collapse() {
	if n := m.selected(); n != nil && n.isDir && n.open {
		n.open = false
		m.rebuild()
		return
	}
	if n := m.selected(); n != nil && n.parent != nil {
		m.cur = m.indexOf(n.parent)
		m.ensureCursorVisible()
	}
}

func (m *fileEditModel) openSelected() {
	n := m.selected()
	if n == nil {
		return
	}
	if n.isDir {
		m.toggleDir(n)
		return
	}
	if err := m.enterEditor(); err != nil {
		m.status = err.Error()
	}
}

func (m *fileEditModel) toggleDir(n *fileNode) {
	if err := m.loadChildren(n); err != nil {
		m.status = err.Error()
		return
	}
	n.open = !n.open
	m.rebuild()
	m.cur = m.indexOf(n)
	m.ensureCursorVisible()
}

// enterEditor opens the selected file in the focused textarea. Re-entering the same unsaved file
// keeps the buffer; a different file reloads. Symlinks are never followed.
func (m *fileEditModel) enterEditor() error {
	n := m.selected()
	if n == nil || n.isDir {
		return nil
	}
	info, err := m.vault.lstat(n.path)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		m.status = n.path + " is a symlink — editing it is refused"
		return nil
	}
	if m.editPath == n.path && (m.dirty || m.editorOpen) {
		m.mode, m.editorOpen = modeEdit, true
		m.editor.Focus()
		m.status = "editing " + n.path
		return nil
	}
	if m.dirty && m.editPath != "" {
		m.confirmPrompt = "Discard unsaved changes to " + filepath.Base(m.editPath) + "?"
		m.pendingOpen = n
		m.mode = modeConfirm
		return nil
	}
	if err := m.previewFile(n.path); err != nil {
		return err
	}
	if !m.previewOK {
		m.editorOpen = false
		m.status = "cannot edit " + n.path
		return nil
	}
	m.mode, m.editorOpen = modeEdit, true
	m.editor.Focus()
	m.status = "editing " + n.path
	return nil
}

// previewFile loads path into the (blurred) editor; binary files and files over maxEditableBytes
// reset the preview and report a status.
func (m *fileEditModel) previewFile(path string) error {
	info, err := m.vault.lstat(path)
	if err != nil {
		return err
	}
	if info.Size() > maxEditableBytes {
		m.previewOK = false
		m.editor.SetValue("")
		m.status = fmt.Sprintf("%s is %.1f MiB — only files up to 2 MiB can be edited", path, float64(info.Size())/(1<<20))
		return nil
	}
	data, err := m.vault.readFile(path)
	if err != nil {
		return err
	}
	if isBinaryContent(data) {
		m.previewOK = false
		m.editor.SetValue("")
		m.status = path + " looks like a binary file — it cannot be edited"
		return nil
	}
	m.editPath = path
	m.original = string(data)
	m.hl = codeedit.NewHighlighter(path, m.original)
	m.find.Close()
	m.sel.Clear()
	codeedit.SetText(&m.editor, m.original)
	m.previewOK = true
	m.dirty = false
	if m.mode != modeEdit {
		m.editor.Blur()
	}
	return nil
}

func (m *fileEditModel) save() error {
	if m.editPath == "" {
		return nil
	}
	if err := m.vault.write(m.editPath, []byte(m.editor.Value())); err != nil {
		return err
	}
	m.original = m.editor.Value()
	m.dirty = false
	m.status = "saved " + m.editPath
	return nil
}

func (m *fileEditModel) savePending() {
	if m.editPath == "" || !m.dirty {
		return
	}
	if err := m.save(); err != nil {
		m.status = err.Error()
	}
}

// beginCreate switches to the new-name input mode; the entry is created in the selected directory
// (or the selected file's parent).
func (m *fileEditModel) beginCreate(kind fileKind) {
	n := m.selected()
	if n == nil {
		return
	}
	if n == m.root && !n.isDir {
		m.status = "cannot create inside a single-file resource"
		return
	}
	dir := n.path
	if !n.isDir {
		dir = n.parent.path
	}
	m.input = textinput.New()
	m.input.Cursor.SetMode(cursor.CursorStatic)
	m.input.Width = max(m.paneWidth()-12, 4)
	m.input.Placeholder = "name"
	m.input.Focus()
	m.inputKind = kind
	m.inputDir = dir
	m.mode = modeInput
	m.status = ""
}

func (m *fileEditModel) createEntry() {
	name := strings.TrimSpace(m.input.Value())
	m.input.Blur()
	m.mode = modeNav
	if name == "" || strings.Contains(name, "/") {
		m.status = "invalid name: must be non-empty and contain no '/'"
		return
	}
	target := filepath.Join(m.inputDir, name)
	var err error
	switch m.inputKind {
	case createFile:
		// O_EXCL|O_NOFOLLOW: an existing symlink at target is refused, never followed.
		err = m.vault.dirOp(target, func(r *os.Root, name string) error {
			f, openErr := r.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_EXCL|unix.O_NOFOLLOW, 0o600)
			if openErr != nil {
				return openErr
			}
			return f.Close()
		})
	case createDir:
		err = m.vault.dirOp(target, func(r *os.Root, name string) error { return r.Mkdir(name, 0o750) })
	}
	if err != nil {
		m.status = err.Error()
		return
	}
	if err := applyNewFileMeta(m.vault, target, m.inputKind == createDir); err != nil {
		m.status = err.Error()
	}
	if parent := m.findDir(m.inputDir); parent != nil {
		parent.open = true
		parent.loaded = false
		if err := m.loadChildren(parent); err != nil {
			m.status = err.Error()
		}
		m.rebuild()
		if n := m.nodeByName(m.inputDir, name); n != nil {
			m.cur = m.indexOf(n)
			m.ensureCursorVisible()
		}
	}
	m.status = "created " + target
}

func (m *fileEditModel) nodeByName(dir, name string) *fileNode {
	for _, row := range m.rows {
		if row.parent != nil && row.parent.path == dir && row.name == name {
			return row
		}
	}
	return nil
}

func (m *fileEditModel) findDir(path string) *fileNode {
	for _, row := range m.rows {
		if row.isDir && row.path == path {
			return row
		}
	}
	return nil
}

// beginChmod prepares the octal mode editor: 3-4 digits, the optional leading digit carrying
// setuid/setgid/sticky (e.g. 4755).
func (m *fileEditModel) beginChmod() {
	n := m.selected()
	if n == nil {
		return
	}
	info, err := m.vault.lstat(n.path)
	if err != nil {
		m.status = err.Error()
		return
	}
	if info.Mode()&os.ModeSymlink != 0 {
		m.status = n.path + " is a symlink — chmod is refused"
		return
	}
	m.chmodPath = n.path
	m.chmodInput = textinput.New()
	m.chmodInput.Cursor.SetMode(cursor.CursorStatic)
	m.chmodInput.Width = max(m.paneWidth()-12, 4)
	m.chmodInput.Prompt = "mode: "
	m.chmodInput.Placeholder = "0644"
	// Prefill the full mode incl. setuid/setgid/sticky so confirming untouched input can't strip
	// them.
	m.chmodInput.SetValue(fmt.Sprintf("%04o", fileModeToUnixMode(info.Mode())))
	m.chmodInput.Focus()
	m.mode = modeChmod
	m.status = ""
}

func (m *fileEditModel) applyChmod() {
	v := strings.TrimSpace(m.chmodInput.Value())
	m.chmodInput.Blur()
	m.mode = modeNav
	if !isOctalMode(v) {
		m.status = "invalid mode: 1-4 octal digits (e.g. 0644)"
		return
	}
	perm, _ := strconv.ParseUint(v, 8, 16)
	if perm > 0o7777 {
		m.status = "invalid mode: must not exceed 07777"
		return
	}
	// Re-validate at apply time on an O_PATH|O_NOFOLLOW descriptor: beginChmod's Lstat check is
	// seconds stale. The value is a unix mode: convert via the converter, or setuid/setgid/sticky
	// are lost (a raw os.FileMode cast maps them onto bits chmod ignores).
	mode := unixModeToFileMode(uint32(perm))
	err := m.vault.withFile(m.chmodPath, unix.O_PATH, func(f *os.File) error {
		if info, statErr := f.Stat(); statErr != nil || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("%s is a symlink — chmod is refused", m.chmodPath)
		}
		// O_PATH descriptors can't fchmod directly; chmod via /proc/self/fd resolves the pinned
		// inode without re-walking the (swappable) path.
		return os.Chmod(fmt.Sprintf("/proc/self/fd/%d", f.Fd()), mode)
	})
	if err != nil {
		m.status = err.Error()
		return
	}
	m.status = "chmod " + m.chmodPath + " = " + v
}

// beginChown opens the owner picker (every login user plus "keep current").
func (m *fileEditModel) beginChown() {
	n := m.selected()
	if n == nil {
		return
	}
	info, err := m.vault.lstat(n.path)
	if err != nil {
		m.status = err.Error()
		return
	}
	if info.Mode()&os.ModeSymlink != 0 {
		m.status = n.path + " is a symlink — chown is refused"
		return
	}
	m.chownPath = n.path
	m.chownUsers = make([]chownUser, 0, 4)
	current := "keep current"
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		current = fmt.Sprintf("keep current (uid %d, gid %d)", stat.Uid, stat.Gid)
	}
	m.chownUsers = append(m.chownUsers, chownUser{label: current, uid: -1, gid: -1})
	users, err := inst.ListUsers()
	if err != nil {
		m.status = err.Error()
		return
	}
	for _, u := range users {
		m.chownUsers = append(m.chownUsers, chownUser{
			label: fmt.Sprintf("%s (uid %d, gid %d)", u.Name, u.UID, u.GID),
			uid:   int(u.UID),
			gid:   int(u.GID),
		})
	}
	m.chownCur = 0
	m.chownTop = 0
	m.mode = modeChown
	m.status = ""
}

func (m *fileEditModel) applyChown() {
	if len(m.chownUsers) == 0 {
		m.mode = modeNav
		return
	}
	u := m.chownUsers[clamp(m.chownCur, len(m.chownUsers)-1)]
	if u.uid < 0 {
		m.status = "ownership of " + m.chownPath + " unchanged"
		m.mode = modeNav
		return
	}
	err := m.vault.withFile(m.chownPath, unix.O_PATH, func(f *os.File) error {
		if info, statErr := f.Stat(); statErr != nil || info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("%s is a symlink — chown is refused", m.chownPath)
		}
		return unix.Fchownat(int(f.Fd()), "", u.uid, u.gid, unix.AT_EMPTY_PATH)
	})
	if err != nil {
		m.status = err.Error()
		m.mode = modeNav
		return
	}
	m.status = "chown " + m.chownPath + " → " + u.label
	m.mode = modeNav
}

// beginDelete asks for confirmation (default no) before removing the entry (recursive for
// directories).
func (m *fileEditModel) beginDelete() {
	n := m.selected()
	if n == nil {
		return
	}
	if n == m.root {
		m.status = "cannot delete the root of the vault"
		return
	}
	prompt := "Delete " + n.path + "?"
	if n.isDir {
		prompt += " (recursively)"
	}
	m.confirmPrompt = prompt
	m.pendingDelete = n.path
	m.mode = modeConfirm
}

func (m *fileEditModel) confirmDelete() {
	if m.pendingDelete == "" {
		m.mode = modeNav
		return
	}
	path := m.pendingDelete
	m.pendingDelete = ""
	m.mode = modeNav
	if err := m.vault.dirOp(path, (*os.Root).RemoveAll); err != nil {
		m.status = err.Error()
		return
	}
	m.status = "deleted " + path
	if parent := m.findDir(filepath.Dir(path)); parent != nil {
		parent.open = true
		parent.loaded = false
		_ = m.loadChildren(parent)
		m.rebuild()
	}
}

// paneWidth is the usable width of the right pane.
func (m *fileEditModel) paneWidth() int {
	return max(m.rightW, 10)
}

func (m *fileEditModel) paneHeight() int {
	if m.height < 6 {
		return 2
	}
	// The margins, the header row and the legend's wrapped rows. A frame taller than the
	// terminal loses its top rows, which shifts every mouse row by as many.
	return max(m.height-appStyle.GetVerticalMargins()-1-m.legendRows(), 2)
}

// contentHeight is a pane's height below its title row.
func (m *fileEditModel) contentHeight() int {
	return max(m.paneHeight()-1, 1)
}

// legendRows: terminal rows the width-wrapped legend occupies in the current mode.
func (m *fileEditModel) legendRows() int {
	return strings.Count(m.legend(), "\n") + 1
}

func (m *fileEditModel) View() string {
	if m.width == 0 {
		return initializingView
	}
	status := m.status
	if m.dirty && m.editPath != "" {
		if status != "" {
			status += " · "
		}
		status += "unsaved " + filepath.Base(m.editPath)
	}
	head := lipgloss.JoinHorizontal(lipgloss.Left,
		headerStyle.Render("app-listener — edit protected directory"),
		"  ",
		infoStyle.Render(status),
	)
	return appStyle.Render(lipgloss.JoinVertical(lipgloss.Left, head, m.body(), m.legend()))
}

// body lays the tree and info/editor pane side by side, each fixed to its column width and titled,
// the focused pane's title underlined.
func (m *fileEditModel) body() string {
	h, ph := m.contentHeight(), m.paneHeight()
	// Fix the right column to rightW so it reaches the terminal edge (renderRight clips to
	// rightW, so Width only pads; MaxWidth guards a stray wide line), and both panes to exactly
	// paneHeight rows: see paneHeight.
	pane := lipgloss.NewStyle().Height(ph).MaxHeight(ph)
	right := pane.Width(m.rightW).MaxWidth(m.rightW).Render(
		paneTitle(m.rightTitle(), m.rightW, m.mode != modeNav) + "\n" +
			strings.TrimSuffix(m.renderRight(m.rightW, h), "\n"))
	if m.treeHidden {
		return right
	}
	left := pane.Render(paneTitle("Files", m.leftW, m.mode == modeNav) + "\n" +
		strings.TrimSuffix(m.renderTree(m.leftW, h), "\n"))
	sep := strings.TrimSuffix(strings.Repeat("│\n", ph), "\n")
	return lipgloss.JoinHorizontal(lipgloss.Top, left, paneSepStyle.Render(sep), right)
}

func paneTitle(text string, width int, active bool) string {
	label := clipLabel(" "+text, width)
	label += strings.Repeat(" ", max(0, width-lipgloss.Width(label)))
	if active {
		return paneTitleActive.Render(label)
	}
	return paneTitleIdle.Render(label)
}

func (m *fileEditModel) rightTitle() string {
	switch {
	case m.mode == modeInput:
		return "New entry"
	case m.mode == modeChmod:
		return "Permissions"
	case m.mode == modeChown:
		return "Owner"
	case m.mode == modeConfirm:
		return "Confirm"
	case !m.editorOpen:
		return "Details"
	}
	t := "Editor · " + filepath.Base(m.editPath)
	if m.dirty {
		t += " ●"
	}
	if m.treeHidden {
		t += "   (ctrl+b shows the files)"
	}
	return t
}

func (m *fileEditModel) renderTree(width, height int) string {
	var b strings.Builder
	for i := m.top; i < len(m.rows) && i < m.top+height; i++ {
		n := m.rows[i]
		marker := "  "
		if n.isDir {
			marker = "▸ "
			if n.open {
				marker = "▾ "
			}
		}
		prefix := strings.Repeat("  ", n.depth)
		// Pad every row to the full pane width: otherwise JoinHorizontal sizes the tree column to
		// its longest line and the right pane never reaches the edge. Padding first also makes the
		// selection bar span the pane.
		label := clipLabel(prefix+marker+n.name, width)
		label += strings.Repeat(" ", max(0, width-lipgloss.Width(label)))
		if i == m.cur {
			label = fileTreeCursor.Render(label)
		}
		b.WriteString(label + "\n")
	}
	return b.String()
}

func (m *fileEditModel) renderRight(width, height int) string {
	switch m.mode {
	case modeNav, modeEdit:
		if m.editorOpen {
			return m.find.View(m.hl, &m.editor, &m.sel)
		}
	case modeInput:
		kind := "file"
		if m.inputKind == createDir {
			kind = "directory"
		}
		return "New " + kind + " in " + clipLabel(m.inputDir, width-20) + "\n\n" + m.input.View()
	case modeChmod:
		return "chmod " + clipLabel(m.chmodPath, width-20) + "\n\n" + m.chmodInput.View()
	case modeChown:
		return m.renderChownPicker(width, height)
	case modeConfirm:
		return m.confirmPrompt + "\n\n  y/↵ confirm · n/esc cancel"
	}
	return m.renderEntryInfo(width, height)
}

// renderChownPicker draws the scrolling list of chown candidates.
func (m *fileEditModel) renderChownPicker(width, height int) string {
	var b strings.Builder
	b.WriteString("chown " + clipLabel(m.chownPath, width-20))
	b.WriteString("\n\n")
	rows := max(height-2, 1)
	for i := m.chownTop; i < len(m.chownUsers) && i < m.chownTop+rows; i++ {
		line := clipLabel(m.chownUsers[i].label, max(1, width-4))
		if i == m.chownCur {
			line = fileTreeCursor.Render(line)
		}
		b.WriteString(line + "\n")
	}
	return b.String()
}

// renderEntryInfo draws the right pane: a directory listing or file metadata.
func (m *fileEditModel) renderEntryInfo(width, height int) string {
	n := m.selected()
	if n == nil {
		return "no selection"
	}
	var b strings.Builder
	b.WriteString(infoStyle.Render(clipLabel(n.path, width-infoStyle.GetHorizontalFrameSize())))
	b.WriteString("\n\n")
	if n.isDir {
		if err := m.loadChildren(n); err != nil {
			b.WriteString(err.Error())
			return b.String()
		}
		for _, e := range listEntries(m.vault, n, height-2) {
			name := e.name
			if e.isDir {
				name += "/"
			}
			fmt.Fprintf(&b, "%8s  %s  %s\n",
				humanSize(e.size),
				e.mtime.Format("2006-01-02 15:04"),
				clipLabel(name, max(1, width-28)),
			)
		}
		return b.String()
	}
	info, err := m.vault.lstat(n.path)
	if err != nil {
		b.WriteString(err.Error())
		return b.String()
	}
	fmt.Fprintf(&b, "size   %s\nmode   %s\nmtime  %s\n\n",
		humanSize(info.Size()), info.Mode().String(), info.ModTime().Format("2006-01-02 15:04"))
	if n.path == m.editPath {
		b.WriteString("e/↵ to edit the file (ctrl+s saves) · m chmod · c chown")
	} else {
		b.WriteString("e/↵ to open the file in the editor · m chmod · c chown")
	}
	return b.String()
}

func (m *fileEditModel) legend() string {
	parts := make([]string, 0, len(legendHints[m.mode]))
	for _, h := range legendHints[m.mode] {
		parts = append(parts, fileEditHintKey.Render(h.keys)+" "+h.label)
	}
	// Wrap to the terminal so the wide nav-mode hint list doesn't spill past the edge; paneHeight
	// accounts for the extra rows.
	usable := max(m.width-appStyle.GetHorizontalMargins(), 20)
	return footerStyle.Width(usable).Render(strings.Join(parts, "  ·  "))
}

var legendHints = map[fileEditMode][]hint{
	modeNav: {
		{"↑/↓ j/k", "move"},
		{"→/↵ l", "open"},
		{"← h", "back"},
		{"alt+→", "editor"},
		{"ctrl+b", "toggle files"},
		{"g/G", "top/bottom"},
		{"a/A", "new file/dir"},
		{"d", "delete"},
		{"m", "chmod"},
		{"c", "chown"},
		{"ctrl+s", "save"},
		{"q/esc", "quit (esc: close the editor first)"},
	},
	modeEdit: {
		{"ctrl+s", "save & close"},
		{"esc", "close"},
		{"alt+←", "files"},
		{"ctrl+b", "toggle files"},
		{"⇧+arrows/mouse", "select"},
		{"ctrl+c/x/v", "copy/cut/paste"},
		{"ctrl+f", "find"},
		{"ctrl+r", "replace"},
		{"ctrl+←/→", "word jump"},
		{"pgup/pgdn", "page"},
		{"tab/⇧tab", "indent"},
	},
	modeInput: {
		{"↵", "create"},
		{"esc", "cancel"},
	},
	modeChmod: {
		{"octal", "mode (e.g. 0644)"},
		{"↵", "apply"},
		{"esc", "cancel"},
	},
	modeChown: {
		{"j/k", "move"},
		{"↵", "apply"},
		{"esc", "cancel"},
	},
	modeConfirm: {
		{"y/↵", "confirm"},
		{"n/esc", "cancel"},
	},
}

func clamp(v, hi int) int {
	if v < 0 {
		return 0
	}
	if v > hi {
		return hi
	}
	return v
}
