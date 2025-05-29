package fcfs

import (
	"context"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/session"
)

/* FCFS Scheduler attached with session.Manager */
type FCFSScheduler struct {
	sessionManager *session.Manager 
}

/* run the fcfs scheduler with context */
func (f *FCFSScheduler) Run(ctx context.Context) error {
	for {
		select {

		/* check if ctx is done - catchable if default is not working hard (ideal scheduler) */
		case <-ctx.Done():
			return nil

		/* in case default is working hard - ctx is passed here so it must attempt to quit */
		default:
			/* RULE: ctx is propogates all over the coming functions */
				
		}
	}	
}
