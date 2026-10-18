-- The edit-protected password: one PBKDF2 record (editprotected.Hash), never the password.
CREATE TABLE IF NOT EXISTS edit_auth (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	hash TEXT NOT NULL,
	updated_at INTEGER NOT NULL
);

-- Password guesses shared by every verifier (daemon control socket, CLI), surviving restarts.
CREATE TABLE IF NOT EXISTS edit_auth_throttle (
	id INTEGER PRIMARY KEY CHECK (id = 1),
	failures INTEGER NOT NULL,
	locked_until INTEGER NOT NULL
);
