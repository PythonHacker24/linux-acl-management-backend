package types

import (
    "time"

    "github.com/google/uuid"
)

/*
	contains shared definations where compete modulation was not possible
	Eg. session and transprocesser need same transaction structure and updating seperate definations
	needs rewriting same code multiple times.
*/

/* request body for scheduling transaction */
type ScheduleTransactionRequest struct {
	Operation  OperationType `json:"operation"`
	TargetPath string        `json:"targetPath"`
	Entries    []ACLEntry    `json:"entries"`
}

/* represents the result of the transaction */
type TxnStatus string

/* defining transactions status types */
const (
    StatusPending   TxnStatus = "pending"
    StatusSuccess   TxnStatus = "success"
    StatusFailed    TxnStatus = "failed"
)

/* represents what kind of ACL operation was performed */
type OperationType string

/* defining operating types */
const (
    OperationGetACL OperationType = "getfacl"
    OperationSetACL OperationType = "setfacl"
)

/* represents an individual ACL rule attempted to be changed */
type ACLEntry struct {
	/* e.g., "user", "group", "mask", "other" */
    EntityType string   `json:"entityType"`

	/* 	
		username, group name, or blank
		blank means it applies to the current owner/group (e.g., user::, group::, other::, mask::) 
	*/
    Entity     string   `json:"entity"`

	/* e.g., "rwx", "rw-", etc. */
    Permissions string  `json:"permissions"`

	/* e.g., "add", "modify", "remove" */
    Action      string  `json:"action"`

	/* whether this is a default ACL (i.e., applies to new files/subdirs) */
	IsDefault bool `json:"isDefault"`

	/* only set if failed */
	Error       string `json:"error,omitempty"` 
    Success     bool   `json:"success"`
}

/* holds the full state of a permission change operation */
type Transaction struct {
    ID          uuid.UUID         `json:"id"`
    SessionID   uuid.UUID         `json:"sessionId"`
    Timestamp   time.Time         `json:"timestamp"`

	/* getfacl/setfacl */
    Operation   OperationType     `json:"operation"` 

	/* File/directory affected */
    TargetPath  string            `json:"targetPath"` 

	/* ACL entries involved */
    Entries     []ACLEntry        `json:"entries"` 

	/* success/failure/pending */
    Status      TxnStatus		  `json:"status"` 

	/* set if failed */
    ErrorMsg    string            `json:"errorMsg,omitempty"` 

	/* stdout or stderr captured */
    Output      string            `json:"output"`

	/* user who triggered this */
    ExecutedBy  string            `json:"executedBy"` 

	/* execution duration in ms */
    DurationMs  int64             `json:"durationMs"`
}
