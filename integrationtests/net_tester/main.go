package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

// connFD returns the raw socket fd of a net.Conn (TCP or UDP).
func connFD(conn net.Conn) int {
	type syscallConn interface {
		SyscallConn() (syscall.RawConn, error)
	}
	sc, ok := conn.(syscallConn)
	if !ok {
		return -1
	}
	raw, err := sc.SyscallConn()
	if err != nil {
		return -1
	}
	fd := -1
	if err := raw.Control(func(f uintptr) { fd = int(f) }); err != nil {
		return -1
	}
	return fd
}

func tcpServer(port string) {
	l, err := (&net.ListenConfig{}).Listen(context.Background(), "tcp", "0.0.0.0:"+port)
	if err != nil {
		fmt.Printf("listen error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("READY")
	conn, err := l.Accept()
	if err != nil {
		fmt.Printf("accept error: %v\n", err)
		os.Exit(1)
	}
	fd := connFD(conn)
	// Use sendto/recvfrom so the sendmsg/sendto and recvfrom/recvmsg
	// tracepoints in network-monitor fire (plain write(2) is not traced).
	_ = unix.Sendto(fd, []byte("ok\n"), 0, nil)
	buf := make([]byte, 64)
	_, _, _ = unix.Recvfrom(fd, buf, 0)
	conn.Close()
	l.Close()
}

// tcpServerDelayed binds first, then listens once the marker file appears.
// Used to observe a LISTEN denied by the guard on a socket that was already
// bound (and therefore allowed) before the guard attached.
func tcpServerDelayed(port, waitFile string) error {
	sock, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, 0)
	if err != nil {
		return fmt.Errorf("socket error: %w", err)
	}
	defer unix.Close(sock)

	var addr unix.SockaddrInet4
	addr.Port, err = strconv.Atoi(port)
	if err != nil || addr.Port <= 0 {
		return fmt.Errorf("invalid port: %s", port)
	}
	addr.Addr = [4]byte{0, 0, 0, 0}

	if err := unix.Bind(sock, &addr); err != nil {
		return fmt.Errorf("bind error: %w", err)
	}
	fmt.Println("BOUND")

	// Wait until the guard is attached (test creates the marker), then listen.
	waitForMarker(waitFile)

	if err := unix.Listen(sock, 16); err != nil {
		return fmt.Errorf("listen error: %w", err)
	}
	fmt.Println("LISTENING")
	time.Sleep(2 * time.Second)
	return nil
}

func tcpClient(host, port string) {
	var d net.Dialer
	d.Timeout = 2 * time.Second
	conn, err := d.DialContext(context.Background(), "tcp", host+":"+port)
	if err != nil {
		fmt.Printf("dial error: %v\n", err)
		os.Exit(1)
	}
	fd := connFD(conn)
	_ = unix.Sendto(fd, []byte("hello\n"), 0, nil)
	buf := make([]byte, 64)
	_, _, _ = unix.Recvfrom(fd, buf, 0)
	conn.Close()
}

func udpServer(port string) {
	conn, err := (&net.ListenConfig{}).ListenPacket(context.Background(), "udp", "0.0.0.0:"+port)
	if err != nil {
		fmt.Printf("listen error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("READY")
	buf := make([]byte, 1024)
	n, addr, err := conn.ReadFrom(buf)
	if err == nil && n > 0 {
		_, _ = conn.WriteTo([]byte("ok\n"), addr)
	}
	conn.Close()
}

// waitForMarker polls until the given file exists (or a timeout). Tests use a
// marker file to time operations relative to the guard attaching its hooks.
func waitForMarker(path string) {
	deadline := time.Now().Add(60 * time.Second)
	for {
		if _, err := os.Stat(path); err == nil {
			return
		}
		if time.Now().After(deadline) {
			fmt.Printf("waitfile timeout: %s\n", path)
			os.Exit(1)
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// udpRecvLoop binds and keeps calling recvmsg in non-blocking mode. Started
// before the guard, its recvmsg calls are denied once the guard attaches.
func udpRecvLoop(port, waitFile string) error {
	sock, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return fmt.Errorf("socket error: %w", err)
	}
	defer func() { _ = unix.Close(sock) }()

	var addr unix.SockaddrInet4
	addr.Port, err = strconv.Atoi(port)
	if err != nil || addr.Port <= 0 {
		return fmt.Errorf("invalid port: %s", port)
	}
	if err := unix.Bind(sock, &addr); err != nil {
		return fmt.Errorf("bind error: %w", err)
	}
	fmt.Println("BOUND")

	// Wait until the guard is attached (test creates the marker), then call
	// recvmsg a few times so the blocked RECV events are logged.
	waitForMarker(waitFile)

	buf := make([]byte, 1024)
	for i := 0; i < 10; i++ {
		_, _, _ = unix.Recvfrom(sock, buf, unix.MSG_DONTWAIT)
		time.Sleep(50 * time.Millisecond)
	}
	_ = buf
	return nil
}

// udpSendLoop dials a UDP socket then keeps sending datagrams on an interval.
// Started before the guard, the sends are denied once the guard attaches.
func udpSendLoop(host, port, waitFile string, intervalMs int) {
	var d net.Dialer
	conn, err := d.DialContext(context.Background(), "udp", host+":"+port)
	if err != nil {
		fmt.Printf("dial error: %v\n", err)
		os.Exit(1)
	}
	fd := connFD(conn)

	// Wait until the guard is attached (test creates the marker), then send a
	// few datagrams so the blocked SEND events are logged.
	waitForMarker(waitFile)

	for i := 0; i < 10; i++ {
		_ = unix.Sendto(fd, []byte("hello\n"), 0, nil)
		time.Sleep(time.Duration(intervalMs) * time.Millisecond)
	}
	conn.Close()
}

func udpClient(host, port string) {
	var d net.Dialer
	d.Timeout = 2 * time.Second
	conn, err := d.DialContext(context.Background(), "udp", host+":"+port)
	if err != nil {
		fmt.Printf("dial error: %v\n", err)
		os.Exit(1)
	}
	_, err = conn.Write([]byte("hello\n"))
	if err != nil {
		fmt.Printf("write error: %v\n", err)
	}
	buf := make([]byte, 1024)
	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, _ = conn.Read(buf)
	conn.Close()
}

func unixServer(path string) {
	l, err := (&net.ListenConfig{}).Listen(context.Background(), "unix", path)
	if err != nil {
		fmt.Printf("listen error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("READY")
	conn, err := l.Accept()
	if err == nil {
		conn.Close()
	}
	l.Close()
}

// needArgs exits with code 1 unless len(args) >= n.
func needArgs(args []string, n int) {
	if len(args) < n {
		os.Exit(1)
	}
}

// commands maps the CLI subcommand to its handler.
var commands = map[string]func(args []string) int{
	"tcp-server": func(args []string) int {
		needArgs(args, 1)
		tcpServer(args[0])
		return 0
	},
	"tcp-server-delayed": func(args []string) int {
		needArgs(args, 2)
		if err := tcpServerDelayed(args[0], args[1]); err != nil {
			fmt.Printf("%v\n", err)
			return 1
		}
		return 0
	},
	"tcp-client": func(args []string) int {
		needArgs(args, 2)
		tcpClient(args[0], args[1])
		return 0
	},
	"udp-server": func(args []string) int {
		needArgs(args, 1)
		udpServer(args[0])
		return 0
	},
	"udp-client": func(args []string) int {
		needArgs(args, 2)
		udpClient(args[0], args[1])
		return 0
	},
	"udp-recv-loop": func(args []string) int {
		needArgs(args, 2)
		if err := udpRecvLoop(args[0], args[1]); err != nil {
			fmt.Printf("%v\n", err)
			return 1
		}
		return 0
	},
	"udp-send-loop": func(args []string) int {
		needArgs(args, 4)
		intervalMs := 0
		_, _ = fmt.Sscanf(args[3], "%d", &intervalMs)
		udpSendLoop(args[0], args[1], args[2], intervalMs)
		return 0
	},
	"unix-server": func(args []string) int {
		needArgs(args, 1)
		unixServer(args[0])
		return 0
	},
	"dns": func(args []string) int {
		needArgs(args, 1)
		r := &net.Resolver{}
		if _, err := r.LookupHost(context.Background(), args[0]); err != nil {
			fmt.Printf("dns error: %v\n", err)
			return 1
		}
		return 0
	},
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing command")
		os.Exit(1)
	}

	run, ok := commands[os.Args[1]]
	if !ok {
		fmt.Printf("unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
	os.Exit(run(os.Args[2:]))
}
