package service

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/zentao"
)

type TestBugService struct {
}

func NewTestBugService() *TestBugService {
	return &TestBugService{}
}

func (s *TestBugService) Submit(bug commDomain.ZtfBug, workspacePath string) (err error) {
	err = zentaoUtils.CommitBug(bug, workspacePath)

	return
}
