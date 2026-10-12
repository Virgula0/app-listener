package ebpf

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/sys/unix"
)

func TestDurableFS(t *testing.T) {
	for _, typ := range []int64{unix.EXT4_SUPER_MAGIC, unix.XFS_SUPER_MAGIC, unix.BTRFS_SUPER_MAGIC, unix.F2FS_SUPER_MAGIC} {
		if !durableFS(typ) {
			t.Fatalf("%#x must be durable", typ)
		}
	}
	for _, typ := range []int64{unix.FUSE_SUPER_MAGIC, unix.NFS_SUPER_MAGIC, unix.OVERLAYFS_SUPER_MAGIC,
		unix.TMPFS_MAGIC, 0xFE534D42 /* smb2 */, 0xFF534D42 /* cifs */} {
		if durableFS(typ) {
			t.Fatalf("%#x must not be durable", typ)
		}
	}
}

func TestStatDurableRefusesAnotherFingerprint(t *testing.T) {
	p := filepath.Join(t.TempDir(), "bin")
	if err := os.WriteFile(p, []byte("x"), 0o755); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(p)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	fp, err := StatBinaryFile(f)
	if err != nil {
		t.Fatal(err)
	}
	fp.CtimeNs++
	if _, ok := StatDurable(f, fp); ok {
		t.Fatal("a fingerprint the fd no longer shows must not be durable")
	}
}

func TestSyncFileOnReadOnlyFd(t *testing.T) {
	p := filepath.Join(t.TempDir(), "bin")
	if err := os.WriteFile(p, []byte("x"), 0o755); err != nil {
		t.Fatal(err)
	}
	f, err := os.Open(p)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := SyncFile(f); err != nil {
		t.Fatal(err)
	}
}
