package action

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	_scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/exec"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/zentao"
	"github.com/aaronchen2k/deeptest/internal/command"
	commandConfig "github.com/aaronchen2k/deeptest/internal/command/config"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/kataras/iris/v12/websocket"
	"path"
	"strconv"
)

func RunZTFTest(files []string, suiteIdStr, taskIdStr string, actionModule *command.IndexModule) error {
	cases := make([]string, 0)
	req := serverDomain.WsReq{
		ProjectPath: commConsts.WorkDir,
	}
	msg := websocket.Message{}

	if suiteIdStr != "" { // run with suite id
		req.SuiteId = suiteIdStr
		req.Act = commConsts.ExecSuite
		cases = getCaseBySuiteId(suiteIdStr, files[0])
	} else if taskIdStr != "" { // run with task id,
		req.TaskId = taskIdStr
		req.Act = commConsts.ExecTask
		cases = getCaseByTaskId(taskIdStr, files[0])
	} else {
		suite, dir := isRunWithSuiteFile(files)
		result := isRunWithResultFile(files)
		req.Act = commConsts.ExecCase
		if suite != "" { // run from suite file
			if dir == "" { // not found dir in files param
				dir = fileUtils.AbsolutePath(".")
			}
			cases = getCaseBySuiteFile(suite, dir)
		} else if result != "" { // run from failed result file
			cases = scriptUtils.GetFailedCasesDirectlyFromTestResult(result)
		} else { // run files
			for _, v1 := range files {
				if fileUtils.IsDir(v1) {
					temps := scriptUtils.LoadScriptByProject(v1)
					for _, v2 := range temps {
						cases = append(cases, v2)
					}
				} else {
					cases = append(cases, v1)
				}
			}
		}

	}

	req.Cases = cases
	_scriptUtils.Exec(nil, nil, nil, req, msg)

	return nil
}

func getCaseBySuiteId(id string, dir string) []string {
	caseIdMap := map[int]string{}
	cases := make([]string, 0)

	suiteId, err := strconv.Atoi(id)
	if err == nil && suiteId > 0 {
		commandConfig.CheckRequestConfig()
		cases, err = zentaoUtils.GetCasesBySuite(0, suiteId, dir)
	}
	scriptUtils.GetScriptByIdsInDir(dir, caseIdMap, &cases)
	return cases
}

func getCaseByTaskId(id string, dir string) []string {
	caseIdMap := map[int]string{}
	cases := make([]string, 0)

	taskId, err := strconv.Atoi(id)
	if err == nil && taskId > 0 {
		commandConfig.CheckRequestConfig()
		cases, err = zentaoUtils.GetCasesByTask(0, taskId, dir)
	}

	scriptUtils.GetScriptByIdsInDir(dir, caseIdMap, &cases)
	return cases
}

func isRunWithSuiteFile(files []string) (suiteFile, dir string) {
	for _, file := range files {
		if path.Ext(file) == "."+consts.ExtNameSuite {
			suiteFile = file
		} else {
			if dir == "" { // only select the first dir
				dir = file
			}
		}

		if suiteFile != "" && dir != "" {
			break
		}
	}

	return
}

func isRunWithResultFile(files []string) string {
	var resultFile string

	for _, file := range files {
		if path.Ext(file) == "."+consts.ExtNameResult || path.Ext(file) == "."+consts.ExtNameJson {
			resultFile = file

			return resultFile
		}
	}

	return ""
}

func getCaseBySuiteFile(file string, dir string) []string {
	caseIdMap := map[int]string{}
	cases := make([]string, 0)

	scriptUtils.GetCaseIdsInSuiteFile(file, &caseIdMap)
	scriptUtils.GetScriptByIdsInDir(dir, caseIdMap, &cases)

	return cases
}
