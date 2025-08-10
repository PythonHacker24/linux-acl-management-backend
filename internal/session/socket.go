package session

import (
	"context"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

/* handle websocket commands from clients */
func (m *Manager) handleWebSocketCommands(conn *websocket.Conn, username, sessionID string, ctxVal context.Context, cancel context.CancelFunc) {
	defer cancel()

	/* infinite loop */
	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				m.errCh <- fmt.Errorf("websocket error: %w", err)
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
					m.errCh <- fmt.Errorf("failed to send pong: %w", err)
					return
				}

			/* refresh content served */
			case "refresh":
				/* client requests fresh data - implement based on current context */
				val := ctxVal.Value("type")

				switch val {
				case CtxStreamUserSession:
					/* push user session */
					if err := m.sendCurrentSession(conn, sessionID); err != nil {
						m.errCh <- fmt.Errorf("failed to send current session on command: %w", err)
					}
				case CtxStreamUserTransactionsResults:
					/* push user transactions results */
					if err := m.sendCurrentUserTransactionsResults(conn, sessionID, 100); err != nil {
						m.errCh <- fmt.Errorf("failed to send current list of results transactions on command: %w", err)
					}
				case CtxStreamUserTransactionsPending:
					/* push user transactions pending */
					if err := m.sendCurrentUserTransactionsPending(conn, sessionID, 100); err != nil {
						m.errCh <- fmt.Errorf("failed to send current list of results transactions on command: %w", err)
					}
				}
			}
		}
	}
}
