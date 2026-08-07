package editprotected

import (
	"testing"
)

// TestPickEncryptedDirSingleSkipsPicker verifies the single-candidate fast
// path: with only one encrypted directory there is no TUI picker to run
// headlessly, and that directory is returned directly.
func TestPickEncryptedDirSingleSkipsPicker(t *testing.T) {
	chosen, err := pickEncryptedDir([]string{"/vault/one"})
	if err != nil {
		t.Fatalf("pickEncryptedDir: %v", err)
	}
	if chosen != "/vault/one" {
		t.Errorf("chosen = %q, want /vault/one", chosen)
	}
}
