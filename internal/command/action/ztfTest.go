package action

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	"github.com/aaronchen2k/deeptest/internal/comm/vari"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverLog "github.com/aaronchen2k/deeptest/internal/server/core/log"
	_scriptUtils "github.com/aaronchen2k/deeptest/internal/server/modules/utils/exec"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/server/modules/utils/script"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/kataras/iris/v12/websocket"
	"path"
	"strconv"
)

func RunZTFTest(files []string, suiteIdStr, taskIdStr string) error {
	serverLog.Init()
	cases := make([]string, 0)
	if suiteIdStr != "" { // run with suite id
		dir := fileUtils.AbsolutePath(".")
		if vari.ServerProjectDir != "" {
			dir = vari.ServerProjectDir
		} else if len(files) > 0 {
			dir = files[0]
		}

		cases = getCaseBySuiteId(suiteIdStr, dir)

	} else if taskIdStr != "" { // run with task id,
		dir := fileUtils.AbsolutePath(".")
		if vari.ServerProjectDir != "" {
			dir = vari.ServerProjectDir
		} else if len(files) > 0 {
			dir = files[0]
		}

		cases = getCaseByTaskId(taskIdStr, dir)

	} else {
		suite, dir := isRunWithSuiteFile(files)
		result := isRunWithResultFile(files)

		if suite != "" { // run from suite file
			if dir == "" { // not found dir in files param
				dir = fileUtils.AbsolutePath(".")
				if vari.ServerProjectDir != "" {
					dir = vari.ServerProjectDir
				}
			}

			cases = getCaseBySuiteFile(suite, dir)
		} else if result != "" { // run from failed result file
			//cases = assertUtils.GetFailedCasesDirectlyFromTestResult(result)
			scriptTree, _ := scriptUtils.LoadScriptTree(commConsts.WorkDir)
			cases = GetCasesFromChildren(scriptTree.Children)
		} else { // run files
			//cases = assertUtils.GetCaseByDirAndFile(files)
			scriptTree, _ := scriptUtils.LoadScriptTree(commConsts.WorkDir)
			cases = GetCasesFromChildren(scriptTree.Children)
		}
	}

	if len(cases) < 1 {
		logUtils.PrintTo("\n" + i118Utils.Sprintf("no_cases"))
		return nil
	}

	//runCases(cases)
	req := serverDomain.WsReq{
		Cases: cases,
		Act:   commConsts.ExecCase,
	}
	msg := websocket.Message{}
	_scriptUtils.Exec(nil, nil, nil, req, msg)

	return nil
}

func getCaseBySuiteId(id string, dir string) []string { // todo
	//caseIdMap := map[int]string{}
	cases := make([]string, 0)

	//suiteId, err := strconv.Atoi(id)

	return cases
}

func getCaseBySuiteFile(file string, dir string) []string {
	//caseIdMap := map[int]string{}
	cases := make([]string, 0)

	//assertUtils.GetCaseIdsInSuiteFile(file, &caseIdMap)
	//assertUtils.GetScriptByIdsInDir(dir, caseIdMap, &cases)

	return cases
}

func getCaseByTaskId(id string, dir string) []string {
	//caseIdMap := map[int]string{}
	cases := make([]string, 0)

	taskId, err := strconv.Atoi(id)
	if err == nil && taskId > 0 {
		//configUtils.CheckRequestConfig()
		//zentaoService.GetCaseIdsByTask(id, &caseIdMap)
	}

	//assertUtils.GetScriptByIdsInDir(dir, caseIdMap, &cases)
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
