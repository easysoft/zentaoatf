package action

import (
	"github.com/easysoft/zentaoatf/src/model"
	scriptService "github.com/easysoft/zentaoatf/src/service/script"
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	"github.com/easysoft/zentaoatf/src/utils/common"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
	"github.com/fatih/color"
	"strconv"
)

func Run(files []string, suite string, task string, result string) {
	caseIdMap := map[int]string{}
	cases := make([]string, 0)

	vari.WorkDir = fileUtils.AbosutePath(".")
	vari.RunDir = zentaoUtils.PathToRunName()

	if suite != "" {
		suiteId, err := strconv.Atoi(suite)
		if err == nil && suiteId > 0 { // load cases from remote by suite id
			zentaoService.GetCaseIdsBySuite(suite, &caseIdMap)
		} else { // load cases in suite file
			scriptService.GetCaseIdsInSuiteFile(suite, &caseIdMap)
		}

		scriptService.GetScriptByIdsInDir(files[0], caseIdMap, &cases)
	} else if task != "" { // load cases from remote by task id
		taskId, err := strconv.Atoi(task)
		if err == nil && taskId > 0 {
			zentaoService.GetCaseIdsByTask(task, &caseIdMap)
		}

		scriptService.GetScriptByIdsInDir(files[0], caseIdMap, &cases)
	} else if result != "" { // load cases result file
		scriptService.GetFailedCasesFromTestResult(result, &caseIdMap)

		scriptService.GetScriptByIdsInDir(files[0], caseIdMap, &cases)
	} else { // find cases in current dir
		for _, file := range files {
			scriptService.GetAllScriptsInDir(file, &cases)
		}
	}

	if len(cases) < 1 {
		logUtils.PrintToCmd(color.RedString("\n" + i118Utils.I118Prt.Sprintf("no_scripts")))
		return
	}

	var report = model.TestReport{Env: commonUtils.GetOs(),
		Pass: 0, Fail: 0, Total: 0, Cases: make([]model.CaseLog, 0)}

	testingService.ExeScripts(cases, &report)

	testingService.CheckResults(cases, &report)
	testingService.Print(report)
}
