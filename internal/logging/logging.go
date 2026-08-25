// Package logging centralizes the verbosity model of the CLI: the mapping
// from the user-facing --verbose levels (constants.VerboseLevel) onto logrus
// levels, the --dump-log file hook and the bridge that routes Go's standard
// library logger (used by google/fscrypt during migrations) into logrus.
package logging

import (
	"bufio"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/constants"
)

// DefaultVerbose is the verbosity applied when --verbose is not given: it
// maps onto the logrus info level, i.e. exactly what the tools display
// today. A dump taken with only --dump-log therefore records precisely the
// displayed stream unless the operator raises --verbose explicitly.
const DefaultVerbose = constants.InternalWarningsLevelTwo

// VerboseToLogrus maps a VerboseLevel onto the logrus level that governs
// both the console stream and the --dump-log file:
//
//	NeededInfo                    -> warn  (essential lines only)
//	InternalWarningsLevelTwo      -> info  (what is printed today)
//	PrintAdditionalInfoLevelOne   -> debug (+ internal details)
//	PrintAdditionalInfoLevelTwo   -> trace (everything)
func VerboseToLogrus(v constants.VerboseLevel) (log.Level, error) {
	switch v {
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

// DumpHook mirrors selected logrus entries into the --dump-log writer in a
// plain, color-free format. It applies its own threshold derived from the
// verbose level, independent of the console logger level, so an interactive
// TUI run keeps its usual console output while the dump captures everything
// the verbose setting allows. Every write is flushed: an abrupt process exit
// (a signal-terminated daemon) never loses already-emitted lines.
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
	// logrus orders levels by severity (panic=0 .. trace=6): an entry is
	// within the threshold when it is at least as severe as min.
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

// formatEntry renders one plain-text line: "timestamp [LEVEL] message
// key=value ...". The console formatter force-enables ANSI colors; the dump
// must stay clean text.
func formatEntry(e *log.Entry) string {
	var b strings.Builder
	b.WriteString(e.Time.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(&b, " [%-5s] ", strings.ToUpper(e.Level.String()))
	b.WriteString(strings.TrimRight(e.Message, "\n"))
	for _, k := range sortedKeys(e.Data) {
		fmt.Fprintf(&b, " %s=%v", k, e.Data[k])
	}
	b.WriteByte('\n')
	return b.String()
}

// sortedKeys returns the entry's field keys in deterministic order so dumps
// are diff-friendly.
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

// BridgeStdLog routes the default Go standard library logger — google/fscrypt
// logs its migration progress through it — into logrus at info level. Those
// lines keep appearing exactly like today on the console AND become part of
// --dump-log, closing the hole where encryption-phase output was invisible
// to the dump. The returned restore puts back the previous writer/flags.
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
