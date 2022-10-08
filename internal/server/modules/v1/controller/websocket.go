package controller

import (
	"encoding/json"
	"fmt"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	watchHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/watch"
	websocketHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/websocket"
	serverConfig "github.com/easysoft/zentaoatf/internal/server/config"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/service"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12/websocket"
)

type WebSocketCtrl struct {
	Namespace         string
	*websocket.NSConn `stateless:"true"`

	WorkspaceService *service.WorkspaceService `inject:""`
	TestExecService  *service.ExecService      `inject:""`
}

func NewWebSocketCtrl() *WebSocketCtrl {
	inst := &WebSocketCtrl{Namespace: serverConfig.WsDefaultNameSpace}
	return inst
}

func (c *WebSocketCtrl) OnNamespaceConnected(wsMsg websocket.Message) error {
	websocketHelper.SetConn(c.Conn)

	logUtils.Infof(i118Utils.Sprintf("ws_namespace_connected", c.Conn.ID(), wsMsg.Room))

	resp := commDomain.WsResp{Msg: "from server: connected to websocket"}
	bytes, _ := json.Marshal(resp)
	mqData := commDomain.MqMsg{Namespace: wsMsg.Namespace, Room: wsMsg.Room, Event: wsMsg.Event, Content: string(bytes)}
	websocketHelper.PubMsg(mqData)
	return nil
}

// OnNamespaceDisconnect
// This will call the "OnVisit" event on all clients, except the current one,
// it can't because it's left but for any case use this type of design.
func (c *WebSocketCtrl) OnNamespaceDisconnect(wsMsg websocket.Message) error {
	logUtils.Infof(i118Utils.Sprintf("ws_namespace_disconnected", c.Conn.ID()))

	resp := commDomain.WsResp{Msg: fmt.Sprintf("ws_connected")}
	bytes, _ := json.Marshal(resp)
	mqData := commDomain.MqMsg{Namespace: wsMsg.Namespace, Room: wsMsg.Room, Event: wsMsg.Event, Content: string(bytes)}
	websocketHelper.PubMsg(mqData)
	return nil
}

// OnChat This will call the "OnVisit" event on all clients,
// including the current one, with the 'newCount' variable.
func (c *WebSocketCtrl) OnChat(wsMsg websocket.Message) (err error) {
	ctx := websocket.GetContext(c.Conn)
	logUtils.Infof(i118Utils.Sprintf("ws_onchat", ctx.RemoteAddr(), wsMsg.Room, string(wsMsg.Body)))

	req := serverDomain.ExecReq{}
	err = json.Unmarshal(wsMsg.Body, &req)
	if err != nil {
		msg := i118Utils.Sprintf("wrong_req_params", err.Error())
		websocketHelper.SendExecMsg(msg, "", commConsts.Error, nil, &wsMsg)
		logUtils.ExecConsole(color.FgRed, msg)
		return
	}

	if req.Act == commConsts.ExecInit {
		msg := i118Utils.Sprintf("success_to_conn")
		//websocketHelper.SendExecMsg(msg, strconv.FormatBool(execHelper.GetRunning()), wsMsg)
		logUtils.ExecConsole(color.FgCyan, msg)

	} else if req.Act == commConsts.Watch {
		watchHelper.WatchFromReq(req.TestSets, &wsMsg)

	} else if req.Act == commConsts.ExecStop {
		c.TestExecService.Stop(&wsMsg)

	} else {
		c.TestExecService.Start(req, &wsMsg)

	}

	return
}
