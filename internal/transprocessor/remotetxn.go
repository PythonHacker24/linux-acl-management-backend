package transprocessor

import (
	"context"
	"fmt"
	"time"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/types"
	protos "github.com/PythonHacker24/linux-acl-management-backend/proto"
)

/* NEED TO ADD DURATION TIME -> EXISTS IN LOCAL */

/* takes a transactions and attempts to execute it via daemons */
func (p *PermProcessor) HandleRemoteTransaction(host string, port int, txn *types.Transaction, absolutePath string) error {

	/* if gRPCPool is nil, return an error */
	if p.gRPCPool == nil { 
		return fmt.Errorf("gRPC pool is nil") 
	}

	/* get connection to the respective daemon */
	address := fmt.Sprintf("%s:%d", host, port) 
	conn, err := p.gRPCPool.GetConn(address, p.errCh)
	if err != nil {
		p.errCh <- err
		return fmt.Errorf("failed to connect with a daemon: %s", address)
	}

	/* make it a for loop for interating all entries */
	aclpayload := &protos.ACLEntry{
		EntityType: txn.Entries.EntityType,
		Entity: txn.Entries.Entity,
		Permissions: txn.Entries.Permissions,
		Action: txn.Entries.Action,
		IsDefault: txn.Entries.IsDefault,
	}

	/* build the request for daemon */
	request := &protos.ApplyACLRequest{
		TransactionID: txn.ID.String(),
		TargetPath: absolutePath,
		Entry: aclpayload,
	}

	/* MAKE IT CONFIGURABLE */
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Minute)	

	aclClient := protos.NewACLServiceClient(conn)
	aclResponse, err := aclClient.ApplyACLEntry(ctx, request)
	if err != nil || aclResponse == nil  {
		p.errCh <- fmt.Errorf("failed to send ACL request to daemon")
		cancel()
		return err
	}

	if aclResponse.Success {

		/* 
			this is a bit crude for now, let daemon set this 
			backend should not have control over execution
		*/

		/* set transaction successful*/
		txn.Output = "ACL executed successfully on filesystem servers"
		
		txn.ExecStatus = true
	} else {
		txn.ErrorMsg = "ACL failed to get executed in the filesystem server"
	}

	cancel()
	return nil
}
