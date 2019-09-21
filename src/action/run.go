package action

import (
	"github.com/easysoft/zentaoatf/src/model"
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	assertUtils "github.com/easysoft/zentaoatf/src/utils/assert"
	"github.com/easysoft/zentaoatf/src/utils/common"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/mattn/go-runewidth"
	"path"
	"strconv"
)

func Run(files []string, suiteIdStr string, taskIdStr string) error {
	logUtils.InitLogger()

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
		result := isRunWithResultFile(files)

		if suite != "" {
			cases = getCaseBySuiteFile(suite, dir)
		} else if result != "" {
			cases = assertUtils.GetFailedCasesDirectlyFromTestResult(result)
		} else {
			cases = assertUtils.GetCaseByDirAndFile(files)
		}
	}

	if len(cases) < 1 {
		logUtils.PrintToCmd("\n"+i118Utils.I118Prt.Sprintf("no_scripts"), -1)
		return nil
	}

	//if commonUtils.IsWin() {
	conf := configUtils.ReadCurrConfig()
	configUtils.InputForScriptInterpreter(cases, &conf)
	configUtils.SaveConfig(conf)
	//}

	runCases(cases)

	return nil
}

func getCaseByTaskId(id string, dir string) []string {
	caseIdMap := map[int]string{}
	cases := make([]string, 0)

	taskId, err := strconv.Atoi(id)
	if err == nil && taskId > 0 {

		configUtils.CheckRequestConfig()
		zentaoService.GetCaseIdsByTask(id, &caseIdMap)
	}

	assertUtils.GetScriptByIdsInDir(dir, caseIdMap, &cases)
	return cases
}

func getCaseBySuiteId(id string, dir string) []string {
	caseIdMap := map[int]string{}
	cases := make([]string, 0)

	suiteId, err := strconv.Atoi(id)
	if err == nil && suiteId > 0 {
		configUtils.CheckRequestConfig()
		zentaoService.GetCaseIdsBySuite(id, &caseIdMap)
	}

	assertUtils.GetScriptByIdsInDir(dir, caseIdMap, &cases)
	return cases
}

func getCaseBySuiteFile(file string, dir string) []string {
	caseIdMap := map[int]string{}
	cases := make([]string, 0)

	assertUtils.GetCaseIdsInSuiteFile(file, &caseIdMap)
	assertUtils.GetScriptByIdsInDir(dir, caseIdMap, &cases)

	return cases
}

func runCases(cases []string) {
	var report = model.TestReport{Env: commonUtils.GetOs(),
		Pass: 0, Fail: 0, Total: 0, Cases: make([]model.CaseLog, 0)}

	pathMaxWidth := 0
	for _, file := range cases {
		lent := runewidth.StringWidth(file)
		if lent > pathMaxWidth {
			pathMaxWidth = lent
		}
	}

	testingService.ExeScripts(cases, &report, pathMaxWidth)
	testingService.Report(report, pathMaxWidth)
}

func isRunWithSuiteFile(files []string) (string, string) {
	var suiteFile string
	var dir string

	for _, file := range files {
		if path.Ext(file) == "."+constant.ExtNameSuite {
			suiteFile = file
		} else {
			if fileUtils.IsDir(file) && dir == "" { // only select the first dir
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

func isRunWithResultFile(files []string) string {
	var resultFile string

	for _, file := range files {
		if path.Ext(file) == "."+constant.ExtNameResult || path.Ext(file) == "."+constant.ExtNameJson {
			resultFile = file

			return resultFile
		}
	}

	return ""
}
