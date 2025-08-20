package transprocessor

import "github.com/PythonHacker24/linux-acl-management-backend/internal/grpcpool"

/*
	transprocessor implements the transactions structure that whole project complies with
*/

/* permissions processor */
type PermProcessor struct {
	gRPCPool *grpcpool.ClientPool
	errCh    chan<- error
}
