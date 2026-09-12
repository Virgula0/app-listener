package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	inst "github.com/Virgula0/app-listener/internal/install"
)

// Limits for the embedded file preview/editor: binary content is never
// edited, and files past the edit limit are shown read-only.
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

// chownUser is one candidate in the chown picker. A uid of -1 keeps the
// current ownership.
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

// fileEditModel is the two-pane protected-file editor: a file tree on the
// left and a directory listing / file editor / chmod / chown pane on the
// right. The vault it edits is expected to be unlocked (provisioned) by the
// caller for the whole session and is re-locked when RunFileEditor returns.
type fileEditModel struct {
	root *fileNode
	rows []*fileNode
	cur  int
	top  int // tree scroll offset

	// singleFile is true when root itself is a guarded regular file (a
	// file-vault resource) rather than a directory: there is nothing to
	// browse, so the tree is a single, childless row and save() writes
	// in place instead of through the tree editor's usual temp-file+rename
	// (see writeFileInPlace).
	singleFile bool

	mode fileEditMode

	editor    textarea.Model
	editPath  string
	original  string
	dirty     bool
	previewOK bool

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

	width  int
	height int
	leftW  int
	rightW int

	status string
}

// RunFileEditor opens a full-screen editor rooted at root, which must already
// be unlocked/plaintext. root is usually a directory (the two-pane tree
// browser); when it is a single guarded regular file (a file-vault resource)
// the tree is skipped and the editor opens straight on that one file. It
// returns when the editor is closed; nothing about the fscrypt state is
// changed, so the caller must re-lock the vault afterwards.
func RunFileEditor(root string) error {
	m := newFileEditModel(root)
	prog := tea.NewProgram(m, tea.WithAltScreen())
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
	m.editor = textarea.New()
	m.editor.ShowLineNumbers = true
	// Word jumps are the vim-style ctrl+left/right the editor promises in
	// its legend (the alt+/alt+f defaults stay available).
	m.editor.KeyMap.WordBackward.SetKeys("ctrl+left", "alt+left", "alt+b")
	m.editor.KeyMap.WordForward.SetKeys("ctrl+right", "alt+right", "alt+f")
	m.editor.Blur()

	// A single-file resource has nothing to browse — os.ReadDir(root) on a
	// regular file fails outright ("open root: not a directory"), which is
	// exactly the bug this branch fixes. Lstat, not Stat: a symlinked root
	// is refused by enterEditor below, never silently followed.
	info, statErr := os.Lstat(root)
	if statErr == nil && info.Mode().IsRegular() {
		m.singleFile = true
		m.root = &fileNode{path: root, name: name, loaded: true}
	} else {
		m.root = &fileNode{path: root, name: name, isDir: true, open: true}
		m.status = "loading " + root
		if err := m.loadChildren(m.root); err != nil {
			m.status = "error: " + err.Error()
		}
	}
	m.rebuild()
	if m.singleFile {
		// Nothing to navigate to first: open straight into the editor,
		// exactly like pressing enter on the single row would.
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
		prevMode := m.mode
		cmd, quit := m.handleKey(msg)
		if quit {
			return m, tea.Quit
		}
		if m.mode != prevMode {
			// A mode change alters the pane split (the editor is a sidebar
			// vs. half vs. focus) — re-fit before the next render.
			m.refit()
		}
		return m, cmd
	}
	return m, nil
}

// treeSidebarWidth is the tree pane width while editing a file: the editor
// is the focus, so the tree shrinks to a context sidebar and the editor
// takes the rest of the terminal.
const treeSidebarWidth = 32

