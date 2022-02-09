package action

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/command"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	_scriptUtils "github.com/aaronchen2k/deeptest/internal/server/modules/helper/exec"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/kataras/iris/v12/websocket"
)

func RunZTFTest(file []string, suiteIdStr, taskIdStr string, actionModule *command.IndexModule) error {
	cases := make([]string, 0)
	req := serverDomain.WsReq{
		ProjectPath: commConsts.WorkDir,
	}
	msg := websocket.Message{}

	if suiteIdStr != "" { // run with suite id
		req.SuiteId = suiteIdStr
		req.Act = commConsts.ExecSuite

	} else if taskIdStr != "" { // run with task id,
		req.TaskId = taskIdStr
		req.Act = commConsts.ExecTask

	} else {
		req.Act = commConsts.ExecCase
	}
	if !fileUtils.IsDir(file[0]) {
		cases = append(cases, file[0])
	} else {
		_, _, _, scriptTree, _ := actionModule.ProjectService.GetByUser(commConsts.WorkDir)
		cases = GetCasesFromChildren(scriptTree.Children)
	}

	req.Cases = cases

	_scriptUtils.Exec(nil, nil, nil, req, msg)

	return nil
}

// 扁平化
func GetCasesFromChildren(scripts []*serverDomain.TestAsset) (cases []string) {
	for _, v := range scripts {
		if v.Path != "" {
			cases = append(cases, v.Path)
		}
		if v.Children != nil {
			GetCasesFromChildren(v.Children)
		}
	}
	return
}
