package binledger

import (
	"context"

	ebpf "github.com/Virgula0/app-listener/internal/infrastructure"
)

// MaxHashCacheRows bounds hash_cache: a whitelist far past any real one.
const MaxHashCacheRows = 16384

// HashCache maps a file's DurableStat to its content hash, as recorded by an earlier run. A row
// only spares reading the file; Lookup still judges the hash.
type HashCache map[ebpf.DurableStat][32]byte

// LoadHashCache returns every well-formed hash_cache row (a malformed one is skipped: its file is
// hashed again).
func (l *Ledger) LoadHashCache() (HashCache, error) {
	rows, err := l.db.QueryContext(context.Background(),
		`SELECT dev, ino, fsid, size, mtime_ns, ctime_ns, btime_ns, sha256 FROM hash_cache LIMIT ?`, MaxHashCacheRows)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := HashCache{}
	for rows.Next() {
		var dev, ino, fsid int64
		var ds ebpf.DurableStat
		var sum []byte
		if err := rows.Scan(&dev, &ino, &fsid, &ds.Size, &ds.MtimeNs, &ds.CtimeNs, &ds.BtimeNs, &sum); err != nil {
			return nil, err
		}
		if len(sum) != 32 {
			continue
		}
		ds.Dev, ds.Ino, ds.Fsid = uint64(dev), uint64(ino), uint64(fsid) //nolint:gosec // stored bit-cast
		out[ds] = [32]byte(sum)
	}
	return out, rows.Err()
}

// SaveHashCache upserts rows in one transaction; with prune the table becomes exactly rows.
func (l *Ledger) SaveHashCache(rows HashCache, prune bool) error {
	tx, err := l.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() //nolint:errcheck // no-op after Commit
	if prune {
		if _, err = tx.ExecContext(context.Background(), `DELETE FROM hash_cache`); err != nil {
			return err
		}
	}
	stmt, err := tx.PrepareContext(context.Background(), `INSERT OR REPLACE INTO hash_cache
		(dev, ino, fsid, size, mtime_ns, ctime_ns, btime_ns, sha256) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	n := 0
	for ds, sum := range rows {
		if n >= MaxHashCacheRows {
			break
		}
		// sqlite integers are signed: unsigned ids are stored bit-cast.
		if _, err = stmt.ExecContext(context.Background(), int64(ds.Dev), int64(ds.Ino), int64(ds.Fsid), //nolint:gosec // bit-cast
			ds.Size, ds.MtimeNs, ds.CtimeNs, ds.BtimeNs, sum[:]); err != nil {
			return err
		}
		n++
	}
	return tx.Commit()
}
