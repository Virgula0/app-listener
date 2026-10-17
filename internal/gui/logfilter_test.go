//go:build gui

package gui

import (
	"bytes"
	"log"
	"testing"
)

func TestFyneLogFilterDropsBenignReports(t *testing.T) {
	var buf bytes.Buffer
	l := log.New(&fyneLogFilter{out: &buf}, "", 0)
	l.Println("Fyne error: ", "Getting favorite locations")
	l.Println("  Cause:", "this computer does not define a VIDEOS folder")
	l.Printf("  At: %s:%d", "dialog/file.go", 399)
	l.Println("Fyne error: ", "Failed to load image")
	l.Println("  Cause:", "boom")
	l.Println("  At: x.go:1")
	l.Println("Fyne error: ", "failed to watch theme settings")
	l.Println("  Cause:", "dbus")
	l.Println("unrelated")
	want := "Fyne error:  Failed to load image\n  Cause: boom\n  At: x.go:1\nunrelated\n"
	if buf.String() != want {
		t.Fatalf("filtered output = %q, want %q", buf.String(), want)
	}
}
