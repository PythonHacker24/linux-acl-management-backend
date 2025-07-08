-- name: CreateResultsTransactionPQ :one
INSERT INTO results_transactions_archive (
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

-- name: GetResultsTransactionPQ :one
SELECT * FROM results_transactions_archive
WHERE id = $1;

-- name: GetResultsTransactionsBySessionPQ :many
SELECT * FROM results_transactions_archive
WHERE session_id = $1
ORDER BY timestamp DESC;

-- name: GetSuccessfulResultsTransactionsPQ :many
SELECT * FROM results_transactions_archive
WHERE session_id = $1 AND status = 'success'
ORDER BY timestamp DESC;

-- name: GetFailedResultsTransactionsPQ :many
SELECT * FROM results_transactions_archive
WHERE session_id = $1 AND status = 'failed'
ORDER BY timestamp DESC;

-- name: GetResultsTransactionsByOperationPQ :many
SELECT * FROM results_transactions_archive
WHERE session_id = $1 AND operation = $2
ORDER BY timestamp DESC;

-- name: GetResultsTransactionsByPathPQ :many
SELECT * FROM results_transactions_archive
WHERE session_id = $1 AND target_path = $2
ORDER BY timestamp DESC;

-- name: UpdateResultsTransactionStatusPQ :one
UPDATE results_transactions_archive
SET 
    status = $2,
    error_msg = $3,
    output = $4,
    duration_ms = $5
WHERE id = $1
RETURNING *;

-- name: DeleteResultsTransactionPQ :exec
DELETE FROM results_transactions_archive
WHERE id = $1;

-- name: DeleteResultsTransactionsBySessionPQ :exec
DELETE FROM results_transactions_archive
WHERE session_id = $1;

-- name: CountResultsTransactionsByStatusPQ :one
SELECT COUNT(*) FROM results_transactions_archive
WHERE session_id = $1 AND status = $2;

-- name: CountResultsTransactionsByOperationPQ :one
SELECT COUNT(*) FROM results_transactions_archive
WHERE session_id = $1 AND operation = $2;

-- name: GetResultsTransactionStatsPQ :one
SELECT 
    COUNT(*) as total_transactions,
    COUNT(CASE WHEN status = 'success' THEN 1 END) as successful_transactions,
    COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed_transactions,
    AVG(duration_ms) as avg_duration_ms
FROM results_transactions_archive
WHERE session_id = $1;
