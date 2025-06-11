package session

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

/* handle websocket commands from clients */
func (m *Manager)handleWebSocketCommands(conn *websocket.Conn, ctxVal context.Context, cancel context.CancelFunc) {
	defer cancel()

	/* infinite loop */
	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				m.errCh<-fmt.Errorf("websocket error: %w", err)	
			}
			break
		}	

		/* handle commands from clients */
		if msgType, ok := msg["type"].(string); ok {
			switch msgType {
			
			/* ping echo test */
			case "ping":
				pongMsg := StreamMessage{
					Type:      "pong",
					Data:      "pong",
					Timestamp: time.Now(),
				}
				if err := conn.WriteJSON(pongMsg); err != nil {
					m.errCh<-fmt.Errorf("failed to send pong: %w", err)
					return
				}

			/* refresh content served */
			case "refresh":
				/* client requests fresh data - implement based on current context */
				val := ctxVal.Value("type")

				switch val {
				case StreamUserSession:
					/* push user session */
				case StreamUserTransactions:
					/* push user transactions */
				case StreamAllSessions:
					/* push all sessions */
				case StreamAllTransactions:
					/* push all transactions */
				}
			}
		}
	}
}
