package action

import (
	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	"github.com/easysoft/zentaoatf/internal/pkg/helper/exec"
	"github.com/easysoft/zentaoatf/internal/pkg/helper/script"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/pkg/consts"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"
	"path"
)

func RunZTFTest(files []string, moduleIdStr, suiteIdStr, taskIdStr string) error {
	req := serverDomain.WsReq{
		ScriptDirParamFromCmdLine: files[0],
	}
	testSet := serverDomain.TestSet{
		WorkspacePath: commConsts.ZtfDir,
	}

	if moduleIdStr != "" { // run with module id
		testSet.ProductId = stringUtils.ParseInt(commConsts.ProductId)
		testSet.ModuleId = stringUtils.ParseInt(moduleIdStr)
		req.Act = commConsts.ExecModule

	} else if suiteIdStr != "" { // run with suite id
		testSet.SuiteId = stringUtils.ParseInt(suiteIdStr)
		req.Act = commConsts.ExecSuite

	} else if taskIdStr != "" { // run with task id,
		testSet.TaskId = stringUtils.ParseInt(taskIdStr)
		req.Act = commConsts.ExecTask

	} else {
		cases := make([]string, 0)

		suite, dir := isRunWithSuiteFile(files)
		result := isRunWithResultFile(files)
		req.Act = commConsts.ExecCase
		if suite != "" { // run from suite file
			if dir == "" { // not found dir in files param
				dir = fileUtils.AbsolutePath(".")
			}
			cases = getCaseBySuiteFile(suite, dir)
		} else if result != "" { // run from failed result file
			cases = scriptHelper.GetFailedCasesDirectlyFromTestResult(result)
		} else { // run files
			for _, v1 := range files {
				if fileUtils.IsDir(v1) {
					temps := scriptHelper.LoadScriptByWorkspace(v1)
					for _, v2 := range temps {
						cases = append(cases, v2)
					}
				} else {
					cases = append(cases, v1)
				}
			}
		}

		testSet.Cases = cases
	}

	req.TestSets = append(req.TestSets, testSet)

	execHelper.Exec(nil, req, nil)

	return nil
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
	caseIds := make([]int, 0)
	caseIdMap := map[int]string{}

	scriptHelper.GetCaseIdsInSuiteFile(file, &caseIds)
	scriptHelper.GetScriptByIdsInDir(dir, &caseIdMap)

	cases := scriptHelper.GetCaseByListInMap(caseIds, caseIdMap)

	return cases
}
