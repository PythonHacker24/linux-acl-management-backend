package session

import (
	"context"
	"fmt"
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
