package service

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
)

type TestExecService struct {
	TestExecRepo *repo.TestExecRepo `inject:""`
}

func NewTestExecService() *TestExecService {
	return &TestExecService{}
}

func (s *TestExecService) Paginate(req serverDomain.TestExecReqPaginate) (ret domain.PageData, err error) {
	ret, err = s.TestExecRepo.Paginate(req)
	if err != nil {
		return
	}

	return
}

func (s *TestExecService) FindById(id uint) (model.TestExec, error) {
	return s.TestExecRepo.FindById(id)
}

func (s *TestExecService) Create(testExecution model.TestExec) (uint, error) {
	return s.TestExecRepo.Create(testExecution)
}

func (s *TestExecService) Update(id uint, testExecution model.TestExec) error {
	return s.TestExecRepo.Update(id, testExecution)
}

func (s *TestExecService) DeleteById(id uint) error {
	return s.TestExecRepo.DeleteById(id)
}
