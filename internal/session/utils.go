package session

import "time"

/* serialize session information to store in Redis */
func (s *Session) serializeSessionForRedis() map[string]interface{} {
	return map[string]interface{} {
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
