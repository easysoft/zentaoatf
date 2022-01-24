package service

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
)

type TestExecService struct {
	TestExecRepo *repo.TestExecRepo `inject:""`
}

func NewTestExecService() *TestExecService {
	return &TestExecService{}
}

func (s *TestExecService) List(projectPath string) (
	ret domain.PageData, err error) {

	return
}

func (s *TestExecService) Get(path string) (exec model.TestExec, err error) {
	return
}

func (s *TestExecService) Delete(path string) (err error) {
	return
}
