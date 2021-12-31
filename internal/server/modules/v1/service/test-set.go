package service

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
)

type TestSetService struct {
	TestSetRepo *repo.TestSetRepo `inject:""`
}

func NewTestSetService() *TestSetService {
	return &TestSetService{}
}

func (s *TestSetService) Paginate(req serverDomain.ReqPaginate) (ret domain.PageData, err error) {
	ret, err = s.TestSetRepo.Paginate(req)
	if err != nil {
		return
	}

	return
}

func (s *TestSetService) FindById(id uint) (model.TestSet, error) {
	return s.TestSetRepo.FindById(id)
}

func (s *TestSetService) Create(testSet model.TestSet) (uint, error) {
	return s.TestSetRepo.Create(testSet)
}

func (s *TestSetService) Update(id uint, testSet model.TestSet) error {
	return s.TestSetRepo.Update(id, testSet)
}

func (s *TestSetService) DeleteById(id uint) error {
	return s.TestSetRepo.DeleteById(id)
}
