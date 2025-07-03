package transprocessor

import (
	"context"

	"go.uber.org/zap"

	"github.com/PythonHacker24/linux-acl-management-backend/config"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/grpcpool"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/session"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/types"
)

/* instanciate new permission processor */
func NewPermProcessor(gRPCPool *grpcpool.ClientPool, errCh chan<-error) *PermProcessor {
	return &PermProcessor{
		gRPCPool: gRPCPool,
		errCh: errCh,
	}
}

/* processor for permissions manager */
func (p *PermProcessor) Process(ctx context.Context, curSession *session.Session, txn *types.Transaction) error {

	/* add complete information here + persistent logging in database */
	zap.L().Info("Processing Transaction",
		zap.String("user", curSession.Username),
	)

	select {
	case <-ctx.Done():
		/* close the processor */
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

		isRemote, host, port, found, absolutePath := FindServerFromPath(config.BackendConfig.FileSystemServers, txn.TargetPath)

		if !found {
			/* filepath is invalid, filesystem doesn't exist */
			txn.ErrorMsg = "filesystem of given path doesn't exist"
		} else {
			if isRemote {
				/* handle through daemons */
				p.HandleRemoteTransaction(host, port, txn, absolutePath)
			} else {
				/* handle locally */

				/* HandleLocalTransactions(txn) */
				p.HandleLocalTransaction(txn, absolutePath)
			}
		}

		zap.L().Info("Completed Transaction", 
			zap.String("ID", txn.ID.String()),
		)
	}

	return nil
}
