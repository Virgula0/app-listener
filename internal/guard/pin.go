package guard

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	cilium "github.com/cilium/ebpf"
	"github.com/cilium/ebpf/link"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

// pinFilePrefix tags every pin file this daemon creates so CleanupStalePins
// never touches another tool's pins sharing the bpffs mount.
const pinFilePrefix = "al-"

// pinSubdir is the one directory this guard creates under the bpffs mount to
// keep its pins tidy; a pre-existing one is reused only when dirIsRootOwnedSafe.
const pinSubdir = "app-listener"

// ResolvePinBase makes link pinning usable and returns the directory pin
// files go in — or "" when this host cannot support pinning, in which case
// the caller runs with live-only enforcement (no SIGKILL survival) and should
// say so loudly. mountpoint must be a bpffs mount (link pins only survive
// process death on bpffs); a systemd host has /sys/fs/bpf mounted very early,
// so this mounts it only where it is missing (a container, a minimal init).
//
// Pin files are kept FLAT — at most one shallow subdirectory, names limited
// to [a-z0-9-]: some hardened kernels reject deep mkdir or unusual characters
// on bpffs, and a flat BPF_OBJ_PIN is the lowest common denominator every
// bpffs supports.
func ResolvePinBase(mountpoint string) string {
	onBpffs, err := isBpffs(mountpoint)
	if err != nil {
		log.Errorf("guard: CRITICAL: cannot inspect %s (%v) — LSM link pinning disabled, guards will NOT "+
			"survive a SIGKILL", mountpoint, err)
		return ""
	}
	if !onBpffs {
		if mErr := mountBpffs(mountpoint); mErr != nil {
			log.Errorf("guard: CRITICAL: %v — LSM link pinning disabled, guards will NOT survive a SIGKILL. "+
				"Mount it yourself: mount -t bpf bpffs %s", mErr, mountpoint)
			return ""
		}
	}
	// Prefer one tidy subdirectory; fall back to the mount root (kernel-created,
	// not adoptable by a local user) if the kernel refuses even a single mkdir
	// on bpffs, or if a subdirectory is already there but is not a plain
	// root-owned, non-world-writable directory — a local user who reached the
	// bpffs mount must not be able to seed or observe our pins.
	sub := filepath.Join(mountpoint, pinSubdir)
	switch mkErr := os.Mkdir(sub, 0o700); {
	case mkErr == nil:
		return sub
	case errors.Is(mkErr, os.ErrExist):
		if safeErr := dirIsRootOwnedSafe(sub); safeErr != nil {
			log.Warnf("guard: not reusing %s (%v) — pinning LSM links directly under %s", sub, safeErr, mountpoint)
			return mountpoint
		}
		return sub
	default:
		log.Warnf("guard: cannot create %s (%v) — pinning LSM links directly under %s", sub, mkErr, mountpoint)
		return mountpoint
	}
}

// dirIsRootOwnedSafe returns why path is unsafe to reuse as the pin base, or
// nil when it is a real directory owned by uid 0 with no group/other write.
func dirIsRootOwnedSafe(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return errors.New("it is a symlink")
	}
	if !info.IsDir() {
		return errors.New("it is not a directory")
	}
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return errors.New("ownership is unreadable")
	}
	if st.Uid != 0 {
		return fmt.Errorf("owned by uid %d, not root", st.Uid)
	}
	if perm := info.Mode().Perm(); perm&0o022 != 0 {
		return fmt.Errorf("mode %04o is group/other-writable", perm)
	}
	return nil
}

func mountBpffs(mountpoint string) error {
	log.Infof("guard: %s is not a bpffs mount — mounting one for LSM link pinning", mountpoint)
	if mkErr := os.MkdirAll(mountpoint, 0o755); mkErr != nil {
		return fmt.Errorf("creating bpffs mountpoint %s: %w", mountpoint, mkErr)
	}
	if mErr := unix.Mount("bpffs", mountpoint, "bpf", 0, ""); mErr != nil {
		return fmt.Errorf("mounting bpffs at %s: %w", mountpoint, mErr)
	}
	if onBpffs, err := isBpffs(mountpoint); err != nil || !onBpffs {
		return fmt.Errorf("%s is still not bpffs after mounting (err=%v)", mountpoint, err)
	}
	return nil
}

