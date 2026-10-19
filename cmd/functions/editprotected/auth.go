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
	"time"
	"unicode"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/binledger"
)

// LegacyHashFile is where releases before the binaries.db edit_auth table kept the hash;
// MigrateLegacyHash imports it once.
const LegacyHashFile = "/etc/app-listener/edit-auth.hash"

// ledgerPath and legacyHashPath are binledger.DefaultPath and LegacyHashFile in production; tests
// point them at temp files.
var (
	ledgerPath     = binledger.DefaultPath
	legacyHashPath = LegacyHashFile
)

// Password strength floor: the password is a local, rate-limited, lockout-protected secret (every
// verifier draws on CheckPassword's guess budget), so the bar is "not trivially guessable", not
// "resists an offline crack of the record in binaries.db".
const (
	minPasswordLen = 12
	minCharClasses = 3
	pbkdf2Iters    = 600_000
	pbkdf2KeyLen   = 32
	pbkdf2SaltLen  = 16
)

// Origin records how the password was set: by the installer or by `edit-protected
// password`. Informational only.
type Origin string

const (
	OriginInstall Origin = "install"
	OriginCLI     Origin = "cli"
)

// ErrNoPassword: no edit-protected password configured, so live mode is unavailable and
// edit-protected falls back to the daemon-stopped flow.
var ErrNoPassword = errors.New("no edit-protected password is configured")

// ErrWrongPassword: the password does not match the stored record.
var ErrWrongPassword = errors.New("authentication failed")

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

// Hash derives the stored form of pw: pbkdf2-sha256$<iter>$<b64salt>$<b64dk>$<origin>. base64 is
// raw-url (no padding, no '$').
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

// parsedHash is a decoded password record.
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

// openLedger opens binaries.db; ErrNoPassword when nothing is installed.
func openLedger() (*binledger.Ledger, error) {
	if _, err := os.Stat(filepath.Dir(ledgerPath)); err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNoPassword
		}
		return nil, err
	}
	return binledger.Open(ledgerPath)
}

// HashConfigured reports whether a password is configured.
func HashConfigured() (bool, error) {
	_, err := LoadHash()
	if errors.Is(err, ErrNoPassword) {
		return false, nil
	}
	return err == nil, err
}

// LoadHash returns the stored record, ErrNoPassword when none is set.
func LoadHash() (string, error) {
	l, err := openLedger()
	if err != nil {
		return "", err
	}
	defer l.Close()
	encoded, err := l.EditAuthHash()
	if err != nil {
		return "", fmt.Errorf("reading the edit-protected password: %w", err)
	}
	if strings.TrimSpace(encoded) == "" {
		return "", ErrNoPassword
	}
	return encoded, nil
}

// StoreHash replaces the stored record. binaries.db is sealed while the daemon runs: only this
// binary can write it, so a non-app-listener root process can't install its own password.
func StoreHash(encoded string) error {
	if _, err := parseHash(encoded); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(ledgerPath), 0o755); err != nil {
		return err
	}
	l, err := binledger.Open(ledgerPath)
	if err != nil {
		return err
	}
	defer l.Close()
	return l.SetEditAuthHash(encoded)
}

// CheckPassword verifies pw against the stored record under the persistent guess budget
// (binledger.BeginAuthAttempt): the daemon's control socket and the CLI draw on the same one, so
// neither restarting the daemon nor running the CLI in a loop is an unthrottled oracle.
func CheckPassword(pw string) error {
	l, err := openLedger()
	if err != nil {
		return err
	}
	defer l.Close()
	encoded, err := l.EditAuthHash()
	if err != nil {
		return fmt.Errorf("reading the edit-protected password: %w", err)
	}
	if strings.TrimSpace(encoded) == "" {
		return ErrNoPassword
	}
	if lockErr := l.BeginAuthAttempt(time.Now()); lockErr != nil {
		return lockErr
	}
	ok, err := Verify(encoded, pw)
	if err != nil {
		return fmt.Errorf("verifying the password: %w", err)
	}
	if !ok {
		return ErrWrongPassword
	}
	return l.AuthSucceeded()
}

// MigrateLegacyHash moves a LegacyHashFile record into binaries.db, then removes the file. It is
// imported only while binaries.db holds no password: a stale file never replaces the current one.
// The file is removed even when the import fails (malformed record, DB error): unlike binaries.db
// it is not sealed, so any root process could read and crack it offline.
func MigrateLegacyHash() error {
	fi, err := os.Lstat(legacyHashPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if !fi.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file: not imported", legacyHashPath)
	}
	importErr := importLegacyHash()
	if rmErr := os.Remove(legacyHashPath); rmErr != nil {
		return errors.Join(importErr, fmt.Errorf("removing %s: %w", legacyHashPath, rmErr))
	}
	if importErr == nil {
		return nil
	}
	if ok, _ := HashConfigured(); ok {
		return fmt.Errorf("%w: %s removed, the password in %s stays", importErr, legacyHashPath, ledgerPath)
	}
	return fmt.Errorf("%w: %s removed and no edit-protected password is set: stop the daemon and "+
		"run `app-listener edit-protected password`", importErr, legacyHashPath)
}

func importLegacyHash() error {
	data, err := os.ReadFile(legacyHashPath)
	if err != nil {
		return err
	}
	rec := strings.TrimSpace(string(data))
	if rec == "" {
		return nil
	}
	if _, parseErr := parseHash(rec); parseErr != nil {
		return fmt.Errorf("%s: %w: not imported", legacyHashPath, parseErr)
	}
	cur, err := LoadHash()
	switch {
	case errors.Is(err, ErrNoPassword):
		if storeErr := StoreHash(rec); storeErr != nil {
			return fmt.Errorf("importing %s: %w", legacyHashPath, storeErr)
		}
		log.Infof("edit-protected password imported from %s into %s", legacyHashPath, ledgerPath)
	case err != nil:
		return fmt.Errorf("importing %s: %w", legacyHashPath, err)
	case cur != rec:
		log.Warnf("%s differs from the password in %s: keeping the latter", legacyHashPath, ledgerPath)
	}
	return nil
}
