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
	username, ok := r.Context().Value(middleware.ContextKeyUsername).(string)
	if !ok {
		http.Error(w, "Invalid user context", http.StatusInternalServerError)
		return
	}

	/* acquire manager lock to access sessions map */
	m.mutex.RLock()
	session := m.sessionsMap[username]
	m.mutex.RUnlock()

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

	tx := types.Transaction{
		ID:         uuid.New(),
		SessionID:  session.ID,
		Timestamp:  time.Now(),
		Operation:  req.Operation,
		TargetPath: req.TargetPath,
		Entries:    req.Entries,
		Status:     types.StatusPending,
		ExecutedBy: username,
	}

	/* add transaction to session - session lock is already held */
	if err := m.AddTransaction(session, &tx); err != nil {
		http.Error(w, "Failed to add transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"message": "Transaction scheduled",
		"txn_id":  tx.ID.String(),
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

type handlerCtxKey string

const (
	CtxStreamUserSession      handlerCtxKey = "stream_user_session"
	CtxStreamUserTransactions handlerCtxKey = "stream_user_transactions"
	CtxStreamAllSessions      handlerCtxKey = "stream_all_sessions"
	CtxStreamAllTransactions  handlerCtxKey = "stream_all_transactions"
)

/*
get single session data
requires user authentication from middleware
user/
*/
func (m *Manager) StreamUserSession(w http.ResponseWriter, r *http.Request) {

	/* get the username */
	username, ok := r.Context().Value(middleware.ContextKeyUsername).(string)
	if !ok {
		http.Error(w, "Invalid user context", http.StatusInternalServerError)
		return
	}

	/* get the session id */
	sessionID, ok := r.Context().Value(middleware.ContextKeySessionID).(string)
	if !ok {
		http.Error(w, "Invalid session context", http.StatusInternalServerError)
		return
	}

	m.mutex.RLock()
	session, exists := m.sessionsMap[username]
	m.mutex.RUnlock()

	if !exists || session.ID.String() != sessionID {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	/* user exists and verified, upgrade the websocket connection */
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
	ctxVal := context.WithValue(ctx, "type", CtxStreamUserSession)

	/* handle web socket instructions from client */
	m.handleWebSocketCommands(conn, username, sessionID, ctxVal, cancel)
}

/*
get user transactions information
requires user authentication from middleware
user/
*/
func (m *Manager) StreamUserTransactions(w http.ResponseWriter, r *http.Request) {

	/* get the username */
	username, ok := r.Context().Value(middleware.ContextKeyUsername).(string)
	if !ok {
		http.Error(w, "Invalid user context", http.StatusInternalServerError)
		return
	}

	/* get the session id */
	sessionID, ok := r.Context().Value(middleware.ContextKeySessionID).(string)
	if !ok {
		http.Error(w, "Invalid session ID context", http.StatusInternalServerError)
		return
	}

	m.mutex.RLock()
	session, exists := m.sessionsMap[username]
	m.mutex.RUnlock()

	if !exists || session.ID.String() != sessionID {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	/* user exists and verified, upgrade the websocket connection */
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

	/* sending initial list of transactions data */
	if err := m.sendCurrentUserTransactions(conn, username, sessionID, 100); err != nil {
		m.errCh <- fmt.Errorf("error sending initial transactions: %w", err)
		return
	}

	/* stream changes in transactions made in redis */
	go m.listenForTransactionsChanges(ctx, conn, sessionID)

	/* specify the handler context */
	ctxVal := context.WithValue(ctx, "type", CtxStreamUserTransactions)

	/* handle web socket instructions from client */
	m.handleWebSocketCommands(conn, username, sessionID, ctxVal, cancel)
}

// /*
// get all sessions in the system
// requires admin authentication from middleware
// admin/
// */
// func (m *Manager) StreamAllSessions(w http.ResponseWriter, r *http.Request) {
//
// 	/* check if the user is admin */
//
// 	/* upgrade the connection if user is admin */
// 	conn, err := m.upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		m.errCh <- fmt.Errorf("websocket upgrade error: %w", err)
// 		return
// 	}
// 	defer conn.Close()
//
// 	/*
// 		context with cancel for web socket handlers
// 		this is the official context for a websocket connection
// 		cancelling this means closing components of the websocket handler
// 	*/
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
//
// 	/* sending initial list of all sessions */
// 	if err := m.sendListofAllSessions(conn, 100); err != nil {
// 		m.errCh <- fmt.Errorf("error sending initial list of all sessions: %w", err)
// 		return
// 	}
//
// 	/* stream changes in transactions made in redis */
// 	go m.listenForAllSessionsChanges(ctx, conn)
//
// 	/* specify the handler context */
// 	ctxVal := context.WithValue(ctx, "type", CtxStreamAllSessions)
//
// 	/* handle web socket instructions from client */
// 	m.handleWebSocketCommands(conn, ctxVal, cancel)
// }
//
// /*
// get all transaction in the system
// requires admin authentication from middleware
// admin/
// */
// func (m *Manager) StreamAllTransactions(w http.ResponseWriter, r *http.Request) {
//
// 	/* check if the user is admin */
//
// 	/* upgrade the connection if user is admin */
// 	conn, err := m.upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		m.errCh <- fmt.Errorf("websocket upgrade error: %w", err)
// 		return
// 	}
// 	defer conn.Close()
//
// 	/*
// 		context with cancel for web socket handlers
// 		this is the official context for a websocket connection
// 		cancelling this means closing components of the websocket handler
// 	*/
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
//
// 	/* sending initial list of all sessions */
// 	if err := m.sendListofAllTransactions(conn, 100); err != nil {
// 		m.errCh <- fmt.Errorf("error sending initial list of all transactions: %w", err)
// 		return
// 	}
//
// 	/* stream changes in transactions made in redis */
// 	go m.listenForAllTransactionsChanges(ctx, conn)
//
// 	/* specify the handler context */
// 	ctxVal := context.WithValue(ctx, "type", CtxStreamAllTransactions)
//
// 	/* handle web socket instructions from client */
// 	m.handleWebSocketCommands(conn, ctxVal, cancel)
// }
