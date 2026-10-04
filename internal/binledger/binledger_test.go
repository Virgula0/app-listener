package binledger

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/Virgula0/app-listener/internal/daemonconfig"
)

func TestLedgerRecordLookupPending(t *testing.T) {
	path := filepath.Join(t.TempDir(), "binaries.db")
	l, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	if ok, err := l.Initialized(); err != nil || ok {
		t.Fatalf("fresh ledger initialized=%v err=%v", ok, err)
	}
	if _, ok, err := l.Lookup("/home/u/.local/bin/gh"); err != nil || ok {
		t.Fatalf("lookup on empty ledger ok=%v err=%v", ok, err)
	}

	sum := [32]byte{1, 2, 3}
	if err := l.AddPending(&Pending{Line: "/home/u/.local/bin/gh", Resolved: "/home/u/.local/bin/gh",
		SHA256: sum, Resources: "/home/u/.config/gh", SeenAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	if p, err := l.ListPending(); err != nil || len(p) != 1 || p[0].SHA256 != sum {
		t.Fatalf("pending = %+v, %v", p, err)
	}
	if err := l.Record("/home/u/.local/bin/gh", sum, SourceConfirmed); err != nil {
		t.Fatal(err)
	}
	if got, ok, err := l.Lookup("/home/u/.local/bin/gh"); err != nil || !ok || got != sum {
		t.Fatalf("lookup = %x %v %v", got, ok, err)
	}
	if p, _ := l.ListPending(); len(p) != 0 {
		t.Fatalf("Record left a pending row: %+v", p)
	}
	if err := l.MarkInitialized(); err != nil {
		t.Fatal(err)
	}
	if ok, err := l.Initialized(); err != nil || !ok {
		t.Fatalf("initialized=%v err=%v", ok, err)
	}
}

// The self guard lets the daemon rewrite existing files in /etc/app-listener, never create one:
// every write must leave the journal in place.
func TestLedgerKeepsJournalFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "binaries.db")
	l, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	if err := l.Record("/x", [32]byte{9}, SourceInstall); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(JournalPath(path)); err != nil {
		t.Fatalf("journal removed after a transaction: %v", err)
	}
	for _, p := range []string{path, JournalPath(path)} {
		if st, err := os.Stat(p); err != nil || st.Mode().Perm() != 0o600 {
			t.Fatalf("%s: mode %v err %v, want 0600", p, st.Mode().Perm(), err)
		}
	}
}

func TestEnsurePlaceholdersNoInstallDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "binaries.db")
	if err := EnsurePlaceholders(path); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("placeholder created without an install dir: %v", err)
	}
}

// A symlink line is keyed on the link: re-pointing it at another binary needs that one's hash.
func TestLinesKeyOnLinkAndIncludePending(t *testing.T) {
	cfg := &daemonconfig.Config{Resources: []daemonconfig.Resource{{
		Binaries:        []daemonconfig.BinaryRule{{Path: "/usr/bin/ssh"}, {Path: "/opt/app/app-1.2", Link: "/opt/app/app"}},
		PendingBinaries: []daemonconfig.BinaryRule{{Path: "/home/u/vault/bin/tool"}},
	}}}
	want := []string{"/usr/bin/ssh", "/opt/app/app", "/home/u/vault/bin/tool"}
	if got := Lines(cfg); !slices.Equal(got, want) {
		t.Fatalf("Lines = %v, want %v", got, want)
	}
}
