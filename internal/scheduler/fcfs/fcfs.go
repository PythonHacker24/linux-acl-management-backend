package fcfs

import (
	"runtime"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/session"
)

/* spawns a new FCFS scheduler */
func NewFCFSScheduler(sm *session.Manager) *FCFSScheduler {
	/* calculate max workers */
	/* TODO: make it configurable and dynamic */
	maxWorkers := runtime.NumCPU()
	return &FCFSScheduler{
		sessionManager: sm,
		maxWorkers: maxWorkers,
		semaphore: make(chan struct{}, maxWorkers),
	}
}
