package service

import (
	"github.com/aaronchen2k/deeptest/internal/comm/consts"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverConfig "github.com/aaronchen2k/deeptest/internal/server/config"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"gopkg.in/ini.v1"
	"path"
	"path/filepath"
)

type ConfigService struct {
}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

func (s *ConfigService) ReadCurrConfig(projectPath string) (config serverDomain.ProjectConfig) {
	pth := filepath.Join(projectPath, commConsts.ConfigDir, commConsts.ConfigFile)

	config.Url = commonUtils.AddSlashForUrl(config.Url)
	ini.MapTo(&config, pth)

	return config
}

func (s *ConfigService) SaveConfig(config serverDomain.ProjectConfig, projectPath string) (err error) {
	pth := filepath.Join(projectPath, commConsts.ConfigDir, commConsts.ConfigFile)
	fileUtils.MkDirIfNeeded(path.Dir(pth))

	config.Version = serverConfig.CONFIG.System.Version

	cfg := ini.Empty()
	cfg.ReflectFrom(&config)

	cfg.SaveTo(pth)
	logUtils.Infof("Successfully update config file %s.", pth)

	return nil
}
