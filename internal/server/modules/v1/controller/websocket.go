package controller

import (
	"encoding/json"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/config"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/exec"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12/websocket"
	"strings"
)

const ()

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
func (c *WsCtrl) OnChat(wsMsg websocket.Message) (err error) {
	ctx := websocket.GetContext(c.Conn)
	logUtils.Infof("WebSocket OnChat: remote address=%s, room=%s, wsMsg=%s", ctx.RemoteAddr(), wsMsg.Room, string(wsMsg.Body))

	req := serverDomain.WsMsg{}
	err = json.Unmarshal(wsMsg.Body, &req)
	if err != nil {
		msg := i118Utils.Sprintf("wrong_req_params", err.Error())
		data := map[string]interface{}{"msg": msg, "isRunning": false}
		c.SendMsg(data, wsMsg)
		logUtils.ExecConsole(color.FgRed, msg)
		return
	}

	act := req.Act

	if act == commConsts.ExecInit {
		msg := i118Utils.Sprintf("success_to_conn")
		data := map[string]interface{}{"msg": msg, "isRunning": scriptUtils.IsRunning}
		c.SendMsg(data, wsMsg)
		logUtils.ExecConsole(color.FgCyan, msg)
		return
	}

	if act == commConsts.ExecStop {
		if !scriptUtils.IsRunning {
			ch = nil
		} else {
			ch <- 1
			ch = nil
		}

		msg := i118Utils.Sprintf("stopping_previous")
		data := map[string]interface{}{"msg": msg, "isRunning": false}
		c.SendMsg(data, wsMsg)
		logUtils.ExecConsole(color.FgCyan, msg)
		return
	}

	if act == commConsts.ExecCase && scriptUtils.IsRunning {
		msg := i118Utils.Sprintf("pls_stop_previous")
		data := map[string]interface{}{"msg": msg, "isRunning": true}
		c.SendMsg(data, wsMsg)
		logUtils.ExecConsole(color.FgRed, msg)

		return
	}

	ch = make(chan int)
	//go shellUtils.ExeShellCallback(ch, "/Users/aaron/work/testing/res/loop.sh", "", c.SendMsg, wsMsg)
	go scriptUtils.Exec(ch, c.SendExecMsg, req, wsMsg)
	scriptUtils.IsRunning = true

	msg := i118Utils.Sprintf("start_to_run")
	data := map[string]interface{}{"msg": msg, "isRunning": true}
	c.SendMsg(data, wsMsg)
	logUtils.ExecConsole(color.FgCyan, msg)

	return
}

func (c *WsCtrl) SendExecMsg(info string, msg websocket.Message) {
	logUtils.Infof("WebSocket SendMsg: room=%s, info=%s, msg=%s", msg.Room, info, string(msg.Body))
	info = strings.Trim(info, "\n")
	data := map[string]interface{}{"msg": info}
	c.WebSocketService.Broadcast(msg.Namespace, msg.Room, msg.Event, data)
}

func (c *WsCtrl) SendMsg(data map[string]interface{}, msg websocket.Message) {
	logUtils.Infof("WebSocket SendMsg: room=%s, msg=%s", msg.Room, string(msg.Body))
	c.WebSocketService.Broadcast(msg.Namespace, msg.Room, msg.Event, data)
}
