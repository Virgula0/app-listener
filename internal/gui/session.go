//go:build gui

package gui

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"syscall"
)

// sessionVars are the display variables sudo's env_reset drops; nothing else is ever imported.
var sessionVars = []string{"XDG_RUNTIME_DIR", "WAYLAND_DISPLAY", "DISPLAY", "XAUTHORITY"}

// wmVars locate the tiling WM's IPC socket (wm.Detect checks its owner). They are read from the
// same ancestor but never exported: only the floating request uses them.
var wmVars = []string{"HYPRLAND_INSTANCE_SIGNATURE", "SWAYSOCK", "I3SOCK"}

var displayRe = regexp.MustCompile(`^[A-Za-z0-9.\-]*:\d+(\.\d+)?$`)

// restoreSessionEnv lets a root GUI reach the invoking user's compositor or X server under
// sudo/pkexec. Only unset variables are filled, from the nearest ancestor running as that user,
// and each must name a socket/dir/file that user owns (see xauthority for the cookie file). It
// returns that ancestor's wmVars.
func restoreSessionEnv() (wmEnv map[string]string) {
	if os.Geteuid() != 0 {
		return nil
	}
	uid, ok := invokingUID()
	if !ok {
		return nil
	}
	found := ancestorEnv(uid)
	wmEnv = map[string]string{}
	for _, k := range wmVars {
		wmEnv[k] = found[k]
	}
	if os.Getenv("XDG_RUNTIME_DIR") == "" {
		dir := found["XDG_RUNTIME_DIR"]
		if dir == "" {
			dir = fmt.Sprintf("/run/user/%d", uid)
		}
		if ownedBy(dir, uid, os.ModeDir) {
			_ = os.Setenv("XDG_RUNTIME_DIR", dir)
		}
	}
	if os.Getenv("WAYLAND_DISPLAY") == "" {
		if name := found["WAYLAND_DISPLAY"]; name != "" && waylandSocketOK(name, uid) {
			_ = os.Setenv("WAYLAND_DISPLAY", name)
		}
	}
	if os.Getenv("DISPLAY") == "" && displayRe.MatchString(found["DISPLAY"]) {
		_ = os.Setenv("DISPLAY", found["DISPLAY"])
	}
	if os.Getenv("XAUTHORITY") == "" {
		if auth := xauthority(found["XAUTHORITY"], uid); auth != "" {
			_ = os.Setenv("XAUTHORITY", auth)
		}
	}
	return wmEnv
}

// invokingUID is the non-root user sudo or pkexec ran us for (both set it themselves).
func invokingUID() (uint32, bool) {
	for _, name := range []string{"SUDO_UID", "PKEXEC_UID"} {
		if v, err := strconv.ParseUint(os.Getenv(name), 10, 32); err == nil && v != 0 {
			return uint32(v), true
		}
	}
	return 0, false
}

// ancestorEnv reads sessionVars and wmVars from the closest ancestor whose real uid is uid (the shell that
// ran sudo; sudo's own processes in between run as root).
func ancestorEnv(uid uint32) map[string]string {
	out := map[string]string{}
	pid := os.Getppid()
	for depth := 0; depth < 8 && pid > 1; depth++ {
		status, err := os.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
		if err != nil {
			return out
		}
		ruid, ppid := statusField(status, "Uid:"), statusField(status, "PPid:")
		if ruid == strconv.FormatUint(uint64(uid), 10) {
			environ, envErr := os.ReadFile(fmt.Sprintf("/proc/%d/environ", pid))
			if envErr != nil {
				return out
			}
			for kv := range bytes.SplitSeq(environ, []byte{0}) {
				k, v, ok := strings.Cut(string(kv), "=")
				if ok && (slices.Contains(sessionVars, k) || slices.Contains(wmVars, k)) {
					out[k] = v
				}
			}
			return out
		}
		next, err := strconv.Atoi(ppid)
		if err != nil {
			return out
		}
		pid = next
	}
	return out
}

// statusField returns the first value of a /proc/<pid>/status line ("Uid:" → real uid).
func statusField(status []byte, key string) string {
	for line := range strings.SplitSeq(string(status), "\n") {
		if rest, ok := strings.CutPrefix(line, key); ok {
			if f := strings.Fields(rest); len(f) > 0 {
				return f[0]
			}
		}
	}
	return ""
}

func waylandSocketOK(name string, uid uint32) bool {
	path := name
	if !filepath.IsAbs(name) {
		if strings.Contains(name, "/") {
			return false
		}
		path = filepath.Join(os.Getenv("XDG_RUNTIME_DIR"), name)
	}
	return ownedBy(path, uid, os.ModeSocket)
}

// xauthority copies the user's cookie file (their XAUTHORITY, else ~/.Xauthority) into a private
// root-owned file, opened once with O_NOFOLLOW: libXau never opens a path the user can swap.
func xauthority(fromEnv string, uid uint32) string {
	candidates := []string{fromEnv}
	if u, err := user.LookupId(strconv.FormatUint(uint64(uid), 10)); err == nil {
		candidates = append(candidates, filepath.Join(u.HomeDir, ".Xauthority"))
	}
	for _, c := range candidates {
		if c == "" || !filepath.IsAbs(c) {
			continue
		}
		if cookie, ok := readOwnedFile(c, uid); ok {
			if copyPath := writePrivate(cookie); copyPath != "" {
				return copyPath
			}
		}
	}
	return ""
}

func readOwnedFile(path string, uid uint32) ([]byte, bool) {
	const maxCookieFile = 1 << 20
	f, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if err != nil {
		return nil, false
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil || !fi.Mode().IsRegular() || fi.Size() > maxCookieFile {
		return nil, false
	}
	if st, ok := fi.Sys().(*syscall.Stat_t); !ok || st.Uid != uid {
		return nil, false
	}
	data, err := io.ReadAll(io.LimitReader(f, maxCookieFile))
	return data, err == nil
}

func writePrivate(data []byte) string {
	f, err := os.CreateTemp("", "app-listener-xauth-*")
	if err != nil {
		return ""
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		_ = os.Remove(f.Name())
		return ""
	}
	xauthCopy = f.Name()
	return f.Name()
}

// xauthCopy is the private cookie copy restoreSessionEnv made, removed by cleanupSessionEnv.
var xauthCopy string

func cleanupSessionEnv() {
	if xauthCopy != "" {
		_ = os.Remove(xauthCopy)
	}
}

// ownedBy reports whether path is (without following a symlink) of the given type (0 = regular
// file) and owned by uid.
func ownedBy(path string, uid uint32, typ os.FileMode) bool {
	fi, err := os.Lstat(path)
	if err != nil || fi.Mode().Type() != typ {
		return false
	}
	st, ok := fi.Sys().(*syscall.Stat_t)
	return ok && st.Uid == uid
}

// invokingHome is the home of the user sudo/pkexec ran us for, else ours.
func invokingHome() string {
	if uid, ok := invokingUID(); ok {
		if u, err := user.LookupId(strconv.FormatUint(uint64(uid), 10)); err == nil {
			return u.HomeDir
		}
	}
	home, _ := os.UserHomeDir()
	return home
}
