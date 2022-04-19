package websocketHelper

import (
	"encoding/json"
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	i118Utils "github.com/easysoft/zentaoatf/internal/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"strings"
)

var (
	wsConn *neffos.Conn
)

func SendOutputMsg(msg, isRunning string, wsMsg *websocket.Message) {
	logUtils.Infof(i118Utils.Sprintf("ws_send_exec_msg", wsMsg.Room,
		strings.ReplaceAll(strings.TrimSpace(msg), `%`, `%%`)))

	msg = strings.Trim(msg, "\n")
	resp := serverDomain.WsResp{Msg: msg, Category: commConsts.Output}

	bytes, _ := json.Marshal(resp)
	mqData := MqMsg{Namespace: wsMsg.Namespace, Room: wsMsg.Room, Event: wsMsg.Event, Content: string(bytes)}
	PubMsg(mqData)
}

func SendExecMsg(msg, isRunning string, category commConsts.WsMsgCategory, wsMsg *websocket.Message) {
	logUtils.Infof(i118Utils.Sprintf("ws_send_exec_msg", wsMsg.Room,
		strings.ReplaceAll(strings.TrimSpace(msg), `%`, `%%`)))

	msg = strings.TrimSpace(msg)
	resp := serverDomain.WsResp{Msg: msg, IsRunning: isRunning, Category: category}

	bytes, _ := json.Marshal(resp)
	mqData := MqMsg{Namespace: wsMsg.Namespace, Room: wsMsg.Room, Event: wsMsg.Event, Content: string(bytes)}
	PubMsg(mqData)
}

func Broadcast(namespace, room, event string, content string) {
	wsConn.Server().Broadcast(nil, websocket.Message{
		Namespace: namespace,
		Room:      room,
		Event:     event,
		Body:      []byte(content),
	})
}

func SetConn(conn *neffos.Conn) {
	wsConn = conn
}

type PrefixedLogger struct {
	Prefix string
}

func (s *PrefixedLogger) Log(msg string) {
	fmt.Printf("%s: %s\n", s.Prefix, msg)
}
