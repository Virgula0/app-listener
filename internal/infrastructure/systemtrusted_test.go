package ebpf

import (
	"os"
	"path/filepath"
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

func TestCheckSystemTrustedRootOwnedBinary(t *testing.T) {
	if err := CheckSystemTrusted(openConfined(t, "/usr/bin/true")); err != nil {
		t.Fatalf("/usr/bin/true is a root-owned file in a root-owned dir: %v", err)
	}
}

func TestCheckSystemTrustedRefusesUserFile(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root owns the temp file")
	}
	p := filepath.Join(t.TempDir(), "app")
	if err := os.WriteFile(p, []byte("x"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := CheckSystemTrusted(openConfined(t, p)); err == nil {
		t.Fatal("a user-owned file must not be a system file")
	}
}

func TestCheckSystemTrustedFollowsLinkToSystemFile(t *testing.T) {
	// The link is followed: the judged inode is /usr/bin/true in its own directory, not the link.
	p := filepath.Join(t.TempDir(), "link")
	if err := os.Symlink("/usr/bin/true", p); err != nil {
		t.Fatal(err)
	}
	if err := CheckSystemTrusted(openConfined(t, p)); err != nil {
		t.Fatalf("a link to a system file reaches the system file: %v", err)
	}
}

func TestCheckSystemTrustedRefusesDirectory(t *testing.T) {
	if err := CheckSystemTrusted(openConfined(t, "/usr/bin")); err == nil {
		t.Fatal("a directory is not a system binary")
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
