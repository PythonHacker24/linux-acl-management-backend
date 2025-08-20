package grpcpool

import (
	"context"
	"fmt"
	"time"

	pb "github.com/PythonHacker24/linux-acl-management-backend/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

/* monitor gRPC connections */
func (p *ClientPool) MonitorHealth(addr string, conn *grpc.ClientConn, errCh chan<- error) {
	/* TODO: make it configurable */
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-p.stopCh:
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			pingClient := pb.NewPingServiceClient(conn)
			_, err := pingClient.Ping(ctx, &pb.PingRequest{})
			cancel()

			if err != nil {
				errCh <- fmt.Errorf("ping failed for daemon at %s: %w", addr, err)

				p.mu.Lock()
				conn.Close()
				delete(p.conns, addr)
				p.mu.Unlock()

				return
			} else {
				zap.L().Info("Ping success",
					zap.String("Address", addr),
				)
			}
		}
	}
}
