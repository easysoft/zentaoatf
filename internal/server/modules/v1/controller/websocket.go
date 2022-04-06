package controller

import (
	"encoding/json"
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/exec"
	websocketUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/websocket"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/config"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/service"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12/websocket"
	"strings"
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

func (c *WebSocketCtrl) OnNamespaceConnected(msg websocket.Message) error {
	websocketUtils.SetConn(c.Conn)

	logUtils.Infof(i118Utils.Sprintf("ws_namespace_connected", c.Conn.ID(), msg.Room))

	data := serverDomain.WsResp{Msg: "from server: connected to websocket"}
	websocketUtils.Broadcast(msg.Namespace, "", "OnVisit", data)
	return nil
}

// OnNamespaceDisconnect
// This will call the "OnVisit" event on all clients, except the current one,
// it can't because it's left but for any case use this type of design.
func (c *WebSocketCtrl) OnNamespaceDisconnect(msg websocket.Message) error {
	logUtils.Infof(i118Utils.Sprintf("ws_namespace_disconnected", c.Conn.ID()))

	data := serverDomain.WsResp{Msg: fmt.Sprintf("ws_connected")}
	websocketUtils.Broadcast(msg.Namespace, "", "OnVisit", data)
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
		websocketUtils.SendExecMsg(msg, "", commConsts.Error, wsMsg)
		logUtils.ExecConsole(color.FgRed, msg)
		return
	}

	act := req.Act

	if act == commConsts.ExecInit {
		msg := i118Utils.Sprintf("success_to_conn")
		//websocketUtils.SendExecMsg(msg, strconv.FormatBool(scriptUtils.GetRunning()), wsMsg)
		logUtils.ExecConsole(color.FgCyan, msg)
		return
	}

	if act == commConsts.ExecStop {
		if ch != nil {
			if !scriptUtils.GetRunning() {
				ch = nil
			} else {
				ch <- 1
				ch = nil
			}
		}

		scriptUtils.SetRunning(false)

		msg := i118Utils.Sprintf("end_task")
		websocketUtils.SendExecMsg(msg, "false", commConsts.Run, wsMsg)
		logUtils.ExecConsole(color.FgCyan, msg)
		return
	}

	if scriptUtils.GetRunning() && (act == commConsts.ExecCase || act == commConsts.ExecModule ||
		act == commConsts.ExecSuite || act == commConsts.ExecTask || act == commConsts.ExecUnit) {
		msg := i118Utils.Sprintf("pls_stop_previous")
		websocketUtils.SendExecMsg(msg, "true", commConsts.Run, wsMsg)
		logUtils.ExecConsole(color.FgRed, msg)

		return
	}

	for idx, _ := range req.TestSets {
		testSet := &req.TestSets[idx]
		if testSet.WorkspaceId != 0 {
			po, _ := c.WorkspaceService.Get(uint(testSet.WorkspaceId))
			testSet.WorkspacePath = po.Path
		}

		testSet.Scope = req.Scope
		testSet.Seq = req.Seq

		if testSet.ProductId == 0 && req.ProductId != 0 {
			testSet.ProductId = req.ProductId
		}

		testSet.ModuleId = req.ModuleId
		testSet.SuiteId = req.SuiteId
		testSet.TaskId = req.TaskId

		testSet.ScriptDirParamFromCmdLine = req.ScriptDirParamFromCmdLine

		setTestTool(testSet, req)
		setBuildTool(testSet, req)

		if testSet.Cmd == "" && req.Cmd != "" {
			testSet.Cmd = req.Cmd
		}

		testSet.SubmitResult = req.SubmitResult

	}

	ch = make(chan int, 1)
	go func() {
		scriptUtils.Exec(ch, req, wsMsg)
		scriptUtils.SetRunning(false)
	}()

	scriptUtils.SetRunning(true)

	msg := i118Utils.Sprintf("start_task")
	websocketUtils.SendExecMsg(msg, "true", commConsts.Run, wsMsg)
	logUtils.ExecConsole(color.FgCyan, msg)

	return
}

func setTestTool(testSet *serverDomain.TestSet, req serverDomain.WsReq) {
	if testSet.TestTool == "" && req.TestTool != "" {
		testSet.TestTool = req.TestTool
	}

	if testSet.TestTool == "" {
		testSet.TestTool = testSet.WorkspaceType
	}
}

func setBuildTool(testSet *serverDomain.TestSet, req serverDomain.WsReq) {
	if testSet.BuildTool == "" && req.BuildTool != "" {
		testSet.BuildTool = req.BuildTool
	}

	if testSet.BuildTool == "" {
		arr := strings.Split(testSet.Cmd, " ")
		testSet.BuildTool = commConsts.UnitBuildToolMap[strings.TrimSpace(arr[0])]
	}
}
