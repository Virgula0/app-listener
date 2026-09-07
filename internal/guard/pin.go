package guard

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cilium/ebpf/link"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

// pinFilePrefix tags every pin file this daemon creates so CleanupStalePins
// never touches another tool's pins sharing the bpffs mount.
const pinFilePrefix = "al-"

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
	// Prefer one tidy subdirectory; fall back to the mount root if the kernel
	// refuses even a single mkdir on bpffs.
	sub := filepath.Join(mountpoint, "app-listener")
	if mkErr := os.Mkdir(sub, 0o700); mkErr == nil || errors.Is(mkErr, os.ErrExist) {
		return sub
	}
	log.Warnf("guard: cannot create %s — pinning LSM links directly under %s", sub, mountpoint)
	return mountpoint
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
