package grpcpool

import (
	"google.golang.org/grpc"
	"sync"
)

/* gRPC connection pool for daemons */
type ClientPool struct {
	mu          sync.RWMutex
	conns       map[string]*grpc.ClientConn
	dialOptions []grpc.DialOption
	stopCh      chan struct{}
}
