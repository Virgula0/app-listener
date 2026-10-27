package maintainer

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestPathRecord(t *testing.T) {
	for _, ok := range []string{"", "\n", "/usr/lib/libc.so.6\n/usr/lib/libc-2.39.so\n", "/a/b_c+d@e~f:g,h-i"} {
		if !pathRecord([]byte(ok)) {
			t.Errorf("%q is a path record", ok)
		}
	}
	for _, bad := range []string{"x", "/**/eval(atob('Zm9v'))", "/usr/bin/curl evil", "/a\x00", "/a;b",
		"/a`b`", "/a$(b)", "/a<script>", "/a\n#!/bin/sh"} {
		if pathRecord([]byte(bad)) {
			t.Errorf("%q is no path record", bad)
		}
	}
}

func write(t *testing.T, p, content string, mode os.FileMode) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), mode); err != nil {
		t.Fatal(err)
	}
}

// fixture: a lib_dir holding lib.so and an archive with member.so, and an upper layer the run wrote.
func fixture(t *testing.T) (lib, upper string) {
	t.Helper()
	base := t.TempDir()
	lib = filepath.Join(base, "root", "lib")
	write(t, filepath.Join(lib, "lib.so"), strings.Repeat("L", 9000), 0o755)
	writeTar(t, filepath.Join(lib, "rt.tar"), map[string]string{"rt/member.so": strings.Repeat("M", 7000)})
	upper = filepath.Join(base, "upper")
	return lib, filepath.Join(upper, "lib")
}

func TestScanAndValidateAcceptsLibDirContent(t *testing.T) {
	lib, up := fixture(t)
	write(t, filepath.Join(up, "rt", "member.so"), strings.Repeat("M", 7000), 0o4755)
	write(t, filepath.Join(up, "old", "lib.so"), strings.Repeat("L", 9000), 0o755)
	write(t, filepath.Join(up, "rt", "pins", "system_libc.so.6"), "/usr/lib/libc.so.6\n/usr/lib/libc.so.6.1\n", 0o755)
	write(t, filepath.Join(up, "rt", "pins", "done"), "", 0o644)
	if err := os.Symlink("/usr/lib/libc.so.6", filepath.Join(up, "rt", "pins", "libc.so.6")); err != nil {
		t.Fatal(err)
	}
	ops, err := scanLib(filepath.Dir(up), "lib")
	if err != nil {
		t.Fatal(err)
	}
	if err := validate(lib, ops); err != nil {
		t.Fatal(err)
	}
	got := map[string]op{}
	for _, o := range ops {
		got[o.rel] = o
	}
	if o := got["rt/member.so"]; o.kind != opFile || o.mode != 0o755 || o.data != nil {
		t.Errorf("member: %+v (setuid must be dropped, content copied from the stage)", o)
	}
	if o := got["rt/pins/system_libc.so.6"]; o.data == nil || o.mode != 0o644 {
		t.Errorf("record: %+v (kept in memory, no exec bits)", o)
	}
	if o := got["rt/pins/libc.so.6"]; o.kind != opLink || o.target != "/usr/lib/libc.so.6" {
		t.Errorf("link: %+v", o)
	}
	if o := got["rt"]; o.kind != opDir || !o.names["member.so"] || !o.names["pins"] {
		t.Errorf("dir: %+v", o)
	}
}

func TestValidateRefusesForeignContent(t *testing.T) {
	for name, content := range map[string]string{
		"a library": strings.Repeat("X", 9000),
		"a script":  "/**/fetch('//evil')",
		"a command": "/usr/bin/curl evil.example",
	} {
		lib, up := fixture(t)
		write(t, filepath.Join(up, "rt", "payload"), content, 0o644)
		ops, err := scanLib(filepath.Dir(up), "lib")
		if err != nil {
			t.Fatal(err)
		}
		if err := validate(lib, ops); err == nil {
			t.Errorf("%s the lib_dir does not hold was accepted", name)
		}
	}
}

