package binledger

import (
	"context"
	"path/filepath"
	"testing"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

func TestHashCacheRoundTripAndPrune(t *testing.T) {
	l, err := Open(filepath.Join(t.TempDir(), "binaries.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	// High-bit ids: sqlite integers are signed.
	a := ebpf.DurableStat{BinaryStat: ebpf.BinaryStat{Dev: 1<<63 | 5, Ino: 1<<63 | 7, Size: 10, MtimeNs: 1, CtimeNs: 2},
		BtimeNs: 3, Fsid: 1<<63 | 9}
	b := a
	b.Ino = 8
	if err := l.SaveHashCache(HashCache{a: {1}, b: {2}}, false); err != nil {
		t.Fatal(err)
	}
	got, err := l.LoadHashCache()
	if err != nil || len(got) != 2 || got[a] != [32]byte{1} || got[b] != [32]byte{2} {
		t.Fatalf("round trip = %v, %v", got, err)
	}

	// Same inode, new fingerprint: replaces the row.
	a2 := a
	a2.CtimeNs = 99
	if err := l.SaveHashCache(HashCache{a2: {3}}, false); err != nil {
		t.Fatal(err)
	}
	if got, _ := l.LoadHashCache(); len(got) != 2 || got[a2] != [32]byte{3} {
		t.Fatalf("upsert = %v", got)
	}
	if err := l.SaveHashCache(HashCache{b: {2}}, true); err != nil {
		t.Fatal(err)
	}
	if got, _ := l.LoadHashCache(); len(got) != 1 || got[b] != [32]byte{2} {
		t.Fatalf("prune kept = %v", got)
	}
}

func TestHashCacheSkipsMalformedRow(t *testing.T) {
	l, err := Open(filepath.Join(t.TempDir(), "binaries.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	if _, err := l.db.ExecContext(context.Background(), `INSERT INTO hash_cache VALUES (1, 2, 3, 4, 5, 6, 7, x'0102')`); err != nil {
		t.Fatal(err)
	}
	if got, err := l.LoadHashCache(); err != nil || len(got) != 0 {
		t.Fatalf("a short hash must be skipped: %v, %v", got, err)
	}
}
