package install

import (
	"fmt"
	"io"
	stdlog "log"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

// bottomBar is a single-line progress bar pinned at the bottom of the
// terminal. The loggers (logrus and the standard log used by the fscrypt
// library) are routed through it: before each log line the bar line is
// cleared, the line is written, and the bar is re-rendered underneath — so
// the fscrypt logs scroll normally above a bar that stays visible. When
// stderr is not a terminal the bar is disabled and everything passes
// through untouched.
type bottomBar struct {
	mu   sync.Mutex
	w    io.Writer
	tty  bool
	bar  string
	show bool
}

func newBottomBar(w io.Writer) *bottomBar {
	b := &bottomBar{w: w}
	_, err := unix.IoctlGetWinsize(int(os.Stderr.Fd()), unix.TIOCGWINSZ)
	b.tty = err == nil
	return b
}

// set renders (or updates) the bar at the bottom of the terminal.
func (b *bottomBar) set(label string, fraction float64) {
	if !b.tty {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if bar := renderBar(label, fraction); bar != b.bar {
		b.bar = bar
		b.render()
	}
}

func (b *bottomBar) render() {
	if !b.tty || b.bar == "" {
		return
	}
	// Carriage return, clear the whole line, print the bar without a
	// trailing newline so the next log line can clear it again.
	fmt.Fprintf(b.w, "\r\x1b[2K%s", b.bar)
	b.show = true
}

// Write implements io.Writer for the loggers: a visible bar line is
// cleared before the log text and re-rendered after every completed line.
func (b *bottomBar) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.tty && b.show {
		fmt.Fprint(b.w, "\r\x1b[2K")
		b.show = false
	}
	n, err := b.w.Write(p)
	if b.tty && len(p) > 0 && p[len(p)-1] == '\n' {
		b.render()
	}
	return n, err
}

// done clears the bar and leaves the cursor at the start of a fresh line.
func (b *bottomBar) done() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.tty {
		fmt.Fprint(b.w, "\r\x1b[2K")
	}
	b.bar = ""
	b.show = false
}

// renderBar builds the bar line: label, a 20-block progress bar and the
// percentage, truncated to the current terminal width (resizing-safe: the
// width is re-queried on every render).
func renderBar(label string, fraction float64) string {
	pct := int(fraction * 100)
	if pct > 100 {
		pct = 100
	}
	const blocks = 20
	filled := blocks * pct / 100
	bar := strings.Repeat("█", filled) + strings.Repeat("░", blocks-filled)
	text := fmt.Sprintf("%s %s %3d%%", label, bar, pct)
	if ws, err := unix.IoctlGetWinsize(int(os.Stderr.Fd()), unix.TIOCGWINSZ); err == nil && ws.Col > 0 {
		if width := int(ws.Col) - 1; len(text) > width {
			text = "…" + text[len(text)-width+1:]
		}
	}
	return text
}

// withBottomBar runs run with the progress bar active and the loggers
// routed through it; the bar and the redirects are cleaned up on return,
// so the terminal is left exactly as it was.
func withBottomBar(run func(bar *bottomBar) error) error {
	bar := newBottomBar(os.Stderr)
	prevStd := stdlog.Writer()
	prevLogrus := log.StandardLogger().Out
	stdlog.SetOutput(bar)
	log.SetOutput(bar)
	defer func() {
		bar.done()
		stdlog.SetOutput(prevStd)
		log.SetOutput(prevLogrus)
	}()
	return run(bar)
}
