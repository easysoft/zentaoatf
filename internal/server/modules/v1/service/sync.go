package service

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/zentao"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/fatih/color"
	"path/filepath"
)

type SyncService struct {
	TestScriptService *TestScriptService `inject:""`
	TestCaseService   *TestCaseService   `inject:""`
}

func NewSyncService() *SyncService {
	return &SyncService{}
}

func (s *SyncService) SyncFromZentao(settings commDomain.SyncSettings, workspacePath string) (err error) {
	productId := settings.ProductId
	moduleId := settings.ModuleId
	suiteId := settings.SuiteId
	taskId := settings.TaskId

	byModule := settings.ByModule
	independentFile := settings.IndependentFile
	lang := settings.Lang

	ok := langUtils.CheckSupportLanguages(lang)
	if !ok {
		return
	}

	cases, loginFail := s.TestCaseService.LoadTestCases(productId, moduleId, suiteId, taskId, workspacePath)

	if cases != nil && len(cases) > 0 {
		productId = cases[0].Product
		targetDir := fileUtils.AddPathSepIfNeeded(filepath.Join(workspacePath, fmt.Sprintf("product%d", productId)))

		count, err := scriptUtils.GenerateScripts(cases, lang, independentFile, byModule, targetDir)
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

func (s *SyncService) SyncToZentao(workspacePath string, commitProductId int) (err error) {
	productPath := ""
	if commConsts.ComeFrom == "cmd" {
		productPath = fileUtils.RemovePathSepIfNeeded(workspacePath)
		workspacePath = commConsts.WorkDir
	} else {
		productPath = filepath.Join(workspacePath, fmt.Sprintf("product%d", commitProductId))
	}
	caseFiles := scriptUtils.LoadScriptByWorkspace(productPath)

	for _, cs := range caseFiles {
		pass, id, _, title := scriptUtils.GetCaseInfo(cs)

		if pass {
			steps, isOldFormat := scriptUtils.GetStepAndExpectMap(cs)
			if commConsts.Verbose {
				logUtils.Infof("isOldFormat = ", isOldFormat)
			}

			zentaoUtils.CommitCase(id, title, steps, workspacePath)
		}
	}

	return
}
