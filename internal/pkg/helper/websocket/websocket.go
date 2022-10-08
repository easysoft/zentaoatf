package websocketHelper

import (
	"encoding/json"
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"strings"
)

var (
	wsConn *neffos.Conn
)

func SendOutputMsg(msg, isRunning string, info iris.Map, wsMsg *websocket.Message) {
	logUtils.Infof(i118Utils.Sprintf("ws_send_exec_msg", wsMsg.Room,
		strings.ReplaceAll(strings.TrimSpace(msg), `%`, `%%`)))

	if wsMsg == nil {
		return
	}

	msg = strings.Trim(msg, "\n")
	resp := commDomain.WsResp{Msg: msg, Category: commConsts.Output, Info: info}

	bytes, _ := json.Marshal(resp)
	mqData := commDomain.MqMsg{Namespace: wsMsg.Namespace, Room: wsMsg.Room, Event: wsMsg.Event, Content: string(bytes)}
	PubMsg(mqData)
}

func SendExecMsg(msg, isRunning string, category commConsts.WsMsgCategory, info iris.Map, wsMsg *websocket.Message) {
	logUtils.Infof(i118Utils.Sprintf("ws_send_exec_msg", wsMsg.Room,
		strings.ReplaceAll(strings.TrimSpace(msg), `%`, `%%`)))

	if wsMsg == nil {
		return
	}

	msg = strings.TrimSpace(msg)
	resp := commDomain.WsResp{Msg: msg, IsRunning: isRunning, Category: category, Info: info}

	bytes, _ := json.Marshal(resp)
	mqData := commDomain.MqMsg{Namespace: wsMsg.Namespace, Room: wsMsg.Room, Event: wsMsg.Event, Content: string(bytes)}
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
