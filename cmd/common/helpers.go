package common

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/charmbracelet/x/term"

	"github.com/Virgula0/app-listener/internal/infrastructure"
)

const defaultServeAddress = "127.0.0.1:9999"

// ServeConfig is the validated CLI configuration for the browser TUI.
type ServeConfig struct {
	Enabled  bool
	Address  string
	Username string
	Password string
}

// AddServeFlags adds the browser TUI flags to commands that expose a TUI.
func AddServeFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.String("serve", "", "Serve the TUI in a browser at host:port (default when omitted: 127.0.0.1:9999)")
	flags.Lookup("serve").NoOptDefVal = defaultServeAddress
	flags.String("user", "", "HTTP Basic Auth username (requires --serve and --password)")
	flags.String("password", "", "HTTP Basic Auth password (requires --serve and --user)")
}

// ParseServeFlags validates flags shared by the five monitoring and guarding modes.
func ParseServeFlags(cmd *cobra.Command) (ServeConfig, error) {
	serveFlag := cmd.Flags().Lookup("serve")
	if serveFlag == nil {
		return ServeConfig{}, nil
	}

	config := ServeConfig{Enabled: serveFlag.Changed}
	config.Address, _ = cmd.Flags().GetString("serve")
	config.Username, _ = cmd.Flags().GetString("user")
	config.Password, _ = cmd.Flags().GetString("password")

	credentialsSet, _ := credentialFlagState(cmd)
	if !config.Enabled {
		if credentialsSet {
			return ServeConfig{}, errors.New("--user and --password require --serve")
		}
		return config, nil
	}

	if err := validateEnabledServe(cmd, config); err != nil {
		return ServeConfig{}, err
	}

	host, _, _ := net.SplitHostPort(config.Address)
	if !isLoopbackHost(host) {
		log.Warn("browser TUI is bound outside loopback; HTTP traffic and Basic Auth credentials are not encrypted, use a TLS reverse proxy")
	}
	return config, nil
}

// validateEnabledServe enforces every rule that only applies when --serve was
// actually given: presentation conflicts, credential pairing, address syntax
// and the interactive-terminal requirement of the dual local+browser TUI.
func validateEnabledServe(cmd *cobra.Command, config ServeConfig) error {
	if requestedBoolFlag(cmd, "headless") {
		return errors.New("--serve and --headless are mutually exclusive")
	}
	if requestedBoolFlag(cmd, "gui") {
		return errors.New("--serve and --gui are mutually exclusive")
	}
	if requestedBoolFlag(cmd, "genkey") {
		return errors.New("--serve is not available with --genkey")
	}
	if requestedBoolFlag(cmd, "blocked-only") {
		return errors.New("--blocked-only is only available with --headless")
	}
	credentialsSet, mixedCredentials := credentialFlagState(cmd)
	if credentialsSet && (mixedCredentials || config.Username == "" || config.Password == "") {
		return errors.New("--user and --password must be non-empty and specified together")
	}
	if err := validateServeAddress(config.Address); err != nil {
		return err
	}
	if !term.IsTerminal(os.Stdin.Fd()) {
		return errors.New("--serve needs an interactive terminal for the local TUI; use --headless without --serve for services")
	}
	return nil
}

// credentialFlagState reports whether any credential flag was set explicitly
// and whether the two flags disagree (one given without the other).
func credentialFlagState(cmd *cobra.Command) (set, mixed bool) {
	userFlag := cmd.Flags().Lookup("user")
	passwordFlag := cmd.Flags().Lookup("password")
	if userFlag == nil || passwordFlag == nil {
		return false, false
	}
	return userFlag.Changed || passwordFlag.Changed, userFlag.Changed != passwordFlag.Changed
}

func requestedBoolFlag(cmd *cobra.Command, name string) bool {
	flag := cmd.Flags().Lookup(name)
	return flag != nil && flag.Value.String() == "true"
}

func validateServeAddress(address string) error {
	host, portText, err := net.SplitHostPort(address)
	if err != nil || host == "" {
		return fmt.Errorf("invalid --serve address %q: expected host:port", address)
	}
	port, err := strconv.Atoi(portText)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invalid --serve address %q: port must be between 1 and 65535", address)
	}
	return nil
}

func isLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

type RawTarget struct {
	AbsPath string
	IsDir   bool
}

func ResolveTargets(paths []string) ([]RawTarget, error) {
	seen := make(map[string]bool)
	var targets []RawTarget

	for _, p := range paths {
		abs, err := filepath.Abs(p)
		if err != nil {
			log.Errorf("Invalid path %q: %v", p, err)
			return nil, err
		}

		if seen[abs] {
			return nil, fmt.Errorf("duplicate watch path: %s", abs)
		}
		seen[abs] = true

		info, err := os.Stat(abs)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				log.Errorf("Path does not exist: %s", abs)
			} else {
				log.Errorf("Cannot access path %s: %v", abs, err)
			}
			return nil, err
		}

		targets = append(targets, RawTarget{AbsPath: abs, IsDir: info.IsDir()})
	}

	return targets, nil
}

func ValidateTargets(targets []RawTarget) error {
	for i, a := range targets {
		for j, b := range targets {
			if i == j {
				continue
			}
			if err := validatePair(a, b); err != nil {
				return err
			}
		}
	}
	return nil
}

