package session

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/types"
	"github.com/google/uuid"
)

/*
	TODO: watchdog for session
	Live sessions and transactions can be montired through Redis and PostgreSQL
	the watchdog here shows the processing happening, which needs to be done in the
	later stages of development
*/

/* frontend safe handler for issuing transaction */
func (m *Manager) IssueTransaction(w http.ResponseWriter, r *http.Request) {
	/* extract username from JWT Token */
	username := r.Context().Value("username")

	/* acquire manager lock to access sessions map */
	m.mutex.Lock()
	session := m.sessionsMap[username.(string)]
	m.mutex.Unlock()

	if session == nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	/* acquire session lock for transaction operations */
	session.Mutex.Lock()
	defer session.Mutex.Unlock()

	var req types.ScheduleTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	tx := &types.Transaction{
		ID:         uuid.New(),
		SessionID:  session.ID,
		Timestamp:  time.Now(),
		Operation:  req.Operation,
		TargetPath: req.TargetPath,
		Entries:    req.Entries,
		Status:     types.StatusPending,
		ExecutedBy: username.(string),
	}

	/* add transaction to session - session lock is already held */
	if err := m.AddTransaction(session, tx); err != nil {
		http.Error(w, "Failed to add transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Transaction scheduled",
		"txn_id":  tx.ID.String(),
	})
}
