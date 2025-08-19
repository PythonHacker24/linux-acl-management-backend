package transprocessor

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/grpcpool"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/session"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/types"
)

/* instanciate new permission processor */
func NewPermProcessor(gRPCPool *grpcpool.ClientPool, errCh chan<- error) *PermProcessor {
	return &PermProcessor{
		gRPCPool: gRPCPool,
		errCh:    errCh,
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

		/* this line decides between systems like BeeGFS and NFS due to difference in ACL execution */
		isRemote, host, port, found, absolutePath := FindServerFromPath(txn.TargetPath)

		zap.L().Info("Found server",
			zap.String("targetPath", txn.TargetPath),
			zap.String("isRemote", fmt.Sprintf("%t", isRemote)),
			zap.String("host", host),
			zap.Int("port", port),
			zap.String("found", fmt.Sprintf("%t", found)),
			zap.String("absolutePath", absolutePath),
		)

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

		/* REMOVE THIS */
		zap.L().Info("Completed Transaction",
			zap.String("ID", txn.ID.String()),
		)
	}

	return nil
}
