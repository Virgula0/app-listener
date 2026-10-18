-- A file's durable stat -> its content hash, so an unchanged binary is not read again.
CREATE TABLE IF NOT EXISTS hash_cache (
	dev INTEGER NOT NULL,
	ino INTEGER NOT NULL,
	fsid INTEGER NOT NULL,
	size INTEGER NOT NULL,
	mtime_ns INTEGER NOT NULL,
	ctime_ns INTEGER NOT NULL,
	btime_ns INTEGER NOT NULL,
	sha256 BLOB NOT NULL,
	PRIMARY KEY (dev, ino)
);
