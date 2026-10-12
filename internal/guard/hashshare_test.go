package guard

import (
	"crypto/sha256"
	"os"
	"path/filepath"
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
	h, _, r := HashStats()
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