// paneLayout returns the left (tree) and right (info/editor) column widths
// for the current mode, from the width the terminal actually gives the
// content (after the app margin) minus one column for the separator. In
// modeEdit the tree is a sidebar; every other mode splits evenly.
func (m *fileEditModel) paneLayout() (left, right int) {
	usable := max(m.width-appStyle.GetHorizontalMargins(), 4)
	if m.mode == modeEdit {
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

// refit recomputes the pane geometry (mode-dependent) and re-sizes the
// editor to it. Call it whenever the terminal size or the mode changes, so
// entering/leaving the editor widens/narrows it immediately.
func (m *fileEditModel) refit() {
	m.leftW, m.rightW = m.paneLayout()
	m.editor.SetWidth(max(m.rightW, 10))
	m.editor.SetHeight(max(m.paneHeight()-1, 3))
	m.ensureCursorVisible()
}

func (m *fileEditModel) handleKey(msg tea.KeyMsg) (tea.Cmd, bool) {
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

// navAction handles the non-movement nav keys; it reports whether the
// editor must quit.
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
	case "q", keyEsc, "ctrl+c":
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
	switch msg.String() {
	case "ctrl+s":
		if err := m.save(); err != nil {
			m.status = err.Error()
		}
		m.editor.Blur()
		m.mode = modeNav
	case keyEsc:
		m.editor.Blur()
		m.mode = modeNav
	default:
		var cmd tea.Cmd
		m.editor, cmd = m.editor.Update(msg)
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
		rows := max(m.paneHeight()-2, 1)
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
		m.confirmDelete()
	case "n", keyEsc, "q":
		m.mode = modeNav
	}
	return nil, false
}

// loadChildren fills a directory node with its (sorted, lazily loaded)
// children. Directory entries come first, then files, each ordered by name
// without respect to case.
func (m *fileEditModel) loadChildren(n *fileNode) error {
	if !n.isDir || n.loaded {
		return nil
	}
	entries, err := os.ReadDir(n.path)
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

// sortNodes orders directory nodes before file nodes, then by name without
// regard to case.
func sortNodes(nodes []*fileNode) {
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].isDir != nodes[j].isDir {
			return nodes[i].isDir
		}
		return strings.ToLower(nodes[i].name) < strings.ToLower(nodes[j].name)
	})
}

// rebuild flattens the visible tree (the root and every loaded, open
// directory descendant) into the row slice, keeping the cursor on the same
// node when possible.
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
	rows := m.paneHeight()
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

// enterEditor opens the selected file in the focused textarea. Re-entering
// the editor for the same, unsaved file keeps the buffer; joining a
// different file reloads from disk. Symlinks are never followed.
func (m *fileEditModel) enterEditor() error {
	n := m.selected()
	if n == nil || n.isDir {
		return nil
	}
	info, err := os.Lstat(n.path)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		m.status = n.path + " is a symlink — editing it is refused"
		return nil
	}
	if m.editPath == n.path && m.dirty {
		m.mode = modeEdit
		m.editor.Focus()
		m.status = "editing " + n.path
		return nil
	}
	if err := m.previewFile(n.path); err != nil {
		return err
	}
	if !m.previewOK {
		m.status = "cannot edit " + n.path
		return nil
	}
	m.mode = modeEdit
	m.editor.Focus()
	m.status = "editing " + n.path
	return nil
}

