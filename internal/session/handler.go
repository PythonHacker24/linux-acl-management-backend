package session

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/PythonHacker24/linux-acl-management-backend/api/middleware"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/types"
)

/*
	TODO: watchdog for session
	Live sessions and transactions can be monitored through Redis and PostgreSQL
	the watchdog here shows the processing happening, which needs to be done in the
	later stages of development
*/

/* frontend safe handler for issuing transaction */
func (m *Manager) IssueTransaction(w http.ResponseWriter, r *http.Request) {
	/* extract username from JWT Token */
	username := r.Context().Value(middleware.ContextKeyUsername)

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

type handlerCtxKey string

const (
	StreamUserSession      handlerCtxKey = "stream_user_session"
	StreamUserTransactions handlerCtxKey = "stream_user_transactions"
	StreamAllSessions      handlerCtxKey = "stream_all_sessions"
	StreamAllTransactions  handlerCtxKey = "stream_all_transactions"
)

/*
get single session data
requires user authentication from middleware
user/
*/
func (m *Manager) StreamUserSession(w http.ResponseWriter, r *http.Request) {

	/* username := r.Context().Value(middleware.ContextKeyUsername) */
	sessionID := r.Context().Value(middleware.ContextKeySessionID).(string)

	/* add a check for sessionID belongs to user */
	conn, err := m.upgrader.Upgrade(w, r, nil)
	if err != nil {
		m.errCh <- fmt.Errorf("websocket upgrade error: %w", err)
		return
	}
	defer conn.Close()

	/* 
		context with cancel for web socket handlers 
		this is the official context for a websocket connection
		cancelling this means closing components of the websocket handler
	*/
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/* sending initial session data */
	if err := m.sendCurrentSession(conn, sessionID); err != nil {
		m.errCh <- fmt.Errorf("error sending initial session: %w", err)
		return
	}

	/* stream changes in session made in Redis */
	go m.listenForSessionChanges(ctx, conn, sessionID)

	/* specify the handler context */
	ctxVal := context.WithValue(ctx, "type", StreamUserSession)

	/* handle web socket instructions from client */
	m.handleWebSocketCommands(conn, sessionID, ctxVal, cancel)
}

/*
// get user transactions information
// requires user authentication from middleware
// user/
// */
// func (m *manager) streamusertransactions(w http.responsewriter, r *http.request) {
// 	/* username := r.context().value(middleware.contextkeyusername) */
// 	sessionid := r.context().value(middleware.contextkeysessionid)
//
// 	/* add a check for sessionid belongs to user */
// 	conn, err := m.upgrader.upgrade(w, r, nil)
// 	if err != nil {
// 		m.errch <- fmt.errorf("websocket upgrade error: %w", err)
// 		return
// 	}
// 	defer conn.close()
//
// 	/* context with cancel for web socket handlers */
// 	ctx, cancel := context.withcancel(context.background())
// 	defer cancel()
//
// 	/* sending initial list of transactions data */
// 	if err := m.sendcurrenttransactions(conn, sessionid); err != nil {
// 		// log.printf("error sending initial session: %v", err)
// 		m.errch <- fmt.errorf("error sending initial session: %w", err)
// 		return
// 	}
//
// 	/* stream changes in transactions made in redis */
// 	go m.listenfortransactionschanges(ctx, conn, sessionid)
//
// 	/* handle web socket instructions from client */
// 	m.handlewebsocketcommands(conn, cancel)
// }
//
// /*
// get all sessions in the system
// requires admin authentication from middleware
// admin/
// */
// func (m *manager) streamallsessions(w http.responsewriter, r *http.request) {
//
// }
//
// /*
// get all transaction in the system
// requires admin authentication from middleware
// admin/
// */
// func (m *manager) streamalltransactions(w http.responsewriter, r *http.request) {
//
// }
