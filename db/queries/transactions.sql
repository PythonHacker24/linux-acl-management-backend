-- name: CreateTransaction :one
INSERT INTO transactions_archive (
    id, session_id, status, output, created_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetTransaction :one
SELECT * FROM transactions_archive 
WHERE id = $1;

-- name: GetTransactionsBySession :many
SELECT * FROM transactions_archive 
WHERE session_id = $1
ORDER BY created_at DESC;

-- name: GetSuccessfulTransactions :many
SELECT * FROM transactions_archive 
WHERE session_id = $1 AND status = 'success'
ORDER BY created_at DESC;

-- name: GetFailedTransactions :many
SELECT * FROM transactions_archive 
WHERE session_id = $1 AND status = 'failure'
ORDER BY created_at DESC;

-- name: DeleteTransaction :exec
DELETE FROM transactions_archive WHERE id = $1;

-- name: DeleteTransactionsBySession :exec
DELETE FROM transactions_archive WHERE session_id = $1;

-- name: CountTransactionsByStatus :one
SELECT COUNT(*) FROM transactions_archive 
WHERE session_id = $1 AND status = $2;
