package install

import (
	"os"
	"path/filepath"
	"testing"
)

// The fscrypt migration copies a user's tree back into a freshly created vault while running as
// root and while no guard is attached (the daemon is stopped for the migration). copyDir hands the
// destination directory to the unprivileged user (preserveMeta chowns it to the source owner)
// BEFORE copying any entry into it, so that user can create entries in the destination mid-copy.
// copyFile must therefore never follow a symlink it finds at a destination name.
func TestCopyTreeRefusesPlantedDestinationSymlink(t *testing.T) {
	base := t.TempDir()
	src := filepath.Join(base, "src")
	dst := filepath.Join(base, "dst")
	victim := filepath.Join(base, "victim")

	if err := os.MkdirAll(src, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(src, "payload"), []byte("attacker chosen bytes"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(victim, []byte("VICTIM"), 0o600); err != nil {
		t.Fatal(err)
	}

	// The destination already exists (it is the recreated vault root) and the attacker, who owns
	// it, has planted a symlink under the name the copy is about to write.
	if err := os.MkdirAll(dst, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(victim, filepath.Join(dst, "payload")); err != nil {
		t.Fatal(err)
	}

	err := CopyTree(src, dst)

	if got, rerr := os.ReadFile(victim); rerr != nil || string(got) != "VICTIM" {
		t.Fatalf("migration wrote through a planted destination symlink: victim = %q (err %v)", got, rerr)
	}
	if err == nil {
		t.Fatal("CopyTree followed a planted destination symlink; expected a refusal")
	}
}
