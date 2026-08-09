package systemd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

func TestVerifyResourcesLocked(t *testing.T) {
	provisioned := map[string]bool{
		"/encrypted-unlocked": true,
		"/encrypted-locked":   false,
	}
	provisionedFn := func(path string) (bool, error) {
		return provisioned[path], nil
	}

	// A locked tree is fine.
	locked := []daemonconfig.Resource{
		{Path: "/encrypted-locked", NeedEncryption: true},
		{Path: "/plain", NeedEncryption: false},
	}
	if err := VerifyResourcesLocked(locked, provisionedFn); err != nil {
		t.Fatalf("VerifyResourcesLocked on locked resources: %v", err)
	}

	// An unlocked (still provisioned) tree is a fatal error.
	unlocked := []daemonconfig.Resource{
		{Path: "/encrypted-unlocked", NeedEncryption: true},
	}
	if err := VerifyResourcesLocked(unlocked, provisionedFn); err == nil {
		t.Fatal("VerifyResourcesLocked accepted an unlocked encrypted resource")
	}
}

func TestEnsureAndRemoveBinSymlink(t *testing.T) {
	dir := t.TempDir()
	binPath := filepath.Join(dir, "sbin", "app-listener")
	if err := os.MkdirAll(filepath.Dir(binPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(binPath, []byte("#!/bin/sh\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	linkPath := filepath.Join(dir, "bin", "app-listener")
	if err := os.MkdirAll(filepath.Dir(linkPath), 0o755); err != nil {
		t.Fatal(err)
	}

	// Missing link: created.
	if err := EnsureBinSymlinkAt(linkPath, binPath); err != nil {
		t.Fatalf("EnsureBinSymlinkAt created the link but reported: %v", err)
	}
	if target, err := os.Readlink(linkPath); err != nil || target != binPath {
		t.Fatalf("link target = %q (err %v), want %q", target, err, binPath)
	}

	// Already correct: no error, link untouched.
	if err := EnsureBinSymlinkAt(linkPath, binPath); err != nil {
		t.Fatalf("EnsureBinSymlinkAt on an existing correct link: %v", err)
	}

	// Foreign symlink: never replaced.
	if err := os.Remove(linkPath); err != nil {
		t.Fatal(err)
	}
	foreign := filepath.Join(dir, "foreign")
	if err := os.WriteFile(foreign, []byte("#!/bin/sh\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(foreign, linkPath); err != nil {
		t.Fatal(err)
	}
	if err := EnsureBinSymlinkAt(linkPath, binPath); err != nil {
		t.Fatalf("EnsureBinSymlinkAt on a foreign symlink: %v", err)
	}
	if target, _ := os.Readlink(linkPath); target != foreign {
		t.Fatalf("foreign symlink was replaced (target %q)", target)
	}

	// Regular file: never replaced.
	if err := os.Remove(linkPath); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(linkPath, []byte("not a link"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := EnsureBinSymlinkAt(linkPath, binPath); err != nil {
		t.Fatalf("EnsureBinSymlinkAt on a regular file: %v", err)
	}
	if info, err := os.Lstat(linkPath); err != nil || !info.Mode().IsRegular() {
		t.Fatalf("regular file was replaced (mode %v, err %v)", info.Mode(), err)
	}

	// RemoveBinSymlinkAt only removes a symlink pointing at binPath: the
	// regular file is left alone.
	RemoveBinSymlinkAt(linkPath, binPath)
	if _, err := os.Lstat(linkPath); err != nil {
		t.Fatalf("RemoveBinSymlinkAt removed the regular file: %v", err)
	}
	if err := os.Remove(linkPath); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(binPath, linkPath); err != nil {
		t.Fatal(err)
	}
	RemoveBinSymlinkAt(linkPath, binPath)
	if _, err := os.Lstat(linkPath); !os.IsNotExist(err) {
		t.Fatalf("RemoveBinSymlinkAt did not remove the link (err %v)", err)
	}
	// Missing link is a no-op.
	RemoveBinSymlinkAt(linkPath, binPath)
	// Foreign link is kept.
	if err := os.Symlink(foreign, linkPath); err != nil {
		t.Fatal(err)
	}
	RemoveBinSymlinkAt(linkPath, binPath)
	if target, err := os.Readlink(linkPath); err != nil || target != foreign {
		t.Fatalf("RemoveBinSymlinkAt removed a foreign link (target %q, err %v)", target, err)
	}
}

func TestEnsureBinSymlinkMissingParent(t *testing.T) {
	dir := t.TempDir()
	// The parent directory does not exist: only a warning, no error.
	if err := EnsureBinSymlinkAt(filepath.Join(dir, "missing", "app-listener"), filepath.Join(dir, "sbin", "app-listener")); err != nil {
		t.Fatalf("EnsureBinSymlinkAt with a missing parent directory: %v", err)
	}
}

func TestReplaceInstalledBinary(t *testing.T) {
	dir := t.TempDir()
	dst := filepath.Join(dir, "app-listener")
	if err := os.WriteFile(dst, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(dir, "new-binary")
	if err := os.WriteFile(src, []byte("new contents"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := ReplaceInstalledBinary(src, dst); err != nil {
		t.Fatalf("ReplaceInstalledBinary: %v", err)
	}
	got, err := os.ReadFile(dst)
	if err != nil || string(got) != "new contents" {
		t.Fatalf("dst content = %q (err %v), want %q", got, err, "new contents")
	}
	info, err := os.Stat(dst)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0o700 {
		t.Fatalf("dst mode = %o, want 700", perm)
	}
}
