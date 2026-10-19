package guard

import (
	"debug/elf"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

// hostLoader is the PT_INTERP of this (cgo, dynamically linked) test binary.
func hostLoader(t *testing.T) string {
	t.Helper()
	self, err := os.Executable()
	if err != nil {
		t.Skip(err)
	}
	e, err := elf.Open(self)
	if err != nil {
		t.Skip(err)
	}
	defer e.Close()
	for _, p := range e.Progs {
		if p.Type == elf.PT_INTERP {
			b, err := io.ReadAll(p.Open())
			if err == nil {
				return strings.TrimRight(string(b), "\x00")
			}
		}
	}
	t.Skip("statically linked test binary")
	return ""
}

func TestLoaderRefusal(t *testing.T) {
	dir := t.TempDir()
	plain := filepath.Join(dir, "tool")
	if err := os.WriteFile(plain, []byte("\x7fELF plain"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{plain, filepath.Join(dir, "missing"), "/usr/lib/libldap.so.2"} {
		if err := LoaderRefusal(p); err != nil {
			t.Errorf("%s: %v", p, err)
		}
	}
	for _, n := range []string{"ld-linux-x86-64.so.2", "ld-2.31.so", "ld-musl-x86_64.so.1", "ld64.so.2", "ld.so"} {
		p := filepath.Join(dir, n)
		if err := os.WriteFile(p, []byte("x"), 0o755); err != nil {
			t.Fatal(err)
		}
		if LoaderRefusal(p) == nil {
			t.Errorf("%s: a loader name must be refused", n)
		}
	}

	ld := hostLoader(t)
	if LoaderRefusal(ld) == nil {
		t.Fatalf("the host loader %s must be refused", ld)
	}
	link := filepath.Join(dir, "runner")
	if err := os.Symlink(ld, link); err != nil {
		t.Fatal(err)
	}
	if LoaderRefusal(link) == nil {
		t.Fatal("a symlink to the loader must be refused")
	}
	b, err := os.ReadFile(ld)
	if err != nil {
		t.Skip(err)
	}
	renamed := filepath.Join(dir, "innocent")
	if err := os.WriteFile(renamed, b, 0o755); err != nil {
		t.Fatal(err)
	}
	if LoaderRefusal(renamed) == nil {
		t.Fatal("a renamed copy of the loader must be refused by its DT_SONAME")
	}
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{{Path: dir,
		Binaries: []daemonconfig.BinaryRule{{Path: plain}, {Path: renamed}}}}}
	if err := RefuseIdentityLines(cfg, nil); err == nil || !strings.Contains(err.Error(), "dynamic loader") {
		t.Fatalf("a new loader line must be refused: %v", err)
	}
}
