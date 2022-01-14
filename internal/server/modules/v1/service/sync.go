package service

import (
	"fmt"
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	i118Utils "github.com/aaronchen2k/deeptest/internal/pkg/lib/i118"
	langUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/lang"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/zentao"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/script"
	"path/filepath"
	"strconv"
)

type SyncService struct {
	ZtfScriptService *ZtfScriptService `inject:""`
	ZtfCaseService   *ZtfCaseService   `inject:""`
	AssetService     *AssetService     `inject:""`
}

func NewSyncService() *SyncService {
	return &SyncService{}
}

func (s *SyncService) SyncFromZentao(settings commDomain.SyncSettings, projectPath string) (err error) {
	productId := settings.ProductId
	moduleId := settings.ModuleId
	suiteId := settings.SuiteId
	taskId := settings.TaskId
	independentFile := settings.IndependentFile
	lang := settings.Lang

	ok := langUtils.CheckSupportLanguages(lang)
	if !ok {
		return
	}

	cases, loginFail := s.ZtfCaseService.LoadTestCases(productId, moduleId, suiteId, taskId, projectPath)

	if cases != nil && len(cases) > 0 {
		productId, _ = strconv.Atoi(cases[0].Product)
		targetDir := fileUtils.AddPathSepIfNeeded(fmt.Sprintf("product%d", productId))
		prefix := ""
		byModule := moduleId > 0

		count, err := s.ZtfScriptService.Generate(cases, lang, independentFile, byModule, targetDir, prefix)
		if err == nil {
			logUtils.Infof(i118Utils.Sprintf("success_to_generate", count, targetDir))
		} else {
			logUtils.Infof(err.Error())
		}
	} else {
		if !loginFail {
			logUtils.Infof(i118Utils.Sprintf("no_cases"))
		}
	}

	return
}

func (s *SyncService) SyncToZentao(projectPath string, commitProductId int) (err error) {
	productPath := filepath.Join(projectPath, fmt.Sprintf("product%d", commitProductId))
	caseFiles := s.AssetService.LoadScriptByProject(productPath)

	for _, cs := range caseFiles {
		pass, id, _, title := zentaoUtils.GetCaseInfo(cs)

		if pass {
			stepMap, stepTypeMap, expectMap, isOldFormat := scriptUtils.GetStepAndExpectMap(cs)
			if commConsts.Verbose {
				logUtils.Infof("isOldFormat = ", isOldFormat)
			}

			isIndependent, expectIndependentContent := zentaoUtils.GetDependentExpect(cs)
			if isIndependent {
				expectMap = scriptUtils.GetExpectMapFromIndependentFileObsolete(expectMap, expectIndependentContent, true)
			}

			s.ZtfCaseService.CommitCase(id, title, stepMap, stepTypeMap, expectMap, projectPath)
		}
	}

	return
}
