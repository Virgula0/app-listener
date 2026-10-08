package ebpf

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

// MulticallKind is how an executable picks the program it runs.
type MulticallKind uint8

const (
	// SingleBinary runs one program: its identity is its inode.
	SingleBinary MulticallKind = iota
	// UutilsMulticall is uutils coreutils: it runs the applet named by the basename it was exec'd as,
	// which the guard attests at exec (guard_exec_applet). Applets lists the build's own.
	UutilsMulticall
	// OpaqueMulticall runs applets the guard can't tell apart (busybox and toybox also dispatch
	// in-process, without an exec; a uutils build that won't list its applets; any other binary
	// calling itself a multi-call binary, e.g. a future uutils diffutils, until it is vetted).
	OpaqueMulticall
)

// MulticallNameMax bounds an applet name, NUL included: the kernel's name key (guard.bpf.c).
const MulticallNameMax = 16

// multicallScanMax: multicall binaries are small (uutils ~11 MiB, busybox ~1 MiB); bigger images
// (Electron, Chromium) are never scanned.
const multicallScanMax = 64 << 20

// Multicall is the classification of one executable inode.
type Multicall struct {
	Kind    MulticallKind
	Family  string   // "uutils", "busybox", "toybox", FamilyUnrecognized; empty for SingleBinary
	Applets []string // UutilsMulticall only: the build's `--list` that fits MulticallNameMax
}

// multicallMarkers are stored reversed: spelled forward, they would be in this binary too, which
// would then classify itself (and any Go binary linking this package) as a multicall.
var multicallMarkers = []struct {
	family string
	all    [][]byte
}{
	{"uutils", [][]byte{reversed(")yranib llac-itlum("), reversed(")slitueroc slituu(")}},
	{"busybox", [][]byte{reversed("v xoBysuB")}},
	{"toybox", [][]byte{reversed("gnol--[ xobyot")}},
}

// FamilyUnrecognized: the binary says it is a multi-call binary but matches no vetted family. It
// is refused (fail closed): a new multicall family must be vetted before its applets get identities.
const FamilyUnrecognized = "unrecognized"

// genericMulticallMarker is what uutils projects print in their usage text; matched after the
// vetted families, so only an unvetted multicall falls through to it.
var genericMulticallMarker = reversed("yranib llac-itlum")

func reversed(s string) []byte {
	b := []byte(s)
	slices.Reverse(b)
	return b
}

type multicallCacheKey struct {
	dev, ino     uint64
	size         int64
	mtime, ctime unix.Timespec
}

var (
	multicallMu    sync.Mutex
	multicallCache = map[multicallCacheKey]Multicall{}
)

// ClassifyMulticall classifies the executable inode f holds (an O_PATH fd is enough). Scanned once
// per inode and content fingerprint. An error is not a SingleBinary: callers fail closed.
func ClassifyMulticall(f *os.File) (Multicall, error) {
	var st unix.Stat_t
	if err := unix.Fstat(int(f.Fd()), &st); err != nil {
		return Multicall{}, &os.PathError{Op: "fstat", Path: f.Name(), Err: err}
	}
	if st.Mode&unix.S_IFMT != unix.S_IFREG || st.Size > multicallScanMax {
		return Multicall{Kind: SingleBinary}, nil
	}
	key := multicallCacheKey{dev: st.Dev, ino: st.Ino, size: st.Size, mtime: st.Mtim, ctime: st.Ctim}
	multicallMu.Lock()
	m, ok := multicallCache[key]
	multicallMu.Unlock()
	if ok {
		return m, nil
	}
	family, err := scanMulticallFamily(f)
	if err != nil {
		return Multicall{}, err
	}
	switch family {
	case "":
		m = Multicall{Kind: SingleBinary}
	case "uutils":
		m = Multicall{Kind: OpaqueMulticall, Family: family}
		if !rootPlacedFile(f) {
			break
		}
		if applets, lerr := listUutilsApplets(f); lerr == nil && len(applets) > 0 {
			m = Multicall{Kind: UutilsMulticall, Family: family, Applets: applets}
		}
	default:
		m = Multicall{Kind: OpaqueMulticall, Family: family}
	}
	multicallMu.Lock()
	if len(multicallCache) > 4096 {
		multicallCache = map[multicallCacheKey]Multicall{}
	}
	multicallCache[key] = m
	multicallMu.Unlock()
	return m, nil
}

