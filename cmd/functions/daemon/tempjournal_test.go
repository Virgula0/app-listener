package daemon

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/Virgula0/app-listener/internal/guard"
)

func TestTempJournalRoundTripInPlace(t *testing.T) {
	old := tempJournalFile
	tempJournalFile = filepath.Join(t.TempDir(), "temp-grants.json")
	t.Cleanup(func() { tempJournalFile = old })

	if err := ensurePlaceholder(tempJournalFile); err != nil {
		t.Fatal(err)
	}
	before, err := os.Stat(tempJournalFile)
	if err != nil {
		t.Fatal(err)
	}
	if rows, err := readTempJournal(); err != nil || rows != nil {
		t.Fatalf("empty placeholder read as %v, %v", rows, err)
	}

	want := []guard.GuardResInodeKey{
		{ResId: 3, Ino: guard.GuardInodeKey{Dev: 1 << 20, Ino: 42}},
		{ResId: 7, Ino: guard.GuardInodeKey{Dev: 2, Ino: 9}},
	}
	if err := writeTempJournal(want); err != nil {
		t.Fatal(err)
	}
	got, err := readTempJournal()
	if err != nil || !slices.Equal(got, want) {
		t.Fatalf("read back %v, %v; want %v", got, err, want)
	}

	if err := writeTempJournal(nil); err != nil {
		t.Fatal(err)
	}
	if rows, err := readTempJournal(); err != nil || rows != nil {
		t.Fatalf("cleared journal read as %v, %v", rows, err)
	}
	after, err := os.Stat(tempJournalFile)
	if err != nil {
		t.Fatal(err)
	}
	// The self-guarded /etc/app-listener only lets the daemon rewrite an existing entry.
	if !os.SameFile(before, after) {
		t.Fatal("the journal was replaced instead of rewritten in place")
	}
}
