package binledger

import (
	"errors"
	"path/filepath"
	"testing"
	"time"
)

func openTemp(t *testing.T) *Ledger {
	t.Helper()
	l, err := Open(filepath.Join(t.TempDir(), "binaries.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { l.Close() })
	return l
}

func TestEditAuthHashRoundTrip(t *testing.T) {
	l := openTemp(t)
	if h, err := l.EditAuthHash(); err != nil || h != "" {
		t.Fatalf("fresh ledger hash=%q err=%v", h, err)
	}
	for _, h := range []string{"first", "second"} {
		if err := l.SetEditAuthHash(h); err != nil {
			t.Fatal(err)
		}
		if got, err := l.EditAuthHash(); err != nil || got != h {
			t.Fatalf("hash=%q err=%v, want %q", got, err, h)
		}
	}
	if err := l.SetEditAuthHash(""); err == nil {
		t.Fatal("an empty record was stored: that would clear the password")
	}
}

func TestBeginAuthAttemptLocksAfterBudget(t *testing.T) {
	l := openTemp(t)
	now := time.Now()
	for i := range MaxAuthFailures {
		if err := l.BeginAuthAttempt(now); err != nil {
			t.Fatalf("attempt %d refused: %v", i+1, err)
		}
	}
	var locked *LockedError
	if err := l.BeginAuthAttempt(now); !errors.As(err, &locked) {
		t.Fatalf("attempt past the budget: %v, want a lockout", err)
	}
	if err := l.BeginAuthAttempt(now.Add(AuthLockout + time.Second)); err != nil {
		t.Fatalf("attempt after the lockout: %v", err)
	}
}

func TestAuthSucceededResetsBudget(t *testing.T) {
	l := openTemp(t)
	now := time.Now()
	for range MaxAuthFailures - 1 {
		if err := l.BeginAuthAttempt(now); err != nil {
			t.Fatal(err)
		}
	}
	if err := l.AuthSucceeded(); err != nil {
		t.Fatal(err)
	}
	for i := range MaxAuthFailures {
		if err := l.BeginAuthAttempt(now); err != nil {
			t.Fatalf("attempt %d after a success refused: %v", i+1, err)
		}
	}
}

// Every Open re-applies the schema: an existing ledger must keep its rows.
func TestReopenKeepsRows(t *testing.T) {
	path := filepath.Join(t.TempDir(), "binaries.db")
	l, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := l.SetEditAuthHash("kept"); err != nil {
		t.Fatal(err)
	}
	l.Close()
	l, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	if h, err := l.EditAuthHash(); err != nil || h != "kept" {
		t.Fatalf("hash after reopen=%q err=%v", h, err)
	}
}
