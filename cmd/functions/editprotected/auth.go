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
func HashFileExists() (bool, error) {
	_, err := os.Stat(hashFilePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// LoadHashFile reads the stored record, mapping absence to ErrNoHashFile.
func LoadHashFile() (string, error) {
	data, err := os.ReadFile(hashFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", ErrNoHashFile
		}
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// WriteHashFile atomically writes encoded to HashFile as 0600 root-owned
// (temp sibling + fsync + rename), mirroring fscrypt.GenerateMasterKey. The
// caller is responsible for reloading the daemon so its self-guard picks up
// the (possibly new) inode.
func WriteHashFile(encoded string) error {
	dir := filepath.Dir(hashFilePath)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("creating %s: %w", dir, err)
	}
	tmp, err := os.CreateTemp(dir, ".edit-auth.hash.*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = tmp.Close()
			_ = os.Remove(tmpPath)
		}
	}()
	if err := tmp.Chmod(0o600); err != nil {
		return err
	}
	if _, err := tmp.WriteString(encoded + "\n"); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	cleanup = false
	return os.Rename(tmpPath, hashFilePath)
}

// RemoveHashFile deletes the hash file if present; a missing file is success.
func RemoveHashFile() error {
	if err := os.Remove(hashFilePath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
