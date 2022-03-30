package service

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	langHelper "github.com/aaronchen2k/deeptest/internal/comm/helper/lang"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/zentao"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	"github.com/fatih/color"
	"path/filepath"
)

type SyncService struct {
	WorkspaceRepo     *repo.WorkspaceRepo `inject:""`
	TestScriptService *TestScriptService  `inject:""`
	TestCaseService   *TestCaseService    `inject:""`
	SiteService       *SiteService        `inject:""`
}

func NewSyncService() *SyncService {
	return &SyncService{}
}

func (s *SyncService) SyncFromZentao(settings commDomain.SyncSettings, config commDomain.WorkspaceConf, workspacePath string) (err error) {
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

	cases, loginFail := s.TestCaseService.LoadTestCases(productId, moduleId, suiteId, taskId, config)

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

func (s *SyncService) SyncToZentao(cases []string, workspacePath string, commitProductId int, config commDomain.WorkspaceConf) (err error) {
	pth := ""
	if commConsts.ComeFrom == "cmd" {
		pth = fileUtils.RemovePathSepIfNeeded(workspacePath)
		workspacePath = commConsts.WorkDir
	} else {
		pth = filepath.Join(workspacePath, fmt.Sprintf("product%d", commitProductId))
	}

	if cases == nil { // from command line
		cases = scriptUtils.LoadScriptByWorkspace(pth)
	}

	for _, cs := range cases {
		pass, id, _, title := scriptUtils.GetCaseInfo(cs)

		if pass {
			steps, _ := scriptUtils.GetStepAndExpectMap(cs)
			zentaoUtils.CommitCase(id, title, steps, config)
		}
	}

	return
}
