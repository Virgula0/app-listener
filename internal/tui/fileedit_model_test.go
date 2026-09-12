package tui

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// testTree builds the temp vault tree used by most tests:
//
//	<root>/bin/tool.sh
//	<root>/alpha.txt  ("hello alpha")
//	<root>/zebra.txt  ("zebra")
func testTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "bin"), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "bin", "tool.sh"), []byte("#!/bin/sh\n"), 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "alpha.txt"), []byte("hello alpha"), 0o640); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "zebra.txt"), []byte("zebra"), 0o640); err != nil {
		t.Fatal(err)
	}
	return root
}

func updateKey(m *fileEditModel, key string) tea.Cmd {
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)})
	return cmd
}

func TestIsBinaryContent(t *testing.T) {
	if isBinaryContent([]byte("plain text\n")) {
		t.Error("plain text classified as binary")
	}
	if !isBinaryContent([]byte{'a', 0, 'b'}) {
		t.Error("NUL byte not classified as binary")
	}
	if !isBinaryContent([]byte{0}) {
		t.Error("leading NUL not classified as binary")
	}
	// Only the first binaryScanBytes window is scanned: a NUL beyond it
	// is not seen, so the file is treated as text.
	late := make([]byte, binaryScanBytes+16) // all-NUL beyond the window
	copy(late, strings.Repeat("x", binaryScanBytes))
	late[binaryScanBytes] = 0
	if isBinaryContent(late) {
		t.Error("NUL beyond the scan window classified the file as binary")
	}
}

