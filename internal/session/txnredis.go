package session

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/types"
)

/* save a pending transaction in Redis as transactionID -> JSON in session:sessionID:txnpending */
func (m *Manager) SavePendingTransaction(session *Session, tx *types.Transaction) error {
	ctx := context.Background()

	/* get the session ID */
	sessionID := session.ID

	/* create the Redis key for pending transactions */
	key := fmt.Sprintf("session:%s:txnpending", sessionID)

	/* marshal transaction to JSON */
	txBytes, err := json.Marshal(tx)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}

	/* use HSET to store transactionID -> JSON */
	return m.redis.HSet(ctx, key, tx.ID.String(), txBytes).Err()
}

/* remove a pending transaction by ID from Redis HASH session:<sessionID>:txnpending */
func (m *Manager) RemovePendingTransaction(session *Session, txnID uuid.UUID) error {
	ctx := context.Background()

	sessionID := session.ID
	key := fmt.Sprintf("session:%s:txnpending", sessionID)

	/* remove the transaction ID field from the hash */
	return m.redis.HDel(ctx, key, txnID.String()).Err()
}

/* returns latest results of processed transactions */
func (m *Manager) getTransactionResultsRedis(session *Session, limit int) ([]types.Transaction, error) {
	ctx := context.Background()

	/* get the session ID */
	sessionID := session.ID

	/* create a key for Redis operation */
	key := fmt.Sprintf("session:%s:txresults", sessionID)

	/* returns transactions in chronological order */
	values, err := m.redis.LRange(ctx, key, int64(-limit), -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction results: %w", err)
	}

	/* converts each JSON string back into a TransactionResult */
	results := make([]types.Transaction, 0, len(values))
	for _, val := range values {
		var result types.Transaction
		if err := json.Unmarshal([]byte(val), &result); err != nil {
			m.errCh <- fmt.Errorf("failed to unmarshal transaction result: %w; value: %s", err, val)
			continue
		}
		results = append(results, result)
	}

	return results, nil
}

/* save transaction results to redis */
func (m *Manager) SaveTransactionRedisList(session *Session, txResult *types.Transaction, list string) error {

	ctx := context.Background()

	/* get the session ID */
	sessionID := session.ID

	/* create a key for Redis operation */
	key := fmt.Sprintf("session:%s:%s", sessionID, list)

	/* marshal transaction result to JSON */
	resultBytes, err := json.Marshal(txResult)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction result: %w", err)
	}

	/* push the transaction result in the back of the list */
	return m.redis.RPush(ctx, key, resultBytes).Err()
}
