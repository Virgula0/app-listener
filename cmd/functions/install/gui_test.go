package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunWithRefusesDeployWithoutBinary(t *testing.T) {
	err := RunWith(nil, Request{Mode: ModeWizard})
	if err == nil || !strings.Contains(err.Error(), "no headless binary") {
		t.Errorf("wizard without a binary: got %v, want a refusal", err)
	}
}

func TestRunWithRefusesNonHeadlessBinary(t *testing.T) {
	script := filepath.Join(t.TempDir(), "app-listener")
	if err := os.WriteFile(script, []byte("#!/bin/sh\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	err := RunWith(nil, Request{Mode: ModeWizard, Binary: script})
	if err == nil || !strings.Contains(err.Error(), "refusing to deploy") {
		t.Errorf("wizard with a non-Go binary: got %v, want a refusal", err)
	}
}

func TestNeedsBinary(t *testing.T) {
	want := map[Mode]bool{
		ModeWizard: true, ModeUpdateCatalog: false,
		ModeDiffCatalog: false, ModeRestoreBackups: false, ModeDeleteBackups: false,
	}
	for m, w := range want {
		if got := needsBinary(m); got != w {
			t.Errorf("needsBinary(%d) = %v, want %v", m, got, w)
		}
	}
}
