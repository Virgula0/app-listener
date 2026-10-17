//go:build gui

package gui

import (
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func TestEncodeTable(t *testing.T) {
	headers := []string{"Process", "Path"}
	rows := [][]string{{"ssh", "./id,ed25519"}, {"cat", `say "hi"`}}
	csvOut, err := encodeTable("CSV", headers, rows)
	if err != nil {
		t.Fatal(err)
	}
	if want := "Process,Path\nssh,\"./id,ed25519\"\ncat,\"say \"\"hi\"\"\"\n"; string(csvOut) != want {
		t.Fatalf("CSV = %q, want %q", csvOut, want)
	}
	jsonOut, err := encodeTable("JSON", headers, rows[:1])
	if err != nil {
		t.Fatal(err)
	}
	if want := "[\n  {\n    \"Path\": \"./id,ed25519\",\n    \"Process\": \"ssh\"\n  }\n]\n"; string(jsonOut) != want {
		t.Fatalf("JSON = %q, want %q", jsonOut, want)
	}
}

func TestWriteUserFileRefusesSymlinks(t *testing.T) {
	dir := t.TempDir()
	victim := filepath.Join(dir, "victim")
	if err := os.WriteFile(victim, []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}

	out := filepath.Join(dir, "out.csv")
	if err := writeUserFile(out, []byte("a,b\n")); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(out); string(got) != "a,b\n" {
		t.Fatalf("wrote %q", got)
	}

	// A symlinked folder component is refused, not followed.
	real := filepath.Join(dir, "real")
	if err := os.Mkdir(real, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(real, filepath.Join(dir, "link")); err != nil {
		t.Fatal(err)
	}
	if err := writeUserFile(filepath.Join(dir, "link", "x.csv"), []byte("x")); err == nil {
		t.Fatal("write through a symlinked folder succeeded")
	}
	if _, err := os.Stat(filepath.Join(real, "x.csv")); err == nil {
		t.Fatal("file landed behind the symlink")
	}

	// A symlink at the file name is replaced, its target untouched.
	planted := filepath.Join(dir, "planted.csv")
	if err := os.Symlink(victim, planted); err != nil {
		t.Fatal(err)
	}
	if err := writeUserFile(planted, []byte("new")); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.ReadFile(victim); string(got) != "keep" {
		t.Fatalf("symlink target overwritten: %q", got)
	}
	if fi, err := os.Lstat(planted); err != nil || !fi.Mode().IsRegular() {
		t.Fatalf("planted path is not a regular file now: %v %v", fi, err)
	}

	if err := writeUserFile("relative.csv", nil); err == nil {
		t.Fatal("relative path accepted")
	}
}

func TestAdaptiveGridDropsColumns(t *testing.T) {
	cell := func() fyne.CanvasObject {
		r := canvas.NewRectangle(nil)
		r.SetMinSize(fyne.NewSize(200, 50))
		return r
	}
	g := newAdaptiveGrid(4, cell(), cell(), cell(), cell())
	g.Resize(fyne.NewSize(1000, 50))
	if g.MinSize().Height != 50 {
		t.Fatalf("4 columns: min height %v, want one row", g.MinSize().Height)
	}
	g.Resize(fyne.NewSize(450, 50))
	if want := float32(50*2 + cardGap); g.MinSize().Height != want {
		t.Fatalf("2 columns: min height %v, want %v", g.MinSize().Height, want)
	}
	if g.MinSize().Width != 200 {
		t.Fatalf("min width %v, want one column", g.MinSize().Width)
	}
	if x := g.Objects[1].Position().X; x != (450-cardGap)/2+cardGap {
		t.Fatalf("second card at x=%v", x)
	}
}
