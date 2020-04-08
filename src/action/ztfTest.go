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
	stringUtils "github.com/easysoft/zentaoatf/src/utils/string"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
	"github.com/mattn/go-runewidth"
	"path"
	"strconv"
)

func RunZTFTest(files []string, suiteIdStr string, taskIdStr string) error {
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
		logUtils.PrintTo("\n" + i118Utils.I118Prt.Sprintf("no_scripts"))
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

func runCases(files []string) {
	casesToRun := make([]string, 0)
	casesToIgnore := make([]string, 0)

	// config interpreter if needed
	if commonUtils.IsWin() {
		conf := configUtils.ReadCurrConfig()
		configChanged := configUtils.InputForScriptInterpreter(files, &conf, "run")
		if configChanged {
			configUtils.SaveConfig(conf)
		}
	}

	conf := configUtils.ReadCurrConfig()
	for _, file := range files {
		if commonUtils.IsWin() {
			if path.Ext(file) == ".sh" { // filter by os
				continue
			}

			ext := path.Ext(file)
			if ext != "" {
				ext = ext[1:]
			}
			lang := vari.ScriptExtToNameMap[ext]
			interpreter := commonUtils.GetFieldVal(conf, stringUtils.Ucfirst(lang))
			if interpreter == "-" && vari.Interpreter == "" { // not to ignore if interpreter set
				interpreter = ""

				casesToIgnore = append(casesToIgnore, file)
			}
			if lang != "bat" && interpreter == "" { // ignore the ones with no interpreter set
				continue
			}
		} else if !commonUtils.IsWin() { // filter by os
			if path.Ext(file) == ".bat" {
				continue
			}
		}

		casesToRun = append(casesToRun, file)
	}

	var report = model.TestReport{Env: commonUtils.GetOs(),
		Pass: 0, Fail: 0, Total: 0, FuncResult: make([]model.FuncResult, 0)}
	report.TestType = "func"
	report.TestFrame = "ztf"

	pathMaxWidth := 0
	numbMaxWidth := 0
	for _, file := range casesToRun {
		lent := runewidth.StringWidth(file)
		if lent > pathMaxWidth {
			pathMaxWidth = lent
		}

		content := fileUtils.ReadFile(file)
		caseId := zentaoUtils.ReadCaseId(content)
		if len(caseId) > numbMaxWidth {
			numbMaxWidth = len(caseId)
		}
	}

	testingService.ExeScripts(casesToRun, casesToIgnore, &report, pathMaxWidth, numbMaxWidth)
	testingService.GenZTFTestReport(report, pathMaxWidth)
}

func isRunWithSuiteFile(files []string) (string, string) {
	var suiteFile string
	var dir string

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

	if len(files) == 1 && suiteFile != "" && dir == "" { // no dir provided, not including a wrong dir param
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
