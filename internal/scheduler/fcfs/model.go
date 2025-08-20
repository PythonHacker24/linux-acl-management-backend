package fcfs

import (
	"github.com/PythonHacker24/linux-acl-management-backend/internal/session"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/transprocessor"
)

/*
	Notes: the structure of scheduler is very modular
	Docs must be updated for replacing a certain scheduler module with another
	This includes installation of prebuilt module or developing a module
*/

/* FCFS Scheduler attached with curSession.Manager */
type FCFSScheduler struct {
	curSessionManager *session.Manager
	maxWorkers        int

	/* for limiting spawning of goroutines */
	semaphore chan struct{}
	processor transprocessor.TransactionProcessor
}
