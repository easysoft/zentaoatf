package service

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/server/modules/v1/utils/script"
)

type ZtfScriptService struct {
	ZtfCaseService *ZtfCaseService   `inject:""`
	ProjectRepo    *repo.ProjectRepo `inject:""`
}

func NewZtfScriptService() *ZtfScriptService {
	return &ZtfScriptService{}
}

func (s *ZtfScriptService) Generate(cases []commDomain.ZtfCase, langType string, independentFile bool,
	byModule bool, targetDir string, prefix string) (int, error) {
	caseIds := make([]string, 0)
	for _, cs := range cases {
		scriptUtils.GenerateScript(cs, langType, independentFile, &caseIds, targetDir, byModule, prefix)
	}

	scriptUtils.GenSuite(caseIds, targetDir)

	return len(cases), nil
}
