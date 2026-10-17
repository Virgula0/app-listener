//go:build gui

package gui

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"testing"
)

func TestStageLocal(t *testing.T) {
	dir := t.TempDir()
	data := []byte("\x7fELF headless")
	b, err := stageLocal(bytes.NewReader(data), dir, func(string) error { return nil })
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	if b.sha != hex.EncodeToString(sum[:]) || b.tag != "" {
		t.Errorf("staged %+v, want sha %x and no tag", b, sum)
	}
	fi, err := os.Stat(b.path)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != 0o700 {
		t.Errorf("mode %v, want 0700", fi.Mode().Perm())
	}
	if got, _ := os.ReadFile(b.path); !bytes.Equal(got, data) {
		t.Errorf("content %q, want %q", got, data)
	}
}

func TestStageLocalRefusedIsRemoved(t *testing.T) {
	dir := t.TempDir()
	refuse := errors.New("the GUI build")
	b, err := stageLocal(bytes.NewReader([]byte("x")), dir, func(string) error { return refuse })
	if !errors.Is(err, refuse) || b != nil {
		t.Fatalf("got %v, %v; want the check's refusal", b, err)
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Errorf("refused build left %d file(s) in staging", len(entries))
	}
}
