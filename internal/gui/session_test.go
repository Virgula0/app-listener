//go:build gui

package gui

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadOwnedFileRefusesSymlink(t *testing.T) {
	dir := t.TempDir()
	real := filepath.Join(dir, "cookie")
	if err := os.WriteFile(real, []byte("MIT-MAGIC"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, ".Xauthority")
	if err := os.Symlink(real, link); err != nil {
		t.Fatal(err)
	}
	uid := uint32(os.Getuid())
	if _, ok := readOwnedFile(link, uid); ok {
		t.Error("a symlinked cookie file must be refused")
	}
	if data, ok := readOwnedFile(real, uid); !ok || string(data) != "MIT-MAGIC" {
		t.Errorf("the user's own cookie file must be read, got %q ok=%v", data, ok)
	}
	if _, ok := readOwnedFile(real, uid+1); ok {
		t.Error("a cookie file owned by another user must be refused")
	}
}

func TestOwnedByChecksTypeAndOwner(t *testing.T) {
	dir := t.TempDir()
	uid := uint32(os.Getuid())
	if !ownedBy(dir, uid, os.ModeDir) {
		t.Error("own directory must pass")
	}
	if ownedBy(dir, uid+1, os.ModeDir) {
		t.Error("another user's directory must fail")
	}
	if ownedBy(dir, uid, os.ModeSocket) {
		t.Error("a directory is not a socket")
	}
}

func TestWaylandSocketNameMustNotTraverse(t *testing.T) {
	t.Setenv("XDG_RUNTIME_DIR", t.TempDir())
	if waylandSocketOK("../wayland-1", uint32(os.Getuid())) {
		t.Error("a relative name with a slash must be refused")
	}
}

func TestStatusField(t *testing.T) {
	status := []byte("Name:\tzsh\nPPid:\t42\nUid:\t1000\t1000\t1000\t1000\n")
	if got := statusField(status, "Uid:"); got != "1000" {
		t.Errorf("Uid = %q", got)
	}
	if got := statusField(status, "PPid:"); got != "42" {
		t.Errorf("PPid = %q", got)
	}
}

func TestDisplayValue(t *testing.T) {
	for v, want := range map[string]bool{":1": true, ":0.0": true, "localhost:10.0": true, "": false, ":1;rm": false} {
		if got := displayRe.MatchString(v); got != want {
			t.Errorf("displayRe(%q) = %v, want %v", v, got, want)
		}
	}
}
