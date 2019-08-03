package action

import (
	"github.com/easysoft/zentaoatf/src/biz"
	"github.com/easysoft/zentaoatf/src/misc"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/script"
	"github.com/easysoft/zentaoatf/src/utils"
	"strings"
)

func Run(scriptDir string, fileNames []string, langType string) {
	if strings.Index(scriptDir, "/") != 0 {
		scriptDir = utils.Prefer.WorkDir + scriptDir
	}

	LangMap := script.LangMap
	var files []string
	if fileNames != nil && len(fileNames) > 0 {
		if len(fileNames) == 1 {
			if strings.Index(fileNames[0], ".suite") > -1 {
				utils.RunMode = misc.SUITE
			} else {
				utils.RunMode = misc.SCRIPT
			}
			utils.RunDir = utils.PathToRunName(fileNames[0])
		} else {
			utils.RunMode = misc.BATCH
			utils.RunDir = utils.PathToRunName("")
		}

		files, _ = utils.GetSpecifiedFiles(scriptDir, fileNames)
	} else {
		files, _ = utils.GetAllFiles(scriptDir, LangMap[langType]["extName"])
		utils.RunMode = misc.DIR
		utils.RunDir = utils.PathToRunName(scriptDir)
	}

	var report = model.TestReport{Path: scriptDir, Env: utils.GetOs(),
		Pass: 0, Fail: 0, Total: 0, Cases: make([]model.CaseLog, 0)}

	biz.ExeScripts(files, scriptDir, langType, &report)

	biz.CheckResults(files, scriptDir, langType, &report)
	biz.Print(report, scriptDir)
}