// previewFile loads path's content into the editor (blurred). Binary files
// and files over maxEditableBytes reset the preview and report a status
// instead.
func (m *fileEditModel) previewFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.Size() > maxEditableBytes {
		m.previewOK = false
		m.editor.SetValue("")
		m.status = fmt.Sprintf("%s is %.1f MiB — only files up to 2 MiB can be edited", path, float64(info.Size())/(1<<20))
		return nil
	}
	data, err := os.ReadFile(path)
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
	m.editor.SetValue(m.original)
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
	// A single-file resource's parent directory is itself guarded against
	// anything created or renamed beside it (the same "rename-over-
	// watchroot" defense CLAUDE.md calls out) — writeFileKeepMeta's usual
	// temp-sibling + rename would be denied there. writeFileInPlace rewrites
	// the file's existing inode directly, same as every other in-place
	// unlock/lock in this codebase.
	writer := writeFileKeepMeta
	if m.singleFile {
		writer = writeFileInPlace
	}
	if err := writer(m.editPath, []byte(m.editor.Value())); err != nil {
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

// beginCreate switches to the input mode used for the new file/directory
// name. The new entry is created inside the currently selected directory
// (or the parent of a selected file).
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
		err = os.WriteFile(target, nil, 0o600)
	case createDir:
		err = os.Mkdir(target, 0o750)
	}
	if err != nil {
		m.status = err.Error()
		return
	}
	if err := applyNewFileMeta(target, m.inputKind == createDir); err != nil {
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

// beginChmod prepares the octal mode editor for the selected entry. The
// value is a conventional unix mode: 3-4 octal digits where the optional
// leading digit carries setuid/setgid/sticky (e.g. 4755).
func (m *fileEditModel) beginChmod() {
	n := m.selected()
	if n == nil {
		return
	}
	info, err := os.Lstat(n.path)
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
	m.chmodInput.Width = max(m.paneWidth()-12, 4)
	m.chmodInput.Prompt = "mode: "
	m.chmodInput.Placeholder = "0644"
	// Prefill the full mode including setuid/setgid/sticky, so confirming
	// an untouched input can never silently strip those bits.
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
	// The typed value is a unix mode; os.FileMode must be derived through
	// the converter or the setuid/setgid/sticky bits are lost (a raw cast
	// maps them onto undefined FileMode bits that chmod ignores).
	if err := os.Chmod(m.chmodPath, unixModeToFileMode(uint32(perm))); err != nil {
		m.status = err.Error()
		return
	}
	m.status = "chmod " + m.chmodPath + " = " + v
}

// beginChown opens the owner picker listing every login user plus a
// "keep current" entry.
func (m *fileEditModel) beginChown() {
	n := m.selected()
	if n == nil {
		return
	}
	info, err := os.Lstat(n.path)
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
	if err := os.Lchown(m.chownPath, u.uid, u.gid); err != nil {
		m.status = err.Error()
		m.mode = modeNav
		return
	}
	m.status = "chown " + m.chownPath + " → " + u.label
	m.mode = modeNav
}

// beginDelete asks for confirmation (default no) before removing the
// selected entry — recursive for directories.
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
	if err := os.RemoveAll(path); err != nil {
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
	// 2 chrome rows (the header + the top margin) plus however many rows the
	// legend wraps to at this width.
	return max(m.height-2-m.legendRows(), 2)
}

// legendRows is the number of terminal rows the (width-wrapped) legend
// occupies for the current mode.
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

// body lays the tree and the info/editor pane side by side, each fixed to
// its column width so the layout fills the terminal.
func (m *fileEditModel) body() string {
	return lipgloss.JoinHorizontal(lipgloss.Top,
		m.renderTree(m.leftW, m.paneHeight()),
		paneSepStyle.Render("│"),
		// Fix the right column to rightW so the pane reaches the terminal
		// edge (renderRight clips its own content to rightW, so Width only
		// pads; MaxWidth is the belt against a stray wide line).
		lipgloss.NewStyle().Width(m.rightW).MaxWidth(m.rightW).
			Render(m.renderRight(m.rightW, m.paneHeight())),
	)
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
		// Pad every row to the full pane width so the block lipgloss lays out
		// is exactly `width` wide — otherwise JoinHorizontal sizes the tree
		// column to its longest line and the right pane never reaches the
		// terminal edge. Padding first also makes the selection bar span the
		// whole pane.
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
	case modeEdit:
		return m.editor.View()
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

// renderEntryInfo draws the right pane for the selected node: a directory
// listing or the file metadata.
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
		for _, e := range listEntries(n, height-2) {
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
	info, err := os.Lstat(n.path)
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
	// Wrap to the terminal so the nav-mode hint list (which is wider than a
	// typical terminal) never spills past the right edge; paneHeight accounts
	// for the extra rows.
	usable := max(m.width-appStyle.GetHorizontalMargins(), 20)
	return footerStyle.Width(usable).Render(strings.Join(parts, "  ·  "))
}

var legendHints = map[fileEditMode][]hint{
	modeNav: {
		{"↑/↓ j/k", "move"},
		{"→/↵ l", "open"},
		{"← h", "back"},
		{"g/G", "top/bottom"},
		{"a/A", "new file/dir"},
		{"d", "delete"},
		{"m", "chmod"},
		{"c", "chown"},
		{"ctrl+s", "save"},
		{"q/esc", "quit"},
	},
	modeEdit: {
		{"ctrl+s", "save & close"},
		{"esc", "close"},
		{"ctrl+←/→", "word jump"},
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
