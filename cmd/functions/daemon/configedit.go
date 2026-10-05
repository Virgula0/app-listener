package daemon

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
	"github.com/Virgula0/app-listener/internal/logging"
)

// maxConfigBytes bounds a PUT body.
const maxConfigBytes = 4 << 20

// configEditor applies `edit-protected --edit-config`: write daemon.conf in place, run the
// in-process reload, and on a failed reload put the previous bytes back, so the file keeps
// describing the configuration that is still running.
type configEditor struct {
	path   string
	reload func() (*daemonconfig.Config, error)
	mu     sync.Mutex
}

// apply replaces base (what the client was sent) with next and reloads.
func (e *configEditor) apply(base, next []byte) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, err := daemonconfig.Parse(next); err != nil {
		return fmt.Errorf("the new configuration does not parse: %w", err)
	}
	if err := swapFile(e.path, base, next); err != nil {
		return err
	}
	log.Warnf("daemon: control: %s replaced via edit-protected --edit-config — reloading", e.path)
	if _, err := e.reload(); err != nil {
		if rerr := swapFile(e.path, next, base); rerr != nil {
			log.Errorf("daemon: control: restoring %s after a failed reload: %v", e.path, rerr)
			return fmt.Errorf("the reload failed (%v) and restoring %s failed too (%v) — the previous "+
				"configuration keeps running; fix the file before the next reload", err, e.path, rerr)
		}
		log.Warnf("daemon: control: the edited configuration did not reload — %s restored", e.path)
		return fmt.Errorf("the reload failed: %w — the previous configuration keeps running and %s was "+
			"restored", err, e.path)
	}
	return nil
}

func (m *controlManager) setConfigEditor(e *configEditor) {
	if m == nil {
		return
	}
	m.configEdit.Store(e)
}

// runConfigEdit serves CONFIG: send daemon.conf, then take either END or one PUT of its new
// content. The client edits offline, between two connections: nothing is widened meanwhile.
func (cs *controlServer) runConfigEdit(conn net.Conn, reader *bufio.Reader) {
	ed := cs.configEditor()
	if ed == nil {
		_ = cs.reply(conn, false, "configuration editing is not available yet — retry in a moment")
		return
	}
	base, err := os.ReadFile(ed.path)
	if err != nil {
		_ = cs.reply(conn, false, "reading the configuration: "+err.Error())
		return
	}
	_ = conn.SetWriteDeadline(time.Now().Add(controlHandshakeTimeout))
	if _, werr := fmt.Fprintf(conn, "OK %d\n%s", len(base), base); werr != nil {
		return
	}

	_ = conn.SetDeadline(time.Now().Add(controlSelectTimeout))
	next, err := readConfigPut(reader)
	if err != nil || next == nil {
		return // END, or the client went away: nothing changed
	}
	if bytes.Equal(next, base) {
		_ = cs.reply(conn, true, "")
		return
	}
	cs.mu.Lock()
	busy := cs.active != nil
	cs.mu.Unlock()
	if busy {
		_ = cs.reply(conn, false, "a live edit-protected session is active — the reload would revoke it; end it first")
		return
	}
	if err := ed.apply(base, next); err != nil {
		_ = cs.reply(conn, false, err.Error())
		log.Warnf("daemon: control: configuration edit refused: %v", logging.SanitizeText(err.Error()))
		return
	}
	_ = cs.reply(conn, true, "")
	log.Warnf("daemon: control: configuration edited and reloaded from %s", ed.path)
}

// readConfigPut parses "PUT <n>\n<n bytes>"; END yields nil.
func readConfigPut(reader *bufio.Reader) ([]byte, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.TrimRight(line, "\r\n")
	if line == ctrlEnd {
		return nil, nil
	}
	rest, ok := strings.CutPrefix(line, ctrlPut+" ")
	if !ok {
		return nil, fmt.Errorf("expected %q or %q", ctrlPut, ctrlEnd)
	}
	n, err := strconv.Atoi(rest)
	if err != nil || n < 1 || n > maxConfigBytes {
		return nil, errors.New("bad PUT size")
	}
	body := make([]byte, n)
	if _, err := io.ReadFull(reader, body); err != nil {
		return nil, err
	}
	return body, nil
}
