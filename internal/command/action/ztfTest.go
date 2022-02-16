package action

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	_scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/exec"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/zentao"
	"github.com/aaronchen2k/deeptest/internal/command"
	commandConfig "github.com/aaronchen2k/deeptest/internal/command/config"
	stringUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/string"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/kataras/iris/v12/websocket"
	"strconv"
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
		cases = getCaseBySuiteId(suiteIdStr, file[0])
	} else if taskIdStr != "" { // run with task id,
		req.TaskId = taskIdStr
		req.Act = commConsts.ExecTask

	} else {
		req.Act = commConsts.ExecCase
	}
	//if !fileUtils.IsDir(file[0]) {
	//	cases = append(cases, file[0])
	//} else {
	//	_, _, _, scriptTree, _ := actionModule.ProjectService.GetByUser(commConsts.WorkDir)
	//	cases = getCasesFromChildren(scriptTree.Children)
	//}

	req.Cases = cases

	_scriptUtils.Exec(nil, nil, nil, req, msg)

	return nil
}

// 扁平化
func getCasesFromChildren(scripts []*serverDomain.TestAsset) (cases []string) {
	for _, v := range scripts {
		if v.Path != "" {
			cases = append(cases, v.Path)
		}
		if v.Children != nil {
			getCasesFromChildren(v.Children)
		}
	}
	return
}
func getCaseBySuiteId(id string, dir string) []string {
	caseIdMap := map[int]string{}
	cases := make([]string, 0)

	suiteId, err := strconv.Atoi(id)
	if err == nil && suiteId > 0 {
		commandConfig.CheckRequestConfig()
		cases, err = zentaoUtils.GetCasesBySuite(0, stringUtils.ParseInt(id), dir)
	}
	scriptUtils.GetScriptByIdsInDir(dir, caseIdMap, &cases)
	return cases
}
