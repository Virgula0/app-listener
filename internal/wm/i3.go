package wm

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
	"strconv"
	"time"
)

// i3 IPC message types (shared by sway).
const (
	i3RunCommand    = 0
	i3GetWorkspaces = 1
	i3GetTree       = 4
)

var i3Magic = []byte("i3-ipc")

// i3Request sends one message: "i3-ipc", payload length and type (native byte order), payload.
func (s Session) i3Request(ctx context.Context, typ uint32, payload string) ([]byte, error) {
	var d net.Dialer
	c, err := d.DialContext(ctx, "unix", s.socket)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(2 * time.Second)
	}
	_ = c.SetDeadline(deadline)
	msg := make([]byte, 0, len(i3Magic)+8+len(payload))
	msg = append(msg, i3Magic...)
	msg = binary.NativeEndian.AppendUint32(msg, uint32(len(payload))) //nolint:gosec // short command strings
	msg = binary.NativeEndian.AppendUint32(msg, typ)
	msg = append(msg, payload...)
	if _, werr := c.Write(msg); werr != nil {
		return nil, werr
	}
	head := make([]byte, len(i3Magic)+8)
	if _, rerr := io.ReadFull(c, head); rerr != nil {
		return nil, rerr
	}
	if !bytes.Equal(head[:len(i3Magic)], i3Magic) {
		return nil, errors.New("i3 ipc: bad reply magic")
	}
	n := binary.NativeEndian.Uint32(head[len(i3Magic):])
	if n > maxReply {
		return nil, fmt.Errorf("i3 ipc: reply of %d bytes", n)
	}
	body := make([]byte, n)
	_, err = io.ReadFull(c, body)
	return body, err
}

// reachable reports whether an i3 socket answers (a crashed i3 can leave its socket behind).
func reachable(path string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err := Session{Kind: I3, socket: path}.i3Request(ctx, i3GetWorkspaces, "")
	return err == nil
}

// i3Criteria selects w: sway knows client PIDs, i3 only X11 window ids.
func (s Session) i3Criteria(w Window) (string, error) {
	if s.Kind == Sway {
		return fmt.Sprintf("[pid=%d]", w.PID), nil
	}
	if w.XID == 0 {
		return "", errors.New("i3: the window has no X11 id")
	}
	return fmt.Sprintf("[id=%d]", w.XID), nil
}

// i3Managed waits until w appears in the layout tree.
func (s Session) i3Managed(ctx context.Context, w Window) error {
	key, id := `"window"`, uint64(w.XID)
	if s.Kind == Sway {
		key, id = `"pid"`, uint64(w.PID) //nolint:gosec // a PID is positive
	}
	re := regexp.MustCompile(key + `:\s*` + strconv.FormatUint(id, 10) + `\s*[,}]`)
	return poll(ctx, func() (bool, error) {
		tree, err := s.i3Request(ctx, i3GetTree, "")
		if err != nil {
			return false, err
		}
		return re.Match(tree), nil
	})
}

type i3Workspace struct {
	Focused bool `json:"focused"`
	Rect    struct {
		X      int `json:"x"`
		Y      int `json:"y"`
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"rect"`
}

func (s Session) i3FocusedArea(ctx context.Context) (Rect, error) {
	out, err := s.i3Request(ctx, i3GetWorkspaces, "")
	if err != nil {
		return Rect{}, err
	}
	var wss []i3Workspace
	if err := json.Unmarshal(out, &wss); err != nil {
		return Rect{}, fmt.Errorf("i3 workspaces: %w", err)
	}
	for _, ws := range wss {
		if ws.Focused {
			return Rect{X: ws.Rect.X, Y: ws.Rect.Y, W: ws.Rect.Width, H: ws.Rect.Height}, nil
		}
	}
	return Rect{}, errors.New("i3: no focused workspace")
}

type i3Result struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func (s Session) i3Command(ctx context.Context, cmd string) error {
	out, err := s.i3Request(ctx, i3RunCommand, cmd)
	if err != nil {
		return err
	}
	var res []i3Result
	if err := json.Unmarshal(out, &res); err != nil {
		return fmt.Errorf("i3 command reply: %w", err)
	}
	for _, r := range res {
		if !r.Success {
			return fmt.Errorf("i3 command %q: %s", cmd, r.Error)
		}
	}
	return nil
}

func (s Session) i3Float(ctx context.Context, w Window, pref, minSize Size) error {
	crit, err := s.i3Criteria(w)
	if err != nil {
		return err
	}
	if werr := s.i3Managed(ctx, w); werr != nil {
		return werr
	}
	area, err := s.i3FocusedArea(ctx)
	if err != nil {
		return err
	}
	size := FitSize(pref, minSize, area)
	return s.i3Command(ctx, fmt.Sprintf("%s floating enable, resize set %d px %d px, move position center",
		crit, size.W, size.H))
}

func (s Session) i3Tile(ctx context.Context, w Window) error {
	crit, err := s.i3Criteria(w)
	if err != nil {
		return err
	}
	return s.i3Command(ctx, crit+" floating disable")
}
