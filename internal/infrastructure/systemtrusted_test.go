package ebpf

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func openConfined(t *testing.T, path string) *os.File {
	t.Helper()
	f, err := OpenConfined(path)
	if err != nil {
		t.Fatalf("OpenConfined(%s): %v", path, err)
	}
	t.Cleanup(func() { f.Close() })
	return f
}

func TestOpenSystemPlacedRootOwnedBinary(t *testing.T) {
	paths := []string{"/usr/bin/true", "/usr/bin/../bin/./true", "//usr/bin/true"}
	if fi, err := os.Lstat("/bin"); err == nil && fi.Mode()&os.ModeSymlink != 0 {
		paths = append(paths, "/bin/true") // merged /usr: a root-owned link in a root-owned dir
	}
	want, err := os.Stat("/usr/bin/true")
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range paths {
		f, err := OpenSystemPlaced(p)
		if err != nil {
			t.Errorf("OpenSystemPlaced(%s): %v", p, err)
			continue
		}
		got, err := f.Stat()
		f.Close()
		if err != nil || !os.SameFile(got, want) {
			t.Errorf("OpenSystemPlaced(%s) opened another inode (%v)", p, err)
		}
	}
}

// The test user owns its temp dir: whatever a name there points at, root did not place the name.
func TestOpenSystemPlacedRefusesUserPlacedNames(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root owns the temp dir")
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "app"), []byte("x"), 0o755); err != nil {
		t.Fatal(err)
	}
	for l, target := range map[string]string{"link": "/usr/bin/true", "bin": "/usr/bin", "rel": "app"} {
		if err := os.Symlink(target, filepath.Join(dir, l)); err != nil {
			t.Fatal(err)
		}
	}
	for _, p := range []string{"app", "link", "bin/true", "rel", "bin/../bin/true"} {
		if f, err := OpenSystemPlaced(filepath.Join(dir, p)); err == nil {
			f.Close()
			t.Errorf("%s: a name in a user-owned directory was judged root-placed", p)
		}
	}
}

func TestOpenSystemPlacedRefusesUserHome(t *testing.T) {
	SetUserHomes([]string{"/usr/lib"})
	t.Cleanup(func() { SetUserHomes(nil) })
	if f, err := OpenSystemPlaced("/usr/bin/true"); err != nil {
		t.Fatalf("/usr/bin/true is outside the home: %v", err)
	} else {
		f.Close()
	}
	SetUserHomes([]string{"/usr"})
	if f, err := OpenSystemPlaced("/usr/bin/true"); err == nil {
		f.Close()
		t.Fatal("a name in a user home was judged root-placed")
	} else if !strings.Contains(err.Error(), "user home") {
		t.Fatalf("refused for %v; want the user-home reason", err)
	}
}

func TestOpenSystemPlacedRefusesNonFiles(t *testing.T) {
	for _, p := range []string{"/usr/bin", "/", "usr/bin/true", "/usr/bin/true/x", "/nonexistent-app-listener"} {
		if f, err := OpenSystemPlaced(p); err == nil {
			f.Close()
			t.Errorf("OpenSystemPlaced(%s) succeeded", p)
		}
	}
}

func TestGlobSystemPlaced(t *testing.T) {
	if got := GlobSystemPlaced("/usr/bin/tru[e]", 8); !slices.Equal(got, []string{"/usr/bin/true"}) {
		t.Errorf("root-placed glob = %v", got)
	}
	if got := GlobSystemPlaced("/usr/bin/*", 3); len(got) != 3 || !slices.IsSorted(got) {
		t.Errorf("limit 3 = %v", got)
	}
	if got := GlobSystemPlaced("relative/*", 8); got != nil {
		t.Errorf("relative pattern = %v", got)
	}
	if os.Geteuid() == 0 {
		return
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "true"), nil, 0o755); err != nil {
		t.Fatal(err)
	}
	if got := GlobSystemPlaced(filepath.Join(dir, "*"), 8); got != nil {
		t.Errorf("a user-owned dir was expanded: %v", got)
	}
}

func TestComputeBinaryEntryFileMatchesPath(t *testing.T) {
	want, err := ComputeBinaryEntry("/usr/bin/true")
	if err != nil {
		t.Fatal(err)
	}
	got, err := ComputeBinaryEntryFile(openConfined(t, "/usr/bin/true"), "/usr/bin/true")
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("entry from the fd %+v != from the path %+v", got, want)
	}
}