func scanMulticallFamily(f *os.File) (string, error) {
	r, err := os.Open(fmt.Sprintf("/proc/self/fd/%d", f.Fd()))
	if err != nil {
		return "", err
	}
	defer r.Close()
	found := make([][]bool, len(multicallMarkers))
	for i, m := range multicallMarkers {
		found[i] = make([]bool, len(m.all))
	}
	generic := false
	const chunk, overlap = 1 << 20, 32
	buf := make([]byte, chunk+overlap)
	carry := 0
	for {
		n, rerr := io.ReadFull(r, buf[carry:])
		window := buf[:carry+n]
		for i, m := range multicallMarkers {
			for j, marker := range m.all {
				found[i][j] = found[i][j] || bytes.Contains(window, marker)
			}
		}
		generic = generic || bytes.Contains(window, genericMulticallMarker)
		if errors.Is(rerr, io.EOF) || errors.Is(rerr, io.ErrUnexpectedEOF) {
			break
		}
		if rerr != nil {
			return "", rerr
		}
		carry = copy(buf, window[len(window)-overlap:])
	}
	return familyOf(found, generic), nil
}

// familyOf is the first vetted family whose markers were all found, else FamilyUnrecognized for a
// binary that only calls itself a multi-call binary, else "".
func familyOf(found [][]bool, generic bool) string {
	for i, m := range multicallMarkers {
		if !slices.Contains(found[i], false) {
			return m.family
		}
	}
	if generic {
		return FamilyUnrecognized
	}
	return ""
}

// rootPlacedFile: only root could have placed the file f holds (SystemPlacedInode), checked before
// listUutilsApplets executes it. Anything else stays OpaqueMulticall: refused, never run.
func rootPlacedFile(f *os.File) bool {
	dev, ino, err := StatFile(f)
	if err != nil {
		return false
	}
	resolved, err := os.Readlink(fmt.Sprintf("/proc/self/fd/%d", f.Fd()))
	return err == nil && SystemPlacedInode(resolved, dev, ino)
}

// listUutilsApplets runs the build's own `coreutils --list`: its link names differ from its applet
// set (Ubuntu ships sha3sum links that run sum). Run through /proc/self/fd so uutils takes the applet
// from argv[0] (`coreutils` = argv[1] dispatch), as nobody with an empty environment.
func listUutilsApplets(f *os.File) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "/proc/self/fd/3", "--list")
	cmd.Args[0] = "coreutils"
	cmd.Env = []string{}
	cmd.Dir = "/"
	cmd.ExtraFiles = []*os.File{f}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{Uid: 65534, Gid: 65534, NoSetGroups: false},
		Setsid:     true,
		Pdeathsig:  syscall.SIGKILL,
	}
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("listing uutils applets: %w", err)
	}
	if len(out) > 64<<10 {
		return nil, errors.New("listing uutils applets: output too large")
	}
	var applets []string
	seen := map[string]bool{}
	sc := bufio.NewScanner(bytes.NewReader(out))
	for sc.Scan() {
		name := sc.Text()
		if !validAppletName(name) || seen[name] {
			continue
		}
		seen[name] = true
		applets = append(applets, name)
	}
	return applets, nil
}

// validAppletName: what the kernel can match (fits its key, no path separator) and a name uutils
// dispatches on as-is (no '.': file_stem would cut it).
func validAppletName(name string) bool {
	if name == "" || len(name) >= MulticallNameMax || name == "coreutils" {
		return false
	}
	for _, c := range []byte(name) {
		ok := c >= 'a' && c <= 'z' || c >= '0' && c <= '9' || c == '-' || c == '_' || c == '['
		if !ok {
			return false
		}
	}
	return true
}
