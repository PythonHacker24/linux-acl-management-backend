package fcfs

import "github.com/PythonHacker24/linux-acl-management-backend/internal/session"

func NewFCFSScheduler(sm *session.Manager) *FCFSScheduler {
	return &FCFSScheduler{
		sessionManager: sm,
	}
}
