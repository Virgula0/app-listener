package daemon

import (
	"path/filepath"
	"testing"

	"github.com/Virgula0/app-listener/internal/binledger"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

func ds(ino uint64) ebpf.DurableStat {
	return ebpf.DurableStat{BinaryStat: ebpf.BinaryStat{Dev: 1, Ino: ino, Size: 1}, BtimeNs: 1, Fsid: 1}
}

func TestHashCacheFlushKeepsOnlyThisRunsRows(t *testing.T) {
	l, err := binledger.Open(filepath.Join(t.TempDir(), "binaries.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	if err := l.SaveHashCache(binledger.HashCache{ds(1): {1}, ds(2): {2}}, false); err != nil {
		t.Fatal(err)
	}

	c := newHashCache(l)
	if h, ok := c.LookupHash(ds(1)); !ok || h != [32]byte{1} {
		t.Fatal("loaded row missed")
	}
	if _, ok := c.LookupHash(ds(3)); ok {
		t.Fatal("unknown row hit")
	}
	c.StoreHash(ds(3), [32]byte{3})
	c.flush(true)
	got, err := l.LoadHashCache()
	if err != nil || len(got) != 2 || got[ds(1)] != [32]byte{1} || got[ds(3)] != [32]byte{3} {
		t.Fatalf("after the startup flush = %v, %v: want the hit and the new row, not the unused one", got, err)
	}

	c.StoreHash(ds(4), [32]byte{4})
	c.flush(false)
	if got, _ := l.LoadHashCache(); len(got) != 3 || got[ds(4)] != [32]byte{4} {
		t.Fatalf("a later flush must add, not prune: %v", got)
	}
	var nilCache *hashCache
	nilCache.flush(true)
}
