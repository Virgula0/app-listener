package guard

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	cilium "github.com/cilium/ebpf"
)

func TestPathWithinMatchesPrefixRule(t *testing.T) {
	paths := []string{"", "/", "//", "/a", "/a/", "/a/b", "/ab", "/ab/c", "a", "a/b", "/a//b", "/home/u/.ssh", "/home/u/.sshx"}
	for _, root := range append(paths, "/home/u", "/home/u/.ssh") {
		for _, path := range paths {
			want := path == root || strings.HasPrefix(path, root+"/")
			if got := pathWithin(path, root); got != want {
				t.Fatalf("pathWithin(%q, %q) = %v, want %v", path, root, got, want)
			}
		}
	}
	if n := testing.AllocsPerRun(100, func() { pathWithin("/home/u/.ssh/id", "/home/u/.ssh") }); n != 0 {
		t.Fatalf("pathWithin allocates %v times", n)
	}
}

// fakeBatches serves rows in fixed-size calls; the last one ends with ErrKeyNotExist alongside its
// rows, as the kernel does.
func fakeBatches(batches [][]uint64, vals [][]uint32, tail error) inodeBatchFn {
	call := 0
	return func(_ *cilium.MapBatchCursor, keys []GuardInodeKey, out []uint32) (int, error) {
		if call >= len(batches) {
			return 0, tail
		}
		b, v := batches[call], vals[call]
		call++
		for i := range b {
			keys[i] = GuardInodeKey{Dev: 1, Ino: b[i]}
			out[i] = v[i]
		}
		if call == len(batches) {
			return len(b), tail
		}
		return len(b), nil
	}
}

func TestBatchInodeRows(t *testing.T) {
	key := func(inos ...uint64) []GuardInodeKey {
		out := make([]GuardInodeKey, 0, len(inos))
		for _, i := range inos {
			out = append(out, GuardInodeKey{Dev: 1, Ino: i})
		}
		return out
	}
	batches := [][]uint64{{1, 2, 3}, {4, 2, 5}}
	vals := [][]uint32{{7, 8, 7}, {7, 7, 7}}

	rows, err := batchInodeRows(fakeBatches(batches, vals, cilium.ErrKeyNotExist), 7)
	if err != nil || !reflect.DeepEqual(rows, key(1, 3, 4, 2, 5)) {
		t.Fatalf("rows %v err %v: want resID 7's keys once each, last batch included", rows, err)
	}

	boom := errors.New("ENOSPC")
	if _, err := batchInodeRows(fakeBatches(batches, vals, boom), 7); !errors.Is(err, boom) {
		t.Fatalf("err %v, want the batch error so the caller falls back", err)
	}
	if _, err := batchInodeRows(fakeBatches(batches, vals, nil), 7); err == nil {
		t.Fatal("a stalled walk (no rows, no end marker) must error, not loop or report done")
	}
}
