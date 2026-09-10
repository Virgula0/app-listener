package editprotected

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAuditTree(t *testing.T) {
	root := t.TempDir()

	// A clean 0600 file: no finding.
	if err := os.WriteFile(filepath.Join(root, "secret"), []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	// Group/other readable: finding.
	if err := os.WriteFile(filepath.Join(root, "leaky"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Executable: finding.
	if err := os.WriteFile(filepath.Join(root, "helper"), []byte("x"), 0o700); err != nil {
		t.Fatal(err)
	}
	// Symlink escaping the tree: finding.
	if err := os.Symlink("/etc/passwd", filepath.Join(root, "escape")); err != nil {
		t.Fatal(err)
	}

	findings := auditTree(root)
	joined := strings.Join(findings, "\n")

	for _, want := range []string{"leaky", "helper", "escape"} {
		if !strings.Contains(joined, want) {
			t.Errorf("expected a finding mentioning %q, got:\n%s", want, joined)
		}
	}
	if strings.Contains(joined, "secret") {
		t.Errorf("0600 file should not be flagged, got:\n%s", joined)
	}
}

func TestWithin(t *testing.T) {
	if !within("/a/b", "/a/b/c") {
		t.Error("child should be within")
	}
	if !within("/a/b", "/a/b") {
		t.Error("self should be within")
	}
	if within("/a/b", "/a/bc") {
		t.Error("sibling prefix should not be within")
	}
	if within("/a/b", "/a") {
		t.Error("parent should not be within")
	}
}
