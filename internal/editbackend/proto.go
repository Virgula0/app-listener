// Package editbackend is the GUI's edit-protected link: the window runs the INSTALLED headless
// binary as `edit-protected --backend` and drives it with JSON lines over its stdin/stdout. Every
// protected-file access, the control-socket session and the password stay in that child, the binary
// the daemon's peer check admits; the GUI process only renders. The backend does nothing the
// `edit-protected` CLI can't: same root check, same daemon-side AUTH, same confined vault code.
package editbackend

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"
)

// Proto is the wire version; both sides refuse a mismatch (GUI and installed binary of different
// releases).
const Proto = 1

// MaxLine bounds one JSON line: a write's content (base64, up to MaxWriteBytes) plus framing.
const MaxLine = 8 << 20

// MaxWriteBytes caps one write's content (the editors load files up to 2 MiB; edits may grow them).
const MaxWriteBytes = 4 << 20

// Ops. Hello comes first; then one workflow per process (edit, forward, config, password).
const (
	OpHello         = "hello"
	OpAuth          = "auth"
	OpSelect        = "select"
	OpOfflineScan   = "offline_scan"
	OpOfflineOpen   = "offline_open"
	OpList          = "list"
	OpRead          = "read"
	OpWrite         = "write"
	OpCreate        = "create"
	OpChmod         = "chmod"
	OpChown         = "chown"
	OpRemove        = "remove"
	OpUsers         = "users"
	OpPing          = "ping"
	OpEnd           = "end"
	OpForwardCheck  = "forward_check"
	OpForward       = "forward"
	OpConfigGet     = "config_get"
	OpConfigPut     = "config_put"
	OpPasswordState = "password_state"
	OpSetPassword   = "set_password"
	OpClearPassword = "clear_password"
)

// Kinds of backend messages: a reply to one request, or an unsolicited event.
const (
	KindReply = "reply"
	// KindEvent carries a FORWARD grant's guard event.
	KindEvent = "event"
	// KindEnded: the daemon ended the session (idle timeout, reload, a granted binary changed).
	KindEnded = "ended"
)

// Request is one GUI -> backend message. Paths are relative to the opened resource ("." = itself).
type Request struct {
	ID          uint64   `json:"id"`
	Op          string   `json:"op"`
	Password    string   `json:"password,omitempty"`
	NewPassword string   `json:"new_password,omitempty"`
	TimeoutSec  int64    `json:"timeout_sec,omitempty"`
	Resource    string   `json:"resource,omitempty"`
	Resources   []string `json:"resources,omitempty"`
	Path        string   `json:"path,omitempty"`
	Name        string   `json:"name,omitempty"`
	Dir         bool     `json:"dir,omitempty"`
	Data        []byte   `json:"data,omitempty"`
	Base        []byte   `json:"base,omitempty"`
	Mode        string   `json:"mode,omitempty"`
	UID         int      `json:"uid,omitempty"`
	GID         int      `json:"gid,omitempty"`
	Binaries    []string `json:"binaries,omitempty"`
	Allow       bool     `json:"allow,omitempty"`
	Events      []string `json:"events,omitempty"`
}

// Message is one backend -> GUI message.
type Message struct {
	Kind string `json:"kind"`
	ID   uint64 `json:"id,omitempty"`
	Err  string `json:"err,omitempty"`

	Hello      *Hello        `json:"hello,omitempty"`
	Resources  []string      `json:"resources,omitempty"`
	SingleFile bool          `json:"single_file,omitempty"`
	Entries    []Entry       `json:"entries,omitempty"`
	File       *File         `json:"file,omitempty"`
	Entry      *Entry        `json:"entry,omitempty"`
	Users      []User        `json:"users,omitempty"`
	Notes      []string      `json:"notes,omitempty"`
	Findings   []string      `json:"findings,omitempty"`
	Binaries   []BinaryCheck `json:"binaries,omitempty"`
	Data       []byte        `json:"data,omitempty"`
	Password   *PwState      `json:"password,omitempty"`
	Event      *Event        `json:"event,omitempty"`
	// Stale: a config_put found daemon.conf changed since Base; Data is the current file.
	Stale bool `json:"stale,omitempty"`
}

