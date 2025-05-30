package fcfs

import (
	"runtime"

	"github.com/PythonHacker24/linux-acl-management-backend/config"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/session"
)

/* spawns a new FCFS scheduler */
func NewFCFSScheduler(sm *session.Manager) *FCFSScheduler {
	/* calculate max workers */
	maxProcs := runtime.GOMAXPROCS(0)
	maxWorkers := config.BackendConfig.AppInfo.MaxWorkers 

	/* 
		incase of maxWorkers set less than or equal to 0, 
		use 75% of GOMAXPROCS to prevent starvation to other processes 
	*/
	if maxWorkers <= 0 {
		maxWorkers = int(float64(maxProcs) * 0.75)
	}

	/* Prevent over-allocation */
	if maxWorkers > maxProcs {
		maxWorkers = maxProcs
	}

	return &FCFSScheduler{
		curSessionManager: sm,
		maxWorkers: maxWorkers,
		semaphore: make(chan struct{}, maxWorkers),
	}
}
