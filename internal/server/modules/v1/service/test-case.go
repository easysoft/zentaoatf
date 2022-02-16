package service

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	configUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/zentao"
)

type TestCaseService struct {
}

func NewTestCaseService() *TestCaseService {
	return &TestCaseService{}
}

func (s *TestCaseService) LoadTestCases(productId, moduleId, suiteId, taskId int, projectPath string) (
	cases []commDomain.ZtfCase, loginFail bool) {

	config := configUtils.LoadByProjectPath(projectPath)

	err := zentaoUtils.Login(config)
	if err != nil {
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
