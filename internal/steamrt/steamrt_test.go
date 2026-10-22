package steamrt

import (
	"archive/tar"
	"bytes"
	"crypto/md5" //nolint:gosec // matches the production check
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/unix"
)

type tarEntry struct {
	name, body, link string
	typ              byte
	mode             int64
}

func buildTar(t *testing.T, entries []tarEntry) *tar.Reader {
	t.Helper()
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for _, e := range entries {
		mode := e.mode
		if mode == 0 {
			mode = 0o644
		}
		h := &tar.Header{Name: e.name, Typeflag: e.typ, Mode: mode, Linkname: e.link, Size: int64(len(e.body)),
			ModTime: time.Unix(1_700_000_000, 0)}
		if e.typ != tar.TypeReg {
			h.Size = 0
		}
		if err := tw.WriteHeader(h); err != nil {
			t.Fatal(err)
		}
		if e.typ == tar.TypeReg {
			if _, err := tw.Write([]byte(e.body)); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	return tar.NewReader(&buf)
}

func extractInto(t *testing.T, entries []tarEntry) (string, error) {
	t.Helper()
	root := t.TempDir()
	fd, err := openDir(root)
	if err != nil {
		t.Fatal(err)
	}
	defer unix.Close(fd)
	e := extractor{root: fd, top: "steam-runtime", uid: os.Getuid(), gid: os.Getgid()}
	return root, e.all(buildTar(t, entries))
}

func TestExtractBuildsTreeWithoutTopEntry(t *testing.T) {
	root, err := extractInto(t, []tarEntry{
		{name: "steam-runtime/setup.sh", typ: tar.TypeReg, body: "#!/bin/sh\n", mode: 0o6755},
		{name: "steam-runtime/usr/lib/", typ: tar.TypeDir, mode: 0o755},
		{name: "steam-runtime/usr/lib/libdemo.so.1", typ: tar.TypeReg, body: "ELF"},
		{name: "steam-runtime/usr/lib/libdemo.so", typ: tar.TypeSymlink, link: "libdemo.so.1"},
		{name: "steam-runtime/usr/lib/libdemo.hard", typ: tar.TypeLink, link: "steam-runtime/usr/lib/libdemo.so.1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(filepath.Join(root, "steam-runtime/setup.sh"))
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != 0o755 || fi.Mode()&(os.ModeSetuid|os.ModeSetgid) != 0 {
		t.Fatalf("setup.sh mode %v: setuid/setgid must be stripped", fi.Mode())
	}
	if got, _ := os.Readlink(filepath.Join(root, "steam-runtime/usr/lib/libdemo.so")); got != "libdemo.so.1" {
		t.Fatalf("symlink target %q", got)
	}
	if b, _ := os.ReadFile(filepath.Join(root, "steam-runtime/usr/lib/libdemo.hard")); string(b) != "ELF" {
		t.Fatalf("hard link content %q", b)
	}
}

func TestExtractRefusesEscapes(t *testing.T) {
	for name, entries := range map[string][]tarEntry{
		"dotdot":        {{name: "steam-runtime/../evil", typ: tar.TypeReg, body: "x"}},
		"absolute":      {{name: "/etc/evil", typ: tar.TypeReg, body: "x"}},
		"outside top":   {{name: "other/evil", typ: tar.TypeReg, body: "x"}},
		"top prefix":    {{name: "steam-runtime-evil/x", typ: tar.TypeReg, body: "x"}},
		"through link":  {{name: "steam-runtime/l", typ: tar.TypeSymlink, link: "/tmp"}, {name: "steam-runtime/l/evil", typ: tar.TypeReg, body: "x"}},
		"device":        {{name: "steam-runtime/dev", typ: tar.TypeChar}},
		"fifo":          {{name: "steam-runtime/fifo", typ: tar.TypeFifo}},
		"hard outside":  {{name: "steam-runtime/h", typ: tar.TypeLink, link: "../../etc/passwd"}},
		"hard to link":  {{name: "steam-runtime/l", typ: tar.TypeSymlink, link: "/etc/passwd"}, {name: "steam-runtime/h", typ: tar.TypeLink, link: "steam-runtime/l"}},
		"over existing": {{name: "steam-runtime/f", typ: tar.TypeSymlink, link: "/tmp/x"}, {name: "steam-runtime/f", typ: tar.TypeReg, body: "x"}},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := extractInto(t, entries); err == nil {
				t.Fatal("extracted")
			}
		})
	}
}

func newRuntimeDir(t *testing.T) Runtime {
	t.Helper()
	return Runtime{Dir: t.TempDir(), Name: "steam-runtime", UID: uint32(os.Getuid()), GID: uint32(os.Getgid())} //nolint:gosec // test uid
}

func write(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestPending(t *testing.T) {
	r := newRuntimeDir(t)
	if p, err := r.Pending(); err != nil || p {
		t.Fatalf("no archive: pending=%v err=%v", p, err)
	}
	write(t, filepath.Join(r.Dir, "steam-runtime.tar.xz"), "archive")
	write(t, filepath.Join(r.Dir, "steam-runtime.tar.xz.checksum"), "new  steam-runtime.tar.xz\n")
	if p, _ := r.Pending(); !p {
		t.Fatal("no unpacked tree: must be pending")
	}
	write(t, filepath.Join(r.Dir, "steam-runtime/checksum"), "old  steam-runtime.tar.xz\n")
	if p, _ := r.Pending(); !p {
		t.Fatal("different checksum: must be pending")
	}
	write(t, filepath.Join(r.Dir, "steam-runtime/checksum"), "new  steam-runtime.tar.xz\n")
	if p, _ := r.Pending(); p {
		t.Fatal("same checksum: not pending")
	}
}

func TestVerifyMD5(t *testing.T) {
	r := newRuntimeDir(t)
	path := filepath.Join(r.Dir, "steam-runtime.tar.xz")
	write(t, path, "archive bytes")
	f, err := os.Open(path) //nolint:gosec // test file
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	good := fmt.Sprintf("%x  steam-runtime.tar.xz\n", md5.Sum([]byte("archive bytes"))) //nolint:gosec // see import
	if err := verifyMD5(f, "steam-runtime.tar.xz", []byte(good)); err != nil {
		t.Fatal(err)
	}
	if err := verifyMD5(f, "steam-runtime.tar.xz", []byte(strings.Replace(good, "a", "b", 1))); err == nil {
		t.Fatal("a wrong checksum verified")
	}
}

func TestSwapKeepsOldTreeAndWritesChecksum(t *testing.T) {
	r := newRuntimeDir(t)
	write(t, filepath.Join(r.Dir, "steam-runtime/version.txt"), "old")
	write(t, filepath.Join(r.Dir, "steam-runtime.old/stale"), "older")
	write(t, filepath.Join(r.Dir, "steam-runtime.tmp/steam-runtime/version.txt"), "new")
	dir, err := openDir(r.Dir)
	if err != nil {
		t.Fatal(err)
	}
	defer unix.Close(dir)
	tmp, err := openBeneath(dir, "steam-runtime.tmp", unix.O_RDONLY|unix.O_DIRECTORY, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer unix.Close(tmp)
	if err := r.swap(dir, tmp, []byte("sum\n")); err != nil {
		t.Fatal(err)
	}
	for path, want := range map[string]string{
		"steam-runtime/version.txt":     "new",
		"steam-runtime/checksum":        "sum\n",
		"steam-runtime.old/version.txt": "old",
	} {
		if b, err := os.ReadFile(filepath.Join(r.Dir, path)); err != nil || string(b) != want {
			t.Fatalf("%s = %q, %v; want %q", path, b, err, want)
		}
	}
	for _, gone := range []string{"steam-runtime.old/stale", "steam-runtime.tmp"} {
		if _, err := os.Lstat(filepath.Join(r.Dir, gone)); !os.IsNotExist(err) {
			t.Fatalf("%s must be gone: %v", gone, err)
		}
	}
}

func TestRemoveAllDoesNotFollowSymlinks(t *testing.T) {
	outside := t.TempDir()
	write(t, filepath.Join(outside, "keep"), "x")
	r := newRuntimeDir(t)
	write(t, filepath.Join(r.Dir, "tree/file"), "x")
	if err := os.Symlink(outside, filepath.Join(r.Dir, "tree/link")); err != nil {
		t.Fatal(err)
	}
	dir, err := openDir(r.Dir)
	if err != nil {
		t.Fatal(err)
	}
	defer unix.Close(dir)
	if err := removeAll(dir, "tree"); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(outside, "keep")); err != nil {
		t.Fatalf("removeAll followed a symlink out of the tree: %v", err)
	}
}
