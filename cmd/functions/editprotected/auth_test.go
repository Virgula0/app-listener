package editprotected

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
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
