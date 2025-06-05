-- name: CreateTransactionPQ :one
INSERT INTO transactions_archive (
    id,
    session_id,
    timestamp,
    operation,
    target_path,
    entries,
    status,
    error_msg,
    output,
    executed_by,
    duration_ms
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetTransactionPQ :one
SELECT * FROM transactions_archive
WHERE id = $1;

-- name: GetTransactionsBySessionPQ :many
SELECT * FROM transactions_archive
WHERE session_id = $1
ORDER BY timestamp DESC;

-- name: GetSuccessfulTransactionsPQ :many
SELECT * FROM transactions_archive
WHERE session_id = $1 AND status = 'success'
ORDER BY timestamp DESC;

-- name: GetFailedTransactionsPQ :many
SELECT * FROM transactions_archive
WHERE session_id = $1 AND status = 'failed'
ORDER BY timestamp DESC;

-- name: GetPendingTransactionsPQ :many
SELECT * FROM transactions_archive
WHERE session_id = $1 AND status = 'pending'
ORDER BY timestamp DESC;

-- name: GetTransactionsByOperationPQ :many
SELECT * FROM transactions_archive
WHERE session_id = $1 AND operation = $2
ORDER BY timestamp DESC;

-- name: GetTransactionsByPathPQ :many
SELECT * FROM transactions_archive
WHERE session_id = $1 AND target_path = $2
ORDER BY timestamp DESC;

-- name: UpdateTransactionStatusPQ :one
UPDATE transactions_archive
SET 
    status = $2,
    error_msg = $3,
    output = $4,
    duration_ms = $5
WHERE id = $1
RETURNING *;

-- name: DeleteTransactionPQ :exec
DELETE FROM transactions_archive
WHERE id = $1;

-- name: DeleteTransactionsBySessionPQ :exec
DELETE FROM transactions_archive
WHERE session_id = $1;

-- name: CountTransactionsByStatusPQ :one
SELECT COUNT(*) FROM transactions_archive
WHERE session_id = $1 AND status = $2;

-- name: CountTransactionsByOperationPQ :one
SELECT COUNT(*) FROM transactions_archive
WHERE session_id = $1 AND operation = $2;

-- name: GetTransactionStatsPQ :one
SELECT 
    COUNT(*) as total_transactions,
    COUNT(CASE WHEN status = 'success' THEN 1 END) as successful_transactions,
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_transactions,
    COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_transactions,
    AVG(duration_ms) as avg_duration_ms
FROM transactions_archive
WHERE session_id = $1;
