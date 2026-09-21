package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/textinput"
)

// Security regressions for the edit-protected editor's write paths. The editor runs as root over a
// tree owned by an unprivileged user, and the offline flow unlocks the vault with the daemon
// STOPPED (no guard attached), so that user can pre-plant entries inside their own directory. Every
// write must therefore refuse to traverse a symlink instead of validating a path and re-opening it.

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func mustRead(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// writeFileKeepMeta Lstat-checks `path` but creates its temp file at the predictable sibling
// path+".app_listener.edit" with O_CREATE|O_TRUNC and no O_EXCL/O_NOFOLLOW. A symlink planted at
// that name is followed, so root truncates and rewrites the link target — and then chmods/chowns
// it to the edited file's owner.
func TestWriteFileKeepMetaRefusesPlantedTempSibling(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "id_rsa")
	victim := filepath.Join(root, "victim")
	mustWrite(t, target, "original key")
	mustWrite(t, victim, "VICTIM")

	if err := os.Symlink(victim, target+".app_listener.edit"); err != nil {
		t.Fatal(err)
	}

	err := writeFileKeepMeta(target, []byte("edited by root"))

	// Security property: the symlink target must be untouched. The write may still succeed by
	// replacing the real target atomically (the temp name is our namespace, cleared with O_EXCL),
	// but it must never follow the planted link.
	if got := mustRead(t, victim); got != "VICTIM" {
		t.Fatalf("root wrote through the planted temp-sibling symlink: victim = %q (err %v)", got, err)
	}
	if err == nil {
		if got := mustRead(t, target); got != "edited by root" {
			t.Fatalf("write reported success but the real target was not updated: %q", got)
		}
	}
}

// createEntry uses os.WriteFile + os.Chmod (applyNewFileMeta), both of which follow an existing
// symlink. Creating an entry whose name is already a planted symlink truncates the link target.
func TestCreateEntryRefusesPlantedSymlink(t *testing.T) {
	root := t.TempDir()
	victim := filepath.Join(root, "victim")
	mustWrite(t, victim, "VICTIM")
	if err := os.Symlink(victim, filepath.Join(root, "notes.txt")); err != nil {
		t.Fatal(err)
	}

	m := &fileEditModel{}
	m.input = textinput.New()
	m.input.SetValue("notes.txt")
	m.inputDir = root
	m.inputKind = createFile
	m.createEntry()

	if got := mustRead(t, victim); got != "VICTIM" {
		t.Fatalf("root truncated the target of a planted symlink: victim = %q", got)
	}
	if strings.HasPrefix(m.status, "created ") {
		t.Fatalf("createEntry created an entry through a planted symlink; expected a refusal, status = %q", m.status)
	}
}
