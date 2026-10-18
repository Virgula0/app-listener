package editprotected

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Virgula0/app-listener/internal/binledger"
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

// useTempStore points the ledger and the legacy file at a temp dir.
func useTempStore(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	oldLedger, oldLegacy := ledgerPath, legacyHashPath
	ledgerPath, legacyHashPath = filepath.Join(dir, "binaries.db"), filepath.Join(dir, "edit-auth.hash")
	t.Cleanup(func() { ledgerPath, legacyHashPath = oldLedger, oldLegacy })
	return dir
}

func mustHash(t *testing.T, pw string, origin Origin) string {
	t.Helper()
	enc, err := Hash(pw, origin)
	if err != nil {
		t.Fatal(err)
	}
	return enc
}

func TestStoreLoadHash(t *testing.T) {
	useTempStore(t)
	if ok, err := HashConfigured(); err != nil || ok {
		t.Fatalf("fresh store configured=%v err=%v", ok, err)
	}
	if _, err := LoadHash(); !errors.Is(err, ErrNoPassword) {
		t.Fatalf("LoadHash on a fresh store = %v, want ErrNoPassword", err)
	}
	enc := mustHash(t, "Correct-Horse-9", OriginInstall)
	if err := StoreHash(enc); err != nil {
		t.Fatal(err)
	}
	if got, err := LoadHash(); err != nil || got != enc {
		t.Fatalf("LoadHash = %q, %v", got, err)
	}
	if fi, err := os.Stat(ledgerPath); err != nil || fi.Mode().Perm() != 0o600 {
		t.Fatalf("ledger %v %v, want mode 0600", fi, err)
	}
}

func TestStoreHashRefusesMalformed(t *testing.T) {
	useTempStore(t)
	for _, bad := range []string{"", "notahash", "pbkdf2-sha256$1000$AAAA$AAAA$bogus"} {
		if err := StoreHash(bad); err == nil {
			t.Fatalf("StoreHash(%q) accepted a record no password matches", bad)
		}
	}
}

func TestCheckPasswordSharesGuessBudget(t *testing.T) {
	useTempStore(t)
	const pw = "Correct-Horse-9"
	if err := StoreHash(mustHash(t, pw, OriginInstall)); err != nil {
		t.Fatal(err)
	}
	if err := CheckPassword(pw); err != nil {
		t.Fatalf("correct password refused: %v", err)
	}
	for i := range binledger.MaxAuthFailures {
		if err := CheckPassword("Wrong-Horse-99"); !errors.Is(err, ErrWrongPassword) {
			t.Fatalf("guess %d: %v, want ErrWrongPassword", i+1, err)
		}
	}
	var locked *binledger.LockedError
	if err := CheckPassword(pw); !errors.As(err, &locked) {
		t.Fatalf("correct password during the lockout: %v, want a lockout (no oracle while locked)", err)
	}
}

func TestMigrateLegacyHashImports(t *testing.T) {
	useTempStore(t)
	enc := mustHash(t, "Correct-Horse-9", OriginInstall)
	if err := os.WriteFile(legacyHashPath, []byte(enc+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := MigrateLegacyHash(); err != nil {
		t.Fatal(err)
	}
	if got, err := LoadHash(); err != nil || got != enc {
		t.Fatalf("after migration LoadHash = %q, %v", got, err)
	}
	if _, err := os.Lstat(legacyHashPath); !os.IsNotExist(err) {
		t.Fatalf("legacy file left behind: %v", err)
	}
}

func TestMigrateLegacyHashNeverReplaces(t *testing.T) {
	useTempStore(t)
	cur := mustHash(t, "Correct-Horse-9", OriginCLI)
	if err := StoreHash(cur); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(legacyHashPath, []byte(mustHash(t, "Other-Horse-77", OriginInstall)), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := MigrateLegacyHash(); err != nil {
		t.Fatal(err)
	}
	if got, _ := LoadHash(); got != cur {
		t.Fatal("a legacy file replaced the stored password")
	}
}

func TestMigrateLegacyHashEmptyPlaceholder(t *testing.T) {
	useTempStore(t)
	if err := os.WriteFile(legacyHashPath, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := MigrateLegacyHash(); err != nil {
		t.Fatal(err)
	}
	if ok, _ := HashConfigured(); ok {
		t.Fatal("an empty placeholder configured a password")
	}
	if _, err := os.Lstat(legacyHashPath); !os.IsNotExist(err) {
		t.Fatalf("placeholder left behind: %v", err)
	}
}

func TestMigrateLegacyHashRefusesOddFiles(t *testing.T) {
	dir := useTempStore(t)
	target := filepath.Join(dir, "target")
	if err := os.WriteFile(target, []byte(mustHash(t, "Correct-Horse-9", OriginInstall)), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(target, legacyHashPath); err != nil {
		t.Fatal(err)
	}
	if err := MigrateLegacyHash(); err == nil {
		t.Fatal("a symlinked legacy file was imported")
	}
	if err := os.Remove(legacyHashPath); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(legacyHashPath, []byte("garbage"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := MigrateLegacyHash(); err == nil {
		t.Fatal("a malformed legacy record was accepted")
	}
	if _, err := os.Lstat(legacyHashPath); err != nil {
		t.Fatalf("a malformed legacy file must stay for the operator to inspect: %v", err)
	}
	if ok, _ := HashConfigured(); ok {
		t.Fatal("a refused legacy file configured a password")
	}
}
