package controller

import (
	"encoding/json"
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/comm/domain"
	execHelper "github.com/easysoft/zentaoatf/internal/comm/helper/exec"
	websocketHelper "github.com/easysoft/zentaoatf/internal/comm/helper/websocket"
	i118Utils "github.com/easysoft/zentaoatf/internal/pkg/lib/i118"
	"github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	"github.com/easysoft/zentaoatf/internal/server/config"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/service"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
)

var (
	ch chan int
)

type WebSocketCtrl struct {
	Namespace         string
	WorkspaceService  *service.WorkspaceService `inject:""`
	*websocket.NSConn `stateless:"true"`
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

	req := serverDomain.WsReq{}
	err = json.Unmarshal(wsMsg.Body, &req)
	if err != nil {
		msg := i118Utils.Sprintf("wrong_req_params", err.Error())
		websocketHelper.SendExecMsg(msg, "", commConsts.Error, nil, &wsMsg)
		logUtils.ExecConsole(color.FgRed, msg)
		return
	}

	act := req.Act

	if act == commConsts.ExecInit {
		msg := i118Utils.Sprintf("success_to_conn")
		//websocketHelper.SendExecMsg(msg, strconv.FormatBool(execHelper.GetRunning()), wsMsg)
		logUtils.ExecConsole(color.FgCyan, msg)
		return
	}

	if act == commConsts.ExecStop {
		if ch != nil {
			if !execHelper.GetRunning() {
				ch = nil
			} else {
				ch <- 1
				ch = nil
			}
		}

		execHelper.SetRunning(false)

		msg := i118Utils.Sprintf("end_task")
		websocketHelper.SendExecMsg(msg, "false", commConsts.Run, nil, &wsMsg)
		logUtils.ExecConsole(color.FgCyan, msg)
		return
	}

	if execHelper.GetRunning() && (act == commConsts.ExecCase || act == commConsts.ExecModule ||
		act == commConsts.ExecSuite || act == commConsts.ExecTask || act == commConsts.ExecUnit) {
		msg := i118Utils.Sprintf("pls_stop_previous")
		websocketHelper.SendExecMsg(msg, "true", commConsts.Run, nil, &wsMsg)
		logUtils.ExecConsole(color.FgRed, msg)

		return
	}

	// populate test set's props with parent
	execHelper.PopulateTestSetProps(&req)
	for idx, _ := range req.TestSets {
		testSet := &req.TestSets[idx]

		if testSet.WorkspaceId != 0 {
			po, _ := c.WorkspaceService.Get(uint(testSet.WorkspaceId))
			testSet.WorkspacePath = po.Path
		}
	}

	ch = make(chan int, 1)
	go func() {
		execHelper.Exec(ch, req, &wsMsg)
		execHelper.SetRunning(false)
	}()

	execHelper.SetRunning(true)

	msg := i118Utils.Sprintf("start_task")
	websocketHelper.SendExecMsg(msg, "true", commConsts.Run,
		iris.Map{"status": "start-task"}, &wsMsg)
	logUtils.ExecConsole(color.FgCyan, msg)

	return
}
