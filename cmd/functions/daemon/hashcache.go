package daemon

import (
	"maps"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/Virgula0/app-listener/internal/binledger"
	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// hashCache is guard.HashStore over the ledger's hash_cache: read whole once at start, written
// back in one transaction per flush, so a start costs one query and one commit, not one per binary.
type hashCache struct {
	ledger *binledger.Ledger

	mu     sync.Mutex
	loaded binledger.HashCache
	// used: this run's rows (hits and new hashes), what a pruning flush keeps. dirty: not yet
	// written.
	used  binledger.HashCache
	dirty binledger.HashCache
}

// newHashCache loads l's rows; a failed read starts empty (every binary is hashed).
func newHashCache(l *binledger.Ledger) *hashCache {
	loaded, err := l.LoadHashCache()
	if err != nil {
		log.Warnf("daemon: reading the binary hash cache: %v — hashing every whitelisted binary", err)
		loaded = binledger.HashCache{}
	}
	return &hashCache{ledger: l, loaded: loaded, used: binledger.HashCache{}, dirty: binledger.HashCache{}}
}

func (c *hashCache) LookupHash(ds ebpf.DurableStat) ([32]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	sum, ok := c.loaded[ds]
	if ok {
		c.used[ds] = sum
	}
	return sum, ok
}

func (c *hashCache) StoreHash(ds ebpf.DurableStat, sum [32]byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.used) >= binledger.MaxHashCacheRows {
		return
	}
	c.loaded[ds] = sum
	c.used[ds] = sum
	c.dirty[ds] = sum
}

// flush writes new rows; prune (once startup admitted the whitelist) drops every row this run
// didn't use: replaced binaries and removed lines.
func (c *hashCache) flush(prune bool) {
	if c == nil {
		return
	}
	c.mu.Lock()
	rows := c.dirty
	if prune {
		rows = c.used
	}
	if len(rows) == 0 && !prune {
		c.mu.Unlock()
		return
	}
	snapshot := make(binledger.HashCache, len(rows))
	maps.Copy(snapshot, rows)
	c.dirty = binledger.HashCache{}
	c.mu.Unlock()

	if err := c.ledger.SaveHashCache(snapshot, prune); err != nil {
		log.Warnf("daemon: writing the binary hash cache: %v — the next start hashes those binaries again", err)
	}
}
