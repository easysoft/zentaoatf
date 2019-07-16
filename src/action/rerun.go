package action

import (
	"github.com/easysoft/zentaoatf/src/biz"
	"github.com/easysoft/zentaoatf/src/model"
	. "github.com/easysoft/zentaoatf/src/utils"
)

func Rerun(resultFile string) {
	files, scriptDir, langType, _ := GetFailedFiles(resultFile)

	var report = model.TestReport{Path: scriptDir, Env: GetOs(),
		Pass: 0, Fail: 0, Total: 0, Cases: make([]model.CaseLog, 0)}

	biz.ExeScripts(files, scriptDir, langType, &report)

	biz.CheckResults(files, scriptDir, langType, &report)
	biz.Print(report, scriptDir)
}
