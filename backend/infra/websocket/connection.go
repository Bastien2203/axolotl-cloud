package websocket

type WebSocketConnection interface {
	Send(data WSMessage[any]) error
	Close() error
	Done() <-chan struct{}
}
