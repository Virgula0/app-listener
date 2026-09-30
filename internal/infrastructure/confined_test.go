package ebpf

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A test runs as its own user, so every temp dir is user-writable and confinement applies.
func TestResolveConfined(t *testing.T) {
	base, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	outside, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"files/real", "files/bin.old", "other"} {
		if err := os.MkdirAll(filepath.Join(base, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for _, f := range []string{"files/real/wineserver", "other/target", filepath.Join(outside, "wineserver")} {
		if !filepath.IsAbs(f) {
			f = filepath.Join(base, f)
		}
		if err := os.WriteFile(f, nil, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	links := map[string]string{
		"files/bin":     outside,                // absolute, escaping
		"files/up":      "../../..",             // relative, escaping
		"files/inside":  "real",                 // stays in its parent
		"files/abs":     base + "/files/real",   // absolute: refused even inside
		"files/final":   base + "/other/target", // final component: followed freely
		"files/chained": "inside/../up",         // escapes only through a second link
	}
	for l, target := range links {
		if err := os.Symlink(target, filepath.Join(base, l)); err != nil {
			t.Fatal(err)
		}
	}

	for _, p := range []string{"files/bin/wineserver", "files/up/x", "files/abs/wineserver", "files/chained/x"} {
		if got, err := ResolveConfined(filepath.Join(base, p)); err == nil {
			t.Errorf("%s resolved to %s through a symlink leaving its user-writable parent", p, got)
		} else if !strings.Contains(err.Error(), "leaving its parent directory") {
			t.Errorf("%s: refused for %v; want the confinement reason", p, err)
		}
	}
	for p, want := range map[string]string{
		"files/inside/wineserver": base + "/files/real/wineserver",
		"files/final":             base + "/other/target",
		"files/real/wineserver":   base + "/files/real/wineserver",
	} {
		if got, err := ResolveConfined(filepath.Join(base, p)); err != nil || got != want {
			t.Errorf("ResolveConfined(%s) = %q, %v; want %q", p, got, err, want)
		}
	}
	if _, _, err := StatConfined(filepath.Join(base, "files/missing")); !os.IsNotExist(err) {
		t.Errorf("a missing path must stay a not-exist error (the caller defers it): %v", err)
	}
}

// Root-owned parents outside a home resolve as usual; inside one (a root-owned catalog tree under
// /root) an escaping symlink is refused all the same.
func TestResolveConfinedRootOwnedParents(t *testing.T) {
	if _, _, err := StatConfined("/proc/self/exe"); err != nil {
		t.Fatalf("StatConfined(/proc/self/exe): %v", err)
	}
	fi, err := os.Lstat("/var/run")
	if err != nil || fi.Mode()&os.ModeSymlink == 0 {
		t.Skip("/var/run is not a symlink here")
	}
	entries, err := os.ReadDir("/var/run/")
	if err != nil || len(entries) == 0 {
		t.Skip("/run is empty or unreadable")
	}
	p := "/var/run/" + entries[0].Name()
	defer SetUserHomes(nil)

	SetUserHomes(nil)
	if _, err := ResolveConfined(p); err != nil {
		t.Fatalf("root-owned /var/run -> /run must be followed outside a home: %v", err)
	}
	SetUserHomes([]string{"/var"})
	if got, err := ResolveConfined(p); err == nil {
		t.Errorf("/var as a home: %s resolved to %s through a symlink leaving /var", p, got)
	}
	SetUserHomes([]string{"/"})
	if _, err := ResolveConfined(p); err == nil {
		t.Errorf(`homes ["/"] (unreadable passwd) must confine every symlinked directory`)
	}
}
