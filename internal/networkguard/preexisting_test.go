package networkguard

import (
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

func TestCodeVouched(t *testing.T) {
	const user, other, ugrp = 1000, 1001, 1000
	sys := func(mode uint32) fileFacts { return fileFacts{uid: 0, gid: 0, mode: mode, sb: 1, vouched: true} }
	usr := func(mode uint32, sb uint64, vouched bool) fileFacts {
		return fileFacts{uid: user, gid: ugrp, mode: mode, sb: sb, vouched: vouched}
	}
	rootExe := sys(0o100755)
	userExe := usr(0o100755, 2, false)
	for _, tc := range []struct {
		name          string
		lib, dir, exe fileFacts
		want          bool
	}{
		{"system lib", sys(0o100644), sys(0o40755), rootExe, true},
		{"system lib, group root may write", sys(0o100664), sys(0o40775), rootExe, true},
		{"system lib on an unvouched mount", fileFacts{mode: 0o100644, sb: 9}, sys(0o40755), rootExe, false},
		{"root lib in a user-writable dir", sys(0o100644), usr(0o40755, 1, true), rootExe, false},
		{"root lib, other-writable", sys(0o100646), sys(0o40755), rootExe, false},
		{"user lib for a root exe", usr(0o100644, 1, true), usr(0o40755, 1, true), rootExe, false},
		{"user lib for that user's exe, same fs", usr(0o100644, 2, false), usr(0o40755, 2, false), userExe, true},
		{"user lib in a root dir", usr(0o100644, 2, false), sys(0o40755), userExe, true},
		{"user lib on another unvouched fs (FUSE image)", usr(0o100644, 7, false), usr(0o40755, 7, false), userExe, false},
		{"user lib on another vouched fs", usr(0o100644, 3, true), usr(0o40755, 3, true), userExe, true},
		{"another user's lib", fileFacts{uid: other, mode: 0o100644, sb: 2}, usr(0o40755, 2, false), userExe, false},
		{"group-writable lib, exe not", usr(0o100664, 2, false), usr(0o40755, 2, false), userExe, false},
		{"group-writable lib and exe (UPG umask 002)", usr(0o100664, 2, false), usr(0o40775, 2, false),
			usr(0o100775, 2, false), true},
		{"group-writable dir, exe not", usr(0o100644, 2, false), usr(0o40775, 2, false), userExe, false},
		{"other-writable dir", usr(0o100644, 2, false), usr(0o41777, 2, false), userExe, false},
	} {
		if got := codeVouched(tc.lib, tc.dir, tc.exe); got != tc.want {
			t.Errorf("%s: codeVouched = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestLoaderEnv(t *testing.T) {
	for env, want := range map[string]bool{
		"PATH=/bin\x00LD_PRELOAD=/tmp/x.so\x00":      true,
		"LD_AUDIT=/home/u/a.so\x00":                  true,
		"LD_PRELOAD=\x00LD_AUDIT=\x00HOME=/root\x00": false,
		"LD_LIBRARY_PATH=/tmp\x00":                   false, // its libraries are judged when mapped
		"XLD_PRELOAD=/tmp/x.so\x00":                  false,
		"":                                           false,
	} {
		if got := loaderEnv([]byte(env)); got != want {
			t.Errorf("loaderEnv(%q) = %v, want %v", env, got, want)
		}
	}
}

func TestReadMountDevs(t *testing.T) {
	p := filepath.Join(t.TempDir(), "mountinfo")
	mi := "22 1 259:2 / / rw,relatime shared:1 - ext4 /dev/nvme0n1p2 rw\n" +
		"27 22 0:50 / /home/u/fuse rw,relatime shared:32 - fuse.bindfs bindfs rw,user_id=0\n" +
		"40 22 0:31 / /home rw,nosuid,nodev shared:9 - btrfs /dev/sda2 rw,subvol=/home\n"
	if err := os.WriteFile(p, []byte(mi), 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := readMountDevs(p)
	if err != nil {
		t.Fatal(err)
	}
	want := map[uint64]mountDev{
		22: {dev: ebpf.KernelDev(259, 2)},
		27: {dev: ebpf.KernelDev(0, 50), fuse: true},
		40: {dev: ebpf.KernelDev(0, 31)},
	}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for id, m := range want {
		if got[id] != m {
			t.Errorf("mount %d: got %+v, want %+v", id, got[id], m)
		}
	}
	if err := os.WriteFile(p, []byte("22 1 259:2 / / rw shared:1 ext4\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := readMountDevs(p); err == nil {
		t.Fatal("a malformed line must fail the table")
	}
}

func TestReasonText(t *testing.T) {
	for r, want := range map[uint32]string{0: "", reasonCode: "code-suspect", reasonEnv: "ld-preload",
		reasonSuperseded: "superseded", 9: "reason-9"} {
		if got := reasonText(r); got != want {
			t.Errorf("reasonText(%d) = %q, want %q", r, got, want)
		}
	}
}

func TestParseProcs(t *testing.T) {
	b := make([]byte, 2*procRecSize)
	binary.LittleEndian.PutUint32(b[0:], 7)
	binary.LittleEndian.PutUint32(b[4:], 101364)
	binary.LittleEndian.PutUint64(b[8:], 1<<40)
	binary.LittleEndian.PutUint32(b[16:], 1)
	binary.LittleEndian.PutUint32(b[20:], 4200)
	got, err := parseProcs(b)
	if err != nil {
		t.Fatal(err)
	}
	want := []procRec{{nspid: 7, tgid: 101364, start: 1 << 40}, {nspid: 1, tgid: 4200}}
	if !slices.Equal(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
	if _, err := parseProcs(b[:20]); err == nil {
		t.Fatal("a torn record must fail")
	}
}

func TestJudgePasses(t *testing.T) {
	parent := procRec{nspid: 5, tgid: 500, start: 1}
	child := procRec{nspid: 6, tgid: 600, start: 2}  // forked before parent was marked
	reused := procRec{nspid: 6, tgid: 600, start: 3} // the same tgid, a later process
	passes := [][]procRec{{parent}, {parent, child}, {parent, child, reused}, {parent, child, reused}}
	suspect := map[procRec]bool{parent: true, child: true, reused: true}
	var listed int
	marks := map[procRec]int{}
	err := judgePasses(func() ([]procRec, error) {
		listed++
		return passes[min(listed, len(passes))-1], nil
	}, func(p procRec) (bool, error) {
		marks[p]++
		return suspect[p], nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if listed != 4 {
		t.Errorf("passes = %d, want 4 (the last one marks nothing)", listed)
	}
	for _, p := range []procRec{parent, child, reused} {
		if marks[p] != 1 {
			t.Errorf("%+v judged %d times, want once", p, marks[p])
		}
	}

	// A process tree that keeps forking suspect children never converges: refuse to start.
	n := uint32(0)
	err = judgePasses(func() ([]procRec, error) {
		n++
		return []procRec{{nspid: n, tgid: n}}, nil
	}, func(procRec) (bool, error) { return true, nil })
	if err == nil {
		t.Fatal("a non-converging scan must fail")
	}
	if !errors.Is(judgePasses(func() ([]procRec, error) { return nil, os.ErrPermission },
		func(procRec) (bool, error) { return false, nil }), os.ErrPermission) {
		t.Fatal("a failed listing must fail the scan")
	}
}
