package ebpf

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"
)

// ResolveLibraryClosure feeds the daemon-wide trusted-library set (guard_trusted_files,
// TRUSTED_LIB), and TRUSTED_LIB is an unconditional pass in trust_mmap — the check that is supposed
// to stop LD_PRELOAD of attacker code into a whitelisted process. The closure resolves DT_NEEDED
// using the binary's OWN DT_RPATH/DT_RUNPATH with $ORIGIN expanded, searched BEFORE the system
// directories, and filters nothing by ownership.
//
// Many whitelisted binaries in the catalog live in user-writable, unguarded directories
// (~/.config/discord/*/Discord, ~/.local/bin/*, ~/.foundry/bin/*) and ship RUNPATH=$ORIGIN, so a
// same-user attacker who drops a .so beside such a binary gets its inode admitted to the trusted
// set — and the set is daemon-wide, so it is then loadable into /usr/bin/ssh too.
//
// The closure must apply the same rule the kernel side uses for auto-trust (is_system_trusted):
// root-owned and not writable by a normal user. Anything else has to be an explicit, reviewed
// allow_lib.
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

	closure, rerr := ResolveLibraryClosure(appBin)
	if rerr != nil {
		t.Fatalf("ResolveLibraryClosure(%s): %v", appBin, rerr)
	}

	if slices.Contains(closure, rogueLib) {
		t.Fatalf("closure trusts a library from a user-writable $ORIGIN directory: %s\n"+
			"it would be added daemon-wide as TRUSTED_LIB and become loadable into every "+
			"whitelisted binary\nfull closure: %v", rogueLib, closure)
	}
}
