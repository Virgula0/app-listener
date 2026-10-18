package wm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

// maxReply caps what is read from a WM socket (the client and tree listings grow per window).
const maxReply = 16 << 20

type hyprClient struct {
	Address  string `json:"address"`
	PID      int    `json:"pid"`
	Floating bool   `json:"floating"`
	Monitor  int    `json:"monitor"`
	Mapped   bool   `json:"mapped"`
}

type hyprMonitor struct {
	ID        int     `json:"id"`
	X         int     `json:"x"`
	Y         int     `json:"y"`
	Width     int     `json:"width"`
	Height    int     `json:"height"`
	Scale     float64 `json:"scale"`
	Transform int     `json:"transform"`
	Reserved  []int   `json:"reserved"`
}

// area is the monitor's usable logical area: scaled, rotated, minus bars ([left, top, right, bottom]).
func (m *hyprMonitor) area() Rect {
	w, h := float64(m.Width), float64(m.Height)
	if m.Scale > 0 {
		w, h = w/m.Scale, h/m.Scale
	}
	if m.Transform%2 == 1 {
		w, h = h, w
	}
	r := Rect{X: m.X, Y: m.Y, W: int(w), H: int(h)}
	if len(m.Reserved) == 4 {
		r.X += m.Reserved[0]
		r.Y += m.Reserved[1]
		r.W -= m.Reserved[0] + m.Reserved[2]
		r.H -= m.Reserved[1] + m.Reserved[3]
	}
	return r
}

func (s Session) hyprRequest(ctx context.Context, cmd string) ([]byte, error) {
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
	if _, err := io.WriteString(c, cmd); err != nil {
		return nil, err
	}
	return io.ReadAll(io.LimitReader(c, maxReply))
}

func (s Session) hyprClient(ctx context.Context, pid int) (*hyprClient, error) {
	var found *hyprClient
	err := poll(ctx, func() (bool, error) {
		out, err := s.hyprRequest(ctx, "j/clients")
		if err != nil {
			return false, err
		}
		var clients []hyprClient
		if err := json.Unmarshal(out, &clients); err != nil {
			return false, fmt.Errorf("hyprland clients: %w", err)
		}
		for i := range clients {
			if clients[i].PID == pid && clients[i].Mapped && strings.HasPrefix(clients[i].Address, "0x") {
				found = &clients[i]
				return true, nil
			}
		}
		return false, nil
	})
	return found, err
}

func (s Session) hyprMonitor(ctx context.Context, id int) (*hyprMonitor, error) {
	out, err := s.hyprRequest(ctx, "j/monitors")
	if err != nil {
		return nil, err
	}
	var mons []hyprMonitor
	if err := json.Unmarshal(out, &mons); err != nil {
		return nil, fmt.Errorf("hyprland monitors: %w", err)
	}
	for i := range mons {
		if mons[i].ID == id {
			return &mons[i], nil
		}
	}
	return nil, fmt.Errorf("hyprland monitor %d not found", id)
}

// hyprDispatch runs dispatchers in one batch; Hyprland answers "ok" per dispatcher.
func (s Session) hyprDispatch(ctx context.Context, cmds ...string) error {
	out, err := s.hyprRequest(ctx, "[[BATCH]]dispatch "+strings.Join(cmds, ";dispatch "))
	if err != nil {
		return err
	}
	if got := strings.Count(string(out), "ok"); got != len(cmds) {
		return fmt.Errorf("hyprland dispatch: %q", strings.TrimSpace(string(out)))
	}
	return nil
}

func (s Session) hyprFloat(ctx context.Context, w Window, pref, minSize Size) error {
	c, err := s.hyprClient(ctx, w.PID)
	if err != nil {
		return err
	}
	mon, err := s.hyprMonitor(ctx, c.Monitor)
	if err != nil {
		return err
	}
	area := mon.area()
	r := centered(FitSize(pref, minSize, area), area)
	sel := "address:" + c.Address
	var cmds []string
	if !c.Floating {
		cmds = append(cmds, "setfloating "+sel)
	}
	cmds = append(cmds,
		fmt.Sprintf("resizewindowpixel exact %d %d,%s", r.W, r.H, sel),
		fmt.Sprintf("movewindowpixel exact %d %d,%s", r.X, r.Y, sel))
	return s.hyprDispatch(ctx, cmds...)
}

func (s Session) hyprTile(ctx context.Context, w Window) error {
	c, err := s.hyprClient(ctx, w.PID)
	if err != nil || !c.Floating {
		return err
	}
	return s.hyprDispatch(ctx, "settiled address:"+c.Address)
}
