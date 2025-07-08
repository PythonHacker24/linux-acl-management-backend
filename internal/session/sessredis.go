package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/types"
)

/* TODO: make the operations below thread safe with mutexes*/

/* store session into Redis database */
func (m *Manager) saveSessionRedis(session *Session) error {
	ctx := context.Background()

	/* session key for redis */
	key := fmt.Sprintf("session:%s", session.ID)

	/* serialize the session with relevant information */
	sessionSerialized := session.serializeSessionForRedis()

	/* hset the session to redis */
	if err := m.redis.HSet(ctx, key, sessionSerialized).Err(); err != nil {
		return fmt.Errorf("failed to save session to Redis: %w", err)
	}

	return nil
}

/* update expiry time in session */
func (m *Manager) updateSessionExpiryRedis(session *Session) error {

	/*
		function expects that new expiry time is already set in the session
	*/

	ctx := context.Background()

	/* create a key for Redis operation */
	key := fmt.Sprintf("session:%s", session.ID)

	/* convert the expiry time to  */
	formattedExpiry := session.Expiry.Format(time.RFC3339)

	/* update just the expiry field */
	err := m.redis.HSet(ctx, key, "expiry", formattedExpiry).Err()
	if err != nil {
		return fmt.Errorf("failed to update session expiry in Redis: %w", err)
	}

	return nil
}

/* update status of the session - update and set expired operations will be done with this */
func (m *Manager) updateSessionStatusRedis(session *Session, status Status) error {

	ctx := context.Background()

	/* create a key for Redis operation */
	key := fmt.Sprintf("session:%s", session.ID)

	/* update the session status */
	err := m.redis.HSet(ctx, key, "status", status).Err()
	if err != nil {
		return fmt.Errorf("failed to mark session as expired in Redis: %w", err)
	}

	return nil
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