// Hello is the backend's mode report, the same decision `edit-protected` makes from its flags.
type Hello struct {
	Proto   int    `json:"proto"`
	Version string `json:"version"`
	// Live: the control socket exists (a running daemon with a password): edits go live, and
	// forward / config are available. Otherwise edits run offline (daemon must be stopped).
	Live bool `json:"live"`
	// HasPassword: a password file is configured (with !Live: the daemon isn't running).
	HasPassword bool `json:"has_password"`
}

// Entry is one directory entry (or a stat). Invalid: its name isn't UTF-8, so it can't be named
// back over the wire; shown, never operated on.
type Entry struct {
	Name    string    `json:"name"`
	Dir     bool      `json:"dir,omitempty"`
	Symlink bool      `json:"symlink,omitempty"`
	Size    int64     `json:"size"`
	Mode    uint32    `json:"mode"`
	ModeStr string    `json:"mode_str"`
	UID     uint32    `json:"uid"`
	GID     uint32    `json:"gid"`
	ModTime time.Time `json:"mtime"`
	Invalid bool      `json:"invalid,omitempty"`
}

// File is a read's result; Data is empty when the file is binary or too large to edit.
type File struct {
	Data     []byte `json:"data,omitempty"`
	Size     int64  `json:"size"`
	Binary   bool   `json:"binary,omitempty"`
	TooLarge bool   `json:"too_large,omitempty"`
}

// User is one chown target.
type User struct {
	Name string `json:"name"`
	UID  int    `json:"uid"`
	GID  int    `json:"gid"`
}

// BinaryCheck is what `--forward -w` confirms before admitting a binary.
type BinaryCheck struct {
	Path         string   `json:"path"`
	SystemPlaced bool     `json:"system_placed"`
	Runtimes     []string `json:"runtimes,omitempty"`
	Err          string   `json:"err,omitempty"`
}

// PwState is the password file's state.
type PwState struct {
	Configured bool `json:"configured"`
	// Origin is "install" or "cli" when Configured.
	Origin        string `json:"origin,omitempty"`
	DaemonRunning bool   `json:"daemon_running"`
}

// Event is a granted binary's access (guard.GuardEvent, flattened for display only).
type Event struct {
	Time    time.Time `json:"time"`
	PID     uint32    `json:"pid"`
	UID     uint32    `json:"uid"`
	Comm    string    `json:"comm"`
	Op      string    `json:"op"`
	Path    string    `json:"path"`
	Dest    string    `json:"dest,omitempty"`
	Allowed bool      `json:"allowed"`
}

// Writer serializes messages as JSON lines.
type Writer struct {
	enc *json.Encoder
}

// NewWriter writes to w.
func NewWriter(w io.Writer) *Writer {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return &Writer{enc: enc}
}

// Write sends one message.
func (w *Writer) Write(v any) error { return w.enc.Encode(v) }

// Reader parses JSON lines of at most MaxLine bytes.
type Reader struct {
	sc *bufio.Scanner
}

// NewReader reads from r.
func NewReader(r io.Reader) *Reader {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 64<<10), MaxLine)
	return &Reader{sc: sc}
}

// Read decodes the next line into v; io.EOF at a clean end.
func (r *Reader) Read(v any) error {
	if !r.sc.Scan() {
		if err := r.sc.Err(); err != nil {
			return err
		}
		return io.EOF
	}
	return json.Unmarshal(r.sc.Bytes(), v)
}

// CleanRel validates a relative path from the wire: UTF-8, no NUL, not absolute, never above
// ".". The vault re-checks containment (os.Root); this only rejects malformed requests early.
func CleanRel(rel string) (string, error) {
	if rel == "" {
		rel = "."
	}
	switch {
	case !utf8.ValidString(rel), strings.ContainsRune(rel, 0):
		return "", errors.New("malformed path")
	case len(rel) > 4096:
		return "", errors.New("path too long")
	case filepath.IsAbs(rel):
		return "", fmt.Errorf("%s: expected a path relative to the resource", rel)
	}
	c := filepath.Clean(rel)
	if c == ".." || strings.HasPrefix(c, "../") {
		return "", fmt.Errorf("%s is outside the resource", rel)
	}
	return c, nil
}
