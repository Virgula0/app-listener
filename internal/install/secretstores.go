package install

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

const ipcDoc = "docs/limitations.md#secrets-served-over-ipc"

// keyringCredHelpers are docker credential helpers backed by the Secret Service keyring.
var keyringCredHelpers = []string{"secretservice", "desktop"}

// SecretStoreNotes lists u's tokens kept where the guard can't judge the caller: the Secret Service
// keyring, or a git credential helper, both answering any of u's programs. A file-storage hint is
// given only when that file sits in a tree cfg guards. Config files are parsed, never the tools run:
// the installer is root.
func SecretStoreNotes(cfg *daemonconfig.Config, u User) []string {
	covered := CoveredPaths(cfg)
	var notes []string
	ghDir := filepath.Join(u.Home, ".config", "gh")
	if PathCovered(ghDir, covered) {
		if hosts, err := readOwnedRegular(filepath.Join(ghDir, "hosts.yml"), u.UID); err == nil &&
			!strings.Contains(hosts, "oauth_token:") && strings.Contains(hosts, "user:") {
			notes = append(notes, "gh keeps its token in the keyring, which any of your programs can ask: "+
				"`gh auth login --insecure-storage` moves it to the guarded ~/.config/gh/hosts.yml")
		}
	}
	dockerDir := filepath.Join(u.Home, ".docker")
	if PathCovered(dockerDir, covered) {
		if helper := dockerKeyringHelper(filepath.Join(dockerDir, "config.json"), u.UID); helper != "" {
			notes = append(notes, fmt.Sprintf("docker keeps registry logins in the keyring "+
				"(docker-credential-%s): drop credsStore/credHelpers from ~/.docker/config.json and log in "+
				"again to keep them in the guarded file", helper))
		}
	}
	if helpers := gitCredentialHelpers(u); len(helpers) > 0 {
		notes = append(notes, fmt.Sprintf("git's credential helper (%s) answers `git credential fill` "+
			"for any of your programs; no file guard changes that", strings.Join(helpers, ", ")))
	}
	return notes
}

// WarnSecretStores logs SecretStoreNotes for each of users.
func WarnSecretStores(cfg *daemonconfig.Config, users []User) {
	for i := range users {
		for _, n := range SecretStoreNotes(cfg, users[i]) {
			log.Warnf("%s: %s (%s)", users[i].Name, n, ipcDoc)
		}
	}
}

func dockerKeyringHelper(path string, uid uint32) string {
	text, err := readOwnedRegular(path, uid)
	if err != nil {
		return ""
	}
	var conf struct {
		CredsStore  string            `json:"credsStore"`
		CredHelpers map[string]string `json:"credHelpers"`
	}
	if json.Unmarshal([]byte(text), &conf) != nil {
		return ""
	}
	if slices.Contains(keyringCredHelpers, conf.CredsStore) {
		return conf.CredsStore
	}
	for _, h := range conf.CredHelpers {
		if slices.Contains(keyringCredHelpers, h) {
			return h
		}
	}
	return ""
}

// gitCredentialHelpers lists the `helper` values of every [credential] section in u's global git
// configs. Includes are not followed.
func gitCredentialHelpers(u User) []string {
	var out []string
	for _, p := range []string{filepath.Join(u.Home, ".gitconfig"), filepath.Join(u.Home, ".config", "git", "config")} {
		text, err := readOwnedRegular(p, u.UID)
		if err != nil {
			continue
		}
		inCredential := false
		for line := range strings.Lines(text) {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "[") {
				inCredential = strings.HasPrefix(strings.ToLower(line), "[credential")
				continue
			}
			key, val, ok := strings.Cut(line, "=")
			if !inCredential || !ok || !strings.EqualFold(strings.TrimSpace(key), "helper") {
				continue
			}
			if v := strings.Trim(strings.TrimSpace(val), `"`); v != "" && !slices.Contains(out, v) {
				out = append(out, v)
			}
		}
	}
	return out
}

// readOwnedRegular reads (at most 1 MiB of) a regular file uid owns. O_NONBLOCK|O_NOFOLLOW: a FIFO
// or symlink planted at the name must neither block root nor point it elsewhere.
func readOwnedRegular(path string, uid uint32) (string, error) {
	f, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NONBLOCK|syscall.O_NOFOLLOW|syscall.O_NOCTTY, 0)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return "", err
	}
	if st, ok := fi.Sys().(*syscall.Stat_t); !fi.Mode().IsRegular() || !ok || st.Uid != uid {
		return "", fmt.Errorf("%s: not a regular file owned by uid %d", path, uid)
	}
	data, err := io.ReadAll(io.LimitReader(f, 1<<20))
	return string(data), err
}
