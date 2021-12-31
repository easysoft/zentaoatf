package service

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
)

type TestSuiteService struct {
	TestSuiteRepo *repo.TestSuiteRepo `inject:""`
}

func NewTestSuiteService() *TestSuiteService {
	return &TestSuiteService{}
}

func (s *TestSuiteService) Paginate(req serverDomain.ReqPaginate) (ret domain.PageData, err error) {
	ret, err = s.TestSuiteRepo.Paginate(req)
	if err != nil {
		return
	}

	return
}

func (s *TestSuiteService) FindById(id uint) (model.TestSuite, error) {
	return s.TestSuiteRepo.FindById(id)
}

func (s *TestSuiteService) Create(testSuite model.TestSuite) (uint, error) {
	return s.TestSuiteRepo.Create(testSuite)
}

func (s *TestSuiteService) Update(id uint, testSuite model.TestSuite) error {
	return s.TestSuiteRepo.Update(id, testSuite)
}

func (s *TestSuiteService) DeleteById(id uint) error {
	return s.TestSuiteRepo.DeleteById(id)
}
