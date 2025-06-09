package scheduler

import "context"

/* scheduler interface */
type Scheduler interface {
	Run(ctx context.Context) error
}
