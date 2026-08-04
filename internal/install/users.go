package install

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// User is a login account that may have directories worth protecting.
type User struct {
	Name string
	UID  uint32
	Home string
}

// minLoginUID is the lower bound for "real" users; system accounts below
// this are skipped. Root (UID 0) is an exception: its home directory
// holds the most sensitive credentials on the system.
const minLoginUID = 1000

// ListUsers returns the local login users — root plus every user with a
// UID >= 1000 — whose home directory exists.
func ListUsers() ([]User, error) {
	data, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return nil, fmt.Errorf("read /etc/passwd: %w", err)
	}
	return parsePasswd(data)
}

// parsePasswd parses /etc/passwd content. Users are kept when the UID is
// 0 (root) or >= minLoginUID, the home directory is a plausible absolute
// path, and the home directory actually exists.
func parsePasswd(data []byte) ([]User, error) {
	var users []User
	for lineno, line := range strings.Split(string(data), "\n") {
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Split(line, ":")
		if len(fields) < 6 {
			return nil, fmt.Errorf("/etc/passwd line %d: expected 7 colon-separated fields, got %d", lineno+1, len(fields))
		}
		uid, err := strconv.ParseUint(fields[2], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("/etc/passwd line %d: bad UID %q: %w", lineno+1, fields[2], err)
		}
		home := fields[5]
		if (uid < minLoginUID && uid != 0) || home == "" || home == "/" || !strings.HasPrefix(home, "/") {
			continue
		}
		if _, err := os.Stat(home); err != nil {
			continue
		}
		users = append(users, User{Name: fields[0], UID: uint32(uid), Home: home})
	}
	return users, nil
}
