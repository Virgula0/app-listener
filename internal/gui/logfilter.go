//go:build gui

package gui

import (
	"bytes"
	"io"
	"log"
	"sync"
)

// benignFyneErrors are fyne.LogError reasons that are expected for a root window under sudo: no
// session D-Bus to watch the desktop theme (the app has its own), no XDG folder for a file
// dialog's favorites.
var benignFyneErrors = [][]byte{
	[]byte("failed to watch theme settings"),
	[]byte("Getting favorite locations"),
}

// fyneLogFilter drops benign fyne.LogError reports from the standard logger: each is a
// "Fyne error:" line followed by "  Cause:" and "  At:" lines, one Write per line.
type fyneLogFilter struct {
	mu       sync.Mutex
	out      io.Writer
	skipping bool
}

func (f *fyneLogFilter) Write(p []byte) (int, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	switch {
	case bytes.Contains(p, []byte("Fyne error: ")):
		f.skipping = false
		for _, reason := range benignFyneErrors {
			if bytes.Contains(p, reason) {
				f.skipping = true
			}
		}
	case f.skipping && (bytes.Contains(p, []byte("  Cause:")) || bytes.Contains(p, []byte("  At:"))):
	default:
		f.skipping = false
	}
	if f.skipping {
		return len(p), nil
	}
	return f.out.Write(p)
}

// filterFyneLogs installs fyneLogFilter on the standard logger until restore is called.
func filterFyneLogs() (restore func()) {
	prev := log.Writer()
	log.SetOutput(&fyneLogFilter{out: prev})
	return func() { log.SetOutput(prev) }
}