func isBpffs(path string) (bool, error) {
	var st unix.Statfs_t
	if err := unix.Statfs(path, &st); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil // not there yet: caller will create + mount
		}
		return false, fmt.Errorf("statfs %s: %w", path, err)
	}
	return st.Type == unix.BPF_FS_MAGIC, nil
}

// PinPrefix returns the flat pin-file prefix for one guard: every hook of
// that guard is pinned at PinPrefix(...)+<hook>. The name encodes the daemon
// generation (so CleanupStalePins can tell live pins from a killed
// predecessor's) and a hash of the resource path (so two guards never
// collide). Deliberately all lowercase [a-z0-9-]: no dots, no underscores,
// no nesting — shapes some hardened bpffs implementations reject.
func PinPrefix(base, generation, resourcePath string) string {
	if base == "" {
		return "" // pinning disabled on this host
	}
	sum := sha256.Sum256([]byte(resourcePath))
	return filepath.Join(base, pinFilePrefix+sanitizeGen(generation)+"-"+hex.EncodeToString(sum[:6])+"-")
}

// genOfPinFile extracts the generation field from a pin file name, or "" when
// the name is not one of ours.
func genOfPinFile(name string) string {
	if !strings.HasPrefix(name, pinFilePrefix) {
		return ""
	}
	parts := strings.SplitN(strings.TrimPrefix(name, pinFilePrefix), "-", 3)
	if len(parts) < 3 {
		return ""
	}
	return parts[0]
}

func sanitizeGen(g string) string {
	var b strings.Builder
	for _, r := range g {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		}
	}
	if b.Len() == 0 {
		return "0"
	}
	return b.String()
}

// CleanupStalePins removes every pin file this daemon left in base whose
// generation is not in keep: it unpins each (detaching the still-running LSM
// program a killed predecessor left attached) and closes it. Used at daemon
// start and after a reload. A missing or empty base is not an error.
func CleanupStalePins(base string, keep map[string]bool) (removed int, err error) {
	if base == "" {
		return 0, nil
	}
	entries, readErr := os.ReadDir(base)
	if readErr != nil {
		if errors.Is(readErr, os.ErrNotExist) {
			return 0, nil
		}
		return 0, fmt.Errorf("reading pin base %s: %w", base, readErr)
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		gen := genOfPinFile(e.Name())
		if gen == "" || keep[gen] {
			continue
		}
		path := filepath.Join(base, e.Name())
		if l, loadErr := link.LoadPinnedLink(path, nil); loadErr == nil {
			if reason := foreignPinnedLink(l); reason != "" {
				log.Warnf("guard: pin file %s holds a %s, not one of our LSM links — leaving it in place", path, reason)
				_ = l.Close()
				continue
			}
			if unpinErr := l.Unpin(); unpinErr != nil {
				log.Warnf("guard: unpinning stale link %s failed: %v", path, unpinErr)
			}
			_ = l.Close()
		} else if rmErr := os.Remove(path); rmErr != nil {
			log.Warnf("guard: removing stale pin file %s failed: %v", path, rmErr)
			continue
		}
		removed++
	}
	if removed > 0 {
		log.Warnf("guard: retired %d stale LSM link pin(s) from a previous daemon instance", removed)
	}
	return removed, nil
}

// foreignPinnedLink returns a short description when l is demonstrably NOT one
// of this guard's LSM links — a non-LSM program, or an LSM program outside our
// guard_* namespace — so CleanupStalePins leaves it untouched even though its
// file name matched the al- pattern. It is deliberately conservative: when the
// link or its program cannot be introspected (older kernel, missing
// capability) it returns "" and the caller falls back to the file-name match,
// which is already app-listener specific.
func foreignPinnedLink(l link.Link) string {
	info, err := l.Info()
	if err != nil {
		return ""
	}
	prog, err := cilium.NewProgramFromID(info.Program)
	if err != nil {
		return ""
	}
	defer prog.Close()
	pInfo, err := prog.Info()
	if err != nil {
		return ""
	}
	if pInfo.Type != cilium.LSM {
		return fmt.Sprintf("%s program", pInfo.Type)
	}
	if !strings.HasPrefix(pInfo.Name, "guard_") {
		return fmt.Sprintf("foreign LSM program %q", pInfo.Name)
	}
	return ""
}
