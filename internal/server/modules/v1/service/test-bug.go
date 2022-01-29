package service

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
)

type TestBugService struct {
}

func NewTestBugService() *TestBugService {
	return &TestBugService{}
}

func (s *TestBugService) Submit(bug commDomain.ZtfBug) (err error) {
	//zentaoUtils.CommitBug(bug)

	return
}
