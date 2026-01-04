CREATE TABLE goose_db_version (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		version_id INTEGER NOT NULL,
		is_applied INTEGER NOT NULL,
		tstamp TIMESTAMP DEFAULT (datetime('now'))
	);

CREATE TABLE sqlite_sequence(name,seq);

CREATE TABLE requests (
  id TEXT PRIMARY KEY,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  name TEXT NOT NULL,
  method TEXT NOT NULL CHECK (method IN ('GET', 'POST', 'PUT', 'DELETE',
'PATCH')),
  url TEXT NOT NULL,
  content_type TEXT,
  body TEXT,
  params TEXT,
  auth TEXT,
  headers TEXT
);