func TestScanRefusesSpecialFiles(t *testing.T) {
	_, up := fixture(t)
	if err := os.MkdirAll(up, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(filepath.Join(up, "fifo"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := scanLib(filepath.Dir(up), "lib"); err == nil {
		t.Fatal("a FIFO must refuse the run")
	}
}

func TestScanUntouchedLibDir(t *testing.T) {
	_, up := fixture(t)
	ops, err := scanLib(filepath.Dir(up), "lib")
	if err != nil || len(ops) != 0 {
		t.Fatalf("no upper copy means no change: %v %v", ops, err)
	}
}

func TestTouches(t *testing.T) {
	p := &Plan{libs: []libPlan{{dir: "/s/lib", ops: []op{
		{kind: opDelete, rel: "gone"},
		{kind: opDir, rel: "kept"},
		{kind: opDir, rel: "replaced", opaque: true},
		{kind: opFile, rel: "kept/file"},
		{kind: opLink, rel: "linked"},
	}}}}
	for target, want := range map[string]bool{
		"/s/lib/gone/bin": true, "/s/lib/replaced/bin": true, "/s/lib/kept/file": true, "/s/lib/linked/bin": true,
		"/s/lib/kept/other": false, "/s/other/gone": false, "/s/lib/gonex": false,
	} {
		if got := p.Touches(target); got != want {
			t.Errorf("Touches(%s) = %v, want %v", target, got, want)
		}
	}
}

func TestApplyWritesValidatedContent(t *testing.T) {
	lib, up := fixture(t)
	write(t, filepath.Join(lib, "rt", "stale"), "old", 0o644)
	write(t, filepath.Join(lib, "rt", "pins", "libc.so.6"), "was a file", 0o644)
	write(t, filepath.Join(up, "rt", "member.so"), strings.Repeat("M", 7000), 0o755)
	write(t, filepath.Join(up, "rt", "pins", "done"), "", 0o644)
	if err := os.Symlink("/usr/lib/libc.so.6", filepath.Join(up, "rt", "pins", "libc.so.6")); err != nil {
		t.Fatal(err)
	}
	ops, err := scanLib(filepath.Dir(up), "lib")
	if err != nil {
		t.Fatal(err)
	}
	if err := validate(lib, ops); err != nil {
		t.Fatal(err)
	}
	for i := range ops {
		if ops[i].rel == "rt" {
			ops[i].opaque = true // what overlayfs marks a replaced directory
		}
	}
	uid, gid := uint32(os.Getuid()), uint32(os.Getgid()) //nolint:gosec // test ids
	p := &Plan{uid: uid, gid: gid, libs: []libPlan{{dir: lib, ops: ops}}}
	grants := 0
	if err := p.Apply(func(_ string, fn func() error) error { grants++; return fn() }); err != nil {
		t.Fatal(err)
	}
	if grants != 1 {
		t.Errorf("one grant per lib_dir, got %d", grants)
	}
	if b, _ := os.ReadFile(filepath.Join(lib, "rt", "member.so")); string(b) != strings.Repeat("M", 7000) {
		t.Error("member not written")
	}
	if _, err := os.Lstat(filepath.Join(lib, "rt", "stale")); !os.IsNotExist(err) {
		t.Error("a replaced directory keeps only what the run left in it")
	}
	if l, _ := os.Readlink(filepath.Join(lib, "rt", "pins", "libc.so.6")); l != "/usr/lib/libc.so.6" {
		t.Errorf("link: %q", l)
	}
	if fi, err := os.Stat(filepath.Join(lib, "rt", "pins", "done")); err != nil || fi.Size() != 0 {
		t.Errorf("marker: %v", err)
	}
	if left, _ := filepath.Glob(filepath.Join(lib, "rt", ".app-listener-*")); len(left) != 0 {
		t.Errorf("temp files left: %v", left)
	}
}

func TestApplyRefusesAStageChangedAfterValidation(t *testing.T) {
	lib, up := fixture(t)
	src := filepath.Join(up, "rt", "member.so")
	write(t, src, strings.Repeat("M", 7000), 0o755)
	ops, err := scanLib(filepath.Dir(up), "lib")
	if err != nil {
		t.Fatal(err)
	}
	if err := validate(lib, ops); err != nil {
		t.Fatal(err)
	}
	write(t, src, strings.Repeat("X", 7000), 0o755)
	p := &Plan{uid: uint32(os.Getuid()), gid: uint32(os.Getgid()), libs: []libPlan{{dir: lib, ops: ops}}} //nolint:gosec // test ids
	if err := p.Apply(func(_ string, fn func() error) error { return fn() }); err == nil {
		t.Fatal("a staged file rewritten after validation was committed")
	}
	if _, err := os.Lstat(filepath.Join(lib, "rt", "member.so")); !os.IsNotExist(err) {
		t.Error("the rewritten content reached the lib_dir")
	}
	if left, _ := filepath.Glob(filepath.Join(lib, "rt", ".app-listener-*")); len(left) != 0 {
		t.Errorf("temp files left: %v", left)
	}
}

func sum(s string) [32]byte { return sha256.Sum256([]byte(s)) }
