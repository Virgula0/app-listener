package daemon

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

const (
	portRangePath    = "/proc/sys/net/ipv4/ip_local_port_range"
	portReservedPath = "/proc/sys/net/ipv4/ip_local_reserved_ports"
)

// ensurePortAnchors makes each [port] section's anchor before its guard claims it: a root-owned,
// read-only, single-link file in a root-owned 0700 directory. An existing one is reused only if it
// is exactly that (the guard already denies replacing it). A port the kernel may hand out as an
// ephemeral one is refused: listen() without bind() autobinds after the listen hook has judged.
func ensurePortAnchors(cfg *daemonconfig.Config) error {
	var ports []uint16
	for i := range cfg.Resources {
		if p := cfg.Resources[i].Port; p != 0 {
			ports = append(ports, p)
		}
	}
	if len(ports) == 0 {
		return nil
	}
	lo, hi, err := ephemeralRange()
	if err != nil {
		return err
	}
	reserved := readReservedPorts()
	for _, p := range ports {
		if p >= lo && p <= hi && !reserved(p) {
			return fmt.Errorf("[port tcp/%d] is in the ephemeral range %d-%d (%s): any program could be "+
				"handed it by listen() without bind(); add it to %s or pick another port", p, lo, hi,
				portRangePath, portReservedPath)
		}
	}
	for _, dir := range []string{filepath.Dir(daemonconfig.PortAnchorDir), daemonconfig.PortAnchorDir} {
		if err := ensureRootDir(dir); err != nil {
			return err
		}
	}
	for _, p := range ports {
		if err := ensureAnchor(daemonconfig.PortAnchor(p)); err != nil {
			return err
		}
	}
	return nil
}

func ensureRootDir(dir string) error {
	if err := os.Mkdir(dir, 0o700); err != nil && !errors.Is(err, os.ErrExist) {
		return fmt.Errorf("creating %s: %w", dir, err)
	}
	var st unix.Stat_t
	if err := unix.Lstat(dir, &st); err != nil {
		return fmt.Errorf("checking %s: %w", dir, err)
	}
	if st.Mode&unix.S_IFMT != unix.S_IFDIR || st.Uid != 0 || st.Mode&0o022 != 0 {
		return fmt.Errorf("%s must be a root-owned directory nobody else can write", dir)
	}
	return nil
}

func ensureAnchor(path string) error {
	fd, err := unix.Open(path, unix.O_CREAT|unix.O_EXCL|unix.O_WRONLY|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0o400)
	if err == nil {
		return unix.Close(fd)
	}
	if !errors.Is(err, unix.EEXIST) {
		return fmt.Errorf("creating port anchor %s: %w", path, err)
	}
	var st unix.Stat_t
	if err := unix.Lstat(path, &st); err != nil {
		return fmt.Errorf("checking port anchor %s: %w", path, err)
	}
	if st.Mode&unix.S_IFMT != unix.S_IFREG || st.Uid != 0 || st.Nlink != 1 || st.Mode&0o222 != 0 {
		return fmt.Errorf("port anchor %s must be a root-owned read-only regular file with one link", path)
	}
	return nil
}

func ephemeralRange() (lo, hi uint16, err error) {
	raw, err := os.ReadFile(portRangePath)
	if err != nil {
		return 0, 0, fmt.Errorf("reading the ephemeral port range: %w", err)
	}
	f := strings.Fields(string(raw))
	if len(f) != 2 {
		return 0, 0, fmt.Errorf("unexpected %s: %q", portRangePath, raw)
	}
	l, errLo := strconv.ParseUint(f[0], 10, 16)
	h, errHi := strconv.ParseUint(f[1], 10, 16)
	if errLo != nil || errHi != nil {
		return 0, 0, fmt.Errorf("unexpected %s: %q", portRangePath, raw)
	}
	return uint16(l), uint16(h), nil
}

// readReservedPorts parses ip_local_reserved_ports ("8332,9000-9010"); unreadable reserves nothing.
func readReservedPorts() func(uint16) bool {
	raw, err := os.ReadFile(portReservedPath)
	if err != nil {
		return func(uint16) bool { return false }
	}
	return parseReservedPorts(string(raw))
}

func parseReservedPorts(raw string) func(uint16) bool {
	type span struct{ lo, hi uint64 }
	var spans []span
	for part := range strings.SplitSeq(strings.TrimSpace(raw), ",") {
		a, b, isRange := strings.Cut(part, "-")
		lo, err := strconv.ParseUint(a, 10, 16)
		if err != nil {
			continue
		}
		hi := lo
		if isRange {
			if hi, err = strconv.ParseUint(b, 10, 16); err != nil {
				continue
			}
		}
		spans = append(spans, span{lo, hi})
	}
	return func(p uint16) bool {
		for _, s := range spans {
			if uint64(p) >= s.lo && uint64(p) <= s.hi {
				return true
			}
		}
		return false
	}
}
