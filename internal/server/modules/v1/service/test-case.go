package service

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/zentao"
)

type TestCaseService struct {
}

func NewTestCaseService() *TestCaseService {
	return &TestCaseService{}
}

func (s *TestCaseService) LoadTestCases(productId, moduleId, suiteId, taskId int, projectPath string) (
	cases []commDomain.ZtfCase, loginFail bool) {

	return zentaoUtils.LoadTestCases(productId, moduleId, suiteId, taskId, projectPath)
}
