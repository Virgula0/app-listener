package editprotected

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

// HashFile is where the installer (and `edit-protected --set-password`)
// persists the PBKDF2 hash of the edit-protected authentication password. It
// lives next to fscrypt.key in the daemon's self-protected config directory:
// the running daemon guards it whitelist-empty (readable/writable only by the
// app-listener binary — see selfProtectSpecs in the daemon package), and on
// disk it is 0600 root-owned.
const HashFile = "/etc/app-listener/edit-auth.hash"

// hashFilePath is HashFile in production; tests point it at a temp file.
var hashFilePath = HashFile

// Password strength floor. The password is a local, rate-limited,
// lockout-protected secret (the daemon throttles guesses over the control
// socket), so the bar is "not trivially guessable", not "resists an offline
// crack of the 0600 root-only hash file".
const (
	minPasswordLen = 12
	minCharClasses = 3
	pbkdf2Iters    = 600_000
	pbkdf2KeyLen   = 32
	pbkdf2SaltLen  = 16
)

// Origin records how the hash file was created. `edit-protected
// --set-password` is allowed to rotate a cli-set password but refuses to
// touch one chosen during `app-listener install` (rotating that one requires
// re-running the installer).
type Origin string

const (
	OriginInstall Origin = "install"
	OriginCLI     Origin = "cli"
)

// ErrNoHashFile means no edit-protected password has been configured: live
// mode is unavailable and edit-protected falls back to the daemon-stopped
// flow.
var ErrNoHashFile = errors.New("no edit-protected password is configured")

// ValidatePassword enforces the strength floor: length, character-class
// diversity, and no single repeated character.
func ValidatePassword(pw string) error {
	if len(pw) < minPasswordLen {
		return fmt.Errorf("password too short: need at least %d characters", minPasswordLen)
	}
	var hasLower, hasUpper, hasDigit, hasSymbol bool
	allSame := true
	for i, r := range pw {
		if i > 0 && byte(r) != pw[0] {
			allSame = false
		}
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsSpace(r):
			return errors.New("password must not contain whitespace")
		default:
			hasSymbol = true
		}
	}
	if allSame {
		return errors.New("password is a single repeated character")
	}
	classes := 0
	for _, ok := range []bool{hasLower, hasUpper, hasDigit, hasSymbol} {
		if ok {
			classes++
		}
	}
	if classes < minCharClasses {
		return fmt.Errorf("password not complex enough: use at least %d of lowercase, uppercase, digits, symbols", minCharClasses)
	}
	return nil
}

// Hash derives the stored representation of pw:
//
//	pbkdf2-sha256$<iter>$<b64salt>$<b64dk>$<origin>
//
// base64 is raw-url (no padding, no '$'). The origin trailer lets
// --set-password distinguish an installer-chosen password from a
// cli-chosen one.
func Hash(pw string, origin Origin) (string, error) {
	if err := ValidatePassword(pw); err != nil {
		return "", err
	}
	salt := make([]byte, pbkdf2SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("generating salt: %w", err)
	}
	dk, err := pbkdf2.Key(sha256.New, pw, salt, pbkdf2Iters, pbkdf2KeyLen)
	if err != nil {
		return "", fmt.Errorf("deriving key: %w", err)
	}
	enc := base64.RawURLEncoding
	return strings.Join([]string{
		"pbkdf2-sha256",
		strconv.Itoa(pbkdf2Iters),
		enc.EncodeToString(salt),
		enc.EncodeToString(dk),
		string(origin),
	}, "$"), nil
}

// parsedHash is a decoded HashFile record.
type parsedHash struct {
	iter   int
	salt   []byte
	dk     []byte
	origin Origin
}

func parseHash(encoded string) (parsedHash, error) {
	parts := strings.Split(strings.TrimSpace(encoded), "$")
	if len(parts) != 5 || parts[0] != "pbkdf2-sha256" {
		return parsedHash{}, errors.New("malformed hash record")
	}
	iter, err := strconv.Atoi(parts[1])
	if err != nil || iter < 1 {
		return parsedHash{}, errors.New("malformed hash record: iteration count")
	}
	enc := base64.RawURLEncoding
	salt, err := enc.DecodeString(parts[2])
	if err != nil {
		return parsedHash{}, errors.New("malformed hash record: salt")
	}
	dk, err := enc.DecodeString(parts[3])
	if err != nil {
		return parsedHash{}, errors.New("malformed hash record: digest")
	}
	origin := Origin(parts[4])
	if origin != OriginInstall && origin != OriginCLI {
		return parsedHash{}, fmt.Errorf("malformed hash record: unknown origin %q", parts[4])
	}
	return parsedHash{iter: iter, salt: salt, dk: dk, origin: origin}, nil
}

