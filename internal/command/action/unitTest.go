package action

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	_scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/exec"
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/kataras/iris/v12/websocket"
)

func RunUnitTest(cmdStr string) {
	req := serverDomain.WsReq{
		ProductId:   commConsts.ProductId,
		ProjectPath: commConsts.WorkDir,
		Act:         commConsts.ExecUnit,
		Cmd:         cmdStr,
		TestTool:    commConsts.UnitTestTool,
		BuildTool:   commConsts.UnitBuildTool,
	}
	if stringUtils.ParseInt(req.ProductId) != 0 {
		req.SubmitResult = true
	}

	msg := websocket.Message{}
	_scriptUtils.Exec(nil, req, msg)
}
