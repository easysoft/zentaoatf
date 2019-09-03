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

	if suiteIdStr != "" { // run with suite id
		dir := fileUtils.AbosutePath(".")
		if len(files) > 0 {
			dir = files[0]
		}

		cases = getCaseBySuiteId(suiteIdStr, dir)

	} else if taskIdStr != "" { // run with task id,
		dir := fileUtils.AbosutePath(".")
		if len(files) > 0 {
			dir = files[0]
		}

		cases = getCaseByTaskId(taskIdStr, dir)

	} else {
		suite, dir := isRunWithSuiteFile(files)

		if suite != "" {
			cases = getCaseBySuiteFile(suite, dir)
		} else {
			cases = getCaseByDirAndFile(files)
		}
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

func isRunWithSuiteFile(files []string) (string, string) {
	var suiteFile string
	var dir string

	for _, file := range files {
		if path.Ext(file) == "."+constant.ExtNameSuite {
			suiteFile = file
		} else {
			if fileUtils.IsDir(file) && dir != "" { // only select the first dir
				dir = file
			}
		}

		if suiteFile != "" && dir != "" {
			break
		}
	}

	if suiteFile != "" && dir == "" {
		dir = fileUtils.AbosutePath(".")
	}

	return suiteFile, dir
}
