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
func (m *Manager) CreateSession(username, ipAddress, userAgent string) (uuid.UUID, error) {

	/* lock the ActiveSessions mutex till the function ends */
	m.mutex.Lock()
	defer m.mutex.Unlock()

	/* check if session exists -> if yes, reset the timer and return the session ID */
	if session, exists := m.sessionsMap[username]; exists {
		if err := m.RefreshTimer(username); err != nil {
			m.errCh <- err
			return uuid.Nil, fmt.Errorf("sessions exists, but failed to refresh the timer")
		}
		return session.ID, nil
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
			func() {
				err := m.ExpireSession(username)
				if err != nil {
					m.errCh <- err
				}
			},
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
	if err := m.saveSessionRedis(session); err != nil {
		m.errCh <- err
		return uuid.Nil, fmt.Errorf("failed to store session to Redis")
	}

	return sessionID, nil
}

/* for expiring a session */
func (m *Manager) ExpireSession(username string) error {
	/* thread safety for the manager */
	m.mutex.Lock()
	defer m.mutex.Unlock()

	/* check if user exists in active sessions */
	session, ok := m.sessionsMap[username]
	if !ok {
		return fmt.Errorf("active user session not found")
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
			txResult, ok := node.Value.(*types.Transaction)
			if !ok {
				continue
			}

			/* make sure to set status to pending (shouldn't it be already set?) */
			txResult.Status = types.StatusPending

			/* convert transactions into PostgreSQL compatible parameters */
			txnPQ, err := ConvertTransactionPendingtoStoreParams(*txResult)
			if err != nil {
				m.errCh <- fmt.Errorf("failed to convert pending transaction to pending archive format: %w", err)
				continue
			}

			/* store transaction in PostgreSQL with retries */
			var storeErr error
			for retries := 0; retries < 3; retries++ {
				if _, err := m.archivalPQ.CreatePendingTransactionPQ(context.Background(), txnPQ); err != nil {
					storeErr = err
					time.Sleep(time.Second * time.Duration(retries+1))
					continue
				}
				storeErr = nil
				break
			}
			if storeErr != nil {
				m.errCh <- fmt.Errorf("failed to archive transaction %s after retries: %w", txResult.ID, storeErr)
				continue
			}
		}

		/* mark session as pending */
		session.Status = StatusPending
	} else {
		/* empty transactions queue; mark the session as expired */
		session.Status = StatusExpired
	}

	/* get transaction results from Redis */
	results, err := m.getTransactionResultsRedis(session, 10000)
	if err != nil {
		m.errCh <- fmt.Errorf("failed to get transaction results from Redis: %w", err)
	} else {
		for _, txResult := range results {
			if txResult.Status == types.StatusSuccess || txResult.Status == types.StatusFailed {
				pqParams, err := ConvertTransactionResulttoStoreParams(txResult)
				if err != nil {
					m.errCh <- fmt.Errorf("failed to convert transaction result to archive format: %w", err)
					continue
				}
				var storeErr error
				for retries := 0; retries < 3; retries++ {
					if _, err := m.archivalPQ.CreateResultsTransactionPQ(context.Background(), pqParams); err != nil {
						storeErr = err
						time.Sleep(time.Second * time.Duration(retries+1))
						continue
					}
					storeErr = nil
					break
				}
				if storeErr != nil {
					m.errCh <- fmt.Errorf("failed to archive transaction result %s after retries: %w", txResult.ID, storeErr)
					continue
				}
			}
		}
	}

	/* remove session from sessionOrder Linked List */
	if session.listElem != nil {
		m.sessionOrder.Remove(session.listElem)
	}

	/* convert all session parameters to PostgreSQL compatible parameters */
	archive, err := ConvertSessionToStoreParams(session)
	if err == nil {
		/* store session to the archive with retries */
		var storeErr error
		for retries := 0; retries < 3; retries++ {
			if _, err := m.archivalPQ.StoreSessionPQ(context.Background(), *archive); err != nil {
				storeErr = err
				time.Sleep(time.Second * time.Duration(retries+1))
				continue
			}
			storeErr = nil
			break
		}
		if storeErr != nil {
			m.errCh <- fmt.Errorf("failed to archive session after retries: %w", storeErr)
		}
	} else {
		/* handle err */
		m.errCh <- fmt.Errorf("failed to convert session to archive format: %w", err)
	}

	/* delete both session and transaction results from Redis */
	sessionKey := fmt.Sprintf("session:%s", session.ID)
	txResultsKey := fmt.Sprintf("session:%s:txresults", session.ID)
	result := m.redis.Del(context.Background(), sessionKey, txResultsKey)
	if result.Err() != nil {
		m.errCh <- fmt.Errorf("failed to delete session from Redis: %w", result.Err())
	}

	/* remove session from sessionsMap */
	delete(m.sessionsMap, username)

	return nil
}

/* add transaction to a session - assumes caller holds necessary locks */
func (m *Manager) AddTransaction(session *Session, txn *types.Transaction) error {
	/* push transaction into the queue from back */
	session.TransactionQueue.PushBack(txn)

	/* store transaction to Redis as a pending transaction */
	if err := m.SavePendingTransaction(session, txn); err != nil {
		return fmt.Errorf("failed to save transaction to Redis: %w", err)
	}

	return nil
}

/* refresh the session timer */
func (m *Manager) RefreshTimer(username string) error {
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
		func() {
			if err := m.ExpireSession(username); err != nil {
				m.errCh <- err
			}
		},
	)

	/* update Redis for session */
	if err := m.saveSessionRedis(session); err != nil {
		m.errCh <- err
		return fmt.Errorf("failed to store session to Redis")
	}

	return nil
}

/* check is a session exists for a username */
func (m *Manager) SessionExistance(username string) (uuid.UUID, bool, error) {
	/* thread safety for the manager */
	m.mutex.Lock()
	defer m.mutex.Unlock()

	/* get session from sessionMap */
	session, exists := m.sessionsMap[username]
	if exists {

		/* thread safety for the session */
		session.Mutex.Lock()
		defer session.Mutex.Unlock()

		if session.Username == username {
			return session.ID, true, nil
		}
	}

	return uuid.Nil, false, nil
}
