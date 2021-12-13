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
	"github.com/easysoft/zentaoatf/src/utils/vari"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
	"github.com/mattn/go-runewidth"
	"path"
	"strconv"
)

func RunZTFTest(files []string, suiteIdStr, taskIdStr string) error {
	logUtils.InitLogger()

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
			cases = assertUtils.GetFailedCasesDirectlyFromTestResult(result)
		} else { // run files
			cases = assertUtils.GetCaseByDirAndFile(files)
		}
	}

	if len(cases) < 1 {
		logUtils.PrintTo("\n" + i118Utils.Sprintf("no_cases"))
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
	casesToRun, casesToIgnore := filterCases(cases)

	var report = model.TestReport{Env: commonUtils.GetOs(),
		Pass: 0, Fail: 0, Total: 0, FuncResult: make([]model.FuncResult, 0)}
	report.TestType = "func"
	report.TestFrame = constant.AppName

	pathMaxWidth := 0
	numbMaxWidth := 0
	for _, cs := range casesToRun {
		lent := runewidth.StringWidth(cs)
		if lent > pathMaxWidth {
			pathMaxWidth = lent
		}

		content := fileUtils.ReadFile(cs)
		caseId := zentaoUtils.ReadCaseId(content)
		if len(caseId) > numbMaxWidth {
			numbMaxWidth = len(caseId)
		}
	}

	testingService.ExeScripts(casesToRun, casesToIgnore, &report, pathMaxWidth, numbMaxWidth)
	testingService.GenZTFTestReport(report, pathMaxWidth)
}

func isRunWithSuiteFile(files []string) (suiteFile, dir string) {
	for _, file := range files {
		if path.Ext(file) == "."+constant.ExtNameSuite {
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
		if path.Ext(file) == "."+constant.ExtNameResult || path.Ext(file) == "."+constant.ExtNameJson {
			resultFile = file

			return resultFile
		}
	}

	return ""
}
