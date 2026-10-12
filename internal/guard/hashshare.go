package guard

import (
	"os"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sys/unix"

	"github.com/Virgula0/app-listener/internal/infrastructure"
)

// sharedHashes dedupes content hashing across every guard, reload and verifier tick of the process
// (one app blob is admitted by ~10 guards). Keyed by inode; a hash is reused only under the full
// fingerprint it was taken under, the rule the verifier already skips unchanged binaries by.
var (
	sharedHashMu sync.Mutex
	sharedHashes = map[GuardInodeKey]sharedHashEntry{}
	hashStore    HashStore // under sharedHashMu
	hashStats    struct{ hashed, reused, cached, bytes atomic.Uint64 }
)

// HashStore persists content hashes across daemon runs (the binary ledger's hash cache), keyed by
// a DurableStat. A hit only replaces reading the file: the ledger still judges the hash. Lookup
// failures are misses.
type HashStore interface {
	LookupHash(ebpf.DurableStat) ([32]byte, bool)
	StoreHash(ebpf.DurableStat, [32]byte)
}

// SetHashStore installs the persistent store; nil keeps sharing in-process only.
func SetHashStore(s HashStore) {
	sharedHashMu.Lock()
	hashStore = s
	sharedHashMu.Unlock()
}

// statDurable and syncFile are the ebpf calls; tests swap them (a durable mount needs root).
var (
	statDurable = ebpf.StatDurable
	syncFile    = ebpf.SyncFile
)

type sharedHashEntry struct {
	stat ebpf.BinaryStat
	hash [32]byte
}

// maxSharedHashes bounds the map; only a pathological whitelist reaches it.
const maxSharedHashes = 8192

// settleAge: a fingerprint vouches for content only once its ctime is older than any timestamp
// granularity. A same-size rewrite within one coarse tick (pre-multigrain kernels) keeps
// size/mtime/ctime, so a younger file's hash is never shared or trusted as unchanged.
var settleAge = time.Second

func settled(fp ebpf.BinaryStat) bool {
	return time.Now().UnixNano()-fp.CtimeNs >= settleAge.Nanoseconds()
}

// computeEntryFile is ebpf.ComputeBinaryEntryFile; tests swap it to write during the read.
var computeEntryFile = ebpf.ComputeBinaryEntryFile

// hashFileShared returns the entry (named path) of the inode f holds, any fd including O_PATH, and
// fp, f's fingerprint before the read. A hash taken under fp is reused. A fresh one is shared only
// when f's fingerprint is unchanged after the read and settled: content written during it is never
// published under either fingerprint (the caller still gets it, as an unshared hash).
func hashFileShared(f *os.File, path string) (entry BinaryEntry, fp ebpf.BinaryStat, err error) {
	if fp, err = ebpf.StatBinaryFile(f); err != nil {
		return entry, fp, err
	}
	key := GuardInodeKey{Dev: fp.Dev, Ino: fp.Ino}
	sharedHashMu.Lock()
	e, ok := sharedHashes[key]
	store := hashStore
	sharedHashMu.Unlock()
	if ok && e.stat == fp {
		hashStats.reused.Add(1)
		return ebpf.EntryWithHash(path, e.hash), fp, nil
	}

	ds, durable := ebpf.DurableStat{}, false
	if store != nil {
		ds, durable = statDurable(f, fp)
	}
	if durable {
		if h, hit := store.LookupHash(ds); hit {
			hashStats.cached.Add(1)
			publishHash(key, fp, h)
			return ebpf.EntryWithHash(path, h), fp, nil
		}
		// Persisted hashes must outlive this run: clean every page first (SyncFile).
		durable = syncFile(f) == nil
	}

	if entry, err = computeEntryFile(f, path); err != nil {
		return entry, fp, err
	}
	hashStats.hashed.Add(1)
	hashStats.bytes.Add(uint64(max(fp.Size, 0))) //nolint:gosec // clamped non-negative
	if !unchangedSince(f, fp) || !settled(fp) {
		return entry, fp, nil
	}
	publishHash(key, fp, entry.Hash)
	if durable {
		store.StoreHash(ds, entry.Hash)
	}
	return entry, fp, nil
}

func publishHash(key GuardInodeKey, fp ebpf.BinaryStat, hash [32]byte) {
	sharedHashMu.Lock()
	if len(sharedHashes) >= maxSharedHashes {
		sharedHashes = map[GuardInodeKey]sharedHashEntry{}
	}
	sharedHashes[key] = sharedHashEntry{stat: fp, hash: hash}
	sharedHashMu.Unlock()
}

func unchangedSince(f *os.File, fp ebpf.BinaryStat) bool {
	after, err := ebpf.StatBinaryFile(f)
	return err == nil && after == fp
}

// hashPathShared is hashFileShared for the file path names now (symlinks followed).
func hashPathShared(path string) (BinaryEntry, ebpf.BinaryStat, error) {
	fd, err := unix.Open(path, unix.O_PATH|unix.O_CLOEXEC, 0)
	if err != nil {
		return BinaryEntry{}, ebpf.BinaryStat{}, &os.PathError{Op: "open", Path: path, Err: err}
	}
	f := os.NewFile(uintptr(fd), path)
	defer f.Close()
	return hashFileShared(f, path)
}

// newVerifyState pins key and its admitted hash for the verifier. When this run hashed that inode
// to that hash, the state starts hashed under the fingerprint it was taken with, so the first tick
// costs a stat, not a re-read.
func newVerifyState(key GuardInodeKey, hash [32]byte) *binaryVerifyState {
	st := &binaryVerifyState{key: key, hash: hash}
	sharedHashMu.Lock()
	e, ok := sharedHashes[key.Real()]
	sharedHashMu.Unlock()
	if ok && e.hash == hash {
		st.stat, st.hashed = e.stat, true
	}
	return st
}

// HashStats reports this process's binary content hashing so far: files read, bytes read, hashes
// reused within the run, and hashes taken from the persistent store.
func HashStats() (hashed, bytes, reused, cached uint64) {
	return hashStats.hashed.Load(), hashStats.bytes.Load(), hashStats.reused.Load(), hashStats.cached.Load()
}
