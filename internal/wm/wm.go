// Package wm floats a window on a tiling window manager (Hyprland, sway, i3) over the
// compositor's own IPC socket: fyne exposes no window-type hint, and these WMs tile every
// window by default.
package wm

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"syscall"
	"time"
)

// Kind is a supported tiling window manager.
type Kind int

const (
	None Kind = iota
	Hyprland
	Sway
	I3
)

func (k Kind) String() string {
	switch k {
	case Hyprland:
		return "Hyprland"
	case Sway:
		return "sway"
	case I3:
		return "i3"
	}
	return "none"
}

// Session is the detected window manager and its IPC socket.
type Session struct {
	Kind   Kind
	socket string
}

// Tiling reports whether a supported tiling window manager was found.
func (s Session) Tiling() bool { return s.Kind != None }

// Window identifies the window to act on: Hyprland and sway match its PID, i3 its X11 id.
type Window struct {
	PID int
	XID uint32
}

// Size is a width and height in the WM's pixels (logical on Wayland).
type Size struct{ W, H int }

// Rect is an area of the WM's layout.
type Rect struct{ X, Y, W, H int }

var hyprSigRe = regexp.MustCompile(`^[A-Za-z0-9_]+$`)

// Detect finds the session's window manager from env (the invoking user's session variables).
// Every socket must be a socket owned by uid: the caller runs as root and talks to the user's WM.
func Detect(env func(string) string, uid uint32) Session {
	if p := hyprlandSocket(env, uid); p != "" {
		return Session{Kind: Hyprland, socket: p}
	}
	if p := env("SWAYSOCK"); filepath.IsAbs(p) && socketOwnedBy(p, uid) {
		return Session{Kind: Sway, socket: p}
	}
	if p := env("I3SOCK"); filepath.IsAbs(p) && socketOwnedBy(p, uid) {
		return Session{Kind: I3, socket: p}
	}
	if p := i3RuntimeSocket(env, uid); p != "" {
		return Session{Kind: I3, socket: p}
	}
	return Session{}
}

func hyprlandSocket(env func(string) string, uid uint32) string {
	sig := env("HYPRLAND_INSTANCE_SIGNATURE")
	if sig == "" || !hyprSigRe.MatchString(sig) {
		return ""
	}
	for _, dir := range []string{filepath.Join(env("XDG_RUNTIME_DIR"), "hypr"), "/tmp/hypr"} {
		if p := filepath.Join(dir, sig, ".socket.sock"); filepath.IsAbs(p) && socketOwnedBy(p, uid) {
			return p
		}
	}
	return ""
}

// i3RuntimeSocket finds i3's socket in the runtime dir: i3 exports I3SOCK only to its children, so
// a terminal started otherwise lacks it.
func i3RuntimeSocket(env func(string) string, uid uint32) string {
	runtime := env("XDG_RUNTIME_DIR")
	if !filepath.IsAbs(runtime) || env("DISPLAY") == "" || env("WAYLAND_DISPLAY") != "" {
		return ""
	}
	matches, _ := filepath.Glob(filepath.Join(runtime, "i3", "ipc-socket.*"))
	for _, p := range matches {
		if socketOwnedBy(p, uid) && reachable(p) {
			return p
		}
	}
	return ""
}

func socketOwnedBy(path string, uid uint32) bool {
	fi, err := os.Lstat(path)
	if err != nil || fi.Mode().Type() != os.ModeSocket {
		return false
	}
	st, ok := fi.Sys().(*syscall.Stat_t)
	return ok && st.Uid == uid
}

// FitSize is pref grown to minSize and capped at 90% of the usable area.
func FitSize(pref, minSize Size, area Rect) Size {
	fit := func(p, m, a int) int {
		v := max(p, m)
		if limit := a * 9 / 10; limit > 0 && v > limit {
			v = limit
		}
		return max(v, 1)
	}
	return Size{W: fit(pref.W, minSize.W, area.W), H: fit(pref.H, minSize.H, area.H)}
}

// centered is s centered in area.
func centered(s Size, area Rect) Rect {
	return Rect{X: area.X + (area.W-s.W)/2, Y: area.Y + (area.H-s.H)/2, W: s.W, H: s.H}
}

// ErrNotFound means the window never showed up in the WM's window list.
var ErrNotFound = errors.New("window not managed by the window manager yet")

// Float makes w floating, sized by FitSize on its output and centered. It waits (until ctx ends)
// for the WM to manage the window, since a just-shown window is mapped asynchronously.
func (s Session) Float(ctx context.Context, w Window, pref, minSize Size) error {
	switch s.Kind {
	case Hyprland:
		return s.hyprFloat(ctx, w, pref, minSize)
	case Sway, I3:
		return s.i3Float(ctx, w, pref, minSize)
	}
	return errors.New("no tiling window manager detected")
}

// Tile puts w back into the tiling layout.
func (s Session) Tile(ctx context.Context, w Window) error {
	switch s.Kind {
	case Hyprland:
		return s.hyprTile(ctx, w)
	case Sway, I3:
		return s.i3Tile(ctx, w)
	}
	return errors.New("no tiling window manager detected")
}

// poll retries find every 100ms until it reports done or ctx ends.
func poll(ctx context.Context, find func() (bool, error)) error {
	tick := time.NewTicker(100 * time.Millisecond)
	defer tick.Stop()
	for {
		done, err := find()
		if err != nil {
			return err
		}
		if done {
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w: %w", ErrNotFound, ctx.Err())
		case <-tick.C:
		}
	}
}