func TestWriteFileKeepMeta(t *testing.T) {
	path := filepath.Join(t.TempDir(), "secret.cfg")
	if err := os.WriteFile(path, []byte("old"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := writeFileKeepMeta(path, []byte("new-content")); err != nil {
		t.Fatalf("writeFileKeepMeta: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "new-content" {
		t.Errorf("content = %q, want new-content", data)
	}
	info, err := os.Lstat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("mode = %v, want 0600 (the original, not the temp default)", info.Mode().Perm())
	}
}

func TestWriteFileKeepMetaRefusesSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	if err := os.WriteFile(target, []byte("t"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "link")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	if err := writeFileKeepMeta(link, []byte("hijack")); err == nil {
		t.Fatal("writing through a symlink was not refused")
	}
}

func TestWriteFileInPlace(t *testing.T) {
	path := filepath.Join(t.TempDir(), "registry.vdf")
	if err := os.WriteFile(path, []byte("old content"), 0o600); err != nil {
		t.Fatal(err)
	}
	before := inodeOf(t, path)

	if err := writeFileInPlace(path, []byte("new")); err != nil {
		t.Fatalf("writeFileInPlace: %v", err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "new" {
		t.Errorf("content = %q, want %q", data, "new")
	}
	if got := inodeOf(t, path); got != before {
		t.Errorf("writeFileInPlace changed the inode: %d -> %d", before, got)
	}
}

func TestWriteFileInPlaceRefusesSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	if err := os.WriteFile(target, []byte("t"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "link")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	if err := writeFileInPlace(link, []byte("hijack")); err == nil {
		t.Fatal("writing through a symlink was not refused")
	}
}

func inodeOf(t *testing.T, path string) uint64 {
	t.Helper()
	info, err := os.Lstat(path)
	if err != nil {
		t.Fatalf("lstat %s: %v", path, err)
	}
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		t.Fatal("Stat_t unavailable on this platform")
	}
	return st.Ino
}

// TestFileEditModelSingleFileTarget is the regression test for the bug where
// RunFileEditor/newFileEditModel unconditionally treated root as a
// directory: os.ReadDir on a single guarded regular file (a file-vault
// resource, e.g. registry.vdf) failed outright with "open root: not a
// directory" before the editor ever opened. A single-file root must skip
// the tree entirely, open straight into the editor, and save in place —
// never through the tree editor's usual temp-file + rename (which a
// single-file watch root's guard denies, since it also protects its own
// parent directory against anything created or renamed beside it).
func TestFileEditModelSingleFileTarget(t *testing.T) {
	path := filepath.Join(t.TempDir(), "registry.vdf")
	if err := os.WriteFile(path, []byte("original content"), 0o600); err != nil {
		t.Fatal(err)
	}
	before := inodeOf(t, path)

	m := newFileEditModel(path)
	if !m.singleFile {
		t.Fatal("singleFile must be true for a regular-file root")
	}
	if m.mode != modeEdit {
		t.Fatalf("mode = %v, want modeEdit (opened straight into the editor)", m.mode)
	}
	if got := m.editor.Value(); got != "original content" {
		t.Errorf("editor = %q, want the file content", got)
	}
	if len(m.rows) != 1 || m.rows[0] != m.root {
		t.Fatalf("rows = %v, want exactly [root] (nothing to browse)", m.rows)
	}

	m.editor.SetValue("edited content")
	updateKey(m, "ctrl+s")
	if m.mode != modeNav {
		t.Fatalf("after save, mode = %v, want modeNav", m.mode)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "edited content" {
		t.Errorf("file on disk = %q, want edited content", data)
	}
	if got := inodeOf(t, path); got != before {
		t.Errorf("save changed the file's inode (must be in-place, never rename): %d -> %d", before, got)
	}
}

// TestFileEditModelSingleFileCreateRefused: "a"/"A" would otherwise
// dereference the (nil) parent of a single-file root's node.
func TestFileEditModelSingleFileCreateRefused(t *testing.T) {
	path := filepath.Join(t.TempDir(), "registry.vdf")
	if err := os.WriteFile(path, []byte("content"), 0o600); err != nil {
		t.Fatal(err)
	}
	m := newFileEditModel(path)
	// enterEditor() at construction leaves the model in modeEdit; back out
	// to modeNav first, same as pressing esc, to reach the create keys.
	m.mode = modeNav
	m.editor.Blur()

	m.beginCreate(createFile)
	if m.mode == modeInput {
		t.Fatal("beginCreate must refuse to enter input mode for a single-file resource")
	}
	if !strings.Contains(m.status, "single-file") {
		t.Errorf("status = %q, want a message about single-file resources", m.status)
	}
}

func TestFileEditNavigateAndEdit(t *testing.T) {
	root := testTree(t)
	m := newFileEditModel(root)
	m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	if got := m.selected(); got != nil && got.path != root {
		t.Fatalf("initial selection = %s, want root %s", got.path, root)
	}

	// root -> bin (directories sort first)
	updateKey(m, "j")
	if got := m.selected(); got == nil || got.name != "bin" {
		t.Fatalf("after j selection = %v, want bin", got)
	}

	// enter bin, cursor returns to it, then j twice -> tool.sh -> alpha.txt
	updateKey(m, "l")
	if got := m.selected(); got == nil || got.name != "bin" {
		t.Fatalf("after expanding, selection = %v, want bin", got)
	}
	updateKey(m, "j")
	if got := m.selected(); got == nil || got.name != "tool.sh" {
		t.Fatalf("after j, selection = %v, want tool.sh", got)
	}
	updateKey(m, "j")
	if got := m.selected(); got == nil || got.name != "alpha.txt" {
		t.Fatalf("after j, selection = %v, want alpha.txt", got)
	}

	// open the file for editing
	updateKey(m, "e")
	if m.mode != modeEdit {
		t.Fatalf("mode = %v, want modeEdit", m.mode)
	}
	if got := m.editor.Value(); got != "hello alpha" {
		t.Errorf("editor = %q, want the file content", got)
	}

	// change it, save, return to nav
	m.editor.SetValue("hello edited")
	updateKey(m, "ctrl+s")
	if m.mode != modeNav {
		t.Fatalf("after save, mode = %v, want modeNav", m.mode)
	}
	data, err := os.ReadFile(filepath.Join(root, "alpha.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello edited" {
		t.Errorf("file on disk = %q, want hello edited", data)
	}
}

func TestFileEditPaneGeometry(t *testing.T) {
	root := testTree(t)
	m := newFileEditModel(root)
	const w = 120
	m.Update(tea.WindowSizeMsg{Width: w, Height: 30})

	usable := w - appStyle.GetHorizontalMargins()

	// Nav mode: even split, and the two panes + separator exactly fill the
	// usable width.
	if m.leftW+1+m.rightW != usable {
		t.Fatalf("nav: leftW(%d)+1+rightW(%d) = %d, want usable %d", m.leftW, m.rightW, m.leftW+1+m.rightW, usable)
	}
	if abs(m.leftW-m.rightW) > 1 {
		t.Fatalf("nav: panes should be roughly even, got left=%d right=%d", m.leftW, m.rightW)
	}

	// The rendered tree block is padded to exactly leftW so the layout does
	// not collapse to the longest filename.
	for _, line := range strings.Split(strings.TrimRight(m.renderTree(m.leftW, m.paneHeight()), "\n"), "\n") {
		if lw := lipgloss.Width(line); lw != m.leftW {
			t.Fatalf("tree line width = %d, want leftW %d (%q)", lw, m.leftW, line)
		}
	}

	// The two panes + separator fill the usable width (the bug: the body used
	// to collapse to the tree's longest line + the editor, leaving dead space).
	if got := maxLineWidth(m.body()); got != usable {
		t.Fatalf("nav: body width = %d, want usable %d", got, usable)
	}

	// And nothing (including the long nav-mode legend) spills past the
	// terminal's right edge.
	if got := maxLineWidth(m.View()); got > w {
		t.Fatalf("nav: View width = %d, overflows terminal %d", got, w)
	}
	// The legend wraps rather than overflowing, so the editor still gets a
	// sensible height.
	if m.paneHeight() < 10 {
		t.Fatalf("nav: paneHeight collapsed to %d after the legend wrap", m.paneHeight())
	}

	// Edit mode: the tree becomes a sidebar and the editor takes most of the
	// width.
	navRight := m.rightW
	updateKey(m, "j") // -> bin
	updateKey(m, "l") // expand
	updateKey(m, "j") // tool.sh
	updateKey(m, "j") // alpha.txt
	updateKey(m, "e") // edit
	if m.mode != modeEdit {
		t.Fatalf("mode = %v, want modeEdit", m.mode)
	}
	if m.leftW > treeSidebarWidth {
		t.Fatalf("edit: tree pane should shrink to a sidebar, got leftW=%d", m.leftW)
	}
	if m.rightW <= navRight {
		t.Fatalf("edit: editor pane should widen (was %d, now %d)", navRight, m.rightW)
	}
	if m.leftW+1+m.rightW != usable {
		t.Fatalf("edit: panes+sep = %d, want usable %d", m.leftW+1+m.rightW, usable)
	}
	if got := maxLineWidth(m.body()); got != usable {
		t.Fatalf("edit: body width = %d, want usable %d", got, usable)
	}

	// Leaving the editor restores the even split.
	updateKey(m, "esc")
	if m.rightW != navRight {
		t.Fatalf("after esc: rightW = %d, want the nav split %d", m.rightW, navRight)
	}
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func maxLineWidth(view string) int {
	w := 0
	for _, line := range strings.Split(view, "\n") {
		if lw := lipgloss.Width(line); lw > w {
			w = lw
		}
	}
	return w
}

func TestFileTreeGotoEdges(t *testing.T) {
	root := testTree(t)
	m := newFileEditModel(root)
	m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	updateKey(m, "G")
	if m.cur != len(m.rows)-1 {
		t.Errorf("G cursor = %d, want last row %d", m.cur, len(m.rows)-1)
	}
	updateKey(m, "g")
	if m.cur != 0 {
		t.Errorf("g cursor = %d, want 0", m.cur)
	}
	// Left on a closed directory is a no-op and stays safe.
	updateKey(m, "h")
	if m.cur != 0 {
		t.Errorf("h cursor = %d, want 0", m.cur)
	}
}

func TestFileCreateAndDeleteFlow(t *testing.T) {
	root := testTree(t)
	m := newFileEditModel(root)
	m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// select alpha.txt then create a sibling file next to it
	updateKey(m, "j") // bin
	updateKey(m, "j") // alpha.txt
	updateKey(m, "a")
	if m.mode != modeInput {
		t.Fatalf("mode = %v, want modeInput", m.mode)
	}
	m.input.SetValue("notes.md")
	updateKey(m, "enter")
	if m.mode != modeNav {
		t.Fatalf("after create, mode = %v, want modeNav", m.mode)
	}
	if got := m.selected(); got == nil || got.name != "notes.md" {
		t.Fatalf("cursor after create = %v, want notes.md", got)
	}
	if _, err := os.Stat(filepath.Join(root, "notes.md")); err != nil {
		t.Fatalf("created file missing: %v", err)
	}

	// delete the created file via the default-no confirm
	updateKey(m, "d")
	if m.mode != modeConfirm {
		t.Fatalf("mode = %v, want modeConfirm", m.mode)
	}
	updateKey(m, "y")
	if m.mode != modeNav {
		t.Fatalf("after confirm, mode = %v, want modeNav", m.mode)
	}
	if _, err := os.Stat(filepath.Join(root, "notes.md")); !os.IsNotExist(err) {
		t.Errorf("deleted file still exists (err=%v)", err)
	}
}

func TestQuitDirtyRequiresConfirm(t *testing.T) {
	root := testTree(t)
	m := newFileEditModel(root)
	m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	// open alpha.txt and dirty it by typing
	updateKey(m, "j") // bin
	updateKey(m, "j") // tool.sh
	updateKey(m, "j") // alpha.txt
	updateKey(m, "e")
	updateKey(m, "z") // typing diverges the buffer
	if !m.dirty {
		t.Fatal("buffer is not marked dirty after typing")
	}
	updateKey(m, "esc") // back to nav, buffer kept
	if m.mode != modeNav {
		t.Fatalf("esc should return to nav, mode = %v", m.mode)
	}

	// q with a dirty buffer must ask, not quit
	if cmd := updateKey(m, "q"); cmd != nil {
		t.Fatal("q quit without confirmation")
	}
	if m.mode != modeConfirm {
		t.Fatalf("mode = %v, want modeConfirm for dirty quit", m.mode)
	}

	// n cancels the quit
	updateKey(m, "n")
	if m.mode != modeNav {
		t.Fatalf("n should cancel, mode = %v", m.mode)
	}

	// q again, then y confirms the quit (returns tea.Quit)
	updateKey(m, "q")
	cmd := updateKey(m, "y")
	if cmd == nil {
		t.Fatal("y on confirm-quit did not quit")
	}
	data, err := os.ReadFile(filepath.Join(root, "alpha.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello alpha" {
		t.Errorf("file changed on disk after discard: %q", data)
	}
}

// TestIsOctalMode covers the chmod validator.
func TestIsOctalMode(t *testing.T) {
	valid := []string{"0", "7", "644", "0644", "0777", "0000", "7000"}
	for _, v := range valid {
		if !isOctalMode(v) {
			t.Errorf("isOctalMode(%q) = false, want true", v)
		}
	}
	invalid := []string{"", "8", "9", "06444", "0x1ff", "644 ", " -1", "abcd", "064A"}
	for _, v := range invalid {
		if isOctalMode(v) {
			t.Errorf("isOctalMode(%q) = true, want false", v)
		}
	}
}

// singleSelection focuses the model on the node named name (walking the
// rows in order).
func selectNamed(t *testing.T, m *fileEditModel, name string) {
	t.Helper()
	for i, row := range m.rows {
		if row.name == name {
			m.cur = i
			m.ensureCursorVisible()
			return
		}
	}
	t.Fatalf("node %q not found in tree", name)
}

// TestFileChmod exercises the modeChmod flow end to end on a regular file.
func TestFileChmod(t *testing.T) {
	root := testTree(t)
	m := newFileEditModel(root)
	m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	selectNamed(t, m, "alpha.txt")
	if err := os.Chmod(filepath.Join(root, "alpha.txt"), 0o600); err != nil {
		t.Fatal(err)
	}
	updateKey(m, "m")
	if m.mode != modeChmod {
		t.Fatalf("mode = %v, want modeChmod", m.mode)
	}
	if got := m.chmodInput.Value(); got != "0600" {
		t.Errorf("prefilled mode = %q, want 0600", got)
	}

	m.chmodInput.SetValue("0640")
	updateKey(m, "enter")
	if m.mode != modeNav {
		t.Fatalf("mode after apply = %v, want modeNav", m.mode)
	}
	info, err := os.Lstat(filepath.Join(root, "alpha.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o640 {
		t.Errorf("mode after chmod = %v, want 0640", info.Mode().Perm())
	}

	// Invalid input cancels without touching the file.
	updateKey(m, "m")
	m.chmodInput.SetValue("99")
	updateKey(m, "enter")
	if m.mode != modeNav {
		t.Fatalf("mode after invalid chmod = %v, want modeNav", m.mode)
	}
	if !strings.Contains(m.status, "invalid") {
		t.Errorf("invalid mode was not rejected (status = %q)", m.status)
	}
	info, err = os.Lstat(filepath.Join(root, "alpha.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o640 {
		t.Errorf("file mutated by invalid chmod: %o", info.Mode().Perm())
	}
}

// TestFileChownFlow verifies the owner picker lists logical users and that
// the "keep current" entry leaves the file alone.
func TestFileChownFlow(t *testing.T) {
	root := testTree(t)
	m := newFileEditModel(root)
	m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	selectNamed(t, m, "alpha.txt")
	updateKey(m, "c")
	if m.mode != modeChown {
		t.Fatalf("mode = %v, want modeChown", m.mode)
	}
	if len(m.chownUsers) == 0 {
		t.Fatal("chown picker has no entries")
	}
	if m.chownUsers[0].uid != -1 {
		t.Errorf("first chown entry must be keep-current (uid -1), got %d", m.chownUsers[0].uid)
	}

	info, err := os.Lstat(filepath.Join(root, "alpha.txt"))
	if err != nil {
		t.Fatal(err)
	}
	sysBefore := info.Sys().(*syscall.Stat_t)

	// Press enter on "keep current": ownership reported unchanged, mode back.
	if cmd := updateKey(m, "enter"); cmd != nil {
		t.Fatal("unexpected ade from keep-current chown")
	}
	if m.mode != modeNav {
		t.Fatalf("mode after keep-current = %v, want modeNav", m.mode)
	}
	info, err = os.Lstat(filepath.Join(root, "alpha.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if sys := info.Sys().(*syscall.Stat_t); sys.Uid != sysBefore.Uid || sys.Gid != sysBefore.Gid {
		t.Errorf("keep-current chown changed ownership: %d:%d -> %d:%d",
			sysBefore.Uid, sysBefore.Gid, sys.Uid, sys.Gid)
	}
}

// TestDefaultNewFileMode models the conventional 0644/0755 defaults.
func TestDefaultNewFileMode(t *testing.T) {
	if got := defaultNewFileMode(false); got != 0o644 {
		t.Errorf("file default = %o, want 0644", got)
	}
	if got := defaultNewFileMode(true); got != 0o755 {
		t.Errorf("dir default = %o, want 0755", got)
	}
}

// TestUnixModeToFileMode covers the conversion from conventional unix mode
// numbers (with special bits) to Go's os.FileMode flags. A raw cast is not
// equivalent: it maps setuid/setgid/sticky onto undefined FileMode bits
// that chmod(2) ignores, silently reducing 4755 to 0755.
func TestUnixModeToFileMode(t *testing.T) {
	cases := []struct {
		in   uint32
		want os.FileMode
	}{
		{0o644, 0o644},
		{0o755, 0o755},
		{0o0000, 0},
		{0o4000, os.ModeSetuid},
		{0o2000, os.ModeSetgid},
		{0o1000, os.ModeSticky},
		{0o4755, os.ModeSetuid | 0o755},
		{0o2755, os.ModeSetgid | 0o755},
		{0o1777, os.ModeSticky | 0o777},
		{0o7777, os.ModeSetuid | os.ModeSetgid | os.ModeSticky | 0o777},
	}
	for _, tc := range cases {
		if got := unixModeToFileMode(tc.in); got != tc.want {
			t.Errorf("unixModeToFileMode(%04o) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

// TestFileModeToUnixModeRoundTrip verifies both converters invert each
// other for every combination that fits in four octal digits.
func TestFileModeToUnixModeRoundTrip(t *testing.T) {
	for _, mode := range []uint32{0o0000, 0o0644, 0o755, 0o1777, 0o2755, 0o4755, 0o6755, 0o7755, 0o7777} {
		if got := fileModeToUnixMode(unixModeToFileMode(mode)); got != mode {
			t.Errorf("round trip of %04o = %04o", mode, got)
		}
		if back := unixModeToFileMode(fileModeToUnixMode(unixModeToFileMode(mode))); back != unixModeToFileMode(mode) {
			t.Errorf("inverse round trip diverged for %04o", mode)
		}
	}
}

// TestChmodSuidEndToEnd reproduces the reported bug through the editor
// flow: chmod 4755 must leave a real suid bit on disk, and reopening the
// chmod dialog must prefill 4755 instead of degrading it to 0755 (which
// would strip suid on a bare confirm).
func TestChmodSuidEndToEnd(t *testing.T) {
	root := testTree(t)
	m := newFileEditModel(root)
	m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	selectNamed(t, m, "alpha.txt")
	updateKey(m, "m")
	m.chmodInput.SetValue("4755")
	updateKey(m, "enter")

	info, err := os.Lstat(filepath.Join(root, "alpha.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode()&os.ModeSetuid == 0 || info.Mode().Perm() != 0o755 {
		t.Errorf("chmod 4755 produced %v, want setuid + 0755", info.Mode())
	}

	// Reopening prefills the full mode including the special bits.
	updateKey(m, "m")
	if got := m.chmodInput.Value(); got != "4755" {
		t.Errorf("prefilled mode = %q, want 4755", got)
	}
	// And confirming the untouched input keeps the file exactly as is.
	updateKey(m, "enter")
	info, err = os.Lstat(filepath.Join(root, "alpha.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode()&os.ModeSetuid == 0 || info.Mode().Perm() != 0o755 {
		t.Errorf("confirming untouched input changed mode to %v", info.Mode())
	}
}
