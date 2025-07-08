-- name: CreatePendingTransactionPQ :one
INSERT INTO pending_transactions_archive (
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

-- name: GetPendingTransactionPQ :one
SELECT * FROM pending_transactions_archive
WHERE id = $1;

-- name: GetPendingTransactionsBySessionPQ :many
SELECT * FROM pending_transactions_archive
WHERE session_id = $1
ORDER BY timestamp DESC;

-- name: GetPendingTransactionsPQ :many
SELECT * FROM pending_transactions_archive
WHERE session_id = $1 AND status = 'pending'
ORDER BY timestamp DESC;

-- name: GetPendingTransactionsByOperationPQ :many
SELECT * FROM pending_transactions_archive
WHERE session_id = $1 AND operation = $2
ORDER BY timestamp DESC;

-- name: GetPendingTransactionsByPathPQ :many
SELECT * FROM pending_transactions_archive
WHERE session_id = $1 AND target_path = $2
ORDER BY timestamp DESC;

-- name: UpdatePendingTransactionStatusPQ :one
UPDATE pending_transactions_archive
SET 
    status = $2,
    error_msg = $3,
    output = $4,
    duration_ms = $5
WHERE id = $1
RETURNING *;

-- name: DeletePendingTransactionPQ :exec
DELETE FROM pending_transactions_archive
WHERE id = $1;

-- name: DeletePendingTransactionsBySessionPQ :exec
DELETE FROM pending_transactions_archive
WHERE session_id = $1;

-- name: CountPendingTransactionsByStatusPQ :one
SELECT COUNT(*) FROM pending_transactions_archive
WHERE session_id = $1 AND status = $2;

-- name: CountPendingTransactionsByOperationPQ :one
SELECT COUNT(*) FROM pending_transactions_archive
WHERE session_id = $1 AND operation = $2;

-- name: GetPendingTransactionStatsPQ :one
SELECT 
    COUNT(*) as total_transactions,
    COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_transactions,
    AVG(duration_ms) as avg_duration_ms
FROM pending_transactions_archive
WHERE session_id = $1;
