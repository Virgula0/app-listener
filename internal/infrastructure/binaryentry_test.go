package ebpf

import (
	"crypto/sha256"
	"strings"
	"testing"
)

func TestBinariesSummary(t *testing.T) {
	summary := BinariesSummary([]BinaryEntry{
		{Path: "/usr/bin/cat", Hash: sha256.Sum256([]byte("cat-content")), Comm: "cat"},
		{Path: "/usr/bin/dog", Hash: sha256.Sum256([]byte("dog-content")), Comm: "dog"},
	})
	for _, want := range []string{"/usr/bin/cat", "/usr/bin/dog", "sha256"} {
		if !strings.Contains(summary, want) {
			t.Errorf("summary %q lacks %q", summary, want)
		}
	}
}
