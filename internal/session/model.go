package session

import (
	"container/list"
	"sync"
	"time"

	"github.com/google/uuid"
)

/* defining Status type for sessions */
type Status string

/* for status field */
const (
	StatusActive  Status = "active"
	StatusExpired Status = "expired"
	StatusPending Status = "pending"
)

/*
session struct for a user
appropriate fields must always be updated when any request is made
*/
type Session struct {
	/* keep count of completed and failed transactions */
	CompletedCount int
	FailedCount    int

	/* session status: active: 1 / expired: 0 */
	Status Status

	/* unique ID of session [will be associated with the user forever in logs] */
	ID uuid.UUID

	/* username of the user */
	Username string

	/*
		IP and UserAgent for security logs
		also can be used for blacklisting and whitelistings
		illegal useragents can be caught as well as unauthorized IP addresses
	*/
	IP        string
	UserAgent string

	/* for logging user activity */
	CreatedAt    time.Time
	LastActiveAt time.Time
	Expiry       time.Time
	Timer        *time.Timer

	/* transactions issued by the user */
	TransactionQueue *list.List

	/*
		listElem stores it's node address in sessionOrder
		this is done to maintain O(1) runtime performance while deleting session
	*/
	listElem *list.Element

	/* mutex for thread safety */
	Mutex sync.Mutex
}

/* SessionStreamData is a frontend-safe representation of a session that goes through websocket */
type SessionStreamData struct {
	ID             string    `json:"id"`
	Username       string    `json:"username"`
	IP             string    `json:"ip"`
	UserAgent      string    `json:"userAgent"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
	LastActiveAt   time.Time `json:"lastActiveAt"`
	Expiry         time.Time `json:"expiry"`
	CompletedCount int       `json:"completed"`
	FailedCount    int       `json:"failed"`
}

/* websocket stream message */
type StreamMessage struct {
	Type      string    `json:"type"`
	Data      any       `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

/* TransactionStreamData is a frontend-safe representation of a transaction sent via websocket */
type TransactionStreamData struct {
	ID         string           `json:"id"`
	SessionID  string           `json:"sessionId"`
	Timestamp  time.Time        `json:"timestamp"`
	Operation  string           `json:"operation"`
	TargetPath string           `json:"targetPath"`
	Entries    []ACLEntryStream `json:"entries"`
	Status     string           `json:"status"`
	ErrorMsg   string           `json:"errorMsg,omitempty"`
	Output     string           `json:"output"`
	ExecutedBy string           `json:"executedBy"`
	DurationMs int64            `json:"durationMs"`
}

/* ACLEntryStream is a frontend-safe version of an individual ACL entry */
type ACLEntryStream struct {
	EntityType  string `json:"entityType"`
	Entity      string `json:"entity"`
	Permissions string `json:"permissions"`
	Action      string `json:"action"`
	Success     bool   `json:"success"`
	Error       string `json:"error,omitempty"`
}

/* archival data fetch requests */
type ArchivalRequest struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}
