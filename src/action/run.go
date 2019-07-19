package action

import (
	"github.com/easysoft/zentaoatf/src/biz"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/utils"
)

func Run(scriptDir string, fileNames []string, langType string) {
	var files []string
	if fileNames != nil && len(fileNames) > 0 {
		files, _ = utils.GetSpecifiedFiles(scriptDir, fileNames)
	} else {
		files, _ = utils.GetAllFiles(scriptDir, langType)
	}

	var report = model.TestReport{Path: scriptDir, Env: utils.GetOs(),
		Pass: 0, Fail: 0, Total: 0, Cases: make([]model.CaseLog, 0)}

	biz.ExeScripts(files, scriptDir, langType, &report)

	biz.CheckResults(files, scriptDir, langType, &report)
	biz.Print(report, scriptDir)
}
