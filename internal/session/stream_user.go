package session

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

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
	keyspacePattern := fmt.Sprintf("__keyspace@0__:session:%s", sessionID)

	/* subscribe to Redis keyspace */
	pubsub, err := m.redis.PSubscribe(ctx, keyspacePattern)
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

	/* send the message payload via websocket */
	return conn.WriteJSON(message)
}
