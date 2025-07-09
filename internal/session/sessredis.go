package session

import (
	"context"
	"fmt"
	"time"
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

/* increment the failed field of the session in Redis */
func (m *Manager) IncrementSessionFailedRedis(session *Session) error {
	ctx := context.Background()
	key := fmt.Sprintf("session:%s", session.ID)

	if err := m.redis.HIncrBy(ctx, key, "failed", 1).Err(); err != nil {
		return fmt.Errorf("failed to increment failed count in Redis: %w", err)
	}
	return nil
}

/* increment the completed field of the session in Redis */
func (m *Manager) IncrementSessionCompletedRedis(session *Session) error {
	ctx := context.Background()
	key := fmt.Sprintf("session:%s", session.ID)

	if err := m.redis.HIncrBy(ctx, key, "completed", 1).Err(); err != nil {
		return fmt.Errorf("failed to increment completed count in Redis: %w", err)
	}
	return nil
}
