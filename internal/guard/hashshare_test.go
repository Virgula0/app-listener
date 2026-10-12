package guard

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// resetSharedHashes isolates a test from the process-wide cache; settle 0 lets fresh temp files
// count as settled.
func resetSharedHashes(t *testing.T, settle time.Duration) {
	t.Helper()
	prevSettle, prevCompute := settleAge, computeEntryFile
	sharedHashes = map[GuardInodeKey]sharedHashEntry{}
	settleAge = settle
	t.Cleanup(func() {
		sharedHashes = map[GuardInodeKey]sharedHashEntry{}
		settleAge, computeEntryFile = prevSettle, prevCompute
	})
}

func writeBin(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "app")
	if err := os.WriteFile(p, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
	return p
}

func hashedReused(t *testing.T) (uint64, uint64) {
	t.Helper()
	h, _, r, _ := HashStats()
	return h, r
}

func TestConfinedEntryHashesOncePerFingerprint(t *testing.T) {
	resetSharedHashes(t, 0)
	p := writeBin(t, "content-a")
	h0, r0 := hashedReused(t)

	e1, v1, err := ConfinedEntry(p)
	if err != nil {
		t.Fatal(err)
	}
	defer v1.Close()
	e2, v2, err := ConfinedEntry(p)
	if err != nil {
		t.Fatal(err)
	}
	defer v2.Close()

	if e1 != e2 || v1.Key != v2.Key {
		t.Fatalf("second open diverged: %+v %v vs %+v %v", e1, v1.Key, e2, v2.Key)
	}
	if e1.Hash != sha256.Sum256([]byte("content-a")) {
		t.Fatal("wrong hash")
	}
	if h, r := hashedReused(t); h-h0 != 1 || r-r0 != 1 {
		t.Fatalf("hashed %d reused %d, want 1 and 1", h-h0, r-r0)
	}
}

func TestSharedHashMissesOnChangedFingerprint(t *testing.T) {
	resetSharedHashes(t, 0)
	p := writeBin(t, "content-a")
	e1, v1, err := ConfinedEntry(p)
	if err != nil {
		t.Fatal(err)
	}
	v1.Close()

	// Same inode, new content and size: rewritten in place.
	if err := os.WriteFile(p, []byte("content-bb"), 0o755); err != nil {
		t.Fatal(err)
	}
	e2, v2, err := ConfinedEntry(p)
	if err != nil {
		t.Fatal(err)
	}
	defer v2.Close()
	if v2.Key != v1.Key {
		t.Fatal("expected the same inode")
	}
	if e2.Hash == e1.Hash || e2.Hash != sha256.Sum256([]byte("content-bb")) {
		t.Fatal("a stale hash was reused after an in-place rewrite")
	}
}

func TestSharedHashNotPublishedWhenWrittenDuringRead(t *testing.T) {
	resetSharedHashes(t, 0)
	p := writeBin(t, "content-a")
	computeEntryFile = func(f *os.File, path string) (BinaryEntry, error) {
		e, err := ebpf.ComputeBinaryEntryFile(f, path)
		if werr := os.WriteFile(p, []byte("content-written-during"), 0o755); werr != nil {
			t.Fatal(werr)
		}
		return e, err
	}
	_, v, err := ConfinedEntry(p)
	if err != nil {
		t.Fatal(err)
	}
	v.Close()
	if len(sharedHashes) != 0 {
		t.Fatal("a hash of a file written during the read was shared")
	}
}

func TestSharedHashNotPublishedBeforeSettled(t *testing.T) {
	resetSharedHashes(t, time.Hour)
	p := writeBin(t, "content-a")
	_, v, err := ConfinedEntry(p)
	if err != nil {
		t.Fatal(err)
	}
	v.Close()
	if len(sharedHashes) != 0 {
		t.Fatal("a just-written file's hash was shared")
	}
}

func TestNewVerifyStateSeededOnlyByMatchingHash(t *testing.T) {
	resetSharedHashes(t, 0)
	p := writeBin(t, "content-a")
	e, v, err := ConfinedEntry(p)
	if err != nil {
		t.Fatal(err)
	}
	v.Close()

	st := newVerifyState(v.Key, e.Hash)
	if !st.hashed || st.stat.Ino != v.Key.Ino {
		t.Fatalf("not seeded: %+v", st)
	}
	if other := newVerifyState(v.Key, [32]byte{1}); other.hashed {
		t.Fatal("seeded although the admitted hash differs from the one taken")
	}
	if other := newVerifyState(GuardInodeKey{Dev: v.Key.Dev, Ino: v.Key.Ino + 1}, e.Hash); other.hashed {
		t.Fatal("seeded for another inode")
	}
}

