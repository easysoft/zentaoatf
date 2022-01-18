package controller

import (
	"encoding/json"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/config"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/exec"
	"github.com/kataras/iris/v12/websocket"
	"strings"
)

const (
	result = "result"
	outPut = "output"
)

var (
	ch chan int
)

type WsCtrl struct {
	Namespace         string
	*websocket.NSConn `stateless:"true"`

	WebSocketService *service.WebSocketService `inject:""`
}

func NewWsCtrl() *WsCtrl {
	inst := &WsCtrl{Namespace: serverConfig.WsDefaultNameSpace}
	return inst
}

func (c *WsCtrl) OnNamespaceConnected(msg websocket.Message) error {
	c.WebSocketService.SetConn(c.Conn)

	logUtils.Infof("WebSocket OnNamespaceConnected: ConnID=%s, Room=%s", c.Conn.ID(), msg.Room)

	data := map[string]string{"msg": "from server: connected to websocket"}
	c.WebSocketService.Broadcast(msg.Namespace, "", "OnVisit", data)
	return nil
}

// OnNamespaceDisconnect
// This will call the "OnVisit" event on all clients, except the current one,
// it can't because it's left but for any case use this type of design
func (c *WsCtrl) OnNamespaceDisconnect(msg websocket.Message) error {
	logUtils.Infof("WebSocket OnNamespaceDisconnect: ConnID=%s", c.Conn.ID())

	data := map[string]string{"msg": "from server: disconnected to websocket"}
	c.WebSocketService.Broadcast(msg.Namespace, "", "OnVisit", data)
	return nil
}

// OnChat This will call the "OnVisit" event on all clients, including the current one, with the 'newCount' variable.
func (c *WsCtrl) OnChat(msg websocket.Message) (err error) {
	ctx := websocket.GetContext(c.Conn)
	logUtils.Infof("WebSocket OnChat: remote address=%s, room=%s, msg=%s", ctx.RemoteAddr(), msg.Room, string(msg.Body))

	if ch != nil {
		ch <- 1
		ch = nil
		c.SendMsgByKey(result, "try to stop previous request...", msg)
	} else {
		ch = make(chan int)

		req := serverDomain.WsMsg{}
		err = json.Unmarshal(msg.Body, &req)
		if err != nil {
			logUtils.Errorf("参数验证失败", err.Error())
			resp := websocket.Message{Body: []byte(err.Error())}
			c.SendMsgByKey(result, "please wait previous request is in process", resp)
			return
		}

		//go shellUtils.ExeShellCallback(ch, "/Users/aaron/work/testing/res/loop.sh", "", c.SendMsg, msg)
		go scriptUtils.Exec(ch, c.SendExecMsg, req, msg)
	}

	return
}

func (c *WsCtrl) SendExecMsg(info string, msg websocket.Message) {
	logUtils.Infof("WebSocket SendMsg: room=%s, info=%s, msg=%s", msg.Room, info, string(msg.Body))
	info = strings.Trim(info, "\n")
	c.WebSocketService.Broadcast(msg.Namespace, msg.Room, msg.Event, info)
}

func (c *WsCtrl) SendMsg(value string, msg websocket.Message) {
	c.SendMsgByKey("", value, msg)
}
func (c *WsCtrl) SendMsgByKey(key, value string, msg websocket.Message) {
	if key == "" {
		key = outPut
	}
	data := map[string]string{key: value}

	logUtils.Infof("WebSocket SendMsg: room=%s, msg=%s", msg.Room, string(msg.Body))
	c.WebSocketService.Broadcast(msg.Namespace, msg.Room, msg.Event, data)
}
