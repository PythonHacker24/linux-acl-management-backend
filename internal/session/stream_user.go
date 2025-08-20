package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/types"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

/* ==== User Session ==== */

/* send current session of a user */
func (m *Manager) sendCurrentSession(conn *websocket.Conn, sessionID string) error {
	ctx := context.Background()

	/* get data for current session from Redis */
	key := fmt.Sprintf("session:%s", sessionID)
	sessionData, err := m.redis.HGetAll(ctx, key).Result()
	if err != nil {
		/* user session doen't exists */
		if err == redis.Nil {
			message := StreamMessage{
				Type: "session_state",
				Data: map[string]any{
					"session_id": sessionID,
					"exists":     false,
				},
				Timestamp: time.Now(),
			}
			return conn.WriteJSON(message)
		}
		/* error cannot be handed, return a error */
		return fmt.Errorf("failed to get session data: %w", err)
	}

	/* session exists; convert Redis hash into session data */
	session := convertRedisHashToSession(sessionData)
	message := StreamMessage{
		Type: "session_state",
		Data: map[string]any{
			"session_id": sessionID,
			"exists":     true,
			"session":    session,
		},
		Timestamp: time.Now(),
	}

	return conn.WriteJSON(message)
}

/* send data regarding current session */
func (m *Manager) listenForSessionChanges(ctx context.Context, conn *websocket.Conn, sessionID string) {
	/* subscribe to both keyspace and keyevent notifications */
	keyspacePattern := fmt.Sprintf("__keyspace@0__:session:%s", sessionID)
	keyeventPattern := fmt.Sprintf("__keyevent@0__:hset:session:%s", sessionID)

	/* subscribe to Redis keyspace and keyevent */
	pubsub, err := m.redis.PSubscribe(ctx, keyspacePattern, keyeventPattern)
	if err != nil {
		m.errCh <- fmt.Errorf("failed to subscribe to redis events: %w", err)
		return
	}

	defer pubsub.Close()

	/* Redis update channel */
	ch := pubsub.Channel()

	/* handling session changes */
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			/* changes in session stored in Redis detected; handle the event */
			if err := m.handleSessionChangeEvent(conn, sessionID, msg); err != nil {
				m.errCh <- fmt.Errorf("error handling session change: %w", err)
			}
		}
	}
}

/* handle session change event */
func (m *Manager) handleSessionChangeEvent(conn *websocket.Conn, sessionID string, msg *redis.Message) error {
	ctx := context.Background()

	/* get session data from Redis */
	key := fmt.Sprintf("session:%s", sessionID)
	sessionData, err := m.redis.HGetAll(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to get updated session data: %w", err)
	}

	/* convert session data from Redis hash to session */
	session := convertRedisHashToSession(sessionData)

	/* prepare the message payload */
	message := StreamMessage{
		Type: "session_update",
		Data: map[string]any{
			"session_id":   sessionID,
			"session":      session,
			"event_type":   msg.Payload,
			"event_source": "redis_keyspace",
		},
		Timestamp: time.Now(),
	}

	/* send the message to the client */
	return conn.WriteJSON(message)
}

/* ==== User Transaction List ==== */

/* send current user results transactions */
func (m *Manager) sendCurrentUserTransactionsResults(conn *websocket.Conn, sessionID string, limit int) error {
	ctx := context.Background()

	/* get latest transactions from Redis */
	key := fmt.Sprintf("session:%s:txresults", sessionID)
	values, err := m.redis.LRange(ctx, key, int64(-limit), -1).Result()
	if err != nil {
		return fmt.Errorf("failed to get transaction results: %w", err)
	}

	/* convert each JSON string back into a Transaction */
	transactions := make([]types.Transaction, 0, len(values))
	for _, val := range values {
		var tx types.Transaction
		if err := json.Unmarshal([]byte(val), &tx); err != nil {
			/* skip malformed results */
			continue
		}
		transactions = append(transactions, tx)
	}

	/* prepare the message payload */
	message := StreamMessage{
		Type: "transaction_update",
		Data: map[string]any{
			"session_id":   sessionID,
			"transactions": transactions,
		},
		Timestamp: time.Now(),
	}

	/* send the message to the client */
	return conn.WriteJSON(message)
}

/* send current user pending transactions */
func (m *Manager) sendCurrentUserTransactionsPending(conn *websocket.Conn, sessionID string, limit int) error {
	ctx := context.Background()

	/* get latest transactions from Redis */
	key := fmt.Sprintf("session:%s:txpending", sessionID)
	values, err := m.redis.LRange(ctx, key, int64(-limit), -1).Result()
	if err != nil {
		return fmt.Errorf("failed to get transaction results: %w", err)
	}

	/* convert each JSON string back into a Transaction */
	transactions := make([]types.Transaction, 0, len(values))
	for _, val := range values {
		var tx types.Transaction
		if err := json.Unmarshal([]byte(val), &tx); err != nil {
			/* skip malformed results */
			continue
		}
		transactions = append(transactions, tx)
	}

	/* prepare the message payload */
	message := StreamMessage{
		Type: "transaction_update",
		Data: map[string]any{
			"session_id":   sessionID,
			"transactions": transactions,
		},
		Timestamp: time.Now(),
	}

	/* send the message to the client */
	return conn.WriteJSON(message)
}

