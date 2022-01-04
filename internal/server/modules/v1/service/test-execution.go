package service

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
)

type TestExecutionService struct {
	TestExecutionRepo *repo.TestExecutionRepo `inject:""`
}

func NewTestExecutionService() *TestExecutionService {
	return &TestExecutionService{}
}

func (s *TestExecutionService) Paginate(req serverDomain.TestExecutionReqPaginate) (ret domain.PageData, err error) {
	ret, err = s.TestExecutionRepo.Paginate(req)
	if err != nil {
		return
	}

	return
}

func (s *TestExecutionService) FindById(id uint) (model.TestExecution, error) {
	return s.TestExecutionRepo.FindById(id)
}

func (s *TestExecutionService) Create(testExecution model.TestExecution) (uint, error) {
	return s.TestExecutionRepo.Create(testExecution)
}

func (s *TestExecutionService) Update(id uint, testExecution model.TestExecution) error {
	return s.TestExecutionRepo.Update(id, testExecution)
}

func (s *TestExecutionService) DeleteById(id uint) error {
	return s.TestExecutionRepo.DeleteById(id)
}
