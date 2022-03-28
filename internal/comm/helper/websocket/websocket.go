package websocketHelper

import (
	"encoding/json"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
	"strings"
)

var (
	wsConn *neffos.Conn
)

func SendOutputMsg(msg, isRunning string, wsMsg websocket.Message) {
	logUtils.Infof(i118Utils.Sprintf("ws_send_exec_msg", wsMsg.Room, msg, string(wsMsg.Body)))

	msg = strings.Trim(msg, "\n")
	data := serverDomain.WsResp{Msg: msg, Category: commConsts.Output}

	Broadcast(wsMsg.Namespace, wsMsg.Room, wsMsg.Event, data)
}

func SendExecMsg(msg, isRunning string, category commConsts.WsMsgCategory, wsMsg websocket.Message) {
	logUtils.Infof("WebSocket SendExecMsg: room=%s, msg=%s", wsMsg.Room, string(wsMsg.Body))

	data := serverDomain.WsResp{Msg: msg, IsRunning: isRunning, Category: category}

	Broadcast(wsMsg.Namespace, wsMsg.Room, wsMsg.Event, data)
}

func Broadcast(namespace, room, event string, data interface{}) {
	bytes, _ := json.Marshal(data)

	wsConn.Server().Broadcast(nil, websocket.Message{
		Namespace: namespace,
		Room:      room,
		Event:     event,
		Body:      bytes,
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
