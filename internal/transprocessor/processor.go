package transprocessor

import (
	"context"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/session"
)

/*
	the job of scheduler was to handle how transactions are allocated to executors
	transprocessor's job is to open the content of transactions and take care of them onwards
	so work of transactions was irrelevant to scheduler - transprocessor is responsible

	also, this archirecture allows us to create mulitple processors in case we plan to extend in future
*/

/* transaction processor - pluggable to any scheduler */
type TransactionProcessor interface {
	Process(ctx context.Context, curSession *session.Session, transaction interface{}) error
}
