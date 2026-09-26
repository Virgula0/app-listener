package guard

import (
	"os"
	"path/filepath"
	"testing"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// Regression: a lib_dir whose directory held only dangling symlinks (Steam
// pressure-vessel `overrides/.../aliases`, pointing at container-only paths)
// failed populate with "all entries failed to stat … fscrypt-encrypted and
// locked", crash-looping the daemon.
func TestReadOnlyWalkToleratesDanglingSymlinks(t *testing.T) {
	root := t.TempDir()
	aliases := filepath.Join(root, "overrides", "aliases")
	if err := os.MkdirAll(aliases, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, n := range []string{"libbz2.so.1", "libcurl-gnutls.so.3", "libOSMesa.so.6"} {
		if err := os.Symlink("/nonexistent/container-only/"+n, filepath.Join(aliases, n)); err != nil {
			t.Fatal(err)
		}
	}
	g := &Guard{mode: ModeReadOnly}
	add := func(p string) error { _, _, err := g.treeInode(p); return err }
	if err := walkInodes(root, true, 0, 0, add); err != nil {
		t.Fatalf("read-only walk over dangling symlinks must succeed, got: %v", err)
	}
}

// A live in-tree link must contribute the LINK's inode, never the host
// library it points at — otherwise the read-only guard would gate writes to
// host system libraries.
func TestReadOnlyTreeInodeDoesNotFollowSymlinks(t *testing.T) {
	outside := filepath.Join(t.TempDir(), "libhost.so")
	if err := os.WriteFile(outside, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(t.TempDir(), "libhost.so")
	if err := os.Symlink(outside, link); err != nil {
		t.Fatal(err)
	}
	_, targetIno, _ := ebpf.StatInode(outside)

	_, roIno, err := (&Guard{mode: ModeReadOnly}).treeInode(link)
	if err != nil {
		t.Fatal(err)
	}
	if roIno == targetIno {
		t.Errorf("read-only guard recorded the symlink TARGET inode (host library) instead of the link")
	}
	// Whitelist (secret) guards keep the historical follow semantics.
	_, wlIno, err := (&Guard{mode: ModeWhitelist}).treeInode(link)
	if err != nil {
		t.Fatal(err)
	}
	if wlIno != targetIno {
		t.Errorf("whitelist guard behavior changed: want target inode %d, got %d", targetIno, wlIno)
	}
}
