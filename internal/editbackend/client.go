package editbackend

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"syscall"
	"time"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
	"github.com/Virgula0/app-listener/internal/protected"
	"github.com/Virgula0/app-listener/internal/systemd"
)

// ErrClosed: the backend exited (or was closed) before answering.
var ErrClosed = errors.New("the edit-protected backend exited")

// Client is one running backend. Every call blocks until its reply; events and the daemon ending
// the session arrive on Events.
type Client struct {
	hello Hello

	cmd   *exec.Cmd
	stdin io.WriteCloser
	w     *Writer
	wmu   sync.Mutex

	mu      sync.Mutex
	nextID  uint64
	pending map[uint64]chan Message
	done    chan struct{}
	exitErr error

	events chan Message
	ended  chan string
	logs   chan string
}

// childEnv is all the backend inherits: the GUI's environment (LD_PRELOAD, GODEBUG, a planted
// PATH) never reaches a process the daemon grants write access to. SUDO_UID/SUDO_GID give new
// entries the invoking user's ownership, as under `sudo app-listener edit-protected`.
func childEnv() []string {
	env := []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin", "HOME=/root", "LANG=C.UTF-8"}
	for _, k := range []string{"SUDO_UID", "SUDO_GID"} {
		if v := os.Getenv(k); v != "" {
			if _, err := strconv.ParseUint(v, 10, 32); err == nil {
				env = append(env, k+"="+v)
			}
		}
	}
	return env
}

// verifiedBinary is the installed binary's inode when only root could have placed it and a running
// daemon executes that same inode (the one its control socket admits). A mismatch fails before
// the backend starts, so no password reaches a binary the daemon would refuse anyway.
func verifiedBinary() (dev, ino uint64, err error) {
	f, err := ebpf.OpenSystemPlaced(systemd.InstallBinaryPath)
	if err != nil {
		return 0, 0, fmt.Errorf("the installed binary %s is not usable: %w", systemd.InstallBinaryPath, err)
	}
	dev, ino, err = ebpf.StatFile(f)
	f.Close()
	if err != nil {
		return 0, 0, err
	}
	pids, err := protected.FindDaemonProcesses("/proc")
	if err != nil {
		return 0, 0, fmt.Errorf("scanning for the daemon: %w", err)
	}
	for _, pid := range pids {
		d, i, err := ebpf.StatInode(fmt.Sprintf("/proc/%d/exe", pid))
		if err != nil {
			continue // exited since the scan
		}
		if d != dev || i != ino {
			return 0, 0, fmt.Errorf("the running daemon (pid %d) runs a different binary than %s (updated "+
				"without a restart?): restart it with systemctl restart %s", pid, systemd.InstallBinaryPath,
				systemd.DaemonServiceName)
		}
	}
	return dev, ino, nil
}

// Start runs `<installed binary> edit-protected --backend` and exchanges hello. Pdeathsig and the
// stdin pipe (closed when this process dies) both end the session if the window goes away.
func Start() (*Client, error) {
	dev, ino, err := verifiedBinary()
	if err != nil {
		return nil, err
	}
	//nolint:gosec // fixed path, inode-checked before and after the exec
	cmd := exec.CommandContext(context.Background(), systemd.InstallBinaryPath, "edit-protected", "--backend")
	cmd.Env = childEnv()
	cmd.Dir = "/"
	cmd.SysProcAttr = &syscall.SysProcAttr{Pdeathsig: syscall.SIGTERM}
	// Start returns after the exec: the child's image is final. Check it, not the path.
	return launch(cmd, func(pid int) error {
		if d, i, err := ebpf.StatInode(fmt.Sprintf("/proc/%d/exe", pid)); err != nil || d != dev || i != ino {
			return fmt.Errorf("%s was replaced while starting it", systemd.InstallBinaryPath)
		}
		return nil
	})
}

