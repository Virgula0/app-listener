package daemon

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

// pinStateFile records which pin generation (and PID) this process's CONFIG guards (the [watch]
// sections from buildGuards, never selfguards.go's, whose generation never rotates) are pinned
// under. A separate process with no live *guard.Guard (`daemon --lockdown`, the ExecStopPost safety
// net) reads it to recompute the bpffs pin paths (guard.PinPrefix / WithPinnedSelfVaultAccess) a
// crashed daemon left, and to tell whether the owner might still be alive before touching them.
//
// Lives beside fscrypt.key and edit-auth.hash in the self-guarded /etc/app-listener (ModeReadOnly,
// selfProtectSpecs). pid+gen aren't secret, but their INTEGRITY matters, which the RO self-guard
// provides once a daemon runs (only the app-listener binary may overwrite it).
const pinStateFileProd = "/etc/app-listener/pin-state.json"

// pinStateFile is pinStateFileProd in production; tests point it at a temp
// file (mirroring cmd/functions/editprotected/auth.go's hashFilePath).
var pinStateFile = pinStateFileProd

// pinStateRecord is the on-disk shape of pinStateFile.
type pinStateRecord struct {
	PID       int       `json:"pid"`
	Gen       string    `json:"gen"`
	UpdatedAt time.Time `json:"updated_at"`
}

// errNoPinState: no trustworthy pin-generation record (pre-upgrade daemon, fresh install never
// started, malformed/empty file). `daemon --lockdown` then does a plain, unwidened lock attempt.
var errNoPinState = errors.New("no usable pin-state record found")

// writePinState records pin's generation for the config guards just brought up or rotated to
// (reload), IN PLACE on the file's existing inode (like editprotected/auth.go's WriteHashFile, same
// reason): once the RO self-guard on /etc/app-listener is attached, only rewriting an EXISTING
// entry is permitted. The first call per host runs before that guard exists (runDaemon: right after
// startGuardedDaemonAbortable succeeds, before `go sg.attach(pin)`), so the create fallback runs at
// most once. Best-effort: failure only means `--lockdown` can't recover this generation, never that
// startup/reload fails.
func writePinState(pin pinCfg) error {
	dir := filepath.Dir(pinStateFile)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("creating %s: %w", dir, err)
	}
	encoded, err := json.Marshal(pinStateRecord{PID: os.Getpid(), Gen: pin.gen, UpdatedAt: time.Now().UTC()})
	if err != nil {
		return fmt.Errorf("encoding pin state: %w", err)
	}

	f, err := os.OpenFile(pinStateFile, os.O_RDWR, 0o600)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("opening %s: %w", pinStateFile, err)
		}
		f, err = os.OpenFile(pinStateFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err != nil {
			return fmt.Errorf("creating %s: %w", pinStateFile, err)
		}
	}
	defer f.Close()
	if err := f.Truncate(0); err != nil {
		return fmt.Errorf("truncating %s: %w", pinStateFile, err)
	}
	if _, err := f.WriteAt(encoded, 0); err != nil {
		return fmt.Errorf("writing %s: %w", pinStateFile, err)
	}
	if err := f.Chmod(0o600); err != nil {
		return fmt.Errorf("chmod %s: %w", pinStateFile, err)
	}
	return f.Sync()
}

// ensurePinStateFilePlaceholder makes sure pinStateFile exists (even empty) before self guards
// attach, for the reason in selfguards.go's ensureHashFilePlaceholder: the live RO guard allows
// rewriting an existing entry but not creating one, and self guards now attach before
// writePinState's create fallback would run. No-op if /etc/app-listener is missing or the file
// exists.
func ensurePinStateFilePlaceholder() error {
	if _, err := os.Stat(pinStateFile); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	if _, statErr := os.Stat(filepath.Dir(pinStateFile)); statErr != nil {
		return nil //nolint:nilerr // no /etc/app-listener yet (nothing installed): nothing to bootstrap, not an error for the caller
	}
	f, err := os.OpenFile(pinStateFile, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	return f.Close()
}

// readPinState loads the last pin generation. Absent, empty or malformed all map to errNoPinState:
// an untrustworthy record is no better than none.
func readPinState() (pinStateRecord, error) {
	data, err := os.ReadFile(pinStateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return pinStateRecord{}, errNoPinState
		}
		return pinStateRecord{}, err
	}
	if len(data) == 0 {
		return pinStateRecord{}, errNoPinState
	}
	var rec pinStateRecord
	if jsonErr := json.Unmarshal(data, &rec); jsonErr != nil || rec.Gen == "" || rec.PID <= 0 {
		return pinStateRecord{}, errNoPinState
	}
	return rec, nil
}

// pinOwnerLikelyAlive reports whether rec.PID is a running process whose exe is this same binary,
// i.e. a live daemon still owns rec.Gen and nothing may touch (or load) its pins (that would race
// its access to the guarded resources). Dead, gone or PID-recycled reads as not alive. In the
// normal ExecStopPost flow this is always false (systemd runs it after the main process exits); the
// check covers manual, repeated or racing --lockdown invocations and relockStaleVaults' reuse of
// this path.
func pinOwnerLikelyAlive(rec pinStateRecord) bool {
	if rec.PID <= 0 {
		return false
	}
	target, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", rec.PID))
	if err != nil {
		return false // no such process, or /proc/<pid>/exe unreadable (exited, zombie)
	}
	self, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return false
	}
	return target == self
}

// recoverPinState reads the last pin generation for lockRootRecovering to widen a file-vault root's
// self-access. Returns a zero record (Gen == "": do not widen) when there's no usable record or its
// PID still looks like a live daemon. logPrefix tags logs ("lockdown" or "daemon"); base is the
// caller's already-resolved bpffs pin base (guard.ResolvePinBase), passed in to avoid re-probing
// bpffs.
func recoverPinState(logPrefix string) (rec pinStateRecord) {
	got, err := readPinState()
	switch {
	case err != nil:
		log.Warnf("%s: %v — a file-vault resource left unlocked by a prior crash cannot be self-access-widened this run", logPrefix, err)
		return pinStateRecord{}
	case pinOwnerLikelyAlive(got):
		log.Warnf("%s: pid %d still looks like a live app-listener daemon (generation %s) — not touching its pinned guards", logPrefix, got.PID, got.Gen)
		return pinStateRecord{}
	default:
		return got
	}
}
