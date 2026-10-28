package bintrust

import (
	"crypto/sha256"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/binledger"
)

func openLedger(t *testing.T) *binledger.Ledger {
	t.Helper()
	l, err := binledger.Open(filepath.Join(t.TempDir(), "binaries.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	return l
}

func write(t *testing.T, path, content string) [32]byte {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
	return sha256.Sum256([]byte(content))
}

func TestInspectRefusesAFIFOWithoutBlocking(t *testing.T) {
	fifo := filepath.Join(t.TempDir(), "fifo")
	if err := unix.Mkfifo(fifo, 0o600); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { _, _, err := Inspect(fifo); done <- err }()
	select {
	case err := <-done:
		if err == nil || !strings.Contains(err.Error(), "not a regular file") {
			t.Fatalf("err = %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Inspect blocked on a FIFO")
	}
}

func TestPendingRowsSkipChangedContent(t *testing.T) {
	l := openLedger(t)
	dir := t.TempDir()
	same, changed := filepath.Join(dir, "same"), filepath.Join(dir, "changed")
	sum := write(t, same, "v1")
	old := write(t, changed, "v1")
	write(t, changed, "v2")
	for _, p := range []*binledger.Pending{
		{Line: same, Resolved: same, SHA256: sum, Resources: "/r", SeenAt: time.Unix(1, 0)},
		{Line: changed, Resolved: changed, SHA256: old, Resources: "/r", SeenAt: time.Unix(2, 0)},
	} {
		if err := l.AddPending(p); err != nil {
			t.Fatal(err)
		}
	}
	rows, err := Pending(l)
	if err != nil || len(rows) != 2 {
		t.Fatalf("rows = %+v, %v", rows, err)
	}
	if why := rows[0].Skip(); why != "" || rows[0].Stat == nil {
		t.Fatalf("unchanged row skipped: %q", why)
	}
	if why := rows[1].Skip(); !strings.Contains(why, "changed again") {
		t.Fatalf("changed row: skip = %q", why)
	}
}

func TestConfirmPendingOnlyWhatWasRefusedAndIsStillThere(t *testing.T) {
	l := openLedger(t)
	bin := filepath.Join(t.TempDir(), "tool")
	sum := write(t, bin, "v1")
	if err := ConfirmPending(l, bin, sum); err == nil {
		t.Fatal("confirmed a line the daemon never refused")
	}
	if err := l.AddPending(&binledger.Pending{Line: bin, Resolved: bin, SHA256: sum, Resources: "/r",
		SeenAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	write(t, bin, "v2")
	if err := ConfirmPending(l, bin, sum); !errors.Is(err, ErrChanged) {
		t.Fatalf("file changed since: err = %v", err)
	}
	write(t, bin, "v1")
	if err := ConfirmPending(l, bin, sum); err != nil {
		t.Fatal(err)
	}
	if got, ok, _ := l.Lookup(bin); !ok || got != sum {
		t.Fatalf("lookup = %x %v", got, ok)
	}
}

func TestConfirmGivenNeedsTheReviewedContent(t *testing.T) {
	l := openLedger(t)
	bin := filepath.Join(t.TempDir(), "tool")
	rows := Given([]string{bin}, "(by hand)")
	if rows[0].Skip() == "" {
		t.Fatal("a missing file is confirmable")
	}
	sum := write(t, bin, "v1")
	rows = Given([]string{bin}, "(by hand)")
	if rows[0].Skip() != "" || rows[0].SHA256 != sum {
		t.Fatalf("row = %+v", rows[0])
	}
	write(t, bin, "v2")
	if err := ConfirmGiven(l, bin, sum); !errors.Is(err, ErrChanged) {
		t.Fatalf("err = %v", err)
	}
	if _, ok, _ := l.Lookup(bin); ok {
		t.Fatal("recorded content that is no longer there")
	}
}

// trust_confirm runs before AUTH: a line the ledger doesn't refuse with that content is never
// opened, so a caller can't learn whether a guessed hash matches a file it names.
func TestConfirmPendingNeverOpensALineTheLedgerDoesNotRefuse(t *testing.T) {
	l := openLedger(t)
	fifo := filepath.Join(t.TempDir(), "fifo")
	if err := unix.Mkfifo(fifo, 0o600); err != nil {
		t.Fatal(err)
	}
	// A FIFO opened for reading would report "not a regular file": the error must be the ledger's.
	if err := ConfirmPending(l, fifo, [32]byte{1}); !errors.Is(err, errNotRefused) {
		t.Fatalf("err = %v", err)
	}
	if err := l.AddPending(&binledger.Pending{Line: fifo, Resolved: fifo, SHA256: [32]byte{2}, Resources: "/r",
		SeenAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	if err := ConfirmPending(l, fifo, [32]byte{1}); !errors.Is(err, errNotRefused) {
		t.Fatalf("refused with other content: err = %v", err)
	}
}
