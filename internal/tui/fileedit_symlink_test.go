package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/textarea"
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

	err := writeFileKeepMeta(openRootT(t, root), "id_rsa", []byte("edited by root"))

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

	m := newFileEditModel(root)
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

// swappedParent lists root/sub/inner/<name> as the editor would, then swaps sub for a symlink to
// an outside dir holding inner/<name>: O_NOFOLLOW covers only the final component, so the editor's
// stored path now walks out of the vault. Returns the model, stored path and victim.
func swappedParent(t *testing.T, name string) (*fileEditModel, string, string) {
	t.Helper()
	root := t.TempDir()
	outside := t.TempDir()
	sub := filepath.Join(root, "sub")
	for _, d := range []string{filepath.Join(sub, "inner"), filepath.Join(outside, "inner")} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	stored := filepath.Join(sub, "inner", name)
	mustWrite(t, stored, "vault content")
	m := newFileEditModel(root)

	victim := filepath.Join(outside, "inner", name)
	mustWrite(t, victim, "VICTIM")
	if err := os.Rename(sub, filepath.Join(root, "sub.aside")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, sub); err != nil {
		t.Fatal(err)
	}
	return m, stored, victim
}

func TestSaveRefusesSymlinkedParent(t *testing.T) {
	m, stored, victim := swappedParent(t, "config")
	m.editPath = stored
	m.editor = textarea.New()
	m.editor.SetValue("edited by root")
	err := m.save()
	if got := mustRead(t, victim); got != "VICTIM" {
		t.Fatalf("save wrote outside the vault through a swapped parent dir: victim = %q (err %v)", got, err)
	}
}

func TestChmodRefusesSymlinkedParent(t *testing.T) {
	m, stored, victim := swappedParent(t, "config")
	if err := os.Chmod(victim, 0o600); err != nil {
		t.Fatal(err)
	}
	m.chmodPath = stored
	m.chmodInput = textinput.New()
	m.chmodInput.SetValue("0666")
	m.applyChmod()
	info, err := os.Lstat(victim)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("chmod reached outside the vault through a swapped parent dir: victim mode %o", info.Mode().Perm())
	}
}

func TestDeleteRefusesSymlinkedParent(t *testing.T) {
	m, stored, victim := swappedParent(t, "config")
	m.pendingDelete = stored
	m.confirmDelete()
	if _, err := os.Lstat(victim); err != nil {
		t.Fatalf("delete removed a file outside the vault through a swapped parent dir: %v", err)
	}
}

func TestCreateEntryRefusesSymlinkedParent(t *testing.T) {
	m, stored, victim := swappedParent(t, "config")
	m.input = textinput.New()
	m.input.SetValue("planted")
	m.inputDir = filepath.Dir(stored)
	m.inputKind = createFile
	m.createEntry()
	if _, err := os.Lstat(filepath.Join(filepath.Dir(victim), "planted")); err == nil {
		t.Fatalf("createEntry created a file outside the vault through a swapped parent dir (status %q)", m.status)
	}
}

// A swap ABOVE the vault root (the user owns its ancestors) must not redirect the session either:
// the save lands in the real, moved vault.
func TestSaveFollowsPinnedVaultNotSwappedAncestor(t *testing.T) {
	base := t.TempDir()
	outside := t.TempDir()
	vault := filepath.Join(base, "cfg", "vault")
	for _, d := range []string{vault, filepath.Join(outside, "vault")} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	mustWrite(t, filepath.Join(vault, "config"), "vault content")
	victim := filepath.Join(outside, "vault", "config")
	mustWrite(t, victim, "VICTIM")
	m := newFileEditModel(vault)
	defer m.vault.close()

	if err := os.Rename(filepath.Join(base, "cfg"), filepath.Join(base, "cfg.aside")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(base, "cfg")); err != nil {
		t.Fatal(err)
	}
	m.editPath = filepath.Join(vault, "config")
	m.editor = textarea.New()
	m.editor.SetValue("edited by root")
	err := m.save()
	if got := mustRead(t, victim); got != "VICTIM" {
		t.Fatalf("save followed a swapped vault ancestor: victim = %q (err %v)", got, err)
	}
	if got := mustRead(t, filepath.Join(base, "cfg.aside", "vault", "config")); got != "edited by root" {
		t.Fatalf("save did not reach the pinned vault: %q (err %v)", got, err)
	}
}