func validatePair(a, b RawTarget) error {
	if a.IsDir && b.IsDir {
		if isSubDir(a.AbsPath, b.AbsPath) {
			return fmt.Errorf(
				"%q is a subdirectory of %q \u2014 remove the more specific path",
				b.AbsPath, a.AbsPath)
		}
		if isSubDir(b.AbsPath, a.AbsPath) {
			return fmt.Errorf(
				"%q is a subdirectory of %q \u2014 remove the more specific path",
				a.AbsPath, b.AbsPath)
		}
	}

	if !a.IsDir && b.IsDir && pathWithinDir(a.AbsPath, b.AbsPath) {
		return fmt.Errorf(
			"%q is already covered by directory %q \u2014 remove the redundant file path",
			a.AbsPath, b.AbsPath)
	}
	if a.IsDir && !b.IsDir && pathWithinDir(b.AbsPath, a.AbsPath) {
		return fmt.Errorf(
			"%q is already covered by directory %q \u2014 remove the redundant file path",
			b.AbsPath, a.AbsPath)
	}
	return nil
}

func isSubDir(child, parent string) bool {
	return strings.HasPrefix(child+"/", parent+"/")
}

func pathWithinDir(path, dir string) bool {
	return strings.HasPrefix(path+"/", dir+"/")
}

func BuildEBPFTargets(rawTargets []RawTarget) []ebpf.Target {
	targets := make([]ebpf.Target, len(rawTargets))
	for i, rt := range rawTargets {
		t := ebpf.Target{Path: rt.AbsPath, IsDir: rt.IsDir}
		if rt.IsDir {
			t.Dir = rt.AbsPath
		} else {
			t.Dir = filepath.Dir(rt.AbsPath)
			t.File = filepath.Base(rt.AbsPath)
		}
		targets[i] = t
	}
	return targets
}

func MakeDisplayPaths(rawTargets []RawTarget) []string {
	paths := make([]string, len(rawTargets))
	for i, rt := range rawTargets {
		paths[i] = rt.AbsPath
	}
	return paths
}

func CheckEBPF() error {
	log.Info("Checking eBPF availability...")
	if err := ebpf.Check(); err != nil {
		log.Errorf("eBPF check failed: %v", err)
		return err
	}
	log.Info("eBPF available")
	return nil
}

// CheckBPFLSM verifies that the active kernel LSM stack includes the BPF
// LSM. Required by every mode that relies on LSM hooks (guard,
// network-guard, daemon); monitor and network-monitor are kprobe-based and
// must not call it.
func CheckBPFLSM() error {
	if err := ebpf.CheckBPFLSM(); err != nil {
		log.Errorf("BPF LSM check failed: %v", err)
		return err
	}
	log.Info("BPF LSM available")
	return nil
}

func ParseEventsFlag(eventsFlag []string) ([]ebpf.EventType, error) {
	if len(eventsFlag) == 0 {
		return ebpf.EventTypes(), nil
	}

	var parsed []ebpf.EventType
	for _, s := range eventsFlag {
		et, ok := ebpf.ParseEventType(strings.TrimSpace(s))
		if !ok {
			return nil, fmt.Errorf("unknown event type %q (valid: OPEN, READ, WRITE, DELETE, RENAME, SYMLINK, HARDLINK, MKDIR, MMAP, ATTR, STAT, MKNOD; for network mode: CONNECT, ACCEPT, SEND, RECV, CLOSE, DNS)", s)
		}
		parsed = append(parsed, et)
	}
	return parsed, nil
}

// ParseNetEventsFlag parses the network --events flag; an empty list means
// every network event type.
func ParseNetEventsFlag(netEventsFlag []string) ([]ebpf.NetEventType, error) {
	if len(netEventsFlag) == 0 {
		return ebpf.NetEventTypes(), nil
	}

	var parsed []ebpf.NetEventType
	for _, s := range netEventsFlag {
		et, ok := ebpf.ParseNetEventType(strings.TrimSpace(s))
		if !ok {
			return nil, fmt.Errorf("unknown network event type %q (valid: CONNECT, ACCEPT, SEND, RECV, CLOSE, DNS, BIND, LISTEN)", s)
		}
		parsed = append(parsed, et)
	}
	return parsed, nil
}

func WarnIgnoredFlags(rawTargets []RawTarget, recursive bool, depth int) {
	for _, rt := range rawTargets {
		if !rt.IsDir && recursive {
			log.Warnf("--recursive is ignored when monitoring a single file (%s)", rt.AbsPath)
		}
		if !rt.IsDir && depth > 0 {
			log.Warnf("--depth is ignored when monitoring a single file (%s)", rt.AbsPath)
		}
	}
}

// UIDResolver caches UID → username lookups to avoid per-event syscalls.
// Create once at startup and reuse for every event.
type UIDResolver struct {
	mu    sync.Mutex
	cache map[uint32]string
}

// NewUIDResolver returns a resolver ready for use.
func NewUIDResolver() *UIDResolver {
	return &UIDResolver{cache: make(map[uint32]string)}
}

// Resolve returns the username for uid, falling back to the numeric string.
func (r *UIDResolver) Resolve(uid uint32) string {
	r.mu.Lock()
	defer r.mu.Unlock()

	if name, ok := r.cache[uid]; ok {
		return name
	}

	name := strconv.FormatUint(uint64(uid), 10)
	if u, err := user.LookupId(strconv.FormatUint(uint64(uid), 10)); err == nil {
		name = u.Username
	}
	r.cache[uid] = name
	return name
}
