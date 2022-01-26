package scriptUtils

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	serverLog "github.com/aaronchen2k/deeptest/internal/server/core/log"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/kataras/iris/v12/websocket"
)

func Exec(ch chan int, sendOutputMsg, sendExecMsg func(info, isRunning string, msg websocket.Message), req serverDomain.WsReq, msg websocket.Message) (
	err error) {

	serverLog.InitExecLog(req.ProjectPath)

	if req.Act == commConsts.ExecCase {
		ExecCase(ch, sendOutputMsg, sendExecMsg, req, msg)
	} else if req.Act == commConsts.ExecModule {
		ExecModule(ch, sendOutputMsg, sendExecMsg, req, msg)
	} else if req.Act == commConsts.ExecSuite {
		ExecSuite(ch, sendOutputMsg, sendExecMsg, req, msg)
	} else if req.Act == commConsts.ExecTask {
		ExecTask(ch, sendOutputMsg, sendExecMsg, req, msg)
	} else if req.Act == commConsts.ExecUnit {
		ExecUnit(ch, sendOutputMsg, sendExecMsg, req, msg)
	}

	return
}