// launch starts cmd, runs verify on its pid before anything is sent, and exchanges hello.
func launch(cmd *exec.Cmd, verify func(pid int) error) (*Client, error) {
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("starting %s: %w", cmd.Path, err)
	}
	if err := verify(cmd.Process.Pid); err != nil {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		return nil, err
	}
	c := &Client{
		cmd: cmd, stdin: stdin, w: NewWriter(stdin), pending: map[uint64]chan Message{},
		done: make(chan struct{}), events: make(chan Message, 256), ended: make(chan string, 1),
		logs: make(chan string, 256),
	}
	go c.readLogs(stderr)
	go c.readMessages(stdout)
	if err := c.handshake(); err != nil {
		c.Close()
		return nil, err
	}
	return c, nil
}

func (c *Client) handshake() error {
	m, err := c.Call(&Request{Op: OpHello})
	if err != nil {
		return err
	}
	if m.Hello == nil || m.Hello.Proto != Proto {
		return fmt.Errorf("the installed binary speaks another edit-protected protocol than this window: " +
			"update it (app-listener update) or use a window of the same release")
	}
	c.hello = *m.Hello
	return nil
}

// Hello is the backend's mode report.
func (c *Client) Hello() Hello { return c.hello }

func (c *Client) readLogs(r io.Reader) {
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		select {
		case c.logs <- sc.Text():
		default: // a slow display drops lines, never stalls the backend
		}
	}
}

func (c *Client) readMessages(r io.Reader) {
	in := NewReader(r)
	for {
		var m Message
		if err := in.Read(&m); err != nil {
			break
		}
		switch m.Kind {
		case KindEnded:
			select {
			case c.ended <- m.Err:
			default: // sent once by the backend
			}
			continue
		case KindEvent:
			select {
			case c.events <- m:
			default: // a slow display drops events, never stalls the reader
			}
			continue
		}
		c.mu.Lock()
		ch := c.pending[m.ID]
		delete(c.pending, m.ID)
		c.mu.Unlock()
		if ch != nil {
			ch <- m
		}
	}
	err := c.cmd.Wait()
	c.mu.Lock()
	c.exitErr = err
	for id, ch := range c.pending {
		close(ch)
		delete(c.pending, id)
	}
	c.mu.Unlock()
	close(c.done)
}

// Call sends req and waits for its reply; a reply carrying an error returns it as err, with the
// message (some replies carry data alongside, e.g. a stale config).
func (c *Client) Call(req *Request) (Message, error) {
	ch := make(chan Message, 1)
	c.mu.Lock()
	select {
	case <-c.done:
		c.mu.Unlock()
		return Message{}, ErrClosed
	default:
	}
	c.nextID++
	r := *req
	r.ID = c.nextID
	c.pending[r.ID] = ch
	c.mu.Unlock()

	c.wmu.Lock()
	err := c.w.Write(&r)
	c.wmu.Unlock()
	if err != nil {
		c.mu.Lock()
		delete(c.pending, r.ID)
		c.mu.Unlock()
		return Message{}, fmt.Errorf("%w: %v", ErrClosed, err)
	}
	m, ok := <-ch
	if !ok {
		return Message{}, ErrClosed
	}
	if m.Err != "" {
		return m, errors.New(m.Err)
	}
	return m, nil
}

// Events streams a FORWARD grant's events.
func (c *Client) Events() <-chan Message { return c.events }

// Ended yields the reason once the daemon ended the session (idle timeout, reload, a granted
// binary changed).
func (c *Client) Ended() <-chan string { return c.ended }

// Logs streams the backend's log lines.
func (c *Client) Logs() <-chan string { return c.logs }

// Done closes when the backend exited.
func (c *Client) Done() <-chan struct{} { return c.done }

// Close ends the backend: EOF on its stdin narrows the grant (or re-locks the vault) and exits it;
// one still running after the grace period is killed.
func (c *Client) Close() {
	_ = c.stdin.Close()
	select {
	case <-c.done:
	case <-time.After(30 * time.Second):
		_ = c.cmd.Process.Signal(syscall.SIGTERM)
		select {
		case <-c.done:
		case <-time.After(10 * time.Second):
			_ = c.cmd.Process.Kill()
			<-c.done
		}
	}
}
