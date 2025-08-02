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
	CtxStreamUserSession					handlerCtxKey = "stream_user_session"
	CtxStreamUserTransactionsResults		handlerCtxKey = "stream_user_transactions_results"
	CtxStreamUserTransactionsPending		handlerCtxKey = "stream_user_transactions_pending"
	CtxStreamAllSessions      				handlerCtxKey = "stream_all_sessions"
	CtxStreamAllTransactions  				handlerCtxKey = "stream_all_transactions"
	CtxStreamUserArchiveSession 			handlerCtxKey = "stream_user_archive_sessions"
	CtxStreamUserArchiveResultsTransactions handlerCtxKey = "stream_user_archive_results_transactions"
	CtxStreamUserArchivePendingTransactions handlerCtxKey = "stream_user_archive_pending_transactions"
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
get user transactions results information
requires user authentication from middleware
user/
*/
func (m *Manager) StreamUserTransactionsResults(w http.ResponseWriter, r *http.Request) {

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
	if err := m.sendCurrentUserTransactionsResults(conn, sessionID, 100); err != nil {
		m.errCh <- fmt.Errorf("error sending initial transactions: %w", err)
		return
	}

	/* stream changes in transactions made in redis */
	go m.listenForTransactionsChangesResults(ctx, conn, sessionID)

	/* specify the handler context */
	ctxVal := context.WithValue(ctx, "type", CtxStreamUserTransactionsResults)

	/* handle web socket instructions from client */
	m.handleWebSocketCommands(conn, username, sessionID, ctxVal, cancel)
}

/*
get user transactions pending information
requires user authentication from middleware
user/
*/
func (m *Manager) StreamUserTransactionsPending(w http.ResponseWriter, r *http.Request) {

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
	if err := m.sendCurrentUserTransactionsPending(conn, sessionID, 100); err != nil {
		m.errCh <- fmt.Errorf("error sending initial transactions: %w", err)
		return
	}

	/* stream changes in transactions made in redis */
	go m.listenForTransactionsChangesPending(ctx, conn, sessionID)

	/* specify the handler context */
	ctxVal := context.WithValue(ctx, "type", CtxStreamUserTransactionsPending)

	/* handle web socket instructions from client */
	m.handleWebSocketCommands(conn, username, sessionID, ctxVal, cancel)
}

/*
get user archived sessions information
requires user authentication from middleware
user/
*/
func (m *Manager) StreamUserArchiveSessions(w http.ResponseWriter, r *http.Request) {
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

	/* sending initial session data */
	if err := m.sendCurrentArchivedSessions(conn, username, 1, 1); err != nil {
		m.errCh <- fmt.Errorf("error sending initial session: %w", err)
		return
	}

	/* 
		don't send updated stream of data
		archive is updated when session expires
		so user will have no access to the dashboard and can't see new
	*/

	/* specify the handler context */
	ctxVal := context.WithValue(ctx, "type", CtxStreamUserArchiveSession)

	/* handle web socket instructions from client */
	m.handleWebSocketCommands(conn, username, sessionID, ctxVal, cancel)
} 

/*
get user archived pending transactions information
requires user authentication from middleware
user/
*/
func (m *Manager) StreamUserArchivePendingTransactions(w http.ResponseWriter, r *http.Request) {
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

	/* sending initial session data */
	if err := m.sendCurrentArchivedPendingTransactions(conn, username, 1, 10); err != nil {
		m.errCh <- fmt.Errorf("error sending initial session: %w", err)
		return
	}

	/* 
		don't send updated stream of data
		archive is updated when session expires
		so user will have no access to the dashboard and can't see new
	*/

	/* specify the handler context */
	ctxVal := context.WithValue(ctx, "type", CtxStreamUserArchivePendingTransactions)

	/* handle web socket instructions from client */
	m.handleWebSocketCommands(conn, username, sessionID, ctxVal, cancel)
} 
