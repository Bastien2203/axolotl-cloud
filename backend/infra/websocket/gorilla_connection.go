package websocket

import (
	"axolotl-cloud/infra/logger"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type GorillaConnection struct {
	conn *websocket.Conn
	send chan message
	done chan struct{}
	once sync.Once
}

type message struct {
	data WSMessage[any]
}

func NewGorillaConnection(w http.ResponseWriter, r *http.Request) (*GorillaConnection, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	conn := &GorillaConnection{
		conn: ws,
		send: make(chan message, 256),
		done: make(chan struct{}),
	}

	go conn.writePump()

	return conn, nil
}

func (c *GorillaConnection) Send(data WSMessage[any]) error {
	msg := message{data: data}
	select {
	case c.send <- msg:
		return nil
	case <-c.done:
		return errors.New("connection closed")
	}
}

func (c *GorillaConnection) Close() error {
	c.once.Do(func() {
		close(c.done)
		c.conn.Close()
	})
	return nil
}

func (c *GorillaConnection) Done() <-chan struct{} {
	return c.done
}

func (c *GorillaConnection) writePump() {
	defer c.Close()

	for {
		select {
		case msg := <-c.send:
			bytes, err := json.Marshal(msg.data)
			if err != nil {
				return
			}
			logger.Debug("Sending message: %s", string(bytes))
			if err := c.conn.WriteMessage(1, bytes); err != nil {
				return
			}
		case <-c.done:
			return
		}
	}
}
