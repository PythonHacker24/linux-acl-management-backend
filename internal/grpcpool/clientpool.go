package grpcpool

import (
	"fmt"

	"google.golang.org/grpc"
)

/* creates a new client pool */
func NewClientPool(opts ...grpc.DialOption) *ClientPool {
    return &ClientPool{
        conns:       make(map[string]*grpc.ClientConn),
        dialOptions: opts,
		stopCh: 	 make(chan struct{}),
    }
}

/* 
	creates a connection to given server 
	allows multiple connections to be established to daemons for any transactions in execution
*/
func (p *ClientPool) GetConn(addr string, errCh chan<-error) (*grpc.ClientConn, error) {
	/* check if connection exists or not */
	p.mu.RLock()
	conn, exists := p.conns[addr]
	p.mu.RUnlock()

	/* return the connection if it exists */
	if exists {
		return conn, nil
	}

	/* so connection doesn't exist, create a new one */
	p.mu.Lock()
	defer p.mu.Unlock()

	/* double check again (might been created between if exists and this line) */
	conn, exists = p.conns[addr]
	if exists {
        return conn, nil
    }

	/* create a new client for gRPC server */
	newConn, err := grpc.NewClient(addr, p.dialOptions...)
    if err != nil {
		return nil, fmt.Errorf("failed to add new connection: %w", err)
    }

	/* add connection to the pool */
    p.conns[addr] = newConn

	/*
		in case of connection issues, it will remove itself from connection pool
		when connection is demanded again, whole logic written above will be executed again 
	*/
	go p.MonitorHealth(addr, newConn, errCh)
	
	/* return connection */
    return newConn, nil
}

/* 
	close all connections in the pool
	call this while error channel exists
*/
func (p *ClientPool) CloseAll(errCh chan<-error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	/* iterate over all the connections and attempt to close them all */
	for _, conn := range p.conns {
		if err := conn.Close(); err != nil {
			errCh <- fmt.Errorf("error while closing gRPC connection: %w", err)
		}
	}

	p.conns = make(map[string]*grpc.ClientConn)
}
