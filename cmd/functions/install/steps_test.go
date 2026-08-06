package install

import (
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/internal/fscrypt"
)

// TestAskEncryptionSkipsNeedEncryptionFalse verifies that a resource
// declared need_encryption: false in the config never triggers the
// encryption question: the directory is not added to toEncrypt, the config
// text is not rewritten, and no interactive prompt is shown (a prompt would
// block on stdin and fail this test instead).
func TestAskEncryptionSkipsNeedEncryptionFalse(t *testing.T) {
	dir := t.TempDir()

	cfgText := "[watch]\npath = " + dir + "\nneed_encryption: false\n"
	cfg, err := validateConfigText(cfgText)
	if err != nil {
		t.Fatalf("parsing config: %v", err)
	}

	text, toEncrypt, err := askEncryption(fscrypt.New(), cfgText, cfg)
	if err != nil {
		t.Fatalf("askEncryption: %v", err)
	}
	if len(toEncrypt) != 0 {
		t.Errorf("toEncrypt = %v, want none for need_encryption: false", toEncrypt)
	}
	if text != cfgText {
		t.Errorf("config text must be left untouched, got:\n%s", text)
	}
	if strings.Contains(text, "need_encryption: true") {
		t.Errorf("config text must keep need_encryption: false, got:\n%s", text)
	}
}
