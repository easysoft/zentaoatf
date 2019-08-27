package action

import (
	"github.com/easysoft/zentaoatf/src/model"
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	zentaoService "github.com/easysoft/zentaoatf/src/service/zentao"
	"github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/file"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"github.com/fatih/color"
	"path"
	"path/filepath"
	"strconv"
)

func Run(files []string, suite string, task string, result string) {
	caseIdMap := map[int]string{}
	cases := make([]string, 0)

	vari.RunDir, _ = filepath.Abs("")

	if suite != "" {
		suiteId, err := strconv.Atoi(suite)
		if err == nil && suiteId > 0 { // load cases from remote by suite id
			zentaoService.GetCaseIdsBySuite(suiteId, &caseIdMap)
		} else { // load cases in suite file
			fileUtils.GetCaseIdsInSuiteFile(suite, &caseIdMap)
		}

		fileUtils.GetScriptByIdsInDir(files[0], caseIdMap, &cases)
	} else if task != "" { // load cases from remote by task id
		taskId, err := strconv.Atoi(suite)
		if err == nil {
			zentaoService.GetCaseIdsByTask(taskId, &caseIdMap)
		}

		fileUtils.GetScriptByIdsInDir(files[0], caseIdMap, &cases)
	} else if result != "" { // load cases result file
		fileUtils.GetFailedCasesFromTestResult(result, &caseIdMap)

		fileUtils.GetScriptByIdsInDir(files[0], caseIdMap, &cases)
	} else { // find cases in current dir
		for _, file := range files {
			if !path.IsAbs(file) {
				file, _ = filepath.Abs(file)
			}

			fileUtils.GetAllScriptsInDir(file, &cases)
		}
	}

	if len(files) < 1 {
		logUtils.PrintToCmd(color.RedString("\n" + i118Utils.I118Prt.Sprintf("no_scripts")))
		return
	}

	var report = model.TestReport{Path: vari.ReportDir, Env: commonUtils.GetOs(),
		Pass: 0, Fail: 0, Total: 0, Cases: make([]model.CaseLog, 0)}

	testingService.ExeScripts(files, &report)

	testingService.CheckResults(files, &report)
	testingService.Print(report)
}
