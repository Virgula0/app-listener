package ebpf

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// RUNPATH=$ORIGIN often points into a user-writable $HOME dir: a library planted there must be
// refused (warnUntrustedLibs), never classed as system-trusted.
func TestResolveLibraryClosureRejectsUserWritableOriginDir(t *testing.T) {
	gcc, err := exec.LookPath("gcc")
	if err != nil {
		t.Skip("gcc not available")
	}

	// Stands in for a whitelisted app installed under the user's home: the attacker owns this
	// directory, so they can create files in it.
	appDir := t.TempDir()

	libSrc := filepath.Join(appDir, "poc.c")
	if werr := os.WriteFile(libSrc, []byte("int poc_entry(void){return 0;}\n"), 0o600); werr != nil {
		t.Fatal(werr)
	}
	rogueLib := filepath.Join(appDir, "libpoc.so")
	if out, cerr := exec.Command(gcc, "-shared", "-fPIC", "-o", rogueLib, libSrc).CombinedOutput(); cerr != nil {
		t.Skipf("cannot build shared object: %v: %s", cerr, out)
	}

	binSrc := filepath.Join(appDir, "app.c")
	if werr := os.WriteFile(binSrc, []byte("int poc_entry(void);\nint main(void){return poc_entry();}\n"), 0o600); werr != nil {
		t.Fatal(werr)
	}
	appBin := filepath.Join(appDir, "app")
	out, cerr := exec.Command(gcc, "-o", appBin, binSrc,
		"-L"+appDir, "-lpoc", "-Wl,-rpath,$ORIGIN").CombinedOutput()
	if cerr != nil {
		t.Skipf("cannot build binary with $ORIGIN rpath: %v: %s", cerr, out)
	}

	closure, refused, rerr := LibraryClosure(appBin)
	if rerr != nil {
		t.Fatalf("LibraryClosure(%s): %v", appBin, rerr)
	}
	// Ownership is judged before the mount, so this holds without root (no /proc/1/root).
	if why := refused[rogueLib]; why == nil || !strings.Contains(why.Error(), "not root-owned") {
		t.Errorf("rogue library refused for %v; want its user ownership", why)
	}

	if slices.Contains(closure, rogueLib) {
		t.Fatalf("closure classes a library from a user-writable $ORIGIN directory as system-trusted: "+
			"%s\nfull closure: %v", rogueLib, closure)
	}
}
