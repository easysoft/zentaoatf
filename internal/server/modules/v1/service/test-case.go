package service

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	configUtils "github.com/aaronchen2k/deeptest/internal/server/modules/helper/config"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/server/modules/helper/zentao"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
)

type TestCaseService struct {
	TestCaseRepo *repo.TestCaseRepo `inject:""`
}

func NewTestCaseService() *TestCaseService {
	return &TestCaseService{}
}

func (s *TestCaseService) LoadTestCases(productId, moduleId, suiteId, taskId int, projectPath string) (
	cases []commDomain.ZtfCase, loginFail bool) {

	config := configUtils.LoadByProjectPath(projectPath)

	ok := zentaoUtils.Login(config)
	if !ok {
		loginFail = true
		return
	}

	if moduleId != 0 {
		cases = zentaoUtils.ListCaseByModule(config.Url, productId, moduleId)
	} else if suiteId != 0 {
		cases = zentaoUtils.ListCaseBySuite(config.Url, 0, suiteId)
	} else if taskId != 0 {
		cases = zentaoUtils.ListCaseByTask(config.Url, 0, taskId)
	} else if productId != 0 {
		cases = zentaoUtils.ListCaseByProduct(config.Url, productId)
	}

	return
}

func (s *TestCaseService) Paginate(req serverDomain.TestCaseReqPaginate) (ret domain.PageData, err error) {
	ret, err = s.TestCaseRepo.Paginate(req)
	if err != nil {
		return
	}

	return
}

func (s *TestCaseService) FindById(id uint) (model.TestCase, error) {
	return s.TestCaseRepo.FindById(id)
}

func (s *TestCaseService) Create(testCase model.TestCase) (uint, error) {
	return s.TestCaseRepo.Create(testCase)
}

func (s *TestCaseService) Update(id uint, testCase model.TestCase) error {
	return s.TestCaseRepo.Update(id, testCase)
}

func (s *TestCaseService) DeleteById(id uint) error {
	return s.TestCaseRepo.DeleteById(id)
}
