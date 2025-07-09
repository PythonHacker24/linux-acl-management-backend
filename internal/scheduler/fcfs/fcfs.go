package fcfs

import (
	"context"
	"runtime"

	"go.uber.org/zap"

	"github.com/PythonHacker24/linux-acl-management-backend/config"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/session"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/transprocessor"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/types"
)

/* spawns a new FCFS scheduler */
func NewFCFSScheduler(sm *session.Manager, processor transprocessor.TransactionProcessor) *FCFSScheduler {
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
		maxWorkers:        maxWorkers,
		semaphore:         make(chan struct{}, maxWorkers),
		processor:         processor,
	}
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
			/* RULE: ctx is propagated all over the coming functions */

			/* get next session in the queue (round robin manner) */
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
			transaction := curSession.TransactionQueue.Remove(curSession.TransactionQueue.Front()).(*types.Transaction)
			curSession.Mutex.Unlock()

			/* block if all workers are busy */
			f.semaphore <- struct{}{}

			/* go routine is available to be spawned */
			go func(curSession *session.Session, transaction *types.Transaction) {
				/* defer clearing the semaphore channel */
				defer func() { <-f.semaphore }()

				/*
					process the transaction
					* processTransaction handles transaction processing completely
					* now it is responsible now responsible to execute it
					* role of scheduler in handling transactions ends here
				*/
				if err := f.processor.Process(ctx, curSession, transaction); err != nil {
					zap.L().Error("Failed to process transaction",
						zap.Error(err),
					)
				}

				/* we assume the transaction has been processed -> updated Redis */
				transaction.Status = types.StatusSuccess

				/* this whole code snippet should be called "Update Session State after transaction execution" */

				/* update the session's completed/failed count */
				curSession.Mutex.Lock()
				if transaction.ExecStatus {
					curSession.CompletedCount++
					f.curSessionManager.IncrementSessionCompletedRedis(curSession)
				} else {
					curSession.FailedCount++
					f.curSessionManager.IncrementSessionFailedRedis(curSession)
				}
				curSession.Mutex.Unlock()

				/* store the result of processed transaction into Redis */
				f.curSessionManager.SaveTransactionRedisList(curSession, transaction, "txresults")

				/* remove the transaction as pending from Redis */
				f.curSessionManager.RemovePendingTransaction(curSession, transaction.ID)

			}(curSession, transaction)
		}
	}
}
