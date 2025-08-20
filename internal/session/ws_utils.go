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

// /* build transaction stream data from a map */
// func buildTransactionStreamDataFromMap(data map[string]interface{}) (TransactionStreamData, error) {
// 	entriesRaw, _ := data["entries"].([]interface{})
// 	entries := make([]ACLEntryStream, 0, len(entriesRaw))
// 	for _, e := range entriesRaw {
// 		if entryMap, ok := e.(map[string]interface{}); ok {
// 			entry := ACLEntryStream{
// 				EntityType:  fmt.Sprintf("%v", entryMap["entityType"]),
// 				Entity:      fmt.Sprintf("%v", entryMap["entity"]),
// 				Permissions: fmt.Sprintf("%v", entryMap["permissions"]),
// 				Action:      fmt.Sprintf("%v", entryMap["action"]),
// 				Success:     entryMap["success"] == true,
// 				Error:       fmt.Sprintf("%v", entryMap["error"]),
// 			}
// 			entries = append(entries, entry)
// 		}
// 	}

// 	var timestamp time.Time
// 	if t, err := time.Parse(time.RFC3339, data["timestamp"].(string)); err == nil {
// 		timestamp = t
// 	}

// 	return TransactionStreamData{
// 		ID:         fmt.Sprintf("%v", data["id"]),
// 		SessionID:  fmt.Sprintf("%v", data["sessionId"]),
// 		Timestamp:  timestamp,
// 		Operation:  fmt.Sprintf("%v", data["operation"]),
// 		TargetPath: fmt.Sprintf("%v", data["targetPath"]),
// 		Entries:    entries,
// 		Status:     fmt.Sprintf("%v", data["status"]),
// 		ErrorMsg:   fmt.Sprintf("%v", data["errorMsg"]),
// 		Output:     fmt.Sprintf("%v", data["output"]),
// 		ExecutedBy: fmt.Sprintf("%v", data["executedBy"]),
// 		DurationMs: int64(mustFloat64(data["durationMs"])),
// 	}, nil
// }

// /* returns if value is of float type */
// func mustFloat64(val interface{}) float64 {
// 	if f, ok := val.(float64); ok {
// 		return f
// 	}
// 	return 0
// }
