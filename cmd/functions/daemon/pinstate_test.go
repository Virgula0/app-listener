package daemon

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"testing"
)

func withTempPinStateFile(t *testing.T) {
	t.Helper()
	old := pinStateFile
	pinStateFile = filepath.Join(t.TempDir(), "pin-state.json")
	t.Cleanup(func() { pinStateFile = old })
}

func pinStateInode(t *testing.T) uint64 {
	t.Helper()
	info, err := os.Lstat(pinStateFile)
	if err != nil {
		t.Fatalf("lstat %s: %v", pinStateFile, err)
	}
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		t.Fatal("Stat_t unavailable on this platform")
	}
	return st.Ino
}

// TestPinStateRoundTrip: a freshly written record reads back with the same
// PID and generation.
func TestPinStateRoundTrip(t *testing.T) {
	withTempPinStateFile(t)

	if _, err := readPinState(); !errors.Is(err, errNoPinState) {
		t.Fatalf("readPinState before any write = %v, want errNoPinState", err)
	}

	if err := writePinState(pinCfg{base: "/sys/fs/bpf/app-listener", gen: "abc123"}); err != nil {
		t.Fatal(err)
	}

	got, err := readPinState()
	if err != nil {
		t.Fatal(err)
	}
	if got.Gen != "abc123" {
		t.Errorf("Gen = %q, want %q", got.Gen, "abc123")
	}
	if got.PID != os.Getpid() {
		t.Errorf("PID = %d, want %d", got.PID, os.Getpid())
	}
}

// TestWritePinStateInPlace is the regression test mirroring
// TestWriteHashFileInPlaceOverPlaceholder: once a record already exists (as
// it always does after the first write in a process's lifetime), a second
// write — a reload minting a fresh generation — must rewrite it on the SAME
// inode, never by creating a new directory entry and swapping it in. Once
// the daemon's ModeReadOnly self-guard on /etc/app-listener is attached,
// only rewriting an EXISTING entry is permitted.
func TestWritePinStateInPlace(t *testing.T) {
	withTempPinStateFile(t)

	if err := writePinState(pinCfg{gen: "gen-one"}); err != nil {
		t.Fatal(err)
	}
	before := pinStateInode(t)

	if err := writePinState(pinCfg{gen: "gen-two"}); err != nil {
		t.Fatal(err)
	}
	if got := pinStateInode(t); got != before {
		t.Errorf("writePinState changed the inode: %d -> %d", before, got)
	}
	got, err := readPinState()
	if err != nil || got.Gen != "gen-two" {
		t.Fatalf("readPinState after second write = %+v, %v", got, err)
	}
}

// TestReadPinStateMalformed: an empty or corrupt file reads as "no state",
// same as absent — a record this process cannot trust is no better than a
// missing one, and callers must fall back to the unwidened lock attempt
// rather than propagate a parse error.
func TestReadPinStateMalformed(t *testing.T) {
	withTempPinStateFile(t)

	for _, content := range []string{"", "not json", `{"pid": 0, "gen": "x"}`, `{"pid": 123, "gen": ""}`} {
		if err := os.WriteFile(pinStateFile, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := readPinState(); !errors.Is(err, errNoPinState) {
			t.Errorf("readPinState(%q) = %v, want errNoPinState", content, err)
		}
	}
}

// TestPinOwnerLikelyAlive: this test process's own PID always reads as
// "alive" (its exe is trivially this same binary); a PID no process can
// plausibly hold (0, negative) never does.
func TestPinOwnerLikelyAlive(t *testing.T) {
	if !pinOwnerLikelyAlive(pinStateRecord{PID: os.Getpid()}) {
		t.Error("this process's own PID should read as alive")
	}
	if pinOwnerLikelyAlive(pinStateRecord{PID: 0}) {
		t.Error("PID 0 must never read as alive")
	}
	if pinOwnerLikelyAlive(pinStateRecord{PID: -1}) {
		t.Error("a negative PID must never read as alive")
	}
}

// TestPinOwnerLikelyAliveDeadPID: a PID that (almost certainly) names no
// running process at all reads as not-alive. Uses a very large PID unlikely
// to be assigned on a test host rather than guessing a specific dead PID.
func TestPinOwnerLikelyAliveDeadPID(t *testing.T) {
	const implausiblePID = 1 << 30
	if _, err := os.Stat(filepath.Join("/proc", strconv.Itoa(implausiblePID))); err == nil {
		t.Skip("implausible PID is somehow live on this host")
	}
	if pinOwnerLikelyAlive(pinStateRecord{PID: implausiblePID}) {
		t.Error("a PID naming no running process must not read as alive")
	}
}

// TestRecoverPinStateSkipsLiveOwner: recoverPinState must return the zero
// record (Gen == "", the "do not widen" signal) when the recorded PID still
// looks alive — touching a live daemon's pins would race its own concurrent
// access to the resource it guards.
func TestRecoverPinStateSkipsLiveOwner(t *testing.T) {
	withTempPinStateFile(t)

	if err := writePinState(pinCfg{gen: "still-live"}); err != nil {
		t.Fatal(err) // PID recorded is this test process's own — reads as alive
	}
	if rec := recoverPinState("test"); rec.Gen != "" {
		t.Errorf("recoverPinState with a live owner = %+v, want zero record", rec)
	}
}

// TestRecoverPinStateNoRecord: with nothing ever written, recoverPinState
// returns the zero record rather than erroring.
func TestRecoverPinStateNoRecord(t *testing.T) {
	withTempPinStateFile(t)

	if rec := recoverPinState("test"); rec.Gen != "" {
		t.Errorf("recoverPinState with no record = %+v, want zero record", rec)
	}
}
