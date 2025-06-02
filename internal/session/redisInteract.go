package session

import (
	"context"
	"fmt"
	"time"
)

/* we make use of Redis hashes for this application */

/* store session into Redis database */
func (m *Manager) saveSession(username string) error {
	ctx := context.Background()

	/* find the session in session map */
	session, ok := m.sessionsMap[username]
	if !ok {
		return fmt.Errorf("username not found in session")
	}

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

	/* find the session in session map */
	session, ok := m.sessionsMap[username]
	if !ok {
		return fmt.Errorf("username not found in session")
	}

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

	/* find the session in session map */
	session, ok := m.sessionsMap[username]
	if !ok {
		return fmt.Errorf("username not found in session")
	}

	/* create a key for Redis operation */
	key := fmt.Sprintf("session:%s", session.ID)

	/* update the session status */
	err := m.redis.HSet(ctx, key, "status", status).Err()
	if err != nil {
		return fmt.Errorf("failed to mark session as expired in Redis: %w", err)
	}

	return nil
}
