package guard

import (
	"testing"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

func TestLocalizePaths(t *testing.T) {
	tests := []struct {
		name               string
		mount              *ebpf.MountInfo
		typ                ebpf.EventType
		path, dest         string
		wantPath, wantDest string
	}{
		{name: "separate /home partition", mount: &ebpf.MountInfo{Root: "/", Mountpoint: "/home"},
			typ: ebpf.EventOpen, path: "/u/.ssh/id_ed25519", wantPath: "/home/u/.ssh/id_ed25519"},
		{name: "btrfs subvolume", mount: &ebpf.MountInfo{Root: "/@home", Mountpoint: "/home"},
			typ: ebpf.EventStat, path: "/@home/u/.gnupg", wantPath: "/home/u/.gnupg"},
		{name: "rename dest walked from a dentry", mount: &ebpf.MountInfo{Root: "/", Mountpoint: "/home"},
			typ: ebpf.EventRename, path: "/u/a", dest: "/u/b", wantPath: "/home/u/a", wantDest: "/home/u/b"},
		{name: "symlink dest is the link body", mount: &ebpf.MountInfo{Root: "/", Mountpoint: "/home"},
			typ: ebpf.EventSymlink, path: "/u/lnk", dest: "/etc/passwd", wantPath: "/home/u/lnk", wantDest: "/etc/passwd"},
		{name: "outside the mounted part stays raw", mount: &ebpf.MountInfo{Root: "/@home", Mountpoint: "/home"},
			typ: ebpf.EventOpen, path: "/@other/f", wantPath: "/@other/f"},
		{name: "root filesystem unchanged", mount: &ebpf.MountInfo{Root: "/", Mountpoint: "/"},
			typ: ebpf.EventOpen, path: "/etc/x", wantPath: "/etc/x"},
		{name: "unknown mount unchanged", typ: ebpf.EventOpen, path: "/u/f", wantPath: "/u/f"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Guard{}
			g.rootMount.Store(tt.mount)
			ev := &GuardEvent{FileEvent: ebpf.FileEvent{Type: tt.typ, Path: tt.path, Dest: tt.dest}}
			g.localizePaths(ev)
			if ev.Path != tt.wantPath || ev.Dest != tt.wantDest {
				t.Fatalf("got (%q, %q), want (%q, %q)", ev.Path, ev.Dest, tt.wantPath, tt.wantDest)
			}
		})
	}
}

func TestNoteRootMountFollowsMountID(t *testing.T) {
	root, err := openRoot(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer root.close()

	g := &Guard{path: root.resolved}
	g.noteRootMount(root)
	m := g.rootMount.Load()
	if m == nil || m.ID != root.mnt {
		t.Fatalf("rootMount = %+v, want mount id %d", m, root.mnt)
	}
	if _, ok := ebpf.PathUnder(root.resolved, m.Mountpoint); !ok {
		t.Fatalf("root %s is not under its mount %s", root.resolved, m.Mountpoint)
	}

	stale := &ebpf.MountInfo{ID: root.mnt + 1}
	g.rootMount.Store(stale)
	g.noteRootMount(root)
	if g.rootMount.Load() == stale {
		t.Fatal("a changed mount id must re-read mountinfo")
	}
}
