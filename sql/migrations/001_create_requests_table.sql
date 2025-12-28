-- +goose Up
CREATE TABLE IF NOT EXISTS requests (
  id TEXT PRIMARY KEY,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  name TEXT NOT NULL,
  method TEXT NOT NULL CHECK (method IN ('GET', 'POST', 'PUT', 'DELETE',
'PATCH')),
  url TEXT NOT NULL,
  content_type TEXT,
  body TEXT,
  params TEXT,     -- JSON stored as TEXT for SQLite compatibility
  auth TEXT,       -- JSON stored as TEXT for SQLite compatibility
  headers TEXT     -- JSON stored as TEXT for SQLite compatibility
);

-- +goose Down
DROP TABLE IF EXISTS requests;
