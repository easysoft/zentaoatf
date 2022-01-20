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
	"strconv"
	"strings"
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

	data := serverDomain.WsResp{Msg: "from server: connected to websocket"}
	c.WebSocketService.Broadcast(msg.Namespace, "", "OnVisit", data)
	return nil
}

// OnNamespaceDisconnect
// This will call the "OnVisit" event on all clients, except the current one,
// it can't because it's left but for any case use this type of design.
func (c *WsCtrl) OnNamespaceDisconnect(msg websocket.Message) error {
	logUtils.Infof("WebSocket OnNamespaceDisconnect: ConnID=%s", c.Conn.ID())

	data := serverDomain.WsResp{Msg: "from server: disconnected to websocket"}
	c.WebSocketService.Broadcast(msg.Namespace, "", "OnVisit", data)
	return nil
}

// OnChat This will call the "OnVisit" event on all clients,
// including the current one, with the 'newCount' variable.
func (c *WsCtrl) OnChat(wsMsg websocket.Message) (err error) {
	ctx := websocket.GetContext(c.Conn)
	logUtils.Infof("WebSocket OnChat: remote address=%s, room=%s, wsMsg=%s", ctx.RemoteAddr(), wsMsg.Room, string(wsMsg.Body))

	req := serverDomain.WsReq{}
	err = json.Unmarshal(wsMsg.Body, &req)
	if err != nil {
		msg := i118Utils.Sprintf("wrong_req_params", err.Error())
		data := serverDomain.WsResp{Msg: msg}
		c.SendExecMsg(data, wsMsg)
		logUtils.ExecConsole(color.FgRed, msg)
		return
	}

	act := req.Act

	if act == commConsts.ExecInit {
		msg := i118Utils.Sprintf("success_to_conn")
		data := serverDomain.WsResp{Msg: msg, IsRunning: strconv.FormatBool(scriptUtils.GetRunning())}
		c.SendExecMsg(data, wsMsg)
		logUtils.ExecConsole(color.FgCyan, msg)
		return
	}

	if act == commConsts.ExecStop {
		if !scriptUtils.GetRunning() {
			ch = nil
		} else {
			ch <- 1
			ch = nil
		}
		scriptUtils.SetRunning(false)

		msg := i118Utils.Sprintf("stopping_previous")
		data := serverDomain.WsResp{Msg: msg, IsRunning: "false"}
		c.SendExecMsg(data, wsMsg)
		logUtils.ExecConsole(color.FgCyan, msg)
		return
	}

	if (act == commConsts.ExecCase || act == commConsts.ExecModule ||
		act == commConsts.ExecSuite || act == commConsts.ExecTask) && scriptUtils.GetRunning() {
		msg := i118Utils.Sprintf("pls_stop_previous")
		data := serverDomain.WsResp{Msg: msg, IsRunning: "true"}
		c.SendExecMsg(data, wsMsg)
		logUtils.ExecConsole(color.FgRed, msg)

		return
	}

	ch = make(chan int, 1)
	//go shellUtils.ExeShellCallback(ch, "/Users/aaron/work/testing/res/loop.sh", "", c.SendExecMsg, wsMsg)
	go scriptUtils.Exec(ch, c.SendOutputMsg, req, wsMsg)
	scriptUtils.SetRunning(true)

	msg := i118Utils.Sprintf("start_to_run")
	data := serverDomain.WsResp{Msg: msg, IsRunning: "true"}
	c.SendExecMsg(data, wsMsg)
	logUtils.ExecConsole(color.FgCyan, msg)

	return
}

func (c *WsCtrl) SendOutputMsg(msg string, wsMsg websocket.Message) {
	logUtils.Infof("WebSocket SendExecMsg: room=%s, info=%s, msg=%s", wsMsg.Room, msg, string(wsMsg.Body))

	msg = strings.Trim(msg, "\n")
	data := serverDomain.WsResp{Msg: msg, Category: commConsts.Output}

	c.WebSocketService.Broadcast(wsMsg.Namespace, wsMsg.Room, wsMsg.Event, data)
}

func (c *WsCtrl) SendExecMsg(data serverDomain.WsResp, msg websocket.Message) {
	logUtils.Infof("WebSocket SendExecMsg: room=%s, msg=%s", msg.Room, string(msg.Body))

	data.Category = commConsts.Exec

	c.WebSocketService.Broadcast(msg.Namespace, msg.Room, msg.Event, data)
}
