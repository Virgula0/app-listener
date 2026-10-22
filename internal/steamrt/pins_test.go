package steamrt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadPinsAcceptsSetupOutput(t *testing.T) {
	upper := t.TempDir()
	write(t, filepath.Join(upper, "pinned_libs_64/done"), "")
	write(t, filepath.Join(upper, "pinned_libs_64/has_pins"), "")
	write(t, filepath.Join(upper, "pinned_libs_64/system_libc.so.6"), "/usr/lib/libc.so.6\n/usr/lib/libc.so.6\n")
	if err := os.Symlink("/usr/lib/libz.so.1", filepath.Join(upper, "pinned_libs_64/libz.so.1")); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(upper, "libcurl_compat_32"), 0o755); err != nil {
		t.Fatal(err)
	}
	ops, err := readPins(upper)
	if err != nil {
		t.Fatal(err)
	}
	if len(ops) != 2 {
		t.Fatalf("ops %+v", ops)
	}
}

func TestReadPinsRefusesAnythingElse(t *testing.T) {
	for name, setup := range map[string]func(upper string){
		"library file":     func(u string) { write(t, filepath.Join(u, "pinned_libs_64/libevil.so"), "\x7fELF") },
		"marker with data": func(u string) { write(t, filepath.Join(u, "pinned_libs_32/done"), "x") },
		"binary record":    func(u string) { write(t, filepath.Join(u, "pinned_libs_32/system_x"), "a\x00b") },
		"big record": func(u string) {
			write(t, filepath.Join(u, "pinned_libs_32/system_x"), strings.Repeat("a", maxPinFile+1))
		},
		"other top dir":  func(u string) { write(t, filepath.Join(u, "usr/lib/libevil.so"), "x") },
		"copied-up file": func(u string) { write(t, filepath.Join(u, "setup.sh"), "x") },
		"subdirectory": func(u string) {
			if err := os.MkdirAll(filepath.Join(u, "pinned_libs_64/sub"), 0o755); err != nil {
				t.Fatal(err)
			}
		},
		"pin dir is a file": func(u string) { write(t, filepath.Join(u, "pinned_libs_64"), "") },
	} {
		t.Run(name, func(t *testing.T) {
			upper := t.TempDir()
			setup(upper)
			if _, err := readPins(upper); err == nil {
				t.Fatal("accepted")
			}
		})
	}
}

func TestApplyPinsReplacesOpaqueAndMergesOthers(t *testing.T) {
	r := newRuntimeDir(t)
	rt := filepath.Join(r.Dir, r.Name)
	write(t, filepath.Join(rt, "pinned_libs_64/stale"), "")
	write(t, filepath.Join(rt, "pinned_libs_32/keep"), "")
	write(t, filepath.Join(rt, "libcurl_compat_64/gone"), "")
	ops := []pinOp{
		{dir: "pinned_libs_64", replace: true, entries: []pinEntry{{name: "done"}, {name: "libz.so.1", link: "/usr/lib/libz.so.1"}}},
		{dir: "pinned_libs_32", entries: []pinEntry{{name: "done"}, {name: "keep", whiteout: true}}},
		{dir: "libcurl_compat_64", removed: true},
	}
	if err := r.applyPins(ops); err != nil {
		t.Fatal(err)
	}
	if got, _ := os.Readlink(filepath.Join(rt, "pinned_libs_64/libz.so.1")); got != "/usr/lib/libz.so.1" {
		t.Fatalf("pin symlink %q", got)
	}
	for _, gone := range []string{"pinned_libs_64/stale", "pinned_libs_32/keep", "libcurl_compat_64"} {
		if _, err := os.Lstat(filepath.Join(rt, gone)); !os.IsNotExist(err) {
			t.Fatalf("%s must be gone: %v", gone, err)
		}
	}
	if _, err := os.Lstat(filepath.Join(rt, "pinned_libs_32/done")); err != nil {
		t.Fatal(err)
	}
}

func TestPinsStale(t *testing.T) {
	r := newRuntimeDir(t)
	rt := filepath.Join(r.Dir, r.Name)
	write(t, filepath.Join(rt, "setup.sh"), "")
	if !r.PinsStale() {
		t.Fatal("no done markers: stale")
	}
	write(t, filepath.Join(rt, "pinned_libs_32/done"), "")
	write(t, filepath.Join(rt, "pinned_libs_64/done"), "")
	if r.PinsStale() {
		t.Fatal("done, no pins: not stale")
	}
	host := t.TempDir()
	write(t, filepath.Join(host, "libx.so.1.2"), "")
	if err := os.Symlink("libx.so.1.2", filepath.Join(host, "libx.so.1")); err != nil {
		t.Fatal(err)
	}
	write(t, filepath.Join(rt, "pinned_libs_64/has_pins"), "")
	record := filepath.Join(rt, "pinned_libs_64/system_libx.so.1")
	write(t, record, filepath.Join(host, "libx.so.1")+"\n"+filepath.Join(host, "libx.so.1.2")+"\n")
	if r.PinsStale() {
		t.Fatal("pin resolves to its record: not stale")
	}
	write(t, record, filepath.Join(host, "libx.so.1")+"\n"+filepath.Join(host, "libx.so.1.1")+"\n")
	if !r.PinsStale() {
		t.Fatal("host library upgraded: stale")
	}
}
