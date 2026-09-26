package daemon

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Virgula0/app-listener/internal/logging"
	"github.com/Virgula0/app-listener/internal/usecase"
)

// gateLogWindow is how long repeats of one metadata-only process-gate denial
// are folded into a single summary line.
const gateLogWindow = time.Minute

// gateLogMaxKeys bounds the limiter's memory. Past it the limiter stops
// folding new keys and lets them log in full — never the other way round.
const gateLogMaxKeys = 4096

// gateLogLimiter folds repeated metadata-only process-gate denials in the daemon's LOG (desktop
// services like Hyprland and pipewire re-inspect the same process every few seconds, each refusal
// used to be a journal line). Presentation only: the kernel still denies every attempt.
// Deliberately narrow: only a blocked PTRACE with mode=READ (/proc/<pid> metadata) is folded;
// memory access (mode=ATTACH), PROC_MEM, TRACED_EXEC and every file-level denial are always logged
// in full, since those matter under real attack.
type gateLogLimiter struct {
	now     func() time.Time
	entries map[string]*gateLogEntry
}

type gateLogEntry struct {
	windowStart time.Time
	suppressed  int
	// sample is the full event of the first occurrence, for the summary.
	sample usecase.DaemonEvent
	// exe is the caller's executable, resolved once at the first occurrence
	// (the caller may be gone by the time the summary is written).
	exe string
}

func newGateLogLimiter() *gateLogLimiter {
	return &gateLogLimiter{now: time.Now, entries: make(map[string]*gateLogEntry)}
}

// foldable reports whether ev is the one kind of denial the limiter may fold.
func foldable(ev *usecase.DaemonEvent) bool {
	return ev.Event.Blocked && ev.Event.Process == "PTRACE" && strings.HasSuffix(ev.Event.Path, " mode=READ")
}

// admitEvent is the headless log filter: with dropMetadata set, a foldable denial is never written
// (not even summarized); otherwise the limiter decides.
func admitEvent(l *gateLogLimiter, ev *usecase.DaemonEvent, dropMetadata bool) bool {
	if dropMetadata && foldable(ev) {
		return false
	}
	return l.Admit(ev)
}

// Admit reports whether ev should be written in full. The first (resource, caller pid, caller comm,
// target) denial in a window is admitted; repeats are counted for the next Flush.
func (l *gateLogLimiter) Admit(ev *usecase.DaemonEvent) bool {
	if !foldable(ev) {
		return true
	}
	// Grouped by resource, caller PROCESS and target PROGRAM. Not caller comm (it names the thread;
	// a thread pool would count every worker as new) and not target pid (every new Steam/Wine
	// process would be a "first occurrence"). A window's first line keeps the full event, pid
	// included, as the example.
	key := fmt.Sprintf("%s\x00%d\x00%s", ev.Resource, ev.Event.PID, targetProgram(ev.Event.Path))
	if e, ok := l.entries[key]; ok {
		e.suppressed++
		return false
	}
	if len(l.entries) >= gateLogMaxKeys {
		return true
	}
	exe := "~"
	if target, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", ev.Event.PID)); err == nil {
		exe = logging.SanitizeText(target)
	}
	l.entries[key] = &gateLogEntry{windowStart: l.now(), sample: *ev, exe: exe}
	return true
}

// targetProgram drops the "pid=N " prefix of a process-gate path ("pid=46804 comm=steamwebhelper
// mode=READ" -> "comm=steamwebhelper mode=READ").
func targetProgram(path string) string {
	if i := strings.Index(path, " comm="); i >= 0 {
		return path[i+1:]
	}
	return path
}

// Flush writes one summary per key whose window elapsed with suppressed repeats, and forgets
// elapsed keys (a still-active pair reappears once per window). force flushes every key regardless
// of age (shutdown).
func (l *gateLogLimiter) Flush(w io.Writer, force bool) {
	now := l.now()
	keys := make([]string, 0, len(l.entries))
	for k := range l.entries {
		keys = append(keys, k)
	}
	sort.Strings(keys) // deterministic output
	for _, k := range keys {
		e := l.entries[k]
		if !force && now.Sub(e.windowStart) < gateLogWindow {
			continue
		}
		if e.suppressed > 0 {
			// A distinct prefix: this line is a count, not one denial, and
			// parsers keyed on "DAEMON DENIED  op=" must not take it for one.
			fmt.Fprintf(w, "%sDAEMON DENIED-REPEAT  op=PTRACE  comm=%s  commFullPath=%s  pid=%d  resource=%s  path=%s  repeats=%d  window=%s\n",
				syslogWarning, logging.SanitizeText(e.sample.Event.Comm), e.exe, e.sample.Event.PID,
				logging.SanitizeText(e.sample.Resource), logging.SanitizeText(e.sample.Event.Path),
				e.suppressed, now.Sub(e.windowStart).Round(time.Second))
		}
		delete(l.entries, k)
	}
}
