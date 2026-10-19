package fscrypt

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeFile(t *testing.T, p string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestRefuseOutsideHardlinks_InTreeLinksAccepted(t *testing.T) {
	root := filepath.Join(t.TempDir(), "vault")
	writeFile(t, filepath.Join(root, "a"))
	writeFile(t, filepath.Join(root, "sub", "plain"))
	if err := os.Link(filepath.Join(root, "a"), filepath.Join(root, "sub", "b")); err != nil {
		t.Fatal(err)
	}
	if err := refuseOutsideHardlinks(root); err != nil {
		t.Fatalf("every link is inside the tree: %v", err)
	}
}

func TestRefuseOutsideHardlinks_OutsideLinkRefused(t *testing.T) {
	base := t.TempDir()
	root := filepath.Join(base, "vault")
	secret := filepath.Join(root, "deep", "secret")
	writeFile(t, secret)
	if err := os.Link(secret, filepath.Join(base, "stolen")); err != nil {
		t.Fatal(err)
	}
	err := refuseOutsideHardlinks(root)
	if err == nil || !strings.Contains(err.Error(), secret) {
		t.Fatalf("want a refusal naming %s, got %v", secret, err)
	}
}

func TestRefuseOutsideHardlinks_FileRoot(t *testing.T) {
	base := t.TempDir()
	secret := filepath.Join(base, "secret")
	writeFile(t, secret)
	if err := refuseOutsideHardlinks(secret); err != nil {
		t.Fatalf("single link: %v", err)
	}
	if err := os.Link(secret, filepath.Join(base, "other")); err != nil {
		t.Fatal(err)
	}
	if err := refuseOutsideHardlinks(secret); err == nil {
		t.Fatal("a file root with a second name must be refused")
	}
}

func TestRefuseOutsideHardlinks_MissingRootFails(t *testing.T) {
	if err := refuseOutsideHardlinks(filepath.Join(t.TempDir(), "gone")); err == nil {
		t.Fatal("an unreadable root must fail closed")
	}
}
