package daemonconfig

import (
	"fmt"
	"strconv"
	"strings"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// PortAnchorDir holds one root-owned anchor file per `[port tcp/<n>]` section, made by the daemon.
// The file is the resource's guarded root, so a port resource is judged, reloaded and pinned like a
// single-file one, while the guard's socket hooks route the port to it.
const PortAnchorDir = "/run/app-listener/ports"

// PortAnchor is port's anchor file.
func PortAnchor(port uint16) string {
	return fmt.Sprintf("%s/tcp-%d", PortAnchorDir, port)
}

// ParsePortSection recognizes a `[port tcp/<n>]` header: ok false for any other line, an error for
// a malformed one.
func ParsePortSection(line string) (port uint16, ok bool, err error) {
	if !strings.HasPrefix(line, "[port ") || !strings.HasSuffix(line, "]") {
		return 0, false, nil
	}
	spec := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "[port "), "]"))
	num, found := strings.CutPrefix(spec, "tcp/")
	if !found {
		return 0, true, fmt.Errorf("malformed section header %q: expected \"[port tcp/<number>]\" (only TCP)", line)
	}
	n, convErr := strconv.ParseUint(num, 10, 16)
	if convErr != nil || n == 0 {
		return 0, true, fmt.Errorf("malformed section header %q: port must be 1-65535", line)
	}
	return uint16(n), true, nil
}

// checkPortGroup validates a `[port]` section: binary lines whose events are CONNECT (reach the
// port) or MKNOD (listen on it), and nothing that names a file or a vault.
func checkPortGroup(g *watchGroup) error {
	if g.runtime || len(g.watchPaths) > 0 || len(g.libs) > 0 || len(g.pendingLibs) > 0 ||
		len(g.libDirs) > 0 || len(g.libBinaries) > 0 || len(g.libPending) > 0 {
		return fmt.Errorf("daemon config line %d: a port section takes only binary lines", g.lineNo)
	}
	if g.needEncryption {
		return fmt.Errorf("daemon config line %d: a port section is never encrypted: drop need_encryption", g.lineNo)
	}
	for _, rules := range [][]BinaryRule{g.binaries, g.pending} {
		for _, b := range rules {
			for _, e := range b.Events {
				if e != ebpf.EventConnect && e != ebpf.EventMknod {
					return fmt.Errorf("daemon config line %d: %s %s: a port section's events are CONNECT "+
						"(connect to the port) and MKNOD (listen on it)", g.lineNo, b.Path, e)
				}
			}
		}
	}
	return nil
}

// materializePortGroup appends the section's Resource. Its anchor isn't statted here: the daemon
// makes it before building guards, and other parsers (config edits, catalog refresh) never need it.
func materializePortGroup(cfg *Config, g *watchGroup) error {
	if err := checkPortGroup(g); err != nil {
		return err
	}
	for i := range cfg.Resources {
		if cfg.Resources[i].Port == g.port {
			return fmt.Errorf("daemon config line %d: duplicate port section: tcp/%d", g.lineNo, g.port)
		}
	}
	cfg.Resources = append(cfg.Resources, Resource{
		Path:            PortAnchor(g.port),
		Port:            g.port,
		Binaries:        g.binaries,
		PendingBinaries: g.pending,
	})
	return nil
}
