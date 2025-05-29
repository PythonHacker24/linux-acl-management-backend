package session

import (
	"container/list"
	"sync"
	"time"
)

/* session struct for a user */
type Session struct {
	Username          string
	Expiry            time.Time
	Timer             *time.Timer
	TransactionQueue  *list.List
	Mutex             sync.Mutex

	/* 
		listElem stores it's node address in sessionOrder 
		this is done to maintain O(1) runtime performance while deleting session
	*/
	listElem 		  *list.Element
}