/* listen for results transaction changes in Redis */
func (m *Manager) listenForTransactionsChangesResults(ctx context.Context, conn *websocket.Conn, sessionID string) {
	/* subscribe to both keyspace and keyevent notifications */
	keyspacePattern := fmt.Sprintf("__keyspace@0__:session:%s:txresults", sessionID)
	keyeventPattern := fmt.Sprintf("__keyevent@0__:rpush:session:%s:txresults", sessionID)

	/* subscribe to Redis keyspace and keyevent */
	pubsub, err := m.redis.PSubscribe(ctx, keyspacePattern, keyeventPattern)
	if err != nil {
		m.errCh <- fmt.Errorf("failed to subscribe to redis events: %w", err)
		return
	}

	defer pubsub.Close()

	/* Redis update channel */
	ch := pubsub.Channel()

	/* handling transaction changes */
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			/* changes in transactions stored in Redis detected; handle the event */
			if err := m.handleTransactionChangeEventResults(conn, sessionID, msg); err != nil {
				m.errCh <- fmt.Errorf("error handling transaction change: %w", err)
			}
		}
	}
}

/* listen for pending transaction changes in Redis */
func (m *Manager) listenForTransactionsChangesPending(ctx context.Context, conn *websocket.Conn, sessionID string) {
	/* subscribe to both keyspace and keyevent notifications */
	keyspacePattern := fmt.Sprintf("__keyspace@0__:session:%s:txpending", sessionID)
	keyeventPattern := fmt.Sprintf("__keyevent@0__:rpush:session:%s:txpending", sessionID)

	/* subscribe to Redis keyspace and keyevent */
	pubsub, err := m.redis.PSubscribe(ctx, keyspacePattern, keyeventPattern)
	if err != nil {
		m.errCh <- fmt.Errorf("failed to subscribe to redis events: %w", err)
		return
	}

	defer pubsub.Close()

	/* Redis update channel */
	ch := pubsub.Channel()

	/* handling transaction changes */
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			/* changes in transactions stored in Redis detected; handle the event */
			if err := m.handleTransactionChangeEventPending(conn, sessionID, msg); err != nil {
				m.errCh <- fmt.Errorf("error handling transaction change: %w", err)
			}
		}
	}
}

/*
	currently, handleTransactionChangeEvent sends the complete JSON package whenever anything is updated.
	The whole frontend will be updated even if one transaction changes it's state (for example, setting active to expired).
*/

/* handle transaction results change event */
func (m *Manager) handleTransactionChangeEventResults(conn *websocket.Conn, sessionID string, msg *redis.Message) error {
	ctx := context.Background()

	/* get latest transactions */
	key := fmt.Sprintf("session:%s:txresults", sessionID)
	values, err := m.redis.LRange(ctx, key, -100, -1).Result()
	if err != nil {
		return fmt.Errorf("failed to get transaction results: %w", err)
	}

	/* convert each JSON string back into a Transaction */
	transactions := make([]types.Transaction, 0, len(values))
	for _, val := range values {
		var tx types.Transaction
		if err := json.Unmarshal([]byte(val), &tx); err != nil {
			/* skip malformed results */
			continue
		}
		transactions = append(transactions, tx)
	}

	/* prepare the message payload */
	message := StreamMessage{
		Type: "transaction_update",
		Data: map[string]any{
			"session_id":   sessionID,
			"transactions": transactions,
			"event_type":   msg.Payload,
			"event_source": "redis_keyspace",
		},
		Timestamp: time.Now(),
	}

	/* send the message to the client */
	return conn.WriteJSON(message)
}

/* handle transaction pending change event */
func (m *Manager) handleTransactionChangeEventPending(conn *websocket.Conn, sessionID string, msg *redis.Message) error {
	ctx := context.Background()

	/* get latest transactions */
	key := fmt.Sprintf("session:%s:txpending", sessionID)
	values, err := m.redis.LRange(ctx, key, -100, -1).Result()
	if err != nil {
		return fmt.Errorf("failed to get pending transactions: %w", err)
	}

	/* convert each JSON string back into a Transaction */
	transactions := make([]types.Transaction, 0, len(values))
	for _, val := range values {
		var tx types.Transaction
		if err := json.Unmarshal([]byte(val), &tx); err != nil {
			/* skip malformed results */
			continue
		}
		transactions = append(transactions, tx)
	}

	/* prepare the message payload */
	message := StreamMessage{
		Type: "transaction_update",
		Data: map[string]any{
			"session_id":   sessionID,
			"transactions": transactions,
			"event_type":   msg.Payload,
			"event_source": "redis_keyspace",
		},
		Timestamp: time.Now(),
	}

	/* send the message to the client */
	return conn.WriteJSON(message)
}
