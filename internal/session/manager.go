package session

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/PythonHacker24/linux-acl-management-backend/config"
)

/*
	session manager
	sessionsMap -> Maps of sessions -> for O(1) access | fast access during deletion
	sessionOrder -> LinkedList of sessions -> for round robin | fair scheduling
	sessionsMap and sessionOrder are always in sync
	both are kept at the same time due to various runtime performance requirements
	trading off space for runtime speed performance
*/
type Manager struct {
	sessionsMap		map[string]*Session
	sessionOrder	*list.List
	mutex 			sync.RWMutex
}

/* for creating a session for user - used by HTTP HANDLERS */
func (m *Manager) CreateSession(username string) error {

	/* lock the ActiveSessions mutex till the function ends */
	m.mutex.Lock()
	defer m.mutex.Unlock()

	/* check if session exists */
	if _, exists := m.sessionsMap[username]; exists {
		return fmt.Errorf("user already exists in active sessions")
	}

	/* create the session */
	session := &Session{
		Username: username,
		Expiry:   time.Now().Add(time.Duration(config.BackendConfig.AppInfo.SessionTimeout) * time.Hour),
		Timer: time.AfterFunc(time.Duration(config.BackendConfig.AppInfo.SessionTimeout)*time.Hour,
			func() { m.ExpireSession(username) },
		),
		TransactionQueue:  list.New(),
	}

	/* add session to active sessions map and list */
	element := 	m.sessionOrder.PushBack(session)
	session.listElem = element

	m.sessionsMap[username] = session

	return nil
}

/* for expiring a session */
func (m *Manager) ExpireSession(username string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	/* check if user exists in active sessions */
	session, ok := m.sessionsMap[username]
	if !ok {
		return
	}

	/* remove session from sessionOrder Linked List */
	if session.listElem != nil {
		m.sessionOrder.Remove(session.listElem)
	}

	/* remove session from sessionsMap */
	delete(m.sessionsMap, username)
}
