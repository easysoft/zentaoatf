package service

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/server/modules/helper/zentao"
)

type TestBugService struct {
}

func NewTestBugService() *TestBugService {
	return &TestBugService{}
}

func (s *TestBugService) Submit(bug commDomain.ZtfBug, projectPath string) (err error) {
	err = zentaoUtils.CommitBug(bug, projectPath)

	return
}
