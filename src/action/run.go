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

func Run(files []string, suiteIdStr string, taskIdStr string) {
	caseIdMap := map[int]string{}
	cases := make([]string, 0)

	vari.WorkDir = fileUtils.AbosutePath(".")
	vari.RunDir = zentaoUtils.RunDateFolder()

	if suiteIdStr != "" {
		suiteId, err := strconv.Atoi(suiteIdStr)
		if err == nil && suiteId > 0 {

			stdinUtils.CheckRequestConfig()
			zentaoService.GetCaseIdsBySuite(suiteIdStr, &caseIdMap)
		}

		scriptService.GetScriptByIdsInDir(files[0], caseIdMap, &cases)
	} else if taskIdStr != "" { // load cases from remote by taskIdStr id
		taskId, err := strconv.Atoi(taskIdStr)
		if err == nil && taskId > 0 {

			stdinUtils.CheckRequestConfig()
			zentaoService.GetCaseIdsByTask(taskIdStr, &caseIdMap)
		}

		scriptService.GetScriptByIdsInDir(files[0], caseIdMap, &cases)
	} else { // no suiteId, taskId param
		if len(files) > 1 && fileUtils.IsDir(files[0]) &&
			path.Ext(files[1]) == "."+constant.ExtNameSuite { // run suite file

			scriptService.GetCaseIdsInSuiteFile(files[1], &caseIdMap)
			scriptService.GetScriptByIdsInDir(files[0], caseIdMap, &cases)

		} else if len(files) > 1 && fileUtils.IsDir(files[0]) &&
			path.Ext(files[1]) == "."+constant.ExtNameResult { // run result file

			scriptService.GetFailedCasesFromTestResult(files[1], &caseIdMap)
			scriptService.GetScriptByIdsInDir(files[0], caseIdMap, &cases)

		} else { // run with dir and script files
			for _, file := range files {
				scriptService.GetAllScriptsInDir(file, &cases)
			}
		}
	}

	if len(cases) < 1 {
		logUtils.PrintToCmd("\n"+i118Utils.I118Prt.Sprintf("no_scripts")+"\n", -1)
		return
	}

	var report = model.TestReport{Env: commonUtils.GetOs(),
		Pass: 0, Fail: 0, Total: 0, Cases: make([]model.CaseLog, 0)}

	testingService.ExeScripts(cases, &report)

	testingService.Report(report)
}
