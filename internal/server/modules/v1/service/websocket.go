package service

import (
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
)

var (
	wsConn *neffos.Conn
)

type WebSocketService struct {
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{}
}

func (s *WebSocketService) Broadcast(namespace, room, event string, data interface{}) {
	bytes, _ := json.Marshal(data)

	wsConn.Server().Broadcast(nil, websocket.Message{
		Namespace: namespace,
		Room:      room,
		Event:     event,
		Body:      bytes,
	})
}

func (s *WebSocketService) SetConn(conn *neffos.Conn) {
	wsConn = conn
}

type PrefixedLogger struct {
	Prefix string
}

func (s *PrefixedLogger) Log(msg string) {
	fmt.Printf("%s: %s\n", s.Prefix, msg)
}
