package install

import (
	"os"
	"path/filepath"
	"strings"
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

// The migration copies a tree the user owns, as root. An entry classified as a regular file (or
// directory) and then replaced by a symlink before it is opened must not pull content from outside
// the source tree into the destination.
func TestCopyTreeSourceSwappedAfterClassification(t *testing.T) {
	const marker = "OUTSIDE-THE-SOURCE-TREE"
	outside := t.TempDir()
	outsideFile := filepath.Join(outside, "f")
	if err := os.WriteFile(outsideFile, []byte(marker), 0o600); err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name, entry, target string
		mk                  func(string) error
	}{
		{"file", "a", outsideFile, func(p string) error { return os.WriteFile(p, []byte("own"), 0o600) }},
		{"dir", "d", outside, func(p string) error { return os.Mkdir(p, 0o700) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			src := filepath.Join(t.TempDir(), "src")
			dst := filepath.Join(t.TempDir(), "dst")
			if err := os.Mkdir(src, 0o700); err != nil {
				t.Fatal(err)
			}
			entry := filepath.Join(src, tc.entry)
			if err := tc.mk(entry); err != nil {
				t.Fatal(err)
			}

			defer func(prev func(string)) { testHookBeforeSourceOpen = prev }(testHookBeforeSourceOpen)
			testHookBeforeSourceOpen = func(p string) {
				if p != entry {
					return
				}
				if err := os.RemoveAll(entry); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(tc.target, entry); err != nil {
					t.Fatal(err)
				}
			}

			_ = CopyTree(src, dst)

			_ = filepath.WalkDir(dst, func(p string, d os.DirEntry, err error) error {
				if err != nil || d.IsDir() || d.Type()&os.ModeSymlink != 0 {
					return nil
				}
				if data, rerr := os.ReadFile(p); rerr == nil && strings.Contains(string(data), marker) {
					t.Errorf("content from outside the source tree was copied to %s", p)
				}
				return nil
			})
		})
	}
}
