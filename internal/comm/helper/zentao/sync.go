package zentaoHelper

import (
	"fmt"
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

	cases, err := LoadTestCaseDetail(productId, moduleId, suiteId, taskId, config)
	if err != nil {
		return
	}

	if cases != nil && len(cases) > 0 {
		productId = cases[0].Product
		targetDir := fileUtils.AddFilePathSepIfNeeded(filepath.Join(workspacePath, fmt.Sprintf("product%d", productId)))

		count, err := scriptHelper.GenerateScripts(cases, lang, independentFile, byModule, targetDir)
		if err == nil {
			logUtils.Infof(i118Utils.Sprintf("success_to_generate", count, targetDir))
		} else {
			logUtils.Infof(color.RedString(err.Error()))
		}
	}

	return
}

func SyncToZentao(cases []string, config commDomain.WorkspaceConf) (err error) {
	count := 0
	for _, cs := range cases {
		pass, id, _, title := scriptHelper.GetCaseInfo(cs)

		if pass && id > 0 {
			steps, _ := scriptHelper.GetStepAndExpectMap(cs)
			CommitCase(id, title, steps, config)

			count++
		}
	}

	logUtils.Infof(i118Utils.Sprintf("commit_cases_result", count) + "\n")

	return
}
