package session

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/PythonHacker24/linux-acl-management-backend/config"
)

var (
	ActiveSessions      = make(map[string]*Session)
	ActiveSessionsMutex sync.RWMutex
)

/* for creating a session for user */
func CreateSession(username string) error {

	/* lock the ActiveSessions mutex till the function ends */
	ActiveSessionsMutex.Lock()
	defer ActiveSessionsMutex.Unlock()

	/* check if session exists */
	if _, exists := ActiveSessions[username]; exists {
		return fmt.Errorf("user already exists in active sessions")
	}

	/* create the session */
	session := &Session{
		Username: username,
		Expiry:   time.Now().Add(time.Duration(config.BackendConfig.AppInfo.SessionTimeout) * time.Hour),
		Timer: time.AfterFunc(time.Duration(config.BackendConfig.AppInfo.SessionTimeout)*time.Hour,
			func() { ExpireSession(username) },
		),
		TransactionQueue:  list.New(),
		CurrentWorkingDir: config.BackendConfig.AppInfo.BasePath,
	}

	/* add session to active sessions */
	ActiveSessions[username] = session

	return nil
}

/* for expiring a session */
func ExpireSession(username string) {
	ActiveSessionsMutex.Lock()
	defer ActiveSessionsMutex.Unlock()

	/* delete if user exists in active sessions */
	if _, exists := ActiveSessions[username]; exists {
		delete(ActiveSessions, username)
	}
}
