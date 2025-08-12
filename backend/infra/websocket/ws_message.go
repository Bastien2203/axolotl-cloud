package websocket

import (
	"sync"
)

type WSMessageType string

const (
	SubscribeMessageType    WSMessageType = "subscribe"
	UnsubscribeMessageType  WSMessageType = "unsubscribe"
	JobLogUpdateMessageType WSMessageType = "job_log_update"
)

type WSMessage[T any] struct {
	Type WSMessageType `json:"type"`
	Data T             `json:"data"`
}

type WSMessageHandler struct {
	conn WebSocketConnection
}

var (
	topics      = make(map[string][]WebSocketConnection)
	topicsMutex = &sync.RWMutex{}
)

func NewWSMessageHandler(conn WebSocketConnection) *WSMessageHandler {
	return &WSMessageHandler{conn: conn}
}

func (h *WSMessageHandler) HandleMessage(message WSMessage[any]) {
	switch message.Type {
	case SubscribeMessageType:
		if data, ok := message.Data.(string); ok {
			h.Subscribe(data)
		}

	case UnsubscribeMessageType:
		if data, ok := message.Data.(string); ok {
			h.Unsubscribe(data)
		}
	}
}

func (h *WSMessageHandler) Subscribe(topicName string) {
	topicsMutex.Lock()
	defer topicsMutex.Unlock()

	if _, exists := topics[topicName]; !exists {
		topics[topicName] = []WebSocketConnection{}
	}
	topics[topicName] = append(topics[topicName], h.conn)
}

func (h *WSMessageHandler) Unsubscribe(topicName string) {
	topicsMutex.Lock()
	defer topicsMutex.Unlock()

	if _, exists := topics[topicName]; exists {
		for i, conn := range topics[topicName] {
			if conn == h.conn {
				topics[topicName] = append(topics[topicName][:i], topics[topicName][i+1:]...)
				break
			}
		}
		if len(topics[topicName]) == 0 {
			delete(topics, topicName)
		}
	}
}

func (h *WSMessageHandler) UnsubscribeAllTopics() {
	topicsMutex.Lock()
	defer topicsMutex.Unlock()

	for topicName := range topics {
		h.Unsubscribe(topicName)
	}
}

func UnsubscribeEveryoneFromTopic(topicName string) {
	topicsMutex.Lock()
	defer topicsMutex.Unlock()
	delete(topics, topicName)
}

func SendMessageToTopic[T any](topicName string, message WSMessage[T]) {
	topicsMutex.Lock()
	defer topicsMutex.Unlock()

	if conns, exists := topics[topicName]; exists {
		for _, conn := range conns {
			msgAny := WSMessage[any]{
				Type: message.Type,
				Data: message.Data,
			}
			if err := conn.Send(msgAny); err != nil {
				conn.Close()
			}
		}
	}
}
