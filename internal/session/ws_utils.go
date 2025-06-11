package session

import (
	"strconv"
	"time"
)

/* convert redis hash to session */
func convertRedisHashToSession(hash map[string]string) SessionStreamData {
	session := SessionStreamData{}

	if val, ok := hash["id"]; ok {
		session.ID = val
	}

	if val, ok := hash["username"]; ok {
		session.Username = val
	}

	if val, ok := hash["ip"]; ok {
		session.IP = val
	}

	if val, ok := hash["user_agent"]; ok {
		session.UserAgent = val
	}

	if val, ok := hash["status"]; ok {
		session.Status = val
	}

	if val, ok := hash["created_at"]; ok {
		if t, err := time.Parse(time.RFC3339, val); err == nil {
			session.CreatedAt = t
		}
	}

	if val, ok := hash["last_active_at"]; ok {
		if t, err := time.Parse(time.RFC3339, val); err == nil {
			session.LastActiveAt = t
		}
	}

	if val, ok := hash["expiry"]; ok {
		if t, err := time.Parse(time.RFC3339, val); err == nil {
			session.Expiry = t
		}
	}

	if val, ok := hash["completed"]; ok {
		if i, err := strconv.Atoi(val); err == nil {
			session.CompletedCount = i
		}
	}

	if val, ok := hash["failed"]; ok {
		if i, err := strconv.Atoi(val); err == nil {
			session.FailedCount = i
		}
	}

	return session
}
