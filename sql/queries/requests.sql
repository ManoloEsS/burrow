-- name: CreateRequest :one
INSERT INTO requests (
  id, name, method, url, content_type, body, params, auth, headers
) VALUES (
    ?, ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: GetRequest :one
SELECT * FROM requests WHERE id = ? LIMIT 1;

-- name: ListRequests :many
SELECT * FROM requests ORDER BY created_at DESC;

-- name: UpdateRequest :one
UPDATE requests
SET name = ?, method = ?, url = ?, content_type = ?, body = ?,
    params = ?, auth = ?, headers = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteRequest :exec
DELETE FROM requests WHERE id = ?;

-- name: ListRequestsByMethod :many
SELECT * FROM requests WHERE method = ? ORDER BY created_at DESC;
