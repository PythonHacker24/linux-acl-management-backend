package fcfs

import (
	"context"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/session"
	"go.uber.org/zap"
)

/* FCFS Scheduler attached with curSession.Manager */
type FCFSScheduler struct {
	curSessionManager *session.Manager 
	maxWorkers int

	/* for limiting spawning of goroutines */
	semaphore chan struct{}
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
			
			curSession := f.curSessionManager.GetNextSession()
			if curSession == nil {
				/* might need a delay of 10 ms */
				continue
			}
			
			/* check if transaction queue of the session is empty */
			curSession.Mutex.Lock()
			if curSession.TransactionQueue.Len() == 0 {
				curSession.Mutex.Unlock()
				continue
			}

			/* get a transaction from the session to process */
			transaction := curSession.TransactionQueue.Remove(curSession.TransactionQueue.Front())
			curSession.Mutex.Unlock()

			/* block if all workers are busy */
			f.semaphore <- struct{}{}

			/* go routine is available to be spawned */
			go func(curSession *session.Session, transaction interface{}) {
				/* defer clearing the semaphore channel */
				defer func() { <-f.semaphore }()

				/* process the transaction */
				if err := f.processTransaction; err != nil {
					zap.L().Error("Faild to process transaction", 
						zap.Error(err),
					)
				}
			}(curSession, transaction)
		}
	}	
}
