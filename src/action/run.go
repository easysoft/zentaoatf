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
				utils.RunType = misc.SUITE
			} else {
				utils.RunType = misc.SCRIPT
			}
		} else {
			utils.RunType = misc.LIST
		}

		files, _ = utils.GetSpecifiedFiles(scriptDir, fileNames)
	} else {
		files, _ = utils.GetAllFiles(scriptDir, LangMap[langType]["extName"])
		utils.RunType = misc.DIR
	}

	var report = model.TestReport{Path: scriptDir, Env: utils.GetOs(),
		Pass: 0, Fail: 0, Total: 0, Cases: make([]model.CaseLog, 0)}

	biz.ExeScripts(files, scriptDir, langType, &report)

	biz.CheckResults(files, scriptDir, langType, &report)
	biz.Print(report, scriptDir)
}
