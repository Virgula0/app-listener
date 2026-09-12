package editprotected

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
)

func TestValidatePassword(t *testing.T) {
	cases := []struct {
		name string
		pw   string
		ok   bool
	}{
		{"good", "Correct-Horse-9", true},
		{"good symbols+digits+case", "aB3$aB3$aB3$", true},
		{"too short", "aB3$aB3", false},
		{"only two classes", "abcdefghijklMNOP", false},
		{"whitespace", "Correct Horse 9x", false},
		{"all same", "aaaaaaaaaaaa", false},
		{"three classes ok", "abcdefghij1K", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ValidatePassword(c.pw)
			if c.ok && err != nil {
				t.Fatalf("expected valid, got %v", err)
			}
			if !c.ok && err == nil {
				t.Fatalf("expected invalid, got nil")
			}
		})
	}
}

func TestHashVerifyRoundTrip(t *testing.T) {
	const pw = "Correct-Horse-9"
	enc, err := Hash(pw, OriginInstall)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(enc, "$") != 4 {
		t.Fatalf("unexpected encoding %q", enc)
	}

	ok, err := Verify(enc, pw)
	if err != nil || !ok {
		t.Fatalf("Verify(correct) = %v, %v", ok, err)
	}
	ok, err = Verify(enc, "Correct-Horse-8")
	if err != nil || ok {
		t.Fatalf("Verify(wrong) = %v, %v", ok, err)
	}

	origin, err := OriginOf(enc)
	if err != nil || origin != OriginInstall {
		t.Fatalf("OriginOf = %v, %v", origin, err)
	}
}

func TestVerifyMalformed(t *testing.T) {
	for _, bad := range []string{
		"", "notahash", "pbkdf2-sha256$abc$def$ghi$install",
		"pbkdf2-sha256$1000$@@@$@@@$install",
		"pbkdf2-sha256$1000$AAAA$AAAA$bogus",
	} {
		if _, err := Verify(bad, "whatever"); err == nil {
			t.Fatalf("Verify(%q) should error", bad)
		}
	}
}

func TestHashFileLifecycle(t *testing.T) {
	dir := t.TempDir()
	old := hashFilePath
	hashFilePath = filepath.Join(dir, "edit-auth.hash")
	t.Cleanup(func() { hashFilePath = old })

	if ex, _ := HashFileExists(); ex {
		t.Fatal("hash file should not exist yet")
	}
	if _, err := LoadHashFile(); !errors.Is(err, ErrNoHashFile) {
		t.Fatalf("LoadHashFile before write = %v, want ErrNoHashFile", err)
	}

	enc, err := Hash("Correct-Horse-9", OriginCLI)
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteHashFile(enc); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(hashFilePath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("hash file mode = %o, want 600", info.Mode().Perm())
	}

	got, err := LoadHashFile()
	if err != nil || got != enc {
		t.Fatalf("LoadHashFile = %q, %v", got, err)
	}

	if err := RemoveHashFile(); err != nil {
		t.Fatal(err)
	}
	if ex, _ := HashFileExists(); ex {
		t.Fatal("hash file should be gone")
	}
	if err := RemoveHashFile(); err != nil {
		t.Fatalf("RemoveHashFile on missing file = %v", err)
	}
}

func hashFileInode(t *testing.T) uint64 {
	t.Helper()
	info, err := os.Lstat(hashFilePath)
	if err != nil {
		t.Fatalf("lstat %s: %v", hashFilePath, err)
	}
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		t.Fatal("Stat_t unavailable on this platform")
	}
	return st.Ino
}

// TestWriteHashFileInPlaceOverPlaceholder is the regression test for the
// daemon-self-denial bug: WriteHashFile must rewrite an EXISTING file (the
// daemon's zero-length bootstrap placeholder — see
// cmd/functions/daemon/selfguards.go ensureHashFilePlaceholder) on its SAME
// inode, never by creating a new directory entry and renaming it over —
// while the daemon runs, that parent directory is guarded ReadOnly and does
// not permit the self binary to create a new entry there, only to rewrite
// an existing one.
func TestWriteHashFileInPlaceOverPlaceholder(t *testing.T) {
	dir := t.TempDir()
	old := hashFilePath
	hashFilePath = filepath.Join(dir, "edit-auth.hash")
	t.Cleanup(func() { hashFilePath = old })

	// Simulate the daemon's bootstrap: an empty placeholder already exists.
	if err := os.WriteFile(hashFilePath, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if ex, err := HashFileExists(); err != nil || ex {
		t.Fatalf("an empty placeholder must read as not-configured: exists=%v err=%v", ex, err)
	}
	before := hashFileInode(t)

	enc, err := Hash("Correct-Horse-9", OriginCLI)
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteHashFile(enc); err != nil {
		t.Fatalf("WriteHashFile over an existing placeholder: %v", err)
	}
	if got := hashFileInode(t); got != before {
		t.Errorf("WriteHashFile changed the inode: %d -> %d (a live daemon's self-guard would now be tracking the WRONG inode)", before, got)
	}
	got, err := LoadHashFile()
	if err != nil || got != enc {
		t.Fatalf("LoadHashFile after in-place write = %q, %v", got, err)
	}

	// A second write (rotation) must also stay on the same inode.
	enc2, err := Hash("Another-Correct-Horse-9", OriginCLI)
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteHashFile(enc2); err != nil {
		t.Fatal(err)
	}
	if got := hashFileInode(t); got != before {
		t.Errorf("second WriteHashFile changed the inode: %d -> %d", before, got)
	}
}

// TestRemoveHashFileInPlace: clearing the password truncates in place
// (same inode) rather than unlinking — for the same reason WriteHashFile
// no longer renames: a live daemon's ReadOnly self-guard on the parent
// directory does not permit creating a replacement entry later.
func TestRemoveHashFileInPlace(t *testing.T) {
	dir := t.TempDir()
	old := hashFilePath
	hashFilePath = filepath.Join(dir, "edit-auth.hash")
	t.Cleanup(func() { hashFilePath = old })

	enc, err := Hash("Correct-Horse-9", OriginCLI)
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteHashFile(enc); err != nil {
		t.Fatal(err)
	}
	before := hashFileInode(t)

	if err := RemoveHashFile(); err != nil {
		t.Fatal(err)
	}
	if got := hashFileInode(t); got != before {
		t.Errorf("RemoveHashFile changed the inode: %d -> %d", before, got)
	}
	if _, err := os.Lstat(hashFilePath); err != nil {
		t.Errorf("RemoveHashFile must not unlink the file, got: %v", err)
	}
	if ex, _ := HashFileExists(); ex {
		t.Error("hash file should read as not-configured after remove")
	}
}
