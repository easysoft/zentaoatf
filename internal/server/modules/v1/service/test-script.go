package service

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/server/modules/helper/script"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
)

type TestScriptService struct {
	TestScriptRepo *repo.TestScriptRepo `inject:""`
}

func NewTestScriptService() *TestScriptService {
	return &TestScriptService{}
}

func (s *TestScriptService) GenerateScripts(cases []commDomain.ZtfCase, langType string, independentFile bool,
	byModule bool, targetDir string, prefix string) (int, error) {
	caseIds := make([]string, 0)

	for _, cs := range cases {
		scriptUtils.GenerateScript(cs, langType, independentFile, &caseIds, targetDir, byModule, prefix)
	}

	scriptUtils.GenSuite(caseIds, targetDir)

	return len(cases), nil
}

func (s *TestScriptService) Paginate(req serverDomain.TestScriptReqPaginate) (ret domain.PageData, err error) {
	ret, err = s.TestScriptRepo.Paginate(req)
	if err != nil {
		return
	}

	return
}

func (s *TestScriptService) FindById(id uint) (model.TestScript, error) {
	return s.TestScriptRepo.FindById(id)
}

func (s *TestScriptService) Create(testScript model.TestScript) (uint, error) {
	return s.TestScriptRepo.Create(testScript)
}

func (s *TestScriptService) Update(id uint, testScript model.TestScript) error {
	return s.TestScriptRepo.Update(id, testScript)
}

func (s *TestScriptService) DeleteById(id uint) error {
	return s.TestScriptRepo.DeleteById(id)
}
