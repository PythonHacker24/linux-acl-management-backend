-- name: StoreTransactionPQ :one
INSERT INTO transactions_archive (
    id, session_id, action, resource, permissions, status, error, output, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetTransactionPQ :one
SELECT * FROM transactions_archive 
WHERE id = $1;

-- name: GetTransactionsBySessionPQ :many
SELECT * FROM transactions_archive 
WHERE session_id = $1
ORDER BY created_at DESC;

-- name: GetSuccessfulTransactionsPQ :many
SELECT * FROM transactions_archive 
WHERE session_id = $1 AND status = 'success'
ORDER BY created_at DESC;

-- name: GetFailedTransactionsPQ :many
select * from transactions_archive 
where session_id = $1 and status = 'failure'
order by created_at desc;

-- name: GetPendingTransactionsPQ :many
select * from transactions_archive 
where session_id = $1 and status = 'pending'
order by created_at desc;

-- name: DeleteTransactionPQ :exec
DELETE FROM transactions_archive WHERE id = $1;

-- name: DeleteTransactionsBySessionPQ :exec
DELETE FROM transactions_archive WHERE session_id = $1;

-- name: CountTransactionsByStatusPQ :one
SELECT COUNT(*) FROM transactions_archive 
WHERE session_id = $1 AND status = $2;
