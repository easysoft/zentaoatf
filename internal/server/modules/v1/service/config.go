package service

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	serverConfig "github.com/aaronchen2k/deeptest/internal/server/config"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
)

type ConfigService struct {
	ProjectRepo *repo.ProjectRepo `inject:""`
}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

func (s *ConfigService) SaveConfig(config commDomain.ProjectConf) (err error) {
	currProject, err := s.ProjectRepo.GetCurrProjectByUser()
	if err != nil {
		return
	}

	commConsts.ProjectConfig = config
	serverConfig.SaveConfig(config, currProject.Path)

	return
}
