package service

import (
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	configUtils "github.com/aaronchen2k/deeptest/internal/server/modules/helper/config"
	zentaoUtils "github.com/aaronchen2k/deeptest/internal/server/modules/helper/zentao"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
)

type ZtfCaseService struct {
	ProjectRepo *repo.ProjectRepo `inject:""`
}

func (s *ZtfCaseService) NewZtfCaseService() *ZtfCaseService {
	return &ZtfCaseService{}
}

func (s *ZtfCaseService) LoadTestCases(productId, moduleId, suiteId, taskId int, projectPath string) (
	cases []commDomain.ZtfCase, loginFail bool) {

	config := configUtils.LoadByProjectPath(projectPath)

	ok := zentaoUtils.Login(config)
	if !ok {
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
