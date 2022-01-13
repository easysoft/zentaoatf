package service

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	"gopkg.in/ini.v1"
	"path"
	"path/filepath"
)

type ConfigService struct {
	ProjectRepo *repo.ProjectRepo `inject:""`
}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

func (s *ConfigService) LoadByProjectPath(projectPath string) (config commDomain.ProjectConf) {
	pth := filepath.Join(projectPath, commConsts.ConfigDir, commConsts.ConfigFile)
	ini.MapTo(&config, pth)

	config.Url = commonUtils.AddSlashForUrl(config.Url)

	return config
}

func (s *ConfigService) SaveConfig(config commDomain.ProjectConf, projectPath string) (err error) {
	s.SaveToFile(config, projectPath)

	return
}

func (s *ConfigService) ReadFromFile(projectPath string) (config commDomain.ProjectConf) {
	pth := filepath.Join(projectPath, commConsts.ConfigDir, commConsts.ConfigFile)
	ini.MapTo(&config, pth)

	config.Url = commonUtils.AddSlashForUrl(config.Url)

	return config
}

func (s *ConfigService) SaveToFile(config commDomain.ProjectConf, projectPath string) (err error) {
	pth := filepath.Join(projectPath, commConsts.ConfigDir, commConsts.ConfigFile)
	fileUtils.MkDirIfNeeded(path.Dir(pth))

	config.Version = commConsts.ConfigVersion

	cfg := ini.Empty()
	cfg.ReflectFrom(&config)

	cfg.SaveTo(pth)
	logUtils.Infof("Successfully update config file %s.", pth)

	return nil
}
