package zentaoHelper

import (
	"fmt"
	"path/filepath"

	commDomain "github.com/easysoft/zentaoatf/internal/comm/domain"
	langHelper "github.com/easysoft/zentaoatf/internal/comm/helper/lang"
	scriptHelper "github.com/easysoft/zentaoatf/internal/comm/helper/script"
	fileUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/file"
	i118Utils "github.com/easysoft/zentaoatf/internal/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	"github.com/fatih/color"
)

func SyncFromZentao(settings commDomain.SyncSettings, config commDomain.WorkspaceConf, workspacePath string) (
	pths []string, err error) {

	productId := settings.ProductId
	moduleId := settings.ModuleId
	suiteId := settings.SuiteId
	taskId := settings.TaskId
	caseId := settings.CaseId

	byModule := settings.SaveByModule
	independentFile := settings.IndependentFile
	lang := "php" // settings.Lang

	ok := langHelper.CheckSupportLanguages(lang)
	if !ok {
		return
	}

	cases := make([]commDomain.ZtfCase, 0)
	if caseId != 0 {
		cs, err := GetTestCaseDetail(caseId, config)
		if err == nil {
			cases = append(cases, cs)
		}
	} else {
		cases, err = LoadTestCasesDetail(productId, moduleId, suiteId, taskId, config)
	}

	if err != nil {
		return
	}

	if cases == nil || len(cases) == 0 {
		return
	}

	if productId == 0 {
		productId = cases[0].Product
	}
	targetDir := fileUtils.AddFilePathSepIfNeeded(filepath.Join(workspacePath, fmt.Sprintf("product%d", productId)))

	pths, targetDir, err = scriptHelper.GenerateScripts(cases, lang, independentFile, byModule, targetDir)
	if err == nil {
		logUtils.Infof(i118Utils.Sprintf("success_to_generate", len(pths), targetDir))
	} else {
		logUtils.Infof(color.RedString(err.Error()))
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
