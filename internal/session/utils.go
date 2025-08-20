package session

import "time"

/* serialize session information to store in Redis */
func (s *Session) serializeSessionForRedis() map[string]any {
	return map[string]any{
		"id":             s.ID,
		"username":       s.Username,
		"ip":             s.IP,
		"user_agent":     s.UserAgent,
		"status":         "active",
		"created_at":     s.CreatedAt.Format(time.RFC3339),
		"last_active_at": s.LastActiveAt.Format(time.RFC3339),
		"expiry":         s.Expiry.Format(time.RFC3339),
		"completed":      s.CompletedCount,
		"failed":         s.FailedCount,
	}
}

/* returns all the usernames in the manager */
func (m *Manager) GetAllUsernames() []string {
	/* thread safety of manager */
	m.mutex.Lock()
	defer m.mutex.Unlock()

	/* create and fill slice for usernames */
	usernames := make([]string, 0, len(m.sessionsMap))
	for _, session := range m.sessionsMap {
		usernames = append(usernames, session.Username)
	}

	return usernames
}
