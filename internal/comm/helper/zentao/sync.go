package zentaoHelper

import (
	"fmt"
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/comm/domain"
	langHelper "github.com/easysoft/zentaoatf/internal/comm/helper/lang"
	scriptHelper "github.com/easysoft/zentaoatf/internal/comm/helper/script"
	fileUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/file"
	i118Utils "github.com/easysoft/zentaoatf/internal/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	"github.com/fatih/color"
	"path/filepath"
)

func SyncFromZentao(settings commDomain.SyncSettings, config commDomain.WorkspaceConf, workspacePath string) (err error) {
	productId := settings.ProductId
	moduleId := settings.ModuleId
	suiteId := settings.SuiteId
	taskId := settings.TaskId

	byModule := settings.ByModule
	independentFile := settings.IndependentFile
	lang := settings.Lang

	ok := langHelper.CheckSupportLanguages(lang)
	if !ok {
		return
	}

	cases, loginFail := LoadTestCases(productId, moduleId, suiteId, taskId, config)

	if cases != nil && len(cases) > 0 {
		productId = cases[0].Product
		targetDir := fileUtils.AddPathSepIfNeeded(filepath.Join(workspacePath, fmt.Sprintf("product%d", productId)))

		count, err := scriptHelper.GenerateScripts(cases, lang, independentFile, byModule, targetDir)
		if err == nil {
			logUtils.Infof(i118Utils.Sprintf("success_to_generate", count, targetDir))
		} else {
			logUtils.Infof(color.RedString(err.Error()))
		}
	} else {
		if !loginFail {
			logUtils.Infof(i118Utils.Sprintf("no_cases"))
		}
	}

	return
}

func SyncToZentao(cases []string, workspacePath string, commitProductId int, config commDomain.WorkspaceConf) (err error) {
	pth := ""
	if commConsts.ExecFrom == commConsts.FromCmd {
		pth = fileUtils.RemovePathSepIfNeeded(workspacePath)
		workspacePath = commConsts.WorkDir
	} else {
		pth = filepath.Join(workspacePath, fmt.Sprintf("product%d", commitProductId))
	}

	if cases == nil { // from command line
		cases = scriptHelper.LoadScriptByWorkspace(pth)
	}

	for _, cs := range cases {
		pass, id, _, title := scriptHelper.GetCaseInfo(cs)

		if pass {
			steps, _ := scriptHelper.GetStepAndExpectMap(cs)
			CommitCase(id, title, steps, config)
		}
	}

	return
}