// Verify reports whether pw matches the stored record. The comparison is
// constant-time; a malformed record is an error, never a silent mismatch.
func Verify(encoded, pw string) (bool, error) {
	p, err := parseHash(encoded)
	if err != nil {
		return false, err
	}
	got, err := pbkdf2.Key(sha256.New, pw, p.salt, p.iter, len(p.dk))
	if err != nil {
		return false, fmt.Errorf("deriving key: %w", err)
	}
	return subtle.ConstantTimeCompare(got, p.dk) == 1, nil
}

// OriginOf returns how the stored record was created.
func OriginOf(encoded string) (Origin, error) {
	p, err := parseHash(encoded)
	if err != nil {
		return "", err
	}
	return p.origin, nil
}

// HashFileExists reports whether an edit-protected password is configured.
// The daemon pre-creates HashFile as an empty placeholder before it ever
// self-guards it (see cmd/functions/daemon/selfguards.go
// ensureHashFilePlaceholder), so existence alone no longer means
// "configured" — an empty (or whitespace-only) file means "not configured",
// same as the file being absent entirely.
func HashFileExists() (bool, error) {
	data, err := os.ReadFile(hashFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return strings.TrimSpace(string(data)) != "", nil
}

// LoadHashFile reads the stored record, mapping absence OR an empty
// placeholder to ErrNoHashFile.
func LoadHashFile() (string, error) {
	data, err := os.ReadFile(hashFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", ErrNoHashFile
		}
		return "", err
	}
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" {
		return "", ErrNoHashFile
	}
	return trimmed, nil
}

// WriteHashFile writes encoded to HashFile IN PLACE, on the file's existing
// inode (open, truncate, write, fsync) — never by creating a new directory
// entry and renaming it over. That distinction matters here specifically:
// while the daemon runs, HashFile sits inside the ReadOnly-guarded
// /etc/app-listener (see selfProtectSpecs), whose self-allow for the daemon's
// own binary covers rewriting an EXISTING entry but not creating a new one.
// The daemon guarantees HashFile always exists (as an empty placeholder if
// no password is set — see ensureHashFilePlaceholder) before that guard ever
// attaches, so this only ever needs the in-place path; the create branch
// below is a fallback for the case nothing has bootstrapped it yet (no
// daemon has run since /etc/app-listener was created).
func WriteHashFile(encoded string) error {
	dir := filepath.Dir(hashFilePath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("creating %s: %w", dir, err)
	}
	f, err := os.OpenFile(hashFilePath, os.O_RDWR, 0o600)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("opening %s: %w", hashFilePath, err)
		}
		f, err = os.OpenFile(hashFilePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err != nil {
			return fmt.Errorf("creating %s: %w", hashFilePath, err)
		}
	}
	defer f.Close()
	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("truncating %s: %w", hashFilePath, err)
	}
	if _, err := f.WriteAt([]byte(encoded+"\n"), 0); err != nil {
		return fmt.Errorf("writing %s: %w", hashFilePath, err)
	}
	if err := f.Chmod(0o600); err != nil {
		return fmt.Errorf("chmod %s: %w", hashFilePath, err)
	}
	return f.Sync()
}

// RemoveHashFile clears the stored password by truncating HashFile to empty
// IN PLACE — see WriteHashFile for why this never unlinks the file (removing
// then recreating it would need a new directory entry in the ReadOnly-
// guarded /etc/app-listener, which the daemon's own self-allow does not
// permit while it is running). A missing file is treated as already-cleared.
func RemoveHashFile() error {
	f, err := os.OpenFile(hashFilePath, os.O_WRONLY, 0)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("truncating %s: %w", hashFilePath, err)
	}
	return f.Sync()
}
