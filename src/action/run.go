package action

import (
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/service/script"
	testingService "github.com/easysoft/zentaoatf/src/service/testing"
	"github.com/easysoft/zentaoatf/src/utils/common"
	"github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/file"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	zentaoUtils "github.com/easysoft/zentaoatf/src/utils/zentao"
	"strings"
)

func Run(scriptDir string, fileNames []string, langType string) {
	LangMap := scriptService.LangMap
	var files []string

	if fileNames != nil && len(fileNames) > 0 { // pass a list, cui always
		if len(fileNames) == 1 {
			if strings.Index(fileNames[0], ".suite") > -1 {
				vari.RunMode = constant.RunModeSuite
			} else {
				vari.RunMode = constant.RunModeScript
			}
			vari.RunDir = zentaoUtils.PathToRunName(fileNames[0])
		} else {
			vari.RunMode = constant.RunModeBatch
			vari.RunDir = zentaoUtils.PathToRunName("")
		}

		files, _ = fileUtils.GetSpecifiedFilesInWorkDir(fileNames)

	} else { // give a dir
		fileUtils.GetAllFilesInDir(scriptDir, LangMap[langType]["extName"], &files)

		vari.RunMode = constant.RunModeDir
		vari.RunDir = zentaoUtils.PathToRunName(scriptDir)
	}

	var report = model.TestReport{Path: vari.Prefer.WorkDir, Env: commonUtils.GetOs(),
		Pass: 0, Fail: 0, Total: 0, Cases: make([]model.CaseLog, 0)}

	testingService.ExeScripts(files, vari.Prefer.WorkDir, langType, &report)

	testingService.CheckResults(files, vari.Prefer.WorkDir, langType, &report)
	testingService.Print(report, vari.Prefer.WorkDir)
}
