package transprocessor

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/session"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/types"
)

/* instanciate new permission processor */
func NewPermProcessor(errCh chan<-error) *PermProcessor {
	return &PermProcessor{
		errCh: errCh,
	}
}

/* processor for permissions manager */
func (p *PermProcessor) Process(ctx context.Context, curSession *session.Session, tx interface{}) error {
	transaction, ok := tx.(*types.Transaction)
	if !ok {
		return fmt.Errorf("invalid transaction type")
	}

	/* add complete information here + persistent logging in database */
	zap.L().Info("Processing Transaction",
		zap.String("user", curSession.Username),
	)

	select {
	case <-ctx.Done():
		/*
			store this into persistent storage too!
			make sure database connections are closed after scheduler shutsdown
		*/
		zap.L().Warn("Transaction process stopped due to shutdown",
			zap.String("user", curSession.Username),
		)
		return ctx.Err()
	default:
		/*
			permprocessor hands over transactions to remoteprocessor/localprocessor depending upon request
			remoteprocessor -> handles permissions on remote servers
			localprocessor -> handles permissions on local system (where this backend is deployed)
		*/
		_ = transaction

		/* for testing purposes only */
		time.Sleep(5 * time.Second)
		zap.L().Info("Completed Transaction", 
			zap.String("ID", transaction.ID.String()),
		)
	}

	return nil
}
