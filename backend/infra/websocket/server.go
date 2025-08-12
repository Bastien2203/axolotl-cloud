package websocket

import (
	"axolotl-cloud/infra/logger"
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	Upgrader     *websocket.Upgrader
	OnConnect    func(conn WebSocketConnection)
	OnMessage    func(conn WebSocketConnection, data WSMessage[any])
	OnDisconnect func(conn WebSocketConnection, err error)
}

func (s *WebSocketServer) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := NewGorillaConnection(w, r)
	if err != nil {
		http.Error(w, "WebSocket upgrade failed", http.StatusBadRequest)
		return
	}

	// Notify new connection
	if s.OnConnect != nil {
		s.OnConnect(conn)
	}

	// Handle disconnection
	go func() {
		<-conn.Done()
		if s.OnDisconnect != nil {
			s.OnDisconnect(conn, err)
		}
	}()

	// Process incoming messages
	go func() {
		for {
			select {
			case <-conn.Done():
				return
			default:
				_, data, err := conn.conn.ReadMessage()
				if err != nil {
					conn.Close()
					return
				}
				if s.OnMessage != nil {
					logger.Debug("Received message: %s", string(data))
					var msg WSMessage[any]
					if err := json.Unmarshal(data, &msg); err != nil {
						conn.Close()
						return
					}
					s.OnMessage(conn, msg)
				}
			}
		}
	}()
}
