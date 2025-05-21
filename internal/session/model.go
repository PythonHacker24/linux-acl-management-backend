package session

import (
	"container/list"
	"sync"
	"time"
)

/* session struct for a user */
type Session struct {
    Username            string
    CurrentWorkingDir   string
    Expiry              time.Time
    Timer               *time.Timer
    TransactionQueue    *list.List
    Mutex               sync.Mutex
}
