// Package logging centralizes the CLI verbosity model: mapping constants.VerboseLevel onto
// logrus levels, the --dump-log file hook, and routing Go's stdlib logger (used by google/fscrypt
// during migrations) into logrus.
package logging

import (
	"bufio"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"strings"
	"sync"
	"unicode"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/constants"
)

// DefaultVerbose is the verbosity applied when --verbose is omitted (logrus info): a dump taken
// with only --dump-log records exactly the displayed stream unless the operator raises --verbose.
const DefaultVerbose = constants.InternalWarningsLevelTwo

// VerboseToLogrus maps a VerboseLevel onto the logrus level governing both the console
// stream and the --dump-log file; see the ladder documented on constants.VerboseLevel.
func VerboseToLogrus(v constants.VerboseLevel) (log.Level, error) {
	switch v {
	case constants.ErrorsOnly:
		return log.ErrorLevel, nil
	case constants.NeededInfo:
		return log.WarnLevel, nil
	case constants.InternalWarningsLevelTwo:
		return log.InfoLevel, nil
	case constants.PrintAdditionalInfoLevelOne:
		return log.DebugLevel, nil
	case constants.PrintAdditionalInfoLevelTwo:
		return log.TraceLevel, nil
	default:
		return 0, fmt.Errorf("invalid verbose level %d (valid: 0..%d)",
			int(v), int(constants.PrintAdditionalInfoLevelTwo))
	}
}

// DumpHook mirrors selected logrus entries into the --dump-log writer as plain, color-free text,
// applying its own verbose-derived threshold independent of the console logger level. Every write
// is flushed (sync-per-write): an abrupt process exit never loses already-emitted lines.
type DumpHook struct {
	mu  sync.Mutex
	w   io.Writer
	min log.Level
}

// NewDumpHook builds a dump hook for the given verbose level.
func NewDumpHook(w io.Writer, v constants.VerboseLevel) (*DumpHook, error) {
	threshold, err := VerboseToLogrus(v)
	if err != nil {
		return nil, err
	}
	return &DumpHook{w: w, min: threshold}, nil
}

// Levels implements logrus.Hook: every level passes through, Fire filters.
func (h *DumpHook) Levels() []log.Level {
	return []log.Level{
		log.PanicLevel, log.FatalLevel, log.ErrorLevel,
		log.WarnLevel, log.InfoLevel, log.DebugLevel, log.TraceLevel,
	}
}

// Fire writes one formatted line per entry below-or-equal the threshold.
func (h *DumpHook) Fire(e *log.Entry) error {
	// Levels are severity-ordered (panic=0 .. trace=6): within-threshold = at least as severe.
	if e.Level > h.min {
		return nil
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	if _, err := io.WriteString(h.w, formatEntry(e)); err != nil {
		return err
	}
	if f, ok := h.w.(*os.File); ok {
		return f.Sync()
	}
	return nil
}

// formatEntry renders one plain-text line: "timestamp [LEVEL] message key=value ...";
// unlike the ANSI-colored console formatter, the dump stays clean text.
func formatEntry(e *log.Entry) string {
	var b strings.Builder
	b.WriteString(e.Time.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(&b, " [%-5s] ", strings.ToUpper(e.Level.String()))
	b.WriteString(strings.TrimRight(SanitizeText(e.Message), "\n"))
	for _, k := range sortedKeys(e.Data) {
		fmt.Fprintf(&b, " %s=%v", k, SanitizeText(fmt.Sprint(e.Data[k])))
	}
	b.WriteByte('\n')
	return b.String()
}

// SanitizeText replaces control characters (which includes terminal escape
// sequences such as ESC/OSC and forged newlines) with '?'. Event fields
// originate from the kernel ring buffer — attacker-controlled filenames,
// process names and destinations — and would otherwise forge audit lines in
// journald or the --dump-log file.
func SanitizeText(value string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return '?'
		}
		return r
	}, value)
}

// sortedKeys returns entry field keys in deterministic, diff-friendly order.
func sortedKeys(data map[string]any) []string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sortStrings(keys)
	return keys
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}

// BridgeStdLog routes the default Go stdlib logger (google/fscrypt logs its migration
// progress through it) into logrus at info level, so those lines reach both the console and
// --dump-log. The returned restore reinstates the previous writer/flags.
func BridgeStdLog() (restore func()) {
	prevWriter := stdlog.Writer()
	prevFlags := stdlog.Flags()

	r, w := io.Pipe()
	done := make(chan struct{})
	go func() {
		defer close(done)
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				log.Infof("%s", line)
			}
		}
	}()

	stdlog.SetOutput(w)
	stdlog.SetFlags(0)
	return func() {
		stdlog.SetOutput(prevWriter)
		stdlog.SetFlags(prevFlags)
		_ = w.Close() // ends the scanner goroutine
		<-done
	}
}
