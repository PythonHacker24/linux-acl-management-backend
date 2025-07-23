package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/postgresql"
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
				Data: map[string]interface{}{
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
		Data: map[string]interface{}{
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
		Data: map[string]interface{}{
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

/* send current user transactions */
func (m *Manager) sendCurrentUserTransactions(conn *websocket.Conn, username, sessionID string, limit int) error {
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
		Data: map[string]interface{}{
			"session_id":   sessionID,
			"transactions": transactions,
		},
		Timestamp: time.Now(),
	}

	/* send the message to the client */
	return conn.WriteJSON(message)
}

/* listen for transaction changes in Redis */
func (m *Manager) listenForTransactionsChanges(ctx context.Context, conn *websocket.Conn, sessionID string) {
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
			if err := m.handleTransactionChangeEvent(conn, sessionID, msg); err != nil {
				m.errCh <- fmt.Errorf("error handling transaction change: %w", err)
			}
		}
	}
}

/*
	currently, handleTransactionChangeEvent sends the complete JSON package whenever anything is updated. 
	The whole frontend will be updated even if one transaction changes it's state (for example, setting active to expired).
*/

/* handle transaction change event */
func (m *Manager) handleTransactionChangeEvent(conn *websocket.Conn, sessionID string, msg *redis.Message) error {
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
		Data: map[string]interface{}{
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

/* ==== User Archived Session ==== */

/* send list of archived session of a user */
func (m *Manager) sendCurrentArchivedSessions(conn *websocket.Conn, username string, page int, pageSize int) error {
	ctx := context.Background()

	/* calculate LIMIT and OFFSET */
	limit := int32(pageSize)
	offset := int32((page - 1) * pageSize)

	pqParams := &postgresql.GetSessionByUsernamePaginatedPQParams{
		Username: username,
		Limit: limit,
		Offset: offset,
	}

	/* query sessions */
	sessions, err := m.archivalPQ.GetSessionByUsernamePaginatedPQ(ctx, *pqParams)
	if err != nil {
		return fmt.Errorf("failed to get sessions for username %s: %w", username, err)
	}

	exists := len(sessions) > 0

	/* convert to plain JSON-compatible slices */
	var outgoing []map[string]interface{}
	for _, session := range sessions {
		outgoing = append(outgoing, map[string]interface{}{
			"id":              session.ID,
			"username":        session.Username,
			"ip":              session.Ip.String,
			"user_agent":      session.UserAgent.String,
			"status":          session.Status,
			"created_at":      session.CreatedAt,
			"last_active_at":  session.LastActiveAt,
			"expiry":          session.Expiry,
			"completed_count": session.CompletedCount,
			"failed_count":    session.FailedCount,
			"archived_at":     session.ArchivedAt,
		})
	}

	message := StreamMessage{
		Type: "session_state",
		Data: map[string]interface{}{
			"username": username,
			"exists":   exists,
			"sessions": outgoing,
			"page":     page,
			"pageSize": pageSize,
		},
		Timestamp: time.Now(),
	}

	return conn.WriteJSON(message)
}

/* ==== User Archived Transactions ==== */

/* send list of archived pending transactions */
func (m *Manager) sendCurrentArchivedPendingTransactions(conn *websocket.Conn, username string, page int, pageSize int) error {
	ctx := context.Background()

	limit := int32(pageSize)
	offset := int32((page - 1) * pageSize)

	pqParams := &postgresql.GetPendingTransactionsByUserPaginatedPQParams{
		ExecutedBy: username,
		Limit: limit,
		Offset: offset,
	}

	/* query transactions by executed_by */
	transactions, err := m.archivalPQ.GetPendingTransactionsByUserPaginatedPQ(ctx, *pqParams)
	if err != nil {
		return fmt.Errorf("failed to get transactions for user %s: %w", username, err)
	}

	exists := len(transactions) > 0

	var outgoing []map[string]interface{}
	for _, tx := range transactions {
		outgoing = append(outgoing, map[string]interface{}{
			"id":           tx.ID,
			"session_id":   tx.SessionID,
			"timestamp":    tx.Timestamp,
			"operation":    tx.Operation,
			"target_path":  tx.TargetPath,
			"entries":      tx.Entries, // JSONB will decode as []byte, so decode if needed
			"status":       tx.Status,
			"error_msg":    tx.ErrorMsg,
			"output":       tx.Output,
			"executed_by":  tx.ExecutedBy,
			"duration_ms":  tx.DurationMs,
			"ExecStatus":   tx.Execstatus,
			"created_at":   tx.CreatedAt,
		})
	}

	message := StreamMessage{
		Type: "transactions_state",
		Data: map[string]interface{}{
			"username":    username,
			"exists":      exists,
			"transactions": outgoing,
			"page":        page,
			"pageSize":    pageSize,
		},
		Timestamp: time.Now(),
	}

	return conn.WriteJSON(message)
}
