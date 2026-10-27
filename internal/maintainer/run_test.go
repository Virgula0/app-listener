package maintainer

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"testing"
)

func TestCopyScriptHoldsOnlyTheConfirmedContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "steam.sh")
	write(t, path, "#!/bin/bash\necho ok\n", 0o755)
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	for _, tc := range []struct {
		sum [32]byte
		ok  bool
	}{{sha256.Sum256([]byte("#!/bin/bash\necho ok\n")), true}, {sha256.Sum256([]byte("#!/bin/bash\necho evil\n")), false}} {
		st := &stage{script: filepath.Join(t.TempDir(), "script")}
		err := st.copyScript(&Spec{Script: f, ScriptSum: tc.sum, ScriptPath: path})
		if (err == nil) != tc.ok {
			t.Fatalf("confirmed=%v: %v", tc.ok, err)
		}
		if b, _ := os.ReadFile(st.script); tc.ok && string(b) != "#!/bin/bash\necho ok\n" {
			t.Fatalf("copy holds %q", b)
		}
	}
}
