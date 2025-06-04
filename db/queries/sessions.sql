-- name: StoreSessionPQ :one
INSERT INTO sessions_archive (
    id, username, ip, user_agent, status, 
    created_at, last_active_at, expiry
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetSessionPQ :one
SELECT * FROM sessions_archive 
WHERE id = $1;

-- name: GetSessionByUsernamePQ :many
SELECT * FROM sessions_archive 
WHERE username = $1
ORDER BY created_at DESC;

-- name: DeleteSessionPQ :exec
DELETE FROM sessions_archive WHERE id = $1;
