package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	/* TODO: fix the cyclic dependencies */

	"github.com/PythonHacker24/linux-acl-management-backend/internal/types"
)

/* we make use of Redis hashes for this application */

/* TODO: make the operations below thread safe with mutexes*/

/* store session into Redis database */
func (m *Manager) saveSession(username string) error {
	ctx := context.Background()

	/* thread safety for the manager */
	m.mutex.Lock()
	defer m.mutex.Unlock()

	/* find the session in session map */
	session, ok := m.sessionsMap[username]
	if !ok {
		return fmt.Errorf("username not found in session")
	}

	/* thread safety for the session */
	session.Mutex.Lock()
	defer session.Mutex.Unlock()

	/* session key for redis */
	key := fmt.Sprintf("session:%s", session.ID)

	/* serialize the session with relevant information */
	sessionSerialized := session.serializeSessionForRedis()

	if err := m.redis.HSet(ctx, key, sessionSerialized).Err(); err != nil {
		return fmt.Errorf("failed to save session to Redis: %w", err)
	}

	return nil
}

/* update expiry time in session */
func (m *Manager) updateSessionExpiry(username string) error {

	/*
		function expects that new expiry time is already set in the session
	*/

	ctx := context.Background()

	/* thread safety for the manager */
	m.mutex.Lock()
	defer m.mutex.Unlock()

	/* find the session in session map */
	session, ok := m.sessionsMap[username]
	if !ok {
		return fmt.Errorf("username not found in session")
	}

	/* thread safety for the session */
	session.Mutex.Lock()
	defer session.Mutex.Unlock()

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

/* TODO: Sessions must be marked expired when main.go exits */

/* update status of the session - update and set expired operations will be done with this */
func (m *Manager) updateSessionStatus(username string, status Status) error {

	ctx := context.Background()

	/* thread safety for the manager */
	m.mutex.Lock()
	defer m.mutex.Unlock()

	/* find the session in session map */
	session, ok := m.sessionsMap[username]
	if !ok {
		return fmt.Errorf("username not found in session")
	}

	/* thread safety for the session */
	session.Mutex.Lock()
	defer session.Mutex.Unlock()

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
func (m *Manager) saveTransactionResults(username string, txResult types.Transaction) error {

	ctx := context.Background()

	/* thread safety for the manager */
	m.mutex.Lock()
	defer m.mutex.Unlock()

	/* find the session in session map */
	session, ok := m.sessionsMap[username]
	if !ok {
		return fmt.Errorf("username not found in session")
	}

	/* thread safety for the session */
	session.Mutex.Lock()
	defer session.Mutex.Unlock()

	/* get the session ID */
	sessionID := session.ID

	/* create a key for Redis operation */
	key := fmt.Sprintf("session:%s:txresults", sessionID)

	/* marshal transaction result to JSON */
	resultBytes, err := json.Marshal(txResult)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction result: %w", err)
	}

	/* push the transaction result in the back of the list */
	return m.redis.RPush(ctx, key, resultBytes).Err()
}

func (m *Manager) getTransactionResults(username string, limit int) ([]types.TransactionResult, error) {
	ctx := context.Background()

	/* thread safety for the manager */
	m.mutex.Lock()
	defer m.mutex.Unlock()

	/* find the session in session map */
	session, ok := m.sessionsMap[username]
	if !ok {
		return nil, fmt.Errorf("username not found in session")
	}

	/* thread safety for the session */
	session.Mutex.Lock()
	defer session.Mutex.Unlock()

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
	results := make([]types.TransactionResult, 0, len(values))
	for _, val := range values {
		var result types.TransactionResult
		if err := json.Unmarshal([]byte(val), &result); err != nil {
			/* skip malformed results */
			continue
		}
		results = append(results, result)
	}

	return results, nil
}
