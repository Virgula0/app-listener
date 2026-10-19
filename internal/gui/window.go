//go:build gui

package gui

import (
	"context"
	"os"
	"time"

	"fyne.io/fyne/v2/driver"
	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/wm"
)

const prefFloatTiling = "floatOnTiling"

// detectWM finds the session's tiling WM. Its variables come from our environment, else from the
// ancestor restoreSessionEnv read (sudo drops them); sockets must belong to the invoking user.
func detectWM(inherited map[string]string) wm.Session {
	uid := uint32(os.Geteuid()) //nolint:gosec // uids fit
	if u, ok := invokingUID(); ok && uid == 0 {
		uid = u
	}
	return wm.Detect(func(k string) string {
		if v := os.Getenv(k); v != "" {
			return v
		}
		return inherited[k]
	}, uid)
}

// floatPref is the "float on tiling window managers" setting (default on).
func (s *shell) floatPref() bool { return s.app.Preferences().BoolWithFallback(prefFloatTiling, true) }

// setFloating floats the window (sized to its content and the output) or tiles it again, on a
// tiling WM only. A failure only logs: the window stays usable either way.
func (s *shell) setFloating(on bool) {
	if !s.wm.Tiling() {
		return
	}
	win := wm.Window{PID: os.Getpid()}
	if nw, ok := s.win.(driver.NativeWindow); ok {
		nw.RunNative(func(ctx any) {
			if x, ok := ctx.(driver.X11WindowContext); ok {
				win.XID = uint32(x.WindowHandle)
			}
		})
	}
	scale := s.win.Canvas().Scale()
	minSize := s.win.Content().MinSize()
	pref := wm.Size{W: int(preferredSize.Width * scale), H: int(preferredSize.Height * scale)}
	minPx := wm.Size{W: int(minSize.Width * scale), H: int(minSize.Height * scale)}
	sess := s.wm
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var err error
		if on {
			err = sess.Float(ctx, win, pref, minPx)
		} else {
			err = sess.Tile(ctx, win)
		}
		if err != nil {
			log.Warnf("gui: could not change the window's floating state on %s: %v", sess.Kind, err)
		}
	}()
}
