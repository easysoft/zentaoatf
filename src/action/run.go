package action

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/biz"
	"github.com/easysoft/zentaoatf/src/misc"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/script"
	"github.com/easysoft/zentaoatf/src/utils"
	"strings"
)

func Run(scriptDir string, fileNames []string, langType string) {
	LangMap := script.LangMap
	var files []string

	if fileNames != nil && len(fileNames) > 0 { // pass a list, cui always
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
	} else { // give a dir
		utils.GetAllFiles(scriptDir, LangMap[langType]["extName"], &files)
		fmt.Printf("%v", scriptDir)
		fmt.Printf("%v", files)
		utils.RunMode = misc.DIR
		utils.RunDir = utils.PathToRunName(scriptDir)
	}

	var report = model.TestReport{Path: utils.Prefer.WorkDir, Env: utils.GetOs(),
		Pass: 0, Fail: 0, Total: 0, Cases: make([]model.CaseLog, 0)}

	biz.ExeScripts(files, utils.Prefer.WorkDir, langType, &report)

	biz.CheckResults(files, utils.Prefer.WorkDir, langType, &report)
	biz.Print(report, utils.Prefer.WorkDir)
}
