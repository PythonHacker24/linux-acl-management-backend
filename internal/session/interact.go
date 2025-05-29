package session

import (
	"container/list"
)

/* create a new session manager */
func NewManager() *Manager {
	return &Manager{
		sessionsMap:  make(map[string]*Session),
		sessionOrder: list.New(),
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