func TestVerifierSkipsSeededBinaryOnFirstTick(t *testing.T) {
	resetSharedHashes(t, 0)
	p := writeBin(t, "content-a")
	e, v, err := ConfinedEntry(p)
	if err != nil {
		t.Fatal(err)
	}
	v.Close()
	g := &Guard{binaryVerifyStates: map[string]*binaryVerifyState{p: newVerifyState(v.Key, e.Hash)}}

	h0, _ := hashedReused(t)
	g.verifyBinaryHashesOnce()
	if h, _ := hashedReused(t); h != h0 {
		t.Fatal("the first tick re-read a binary hashed at admission")
	}
}

// The verifier pins the inode the hashed fd held, not one stat'ed by path before the open.
func TestVerifierHashesAndPinsFromOneFd(t *testing.T) {
	resetSharedHashes(t, 0)
	p := writeBin(t, "content-a")
	e, v, err := ConfinedEntry(p)
	if err != nil {
		t.Fatal(err)
	}
	v.Close()
	st := &binaryVerifyState{key: v.Key, hash: e.Hash}
	g := &Guard{binaryVerifyStates: map[string]*binaryVerifyState{p: st}}
	g.verifyBinaryHashesOnce()
	if !st.hashed || st.stat.Ino != v.Key.Ino || st.demoted {
		t.Fatalf("unchanged binary not pinned: %+v", st)
	}
}

type fakeHashStore struct {
	rows    map[ebpf.DurableStat][32]byte
	lookups int
}

func (s *fakeHashStore) LookupHash(ds ebpf.DurableStat) ([32]byte, bool) {
	s.lookups++
	h, ok := s.rows[ds]
	return h, ok
}

func (s *fakeHashStore) StoreHash(ds ebpf.DurableStat, h [32]byte) { s.rows[ds] = h }

// withStore installs a fake persistent store whose files all count as durable (a real durable
// mount needs root), recording syncFile calls.
func withStore(t *testing.T, durable bool, syncErr error) (*fakeHashStore, *int) {
	t.Helper()
	store := &fakeHashStore{rows: map[ebpf.DurableStat][32]byte{}}
	syncs := 0
	prevStat, prevSync := statDurable, syncFile
	statDurable = func(_ *os.File, fp ebpf.BinaryStat) (ebpf.DurableStat, bool) {
		return ebpf.DurableStat{BinaryStat: fp, BtimeNs: 7, Fsid: 9}, durable
	}
	syncFile = func(*os.File) error { syncs++; return syncErr }
	SetHashStore(store)
	t.Cleanup(func() {
		SetHashStore(nil)
		statDurable, syncFile = prevStat, prevSync
	})
	return store, &syncs
}

func TestPersistentStoreFilledAfterSyncAndHit(t *testing.T) {
	resetSharedHashes(t, 0)
	store, syncs := withStore(t, true, nil)
	p := writeBin(t, "content-a")

	e1, v1, err := ConfinedEntry(p)
	if err != nil {
		t.Fatal(err)
	}
	v1.Close()
	if *syncs != 1 || len(store.rows) != 1 {
		t.Fatalf("a fresh durable hash must be synced then stored: syncs %d rows %d", *syncs, len(store.rows))
	}

	// A new run: nothing in memory, the store answers without a read.
	sharedHashes = map[GuardInodeKey]sharedHashEntry{}
	h0, _ := hashedReused(t)
	_, _, _, c0 := HashStats()
	e2, v2, err := ConfinedEntry(p)
	if err != nil {
		t.Fatal(err)
	}
	v2.Close()
	_, _, _, c1 := HashStats()
	if h, _ := hashedReused(t); h != h0 || c1-c0 != 1 || e2.Hash != e1.Hash {
		t.Fatalf("store hit re-read the file (hashed +%d, cached +%d)", h-h0, c1-c0)
	}
	if *syncs != 1 {
		t.Fatal("a hit must not sync")
	}
	if st := newVerifyState(v2.Key, e2.Hash); !st.hashed {
		t.Fatal("a store hit must seed the verifier like a fresh hash")
	}
}

func TestPersistentStoreMissesOnRewrite(t *testing.T) {
	resetSharedHashes(t, 0)
	_, _ = withStore(t, true, nil)
	p := writeBin(t, "content-a")
	if _, v, err := ConfinedEntry(p); err == nil {
		v.Close()
	}
	if err := os.WriteFile(p, []byte("content-bb"), 0o755); err != nil {
		t.Fatal(err)
	}
	sharedHashes = map[GuardInodeKey]sharedHashEntry{}
	e, v, err := ConfinedEntry(p)
	if err != nil {
		t.Fatal(err)
	}
	v.Close()
	if e.Hash != sha256.Sum256([]byte("content-bb")) {
		t.Fatal("a stale persisted hash was used after a rewrite")
	}
}

