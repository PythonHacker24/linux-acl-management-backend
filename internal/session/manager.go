package session

import (
	"container/list"
	"sync"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/session/redis"
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
	redis 			redis.RedisClient
}

/* create a new session manager */
func NewManager(redis redis.RedisClient) *Manager {
	return &Manager{
		sessionsMap:  make(map[string]*Session),
		sessionOrder: list.New(),
		redis:	redis,
	}
}

/* get next session for round robin */
func (m *Manager) GetNextSession() *Session {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	/* check if sessionOrder is empty */
	if m.sessionOrder.Len() == 0 {
		return nil 
	}

	element := m.sessionOrder.Front()
	session := element.Value.(*Session)
	
	m.sessionOrder.MoveToBack(element)
	return session 
}
