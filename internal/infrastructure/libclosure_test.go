package ebpf

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveLibraryClosure(t *testing.T) {
	// Resolve the closure of a known dynamically linked binary. /bin/cat (or
	// its coreutils path) links libc and has an interpreter.
	bin := ""
	for _, cand := range []string{"/bin/cat", "/usr/bin/cat", "/bin/ls", "/usr/bin/ls"} {
		if fileExists(cand) {
			bin = cand
			break
		}
	}
	if bin == "" {
		t.Skip("no known dynamic binary available")
	}

	closure, err := ResolveLibraryClosure(bin)
	if err != nil {
		t.Fatalf("ResolveLibraryClosure(%s): %v", bin, err)
	}
	if len(closure) == 0 {
		t.Fatalf("empty closure for %s (expected libc + interpreter)", bin)
	}

	var haveInterp, haveLibc bool
	for _, p := range closure {
		if !filepath.IsAbs(p) {
			t.Errorf("closure entry not absolute: %q", p)
		}
		if !fileExists(p) {
			t.Errorf("closure entry does not exist: %q", p)
		}
		base := filepath.Base(p)
		if strings.HasPrefix(base, "ld-linux") || base == "ld.so" || strings.HasPrefix(base, "ld64") {
			haveInterp = true
		}
		if strings.HasPrefix(base, "libc.so") {
			haveLibc = true
		}
	}
	if !haveInterp {
		t.Errorf("closure missing the program interpreter: %v", closure)
	}
	if !haveLibc {
		t.Errorf("closure missing libc: %v", closure)
	}
}

func TestResolveLibraryClosureStaticIsEmpty(t *testing.T) {
	// A non-ELF or static object must not error hard; craft a tiny non-ELF
	// file and confirm elf.Open fails cleanly (error path), not a panic.
	dir := t.TempDir()
	notELF := filepath.Join(dir, "notelf")
	if err := writeFileHelper(notELF, "not an elf"); err != nil {
		t.Fatal(err)
	}
	if _, err := ResolveLibraryClosure(notELF); err == nil {
		t.Errorf("expected an error opening a non-ELF file")
	}
}

func writeFileHelper(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}
