package session

import (
	"container/list"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/PythonHacker24/linux-acl-management-backend/config"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/types"
)

/* for creating a session for user - used by HTTP HANDLERS */
func (m *Manager) CreateSession(username, ipAddress, userAgent string) error {

	/* lock the ActiveSessions mutex till the function ends */
	m.mutex.Lock()
	defer m.mutex.Unlock()

	/* check if session exists */
	if _, exists := m.sessionsMap[username]; exists {
		return fmt.Errorf("user already exists in active sessions")
	}

	/* Generate session metadata */
	sessionID := uuid.New()
	now := time.Now()

	/* create the session */
	session := &Session{
		ID:           sessionID,
		Status:       StatusActive,
		Username:     username,
		IP:           ipAddress,
		UserAgent:    userAgent,
		Expiry:       time.Now().Add(time.Duration(config.BackendConfig.AppInfo.SessionTimeout) * time.Hour),
		CreatedAt:    now,
		LastActiveAt: now,
		Timer: time.AfterFunc(time.Duration(config.BackendConfig.AppInfo.SessionTimeout)*time.Hour,
			func() { m.ExpireSession(username) },
		),
		CompletedCount:   0,
		FailedCount:      0,
		TransactionQueue: list.New(),
	}

	/* add session to active sessions map and list */
	element := m.sessionOrder.PushBack(session)
	session.listElem = element

	/* store session into the manager */
	m.sessionsMap[username] = session

	/* store session to Redis */
	m.saveSessionRedis(session)

	return nil
}

/* for expiring a session */
func (m *Manager) ExpireSession(username string) {
	/* thread safety for the manager */
	m.mutex.Lock()
	defer m.mutex.Unlock()

	/* check if user exists in active sessions */
	session, ok := m.sessionsMap[username]
	if !ok {
		return
	}

	session.Mutex.Lock()
	defer session.Mutex.Unlock()

	/*
		delete the session from Redis
		check if any transactions are remaining in the queue
		if yes, label transactions and sessions pending
		if no, label session expired
		push session and transactions to archive
	*/

	/* check if transactions are remaining in the session queue */
	if session.TransactionQueue.Len() != 0 {
		/* transactions are pending, mark them pending */
		for node := session.TransactionQueue.Front(); node != nil; node = node.Next() {
			/* work on transaction structure for *list.List() */
			txResult, ok := node.Value.(*types.Transaction)
			if !ok {
				continue
			}
			txResult.Status = types.StatusPending

			txnPQ, err := ConvertTransactiontoStoreParams(*txResult)
			if err != nil {
				/* error is conversion, continue the loop in good faith */
				/* need to handle these errors later */
				fmt.Printf("Failed to convert transaction to archive format: %v\n", err)
				continue
			}

			/* store transaction in PostgreSQL */
			if _, err := m.archivalPQ.CreateTransactionPQ(context.Background(), txnPQ); err != nil {
				/* log error but continue processing other transactions */
				fmt.Printf("Failed to archive transaction %s: %v\n", txResult.ID, err)
				continue
			}
		}

		/* mark session as pending */
		session.Status = StatusPending
	} else {
		/* empty transactions queue; mark the session as expired */
		session.Status = StatusExpired
	}

	/* remove session from sessionOrder Linked List */
	if session.listElem != nil {
		m.sessionOrder.Remove(session.listElem)
	}

	/* convert all session parameters to PostgreSQL compatible parameters */
	archive, err := ConvertSessionToStoreParams(session)
	if err != nil {
		/* session conversion failed, leave it in good faith */
		/* handle err later */
		fmt.Printf("Failed to convert session to archive format: %v\n", err)
		return
	}

	/* store session to the archive */
	if _, err := m.archivalPQ.StoreSessionPQ(context.Background(), *archive); err != nil {
		fmt.Printf("Failed to archive session: %v\n", err)
		return
	}

	/* delete both session and transaction results from Redis */
	sessionKey := fmt.Sprintf("session:%s", session.ID)
	txResultsKey := fmt.Sprintf("session:%s:txresults", session.ID)
	result := m.redis.Del(context.Background(), sessionKey, txResultsKey)
	if result.Err() != nil {
		fmt.Printf("Failed to delete session from Redis: %v\n", result.Err())
	} else {
		// Log the number of keys deleted
		deleted, _ := result.Result()
		fmt.Printf("Successfully deleted %d keys from Redis for session %s\n", deleted, session.ID)
	}

	/* remove session from sessionsMap */
	delete(m.sessionsMap, username)
}

/* add transaction to a session - assumes caller holds necessary locks */
func (m *Manager) AddTransaction(session *Session, txn interface{}) error {
	/* push transaction into the queue from back */
	session.TransactionQueue.PushBack(txn)

	/* convert transaction to correct type and save to Redis */
	if tx, ok := txn.(*types.Transaction); ok {
		if err := m.saveTransactionResultsRedis(session, *tx); err != nil {
			return fmt.Errorf("failed to save transaction to Redis: %w", err)
		}
	} else {
		return fmt.Errorf("invalid transaction type: expected *types.Transaction")
	}

	return nil
}

/* refresh the session timer */
func (m *Manager) refreshTimer(username string) error {
	/* thread safety for the manager */
	m.mutex.Lock()
	defer m.mutex.Unlock()

	/* get session from sessionMap */
	session, exists := m.sessionsMap[username]
	if !exists {
		return fmt.Errorf("session not found")
	}

	/* thread safety for the session */
	session.Mutex.Lock()
	defer session.Mutex.Unlock()

	/* reset the expiry time and last active time */
	session.Expiry = time.Now().Add(time.Duration(config.BackendConfig.AppInfo.SessionTimeout) * time.Hour)
	session.LastActiveAt = time.Now()

	/* stop the session timer */
	if session.Timer != nil {
		session.Timer.Stop()
	}

	/* reset the session timer */
	session.Timer = time.AfterFunc(time.Duration(config.BackendConfig.AppInfo.SessionTimeout)*time.Hour,
		func() { m.ExpireSession(username) },
	)

	/* update Redis for session */

	return nil
}

/* TODO: toDashoardView must be changed to fetch data from Redis only */

/* convert session information into frontend safe structure */
func (m *Manager) toDashboardView(username string) (SessionView, error) {
	/* thread safety for the manager */
	m.mutex.Lock()
	defer m.mutex.Unlock()

	/* get session from sessionMap */
	session, exists := m.sessionsMap[username]
	if !exists {
		return SessionView{}, fmt.Errorf("session not found")
	}

	/* thread safety for the session */
	session.Mutex.Lock()
	defer session.Mutex.Unlock()

	/* can be directly served as JSON in handler */
	return SessionView{
		ID:             session.ID.String(),
		Username:       session.Username,
		IP:             session.IP,
		UserAgent:      session.UserAgent,
		CreatedAt:      session.CreatedAt,
		LastActiveAt:   session.LastActiveAt,
		Expiry:         session.Expiry,
		CompletedCount: session.CompletedCount,
		FailedCount:    session.FailedCount,
		PendingCount:   session.TransactionQueue.Len(),
	}, nil
}
