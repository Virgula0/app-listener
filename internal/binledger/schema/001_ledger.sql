-- The binary ledger: confirmed whitelist hashes, install grants and refused binaries.
CREATE TABLE IF NOT EXISTS meta (key TEXT PRIMARY KEY, value TEXT NOT NULL);

CREATE TABLE IF NOT EXISTS binaries (
	line TEXT PRIMARY KEY,
	sha256 BLOB NOT NULL,
	source TEXT NOT NULL,
	updated_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS grants (line TEXT PRIMARY KEY);

CREATE TABLE IF NOT EXISTS pending (
	line TEXT PRIMARY KEY,
	resolved TEXT NOT NULL,
	sha256 BLOB NOT NULL,
	resources TEXT NOT NULL,
	seen_at INTEGER NOT NULL
);
