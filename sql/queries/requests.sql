-- name: CreateRequest :one
INSERT INTO request_blobs (
  name, request_json
) VALUES (
    ?, ?
)
RETURNING *;

-- name: GetRequest :one
SELECT * FROM request_blobs WHERE name = ? LIMIT 1;

-- name: ListRequests :many
SELECT * FROM request_blobs ORDER BY created_at DESC;

-- name: UpdateRequest :one
UPDATE request_blobs
SET request_json = ?, updated_at = CURRENT_TIMESTAMP
WHERE name = ?
RETURNING *;

-- name: DeleteRequest :exec
DELETE FROM request_blobs WHERE name = ?;

