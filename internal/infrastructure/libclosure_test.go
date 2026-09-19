package ebpf

import (
	"debug/elf"
	"encoding/binary"
	"os"
	"path/filepath"
	"runtime"
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

// TestResolveLibraryClosureScript: a whitelisted script runs as its
// interpreter, so it has no closure of its own — and must not be reported as
// a broken ELF (it filled the daemon log with "bad magic number" warnings).
func TestResolveLibraryClosureScript(t *testing.T) {
	p := filepath.Join(t.TempDir(), "launcher.sh")
	if err := os.WriteFile(p, []byte("#!/bin/bash\nexec true\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	libs, err := ResolveLibraryClosure(p)
	if err != nil || len(libs) != 0 {
		t.Fatalf("ResolveLibraryClosure(script) = %v, %v; want empty, nil", libs, err)
	}
}

// TestResolveLibraryClosureForeignArch: an ELF for an architecture the host
// cannot run natively (Steam's pressure-vessel-arm64, meant for FEX) has an
// interpreter that does not exist here — it is skipped, not warned about.
func TestResolveLibraryClosureForeignArch(t *testing.T) {
	foreign := elf.EM_AARCH64
	if runtime.GOARCH == "arm64" {
		foreign = elf.EM_X86_64
	}
	// Minimal ELF64 little-endian header: no sections, no program headers.
	hdr := make([]byte, 64)
	copy(hdr, []byte{0x7f, 'E', 'L', 'F', byte(elf.ELFCLASS64), byte(elf.ELFDATA2LSB), byte(elf.EV_CURRENT)})
	binary.LittleEndian.PutUint16(hdr[16:], uint16(elf.ET_EXEC))
	binary.LittleEndian.PutUint16(hdr[18:], uint16(foreign))
	binary.LittleEndian.PutUint32(hdr[20:], uint32(elf.EV_CURRENT))
	binary.LittleEndian.PutUint16(hdr[52:], 64) // e_ehsize
	p := filepath.Join(t.TempDir(), "pressure-vessel-wrap")
	if err := os.WriteFile(p, hdr, 0o755); err != nil {
		t.Fatal(err)
	}
	libs, err := ResolveLibraryClosure(p)
	if err != nil || len(libs) != 0 {
		t.Fatalf("ResolveLibraryClosure(foreign ELF) = %v, %v; want empty, nil", libs, err)
	}
}

func TestRunsNatively(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skip("architecture table asserted for amd64 only")
	}
	if !runsNatively(elf.EM_X86_64) || !runsNatively(elf.EM_386) {
		t.Errorf("x86_64 and i386 (Steam's 32-bit client) must count as native on amd64")
	}
	if runsNatively(elf.EM_AARCH64) {
		t.Errorf("aarch64 must not count as native on amd64")
	}
}
