package headless

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// buildTiny builds an empty main of module mod with tags.
func buildTiny(t *testing.T, mod, tags string) string {
	t.Helper()
	dir := t.TempDir()
	must := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	must(os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module "+mod+"\n\ngo 1.21\n"), 0o600))
	must(os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n\nfunc main() {}\n"), 0o600))
	out := filepath.Join(dir, "tiny")
	cmd := exec.Command("go", "build", "-tags", tags, "-o", out, ".")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0", "GOFLAGS=")
	if b, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, b)
	}
	return out
}

func TestCheck(t *testing.T) {
	if err := Check(buildTiny(t, mainPath, "netgo,ci")); err != nil {
		t.Fatalf("a headless build must pass: %v", err)
	}
	if err := Check(buildTiny(t, mainPath, "netgo,gui")); err == nil {
		t.Fatal("a -tags gui build must be refused")
	}
	if err := Check(buildTiny(t, "example.com/other", "netgo")); err == nil {
		t.Fatal("another program must be refused")
	}
	notELF := filepath.Join(t.TempDir(), "x")
	if err := os.WriteFile(notELF, []byte("#!/bin/sh\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := Check(notELF); err == nil {
		t.Fatal("a file without build info must be refused")
	}
}

func TestCheckGUI(t *testing.T) {
	if err := CheckGUI(buildTiny(t, mainPath, "netgo,gui")); err != nil {
		t.Fatalf("a -tags gui build must pass: %v", err)
	}
	if err := CheckGUI(buildTiny(t, mainPath, "netgo,ci")); err == nil {
		t.Fatal("a headless build must be refused")
	}
	if err := CheckGUI(buildTiny(t, "example.com/other", "gui")); err == nil {
		t.Fatal("another program must be refused")
	}
}