func TestPersistentStoreSkipped(t *testing.T) {
	for name, tc := range map[string]struct {
		durable bool
		syncErr error
		settle  time.Duration
		racer   bool
		lookups int
	}{
		"not durable":         {durable: false, lookups: 0},
		"sync failed":         {durable: true, syncErr: os.ErrInvalid, lookups: 1},
		"not settled":         {durable: true, settle: time.Hour, lookups: 1},
		"written during read": {durable: true, racer: true, lookups: 1},
	} {
		t.Run(name, func(t *testing.T) {
			resetSharedHashes(t, tc.settle)
			store, _ := withStore(t, tc.durable, tc.syncErr)
			p := writeBin(t, "content-a")
			if tc.racer {
				computeEntryFile = func(f *os.File, path string) (BinaryEntry, error) {
					e, err := ebpf.ComputeBinaryEntryFile(f, path)
					if werr := os.WriteFile(p, []byte("content-written-during"), 0o755); werr != nil {
						t.Fatal(werr)
					}
					return e, err
				}
			}
			_, v, err := ConfinedEntry(p)
			if err != nil {
				t.Fatal(err)
			}
			v.Close()
			if len(store.rows) != 0 || store.lookups != tc.lookups {
				t.Fatalf("rows %d lookups %d (want 0, %d)", len(store.rows), store.lookups, tc.lookups)
			}
		})
	}
}

// Concurrent callers of one inode read it once: the others wait for the published hash. The
// leader's read is held until every caller has started, so the count doesn't depend on timing.
func TestConcurrentCallersHashOnce(t *testing.T) {
	resetSharedHashes(t, 0)
	p := writeBin(t, "content-a")
	started, release := make(chan struct{}), make(chan struct{})
	var once sync.Once
	computeEntryFile = func(f *os.File, path string) (BinaryEntry, error) {
		once.Do(func() { close(started) })
		<-release
		return ebpf.ComputeBinaryEntryFile(f, path)
	}
	h0, r0 := hashedReused(t)

	const callers = 8
	var wg sync.WaitGroup
	hashes := make([][32]byte, callers)
	run := func(i int) {
		defer wg.Done()
		e, v, err := ConfinedEntry(p)
		if err != nil {
			t.Error(err)
			return
		}
		v.Close()
		hashes[i] = e.Hash
	}
	wg.Add(1)
	go run(0)
	<-started
	for i := 1; i < callers; i++ {
		wg.Add(1)
		go run(i)
	}
	time.Sleep(20 * time.Millisecond)
	close(release)
	wg.Wait()

	if h, r := hashedReused(t); h-h0 != 1 || r-r0 != callers-1 {
		t.Fatalf("hashed %d reused %d, want 1 and %d", h-h0, r-r0, callers-1)
	}
	for _, h := range hashes {
		if h != sha256.Sum256([]byte("content-a")) {
			t.Fatal("a waiter got a wrong hash")
		}
	}
	if len(hashing) != 0 {
		t.Fatal("an in-flight marker was left behind")
	}
}

// A leader whose hash isn't published (here: failed) lets a waiter hash for itself.
func TestWaiterHashesWhenLeaderFails(t *testing.T) {
	resetSharedHashes(t, 0)
	p := writeBin(t, "content-a")
	started, release := make(chan struct{}), make(chan struct{})
	var calls sync.Mutex
	n := 0
	computeEntryFile = func(f *os.File, path string) (BinaryEntry, error) {
		calls.Lock()
		n++
		first := n == 1
		calls.Unlock()
		if first {
			close(started)
			<-release
			return BinaryEntry{}, os.ErrInvalid
		}
		return ebpf.ComputeBinaryEntryFile(f, path)
	}
	errc := make(chan error, 1)
	go func() {
		_, v, err := ConfinedEntry(p)
		if err == nil {
			v.Close()
		}
		errc <- err
	}()
	<-started
	type result struct {
		e   BinaryEntry
		err error
	}
	waiter := make(chan result, 1)
	go func() {
		e, v, err := ConfinedEntry(p)
		if err == nil {
			v.Close()
		}
		waiter <- result{e, err}
	}()
	time.Sleep(20 * time.Millisecond)
	close(release)
	if err := <-errc; err == nil {
		t.Fatal("the leader's failure was lost")
	}
	r := <-waiter
	if r.err != nil || r.e.Hash != sha256.Sum256([]byte("content-a")) {
		t.Fatalf("waiter: %v %x", r.err, r.e.Hash)
	}
}
