package session

import (
	"container/list"
	"time"

	"github.com/PythonHacker24/linux-acl-management-backend/config"
)

/* for creating a session for user */
func CreateSession(username string) {
	session := Session{
		Username: username,
		Expiry:   time.Now().Add(time.Duration(config.BackendConfig.AppInfo.SessionTimeout) * time.Hour),
		Timer: time.AfterFunc(time.Duration(config.BackendConfig.AppInfo.SessionTimeout)*time.Hour,
			func() { ExpireSession(username) },
		),
		TransactionQueue:  list.New(),
		CurrentWorkingDir: ".",
	}

}

/* for expiring a session */
func ExpireSession(username string) {

}
