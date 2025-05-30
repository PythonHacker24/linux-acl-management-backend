package transprocessor

import (
	"context"
	"fmt"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/session"
	"go.uber.org/zap"
)

/* instanciate new permission processor */
func NewPermProcessor() *PermProcessor {
	return &PermProcessor{}	
}

/* processor for permissions manager */
func (p *PermProcessor) Process(ctx context.Context, curSession *session.Session, tx interface{}) error {
	transaction, ok := tx.(*Transaction) 
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
	}

	return nil	
}
