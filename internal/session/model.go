package session

import (
	"container/list"
	"sync"
	"time"
)

/* 
	session struct for a user 
	appropriate fields must always be updated when any request is made
*/
type Session struct {
	/* keep count of completed and failed transactions */
	CompletedCount    int
	FailedCount       int

	/* unique ID of session [will be associated with the user forever in logs] */
	ID 				  string

	/* username of the user */
	Username          string

	/* 
		IP and UserAgent for security logs
		also can be used for blacklisting and whitelistings 
		illegal useragents can be caught as well as unauthorized IP addresses
	*/
	IP 				  string
	UserAgent		  string

	/* for logging user activity */
	CreatedAt		  time.Time
	LastActiveAt      time.Time
	Expiry            time.Time
	Timer             *time.Timer

	/* transactions issued by the user */
	TransactionQueue  *list.List

	/* mutex for thread safety */
	Mutex             sync.Mutex

	/* 
		listElem stores it's node address in sessionOrder 
		this is done to maintain O(1) runtime performance while deleting session
	*/
	listElem 		  *list.Element
}
