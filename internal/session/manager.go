package session

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/PythonHacker24/linux-acl-management-backend/config"
)

var (
	activeSessions = make(map[string]*Session)
	activeSessionsMutex sync.RWMutex
)

/* for creating a session for user */
func CreateSession(username string) error {

	/* lock the activeSessions mutex till the function ends */
	activeSessionsMutex.Lock()
	defer activeSessionsMutex.Unlock()

	/* check if session exists */
	if _, exists := activeSessions[username]; exists {
		return fmt.Errorf("user already exists in active sessions")
	}

	/* create the session */
	session := &Session{
		Username: username,
		Expiry:   time.Now().Add(time.Duration(config.BackendConfig.AppInfo.SessionTimeout) * time.Hour),
		Timer: time.AfterFunc(time.Duration(config.BackendConfig.AppInfo.SessionTimeout) * time.Hour,
			func() { ExpireSession(username) },
		),
		TransactionQueue:  list.New(),
		CurrentWorkingDir: config.BackendConfig.AppInfo.BasePath,
	}

	/* add session to active sessions */
	activeSessions[username] = session 

	return nil
}

/* for expiring a session */
func ExpireSession(username string) {
	activeSessionsMutex.Lock()
	defer activeSessionsMutex.Unlock()

	/* delete if user exists in active sessions */
	if _, exists := activeSessions[username]; exists {
		delete(activeSessions, username)
	}
}
