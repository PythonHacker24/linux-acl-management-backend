package scheduler

/*
	laclm uses FCFS scheduling algorithm
	the scheduler module is highly modular and can be attached/detached/replaced quickly on-demand
*/

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

/*
initialized a scheduler of Scheduler type as a goroutine with context
when ctx.Done() is recieved, scheduler starts shutting down
the main function waits till wg.Done() is not called ensuring complete shutdown of scheduler
in case of any error, the errCh is used to propogate it back to main function where it's handled
*/
func InitScheduler(ctx context.Context, sched Scheduler, wg *sync.WaitGroup, errCh chan<- error) {
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		zap.L().Info("Scheduler Initialization Started")

		/* the context is used here for gracefully stopping the scheduler */
		if err := sched.Run(ctx); err != nil {
			zap.L().Error("Scheduler running error",
				zap.Error(err),
			)
		} else {
			zap.L().Info("Scheduler Stopped Gracefully")
		}
	}(ctx)
}
