package execHelper

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	serverConfig "github.com/aaronchen2k/deeptest/internal/server/config"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/kataras/iris/v12/websocket"
)

func Exec(ch chan int, req serverDomain.WsReq, msg websocket.Message) (
	err error) {

	serverConfig.InitExecLog(req.WorkspacePath)

	if req.ScriptDirParamFromCmdLine == "" {
		req.ScriptDirParamFromCmdLine = req.WorkspacePath
	}

	if req.Act == commConsts.ExecCase {
		ExecCase(ch, req, msg)
	} else if req.Act == commConsts.ExecModule {
		ExecModule(ch, req, msg)
	} else if req.Act == commConsts.ExecSuite {
		ExecSuite(ch, req, msg)
	} else if req.Act == commConsts.ExecTask {
		ExecTask(ch, req, msg)
	} else if req.Act == commConsts.ExecUnit {
		ExecUnit(ch, req, msg)
	}

	return
}
