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

// pinStateFile records which pin generation (and PID) this daemon process's
// CONFIG guards (the [watch] sections, built by buildGuards — never the
// self-guards in selfguards.go, whose generation never rotates) are
// currently pinned under. A later, separate process with no live *guard.Guard
// object — `daemon --lockdown`, the systemd ExecStopPost safety net — reads
// it to recompute the exact bpffs pin paths (see guard.PinPrefix /
// guard.WithPinnedSelfVaultAccess) a crashed daemon left behind, and to tell
// whether the daemon that owns them might still be alive before ever
// touching them.
//
// Lives next to fscrypt.key and edit-auth.hash inside the daemon's
// self-guarded /etc/app-listener (ModeReadOnly — see selfProtectSpecs).
// pid+gen are not secrets, so this file's confidentiality does not matter;
// its INTEGRITY does, which the ReadOnly self-guard already provides once a
// daemon is running (only the app-listener binary may overwrite it).
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

// errNoPinState means there is no trustworthy pin-generation record to act
// on — a pre-upgrade daemon, a fresh install that never started, or a
// malformed/empty file. `daemon --lockdown` treats this exactly like "cannot
// recover the generation" and falls back to a plain, unwidened lock attempt.
var errNoPinState = errors.New("no usable pin-state record found")

// writePinState records pin's generation for the config guards this process
// just brought up (startup) or rotated to (a successful reload), IN PLACE on
// pinStateFile's existing inode once it exists — mirroring
// cmd/functions/editprotected/auth.go's WriteHashFile, and for the identical
// reason: once the daemon's own ModeReadOnly self-guard on /etc/app-listener
// is attached, only rewriting an EXISTING directory entry is permitted, not
// creating a new one. The very first call in a host's lifetime always runs
// before that guard exists (see runDaemon: this is called right after
// startGuardedDaemonAbortable succeeds, before `go sg.attach(pin)`), so the
// create-fallback below is exercised at most once per host. Best-effort: a
// failure here only means `--lockdown` cannot recover this generation later,
// never that startup/reload itself should fail.
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

// readPinState loads the last-recorded pin generation. Absent, empty, or
// malformed content all map to errNoPinState: a record this process cannot
// trust is no better than a missing one.
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

// pinOwnerLikelyAlive reports whether rec.PID names a still-running process
// whose own executable is this same binary — i.e. a live app-listener daemon
// still owns pin generation rec.Gen, and nothing should touch (or even load)
// any pin under it: doing so would race that daemon's own concurrent access
// to the resources it guards. A dead, gone, or PID-recycled-by-an-unrelated-
// process reads as "not alive". In the normal ExecStopPost flow this always
// reads false: systemd runs ExecStopPost only once the unit's main process
// has fully exited. The check exists for every OTHER way `--lockdown` (or
// relockStaleVaults's reuse of the same recovery path at the next startup)
// can be invoked — manually, twice, or racing a slow-to-exit process.
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

// recoverPinState reads the last-recorded pin generation for lockRootRecovering
// to widen a file-vault root's self-access with. Returns a zero
// pinStateRecord (Gen == "") — meaning "cannot recover, do not attempt any
// widening" — when there is no usable record, or when the record's PID still
// looks like a live app-listener daemon (touching its pins would race it).
// logPrefix tags the log lines with the caller's context ("lockdown" or
// "daemon"); base is the caller's already-resolved bpffs pin base
// (guard.ResolvePinBase), passed in rather than re-resolved here so a caller
// that already has it (prepareDaemonStart's pin.base) does not redundantly
// re-probe bpffs.
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
