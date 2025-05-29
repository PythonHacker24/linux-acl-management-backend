package scheduler

import "context"

/* schedular interface */
type Scheduler interface {
	Run(ctx context.Context) error
}
