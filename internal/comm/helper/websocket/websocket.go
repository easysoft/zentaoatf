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
	"time"
)

var (
	wsConn *neffos.Conn
)

func SendOutputMsg(msg, isRunning string, wsMsg *websocket.Message) {
	logUtils.Infof(i118Utils.Sprintf("ws_send_exec_msg", wsMsg.Room,
		strings.ReplaceAll(strings.TrimSpace(msg), `%`, `%%`)))

	msg = strings.Trim(msg, "\n")
	data := serverDomain.WsResp{Msg: msg, Category: commConsts.Output}

	time.Sleep(time.Millisecond * 50)
	Broadcast(wsMsg.Namespace, wsMsg.Room, wsMsg.Event, data)
}

func SendExecMsg(msg, isRunning string, category commConsts.WsMsgCategory, wsMsg *websocket.Message) {
	logUtils.Infof(i118Utils.Sprintf("ws_send_exec_msg", wsMsg.Room,
		strings.ReplaceAll(strings.TrimSpace(msg), `%`, `%%`)))

	msg = strings.TrimSpace(msg)
	data := serverDomain.WsResp{Msg: msg, IsRunning: isRunning, Category: category}

	time.Sleep(time.Millisecond * 50)
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
