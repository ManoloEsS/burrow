-- +migrate Up
CREATE TABLE IF NOT EXISTS request_blobs (
  name TEXT PRIMARY KEY,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  request_json TEXT NOT_NULL
);

-- +migrate Down
DROP TABLE IF EXISTS request_blobs;
