package action

import (
	"github.com/easysoft/zentaoatf/src/model"
	scriptService "github.com/easysoft/zentaoatf/src/service/script"
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	"github.com/easysoft/zentaoatf/src/utils/common"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	stdinUtils "github.com/easysoft/zentaoatf/src/utils/stdin"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
	"path"
	"strconv"
)

func Run(files []string, suiteIdStr string, taskIdStr string) error {
	vari.WorkDir = fileUtils.AbosutePath(".")
	vari.RunDir = zentaoUtils.RunDateFolder()

	cases := make([]string, 0)

	if suiteIdStr != "" { // run with suite id,
		if len(files) == 0 { // no dir
			files[0] = fileUtils.AbosutePath(".")
		}
		cases = getCaseBySuiteId(suiteIdStr, files[0])
	} else if taskIdStr != "" { // run with task id,
		if len(files) == 0 { // no dir
			files[0] = fileUtils.AbosutePath(".")
		}
		cases = getCaseByTaskId(taskIdStr, files[0])
	} else if path.Ext(files[0]) == "."+constant.ExtNameSuite {
		if len(files) == 1 { // run suite file, but no dir
			files[1] = fileUtils.AbosutePath(".")
		}

		cases = getCaseBySuiteFile(files[0], files[1])
	} else {
		cases = getCaseByDirAndFile(files)
	}

	if len(cases) < 1 {
		logUtils.PrintToCmd("\n"+i118Utils.I118Prt.Sprintf("no_scripts"), -1)
		return nil
	}

	runCases(cases)

	return nil
}

func getCaseByTaskId(id string, dir string) []string {
	caseIdMap := map[int]string{}
	cases := make([]string, 0)

	taskId, err := strconv.Atoi(id)
	if err == nil && taskId > 0 {

		stdinUtils.CheckRequestConfig()
		zentaoService.GetCaseIdsByTask(id, &caseIdMap)
	}

	scriptService.GetScriptByIdsInDir(dir, caseIdMap, &cases)
	return cases
}

func getCaseBySuiteId(id string, dir string) []string {
	caseIdMap := map[int]string{}
	cases := make([]string, 0)

	suiteId, err := strconv.Atoi(id)
	if err == nil && suiteId > 0 {
		stdinUtils.CheckRequestConfig()
		zentaoService.GetCaseIdsBySuite(id, &caseIdMap)
	}

	scriptService.GetScriptByIdsInDir(dir, caseIdMap, &cases)
	return cases
}

func getCaseBySuiteFile(file string, dir string) []string {
	caseIdMap := map[int]string{}
	cases := make([]string, 0)

	scriptService.GetCaseIdsInSuiteFile(file, &caseIdMap)
	scriptService.GetScriptByIdsInDir(dir, caseIdMap, &cases)

	return cases
}

func getCaseByDirAndFile(files []string) []string {
	cases := make([]string, 0)

	for _, file := range files {
		scriptService.GetAllScriptsInDir(file, &cases)
	}

	return cases
}

func runCases(cases []string) {
	var report = model.TestReport{Env: commonUtils.GetOs(),
		Pass: 0, Fail: 0, Total: 0, Cases: make([]model.CaseLog, 0)}

	testingService.ExeScripts(cases, &report)
	testingService.Report(report)
}
